package sync

import (
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var SyncCommand = cli.Command{
	Name: "sync",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		_, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		// _, err := dc.LoadOutput(ctx)
		// if err != nil {
		// 	return err
		// }
		//  @TODO call sync lambda

		return nil
	}}
