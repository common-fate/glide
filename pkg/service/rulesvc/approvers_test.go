package rulesvc

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/stretchr/testify/assert"
)

func TestGetApprovers(t *testing.T) {
	type testcase struct {
		name         string
		giveRule     rule.AccessRule
		mockGetGroup *identity.Group
		want         []string
	}

	testcases := []testcase{
		{
			name: "users only",
			giveRule: rule.AccessRule{
				Approval: rule.Approval{
					Users: []string{"usr_1"},
				},
			},
			want: []string{"usr_1"},
		},
		{
			name: "users and groups",
			giveRule: rule.AccessRule{
				Approval: rule.Approval{
					Users:  []string{"usr_1"},
					Groups: []string{"grp_1"},
				},
			},
			mockGetGroup: &identity.Group{
				Users: []string{"usr_2"},
			},
			want: []string{"usr_1", "usr_2"},
		},
		{
			name: "users and groups duplicate",
			giveRule: rule.AccessRule{
				Approval: rule.Approval{
					Users:  []string{"usr_1", "usr_2"},
					Groups: []string{"grp_1"},
				},
			},
			mockGetGroup: &identity.Group{
				Users: []string{"usr_2"},
			},
			want: []string{"usr_1", "usr_2"},
		},
		{
			name: "group only",
			giveRule: rule.AccessRule{
				Approval: rule.Approval{
					Groups: []string{"grp_1"},
				},
			},
			mockGetGroup: &identity.Group{
				Users: []string{"usr_2"},
			},
			want: []string{"usr_2"},
		},
		// returning an empty array rather than nil ensures that our API endpoints
		// that use this method don't return null when the frontend is expecting an array.
		{
			name:     "no approvers",
			giveRule: rule.AccessRule{},
			want:     []string{},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQuery(&storage.GetGroup{Result: tc.mockGetGroup})

			ctx := context.Background()
			got, err := GetApprovers(ctx, db, tc.giveRule)
			if err != nil {
				t.Fatal(err)
			}
			assert.ElementsMatch(t, tc.want, got)
		})
	}

}
