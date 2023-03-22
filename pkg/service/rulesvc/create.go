package rulesvc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/pkg/errors"
)

// validateTargetAgainstSchema checks that all the arguments match the schema of the provider
// It validates that all required arguments were provided with at least 1 value
// returns apio.APIError so it will bubble up as a 400 error from api usage
func validateTargetAgainstSchema(in types.CreateAccessRuleTarget, providerArgSchema *types.ArgSchema) error {
	if len(providerArgSchema.AdditionalProperties) != len(in.With.AdditionalProperties) {
		return apio.NewRequestError(errors.New("target is missing required arguments from the provider schema"), http.StatusBadRequest)
	}
	for argumentID, argument := range in.With.AdditionalProperties {
		hasAtLeastOneValue := len(argument.Values) != 0
		argumentSchema, ok := providerArgSchema.AdditionalProperties[argumentID]
		if !ok {
			return apio.NewRequestError(errors.New("argument does not match schema for provider"), http.StatusBadRequest)
		}
		// filter any group options which do not have any values
		for groupId, group := range argument.Groupings.AdditionalProperties {
			if _, ok := argumentSchema.Groups.AdditionalProperties[groupId]; !ok {
				return apio.NewRequestError(errors.New("argument group does not match schema for provider"), http.StatusBadRequest)
			}
			if len(group) != 0 {
				hasAtLeastOneValue = true
			}
		}
		if !hasAtLeastOneValue {
			return apio.NewRequestError(errors.New("arguments must have at least 1 value or group value"), http.StatusBadRequest)
		}
	}
	return nil
}

// validateTargetArgumentAgainstCachedOptions checks that all the argument values and argument group values currently exist in the cache.
// this prevents being able to create an access rule with arguments which are invalid for the provider.
// returns apio.APIError so it will bubble up as a 400 error from api usage
func (s *Service) validateTargetArgumentAgainstCachedOptions(ctx context.Context, in types.CreateAccessRuleTarget, providerArgSchema *types.ArgSchema) error {
	for argumentID, argument := range in.With.AdditionalProperties {
		if providerArgSchema.AdditionalProperties[argumentID].RuleFormElement != types.ArgumentRuleFormElementINPUT {
			_, argOptions, groupOptions, err := s.Cache.RefreshCachedProviderArgOptions(ctx, in.ProviderId, argumentID)
			// _, argOptions, groupOptions, err := s.Cache.LoadCachedProviderArgOptions(ctx, in.ProviderId, argumentID)
			if err != nil {
				return err
			}
			groupOptionsValueMap := make(map[string]map[string]string)
			argOptionsValueMap := make(map[string]string)
			for _, arg := range argOptions {
				argOptionsValueMap[arg.Value] = arg.Value
			}
			for _, group := range groupOptions {
				options := groupOptionsValueMap[group.Group]
				if options == nil {
					options = make(map[string]string)
				}
				options[group.Value] = group.Value
				groupOptionsValueMap[group.Group] = options
			}

			for groupId, groupValues := range argument.Groupings.AdditionalProperties {
				if len(groupValues) > 0 {
					if _, ok := groupOptionsValueMap[groupId]; !ok {
						return apio.NewRequestError(errors.New("argument group values do not match available options for provider"), http.StatusBadRequest)
					}
					for _, value := range groupValues {
						if _, ok := groupOptionsValueMap[groupId][value]; !ok {
							return apio.NewRequestError(errors.New("argument group values do not match available options for provider"), http.StatusBadRequest)
						}
					}
				}
			}
			for _, value := range argument.Values {
				if _, ok := argOptionsValueMap[value]; !ok {
					return apio.NewRequestError(errors.New("argument values do not match available options for provider"), http.StatusBadRequest)
				}
			}
		}
	}
	return nil
}
func (s *Service) ProcessTarget(ctx context.Context, in types.CreateAccessRuleTarget) (rule.Target, error) {

	targetgroup := rule.Target{
		ProviderID:               in.ProviderId,
		BuiltInProviderType:      "",
		TargetGroupID:            in.ProviderId,
		With:                     make(map[string]string),
		WithSelectable:           make(map[string][]string),
		WithArgumentGroupOptions: make(map[string]map[string][]string),
	}

	q := storage.GetTargetGroup{ID: in.ProviderId}
	_, err := s.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		return rule.Target{}, err
	}

	targetgroup.TargetGroupFrom = q.Result.From

	for argumentID, argument := range in.With.AdditionalProperties {
		// check if the provided argId is a valid argument id in TargetGroup's schema.
		arg, ok := q.Result.Schema.Properties[argumentID]
		if !ok {
			return rule.Target{}, apio.NewRequestError(fmt.Errorf("argument '%s' does not match schema for targetgroup '%s'", argumentID, in.ProviderId), http.StatusBadRequest)
		}

		if len(argument.Values) < 1 {
			return rule.Target{}, apio.NewRequestError(errors.New("argument must have associated value with it"), http.StatusBadRequest)
		}

		if arg.Resource != nil {
			qGetResourcesForTG := storage.ListCachedTargetGroupResource{TargetGroupID: in.ProviderId, ResourceType: *arg.Resource}
			_, err := s.DB.Query(ctx, &qGetResourcesForTG)
			if err != nil {
				return rule.Target{}, err
			}

			// check if the provided value is a valid resource id in cached resources.
			for _, providedValue := range argument.Values {
				isValidArgValue := false
				for _, cachedResource := range qGetResourcesForTG.Result {
					if cachedResource.Resource.ID == providedValue {
						isValidArgValue = true
					}
				}

				if !isValidArgValue {
					return rule.Target{}, apio.NewRequestError(fmt.Errorf("invalid argument value '%s' provided for argument '%s'", providedValue, argumentID), http.StatusBadRequest)
				}
			}
		}

		if len(argument.Values) == 1 {
			targetgroup.With[argumentID] = argument.Values[0]
		} else {
			targetgroup.WithSelectable[argumentID] = argument.Values
		}
	}

	return targetgroup, nil

}

func (s *Service) CreateAccessRule(ctx context.Context, userID string, in types.CreateAccessRuleRequest) (*rule.AccessRule, error) {
	id := types.NewAccessRuleID()

	log := logger.Get(ctx).With("user.id", userID, "access_rule.id", id)
	now := s.Clock.Now()

	q := storage.GetTargetGroup{ID: in.Target.ProviderId}
	_, err := s.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		return nil, err
	}

	target, err := s.ProcessTarget(ctx, in.Target)
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
		Target:          target,
		TimeConstraints: in.TimeConstraints,
		Version:         types.NewVersionID(),
		Current:         true,
	}

	log.Debugw("saving access rule", "rule", rul)

	// save the request.
	err = s.DB.Put(ctx, &rul)
	if err != nil {
		return nil, err
	}

	// analytics event
	analytics.FromContext(ctx).Track(&analytics.RuleCreated{
		CreatedBy:             userID,
		RuleID:                rul.ID,
		BuiltInProvider:       rul.Target.BuiltInProviderType,
		Provider:              rul.Target.TargetGroupFrom.ToAnalytics(),
		PDKProvider:           rul.Target.IsForTargetGroup(),
		MaxDurationSeconds:    in.TimeConstraints.MaxDurationSeconds,
		UsesSelectableOptions: rul.Target.UsesSelectableOptions(),
		UsesDynamicOptions:    rul.Target.UsesDynamicOptions(),
		RequiresApproval:      rul.Approval.IsRequired(),
	})

	return &rul, nil
}
