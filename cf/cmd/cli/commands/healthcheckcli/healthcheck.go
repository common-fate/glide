package healthcheckcli

import (
	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/healthcheck"
	"github.com/common-fate/common-fate/pkg/service/healthchecksvc"
	"github.com/common-fate/ddb"
	"github.com/urfave/cli/v2"
)

// this command can be run in dev with:
// go run cf/cmd/cli/main.go healthcheck

var Command = cli.Command{
	Name:        "healthcheck",
	Description: "healthcheck a deployment",
	Usage:       "healthcheck a deployment",
	Action: cli.ActionFunc(func(c *cli.Context) error {

		ctx := c.Context

		// opts := []types.ClientOption{}

		do, err := deploy.LoadConfig(deploy.DefaultFilename)
		if err != nil {
			return err
		}
		o, err := do.LoadOutput(ctx)
		if err != nil {
			return err
		}
		db, err := ddb.New(ctx, o.DynamoDBTable)
		if err != nil {
			return err
		}

		hc := healthcheck.HealthChecker{
			DB: db,
			HealthCheck: &healthchecksvc.Service{
				DB: db,
			},
		}

		// cfApi, err := types.NewClientWithResponses("http://0.0.0.0:8080", opts...)
		// if err != nil {
		// 	return err
		// }

		// run the health check synchronously ⭐️⭐️⭐️⭐️
		// find a way to invoke the health check and determin an api response
		// await the result

		// res, err := cfApi.AdminHealthCheckTargetGroupsWithResponse(ctx)
		// if err != nil {
		// 	return err
		// }

		err = hc.Check(ctx)
		if err != nil {
			return err
		}

		// now run a fetch
		// listRes, err := cfApi.ListTargetGroupDeploymentsWithResponse(ctx, nil)

		healthyCount := 0
		unhealthyCount := 0

		// for _, deployment := range listRes.JSON200.Res {
		// 	if deployment.Healthy {
		// 		healthyCount++
		// 	} else {
		// 		unhealthyCount++
		// 	}
		// }

		clio.Log("healthcheck result")
		clio.Log("healthy: %d", healthyCount)
		clio.Log("unhealthy: %d", unhealthyCount)

		return nil
	}),
}
