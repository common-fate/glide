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
	GroupTargets(ctx context.Context, targets []types.Target) (requests.Requestv2, error)
}

func (s *Service) GroupTargets(ctx context.Context, targets []types.Target) (*requests.Requestv2, error) {

	u := auth.UserFromContext(ctx)

	preflight := requests.Requestv2{
		ID: types.NewRequestID(),
		// Groups:      map[string]requests.AccessGroup{},
		RequestedBy: *u,
	}

	items := []ddb.Keyer{}

	//Go through each1 target in the request and group them up based on access rule

	//eg. a preflight request could have targets from multiple fields in the same access rule
	//as well as targets from a different access rule. Eg. Aws sso and OKTA groups
	//goal here is to have a list of these groups which can be saved to the database and be easily read back
	//to be processed into grants on submission
	for _, target := range targets {

		//Grouping up targets based on which access rule they are apart of
		_, ok := preflight.Groups[target.AccessRule]
		if !ok {
			newTarget := map[string]string{}

			for key, val := range target.With.AdditionalProperties {
				newTarget[key] = val
			}

			//lookup access rule

			ac := storage.GetAccessRuleCurrent{ID: target.AccessRule}

			_, err := s.DB.Query(ctx, &ac)
			if err != nil {
				return nil, err
			}

			now := time.Now()
			preflight.Groups[target.AccessRule] = requests.AccessGroup{
				AccessRule:      *ac.Result,
				Reason:          target.Reason,
				TimeConstraints: requests.Timing{Duration: time.Duration(target.TimeConstraints.MaxDurationSeconds), StartTime: &now},
				With:            []map[string]string{newTarget},
				Request:         preflight.ID,
				UpdatedAt:       now,
				Status:          requests.PENDING,
			}

		} else {

			newTarget := map[string]string{}

			for key, val := range target.With.AdditionalProperties {
				newTarget[key] = val
			}

			if thisRequest, ok := preflight.Groups[target.AccessRule]; ok {
				thisRequest.With = append(thisRequest.With, newTarget)
				preflight.Groups[target.AccessRule] = thisRequest
			}

		}

	}

	//save the group items and the preflight request
	for _, group := range preflight.Groups {
		items = append(items, &group)
	}

	items = append(items, &preflight)

	err := s.DB.PutBatch(ctx, items...)
	if err != nil {
		return nil, err
	}
	//create a preflight object in the db
	return &preflight, nil
}
