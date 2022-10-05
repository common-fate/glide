package members

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var removeCommand = cli.Command{
	Name: "remove",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "username", Aliases: []string{"u"}, Usage: "The username of the user to add", Required: true},
		&cli.StringFlag{Name: "group", Aliases: []string{"g"}, Usage: "The group ID to add the user to", Required: true},
	},
	Description: "remove a Cognito user from a group",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		username := c.String("username")
		group := c.String("group")

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

		_, err = cog.AdminRemoveUserFromGroup(ctx, &cognitoidentityprovider.AdminRemoveUserFromGroupInput{
			GroupName:  &group,
			Username:   &username,
			UserPoolId: &o.UserPoolID,
		})

		if err != nil {
			return err
		}

		clio.Success("Removed user %s from group '%s'", username, group)
		clio.Warn("Run 'gdeploy identity sync' to sync your changes now.")
		return nil
	},
}
