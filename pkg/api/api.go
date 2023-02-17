// The api package defines all of our REST API endpoints.
package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/benbjohnson/clock"
	registry_types "github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"

	"github.com/common-fate/common-fate/accesshandler/pkg/providerregistry"
	"github.com/common-fate/common-fate/accesshandler/pkg/psetup"
	ahtypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/auth"
	"github.com/common-fate/common-fate/pkg/cache"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/common-fate/common-fate/pkg/gconfig"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/identity/identitysync"
	"github.com/common-fate/common-fate/pkg/providersetup"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/common-fate/pkg/service/accesssvc"
	"github.com/common-fate/common-fate/pkg/service/cachesvc"
	"github.com/common-fate/common-fate/pkg/service/cognitosvc"

	"github.com/common-fate/common-fate/pkg/service/internalidentitysvc"
	"github.com/common-fate/common-fate/pkg/service/psetupsvc"
	"github.com/common-fate/common-fate/pkg/service/rulesvc"
	"github.com/common-fate/common-fate/pkg/service/targetdeploymentsvc"
	"github.com/common-fate/common-fate/pkg/service/targetgroupsvc"
	"github.com/common-fate/common-fate/pkg/service/workflowsvc"
	"github.com/common-fate/common-fate/pkg/service/workflowsvc/runtimes/live"
	"github.com/common-fate/common-fate/pkg/targetgroup"

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
	Access              AccessService
	Rules               AccessRuleService
	ProviderSetup       ProviderSetupService
	AccessHandlerClient ahtypes.ClientWithResponsesInterface
	AdminGroup          string
	IdentityProvider    string
	FrontendURL         string

	Cache          CacheService
	IdentitySyncer auth.IdentitySyncer
	// Set this to nil if cognito is not configured as the IDP for the deployment
	Cognito                      CognitoService
	InternalIdentity             InternalIdentityService
	TargetGroupService           TargetGroupService
	TargetGroupDeploymentService TargetGroupDeploymentService
	Workflow                     Workflow
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_cognito_service.go -package=mocks . CognitoService
type CognitoService interface {
	AdminCreateUser(ctx context.Context, in cognitosvc.CreateUserOpts) (*identity.User, error)
	AdminUpdateUserGroups(ctx context.Context, in cognitosvc.UpdateUserGroupsOpts) (*identity.User, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_providersetup_service.go -package=mocks . ProviderSetupService

// ProviderSetupService contains business logic for managing the guided provider setup workflows.
type ProviderSetupService interface {
	Create(ctx context.Context, providerType string, existingProviders deploy.ProviderMap, r providerregistry.ProviderRegistry) (*providersetup.Setup, error)
	CompleteStep(ctx context.Context, setupID string, stepIndex int, body types.ProviderSetupStepCompleteRequest) (*providersetup.Setup, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_access_service.go -package=mocks . AccessService

// RequestServices can create Access Requests.
type AccessService interface {
	CreateRequests(ctx context.Context, in accesssvc.CreateRequestsOpts) ([]accesssvc.CreateRequestResult, error)
	AddReviewAndGrantAccess(ctx context.Context, opts accesssvc.AddReviewOpts) (*accesssvc.AddReviewResult, error)
	CancelRequest(ctx context.Context, opts accesssvc.CancelRequestOpts) error
	CreateFavorite(ctx context.Context, in accesssvc.CreateFavoriteOpts) (*access.Favorite, error)
	UpdateFavorite(ctx context.Context, in accesssvc.UpdateFavoriteOpts) (*access.Favorite, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_accessrule_service.go -package=mocks . AccessRuleService

// AccessRuleService can create and get rules
type AccessRuleService interface {
	ArchiveAccessRule(ctx context.Context, userID string, in rule.AccessRule) (*rule.AccessRule, error)
	CreateAccessRule(ctx context.Context, userID string, in types.CreateAccessRuleRequest) (*rule.AccessRule, error)
	LookupRule(ctx context.Context, opts rulesvc.LookupRuleOpts) ([]rulesvc.LookedUpRule, error)
	GetRule(ctx context.Context, ID string, user *identity.User, isAdmin bool) (*rule.GetAccessRuleResponse, error)
	UpdateRule(ctx context.Context, in *rulesvc.UpdateOpts) (*rule.AccessRule, error)
	RequestArguments(ctx context.Context, accessRuleTarget rule.Target) (map[string]types.RequestArgument, error)
}

type CacheService interface {
	RefreshCachedProviderArgOptions(ctx context.Context, providerId string, argId string) (bool, []cache.ProviderOption, []cache.ProviderArgGroupOption, error)
	LoadCachedProviderArgOptions(ctx context.Context, providerId string, argId string) (bool, []cache.ProviderOption, []cache.ProviderArgGroupOption, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_internalidentity_service.go -package=mocks . InternalIdentityService

type InternalIdentityService interface {
	UpdateGroup(ctx context.Context, group identity.Group, in types.CreateGroupRequest) (*identity.Group, error)
	CreateGroup(ctx context.Context, in types.CreateGroupRequest) (*identity.Group, error)
	UpdateUserGroups(ctx context.Context, user identity.User, groups []string) (*identity.User, error)
	DeleteGroup(ctx context.Context, group identity.Group) error
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_target_group_service.go -package=mocks . TargetGroupService
type TargetGroupService interface {
	CreateTargetGroup(ctx context.Context, targetGroup types.CreateTargetGroupRequest) (*targetgroup.TargetGroup, error)
	CreateTargetGroupLink(ctx context.Context, req types.CreateTargetGroupLink, targetGroupId string) (*targetgroup.TargetGroup, error)
	RemoveTargetGroupLink(ctx context.Context, deploymentID string, targetGroupId string) error
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_target_group_deployment_service.go -package=mocks . TargetGroupDeploymentService
type TargetGroupDeploymentService interface {
	CreateTargetGroupDeployment(ctx context.Context, req types.CreateTargetGroupDeploymentRequest) (*targetgroup.Deployment, error)
}
type Workflow interface {
	Revoke(ctx context.Context, request access.Request, revokerID string) (*access.Request, error)
}

// API must meet the generated REST API interface.
var _ types.ServerInterface = &API{}

type Opts struct {
	Log                    *zap.SugaredLogger
	AccessHandlerClient    ahtypes.ClientWithResponsesInterface
	ProviderRegistryClient registry_types.ClientWithResponsesInterface
	EventSender            *gevent.Sender
	IdentitySyncer         auth.IdentitySyncer
	DeploymentConfig       deploy.DeployConfigReader
	DynamoTable            string
	PaginationKMSKeyARN    string
	AdminGroup             string
	TemplateData           psetup.TemplateData
	DeploymentSuffix       string
	CognitoUserPoolID      string
	IDPType                string
	AdminGroupID           string
	StateMachineARN        string
	FrontendURL            string
}

// New creates a new API.
func New(ctx context.Context, opts Opts) (*API, error) {
	if opts.Log == nil {
		return nil, errors.New("opts.Log must be provided")
	}
	if opts.AccessHandlerClient == nil {
		return nil, errors.New("AccessHandlerClient must be provided")
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

	if err != nil {
		return nil, err
	}

	a := API{
		DeploymentConfig: opts.DeploymentConfig,
		AdminGroup:       opts.AdminGroup,
		FrontendURL:      opts.FrontendURL,
		InternalIdentity: &internalidentitysvc.Service{
			DB:    db,
			Clock: clk,
		},
		Access: &accesssvc.Service{
			Clock:       clk,
			DB:          db,
			EventPutter: opts.EventSender,
			Cache: &cachesvc.Service{
				ProviderConfigReader: opts.DeploymentConfig,
				DB:                   db,
				AccessHandlerClient:  opts.AccessHandlerClient,
			},
			Rules: &rulesvc.Service{
				Clock:    clk,
				DB:       db,
				AHClient: opts.AccessHandlerClient,
				Cache: &cachesvc.Service{
					ProviderConfigReader: opts.DeploymentConfig,
					DB:                   db,
					AccessHandlerClient:  opts.AccessHandlerClient,
				},
			},
			AHClient: opts.AccessHandlerClient,
			Workflow: &workflowsvc.Service{
				Runtime: &live.Runtime{
					StateMachineARN: opts.StateMachineARN,
					AHClient:        opts.AccessHandlerClient,
					Eventbus:        opts.EventSender,
				},
				DB:       db,
				Clk:      clk,
				Eventbus: opts.EventSender,
			},
		},
		Cache: &cachesvc.Service{
			ProviderConfigReader: opts.DeploymentConfig,
			DB:                   db,
			AccessHandlerClient:  opts.AccessHandlerClient,
		},
		Rules: &rulesvc.Service{
			Clock:    clk,
			DB:       db,
			AHClient: opts.AccessHandlerClient,
			Cache: &cachesvc.Service{
				ProviderConfigReader: opts.DeploymentConfig,
				DB:                   db,
				AccessHandlerClient:  opts.AccessHandlerClient,
			},
		},
		ProviderSetup: &psetupsvc.Service{
			DB:               db,
			TemplateData:     opts.TemplateData,
			DeploymentSuffix: opts.DeploymentSuffix,
		},
		AccessHandlerClient: opts.AccessHandlerClient,
		DB:                  db,
		IdentitySyncer:      opts.IdentitySyncer,
		IdentityProvider:    opts.IDPType,
		TargetGroupService: &targetgroupsvc.Service{
			DB:                     db,
			Clock:                  clk,
			ProviderRegistryClient: opts.ProviderRegistryClient,
		},
		TargetGroupDeploymentService: &targetdeploymentsvc.Service{
			DB:    db,
			Clock: clk,
		},
		Workflow: &workflowsvc.Service{
			Runtime: &live.Runtime{
				StateMachineARN: opts.StateMachineARN,
				AHClient:        opts.AccessHandlerClient,
				Eventbus:        opts.EventSender,
			},
			DB:       db,
			Clk:      clk,
			Eventbus: opts.EventSender,
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
