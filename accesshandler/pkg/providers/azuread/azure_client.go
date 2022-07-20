package azuread

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/common-fate/granted-approvals/pkg/deploy"
)

//making our own azure client to interact with in access handler
type AzureClient struct {
	NewClient *http.Client
	token     string
}

type ClientSecretCredential struct {
	client confidential.Client
}

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

func (c *AzureClient) ListGroups(context.Context) ([]AzureGroup, error) {

	return nil, nil
}

func (c *AzureClient) GetGroup(context.Context, string) (*AzureGroup, error) {

	return nil, nil
}

func (c *AzureClient) GetUser(context.Context, string) (*AzureUser, error) {

	return nil, nil
}

func (c *AzureClient) AddUserToGroup(context.Context, string, string) (*AzureUser, error) {

	return nil, nil
}

func (c *AzureClient) RemoveUserFromGroup(context.Context, string, string) (*AzureUser, error) {

	return nil, nil
}

func (c *AzureClient) ListGroupUsers(context.Context, string) ([]AzureUser, error) {

	return nil, nil
}

// NewAzure will fail if the Azure settings are not configured
func NewAzure(ctx context.Context, settings deploy.Azure) (*AzureClient, error) {
	azAuth, err := NewClientSecretCredential(settings, http.DefaultClient)
	if err != nil {
		return nil, err
	}
	token, err := azAuth.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	return &AzureClient{NewClient: http.DefaultClient, token: token}, nil
}

func NewClientSecretCredential(s deploy.Azure, httpClient *http.Client) (*ClientSecretCredential, error) {
	cred, err := confidential.NewCredFromSecret(s.ClientSecret)
	if err != nil {
		return nil, err
	}
	c, err := confidential.New(s.ClientID, cred,
		confidential.WithAuthority(fmt.Sprintf("%s/%s", ADAuthorityHost, s.TenantID)),
		confidential.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return &ClientSecretCredential{client: c}, nil
}

// GetToken requests an access token from Azure Active Directory. This method is called automatically by Azure SDK clients.
func (c *ClientSecretCredential) GetToken(ctx context.Context) (string, error) {
	ar, err := c.client.AcquireTokenByCredential(ctx, []string{"https://graph.microsoft.com/.default"})
	return ar.AccessToken, err
}
