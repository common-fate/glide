package access

import (
	"testing"

	"github.com/common-fate/ddb"
	"github.com/stretchr/testify/assert"
)

func TestReviewerDDBKeys(t *testing.T) {
	r := Reviewer{
		ReviewerID: "1",
		Request: Request{
			ID:     "req_1",
			Status: APPROVED,
		},
	}

	want := ddb.Keys{
		PK:     "REQUEST_REVIEWER#",
		SK:     "req_1#1",
		GSI1PK: "REQUEST_REVIEWER#1",
		GSI1SK: "req_1",
		GSI2PK: "REQUEST_REVIEWER#1",
		GSI2SK: "APPROVED#req_1",
	}
	got, err := r.DDBKeys()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}
