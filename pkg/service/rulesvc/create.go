package rulesvc

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/common-fate/apikit/logger"
	ahTypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/pkg/errors"
)

func (s *Service) ProcessTarget(ctx context.Context, in types.AccessRuleTarget) (rule.Target, error) {
	// After verifying the provider, we can save the provider type to the rule for convenience
	p, err := s.verifyRuleTarget(ctx, in)
	if err != nil {
		return rule.Target{}, err
	}
	target := rule.Target{
		ProviderID:               in.Provider.Id,
		ProviderType:             p.Type,
		With:                     make(map[string]string),
		WithSelectable:           make(map[string][]string),
		WithArgumentGroupOptions: make(map[string]map[string][]string),
	}

	for k, obj := range in.With.AdditionalProperties {
		nestedMap := make(map[string][]string)
		for groupId, groupValues := range obj.Groupings.AdditionalProperties {

			// only add key to map if there are values for the key.
			if len(groupValues) > 0 {
				nestedMap[groupId] = groupValues
				target.WithArgumentGroupOptions[k] = nestedMap
			}
		}

		// min length 1 is configured in the api spec so len(0) is handled by builtin validation
		// if there is no dynamic fields and only 1 value then add to `With` field else add values to `WithSelectable`.
		if len(obj.Values) == 1 && len(target.WithArgumentGroupOptions[k]) == 0 {
			target.With[k] = obj.Values[0]
		} else {
			target.WithSelectable[k] = obj.Values
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

	return &rul, nil
}

// verifyRuleTarget fetches the provider and returns it if it exists
func (s *Service) verifyRuleTarget(ctx context.Context, target types.AccessRuleTarget) (*ahTypes.Provider, error) {
	p, err := s.AHClient.GetProviderWithResponse(ctx, target.Provider.Id)
	if err != nil {
		return nil, err
	}
	switch p.StatusCode() {
	case http.StatusOK:
		return p.JSON200, nil
	case http.StatusNotFound:
		return nil, ErrProviderNotFound
	case http.StatusInternalServerError:
		return nil, errors.Wrap(errors.New(aws.ToString(p.JSON500.Error)), "error while verifying target exists in access handler")
	}
	return nil, ErrUnhandledResponseFromAccessHandler
}
