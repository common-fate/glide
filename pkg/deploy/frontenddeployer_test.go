package deploy

import (
	"context"
	"testing"

	"github.com/common-fate/granted-approvals/pkg/config"
	"github.com/joho/godotenv"
)

func TestFrontendDeployer(t *testing.T) {
	t.Skip("skipping test because its not ready to be integration tested in CI")
	// this is handy to test the file sync process however we will need to work out a good integration test strategy
	// Maybe split out the cloudfront invalidation then the file sync and aws exports can be tested against some testing s3 buckets
	ctx := context.Background()
	_ = godotenv.Load("../../.env")
	err := DeployProductionFrontend(ctx, config.FrontendDeployerConfig{
		CFReleasesBucket:                     "cf-example-customer-releases-ap-southest-2",
		CFReleasesFrontendAssetsObjectPrefix: "out",
		FrontendBucket:                       "test-bucket-josh2",
		APIURL:                               "https://test.execute-api.us-east-1.amazonaws.com/prod/",
		UserPoolDomain:                       "test.auth.us-east-1.amazoncognito.com",
		CognitoClientID:                      "2aqedb08vdqnabcdeo5u51udlvg",
		FrontendDomain:                       "aaaaaaaaaa.cloudfront.net",
		UserPoolID:                           "us-east-1_abcdefg",
		Region:                               "us-east-1",
		CloudFrontDistributionID:             "1234",
	})
	if err != nil {
		t.Fatal(err)
	}
}
