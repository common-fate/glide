package targetgroupsvc

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

	TargetDeployments TargetGroup
	AdminGroupID      string
}

type TargetGroup interface {
	CreateTargetGroup(context.Context, CreateTargetGroupOpts) (targetgroup.TargetGroup, error)
	UpdateTargetGroup(context.Context, UpdateTargetGroupOpts) (targetgroup.TargetGroup, error)
	ArchiveTargetGroup(context.Context, ArchiveTargetGroupOpts) error
}
