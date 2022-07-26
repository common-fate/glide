package accesssvc

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/gevent"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/service/grantsvc"
	"github.com/common-fate/granted-approvals/pkg/service/rulesvc"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/storage/dbupdate"
	"github.com/common-fate/granted-approvals/pkg/types"
)

type CreateRequestResult struct {
	Request   access.Request
	Reviewers []access.Reviewer
}

// CreateRequest creates a new request and saves it in the database.
// Returns an error if the request is invalid.
func (s *Service) CreateRequest(ctx context.Context, user *identity.User, in types.CreateRequestRequest) (*CreateRequestResult, error) {
	log := logger.Get(ctx).With("user.id", user.ID)
	q := storage.GetAccessRuleCurrent{ID: in.AccessRuleId}
	_, err := s.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		return nil, ErrRuleNotFound
	}
	if err != nil {
		// we don't know how to handle the error from the rule getter, so just return it to the caller.
		return nil, err
	}
	rule := q.Result

	log.Debugw("verifying user belongs to access rule groups", "rule.groups", rule.Groups, "user.groups", user.Groups)
	err = groupMatches(rule.Groups, user.Groups)
	if err != nil {
		return nil, err
	}

	now := s.Clock.Now()
	err = requestIsValid(in, rule)
	if err != nil {
		return nil, err
	}
	// the request is valid, so create it.
	req := access.Request{
		ID:          types.NewRequestID(),
		RequestedBy: user.ID,
		Data: access.RequestData{
			Reason: in.Reason,
		},
		CreatedAt:       now,
		UpdatedAt:       now,
		Status:          access.PENDING,
		RequestedTiming: access.TimingFromRequestTiming(in.Timing),
		Rule:            rule.ID,
		RuleVersion:     rule.Version,
	}

	// If the approval is not required, auto-approve the request
	auto := types.AUTOMATIC
	revd := types.REVIEWED

	if !rule.Approval.IsRequired() {
		req.Status = access.APPROVED
		req.ApprovalMethod = &auto
	} else {
		req.ApprovalMethod = &revd
	}

	approvers, err := rulesvc.GetApprovers(ctx, s.DB, *rule)
	if err != nil {
		return nil, err
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
			// @TODO: add me here????
			// SlackMessageID: ,
		}

		reviewers = append(reviewers, r)
		items = append(items, &r)
	}

	log.Debugw("saving request", "request", req, "reviewers", reviewers)

	// audit log event
	reqEvent := access.NewRequestCreatedEvent(req.ID, req.CreatedAt, &req.RequestedBy)
	items = append(items, &reqEvent)
	// save the request.
	err = s.DB.PutBatch(ctx, items...)
	if err != nil {
		return nil, err
	}

	err = s.EventPutter.Put(ctx, gevent.RequestCreated{Request: req})
	// in a future PR we will shift these events out to be triggered by dynamo db streams
	// This will currently put the app in a strange state if this fails
	if err != nil {
		return nil, err
	}

	// check to see if it valid for instant approval
	if !rule.Approval.IsRequired() {
		log.Debugw("auto-approving", "request", req, "reviewers", reviewers)
		updatedReq, err := s.Granter.CreateGrant(ctx, grantsvc.CreateGrantOpts{Request: req, AccessRule: *rule})
		if err != nil {
			return nil, err
		}
		req = *updatedReq
		items, err := dbupdate.GetUpdateRequestItems(ctx, s.DB, req, dbupdate.WithReviewers(reviewers))
		if err != nil {
			return nil, err
		}
		err = s.DB.PutBatch(ctx, items...)
		if err != nil {
			return nil, err
		}
	}

	res := CreateRequestResult{
		Request:   req,
		Reviewers: reviewers,
	}

	return &res, nil
}

func groupMatches(ruleGroups []string, userGroups []string) error {
	for _, rg := range ruleGroups {
		for _, ug := range userGroups {
			if rg == ug {
				return nil
			}
		}
	}
	return ErrNoMatchingGroup
}

// requestIsValid checks that the request meets the constraints of the rule
// Add additional constraint checks here in this method.
func requestIsValid(request types.CreateRequestRequest, rule *rule.AccessRule) error {
	if request.Timing.DurationSeconds > rule.TimeConstraints.MaxDurationSeconds {
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
	return nil
}
