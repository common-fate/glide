package sso

import (
	"fmt"
	"os"
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
		googleSelected := idpType == "Google"
		oktaSelected := idpType == "Okta"
		googleConfigured := dc.Identity != nil && dc.Identity.Google != nil
		oktaConfigured := dc.Identity != nil && dc.Identity.Okta != nil
		isCurrentIDP := dc.Deployment.Parameters.IdentityProviderType == strings.ToUpper(idpType)
		update := true
		//if there are already params for that idp then ask if they want to update
		if dc.Identity != nil {
			if (googleSelected && googleConfigured) || (oktaSelected && oktaConfigured) {
				if isCurrentIDP {
					p3 := &survey.Confirm{Message: fmt.Sprintf("%s is currently set as your identity provider, do you want to update the configuration?", idpType)}
					err = survey.AskOne(p3, &update)
					if err != nil {
						return err
					}
					if !update {
						clio.Info("Closing SSO setup")
						return nil
					}
				} else {
					clio.Info("You already have configuration for %s but it's not currently set as your identity provider", idpType)
					p3 := &survey.Confirm{Message: fmt.Sprintf("Do you need to update the configuration for %s as well as setting it as your identity provider?", idpType)}
					var update bool
					err = survey.AskOne(p3, &update)
					if err != nil {
						return err
					}
				}
			}
		}
		if update {
			if googleSelected {
				docs := "https://docs.commonfate.io/granted-approvals/sso/google"
				clio.Info("You can follow along with the Google setup guide in our docs: %s", docs)
				var google deploy.Google
				if googleConfigured {
					google = *dc.Identity.Google
				}
				var token string
				p1 := &survey.Password{Message: "API Token:"}
				err = survey.AskOne(p1, &token)
				if err != nil {
					return err
				}
				p2 := &survey.Input{Message: "Google Workspace Domain:", Default: google.Domain}
				err = survey.AskOne(p2, &google.Domain)
				if err != nil {
					return err
				}
				p3 := &survey.Input{Message: "Google Admin Email", Default: google.AdminEmail}
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
				clio.Success("SSM Parameters set successfully")
				if dc.Identity == nil {
					dc.Identity = &deploy.IdentityConfig{
						Google: &google,
					}
				} else {
					dc.Identity.Google = &google
				}
			} else if oktaSelected {
				docs := "https://docs.commonfate.io/granted-approvals/sso/okta"
				clio.Info("You can follow along with the Okta setup guide in our docs: %s", docs)
				var okta deploy.Okta
				if oktaConfigured {
					okta = *dc.Identity.Okta
				}
				var token string
				p1 := &survey.Password{Message: "API Token:"}
				err = survey.AskOne(p1, &token)
				if err != nil {
					return err
				}
				p2 := &survey.Input{Message: "Okta Org URL:", Default: okta.OrgURL}
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
				clio.Success("SSM Parameters set successfully")
				if dc.Identity == nil {
					dc.Identity = &deploy.IdentityConfig{
						Okta: &okta,
					}
				} else {
					dc.Identity.Okta = &okta
				}
			}
			clio.Info("The following parameters are required to setup a SAML app in your identity provider")
			o, err := dc.LoadSAMLOutput(ctx)
			if err != nil {
				return err
			}
			o.PrintSAMLTable()
			var (
				fromUrl    = "URL"
				fromString = "String"
				fromFile   = "File"
			)
			p4 := &survey.Select{Message: "Would you like to use a metadata URL, an XML string, or load XML from a file?", Options: []string{fromUrl, fromString, fromFile}}
			metadataChoice := fromUrl
			err = survey.AskOne(p4, &metadataChoice)
			if err != nil {
				return err
			}
			switch metadataChoice {
			case fromUrl:
				p5 := &survey.Input{Message: "Metadata URL"}
				err = survey.AskOne(p5, &dc.Deployment.Parameters.SamlSSOMetadataURL)
				if err != nil {
					return err
				}
			case fromString:
				p5 := &survey.Input{Message: "Metadata XML String"}
				err = survey.AskOne(p5, &dc.Deployment.Parameters.SamlSSOMetadataURL)
				if err != nil {
					return err
				}
			case fromFile:
				p5 := &survey.Input{Message: "Metadata XML file"}
				var res string
				err := survey.AskOne(p5, &res, func(options *survey.AskOptions) error {
					options.Validators = append(options.Validators, func(ans interface{}) error {
						p := ans.(string)
						fileInfo, err := os.Stat(p)
						if err != nil {
							return err
						}
						if fileInfo.IsDir() {
							return fmt.Errorf("path is a directory, must be a file")
						}
						return nil
					})
					return nil
				})
				if err != nil {
					return err
				}
				b, err := os.ReadFile(res)
				if err != nil {
					return err
				}
				dc.Deployment.Parameters.SamlSSOMetadata = string(b)
			}
		}
		dc.Deployment.Parameters.IdentityProviderType = strings.ToUpper(idpType)
		clio.Info("Updating your deployment config")
		err = dc.Save(f)
		if err != nil {
			return err
		}
		clio.Success("Successfully completed SSO configuration")
		clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")
		return nil

	},
}
