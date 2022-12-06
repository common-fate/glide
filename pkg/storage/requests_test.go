package storage

import (
	"context"
	"testing"
	"time"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbtest"
	"github.com/stretchr/testify/assert"
)

func exampleRequest() access.Request {
	reason := "example"
	return access.Request{
		ID:          types.NewRequestID(),
		RequestedBy: "testuser",
		Rule:        "testrule",
		Status:      access.PENDING,
		Data: access.RequestData{
			Reason: &reason,
		},
		RequestedTiming: access.Timing{
			Duration: time.Minute * 5,
		},
		CreatedAt: time.Now().In(time.UTC),
	}
}

func TestGetRequest(t *testing.T) {
	s := newTestingStorage(t)
	// set up test fixture data.
	req1 := exampleRequest()

	requests := []access.Request{req1}
	ddbtest.PutFixtures(t, s, requests)

	type testcase struct {
		name    string
		giveID  string
		want    *access.Request
		wantErr error
	}

	testcases := []testcase{
		{
			name:   "ok",
			giveID: req1.ID,
			want:   &req1,
		},
		{
			name:    "request not exist",
			giveID:  types.NewRequestID(),
			wantErr: ddb.ErrNoItems,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			q := GetRequest{ID: tc.giveID}
			_, err := s.Query(ctx, &q)
			if err != nil && tc.wantErr == nil {
				t.Fatal(err)
			}
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			}

			assert.Equal(t, tc.want, q.Result)
		})
	}
}
func TestListRequests(t *testing.T) {
	s := newTestingStorage(t)
	// set up test fixture data.
	req1 := exampleRequest()
	req2 := exampleRequest()

	requests := []access.Request{req1, req2}
	ddbtest.PutFixtures(t, s, requests)

	type testcase struct {
		name    string
		want    []access.Request
		wantErr error
	}

	testcases := []testcase{
		{
			name: "ok",
			want: requests,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			q := ListRequests{}
			_, err := s.Query(ctx, &q)
			if err != nil && tc.wantErr == nil {
				t.Fatal(err)
			}
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			}

			for _, item := range tc.want {
				assert.Contains(t, q.Result, item)
			}
		})
	}
}
func TestListRequestsForStatus(t *testing.T) {
	s := newTestingStorage(t)
	// set up test fixture data.
	req1 := exampleRequest()
	req2 := exampleRequest()
	req3 := exampleRequest()
	req4 := exampleRequest()
	req1.Status = access.APPROVED
	req2.Status = access.CANCELLED
	req3.Status = access.DECLINED
	req4.Status = access.PENDING

	requests := []access.Request{req1, req2, req3, req4}
	ddbtest.PutFixtures(t, s, requests)

	type testcase struct {
		name    string
		status  access.Status
		want    []access.Request
		notWant []access.Request
		wantErr error
	}

	testcases := []testcase{
		{
			name:    "pending",
			status:  access.PENDING,
			want:    []access.Request{req4},
			notWant: []access.Request{req1, req2, req3},
		},
		{
			name:    "approved",
			status:  access.APPROVED,
			want:    []access.Request{req1},
			notWant: []access.Request{req2, req4, req3},
		},
		{
			name:    "cancelled",
			status:  access.CANCELLED,
			want:    []access.Request{req2},
			notWant: []access.Request{req1, req4, req3},
		},
		{
			name:    "declined",
			status:  access.DECLINED,
			want:    []access.Request{req3},
			notWant: []access.Request{req1, req2, req4},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			q := ListRequestsForStatus{Status: tc.status}
			_, err := s.Query(ctx, &q)
			if err != nil && tc.wantErr == nil {
				t.Fatal(err)
			}

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			}

			for _, item := range tc.want {
				assert.Contains(t, q.Result, item)
			}
			for _, item := range tc.notWant {
				assert.NotContains(t, q.Result, item)
			}

		})
	}
}
func TestListRequestsForUser(t *testing.T) {
	s := newTestingStorage(t)
	// set up test fixture data.
	req1 := exampleRequest()
	req2 := exampleRequest()
	req3 := exampleRequest()

	user1 := types.NewUserID()
	user2 := types.NewUserID()
	req1.RequestedBy = user1
	req2.RequestedBy = user1
	req3.RequestedBy = user2

	requests := []access.Request{req1, req2, req3}
	ddbtest.PutFixtures(t, s, requests)

	type testcase struct {
		name       string
		giveUserID string
		want       []access.Request
		notWant    []access.Request
		wantErr    error
	}

	testcases := []testcase{
		{
			name:       "user1",
			giveUserID: user1,
			want:       []access.Request{req1, req2},
			notWant:    []access.Request{req3},
		},
		{
			name:       "user2",
			giveUserID: user2,
			want:       []access.Request{req3},
			notWant:    []access.Request{req2, req1},
		},
		{
			name:       "none",
			giveUserID: types.NewUserID(),
			want:       []access.Request{},
			notWant:    []access.Request{req1, req2, req3},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			q := ListRequestsForUser{UserId: tc.giveUserID}
			_, err := s.Query(ctx, &q)
			if err != nil && tc.wantErr == nil {
				t.Fatal(err)
			}

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			}

			for _, item := range tc.want {
				assert.Contains(t, q.Result, item)
			}
			for _, item := range tc.notWant {
				assert.NotContains(t, q.Result, item)
			}

		})
	}
}
func TestListRequestsForUserAndStatus(t *testing.T) {
	s := newTestingStorage(t)
	// set up test fixture data.
	req1 := exampleRequest()
	req2 := exampleRequest()

	user1 := types.NewUserID()
	req1.RequestedBy = user1
	req1.Status = access.PENDING
	req2.RequestedBy = user1
	req2.Status = access.APPROVED

	requests := []access.Request{req1, req2}
	ddbtest.PutFixtures(t, s, requests)

	type testcase struct {
		name       string
		giveUserID string
		giveStatus access.Status
		want       []access.Request
		notWant    []access.Request
		wantErr    error
	}

	testcases := []testcase{
		{
			name:       "approved",
			giveUserID: user1,
			giveStatus: access.APPROVED,
			want:       []access.Request{req2},
			notWant:    []access.Request{req1},
		},
		{
			name:       "pending",
			giveUserID: user1,
			giveStatus: access.PENDING,
			want:       []access.Request{req1},
			notWant:    []access.Request{req2},
		},
		{
			name:       "no declined requests",
			giveUserID: user1,
			giveStatus: access.DECLINED,
			want:       []access.Request{},
			notWant:    []access.Request{req1, req2},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			q := ListRequestsForUserAndStatus{UserId: tc.giveUserID, Status: tc.giveStatus}
			_, err := s.Query(ctx, &q)
			if err != nil && tc.wantErr == nil {
				t.Fatal(err)
			}

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
			}

			for _, item := range tc.want {
				assert.Contains(t, q.Result, item)
			}
			for _, item := range tc.notWant {
				assert.NotContains(t, q.Result, item)
			}

		})
	}
}
func TestListRequestReviewers(t *testing.T) {
	s := newTestingStorage(t)

	// define the function to test
	testListRequestReviewers := func(requestID string, want []access.Reviewer) func(t *testing.T) {
		return func(t *testing.T) {
			ctx := context.Background()
			q := ListRequestReviewers{RequestID: requestID}
			_, err := s.Query(ctx, &q)
			if err != nil {
				t.Fatal(err)
			}
			assert.ElementsMatch(t, want, q.Result)
		}
	}

	// set up test fixture data.
	req1 := types.NewRequestID()
	req1Reviews := []access.Reviewer{
		{
			ReviewerID: types.NewUserID(),
			Request: access.Request{
				ID:     req1,
				Status: access.PENDING,
			},
		},
		{
			ReviewerID: types.NewUserID(),
			Request: access.Request{
				ID:     req1,
				Status: access.PENDING,
			},
		},
	}

	req2 := types.NewRequestID()
	req2Reviews := []access.Reviewer{
		{
			ReviewerID: types.NewUserID(),
			Request: access.Request{
				ID:     req2,
				Status: access.PENDING,
			},
		},
	}

	reviewers := append(req1Reviews, req2Reviews...)
	ddbtest.PutFixtures(t, s, reviewers)

	t.Run("req1", testListRequestReviewers(req1, req1Reviews))
	t.Run("req2", testListRequestReviewers(req2, req2Reviews))
	t.Run("request not exist", testListRequestReviewers(types.NewRequestID(), []access.Reviewer{}))
}

