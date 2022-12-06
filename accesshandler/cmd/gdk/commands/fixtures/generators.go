package fixtures

import (
	"context"
	"fmt"
	"strings"

	ssofv2 "github.com/common-fate/common-fate/accesshandler/pkg/providers/aws/sso-v2/fixtures"
	adf "github.com/common-fate/common-fate/accesshandler/pkg/providers/azure/ad/fixtures"
	oktaf "github.com/common-fate/common-fate/accesshandler/pkg/providers/okta/fixtures"
)

type GeneratorDestroyer interface {
	// Generate the fixture data. Returns a JSON encoding of the fixture data.
	Generate(ctx context.Context) ([]byte, error)
	// Destroy the fixture. 'data' is the JSON encoded fixture data.
	Destroy(ctx context.Context, data []byte) error
}

var FixtureRegistry = map[string]GeneratorDestroyer{
	"aws-sso-v2": &ssofv2.Generator{},
	"okta":       &oktaf.Generator{},
	"azure":      &adf.Generator{},
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
