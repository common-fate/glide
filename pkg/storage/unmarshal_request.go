package storage

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/storage/keys"
	"github.com/common-fate/ddb"
)

func UnmarshalRequestGroup(items []map[string]types.AttributeValue) (*access.GroupWithTargets, error) {
	if len(items) == 0 {
		return nil, ddb.ErrNoItems
	}
	var result access.GroupWithTargets
	err := attributevalue.UnmarshalMap(items[0], &result.Group)
	if err != nil {
		return nil, err
	}

	for _, v := range items[1:] {
		var t access.GroupTarget
		err := attributevalue.UnmarshalMap(v, &t)
		if err != nil {
			return nil, err
		}

		result.Targets = append(result.Targets, t)
	}

	return &result, nil
}

func UnmarshalRequest(items []map[string]types.AttributeValue) (*access.RequestWithGroupsWithTargets, error) {
	if len(items) == 0 {
		return nil, ddb.ErrNoItems
	}
	var result access.RequestWithGroupsWithTargets
	err := attributevalue.UnmarshalMap(items[0], &result.Request)
	if err != nil {
		return nil, err
	}

	groups := make(map[string]access.GroupWithTargets)

	for _, v := range items[1:] {
		// items will come out in order, groups first, then targets
		// The process here is to assert which type the item is, then unmarshal it to the correct type.
		// targets need to be assigned onto the correct group struct, so we use a map to track them
		sk := v["SK"]
		skval, ok := sk.(*types.AttributeValueMemberS)
		_ = ok
		if !strings.Contains(skval.Value, keys.AccessRequestGroupTargetKey) {
			var g access.Group
			err := attributevalue.UnmarshalMap(v, &g)
			if err != nil {
				return nil, err
			}
			g.RequestStatus = result.Request.RequestStatus
			groups[g.ID] = access.GroupWithTargets{
				Group: g,
			}
		} else {
			var t access.GroupTarget
			err := attributevalue.UnmarshalMap(v, &t)
			if err != nil {
				return nil, err
			}
			g := groups[t.GroupID]
			g.Targets = append(g.Targets, t)

			groups[t.GroupID] = g
		}
	}
	for _, grp := range groups {
		result.Groups = append(result.Groups, grp)
	}

	return &result, nil
}

func UnmarshalRequests(items []map[string]types.AttributeValue) ([]access.RequestWithGroupsWithTargets, map[string]types.AttributeValue, error) {
	if len(items) == 0 {
		return nil, nil, nil
	}

	var result []access.RequestWithGroupsWithTargets
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

		if foundTargetCount != request.Request.GroupTargetCount {
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
		result = append(result, request)
		request = access.RequestWithGroupsWithTargets{}
		groups = make(map[string]access.GroupWithTargets)
		return nil, nil
	}
	for i, item := range items {
		// items will come out in order, groups first, then targets
		// The process here is to assert which type the item is, then unmarshal it to the correct type.
		// targets need to be assigned onto the correct group struct, so we use a map to track them
		if !strings.Contains((item["SK"].(*types.AttributeValueMemberS).Value), keys.AccessRequestGroupKey) {
			// we have found the start of a new request, so save the previous completely unmarshalled request to the output and reset the request type
			if i > 0 {
				o, err := completeUnmarshallingRequest()
				if err != nil {
					return nil, o, err
				}
			}
			// it is a request
			err := attributevalue.UnmarshalMap(item, &request.Request)
			if err != nil {
				return nil, nil, err
			}
		} else if !strings.Contains((item["SK"].(*types.AttributeValueMemberS).Value), keys.AccessRequestGroupTargetKey) {
			// it is a group
			var g access.Group
			err := attributevalue.UnmarshalMap(item, &g)
			if err != nil {
				return nil, nil, err
			}
			groups[g.ID] = access.GroupWithTargets{
				Group: g,
			}
		} else {
			// it is a target
			var t access.GroupTarget
			err := attributevalue.UnmarshalMap(item, &t)
			if err != nil {
				return nil, nil, err
			}
			g := groups[t.GroupID]
			g.Targets = append(g.Targets, t)
			groups[t.GroupID] = g
		}

	}
	pagination, err := completeUnmarshallingRequest()
	if err != nil {
		return nil, nil, err
	}
	return result, pagination, nil
}

// UnmarshalRequestsBottomToTop is used to unmarshal requests starting at the targets and ending on the request
// useful so that you can read items on of dynamo in either scan forward or scan reversed
func UnmarshalRequestsBottomToTop(items []map[string]types.AttributeValue) ([]access.RequestWithGroupsWithTargets, map[string]types.AttributeValue, error) {
	if len(items) == 0 {
		return nil, nil, nil
	}

	var result []access.RequestWithGroupsWithTargets
	groups := make(map[string]access.GroupWithTargets)

	for _, item := range items {
		if strings.Contains((item["SK"].(*types.AttributeValueMemberS).Value), keys.AccessRequestGroupTargetKey) {
			// it is a target
			var t access.GroupTarget
			err := attributevalue.UnmarshalMap(item, &t)
			if err != nil {
				return nil, nil, err
			}
			g := groups[t.GroupID]
			g.Targets = append(g.Targets, t)
			groups[t.GroupID] = g
		} else if strings.Contains((item["SK"].(*types.AttributeValueMemberS).Value), keys.AccessRequestGroupKey) {
			// it is a group
			var g access.Group
			err := attributevalue.UnmarshalMap(item, &g)
			if err != nil {
				return nil, nil, err
			}
			group := groups[g.ID]
			group.Group = g
			groups[g.ID] = group

		} else {
			// it is a request
			var r access.RequestWithGroupsWithTargets
			err := attributevalue.UnmarshalMap(item, &r.Request)
			if err != nil {
				return nil, nil, err
			}
			for _, v := range groups {
				r.Groups = append(r.Groups, v)
			}
			result = append(result, r)
			groups = make(map[string]access.GroupWithTargets)
		}
	}

	if len(groups) != 0 {
		if len(result) != 0 {
			keys, err := result[len(result)-1].Request.DDBKeys()
			if err != nil {
				return nil, nil, err
			}
			return result, map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: keys.PK},
				"SK": &types.AttributeValueMemberS{Value: keys.SK},
			}, nil
		} else {
			return nil, nil, errors.New("failed to unmarshal requests, this could happen if the data for the request exceeds the 1mb limit for a ddb query")
		}
	}

	return result, nil, nil
}
