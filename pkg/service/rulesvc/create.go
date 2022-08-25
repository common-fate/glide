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

func (s *Service) CreateAccessRule(ctx context.Context, user *identity.User, in types.CreateAccessRuleRequest) (*rule.AccessRule, error) {
	id := types.NewAccessRuleID()

	log := logger.Get(ctx).With("user.id", user.ID, "access_rule.id", id)
	now := s.Clock.Now()

	// After verifying the provider, we can save the provider type to the rule for convenience
	p, err := s.verifyRuleTarget(ctx, in.Target)
	if err != nil {
		return nil, err
	}

	target := rule.Target{
		ProviderID:     in.Target.ProviderId,
		ProviderType:   p.Type,
		With:           make(map[string]string),
		WithSelectable: make(map[string][]rule.Selectable),
	}

	for k, values := range in.Target.With.AdditionalProperties {
		// min length 1 is configured in the api spec so len(0) is handled by builtin validation
		if len(values) == 1 {
			target.With[k] = values[0]
		} else {
			// store the selectables with value and label
			target.WithSelectable[k] = make([]rule.Selectable, len(values))
			for i, value := range values {
				target.WithSelectable[k][i] = rule.Selectable{
					// @todo need to fetch the options again and map values to options to get the label, but its really shlow so leaving it out for now
					Option: rule.Option{Value: value, Label: value},
					Valid:  true,
				}
			}
		}
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
func (s *Service) verifyRuleTarget(ctx context.Context, target types.CreateAccessRuleTarget) (*ahTypes.Provider, error) {
	p, err := s.AHClient.GetProviderWithResponse(ctx, target.ProviderId)
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
