package rulesvc

import (
	"context"

	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/types"
)

type UpdateOpts struct {
	UpdaterID      string
	Rule           rule.AccessRule
	UpdateRequest  types.UpdateAccessRuleRequest
	ApprovalGroups []rule.Approval
}

func (s *Service) UpdateRule(ctx context.Context, in *UpdateOpts) (*rule.AccessRule, error) {
	clk := s.Clock
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
	newVersion.Version = types.NewVersionID()

	// Set the existing version to not current
	in.Rule.Current = false

	// updated the previous version to be a version and inserts the new one as current
	err := s.DB.PutBatch(ctx, &newVersion, &in.Rule)
	if err != nil {
		return nil, err
	}

	return &newVersion, nil
}
