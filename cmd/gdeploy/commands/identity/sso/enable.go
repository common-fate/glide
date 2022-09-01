package sso

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/identity/identitysync"

	"github.com/urfave/cli/v2"
)

var enableCommand = cli.Command{
	Name:        "enable",
	Description: "Set up SSO for a deployment",
	Usage:       "Configure SSO for a deployment",
	Action: func(c *cli.Context) error {
		ctx := c.Context

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

		idpType, _, err := registry.FromCLIOption(selected)
		if err != nil {
			return err
		}

		currentConfig := dc.Deployment.Parameters.IdentityConfiguration[idpType]

		// if there is already config for this idpType, use `sso update` cmd flow.
		if currentConfig != nil {
			update := true
			p3 := &survey.Confirm{Message: fmt.Sprintf("%s is currently set as your identity provider, do you want to update the configuration?", idpType)}
			err := survey.AskOne(p3, &update)
			if err != nil {
				return err
			}
			if !update {
				clio.Info("Closing SSO enable")
				return nil
			}
		}

		return updateOrAddSSO(c, idpType)
	},
}

func updateOrAddSSO(c *cli.Context, idpType string) error {
	ctx := c.Context

	dc, err := deploy.ConfigFromContext(ctx)
	if err != nil {
		return err
	}
	idp, ok := identitysync.Registry().IdentityProviders[idpType]
	if !ok {
		// Should never happen
		return errors.New("no matching identity provider found")
	}
	clio.Info("You can follow our %s setup guide at: https://docs.commonfate.io/granted-approvals/sso/%s for detailed instruction on setting up SSO", idpType, idp.DocsID)

	cfg := idp.IdentityProvider.Config()

	// if existing config then CLI prompt will have defaults loaded.
	currentConfig := dc.Deployment.Parameters.IdentityConfiguration[idpType]
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
	clio.Info("Instructions for setting up SAML SSO for %s can be found here: https://docs.commonfate.io/granted-approvals/sso/%s/#setting-up-saml-sso", idpType, idp.DocsID)
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
	clio.Warn("Don't forget to assign your users to the SAML app in %s so that they can login after setup is complete.", idpType)

	dc.Deployment.Parameters.IdentityProviderType = idpType
	clio.Info(`When using SSO, administrators for Granted are managed in your identity provider.
	Create a group called 'Granted Administrators' in your identity provider and copy the group's ID.
	Users in this group will be able to manage Access Rules.
	`)

	err = idp.IdentityProvider.Init(ctx)
	if err != nil {
		return err
	}

	grps, err := idp.IdentityProvider.ListGroups(ctx)
	if err != nil {
		return err
	}

	// convert groups to a string map
	groupMap := make(map[string]identity.IdpGroup)
	groupNames := []string{}
	chosenKey := ""
	for _, g := range grps {
		key := fmt.Sprintf("%s: %s", g.Name, g.Description)
		groupMap[key] = g
		groupNames = append(groupNames, key)
	}

	// sort groupNames alphabetically
	sort.Strings(groupNames)

	err = survey.AskOne(&survey.Select{
		Message: "The ID of the Granted Administrators group in your identity provider:",
		Options: groupNames,
	}, &chosenKey)

	if err != nil {
		return err
	}

	dc.Deployment.Parameters.AdministratorGroupID = groupMap[chosenKey].ID

	clio.Info("Updating your deployment config")

	f := c.Path("file")
	err = dc.Save(f)
	if err != nil {
		return err
	}
	clio.Success("Successfully completed SSO configuration")
	clio.Warn(`Users and groups will be synced every 5 minutes from your identity provider. To finish enabling SSO, follow these steps:

	  1) Run 'gdeploy update' to apply the changes to your CloudFormation deployment.
	  2) Run 'gdeploy identity sync' to trigger an immediate sync of your user directory.
	`)
	return nil
}
