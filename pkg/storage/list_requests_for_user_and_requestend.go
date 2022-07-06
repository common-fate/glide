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

// See the access.Request.DDBKeys for a comment explaining what the endtime represents for requests
type ListRequestsForUserAndRequestend struct {
	UserID               string
	RequestEndComparator RequestEndComparator
	CompareTo            time.Time
	Result               []access.Request `ddb:"result"`
}

func (l *ListRequestsForUserAndRequestend) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := dynamodb.QueryInput{
		IndexName:              &keys.IndexNames.GSI3,
		ScanIndexForward:       aws.Bool(false),
		KeyConditionExpression: aws.String(fmt.Sprintf("GSI3PK = :pk and GSI3SK %s :sk", l.RequestEndComparator)),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI3PK(l.UserID)},
			":sk": &types.AttributeValueMemberS{Value: keys.AccessRequest.GSI3SK(l.CompareTo)},
		},
	}
	return &qi, nil
}
