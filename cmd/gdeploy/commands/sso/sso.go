package sso

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/deploy"

	"github.com/urfave/cli/v2"
)

// AvailableSSOProviders are the currently implemented SSO providers.
var AvailableSSOProviders = []string{
	"Google",
	"Okta",
}

var SSOCommand = cli.Command{
	Name:        "sso",
	Subcommands: []*cli.Command{&configureCommand},
	Action:      cli.ShowSubcommandHelp,
}

var configureCommand = cli.Command{
	Name:        "configure",
	Flags:       []cli.Flag{&cli.BoolFlag{Name: "overwrite", Aliases: []string{"o"}, Usage: "if provided, will prompt to override parameter values"}},
	Description: "Set up SSO for a deployment",
	Action: func(c *cli.Context) error {

		overwrite := c.Bool("overwrite")
		ctx := c.Context

		f := c.Path("file")

		dc, err := deploy.LoadConfig(f)
		if err != nil {
			return err
		}

		var ssoEnable string
		p2 := &survey.Select{Message: "The SSO provider to deploy with", Options: AvailableSSOProviders}
		err = survey.AskOne(p2, &ssoEnable)
		if err != nil {
			return err
		}

		//if there are already params for that idp then ask if they want to update
		if dc.Identity != nil {
			if (dc.Identity.Google != nil && ssoEnable == "Google") ||
				(dc.Identity.Okta != nil && ssoEnable == "Okta") {
				clio.Info("You already have params set for %s", ssoEnable)
				p3 := &survey.Confirm{Message: "Would you like to update the current parameters?"}
				err = survey.AskOne(p3, &overwrite)
				if err != nil {
					return err
				}

				//if both google and okta have config update the deployment param here and exit
				if dc.Identity.Google != nil && dc.Identity.Okta != nil && !overwrite {
					switch ssoEnable {
					case "Google":
						dc.Deployment.Parameters.IdentityProviderType = "GOOGLE"
					case "Okta":
						dc.Deployment.Parameters.IdentityProviderType = "OKTA"
					}
					err = dc.Save(f)
					if err != nil {
						return err
					}
					clio.Info("Successfully updated SSO configuration")
					clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")
					return nil

				}

				//if they dont want to update current param then exit
				if !overwrite {

					clio.Info("Closing SSO setup")
					return nil
				}
			}
		}

		switch ssoEnable {
		case "Google":
			docs := "https://docs.commonfate.io/granted-approvals/sso/google"
			clio.Info("Find documentation for setting up Google Workspace in our setup docs: %s", docs)

			var google deploy.Google
			if dc.Identity != nil && dc.Identity.Google != nil {
				google = *dc.Identity.Google
			}
			var token string

			p1 := &survey.Password{Message: "API Token:"}
			err = survey.AskOne(p1, &token)
			if err != nil {
				return err
			}

			p2 := &survey.Input{Message: "Google Workspace Domain:"}
			err = survey.AskOne(p2, &google.Domain)
			if err != nil {
				return err
			}

			p3 := &survey.Input{Message: "Google Admin Email"}
			err = survey.AskOne(p3, &google.AdminEmail)
			if err != nil {
				return err
			}
			path, version, err := config.PutSecretVersion(ctx, config.GoogleTokenPath, dc.Deployment.Parameters.DeploymentSuffix, token)
			if err != nil {
				return err
			}
			google.APIToken = config.AWSSSMParamToken(path, version)
			if err != nil {
				return err
			}
			clio.Success("SSM Parameters Set Successfully\n")

			clio.Warn("SAML outputs:\n")
			o, err := dc.LoadSAMLOutput(ctx)

			if err != nil {
				return err
			}
			o.PrintSAMLTable()
			if dc.Identity == nil {
				dc.Identity = &deploy.Identity{
					Google: &google,
				}
			} else {
				dc.Identity.Google = &google
			}
			clio.Info("Find documentation for setting up SAML SSO here: %s", docs)

			dc.Deployment.Parameters.IdentityProviderType = "GSUITE"
			//complete the setup with the saml metadata
			var metadata string
			p4 := &survey.Input{Message: "SAML Metadata String:"}
			err = survey.AskOne(p4, &metadata)
			if err != nil {
				return err
			}
			dc.Deployment.Parameters.SamlSSOMetadata = metadata

		case "Okta":
			docs := "https://docs.commonfate.io/granted-approvals/sso/okta"
			clio.Info("Find documentation for setting up Okta in our setup docs: %s", docs)

			var okta deploy.Okta
			if dc.Identity != nil && dc.Identity.Google != nil {
				okta = *dc.Identity.Okta
			}
			var token string
			p1 := &survey.Password{Message: "API Token:"}
			err = survey.AskOne(p1, &token)
			if err != nil {
				return err
			}

			p2 := &survey.Input{Message: "Okta Org URL:"}
			err = survey.AskOne(p2, &okta.OrgURL)
			if err != nil {
				return err
			}
			path, version, err := config.PutSecretVersion(ctx, config.OktaTokenPath, dc.Deployment.Parameters.DeploymentSuffix, token)
			if err != nil {
				return err
			}
			okta.APIToken = config.AWSSSMParamToken(path, version)
			if err != nil {
				return err
			}
			clio.Success("SSM Parameters Set Successfully\n")
			if dc.Identity == nil {
				dc.Identity = &deploy.Identity{
					Okta: &okta,
				}
			} else {
				dc.Identity.Okta = &okta
			}

			clio.Warn("SAML outputs:\n")
			o, err := dc.LoadSAMLOutput(ctx)

			if err != nil {
				return err
			}
			o.PrintSAMLTable()

			clio.Info("Find documentation for setting up SAML SSO here: %s", docs)

			dc.Deployment.Parameters.IdentityProviderType = "OKTA"

			//complete the setup with the saml metadata
			var metadata string
			t := &survey.Input{Message: "Okta SAML metadata string:"}
			err = survey.AskOne(t, &metadata)
			if err != nil {
				return err
			}
			dc.Deployment.Parameters.SamlSSOMetadataURL = metadata

		}

		err = dc.Save(f)
		if err != nil {
			return err
		}
		clio.Success("completed SSO setup")
		clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")
		return nil

	},
}
