package targetgroup

import (
	"errors"
	"net/http"

	"github.com/common-fate/clio"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "targetgroup",
	Description: "Manage Target Groups",
	Usage:       "Manage Target Groups",
	Subcommands: []*cli.Command{
		&CreateCommand,
		&LinkCommand,
		&UnlinkCommand,
		&ListCommand,
	},
}

var CreateCommand = cli.Command{
	Name:        "create",
	Description: "Create a target group",
	Usage:       "Create a target group",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "id", Required: true},
		&cli.StringFlag{Name: "schema-from", Required: true, Usage: "publisher/name@version"},
		&cli.BoolFlag{Name: "ok-if-exists"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		id := c.String("id")
		schemaFrom := c.String("schema-from")
		cfApi, err := types.NewClientWithResponses("http://0.0.0.0:8080")
		if err != nil {
			return err
		}

		res, err := cfApi.AdminCreateTargetGroupWithResponse(ctx, types.AdminCreateTargetGroupJSONRequestBody{
			ID:           id,
			TargetSchema: schemaFrom,
		})
		if err != nil {
			return err
		}
		switch res.StatusCode() {
		case http.StatusCreated:
			clio.Successf("Successfully created the targetgroup: %s", id)
		case http.StatusUnauthorized:
			return errors.New(res.JSON401.Error)
		case http.StatusInternalServerError:
			return errors.New(res.JSON500.Error)
		default:
			return clierr.New("Unhandled response from the Common Fate API", clierr.Infof("Status Code: %d", res.StatusCode()), clierr.Error(string(res.Body)))
		}
		return nil

	},
}

var LinkCommand = cli.Command{
	Name:        "link",
	Description: "Link a handler to a target group",
	Usage:       "Link a handler to a target group",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "target-group", Required: true},
		&cli.StringFlag{Name: "handler", Required: true},
		&cli.StringFlag{Name: "kind", Required: true},
		&cli.IntFlag{Name: "priority", Value: 100},
	},
	Action: func(c *cli.Context) error {

		ctx := c.Context
		cfApi, err := types.NewClientWithResponses("http://0.0.0.0:8080")
		if err != nil {
			return err
		}

		res, err := cfApi.AdminCreateTargetGroupLinkWithResponse(ctx, c.String("target-group"), types.AdminCreateTargetGroupLinkJSONRequestBody{
			DeploymentId: c.String("handler"),
			Priority:     c.Int("priority"),
			Kind:         c.String("kind"),
		})
		if err != nil {
			return err
		}
		switch res.StatusCode() {
		case http.StatusCreated:
			clio.Successf("Successfully linked the handler '%s' with target group '%s' using kind: '%s'", c.String("handler"), c.String("target-group"), c.String("kind"))
		case http.StatusUnauthorized:
			return errors.New(res.JSON401.Error)
		case http.StatusInternalServerError:
			return errors.New(res.JSON500.Error)
		default:
			return clierr.New("Unhandled response from the Common Fate API", clierr.Infof("Status Code: %d", res.StatusCode()), clierr.Error(string(res.Body)))
		}

		return nil
	},
}
