package providertest

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/accesshandler/pkg/config"
	"github.com/common-fate/common-fate/accesshandler/pkg/providers"
	"github.com/common-fate/common-fate/pkg/deploy"
	"github.com/joho/godotenv"
)

func BenchmarkOrganizationGraph(b *testing.B) {
	ctx := context.Background()
	_ = godotenv.Load("../../../.env")
	dc, err := deploy.GetDeploymentConfig()
	if err != nil {
		b.Fatal(err)
	}
	ps, err := dc.ReadProviders(ctx)
	if err != nil {
		b.Fatal(err)
	}
	err = config.ConfigureProviders(ctx, ps)
	if err != nil {
		b.Fatal(err)
	}
	sso := config.Providers["aws-sso-v2"]
	op := sso.Provider.(providers.ArgOptioner)
	for i := 0; i < b.N; i++ {
		options, err := op.Options(ctx, "accountId")
		if err != nil {
			b.Fatal(err)
		}
		_ = options
	}

}
