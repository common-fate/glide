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
*
Healthcheck
The service has a mildly complex task, it needs to call out to all the deployments, and update the healthiness.

Then for each target group, update the validity of the deployments which are registered to it.

# Then, save it all to the database

# The healthcheck should save the response from the describe endpoint to the deployment item in the database

The validations should be applied as described in the milestone, and tests should be implemented to assert that the healthcheck service works as expected with different responses.

Procedure described:
- initialise a map keyed by target group id
- go over each target group deployment
- run a healthcheck
- update the healthiness of the deployment
- if the deployment is healthy, update the target group map to be healthy
- if the deployment is unhealthy, record the error and continue
- finally, iterrate over the target group map and update the target group healthiness
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

	targetgroupMap := make(map[string]bool)
	upsertItems := []ddb.Keyer{}

	// for each deployment, run a healthcheck
	// update the healthiness of the deployment
	for _, deploymentItem := range listTargetGroupDeployments.Result {
		// run a healthcheck
		// update the healthiness of the deployment
		log.Infof("Running healthcheck for deployment: %s", deploymentItem.ID)

		// Determine requirements/api for querying the deployment ⭐️⭐️⭐️⭐️
		// "The deployment lambda should respond with some data"
		runtime, err := pdk.GetRuntime(ctx, deploymentItem.FunctionARN)
		if err != nil {
			return err
		}
		describeRes, err := runtime.Describe(ctx)
		if err != nil {
			return err
		}

		// healthy, err := CheckIfHealthy(deploymentItem)

		healthy := true
		for _, diagnostic := range describeRes.ConfigValidation {

			// @TODO: determine whether we need to push the config validation to the deployment item ⭐️⭐️⭐️ i.e. deploymentItem.ActiveConfig

			if !diagnostic.Success {
				healthy = false
				break
			}
		}

		if err != nil {
			return err
		}

		// update the deployment
		deploymentItem.Healthy = healthy

		if healthy {
			// update the target group map to be healthy
			targetgroupMap[deploymentItem.ID] = healthy
		} else {
			// @TODO: record the error Diagnostics and continue  ⭐️⭐️⭐️⭐️
			deploymentItem.Diagnostics = []targetgroup.Diagnostic{}
		}

		upsertItems = append(upsertItems, &deploymentItem)
	}

	s.DB.PutBatch(ctx, upsertItems...)

	// now we can iterate over the target group map and update the target group healthiness
	for targetGroupID, healthy := range targetgroupMap {
		// get the target group
		getTargetGroup := storage.GetTargetGroup{ID: targetGroupID}
		_, err := s.DB.Query(ctx, &getTargetGroup)
		if err != nil {
			return err
		}
		if !healthy {
			// update the target group...
			// Q: target group doesn't store health?
		}
	}

	log.Info("completed checking health")

	return nil
}

// func CheckIfHealthy(deploymentItem targetgroup.Deployment) (bool, error) {

// 	// TODO: use me as input when polling the deployment
// 	// deployment.AWSAccount
// 	// deployment.AwsRegion
// 	// deployment.FunctionARN

// 	return true, nil
// }
