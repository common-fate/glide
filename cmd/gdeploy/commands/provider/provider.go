package provider

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "providers",
	Aliases:     []string{"provider"},
	Description: "Manage your Access Providers",
	Usage:       "Manage your Access Providers",
	Subcommands: []*cli.Command{
		&addCommand, &removeCommand, &updateCommand,
	},
}

func DeprecatedWarn(c *cli.Context) error {
	// add a deprecated warning here
	clio.Warn("Warning: Configuring access providers via gdeploy has been deprecated in favour of an interactive setup flow in your admin dashboard. Head there now get started setting up your provider.")
	clio.Warn("Providers can now be setup using the deployed frontend found below:")

	// attempt to fetch the CloudFrontDomain from the context
	// if it's not there, then we can't do anything
	ctx := c.Context
	dc, err := deploy.ConfigFromContext(ctx)
	if err != nil {
		return err
	}
	o, err := dc.LoadOutput(ctx)
	if err != nil {
		return err
	}
	// If we do find it, then we can prompt the user
	feDomain := o.FrontendDomainOutput
	if feDomain != "" {
		url := "https://" + feDomain + "/admin/providers/setup"
		clio.Warn("Docs: https://docs.commonfate.io/granted-approvals/providers/access-providers")
		clio.Warn("Provider Setup Page: " + url)
	}
	return nil
}
