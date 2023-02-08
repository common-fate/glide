package targetdeploymentsvc

import (
	"context"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	registry_types "github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
)

// Service holds business logic relating to Cognito user management.
type Service struct {
	Clock clock.Clock
	DB    ddb.Storage
	// TargetDeployments TargetDeployments
	// AdminGroupID      string
	ProviderRegistryClient registry_types.ClientWithResponsesInterface
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/targetgroupdeployment.go -package=mocks . TargetGroupDeploymentService
type TargetGroupDeploymentService interface {
	CreateTargetGroupDeployment(ctx context.Context, req types.CreateTargetGroupDeploymentRequest) (*targetgroup.Deployment, error)
	// CreateDeployment(context.Context, CreateTargetDeploymentOpts) (targetgroup.Deployment, error)
	// UpdateDeployment(context.Context, UpdateTargetDeploymentOpts) (targetgroup.Deployment, error)
	// ArchiveDeployment(context.Context, ArchiveTargetDeploymentOpts) error
}
