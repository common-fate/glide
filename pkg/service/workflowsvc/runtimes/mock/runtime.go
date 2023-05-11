package mock

import (
	"context"
	"encoding/json"

	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/gevent"
	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/service/requestroutersvc"
	"github.com/common-fate/common-fate/pkg/service/workflowsvc/runtimes/local"
	"github.com/common-fate/common-fate/pkg/targetgroupgranter"
	"github.com/common-fate/ddb"
	"github.com/common-fate/provider-registry-sdk-go/pkg/handlerclient"
	"github.com/common-fate/provider-registry-sdk-go/pkg/msg"
)

type EventPutter interface {
	Put(ctx context.Context, detail gevent.EventTyper) error
}

type MockRuntimeGetter struct{}

func (m *MockRuntimeGetter) GetRuntime(ctx context.Context, handler handler.Handler) (*handlerclient.Client, error) {
	return &handlerclient.Client{Executor: &MockRuntimeGetter{}}, nil
}

func (m *MockRuntimeGetter) Execute(ctx context.Context, request msg.Request) (*msg.Result, error) {
	if request.Type() == msg.RequestTypeGrant {
		b, err := json.Marshal(msg.GrantResponse{})
		if err != nil {
			return nil, err
		}
		return &msg.Result{Response: b}, nil
	}
	b, err := json.Marshal(struct{}{})
	if err != nil {
		return nil, err
	}
	return &msg.Result{Response: b}, nil
}

type Runtime struct {
	runtime *local.Runtime
}

func NewRuntime(db ddb.Storage, eventBus EventPutter, router *requestroutersvc.Service) *Runtime {
	return &Runtime{local.NewRuntime(db, &targetgroupgranter.Granter{
		DB: db, EventPutter: eventBus, RequestRouter: router,
		RuntimeGetter: &MockRuntimeGetter{},
	}, router)}
}

func (r *Runtime) Grant(ctx context.Context, grant access.GroupTarget) error {
	return r.runtime.Grant(ctx, grant)
}

func (r *Runtime) Revoke(ctx context.Context, grantID string) error {
	return r.runtime.Revoke(ctx, grantID)
}
