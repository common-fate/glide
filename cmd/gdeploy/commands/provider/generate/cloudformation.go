package generate

import (
	"context"
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/fmtconvert"
	"github.com/common-fate/common-fate/pkg/ssmkey"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	registryclient "github.com/common-fate/provider-registry-sdk-go/pkg/registryclient"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type Presigner struct {
	PresignClient *s3.PresignClient
}

// GetObject makes a presigned request that can be used to get an object from a bucket.
// The presigned request is valid for the specified number of seconds.
func (presigner Presigner) GetObject(
	bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs)
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

func convertValuesToCloudformationParameter(m map[string]string) string {
	parameters := "--parameters "

	for k, v := range m {
		parameters = parameters + strings.Join([]string{fmt.Sprintf("ParameterKey=\"%s\"", k), ",", fmt.Sprintf("ParameterValue=\"%s\"", v)}, "") + " "
	}

	return parameters
}

var cloudFormationCreate = cli.Command{
	Name:    "cloudformation-create",
	Aliases: []string{"cfn-create"},
	Usage:   "Generate an 'aws cloudformation create-stack' command",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "provider-id", Required: true, Usage: "publisher/name@version"},
		&cli.StringFlag{Name: "handler-id", Required: true, Usage: "The ID of the Handler (for example, 'cf-handler-aws')"},
		&cli.StringFlag{Name: "bootstrap-bucket", Required: true},
		&cli.StringFlag{Name: "common-fate-aws-account", Usage: "The AWS account where Common Fate is deployed"},
		&cli.StringFlag{Name: "region", Usage: "The region to deploy the handler"},
	},
	Action: func(c *cli.Context) error {
		ctx := c.Context
		bootstrapBucket := c.String("bootstrap-bucket")
		handlerID := c.String("handler-id")
		commonFateAWSAccountID := c.String("common-fate-aws-account")
		registry, err := registryclient.New(ctx)
		if err != nil {
			return errors.Wrap(err, "configuring registry client")
		}

		providerString := c.String("provider-id")
		provider, err := providerregistrysdk.ParseProvider(providerString)
		if err != nil {
			return err
		}

		// check that the provider type matches one in our registry
		res, err := registry.GetProviderWithResponse(ctx, provider.Publisher, provider.Name, provider.Version)
		if err != nil {
			return err
		}

		var stackname = c.String("handler-id")
		if stackname == "" {
			err = survey.AskOne(&survey.Input{Message: "enter the cloudformation stackname:", Default: handlerID}, &stackname)
			if err != nil {
				return err
			}
		}

		var region = c.String("region")
		if region == "" {
			err = survey.AskOne(&survey.Input{Message: "enter the region of cloudformation stack deployment"}, &region)
			if err != nil {
				return err
			}
		}

		values := make(map[string]string)

		values["BootstrapBucketName"] = bootstrapBucket
		values["HandlerID"] = handlerID
		values["CommonFateAWSAccountID"] = commonFateAWSAccountID
		lambdaAssetPath := path.Join("registry.commonfate.io", "v1alpha1", "providers", provider.Publisher, provider.Name, provider.Version)
		values["AssetPath"] = path.Join(lambdaAssetPath, "handler.zip")

		awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
		if err != nil {
			return err
		}

		config := res.JSON200.Schema.Config
		if config != nil {
			clio.Info("Enter the values for your configurations:")
			for k, v := range *config {
				if v.Secret != nil && *v.Secret {
					client := ssm.NewFromConfig(awsCfg)

					var secret string
					name := ssmkey.SSMKey(ssmkey.SSMKeyOpts{
						HandlerID:    handlerID,
						Key:          k,
						Publisher:    provider.Publisher,
						ProviderName: provider.Name,
					})

					helpMsg := fmt.Sprintf("This will be stored in AWS SSM Parameter Store with name '%s'", name)
					err = survey.AskOne(&survey.Password{Message: k + ":", Help: helpMsg}, &secret)
					if err != nil {
						return err
					}

					_, err = client.PutParameter(ctx, &ssm.PutParameterInput{
						Name:      aws.String(name),
						Value:     aws.String(secret),
						Type:      types.ParameterTypeSecureString,
						Overwrite: aws.Bool(true),
					})
					if err != nil {
						return err
					}

					clio.Successf("Added to AWS SSM Parameter Store with name '%s'", name)

					// secret config should have "Secret" prefix to the config key name.
					values[fmtconvert.PascalCase(k)+"Secret"] = name

					continue
				}

				var v string
				err = survey.AskOne(&survey.Input{Message: k + ":"}, &v)
				if err != nil {
					return err
				}
				values[fmtconvert.PascalCase(k)] = v

			}
		}

		if commonFateAWSAccountID == "" {
			var v string
			err = survey.AskOne(&survey.Input{Message: "The ID of the AWS account where Common Fate is deployed:"}, &v)
			if err != nil {
				return err
			}
			values["CommonFateAWSAccountID"] = v
		}

		parameterKeys := convertValuesToCloudformationParameter(values)

		s3client := s3.NewFromConfig(awsCfg)
		preSignedClient := s3.NewPresignClient(s3client)

		presigner := Presigner{
			PresignClient: preSignedClient,
		}

		req, err := presigner.GetObject(bootstrapBucket, path.Join(lambdaAssetPath, "cloudformation.json"), int64(time.Hour))
		if err != nil {
			return nil
		}

		templateUrl := fmt.Sprintf(" --template-url \"%s\" ", req.URL)
		stackNameFlag := fmt.Sprintf(" --stack-name %s ", stackname)
		regionFlag := fmt.Sprintf(" --region %s ", region)

		output := strings.Join([]string{"aws cloudformation create-stack", stackNameFlag, regionFlag, templateUrl, parameterKeys, "--capabilities CAPABILITY_NAMED_IAM"}, "")

		fmt.Printf("%v \n", output)

		return nil
	},
}

