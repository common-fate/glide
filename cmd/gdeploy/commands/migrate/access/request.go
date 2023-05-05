package access

import (
	"time"
)

// Status of an Access Request.
type Status string

const (
	APPROVED  Status = "APPROVED"
	DECLINED  Status = "DECLINED"
	CANCELLED Status = "CANCELLED"
	PENDING   Status = "PENDING"
)

type Grant_With struct {
	AdditionalProperties map[string]string `json:"-"`
}

// The current state of the grant.
type GrantStatus string

// Defines values for GrantStatus.
const (
	GrantStatusACTIVE  GrantStatus = "ACTIVE"
	GrantStatusERROR   GrantStatus = "ERROR"
	GrantStatusEXPIRED GrantStatus = "EXPIRED"
	GrantStatusPENDING GrantStatus = "PENDING"
	GrantStatusREVOKED GrantStatus = "REVOKED"
)

type Grant struct {
	Provider string     `json:"provider" dynamodbav:"provider"`
	Subject  string     `json:"subject" dynamodbav:"subject"`
	With     Grant_With `json:"with" dynamodbav:"with"`
	//the time which the grant starts
	Start time.Time `json:"start" dynamodbav:"start"`
	//the time the grant is scheduled to end
	End       time.Time   `json:"end" dynamodbav:"end"`
	Status    GrantStatus `json:"status" dynamodbav:"status"`
	CreatedAt time.Time   `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt" dynamodbav:"updatedAt"`
}

type Option struct {
	Value       string  `json:"value" dynamodbav:"value"`
	Label       string  `json:"label" dynamodbav:"label"`
	Description *string `json:"description" dynamodbav:"description"`
}

// Describes whether a request has been approved automatically or from a review
type ApprovalMethod string

// Defines values for ApprovalMethod.
const (
	AUTOMATIC ApprovalMethod = "AUTOMATIC"
	REVIEWED  ApprovalMethod = "REVIEWED"
)

type Request struct {
	// ID is a read-only field after the request has been created.
	ID string `json:"id" dynamodbav:"id"`

	// RequestedBy is the ID of the user who has made the request.
	RequestedBy string `json:"requestedBy" dynamodbav:"requestedBy"`

	// Rule is the ID of the Access Rule which the request relates to.
	Rule string `json:"rule" dynamodbav:"rule"`
	// RuleVersion is the version string of the rule that this request relates to
	RuleVersion string `json:"ruleVersion" dynamodbav:"ruleVersion"`
	// SelectedWith stores a denormalised version of the option with a label at the time the request was created
	// Allowing it to be easily displayed in the frontend for context and reducing latency on loading requests
	SelectedWith    map[string]Option `json:"selectedWith"  dynamodbav:"selectedWith"`
	Status          Status            `json:"status" dynamodbav:"status"`
	Data            RequestData       `json:"data" dynamodbav:"data"`
	RequestedTiming Timing            `json:"requestedTiming" dynamodbav:"requestedTiming"`
	// When a request is approver, the approver has the option to override the timing, if they do so, this will be populated.
	// If the timing was not overriden, then the original request timing should be used.
	// Override timing should only be set by an approving review
	OverrideTiming *Timing `json:"overrideTiming,omitempty" dynamodbav:"overrideTiming,omitempty"`
	// Grant is the ID of the grant when it is created by the access handler
	Grant *Grant `json:"grant,omitempty" dynamodbav:"grant,omitempty"`
	// ApprovalMethod explains whether an approval was AUTOMATIC, or REVIEWED
	ApprovalMethod *ApprovalMethod `json:"approvalMethod,omitempty" dynamodbav:"approvalMethod,omitempty"`
	// CreatedAt is a read-only field after the request has been created.
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}

type GetIntervalOpts struct {
	Now time.Time
}

// Timing represents all the timing options available
// Duration should always be set
// StartTime should be set if this is a scheduled access
// The combination of startTime and duration make up the start and end times of a grant
type Timing struct {
	Duration time.Duration `json:"duration" dynamodbav:"duration"`
	// If the start time is not nil, this request is for scheduled access, if it is nil, then the request is for asap access
	StartTime *time.Time `json:"start,omitempty" dynamodbav:"start,omitempty"`
}

// RequestData is information provided by the user when they make the request,
// through filling in form fields in the web application.
type RequestData struct {
	Reason *string `json:"reason,omitempty" dynamodbav:"reason,omitempty"`
}
