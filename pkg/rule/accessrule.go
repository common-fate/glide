package rule

import (
	"time"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// AccessRule is a rule governing access to something in Common Fate.
//
// Access Rules have versions.
// When updating an access rule, you need to update the current version with Current = false
// and then insert the new version with Current = true
// This will correctly set the keys and enable the access patterns
type AccessRule struct {

	// Approver config for access rules
	Approval Approval `json:"approval" dynamodbav:"approval"`

	Status      Status `json:"status" dynamodbav:"status"`
	Description string `json:"description" dynamodbav:"description"`

	// Array of group names that the access rule applies to
	Groups   []string           `json:"groups" dynamodbav:"groups"`
	ID       string             `json:"id" dynamodbav:"id"`
	Metadata AccessRuleMetadata `json:"metadata" dynamodbav:"metadata"`
	Name     string             `json:"name" dynamodbav:"name"`
	Target   Target             `json:"target" dynamodbav:"target"`
	// @TODO make a single field for targets
	Targets         map[string]Target
	TimeConstraints types.TimeConstraints `json:"timeConstraints" dynamodbav:"timeConstraints"`
}

// Inherit rule and include `canRequest` field
// which is used to determine if the approval can request the rule or not.
type GetAccessRuleResponse struct {
	Rule       *AccessRule
	CanRequest bool
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

		Status: status,
	}
}

// served basic detail of the access rule
func (a AccessRule) ToAPI() types.AccessRule {

	return types.AccessRule{
		ID:          a.ID,
		Description: a.Description,
		Name:        a.Name,
		TimeConstraints: types.TimeConstraints{
			MaxDurationSeconds: a.TimeConstraints.MaxDurationSeconds,
		},

		Target:    a.Target.ToAPI(),
		CreatedAt: a.Metadata.CreatedAt,
		UpdatedAt: a.Metadata.UpdatedAt,
	}
}

// This is used to serve a user making a request, it contains all the available arguments and options with title, description and labels
func (a AccessRule) ToRequestAccessRuleAPI(requestArguments map[string]types.RequestArgument, canRequest bool) types.RequestAccessRule {
	return types.RequestAccessRule{
		Description: a.Description,
		Name:        a.Name,
		ID:          a.ID,
		Target: types.RequestAccessRuleTarget{
			Arguments: types.RequestAccessRuleTarget_Arguments{
				AdditionalProperties: requestArguments,
			},
		},
		TimeConstraints: a.TimeConstraints,
		CanRequest:      canRequest,
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
	TargetGroupID string `json:"targetGroupId" dynamodbav:"targetGroupId"`

	// TargetGroupFrom is only used for PDK providers and is a denormalised copy of the
	// 'From' field in a Target Group.
	TargetGroupFrom target.From `json:"targetGroupFrom"  dynamodbav:"targetGroupFrom"`
	// Schema is denomalised and saved here for efficiency
	Schema providerregistrysdk.Target `json:"schema" dynamodbav:"schema"`
	With   map[string]string          `json:"with"  dynamodbav:"with"`
}

// converts to basic api type
func (t Target) ToAPI() types.AccessRuleTarget {
	return types.AccessRuleTarget{}
}

func (t Target) ToAPIDetail() types.AccessRuleTargetDetail {

	at := types.AccessRuleTargetDetail{

		With: types.AccessRuleTargetDetail_With{
			AdditionalProperties: make(map[string]types.AccessRuleTargetDetailArguments),
		},
	}

	at.TargetGroup = types.TargetGroup{Id: t.TargetGroupID, From: t.TargetGroupFrom.ToAPI()}

	for k, v := range t.With {
		argument := at.With.AdditionalProperties[k]
		argument.Values = append(argument.Values, v)

		at.With.AdditionalProperties[k] = argument
	}

	// It is essential that all slices be initialised for the apitypes otherwise it will be serialised as null instead of empty
	for k, v := range at.With.AdditionalProperties {
		if v.Values == nil {
			v.Values = make([]string, 0)
		}
		at.With.AdditionalProperties[k] = v
	}

	return at

}

func (r *AccessRule) DDBKeys() (ddb.Keys, error) {
	// If this is a current version of the rule, then the GSI keys are used

	return ddb.Keys{
		PK: keys.AccessRule.PK1,
		SK: keys.AccessRule.SK1(r.ID),
	}, nil

}
