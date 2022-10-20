package rule

import (
	"time"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providerregistry"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	"github.com/common-fate/granted-approvals/pkg/types"
)

// AccessRule is a rule governing access to something in Granted.
//
// Access Rules have versions.
// When updating an access rule, you need to update the current version with Current = false
// and then insert the new version with Current = true
// This will correctly set the keys and enable the access patterns
type AccessRule struct {
	// Current is true if this is the current version
	// When a new version is added, the previous version should be updated to set Current to false
	Current bool `json:"current" dynamodbav:"current"`
	// Approver config for access rules
	Approval    Approval `json:"approval" dynamodbav:"approval"`
	Version     string   `json:"version" dynamodbav:"version"`
	Status      Status   `json:"status" dynamodbav:"status"`
	Description string   `json:"description" dynamodbav:"description"`

	// Array of group names that the access rule applies to
	Groups          []string              `json:"groups" dynamodbav:"groups"`
	ID              string                `json:"id" dynamodbav:"id"`
	Metadata        AccessRuleMetadata    `json:"metadata" dynamodbav:"metadata"`
	Name            string                `json:"name" dynamodbav:"name"`
	Target          Target                `json:"target" dynamodbav:"target"`
	TimeConstraints types.TimeConstraints `json:"timeConstraints" dynamodbav:"timeConstraints"`
}

// ised for admin apis, this contains the access rule target in a format for updating the access rule provider target
func (a AccessRule) ToAPIDetail() types.AccessRuleDetail {
	status := types.AccessRuleStatusACTIVE
	if a.Status == ARCHIVED {
		status = types.AccessRuleStatusARCHIVED
	}
	approval := types.ApproverConfig{
		Groups: []string{},
		Users:  []string{},
	}
	if a.Approval.Groups != nil {
		approval.Groups = a.Approval.Groups
	}
	if a.Approval.Users != nil {
		approval.Users = a.Approval.Users
	}
	return types.AccessRuleDetail{
		ID:          a.ID,
		Description: a.Description,
		Name:        a.Name,
		Metadata: types.AccessRuleMetadata{
			CreatedAt:     a.Metadata.CreatedAt,
			UpdatedAt:     a.Metadata.UpdatedAt,
			UpdateMessage: a.Metadata.UpdateMessage,
			CreatedBy:     a.Metadata.CreatedBy,
			UpdatedBy:     a.Metadata.UpdatedBy,
		},
		Groups: a.Groups,
		TimeConstraints: types.TimeConstraints{
			MaxDurationSeconds: a.TimeConstraints.MaxDurationSeconds,
		},
		Approval: approval,

		Target: a.Target.ToAPIDetail(),

		Status:    status,
		Version:   a.Version,
		IsCurrent: a.Current,
	}
}

// served basic detail of the access rule
func (a AccessRule) ToAPI() types.AccessRule {

	return types.AccessRule{
		ID:          a.ID,
		Version:     a.Version,
		Description: a.Description,
		Name:        a.Name,
		TimeConstraints: types.TimeConstraints{
			MaxDurationSeconds: a.TimeConstraints.MaxDurationSeconds,
		},

		Target:    a.Target.ToAPI(),
		IsCurrent: a.Current,
	}
}

// This is used to serve a user making a request, it contains all the available arguments and options with title, description and labels
func (a AccessRule) ToRequestAccessRuleAPI(requestArguments map[string]types.RequestArgument) types.RequestAccessRule {
	return types.RequestAccessRule{
		Version:     a.Version,
		Description: a.Description,
		Name:        a.Name,
		IsCurrent:   a.Current,
		ID:          a.ID,
		Target: types.RequestAccessRuleTarget{
			Provider: a.Target.ProviderToAPI(),
			Arguments: types.RequestAccessRuleTarget_Arguments{
				AdditionalProperties: requestArguments,
			},
		},
		TimeConstraints: a.TimeConstraints,
	}
}

// AccessRuleMetadata defines model for AccessRuleMetadata.
type AccessRuleMetadata struct {
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	// userID
	CreatedBy      string                  `json:"createdBy" dynamodbav:"createdBy"`
	UpdateMessage  *string                 `json:"updateMessage,omitempty" dynamodbav:"updateMessage,omitempty"`
	UpdateMetadata *map[string]interface{} `json:"updateMetadata,omitempty" dynamodbav:"updateMetadata,omitempty"`
	UpdatedAt      time.Time               `json:"updatedAt" dynamodbav:"updatedAt"`
	// userID
	UpdatedBy string `json:"updatedBy" dynamodbav:"updatedBy"`
}

// Approver config for access rules
type Approval struct {
	// List of group ids represents the groups whos members may approver requests for this rule
	Groups []string `json:"groups" dynamodbav:"groups"`
	//List of users ids represents the individual users who may approve requests for this rule.
	// This does not represent members of the approval groups
	Users []string `json:"users" dynamodbav:"users"`
}

