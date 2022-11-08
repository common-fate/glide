package internalidentitysvc

import (
	"context"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
)

func (s *Service) CreateGroup(ctx context.Context, in types.CreateGroupRequest) (*identity.Group, error) {
	err := s.validateMembers(ctx, in.Members)
	if err != nil {
		return nil, err
	}
	group := identity.Group{
		ID:        types.NewGroupID(),
		IdpID:     in.Name,
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
	err = s.DB.Put(ctx, &group)
	if err != nil {
		return nil, err
	}
	err = s.IdentitySyncer.Sync(ctx)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (s *Service) UpdateGroup(ctx context.Context, group identity.Group, in types.CreateGroupRequest) (*identity.Group, error) {
	if group.Source != identity.INTERNAL {
		return nil, ErrNotInternal
	}
	err := s.validateMembers(ctx, in.Members)
	if err != nil {
		return nil, err
	}
	group.UpdatedAt = s.Clock.Now()
	if in.Description == nil {
		group.Description = ""
	} else {
		group.Description = *in.Description
	}
	group.Name = in.Name
	group.Users = in.Members

	err = s.DB.Put(ctx, &group)
	if err != nil {
		return nil, err
	}
	err = s.IdentitySyncer.Sync(ctx)
	if err != nil {
		return nil, err
	}
	return &group, nil
}
func (s *Service) validateMembers(ctx context.Context, members []string) error {
	// validate that the members exist
	for _, newMemberID := range members {
		u := storage.GetUser{ID: newMemberID}
		_, err := s.DB.Query(ctx, &u)
		if err == ddb.ErrNoItems {
			return UserNotFoundError{UserID: newMemberID}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
