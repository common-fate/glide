package identitysync

import (
	"testing"
	"time"

	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/stretchr/testify/assert"
)

// The processor contains all the mapping logic for create/update/map for users and groups
func TestIdentitySyncProcessor(t *testing.T) {

	type testcase struct {
		name               string
		giveIdpUsers       []identity.IDPUser
		giveIdpGroups      []identity.IDPGroup
		giveInternalUsers  []identity.User
		giveInternalGroups []identity.Group
		wantUserMap        map[string]identity.User
		wantGroupMap       map[string]identity.Group
	}
	now := time.Now()
	testcases := []testcase{
		{
			name: "create group and user",
			giveIdpUsers: []identity.IDPUser{
				{
					ID:        "user1",
					FirstName: "josh",
					LastName:  "wilkes",
					Email:     "josh@test.go",
					Groups: []string{
						"internalEveryoneId",
					},
				},
			},
			giveIdpGroups: []identity.IDPGroup{
				{
					ID:          "internalEveryoneId",
					Name:        "everyone",
					Description: "a description",
				},
			},
			giveInternalUsers:  []identity.User{},
			giveInternalGroups: []identity.Group{},
			wantUserMap: map[string]identity.User{
				"josh@test.go": {
					ID:        "_",
					FirstName: "josh",
					LastName:  "wilkes",
					Email:     "josh@test.go",
					Status:    types.IdpStatusACTIVE,
				},
			},
			wantGroupMap: map[string]identity.Group{
				"internalEveryoneId": {
					ID:          "_",
					IdpID:       "internalEveryoneId",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
				},
			},
		},
		{
			// Archiving should set the status to archived and remove goup and user associations
			name:          "groups and users archived correctly",
			giveIdpUsers:  []identity.IDPUser{},
			giveIdpGroups: []identity.IDPGroup{},
			giveInternalUsers: []identity.User{{
				ID:        "abcd",
				FirstName: "josh",
				LastName:  "wilkes",
				Email:     "josh@test.go",
				Groups:    []string{"1234"},
				Status:    types.IdpStatusACTIVE,
				CreatedAt: now,
				UpdatedAt: now,
			}},
			giveInternalGroups: []identity.Group{{
				ID:          "1234",
				IdpID:       "internalEveryoneId",
				Name:        "everyone",
				Description: "a description",
				Status:      types.IdpStatusACTIVE,
				Users:       []string{"abcd"},
				CreatedAt:   now,
				UpdatedAt:   now,
			}},
			wantUserMap: map[string]identity.User{
				"josh@test.go": {
					ID:        "abcd",
					FirstName: "josh",
					LastName:  "wilkes",
					Email:     "josh@test.go",
					Groups:    []string{},
					Status:    types.IdpStatusARCHIVED,
				},
			},
			wantGroupMap: map[string]identity.Group{
				"internalEveryoneId": {
					ID:          "1234",
					IdpID:       "internalEveryoneId",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusARCHIVED,
					Users:       []string{},
					CreatedAt:   now,
				},
			},
		},
		{
			// This test case tests updating attributes of users and groups
			// it also tests archiving users and groups
			// it tests that existing users and groups are updated correctly when groups are archived
			name: "archive group and user when they no longer exist update other users and groups",
			giveIdpUsers: []identity.IDPUser{{
				ID:        "user2",
				FirstName: "larry",
				LastName:  "browner",
				Email:     "larry@test.go",
				Groups: []string{
					"larrysGroupId",
				},
			}},
			giveIdpGroups: []identity.IDPGroup{{
				ID:          "larrysGroupId",
				Name:        "larrysGroupNewName",
				Description: "a different description",
				Source:      "OKTA",
			}},
			giveInternalUsers: []identity.User{
				{
					ID:        "efgh",
					FirstName: "larry",
					LastName:  "brown",
					Email:     "larry@test.go",
					Groups:    []string{"1234", "5678"},
					Status:    types.IdpStatusACTIVE,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "abcd",
					FirstName: "josh",
					LastName:  "wilkes",
					Email:     "josh@test.go",
					Groups:    []string{"1234"},
					Status:    types.IdpStatusACTIVE,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			giveInternalGroups: []identity.Group{
				{
					ID:          "1234",
					IdpID:       "internalEveryoneId",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
					Users:       []string{"abcd", "efgh"},
					CreatedAt:   now,
					UpdatedAt:   now,
					Source:      "OKTA",
				},
				{
					ID:          "5678",
					IdpID:       "larrysGroupId",
					Name:        "larrysGroup",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
					Users:       []string{"efgh"},
					CreatedAt:   now,
					UpdatedAt:   now,
					Source:      "OKTA",
				},
			},
			wantUserMap: map[string]identity.User{
				"josh@test.go": {
					ID:        "abcd",
					FirstName: "josh",
					LastName:  "wilkes",
					Email:     "josh@test.go",
					Groups:    []string{},
					Status:    types.IdpStatusARCHIVED,
				},
				"larry@test.go": {
					ID:        "efgh",
					FirstName: "larry",
					LastName:  "browner",
					Email:     "larry@test.go",
					Groups:    []string{"5678"},
					Status:    types.IdpStatusACTIVE,
					CreatedAt: now,
				},
			},
			wantGroupMap: map[string]identity.Group{
				"internalEveryoneId": {
					ID:          "1234",
					IdpID:       "internalEveryoneId",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusARCHIVED,
					Users:       []string{},
					CreatedAt:   now,
					Source:      "OKTA",
				},
				"larrysGroupId": {
					ID:          "5678",
					IdpID:       "larrysGroupId",
					Name:        "larrysGroupNewName",
					Description: "a different description",
					Status:      types.IdpStatusACTIVE,
					Users:       []string{"efgh"},
					CreatedAt:   now,
					Source:      "OKTA",
				},
			},
		},
		{
			// a group should be dearchived if it exists again
			name:         "dearchive group when it exists again",
			giveIdpUsers: []identity.IDPUser{},
			giveIdpGroups: []identity.IDPGroup{{
				ID:          "granted_administrators",
				Name:        "granted_administrators",
				Description: "admin group",
			}},
			giveInternalUsers: []identity.User{},
			giveInternalGroups: []identity.Group{
				{
					ID:          "1234",
					IdpID:       "granted_administrators",
					Name:        "granted_administrators",
					Description: "admin group",
					Status:      types.IdpStatusARCHIVED,
					Users:       []string{},
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantUserMap: map[string]identity.User{},
			wantGroupMap: map[string]identity.Group{
				"granted_administrators": {
					ID:          "1234",
					IdpID:       "granted_administrators",
					Name:        "granted_administrators",
					Description: "admin group",
					Status:      types.IdpStatusACTIVE,
					Users:       []string{},
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotUsers, gotGroups := processUsersAndGroups(tc.giveIdpUsers, tc.giveIdpGroups, tc.giveInternalUsers, tc.giveInternalGroups)
			for k, u := range tc.wantUserMap {
				got := gotUsers[k]
				u.ID = got.ID
				if u.Groups == nil {
					u.Groups = got.Groups
				}
				if u.CreatedAt.IsZero() {
					u.CreatedAt = got.CreatedAt
				}
				if u.UpdatedAt.IsZero() {
					u.UpdatedAt = got.UpdatedAt
				}

				tc.wantUserMap[k] = u
			}

			for k, g := range tc.wantGroupMap {
				got := gotGroups[k]
				g.ID = got.ID
				if g.CreatedAt.IsZero() {
					g.CreatedAt = got.CreatedAt
				}
				if g.UpdatedAt.IsZero() {
					g.UpdatedAt = got.UpdatedAt
				}
				if g.Users == nil {
					g.Users = got.Users
				}
				tc.wantGroupMap[k] = g
			}

			assert.Exactly(t, tc.wantUserMap, gotUsers)
			assert.Exactly(t, tc.wantGroupMap, gotGroups)
		})
	}
}
