// The api package defines all of our REST API endpoints.
package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/benbjohnson/clock"

	ahtypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gevent"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/service/accesssvc"
	"github.com/common-fate/granted-approvals/pkg/service/grantsvc"
	"github.com/common-fate/granted-approvals/pkg/service/rulesvc"

	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/types"

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
	DB ddb.Storage
	// Requests is the service which provides business logic for Access Requests.
	Access              AccessService
	Rules               AccessRuleService
	AccessHandlerClient ahtypes.ClientWithResponsesInterface
	AdminGroup          string
	Granter             accesssvc.Granter
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_access_service.go -package=mocks . AccessService

// RequestServices can create Access Requests.
type AccessService interface {
	CreateRequest(ctx context.Context, user *identity.User, in types.CreateRequestRequest) (*accesssvc.CreateRequestResult, error)
	AddReviewAndGrantAccess(ctx context.Context, opts accesssvc.AddReviewOpts) (*accesssvc.AddReviewResult, error)
	CancelRequest(ctx context.Context, opts accesssvc.CancelRequestOpts) error
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_accessrule_service.go -package=mocks . AccessRuleService

// AccessRuleService can create and get rules
type AccessRuleService interface {
	ArchiveAccessRule(ctx context.Context, user *identity.User, in rule.AccessRule) (*rule.AccessRule, error)
	CreateAccessRule(ctx context.Context, user *identity.User, in types.CreateAccessRuleRequest) (*rule.AccessRule, error)
	GetRule(ctx context.Context, ID string, user *identity.User, isAdmin bool) (*rule.AccessRule, error)
	UpdateRule(ctx context.Context, in *rulesvc.UpdateOpts) (*rule.AccessRule, error)
}

// API must meet the generated REST API interface.
var _ types.ServerInterface = &API{}

type Opts struct {
	Log                 *zap.SugaredLogger
	AccessHandlerClient ahtypes.ClientWithResponsesInterface
	EventSender         *gevent.Sender
	DynamoTable         string
	AdminGroup          string
}

// New creates a new API.
func New(ctx context.Context, opts Opts) (*API, error) {
	if opts.Log == nil {
		return nil, errors.New("opts.Log must be provided")
	}
	if opts.AccessHandlerClient == nil {
		return nil, errors.New("AccessHandlerClient must be provided")
	}

	db, err := ddb.New(ctx, opts.DynamoTable)
	if err != nil {
		return nil, err
	}

	clk := clock.New()

	a := API{
		AdminGroup: opts.AdminGroup,
		Access: &accesssvc.Service{
			Clock: clk,
			DB:    db,
			Granter: &grantsvc.Granter{
				AHClient: opts.AccessHandlerClient,
				DB:       db,
				Clock:    clk,
				EventBus: opts.EventSender,
			},
			EventPutter: opts.EventSender,
		},
		Rules: &rulesvc.Service{
			Clock:    clk,
			DB:       db,
			AHClient: opts.AccessHandlerClient,
		},

		AccessHandlerClient: opts.AccessHandlerClient,
		DB:                  db,

		Granter: &grantsvc.Granter{
			AHClient: opts.AccessHandlerClient,
			DB:       db,
			Clock:    clk,
			EventBus: opts.EventSender,
		},
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
