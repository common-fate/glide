package providers

import (
	"testing"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/stretchr/testify/assert"
)

func TestGrantValidationsResults(t *testing.T) {
	a := GrantValidationResults{"a": GrantValidationResult{Name: "a", Logs: diagnostics.Logs{
		diagnostics.Log{
			Level: diagnostics.ErrorLevel,
			Msg:   "failed",
		},
	}}}
	b := GrantValidationResults{"b": GrantValidationResult{Name: "b", Logs: diagnostics.Logs{
		diagnostics.Log{
			Level: diagnostics.InfoLevel,
			Msg:   "did not fail",
		},
	}}}

	assert.True(t, a.Failed())
	assert.Equal(t, "1 error occurred:\n\t* a\n\n", a.FailureMessage())

	assert.False(t, b.Failed())
	assert.Equal(t, "", b.FailureMessage())
}
