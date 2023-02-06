package targetdeploymentsvc

import (
	"context"

	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/ddb"
	"k8s.io/utils/clock"
)

// Service holds business logic relating to Cognito user management.
type Service struct {
	Clock clock.Clock
	DB    ddb.Storage

	TargetDeployments TargetDeployments
	AdminGroupID      string
}

type TargetDeployments interface {
	CreateDeployment(context.Context, CreateTargetDeploymentOpts) (targetgroup.Deployment, error)
	UpdateDeployment(context.Context, UpdateTargetDeploymentOpts) (targetgroup.Deployment, error)
	ArchiveDeployment(context.Context, ArchiveTargetDeploymentOpts) error
}
