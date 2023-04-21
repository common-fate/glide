package storage

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/segmentio/ksuid"

	"github.com/stretchr/testify/assert"
)

func TestListGroupsForStatus(t *testing.T) {
	type testcase struct {
		name         string
		status       types.IdpStatus
		insertBefore []identity.Group
		want         []identity.Group
		notWant      []identity.Group
		wantErr      error
	}

	gACTIVE := identity.Group{
		ID:     ksuid.New().String(),
		Name:   "a",
		IdpID:  "a",
		Users:  []string{"a"},
		Status: types.IdpStatusACTIVE,
	}

	gARCHIVED := identity.Group{
		ID:     ksuid.New().String(),
		Name:   "a",
		IdpID:  "a",
		Users:  []string{"a"},
		Status: types.IdpStatusARCHIVED,
	}

	testcases := []testcase{
		{

			name:         "get active",
			insertBefore: []identity.Group{gACTIVE, gARCHIVED},
			want:         []identity.Group{gACTIVE},
			status:       types.IdpStatusACTIVE,
			notWant:      []identity.Group{},
		},
		{
			name:         "get archived",
			insertBefore: []identity.Group{gACTIVE, gARCHIVED},
			want:         []identity.Group{gARCHIVED},
			status:       types.IdpStatusARCHIVED,
			notWant:      []identity.Group{},
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

			q := ListGroupsForStatus{Status: tc.status}
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
