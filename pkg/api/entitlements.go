package api

import (
	"net/http"
	"strings"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

// List Entitlements
// (GET /api/v1/entitlements)
func (a *API) UserListEntitlements(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	res := types.ListEntitlementsResponse{
		Entitlements: []types.TargetKind{},
	}

	q := storage.ListTargetGroups{}
	err := a.DB.All(ctx, &q)
	if err == ddb.ErrNoItems {
		apio.JSON(ctx, w, res, http.StatusOK)
		return
	}
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}

	kinds := make(map[string]types.TargetKind)
	for _, t := range q.Result {
		key := t.From.Publisher + "#" + t.From.Name + "#" + t.From.Kind + "#"
		kinds[key] = types.TargetKind{
			Publisher: t.From.Publisher,
			Name:      t.From.Name,
			Kind:      t.From.Kind,
			Icon:      t.Icon,
		}
	}

	for _, k := range kinds {
		res.Entitlements = append(res.Entitlements, k)
	}
	apio.JSON(ctx, w, res, http.StatusOK)
}

// (GET /api/v1/entitlements/targets)
func (a *API) UserListEntitlementTargets(w http.ResponseWriter, r *http.Request, params types.UserListEntitlementTargetsParams) {
	ctx := r.Context()
	var opts []func(*ddb.QueryOpts)
	if params.NextToken != nil {
		opts = append(opts, ddb.Page(*params.NextToken))
	}

	res := types.ListTargetsResponse{
		Targets: []types.Target{},
	}

	var results []cache.Target
	var qo *ddb.QueryResult
	var err error
	if params.Kind != nil {
		// validation is handled for the kind param my a regex in the open API spec
		parts := strings.Split(*params.Kind, "/")
		q := storage.ListCachedTargetsForKind{
			Publisher: parts[0],
			Name:      parts[1],
			Kind:      parts[2],
		}
		qo, err = a.DB.Query(ctx, &q, opts...)
		if err == ddb.ErrNoItems {
			apio.JSON(ctx, w, res, http.StatusOK)

			return
		}
		if err != nil && err != ddb.ErrNoItems {
			apio.Error(ctx, w, err)
			return
		}
		results = q.Result
	} else {
		q := storage.ListCachedTargets{}
		qo, err = a.DB.Query(ctx, &q, opts...)
		if err == ddb.ErrNoItems {
			apio.JSON(ctx, w, res, http.StatusOK)

			return
		}
		if err != nil && err != ddb.ErrNoItems {
			apio.Error(ctx, w, err)
			return
		}
		results = q.Result
	}

	if qo.NextPage != "" {
		res.Next = &qo.NextPage
	}
	user := auth.UserFromContext(ctx)

	// Filtering needs to be done in the application layer because of limits with dynamoDB filters
	// in the end, the same amount of read units will be consumed

	for _, target := range cache.Filter(results, user.Groups) {
		res.Targets = append(res.Targets, target.ToAPI())
	}

	apio.JSON(ctx, w, res, http.StatusOK)

}
