package provider

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
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "bootstrapBucket", Aliases: []string{"bb"}, Usage: "The name of the bootstrap bucket provider is deployed to", EnvVars: []string{"DEPLOYMENT_BUCKET"}},
	},
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

		keys := strings.Split(id, "@")

		if len(keys) != 2 {
			return errors.New("incorrect provider id given")
		}

		version := keys[1]

		teamAndName := strings.Split(keys[0], "/")

		team := teamAndName[0]
		name := teamAndName[1]

		//check that the provider type matches one in our registry
		res, err := registryClient.GetProviderWithResponse(ctx, team, name, version)
		if err != nil {
			return err
		}
		if res.StatusCode() != http.StatusOK {
			return fmt.Errorf("provider for that version does not exist: %s", team+name+version)
		}
		//get the bootstrap bucket name
		bs, err := bootstrapper.New(ctx)
		if err != nil {
			return err
		}

		//get bootstrap bucket

		//read from flag
		bootstrapBucket := c.String("bootstrapBucket")

		//work out the lambda asset path
		lambdaAssetPath := path.Join(team, name, version)

		//copy the provider assets into the bucket (this will also copy the cloudformation template too)
		err = bs.CopyProviderAsset(ctx, res.JSON200.LambdaAssetS3Arn, lambdaAssetPath, bootstrapBucket)

		if err != nil {
			return err
		}

		clio.Log(fmt.Sprintf("copied %s into %s", id, path.Join(bootstrapBucket, lambdaAssetPath)))
		return nil
	},
}
