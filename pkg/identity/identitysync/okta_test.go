package identitysync

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"github.com/stretchr/testify/assert"
)

type DefaultResponse struct{}

func (d DefaultResponse) HasNextPage() bool {
	return false
}

func (d DefaultResponse) Next(ctx context.Context, v interface{}) (*okta.Response, error) {
	return nil, errors.New("you can not call the method if HasNextPage returns false")
}

var defaultResponse DefaultResponse

type mockOktaClient struct {
	listGroups     func() []*okta.Group
	listUserGroups func(userId string) []*okta.Group
	listUsers      func() []*okta.User
	listGroupUsers func(groupId string) []*okta.User
}

func (m mockOktaClient) ListUserGroups(ctx context.Context, userId string) ([]*okta.Group, OktaResponse, error) {
	return m.listUserGroups(userId), defaultResponse, nil
}

func (m mockOktaClient) ListUsers(ctx context.Context, qp *query.Params) ([]*okta.User, OktaResponse, error) {
	return m.listUsers(), defaultResponse, nil
}

func (m mockOktaClient) ListGroupUsers(ctx context.Context, groupId string, qp *query.Params) ([]*okta.User, OktaResponse, error) {
	return m.listGroupUsers(groupId), defaultResponse, nil
}

func (m mockOktaClient) ListGroups(ctx context.Context, qp *query.Params) ([]*okta.Group, OktaResponse, error) {
	return m.listGroups(), defaultResponse, nil
}

var listOfGroups = []*okta.Group{
	{
		Id: "group1",
		Profile: &okta.GroupProfile{
			Description: "desc1",
			Name:        "name1",
		},
	},
	{
		Id: "group2",
		Profile: &okta.GroupProfile{
			Description: "desc2",
			Name:        "name2",
		},
	},
	{
		Id: "group3",
		Profile: &okta.GroupProfile{
			Description: "desc3",
			Name:        "name3",
		},
	},
}

var listOfUsers = []*okta.User{
	{
		Id: "user1",
		Profile: &okta.UserProfile{
			"firstName": "name1",
			"lastName":  "surname1",
			"email":     "email1@email.com",
		},
	},
	{
		Id: "user2",
		Profile: &okta.UserProfile{
			"firstName": "name2",
			"lastName":  "surname2",
			"email":     "email2@email.com",
		},
	},
	{
		Id: "user3",
		Profile: &okta.UserProfile{
			"firstName": "name3",
			"lastName":  "surname3",
			"email":     "email3@email.com",
		},
	},
}

func TestOktaSync_ListGroups(t *testing.T) {

	oktaSync := OktaSync{
		client: &mockOktaClient{
			listGroups: func() []*okta.Group {
				return listOfGroups
			},
		},
	}

	groups, err := oktaSync.ListGroups(context.Background())
	assert.NoError(t, err)
	assert.Len(t, groups, 3)

	for i := 0; i < 3; i++ {
		assert.Equal(t, listOfGroups[i].Id, groups[i].ID)
		assert.Equal(t, listOfGroups[i].Profile.Name, groups[i].Name)
		assert.Equal(t, listOfGroups[i].Profile.Description, groups[i].Description)
	}

	// group config should NOT affect ListGroups
	oktaSync.groups.Set("group1")
	groups, err = oktaSync.ListGroups(context.Background())
	assert.NoError(t, err)
	assert.Len(t, groups, 3)

	for i := 0; i < 3; i++ {
		assert.Equal(t, listOfGroups[i].Id, groups[i].ID)
		assert.Equal(t, listOfGroups[i].Profile.Name, groups[i].Name)
		assert.Equal(t, listOfGroups[i].Profile.Description, groups[i].Description)
	}
}

