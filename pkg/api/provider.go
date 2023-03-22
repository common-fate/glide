package api

import (
	"context"
	"net/http"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

func (a *API) fetchTargetGroups(ctx context.Context) []types.TargetGroup {
	q := storage.ListTargetGroups{}

	_, err := a.DB.Query(ctx, &q)

	var targetGroups []types.TargetGroup
	// return empty slice if error
	if err != nil {
		return nil
	}

	for _, tg := range q.Result {
		targetGroups = append(targetGroups, tg.ToAPI())
	}

	return targetGroups
}

func (a *API) AdminListProviders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	targetGroups := a.fetchTargetGroups(ctx)

	combinedResponse := []types.Provider{}

	for _, target := range targetGroups {
		combinedResponse = append(combinedResponse, types.Provider{
			Id:   target.Id,
			Type: target.Icon,
		})
	}
	apio.JSON(ctx, w, combinedResponse, http.StatusOK)
	return

}

func (a *API) AdminGetProvider(w http.ResponseWriter, r *http.Request, providerId string) {
	ctx := r.Context()

	q := storage.GetTargetGroup{ID: providerId}
	_, err := a.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
		return
	}

	if q.Result != nil {
		apio.JSON(ctx, w,
			&types.Provider{
				Id:   q.Result.ID,
				Type: q.Result.Icon,
			}, http.StatusOK)
		return
	}

}

// helper method to check if the provided id is a valid target group.
func (a *API) isTargetGroup(ctx context.Context, targetGroupId string) bool {
	q := storage.GetTargetGroup{ID: targetGroupId}
	_, _ = a.DB.Query(ctx, &q)
	return q.Result != nil
}

func (a *API) AdminGetProviderArgs(w http.ResponseWriter, r *http.Request, providerId string) {
	ctx := r.Context()

	q := storage.GetTargetGroup{ID: providerId}
	_, err := a.DB.Query(ctx, &q)
	if err != nil && err != ddb.ErrNoItems {
		apio.Error(ctx, w, err)
	}

	// Convert the registry schema to the type required for the API
	if q.Result != nil {
		schema := types.ArgSchema{
			AdditionalProperties: map[string]types.Argument{},
		}
		for k, v := range q.Result.Schema.Properties {
			a := types.Argument{
				Id:           k,
				Description:  v.Description,
				ResourceName: v.Resource,
				Groups: &types.Argument_Groups{
					AdditionalProperties: map[string]types.Group{},
				},
				RuleFormElement: types.ArgumentRuleFormElementINPUT,
			}
			if v.Title != nil {
				a.Title = *v.Title
			}

			if v.Resource != nil {
				a.RuleFormElement = types.ArgumentRuleFormElementMULTISELECT
			}
			schema.AdditionalProperties[k] = a
		}

		apio.JSON(ctx, w, schema, http.StatusCreated)
		return
	}

}

func (a *API) fetchProviderResourcesByResourceType(ctx context.Context, providerId string, resourceType string) ([]types.Option, error) {
	cachedResources := storage.ListCachedTargetGroupResource{
		TargetGroupID: providerId,
		ResourceType:  resourceType,
	}

	_, err := a.DB.Query(ctx, &cachedResources)
	if err != nil && err != ddb.ErrNoItems {
		return []types.Option{}, err
	}

	opts := []types.Option{}
	for _, k := range cachedResources.Result {
		opts = append(opts, types.Option{
			Label: k.Resource.Name,
			Value: k.Resource.ID,
		})
	}

	return opts, nil
}

// List provider arg options
// (GET /api/v1/admin/providers/{providerId}/args/{argId}/options)
func (a *API) AdminListProviderArgOptions(w http.ResponseWriter, r *http.Request, providerId string, argId string, params types.AdminListProviderArgOptionsParams) {
	ctx := r.Context()

	res := types.ArgOptionsResponse{
		Options: []types.Option{},
		Groups:  &types.Groups{AdditionalProperties: make(map[string][]types.GroupOption)},
	}

	var err error
	if a.isTargetGroup(ctx, providerId) {
		// argId is either an argument's Id or resource's Name
		res.Options, err = a.fetchProviderResourcesByResourceType(ctx, providerId, argId)
		if err != nil {
			apio.Error(ctx, w, err)
			return
		}

		apio.JSON(ctx, w, res, http.StatusOK)
		return
	}

}

type ListProvidersArgFilterResponse struct {
	Options []types.Option `json:"options"`
}
