package gconfig

import "context"

// RunConfigTest runs ConfigTest() if it is implemented on the interface
func RunConfigTest(ctx context.Context, testable interface{}) error {
	if tester, ok := testable.(Tester); ok {
		if initer, ok := testable.(Initer); ok {
			err := initer.Init(ctx)
			if err != nil {
				return err
			}
		}
		err := tester.TestConfig(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
