package fixtures

import (
	"context"
	"fmt"
	"strings"

	ssof "github.com/common-fate/granted-approvals/accesshandler/pkg/providers/aws/sso/fixtures"
	adf "github.com/common-fate/granted-approvals/accesshandler/pkg/providers/azure/ad/fixtures"
	oktaf "github.com/common-fate/granted-approvals/accesshandler/pkg/providers/okta/fixtures"
)

type GeneratorDestroyer interface {
	// Generate the fixture data. Returns a JSON encoding of the fixture data.
	Generate(ctx context.Context) ([]byte, error)
	// Destroy the fixture. 'data' is the JSON encoded fixture data.
	Destroy(ctx context.Context, data []byte) error
}

var FixtureRegistry = map[string]GeneratorDestroyer{
	"aws_sso": &ssof.Generator{},
	"okta":    &oktaf.Generator{},
	"azure":   &adf.Generator{},
}

func LookupGenerator(name string) (GeneratorDestroyer, error) {
	g, ok := FixtureRegistry[name]
	if !ok {
		return nil, fmt.Errorf("unknown generator %s. Allowed generators: %s", name, allowedGenerators())
	}
	return g, nil
}

// allowedGenerators returns a comma-separated list of the allowed fixture generators.
func allowedGenerators() string {
	var generators []string
	for n := range FixtureRegistry {
		generators = append(generators, n)
	}
	return strings.Join(generators, ", ")
}
