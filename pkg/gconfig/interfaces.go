package gconfig

import "context"

type Configer interface {
	Config() Config
}

// Initers perform some initialisation behaviour such as setting up API clients.
type Initer interface {
	Init(ctx context.Context) error
}

// Tester interface is used to run tests on config
type Tester interface {
	// TestConfig is expected to be called on a loaded config after Init has been called.
	TestConfig(ctx context.Context) error
}
