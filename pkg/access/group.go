package access

import (
	"time"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type AccessRule struct {
	ID string `json:"id" dynamodbav:"id"`
}

type Group struct {
	ID         string                         `json:"id" dynamodbav:"id"`
	RequestID  string                         `json:"requestId" dynamodbav:"request"`
	AccessRule AccessRule                     `json:"accessRule" dynamodbav:"accessRule"`
	Status     types.RequestAccessGroupStatus `json:"status" dynamodbav:"status"`
	// Also denormalised across all the request items
	RequestStatus   types.RequestStatus `json:"requestStatus" dynamodbav:"requestStatus"`
	TimeConstraints Timing              `json:"timeConstraints" dynamodbav:"timeConstraints"`
	OverrideTiming  *Timing             `json:"overrideTimings,omitempty" dynamodbav:"overrideTimings,omitempty"`
	RequestedBy     string              `json:"requestedBy" dynamodbav:"requestedBy"`
	CreatedAt       time.Time           `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt" dynamodbav:"updatedAt"`
	// request reviewers are users who have one or more groups to review on the request as a whole
	RequestReviewers []string `json:"requestReviewers" dynamodbav:"requestReviewers, set"`
	// groupReviewers are the users who are able to review this access group
	GroupReviewers []string `json:"groupReviewers" dynamodbav:"groupReviewers, set"`
}

type GroupWithTargets struct {
	Group
	Targets []GroupTarget
}

func (g *GroupWithTargets) ToAPI() types.RequestAccessGroup {
	out := types.RequestAccessGroup{
		Id:        g.ID,
		RequestId: g.RequestID,
		Status:    g.Status,
		Time:      g.TimeConstraints.ToAPI(),
		Targets:   []types.RequestAccessGroupTarget{},
		CreatedAt: g.CreatedAt,
		UpdatedAt: g.UpdatedAt,
	}
	if g.OverrideTiming != nil {
		out.OverrideTiming = g.OverrideTiming.ToAPI()
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
	return r.TimeConstraints.GetInterval(opts...)
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
		Duration:  time.Second * time.Duration(r.DurationSeconds),
		StartTime: r.StartTime,
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
		GSI1PK: keys.AccessRequestGroup.GSI1PK(i.RequestedBy),
		GSI1SK: keys.AccessRequestGroup.GSI1SK(i.RequestID, i.ID),
		GSI2PK: keys.AccessRequestGroup.GSI2PK(i.RequestedBy, RequestStatusToPastOrUpcoming(i.RequestStatus)),
		GSI2SK: keys.AccessRequestGroup.GSI2SK(i.RequestID, i.ID),
	}
	return keys, nil
}