func TestOktaSync_ListUsersWithoutGroupConfig(t *testing.T) {
	oktaSync := OktaSync{
		client: &mockOktaClient{
			listUserGroups: func(userId string) []*okta.Group {
				switch userId {
				case listOfUsers[0].Id:
					return []*okta.Group{listOfGroups[0]}
				case listOfUsers[1].Id:
					return []*okta.Group{listOfGroups[0], listOfGroups[1]}
				case listOfUsers[2].Id:
					return []*okta.Group{listOfGroups[0], listOfGroups[1], listOfGroups[2]}
				default:
					return []*okta.Group{}
				}
			},
			listUsers: func() []*okta.User {
				return listOfUsers
			},
		},
	}

	users, err := oktaSync.ListUsers(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 3)

	for i := 0; i < 3; i++ {
		assert.Equal(t, listOfUsers[i].Id, users[i].ID)
		assert.Equal(t, (*listOfUsers[i].Profile)["firstName"].(string), users[i].FirstName)
		assert.Equal(t, (*listOfUsers[i].Profile)["lastName"].(string), users[i].LastName)
		assert.Equal(t, (*listOfUsers[i].Profile)["email"].(string), users[i].Email)
		assert.Equal(t, i+1, len(users[i].Groups))
	}
}

func TestOktaSync_ListUsersWithGroupConfig(t *testing.T) {
	oktaSync := OktaSync{
		client: &mockOktaClient{
			listGroupUsers: func(groupId string) []*okta.User {
				switch groupId {
				case listOfGroups[0].Id:
					return []*okta.User{listOfUsers[0]}
				case listOfGroups[1].Id:
					return []*okta.User{listOfUsers[1]}
				case listOfGroups[2].Id:
					return []*okta.User{listOfUsers[2]}
				default:
					return []*okta.User{}
				}
			},
		},
	}

	// any combination of 1 group should return 1 user
	for i := 0; i < 3; i++ {
		// group config should NOT affect ListGroups
		oktaSync.groups.Set(listOfGroups[i].Id)
		users, err := oktaSync.ListUsers(context.Background())
		assert.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, 1, len(users[0].Groups))
		assert.Equal(t, listOfUsers[i].Id, users[0].ID)
		assert.Equal(t, listOfGroups[i].Id, users[0].Groups[0])
	}

	// any combination of 2 groups should return 2 users
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			oktaSync.groups.Set(fmt.Sprintf("%s,%s", listOfGroups[i].Id, listOfGroups[j].Id))
			users, err := oktaSync.ListUsers(context.Background())
			assert.NoError(t, err)
			assert.Len(t, users, 2)
		}
	}
}

func TestOktaSync_Wildcard(t *testing.T) {
	oktaSync := OktaSync{
		client: &mockOktaClient{
			listGroupUsers: func(groupId string) []*okta.User {
				switch groupId {
				case listOfGroups[0].Id:
					return []*okta.User{listOfUsers[0]}
				case listOfGroups[1].Id:
					return []*okta.User{listOfUsers[0], listOfUsers[1]}
				case listOfGroups[2].Id:
					return []*okta.User{listOfUsers[0], listOfUsers[1], listOfUsers[2]}
				default:
					return []*okta.User{}
				}
			},
			listGroups: func() []*okta.Group {
				return listOfGroups
			},
		},
	}

	oktaSync.groups.Set("*")
	users, err := oktaSync.ListUsers(context.Background())

	assert.NoError(t, err)
	assert.Len(t, users, 3)

	// Check users groups by ids, ignoring the order
	for _, u := range users {
		switch u.ID {
		case listOfUsers[0].Id:
			assert.ElementsMatch(t, u.Groups, []string{listOfGroups[0].Id, listOfGroups[1].Id, listOfGroups[2].Id})
		case listOfUsers[1].Id:
			assert.ElementsMatch(t, u.Groups, []string{listOfGroups[1].Id, listOfGroups[2].Id})
		case listOfUsers[2].Id:
			assert.ElementsMatch(t, u.Groups, []string{listOfGroups[2].Id})
		}
	}
}
