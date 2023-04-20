package internalidentitysvc

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUserGroups(t *testing.T) {
	type testcase struct {
		name           string
		groups         []string
		user           identity.User
		withListGroups storage.ListGroupsForSourceAndStatus
		want           *identity.User
		wantError      error
	}

	clk := clock.NewMock()

	testcases := []testcase{
		{
			name:   "ok",
			groups: []string{"a", "b"},
			user: identity.User{
				ID:        "abcd",
				FirstName: "name",
				LastName:  "name",
				Email:     "email",
				Groups:    []string{"c"},
				Status:    types.ACTIVE,
				CreatedAt: clk.Now(),
				UpdatedAt: clk.Now().Add(-time.Second),
			},
			withListGroups: storage.ListGroupsForSourceAndStatus{Result: []identity.Group{
				{ID: "a", Source: identity.INTERNAL}, {ID: "b", Source: identity.INTERNAL},
			}},
			want: &identity.User{
				ID:        "abcd",
				FirstName: "name",
				LastName:  "name",
				Email:     "email",
				Groups:    []string{"a", "b", "c"},
				Status:    types.ACTIVE,
				CreatedAt: clk.Now(),
				UpdatedAt: clk.Now(),
			},
		},
		{
			name:   "ok",
			groups: []string{"a", "b"},
			user: identity.User{
				Groups: []string{"c"},
			},
			withListGroups: storage.ListGroupsForSourceAndStatus{Result: []identity.Group{
				{ID: "b", Source: identity.INTERNAL},
			}},
			wantError: ErrGroupNotFoundOrNotInternal,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQuery(&tc.withListGroups)
			s := Service{
				Clock: clk,
				DB:    db,
			}
			got, err := s.UpdateUserGroups(context.Background(), tc.user, tc.groups)
			if tc.want != nil {
				assert.Equal(t, tc.want, got)
			}
			if tc.wantError != nil {
				assert.EqualError(t, err, tc.wantError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
