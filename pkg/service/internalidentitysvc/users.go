package internalidentitysvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

func (s *Service) UpdateUserGroups(ctx context.Context, user identity.User, groups []string) (*identity.User, error) {
	// update the internal groups
	// leave the external groups unchanged
	internalGroups := make(map[string]identity.Group)
	hasMore := true
	var nextToken *string
	for hasMore {
		gq := storage.ListGroupsForSourceAndStatus{Source: identity.INTERNAL, Status: types.ACTIVE}
		r, err := s.DB.Query(ctx, &gq)
		if err != nil {
			return nil, err
		}
		if r.NextPage != "" {
			nextToken = &r.NextPage
		}
		hasMore = nextToken != nil
		for _, g := range gq.Result {
			internalGroups[g.ID] = g
		}
	}

	// Append the user to all the groups that the
	var itemsToUpdate []ddb.Keyer
	now := s.Clock.Now()
	// add user to all these groups
	for _, g := range groups {
		if ig, ok := internalGroups[g]; !ok {
			return nil, ErrGroupNotFoundOrNotInternal
		} else if !contains(ig.Users, user.ID) {
			ig.Users = append(ig.Users, user.ID)
			ig.UpdatedAt = now
			itemsToUpdate = append(itemsToUpdate, &ig)
		}
	}

	// remove the user from any internal groups if they are no longer in them
	for _, g := range user.Groups {
		if ig, ok := internalGroups[g]; ok {
			// this group has been removed from the user
			if !contains(groups, g) {
				var newUsers []string
				for _, u := range ig.Users {
					if u != user.ID {
						newUsers = append(newUsers, u)
					}
				}
				ig.Users = newUsers
				ig.UpdatedAt = now
				itemsToUpdate = append(itemsToUpdate, &ig)
			}
		}
	}

	updatedUserGroups := groups
	for _, g := range user.Groups {
		if _, ok := internalGroups[g]; !ok {
			updatedUserGroups = append(updatedUserGroups, g)
		}
	}
	user.Groups = updatedUserGroups
	user.UpdatedAt = s.Clock.Now()
	itemsToUpdate = append(itemsToUpdate, &user)
	err := s.DB.PutBatch(ctx, itemsToUpdate...)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func contains(set []string, str string) bool {
	for _, s := range set {
		if s == str {
			return true
		}
	}

	return false
}
