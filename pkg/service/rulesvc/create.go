package rulesvc

import (
	"context"
	"net/http"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/provider"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/pkg/errors"
)

// validateTargetAgainstSchema checks that all the arguments match the schema of the provider
// It validates that all required arguments were provided with at least 1 value
// returns apio.APIError so it will bubble up as a 400 error from api usage
func validateTargetAgainstSchema(in types.CreateAccessRuleTarget, providerArgSchema map[string]provider.Argument) error {
	if len(providerArgSchema) != len(in.With) {
		return apio.NewRequestError(errors.New("target is missing required arguments from the provider schema"), http.StatusBadRequest)
	}
	for argumentID, argument := range in.With {
		hasAtLeastOneValue := len(argument.Values) != 0
		argumentSchema, ok := providerArgSchema[argumentID]
		if !ok {
			return apio.NewRequestError(errors.New("argument does not match schema for provider"), http.StatusBadRequest)
		}
		// filter any group options which do not have any values
		for groupId, group := range argument.Groupings {
			if _, ok := (argumentSchema.Groups)[groupId]; !ok {
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
func (s *Service) validateTargetArgumentAgainstCachedOptions(ctx context.Context, in types.CreateAccessRuleTarget, providerArgSchema map[string]provider.Argument) error {
	for argumentID, argument := range in.With {
		if providerArgSchema[argumentID].RuleFormElement != types.ArgumentRuleFormElementINPUT {
			_, argOptions, groupOptions, err := s.Cache.LoadCachedProviderArgOptions(ctx, in.ProviderId, argumentID)
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

			for groupId, groupValues := range argument.Groupings {
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
	// After verifying the provider, we can save the provider type to the rule for convenience

	q := storage.GetProvider{
		ID: in.ProviderId,
	}
	_, err := s.DB.Query(ctx, &q)
	if err != nil {
		return rule.Target{}, err
	}
	err = validateTargetAgainstSchema(in, q.Result.Schema)
	if err != nil {
		return rule.Target{}, err
	}
	err = s.validateTargetArgumentAgainstCachedOptions(ctx, in, q.Result.Schema)
	if err != nil {
		return rule.Target{}, err
	}
	target := rule.Target{
		ProviderID:               in.ProviderId,
		ProviderType:             q.Result.Type,
		With:                     make(map[string]string),
		WithSelectable:           make(map[string][]string),
		WithArgumentGroupOptions: make(map[string]map[string][]string),
	}

	for argumentID, argument := range in.With {
		for groupId, groupValues := range argument.Groupings {
			if len(groupValues) > 0 {

				argumentGroupOptions := target.WithArgumentGroupOptions[argumentID]
				if argumentGroupOptions == nil {
					argumentGroupOptions = make(map[string][]string)
				}
				argumentGroupOptions[groupId] = groupValues
				target.WithArgumentGroupOptions[argumentID] = argumentGroupOptions
			}
		}
		if len(argument.Values) == 1 {
			target.With[argumentID] = argument.Values[0]
		} else {
			target.WithSelectable[argumentID] = argument.Values
		}
	}

	return target, nil
}

func (s *Service) CreateAccessRule(ctx context.Context, user *identity.User, in types.CreateAccessRuleRequest) (*rule.AccessRule, error) {
	id := types.NewAccessRuleID()

	log := logger.Get(ctx).With("user.id", user.ID, "access_rule.id", id)
	now := s.Clock.Now()

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
			CreatedBy: user.ID,
			UpdatedAt: now,
			UpdatedBy: user.ID,
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
		CreatedBy:             user.ID,
		RuleID:                rul.ID,
		Provider:              rul.Target.ProviderType,
		MaxDurationSeconds:    in.TimeConstraints.MaxDurationSeconds,
		UsesSelectableOptions: rul.Target.UsesSelectableOptions(),
		UsesDynamicOptions:    rul.Target.UsesDynamicOptions(),
		RequiresApproval:      rul.Approval.IsRequired(),
	})

	return &rul, nil
}
