package deploy

import (
	"encoding/json"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test ensures that the TS type for outputs in CDK matches our go type for outputs
// making a change to either will cause this test to fail.
// When you add a new output to the TS type, you will need to add it to the Go struct Output and update this test
func TestOutputStructMatchesTSType(t *testing.T) {
	cwd := "../../deploy/infra"
	// This command doesn't seem to work in github actions, the file ends up not existing.
	// Instead, I added a step in github actions which runs this same pnpm command which makes this test work correctly
	cmd := exec.Command("pnpm", "ts-node", "./test/stack-outputs.ts")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = cwd
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	testOutputs, err := os.ReadFile(path.Join(cwd, "testOutputs.json"))
	if err != nil {
		t.Fatal(err)
	}

	var rawObject map[string]string
	err = json.Unmarshal(testOutputs, &rawObject)
	if err != nil {
		t.Fatal(err)
	}
	output := Output{
		CognitoClientID:               "abcdefg",
		CloudFrontDomain:              "abcdefg",
		FrontendDomainOutput:          "abcdefg",
		CloudFrontDistributionID:      "abcdefg",
		S3BucketName:                  "abcdefg",
		UserPoolID:                    "abcdefg",
		UserPoolDomain:                "abcdefg",
		APIURL:                        "abcdefg",
		WebhookURL:                    "abcdefg",
		APILogGroupName:               "abcdefg",
		IDPSyncLogGroupName:           "abcdefg",
		AccessHandlerLogGroupName:     "abcdefg",
		EventBusLogGroupName:          "abcdefg",
		EventsHandlerLogGroupName:     "abcdefg",
		GranterLogGroupName:           "abcdefg",
		SlackNotifierLogGroupName:     "abcdefg",
		DynamoDBTable:                 "abcdefg",
		GranterStateMachineArn:        "abcdefg",
		EventBusArn:                   "abcdefg",
		EventBusSource:                "abcdefg",
		IdpSyncFunctionName:           "abcdefg",
		Region:                        "abcdefg",
		PaginationKMSKeyARN:           "abcdefg",
		AccessHandlerExecutionRoleARN: "abcdefg",
	}
	b, err := json.Marshal(output)
	if err != nil {
		t.Fatal(err)
	}
	var goStruct map[string]string
	err = json.Unmarshal(b, &goStruct)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, rawObject, goStruct)

	// Now do a sanity check that if something is different it does fail
	output.APILogGroupName = "wrong"
	b, err = json.Marshal(output)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(b, &goStruct)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, rawObject, goStruct)
}
