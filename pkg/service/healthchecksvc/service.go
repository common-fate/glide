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
	"github.com/pkg/errors"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/runtime.go -package=mocks . Runtime
type Runtime interface {
	Describe(ctx context.Context) (describeResponse *providerregistrysdk.DescribeResponse, err error)
}

type RuntimeGetter interface {
	GetRuntime(ctx context.Context, handler handler.Handler) (Runtime, error)
}

type DefaultGetter struct {
}

func (DefaultGetter) GetRuntime(ctx context.Context, h handler.Handler) (Runtime, error) {
	return handler.GetRuntime(ctx, h)
}

// Service holds business logic relating to Access Requests.
type Service struct {
	DB ddb.Storage
	// Use DefaultGetter{}
	// This is interfaced so it can be mocked for testing
	RuntimeGetter RuntimeGetter
}
type groupRoute struct {
	group target.Group
	route target.Route
}
type handlerRouteMapping struct {
	handler     handler.Handler
	groupRoutes []groupRoute
}

// Handler routes is a helper method which maps the handlers, groups, routes together in a convenient data structure by fetching them from dynamo
// This method assumes that there is never a case where a route exists but the matching handler or group does not exist
// The services for deleteing a target group and a handler are written to delete associated routes so this case should not be possible
// handling for this case has therefor been omitted
func (s *Service) handlerRoutes(ctx context.Context) (map[string]handlerRouteMapping, error) {
	gm := make(map[string]target.Group)
	handlerRoutes := make(map[string]handlerRouteMapping)
	// get all deployments
	listHandlers := storage.ListHandlers{}
	_, err := s.DB.Query(ctx, &listHandlers)
	if err != nil {
		return nil, err
	}
	for _, h := range listHandlers.Result {
		handlerRoutes[h.ID] = handlerRouteMapping{
			handler:     h,
			groupRoutes: []groupRoute{},
		}
	}

	listGroups := storage.ListTargetGroups{}
	_, err = s.DB.Query(ctx, &listGroups)
	if err != nil {
		return nil, err
	}
	for _, g := range listGroups.Result {
		gm[g.ID] = g
	}
	listRoutes := storage.ListTargetRoutes{}
	_, err = s.DB.Query(ctx, &listRoutes)
	if err != nil {
		return nil, err
	}

	for _, r := range listRoutes.Result {
		g := gm[r.Group]
		hgr := handlerRoutes[r.Handler]
		hgr.groupRoutes = append(hgr.groupRoutes, groupRoute{
			group: g,
			route: r,
		})
		handlerRoutes[r.Handler] = hgr
	}
	return handlerRoutes, nil
}
func NewDiagFailedToInitialiseRuntime(err error) handler.Diagnostic {
	return handler.Diagnostic{
		Level:   types.LogLevelERROR,
		Message: errors.Wrap(err, "failed to initialise runtime for handler").Error(),
	}
}

// getRuntime will attempt to load the runtime for this handler
// it can possibly return a credential or other aws error when initialising the lambda client
// runtime will be nil if there was an error and the handler will be updated with the diagnostic logs
func (s *Service) getRuntime(ctx context.Context, h handler.Handler) (handler.Handler, Runtime) {
	log := logger.Get(ctx)
	// get the lambda runtime
	runtime, err := s.RuntimeGetter.GetRuntime(ctx, h)
	if err != nil {
		h.Healthy = false
		log.Errorw("Error getting runtime for handler", "handlerId", h.ID, "error", err)
		h.Diagnostics = append(h.Diagnostics, NewDiagFailedToInitialiseRuntime(err))
		return h, nil
	}
	return h, runtime
}

func NewDiagFailedToDescribe(err error) handler.Diagnostic {
	return handler.Diagnostic{
		Level:   types.LogLevelERROR,
		Message: errors.Wrap(err, "failed to describe handler").Error(),
	}
}

// Describe attempts to invoke the handler and process the response, if describing fails, the handler will be returned as unhealthy with a diagnostic log containing the invocation error
// If describing is successful, the handler will be returned and will reflect the actual state of the handlers health and any diagnostics logs
func describe(ctx context.Context, h handler.Handler, runtime Runtime) handler.Handler {
	log := logger.Get(ctx)
	// the runtime is available so now we can call the describe endpoint
	describeRes, err := runtime.Describe(ctx)
	if err != nil {
		h.Healthy = false
		log.Errorw("Error running healthcheck for handler", "id", h.ID, "error", err)
		h.Diagnostics = append(h.Diagnostics, NewDiagFailedToDescribe(err))
		return h
	}
	for _, diagnostic := range describeRes.Diagnostics {
		h.Diagnostics = append(h.Diagnostics, handler.Diagnostic{
			Level:   types.LogLevel(diagnostic.Level),
			Message: diagnostic.Msg,
		})
	}
	h.ProviderDescription = describeRes
	h.Healthy = describeRes.Healthy
	return h
}

