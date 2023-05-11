package deploy

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/common-fate/common-fate/pkg/config"
	"github.com/magefile/mage/sh"
	"github.com/segmentio/ksuid"
	"github.com/sethvargo/go-retry"
	"go.uber.org/zap"
)

//go:embed templates
var templateFiles embed.FS

func (o Output) ToRenderFrontendConfig() RenderFrontendConfig {
	return RenderFrontendConfig{
		Region:          o.Region,
		UserPoolID:      o.UserPoolID,
		CognitoClientID: o.CognitoClientID,
		UserPoolDomain:  o.UserPoolDomain,
		FrontendDomain:  o.FrontendDomainOutput,
		APIURL:          o.APIURL,
		CLIAppClientID:  o.CLIAppClientID,
	}
}

// WriteAWSExports writes aws exports files
func (o Output) WriteAWSExports() error {
	frontendConfig := o.ToRenderFrontendConfig()

	// first, render the application config JSON so that it can be uploaded to S3
	awsExportsProd, err := RenderProductionFrontendConfig(frontendConfig)
	if err != nil {
		return err
	}

	prodPath := "web/public/aws-exports.json"
	err = os.WriteFile(prodPath, []byte(awsExportsProd), 0666)
	if err != nil {
		return err
	}

	zap.S().Infow("wrote production frontend config", "path", prodPath)

	_, err = os.Stat("web/dist")
	if err == nil {
		// also write the application config JSON to the NextJS output folder,
		// in case the variables have changed since the last frontend build.
		outPath := "web/dist/aws-exports.json"
		err = os.WriteFile(outPath, []byte(awsExportsProd), 0666)
		if err != nil {
			return err
		}
		zap.S().Infow("wrote production frontend config", "path", outPath)
	}

	// also render the local development config JSON
	awsExportsDev, err := RenderLocalFrontendConfig(frontendConfig)
	if err != nil {
		return err
	}

	devPath := "web/src/utils/aws-exports.js"
	err = os.WriteFile(devPath, []byte(awsExportsDev), 0666)
	if err != nil {
		return err
	}

	zap.S().Infow("wrote development frontend config", "path", devPath)
	return nil
}

// DeployFrontend uploads the frontend to S3 and invalidates CloudFront
func (o Output) DeployFrontend() error {
	err := o.WriteAWSExports()
	if err != nil {
		return err
	}

	zap.S().Infow("clearing old resources", "bucket", o.S3BucketName)
	err = sh.Run("aws", "s3", "rm", fmt.Sprintf("s3://%s/", o.S3BucketName), "--recursive")
	if err != nil {
		return err
	}

	zap.S().Infow("uploading to s3", "bucket", o.S3BucketName)
	err = sh.Run("aws", "s3", "cp", "./web/dist", fmt.Sprintf("s3://%s/", o.S3BucketName), "--recursive")
	if err != nil {
		return err
	}

	zap.S().Infow("invalidating cloudfront distribution", "distributionId", o.CloudFrontDistributionID)
	err = sh.Run("aws", "cloudfront", "create-invalidation", "--distribution-id", o.CloudFrontDistributionID, "--paths", "/*")
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s", o.FrontendDomainOutput)
	zap.S().Infow("deployed frontend", "url", url)
	return nil
}

// RenderFrontendConfig contains all the required mappings for the templates
type RenderFrontendConfig struct {
	Region          string
	UserPoolID      string
	CognitoClientID string
	UserPoolDomain  string
	FrontendDomain  string
	APIURL          string
	CLIAppClientID  string
}

