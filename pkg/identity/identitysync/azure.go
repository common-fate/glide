package identitysync

import (
	"context"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/identity"
	a "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

type AzureSync struct {
	Client  *msgraphsdk.GraphServiceClient
	Adapter *msgraphsdk.GraphRequestAdapter
}

// NewAzure will fail if the Azure settings are not configured
func NewAzure(ctx context.Context, settings deploy.Azure) (*AzureSync, error) {

	//TODO: For applications calling graph api, 'On-behalf-of' provider is not yet available for the go graph sdk
	// for the time being will be using a client credentials provider which uses application permissions
	cred, err := azidentity.NewClientSecretCredential(settings.TenantID, settings.ClientID, settings.ClientSecret, nil)
	if err != nil {
		return nil, err
	}

	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		return nil, err
	}
	adapter, err := msgraphsdk.NewGraphRequestAdapter(auth)
	if err != nil {
		return nil, err
	}
	client := msgraphsdk.NewGraphServiceClient(adapter)

	return &AzureSync{Client: client, Adapter: adapter}, nil
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

	result, err := a.Client.UsersById(u.ID).MemberOf().Get()
	if err != nil {
		return identity.IdpUser{}, err
	}

	// Use PageIterator to iterate through all groups
	pageIterator, err := msgraphcore.NewPageIterator(result, a.Adapter, models.CreateDirectoryFromDiscriminatorValue)
	if err != nil {
		return identity.IdpUser{}, err
	}
	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		if _, ok := pageItem.(models.Groupable); ok {
			graphGroup := pageItem.(models.Groupable)

			u.Groups = append(u.Groups, aws.ToString(graphGroup.GetId()))

			return true
		}
		return false
	})
	if err != nil {
		return identity.IdpUser{}, err
	}

	return u, nil
}

func (a *AzureSync) ListUsers(ctx context.Context) ([]identity.IdpUser, error) {
	//get all users
	idpUsers := []identity.IdpUser{}
	result, err := a.Client.Users().Get()
	if err != nil {
		return nil, err
	}

	// Use PageIterator to iterate through all users
	pageIterator, err := msgraphcore.NewPageIterator(result, a.Adapter, models.CreateUserCollectionResponseFromDiscriminatorValue)
	if err != nil {
		return nil, err
	}
	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		graphUser := pageItem.(models.Userable)

		user, err := a.idpUserFromAzureUser(ctx, graphUser)
		if err != nil {
			return false
		}
		idpUsers = append(idpUsers, user)

		return true
	})
	if err != nil {
		return nil, err
	}
	return idpUsers, nil
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
	idpGroups := []identity.IdpGroup{}
	result, err := a.Client.Groups().Get()
	if err != nil {
		return nil, err
	}

	// Use PageIterator to iterate through all users
	pageIterator, err := msgraphcore.NewPageIterator(result, a.Adapter, models.CreateGroupCollectionResponseFromDiscriminatorValue)
	if err != nil {
		return nil, err
	}
	err = pageIterator.Iterate(func(pageItem interface{}) bool {
		graphGroup := pageItem.(models.Groupable)

		user := idpGroupFromAzureGroup(graphGroup)

		idpGroups = append(idpGroups, user)

		return true
	})
	if err != nil {
		return nil, err
	}
	return idpGroups, nil
}
