package provider

import (
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/cmd/gdeploy/commands/provider/generate"
	mw "github.com/common-fate/common-fate/cmd/gdeploy/middleware"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/table"
	"github.com/common-fate/provider-registry-sdk-go/pkg/bootstrapper"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	registryclient "github.com/common-fate/provider-registry-sdk-go/pkg/registryclient"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "provider-v2",
	Description: "Explore and manage Providers from the Provider Registry",
	Usage:       "Explore and manage Providers from the Provider Registry",
	Subcommands: []*cli.Command{
		mw.WithBeforeFuncs(&BootstrapCommand, mw.RequireAWSCredentials()),
		&ListCommand,
		&generate.Command,
		mw.WithBeforeFuncs(&deployCommand, mw.RequireAWSCredentials()),
		mw.WithBeforeFuncs(&destroyCommand, mw.RequireAWSCredentials()),
	},
}

var BootstrapCommand = cli.Command{
	Name:        "bootstrap",
	Description: "Before you can deploy a Provider, you will need to bootstrap it. This process will copy the files from the Provider Registry to your bootstrap bucket.",
	Usage:       "Copy a Provider into your AWS account",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true, Usage: "publisher/name@version"},
		&cli.BoolFlag{Name: "force", Usage: "force copy provider assets to bootstrap bucket"},
	},

	Action: func(c *cli.Context) error {
		ctx := c.Context

		registry, err := registryclient.New(ctx)
		if err != nil {
			return errors.Wrap(err, "configuring provider registry client")
		}

		id := c.String("id")

		provider, err := providerregistrysdk.ParseProvider(id)
		if err != nil {
			return err
		}

		//check that the provider type matches one in our registry
		res, err := registry.GetProviderWithResponse(ctx, provider.Publisher, provider.Name, provider.Version)
		if err != nil {
			return err
		}

		clio.Success("Provider exists in the registry")

		awsConfig, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}

		bs := bootstrapper.NewFromConfig(awsConfig)
		if err != nil {
			return err
		}

		clio.Info("Copying provider assets...")

		err = bs.CopyProviderFiles(ctx, *res.JSON200, bootstrapper.WithForceCopy(c.Bool("force")))
		if err != nil {
			return err
		}

		return nil
	},
}

func getProviderId(publisher, name, version string) string {
	return publisher + "/" + name + "@" + version
}

var ListCommand = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Description: "List providers",
	Usage:       "List providers",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		registry, err := registryclient.New(ctx)
		if err != nil {
			return errors.Wrap(err, "configuring provider registry client")
		}

		res, err := registry.ListAllProvidersWithResponse(ctx, &providerregistrysdk.ListAllProvidersParams{
			WithDev: aws.Bool(false),
		})
		if err != nil {
			return err
		}
		tbl := table.New(os.Stderr)
		tbl.Columns("ID", "Name", "Publisher", "Version", "Kinds")
		for _, d := range res.JSON200.Providers {
			var kinds []string
			if d.Schema.Targets != nil {
				for kind := range *d.Schema.Targets {
					kinds = append(kinds, kind)
				}
			}
			tbl.Row(getProviderId(d.Publisher, d.Name, d.Version), d.Name, d.Publisher, d.Version, strings.Join(kinds, ", "))
		}
		return tbl.Flush()
	},
}
