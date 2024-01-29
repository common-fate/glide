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
	"github.com/common-fate/common-fate/pkg/storage/ddbhelpers"
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
	syncMutex   sync.Mutex
	groupFilter string
}

type SyncOpts struct {
	TableName           string
	IdpType             string
	UserPoolId          string
	IdentityConfig      deploy.FeatureMap
	IdentityGroupFilter string
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
		db:          db,
		idp:         idp.IdentityProvider,
		idpType:     opts.IdpType,
		groupFilter: opts.IdentityGroupFilter,
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
	filter := s.groupFilter
	useIdpGroupsAsFilter := filter != ""

	if useIdpGroupsAsFilter {
		// overwrite the existing groups with the filtered groups
		log.Infow("filtering groups", "filter", filter)
		idpGroups, err = FilterGroups(idpGroups, filter)
		if err != nil {
			return err
		}
	}

	log.Infow("fetched users and groups from IDP", "users.count", len(idpUsers), "groups.count", len(idpGroups))

	s.setDeploymentInfo(ctx, log, depid.UserInfo{UserCount: len(idpUsers), GroupCount: len(idpGroups), IDP: s.idpType})

	dbUsers, err := listAllDbUsers(ctx, s.db)
	if err != nil {
		return err
	}
	gq := &storage.ListGroups{}
	_, err = s.db.Query(ctx, gq)
	if err != nil {
		return err
	}
	dbGroups, err := listAllDbGroups(ctx, s.db)
	if err != nil {
		return err
	}

	usersMap, duplicateUsersMap, groupsMap := processUsersAndGroups(s.idpType, idpUsers, idpGroups, dbUsers, dbGroups, useIdpGroupsAsFilter)

	if len(duplicateUsersMap) > 0 {
		duplicatedEmails := []string{}
		for k, _ := range duplicateUsersMap {
			duplicatedEmails = append(duplicatedEmails, k)
		}
		log.Errorw("error found duplicate entries in DB", "users.duplicate-count", len(duplicateUsersMap), "users.dupliated-emails", duplicatedEmails)
	}

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

func listAllDbUsers(ctx context.Context, db ddb.Storage) ([]identity.User, error) {
	dbUsers := []identity.User{}
	uq := &storage.ListUsers{}
	err := ddbhelpers.QueryPages(ctx, db, uq,
		func(pageResult *ddb.QueryResult, pageQueryBuilder ddb.QueryBuilder, lastPage bool) bool {
			if qb, ok := pageQueryBuilder.(*storage.ListUsers); ok {
				dbUsers = append(dbUsers, qb.Result...)
			} else {
				panic("Unknown type for QueryBuilder")
			}
			return true
		},
	)
	return dbUsers, err
}

func listAllDbGroups(ctx context.Context, db ddb.Storage) ([]identity.Group, error) {
	dbGroups := []identity.Group{}
	uq := &storage.ListGroups{}
	err := ddbhelpers.QueryPages(ctx, db, uq,
		func(pageResult *ddb.QueryResult, pageQueryBuilder ddb.QueryBuilder, lastPage bool) bool {
			if qb, ok := pageQueryBuilder.(*storage.ListGroups); ok {
				dbGroups = append(dbGroups, qb.Result...)
			} else {
				panic("Unknown type for QueryBuilder")
			}
			return true
		},
	)
	return dbGroups, err
}

// processUsersAndGroups
// contains all the logic for create/update/archive for users and groups
// It returns a map of users and groups ready to be inserted to the database
//
// Expected Behavior:
// useIdpGroupsAsFilter == true: only users with groups that exist in the IDP will be returned, this is used conditionally with a regex filter that prefilters any groups. Side effects: users with no groups are removed, only filtered groups show in the UI (loss of information; for better or worse)
//
// useIdpGroupsAsFilter == false: users with no groups remain, all groups show in the UI. If a user/group is removed from the IDP, it will be archived in the DB
func processUsersAndGroups(idpType string, idpUsers []identity.IDPUser, idpGroups []identity.IDPGroup, internalUsers []identity.User, internalGroups []identity.Group, useIdpGroupsAsFilter bool) (map[string]identity.User, map[string][]identity.User, map[string]identity.Group) {

	idpGroupMap := make(map[string]identity.IDPGroup)
	for _, g := range idpGroups {
		idpGroupMap[g.ID] = g
	}
	idpUserMap := make(map[string]identity.IDPUser)
	for _, u := range idpUsers {
		idpUserMap[u.Email] = u
	}
	ddbUserMap := make(map[string]identity.User)
	ddbUserDuplicatedMap := make(map[string][]identity.User)
	for _, u := range internalUsers {
		// Collect the latest updated entry with same email
		// This is to deal in the bogus case of duplicates in storage
		if u2, ok := ddbUserMap[u.Email]; !ok {
			ddbUserMap[u.Email] = u
		} else {
			if _, ok := ddbUserDuplicatedMap[u.Email]; !ok {
				ddbUserDuplicatedMap[u.Email] = []identity.User{}
			}
			if u.CreatedAt.Before(u2.CreatedAt) {
				ddbUserMap[u.Email] = u
				ddbUserDuplicatedMap[u.Email] = append(ddbUserDuplicatedMap[u.Email], u2)
			} else {
				ddbUserDuplicatedMap[u.Email] = append(ddbUserDuplicatedMap[u.Email], u)
			}
		}
	}
	ddbGroupMap := make(map[string]identity.Group)
	// This map ensures we have a distinct list of ids
	internalGroupUsers := make(map[string]map[string]string)
	for _, g := range internalGroups {
		ddbGroupMap[g.IdpID] = g
		internalGroupUsers[g.ID] = make(map[string]string)
	}

	// update/create users
	for _, u := range idpUserMap {
		//update
		if existing, ok := ddbUserMap[u.Email]; ok {
			existing.FirstName = u.FirstName
			existing.LastName = u.LastName
			ddbUserMap[u.Email] = existing
		} else {
			// create
			ddbUserMap[u.Email] = u.ToInternalUser()
		}
	}
	// update/create groups
	for _, idpGroup := range idpGroups {
		if existingGroup, ok := ddbGroupMap[idpGroup.ID]; ok { //update
			existingGroup.Description = idpGroup.Description
			existingGroup.Name = idpGroup.Name
			existingGroup.Status = types.IdpStatusACTIVE
			existingGroup.Source = idpType
			ddbGroupMap[idpGroup.ID] = existingGroup
		} else { // create
			newGroup := idpGroup.ToInternalGroup(idpType)
			ddbGroupMap[idpGroup.ID] = newGroup
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

				g.Status = types.IdpStatusARCHIVED
				g.Users = []string{}
				ddbGroupMap[k] = g

				continue // not covered by tests
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

		// This map ensures we have a distinct list of ids
		internalGroupIds := map[string]string{}
		for _, gid := range idpUser.Groups {

			// If we are using the IDP groups as a filter, then we only want to add the groups that exist in the IDP
			if useIdpGroupsAsFilter {
				if _, ok := idpGroupMap[gid]; !ok {
					// continue i.e. skip adding this group since it doesn't exist in the IDP
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
		// for every internal user group
		for _, internalGroupId := range internalUser.Groups {

			// for each internalGroupId on an internal users groups
			// get the source of the group from the db
			// if the group is internal, add it to the list of groups

			source := ddbGroupMap[internalGroupId].Source
			// if the group is internal, add it to the list of groups
			if source == identity.INTERNAL {
				gid := ddbGroupMap[internalGroupId].ID // not covered by tests
				internalGroupIds[gid] = gid            // not covered by tests
			}
		}

		groupKeys := make([]string, 0, len(internalGroupIds))
		for k := range internalGroupIds {
			groupKeys = append(groupKeys, k)
		}

		internalUser.Groups = groupKeys
		// if the user is not in any groups, archive them
		if len(internalUser.Groups) == 0 && useIdpGroupsAsFilter {
			internalUser.Status = types.IdpStatusARCHIVED
		}
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

	return ddbUserMap, ddbUserDuplicatedMap, ddbGroupMap
}
