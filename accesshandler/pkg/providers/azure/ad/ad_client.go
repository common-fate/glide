package ad

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type ListUsersResponse struct {
	OdataContext  string      `json:"@odata.context"`
	OdataNextLink *string     `json:"@odata.nextLink,omitempty"`
	Value         []AzureUser `json:"value"`
}

type AzureUser struct {
	GivenName string `json:"givenName"`
	Mail      string `json:"mail"`
	Surname   string `json:"surname"`
	ID        string `json:"id"`
}

type AzureGroup struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	DisplayName string `json:"displayName"`
}

type UserGroups struct {
	OdataNextLink *string  `json:"@odata.nextLink,omitempty"`
	OdataContext  string   `json:"@odata.context"`
	Value         []string `json:"value"`
}

type ListGroupsResponse struct {
	OdataContext  string       `json:"@odata.context"`
	OdataNextLink *string      `json:"@odata.nextLink,omitempty"`
	Value         []AzureGroup `json:"value"`
}

type GroupMembers struct {
	OdataNextLink *string  `json:"@odata.nextLink,omitempty"`
	OdataContext  string   `json:"@odata.context"`
	Value         []string `json:"value"`
}

func (c *Provider) ListUsers(ctx context.Context) ([]AzureUser, error) {

	//get all users
	idpUsers := []AzureUser{}
	hasMore := true
	var nextToken *string
	url := MSGraphBaseURL + "/users"

	for hasMore {

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+c.token.Get())
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		//return the error if its anything but a 200
		if res.StatusCode != 200 {
			return nil, fmt.Errorf(string(b))
		}

		var lu ListUsersResponse
		err = json.Unmarshal(b, &lu)
		if err != nil {
			return nil, err
		}

		idpUsers = append(idpUsers, lu.Value...)

		nextToken = lu.OdataNextLink
		if nextToken != nil {
			url = *nextToken
		} else {
			hasMore = false
		}

	}

	return idpUsers, nil
}

func (c *Provider) ListGroups(context.Context) ([]AzureGroup, error) {
	idpGroups := []AzureGroup{}
	hasMore := true
	var nextToken *string
	url := MSGraphBaseURL + "/groups"
	for hasMore {

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+c.token.Get())
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		//return the error if its anything but a 200
		if res.StatusCode != 200 {
			return nil, fmt.Errorf(string(b))
		}

		var lu ListGroupsResponse
		err = json.Unmarshal(b, &lu)
		if err != nil {
			return nil, err
		}

		idpGroups = append(idpGroups, lu.Value...)

		nextToken = lu.OdataNextLink
		if nextToken != nil {
			url = *nextToken
		} else {
			hasMore = false
		}
	}
	return idpGroups, nil
}

func (c *Provider) GetGroup(ctx context.Context, groupID string) (*AzureGroup, error) {

	url := MSGraphBaseURL + "/groups/" + groupID

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+c.token.Get())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//return the error if its anything but a 200
	if res.StatusCode != 200 {
		return nil, fmt.Errorf(string(b))
	}

	var g AzureGroup
	err = json.Unmarshal(b, &g)
	if err != nil {
		return nil, err
	}

	return &g, nil

}

func (c *Provider) GetUser(ctx context.Context, userID string) (*AzureUser, error) {

	url := MSGraphBaseURL + "/users/" + userID

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+c.token.Get())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	//return the error if its anything but a 200
	if res.StatusCode != 200 {
		return nil, fmt.Errorf(string(b))
	}

	var u AzureUser
	err = json.Unmarshal(b, &u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

type AddUser struct {
	Key string `json:"@odata.id"`
}

//GroupMember.ReadWrite.All
func (c *Provider) AddUserToGroup(ctx context.Context, userID string, groupID string) error {

	url := MSGraphBaseURL + "/groups/" + groupID + "/members/$ref"

	a := AddUser{Key: "https://graph.microsoft.com/v1.0/directoryObjects/" + userID}
	out, err := json.Marshal(a)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(out))
	req.Header.Add("Authorization", "Bearer "+c.token.Get())
	req.Header.Add("Content-Type", "application/json")

	zap.S().Info("adding user to group")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	//return the error if its anything but a 204
	if res.StatusCode != 204 {
		return fmt.Errorf(string(b))
	}

	return nil
}

