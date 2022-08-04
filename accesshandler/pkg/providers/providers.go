package providers

import (
	"context"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/invopop/jsonschema"
)

// Accessors know how to grant and revoke access to something.
// Accessors are considered the 'bare minimum' Granted providers.
// When writing a provider you must implement this interface.
type Accessor interface {
	// Grant the access.
	Grant(ctx context.Context, subject string, args []byte) error

	// Revoke the access.
	Revoke(ctx context.Context, subject string, args []byte) error
}

// Validators know how to validate access without actually granting it.
type Validator interface {
	// Validate arguments and a subject for access without actually granting it.
	Validate(ctx context.Context, subject string, args []byte) error
}

// ArgSchemarers provide a JSON Schema for the arguments they accept.
type ArgSchemarer interface {
	ArgSchema() *jsonschema.Schema
}

// ArgOptioner provides a list of options for an argument.
type ArgOptioner interface {
	Options(ctx context.Context, arg string) ([]types.Option, error)
}

// Instructioners provide instructions on how a user can access a role or
// resource that we've granted access to
type Instructioner interface {
	Instructions(ctx context.Context, subject string, args []byte) (string, error)
}
