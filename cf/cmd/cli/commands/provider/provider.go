package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/internal/build"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
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
		&cli.StringFlag{Name: "bootstrapBucket", Aliases: []string{"bb"}, Usage: "The name of the bootstrap bucket to copy assets into", EnvVars: []string{"DEPLOYMENT_BUCKET"}},
		&cli.StringFlag{Name: "registryAPIURL", Value: build.ProviderRegistryAPIURL, EnvVars: []string{"COMMONFATE_PROVIDER_REGISTRY_API_URL"}, Hidden: true},
	},

	Action: func(c *cli.Context) error {

		id := c.Args().First()
		if id == "" {
			return errors.New("id argument must be provided")
		}

		ctx := context.Background()
		registryClient, err := providerregistrysdk.NewClientWithResponses(c.String("registryAPIURL"))
		if err != nil {
			return errors.New("error configuring provider registry client")
		}

		keys := strings.Split(id, "@")

		if len(keys) != 2 {
			return errors.New("incorrect provider id given")
		}

		version := keys[1]

		teamAndName := strings.Split(keys[0], "/")
		if len(teamAndName) != 2 {
			return errors.New("incorrect provider id given")
		}

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

		//get bootstrap bucket

		//read from flag
		bootstrapBucket := c.String("bootstrapBucket")

		//work out the lambda asset path
		lambdaAssetPath := path.Join(team, name, version)

		//copy the provider assets into the bucket (this will also copy the cloudformation template too)
		awsCfg, err := aws_config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}
		client := s3.NewFromConfig(awsCfg)
		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bootstrapBucket),
			Key:        aws.String(lambdaAssetPath),
			CopySource: aws.String(strings.TrimPrefix(res.JSON200.LambdaAssetS3Arn, "arn:aws:s3:::")),
		})
		if err != nil {
			return err
		}

		clio.Log(fmt.Sprintf("copied %s into %s", id, path.Join(bootstrapBucket, lambdaAssetPath)))
		return nil
	},
}
