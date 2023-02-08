package auth

import (
	"context"
	"io"

	"errors"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/common-fate/apikit/logger"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/go-chi/chi/v5"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestAdminAuthorizer(t *testing.T) {
	type testcase struct {
		name              string
		adminGroup        string
		user              identity.User
		noClaimMiddleware bool
		path              string
		wantBody          string
		wantPanic         string
		wantCode          int
	}

	testcases := []testcase{
		{
			name:       "ok",
			adminGroup: "admins",
			user: identity.User{
				Groups: []string{"admins"},
			},
			path:     "/api/v1/admin/test",
			wantBody: "ok",
			wantCode: http.StatusOK,
		},
		{
			name:       "not allowed",
			adminGroup: "admins",
			user: identity.User{
				Groups: []string{"other"},
			},
			path:     "/api/v1/admin/test",
			wantBody: `{"error":"Unauthorized"}`,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:       "not allowed with empty claims",
			adminGroup: "admins",
			path:       "/api/v1/admin/test",
			wantBody:   `{"error":"Unauthorized"}`,
			wantCode:   http.StatusUnauthorized,
		},
		{
			name:       "invalid path",
			adminGroup: "admins",
			path:       "/api/v2/admin/test",
			wantBody:   `{"error":"Internal Server Error"}`,
			wantCode:   http.StatusInternalServerError,
		},
		{
			name:              "no claim middleware",
			adminGroup:        "admins",
			noClaimMiddleware: true,
			path:              "/api/v1/admin/test",
			wantBody:          `{"error":"Internal Server Error"}`,
			wantCode:          http.StatusInternalServerError,
		},
		{
			name:       "internal server error on empty adminGroup",
			adminGroup: "",
			wantBody:   `{"error":"The Common Fate administrator group is empty. Update the administrator group in your deployment configuration and redeploy."}`,
			wantCode:   http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			r := chi.NewRouter()

			log := zaptest.NewLogger(t)
			r.Use(logger.Middleware(log))
			// test middleware to set claims
			if !tc.noClaimMiddleware {
				r.Use(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						ctx := r.Context()
						ctx = context.WithValue(ctx, userContext, &tc.user)
						r = r.WithContext(ctx)
						next.ServeHTTP(w, r)
					})
				})
			}

			if tc.wantPanic != "" {
				defer func() {
					err := recover()

					if err != tc.wantPanic {
						t.Fatalf("Wrong panic message: %s", err)
					}
				}()

			}

			r.Use(AdminAuthorizer(tc.adminGroup))
			r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("ok"))
				w.WriteHeader(http.StatusOK)
			})

			req, err := http.NewRequest("GET", tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantCode, rr.Code)

			data, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.wantBody, string(data))
		})
	}
}

func TestMiddleware(t *testing.T) {
	type testcase struct {
		name       string
		claims     *Claims
		authErr    error
		getUserErr error
		idpSyncErr error
		wantBody   string
		wantCode   int
	}

	testcases := []testcase{
		{
			name: "ok",
			claims: &Claims{
				Sub:   "123",
				Email: "test@test.com",
			},
			wantBody: `ok`,
			wantCode: http.StatusOK,
		},
		{
			name:     "authenticator error",
			authErr:  errors.New("error"),
			wantBody: `{"error":"Unauthorized"}`,
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "IDP sync error",
			claims: &Claims{
				Sub:   "123",
				Email: "test@test.com",
			},
			getUserErr: ddb.ErrNoItems,
			idpSyncErr: errors.New("error syncing idp"),
			wantBody:   `{"error":"Unauthorized"}`,
			wantCode:   http.StatusUnauthorized,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			r := chi.NewRouter()
			c := ddbmock.New(t)
			c.MockQueryWithErr(&storage.GetUserByEmail{Email: "test@test.com", Result: &identity.User{}}, tc.getUserErr)

			log := zaptest.NewLogger(t)
			r.Use(logger.Middleware(log))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockAuthenticator(ctrl)
			m.EXPECT().Authenticate(gomock.Any()).Return(tc.claims, tc.authErr)

			mis := NewMockIdentitySyncer(ctrl)
			mis.EXPECT().Sync(gomock.Any()).Return(tc.idpSyncErr).AnyTimes()

			r.Use(Middleware(m, c, mis))
			r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte("ok"))
				w.WriteHeader(http.StatusOK)
			})

			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			assert.Equal(t, tc.wantCode, rr.Code)

			data, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.wantBody, string(data))
		})
	}
}