var NewDiagHandlerUnreachable target.Diagnostic = target.Diagnostic{
	Level:   types.LogLevelERROR,
	Message: "handler is unreachable and route validity cannot be checked",
}

func NewDiagKindSchemaNotExist(route target.Route) target.Diagnostic {
	return target.Diagnostic{
		Level:   types.LogLevelERROR,
		Message: fmt.Sprintf("kind schema '%s' does not exist for the handler '%s'", route.Kind, route.Handler),
	}
}

// Validate route will assert that the handler description is available and that the schema of the handler is compatible with the schema of the target group for the route
func validateRoute(route target.Route, group target.Group, dr *providerregistrysdk.DescribeResponse) target.Route {
	// clear existing diagnostics
	route.Diagnostics = []target.Diagnostic{}

	if dr == nil {
		route.Valid = false
		route.Diagnostics = append(route.Diagnostics, NewDiagHandlerUnreachable)
		return route
	}

	if dr.Schema.Targets == nil {
		// provider doesn't provide any targets

		route.Valid = false
		route.Diagnostics = append(route.Diagnostics, NewDiagHandlerUnreachable)
		return route
	}

	// Check first that the target schema defines the Kind of the route, if not then the route is invalid
	targets := *dr.Schema.Targets
	kindSchema, ok := targets[route.Kind]
	if ok {
		route.Valid = validateProviderSchema(group.TargetSchema.Schema.Properties, kindSchema.Properties)
	} else {
		// invalid route because the kind does not exist in the schema
		route.Valid = false
		route.Diagnostics = append(route.Diagnostics, NewDiagKindSchemaNotExist(route))
	}
	return route
}

// validateProviderSchema asserts that the target schemas are equivalent in structure, comparing only the keys and value types
func validateProviderSchema(schema1 map[string]providerregistrysdk.TargetField, schema2 map[string]providerregistrysdk.TargetField) bool {
	var in = []map[string]providerregistrysdk.TargetField{schema1, schema2}
	var compare = make([]map[string]*string, 2)
	for i := range compare {
		m := make(map[string]*string)
		for key, arg := range in[i] {
			m[key] = arg.Resource
		}
		compare[i] = m
	}
	return reflect.DeepEqual(compare[0], compare[1])
}

// The logic in the Check method is split out into small steps, these steps can be more easily unit tested
// testing of this method as a whole should be implemented as an integration test
func (s *Service) Check(ctx context.Context) error {
	log := logger.Get(ctx)
	log.Info("starting to check health")
	handlerRoutes, err := s.handlerRoutes(ctx)
	if err != nil {
		return err
	}
	upsertItems := []ddb.Keyer{}

	// for each deployment, run a healthcheck
	// update the healthiness of the deployment
	for _, hr := range handlerRoutes {
		h := hr.handler
		// run a healthcheck
		// update the healthiness of the deployment
		log.Infof("Running healthcheck for deployment: %s", h.ID)

		// clear previous diagnostics
		h.Diagnostics = []handler.Diagnostic{}
		h.ProviderDescription = nil
		// Get the handler to describe the lambda
		var runtime Runtime

		h, runtime = s.getRuntime(ctx, h)
		// If the runtime is not nil, then we can describe
		// else, when the routes are validated, they will be marked as invalid
		if runtime != nil {
			// Next describe the provider, if there is an error describing, then the handler will be returned with diagnostics logs and no providerDescription
			h = describe(ctx, h, runtime)
		}

		upsertItems = append(upsertItems, &h)

		// Next validate the routes against the description, if it is nil, then the routes will all be marked invalid
		for _, groupRoute := range hr.groupRoutes {
			route := validateRoute(groupRoute.route, groupRoute.group, h.ProviderDescription)
			// add the route item to be updated
			upsertItems = append(upsertItems, &route)
		}
	}

	err = s.DB.PutBatch(ctx, upsertItems...)
	if err != nil {
		return err
	}
	log.Info("completed checking health")
	return nil
}
