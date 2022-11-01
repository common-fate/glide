package users

import (
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/common-fate/clio"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var DeleteCommand = cli.Command{
	Name: "delete",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "username", Aliases: []string{"u"}, Usage: "The username of the user to delete (should be an email address)", Required: true},
	},
	Description: "Delete a Cognito user",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		username := c.String("username")
		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}
		o, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}

		cog := cognitoidentityprovider.NewFromConfig(cfg)

		_, err = cog.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
			UserPoolId: &o.UserPoolID,
			Username:   &username,
		})

		if err != nil {
			return err
		}

		clio.Successf("Deleted user %s", username)
		clio.Warn("Run 'gdeploy identity sync' to sync your changes now.")
		return nil
	}}
