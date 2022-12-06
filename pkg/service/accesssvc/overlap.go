package accesssvc

import (
	"context"
	"fmt"
	"reflect"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
)

type requestAndRule struct {
	request access.Request
	rule    rule.AccessRule
}

func overlapsExistingGrantCheck(req access.Request, upcomingRequests []access.Request, currentRequestRule rule.AccessRule, allRules []rule.AccessRule, clock clock.Clock) (bool, error) {
	start, end := req.GetInterval(access.WithNow(clock.Now()))
	var upcomingRequestAndRules []requestAndRule

	ruleMap := make(map[string]rule.AccessRule)

	for _, accessRule := range allRules {
		ruleMap[accessRule.ID] = accessRule
	}

	//make a map of requests mapped to their relative access rules
	for _, upcomingRequest := range upcomingRequests {

		if accessRule, ok := ruleMap[upcomingRequest.Rule]; ok {
			upcomingRequestAndRules = append(upcomingRequestAndRules, requestAndRule{request: upcomingRequest, rule: accessRule})

		} else {
			return false, fmt.Errorf("request contains access rule that does not exist")
		}
	}

	currentRequestArguments := make(map[string]string)
	for k, v := range currentRequestRule.Target.With {
		currentRequestArguments[k] = v
	}
	for k, v := range req.SelectedWith {
		currentRequestArguments[k] = v.Value
	}

	for _, r := range upcomingRequestAndRules {

		//check provider is the same
		if r.rule.Target.ProviderID != currentRequestRule.Target.ProviderID {
			continue
		}

		upcomingStart, upcomingEnd := r.request.GetInterval(access.WithNow(clock.Now()))
		if (start.Before(upcomingEnd) || start.Equal(upcomingEnd)) && (end.After(upcomingStart) || end.Equal(upcomingStart)) {

			//check the arguments overlap
			upcomingRequestArguments := make(map[string]string)
			for k, v := range r.rule.Target.With {
				upcomingRequestArguments[k] = v
			}
			for k, v := range r.request.SelectedWith {
				upcomingRequestArguments[k] = v.Value
			}
			//check if the grant is actually active
			if r.request.Grant != nil {
				if r.request.Grant.Status == "ACTIVE" || r.request.Grant.Status == "PENDING" {
					if reflect.DeepEqual(currentRequestArguments, upcomingRequestArguments) {
						return true, nil
					}
				}

			}

		}

	}
	return false, nil
}

func (s *Service) overlapsExistingGrant(ctx context.Context, req access.Request) (bool, error) {
	start, _ := req.GetInterval(access.WithNow(s.Clock.Now()))

	rq := storage.ListRequestsForUserAndRequestend{
		UserID:               req.RequestedBy,
		RequestEndComparator: storage.GreaterThanEqual,
		CompareTo:            start,
	}
	_, err := s.DB.Query(ctx, &rq)
	if err != nil && err != ddb.ErrNoItems {
		return false, err
	}
	upcomingRequests := rq.Result
	if len(upcomingRequests) == 0 {
		return false, nil
	}

	ruleq := storage.GetAccessRuleCurrent{ID: req.Rule}
	_, err = s.DB.Query(ctx, &ruleq)
	if err != nil {
		return false, err
	}

	allRules := storage.ListCurrentAccessRules{}
	_, err = s.DB.Query(ctx, &allRules)
	if err != nil {
		return false, err
	}

	isOverlapping, err := overlapsExistingGrantCheck(req, upcomingRequests, *ruleq.Result, allRules.Result, s.Clock)
	if err != nil {
		return false, err
	}
	return isOverlapping, nil
}
