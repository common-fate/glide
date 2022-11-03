package slack

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	slacknotifier "github.com/common-fate/granted-approvals/pkg/notifiers/slack"
	"github.com/urfave/cli/v2"
)

var disableSlackCommand = cli.Command{
	Name:        "disable",
	Description: "disable slack integration",
	Action: func(c *cli.Context) error {
		ctx := c.Context
		f := c.Path("file")

		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}

		dc.Deployment.Parameters.NotificationsConfiguration.Remove(slacknotifier.NotificationsTypeSlack)
		err = dc.Save(f)
		if err != nil {
			return err
		}
		clio.Success("Successfully disabled Slack")
		clio.Warn("Your changes won't be applied until you redeploy. Run 'gdeploy update' to apply the changes to your CloudFormation deployment.")
		return nil
	},
}
