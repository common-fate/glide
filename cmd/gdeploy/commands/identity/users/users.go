package users

import (
	"fmt"

	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"
	"github.com/urfave/cli/v2"
)

var UsersCommand = cli.Command{
	Name:        "users",
	Description: "Add or remove users from the default cognito user pool.\nThese commands are only available when you are using the default Cognito user pool. If you have connected an SSO provider, like Okta, Google or AzureAD, use those tools to manage your users and groups instead.",
	Subcommands: []*cli.Command{&CreateCommand, &DeleteCommand},
	Action:      cli.ShowSubcommandHelp,
	Before: func(c *cli.Context) error {
		args := c.Args().Slice()
		// if help argument is provided then skip this check.
		for _, arg := range args {
			if arg == "-h" || arg == "--help" || arg == "help" {
				return nil
			}
		}
		dc, err := deploy.ConfigFromContext(c.Context)
		if err != nil {
			return err
		}
		idpType := dc.Deployment.Parameters.IdentityProviderType
		if idpType != "" && idpType != identitysync.IDPTypeCognito {
			return clio.NewCLIError(
				fmt.Sprintf("This command is only available when you are using the default Cognito identity provider, it looks like you are using %s", idpType),
				clio.InfoMsg("If you would like to add or remove a user, manage them in your identity provider, then wait 5 minutes or run 'gdeploy identity sync' to sync the changes immediately"),
				clio.InfoMsg("If you would like to make a user an administrator of Granted Approvals, add them to the %s group in your identity provider.", dc.Deployment.Parameters.AdministratorGroupID),
			)
		}
		return nil
	},
}
