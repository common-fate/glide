package commands

import (
	"fmt"

	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/internal"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var InitCommand = cli.Command{
	Name:        "init",
	Description: "Set up a new Common Fate deployment configuration file",
	Usage:       "Set up a new Common Fate deployment configuration file",
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "overwrite", Usage: "Force an existing deployment configuration file to be overwritten"},
		&cli.StringFlag{Name: "name", Usage: "The name of the CloudFormation stack to create"},
		&cli.StringFlag{Name: "account", Usage: "The account ID to deploy to"},
		&cli.StringFlag{Name: "region", Usage: "The region to deploy to"},
		&cli.StringFlag{Name: "version", Usage: "The version to deploy"},
		&cli.StringFlag{Name: "cognito-domain-prefix", Usage: "The prefix for the Cognito Sign in URL"},
	},
	Action: func(c *cli.Context) error {
		err := ensureConfigDoesntExist(c)
		if err != nil {
			return err
		}

		err = internal.PrintAnalyticsNotice(true)
		if err != nil {
			clio.Debugf("error printing analytics notice: %s", err)
		}

		cfg, err := deploy.SetupReleaseConfig(c)
		if err != nil {
			return err
		}

		f := c.Path("file")

		err = cfg.Save(f)
		if err != nil {
			return err
		}

		clio.Successf("Wrote config to %s", f)
		clio.Warn("Nothing has been deployed yet. To finish deploying Common Fate, run 'gdeploy create' to create the CloudFormation stack in AWS.")
		return nil
	},
}

// sanity check: verify that a config file doesn't already exist.
// if it does, the user may have run this command by mistake.
func ensureConfigDoesntExist(c *cli.Context) error {
	f := c.Path("file")
	overwrite := c.Bool("overwrite")

	_, err := deploy.LoadConfig(f)
	if err != nil && err != deploy.ErrConfigNotExist {
		// we don't know how to handle this, so return it.
		return err
	}
	var deprecatedErr error
	if f == deploy.DefaultFilename {
		_, deprecatedErr = deploy.LoadConfig(deploy.DeprecatedDefaultFilename)
		if deprecatedErr != nil && deprecatedErr != deploy.ErrConfigNotExist {
			// we don't know how to handle this, so return it.
			return deprecatedErr
		}
	}

	if err == deploy.ErrConfigNotExist && deprecatedErr == deploy.ErrConfigNotExist {
		// no config file at risk of being overwritten, so return
		return nil
	}
	if deprecatedErr != nil {
		f = deploy.DeprecatedDefaultFilename
	}
	if overwrite {
		clio.Warnf("--overwrite has been set, the config file %s will be overwritten", f)
		return nil
	}

	// if we get here, the config file exists and is at risk of being overwritten.
	return clierr.New(fmt.Sprintf("A deployment config file %s already exists in this folder.\ngdeploy will exit to avoid overwriting this file, in case you've run this command by mistake.", f),
		deploy.DeprecatedDefaultFilenameWarning,
		clierr.Info(`
To fix this, take one of the following actions:
  a) run 'gdeploy init' from a different folder
  b) run 'gdeploy -f other-config.toml init' to use a separate config file
  c) run 'gdeploy init --overwrite' to force overwriting the existing config
`))

}
