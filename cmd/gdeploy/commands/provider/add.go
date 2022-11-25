package provider

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/clio"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providerregistry"
	"github.com/common-fate/granted-approvals/cmd/gdeploy/middleware"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/urfave/cli/v2"
)

var addCommand = cli.Command{
	Name:        "add",
	Description: "Add an access provider",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Usage: "An identifier for the provider"},
		&cli.StringFlag{Name: "uses", Usage: "The provider type and version"},
		&cli.StringSliceFlag{Name: "with", Usage: "Configuration settings for the provider, in key=value pairs"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		f := c.Path("file")
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		r := providerregistry.Registry()

		uses := c.String("uses")
		var provider providerregistry.RegisteredProvider

		if uses == "" {
			o, err := dc.LoadOutput(ctx)
			if err != nil {
				return err
			}
			setupURL := o.FrontendURL() + "/admin/providers/setup"
			clio.Warn("Configuring access providers interactively via gdeploy has been deprecated.")
			clio.Warn("For the best experience, we recommend using the interactive setup documentation available as part of your deployment of Common Fate.")
			clio.Warnf("Open this link to the provider setup page then select the provider you want to configure: %s", setupURL)
			clio.Warn("At the end of the interactive setup, your configuration will be tested with helpful validation errors should anything go wrong.")
			clio.Warn("When everything is working, you will be given a gdeploy command to run to get everything setup.")
			clio.Warn("If you want to read more about our providers and getting setup, checkout our documentation here: https://docs.commonfate.io/granted-approvals/providers/access-providers")
			clio.NewLine()
			p := survey.Select{
				Message: "What are you trying to grant access to?",
				Options: r.CLIOptions(),
			}
			var chosen string
			err = survey.AskOne(&p, &chosen)
			if err != nil {
				return err
			}
			uses, provider, err = r.FromCLIOption(chosen)
			if err != nil {
				return err
			}
			clio.Infof("Follow the documentation for setting up the %s provider here: https://docs.commonfate.io/granted-approvals/providers/%s", provider.Description, provider.DefaultID)
		} else {
			p, err := r.LookupByUses(uses)
			if err != nil {
				return err
			}
			provider = *p
		}

		id := c.String("id")
		if id == "" {
			id = dc.Deployment.Parameters.ProviderConfiguration.GetIDForNewProvider(provider.DefaultID)
		}

		if _, ok := dc.Deployment.Parameters.ProviderConfiguration[id]; ok {
			return fmt.Errorf("provider %s already exists in %s", id, f)
		}

		with := map[string]string{}

		withArgs := c.StringSlice("with")

		if len(withArgs) == 0 {
			// we need AWS credentials to dump any secrets and test the provider.
			// grab the AWS credential checker middleware and execute it manually.
			checkCreds := middleware.RequireAWSCredentials()
			err := checkCreds(c)
			if err != nil {
				return err
			}

			var pcfg gconfig.Config
			// set up the config for the specific provider by prompting the user.
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
			with, err = pcfg.Dump(ctx, gconfig.SSMDumper{Suffix: dc.Deployment.Parameters.DeploymentSuffix, SecretPathArgs: []interface{}{id}})
			if err != nil {
				return err
			}

		} else {
			// The user has provided some config via the --with arguments, so skip the interactive flow. The user has likely used the
			// guided setup UI, or they know what they're doing.
			//
			// We also skip writing secrets to SSM, as we assume they have already been written. We don't support passing secrets via '--with'
			// as it pollutes the user's shell history with their credentials.
			//
			// We could perform a configuration test here in future to provide some extra assurance to the user that their values are correct.

			// parse the key=value pairs in the 'with' argument.
			for _, kv := range withArgs {
				segments := strings.Split(kv, "=")
				if len(segments) != 2 {
					return fmt.Errorf("could not parse 'with' argument %s: must be in key=value format", kv)
				}

				key, val := segments[0], segments[1]
				with[key] = val
			}
		}

		err = dc.Deployment.Parameters.ProviderConfiguration.Add(id, deploy.Provider{Uses: uses, With: with})
		if err != nil {
			return err
		}

		err = dc.Save(f)
		if err != nil {
			return err
		}

		clio.Successf("wrote config to %s", f)
		clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")
		return nil
	},
}
