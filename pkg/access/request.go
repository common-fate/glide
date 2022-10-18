package access

import (
	"time"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providerregistry"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	ac_types "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/cache"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/common-fate/iso8601"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

// Status of an Access Request.
type Status string

const (
	APPROVED  Status = "APPROVED"
	DECLINED  Status = "DECLINED"
	CANCELLED Status = "CANCELLED"
	PENDING   Status = "PENDING"
)

type Grant struct {
	Provider string              `json:"provider" dynamodbav:"provider"`
	Subject  string              `json:"subject" dynamodbav:"subject"`
	With     ac_types.Grant_With `json:"with" dynamodbav:"with"`
	//the time which the grant starts
	Start time.Time `json:"start" dynamodbav:"start"`
	//the time the grant is scheduled to end
	End       time.Time            `json:"end" dynamodbav:"end"`
	Status    ac_types.GrantStatus `json:"status" dynamodbav:"status"`
	CreatedAt time.Time            `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt" dynamodbav:"updatedAt"`
}

func (g *Grant) ToAHGrant(requestID string) ac_types.Grant {
	return ac_types.Grant{
		ID:       requestID,
		Start:    iso8601.New(g.Start),
		End:      iso8601.New(g.End),
		Provider: g.Provider,
		Subject:  openapi_types.Email(g.Subject),
		Status:   g.Status,
		With:     g.With,
	}
}
func (g *Grant) ToAPI() types.Grant {
	req := types.Grant{
		Start:    g.Start,
		End:      g.End,
		Provider: g.Provider,
		Subject:  openapi_types.Email(g.Subject),
		Status:   types.GrantStatus(g.Status),
	}

	return req
}

type Option struct {
	Value       string  `json:"value" dynamodbav:"value"`
	Label       string  `json:"label" dynamodbav:"label"`
	Description *string `json:"description" dynamodbav:"description"`
}

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
	ApprovalMethod *types.ApprovalMethod `json:"approvalMethod,omitempty" dynamodbav:"approvalMethod,omitempty"`
	// CreatedAt is a read-only field after the request has been created.
	CreatedAt time.Time `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dynamodbav:"updatedAt"`
}
type GetIntervalOpts struct {
	Now time.Time
}

// WithNow allows you to override the now time used by getInterval
func WithNow(t time.Time) func(o *GetIntervalOpts) {
	return func(o *GetIntervalOpts) { o.Now = t }
}

// GetInterval will return the interval for either the requested timing or for the override timing if it is present
func (r *Request) GetInterval(opts ...func(o *GetIntervalOpts)) (start time.Time, end time.Time) {
	if r.OverrideTiming != nil {
		return r.OverrideTiming.GetInterval(opts...)
	}
	return r.RequestedTiming.GetInterval(opts...)
}

// IsScheduled will return true if this request is scheduled, first checking for override timing, then for original timing
func (r *Request) IsScheduled() bool {
	if r.OverrideTiming != nil {
		return r.OverrideTiming.IsScheduled()
	}
	return r.RequestedTiming.IsScheduled()
}

func (r *Request) ToAPI() types.Request {
	req := types.Request{
		AccessRule: types.RequestAccessRule{
			Id:      r.Rule,
			Version: r.RuleVersion,
		},
		Timing:         r.RequestedTiming.ToAPI(),
		Reason:         r.Data.Reason,
		ID:             r.ID,
		RequestedAt:    r.CreatedAt,
		Requestor:      r.RequestedBy,
		Status:         types.RequestStatus(r.Status),
		UpdatedAt:      r.UpdatedAt,
		ApprovalMethod: r.ApprovalMethod,
	}
	if r.Grant != nil {
		g := r.Grant.ToAPI()
		req.Grant = &g
	}

	// show the updated timing rather than the requested timing if it's been overridden by an approver.
	if r.OverrideTiming != nil {
		req.Timing = r.OverrideTiming.ToAPI()
	}

	return req
}

