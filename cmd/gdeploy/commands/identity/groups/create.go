package groups

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var createCommand = cli.Command{
	Name: "create",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "name", Aliases: []string{"n"}, Usage: "the group name to create", Required: true},
	},
	Description: "Create a Cognito group",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		name := c.String("name")

		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		o, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}
		cog := cognitoidentityprovider.NewFromConfig(cfg)
		_, err = cog.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
			GroupName:  &name,
			UserPoolId: &o.UserPoolID,
		})
		if err != nil {
			return err
		}

		clio.Success("created group %s", name)

		return nil
	},
}
