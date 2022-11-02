package server

import (
	"net/http"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/ddb"
	"github.com/common-fate/granted-approvals/pkg/depid"
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
			} else {
				client.SetDeployment(dep.ToAnalytics())
			}

			defer client.Close()
			next.ServeHTTP(w, r)
		})
	}
}