func (r *Request) ToAPIDetail(accessRule rule.AccessRule, canReview bool, argOptions []cache.ProviderOption) types.RequestDetail {
	req := types.RequestDetail{
		AccessRule:     accessRule.ToRequestAccessRuleDetailAPI(argOptions),
		Timing:         r.RequestedTiming.ToAPI(),
		Reason:         r.Data.Reason,
		ID:             r.ID,
		RequestedAt:    r.CreatedAt,
		Requestor:      r.RequestedBy,
		Status:         types.RequestStatus(r.Status),
		UpdatedAt:      r.UpdatedAt,
		CanReview:      canReview,
		ApprovalMethod: r.ApprovalMethod,
		SelectedWith: types.RequestDetail_SelectedWith{
			AdditionalProperties: make(map[string]types.With),
		},
	}
	// Lookup the provider, ignore errors
	// if provider is not found, fallback to using the argument key as the title
	_, provider, _ := providerregistry.Registry().GetLatestByShortType(accessRule.Target.ProviderType)
	for k, v := range r.SelectedWith {
		with := types.With{
			Label:       v.Label,
			Value:       v.Value,
			Title:       k,
			Description: v.Description,
		}
		// attempt to get the title for the argument from the provider arg schema
		if provider != nil {
			if s, ok := provider.Provider.(providers.ArgSchemarer); ok {
				schema := s.ArgSchema()
				if arg, ok := schema[k]; ok {
					with.Title = arg.Title
				}
			}
		}
		req.SelectedWith.AdditionalProperties[k] = with
	}
	if r.Grant != nil {
		g := r.Grant.ToAPI()
		req.Grant = &g
	}
	// show the updated timing rather than the requested timing if it's been overridden by an approver.
	if r.OverrideTiming != nil {
		req.Timing = r.OverrideTiming.ToAPI()
	}

	return req
}

func (r *Request) DDBKeys() (ddb.Keys, error) {
	// - APPROVED requests have an end time on the grant
	// - PENDING Scheduled requests have a request end time
	// - PENDING asap requests should have MAXIMUM endtime
	// - Declined and Cancelled requests should have an end time = createdAt so they get a somewhat natural order in the results
	// - REVOKED grants should have end time = created at
	// - ERROR grants should have end times = created at
	end := r.CreatedAt
	if r.Status == APPROVED || r.Status == PENDING {
		if r.Grant != nil {
			//any grant status other than revoked or error should be equal to grant.end.
			//this is to make sure the error and revoke grants are pushed to the past column in the frontend
			if !(r.Grant.Status == ac_types.GrantStatusREVOKED || r.Grant.Status == ac_types.GrantStatusERROR) {
				end = r.Grant.End
			}
		} else if r.IsScheduled() {
			_, end = r.GetInterval()
		} else {
			// maximum time value in Go
			// this means that asap requests which are not approved will always be the first in results because the end time in unknown until approval
			end = time.Unix(1<<63-1, 0)
		}
	}

	keys := ddb.Keys{
		PK:     keys.AccessRequest.PK1,
		SK:     r.ID,
		GSI1PK: keys.AccessRequest.GSI1PK(r.RequestedBy),
		GSI1SK: keys.AccessRequest.GSI1SK(r.ID),
		GSI2PK: keys.AccessRequest.GSI2PK(string(r.Status)),
		GSI2SK: keys.AccessRequest.GSI2SK(r.RequestedBy, r.ID),
		GSI3PK: keys.AccessRequest.GSI3PK(r.RequestedBy),
		GSI3SK: keys.AccessRequest.GSI3SK(end),
		GSI4PK: keys.AccessRequest.GSI4PK(r.RequestedBy, r.Rule),
		GSI4SK: keys.AccessRequest.GSI4SK(end),
	}

	return keys, nil
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

// TimingFromRequestTiming converts from the api type to the internal type
func TimingFromRequestTiming(r types.RequestTiming) Timing {
	return Timing{
		Duration:  time.Second * time.Duration(r.DurationSeconds),
		StartTime: r.StartTime,
	}
}

// IsScheduled is true if the startTime is not nil
func (t *Timing) IsScheduled() bool {
	return t.StartTime != nil
}

// ToAPI returns the api representation of the timing information
func (t *Timing) ToAPI() types.RequestTiming {
	return types.RequestTiming{
		DurationSeconds: int(t.Duration.Seconds()),
		StartTime:       t.StartTime,
	}
}

// GetInterval returns a start and end time for this timing information
// it will either return times for scheduled access if the timing represents scheduled access.
// Or it will use the time.Now() as the start time.
//
// To override the start time for asap timing, pass in the WithNow(t time.Time) function
func (t *Timing) GetInterval(opts ...func(o *GetIntervalOpts)) (start time.Time, end time.Time) {
	if t.IsScheduled() {
		return *t.StartTime, t.StartTime.Add(t.Duration)
	}
	cfg := GetIntervalOpts{
		Now: time.Now(),
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg.Now, cfg.Now.Add(t.Duration)
}

// RequestData is information provided by the user when they make the request,
// through filling in form fields in the web application.
type RequestData struct {
	Reason *string `json:"reason,omitempty" dynamodbav:"reason,omitempty"`
}
