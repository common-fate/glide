package server

import (
	"context"
	"net/http"
	"time"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/errhandler"
	"github.com/common-fate/apikit/userid"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

// errorHandler meets the errhandler.Handler interface.
type errorHandler struct {
	hub *sentry.Hub
}

func (h *errorHandler) HandleError(err error) {
	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	webErr, ok := errors.Cause(err).(*apio.APIError)

	// we only want to send unhandled errors to Sentry,
	// so if we get a 3xx or 4xx status code on our
	// error we return early and don't dispatch it.
	if ok && webErr.Status < 500 {
		return
	}

	// if we get here, the error is unhandled and we want to send it to Sentry.
	h.hub.CaptureException(err)
}

// sentryMiddleware sends unhandled errors with a 5xx code to Sentry.
func sentryMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		hub := sentry.GetHubFromContext(ctx)
		if hub == nil {
			hub = sentry.CurrentHub().Clone()
			ctx = sentry.SetHubOnContext(ctx, hub)
		}
		hub.Scope().SetRequest(r)

		uid := userid.Get(ctx)
		if uid != "" {
			hub.Scope().SetUser(sentry.User{ID: uid})
		}
		reqID := middleware.GetReqID(ctx)
		if reqID != "" {
			hub.Scope().SetContext("request", map[string]interface{}{
				"id": reqID,
			})
		}

		defer recoverWithSentry(hub, r)

		// add the error handler to the request context,
		// so that calls to apio.Error() are sent to sentry.
		h := errorHandler{hub}
		ctx = errhandler.Set(ctx, &h)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// recoverWithSentry sends panic events to Sentry.
// It's modified from https://github.com/getsentry/sentry-go/blob/aed9115503d3bc31fd9be452423831c2bc7d19c7/http/sentryhttp.go
func recoverWithSentry(hub *sentry.Hub, r *http.Request) {
	if err := recover(); err != nil {
		eventID := hub.RecoverWithContext(
			context.WithValue(r.Context(), sentry.RequestContextKey, r),
			err,
		)
		if eventID != nil {
			// most of the time this is running in a serverless environment,
			// so always flush sentry to avoid us missing an event.
			hub.Flush(2 * time.Second)
		}

		panic(err)
	}
}
