package healthchecksvc

import (
	"context"
	"fmt"
	"reflect"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// Service holds business logic relating to Access Requests.
type Service struct {
	DB ddb.Storage
}
type groupHandlerRoute struct {
	handler              handler.Handler
	group                target.Group
	route                target.Route
	groupAndHandlerExist bool
}

func (s *Service) Check(ctx context.Context) error {
	log := logger.Get(ctx)
	log.Info("starting to check health")
	hm := make(map[string]handler.Handler)
	gm := make(map[string]target.Group)
	handlerRoutes := make(map[string][]groupHandlerRoute)
	// get all deployments
	listHandlers := storage.ListHandlers{}
	_, err := s.DB.Query(ctx, &listHandlers)
	if err != nil {
		return err
	}
	for _, h := range listHandlers.Result {
		hm[h.ID] = h
	}

	listGroups := storage.ListTargetGroups{}
	_, err = s.DB.Query(ctx, &listGroups)
	if err != nil {
		return err
	}
	for _, g := range listGroups.Result {
		gm[g.ID] = g
	}
	listRoutes := storage.ListTargetRoutes{}
	_, err = s.DB.Query(ctx, &listRoutes)
	if err != nil {
		return err
	}

	for _, r := range listRoutes.Result {
		// A route is invalid if the group or the handler does not exist, this should only happen if we have an error in our API
		// allowing the handler to be deleted without deleting the routes as well.
		// or the group to be deleted without deleting the routes
		groupAndHandlerExist := true
		h, ok := hm[r.Handler]
		if !ok {
			groupAndHandlerExist = false
		}
		g, ok := gm[r.Group]
		if !ok {
			groupAndHandlerExist = false
		}
		handlerRoutes[r.Handler] = append(handlerRoutes[r.Handler], groupHandlerRoute{
			handler:              h,
			group:                g,
			route:                r,
			groupAndHandlerExist: groupAndHandlerExist,
		})
	}

	upsertItems := []ddb.Keyer{}

	// for each deployment, run a healthcheck
	// update the healthiness of the deployment
	for _, h := range listHandlers.Result {
		// run a healthcheck
		// update the healthiness of the deployment
		log.Infof("Running healthcheck for deployment: %s", h.ID)

		// clear previous diagnostics
		h.Diagnostics = []handler.Diagnostic{}

		// get the lambda runtime
		runtime, err := handler.GetRuntime(ctx, h)
		if err != nil {
			h.Healthy = false
			log.Warnf("Error getting lambda runtime: %s", h.ID)
			h.Diagnostics = append(h.Diagnostics, handler.Diagnostic{
				Level:   types.LogLevelERROR,
				Message: err.Error(),
			})
			upsertItems = append(upsertItems, &h)
			continue
		}
		// now we can call the describe endpoint
		describeRes, err := runtime.Describe(ctx)
		if err != nil {
			h.Healthy = false
			log.Warnf("Error running healthcheck for deployment: %s", h.ID)
			h.Diagnostics = append(h.Diagnostics, handler.Diagnostic{
				Level:   types.LogLevelERROR,
				Message: err.Error(),
			})
			upsertItems = append(upsertItems, &h)
			continue
		}

		for _, diagnostic := range describeRes.Diagnostics {
			h.Diagnostics = append(h.Diagnostics, handler.Diagnostic{
				Level:   types.LogLevel(diagnostic.Level),
				Message: diagnostic.Msg,
			})
		}

		h.ProviderDescription = describeRes
		h.Healthy = describeRes.Healthy

		for _, handlerRoute := range handlerRoutes[h.ID] {
			route := handlerRoute.route
			// clear existing diagnostics
			route.Diagnostics = []target.Diagnostic{}
			if handlerRoute.groupAndHandlerExist {
				kindSchema, ok := h.ProviderDescription.Schema.Target.AdditionalProperties[handlerRoute.route.Kind]
				if ok {
					route.Valid = s.validateProviderSchema(handlerRoute.group.TargetSchema.Schema.AdditionalProperties, kindSchema.Schema.AdditionalProperties)
				} else {
					// invalid route because the kind does not exist in the schema
					route.Valid = false
					route.Diagnostics = append(route.Diagnostics, target.Diagnostic{
						Level:   "error",
						Message: fmt.Sprintf("kind schema '%s' does not exist for this handler", route.Kind),
					})
				}
			} else {
				route.Diagnostics = append(route.Diagnostics, target.Diagnostic{
					Level:   "error",
					Message: "route group does not exist",
				})
			}
			// add the route item to be updated
			upsertItems = append(upsertItems, &route)
		}

		upsertItems = append(upsertItems, &h)
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
