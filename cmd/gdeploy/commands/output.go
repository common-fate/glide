package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

var Output = cli.Command{
	Name:        "output",
	Aliases:     []string{"outputs"},
	Description: "Get an output value from your Common Fate deployment",
	Usage:       "Get an output value from your Common Fate deployment",
	ArgsUsage:   "<key>",
	Action: func(c *cli.Context) error {
		ctx := c.Context

		dc, err := deploy.ConfigFromContext(ctx)
		if err != nil {
			return err
		}
		o, err := dc.LoadOutput(ctx)
		if err != nil {
			return err
		}

		key := c.Args().First()
		if key == "" {
			return errors.New("usage: gdeploy output <key>")
		}
		val, err := o.Get(key)
		if err != nil {
			return fmt.Errorf("%s. available keys: %s", err.Error(), strings.Join(o.Keys(), ", "))
		}
		fmt.Println(val)
		return nil
	},
}
