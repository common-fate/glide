package storage

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/common-fate/ddb"
	"github.com/joho/godotenv"
)

type testClient struct {
	db *ddb.Client
}

// newTestingStorage creates a new testing storage.
// It skips the tests if the TESTING_DYNAMODB_TABLE env var isn't set.
func newTestingStorage(t *testing.T) *testClient {
	ctx := context.Background()
	_ = godotenv.Load("../../.env")
	table := os.Getenv("TESTING_DYNAMODB_TABLE")
	if table == "" {
		t.Skip("TESTING_DYNAMODB_TABLE is not set")
	}
	s, err := ddb.New(ctx, table)
	if err != nil {
		t.Fatal(err)
	}
	return &testClient{
		db: s,
	}
}

type item struct {
	PK string
	SK string
}

func (i *item) DDBKeys() (ddb.Keys, error) {
	return ddb.Keys{
		PK: i.PK,
		SK: i.SK,
	}, nil
}

// type listAll struct {
// 	PK     string
// 	Result []item `ddb:"result"`
// }

// func (l *listAll) BuildQuery() (*dynamodb.QueryInput, error) {
// 	qi := dynamodb.QueryInput{
// 		KeyConditionExpression: aws.String("PK = :pk"),
// 		ExpressionAttributeValues: map[string]types.AttributeValue{
// 			":pk": &types.AttributeValueMemberS{Value: l.PK},
// 		},
// 	}
// 	return &qi, nil
// }

// // deletePartition deletes all items for a given partition key
// func (t *testClient) deletePartition(pk string) error {
// 	q := listAll{
// 		PK: pk,
// 	}

// 	ctx := context.Background()
// 	err := t.db.All(ctx, &q)
// 	if err != nil {
// 		return err
// 	}
// 	var items []ddb.Keyer
// 	for i := range q.Result {
// 		items = append(items, &q.Result[i])
// 	}

// 	return t.db.DeleteBatch(ctx, items...)
// }

// deleteAll scans dynamo then deletes all items found
func (t *testClient) deleteAll() error {
	ctx := context.Background()
	out, err := t.db.Client().Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(t.db.Table()),
	})
	if err != nil {
		return err
	}
	var items []ddb.Keyer
	for _, v := range out.Items {
		var i item
		err := attributevalue.UnmarshalMap(v, &i)
		if err != nil {
			return err
		}
		items = append(items, &i)
	}
	return t.db.DeleteBatch(ctx, items...)
}
