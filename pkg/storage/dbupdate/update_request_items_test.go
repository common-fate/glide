package dbupdate

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUpdateRequestItems(t *testing.T) {
	type testcase struct {
		name          string
		give          access.Request
		giveOpts      []func(*UpdateRequestOpts)
		withReviewers []access.Reviewer
		want          []ddb.Keyer
		wantErr       error
	}
	request := access.Request{ID: "abcd", Status: access.PENDING}
	requestUpdated := access.Request{ID: "abcd", Status: access.APPROVED}
	reviewers := []access.Reviewer{{ReviewerID: "1", Request: request}, {ReviewerID: "2", Request: request}}
	reviewersUpdated := []access.Reviewer{{ReviewerID: "1", Request: requestUpdated}, {ReviewerID: "2", Request: requestUpdated}}
	testcases := []testcase{
		{
			name:          "ok",
			give:          requestUpdated,
			withReviewers: reviewers,
			want:          []ddb.Keyer{&requestUpdated, &reviewersUpdated[0], &reviewersUpdated[1]},
		},
		{
			name:     "supply reviewers",
			give:     requestUpdated,
			giveOpts: []func(*UpdateRequestOpts){WithReviewers(reviewers)},
			want:     []ddb.Keyer{&requestUpdated, &reviewersUpdated[0], &reviewersUpdated[1]},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			db := ddbmock.New(t)
			db.MockQuery(&storage.ListRequestReviewers{Result: tc.withReviewers})
			got, gotErr := GetUpdateRequestItems(ctx, db, tc.give, tc.giveOpts...)
			if tc.wantErr == nil {
				assert.NoError(t, gotErr)
			}
			if gotErr != nil {
				assert.EqualError(t, gotErr, tc.wantErr.Error())
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