func TestListRequestsForReviewer(t *testing.T) {
	s := newTestingStorage(t)

	// define the function to test
	testListRequestsForReviewer := func(reviewerID string, want []access.Reviewer) func(t *testing.T) {
		return func(t *testing.T) {
			var wantRequests []access.Request
			for _, r := range want {
				wantRequests = append(wantRequests, r.Request)
			}
			ctx := context.Background()
			q := ListRequestsForReviewer{ReviewerID: reviewerID}
			_, err := s.Query(ctx, &q)
			if err != nil {
				t.Fatal(err)
			}
			assert.ElementsMatch(t, wantRequests, q.Result)
		}
	}

	// set up test fixture data.
	user1 := types.NewUserID()
	user1Reviewers := []access.Reviewer{
		{
			ReviewerID: user1,
			Request: access.Request{
				ID: types.NewRequestID(),
			},
		},
		{
			ReviewerID: user1,
			Request: access.Request{
				ID: types.NewRequestID(),
			},
		},
	}

	user2 := types.NewUserID()
	user2Reviewers := []access.Reviewer{
		{
			ReviewerID: user2,
			Request: access.Request{
				ID: types.NewRequestID(),
			},
		},
	}

	reviewers := append(user1Reviewers, user2Reviewers...)
	ddbtest.PutFixtures(t, s, reviewers)

	t.Run("user1", testListRequestsForReviewer(user1, user1Reviewers))
	t.Run("user2", testListRequestsForReviewer(user2, user2Reviewers))
	t.Run("request not exist", testListRequestsForReviewer(types.NewRequestID(), []access.Reviewer{}))
}

