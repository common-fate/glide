package preflightsvc

import (
	"context"
	"errors"
	"time"

	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/requests"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Service struct {
	DB ddb.Storage
}

type PreflightService interface {
	GroupTargets(ctx context.Context, preflightRequest types.CreatePreflightRequest) (*requests.Preflight, error)
}

// Takes in a list of targets and groups them by access rule
// Creates an access group per access rule
// Creates a grant for each target and links back to the access group
// creates and returns a preflight object containing all the data of the request
// Creates all necissary componants which makes up a grant. Request, Access Groups and Grants. Returns a Preflight object
func (s *Service) GroupTargets(ctx context.Context, preflightRequest types.CreatePreflightRequest) (*requests.Preflight, error) {
	now := time.Now()
	preflight := requests.Preflight{}
	u := auth.UserFromContext(ctx)

	request := requests.Requestv2{
		ID:          types.NewRequestID(),
		RequestedBy: *u,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	items := []ddb.Keyer{}

	//use a map to keep track on what acecss groups are made, and not make any duplicates
	accessGroups := map[string]requests.AccessGroup{}

	//organise each grant from the list of targets and create access group / grant for each entitlement
	for _, target := range preflightRequest.Targets {

		//todo: validate that resource is apart of access rule with a look up.
		//lookup each target

		targetItem := storage.GetTarget{ID: target.Id}
		_, err := s.DB.Query(ctx, &targetItem)
		if err != nil {
			return nil, err
		}

		//returns the best suited access rule for the current user
		accessRule, err := s.getAccessRuleForTarget(ctx, targetItem.Result.AccessRules)
		if err != nil {
			return nil, err
		}
		//Grouping up targets based on which access rule they are apart of
		_, ok := accessGroups[accessRule.ID]
		if !ok {
			newTarget := requests.Target{}

			for _, field := range target.Fields {
				newTarget.Fields[field.Id] = requests.Field{Value: requests.FieldValue{Type: "", Value: field.Value}}
			}

			ag := requests.AccessGroup{
				ID:              types.NewAccessGroupID(),
				AccessRule:      *accessRule,
				TimeConstraints: requests.Timing{Duration: time.Duration(accessRule.TimeConstraints.MaxDurationSeconds), StartTime: &now},
				Request:         request.ID,
				UpdatedAt:       now,
				Status:          requests.PENDING,
			}
			accessGroups[accessRule.ID] = ag

			newGrant := requests.Grantv2{
				ID:          types.NewGrantID(),
				AccessGroup: ag.ID,
				Target:      newTarget,
				Subject:     u.Email,
				Status:      types.GrantStatus(requests.PENDING),
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			//add all grants onto the preflight
			preflight.Grants = append(preflight.Grants, newGrant)

			items = append(items, &newGrant)

		} else {

			newTarget := requests.Target{}

			for _, field := range target.Fields {
				newTarget.Fields[field.Id] = requests.Field{Value: requests.FieldValue{Type: "", Value: field.Value}}
			}

			ag := accessGroups[accessRule.ID]

			newGrant := requests.Grantv2{
				ID:          types.NewGrantID(),
				AccessGroup: ag.ID,
				Target:      newTarget,
				Subject:     u.Email,
				Status:      types.GrantStatus(requests.PENDING),
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			items = append(items, &newGrant)

			//add all grants onto the preflight
			preflight.Grants = append(preflight.Grants, newGrant)

		}

	}

	//save the final group items without any dupes and the preflight request
	for _, group := range accessGroups {
		items = append(items, &group)

		//add all groups onto the preflight
		preflight.AccessGroups = append(preflight.AccessGroups, group)
	}

	//add the request item to be saved
	items = append(items, &request)

	err := s.DB.PutBatch(ctx, items...)
	if err != nil {
		return nil, err
	}
	//create a preflight object in the db
	return &preflight, nil
}

func (s *Service) getAccessRuleForTarget(ctx context.Context, accessRules []string) (*rule.AccessRule, error) {
	if len(accessRules) <= 0 {
		return nil, errors.New("no access groups found")
	}

	var currentRule rule.AccessRule

	for _, rule := range accessRules {
		ar := storage.GetAccessRule{ID: rule}
		_, err := s.DB.Query(ctx, &ar)
		if err != nil {
			return nil, err
		}

		//from a list of many rules we should return the rule with the lowest barrier for entry for the user...
		//pick the one with no approval needed
		//with the longest duration

		//if first iteration just make the current rule = first
		//todo: make test cases for these
		if currentRule.ID == "" {
			currentRule = *ar.Result
			continue
		}

		//if new rule doesnt require approval, override it
		if currentRule.Approval.IsRequired() && !ar.Result.Approval.IsRequired() {
			currentRule = *ar.Result
		}

		//if both rules dont require access, but new rule has longer duration. Override it
		if !currentRule.Approval.IsRequired() && !ar.Result.Approval.IsRequired() && ar.Result.TimeConstraints.MaxDurationSeconds > currentRule.TimeConstraints.MaxDurationSeconds {
			currentRule = *ar.Result
		}

		//if both rules require approval, but new rule has longer duration. Override it.
		if currentRule.Approval.IsRequired() && ar.Result.Approval.IsRequired() && ar.Result.TimeConstraints.MaxDurationSeconds > currentRule.TimeConstraints.MaxDurationSeconds {
			currentRule = *ar.Result

		}

	}

	return nil, nil

}
