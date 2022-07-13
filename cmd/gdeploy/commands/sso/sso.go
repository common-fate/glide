package sso

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/deploy"

	"github.com/urfave/cli/v2"
)

// idpTypes are the currently implemented SSO providers.
var idpTypes = []string{
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
	Description: "Set up SSO for a deployment",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		f := c.Path("file")
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		var idpType string
		p2 := &survey.Select{Message: "The SSO provider to deploy with", Options: idpTypes}
		err = survey.AskOne(p2, &idpType)
		if err != nil {
			return err
		}

		googleConfigured := dc.Identity.Google != nil
		googleSelected := idpType == "Google"
		oktaConfigured := dc.Identity.Okta != nil
		oktaSelected := idpType == "Okta"
		isCurrentIDP := dc.Deployment.Parameters.IdentityProviderType == strings.ToUpper(idpType)

		//if there are already params for that idp then ask if they want to update
		if dc.Identity != nil {
			if (googleSelected && googleConfigured) || (oktaSelected && oktaConfigured) {
				if isCurrentIDP {
					p3 := &survey.Confirm{Message: fmt.Sprintf("%s is currently set as your identity provider, do you want to update the configuration?",idpType)}
					var update bool
					err = survey.AskOne(p3, &update)
					if err != nil {
						return err
					}
					if !update {
						clio.Info("Closing SSO setup")
						return nil
					}
				}else {
					clio.Info("You already have configuration for %s but it's not currently set as your identity provider",idpType)
					p3 := &survey.Confirm{Message: "Do you need to update the configuration for %s?")}
					err = survey.AskOne(p3, &overwrite)
					if err != nil {
						return err
					}
				}

				if !overwrite {
					p3 := &survey.Confirm{Message: "This process will overwrite your existing configuration for %s, are you sure?"}
					err = survey.AskOne(p3, &overwrite)
					if err != nil {
						return err
					}
					//if they dont want to update current param then exit
					if !overwrite {
						clio.Info("Closing SSO setup")
						return nil
					}
				}
			}
		}

		if googleSelected {
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
				dc.Identity = &deploy.IdentityConfig{
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
		} else if oktaSelected {
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
				dc.Identity = &deploy.IdentityConfig{
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
			t := &survey.Input{Message: "Okta SAML metadata URL:"}
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
