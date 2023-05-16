package access

import (
	"time"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type AccessRule struct {
	ID string `json:"id" dynamodbav:"id"`
}

type Group struct {
	ID        string `json:"id" dynamodbav:"id"`
	RequestID string `json:"requestId" dynamodbav:"request"`
	// This is a snapshot of the access rule as it was configured when the request was submitted
	AccessRuleSnapshot rule.AccessRule                         `json:"accessRuleSnapshot" dynamodbav:"accessRuleSnapshot"`
	Status             types.RequestAccessGroupStatus          `json:"status" dynamodbav:"status"`
	ApprovalMethod     *types.RequestAccessGroupApprovalMethod `json:"approvalMethod" dynamodbav:"approvalMethod"`
	// Also denormalised across all the request items
	RequestPurposeReason string              `json:"requestPurposeReason" dynamodbav:"requestPurposeReason"`
	RequestStatus        types.RequestStatus `json:"requestStatus" dynamodbav:"requestStatus"`
	RequestedTiming      Timing              `json:"requestedTiming" dynamodbav:"requestedTiming"`
	FinalTiming          *FinalTiming        `json:"finalTiming" dynamodbav:"finalTiming"`
	OverrideTiming       *Timing             `json:"overrideTimings,omitempty" dynamodbav:"overrideTimings,omitempty"`
	RequestedBy          RequestedBy         `json:"requestedBy" dynamodbav:"requestedBy"`
	CreatedAt            time.Time           `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt            time.Time           `json:"updatedAt" dynamodbav:"updatedAt"`
	// request reviewers are users who have one or more groups to review on the request as a whole
	RequestReviewers []string `json:"requestReviewers" dynamodbav:"requestReviewers, set"`
	// groupReviewers are the users who are able to review this access group; id = access.Reviewer.ID
	GroupReviewers []string `json:"groupReviewers" dynamodbav:"groupReviewers, set"`
}

type FinalTiming struct {
	Start time.Time `json:"start" dynamodbav:"start"`
	End   time.Time `json:"end" dynamodbav:"end"`
}

type GroupWithTargets struct {
	Group   Group         `json:"group"`
	Targets []GroupTarget `json:"targets"`
}

func (r *GroupWithTargets) DBItems() []ddb.Keyer {
	var items []ddb.Keyer
	items = append(items, &r.Group)
	for i := range r.Targets {
		items = append(items, &r.Targets[i])
	}
	return items
}
func (g *GroupWithTargets) ToAPI() types.RequestAccessGroup {
	out := types.RequestAccessGroup{
		Id:              g.Group.ID,
		RequestId:       g.Group.RequestID,
		Status:          g.Group.Status,
		RequestedTiming: g.Group.RequestedTiming.ToAPI(),
		Targets:         []types.RequestAccessGroupTarget{},
		ApprovalMethod:  g.Group.ApprovalMethod,
		CreatedAt:       g.Group.CreatedAt,
		UpdatedAt:       g.Group.UpdatedAt,
		RequestedBy:     types.RequestRequestedBy(g.Group.RequestedBy),
		RequestStatus:   g.Group.RequestStatus,
	}
	if g.Group.FinalTiming != nil {
		out.FinalTiming = &types.RequestAccessGroupFinalTiming{
			StartTime: g.Group.FinalTiming.Start,
			EndTime:   g.Group.FinalTiming.End,
		}
	}
	if g.Group.GroupReviewers != nil {
		out.GroupReviewers = &g.Group.GroupReviewers
	}
	if g.Group.RequestReviewers != nil {
		out.GroupReviewers = &g.Group.RequestReviewers
	}
	if g.Group.OverrideTiming != nil {
		ot := g.Group.OverrideTiming.ToAPI()
		out.OverrideTiming = &ot
	}
	for _, target := range g.Targets {
		out.Targets = append(out.Targets, target.ToAPI())
	}

	return out

}

func (r *Group) GetInterval(opts ...func(o *GetIntervalOpts)) (start time.Time, end time.Time) {
	if r.OverrideTiming != nil {
		return r.OverrideTiming.GetInterval(opts...)
	}
	return r.RequestedTiming.GetInterval(opts...)
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

func (t Timing) ToAnalytics() analytics.Timing {
	mode := analytics.TimingModeASAP
	if t.IsScheduled() {
		mode = analytics.TimingModeScheduled
	}

	return analytics.Timing{
		Mode:            mode,
		DurationSeconds: t.Duration.Seconds(),
	}
}

// TimingFromRequestTiming converts from the api type to the internal type
func TimingFromRequestTiming(r types.RequestAccessGroupTiming) Timing {

	return Timing{
		Duration: time.Second * time.Duration(r.DurationSeconds),
	}
}

// IsScheduled is true if the startTime is not nil
func (t *Timing) IsScheduled() bool {
	return t.StartTime != nil
}

// ToAPI returns the api representation of the timing information
func (t *Timing) ToAPI() types.RequestAccessGroupTiming {
	return types.RequestAccessGroupTiming{
		DurationSeconds: int(t.Duration.Seconds()),
		StartTime:       t.StartTime,
	}
}

// WithNow allows you to override the now time used by getInterval
func WithNow(t time.Time) func(o *GetIntervalOpts) {
	return func(o *GetIntervalOpts) { o.Now = t }
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

type GetIntervalOpts struct {
	Now time.Time
}

func (i *Group) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK:     keys.AccessRequestGroup.PK1,
		SK:     keys.AccessRequestGroup.SK1(i.RequestID, i.ID),
		GSI1PK: keys.AccessRequestGroup.GSI1PK(i.RequestedBy.ID),
		GSI1SK: keys.AccessRequestGroup.GSI1SK(RequestStatusToPastOrUpcoming(i.RequestStatus), i.RequestID, i.ID),
		GSI2PK: keys.AccessRequestGroup.GSI2PK(i.RequestStatus),
		GSI2SK: keys.AccessRequestGroup.GSI2SK(i.RequestID, i.ID),
	}
	return keys, nil
}
