package preflightsvc

import (
	"context"

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
		ID:     types.NewRequestID(),
		Groups: map[string]requests.AccessGroup{},
		User:   *u,
	}

	//Go through each target in the request and group them up based on access rule

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

			preflight.Groups[target.AccessRule] = requests.AccessGroup{
				AccessRule:      *ac.Result,
				Reason:          target.Reason,
				TimeConstraints: target.TimeConstraints,
				With:            []map[string]string{newTarget},
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
	//validate current user has access to access rules

	//group requests based on duration and purpose

	//create a preflight object in the db
	return &preflight, nil
}
