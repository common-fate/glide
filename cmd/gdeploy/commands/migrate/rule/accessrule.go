package rule

import (
	"time"
)

// Status is the status of an Access Rule.
type Status string

const (
	ACTIVE   Status = "ACTIVE"
	ARCHIVED Status = "ARCHIVED"
)

// Time configuration for an Access Rule.
type TimeConstraints struct {
	// The maximum duration in seconds the access is allowed for.
	MaxDurationSeconds int `json:"maxDurationSeconds"`
}

// AccessRule is a rule governing access to something in Common Fate.
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
	Groups          []string           `json:"groups" dynamodbav:"groups"`
	ID              string             `json:"id" dynamodbav:"id"`
	Metadata        AccessRuleMetadata `json:"metadata" dynamodbav:"metadata"`
	Name            string             `json:"name" dynamodbav:"name"`
	Target          Target             `json:"target" dynamodbav:"target"`
	TimeConstraints TimeConstraints    `json:"timeConstraints" dynamodbav:"timeConstraints"`
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
type From struct {
	Publisher string `json:"publisher" dynamodbav:"publisher"`
	Name      string `json:"name" dynamodbav:"name"`
	Version   string `json:"version" dynamodbav:"version"`
	Kind      string `json:"kind" dynamodbav:"kind"`
}

// Provider defines model for Provider.
// I expect this will be different to what gets returned in the api response
type Target struct {
	// References the provider's unique ID
	ProviderID    string `json:"providerId"  dynamodbav:"providerId"`
	TargetGroupID string `json:"targetGroupId" dynamodbav:"targetGroupId"`

	// BuiltInProviderType is only used for built-in providers
	BuiltInProviderType string `json:"providerType"  dynamodbav:"providerType"`

	// TargetGroupFrom is only used for PDK providers and is a denormalised copy of the
	// 'From' field in a Target Group.
	TargetGroupFrom From `json:"targetGroupFrom"  dynamodbav:"targetGroupFrom"`

	With map[string]string `json:"with"  dynamodbav:"with"`
	// when target can have multiple values
	WithSelectable map[string][]string `json:"withSelectable"  dynamodbav:"withSelectable"`
	// when target doesn't have values but instead belongs to a group
	// which can be dynamically fetched at access request time.
	WithArgumentGroupOptions map[string]map[string][]string `json:"withArgumentGroupOptions"  dynamodbav:"withArgumentGroupOptions"`
}
