package storage

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/stretchr/testify/assert"
)

func TestListAccessRuleVersions(t *testing.T) {
	type testcase struct {
		name          string
		insertBefore  []rule.AccessRule
		listForRuleID string
		want          []rule.AccessRule
	}

	rul := rule.TestAccessRule()
	testcases := []testcase{
		{
			name:          "ok",
			insertBefore:  []rule.AccessRule{rul},
			listForRuleID: rul.ID,
			want:          []rule.AccessRule{rul},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			s := newTestingStorage(t)

			// insert any required fixture data
			for _, r := range tc.insertBefore {
				err := s.Put(ctx, &r)
				if err != nil {
					t.Fatal(err)
				}
				err = s.Put(ctx, &r)
				if err != nil {
					t.Fatal(err)
				}
			}

			q := ListAccessRuleVersions{ID: tc.listForRuleID}
			_, err := s.Query(ctx, &q)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tc.want, q.Result)
		})
	}
}
