package provider

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var removeCommand = cli.Command{
	Name:        "remove",
	Description: "Remove an existing access provider",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Usage: "An identifier of the provider to be removed."},
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
				Message: "Select the id of the Provider that you want to remove",
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
			clio.Error("Provider configuration doesn't exist. Unable to remove provider '%s'", chosen)
			return nil
		}

		var confirm bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to remove '%s'", chosen),
		}

		err = survey.AskOne(prompt, &confirm)
		if err != nil {
			return err
		}

		if !confirm {
			clio.Info("Removing provider aborted.")
			return nil
		}

		delete(dc.Deployment.Parameters.ProviderConfiguration, chosen)

		f := c.Path("file")
		err = dc.Save(f)
		if err != nil {
			return err
		}

		clio.Success("Successfully removed provider %s", chosen)
		clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")
		return nil
	},
}
