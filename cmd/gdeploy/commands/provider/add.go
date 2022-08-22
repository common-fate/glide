package provider

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providerregistry"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/urfave/cli/v2"
)

var addCommand = cli.Command{
	Name:        "add",
	Description: "Add an access provider",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "overwrite", Usage: "Force SSM parameters to be overwritten if they exist"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		f := c.Path("file")
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		r := providerregistry.Registry()
		p := survey.Select{
			Message: "What are you trying to grant access to?",
			Options: r.CLIOptions(),
		}
		var chosen string
		err = survey.AskOne(&p, &chosen)
		if err != nil {
			return err
		}
		var provider providerregistry.RegisteredProvider
		uses, provider, err := r.FromCLIOption(chosen)
		if err != nil {
			return err
		}

		clio.Info("Follow the documentation for setting up the %s provider here: https://docs.commonfate.io/granted-approvals/providers/%s", provider.Description, provider.DefaultID)

		var id string
		err = survey.AskOne(&survey.Input{
			Message: "The ID for the provider",
			Default: provider.DefaultID,
		}, &id, survey.WithValidator(func(ans interface{}) error {
			str, ok := ans.(string)
			if !ok {
				return errors.New("couldn't validate non-string answer")
			}
			if dc.Deployment.Parameters.ProviderConfiguration == nil {
				return nil
			}
			if _, ok := dc.Deployment.Parameters.ProviderConfiguration[str]; ok {
				return fmt.Errorf("provider %s already exists in %s", str, f)
			}
			return nil
		}))

		if err != nil {
			return err
		}

		// set up the config for the specific provider by prompting the user.
		var pcfg gconfig.Config
		if configer, ok := provider.Provider.(gconfig.Configer); ok {
			pcfg = configer.Config()
			for _, v := range pcfg {
				err := deploy.CLIPrompt(v)
				if err != nil {
					return err
				}
			}
		}

		err = deploy.RunConfigTest(ctx, provider.Provider)
		if err != nil {
			return err
		}

		// if tests pass, dump the config and update in the deployment config
		// secret path args requires the id, all provider config includes the provider ID in the path
		idpWith, err := pcfg.Dump(ctx, gconfig.SSMDumper{Suffix: dc.Deployment.Parameters.DeploymentSuffix, SecretPathArgs: []interface{}{id}})
		if err != nil {
			return err
		}
		err = dc.Deployment.Parameters.ProviderConfiguration.Add(id, deploy.Provider{Uses: uses, With: idpWith})
		if err != nil {
			return err
		}

		err = dc.Save(f)
		if err != nil {
			return err
		}

		clio.Success("wrote config to %s", f)
		clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")
		return nil
	},
}
