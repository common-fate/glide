package deploy

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCDKOutput = Output{
	APIURL:               "https://test.execute-api.us-east-1.amazonaws.com/prod/",
	UserPoolDomain:       "test.auth.us-east-1.amazoncognito.com",
	CognitoClientID:      "2aqedb08vdqnabcdeo5u51udlvg",
	CloudFrontDomain:     "aaaaaaaaaa.cloudfront.net",
	FrontendDomainOutput: "example.granted.run",
	UserPoolID:           "us-east-1_abcdefg",
	S3BucketName:         "common-fate-test-us-east-1-12345567890",
	DynamoDBTable:        "common-fate-test",
	Region:               "us-east-1",
	CLIAppClientID:       "cli-client-app-id",
}

func TestRenderLocalFrontend(t *testing.T) {
	want, err := os.ReadFile("testdata/aws-exports.js.snapshot")
	if err != nil {
		t.Fatal(err)
	}
	got, err := RenderLocalFrontendConfig(testCDKOutput.ToRenderFrontendConfig())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(want), got)
}

func TestRenderProductionFrontend(t *testing.T) {
	want, err := os.ReadFile("testdata/aws-exports.json.snapshot")
	if err != nil {
		t.Fatal(err)
	}
	got, err := RenderProductionFrontendConfig(testCDKOutput.ToRenderFrontendConfig())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(want), got)
}
