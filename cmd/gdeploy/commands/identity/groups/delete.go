package groups

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/common-fate/clio"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var DeleteCommand = cli.Command{
	Name: "delete",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "group-name", Aliases: []string{"n"}, Usage: "The name of the group to be deleted", Required: true},
	},
	Description: "Delete Cognito user group if exists",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		group := c.String("group-name")

		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		// prevent the user deleting the administrators group
		// it is created by the stack deployment automatically
		if group == deploy.DefaultCommonFateAdministratorsGroup || group == dc.Deployment.Parameters.AdministratorGroupID {
			return errors.New("you cannot delete the administrators group")
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
		_, err = cog.DeleteGroup(ctx, &cognitoidentityprovider.DeleteGroupInput{
			GroupName:  &group,
			UserPoolId: &o.UserPoolID,
		})
		if err != nil {
			return err
		}

		clio.Successf("Successfully deleted group '%s'", group)
		clio.Warn("Run 'gdeploy identity sync' to sync your changes now.")
		return nil
	},
}
