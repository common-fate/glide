package gconfig

import "context"

type Configer interface {
	Config() Config
}

// Initers perform some initialisation behaviour such as setting up API clients.
type Initer interface {
	Init(ctx context.Context) error
}
