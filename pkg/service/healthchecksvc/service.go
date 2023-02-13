package healthchecksvc

import (
	"context"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Access Requests.
type Service struct {
	DB ddb.Storage
}

func (s *Service) Check(ctx context.Context) error {
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
