package rulesvc

import (
	"context"
	"errors"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"golang.org/x/sync/errgroup"
)

type UpdateOpts struct {
	UpdaterID      string
	Rule           rule.AccessRule
	UpdateRequest  types.CreateAccessRuleRequest
	ApprovalGroups []rule.Approval
}

func (s *Service) UpdateRule(ctx context.Context, in *UpdateOpts) (*rule.AccessRule, error) {

	//check if user and group exists
	g, gctx := errgroup.WithContext(ctx)
	if in.UpdateRequest.Approval.Users != nil {
		g.Go(func() error {
			for _, u := range *in.UpdateRequest.Approval.Users {

				userLookup := storage.GetUser{ID: u}

				_, err := s.DB.Query(gctx, &userLookup)

				return err
			}
			return nil
		})

	}

	if in.UpdateRequest.Approval.Groups != nil {
		g.Go(func() error {
			for _, u := range *in.UpdateRequest.Approval.Groups {

				groupLookup := storage.GetGroup{ID: u}

				_, err := s.DB.Query(gctx, &groupLookup)

				return err
			}
			return nil
		})
	}

	targets, err := s.ProcessTargets(ctx, in.UpdateRequest.Targets)
	if err != nil {
		return nil, err
	}

	// validate it is under 6 months
	if in.UpdateRequest.TimeConstraints.MaxDurationSeconds > 26*7*24*3600 {
		return nil, errors.New("access rule cannot be longer than 6 months")
	}

	// validate it is under 6 months
	if in.UpdateRequest.TimeConstraints.DefaultDurationSeconds > 26*7*24*3600 {
		return nil, errors.New("access rule cannot be longer than 6 months")
	}

	approvals := rule.Approval{}

	if in.UpdateRequest.Approval.Groups != nil {
		approvals.Groups = *in.UpdateRequest.Approval.Groups
	}

	if in.UpdateRequest.Approval.Users != nil {
		approvals.Users = *in.UpdateRequest.Approval.Users
	}

	meta := in.Rule.Metadata
	meta.UpdatedAt = s.Clock.Now()
	meta.UpdatedBy = in.UpdaterID
	rul := rule.AccessRule{
		ID:              in.Rule.ID,
		Approval:        approvals,
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

	hasFilterExpression := false
	selectedTargets := []string{}
	for _, target := range in.UpdateRequest.Targets {
		selectedTargets = append(selectedTargets, target.TargetGroupId)
		// if len(target.FieldFilterExpessions) > 0 {
		// 	hasFilterExpression = true
		// }
	}

	// analytics event Update access rule
	analytics.FromContext(ctx).Track(&analytics.RuleUpdated{
		UpdatedBy:           in.UpdaterID,
		RuleID:              in.Rule.ID,
		TargetsCount:        len(in.UpdateRequest.Targets),
		HasFilterExpression: hasFilterExpression,
		Targets:             selectedTargets,
		MaxDurationSeconds:  in.Rule.TimeConstraints.MaxDurationSeconds,
		RequiresApproval:    in.Rule.Approval.IsRequired(),
	})

	err = s.Cache.RefreshCachedTargets(ctx)
	if err != nil {
		return nil, err
	}
	return &rul, nil
}
