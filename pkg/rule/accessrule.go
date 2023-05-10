package rule

import (
	"time"

	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// AccessRules define policies for requesting access to entitlements
type AccessRule struct {
	ID       string             `json:"id" dynamodbav:"id"`
	Priority int                `json:"priority" dynamodbav:"priority"`
	Metadata AccessRuleMetadata `json:"metadata" dynamodbav:"metadata"`

	Name        string `json:"name" dynamodbav:"name"`
	Description string `json:"description" dynamodbav:"description"`

	// The map key is TargetGroupID
	Targets []Target `json:"target" dynamodbav:"target"`

	TimeConstraints types.AccessRuleTimeConstraints `json:"timeConstraints" dynamodbav:"timeConstraints"`
	// Array of group names that the access rule applies to
	Groups []string `json:"groups" dynamodbav:"groups"`
	// Approver config for access rules
	Approval Approval `json:"approval" dynamodbav:"approval"`
}

// AccessRuleMetadata defines model for AccessRuleMetadata.
type AccessRuleMetadata struct {
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
	// userID
	CreatedBy string `json:"createdBy" dynamodbav:"createdBy"`
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

type FieldFilterExpessions struct{}
type Target struct {
	TargetGroup           target.Group                     `json:"targetGroup" dynamodbav:"targetGroup"`
	FieldFilterExpessions map[string]FieldFilterExpessions `json:"fieldFilterExpessions" dynamodbav:"fieldFilterExpessions"`
}

// // ised for admin apis, this contains the access rule target in a format for updating the access rule provider target
func (a AccessRule) ToAPI() types.AccessRule {
	approval := types.AccessRuleApproverConfig{
		Groups: []string{},
		Users:  []string{},
	}
	if a.Approval.Groups != nil {
		approval.Groups = a.Approval.Groups
	}
	if a.Approval.Users != nil {
		approval.Users = a.Approval.Users
	}

	targets := []types.AccessRuleTarget{}

	for _, target := range a.Targets {
		targets = append(targets, target.ToAPI())
	}
	return types.AccessRule{
		ID:          a.ID,
		Description: a.Description,
		Name:        a.Name,
		Metadata: types.AccessRuleMetadata{
			CreatedAt: a.Metadata.CreatedAt,
			UpdatedAt: a.Metadata.UpdatedAt,
			CreatedBy: a.Metadata.CreatedBy,
			UpdatedBy: a.Metadata.UpdatedBy,
		},
		Groups: a.Groups,
		TimeConstraints: types.AccessRuleTimeConstraints{
			MaxDurationSeconds: a.TimeConstraints.MaxDurationSeconds,
		},
		Approval: approval,
		Targets:  targets,
		Priority: a.Priority,
	}
}

// converts to basic api type
func (t Target) ToAPI() types.AccessRuleTarget {
	return types.AccessRuleTarget{
		FieldFilterExpessions: map[string]interface{}{},
		TargetGroup:           t.TargetGroup.ToAPI(),
	}
}

func (r *AccessRule) DDBKeys() (ddb.Keys, error) {
	return ddb.Keys{
		PK:     keys.AccessRule.PK1,
		SK:     keys.AccessRule.SK1(r.ID),
		GSI1PK: keys.AccessRule.GSI1PK,
		GSI1SK: keys.AccessRule.GSI1SK(r.Priority, r.ID),
	}, nil

}
