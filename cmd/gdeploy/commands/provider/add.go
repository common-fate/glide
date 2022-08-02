package provider

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/briandowns/spinner"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/lookup"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/urfave/cli/v2"
)

var addCommand = cli.Command{
	Name:        "add",
	Description: "Add an access provider",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "overwrite", Usage: "force SSM parameters to be overwritten if they exist"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		r := lookup.Registry()
		p := survey.Select{
			Message: "What are you trying to grant access to?",
			Options: r.CLIOptions(),
		}
		var chosen string
		err := survey.AskOne(&p, &chosen)
		if err != nil {
			return err
		}
		var provider lookup.RegisteredProvider
		uses, provider, err := r.FromCLIOption(chosen)
		if err != nil {
			return err
		}

		f := c.Path("file")
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

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
		if configer, ok := provider.Provider.(providers.Configer); ok {
			pcfg = configer.Config()

			// ask the user for config values, if the provider supports it.
			err := promptForConfig(c, id, pcfg)
			if err != nil {
				return err
			}
		}

		err = dc.AddProvider(id, deployProvider)
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

func promptForConfig(c *cli.Context, providerID string, pcfg gconfig.Config) error {
	ctx := c.Context
	// split the config into regular config, and secrets.
	// this will allow us to display user-friendly prompts for each section.
	var cfg, secrets gconfig.Config
	for _, v := range pcfg {
		if s, ok := v.(gconfig.Secret); ok && s.IsSecret() {
			secrets = append(secrets, v)
		} else {
			cfg = append(cfg, v)
		}
	}

	for _, v := range cfg {
		err := gconfig.CLIPrompt(v)
		if err != nil {
			return err
		}
	}

	if len(secrets) == 0 {
		// return early if there's no secrets for us to set up.
		return nil
	}

	awsCfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return err
	}

	client := ssm.NewFromConfig(awsCfg)

	clio.Info("This provider requires some sensitive credentials. These will be uploaded to AWS SSM using the path '/granted/%s/<value>'. A reference to the secrets will be stored in your configuration file.", providerID)

	for _, v := range secrets {
		err := gconfig.CLIPrompt(v)
		if err != nil {
			return err
		}
		ssmKey := fmt.Sprintf("/granted/providers/%s/%s", providerID, v.Key())

		si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		si.Suffix = " writing secret to " + ssmKey
		si.Writer = os.Stderr
		si.Start()

		val := v.Get()

		overwrite := c.Bool("overwrite")

		_, err = client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:      &ssmKey,
			Value:     &val,
			Type:      types.ParameterTypeSecureString,
			Overwrite: overwrite,
		})

		si.Stop()

		var pae *types.ParameterAlreadyExists
		if errors.As(err, &pae) {
			clio.Warn(`The SSM parameter %s has already been set. The parameter has not been updated to avoid overwriting an existing value.

If you want to overwrite the parameter, take one of the following actions:
  a) Update the parameter manually using the AWS CLI or the AWS console.
  b) Run 'gdeploy provider add --overwrite' to force overwriting existing parameters.
`, ssmKey)
			continue
		}

		if err != nil {
			return err
		}

		clio.Info("Wrote %s to AWS SSM parameter %s", v.Key(), ssmKey)
	}

	return nil
}
