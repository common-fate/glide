package groups

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var CreateCommand = cli.Command{
	Name: "create",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "group-name", Aliases: []string{"n"}, Usage: "The group ID", Required: true},
		&cli.StringFlag{Name: "description", Aliases: []string{"desc"}, Usage: "The description of the group"},
	},
	Description: "Create a new Cognito user group",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		group := c.String("group-name")
		desc := c.String("description")

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

		o, err = dc.LoadOutput(ctx)
		if err != nil {
			return err
		}

		cog := cognitoidentityprovider.NewFromConfig(cfg)
		_, err = cog.CreateGroup(ctx, &cognitoidentityprovider.CreateGroupInput{
			GroupName:   &group,
			UserPoolId:  &o.UserPoolID,
			Description: &desc,
		})
		if err != nil {
			return err
		}

		clio.Success("Successfully created group '%s'", group)
		clio.Warn("Run 'gdeploy identity sync' to sync your changes now.")
		return nil
	},
}
