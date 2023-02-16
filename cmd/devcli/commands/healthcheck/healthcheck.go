package healthcheck

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/common-fate/clio"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/service/healthchecksvc"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/urfave/cli/v2"
)

var a = `{"provider": {"publisher": "", "name": "", "version": ""}, "config": {"api_url": "https://prod.testvault.granted.run", "unique_vault_id": "2FeRHElazlJsHYmkaV5Xtg53r8R"}, "configValidation": {}, "schema": {"target": {"Default": {"schema": {"vault": {"id": "vault", "type": "string", "title": "Vault", "options": false, "groups": null, "ruleFormElement": "INPUT", "resourceName": null, "description": "The name of an example vault to grant access to (can be any string)"}}}}, "audit": {"resourceLoaders": {}, "resources": {}}, "config": {"api_url": {"type": "string", "usage": "the test vault URL", "secret": false}, "unique_vault_id": {"type": "string", "usage": "the unique vault ID", "secret": false}}}}`
var cc = `{"provider": {"publisher": "josh", "name": "example", "version": "v0.1.4"}, "config": {"api_url": "https://prod.testvault.granted.run", "unique_vault_id": "2FeRHElazlJsHYmkaV5Xtg53r8R"}, "configValidation": {}, "schema": {"target": {"Default": {"schema": {"vault": {"id": "vault", "type": "string", "title": "Vault", "options": false, "groups": null, "ruleFormElement": "INPUT", "resourceName": null, "description": "The name of an example vault to grant access to (can be any string)"}}}}, "audit": {"resourceLoaders": {}, "resources": {}}, "config": {"api_url": {"type": "string", "usage": "the test vault URL", "secret": false}, "unique_vault_id": {"type": "string", "usage": "the unique vault ID", "secret": false}}}}`

// this command can be run in dev with:
// go run cf/cmd/cli/main.go healthcheck

var Command = cli.Command{
	Name:        "healthcheck",
	Description: "healthcheck a deployment",
	Usage:       "healthcheck a deployment",
	Flags:       []cli.Flag{&cli.StringSliceFlag{Name: "deployment-mappings"}},
	Action: cli.ActionFunc(func(c *cli.Context) error {
		var b providerregistrysdk.DescribeResponse
		err := json.Unmarshal([]byte(cc), &b)
		if err != nil {
			return err
		}
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
