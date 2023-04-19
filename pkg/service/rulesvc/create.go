package rulesvc

import (
	"context"
	"errors"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/types"
)

func (s *Service) ProcessTargets(ctx context.Context, in []types.CreateAccessRuleTarget) ([]rule.Target, error) {
	// @TODO

	// Check for duplicate target groups, there can only be one instance of a target group per access rule.

	// TODO when filter expressions are implemented
	// Check validity of filter expressions (the structure of a filter expression shoudl be validated in the API layer automatically)
	// - check that attributes in the filter expression exist on the schema of the target/resource type

	return nil, nil
}

func (s *Service) CreateAccessRule(ctx context.Context, userID string, in types.CreateAccessRuleRequest) (*rule.AccessRule, error) {

	id := types.NewAccessRuleID()

	log := logger.Get(ctx).With("user.id", userID, "access_rule.id", id)
	now := s.Clock.Now()

	targets, err := s.ProcessTargets(ctx, in.Targets)
	if err != nil {
		return nil, err
	}

	// validate it is under 6 months
	if in.TimeConstraints.MaxDurationSeconds > 26*7*24*3600 {
		return nil, errors.New("access rule cannot be longer than 6 months")
	}

	rul := rule.AccessRule{
		ID:          id,
		Approval:    rule.Approval(in.Approval),
		Status:      rule.ACTIVE,
		Description: in.Description,
		Name:        in.Name,
		Groups:      in.Groups,
		Metadata: rule.AccessRuleMetadata{
			CreatedAt: now,
			CreatedBy: userID,
			UpdatedAt: now,
			UpdatedBy: userID,
		},
		Targets:         targets,
		TimeConstraints: in.TimeConstraints,
	}

	log.Debugw("saving access rule", "rule", rul)

	// save the request.
	err = s.DB.Put(ctx, &rul)
	if err != nil {
		return nil, err
	}

	// analytics event
	analytics.FromContext(ctx).Track(&analytics.RuleCreated{
		CreatedBy: userID,
		RuleID:    rul.ID,
		// Provider:           rul.Target.TargetGroupFrom.ToAnalytics(),
		MaxDurationSeconds: in.TimeConstraints.MaxDurationSeconds,

		RequiresApproval: rul.Approval.IsRequired(),
	})

	return &rul, nil
}
