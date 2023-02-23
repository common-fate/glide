package handler

import (
	"errors"
	"net/http"
	"os"

	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

var ListCommand = cli.Command{
	Name:        "list",
	Description: "List handlers",
	Usage:       "List handlers",
	Action: cli.ActionFunc(func(c *cli.Context) error {
		ctx := c.Context
		cfApi, err := types.NewClientWithResponses("http://0.0.0.0:8080")
		if err != nil {
			return err
		}
		res, err := cfApi.AdminListHandlersWithResponse(ctx)
		if err != nil {
			return err
		}
		switch res.StatusCode() {
		case http.StatusOK:
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
				health := "healthy"
				if !d.Healthy {
					health = "unhealthy"
				}
				table.Append([]string{
					d.Id, d.AwsAccount, d.AwsRegion, health,
				})
			}
			table.Render()
		case http.StatusUnauthorized:
			return errors.New(res.JSON401.Error)
		case http.StatusInternalServerError:
			return errors.New(res.JSON500.Error)
		default:
			return clierr.New("Unhandled response from the Common Fate API", clierr.Infof("Status Code: %d", res.StatusCode()), clierr.Error(string(res.Body)))
		}
		return nil
	}),
}
