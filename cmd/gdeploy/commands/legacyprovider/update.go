package legacyprovider

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/accesshandler/pkg/providerregistry"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/urfave/cli/v2"
)

var updateCommand = cli.Command{
	Name:        "update",
	Description: "Update configuration for existing provider",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Usage: "An identifier for the provider"},
		&cli.StringSliceFlag{Name: "with", Usage: "Configuration settings for the provider, in key=value pairs"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		var chosen string
		id := c.String("id")

		// if no '-id' flag provided then stdout current providers as options to choose from.
		if id == "" {
			userProvidersMap := dc.Deployment.Parameters.ProviderConfiguration

			var CLIOptions []string
			for k := range userProvidersMap {
				CLIOptions = append(CLIOptions, k)
			}

			p := survey.Select{
				Message: "Select the id of the Provider that you want to update",
				Options: CLIOptions,
			}

			err = survey.AskOne(&p, &chosen)
			if err != nil {
				return err
			}

		} else {
			chosen = id
		}

		if _, ok := dc.Deployment.Parameters.ProviderConfiguration[chosen]; !ok {
			return clierr.New(fmt.Sprintf("Provider configuration doesn't exist. Unable to remove provider '%s'", chosen), clierr.Info("Try using 'gdeploy providers add' to add a new provider."))
		}

		with := map[string]string{}
		withArgs := c.StringSlice("with")
		uses := ""

		r := providerregistry.Registry()

		// Need to figure out the registeredProvider from chosen option.
		// Parse uses to see provider type and use lookup func to get the actual registered provider.
		currentConfig := dc.Deployment.Parameters.ProviderConfiguration[chosen]
		uses = currentConfig.Uses
		providerType, version, err := providerregistry.ParseUses(uses)
		if err != nil {
			return err
		}

		registeredProvider, err := r.Lookup(providerType, version)
		if err != nil {
			return err
		}

		if len(withArgs) == 0 {
			var pcfg gconfig.Config
			if configer, ok := registeredProvider.Provider.(gconfig.Configer); ok {
				pcfg = configer.Config()

				// load the current values into this config
				err = pcfg.Load(ctx, &gconfig.MapLoader{Values: currentConfig.With})
				if err != nil {
					return err
				}

				// CLI prompt will have defaults values loaded.
				for _, v := range pcfg {
					err := deploy.CLIPrompt(v)
					if err != nil {
						return err
					}
				}
			}

			err = deploy.RunConfigTest(ctx, registeredProvider.Provider)
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

		err = dc.Deployment.Parameters.ProviderConfiguration.Update(chosen, deploy.Provider{Uses: uses, With: with})
		if err != nil {
			return err
		}

		f := c.Path("file")
		err = dc.Save(f)
		if err != nil {
			return err
		}

		clio.Successf("wrote config to %s", f)
		clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")
		return nil
	},
}
