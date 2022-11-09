package internalidentitysvc

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateGroup(t *testing.T) {
	type testcase struct {
		name          string
		in            types.CreateGroupRequest
		withListUsers storage.ListUsers
		want          *identity.Group
		wantError     error
	}

	clk := clock.NewMock()

	testcases := []testcase{
		{
			name: "ok",
			in: types.CreateGroupRequest{
				Name:        "test",
				Description: aws.String("test"),
				Members:     []string{"a", "b"},
			},
			withListUsers: storage.ListUsers{Result: []identity.User{
				{ID: "a"}, {ID: "b"},
			}},
			want: &identity.Group{
				Name:        "test",
				Description: "test",
				Users:       []string{"a", "b"},
				CreatedAt:   clk.Now(),
				UpdatedAt:   clk.Now(),
				Source:      identity.INTERNAL,
				Status:      types.IdpStatusACTIVE,
			},
		},
		{
			name: "user not found",
			in: types.CreateGroupRequest{
				Name:        "test",
				Description: aws.String("test"),
				Members:     []string{"a", "b"},
			},
			withListUsers: storage.ListUsers{Result: []identity.User{
				{ID: "a"},
			}},
			wantError: UserNotFoundError{UserID: "b"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQuery(&tc.withListUsers)
			s := Service{
				Clock: clk,
				DB:    db,
			}
			got, err := s.CreateGroup(context.Background(), tc.in)
			if tc.want != nil {
				tc.want.ID = got.ID
				tc.want.IdpID = got.IdpID
				assert.Equal(t, tc.want, got)
				assert.Equal(t, got.ID, got.IdpID)
			}
			if tc.wantError != nil {
				assert.EqualError(t, err, tc.wantError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestUpdateGroup(t *testing.T) {
	type testcase struct {
		name          string
		in            types.CreateGroupRequest
		group         identity.Group
		withListUsers storage.ListUsers
		want          *identity.Group
		wantError     error
	}

	clk := clock.NewMock()

	testcases := []testcase{
		{
			name: "ok add user",
			in: types.CreateGroupRequest{
				Name:        "test",
				Description: aws.String("test"),
				Members:     []string{"a", "b", "c"},
			},
			group: identity.Group{
				ID:          "abcd",
				IdpID:       "abcd",
				Name:        "test",
				Description: "test",
				Users:       []string{"a", "b"},
				CreatedAt:   clk.Now(),
				UpdatedAt:   clk.Now(),
				Source:      identity.INTERNAL,
				Status:      types.IdpStatusACTIVE,
			},
			withListUsers: storage.ListUsers{Result: []identity.User{
				{ID: "a"}, {ID: "b"}, {ID: "c"},
			}},
			want: &identity.Group{
				ID:          "abcd",
				IdpID:       "abcd",
				Name:        "test",
				Description: "test",
				Users:       []string{"a", "b", "c"},
				CreatedAt:   clk.Now(),
				UpdatedAt:   clk.Now(),
				Source:      identity.INTERNAL,
				Status:      types.IdpStatusACTIVE,
			},
		},
		{
			name: "not internal group",
			group: identity.Group{
				Source: "okta",
			},
			wantError: ErrNotInternal,
		},
		{
			name: "user doesn't exist",
			in: types.CreateGroupRequest{
				Name:        "test",
				Description: aws.String("test"),
				Members:     []string{"a", "b", "c"},
			},
			group: identity.Group{
				Users:  []string{"a", "b"},
				Source: identity.INTERNAL,
			},
			withListUsers: storage.ListUsers{Result: []identity.User{
				{ID: "a"}, {ID: "c"},
			}},
			wantError: UserNotFoundError{UserID: "b"},
		},
		{
			name: "remove user",
			in: types.CreateGroupRequest{
				Name:        "test",
				Description: aws.String("test"),
				Members:     []string{"a", "b"},
			},
			group: identity.Group{
				ID:          "abcd",
				IdpID:       "abcd",
				Name:        "test",
				Description: "test",
				Users:       []string{"a", "b", "c"},
				CreatedAt:   clk.Now(),
				UpdatedAt:   clk.Now(),
				Source:      identity.INTERNAL,
				Status:      types.IdpStatusACTIVE,
			},
			withListUsers: storage.ListUsers{Result: []identity.User{
				{ID: "a"}, {ID: "b"}, {ID: "c"},
			}},
			want: &identity.Group{
				ID:          "abcd",
				IdpID:       "abcd",
				Name:        "test",
				Description: "test",
				Users:       []string{"a", "b"},
				CreatedAt:   clk.Now(),
				UpdatedAt:   clk.Now(),
				Source:      identity.INTERNAL,
				Status:      types.IdpStatusACTIVE,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			db := ddbmock.New(t)
			db.MockQuery(&tc.withListUsers)
			s := Service{
				Clock: clk,
				DB:    db,
			}

			got, err := s.UpdateGroup(context.Background(), tc.group, tc.in)
			if tc.want != nil {
				tc.want.ID = got.ID
				tc.want.IdpID = got.IdpID
				assert.Equal(t, tc.want, got)
				assert.Equal(t, got.IdpID, got.ID)
			}
			if tc.wantError != nil {
				assert.EqualError(t, err, tc.wantError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
