package identitysync

import (
	"context"
	"net/http"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

const MSGraphBaseURL = "https://graph.microsoft.com/v1.0"

type AzureSync struct {
	// Client    *msgraphsdk.GraphServiceClient
	// Adapter   *msgraphsdk.GraphRequestAdapter
	NewClient *http.Client
	token     string
}
type ClientSecretCredential struct {
	client confidential.Client
}

// GetToken requests an access token from Azure Active Directory. This method is called automatically by Azure SDK clients.
func (c *ClientSecretCredential) GetToken(ctx context.Context) (string, error) {
	ar, err := c.client.AcquireTokenByCredential(ctx, []string{"https://graph.microsoft.com/.default"})
	return ar.AccessToken, err
}

func NewClientSecretCredential(s deploy.Azure) (*ClientSecretCredential, error) {
	cred, err := confidential.NewCredFromSecret(s.ClientSecret)
	if err != nil {
		return nil, err
	}
	c, err := confidential.New(s.ClientID, cred)
	if err != nil {
		return nil, err
	}
	return &ClientSecretCredential{client: c}, nil
}

// NewAzure will fail if the Azure settings are not configured
func NewAzure(ctx context.Context, settings deploy.Azure) (*AzureSync, error) {
	azAuth, err := NewClientSecretCredential(settings)
	if err != nil {
		return nil, err
	}
	token, err := azAuth.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	return &AzureSync{NewClient: &http.Client{}, token: token}, nil
}

// idpUserFromAzureUser converts a azure user to the identityprovider interface user type
func (a *AzureSync) idpUserFromAzureUser(ctx context.Context, azureUser models.Userable) (identity.IdpUser, error) {
	u := identity.IdpUser{
		ID:        aws.ToString(azureUser.GetId()),
		FirstName: aws.ToString(azureUser.GetGivenName()),
		LastName:  aws.ToString(azureUser.GetSurname()),
		Email:     aws.ToString(azureUser.GetMail()),
		Groups:    []string{},
	}

	// result, err := a.Client.UsersById(u.ID).MemberOf().Get()
	// if err != nil {
	// 	return identity.IdpUser{}, err
	// }

	// // Use PageIterator to iterate through all groups
	// pageIterator, err := msgraphcore.NewPageIterator(result, a.Adapter, models.CreateDirectoryFromDiscriminatorValue)
	// if err != nil {
	// 	return identity.IdpUser{}, err
	// }
	// err = pageIterator.Iterate(func(pageItem interface{}) bool {
	// 	if _, ok := pageItem.(models.Groupable); ok {
	// 		graphGroup := pageItem.(models.Groupable)

	// 		u.Groups = append(u.Groups, aws.ToString(graphGroup.GetId()))

	// 		return true
	// 	}
	// 	return false
	// })
	// if err != nil {
	// 	return identity.IdpUser{}, err
	// }

	return u, nil
}

func (a *AzureSync) ListUsers(ctx context.Context) ([]identity.IdpUser, error) {
	// //get all users
	// idpUsers := []identity.IdpUser{}
	// result, err := a.Client.Users().Get()
	// if err != nil {
	// 	return nil, err
	// }

	// // Use PageIterator to iterate through all users
	// pageIterator, err := msgraphcore.NewPageIterator(result, a.Adapter, models.CreateUserCollectionResponseFromDiscriminatorValue)
	// if err != nil {
	// 	return nil, err
	// }
	// err = pageIterator.Iterate(func(pageItem interface{}) bool {
	// 	graphUser := pageItem.(models.Userable)

	// 	user, err := a.idpUserFromAzureUser(ctx, graphUser)
	// 	if err != nil {
	// 		return false
	// 	}
	// 	idpUsers = append(idpUsers, user)

	// 	return true
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// return idpUsers, nil

	req, _ := http.NewRequest("GET", MSGraphBaseURL+"/users", nil)
	req.Header.Add("Authorization", "Bearer "+a.token)
	res, _ := a.NewClient.Do(req)
	_ = res
	return nil, nil
}

// idpGroupFromAzureGroup converts a azure group to the identityprovider interface group type
func idpGroupFromAzureGroup(azureGroup models.Groupable) identity.IdpGroup {
	return identity.IdpGroup{
		ID:          aws.ToString(azureGroup.GetId()),
		Name:        aws.ToString(azureGroup.GetDisplayName()),
		Description: aws.ToString(azureGroup.GetDescription()),
	}
}
func (a *AzureSync) ListGroups(ctx context.Context) ([]identity.IdpGroup, error) {
	// idpGroups := []identity.IdpGroup{}
	// result, err := a.Client.Groups().Get()
	// if err != nil {
	// 	return nil, err
	// }

	// // Use PageIterator to iterate through all users
	// pageIterator, err := msgraphcore.NewPageIterator(result, a.Adapter, models.CreateGroupCollectionResponseFromDiscriminatorValue)
	// if err != nil {
	// 	return nil, err
	// }
	// err = pageIterator.Iterate(func(pageItem interface{}) bool {
	// 	graphGroup := pageItem.(models.Groupable)

	// 	user := idpGroupFromAzureGroup(graphGroup)

	// 	idpGroups = append(idpGroups, user)

	// 	return true
	// })
	// if err != nil {
	// 	return nil, err
	// }
	// return idpGroups, nil
	return nil, nil
}
