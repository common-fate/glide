package access

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRequestEventMarshalDDB(t *testing.T) {
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
