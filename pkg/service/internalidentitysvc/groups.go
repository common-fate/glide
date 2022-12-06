package internalidentitysvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

func (s *Service) CreateGroup(ctx context.Context, in types.CreateGroupRequest) (*identity.Group, error) {
	id := types.NewGroupID()
	group := identity.Group{
		ID:        id,
		IdpID:     id,
		Name:      in.Name,
		Status:    types.IdpStatusACTIVE,
		Source:    identity.INTERNAL,
		Users:     in.Members,
		CreatedAt: s.Clock.Now(),
		UpdatedAt: s.Clock.Now(),
	}
	if in.Description != nil {
		group.Description = *in.Description
	}

	users := make(map[string]identity.User)
	hasMore := true
	var nextToken *string
	for hasMore {
		uq := storage.ListUsers{}
		r, err := s.DB.Query(ctx, &uq)
		if err != nil {
			return nil, err
		}
		if r.NextPage != "" {
			nextToken = &r.NextPage
		}
		hasMore = nextToken != nil
		for _, u := range uq.Result {
			users[u.ID] = u
		}
	}

	itemsToUpdate := []ddb.Keyer{&group}
	// validate that the members exist
	for _, newMemberID := range in.Members {
		if user, ok := users[newMemberID]; !ok {
			return nil, UserNotFoundError{UserID: newMemberID}
		} else {
			user.AddGroup(group.ID)
			itemsToUpdate = append(itemsToUpdate, &user)
		}
	}
	err := s.DB.PutBatch(ctx, itemsToUpdate...)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (s *Service) UpdateGroup(ctx context.Context, group identity.Group, in types.CreateGroupRequest) (*identity.Group, error) {
	if group.Source != identity.INTERNAL {
		return nil, ErrNotInternal
	}

	users := make(map[string]identity.User)
	hasMore := true
	var nextToken *string
	for hasMore {
		uq := storage.ListUsers{}
		r, err := s.DB.Query(ctx, &uq)
		if err != nil {
			return nil, err
		}
		if r.NextPage != "" {
			nextToken = &r.NextPage
		}
		hasMore = nextToken != nil
		for _, u := range uq.Result {
			users[u.ID] = u
		}
	}
	// validate that the members exist
	for _, newMemberID := range in.Members {
		if _, ok := users[newMemberID]; !ok {
			return nil, UserNotFoundError{UserID: newMemberID}
		}
	}

	var itemsToUpdate []ddb.Keyer

	for _, u := range group.Users {
		if !contains(in.Members, u) {
			user := users[u]
			user.RemoveGroup(group.ID)
			itemsToUpdate = append(itemsToUpdate, &user)
		}
	}
	for _, u := range in.Members {
		if !contains(group.Users, u) {
			user := users[u]
			user.AddGroup(group.ID)
			itemsToUpdate = append(itemsToUpdate, &user)
		}
	}

	group.UpdatedAt = s.Clock.Now()
	if in.Description == nil {
		group.Description = ""
	} else {
		group.Description = *in.Description
	}
	group.Name = in.Name
	group.Users = in.Members

	itemsToUpdate = append(itemsToUpdate, &group)
	err := s.DB.PutBatch(ctx, itemsToUpdate...)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *Service) DeleteGroup(ctx context.Context, group identity.Group) error {
	if group.Source != identity.INTERNAL {
		return ErrNotInternal
	}

	users := make(map[string]identity.User)
	hasMore := true
	var nextToken *string
	for hasMore {
		uq := storage.ListUsers{}
		r, err := s.DB.Query(ctx, &uq)
		if err != nil {
			return err
		}
		if r.NextPage != "" {
			nextToken = &r.NextPage
		}
		hasMore = nextToken != nil
		for _, u := range uq.Result {
			users[u.ID] = u
		}
	}

	var itemsToUpdate []ddb.Keyer

	// remove group from users
	for _, u := range group.Users {
		user := users[u]
		user.RemoveGroup(group.ID)
		itemsToUpdate = append(itemsToUpdate, &user)
	}

	group.UpdatedAt = s.Clock.Now()
	group.Users = []string{}
	group.Status = types.IdpStatusARCHIVED

	itemsToUpdate = append(itemsToUpdate, &group)
	return s.DB.PutBatch(ctx, itemsToUpdate...)
}
