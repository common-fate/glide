package storage

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestListUsers(t *testing.T) {
	type testcase struct {
		name string

		insertBefore []identity.User
		want         []identity.User
		notWant      []identity.User
		wantErr      error
	}

	g := identity.User{
		ID: ksuid.New().String(),
	}

	testcases := []testcase{
		{
			name:         "ok",
			insertBefore: []identity.User{g},
			want:         []identity.User{g},

			notWant: []identity.User{},
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

			q := ListUsers{}
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
