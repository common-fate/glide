package grants

import (
	"net/http"
	"time"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/iso8601"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

var CreateCommand = cli.Command{
	Name: "create",
	Action: func(c *cli.Context) error {
		api, err := types.NewClientWithResponses(c.String("api-url"))
		if err != nil {
			return err
		}

		b := types.PostGrantsJSONRequestBody{
			Subject:  "chris@commonfate.io",
			Start:    iso8601.New(time.Now().Add(time.Second * 2)),
			End:      iso8601.New(time.Now().Add(time.Hour)),
			Provider: "cf-dev",
			With: types.CreateGrant_With{
				AdditionalProperties: map[string]string{"accountId": "123451234512"},
			},
		}

		res, err := api.PostGrantsWithResponse(c.Context, b)
		if err != nil {
			return err
		}

		if res.StatusCode() == http.StatusCreated {
			zap.S().Infow("created grant", "grant", res.JSON201.Grant)
		} else {
			zap.S().Infow("error creating grant", "error", res.JSON400.Error)
		}

		return nil
	},
}
