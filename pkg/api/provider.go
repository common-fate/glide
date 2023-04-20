package api

import (
	"context"
	"net/http"

	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/types"
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
	// ctx := r.Context()

	// targetGroups := a.fetchTargetGroups(ctx)

	// combinedResponse := []types.Provider{}

	// for _, target := range targetGroups {
	// 	combinedResponse = append(combinedResponse, types.Provider{
	// 		Id:   target.Id,
	// 		Type: target.Icon,
	// 	})
	// }
	// apio.JSON(ctx, w, combinedResponse, http.StatusOK)
	// return

}
