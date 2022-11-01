package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/service/cognitosvc"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
)

// Your GET endpoint
// (GET /api/v1/groups/)
func (a *API) GetGroups(w http.ResponseWriter, r *http.Request, params types.GetGroupsParams) {
	ctx := r.Context()

	queryOpts := []func(*ddb.QueryOpts){ddb.Limit(50)}
	if params.NextToken != nil {
		queryOpts = append(queryOpts, ddb.Page(*params.NextToken))
	}

	status := types.IdpStatusACTIVE

	if params.Status != nil {
		status = types.IdpStatus(*params.Status)
	}

	q := storage.ListGroupsForStatus{
		Status: status,
	}

	qr, err := a.DB.Query(ctx, &q, queryOpts...)

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	res := types.ListGroupsResponse{
		Groups: make([]types.Group, len(q.Result)),
	}
	if qr != nil && qr.NextPage != "" {
		res.Next = &qr.NextPage
	}

	for i, g := range q.Result {
		res.Groups[i] = g.ToAPI()
	}

	apio.JSON(ctx, w, res, http.StatusOK)
}

// Get Group Details
// (GET /api/v1/groups/{groupId})
func (a *API) GetGroup(w http.ResponseWriter, r *http.Request, groupId string) {
	ctx := r.Context()

	q := storage.GetGroup{ID: groupId}

	_, err := a.DB.Query(ctx, &q)
	// return a 404 if the user was not found.
	if errors.As(err, &identity.UserNotFoundError{}) {
		err = apio.NewRequestError(err, http.StatusNotFound)
	}

	if err != nil {
		apio.Error(ctx, w, err)
		return
	}

	apio.JSON(ctx, w, q.Result.ToAPI(), http.StatusOK)
}

// Create Group
// (POST /api/v1/admin/groups)
// Creates a group in cognito if cognito is the selected idp, otherwise will create an Internal Granted group
func (a *API) CreateGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var createGroupRequest types.CreateGroupJSONRequestBody
	err := apio.DecodeJSONBody(w, r, &createGroupRequest)
	if err != nil {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
		return
	}
	if a.Cognito == nil {
		//create cognito group

		group, err := a.Cognito.CreateGroup(ctx, cognitosvc.CreateGroupOpts{
			Name:        createGroupRequest.Name,
			Description: aws.ToString(createGroupRequest.Description),
		})
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		apio.JSON(ctx, w, group.ToAPI(), http.StatusCreated)
	} else {
		// Create internal group
		group := identity.Group{
			ID:          createGroupRequest.Name,
			IdpID:       createGroupRequest.Name,
			Name:        createGroupRequest.Name,
			Description: *createGroupRequest.Description,
			Status:      types.IdpStatusINTERNAL,
			Users:       *createGroupRequest.Members,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err = a.DB.Put(ctx, &group)
		if err != nil {
			apio.Error(ctx, w, apio.NewRequestError(err, http.StatusBadRequest))
			return
		}
	}

}
