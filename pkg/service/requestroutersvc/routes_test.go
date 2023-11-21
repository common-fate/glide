package requestroutersvc

import (
	"context"
	"testing"

	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/common-fate/pkg/target"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/common-fate/provider-registry-sdk-go/pkg/providerregistrysdk"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestRoute(t *testing.T) {
	type testcase struct {
		name                  string
		give                  target.Group
		tgLookupwantErr       error
		wantErr               error
		wantRoutesForGroup    []target.Route
		wantHighPriorityRoute []target.Route
		wantHandler           *handler.Handler
		want                  *RouteResult
	}

	testcases := []testcase{
		{
			name:                  "ok",
			wantRoutesForGroup:    []target.Route{{Group: "123", Priority: 1}, {Group: "abc", Priority: 999}},
			wantHighPriorityRoute: []target.Route{{Group: "abc", Priority: 999}},
			wantHandler:           &handler.Handler{ID: "xyz"},
			wantErr:               nil,
			tgLookupwantErr:       nil,
			want:                  &RouteResult{Route: target.Route{Group: "abc", Handler: "", Kind: "", Priority: 999, Valid: false, Diagnostics: []target.Diagnostic(nil)}, Handler: handler.Handler{ID: "xyz", Runtime: "", AWSAccount: "", AWSRegion: "", Healthy: false, Diagnostics: []handler.Diagnostic(nil), ProviderDescription: (*providerregistrysdk.DescribeResponse)(nil)}},
		},
		{
			name:                  "No routes ",
			wantRoutesForGroup:    []target.Route{},
			wantHighPriorityRoute: []target.Route{},
			wantHandler:           &handler.Handler{},
			wantErr:               ErrNoRoutes,
			tgLookupwantErr:       nil,
			want:                  nil,
		},
		{
			name:                  "No valid routes ",
			wantRoutesForGroup:    []target.Route{{Group: "123", Priority: 1}, {Group: "abc", Priority: 999}},
			wantHighPriorityRoute: []target.Route{},
			wantHandler:           &handler.Handler{},
			wantErr:               ErrCannotRoute,
			tgLookupwantErr:       nil,
			want:                  nil,
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {

			db := ddbmock.New(t)

			db.MockQueryWithErr(&storage.ListTargetRoutesForGroup{Result: tc.wantRoutesForGroup}, tc.tgLookupwantErr)
			db.MockQueryWithErr(&storage.ListValidTargetRoutesForGroupByPriority{Result: tc.wantHighPriorityRoute}, tc.tgLookupwantErr)
			db.MockQueryWithErr(&storage.GetHandler{Result: tc.wantHandler}, tc.tgLookupwantErr)

			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			s := Service{
				DB: db,
			}

			got, err := s.Route(context.Background(), tc.give)

			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr.Error(), err.Error())
			}
			assert.Equal(t, tc.want, got)

		})
	}

}
