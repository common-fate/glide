package preflightsvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/requestsv2.go"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Service struct {
	DB ddb.Storage
}

type PreflightService interface {
	GroupTargets(ctx context.Context, targets []types.Target) (requestsv2.Preflight, error)
}

func (s *Service) GroupTargets(ctx context.Context, targets []types.Target) (requestsv2.Preflight, error) {

	u := auth.UserFromContext(ctx)

	preflight := requestsv2.Preflight{
		ID:       types.NewPreflightID(),
		Requests: map[string]requestsv2.PreflightRequest{},
		User:     u.ID,
	}

	for _, target := range targets {

		//Look up where
		//does this access rule exist in the preflight request map?

		_, ok := preflight.Requests[target.AccessRule]

		//if exists add to the array of targets

		//if not exists create the entry in the map and add target

		if !ok {
			newTarget := map[string]string{}

			for key, val := range target.With.AdditionalProperties {
				newTarget[key] = val
			}

			preflight.Requests[target.AccessRule] = requestsv2.PreflightRequest{
				AccessRule:      target.AccessRule,
				Reason:          target.Reason,
				TimeConstraints: target.TimeConstraints,
				With:            []map[string]string{newTarget},
			}

		} else {

			newTarget := map[string]string{}

			for key, val := range target.With.AdditionalProperties {
				newTarget[key] = val
			}

			if thisRequest, ok := preflight.Requests[target.AccessRule]; ok {
				thisRequest.With = append(thisRequest.With, newTarget)
				preflight.Requests[target.AccessRule] = thisRequest
			}

		}

	}
	//validate current user has access to access rules

	//group requests based on duration and purpose

	//create a preflight object in the db
	return preflight, nil
}
