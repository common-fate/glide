package groups

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var GroupsCommand = cli.Command{
	Name:        "group",
	Subcommands: []*cli.Command{&membersCommand},
	Action:      cli.ShowSubcommandHelp,
}

var membersCommand = cli.Command{
	Name:        "members",
	Subcommands: []*cli.Command{&addCommand},
	Action:      cli.ShowSubcommandHelp,
}

var addCommand = cli.Command{
	Name: "add",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "username", Aliases: []string{"u"}, Usage: "The username of the user to add", Required: true},
		&cli.StringFlag{Name: "group", Aliases: []string{"g"}, Usage: "The group ID to add the user to", Required: true},
	},
	Description: "Add a Cognito user to a group",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}

		username := c.String("username")
		group := c.String("group")

		do, err := deploy.LoadConfig(deploy.DefaultFilename)
		if err != nil {
			return err
		}
		o, err := do.LoadOutput(ctx)
		if err != nil {
			return err
		}
		cog := cognitoidentityprovider.NewFromConfig(cfg)
		_, err = cog.AdminAddUserToGroup(ctx, &cognitoidentityprovider.AdminAddUserToGroupInput{
			GroupName:  &group,
			Username:   &username,
			UserPoolId: &o.UserPoolID,
		})
		if err != nil {
			return err
		}

		zap.S().Infow("added user to group", "username", username, "group", group)

		return nil
	},
}
