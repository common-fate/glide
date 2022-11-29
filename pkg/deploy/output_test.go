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
		WebhookLogGroupName:           "abcdefg",
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
		SAMLIdentityProviderName:      "abcdefg",
		Region:                        "abcdefg",
		PaginationKMSKeyARN:           "abcdefg",
		AccessHandlerExecutionRoleARN: "abcdefg",
		CacheSyncLogGroupName:         "abcdefg",
		RestAPIExecutionRoleARN:       "abcdefg",
		IDPSyncExecutionRoleARN:       "abcdefg",
		CacheSyncFunctionName:         "abcdefg",
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
	type fields struct {
		CognitoClientID               string
		CloudFrontDomain              string
		FrontendDomainOutput          string
		CloudFrontDistributionID      string
		S3BucketName                  string
		UserPoolID                    string
		UserPoolDomain                string
		APIURL                        string
		WebhookURL                    string
		APILogGroupName               string
		IDPSyncLogGroupName           string
		AccessHandlerLogGroupName     string
		EventBusLogGroupName          string
		EventsHandlerLogGroupName     string
		GranterLogGroupName           string
		SlackNotifierLogGroupName     string
		DynamoDBTable                 string
		GranterStateMachineArn        string
		EventBusArn                   string
		EventBusSource                string
		IdpSyncFunctionName           string
		SAMLIdentityProviderName      string
		Region                        string
		PaginationKMSKeyARN           string
		AccessHandlerExecutionRoleARN string
		CacheSyncFunctionName         string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
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
				CognitoClientID:               tt.fields.CognitoClientID,
				CloudFrontDomain:              tt.fields.CloudFrontDomain,
				FrontendDomainOutput:          tt.fields.FrontendDomainOutput,
				CloudFrontDistributionID:      tt.fields.CloudFrontDistributionID,
				S3BucketName:                  tt.fields.S3BucketName,
				UserPoolID:                    tt.fields.UserPoolID,
				UserPoolDomain:                tt.fields.UserPoolDomain,
				APIURL:                        tt.fields.APIURL,
				WebhookURL:                    tt.fields.WebhookURL,
				APILogGroupName:               tt.fields.APILogGroupName,
				IDPSyncLogGroupName:           tt.fields.IDPSyncLogGroupName,
				AccessHandlerLogGroupName:     tt.fields.AccessHandlerLogGroupName,
				EventBusLogGroupName:          tt.fields.EventBusLogGroupName,
				EventsHandlerLogGroupName:     tt.fields.EventsHandlerLogGroupName,
				GranterLogGroupName:           tt.fields.GranterLogGroupName,
				SlackNotifierLogGroupName:     tt.fields.SlackNotifierLogGroupName,
				DynamoDBTable:                 tt.fields.DynamoDBTable,
				GranterStateMachineArn:        tt.fields.GranterStateMachineArn,
				EventBusArn:                   tt.fields.EventBusArn,
				EventBusSource:                tt.fields.EventBusSource,
				IdpSyncFunctionName:           tt.fields.IdpSyncFunctionName,
				SAMLIdentityProviderName:      tt.fields.SAMLIdentityProviderName,
				Region:                        tt.fields.Region,
				PaginationKMSKeyARN:           tt.fields.PaginationKMSKeyARN,
				AccessHandlerExecutionRoleARN: tt.fields.AccessHandlerExecutionRoleARN,
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