var cloudformationUpdate = cli.Command{
	Name:    "cloudformation-update",
	Aliases: []string{"cfn-update"},
	Usage:   "Generate an 'aws cloudformation update-stack' command",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "handler-id", Usage: "The Handler ID and name of the CloudFormation stack", Required: true},
		&cli.StringFlag{Name: "region", Usage: "The region to deploy the handler", Required: true},
		&cli.StringFlag{Name: "provider-id", Usage: "Update the provider-id for the current stack"},
		&cli.BoolFlag{Name: "use-previous-value", Usage: "use the previous stack values for the parameters"},
	},
	Action: func(c *cli.Context) error {
		stackname := c.String("handler-id")
		region := c.String("region")

		ctx := c.Context

		awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
		if err != nil {
			return err
		}

		cfn := cloudformation.NewFromConfig(awsCfg)

		out, err := cfn.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
			StackName: &stackname,
		})
		if err != nil {
			return err
		}

		if len(out.Stacks) > 0 {
			stack := out.Stacks[0]
			values := make(map[string]string)

			providerID := c.String("provider-id")

			// if the provider-id is provided then update the lambda-assets-handler path to the new version.
			if providerID != "" {
				registry, err := registryclient.New(ctx)
				if err != nil {
					return errors.Wrap(err, "configuring provider registry client")
				}

				provider, err := providerregistrysdk.ParseProvider(providerID)
				if err != nil {
					return err
				}

				_, err = registry.GetProviderWithResponse(ctx, provider.Publisher, provider.Name, provider.Version)
				if err != nil {
					return err
				}

				values["AssetPath"] = path.Join("registry.commonfate.io", "v1alpha1", "providers", provider.Publisher, provider.Name, provider.Version, "handler.zip")
			}

			for _, parameter := range stack.Parameters {

				// update-stack shouldn't ask to update the handler-id, bootstrapBucketName
				// if the use-previous-value flag is provided then don't prompt users to reconfigure the required parameters
				if contains([]string{"HandlerID", "BootstrapBucketName"}, *parameter.ParameterKey) || c.Bool("use-previous-value") {
					values[*parameter.ParameterKey] = *parameter.ParameterValue
					continue
				}

				// secret values have this prefix so need to update the SSM parameter store for these keys
				if strings.HasPrefix(*parameter.ParameterValue, "awsssm:///common-fate/provider/") {
					var shouldUpdate bool

					err = survey.AskOne(&survey.Confirm{Message: "Do you want to update value for " + *parameter.ParameterKey + " in AWS parameter store?"}, &shouldUpdate)
					if err != nil {
						return err
					}

					if shouldUpdate {
						client := ssm.NewFromConfig(awsCfg)

						var secret string
						name := *parameter.ParameterValue
						helpMsg := fmt.Sprintf("This will be stored in aws system manager parameter store with name '%s'", name)
						err = survey.AskOne(&survey.Password{Message: *parameter.ParameterKey + ":", Help: helpMsg}, &secret)
						if err != nil {
							return err
						}

						_, err = client.PutParameter(ctx, &ssm.PutParameterInput{
							Name:      aws.String(name),
							Value:     aws.String(secret),
							Type:      types.ParameterTypeSecureString,
							Overwrite: aws.Bool(true),
						})
						if err != nil {
							return err
						}

						clio.Successf("Updated value in AWS System Manager Parameter Store for key with name '%s'", name)
					}

					continue
				}

				// overwrite the value with provided provider-id
				if c.String("provider-id") != "" && *parameter.ParameterKey == "AssetPath" {
					continue
				}

				var v string
				err = survey.AskOne(&survey.Input{Message: *parameter.ParameterKey + ":", Default: *parameter.ParameterValue}, &v)
				if err != nil {
					return err
				}

				if v != *parameter.ParameterValue {
					values[*parameter.ParameterKey] = v
				} else {
					values[*parameter.ParameterKey] = *parameter.ParameterValue
				}
			}

			parameterKeys := convertValuesToCloudformationParameter(values)

			s3client := s3.NewFromConfig(awsCfg)
			preSignedClient := s3.NewPresignClient(s3client)

			presigner := Presigner{
				PresignClient: preSignedClient,
			}

			bootstrapBucket := values["BootstrapBucketName"]
			lambdaAssetPath := values["AssetPath"]

			req, err := presigner.GetObject(bootstrapBucket, strings.Replace(lambdaAssetPath, "handler.zip", "cloudformation.json", 1), int64(time.Hour)*1)
			if err != nil {
				return nil
			}

			templateUrl := fmt.Sprintf(" --template-url \"%s\" ", req.URL)
			stackNameFlag := fmt.Sprintf(" --stack-name %s ", stackname)
			regionFlag := fmt.Sprintf(" --region %s ", region)

			output := strings.Join([]string{"aws cloudformation update-stack", stackNameFlag, regionFlag, templateUrl, parameterKeys, "--capabilities CAPABILITY_NAMED_IAM"}, "")

			fmt.Printf("%v \n", output)
		}

		return nil
	},
}

// utility function to check if the string belongs to the slice.
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
