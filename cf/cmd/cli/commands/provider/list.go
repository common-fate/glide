package provider

import (
	"context"
	"errors"
	"os"

	"github.com/common-fate/common-fate/internal/build"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

var ListProvidersCommand = cli.Command{
	Name:        "list",
	Description: "list all providers in the registry",
	Usage:       "list all providers in the registry",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "registry-api-url", Value: build.ProviderRegistryAPIURL, EnvVars: []string{"COMMONFATE_PROVIDER_REGISTRY_API_URL"}, Hidden: true},
	},
	Action: func(c *cli.Context) error {

		ctx := context.Background()
		registryClient, err := providerregistrysdk.NewClientWithResponses(c.String("registry-api-url"))
		if err != nil {
			return errors.New("error configuring provider registry client")
		}

		//check that the provider type matches one in our registry
		res, err := registryClient.ListAllProvidersWithResponse(ctx)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"name", "team", "version", "s3ARN"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)

		for _, d := range res.JSON200.Providers {

			table.Append([]string{
				d.Name, d.Publisher, d.Version, d.LambdaAssetS3Arn,
			})
		}
		table.Render()
		return nil
	},
}
