package storage

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

type ListRequestWithGroupsWithTargets struct {
	Result []access.RequestWithGroupsWithTargets
}

var _ ddb.QueryOutputUnmarshalerWithPagination = &ListRequestWithGroupsWithTargets{}

func (g *ListRequestWithGroupsWithTargets) BuildQuery() (*dynamodb.QueryInput, error) {
	qi := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk1"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk1": &types.AttributeValueMemberS{Value: keys.AccessRequest.PK1},
		},
	}

	return qi, nil
}

func (g *ListRequestWithGroupsWithTargets) UnmarshalQueryOutputWithPagination(out *dynamodb.QueryOutput) (map[string]types.AttributeValue, error) {
	if len(out.Items) == 0 {
		return nil, nil
	}
	var request access.RequestWithGroupsWithTargets
	groups := make(map[string]access.GroupWithTargets)
	var lastTargetForLastCompleteRequest *access.GroupTarget

	completeUnmarshallingRequest := func() (map[string]types.AttributeValue, error) {
		var lastTargetForCurrentRequest *access.GroupTarget
		var foundTargetCount int
		for _, grp := range groups {
			foundTargetCount += len(grp.Targets)
			if len(grp.Targets) > 0 {
				lastTargetForCurrentRequest = &grp.Targets[len(grp.Targets)-1]
			}
			request.Groups = append(request.Groups, grp)
		}

		if foundTargetCount != request.GroupTargetCount {
			// The full request must have been paginated, so instead of saving it, use the last target from the last full request as the pagination key.
			if lastTargetForLastCompleteRequest != nil {
				keys, err := lastTargetForLastCompleteRequest.DDBKeys()
				if err != nil {
					return nil, err
				}
				return map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: keys.PK},
					"SK": &types.AttributeValueMemberS{Value: keys.SK},
				}, nil
			}
			return nil, errors.New("failed to unmarshal requests, this could happen if the data for the request exceeds the 1mb limit for a ddb query")
		}
		lastTargetForLastCompleteRequest = lastTargetForCurrentRequest
		g.Result = append(g.Result, request)
		request = access.RequestWithGroupsWithTargets{}
		groups = make(map[string]access.GroupWithTargets)
		return nil, nil
	}
	for i, item := range out.Items {
		// items will come out in order, groups first, then targets
		// The process here is to assert which type the item is, then unmarshal it to the correct type.
		// targets need to be assigned onto the correct group struct, so we use a map to track them
		if !strings.Contains((item["SK"].(*types.AttributeValueMemberS).Value), keys.AccessRequestGroupKey) {
			// we have found the start of a new request, so save the previous completely unmarshalled request to the output and reset the request type
			if i > 0 {
				o, err := completeUnmarshallingRequest()
				if err != nil {
					return o, err
				}
			}
			// it is a request
			err := attributevalue.UnmarshalMap(item, &request)
			if err != nil {
				return nil, err
			}
		} else if !strings.Contains((item["SK"].(*types.AttributeValueMemberS).Value), keys.AccessRequestGroupTargetKey) {
			// it is a group
			var g access.Group
			err := attributevalue.UnmarshalMap(item, &g)
			if err != nil {
				return nil, err
			}
			groups[g.ID] = access.GroupWithTargets{
				Group: g,
			}
		} else {
			// it is a target
			var t access.GroupTarget
			err := attributevalue.UnmarshalMap(item, &t)
			if err != nil {
				return nil, err
			}
			g := groups[t.GroupID]
			g.Targets = append(g.Targets, t)
			groups[t.GroupID] = g
		}

	}
	return completeUnmarshallingRequest()
}
