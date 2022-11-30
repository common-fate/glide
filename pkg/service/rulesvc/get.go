package rulesvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
)

func (s *Service) GetRule(ctx context.Context, ID string, user *identity.User, isAdmin bool) (*rule.GetAccessRuleResponse, error) {
	q := storage.GetAccessRuleCurrent{ID: ID}
	_, err := s.DB.Query(ctx, &q)
	// Throw storage errors if they occur
	if err != nil {
		return nil, err
	}

	if canGet(user, q.Result, isAdmin) {
		canRequest := false

		// even though admins can see the rule, requesting access to rule is only possible
		// if they are explicitly selected as request groups when creating the rule.
		for _, group := range user.Groups {
			for _, g := range q.Result.Groups {
				if g == group {
					canRequest = true
				}
			}
		}

		return &rule.GetAccessRuleResponse{
			Rule:       q.Result,
			CanRequest: canRequest,
		}, nil
	}

	// Otherwise not allowed
	return nil, ErrUserNotAuthorized
}

// canGet checks if the user is in the list of groups,
// or if the user is an approver of the rule
func canGet(user *identity.User, rule *rule.AccessRule, isAdmin bool) bool {
	// Admins can always access a rule
	if isAdmin {
		return true
	}
	// DE = User can see a rule they're an approver for
	for _, au := range rule.Approval.Users {
		if au == user.ID {
			return true
		}
	}
	// DE = User can see a rule they're an approver of (via groups)
	for _, group := range user.Groups {
		for _, g := range rule.Approval.Groups {
			if g == group {
				return true
			}
		}
	}
	// DE = User can see a rule they're assigned to (via the groups)
	for _, group := range user.Groups {
		for _, g := range rule.Groups {
			if g == group {
				return true
			}
		}
	}

	return false
}
