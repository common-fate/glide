package identitysync

import (
	"context"
	"fmt"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
)

type IdentityProvider interface {
	ListUsers(ctx context.Context) ([]identity.IdpUser, error)
	ListGroups(ctx context.Context) ([]identity.IdpGroup, error)
}

const (
	COGNITO = "COGNITO"
	GOOGLE  = "GOOGLE"
	OKTA    = "OKTA"
	AZURE   = "AZURE"
)

type IdentitySyncer struct {
	db  ddb.Storage
	idp IdentityProvider
}

type SyncOpts struct {
	TableName      string
	IdpType        string
	UserPoolId     string
	IdentityConfig deploy.IdentityConfig
}

func NewIdentitySyncer(ctx context.Context, opts SyncOpts) (*IdentitySyncer, error) {
	db, err := ddb.New(ctx, opts.TableName)
	if err != nil {
		return nil, err
	}

	var idp IdentityProvider

	switch opts.IdpType {
	case COGNITO:
		idp, err = NewCognito(ctx, Opts{UserPoolID: opts.UserPoolId})
		if err != nil {
			return nil, err
		}
	case OKTA:
		if opts.IdentityConfig.Okta == nil {
			return nil, fmt.Errorf("okta settings not configured")
		}
		clone := *opts.IdentityConfig.Okta
		err = config.LoadAndReplaceSSMValues(ctx, &clone)
		if err != nil {
			return nil, err
		}
		idp, err = NewOkta(ctx, clone)
		if err != nil {
			return nil, err
		}

	case GOOGLE:
		if opts.IdentityConfig.Google == nil {
			return nil, fmt.Errorf("google settings not configured")
		}
		clone := *opts.IdentityConfig.Google
		err = config.LoadAndReplaceSSMValues(ctx, &clone)
		if err != nil {
			return nil, err
		}
		idp, err = NewGcp(ctx, clone)
		if err != nil {
			return nil, err
		}

	// case AZURE:
	default:
		return nil, fmt.Errorf("unsupported provider type %s", opts.IdpType)

	}

	return &IdentitySyncer{
		db:  db,
		idp: idp,
	}, nil
}

func (s *IdentitySyncer) Sync(ctx context.Context) error {
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
	usersMap, groupsMap := processUsersAndGroups(idpUsers, idpGroups, uq.Result, gq.Result)
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

// processUsersAndGroups conatins all the logic for create/update/archive for users and groups
//
// It returns a map of users and groups ready to be inserted to the database
func processUsersAndGroups(idpUsers []identity.IdpUser, idpGroups []identity.IdpGroup, internalUsers []identity.User, internalGroups []identity.Group) (map[string]identity.User, map[string]identity.Group) {
	idpUserMap := make(map[string]identity.IdpUser)
	for _, u := range idpUsers {
		idpUserMap[u.Email] = u
	}
	idpGroupMap := make(map[string]identity.IdpGroup)
	for _, g := range idpGroups {
		idpGroupMap[g.ID] = g
	}
	ddbUserMap := make(map[string]identity.User)
	for _, u := range internalUsers {
		ddbUserMap[u.Email] = u
	}
	ddbGroupMap := make(map[string]identity.Group)
	// This map ensures we have a distinct list of ids
	internalGroupUsers := make(map[string]map[string]string)
	for _, g := range internalGroups {
		ddbGroupMap[g.IdpID] = g
		internalGroupUsers[g.ID] = make(map[string]string)
	}

	// update/create users
	for _, u := range idpUsers {
		if existing, ok := ddbUserMap[u.Email]; ok { //update
			existing.FirstName = u.FirstName
			existing.LastName = u.LastName
			ddbUserMap[u.Email] = existing
		} else { // create
			ddbUserMap[u.Email] = u.ToInternalUser()
		}
	}
	// update/create groups
	for _, g := range idpGroups {
		if existing, ok := ddbGroupMap[g.ID]; ok { //update
			existing.Description = g.Description
			existing.Name = g.Name
			ddbGroupMap[g.ID] = existing
		} else { // create
			newGroup := g.ToInternalGroup()
			ddbGroupMap[g.ID] = newGroup
			internalGroupUsers[newGroup.ID] = make(map[string]string)
		}
	}

	// archive deleted users
	for k, u := range ddbUserMap {
		if _, ok := idpUserMap[k]; !ok {
			u.Status = types.ARCHIVED
			// Remove all group associations from archived users
			u.Groups = []string{}
			ddbUserMap[k] = u
		}
	}
	// archive deleted groups
	for k, g := range ddbGroupMap {
		if _, ok := idpGroupMap[k]; !ok {
			g.Status = types.ARCHIVED
			// Remove all user associations from archived groups
			g.Users = []string{}
			ddbGroupMap[k] = g
		}
	}

	for _, idpUser := range idpUserMap {

		// This map ensures we have a distinct list of ids
		internalGroupIds := map[string]string{}
		for _, idpGroupId := range idpUser.Groups {
			gid := ddbGroupMap[idpGroupId].ID
			internalGroupIds[gid] = gid
			uid := ddbUserMap[idpUser.Email].ID
			internalGroupUsers[gid][uid] = uid
		}
		internalUser := ddbUserMap[idpUser.Email]
		keys := make([]string, 0, len(internalGroupIds))
		for k := range internalGroupIds {
			keys = append(keys, k)
		}
		internalUser.Groups = keys
		ddbUserMap[idpUser.Email] = internalUser
	}

	// Updates the internal groups with new user mappings
	for k, v := range ddbGroupMap {
		um := internalGroupUsers[v.ID]
		keys := make([]string, 0, len(um))
		for k2 := range um {
			keys = append(keys, k2)
		}
		v.Users = keys
		ddbGroupMap[k] = v
	}

	return ddbUserMap, ddbGroupMap
}
