package rulesvc

import (
	"context"
	"errors"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/types"
)

type UpdateOpts struct {
	UpdaterID      string
	Rule           rule.AccessRule
	UpdateRequest  types.CreateAccessRuleRequest
	ApprovalGroups []rule.Approval
}

func (s *Service) UpdateRule(ctx context.Context, in *UpdateOpts) (*rule.AccessRule, error) {

	targets, err := s.ProcessTargets(ctx, in.UpdateRequest.Targets)
	if err != nil {
		return nil, err
	}

	// validate it is under 6 months
	if in.UpdateRequest.TimeConstraints.MaxDurationSeconds > 26*7*24*3600 {
		return nil, errors.New("access rule cannot be longer than 6 months")
	}

	meta := in.Rule.Metadata
	meta.UpdatedAt = s.Clock.Now()
	meta.UpdatedBy = in.UpdaterID
	rul := rule.AccessRule{
		ID:              in.Rule.ID,
		Approval:        rule.Approval(in.UpdateRequest.Approval),
		Description:     in.UpdateRequest.Description,
		Name:            in.UpdateRequest.Name,
		Groups:          in.UpdateRequest.Groups,
		Metadata:        meta,
		Targets:         targets,
		TimeConstraints: in.UpdateRequest.TimeConstraints,
		Priority:        in.UpdateRequest.Priority,
	}

	// updated the previous version to be a version and inserts the new one as current
	err = s.DB.Put(ctx, &rul)
	if err != nil {
		return nil, err
	}

	// analytics event
	analytics.FromContext(ctx).Track(&analytics.RuleUpdated{
		UpdatedBy: in.UpdaterID,
		RuleID:    in.Rule.ID,
		// Provider:           in.Rule.Target.TargetGroupFrom.ToAnalytics(),
		MaxDurationSeconds: in.Rule.TimeConstraints.MaxDurationSeconds,
		RequiresApproval:   in.Rule.Approval.IsRequired(),
	})
	err = s.Cache.RefreshCachedTargets(ctx)
	if err != nil {
		return nil, err
	}
	return &rul, nil
}
