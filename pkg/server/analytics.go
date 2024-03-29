package server

import (
	"net/http"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/depid"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

func analyticsMiddleware(db ddb.Storage, log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	dl := depid.New(db, log)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			client := analytics.New(analytics.Env())
			ctx = analytics.SetContext(ctx, client)
			r = r.WithContext(ctx)

			dep, err := dl.GetDeployment(ctx)
			if err != nil {
				logger.Get(ctx).Errorw("error getting deployment", zap.Error(err))
			}
			if err == nil && dep != nil {
				client.SetDeploymentID(dep.ID)
			}

			defer client.Close()
			next.ServeHTTP(w, r)
		})
	}
}
