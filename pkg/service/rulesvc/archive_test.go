package rulesvc

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/stretchr/testify/assert"
)

func TestArchiveAccessRule(t *testing.T) {
	type testcase struct {
		name      string
		givenRule rule.AccessRule
		wantErr   error
		want      *rule.AccessRule
	}

	clk := clock.NewMock()
	now := clk.Now()
	mockRule := rule.AccessRule{
		ID: "rule",

		Status: rule.ACTIVE,
		Metadata: rule.AccessRuleMetadata{
			CreatedAt: now.Add(time.Minute),
			CreatedBy: "hello",
			UpdatedAt: now.Add(time.Minute),
			UpdatedBy: "hello",
		},
		Current: true,
	}
	want := mockRule
	want.Status = rule.ARCHIVED
	want.Metadata.UpdatedAt = now
	want.Current = true

	testcases := []testcase{
		{
			name:      "ok",
			givenRule: mockRule,
			want:      &want,
		},
		{
			name: "already archived",
			givenRule: rule.AccessRule{
				Status: rule.ACTIVE,
			},
			wantErr: ErrAccessRuleAlreadyArchived,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)
			db.PutBatchErr = tc.wantErr
			db.MockQueryWithErrWithResult(&storage.ListRequestsForStatus{Status: access.PENDING, Result: []access.Request{}}, &ddb.QueryResult{}, nil)
			db.MockQueryWithErrWithResult(&storage.ListRequestReviewers{Result: []access.Reviewer{}}, &ddb.QueryResult{}, nil)

			s := Service{
				Clock: clk,
				DB:    db,
			}

			got, err := s.ArchiveAccessRule(context.Background(), "", tc.givenRule)

			// This is the only thing from service layer that we can't mock yet, hence the override
			if err == nil {
				// Rule id and version id must not be empty strings, we check this prior to overwriting them
				assert.NotEmpty(t, got.Version)
				got.Version = tc.want.Version

			}

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
