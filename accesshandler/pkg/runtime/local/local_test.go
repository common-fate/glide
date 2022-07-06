package local

import (
	"context"
	"testing"
)

func TestInit(t *testing.T) {
	ctx := context.Background()
	r := Runtime{}

	err := r.Init(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
