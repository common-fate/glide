package requests

import (
	"time"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Status string

const (
	APPROVED  Status = "APPROVED"
	DECLINED  Status = "DECLINED"
	CANCELLED Status = "CANCELLED"
	PENDING   Status = "PENDING"
)

type AccessGroup struct {
	AccessRule      rule.AccessRule `json:"accessRule" dynamodbav:"accessRule"`
	ID              string          `json:"id" dynamodbav:"id"`
	Request         string          `json:"request" dynamodbav:"request"`
	TimeConstraints Timing          `json:"timeConstraints" dynamodbav:"timeConstraints"`
	OverrideTiming  *Timing         `json:"overrideTimings,omitempty" dynamodbav:"overrideTimings,omitempty"`
	CreatedAt       time.Time       `json:"createdAt" dynamodbav:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt" dynamodbav:"updatedAt"`
	Status          Status          `json:"status" dynamodbav:"status"`
	// ApprovalMethod explains whether an approval was AUTOMATIC, or REVIEWED
	ApprovalMethod *types.ApprovalMethod `json:"approvalMethod,omitempty" dynamodbav:"approvalMethod,omitempty"`
}

func (i *AccessGroup) DDBKeys() (ddb.Keys, error) {
	keys := ddb.Keys{
		PK: keys.AccessGroup.PK1,
		SK: keys.AccessGroup.SK1(i.Request, i.ID),
	}
	return keys, nil
}

func (i *AccessGroup) ToAPI() types.RequestAccessGroup {
	out := types.RequestAccessGroup{
		Id:     i.ID,
		Status: types.RequestStatus(i.Status),
		Time: types.AccessRuleTimeConstraints{
			MaxDurationSeconds: int(i.TimeConstraints.Duration),
		},
		RequestId:      i.Request,
		OverrideTiming: i.OverrideTiming.ToAPI(),
		CreatedAt:      i.CreatedAt,
		UpdatedAt:      i.UpdatedAt,
	}

	if i.OverrideTiming != nil {
		out.OverrideTiming = i.OverrideTiming.ToAPI()
	}

	return out

}

func (r *AccessGroup) GetInterval(opts ...func(o *GetIntervalOpts)) (start time.Time, end time.Time) {
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
