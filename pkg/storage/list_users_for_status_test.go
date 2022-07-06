package storage

import (
	"context"
	"testing"

	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestListUsersStatus(t *testing.T) {
	type testcase struct {
		name         string
		status       types.IdpStatus
		insertBefore []identity.User
		want         []identity.User
		notWant      []identity.User
		wantErr      error
	}

	gACTIVE := identity.User{
		ID: types.NewUserID(),

		Status: types.IdpStatusACTIVE,
	}

	gARCHIVED := identity.User{
		ID: ksuid.New().String(),

		Status: types.IdpStatusARCHIVED,
	}

	testcases := []testcase{
		{
			name:         "get active",
			insertBefore: []identity.User{gACTIVE, gARCHIVED},
			want:         []identity.User{gACTIVE},
			status:       types.IdpStatusACTIVE,
			notWant:      []identity.User{},
		},
		{
			name:         "get archived",
			insertBefore: []identity.User{gACTIVE, gARCHIVED},
			want:         []identity.User{gARCHIVED},
			status:       types.IdpStatusARCHIVED,
			notWant:      []identity.User{},
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
			}

			q := ListUsersForStatus{Status: tc.status}
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
