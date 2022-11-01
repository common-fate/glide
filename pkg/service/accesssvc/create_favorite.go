package accesssvc

import (
	"context"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
)

// CreateRequest creates a new request and saves it in the database.
// Returns an error if the request is invalid.
func (s *Service) CreateFavorite(ctx context.Context, user *identity.User, in types.CreateFavoriteRequest) (*access.Favorite, error) {
	log := logger.Get(ctx).With("user.id", user.ID)
	q := storage.GetAccessRuleCurrent{ID: in.AccessRuleId}
	_, err := s.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		return nil, ErrRuleNotFound
	}
	if err != nil {
		return nil, err
	}
	rule := q.Result

	log.Debugw("verifying user belongs to access rule groups", "rule.groups", rule.Groups, "user.groups", user.Groups)
	err = groupMatches(rule.Groups, user.Groups)
	if err != nil {
		return nil, err
	}

	now := s.Clock.Now()

	requestArguments, err := s.Rules.RequestArguments(ctx, rule.Target)
	if err != nil {
		return nil, err
	}
	var favoriteWith []map[string][]string
	if in.With != nil {
		for _, v := range *in.With {
			for _, argumentcombo := range v.ArgumentCombinations() {
				err = validateRequest(CreateRequest{
					AccessRuleId: in.AccessRuleId,
					Reason:       in.Reason,
					Timing:       in.Timing,
					With:         argumentcombo,
				}, rule, requestArguments)
				if err != nil {
					return nil, err
				}
			}
			favoriteWith = append(favoriteWith, v.AdditionalProperties)
		}
	} else {
		err = validateRequest(CreateRequest{
			AccessRuleId: in.AccessRuleId,
			Reason:       in.Reason,
			Timing:       in.Timing,
			With:         make(map[string]string),
		}, rule, requestArguments)
		if err != nil {
			return nil, err
		}
	}

	favorite := access.Favorite{
		ID:     types.NewRequestFavoriteID(),
		UserID: user.ID,
		Name:   in.Name,
		Rule:   in.AccessRuleId,
		Data: access.RequestData{
			Reason: in.Reason,
		},
		With:            favoriteWith,
		RequestedTiming: access.TimingFromRequestTiming(in.Timing),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = s.DB.Put(ctx, &favorite)
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}
