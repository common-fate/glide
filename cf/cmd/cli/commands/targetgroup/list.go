package targetgroup

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
	Description: "list target groups",
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
		res, err := cfApi.ListTargetGroupsWithResponse(ctx)
		if err != nil {
			return err
		}

		if res.JSON200 != nil {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "From"})
			table.SetAutoWrapText(false)
			table.SetAutoFormatHeaders(true)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")
			table.SetHeaderLine(false)
			table.SetBorder(false)

			for _, d := range res.JSON200.TargetGroups {

				table.Append([]string{
					d.Id, d.TargetSchema.From,
				})
			}
			table.Render()
		} else {
			clio.Error("no deployments found")
		}

		return nil
	}),
}
