package provider

import (
	"context"
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/common-fate/cloudform/deployer"

	"github.com/pkg/errors"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/cliconfig"
	"github.com/common-fate/common-fate/pkg/client"
	"github.com/common-fate/common-fate/pkg/fmtconvert"
	"github.com/common-fate/common-fate/pkg/prompt"
	"github.com/common-fate/common-fate/pkg/ssmkey"
	cftypes "github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/provider-registry-sdk-go/pkg/bootstrapper"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	registryclient "github.com/common-fate/provider-registry-sdk-go/pkg/registryclient"
	"github.com/sethvargo/go-retry"
	"github.com/urfave/cli/v2"
)

var deployCommand = cli.Command{
	Name:        "deploy",
	Description: "Quickstart command to deploy a provider",
	Usage:       "Quickstart command to deploy a provider",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "provider", Aliases: []string{"p"}, Usage: "The provider to deploy (for example, 'common-fate/aws@v0.4.0')"},
		&cli.StringFlag{Name: "handler-id", Usage: "The Handler ID and CloudFormation stack name to use (by convention, this is 'cf-handler-[provider publisher]-[provider name]')"},
		&cli.StringFlag{Name: "target-group-id", Usage: "Override the ID of the Target Group which will be created"},
		&cli.StringFlag{Name: "common-fate-aws-account", Usage: "Override the Common Fate AWS Account ID (by default the current AWS account ID is used)"},
		&cli.StringFlag{Name: "target", Aliases: []string{"t"}, Usage: "The target kind to use with the provider (only required if the provider grants access to multiple kinds of targets)"},
		&cli.BoolFlag{Name: "confirm-bootstrap", Usage: "Confirm creating a bootstrap bucket if it doesn't exist in the account and region you are deploying to"},
		&cli.StringSliceFlag{Name: "config", Usage: "Provide config values for the provider in key=value format"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context

		awsConfig, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return err
		}

		bs := bootstrapper.NewFromConfig(awsConfig)

		cfg, err := cliconfig.Load()
		if err != nil {
			return err
		}

		registry, err := registryclient.New(ctx)
		if err != nil {
			return errors.Wrap(err, "configuring provider registry client")
		}

		var provider *providerregistrysdk.ProviderDetail

		// validate this as early as possible
		providerArg := c.String("provider")
		if providerArg != "" {
			p, err := providerregistrysdk.ParseProvider(providerArg)
			if err != nil {
				return err
			}
			clio.Infof("Retrieving provider details for '%s' from the Provider Registry...", providerArg)
			res, err := registry.GetProviderWithResponse(ctx, p.Publisher, p.Name, p.Version)
			if err != nil {
				return err
			}
			provider = res.JSON200
		}

		configArgs := map[string]string{}

		for _, arg := range c.StringSlice("config") {
			parts := strings.SplitN(arg, "=", 2) // args are in key=value format

			if len(parts) != 2 {
				return fmt.Errorf("invalid config argument (expected format is --config key=value): %s", arg)
			}

			key := parts[0]
			val := parts[1]
			configArgs[key] = val
		}

		// the client needs to be constructed as early as possible in the
		// CLI command, because client.FromConfig() returns an error
		// prompting the user to run 'cf login' if they are unauthenticated.
		cf, err := client.FromConfig(ctx, cfg)
		if err != nil {
			return err
		}

		// make an admin API call. Even though we don't use the response,
		// this will cause the CLI wizard to fail early if the auth token
		// is expired, or if the user is not an administrator.
		_, err = cf.AdminListHandlersWithResponse(ctx)
		if err != nil {
			return err
		}

		// find the AWS account ID as early as possible, as it will return
		// an error if credentials are expired.
		stsClient := sts.NewFromConfig(awsConfig)
		// Use the sts api to check if these credentials are valid
		stsOut, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			return err
		}

		awsAccount := *stsOut.Account

		cfAccountID := c.String("common-fate-aws-account")
		if cfAccountID == "" {
			clio.Warnf("using the current AWS account (%s) as the Common Fate account (use --common-fate-aws-account to override)", awsAccount)
			cfAccountID = awsAccount
		}

		if provider == nil {
			provider, err = prompt.Provider(ctx, registry)
			if err != nil {
				return err
			}
		}

		bootstrapStackOutput, err := bs.GetOrDeployBootstrapBucket(ctx, deployer.WithConfirm(c.Bool("confirm-bootstrap")))
		if err != nil {
			return err
		}

		selectedProviderKind, err := prompt.Kind(*provider)
		if err != nil {
			return err
		}

		clio.Info("Copying provider assets from the registry to the bootstrap bucket...")
		err = bs.CopyProviderFiles(ctx, *provider)
		if err != nil {
			return err
		}
		clio.Success("Provider assets copied to the bootstrap bucket")

		handlerID := c.String("handler-id")
		if handlerID == "" {
			handlerID = strings.Join([]string{"cf-handler", provider.Publisher, provider.Name}, "-")
		}

		var uniqueHandlerIDFound bool

		for !uniqueHandlerIDFound {
			// check if a lambda role already exists with the given ID
			exists, err := checkIfLambdaRoleExists(ctx, awsConfig, handlerID)
			if err != nil {
				return err
			}
			if !exists {
				uniqueHandlerIDFound = true
			} else {
				clio.Warnf("A Lambda function named '%s' already exists in the account. You will need to set a custom Handler ID.\nBy convention, we use 'cf-handler-[publisher]-[name]-[suffix]' as Handler IDs, for example: 'cf-handler-common-fate-aws-dev'.", handlerID)
				err = survey.AskOne(&survey.Input{Message: "Unique Handler ID:"}, &handlerID)
				if err != nil {
					return err
				}
			}
		}

		lambdaAssetPath := path.Join("registry.commonfate.io", "v1alpha1", "providers", provider.Publisher, provider.Name, provider.Version)

		var oneLinerConfigArgs []string

		var parameters []types.Parameter

		config := provider.Schema.Config
		if config != nil {

			// sort keys alphabetically so they appear in a consistent order between CLI runs.
			var keys []string
			for k := range *config {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			clio.Info("This Provider requires configuration")
			for _, k := range keys {
				v := (*config)[k]

				paramName := fmtconvert.PascalCase(k)

				var isSecret bool
				if v.Secret != nil && *v.Secret {
					isSecret = true
					paramName += "Secret"
				}

				// the CloudFormation parameter value
				var paramVal string

				// check if it was provided as a CLI argument using --config key=value
				if flagVal, ok := configArgs[k]; ok {
					paramVal = flagVal
				} else {
					// prompt the user interactively for the config values
					paramVal, err = promptForConfig(ctx, promptForConfigOpts{
						AWSCfg:    awsConfig,
						IsSecret:  isSecret,
						Key:       k,
						HandlerID: handlerID,
						Provider:  *provider,
					})
					if err != nil {
						return err
					}
				}

				parameters = append(parameters, types.Parameter{
					ParameterKey:   &paramName,
					ParameterValue: &paramVal,
				})

				oneLinerConfigArgs = append(oneLinerConfigArgs, fmt.Sprintf("--config %s=%s", k, paramVal))

				clio.Infof("Setting CloudFormation parameter %s=%s", paramName, paramVal)
			}
		}

		parameters = append(parameters, types.Parameter{
			ParameterKey:   aws.String("CommonFateAWSAccountID"),
			ParameterValue: &cfAccountID,
		})

		parameters = append(parameters, types.Parameter{
			ParameterKey:   aws.String("AssetPath"),
			ParameterValue: aws.String(path.Join(lambdaAssetPath, "handler.zip")),
		})

		parameters = append(parameters, types.Parameter{
			ParameterKey:   aws.String("BootstrapBucketName"),
			ParameterValue: aws.String(bootstrapStackOutput.AssetsBucket),
		})

		parameters = append(parameters, types.Parameter{
			ParameterKey:   aws.String("HandlerID"),
			ParameterValue: aws.String(handlerID),
		})

		targetgroupID := c.String("target-group-id")
		if targetgroupID == "" {
			targetgroupID = strings.TrimPrefix(handlerID, "cf-handler-")
		}

		oneLinerCommand := fmt.Sprintf("cf provider deploy --common-fate-aws-account %s --handler-id %s --target-group-id %s --provider %s %s", cfAccountID, handlerID, targetgroupID, provider, strings.Join(oneLinerConfigArgs, " "))

		clio.NewLine()
		clio.Infof("You can use the following one-liner command to redeploy this Provider in future:\n%s", oneLinerCommand)
		clio.NewLine()

		d := deployer.NewFromConfig(awsConfig)

		clio.Infof("Deploying CloudFormation stack for Handler '%s'", handlerID)

		templateURL := bootstrapStackOutput.CloudFormationURL(provider.Base())

		out, err := d.Deploy(ctx, deployer.DeployOpts{
			Template:  templateURL,
			Params:    parameters,
			StackName: handlerID,
			Confirm:   true,
		})
		if err != nil {
			return err
		}

		// if the output of cloudformation deploy is not 'CREATE_COMPLETE'
		// then should return error
		if out.FinalStatus != "CREATE_COMPLETE" {
			return fmt.Errorf("failed to deploy CloudFormation stack for Handler '%s", handlerID)
		}

		clio.Infof("Deployment completed for HandlerID %s", handlerID)

		clio.Infof("Creating a Target Group '%s' to route Access Requests to the Handler", targetgroupID)

		_, err = cf.AdminCreateTargetGroupWithResponse(ctx, cftypes.AdminCreateTargetGroupJSONRequestBody{
			Id: targetgroupID,
			From: cftypes.TargetGroupFrom{
				Kind:      selectedProviderKind,
				Name:      provider.Name,
				Publisher: provider.Publisher,
				Version:   provider.Version,
			},
		})
		if err != nil {
			return err
		}
		clio.Successf("Target Group created: %s", targetgroupID)

		// register the targetgroup with handler

		reqBody := cftypes.AdminRegisterHandlerJSONRequestBody{
			AwsAccount: awsAccount,
			AwsRegion:  awsConfig.Region,
			Runtime:    "aws-lambda",
			Id:         handlerID,
		}

		_, err = cf.AdminRegisterHandlerWithResponse(ctx, reqBody)
		if err != nil {
			return err
		}

		clio.Successf("Successfully registered Handler '%s' with Common Fate", handlerID)

		_, err = cf.AdminCreateTargetGroupLinkWithResponse(ctx, targetgroupID, cftypes.AdminCreateTargetGroupLinkJSONRequestBody{
			DeploymentId: handlerID,
			Priority:     100,
			Kind:         selectedProviderKind,
		})
		if err != nil {
			return err
		}

		clio.Successf("Successfully linked Handler '%s' with Target Group '%s'", handlerID, targetgroupID)

		clio.Info("Waiting for Handler to become healthy...")

		// retry every 5 seconds for a maximum of two minutes
		err = retry.Do(ctx, retry.WithMaxDuration(time.Minute*2, retry.NewConstant(time.Second*5)), func(ctx context.Context) error {
			ghr, err := cf.AdminGetHandlerWithResponse(ctx, handlerID)
			if err != nil && ghr.StatusCode() < 500 {
				return retry.RetryableError(err)
			}
			if err != nil {
				return err
			}
			if ghr.JSON200.Healthy {
				clio.Successf("Handler '%s' is healthy", handlerID)
				return nil
			}

			clio.Warnw("Handler is not healthy yet", "diagnostics", ghr.JSON200.Diagnostics)

			// the below error will be shown to the user if the time limit is exceeded
			return retry.RetryableError(errors.New("timed out waiting for Handler to become healthy"))
		})
		if err != nil {
			return err
		}

		return nil
	},
}

