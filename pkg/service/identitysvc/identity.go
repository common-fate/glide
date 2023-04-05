package identitysvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
)

func (s *Service) UpdateUserAccessRules(ctx context.Context, users map[string]identity.User, groups map[string]identity.Group) (map[string]identity.User, error) {
	//get all access rules

	//reset user access rules, as we will be overriding them

	for _, u := range users {
		u.AccessRules = []string{}
	}

	q := storage.ListCurrentAccessRules{}

	_, err := s.DB.Query(ctx, &q)
	if err != nil {
		return nil, err
	}

	for _, ar := range q.Result {

		for _, g := range ar.Groups {
			group := groups[g]

			for _, u := range group.Users {
				user := users[u]
				user.AccessRules = append(user.AccessRules, ar.ID)
				users[u] = user
			}
		}
	}

	return users, nil
}
