package api

import (
	"errors"
	"net/http"
	"sync"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/auth"
	"github.com/common-fate/granted-approvals/pkg/cache"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/common-fate/granted-approvals/pkg/types"
	"golang.org/x/sync/errgroup"
)

// "/api/v1/admin/requests"
func (a *API) AdminListRequests(w http.ResponseWriter, r *http.Request, params types.AdminListRequestsParams) {
	ctx := r.Context()

	var err error
	var dbRes []access.Request
	var qR *ddb.QueryResult
	var next *string

	queryOpts := []func(*ddb.QueryOpts){ddb.Limit(50)}
	if params.NextToken != nil {
		queryOpts = append(queryOpts, ddb.Page(*params.NextToken))
	}
	if params.Status != nil {
		q := storage.ListRequestsForStatus{Status: access.Status(*params.Status)}
		qR, err := a.DB.Query(ctx, &q, queryOpts...)
		if err == ddb.ErrNoItems {
			apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
			return
		}
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}
		if qR.NextPage != "" {
			next = &qR.NextPage
		}
		dbRes = q.Result
	} else {
		q := storage.ListRequests{}
		qR, err = a.DB.Query(ctx, &q, queryOpts...)

		if err == ddb.ErrNoItems {
			apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
			return
		}
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}

		if qR.NextPage != "" {
			next = &qR.NextPage
		}

		dbRes = q.Result
	}

	// var endToken int
	res := types.ListRequestsResponse{
		Requests: make([]types.Request, len(dbRes)),
	}

	for i, r := range dbRes {
		res.Requests[i] = r.ToAPI()
	}

	res.Next = next

	apio.JSON(ctx, w, res, http.StatusOK)
}

// Get a request
// (GET /api/v1/admin/requests/{requestId})
func (a *API) AdminGetRequest(w http.ResponseWriter, r *http.Request, requestId string) {
	ctx := r.Context()
	u := auth.UserFromContext(ctx)
	q := storage.GetRequest{ID: requestId}
	_, err := a.DB.Query(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.Error(ctx, w, apio.NewRequestError(err, http.StatusNotFound))
		return
	} else if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	qr := storage.GetAccessRuleVersion{ID: q.Result.Rule, VersionID: q.Result.RuleVersion}
	_, err = a.DB.Query(ctx, &qr)
	// Any error fetching the access rule is an internal server error because it should exist if the request exists
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	if qr.Result == nil {
		apio.Error(ctx, w, errors.New("access rule result was nil"))
		return
	}
	var options []cache.ProviderOption
	var mu sync.Mutex
	g, gctx := errgroup.WithContext(ctx)
	for k := range qr.Result.Target.WithSelectable {
		kCopy := k
		g.Go(func() error {
			// load from the cache, if the user has requested it, the cache is very likely to be valid
			_, opts, err := a.Cache.LoadCachedProviderArgOptions(gctx, qr.Result.Target.ProviderID, kCopy)
			if err != nil {
				return err
			}
			mu.Lock()
			defer mu.Unlock()
			options = append(options, opts...)
			return nil
		})
	}
	for k := range qr.Result.Target.With {
		kCopy := k
		g.Go(func() error {
			// load from the cache, if the user has requested it, the cache is very likely to be valid
			_, opts, err := a.Cache.LoadCachedProviderArgOptions(gctx, qr.Result.Target.ProviderID, kCopy)
			if err != nil {
				return err
			}
			mu.Lock()
			defer mu.Unlock()
			options = append(options, opts...)
			return nil
		})
	}
	err = g.Wait()
	if err != nil {
		apio.Error(ctx, w, err)
		return
	}
	apio.JSON(ctx, w, q.Result.ToAPIDetail(*qr.Result, q.Result.RequestedBy != u.ID, options), http.StatusOK)
}
