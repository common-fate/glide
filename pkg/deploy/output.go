package deploy

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go"
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/cfaws"
	"github.com/mitchellh/mapstructure"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

// Output is the output from deploying the CDK stack to AWS.
type Output struct {
	APILogGroupName                     string `json:"APILogGroupName"`
	APIURL                              string `json:"APIURL"`
	CacheSyncFunctionName               string `json:"CacheSyncFunctionName"`
	CacheSyncLogGroupName               string `json:"CacheSyncLogGroupName"`
	CLIAppClientID                      string `json:"CLIAppClientID"`
	CloudFrontDistributionID            string `json:"CloudFrontDistributionID"`
	CloudFrontDomain                    string `json:"CloudFrontDomain"`
	CognitoClientID                     string `json:"CognitoClientID"`
	DynamoDBTable                       string `json:"DynamoDBTable"`
	EventBusArn                         string `json:"EventBusArn"`
	EventBusLogGroupName                string `json:"EventBusLogGroupName"`
	EventBusSource                      string `json:"EventBusSource"`
	EventsHandlerConcurrentLogGroupName string `json:"EventsHandlerConcurrentLogGroupName"`
	EventsHandlerSequentialLogGroupName string `json:"EventsHandlerSequentialLogGroupName"`
	FrontendDomainOutput                string `json:"FrontendDomainOutput"`
	GovernanceURL                       string `json:"GovernanceURL"`
	GranterV2StateMachineArn            string `json:"GranterV2StateMachineArn"`
	GranterLogGroupName                 string `json:"GranterLogGroupName"`
	HealthcheckFunctionName             string `json:"HealthcheckFunctionName"`
	HealthcheckLogGroupName             string `json:"HealthcheckLogGroupName"`
	IDPSyncExecutionRoleARN             string `json:"IDPSyncExecutionRoleARN"`
	IDPSyncFunctionName                 string `json:"IDPSyncFunctionName"`
	IDPSyncLogGroupName                 string `json:"IDPSyncLogGroupName"`
	PaginationKMSKeyARN                 string `json:"PaginationKMSKeyARN"`
	Region                              string `json:"Region"`
	RestAPIExecutionRoleARN             string `json:"RestAPIExecutionRoleARN"`
	S3BucketName                        string `json:"S3BucketName"`
	SAMLIdentityProviderName            string `json:"SAMLIdentityProviderName"`
	SlackNotifierLogGroupName           string `json:"SlackNotifierLogGroupName"`
	UserPoolDomain                      string `json:"UserPoolDomain"`
	UserPoolID                          string `json:"UserPoolID"`
	WebhookLogGroupName                 string `json:"WebhookLogGroupName"`
	WebhookURL                          string `json:"WebhookURL"`
}

func (c Output) FrontendURL() string {
	if c.FrontendDomainOutput == "" {
		return "https://" + c.CloudFrontDomain
	}

	return "https://" + c.FrontendDomainOutput
}

func (c Output) PrintTable() {
	v := reflect.ValueOf(c)
	t := v.Type()

	data := [][]string{}
	for i := 0; i < v.NumField(); i++ {
		val := fmt.Sprintf("%v", v.Field(i).Interface())
		data = append(data, []string{t.Field(i).Name, val})
	}

	table := tablewriter.NewWriter(os.Stderr)
	table.SetHeader([]string{"Output Parameter", "Value"})

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.Bold},
	)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

// Keys returns the names of the output variables.
func (o Output) Keys() []string {
	v := reflect.ValueOf(o)
	t := v.Type()
	var keys []string

	for i := 0; i < v.NumField(); i++ {
		keys = append(keys, t.Field(i).Name)
	}
	return keys
}

// Get a value by it's key in the output struct
func (o Output) Get(key string) (string, error) {
	r := reflect.ValueOf(o)
	f := r.FieldByName(key)

	if !f.IsValid() {
		return "", fmt.Errorf("key %s not found in output", key)
	}

	return f.String(), nil
}

func (o Output) PrintSAMLTable() {
	data := [][]string{
		{"SAML SSO URL (ACS URL)", fmt.Sprintf("https://%s/saml2/idpresponse", o.UserPoolDomain)},
		{"Audience URI", fmt.Sprintf("urn:amazon:cognito:sp:%s", o.UserPoolID)},
		{"Dashboard URL", o.FrontendURL()},
	}

	table := tablewriter.NewWriter(os.Stderr)
	table.SetHeader([]string{"Parameter", "Value"})

	table.SetColumnColor(
		tablewriter.Colors{tablewriter.FgHiBlackColor},
		tablewriter.Colors{tablewriter.FgGreenColor, tablewriter.Bold},
	)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

// GetStackStatus indicates whether the Cloud Formation stack is online (via "CREATE_COMPLETE")
func (c *Config) GetStackStatus(ctx context.Context) (types.StackStatus, error) {
	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return "", err
	}
	client := cloudformation.NewFromConfig(cfg)
	res, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: &c.Deployment.StackName,
	})
	if err != nil {
		return "", err
	}
	if len(res.Stacks) != 1 {
		return "", fmt.Errorf("expected 1 stack but got %d", len(res.Stacks))
	}

	stack := res.Stacks[0]

	return stack.StackStatus, nil
}

// LoadOutput loads the outputs for the current deployment.
func (c *Config) LoadOutput(ctx context.Context) (Output, error) {
	if c.cachedOutput != nil {
		return *c.cachedOutput, nil
	}

	cfg, err := cfaws.ConfigFromContextOrDefault(ctx)
	if err != nil {
		return Output{}, err
	}
	client := cloudformation.NewFromConfig(cfg)
	res, err := client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
		StackName: &c.Deployment.StackName,
	})

	var ve *smithy.GenericAPIError
	if errors.As(err, &ve) && ve.Code == "ValidationError" {
		clio.Errorf("We couldn't find a CloudFormation stack '%s' in region '%s'.", c.Deployment.StackName, c.Deployment.Region)
		clio.Infof(`
To fix this, take one of the following actions:
  a) verify that your AWS credentials match the account you're trying to deploy to (%s). You can check this by calling 'aws sts get-caller-identity'.
  b) your stack may not have been deployed yet. Run 'gdeploy create' to deploy it using CloudFormation.
`, c.Deployment.Account)
		return Output{}, err
	}

	if err != nil {
		return Output{}, err
	}

	if len(res.Stacks) != 1 {
		return Output{}, fmt.Errorf("expected 1 stack but got %d", len(res.Stacks))
	}

	stack := res.Stacks[0]

	// decode the output variables into the Go struct.
	outputMap := make(map[string]string)
	for _, o := range stack.Outputs {
		outputMap[*o.OutputKey] = *o.OutputValue
	}

	var out Output
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{TagName: "json", Result: &out})
	if err != nil {
		return Output{}, errors.Wrap(err, "setting up decoder")
	}
	err = decoder.Decode(outputMap)
	if err != nil {
		return Output{}, errors.Wrap(err, "decoding CloudFormation outputs")
	}

	// set the cached output value in case this method is called again.
	c.cachedOutput = &out

	return out, nil
}
