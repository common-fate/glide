package identitysync

import (
	"sort"
	"testing"
	"time"

	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/stretchr/testify/assert"
)

// The processor contains all the mapping logic for create/update/map for users and groups
func TestIdentitySyncProcessor(t *testing.T) {

	type testcase struct {
		name                 string
		giveIdpUsers         []identity.IDPUser
		giveIdpGroups        []identity.IDPGroup
		giveInternalUsers    []identity.User
		giveInternalGroups   []identity.Group
		wantUserMap          map[string]identity.User
		wantGroupMap         map[string]identity.Group
		withIdpType          string
		useIdpGroupsAsFilter bool
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
			}},
			giveInternalUsers: []identity.User{
				{
					ID:        "efgh",
					FirstName: "larry",
					LastName:  "brown",
					Email:     "larry@test.go",
					Groups:    []string{"internalEveryoneId", "larrysGroupId"},
					Status:    types.IdpStatusACTIVE,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "abcd",
					FirstName: "josh",
					LastName:  "wilkes",
					Email:     "josh@test.go",
					Groups:    []string{"internalEveryoneId"},
					Status:    types.IdpStatusACTIVE,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			giveInternalGroups: []identity.Group{
				{
					ID:          "internalEveryoneId",
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
					ID:          "larrysGroupId",
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
					Groups:    []string{"larrysGroupId"},
					Status:    types.IdpStatusACTIVE,
					CreatedAt: now,
				},
			},
			wantGroupMap: map[string]identity.Group{
				"internalEveryoneId": {
					ID:          "internalEveryoneId",
					IdpID:       "internalEveryoneId",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusARCHIVED,
					Users:       []string{},
					CreatedAt:   now,
					Source:      "OKTA",
				},
				"larrysGroupId": {
					ID:          "larrysGroupId",
					IdpID:       "larrysGroupId",
					Name:        "larrysGroupNewName",
					Description: "a different description",
					Status:      types.IdpStatusACTIVE,
					Users:       []string{"efgh"},
					CreatedAt:   now,
					Source:      "OKTA",
				},
			},
			withIdpType: "OKTA",
		},
		{
			// a group should be dearchived if it exists again
			name:         "dearchive group when it exists again",
			giveIdpUsers: []identity.IDPUser{},
			giveIdpGroups: []identity.IDPGroup{{
				ID:          "common_fate_administrators",
				Name:        "common_fate_administrators",
				Description: "admin group",
			}},
			giveInternalUsers: []identity.User{},
			giveInternalGroups: []identity.Group{
				{
					ID:          "1234",
					IdpID:       "common_fate_administrators",
					Name:        "common_fate_administrators",
					Description: "admin group",
					Status:      types.IdpStatusARCHIVED,
					Users:       []string{},
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
			wantUserMap: map[string]identity.User{},
			wantGroupMap: map[string]identity.Group{
				"common_fate_administrators": {
					ID:          "1234",
					IdpID:       "common_fate_administrators",
					Name:        "common_fate_administrators",
					Description: "admin group",
					Status:      types.IdpStatusACTIVE,
					Users:       []string{},
					CreatedAt:   now,
					UpdatedAt:   now,
				},
			},
		},
		{
			name: "user with non-existent group",
			giveIdpUsers: []identity.IDPUser{
				{
					ID:        "user1",
					FirstName: "josh",
					LastName:  "wilkes",
					Email:     "josh@test.go",
					Groups: []string{
						"internalEveryoneId",
						"doesnt-exist",
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
			giveInternalUsers: []identity.User{
				{
					ID:        "user1",
					FirstName: "josh",
					LastName:  "wilkes",
					Email:     "josh@test.go",
					Status:    types.IdpStatusACTIVE,
				},
			},
			giveInternalGroups: []identity.Group{
				{
					ID:          "internalEveryoneId",
					IdpID:       "internalEveryoneId",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
				},
			},
			wantUserMap: map[string]identity.User{
				"josh@test.go": {
					ID:        "user1",
					FirstName: "josh",
					LastName:  "wilkes",
					Email:     "josh@test.go",
					Status:    types.IdpStatusACTIVE,
				},
			},
			wantGroupMap: map[string]identity.Group{
				"internalEveryoneId": {
					ID:          "internalEveryoneId",
					IdpID:       "internalEveryoneId",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
					Users: []string{
						"user1",
					},
				},
			},
		},
		{
			/*
				groups
				{name: "admins", id: "1"}
				{name: "dev_ops", id: "2"}

				users
				{name: "bob", id: "1", groups: ["1", "2"]}  // grant access
				{name: "alice", id: "2", groups: ["3"]}	// deny access
				{name: "joe", id: "3", groups: ["1", "3"]} // grant access

				{{{ POST processUsersAndGroups }}}

				groups
				{name: "admins", id: "1"}
				{name: "dev_ops", id: "2"}

				users
				{name: "bob", id: "1", groups: ["1", "2"]}  // grant access
				{name: "joe", id: "3", groups: ["1"]} // grant access

			*/
			name: "user with non-existent group",
			giveIdpUsers: []identity.IDPUser{
				{
					ID:        "user1",
					FirstName: "bob",
					Email:     "bob@mail.com",
					Groups: []string{
						"admins",
						"dev_ops",
					},
				},
				{
					ID:        "user2",
					FirstName: "alice",
					Email:     "alice@mail.com",
					Groups: []string{
						"accounting-ie-not-included-filtered-groups",
					},
				},
				{
					ID:        "user3",
					FirstName: "joe",
					Email:     "joe@mail.com",
					Groups: []string{
						"admins",
						"accounting-ie-not-included-filtered-groups",
					},
				},
			},
			giveIdpGroups: []identity.IDPGroup{
				{
					ID:          "admins",
					Name:        "everyone",
					Description: "a description",
				},
				{
					ID:          "dev_ops",
					Name:        "everyone",
					Description: "a description",
				},
			},
			giveInternalUsers: []identity.User{
				{
					ID:        "user1",
					FirstName: "bob",
					Email:     "bob@mail.com",
					Status:    types.IdpStatusACTIVE,
				},
				{
					ID:        "user3",
					FirstName: "joe",
					Email:     "joe@mail.com",
					Status:    types.IdpStatusACTIVE,
				},
			},
			giveInternalGroups: []identity.Group{
				{
					ID:          "admins",
					IdpID:       "admins",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
				},
				{
					ID:          "dev_ops",
					IdpID:       "dev_ops",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
				},
			},
			wantUserMap: map[string]identity.User{
				"bob@mail.com": {
					ID:        "user1",
					FirstName: "bob",
					Email:     "bob@mail.com",
					Groups: []string{
						"admins",
						"dev_ops",
					},
					Status: types.IdpStatusACTIVE,
				},
				"joe@mail.com": {
					ID:        "user3",
					FirstName: "joe",
					Email:     "joe@mail.com",
					Groups:    []string{"admins"},
					Status:    types.IdpStatusACTIVE,
				},
				"alice@mail.com": {
					ID:        "user3",
					FirstName: "alice",
					Email:     "alice@mail.com",
					Groups:    []string{},
					Status:    types.IdpStatusARCHIVED,
				},
			},
			wantGroupMap: map[string]identity.Group{
				"admins": {
					ID:          "admins",
					IdpID:       "admins",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
					Users: []string{
						"user1",
						"user3",
					},
				},
				"dev_ops": {
					ID:          "dev_ops",
					IdpID:       "dev_ops",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
					Users: []string{
						"user1",
					},
				},
			},
			useIdpGroupsAsFilter: true,
		},
		{
			name: "user must be updated it ACTIVE if present in IDP, but archived internally",
			giveIdpUsers: []identity.IDPUser{
				{
					ID:        "user1",
					FirstName: "bob",
					Email:     "bob@mail.com",
					Groups: []string{
						"admins",
					},
				},
			},
			giveIdpGroups: []identity.IDPGroup{
				{
					ID:          "admins",
					Name:        "everyone",
					Description: "a description",
				},
			},
			giveInternalUsers: []identity.User{
				{
					ID:        "user1",
					FirstName: "bob",
					Email:     "bob@mail.com",
					Status:    types.IdpStatusARCHIVED,
				},
			},
			giveInternalGroups: []identity.Group{
				{
					ID:          "admins",
					IdpID:       "admins",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
				},
			},
			wantUserMap: map[string]identity.User{
				"bob@mail.com": {
					ID:        "user1",
					FirstName: "bob",
					Email:     "bob@mail.com",
					Groups: []string{
						"admins",
					},
					Status: types.IdpStatusACTIVE,
				},
			},
			wantGroupMap: map[string]identity.Group{
				"admins": {
					ID:          "admins",
					IdpID:       "admins",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
					Users: []string{
						"user1",
					},
				},
			},
			useIdpGroupsAsFilter: true,
		},
		{
			name: "group must be updated to ARCHIVED if NOT present in IDP, but ACTIVE internally",
			giveIdpUsers: []identity.IDPUser{
				{
					ID:        "user1",
					FirstName: "bob",
					Email:     "bob@mail.com",
					Groups: []string{
						"admins",
						"group_no_longer_in_idp",
					},
				},
			},
			giveIdpGroups: []identity.IDPGroup{
				{
					ID:          "admins",
					Name:        "everyone",
					Description: "a description",
				},
			},
			giveInternalUsers: []identity.User{
				{
					ID:        "user1",
					FirstName: "bob",
					Email:     "bob@mail.com",
					Status:    types.IdpStatusACTIVE,
				},
			},
			giveInternalGroups: []identity.Group{
				{
					ID:          "admins",
					IdpID:       "admins",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
				},
				{
					ID:          "group_no_longer_in_idp",
					IdpID:       "group_no_longer_in_idp",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
				},
			},
			wantUserMap: map[string]identity.User{
				"bob@mail.com": {
					ID:        "user1",
					FirstName: "bob",
					Email:     "bob@mail.com",
					Groups: []string{
						"admins",
					},
					Status: types.IdpStatusACTIVE,
				},
			},
			wantGroupMap: map[string]identity.Group{
				"admins": {
					ID:          "admins",
					IdpID:       "admins",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
					Users: []string{
						"user1",
					},
				},
				"group_no_longer_in_idp": {
					ID:          "group_no_longer_in_idp",
					IdpID:       "group_no_longer_in_idp",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusARCHIVED,
					Users:       []string{},
				},
			},
			useIdpGroupsAsFilter: true,
		},
		{
			name: "user must be ARCHIVED if NOT in ACTIVE idp group",
			giveIdpUsers: []identity.IDPUser{
				{
					ID:        "user1",
					FirstName: "bob",
					Email:     "bob@mail.com",
					Groups: []string{
						"group_no_longer_in_idp",
					},
				},
			},
			giveIdpGroups: []identity.IDPGroup{
				{
					ID:          "admins",
					Name:        "everyone",
					Description: "a description",
				},
			},
			giveInternalUsers: []identity.User{
				{
					ID:        "user1",
					FirstName: "bob",
					Email:     "bob@mail.com",
					Status:    types.IdpStatusACTIVE,
				},
			},
			giveInternalGroups: []identity.Group{
				{
					ID:          "admins",
					IdpID:       "admins",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
					Source:      identity.INTERNAL,
					Users:       []string{},
				},
				{
					ID:          "group_no_longer_in_idp",
					IdpID:       "group_no_longer_in_idp",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
					Users:       []string{},
					Source:      identity.INTERNAL,
				},
			},
			wantUserMap: map[string]identity.User{
				"bob@mail.com": {
					ID:        "user1",
					FirstName: "bob",
					Email:     "bob@mail.com",
					Groups:    []string{},
					Status:    types.IdpStatusARCHIVED,
				},
			},
			wantGroupMap: map[string]identity.Group{
				"admins": {
					ID:          "admins",
					IdpID:       "admins",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusACTIVE,
					Users:       []string{},
					Source:      identity.INTERNAL,
				},
				"group_no_longer_in_idp": {
					ID:          "group_no_longer_in_idp",
					IdpID:       "group_no_longer_in_idp",
					Name:        "everyone",
					Description: "a description",
					Status:      types.IdpStatusARCHIVED,
					Users:       []string{},
					Source:      identity.INTERNAL,
				},
			},
			withIdpType:          identity.INTERNAL,
			useIdpGroupsAsFilter: true,
		},
	}
	for _, tc := range testcases {

		t.Run(tc.name, func(t *testing.T) {
			gotUsers, gotGroups := processUsersAndGroups(tc.withIdpType, tc.giveIdpUsers, tc.giveIdpGroups, tc.giveInternalUsers, tc.giveInternalGroups, tc.useIdpGroupsAsFilter)
			for k, u := range tc.wantUserMap {
				got := gotUsers[k]
				u.ID = got.ID
				if u.Groups == nil {
					u.Groups = got.Groups
					sort.Strings(u.Groups)
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
					// sort got.Users
					sort.Strings(g.Users)
				}
				tc.wantGroupMap[k] = g
			}

			for _, g := range gotGroups {
				sort.Strings(g.Users)
			}

			for _, u := range gotUsers {
				sort.Strings(u.Groups)
			}

			assert.Exactly(t, tc.wantUserMap, gotUsers)
			assert.Exactly(t, tc.wantGroupMap, gotGroups)
		})
	}
}
