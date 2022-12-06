package lambda

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

// Runtime is a runtime which initiates a stepfunctions workflow
type Runtime struct {
	StateMachineARN        string `env:"COMMONFATE_STATE_MACHINE_ARN"`
	LogLevel               string `env:"LOG_LEVEL,default=info"`
	EventBusArn            string `env:"COMMONFATE_EVENT_BUS_ARN"`
	EventBusSource         string `env:"COMMONFATE_EVENT_BUS_SOURCE"`
	GranterStateMachineARN string `env:"COMMONFATE_STATE_MACHINE_ARN"`
}

// Init initialises the runtime
func (r *Runtime) Init(ctx context.Context) error {
	return envconfig.Process(ctx, r)
}
