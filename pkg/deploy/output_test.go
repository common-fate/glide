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
	//@todo fix this test in CI
	t.Skip()
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
		CognitoClientID:                     "abcdefg",
		CloudFrontDomain:                    "abcdefg",
		FrontendDomainOutput:                "abcdefg",
		CloudFrontDistributionID:            "abcdefg",
		S3BucketName:                        "abcdefg",
		UserPoolID:                          "abcdefg",
		UserPoolDomain:                      "abcdefg",
		APIURL:                              "abcdefg",
		WebhookURL:                          "abcdefg",
		WebhookLogGroupName:                 "abcdefg",
		APILogGroupName:                     "abcdefg",
		IDPSyncLogGroupName:                 "abcdefg",
		EventBusLogGroupName:                "abcdefg",
		EventsHandlerConcurrentLogGroupName: "abcdefg",
		EventsHandlerSequentialLogGroupName: "abcdefg",
		GranterLogGroupName:                 "abcdefg",
		SlackNotifierLogGroupName:           "abcdefg",
		DynamoDBTable:                       "abcdefg",
		EventBusArn:                         "abcdefg",
		EventBusSource:                      "abcdefg",
		IDPSyncFunctionName:                 "abcdefg",
		SAMLIdentityProviderName:            "abcdefg",
		Region:                              "abcdefg",
		PaginationKMSKeyARN:                 "abcdefg",
		CacheSyncLogGroupName:               "abcdefg",
		RestAPIExecutionRoleARN:             "abcdefg",
		IDPSyncExecutionRoleARN:             "abcdefg",
		CacheSyncFunctionName:               "abcdefg",
		GovernanceURL:                       "abcdefg",
		CLIAppClientID:                      "abcdefg",
		HealthcheckFunctionName:             "abcdefg",
		HealthcheckLogGroupName:             "abcdefg",
		GranterV2StateMachineArn:            "abcdefg",
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

func TestOutput_Get(t *testing.T) {

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  Output
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			fields: Output{
				CognitoClientID: "test",
			},
			args: args{
				key: "CognitoClientID",
			},
			want: "test",
		},
		{
			name: "field not exist",
			args: args{
				key: "somethingelse",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Output{
				APILogGroupName:                     tt.fields.APILogGroupName,
				APIURL:                              tt.fields.APIURL,
				CacheSyncFunctionName:               tt.fields.CacheSyncFunctionName,
				CacheSyncLogGroupName:               tt.fields.CacheSyncLogGroupName,
				CLIAppClientID:                      tt.fields.CLIAppClientID,
				CloudFrontDistributionID:            tt.fields.CloudFrontDistributionID,
				CloudFrontDomain:                    tt.fields.CloudFrontDomain,
				CognitoClientID:                     tt.fields.CognitoClientID,
				DynamoDBTable:                       tt.fields.DynamoDBTable,
				EventBusArn:                         tt.fields.EventBusArn,
				EventBusLogGroupName:                tt.fields.EventBusLogGroupName,
				EventBusSource:                      tt.fields.EventBusSource,
				EventsHandlerConcurrentLogGroupName: tt.fields.EventsHandlerConcurrentLogGroupName,
				EventsHandlerSequentialLogGroupName: tt.fields.EventsHandlerSequentialLogGroupName,
				FrontendDomainOutput:                tt.fields.FrontendDomainOutput,
				GovernanceURL:                       tt.fields.GovernanceURL,
				GranterV2StateMachineArn:            tt.fields.GranterV2StateMachineArn,
				GranterLogGroupName:                 tt.fields.GranterLogGroupName,
				HealthcheckFunctionName:             tt.fields.HealthcheckFunctionName,
				HealthcheckLogGroupName:             tt.fields.HealthcheckLogGroupName,
				IDPSyncExecutionRoleARN:             tt.fields.IDPSyncExecutionRoleARN,
				IDPSyncFunctionName:                 tt.fields.IDPSyncFunctionName,
				IDPSyncLogGroupName:                 tt.fields.IDPSyncLogGroupName,
				PaginationKMSKeyARN:                 tt.fields.PaginationKMSKeyARN,
				Region:                              tt.fields.Region,
				RestAPIExecutionRoleARN:             tt.fields.RestAPIExecutionRoleARN,
				S3BucketName:                        tt.fields.S3BucketName,
				SAMLIdentityProviderName:            tt.fields.SAMLIdentityProviderName,
				SlackNotifierLogGroupName:           tt.fields.SlackNotifierLogGroupName,
				UserPoolDomain:                      tt.fields.UserPoolDomain,
				UserPoolID:                          tt.fields.UserPoolID,
				WebhookLogGroupName:                 tt.fields.WebhookLogGroupName,
				WebhookURL:                          tt.fields.WebhookURL,
			}
			got, err := o.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Output.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Output.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
