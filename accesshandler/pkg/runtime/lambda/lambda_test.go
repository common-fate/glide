package lambda

import (
	"context"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	ctx := context.Background()
	r := Runtime{}

	os.Setenv("STATE_MACHINE_ARN", "test:arn")

	err := r.Init(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
