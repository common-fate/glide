package api

import (
	"context"
	"strings"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/runtime/lambda"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/runtime/local"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
)

// A runtime is responsible for the actual execution of a grant and are tied to the
// hosting environment the Access Handler is running in.
//
// Example runtimes are local (for testing only), and AWS Lambda with Step Functions.
type Runtime interface {
	// Init contains any runtime-specific initialisation logic.
	Init(ctx context.Context) error

	// CreateGrant creates a grant by executing runtime-specific workflow logic, such as
	// initiating an AWS Step Functions workflow.
	CreateGrant(ctx context.Context, grant types.ValidCreateGrant) (*types.Grant, error)

	// RevokeGrant revokes a grant by executing runtime-specific workflow logic, such as
	// initiating an AWS Step Functions workflow.
	// Revokes a grant and terminates the previous create grant workflow
	RevokeGrant(ctx context.Context, grantID string) (*types.Grant, error)
}

// runtimes is a map of the supported runtime environments
// for the API.
var runtimes = map[string]Runtime{
	"local":  &local.Runtime{},
	"lambda": &lambda.Runtime{},
}

// validRuntimes returns a comma-separated list of accepted runtime arguments.
func validRuntimes() string {
	var names []string
	for n := range runtimes {
		names = append(names, n)
	}

	return strings.Join(names, ", ")
}
