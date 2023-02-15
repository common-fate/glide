package deployment

import (
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/internal/build"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

var ListCommand = cli.Command{
	Name:        "list",
	Description: "list deployments",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "commonfate-api", Value: build.CommonFateAPIURL, EnvVars: []string{"COMMONFATE_API_URL"}, Hidden: true},
	},
	Action: cli.ActionFunc(func(c *cli.Context) error {

		opts := []types.ClientOption{}
		ctx := c.Context

		cfApi, err := types.NewClientWithResponses(c.String("commonfate-api"), opts...)
		if err != nil {
			return err
		}
		res, err := cfApi.ListTargetGroupDeploymentsWithResponse(ctx)
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
