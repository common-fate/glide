package accesssvc

import (
	"context"

	"github.com/benbjohnson/clock"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/rule"
	"github.com/common-fate/ddb"
)

// Service holds business logic relating to Access Requests.
type Service struct {
	Clock       clock.Clock
	DB          ddb.Storage
	EventPutter EventPutter
	Rules       AccessRuleService
}

type CreateGrantOpts struct {
	ID          string
	With        map[string]string
	AccessRule  rule.AccessRule
	RequestedBy identity.User
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/workflow.go -package=mocks . Workflow
type Workflow interface {
	Revoke(ctx context.Context, group access.GroupWithTargets, revokerID string, revokerEmail string) (*access.Group, error)
	Grant(ctx context.Context, group access.GroupWithTargets, subject string) ([]access.GroupTarget, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/eventputter.go -package=mocks . EventPutter
type EventPutter interface {
	Put(ctx context.Context, detail gevent.EventTyper) error
}

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_accessrule_service.go -package=mocks . AccessRuleService

// AccessRuleService can create and get rules
type AccessRuleService interface {
	// RequestArguments(ctx context.Context, accessRuleTarget rule.Target) (map[string]types.RequestArgument, error)
}
