package storage

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/storage/keys"
)

type RequestEndComparator string

const (
	LessThan         RequestEndComparator = "<"
	LessThanEqual    RequestEndComparator = "<="
	GreaterThan      RequestEndComparator = ">"
	GreaterThanEqual RequestEndComparator = ">="
)

// See the access.Request.DDBKeys for a comment explaining what the endtime represents for requests
type ListRequestsForUserAndRuleAndRequestend struct {
	UserID               string
	RuleID               string
	RequestEndComparator RequestEndComparator
	CompareTo            time.Time
	Result               []access.Request `ddb:"result"`
}

func (l *ListRequestsForUserAndRuleAndRequestend) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		IndexName:              &keys.IndexNames.GSI4,
		ScanIndexForward:       aws.Bool(false),
		KeyConditionExpression: aws.String(fmt.Sprintf("GSI4PK = :pk and GSI4SK %s :sk", l.RequestEndComparator)),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI4PK(l.UserID, l.RuleID)},
			":sk": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI4SK(l.CompareTo)},
		},
	}
	return &qi, nil
}
