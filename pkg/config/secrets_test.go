package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanSuffix(t *testing.T) {
	type testcase struct {
		name       string
		givePath   SecretPath
		giveSuffix string
		want       string
	}
	testcases := []testcase{
		{name: "with suffix", givePath: GoogleTokenPath, giveSuffix: "!@#$%^&*()_=-helloHELLO  space:1", want: "/granted/secrets/identity/google/token-----------_--helloHELLO--space-1"},
		{name: "no suffix", givePath: GoogleTokenPath, want: "/granted/secrets/identity/google/token"},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := suffixedPath(tc.givePath, tc.giveSuffix)
			assert.Equal(t, tc.want, got)
		})
	}
}
