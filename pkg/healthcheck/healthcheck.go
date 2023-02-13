package healthcheck

import (
	"context"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/service/healthchecksvc"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/ddb"
)

type HealthChecker struct {
	DB          ddb.Storage
	HealthCheck *healthchecksvc.Service
}

/*

Healthcheck
The service has a mildly complex task, it needs to call out to all the deployments, and update the healthiness.

Then for each target group, update the validity of the deployments which are registered to it.

# Then, save it all to the database

# The healthcheck should save the response from the describe endpoint to the deployment item in the database

The validations should be applied as described in the milestone, and tests should be implemented to assert that the healthcheck service works as expected with different responses.

Procedure described:
- go over each target group deployment
- run a healthcheck
- update the healthiness of the deployment
- if the deployment is unhealthy, record the error and continue
*/

func (s *HealthChecker) Check(ctx context.Context) error {
	log := logger.Get(ctx)
	log.Info("starting to check health")

	// get all deployments
	listTargetGroupDeployments := storage.ListTargetGroupDeployments{}

	_, err := s.DB.Query(ctx, &listTargetGroupDeployments)
	if err != nil {
		return err
	}

	upsertItems := []ddb.Keyer{}

	// for each deployment, run a healthcheck
	// update the healthiness of the deployment
	for _, deploymentItem := range listTargetGroupDeployments.Result {
		// run a healthcheck
		// update the healthiness of the deployment
		log.Infof("Running healthcheck for deployment: %s", deploymentItem.ID)

		// get the lambda runtime
		runtime, err := pdk.GetRuntime(ctx, deploymentItem.FunctionARN)
		if err != nil {
			return err
		}
		// now we can call the describe endpoint
		describeRes, err := runtime.Describe(ctx)
		if err != nil {
			/**
			[✘] operation error Lambda: Invoke, https response error StatusCode: 404, RequestID: e630c31d-e611-4e31-a02d-6aef0aa7ac7f, ResourceNotFoundException: Functions from 'us-east-1' are not reachable in this region ('ap-southeast-2')
			*/
			return err
		}

		/**
		What we have here:
		- healthy response that defaults to any error
		- every config validation diagnostic stacked onto the one deploymentItem.Diagnostics field

		What we probably want:
		- an improved deploymentItem.Diagnostics field that is a map data type?? ⭐️⭐️
		*/

		// if there is an unhealthy config validation, then the deployment is unhealthy
		healthy := true
		for _, diagnostic := range describeRes.ConfigValidation {
			for _, d := range diagnostic.Logs {
				deploymentItem.Diagnostics = append(deploymentItem.Diagnostics, targetgroup.Diagnostic{
					Level:   d.Level,
					Message: d.Message,
				})
			}
			if !diagnostic.Success {
				healthy = false
				break
			}
		}

		// update the deployment
		deploymentItem.Healthy = healthy

		upsertItems = append(upsertItems, &deploymentItem)
	}

	s.DB.PutBatch(ctx, upsertItems...)

	log.Info("completed checking health")

	return nil
}
