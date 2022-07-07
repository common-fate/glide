package identitysync

import (
	"context"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
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
	cred, err := azidentity.NewClientSecretCredential(settings.TenantID, settings.ClientID, settings.ClientSecret, &azidentity.ClientSecretCredentialOptions{})
	if err != nil {
		return nil, err
	}

	auth, err := a.NewAzureIdentityAuthenticationProviderWithScopes(cred, []string{"User.Read"})
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

// func (a *AzureSync) GetUserByEmail(ctx context.Context, email string) (user *identity.IdpUser, usergroups []identity.IdpGroup, err error) {
// 	oktaUser, _, err := o.Client.User.GetUser(ctx, email)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	user = &identity.IdpUser{
// 		ID:        oktaUser.Id,
// 		FirstName: (*oktaUser.Profile)["firstName"].(string),
// 		LastName:  (*oktaUser.Profile)["lastName"].(string),
// 		Email:     (*oktaUser.Profile)["email"].(string),
// 		Groups:    []string{},
// 	}

// 	userGroups, _, err := a.Client.User.ListUserGroups(ctx, oktaUser.Id)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	for _, g := range userGroups {
// 		user.Groups = append(user.Groups, g.Id)
// 		usergroups = append(usergroups, idpGroupFromOktaGroup(g))
// 	}
// 	return user, usergroups, nil
// }

// idpUserFromAzureUser converts a azure user to the identityprovider interface user type
func (a *AzureSync) idpUserFromAzureUser(ctx context.Context, azureUser models.Userable) (identity.IdpUser, error) {
	u := identity.IdpUser{
		ID:        *azureUser.GetEmployeeId(),
		FirstName: *azureUser.GetGivenName(),
		// LastName:  *azureUser.GetGivenName(),
		// Email:     *azureUser.getemai(),
		Groups: []string{},
	}

	// userGroups, _, err := azureUser.group
	// if err != nil {
	// 	return u, err
	// }
	// for _, g := range userGroups {
	// 	u.Groups = append(u.Groups, g.Id)
	// }

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
func idpGroupFromAzureGroup(oktaGroup models.Groupable) identity.IdpGroup {
	return identity.IdpGroup{
		ID:          *oktaGroup.GetId(),
		Name:        *oktaGroup.GetDisplayName(),
		Description: *oktaGroup.GetDescription(),
	}
}
func (a *AzureSync) ListGroups(ctx context.Context) ([]identity.IdpGroup, error) {
	idpGroups := []identity.IdpGroup{}
	result, err := a.Client.Groups().Get()
	if err != nil {
		return nil, err
	}

	// Use PageIterator to iterate through all users
	pageIterator, err := msgraphcore.NewPageIterator(result, a.Adapter, models.CreateUserCollectionResponseFromDiscriminatorValue)
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
