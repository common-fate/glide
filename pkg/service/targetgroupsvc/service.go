package targetgroupsvc

import (
	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Cognito user management.
type Service struct {
	Clock clock.Clock
	DB    ddb.Storage

	// TargetDeployments TargetGroup
}

type TargetGroup interface {
	// CreateTargetGroup(context.Context, CreateTargetGroupOpts) (targetgroup.TargetGroup, error)
	// UpdateTargetGroup(context.Context, UpdateTargetGroupOpts) (targetgroup.TargetGroup, error)
	// ArchiveTargetGroup(context.Context, ArchiveTargetGroupOpts) error
}
