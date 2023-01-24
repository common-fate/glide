// The api package defines all of our REST API endpoints.
package api

import (
	"context"
	"net/http"

	"github.com/common-fate/common-fate/pkg/deploymentcli/types"
	"github.com/go-chi/chi/v5"
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
	// db     *gorm.DB
	// PSetup *psetupsvc.Service
}

// API must meet the generated REST API interface.
var _ types.ServerInterface = &API{}

type Opts struct {
	ProviderRegistryAPIURL string
}

// New creates a new API. You can add any additional constructor logic here.
func New(ctx context.Context, o Opts) (*API, error) {
	// db, err := storage.NewLocal()
	// if err != nil {
	// 	return nil, err
	// }
	// ahc, err := registryTypes.NewClientWithResponses(o.ProviderRegistryAPIURL)
	// if err != nil {
	// 	return nil, err
	// }
	a := API{
		// db: db, PSetup: &psetupsvc.Service{
		// DB:               db,
		// DeploymentSuffix: "",
		// TemplateData: psetup.TemplateData{
		// 	AccessHandlerExecutionRoleARN: "",
		// },
		// Registry: ahc,
		// }
	}
	return &a, nil
}

// Handler returns a HTTP handler.
// Hander doesn't add any middleware. It is the caller's
// responsibility to add any middleware.
func (a *API) Handler(r chi.Router) http.Handler {
	return types.HandlerFromMux(a, r)
}
