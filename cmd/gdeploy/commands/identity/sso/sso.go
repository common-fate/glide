package sso

import (
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var SSOCommand = cli.Command{
	Name:        "sso",
	Subcommands: []*cli.Command{&enableCommand, &disableCommand, &updateCommand, &samlParametersCommand},
	Action:      cli.ShowSubcommandHelp,
}

var samlParametersCommand = cli.Command{
	Name:        "saml-parameters",
	Description: "Prints a table with parameters required when setting up SAML",
	Usage:       "Prints a table with parameters required when setting up SAML",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		o, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}
		o.PrintSAMLTable()
		return nil
	},
}
