package middleware

import (
	"errors"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.com/briandowns/spinner"
	"github.com/common-fate/clio/clierr"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/urfave/cli/v2"
)

// RequireAWSCredentials attempts to load aws credentials, if they don't exist, iot returns a clio.CLIError
// This function will set the AWS config in context under the key cfaws.AWSConfigContextKey
// use cfaws.ConfigFromContextOrDefault(ctx) to retrieve the value
func RequireAWSCredentials() cli.BeforeFunc {
	return func(c *cli.Context) error {
		ctx := c.Context
		si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		si.Suffix = " loading AWS credentials from your current profile"
		si.Writer = os.Stderr
		si.Start()
		defer si.Stop()
		needCredentialsLog := clierr.Info(`Please export valid AWS credentials to run this command.
For more information see:
https://docs.commonfate.io/granted-approvals/troubleshooting/aws-credentials
`)
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return clierr.New("Failed to load AWS credentials.", clierr.Debugf("Encountered error while loading default aws config: %s", err), needCredentialsLog)
		}

		creds, err := cfg.Credentials.Retrieve(ctx)
		if err != nil {
			return clierr.New("Failed to load AWS credentials.", clierr.Debugf("Encountered error while loading default aws config: %s", err), needCredentialsLog)
		}

		if !creds.HasKeys() {
			return clierr.New("Failed to load AWS credentials.", needCredentialsLog)
		}

		stsClient := sts.NewFromConfig(cfg)
		// Use the sts api to check if these credentials are valid
		_, err = stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			var ae smithy.APIError
			// the aws sdk doesn't seem to have a concrete type for ExpiredToken so instead we check the error code
			if errors.As(err, &ae) && ae.ErrorCode() == "ExpiredToken" {
				return clierr.New("AWS credentials are expired.", needCredentialsLog)
			}
			return clierr.New("Failed to call AWS get caller identity. ", clierr.Debug(err.Error()), needCredentialsLog)
		}

		c.Context = cfaws.SetConfigInContext(ctx, cfg)
		return nil
	}
}
