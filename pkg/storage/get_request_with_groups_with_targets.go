package storage

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type GetRequestWithGroupsWithTargets struct {
	ID     string
	Result *access.RequestWithGroupsWithTargets
}

func (g *GetRequestWithGroupsWithTargets) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1 and begins_with(SK, :sk1)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.AccessRequest.PK1},
			":sk1": &types.AttributeValueMemberS{Value: keys.AccessRequest.SK1(g.ID)},
		},
	}

	return qi, nil
}

func (g *GetRequestWithGroupsWithTargets) UnmarshalQueryOutput(out *dynamodb.QueryOutput) error {
	if len(out.Items) == 0 {
		return ddb.ErrNoItems
	}
	err := attributevalue.UnmarshalMap(out.Items[0], &g.Result)
	if err != nil {
		return err
	}

	groups := make(map[string]access.GroupWithTargets)

	for _, v := range out.Items[1:] {
		// items will come out in order, groups first, then targets
		// The process here is to assert which type the item is, then unmarshal it to the correct type.
		// targets need to be assigned onto the correct group struct, so we use a map to track them
		if !strings.Contains((v["SK"].(*types.AttributeValueMemberS).Value), keys.AccessRequestGroupTargetKey) {
			var g access.Group
			err := attributevalue.UnmarshalMap(v, &g)
			if err != nil {
				return err
			}
			groups[g.ID] = access.GroupWithTargets{
				Group: g,
			}
		} else {
			var t access.GroupTarget
			err := attributevalue.UnmarshalMap(v, &t)
			if err != nil {
				return err
			}
			g := groups[t.GroupID]
			g.Targets = append(g.Targets, t)
			groups[t.GroupID] = g
		}
	}
	for _, grp := range groups {
		g.Result.Groups = append(g.Result.Groups, grp)
	}

	return nil
}
