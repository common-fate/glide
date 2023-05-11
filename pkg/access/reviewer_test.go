package access

import (
	"testing"

	"github.com/common-fate/ddb"
	"github.com/stretchr/testify/assert"
)

func TestReviewerDDBKeys(t *testing.T) {
	r := Reviewer{
		RequestID:  "req_1",
		ReviewerID: "1",
	}

	want := ddb.Keys{
		PK: "REQUEST_REVIEWERV2#",
		SK: "req_1#1",
	}
	got, err := r.DDBKeys()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}
