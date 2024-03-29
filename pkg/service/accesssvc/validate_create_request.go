package accesssvc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

type validateCreateRequestsResponse struct {
	argumentCombinations types.RequestArgumentCombinations
	rule                 rule.AccessRule
	requestArguments     map[string]types.RequestArgument
}

// validateCreateRequests returns APIO errors for bad request errors relating to the whole request
func (s *Service) validateCreateRequests(ctx context.Context, in CreateRequestsOpts) (*validateCreateRequestsResponse, error) {
	// If a request was not submitted with any arguments, then combinations to create will return a single empty map
	// this will be validated against the rule later to ensure that it is expected
	combinationsToCreate, err := in.argumentCombinations()
	if err != nil {
		return nil, err
	}
	q := storage.GetAccessRuleCurrent{ID: in.Create.AccessRuleId}
	_, err = s.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		return nil, apio.NewRequestError(ErrRuleNotFound, http.StatusBadRequest)
	}
	if err != nil {
		// we don't know how to handle the error from the rule getter, so just return it to the caller.
		return nil, err
	}
	rule := *q.Result

	err = groupMatches(rule.Groups, in.User.Groups)
	if err != nil {
		return nil, err
	}

	requestArguments, err := s.Rules.RequestArguments(ctx, rule.Target)
	if err != nil {
		return nil, err
	}
	// validate all the requests for basic errors before attempting to create grants or object in the DB
	for _, combinationToCreate := range combinationsToCreate {
		err = validateCreateRequest(ctx, CreateRequest{
			AccessRuleId: in.Create.AccessRuleId,
			Reason:       in.Create.Reason,
			Timing:       in.Create.Timing,
			With:         combinationToCreate,
		}, rule, requestArguments)
		if err != nil {

			return nil, err
		}
	}
	return &validateCreateRequestsResponse{
		argumentCombinations: combinationsToCreate,
		rule:                 rule,
		requestArguments:     requestArguments,
	}, nil
}

func (cro CreateRequestsOpts) argumentCombinations() (types.RequestArgumentCombinations, error) {
	// create the request. The RequestCreator handles the validation
	// and saving the request to the database.
	var combinationsToCreate types.RequestArgumentCombinations

	// A request which contains no sub requests or a requst which contains one empty sub request is treated as having no arguments.
	// return a single empty combination.
	// this is to support requests for a rule where no arguments are selectable by the user
	if cro.Create.With == nil {
		combinationsToCreate = append(combinationsToCreate, make(map[string]string))
		return combinationsToCreate, nil
	} else {
		arr := *cro.Create.With
		if len(arr) == 0 || (len(arr) == 1 && len(arr[0].AdditionalProperties) == 0) {
			combinationsToCreate = append(combinationsToCreate, make(map[string]string))
			return combinationsToCreate, nil
		}
	}

	for _, v := range *cro.Create.With {
		// a request which contains subrequests with no arguments is invalid
		if len(v.AdditionalProperties) == 0 {
			return nil, apio.NewRequestError(errors.New("request contains subrequest with no arguments"), http.StatusBadRequest)
		}
		combinations, err := v.ArgumentCombinations()
		if errors.As(err, &types.ArgumentHasNoValuesError{}) {
			return nil, apio.NewRequestError(err, http.StatusBadRequest)
		}
		if err != nil {
			return nil, err
		}
		combinationsToCreate = append(combinationsToCreate, combinations...)
	}
	if combinationsToCreate.HasDuplicates() {
		return nil, apio.NewRequestError(errors.New("request contains duplicate subrequest value combinations"), http.StatusBadRequest)
	}

	return combinationsToCreate, nil
}

// requestIsValid checks that the request meets the constraints of the rule
// Add additional constraint checks here in this method.
func validateCreateRequest(ctx context.Context, request CreateRequest, rule rule.AccessRule, requestArguments map[string]types.RequestArgument) error {
	if request.Timing.DurationSeconds > rule.TimeConstraints.MaxDurationSeconds {
		logger.Get(ctx).Errorw("error validating request", zap.Error(fmt.Errorf(fmt.Sprintf("durationSeconds: %d exceeds the maximum duration seconds: %d", request.Timing.DurationSeconds, rule.TimeConstraints.MaxDurationSeconds))))
		return &apio.APIError{
			Err:    errors.New("request validation failed"),
			Status: http.StatusBadRequest,
			Fields: []apio.FieldError{
				{
					Field: "timing.durationSeconds",
					Error: fmt.Sprintf("durationSeconds: %d exceeds the maximum duration seconds: %d", request.Timing.DurationSeconds, rule.TimeConstraints.MaxDurationSeconds),
				},
			},
		}
	}

	given := make(map[string]string)
	expected := make(map[string][]string)
	if request.With != nil {
		given = request.With
	}
	for k, v := range requestArguments {
		if v.RequiresSelection {
			options := make([]string, len(v.Options))
			for _, o := range v.Options {
				if o.Valid {
					options = append(options, o.Value)
				}
			}
			expected[k] = options
		}
	}

	// assert they are the same length.
	// the user provided the expected number of values based on the requestArguments
	if len(given) != len(expected) {
		logger.Get(ctx).Errorw("error validating request", zap.Error(fmt.Errorf("unexpected number of arguments in 'with' field")))

		return &apio.APIError{
			Err:    errors.New("request validation failed"),
			Status: http.StatusBadRequest,
			Fields: []apio.FieldError{
				{
					Field: "with",
					Error: "unexpected number of arguments in 'with' field",
				},
			},
		}
	}
	// assert that the given argument ids are expected and the the value is an allowed value
	for argumentId, allowedValues := range expected {
		givenArgumentValue, ok := given[argumentId]
		if !ok || !contains(allowedValues, givenArgumentValue) {
			logger.Get(ctx).Errorw("error validating request", zap.Error(fmt.Errorf(fmt.Sprintf("unexpected value given for argument %s in with field", argumentId))))

			return &apio.APIError{
				Err:    errors.New("request validation failed"),
				Status: http.StatusBadRequest,
				Fields: []apio.FieldError{
					{
						Field: "with",
						Error: fmt.Sprintf("unexpected value given for argument %s in with field", argumentId),
					},
				},
			}
		}
	}
	return nil
}

func contains(set []string, str string) bool {
	for _, s := range set {
		if s == str {
			return true
		}
	}
	return false
}
