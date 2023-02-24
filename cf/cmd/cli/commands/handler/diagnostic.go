package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

var DiagnosticCommand = cli.Command{
	Name:        "diagnostic",
	Description: "List diagnostic logs for a handler",
	Usage:       "List diagnostic logs for a handler",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true},
	},
	Action: cli.ActionFunc(func(c *cli.Context) error {
		ctx := c.Context
		id := c.String("id")
		cfApi, err := types.NewClientWithResponses("http://0.0.0.0:8080")
		if err != nil {
			return err
		}
		res, err := cfApi.AdminGetHandlerWithResponse(ctx, id)
		if err != nil {
			return err
		}

		switch res.StatusCode() {
		case http.StatusOK:
			health := "healthy"
			if !res.JSON200.Healthy {
				health = "unhealthy"
			}
			fmt.Println("Diagnostic Logs:")
			fmt.Printf("%s %s %s %s\n", res.JSON200.Id, res.JSON200.AwsAccount, res.JSON200.AwsRegion, health)

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

			for _, d := range res.JSON200.Diagnostics {
				table.Append([]string{
					d.Level, d.Message,
				})
			}
			table.Render()
		case http.StatusUnauthorized:
			return errors.New(res.JSON401.Error)
		case http.StatusNotFound:
			return errors.New(res.JSON404.Error)
		case http.StatusInternalServerError:
			return errors.New(res.JSON500.Error)
		default:
			return clierr.New("Unhandled response from the Common Fate API", clierr.Infof("Status Code: %d", res.StatusCode()), clierr.Error(string(res.Body)))
		}
		return nil
	}),
}
