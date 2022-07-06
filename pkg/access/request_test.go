package access

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestRequestMarshalDDB(t *testing.T) {
	type testcase struct {
		name string
		give Request
		want string
	}

	reason := "test reason"

	testcases := []testcase{
		{
			name: "basic",
			give: Request{
				ID:          "req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J",
				RequestedBy: "user",
				Rule:        "rul_123",
				RuleVersion: "2022-01-01T10:00:00Z",
				Status:      "PENDING",
				Data: RequestData{
					Reason: &reason,
				},
				RequestedTiming: Timing{
					Duration: time.Minute * 5,
				},
				CreatedAt: time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			want: `{"PK":"ACCESS_REQUEST#","SK":"req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J","GSI1PK":"ACCESS_REQUEST#user","GSI1SK":"req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J","GSI2PK":"ACCESS_REQUEST#PENDING","GSI2SK":"user#req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J","GSI3PK":"ACCESS_REQUEST#user","GSI3SK":"292277026596-12-04T15:30:07Z","GSI4PK":"ACCESS_REQUEST#user#rul_123","GSI4SK":"292277026596-12-04T15:30:07Z"}`,
		},
		{
			name: "grant revoked",
			give: Request{
				ID:          "req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J",
				RequestedBy: "user",
				Rule:        "rul_123",
				RuleVersion: "2022-01-01T10:00:00Z",
				Status:      APPROVED,
				Data: RequestData{
					Reason: &reason,
				},
				RequestedTiming: Timing{
					Duration: time.Minute * 5,
				},
				CreatedAt: time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
				Grant: &Grant{
					Status: types.REVOKED,
					End:    time.Date(2022, 1, 1, 10, 1, 0, 0, time.UTC),
				},
			},
			want: `{"PK":"ACCESS_REQUEST#","SK":"req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J","GSI1PK":"ACCESS_REQUEST#user","GSI1SK":"req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J","GSI2PK":"ACCESS_REQUEST#APPROVED","GSI2SK":"user#req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J","GSI3PK":"ACCESS_REQUEST#user","GSI3SK":"2022-01-01T10:00:00Z","GSI4PK":"ACCESS_REQUEST#user#rul_123","GSI4SK":"2022-01-01T10:00:00Z"}`,
		},
		{
			name: "approved grant active",
			give: Request{
				ID:          "req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J",
				RequestedBy: "user",
				Rule:        "rul_123",
				RuleVersion: "2022-01-01T10:00:00Z",
				Status:      APPROVED,
				Data: RequestData{
					Reason: &reason,
				},
				RequestedTiming: Timing{
					Duration: time.Minute * 5,
				},
				CreatedAt: time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC),
				Grant: &Grant{
					Status: types.ACTIVE,
					End:    time.Date(2022, 1, 1, 10, 1, 0, 0, time.UTC),
				},
			},
			want: `{"PK":"ACCESS_REQUEST#","SK":"req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J","GSI1PK":"ACCESS_REQUEST#user","GSI1SK":"req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J","GSI2PK":"ACCESS_REQUEST#APPROVED","GSI2SK":"user#req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J","GSI3PK":"ACCESS_REQUEST#user","GSI3SK":"2022-01-01T10:01:00Z","GSI4PK":"ACCESS_REQUEST#user#rul_123","GSI4SK":"2022-01-01T10:01:00Z"}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			item, err := tc.give.DDBKeys()
			if err != nil {
				t.Fatal(err)
			}
			got, err := json.Marshal(item)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.want, string(got))
		})
	}
}

func TestRequestGetInterval(t *testing.T) {
	type testcase struct {
		name      string
		give      Request
		withNow   *time.Time
		wantStart time.Time
		wantEnd   time.Time
	}
	now := time.Now()
	scheduledStart := time.Now().Add(time.Minute)
	nowOverride := time.Now().Add(time.Hour)
	scheduledStartOverride := time.Now().Add(time.Hour * 2)
	duration := time.Minute * 5
	testcases := []testcase{
		{
			name: "requestedTimingASAP",
			give: Request{
				RequestedTiming: Timing{
					Duration: duration,
				},
			},
			withNow:   &now,
			wantStart: now,
			wantEnd:   now.Add(duration),
		},
		{
			name: "requestedTimingScheduled",
			give: Request{
				RequestedTiming: Timing{
					Duration:  duration,
					StartTime: &scheduledStart,
				},
			},
			wantStart: scheduledStart,
			wantEnd:   scheduledStart.Add(duration),
		},
		{
			name: "requestedTimingScheduledIgnoredWithNow",
			give: Request{
				RequestedTiming: Timing{
					Duration:  duration,
					StartTime: &scheduledStart,
				},
			},
			withNow:   &now,
			wantStart: scheduledStart,
			wantEnd:   scheduledStart.Add(duration),
		},
		{
			name: "overrideTimingASAP",
			give: Request{
				RequestedTiming: Timing{
					Duration: duration,
				},
			},
			withNow:   &nowOverride,
			wantStart: nowOverride,
			wantEnd:   nowOverride.Add(duration),
		},
		{
			name: "overrideTimingScheduled",
			give: Request{
				RequestedTiming: Timing{
					Duration:  duration,
					StartTime: &scheduledStartOverride,
				},
			},
			wantStart: scheduledStartOverride,
			wantEnd:   scheduledStartOverride.Add(duration),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			opts := []func(*GetIntervalOpts){}
			if tc.withNow != nil {
				opts = append(opts, WithNow(*tc.withNow))
			}
			gotStart, gotEnd := tc.give.GetInterval(opts...)
			assert.Equal(t, tc.wantStart, gotStart)
			assert.Equal(t, tc.wantEnd, gotEnd)
		})
	}
}
