package sso

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers/azure/ad"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/common-fate/granted-approvals/pkg/deploy"

	"github.com/urfave/cli/v2"
)

// idpTypes are the currently implemented SSO providers.
var idpTypes = []string{
	"Google",
	"Okta",
	"Azure",
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
		azureSelected := idpType == "Azure"
		googleConfigured := dc.Identity != nil && dc.Identity.Google != nil
		oktaConfigured := dc.Identity != nil && dc.Identity.Okta != nil
		azureConfigured := dc.Identity != nil && dc.Identity.Azure != nil
		isCurrentIDP := dc.Deployment.Parameters.IdentityProviderType == strings.ToUpper(idpType)
		update := true
		//if there are already params for that idp then ask if they want to update
		if dc.Identity != nil {
			if (googleSelected && googleConfigured) || (oktaSelected && oktaConfigured) || (azureSelected && azureConfigured) {
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
			} else if azureSelected {
				docs := "https://docs.commonfate.io/granted-approvals/sso/azure"
				clio.Info("You can follow along with the Azure setup guide in our docs: %s", docs)
				var azure deploy.Azure
				if azureConfigured {
					azure = *dc.Identity.Azure
				}
				p1 := &survey.Input{Message: "Tenant ID:", Default: azure.TenantID}
				err = survey.AskOne(p1, &azure.TenantID)
				if err != nil {
					return err
				}

				p2 := &survey.Input{Message: "Client ID:", Default: azure.ClientID}
				err = survey.AskOne(p2, &azure.ClientID)
				if err != nil {
					return err
				}

				var token string
				p3 := &survey.Password{Message: "Client Secret:"}
				err = survey.AskOne(p3, &token)
				if err != nil {
					return err
				}
				//Try using the api to call using the creds to test if they are setup correctly.
				testprovider := ad.Provider{
					TenantID:     azure.TenantID,
					ClientID:     azure.ClientID,
					ClientSecret: token,
				}

				path, version, err := config.PutSecretVersion(ctx, config.AzureSecretPath, dc.Deployment.Parameters.DeploymentSuffix, token)
				if err != nil {
					return err
				}
				azure.ClientSecret = config.AWSSSMParamToken(path, version)
				if err != nil {
					return err
				}

				err = testprovider.Init(ctx)
				if err != nil {

					return err
				}
				clio.Success("Verifying credentials...")

				groups, err := testprovider.Client.ListGroups(ctx)
				if err != nil {
					return fmt.Errorf("Something went wrong calling Azure with provided credentials: %s", err)
				}
				clio.Success("List groups works")

				_, err = testprovider.Client.ListUsers(ctx)
				if err != nil {
					return fmt.Errorf("Something went wrong calling Azure with provided credentials: %s", err)
				}
				clio.Success("List users works")

				_, err = testprovider.Client.ListGroupUsers(ctx, groups[0].ID)
				if err != nil {
					return fmt.Errorf("Something went wrong calling Azure with provided credentials: %s", err)
				}
				clio.Success("List group members works")

				clio.Success("SSM Parameters set successfully")

				if dc.Identity == nil {
					dc.Identity = &deploy.IdentityConfig{
						Azure: &azure,
					}
				} else {
					dc.Identity.Azure = &azure
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

		clio.Info(`When using SSO, administrators for Granted are managed in your identity provider.
Create a group called 'Granted Administrators' in your identity provider and copy the group's ID.
Users in this group will be able to manage Access Rules.
`)
		adminGroupPrompt := &survey.Input{Message: "The ID of the Granted Administrators group in your identity provider:", Default: dc.Deployment.Parameters.AdministratorGroupID}
		err = survey.AskOne(adminGroupPrompt, &dc.Deployment.Parameters.AdministratorGroupID)
		if err != nil {
			return err
		}

		clio.Info("Updating your deployment config")
		err = dc.Save(f)
		if err != nil {
			return err
		}
		clio.Success("Successfully completed SSO configuration")
		clio.Warn(`Users and will be synced every 5 minutes from your identity provider. To finish enabling SSO, follow these steps:

  1) Run 'gdeploy update' to apply the changes to your CloudFormation deployment.
  2) Run 'gdeploy users sync' to trigger an immediate sync of your user directory.
`)
		return nil
	},
}
