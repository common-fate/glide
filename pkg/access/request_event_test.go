package access

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestEventMarshalDDB(t *testing.T) {
	type testcase struct {
		name string
		give RequestEvent
		want string
	}

	testcases := []testcase{
		{
			name: "basic",
			give: RequestEvent{
				ID:        "his_28w2Eebt2Q8nFQJ2dKa1FTE9X0J",
				RequestID: "req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J",
			},
			want: `{"PK":"ACCESS_REQUEST_EVENT#","SK":"req_28w2Eebt2Q8nFQJ2dKa1FTE9X0J#his_28w2Eebt2Q8nFQJ2dKa1FTE9X0J"}`,
		},
	}

	for i := range testcases {
		tc := testcases[i]
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