func TestListRequestsForReviewerAndStatus(t *testing.T) {
	s := newTestingStorage(t)

	// define the function to test
	testListRequestsForReviewerAndStatus := func(reviewerID string, status access.Status, want []access.Reviewer) func(t *testing.T) {
		return func(t *testing.T) {
			var wantRequests []access.Request
			for _, r := range want {
				wantRequests = append(wantRequests, r.Request)
			}
			ctx := context.Background()
			q := ListRequestsForReviewerAndStatus{ReviewerID: reviewerID, Status: status}
			_, err := s.Query(ctx, &q)
			if err != nil {
				t.Fatal(err)
			}
			assert.ElementsMatch(t, wantRequests, q.Result)
		}
	}

	// set up test fixture data.
	user1 := types.NewUserID()
	user1Reviewers := []access.Reviewer{
		{
			ReviewerID: user1,
			Request: access.Request{
				ID:     types.NewRequestID(),
				Status: access.APPROVED,
			},
		},
		{
			ReviewerID: user1,
			Request: access.Request{
				ID:     types.NewRequestID(),
				Status: access.PENDING,
			},
		},
	}

	user2 := types.NewUserID()
	user2Reviewers := []access.Reviewer{
		{
			ReviewerID: user2,
			Request: access.Request{
				ID:     types.NewRequestID(),
				Status: access.CANCELLED,
			},
		},
	}

	reviewers := append(user1Reviewers, user2Reviewers...)
	ddbtest.PutFixtures(t, s, reviewers)

	t.Run("user1 approved", testListRequestsForReviewerAndStatus(user1, access.APPROVED, user1Reviewers[:1]))
	t.Run("user1 pending", testListRequestsForReviewerAndStatus(user1, access.PENDING, user1Reviewers[1:]))
	t.Run("user2", testListRequestsForReviewerAndStatus(user2, access.CANCELLED, user2Reviewers))
	t.Run("user not exist", testListRequestsForReviewerAndStatus(types.NewRequestID(), access.PENDING, []access.Reviewer{}))
}
