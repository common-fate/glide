package identitysync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/granted-approvals/pkg/identity"
)

const MSGraphBaseURL = "https://graph.microsoft.com/v1.0"
const ADAuthorityHost = "https://login.microsoftonline.com"

type AzureSync struct {
	// This is initialised during the Init function call and is not saved in config
	token        gconfig.SecretStringValue
	tenantID     gconfig.StringValue
	clientID     gconfig.StringValue
	clientSecret gconfig.SecretStringValue
}

func (s *AzureSync) Config() gconfig.Config {
	return gconfig.Config{
		gconfig.StringField("tenantId", &s.tenantID, "the Azure AD tenant ID"),
		gconfig.StringField("clientId", &s.clientID, "the Azure AD client ID"),
		gconfig.SecretStringField("clientSecret", &s.clientSecret, "the Azure AD client secret", gconfig.WithNoArgs("/granted/secrets/identity/azure/secret")),
	}
}

func (s *AzureSync) Init(ctx context.Context) error {
	cred, err := confidential.NewCredFromSecret(s.clientSecret.Get())
	if err != nil {
		return err
	}
	c, err := confidential.New(s.clientID.Get(), cred,
		confidential.WithAuthority(fmt.Sprintf("%s/%s", ADAuthorityHost, s.tenantID.Get())))
	if err != nil {
		return err
	}
	token, err := c.AcquireTokenByCredential(ctx, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return err
	}
	s.token.Set(token.AccessToken)
	return nil
}

type ListUsersResponse struct {
	OdataContext  string      `json:"@odata.context"`
	OdataNextLink *string     `json:"@odata.nextLink,omitempty"`
	Value         []AzureUser `json:"value"`
}

// properties of a user in the graph API
//
// https://docs.microsoft.com/en-us/graph/api/resources/user?view=graph-rest-1.0#properties
type AzureUser struct {
	GivenName string `json:"givenName"`
	// this maps to a users email by convention
	// see the graph API spec for details
	// in practive all users have a principal name but some users may not have the "mail" property for different reasons.
	// we use this for the email
	UserPrincipalName string `json:"userPrincipalName"`
	Surname           string `json:"surname"`
	ID                string `json:"id"`
}

type ListGroupsResponse struct {
	OdataContext  string       `json:"@odata.context"`
	OdataNextLink *string      `json:"@odata.nextLink,omitempty"`
	Value         []AzureGroup `json:"value"`
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

// idpUserFromAzureUser converts a azure user to the identityprovider interface user type
func (a *AzureSync) idpUserFromAzureUser(ctx context.Context, azureUser AzureUser) (identity.IdpUser, error) {
	u := identity.IdpUser{
		ID:        azureUser.ID,
		FirstName: azureUser.GivenName,
		LastName:  azureUser.Surname,
		Email:     azureUser.UserPrincipalName,
		Groups:    []string{},
	}

	g, err := a.GetMemberGroups(u.ID)
	if err != nil {
		return identity.IdpUser{}, err
	}
	u.Groups = g

	return u, nil
}

func (a *AzureSync) GetMemberGroups(userID string) ([]string, error) {
	var userGroups []string

	hasMore := true
	var nextToken *string
	url := MSGraphBaseURL + fmt.Sprintf("/directoryObjects/%s/getMemberGroups", userID)

	for hasMore {
		var jsonStr = []byte(`{ "securityEnabledOnly": false}`)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Add("Authorization", "Bearer "+a.token.Get())
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

		var lu UserGroups
		err = json.Unmarshal(b, &lu)
		if err != nil {
			return nil, err
		}

		userGroups = append(userGroups, lu.Value...)

		nextToken = lu.OdataNextLink
		if nextToken != nil {
			url = *nextToken
		} else {
			hasMore = false
		}

	}
	return userGroups, nil
}

func (a *AzureSync) ListUsers(ctx context.Context) ([]identity.IdpUser, error) {

	//get all users
	idpUsers := []identity.IdpUser{}
	hasMore := true
	var nextToken *string
	url := MSGraphBaseURL + "/users"

	for hasMore {

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+a.token.Get())
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

		for _, u := range lu.Value {

			user, err := a.idpUserFromAzureUser(ctx, u)
			if err != nil {
				return nil, err
			}
			idpUsers = append(idpUsers, user)
		}
		nextToken = lu.OdataNextLink
		if nextToken != nil {
			url = *nextToken
		} else {
			hasMore = false
		}

	}

	return idpUsers, nil
}

// idpGroupFromAzureGroup converts a azure group to the identityprovider interface group type
func idpGroupFromAzureGroup(azureGroup AzureGroup) identity.IdpGroup {
	return identity.IdpGroup{
		ID:          azureGroup.ID,
		Name:        azureGroup.DisplayName,
		Description: string(azureGroup.Description),
	}
}
func (a *AzureSync) ListGroups(ctx context.Context) ([]identity.IdpGroup, error) {
	idpGroups := []identity.IdpGroup{}
	hasMore := true
	var nextToken *string
	url := MSGraphBaseURL + "/groups"
	for hasMore {

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Authorization", "Bearer "+a.token.Get())
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

		for _, u := range lu.Value {

			group := idpGroupFromAzureGroup(u)
			if err != nil {
				return nil, err
			}
			idpGroups = append(idpGroups, group)
		}
		nextToken = lu.OdataNextLink
		if nextToken != nil {
			url = *nextToken
		} else {
			hasMore = false
		}
	}
	return idpGroups, nil
}
