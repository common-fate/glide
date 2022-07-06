package providers

import (
	"github.com/common-fate/granted-approvals/accesshandler/pkg/genv"
	"go.uber.org/zap"
)

// LogConfig returns a zap field for the provider's configuration.
// Under the hood, it calls `genv.Dump()` to dump all of the config
// to a map. Any secrets are redacted.
func LogConfig(c Configer) zap.Field {
	cfg := genv.Dump(c.Config())
	return zap.Any("provider.config", cfg)
}
