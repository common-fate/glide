// The api package defines all of our REST API endpoints.
package api

import (
	"context"
	"errors"
	"net/http"

	registry_types "github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/eventhandler"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/identity/identitysync"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/accesssvc"
	"github.com/common-fate/common-fate/pkg/service/cognitosvc"
	"github.com/common-fate/common-fate/pkg/service/handlersvc"
	"github.com/common-fate/common-fate/pkg/service/healthchecksvc"
	"github.com/common-fate/common-fate/pkg/service/internalidentitysvc"
	"github.com/common-fate/common-fate/pkg/service/preflightsvc"
	"github.com/common-fate/common-fate/pkg/service/rulesvc"
	"github.com/common-fate/common-fate/pkg/service/targetsvc"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// API holds all of our API endpoint handlers.
// We use a schema-first approach to ensure that the
// API meets our OpenAPI specification.
//
// To add a new endpoint, follow the below steps:
//
// 1. Edit `openapi.yaml` in this repository.
//
// 2. Run `make generate` to update the generated handler code.
// The code is generated into types.gen.go, and the function
// signatures can be found on the ServerInterface interface.
//
// 3. You'll get a compilation error because API no longer meets
// the ServerInterface interface. The missing function will be your
// new endpoint. Implement the function on API, ensuring that the function
// signature matches the ServerInterface interface.
type API struct {
	// DB is the DynamoDB client which provides direct storage access.
	DB               ddb.Storage
	DeploymentConfig deploy.DeployConfigReader
	// Requests is the service which provides business logic for Access Requests.
	Access           AccessService
	Rules            AccessRuleService
	AdminGroup       string
	IdentityProvider string
	FrontendURL      string

	IdentitySyncer auth.IdentitySyncer
	// Set this to nil if cognito is not configured as the IDP for the deployment
	Cognito            CognitoService
	InternalIdentity   InternalIdentityService
	TargetService      TargetService
	HandlerService     HandlerService
	HealthcheckService HealthcheckService
	PreflightService   PreflightService
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_cognito_service.go -package=mocks . CognitoService
type CognitoService interface {
	AdminCreateUser(ctx context.Context, in cognitosvc.CreateUserOpts) (*identity.User, error)
	AdminUpdateUserGroups(ctx context.Context, in cognitosvc.UpdateUserGroupsOpts) (*identity.User, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_access_service.go -package=mocks . AccessService

// RequestServices can create Access Requests.
type AccessService interface {
	CreateRequest(ctx context.Context, user identity.User, in types.CreateAccessRequestRequest) (*access.RequestWithGroupsWithTargets, error)
	RevokeRequest(ctx context.Context, in access.RequestWithGroupsWithTargets) (*access.RequestWithGroupsWithTargets, error)
	Review(ctx context.Context, user identity.User, isAdmin bool, requestID string, groupID string, in types.ReviewRequest) error
	CancelRequest(ctx context.Context, opts accesssvc.CancelRequestOpts) error
	// CreateFavorite(ctx context.Context, in accesssvc.CreateFavoriteOpts) (*access.Favorite, error)
	// UpdateFavorite(ctx context.Context, in accesssvc.UpdateFavoriteOpts) (*access.Favorite, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_accessrule_service.go -package=mocks . AccessRuleService

// AccessRuleService can create and get rules
type AccessRuleService interface {
	DeleteRule(ctx context.Context, id string) error
	CreateAccessRule(ctx context.Context, userID string, in types.CreateAccessRuleRequest) (*rule.AccessRule, error)
	UpdateRule(ctx context.Context, in *rulesvc.UpdateOpts) (*rule.AccessRule, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_internalidentity_service.go -package=mocks . InternalIdentityService

type InternalIdentityService interface {
	UpdateGroup(ctx context.Context, group identity.Group, in types.CreateGroupRequest) (*identity.Group, error)
	CreateGroup(ctx context.Context, in types.CreateGroupRequest) (*identity.Group, error)
	UpdateUserGroups(ctx context.Context, user identity.User, groups []string) (*identity.User, error)
	DeleteGroup(ctx context.Context, group identity.Group) error
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_target_service.go -package=mocks . TargetService
type TargetService interface {
	CreateGroup(ctx context.Context, targetGroup types.CreateTargetGroupRequest) (*target.Group, error)
	CreateRoute(ctx context.Context, group string, req types.CreateTargetGroupLink) (*target.Route, error)
	DeleteGroup(ctx context.Context, group *target.Group) error
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_handler_service.go -package=mocks . HandlerService
type HandlerService interface {
	RegisterHandler(ctx context.Context, req types.RegisterHandlerRequest) (*handler.Handler, error)
	DeleteHandler(ctx context.Context, handler *handler.Handler) error
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_preflight_service.go -package=mocks . PreflightService
type PreflightService interface {
	ProcessPreflight(ctx context.Context, user identity.User, preflightRequest types.CreatePreflightRequest) (*access.Preflight, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_healthcheck_service.go -package=mocks . HealthcheckService
type HealthcheckService interface {
	Check(ctx context.Context) error
}

// API must meet the generated REST API interface.
var _ types.ServerInterface = &API{}

type Opts struct {
	Log                    *zap.SugaredLogger
	ProviderRegistryClient registry_types.ClientWithResponsesInterface
	UseLocalEventHandler   bool
	IdentitySyncer         auth.IdentitySyncer
	DeploymentConfig       deploy.DeployConfigReader
	DynamoTable            string
	PaginationKMSKeyARN    string
	AdminGroup             string
	DeploymentSuffix       string
	CognitoUserPoolID      string
	IDPType                string
	AdminGroupID           string
	FrontendURL            string
	EventBusArn            string
}

// New creates a new API.
func New(ctx context.Context, opts Opts) (*API, error) {
	if opts.Log == nil {
		return nil, errors.New("opts.Log must be provided")
	}

	if opts.ProviderRegistryClient == nil {
		return nil, errors.New("ProviderRegistryClient must be provided")
	}

	tokenizer, err := ddb.NewKMSTokenizer(ctx, opts.PaginationKMSKeyARN)
	if err != nil {
		return nil, err
	}
	db, err := ddb.New(ctx, opts.DynamoTable, ddb.WithPageTokenizer(tokenizer))
	if err != nil {
		return nil, err
	}

	clk := clock.New()
	var eventBus gevent.EventPutter
	if !opts.UseLocalEventHandler {
		eventBus, err = gevent.NewSender(ctx, gevent.SenderOpts{
			EventBusARN: opts.EventBusArn,
		})
		if err != nil {
			return nil, err
		}
	} else {
		eventBus = eventhandler.NewLocalDevEventHandler(ctx, db, clk)
	}

	a := API{
		DeploymentConfig: opts.DeploymentConfig,
		AdminGroup:       opts.AdminGroup,
		FrontendURL:      opts.FrontendURL,
		InternalIdentity: &internalidentitysvc.Service{
			DB:    db,
			Clock: clk,
		},
		PreflightService: &preflightsvc.Service{
			DB:    db,
			Clock: clk,
		},
		Access: &accesssvc.Service{
			Clock:       clk,
			DB:          db,
			EventPutter: eventBus,
			Rules: &rulesvc.Service{
				Clock: clk,
				DB:    db,
			},
		},
		Rules: &rulesvc.Service{
			Clock: clk,
			DB:    db,
		},

		DB:               db,
		IdentitySyncer:   opts.IdentitySyncer,
		IdentityProvider: opts.IDPType,
		TargetService: &targetsvc.Service{
			DB:                     db,
			Clock:                  clk,
			ProviderRegistryClient: opts.ProviderRegistryClient,
		},
		HandlerService: &handlersvc.Service{
			DB:    db,
			Clock: clk,
		},
		HealthcheckService: &healthchecksvc.Service{
			DB:            db,
			RuntimeGetter: healthchecksvc.DefaultGetter{},
		},
	}

	// only initialise this if cognito is the IDP
	if opts.IDPType == identitysync.IDPTypeCognito {
		cog := &identitysync.CognitoSync{}
		err = cog.Config().Load(ctx, &gconfig.MapLoader{Values: map[string]string{"userPoolId": opts.CognitoUserPoolID}})
		if err != nil {
			return nil, err
		}
		err = cog.Init(ctx)
		if err != nil {
			return nil, err
		}
		a.Cognito = &cognitosvc.Service{
			Clock:        clk,
			DB:           db,
			Syncer:       opts.IdentitySyncer,
			Cognito:      cog,
			AdminGroupID: opts.AdminGroupID,
		}

	}

	return &a, nil
}

// Handler returns a HTTP handler.
// Hander doesn't add any middleware. It is the caller's
// responsibility to add any middleware.
func (a *API) Handler(r chi.Router) http.Handler {
	return types.HandlerWithOptions(a, types.ChiServerOptions{
		BaseRouter: r,
	})
}
