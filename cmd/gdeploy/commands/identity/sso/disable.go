package sso

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"

	"github.com/urfave/cli/v2"
)

var disableCommand = cli.Command{
	Name:        "disable",
	Description: "Clear current sso configuration. Defaults back to cognito.",
	Usage:       "Clear current sso configuration. Defaults back to cognito.",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		f := c.Path("file")
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		if dc.Deployment.Parameters.IdentityProviderType == "" || dc.Deployment.Parameters.IdentityProviderType == identitysync.IDPTypeCognito {
			clio.Info("You don't currently have SSO configured so this command will not make any changes.")
			return nil
		}

		var confirm bool
		prompt := &survey.Confirm{
			Message: "Are you sure you want to disable SSO?",
		}

		err = survey.AskOne(prompt, &confirm)
		if err != nil {
			return err
		}
		if !confirm {
			clio.Info("Cancelled disabling SSO")
			return nil
		}
		if err := dc.ResetIdentityProviderToCognito(f); err != nil {
			return err
		}

		clio.Success("Successfully disabled SSO")
		clio.Warn(`SSO has been disabled and your deployment will now use the default Cognito user pool for logins. To finish disabling SSO, follow these steps:

		1) Run 'gdeploy update' to apply the changes to your CloudFormation deployment.
		2) Run 'gdeploy identity sync' to trigger an immediate sync of your cognito user pool.
	  `)

		return nil
	},
}
