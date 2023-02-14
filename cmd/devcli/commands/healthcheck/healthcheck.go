package healthcheck

import (
	"errors"
	"net/http"
	"strings"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/service/healthchecksvc"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/urfave/cli/v2"
)

// this command can be run in dev with:
// go run cf/cmd/cli/main.go healthcheck

var Command = cli.Command{
	Name:        "healthcheck",
	Description: "healthcheck a deployment",
	Usage:       "healthcheck a deployment",
	Flags:       []cli.Flag{&cli.StringSliceFlag{Name: "deployment-mappings"}},
	Action: cli.ActionFunc(func(c *cli.Context) error {

		ctx := c.Context

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
		for _, dm := range c.StringSlice("deployment-mappings") {
			kv := strings.Split(dm, ":")
			if len(kv) != 2 {
				return errors.New("deployment-mapping is invalid")
			}
			pdk.LocalDeploymentMap[kv[0]] = kv[1]
		}

		hc := healthchecksvc.Service{
			DB: db,
		}

		err = hc.Check(ctx)
		if err != nil {
			return err
		}

		opts := []types.ClientOption{}

		cfApi, err := types.NewClientWithResponses("http://0.0.0.0:8080", opts...)
		if err != nil {
			return err
		}

		// now run a fetch
		listRes, err := cfApi.ListTargetGroupDeploymentsWithResponse(ctx)
		if err != nil {
			return err
		}

		healthyCount := 0
		unhealthyCount := 0

		if listRes.StatusCode() != http.StatusOK {
			clio.Error(err)
			return errors.New("failed to list deployments from API")
		}

		for _, deployment := range listRes.JSON200.Res {
			if deployment.Healthy {
				healthyCount++
			} else {
				unhealthyCount++
			}
		}

		clio.Log("healthcheck result")
		clio.Logf("healthy: %d", healthyCount)
		clio.Logf("unhealthy: %d", unhealthyCount)

		return nil
	}),
}
