package rulesvc

import (
	"context"

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
	clk := s.Clock

	target, err := s.ProcessTarget(ctx, in.UpdateRequest.Target)
	if err != nil {
		return nil, err
	}

	// makes a copy of the existing version which will be mutated
	newVersion := in.Rule

	// fields to be updated
	newVersion.Description = in.UpdateRequest.Description
	newVersion.Name = in.UpdateRequest.Name
	newVersion.Approval.Users = in.UpdateRequest.Approval.Users
	newVersion.Approval.Groups = in.UpdateRequest.Approval.Groups
	newVersion.Groups = in.UpdateRequest.Groups
	newVersion.Metadata.UpdatedBy = in.UpdaterID
	newVersion.Metadata.UpdatedAt = clk.Now()
	newVersion.TimeConstraints = in.UpdateRequest.TimeConstraints
	newVersion.Target = target

	// updated the previous version to be a version and inserts the new one as current
	err = s.DB.PutBatch(ctx, &newVersion, &in.Rule)
	if err != nil {
		return nil, err
	}

	// analytics event
	analytics.FromContext(ctx).Track(&analytics.RuleUpdated{
		UpdatedBy:          in.UpdaterID,
		RuleID:             in.Rule.ID,
		Provider:           in.Rule.Target.TargetGroupFrom.ToAnalytics(),
		MaxDurationSeconds: in.Rule.TimeConstraints.MaxDurationSeconds,
		RequiresApproval:   in.Rule.Approval.IsRequired(),
	})

	return &newVersion, nil
}
