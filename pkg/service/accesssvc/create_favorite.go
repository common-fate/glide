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
		// we don't know how to handle the error from the rule getter, so just return nil,it to the caller.
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
			if v.AdditionalProperties != nil {
				argumentCombos := combinations(v.AdditionalProperties)
				for _, argumentcombo := range argumentCombos {
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

func combinations(subRequest map[string][]string) []map[string]string {
	keys := make([]string, 0, len(subRequest))
	for k := range subRequest {
		keys = append(keys, k)
	}
	var combinations []map[string]string
	if len(keys) > 0 {
		for _, value := range subRequest[keys[0]] {
			combinations = append(combinations, branch(subRequest, keys, map[string]string{keys[0]: value}, 1)...)
		}
	}
	return combinations
}

func branch(subRequest map[string][]string, keys []string, combination map[string]string, keyIndex int) []map[string]string {
	var combos []map[string]string
	key := keys[keyIndex]
	for _, value := range subRequest[key] {
		// Create the target map
		next := map[string]string{key: value}
		// Copy from the original map to the target map
		for k, v := range combination {
			next[k] = v
		}
		if len(keys) == keyIndex+1 {
			combos = append(combos, next)
		} else {
			combos = append(combos, branch(subRequest, keys, next, keyIndex+1)...)
		}
	}
	return combos
}
