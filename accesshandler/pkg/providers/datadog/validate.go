package datadog

import (
	"context"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/common-fate/common-fate/accesshandler/pkg/diagnostics"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
)

func (p *Provider) TestConfig(ctx context.Context) error {
	return nil
}

func (p *Provider) ValidateConfig() map[string]providers.ConfigValidationStep {
	return map[string]providers.ConfigValidationStep{
		"list-dashboards": {
			Name: "List dashboards",
			Run: func(ctx context.Context) diagnostics.Logs {
				dapi := datadogV1.NewDashboardsApi(p.apiClient)
				res, _, err := dapi.ListDashboards(p.DDContext(ctx))
				if err != nil {
					return diagnostics.Error(err)
				}
				return diagnostics.Info("Datadog returned %d dashboards (more may exist, pagination has been ignored)", len(res.Dashboards))
			},
		},
	}
}
