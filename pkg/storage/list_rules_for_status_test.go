package storage

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/stretchr/testify/assert"
)

func TestListAccessRules(t *testing.T) {

	type testcase struct {
		name         string
		status       rule.Status
		insertBefore []rule.AccessRule
		want         []rule.AccessRule
		notWant      []rule.AccessRule
		wantErr      error
	}

	activeRule := rule.TestAccessRule()
	archivedRule := rule.TestAccessRule()
	archivedRule.Status = rule.ARCHIVED
	archived := rule.ARCHIVED
	active := rule.ACTIVE

	testcases := []testcase{
		// {
		// 	name:         "ok",
		// 	insertBefore: []rule.AccessRule{activeRule, archivedRule},
		// 	want:         []rule.AccessRule{activeRule},
		// },
		{
			name:         "archived",
			insertBefore: []rule.AccessRule{activeRule, archivedRule},
			want:         []rule.AccessRule{archivedRule},
			status:       archived,
			notWant:      []rule.AccessRule{activeRule},
		},
		{
			name:         "active",
			insertBefore: []rule.AccessRule{activeRule, archivedRule},
			status:       active,
			want:         []rule.AccessRule{activeRule},
			notWant:      []rule.AccessRule{archivedRule},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			s := newTestingStorage(t)

			// insert any required fixture data
			for _, r := range tc.insertBefore {
				err := s.Put(ctx, &r)
				if err != nil {
					t.Fatal(err)
				}
			}

			q := ListAccessRulesForStatus{Status: tc.status}
			_, err := s.Query(ctx, &q)
			if err != nil && tc.wantErr == nil {
				t.Fatal(err)
			}
			got := q.Result

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			}
			for _, item := range tc.want {
				assert.Contains(t, got, item)
			}
			for _, item := range tc.notWant {
				assert.NotContains(t, got, item, "expected item to not be in results")
			}
		})
	}
}
