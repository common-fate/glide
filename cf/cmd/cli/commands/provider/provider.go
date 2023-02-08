package targetgroup

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/cf/pkg/bootstrapper"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
	Name:        "provider",
	Description: "manage provider",
	Usage:       "manage provider",
	Subcommands: []*cli.Command{
		&BootstrapCommand,
	},
}

var BootstrapCommand = cli.Command{
	Name:        "bootstrap",
	Description: "bootstrap a provider from the registry",
	Usage:       "bootstrap a provider from the registry",
	Action: func(c *cli.Context) error {
		id := c.Args().First()
		if id == "" {
			return errors.New("id argument must be provided")
		}

		var cfg config.ProviderDeploymentCLI
		ctx := context.Background()
		_ = godotenv.Load()

		err := envconfig.Process(ctx, &cfg)
		if err != nil {
			return err
		}

		registryClient, err := providerregistrysdk.NewClientWithResponses(cfg.ProviderRegistryAPIURL)
		if err != nil {
			return errors.New("error configuring provider registry client")
		}

		keys := strings.Split(id, "/")

		if len(keys) != 3 {
			return errors.New("incorrect provider id given")
		}

		team := keys[0]
		name := keys[1]
		version := keys[2]

		//check that the provider type matches one in our registry
		res, err := registryClient.GetProviderWithResponse(ctx, team, name, version)
		if err != nil {
			return err
		}
		if res.StatusCode() != http.StatusOK {
			return errors.New("provider for that version does not exist")
		}
		//get the bootstrap bucket name
		bs, err := bootstrapper.New(ctx)
		if err != nil {
			return err
		}
		bootstrapBucket, err := bs.GetOrDeployBootstrapBucket(ctx)
		if err != nil {
			return err
		}
		//work out the lambda asset path
		lambdaAssetPath := path.Join(team, name, version)

		//copy the provider assets into the bucket
		err = bs.CopyProviderAsset(ctx, res.JSON200.LambdaAssetS3Arn, lambdaAssetPath, bootstrapBucket)

		if err != nil {
			return err
		}
		clio.Log(fmt.Sprintf("copied %s into %s", id, path.Join(res.JSON200.LambdaAssetS3Arn, lambdaAssetPath)))
		return nil
	},
}
