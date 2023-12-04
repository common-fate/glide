package bootstrap

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/provider-registry-sdk-go/pkg/bootstrapper"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name: "bootstrap",
	Description: `The bootstrap command will create a cloudformation stack that deploys an S3 bucket in your account and return the bucket name.
Bootstrapping is required because Cloudformation requires that resources from S3 be in the same region as the cloudfromation stack.
When deploying a provider you must first copy the provider resources from the Provider Registry to your AWS account in the region that you will be deploying the provider.`,
	Usage: "Bootstrap your AWS account to deploy a provider",

	Action: func(c *cli.Context) error {
		ctx := c.Context
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}
		bs := bootstrapper.NewFromConfig(cfg)
		if err != nil {
			return err
		}

		bucket, err := bs.GetOrDeployBootstrapBucket(ctx)
		if err != nil {
			return err
		}
		clio.Successf("bootstrap bucket created: %s", bucket.AssetsBucket)
		return nil
	},
}
