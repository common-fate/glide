package users

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"
	"github.com/urfave/cli/v2"
)

var DeleteCommand = cli.Command{
	Name: "delete",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "username", Aliases: []string{"u"}, Usage: "The username of the user to delete (should be an email address)", Required: true},
	},
	Description: "Delete a Cognito user",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		username := c.String("username")
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}
		adminGroup := dc.Deployment.Parameters.AdministratorGroupID

		idpType := dc.Deployment.Parameters.IdentityProviderType
		if idpType != "" && idpType != identitysync.IDPTypeCognito {
			return clio.NewCLIError(fmt.Sprintf("Your Granted Approvals deployment uses the %s identity provider. Remove users inside of your identity provider rather than using this CLI.\n\nIf you would like to make a user an administrator of Granted Approvals, add them to the %s group in your identity provider.", idpType, adminGroup))
		}

		o, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}

		cog := cognitoidentityprovider.NewFromConfig(cfg)

		_, err = cog.AdminRemoveUserFromGroup(ctx, &cognitoidentityprovider.AdminRemoveUserFromGroupInput{
			GroupName:  &adminGroup,
			UserPoolId: &o.UserPoolID,
			Username:   &username,
		})
		if err != nil {
			var unfe *types.UserNotFoundException

			if errors.As(err, &unfe) {
				return clio.NewCLIError(fmt.Sprintf("Failed to delete %s with error: %s", username, unfe.ErrorMessage()))
			}
			return err
		}

		_, err = cog.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: &o.UserPoolID,
			Username:   &username,
		})

		if err != nil {
			return err
		}

		clio.Success("Deleted user %s", username)

		return nil
	}}
