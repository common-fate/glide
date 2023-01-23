package datadog

import (
	"context"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"go.uber.org/zap"
)

// List options for arg
func (p *Provider) Options(ctx context.Context, arg string) (*types.ArgOptionsResponse, error) {
	switch arg {
	case "dashboard":
		log := zap.S().With("arg", arg)
		log.Info("getting dashboard options")

		dapi := datadogV1.NewDashboardsApi(p.apiClient)
		res, _, err := dapi.ListDashboards(p.DDContext(ctx))
		if err != nil {
			return nil, err
		}

		var opts types.ArgOptionsResponse
		for _, d := range res.Dashboards {
			opts.Options = append(opts.Options, types.Option{Label: d.GetTitle(), Value: d.GetId()})
		}
		return &opts, nil
	}
	return nil, &providers.InvalidArgumentError{Arg: arg}
}
