package sso

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/identity/identitysync"

	"github.com/urfave/cli/v2"
)

var updateCommand = cli.Command{
	Name:        "update",
	Description: "Update SSO configuration",
	Usage:       "Updage SSO configuration ",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		currentIdpType := dc.Deployment.Parameters.IdentityProviderType
		// If SSO is not configured, push the user into the `gconfig identity sso enable` flow
		if currentIdpType == "" || currentIdpType == identitysync.IDPTypeCognito {
			clio.Info("You are currently using cognito as your identity provider. If you were trying to sync users and groups. run `gdeploy identity sync` otherwise, you can setup SSO now.")
			return enableCommand.Run(c)
		}
		clio.Infof("Updating configuration for %s", currentIdpType)
		return updateOrAddSSO(c, currentIdpType)
	},
}
