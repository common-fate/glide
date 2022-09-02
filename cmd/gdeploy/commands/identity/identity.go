package identity

import (
	"fmt"

	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/groups"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/sso"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/sync"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/commands/identity/users"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/middleware"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "identity",
	Description: "Identity commands are used to manage how your users login to Granted Approvals.\nYou can manage users and groups in the default Cognito user pool or configure your corporate SSO provider.",
	Usage:       "Configure how your users login to Granted Approvals",
	Action:      cli.ShowSubcommandHelp,
	Subcommands: []*cli.Command{
		&sso.SSOCommand,
		&sync.SyncCommand,
		middleware.WithBeforeFuncs(&users.UsersCommand, PreventNonCognitoUsage()),
		middleware.WithBeforeFuncs(&groups.GroupsCommand, PreventNonCognitoUsage()),
	},
}

func PreventNonCognitoUsage() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		dc, err := deploy.ConfigFromContext(c.Context)
		if err != nil {
			return err
		}
		idpType := dc.Deployment.Parameters.IdentityProviderType
		if idpType != "" && idpType != identitysync.IDPTypeCognito {
			return clio.NewCLIError(
				fmt.Sprintf("This command is only available when you are using the default Cognito identity provider, it looks like you are using %s", idpType),
				clio.InfoMsg("If you would like to add or remove a user or group, manage them in your identity provider, then wait 5 minutes or run 'gdeploy identity sync' to sync the changes immediately"),
				clio.InfoMsg("If you would like to make a user an administrator of Granted Approvals, add them to the %s group in your identity provider.", dc.Deployment.Parameters.AdministratorGroupID),
			)
		}
		return nil
	}
}
