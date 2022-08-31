package middleware

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.com/briandowns/spinner"
	"github.com/common-fate/granted-approvals/pkg/cfaws"
	"github.com/common-fate/granted-approvals/pkg/clio"
	"github.com/common-fate/granted-approvals/pkg/deploy"
	"github.com/urfave/cli/v2"
)

// RequireAWSCredentials attempts to load aws credentials, if they don't exist, iot returns a clio.CLIError
// This function will set the AWS config in context under the key cfaws.AWSConfigContextKey
// use cfaws.ConfigFromContextOrDefault(ctx) to retrieve the value
// If RequireDeploymentConfig has already run, this function will use the region value from the deployment config when setting the AWS config in context
func RequireAWSCredentials() cli.BeforeFunc {
	return func(c *cli.Context) error {
		ctx := c.Context
		si := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		si.Suffix = " loading AWS credentials from your current profile"
		si.Writer = os.Stderr
		si.Start()
		defer si.Stop()
		needCredentialsLog := clio.LogMsg("Please export valid AWS credentials to run this command.")
		cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
		if err != nil {
			return clio.NewCLIError("Failed to load AWS credentials.", clio.DebugMsg("Encountered error while loading default aws config: %s", err), needCredentialsLog)
		}

		// Use the deployment region if it is available
		var configExists bool
		dc, err := deploy.ConfigFromContext(ctx)
		if err == nil {
			configExists = true
			if dc.Deployment.Region != "" {
				cfg.Region = dc.Deployment.Region
			}
			if dc.Deployment.Account != "" {
				// include the account id in the log message if available
				needCredentialsLog = clio.LogMsg("Please export valid AWS credentials for account %s to run this command.", dc.Deployment.Account)
			}
		}

		creds, err := cfg.Credentials.Retrieve(ctx)
		if err != nil {
			return clio.NewCLIError("Failed to load AWS credentials.", clio.DebugMsg("Encountered error while loading default aws config: %s", err), needCredentialsLog)
		}

		if !creds.HasKeys() {
			return clio.NewCLIError("Failed to load AWS credentials.", needCredentialsLog)
		}

		stsClient := sts.NewFromConfig(cfg)
		// Use the sts api to check if these credentials are valid
		identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			var ae smithy.APIError
			// the aws sdk doesn't seem to have a concrete type for ExpiredToken so instead we check the error code
			if errors.As(err, &ae) && ae.ErrorCode() == "ExpiredToken" {
				return clio.NewCLIError("AWS credentials are expired.", needCredentialsLog)
			}
			return clio.NewCLIError("Failed to call AWS get caller identity. ", clio.DebugMsg(err.Error()), needCredentialsLog)
		}

		//check to see that account number in config is the same account that is assumed
		if configExists && *identity.Account != dc.Deployment.Account {
			return clio.NewCLIError(fmt.Sprintf("AWS account in your deployment config %s does not match the account of your current AWS credentials %s", dc.Deployment.Account, *identity.Account), needCredentialsLog)
		}
		c.Context = cfaws.SetConfigInContext(ctx, cfg)
		return nil
	}
}
