package deployment

import (
	"errors"
	"fmt"
	"os"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

var DiagnosticCommand = cli.Command{
	Name:        "diagnostic",
	Description: "query a deployment by ID and fetch it's diagnostic information",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true},
	},
	Action: cli.ActionFunc(func(c *cli.Context) error {

		opts := []types.ClientOption{}
		ctx := c.Context

		ID := c.String("id")
		if ID == "" {
			return errors.New("id is required, it can be found by referencing the `deployment list` output")
		}

		cfApi, err := types.NewClientWithResponses("http://0.0.0.0:8080", opts...)
		if err != nil {
			return err
		}
		res, err := cfApi.GetTargetGroupDeploymentWithResponse(ctx, ID)
		if err != nil {
			return err
		}

		if res.JSON200 != nil {
			dep := res.JSON200
			healthEmoji := "🟢"
			if !dep.Healthy {
				healthEmoji = "🔴"
			}
			fmt.Println("Diagnostic Logs:")
			fmt.Printf("%s %s %s %s\n", dep.Id, dep.AwsAccount, dep.AwsRegion, healthEmoji)

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{
				"Level", "Message"})
			table.SetAutoWrapText(false)
			table.SetAutoFormatHeaders(true)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")
			table.SetHeaderLine(false)
			table.SetBorder(false)

			for _, d := range dep.Diagnostics {
				table.Append([]string{
					d.Level, d.Message,
				})
			}
			table.Render()
		} else {
			clio.Error("no deployments found")
		}

		return nil
	}),
}
