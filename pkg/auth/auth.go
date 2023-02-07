package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/common-fate/apikit/apio"
	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/apikit/userid"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
	"go.uber.org/zap"
)

type contextKey struct {
	name string
}

var userIDContext = contextKey{name: "userIDContext"}
var userContext = contextKey{name: "userContext"}
var adminContext = contextKey{name: "adminContext"}

// Claims stores the relevant claims from a user's provided auth token.
// The identity token contains more claims, but we only parse the ones that we need.
type Claims struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

//go:generate go run github.com/golang/mock/mockgen -destination=mock_authenticator.go -package=auth . Authenticator

// Authenticators can extract Claims representing a user's authentication from an incoming request.
type Authenticator interface {
	Authenticate(r *http.Request) (*Claims, error)
}

//go:generate go run github.com/golang/mock/mockgen -destination=mock_identitysyncer.go -package=auth . IdentitySyncer

// IdentitySyncer syncs the users with the external identity provider, like Okta or Google Workspaces.
type IdentitySyncer interface {
	Sync(ctx context.Context) error
}

// Middleware is authentication middleware for the Common Fate API.
//
// It takes an Authenticator which knows how to extract the user's identity from the incoming request.
// If the user doesn't exist in the database the middleware will attempt to sync it from the
// connected identity provider.
func Middleware(authenticator Authenticator, db ddb.Storage, idp IdentitySyncer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.Get(ctx)
			claims, err := authenticator.Authenticate(r)
			if err != nil {
				log.Errorw("authentication error", zap.Error(err))
				apio.ErrorString(ctx, w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// get email from claims
			// lookup user by email
			// add user to context
			q := &storage.GetUserByEmail{
				Email: claims.Email,
			}
			_, err = db.Query(ctx, q)

			// log an error and return an unauthorized response if we can't handle the error from the DB.
			if err != nil && err != ddb.ErrNoItems {
				log.Errorw("authentication error", zap.Error(err))
				apio.ErrorString(ctx, w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// if we get ddb.ErrNoItems, the user may not have been synced yet from the IDP.
			// try and sync them now.
			// Note: this approach isn't very performant. In future this can be replaced with an incremental sync,
			// which we can use for handling IDP webhook notifications rather than polling.
			if err == ddb.ErrNoItems {
				log.Info("user does not exist in database - running an IDP sync and trying again", "user", claims)
				err = idp.Sync(ctx)
				if err != nil {
					log.Errorw("error syncing IDP", zap.Error(err))
					apio.ErrorString(ctx, w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
				log.Info("looking up user again")
				// reuse the same query, so that we can access the results later if it's successful.
				_, err = db.Query(ctx, q)
				if err != nil {
					log.Errorw("authentication error", zap.Error(err))
					apio.ErrorString(ctx, w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
			}

			ctx = context.WithValue(ctx, userContext, q.Result)
			ctx = context.WithValue(ctx, userIDContext, q.Result.ID)
			ctx = userid.Set(ctx, q.Result.ID)

			log.Debugw("user is authenticated", "claims", claims)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// AdminAuthorizer only allows users belonging to adminGroup to access administrative endpoints.
// The middleware currently gates all endpoints in the format /api/v1/admin/*
func AdminAuthorizer(adminGroup string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			// the admin group should always be non-empty.
			// If it is empty, we return an internal server error and don't allow any access,
			// as we can't determine whether the user should be authorized to access them.
			if adminGroup == "" {
				apio.ErrorString(ctx, w, "The Common Fate administrator group is empty. Update the administrator group in your deployment configuration and redeploy.", http.StatusInternalServerError)
				return
			}

			usr, ok := ctx.Value(userContext).(*identity.User)
			if !ok {
				apio.Error(ctx, w, errors.New("could not parse auth claims from context"))
				return
			}

			if !strings.HasPrefix(r.URL.Path, "/api/v1") {
				// we can't handle /api/v2 or other future versions of the API,
				// so rather than passing through we return an error here to ensure
				// this middleware is updated in future.
				apio.Error(ctx, w, fmt.Errorf("no authorization logic implemented for %s as it doesn't begin with /api/v1", r.URL.Path))
				return
			}

			isAdminRoute := strings.HasPrefix(r.URL.Path, "/api/v1/admin")

			// a user is an admin if they belong to the adminGroup.
			// The user's groups are set in their identity provider such as Okta or Google Workspace.
			isAdmin := contains(usr.Groups, adminGroup)
			ctx = context.WithValue(ctx, adminContext, isAdmin)
			r = r.WithContext(ctx)

			if isAdminRoute && !isAdmin {
				// the user is trying to access an admin route, but they're not authorized.
				// return a HTTP401 unauthorized response.
				apio.ErrorString(ctx, w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// UserIDFromContext returns the current user's ID.
// It requires that auth.Middleware has run.
func UserIDFromContext(ctx context.Context) string {
	return ctx.Value(userIDContext).(string)
}

// IsAdmin returns whether the user is an administrator or not.
// It requires that the AdminAuthorizer middleware has run.
func IsAdmin(ctx context.Context) bool {
	return ctx.Value(adminContext).(bool)
}

// UserIDFromContext returns the current user's ID.
// It requires that auth.Middleware has run.
func UserFromContext(ctx context.Context) *identity.User {
	usr := ctx.Value(userContext)
	return usr.(*identity.User)
}

// contains is a helper function to check if a string slice
// contains a particular string.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
