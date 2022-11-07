package accesssvc

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/gevent"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/service/grantsvc"
	"github.com/common-fate/granted-approvals/pkg/service/rulesvc"
	"github.com/common-fate/granted-approvals/pkg/storage/dbupdate"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/hashicorp/go-multierror"
)

type CreateRequestResult struct {
	Request   access.Request
	Reviewers []access.Reviewer
}

type CreateRequests struct {
	AccessRuleId string
	Reason       *string
	Timing       types.RequestTiming
	With         *types.CreateRequestWithSubRequest
}

type CreateRequestsOpts struct {
	User   identity.User
	Create CreateRequests
}

// CreateRequests splits the multi request into invividual request after checking for some basic validation errors
// individual requests may fail, these will be returned via a multi error and any requests which were successful will be returned as well
// so be sure to check both teh error and the response
func (s *Service) CreateRequests(ctx context.Context, in CreateRequestsOpts) ([]CreateRequestResult, error) {
	validated, err := s.validateCreateRequests(ctx, in)
	if err != nil {
		return nil, err
	}
	var createError multierror.Error
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []CreateRequestResult
	for _, combinationToCreate := range validated.argumentCombinations {
		wg.Add(1)
		go func(c map[string]string) {
			defer wg.Done()
			res, err := s.createRequest(ctx, createRequestOpts{
				User: in.User,
				Request: CreateRequest{
					AccessRuleId: in.Create.AccessRuleId,
					Reason:       in.Create.Reason,
					Timing:       in.Create.Timing,
					With:         c,
				},
				Rule:             validated.rule,
				RequestArguments: validated.requestArguments,
			})
			mu.Lock()
			if err != nil {
				createError.Errors = append(createError.Errors, err)
			} else {
				results = append(results, res)
			}
			mu.Unlock()
		}(combinationToCreate)
	}
	wg.Wait()
	if createError.Errors == nil {
		return results, nil
	}
	return results, &createError
}

type CreateRequest struct {
	AccessRuleId string
	Reason       *string
	Timing       types.RequestTiming
	With         map[string]string
}
type createRequestOpts struct {
	User             identity.User
	Request          CreateRequest
	Rule             rule.AccessRule
	RequestArguments map[string]types.RequestArgument
}

// createRequest creates a new request and saves it in the database.
func (s *Service) createRequest(ctx context.Context, in createRequestOpts) (CreateRequestResult, error) {
	now := s.Clock.Now()
	log := logger.Get(ctx).With("user.id", in.User.ID)
	// the request is valid, so create it.
	req := access.Request{
		ID:          types.NewRequestID(),
		RequestedBy: in.User.ID,
		Data: access.RequestData{
			Reason: in.Request.Reason,
		},
		CreatedAt:       now,
		UpdatedAt:       now,
		Status:          access.PENDING,
		RequestedTiming: access.TimingFromRequestTiming(in.Request.Timing),
		Rule:            in.Rule.ID,
		RuleVersion:     in.Rule.Version,
		SelectedWith:    make(map[string]access.Option),
	}
	if in.Request.With != nil {
		for k, v := range in.Request.With {
			argument := in.RequestArguments[k]
			found := false
			for _, option := range argument.Options {
				// because validation has passed, we can have certainty that the matching value will be found here
				// as a fallback, return an error if it is not found because something has gone seriously wrong with validation
				if option.Value == v {
					req.SelectedWith[k] = access.Option{
						Value:       option.Value,
						Label:       option.Label,
						Description: option.Description,
					}
					found = true
					break
				}
			}
			if !found {
				// this should never happen but here just in case
				return CreateRequestResult{}, errors.New("unexpected error, failed to find a matching option for a with argument when creating a new access request")
			}
		}
	}

	//validate the request against the access handler - make sure that access will be able to be provisioned
	//validating the grant before the request was made so that the request object does not get created.
	err := s.Granter.ValidateGrant(ctx, grantsvc.CreateGrantOpts{Request: req, AccessRule: in.Rule})
	if err != nil {
		return CreateRequestResult{}, err
	}

	// If the approval is not required, auto-approve the request
	auto := types.AUTOMATIC
	revd := types.REVIEWED

	if !in.Rule.Approval.IsRequired() {
		req.Status = access.APPROVED
		req.ApprovalMethod = &auto
	} else {
		req.ApprovalMethod = &revd
	}

	approvers, err := rulesvc.GetApprovers(ctx, s.DB, in.Rule)
	if err != nil {
		return CreateRequestResult{}, err
	}

	// track items to insert in the database.
	items := []ddb.Keyer{&req}

	// create Reviewers for each approver in the Access Rule. Reviewers will see the request in the End User portal.
	var reviewers []access.Reviewer
	for _, u := range approvers {
		// users cannot approve their own requests.
		// We don't create a Reviewer for them, even if they are an approver on the Access Rule.
		if u == req.RequestedBy {
			continue
		}

		r := access.Reviewer{
			ReviewerID: u,
			Request:    req,
		}

		reviewers = append(reviewers, r)
		items = append(items, &r)
	}

	log.Debugw("saving request", "request", req, "reviewers", reviewers)

	// audit log event
	reqEvent := access.NewRequestCreatedEvent(req.ID, req.CreatedAt, &req.RequestedBy)

	//before saving the request check to see if there already is a active approved rule
	if !in.Rule.Approval.IsRequired() {

		// This will check against the requests which do have grants already
		overlaps, err := s.overlapsExistingGrant(ctx, req)
		if err != nil {
			return CreateRequestResult{}, err
		}
		if overlaps {
			return CreateRequestResult{}, ErrRequestOverlapsExistingGrant
		}

	}

	items = append(items, &reqEvent)
	// save the request.
	err = s.DB.PutBatch(ctx, items...)
	if err != nil {
		return CreateRequestResult{}, err
	}

	err = s.EventPutter.Put(ctx, gevent.RequestCreated{Request: req})
	// in a future PR we will shift these events out to be triggered by dynamo db streams
	// This will currently put the app in a strange state if this fails
	if err != nil {
		return CreateRequestResult{}, err
	}

	// check to see if it valid for instant approval
	if !in.Rule.Approval.IsRequired() {
		log.Debugw("auto-approving", "request", req, "reviewers", reviewers)
		updatedReq, err := s.Granter.CreateGrant(ctx, grantsvc.CreateGrantOpts{Request: req, AccessRule: in.Rule})
		if err != nil {
			return CreateRequestResult{}, err
		}
		req = *updatedReq
		items, err := dbupdate.GetUpdateRequestItems(ctx, s.DB, req, dbupdate.WithReviewers(reviewers))
		if err != nil {
			return CreateRequestResult{}, err
		}
		err = s.DB.PutBatch(ctx, items...)
		if err != nil {
			return CreateRequestResult{}, err
		}
	}

	// analytics event
	analytics.FromContext(ctx).Track(&analytics.RequestCreated{
		RequestedBy:      req.RequestedBy,
		Provider:         rule.Target.ProviderType,
		RuleID:           req.Rule,
		Timing:           req.RequestedTiming.ToAnalytics(),
		HasReason:        req.HasReason(),
		RequiresApproval: in.Rule.Approval.IsRequired(),
	})

	return CreateRequestResult{
		Request:   req,
		Reviewers: reviewers,
	}, nil
}

func groupMatches(ruleGroups []string, userGroups []string) error {
	for _, rg := range ruleGroups {
		for _, ug := range userGroups {
			if rg == ug {
				return nil
			}
		}
	}
	return apio.NewRequestError(ErrNoMatchingGroup, http.StatusBadRequest)
}
