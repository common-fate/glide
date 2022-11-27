package local

import (
	"context"
	"testing"
	"time"

	"github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/iso8601"
)

func TestCreateGrant(t *testing.T) {
	ctx := context.Background()
	r := Runtime{}

	err := r.Init(ctx)
	if err != nil {
		t.Fatal(err)
	}

	g := types.CreateGrant{
		Id:       "abcd",
		Provider: "test",
		Subject:  "test@acme.com",
		Start:    iso8601.New(time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC)),
		End:      iso8601.New(time.Date(2022, 1, 1, 10, 10, 0, 0, time.UTC)),
	}

	// testing time is 1st Jan 2022, 10:00am UTC
	now := time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC)

	vcg, err := g.Validate(ctx, now)
	if err != nil {
		t.Fatal(err)
	}

	_, err = r.CreateGrant(ctx, *vcg)
	if err != nil {
		t.Fatal(err)
	}
}
