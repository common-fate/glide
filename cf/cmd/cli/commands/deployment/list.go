package deployment

import (
	"os"

	cf_cli_client "github.com/common-fate/cli/pkg/client"
	cf_cli_config "github.com/common-fate/cli/pkg/config"

	"github.com/common-fate/clio"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

var ListCommand = cli.Command{
	Name:        "list",
	Description: "list deployments",
	Usage:       "list deployments",
	Action: cli.ActionFunc(func(c *cli.Context) error {

		ctx := c.Context

		cfg, err := cf_cli_config.Load()
		if err != nil {
			return err
		}

		cf, err := cf_cli_client.FromConfig(ctx, cfg)
		if err != nil {
			return err
		}

		res, err := cf.ListTargetGroupDeploymentsWithResponse(ctx)
		if err != nil {
			return err
		}

		if res.JSON200 != nil {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Account", "Region", "Health"})
			table.SetAutoWrapText(false)
			table.SetAutoFormatHeaders(true)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")
			table.SetHeaderLine(false)
			table.SetBorder(false)

			for _, d := range res.JSON200.Res {
				healthEmoji := "🟢"
				if !d.Healthy {
					healthEmoji = "🔴"
				}
				table.Append([]string{
					d.Id, d.AwsAccount, d.AwsRegion, healthEmoji,
				})
			}
			table.Render()
		} else {
			clio.Error("no deployments found")
		}

		return nil
	}),
}