func (a *Approval) IsRequired() bool {
	return len(a.Users) > 0 || len(a.Groups) > 0
}

// Provider defines model for Provider.
// I expect this will be different to what gets returned in the api response
type Target struct {
	// References the provider's unique ID
	ProviderID   string            `json:"providerId"  dynamodbav:"providerId"`
	ProviderType string            `json:"providerType"  dynamodbav:"providerType"`
	With         map[string]string `json:"with"  dynamodbav:"with"`
	// when target can have multiple values
	WithSelectable map[string][]string `json:"withSelectable"  dynamodbav:"withSelectable"`
	// when target doesn't have values but instead belongs to a group
	// which can be dynamically fetched at access request time.
	WithArgumentGroupOptions map[string]map[string][]string `json:"withArgumentGroupOptions"  dynamodbav:"withArgumentGroupOptions"`
}

func (t Target) ProviderToAPI() types.Provider {
	return types.Provider{
		Id:   t.ProviderID,
		Type: t.ProviderType,
	}
}

// converts to basic api type
func (t Target) ToAPI() types.AccessRuleTarget {
	return types.AccessRuleTarget{
		Provider: t.ProviderToAPI(),
	}
}

func (t Target) ToAPIDetail() types.AccessRuleTargetDetail {
	at := types.AccessRuleTargetDetail{
		Provider: types.Provider{
			Id:   t.ProviderID,
			Type: t.ProviderType,
		},
		With: types.AccessRuleTargetDetail_With{
			AdditionalProperties: make(map[string]types.AccessRuleTargetDetailArguments),
		},
	}
	// Lookup the provider, ignore errors
	// if provider is not found, fallback to using the argument key as the title
	_, provider, _ := providerregistry.Registry().GetLatestByShortType(t.ProviderType)

	for k, v := range t.With {
		argument := at.With.AdditionalProperties[k]
		argument.Values = append(argument.Values, v)

		// FormElement is used here when loading the UpdateAccessRule form for the first time
		if provider != nil {
			if s, ok := provider.Provider.(providers.ArgSchemarer); ok {
				schema := s.ArgSchema()
				if arg, ok := schema[k]; ok {
					argument.FormElement = types.AccessRuleTargetDetailArgumentsFormElement(arg.FormElement)
				} else {
					// I don't expect this should ever fail to find a match, however if it does, default to input.
					argument.FormElement = types.INPUT
				}
			}
		}
		at.With.AdditionalProperties[k] = argument
	}
	for k, v := range t.WithSelectable {
		argument := at.With.AdditionalProperties[k]
		argument.Values = append(argument.Values, v...)
		// FormElement is used here when loading the UpdateAccessRule form for the first time
		if provider != nil {
			if s, ok := provider.Provider.(providers.ArgSchemarer); ok {
				schema := s.ArgSchema()
				if arg, ok := schema[k]; ok {
					argument.FormElement = types.AccessRuleTargetDetailArgumentsFormElement(arg.FormElement)
				} else {
					// I don't expect this should ever fail to find a match, however if it does, default to input.
					argument.FormElement = types.INPUT
				}
			}
		}
		at.With.AdditionalProperties[k] = argument
	}
	for k, v := range t.WithArgumentGroupOptions {
		argument := at.With.AdditionalProperties[k]
		argument.Groupings.AdditionalProperties = make(map[string][]string)
		for k2, v2 := range v {
			group := argument.Groupings.AdditionalProperties[k2]
			group = append(group, v2...)
			argument.Groupings.AdditionalProperties[k2] = group
		}
		// FormElement is used here when loading the UpdateAccessRule form for the first time
		if provider != nil {
			if s, ok := provider.Provider.(providers.ArgSchemarer); ok {
				schema := s.ArgSchema()
				if arg, ok := schema[k]; ok {
					argument.FormElement = types.AccessRuleTargetDetailArgumentsFormElement(arg.FormElement)
				} else {
					// I don't expect this should ever fail to find a match, however if it does, default to input.
					argument.FormElement = types.INPUT
				}
			}
		}
		at.With.AdditionalProperties[k] = argument
	}

	return at
}

func (r *AccessRule) DDBKeys() (ddb.Keys, error) {
	// If this is a current version of the rule, then the GSI keys are used
	if r.Current {
		return ddb.Keys{
			PK:     keys.AccessRule.PK1,
			SK:     keys.AccessRule.SK1(r.ID, r.Version),
			GSI1PK: keys.AccessRule.GSI1PK(string(r.Status)),
			GSI1SK: keys.AccessRule.GSI1SK(r.ID),
			GSI2PK: keys.AccessRule.GSI2PK,
			GSI2SK: keys.AccessRule.GSI2SK(r.ID),
		}, nil
	}
	// If this is not a current version of the rule, only the primary keys are used
	return ddb.Keys{
		PK: keys.AccessRule.PK1,
		SK: keys.AccessRule.SK1(r.ID, r.Version),
	}, nil
}
