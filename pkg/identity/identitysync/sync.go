package identitysync

import (
	"context"
	"errors"
	"sync"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/depid"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

type IdentityProvider interface {
	ListUsers(ctx context.Context) ([]identity.IDPUser, error)
	ListGroups(ctx context.Context) ([]identity.IDPGroup, error)
	gconfig.Configer
	gconfig.Initer
}

type IdentitySyncer struct {
	db      ddb.Storage
	idp     IdentityProvider
	idpType string
	// used to prevent concurrent calls to sync
	// prevents unexpected duplication of users and groups when used asyncronously
	syncMutex sync.Mutex
}

type SyncOpts struct {
	TableName      string
	IdpType        string
	UserPoolId     string
	IdentityConfig deploy.FeatureMap
}

func NewIdentitySyncer(ctx context.Context, opts SyncOpts) (*IdentitySyncer, error) {
	db, err := ddb.New(ctx, opts.TableName)
	if err != nil {
		return nil, err
	}

	idp, err := Registry().Lookup(opts.IdpType)
	if err != nil {
		return nil, err
	}
	cfg := idp.IdentityProvider.Config()
	var found bool
	if opts.IdpType == IDPTypeCognito {
		// Cognito has slightly different loading behaviour becauae it is the default provider
		// config is provided directly via env vars when the stack is deployed, rather than via a cloudformation parameter
		found = true
		err = cfg.Load(ctx, &gconfig.MapLoader{Values: map[string]string{
			"userPoolId": opts.UserPoolId,
		}})
		if err != nil {
			return nil, err
		}
	} else {
		if idpCfg, ok := opts.IdentityConfig[opts.IdpType]; ok {
			found = true
			err = cfg.Load(ctx, &gconfig.MapLoader{Values: idpCfg})
			if err != nil {
				return nil, err
			}
		}
	}
	if !found {
		return nil, errors.New("no matching configuration found for idp type")
	}

	err = idp.IdentityProvider.Init(ctx)
	if err != nil {
		return nil, err
	}
	return &IdentitySyncer{
		db:      db,
		idp:     idp.IdentityProvider,
		idpType: opts.IdpType,
	}, nil
}

func (s *IdentitySyncer) Sync(ctx context.Context) error {
	// prevent concurrent calls to sync
	s.syncMutex.Lock()
	defer s.syncMutex.Unlock()
	log := logger.Get(ctx)

	//Fetch all users from IDP
	// The IDP should return the group mappings for users, these group IDs will be internal to the IDP
	idpUsers, err := s.idp.ListUsers(ctx)
	if err != nil {
		return err
	}
	// Fetch all groups from IDP
	idpGroups, err := s.idp.ListGroups(ctx)
	if err != nil {
		return err
	}

	/*

		example regex filter: "admins|devops"

		{{{ ORIGINAL INPUT }}}

		groups
		{name: "admins", id: "1"}
		{name: "dev_ops", id: "2"}
		{name: "accounting", id: "3"}

		users
		{name: "bob", id: "1", groups: ["1", "2"]}  // grant access
		{name: "alice", id: "2", groups: ["3"]}	// deny access
		{name: "joe", id: "3", groups: ["1", "3"]} // grant access
		{name: "jane", id: "4", groups: ["2"]} // grant access

		{{{ POST GROUP FILTER }}}

		groups
		{name: "admins", id: "1"}
		{name: "dev_ops", id: "2"}

		users
		{name: "bob", id: "1", groups: ["1", "2"]}  // grant access
		{name: "alice", id: "2", groups: ["3"]}	// deny access
		{name: "joe", id: "3", groups: ["1", "3"]} // grant access
		{name: "jane", id: "4", groups: ["2"]} // grant access

		{{{ POST processUsersAndGroups }}}

		groups
		{name: "admins", id: "1"}
		{name: "dev_ops", id: "2"}

		users
		{name: "bob", id: "1", groups: ["1", "2"]}  // grant access
		{name: "joe", id: "3", groups: ["1"]} // grant access
		{name: "jane", id: "4", groups: ["2"]} // grant access


		SIDE EFFECTS
		- users with no groups are removed
		- only filtered groups show in the UI (loss of information; for better or worse)

		CURRENT STATE
		- users with no groups remain
		- all groups show in the UI

		APPROACH
		To maintain current state and introduce new filtering,
		We pass processUsersAndGroups an optional prop `useIdpGroupsAsFilter`. If true, then only users with groups that exist in the IDP will be returned, this is used conditionally with a regex filter that *prefilters any groups*. Side effects: users with no groups are removed, only filtered groups show in the UI


	*/
	useIdpGroupsAsFilter := true

	if useIdpGroupsAsFilter {
		// overwrite the existing groups with the filtered groups
		idpGroups, err = FilterGroups(idpGroups, "admins|granted-admins")
		if err != nil {
			return err
		}
	}

	log.Infow("fetched users and groups from IDP", "users.count", len(idpUsers), "groups.count", len(idpGroups))

	s.setDeploymentInfo(ctx, log, depid.UserInfo{UserCount: len(idpUsers), GroupCount: len(idpGroups), IDP: s.idpType})

	uq := &storage.ListUsers{}
	_, err = s.db.Query(ctx, uq)
	if err != nil {
		return err
	}
	gq := &storage.ListGroups{}
	_, err = s.db.Query(ctx, gq)
	if err != nil {
		return err
	}
	usersMap, groupsMap := processUsersAndGroups(s.idpType, idpUsers, idpGroups, uq.Result, gq.Result, useIdpGroupsAsFilter)
	items := make([]ddb.Keyer, 0, len(usersMap)+len(groupsMap))
	for _, v := range usersMap {
		vi := v
		items = append(items, &vi)
	}
	for _, v := range groupsMap {
		vi := v
		items = append(items, &vi)
	}

	return s.db.PutBatch(ctx, items...)
}

// analytics event
func (s *IdentitySyncer) setDeploymentInfo(ctx context.Context, log *zap.SugaredLogger, info depid.UserInfo) {
	ac := analytics.New(analytics.Env())
	dep, err := depid.New(s.db, log).SetUserInfo(ctx, info)
	if err != nil {
		log.Errorw("error setting deployment info", zap.Error(err))
	}
	if dep != nil {
		ac.Track(dep.ToAnalytics())
	}
}

// processUsersAndGroups
// contains all the logic for create/update/archive for users and groups
// It returns a map of users and groups ready to be inserted to the database
//
// If useIdpGroupsAsFilter is true, then only users with groups that exist in the IDP will be returned, this is used conditionally with a regex filter that prefilters any groups. Side effects: users with no groups are removed, only filtered groups show in the UI (loss of information; for better or worse)
func processUsersAndGroups(idpType string, idpUsers []identity.IDPUser, idpGroups []identity.IDPGroup, internalUsers []identity.User, internalGroups []identity.Group, useIdpGroupsAsFilter bool) (map[string]identity.User, map[string]identity.Group) {

	idpGroupMap := make(map[string]identity.IDPGroup)
	for _, g := range idpGroups {
		idpGroupMap[g.ID] = g
	}
	idpUserMap := make(map[string]identity.IDPUser)
	for _, u := range idpUsers {

		if useIdpGroupsAsFilter {
			idpUserHasMatchingGroup := false
			for _, g := range u.Groups {
				if _, ok := idpGroupMap[g]; ok {
					idpUserHasMatchingGroup = true
					break
				} else {
					continue
				}
			}
			if idpUserHasMatchingGroup {
				idpUserMap[u.Email] = u
			}
		} else {
			idpUserMap[u.Email] = u
		}
	}
	ddbUserMap := make(map[string]identity.User)
	for _, u := range internalUsers {
		ddbUserMap[u.Email] = u
	}
	ddbGroupMap := make(map[string]identity.Group)
	// This map ensures we have a distinct list of ids
	internalGroupUsers := make(map[string]map[string]string)
	for _, g := range internalGroups {
		if useIdpGroupsAsFilter {
			if _, ok := idpGroupMap[g.IdpID]; !ok {
				continue
			}
		}
		ddbGroupMap[g.IdpID] = g
		internalGroupUsers[g.ID] = make(map[string]string)
	}

	// update/create users
	for _, u := range idpUserMap {
		if existing, ok := ddbUserMap[u.Email]; ok { //update
			existing.FirstName = u.FirstName
			existing.LastName = u.LastName
			ddbUserMap[u.Email] = existing
		} else { // create

			if useIdpGroupsAsFilter {
				for _, g := range u.Groups {
					if _, ok := idpGroupMap[g]; ok {
						idpUserMap[u.Email] = u
					} else {
						continue
					}
				}
			}

			ddbUserMap[u.Email] = u.ToInternalUser()
		}
	}
	// update/create groups
	for _, g := range idpGroups {
		if useIdpGroupsAsFilter {
			if _, ok := idpGroupMap[g.ID]; !ok {
				continue
			}
		}

		if existing, ok := ddbGroupMap[g.ID]; ok { //update
			existing.Description = g.Description
			existing.Name = g.Name
			existing.Status = types.IdpStatusACTIVE
			existing.Source = idpType
			ddbGroupMap[g.ID] = existing
		} else { // create
			newGroup := g.ToInternalGroup(idpType)
			ddbGroupMap[g.ID] = newGroup
			internalGroupUsers[newGroup.ID] = make(map[string]string)
		}
	}

	// archive deleted users
	for k, u := range ddbUserMap {
		if _, ok := idpUserMap[k]; !ok {
			u.Status = types.IdpStatusARCHIVED
			// Remove all group associations from archived users
			u.Groups = []string{}
			ddbUserMap[k] = u
		} else {
			u.Status = types.IdpStatusACTIVE
			ddbUserMap[k] = u
		}
	}
	// archive deleted groups
	for k, g := range ddbGroupMap {

		if useIdpGroupsAsFilter {
			if _, ok := idpGroupMap[g.ID]; !ok {
				continue
			}
		}

		if _, ok := idpGroupMap[k]; !ok {
			if g.Source != identity.INTERNAL {

				g.Status = types.IdpStatusARCHIVED
				// Remove all user associations from archived groups
				g.Users = []string{}
				ddbGroupMap[k] = g
			}

		}
	}

	for _, idpUser := range idpUserMap {

		// If we are using the IDP groups as a filter, then we only want to add the groups that exist in the IDP
		if useIdpGroupsAsFilter {
			for _, g := range idpUser.Groups {
				if _, ok := idpGroupMap[g]; !ok {
					continue
				}
			}
		}

		// This map ensures we have a distinct list of ids
		internalGroupIds := map[string]string{}
		for _, gid := range idpUser.Groups {

			// If we are using the IDP groups as a filter, then we only want to add the groups that exist in the IDP
			if useIdpGroupsAsFilter {
				if _, ok := idpGroupMap[gid]; !ok {
					continue
				}
			}

			internalGroupIds[gid] = gid
			uid := ddbUserMap[idpUser.Email].ID

			// prevents a panic due to an assignment to a nil map
			if internalGroupUsers[gid] == nil {
				internalGroupUsers[gid] = map[string]string{}
			}

			internalGroupUsers[gid][uid] = uid
		}

		internalUser := ddbUserMap[idpUser.Email]

		//make sure we are saving the internal groups that the user is apart of
		for _, internalGroupId := range internalUser.Groups {
			// If we are using the IDP groups as a filter, then we only want to add the groups that exist in the IDP
			if useIdpGroupsAsFilter {
				if _, ok := idpGroupMap[internalGroupId]; !ok {
					continue
				}
			}

			source := ddbGroupMap[internalGroupId].Source
			if source == identity.INTERNAL {
				gid := ddbGroupMap[internalGroupId].ID
				internalGroupIds[gid] = gid
			}

		}

		keys := make([]string, 0, len(internalGroupIds))
		for k := range internalGroupIds {
			keys = append(keys, k)
		}

		internalUser.Groups = keys
		ddbUserMap[idpUser.Email] = internalUser
	}

	// Updates the internal groups with new user mappings
	for k, v := range ddbGroupMap {
		if v.Source != identity.INTERNAL {
			um := internalGroupUsers[v.ID]
			keys := make([]string, 0, len(um))
			for k2 := range um {
				keys = append(keys, k2)
			}
			v.Users = keys
			ddbGroupMap[k] = v
		}
	}

	return ddbUserMap, ddbGroupMap
}
