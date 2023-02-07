package targetgroupsvc

import (
	"context"

	registry_types "github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/targetgroup"
	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
)

type Service struct {
	Clock                  clock.Clock
	DB                     ddb.Storage
	ProviderRegistryClient registry_types.ClientWithResponsesInterface
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/targetgroup.go -package=mocks . TargetGroupService
type TargetGroupService interface {
	CreateTargetGroup(ctx context.Context, req types.CreateTargetGroupRequest) (*targetgroup.TargetGroup, error)
	UpdateTargetGroup(ctx context.Context, req UpdateOpts) (*targetgroup.TargetGroup, error)
	// ArchiveTargetGroup(context.Context, ArchiveTargetGroupOpts) error
}
