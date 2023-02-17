package healthchecksvc

import (
	"context"
	"reflect"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/pdk"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
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

		// clear previous diagnostics
		deploymentItem.Diagnostics = []targetgroup.Diagnostic{}
		if deploymentItem.TargetGroupAssignment != nil {
			deploymentItem.TargetGroupAssignment.Diagnostics = []targetgroup.Diagnostic{}
		}

		// get the lambda runtime
		runtime, err := pdk.GetRuntime(ctx, deploymentItem)
		if err != nil {
			deploymentItem.Healthy = false
			log.Warnf("Error getting lambda runtime: %s", deploymentItem.ID)
			deploymentItem.Diagnostics = append(deploymentItem.Diagnostics, targetgroup.Diagnostic{
				Level:   string(types.ProviderSetupDiagnosticLogLevelERROR),
				Message: err.Error(),
			})
			upsertItems = append(upsertItems, &deploymentItem)
			continue
		}
		// now we can call the describe endpoint
		describeRes, err := runtime.Describe(ctx)
		if err != nil {
			deploymentItem.Healthy = false
			log.Warnf("Error running healthcheck for deployment: %s", deploymentItem.ID)
			deploymentItem.Diagnostics = append(deploymentItem.Diagnostics, targetgroup.Diagnostic{
				Level:   string(types.ProviderSetupDiagnosticLogLevelERROR),
				Message: err.Error(),
			})
			upsertItems = append(upsertItems, &deploymentItem)
			continue
		}

		/**
		What we have here:
		- healthy response that defaults to any error
		- every config validation diagnostic stacked onto the one deploymentItem.Diagnostics field

		What we probably want:
		- an improved deploymentItem.Diagnostics field that is a map data type??
		- break this down in the future
		*/

		// if there is an unhealthy config validation, then the deployment is unhealthy
		healthy := true
		for _, diagnostic := range describeRes.ConfigValidation.AdditionalProperties {
			for _, d := range diagnostic.Logs {
				deploymentItem.Diagnostics = append(deploymentItem.Diagnostics, targetgroup.Diagnostic{
					Level:   string(d.Level),
					Message: d.Msg,
				})
			}
			if !diagnostic.Success {
				healthy = false
			}
		}

		deploymentItem.ProviderDescription = describeRes
		deploymentItem.Healthy = healthy

		// @TODO validate target schema against targetgroup
		if deploymentItem.TargetGroupAssignment != nil {
			// This needs to be replaced with an actual check

			// need to do a structural comparison of the scheam for the provider in the describe against the schema stored on the target group

			//lookup target group
			targetGroup := storage.GetTargetGroup{ID: deploymentItem.TargetGroupAssignment.TargetGroupID}

			_, err := s.DB.Query(ctx, &targetGroup)
			if err != nil {
				return err
			}

			valid := s.validateProviderSchema(targetGroup.Result.TargetSchema.Schema.AdditionalProperties, describeRes.Schema.Target.AdditionalProperties["Default"].Schema.AdditionalProperties)
			deploymentItem.TargetGroupAssignment.Valid = valid

		}

		// update the deployment

		upsertItems = append(upsertItems, &deploymentItem)
	}

	err = s.DB.PutBatch(ctx, upsertItems...)
	if err != nil {
		return err
	}
	log.Info("completed checking health")
	return nil
}

func (s *Service) validateProviderSchema(schema1 map[string]providerregistrysdk.TargetArgument, schema2 map[string]providerregistrysdk.TargetArgument) bool {

	targetGroupSchemaMap := make(map[string]string)
	for _, arg := range schema1 {

		if arg.ResourceName == nil {
			targetGroupSchemaMap[arg.Id] = "string"

		} else {
			targetGroupSchemaMap[arg.Id] = *arg.ResourceName

		}
	}
	describeSchemaMap := make(map[string]string)
	for _, arg := range schema2 {
		if arg.ResourceName == nil {
			describeSchemaMap[arg.Id] = "string"

		} else {
			describeSchemaMap[arg.Id] = *arg.ResourceName

		}
	}

	return reflect.DeepEqual(describeSchemaMap, targetGroupSchemaMap)

	//do some sort of check here to validate that the schemas are the same and valid.
}
