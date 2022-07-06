package groups

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
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
		&cli.StringFlag{Name: "username", Aliases: []string{"u"}, Usage: "the username of the user to add", Required: true},
		&cli.StringFlag{Name: "group", Aliases: []string{"g"}, Usage: "the group ID to add the user to", Required: true},
	},
	Description: "Add a Cognito user to a group",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}

		username := c.String("username")
		group := c.String("group")

		f := c.Path("file")
		do, err := deploy.LoadConfig(f)
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

		clio.Success("added user %s to group '%s'", username, group)

		return nil
	},
}
