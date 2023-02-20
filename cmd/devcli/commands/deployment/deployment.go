package deployment

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/briandowns/spinner"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"

	"github.com/AlecAivazis/survey/v2"
	aws_types "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/common-fate/clio"
	"github.com/common-fate/cloudform/cfn"
	"github.com/common-fate/cloudform/console"
	"github.com/common-fate/cloudform/ui"
	"github.com/common-fate/common-fate/cf/pkg/bootstrapper"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/common-fate/common-fate/pkg/types"

	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/urfave/cli/v2"
)

//eg command
//go run cmd/devcli/main.go deploy --runtime=aws --publisher=jack --version=v0.1.4 --name=testvault --accountId=12345678912 --aws-region=ap-southeast-2 --suffix=jacktest6

var Command = cli.Command{
	Name:        "deploy",
	Description: "make a new target group and deployment",
	Usage:       "make a new target group and deployment",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "publisher", Usage: "The publisher"},
		&cli.StringFlag{Name: "name", Usage: "The name"},
		&cli.StringFlag{Name: "version", Usage: "The version"},
		&cli.StringFlag{Name: "accountId", Usage: "The accountId"},
		&cli.StringFlag{Name: "aws-region", Usage: "The aws-region"},
		&cli.StringFlag{Name: "suffix", Usage: "The suffix"},
		&cli.StringFlag{Name: "runtime", Usage: "The runtime"},

		&cli.StringFlag{Name: "target-group-override", Usage: "If a target group already exists, pass the name here to use it"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		clio.Info("This command will setup a deployment of a provider\n")

		var providerPublisher string
		if c.String("publisher") != "" {
			providerPublisher = c.String("publisher")
		} else {
			p2 := &survey.Input{Message: "The publisher of the provider"}
			err := survey.AskOne(p2, &providerPublisher)
			if err != nil {
				return err
			}
		}

		var providerName string
		if c.String("name") != "" {
			providerName = c.String("name")
		} else {
			p2 := &survey.Input{Message: "The name of the provider"}
			err := survey.AskOne(p2, &providerName)
			if err != nil {
				return err
			}
		}

		var providerVersion string
		if c.String("version") != "" {
			providerVersion = c.String("version")
		} else {
			p2 := &survey.Input{Message: "The version of the provider"}
			err := survey.AskOne(p2, &providerVersion)
			if err != nil {
				return err
			}
		}

		var accountId string
		if c.String("accountId") != "" {
			accountId = c.String("accountId")
		} else {
			p2 := &survey.Input{Message: "The AWS account to deploy into"}
			err := survey.AskOne(p2, &accountId)
			if err != nil {
				return err
			}
		}

		var awsRegion string
		if c.String("aws-region") != "" {
			awsRegion = c.String("aws-region")
		} else {
			p2 := &survey.Input{Message: "The AWS region to deploy into"}
			err := survey.AskOne(p2, &awsRegion)
			if err != nil {
				return err
			}
		}

		var suffix string
		if c.String("suffix") != "" {
			suffix = c.String("suffix")
		} else {
			p2 := &survey.Input{Message: "A unique suffix for creating resources"}
			err := survey.AskOne(p2, &suffix)
			if err != nil {
				return err
			}
		}

		var cfg config.ProviderDeploymentCLI
		_ = godotenv.Load()

		err := envconfig.Process(ctx, &cfg)
		if err != nil {
			return err
		}
		assetString := strings.Join([]string{providerPublisher, providerName, providerVersion}, "/")

		clio.Success(fmt.Sprintf("Creating deployment for provider: %s", assetString))

		//bootstrap bucket
		bs, err := bootstrapper.New(ctx)
		if err != nil {
			return err
		}
		bootstrapBucket, err := bs.GetOrDeployBootstrapBucket(ctx)
		if err != nil {
			return err
		}

		//bootstrap provider

		registryClient, err := providerregistrysdk.NewClientWithResponses(cfg.ProviderRegistryAPIURL)
		if err != nil {
			return errors.New("error configuring provider registry client")
		}

		//check that the provider type matches one in our registry
		res, err := registryClient.GetProviderWithResponse(ctx, providerPublisher, providerName, providerVersion)
		if err != nil {
			return err
		}
		if res.StatusCode() != http.StatusOK {
			return fmt.Errorf("provider for that version does not exist: %s", assetString)
		}

		//copy the provider assets into the bucket (this will also copy the cloudformation template too)
		awsCfg, err := aws_config.LoadDefaultConfig(ctx)
		if err != nil {
			return err
		}
		client := s3.NewFromConfig(awsCfg)
		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bootstrapBucket),
			Key:        aws.String(path.Join(assetString, "handler.zip")),
			CopySource: aws.String(url.QueryEscape(path.Join(res.JSON200.LambdaAssetS3Arn, "handler.zip"))),
		})
		if err != nil {
			return err
		}

		_, err = client.CopyObject(ctx, &s3.CopyObjectInput{
			Bucket:     aws.String(bootstrapBucket),
			Key:        aws.String(path.Join(assetString, "cloudformation.json")),
			CopySource: aws.String(url.QueryEscape(path.Join(res.JSON200.LambdaAssetS3Arn, "cloudformation.json"))),
		})
		if err != nil {
			return err
		}

		clio.Success(fmt.Sprintf("copied %s into %s\n", assetString, path.Join(bootstrapBucket, assetString)))

		//lambda that is created from the cloudformation should have the same name of deployment we register below
		deploymentName := providerName + "-deployment" + "-" + suffix

		template := "https://" + bootstrapBucket + ".s3." + awsRegion + ".amazonaws.com/" + assetString + "/cloudformation.json"
		clio.Info(template)
		ccfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}
		cfnClient := cfn.New(ccfg)

		si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		si.Suffix = " creating CloudFormation change set"
		si.Writer = os.Stderr
		si.Start()

		//todo update this from being hardcoded parameters for testvault
		params := []aws_types.Parameter{
			{
				ParameterKey:   aws.String("UniqueVaultId"),
				ParameterValue: aws.String("2CBsuomHFRE3mrpLGWFaxbyKXG6"),
			},
			{
				ParameterKey:   aws.String("Id"),
				ParameterValue: &deploymentName,
			},
			{
				ParameterKey:   aws.String("BootstrapBucketName"),
				ParameterValue: &bootstrapBucket,
			},
			{
				ParameterKey:   aws.String("AssetPath"),
				ParameterValue: aws.String(path.Join(providerPublisher, providerName, providerVersion, "handler.zip")),
			},
			{
				ParameterKey:   aws.String("ApiUrl"),
				ParameterValue: aws.String("https://prod.testvault.granted.run"),
			},
			{
				ParameterKey:   aws.String("AccountId"),
				ParameterValue: &accountId,
			},
		}

		stackName := providerName + suffix

		changeSetName, createErr := cfnClient.CreateChangeSet(ctx, template, params, nil, stackName, "")

		si.Stop()

		if createErr != nil {

			return createErr

		}

		uiClient := ui.New(ccfg)

		status, err := uiClient.FormatChangeSet(ctx, stackName, changeSetName)
		if err != nil {
			return err
		}
		clio.Info("The following CloudFormation changes will be made:")
		fmt.Println(status)

		err = cfnClient.ExecuteChangeSet(ctx, stackName, changeSetName)
		if err != nil {
			return err
		}

		status, messages := uiClient.WaitForStackToSettle(ctx, stackName)

		fmt.Println("Final stack status:", ui.ColouriseStatus(status))

		if len(messages) > 0 {
			fmt.Println(console.Yellow("Messages:"))
			for _, message := range messages {
				fmt.Printf("  - %s\n", message)
			}
		}

		//create target group

		cfApi, err := types.NewClientWithResponses(cfg.CommonFateAPIURL)
		if err != nil {
			return err
		}

		targetGroupId := c.String("target-group-override")

		//if there was an override dont create a new target group but link with an old one
		if targetGroupId == "" {
			tgCreateReq := types.AdminCreateTargetGroupJSONRequestBody{
				ID: providerName + "-" + providerVersion + "-" + suffix,
				// will create the target group as eg. commonfate-v0.1.4-jacktest
				TargetSchema: assetString,
			}
			_, err = cfApi.AdminCreateTargetGroupWithResponse(ctx, tgCreateReq)
			if err != nil {
				return err
			}
			clio.Successf("Successfully created target group '%s'", tgCreateReq.ID)
			targetGroupId = tgCreateReq.ID

		}

		reqBody := types.AdminCreateTargetGroupDeploymentJSONRequestBody{
			AwsAccount: accountId,
			AwsRegion:  awsRegion,
			Id:         deploymentName,
			Runtime:    c.String("runtime"),
		}

		_, err = cfApi.AdminCreateTargetGroupDeploymentWithResponse(ctx, reqBody)
		if err != nil {
			return err
		}
		clio.Successf("Successfully created target group deployment '%s'", reqBody.Id)

		//link deployment to target group
		_, err = cfApi.AdminCreateTargetGroupLinkWithResponse(ctx, targetGroupId, types.AdminCreateTargetGroupLinkJSONRequestBody{
			DeploymentId: reqBody.Id,
			Priority:     100,
		})
		if err != nil {
			return err
		}
		clio.Successf("linked deployment '%s' with target group '%s'", reqBody.Id, targetGroupId)

		//run health check
		clio.Successf("Completed deploy")

		return nil
	},
}
