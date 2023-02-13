package targetgroup

import (
	"errors"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "targetgroup",
	Description: "target group",
	Usage:       "target group",
	Subcommands: []*cli.Command{
		&CreateCommand,
		&LinkCommand,
	},
}

var CreateCommand = cli.Command{
	Name:        "create",
	Description: "create a target group",
	Usage:       "create a target group",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "schema-from", Required: true},
		&cli.BoolFlag{Name: "ok-if-exists"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		id := c.Args().First()
		if id == "" {
			return errors.New("id argument must be provided")
		}

		schemaFrom := c.String("schema-from")
		cfApi, err := types.NewClientWithResponses("http://0.0.0.0:8080")
		if err != nil {
			return err
		}

		result, err := cfApi.CreateTargetGroupWithResponse(ctx, types.CreateTargetGroupJSONRequestBody{
			ID:           id,
			TargetSchema: schemaFrom,
		})
		if err != nil {
			return err
		}

		switch result.StatusCode() {
		case 201:
			clio.Successf("created target group '%s'", id)
			return nil
		default:
			return errors.New(string(result.Body))
		}

	},
}

var LinkCommand = cli.Command{
	Name:        "link",
	Description: "link a deployment to a target group",
	Usage:       "link a deployment to a target group",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "group", Required: true},
		&cli.StringFlag{Name: "deployment", Required: true},
		&cli.IntFlag{Name: "priority", Value: 100},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		group := c.String("group")
		deployment := c.String("deployment")
		priority := c.Int("priority")

		if priority < 0 || priority > 999 {
			return errors.New("priority must be a number between 0 and 999")

		}

		cfApi, err := types.NewClientWithResponses("http://0.0.0.0:8080")
		if err != nil {
			return err
		}

		result, err := cfApi.CreateTargetGroupLinkWithResponse(ctx, group, types.CreateTargetGroupLinkJSONRequestBody{
			DeploymentId: deployment,
			Priority:     priority,
		})
		if err != nil {
			return err
		}

		if result.StatusCode() != 200 {
			return errors.New(string(result.Body))
		}

		clio.Successf("linked deployment '%s' with target group '%s'", c.String("deployment"), c.String("name"))
		return nil
	},
}