// RenderLocalFrontendConfig renders the aws-exports.js file
// to be used in local development.
// This accepts a specific config so this function can be reused easily
func RenderLocalFrontendConfig(rfc RenderFrontendConfig) (string, error) {
	tmpl, err := template.ParseFS(templateFiles, "templates/*")
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, "aws-exports.js.tmpl", rfc)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderProductionFrontendConfig renders the aws-exports.json file
// to be used in a production deployment of the frontend to AWS S3
// This accepts a specific config so this function can be reused easily in a custom resource lamda
func RenderProductionFrontendConfig(rfc RenderFrontendConfig) (string, error) {
	tmpl, err := template.ParseFS(templateFiles, "templates/*")
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, "aws-exports.json.tmpl", rfc)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func DeployProductionFrontend(ctx context.Context, cfg config.FrontendDeployerConfig) error {
	defaultAwsConfig, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return err
	}
	defaultS3Client := s3.NewFromConfig(defaultAwsConfig)

	rfc := RenderFrontendConfig{
		Region:          cfg.Region,
		UserPoolID:      cfg.UserPoolID,
		CognitoClientID: cfg.CognitoClientID,
		UserPoolDomain:  cfg.UserPoolDomain,
		FrontendDomain:  cfg.FrontendDomain,
		APIURL:          cfg.APIURL,
		CLIAppClientID:  cfg.CLIAppClientID,
	}

	zap.S().Infow("rendered frontend config", "config", rfc)

	// Read this first and fail before deleting any webapp files
	awsExports, err := RenderProductionFrontendConfig(rfc)
	if err != nil {
		return err
	}

	// delete contents of frontendBucket so it can be overwritten
	// how can we have a 0 downtime deployment?
	// we can't if we need to delete and upload the frontend.
	// might need to use 2 buckets or push frontend updates to prefixed paths in s3 then switch the cloudfront target only once the update is sucessful in a dependent step
	err = paginatedListObjects(ctx, defaultS3Client, &s3.ListObjectsV2Input{
		Bucket: &cfg.FrontendBucket,
	}, func(objects []types.Object) error {
		toDelete := make([]types.ObjectIdentifier, len(objects))
		for i, o := range objects {
			toDelete[i] = types.ObjectIdentifier{
				Key: o.Key,
			}
		}
		_, err := defaultS3Client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: &cfg.FrontendBucket,
			Delete: &types.Delete{
				Objects: toDelete,
				Quiet:   true,
			},
		})
		return err
	})
	if err != nil {
		return err
	}

	zap.S().Infow("deleted old objects", "bucket", cfg.FrontendBucket)

	err = paginatedListObjects(ctx, defaultS3Client, &s3.ListObjectsV2Input{
		Bucket: &cfg.CFReleasesBucket,
		Prefix: &cfg.CFReleasesFrontendAssetsObjectPrefix}, func(objects []types.Object) error {
		for _, o := range objects {
			copySource := cfg.CFReleasesBucket + "/" + aws.ToString(o.Key)
			destKey := strings.TrimPrefix(*o.Key, cfg.CFReleasesFrontendAssetsObjectPrefix+"/")
			// Could possibly use Go routines to split this up, not sure of the effectiveness in a lambda env though
			_, err = defaultS3Client.CopyObject(ctx, &s3.CopyObjectInput{
				Bucket:     &cfg.FrontendBucket,
				CopySource: &copySource,
				Key:        &destKey,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	zap.S().Infow("copied new objects", "bucket", cfg.FrontendBucket)

	// Finally, upload the frontend config
	key := "aws-exports.json"
	_, err = defaultS3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &cfg.FrontendBucket,
		Key:    &key,
		Body:   strings.NewReader(awsExports),
	})
	if err != nil {
		return err
	}

	// Invalidate the distribution cache so the new files are used
	cfClient := cloudfront.NewFromConfig(defaultAwsConfig)
	var attempt int
	// https://github.com/aws/aws-cdk/issues/15891#issuecomment-966456154
	// See this issue which states that the cloudfront API can fail randomly during peak times due to internal AWS API limits.
	// so try invalidating for up to 5 minutes with a backoff to ensure it has the best chance of working
	b := retry.WithMaxDuration(time.Minute*5, retry.NewFibonacci(time.Second))
	err = retry.Do(ctx, b, func(ctx context.Context) (err error) {
		attempt += 1
		// See the aws docs for caller reference, it just needs to be unique for every invalidation request
		callerReference := ksuid.New().String()
		res, err := cfClient.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
			DistributionId: &cfg.CloudFrontDistributionID,
			InvalidationBatch: &cfTypes.InvalidationBatch{
				CallerReference: &callerReference,
				Paths: &cfTypes.Paths{
					Quantity: aws.Int32(1),
					Items:    []string{"/*"},
				},
			},
		})
		if err != nil {
			zap.S().Errorw("failed to invalidate cloudfront distribution due to an AWS API error, retrying", "attempt", attempt, "error", err)
			return retry.RetryableError(err)
		}
		zap.S().Infow("successfully invalidated cloudfront distribution", "attempt", attempt, "cloudfront.id", cfg.CloudFrontDistributionID, "invalidation.id", res.Invalidation.Id)
		return nil
	})
	return err
}

func paginatedListObjects(ctx context.Context, client *s3.Client, in *s3.ListObjectsV2Input, callback func(objects []types.Object) error) error {
	var nextToken *string
	hasMore := true
	for hasMore {
		in.ContinuationToken = nextToken
		objectsRes, err := client.ListObjectsV2(ctx, in)
		if err != nil {
			return err
		}
		hasMore = objectsRes.IsTruncated
		nextToken = objectsRes.NextContinuationToken
		if len(objectsRes.Contents) > 0 {
			err = callback(objectsRes.Contents)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
