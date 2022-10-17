package providers

import (
	"context"
	"embed"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// Accessors know how to grant and revoke access to something.
// Accessors are considered the 'bare minimum' Granted providers.
// When writing a provider you must implement this interface.
type Accessor interface {
	// Grant the access.
	Grant(ctx context.Context, subject string, args []byte, grantID string) error

	// Revoke the access.
	Revoke(ctx context.Context, subject string, args []byte, grantID string) error
}

// AccessTokeners can indicate whether they need an access token to be generated
// as part of the access workflow.
//
// Access Tokens are used in Access Providers to tie a particular session in the
// downstream service back to the access request. In our ECS Shell provider,
// access tokens are enabled for audited Python shell access.
type AccessTokener interface {
	RequiresAccessToken() bool
}

// GrantValidator know how to validate access without actually granting it.
type GrantValidator interface {
	// ValidateGrant arguments and a subject for access without actually granting it.

	ValidateGrant() GrantValidationSteps
}

type ConfigValidationStep struct {
	Name            string
	FieldsValidated []string
	Run             func(ctx context.Context) diagnostics.Logs
}

// ConfigValues can validate the configuration of the Access Provider,
// such as checking whether API keys are valid and if roles can be assumed.
type ConfigValidator interface {
	ValidateConfig() map[string]ConfigValidationStep
}

type ArgSchema map[string]types.Argument

func (a ArgSchema) ToAPI() types.ArgSchema {
	argSchema := types.ArgSchema{
		AdditionalProperties: make(map[string]types.Argument),
	}
	for k, v := range a {
		argSchema.AdditionalProperties[k] = v
	}
	return argSchema
}

type ArgSchemarer interface {
	ArgSchemaV2() ArgSchema
}
type ArgOptionGroupValueser interface {
	ArgOptionGroupValues(ctx context.Context, argId string, groupingName string, groupingValues []string) ([]string, error)
}

// ArgOptioner provides a list of options for an argument and groupings if available.
type ArgOptioner interface {
	Options(ctx context.Context, arg string) (*types.ArgOptionsResponse, error)
}

// Instructioners provide instructions on how a user can access a role or
// resource that we've granted access to
type Instructioner interface {
	Instructions(ctx context.Context, subject string, args []byte, grantId string) (string, error)
}

// SetupDocers return an embedded filesystem containing setup documentation.
type SetupDocer interface {
	SetupDocs() embed.FS
}
