package bootstrap

import (
	"errors"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/cf/pkg/bootstrapper"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "bootstrap",
	Description: "Bootstrap a cloud account for deploying access providers",
	Usage:       "Bootstrap a cloud account for deploying access providers",

	Action: func(c *cli.Context) error {
		ctx := c.Context
		cloud := c.Args().First()
		if cloud == "" || cloud != "aws" {
			return errors.New("cloud argument must be supplied, supports clouds are [aws]")
		}

		bs, err := bootstrapper.New(ctx)
		if err != nil {
			return err
		}
		bootstrapBucket, err := bs.GetOrDeployBootstrapBucket(ctx)
		if err != nil {
			return err
		}
		clio.Log(bootstrapBucket)
		return nil
	},
}