//GroupMember.ReadWrite.All
func (c *Provider) RemoveUserFromGroup(ctx context.Context, userID string, groupID string) error {

	url := MSGraphBaseURL + "/groups/" + groupID + "/members/" + userID + "/$ref"

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Bearer "+c.token.Get())

	zap.S().Info("removing user to group")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	//return the error if its anything but a 204
	if res.StatusCode != 204 {
		return fmt.Errorf(string(b))
	}
	return nil
}

//GroupMember.Read.All
func (c *Provider) ListGroupUsers(ctx context.Context, groupID string) ([]AzureUser, error) {

	var groupMembers []AzureUser

	hasMore := true
	var nextToken *string
	url := MSGraphBaseURL + fmt.Sprintf("/groups/%s/members", groupID)

	for hasMore {
		var jsonStr = []byte(`{ "securityEnabledOnly": false}`)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Add("Authorization", "Bearer "+c.token.Get())
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		//return the error if its anything but a 200
		if res.StatusCode != 200 {
			return nil, fmt.Errorf(string(b))
		}

		var lu ListUsersResponse
		err = json.Unmarshal(b, &lu)
		if err != nil {
			return nil, err
		}

		groupMembers = append(groupMembers, lu.Value...)

		nextToken = lu.OdataNextLink
		if nextToken != nil {
			url = *nextToken
		} else {
			hasMore = false
		}

	}
	return groupMembers, nil
}

type CreateADUser struct {
	AccountEnabled    bool            `json:"accountEnabled"`
	DisplayName       string          `json:"displayName"`
	MailNickname      string          `json:"mailNickname"`
	UserPrincipalName string          `json:"userPrincipalName"`
	PasswordProfile   PasswordProfile `json:"passwordProfile"`
}

type PasswordProfile struct {
	ForceChangePasswordNextSignIn bool   `json:"forceChangePasswordNextSignIn"`
	Password                      string `json:"password"`
}

func (c *Provider) CreateUser(ctx context.Context, user CreateADUser) error {
	url := MSGraphBaseURL + "/users"

	out, err := json.Marshal(user)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(out))
	req.Header.Add("Authorization", "Bearer "+c.token.Get())
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	//return the error if its anything but a 201
	if res.StatusCode != 201 {
		return fmt.Errorf(string(b))
	}
	return nil

}

func (c *Provider) DeleteUser(ctx context.Context, userID string) error {
	url := MSGraphBaseURL + "/users/" + userID

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Bearer "+c.token.Get())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	//return the error if its anything but a 204
	if res.StatusCode != 204 {
		return fmt.Errorf(string(b))
	}
	return nil
}

type CreateADGroup struct {
	Description     string   `json:"description"`
	DisplayName     string   `json:"displayName"`
	GroupTypes      []string `json:"groupTypes"`
	MailEnabled     bool     `json:"mailEnabled"`
	MailNickname    string   `json:"mailNickname"`
	SecurityEnabled bool     `json:"securityEnabled"`
}

type CreateADGroupResponse struct {
	ID string `json:"id"`

	Description string `json:"description"`
	DisplayName string `json:"displayName"`
}

func (c *Provider) CreateGroup(ctx context.Context, group CreateADGroup) (*CreateADGroupResponse, error) {
	url := MSGraphBaseURL + "/groups"

	out, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(out))
	req.Header.Add("Authorization", "Bearer "+c.token.Get())
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var groupRes CreateADGroupResponse
	err = json.Unmarshal(b, &groupRes)
	if err != nil {
		return nil, err
	}
	//return the error if its anything but a 201
	if res.StatusCode != 201 {
		return nil, fmt.Errorf(string(b))
	}
	return &groupRes, nil
}

func (c *Provider) DeleteGroup(ctx context.Context, groupID string) error {
	url := MSGraphBaseURL + "/groups/" + groupID

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Bearer "+c.token.Get())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	//return the error if its anything but a 204
	if res.StatusCode != 204 {
		return fmt.Errorf(string(b))
	}
	return nil
}
