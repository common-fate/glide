package healthcheck

import (
	"context"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/service/healthchecksvc"
	"github.com/common-fate/ddb"
)

type HealthChecker struct {
	DB          ddb.Storage
	HealthCheck *healthchecksvc.Service
}

func (s *HealthChecker) Check(ctx context.Context) error {
	log := logger.Get(ctx)
	log.Info("starting to check health")

	log.Info("completed checking health")

	return nil
}
