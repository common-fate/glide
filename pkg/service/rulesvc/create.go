package rulesvc

import (
	"context"
	"errors"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"golang.org/x/sync/errgroup"
)

func (s *Service) ProcessTargets(ctx context.Context, in []types.CreateAccessRuleTarget) ([]rule.Target, error) {
	// @TODO

	// Check for duplicate target groups, there can only be one instance of a target group per access rule.
	deduplicateTargets := map[string]rule.Target{}

	for _, targetGroup := range in {
		//check if target group exists
		_, ok := deduplicateTargets[targetGroup.TargetGroupId]
		if !ok {
			targetGroupQ := storage.GetTargetGroup{ID: targetGroup.TargetGroupId}
			_, err := s.DB.Query(ctx, &targetGroupQ)
			if err != nil {
				return nil, err
			}

			deduplicateTargets[targetGroupQ.Result.ID] = rule.Target{
				TargetGroup:           *targetGroupQ.Result,
				FieldFilterExpessions: map[string]rule.FieldFilterExpessions{},
			}
		} else {
			//do we want to error out here or just deduplicate it automatically?
			return nil, errors.New("duplicate target in access rule")
		}
	}

	targets := []rule.Target{}
	for _, tg := range deduplicateTargets {
		targets = append(targets, tg)
	}

	// TODO when filter expressions are implemented
	// Check validity of filter expressions (the structure of a filter expression shoudl be validated in the API layer automatically)
	// - check that attributes in the filter expression exist on the schema of the target/resource type

	return targets, nil
}

func (s *Service) CreateAccessRule(ctx context.Context, userID string, in types.CreateAccessRuleRequest) (*rule.AccessRule, error) {

	id := types.NewAccessRuleID()

	log := logger.Get(ctx).With("user.id", userID, "access_rule.id", id)
	now := s.Clock.Now()
	g, gctx := errgroup.WithContext(ctx)

	// check if user and group exists
	if in.Approval.Users != nil {
		g.Go(func() error {
			for _, u := range *in.Approval.Users {

				userLookup := storage.GetUser{ID: u}

				_, err := s.DB.Query(gctx, &userLookup)

				return err
			}
			return nil
		})
	}

	if in.Approval.Groups != nil {
		g.Go(func() error {
			for _, u := range *in.Approval.Groups {

				groupLookup := storage.GetGroup{ID: u}

				_, err := s.DB.Query(ctx, &groupLookup)

				return err
			}
			return nil
		})
	}

	targets, err := s.ProcessTargets(ctx, in.Targets)
	if err != nil {
		return nil, err
	}

	// validate it is under 6 months
	if in.TimeConstraints.MaxDurationSeconds > 26*7*24*3600 {
		return nil, errors.New("access rule cannot be longer than 6 months")
	}
	if in.TimeConstraints.DefaultDurationSeconds > 26*7*24*3600 {
		return nil, errors.New("access rule cannot be longer than 6 months")
	}

	approvals := rule.Approval{}

	if in.Approval.Groups != nil {
		approvals.Groups = *in.Approval.Groups
	}

	if in.Approval.Users != nil {
		approvals.Users = *in.Approval.Users
	}

	rul := rule.AccessRule{
		ID:          id,
		Approval:    approvals,
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
		Priority:        in.Priority,
	}

	log.Debugw("saving access rule", "rule", rul)

	// save the request.
	err = s.DB.Put(ctx, &rul)
	if err != nil {
		return nil, err
	}

	hasFilterExpression := false
	selectedTargets := []string{}
	for _, target := range in.Targets {
		selectedTargets = append(selectedTargets, target.TargetGroupId)
		if len(target.FieldFilterExpessions) > 0 {
			hasFilterExpression = true
		}
	}

	// analytics event Create Access Rule
	analytics.FromContext(ctx).Track(&analytics.RuleCreated{
		CreatedBy:           userID,
		RuleID:              rul.ID,
		RequiresApproval:    rul.Approval.IsRequired(),
		HasFilterExpression: hasFilterExpression,
		TargetsCount:        len(in.Targets),
		Targets:             selectedTargets,
		MaxDurationSeconds:  in.TimeConstraints.MaxDurationSeconds,
	})

	err = s.Cache.RefreshCachedTargets(ctx)
	if err != nil {
		return nil, err
	}

	return &rul, nil
}