func checkIfLambdaRoleExists(ctx context.Context, cfg aws.Config, handlerID string) (exists bool, err error) {
	client := iam.NewFromConfig(cfg)
	_, err = client.GetRole(ctx, &iam.GetRoleInput{
		RoleName: &handlerID,
	})
	var rnf *iamTypes.NoSuchEntityException
	if errors.As(err, &rnf) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// if we get here, the role exists because the API call succeeded.
	return true, nil
}

type promptForConfigOpts struct {
	AWSCfg    aws.Config
	IsSecret  bool
	Key       string
	HandlerID string
	Provider  providerregistrysdk.ProviderDetail
}

// promptForConfig prompts a user interactively for a config value.
// If the value is a secret, it is written to SSM Parameter Store.
//
// For non-secret config values, the value is returned.
// For secrets, the secret reference is returned in the format "awsssm://<SECRET PATH>"
func promptForConfig(ctx context.Context, opts promptForConfigOpts) (string, error) {
	var paramVal string
	if !opts.IsSecret {
		// not a secret, so use the value directly
		err := survey.AskOne(&survey.Input{Message: opts.Key + ":"}, &paramVal)
		if err != nil {
			return "", err
		}
		return paramVal, nil
	}

	// if we get here, its a secret, so write it to SSM Parameter store and return the key
	client := ssm.NewFromConfig(opts.AWSCfg)

	var secret string
	ssmKey := ssmkey.SSMKey(ssmkey.SSMKeyOpts{
		HandlerID:    opts.HandlerID,
		Key:          opts.Key,
		Publisher:    opts.Provider.Publisher,
		ProviderName: opts.Provider.Name,
	})
	helpMsg := fmt.Sprintf("This will be stored in AWS SSM Parameter Store with name '%s'", ssmKey)
	err := survey.AskOne(&survey.Password{Message: opts.Key + ":", Help: helpMsg}, &secret)
	if err != nil {
		return "", err
	}

	_, err = client.PutParameter(ctx, &ssm.PutParameterInput{
		Name:      &ssmKey,
		Value:     &secret,
		Type:      ssmTypes.ParameterTypeSecureString,
		Overwrite: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}

	clio.Successf("Added to AWS SSM Parameter Store with name '%s'", ssmKey)
	paramVal = "awsssm://" + ssmKey
	return paramVal, nil
}
