package sso

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"

	"github.com/urfave/cli/v2"
)

var SSOCommand = cli.Command{
	Name:        "sso",
	Subcommands: []*cli.Command{&configureCommand, &testCommand},
	Action:      cli.ShowSubcommandHelp,
}

var testCommand = cli.Command{
	Name:  "test",
	Usage: "Test SSO configuration",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		f := c.Path("file")
		dc, err := deploy.ConfigFromContext(ctx)

		// Prompt: Do you want to enter an group id, or select frokm a list of groups?
		surveyQ := []string{"Enter group id", "Select from a list of groups"}
		var surveryA string
		err := survey.AskOne(&survey.Select{
			Message: "Do you want to enter an group id, or select from a list of groups?",
			Options: surveyQ,
		}, &surveryA)
		if err != nil {
			return err
		}

		var groupID string

		if surveryA == surveyQ[0] {
			// Prompt: Enter group id
			err := survey.AskOne(&survey.Input{
				Message: "Enter group id",
			}, &groupID)
			if err != nil {
				return err
			}
		} else {
			// Prompt: Select from a list of groups
			var groupIDs []string

			switch dc.Deployment.Parameters.IdentityProviderType {
			case "okta":

			}

			// Fetch group IDs

			// instantiate an

			// @TODO: this will have to be conditionally implemented based on the provider
			// identitysync.IdentityProvider.ListGroups()

			err := survey.AskOne(&survey.MultiSelect{
				Message: "Select from a list of groups",
				Options: groupIDs,
			}, &groupID)

			if err != nil {
				return err
			}
		}

		// f := c.Path("file")
		// dc, err := deploy.ConfigFromContext(ctx)
		// if err != nil {
		// 	return err
		// }
		clio.Info(" tesssst :) ")
		// registry := identitysync.Registry()
		return nil
	},
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
		clio.Info("Follow the documentation for setting up SSO here: https://docs.commonfate.io/granted-approvals/sso/overview")
		registry := identitysync.Registry()

		var selected string
		p2 := &survey.Select{Message: "The SSO provider to deploy with", Options: registry.CLIOptions()} //Default: i
		err = survey.AskOne(p2, &selected)
		if err != nil {
			return err
		}
		idpType, idp, err := registry.FromCLIOption(selected)
		if err != nil {
			return err
		}
		currentConfig := dc.Deployment.Parameters.IdentityConfiguration[idpType]
		update := true

		TEMP_SKIP := false
		//if there are already params for that idp then ask if they want to update
		if currentConfig != nil && TEMP_SKIP {
			if idpType == dc.Deployment.Parameters.IdentityProviderType {
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
				p3 := &survey.Confirm{Message: fmt.Sprintf("Do you need to update the configuration for %s before setting it as your identity provider?", idpType)}
				var update bool
				err = survey.AskOne(p3, &update)
				if err != nil {
					return err
				}
			}
		}

		cfg := idp.IdentityProvider.Config()
		if update && TEMP_SKIP {
			clio.Info("Don't know where to find an SSO credential? Best place to find out would be our docs!")
			clio.Info("Follow our %s setup guide at: https://docs.commonfate.io/granted-approvals/sso/%s", idpType, idpType)
			// if there is existing config, process it into the idp struct
			// This way, the cli prompt will have defaults loaded
			if currentConfig != nil {
				err := cfg.Load(ctx, &gconfig.MapLoader{Values: currentConfig})
				if err != nil {
					return err
				}
			}

			for _, v := range cfg {
				err := deploy.CLIPrompt(v)
				if err != nil {
					return err
				}
			}

			err = deploy.RunConfigTest(ctx, idp.IdentityProvider)
			if err != nil {
				return err
			}
			// if tests pass, dump the config and update in the deployment config
			newConfig, err := cfg.Dump(ctx, gconfig.SSMDumper{Suffix: dc.Deployment.Parameters.DeploymentSuffix})
			if err != nil {
				return err
			}
			dc.Deployment.Parameters.IdentityConfiguration.Upsert(idpType, newConfig)

			clio.Info("The following parameters are required to setup a SAML app in your identity provider")
			clio.Info("Instructions for setting up SAML SSO for %s can be found here: https://docs.commonfate.io/granted-approvals/sso/%s/#setting-up-saml-sso", idpType, idpType)
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

			updateMetadata := true
			if dc.Deployment.Parameters.SamlSSOMetadata != "" || dc.Deployment.Parameters.SamlSSOMetadataURL != "" {

				p5 := &survey.Confirm{Message: "You already have a metadata string/url set, would you like to update it?"}
				err = survey.AskOne(p5, &updateMetadata)
				if err != nil {
					return err
				}
			}
			if updateMetadata {
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
		}

		dc.Deployment.Parameters.IdentityProviderType = idpType
		clio.Warn("Don't forget to assign your users to the SAML app in %s so that they can login after setup is complete.", idpType)
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
