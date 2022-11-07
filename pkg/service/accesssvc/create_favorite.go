package accesssvc

import (
	"context"

	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/types"
)

type CreateFavoriteOpts struct {
	User   identity.User
	Create types.CreateFavoriteRequest
}

// CreateRequest creates a new request and saves it in the database.
// Returns an error if the request is invalid.
func (s *Service) CreateFavorite(ctx context.Context, in CreateFavoriteOpts) (*access.Favorite, error) {
	favorite, err := s.validateFavorite(ctx, in.User, in.Create)
	if err != nil {
		return nil, err
	}
	err = s.DB.Put(ctx, favorite)
	if err != nil {
		return nil, err
	}
	return favorite, nil
}

type UpdateFavoriteOpts struct {
	User     identity.User
	Favorite access.Favorite
	Update   types.CreateFavoriteRequest
}

// UpdateFavorite validates the input then updates the favorite
func (s *Service) UpdateFavorite(ctx context.Context, in UpdateFavoriteOpts) (*access.Favorite, error) {
	favorite, err := s.validateFavorite(ctx, in.User, in.Update)
	if err != nil {
		return nil, err
	}
	favorite.ID = in.Favorite.ID
	err = s.DB.Put(ctx, favorite)
	if err != nil {
		return nil, err
	}
	return favorite, nil
}

// validateFavorite validates the favorite and returns it, ready to be saved in the database
// Returns an error if the favorite is invalid.
func (s *Service) validateFavorite(ctx context.Context, user identity.User, in types.CreateFavoriteRequest) (*access.Favorite, error) {
	_, err := s.validateCreateRequests(ctx, CreateRequestsOpts{
		User: user,
		Create: CreateRequests{
			AccessRuleId: in.AccessRuleId,
			Reason:       in.Reason,
			Timing:       in.Timing,
			With:         in.With,
		},
	})
	if err != nil {
		return nil, err
	}
	now := s.Clock.Now()

	favorite := access.Favorite{
		ID:     types.NewRequestFavoriteID(),
		UserID: user.ID,
		Name:   in.Name,
		Rule:   in.AccessRuleId,
		Data: access.RequestData{
			Reason: in.Reason,
		},
		RequestedTiming: access.TimingFromRequestTiming(in.Timing),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if in.With != nil {
		for _, w := range *in.With {
			// omit any empty sections
			if w.AdditionalProperties != nil {
				favorite.With = append(favorite.With, w.AdditionalProperties)
			}
		}
	}

	return &favorite, nil
}
