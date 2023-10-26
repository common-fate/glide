package autoapproval

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
)

type Status string

const (
	AUTO_APPROVED     Status = "AUTO_APPROVED"
	REQUIRES_APPROVAL Status = "REQUIRES_APPROVAL"
)

type ResponseBody struct {
	Decision      Status `json:"decision"`
	Justification string `json:"justification,omitempty"`
}

type RequestBody struct {
	User identity.User   `json:"user"`
	Rule rule.AccessRule `json:"rule"`
}

type Service struct {
}

func (s Service) Autoapprove(user identity.User, rule rule.AccessRule, lambdaArn string) (bool, error) {
	sess, err := session.NewSession()
	//	&aws.Config{
	//	Region: aws.String("us-west-2"), // Replace with your region
	//}

	if err != nil {
		return false, err
	}

	req := RequestBody{User: user, Rule: rule}

	payload, err := json.Marshal(req)

	if err != nil {
		return false, err
	}

	params := &lambda.InvokeInput{
		FunctionName:   aws.String(lambdaArn), // Lambda arn
		Payload:        payload,
		InvocationType: aws.String("RequestResponse"), // Get synchronous output
	}

	svc := lambda.New(sess)

	resp, err := svc.Invoke(params)

	if err != nil {
		fmt.Println("Error happened when Invoke lambda ", err)
	}

	var output ResponseBody
	err = json.Unmarshal(resp.Payload, &output)

	if err != nil {
		return false, err
	}

	if resp.FunctionError != nil {
		return false, errors.New("Error happened when calling lambda: " + *resp.FunctionError)
	}

	return output.Decision == AUTO_APPROVED, nil
}
