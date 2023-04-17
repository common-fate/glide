package preflightsvc

import (
	"context"
	"time"

	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/requests"
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

		//Grouping up targets based on which access rule they are apart of
		_, ok := accessGroups[target.AccessRule]
		if !ok {
			newTarget := requests.Target{}

			for key, val := range target.With.AdditionalProperties {
				newTarget.Fields[key] = requests.Field{Value: requests.FieldValue{Type: "", Value: val}}
			}

			//lookup access rule
			ac := storage.GetAccessRuleCurrent{ID: target.AccessRule}
			_, err := s.DB.Query(ctx, &ac)
			if err != nil {
				return nil, err
			}

			ag := requests.AccessGroup{
				ID:              types.NewAccessGroupID(),
				AccessRule:      *ac.Result,
				TimeConstraints: requests.Timing{Duration: time.Duration(target.TimeConstraints.MaxDurationSeconds), StartTime: &now},
				Request:         request.ID,
				UpdatedAt:       now,
				Status:          requests.PENDING,
			}
			accessGroups[target.AccessRule] = ag

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

			for key, val := range target.With.AdditionalProperties {
				newTarget.Fields[key] = requests.Field{Value: requests.FieldValue{Type: "", Value: val}}
			}

			ag := accessGroups[target.AccessRule]

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
