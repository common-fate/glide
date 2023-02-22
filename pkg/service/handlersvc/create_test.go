package handlersvc

import (
	"context"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/common-fate/pkg/handler"
	"github.com/common-fate/common-fate/pkg/storage"

	"github.com/common-fate/common-fate/pkg/types"
	"github.com/common-fate/ddb"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateTargetGroupDeployment(t *testing.T) {
	type testcase struct {
		name string
		// database lookup return object (used to mock ErrTargetGroupDeploymentIdAlreadyExists)
		mockGet *storage.GetHandler
		// database put object (used to mock ok response)
		mockPut *handler.Handler
		// input to CreateTargetGroupDeployment
		give    types.CreateTargetGroupDeploymentRequest
		wantErr error
		want    *handler.Handler
	}

	testcases := []testcase{
		{
			name: "bad aws account number",
			give: types.CreateTargetGroupDeploymentRequest{
				AwsAccount: "123_bad_123",
			},
			wantErr: ErrInvalidAwsAccountNumber,
		},
		{
			name: "existing deployment found",
			mockGet: &storage.GetHandler{
				ID: "test1",
				Result: &handler.Handler{
					ID: "test1",
				},
			},
			give: types.CreateTargetGroupDeploymentRequest{
				Id:         "test1",
				AwsAccount: "123456789012",
			},
			wantErr: ErrHandlerIdAlreadyExists,
		},
		{
			name: "ok",
			mockPut: &handler.Handler{
				ID:         "test1",
				AWSAccount: "123456789011",
			},
			give: types.CreateTargetGroupDeploymentRequest{
				Id:         "test1",
				AwsAccount: "123456789012",
			},
			want: &handler.Handler{
				ID:         "test1",
				AWSAccount: "123456789012",
				Diagnostics: []handler.Diagnostic{
					{
						Level:   string(types.ProviderSetupDiagnosticLogLevelINFO),
						Message: "offline: lambda cannot be reached/invoked",
					},
				},
			},
		},
	}

	for _, tc := range testcases {

		tc := tc

		t.Run(tc.name, func(t *testing.T) {

			dbMock := ddbmock.New(t)

			if tc.mockGet != nil {
				// this is used to mock the db lookup for coverage of ErrTargetGroupDeploymentIdAlreadyExists
				dbMock.MockQuery(tc.mockGet)
			} else {
				// this is used to mock s.DB.Put
				dbMock.MockQueryWithErr(&storage.GetHandler{}, ddb.ErrNoItems)
			}
			if tc.mockPut != nil {
				ctx := context.Background()
				err := dbMock.Put(ctx, tc.mockPut)
				if err != nil {
					t.Fatal(err)
				}
			}

			clk := clock.NewMock()

			s := Service{
				Clock: clk,
				DB:    dbMock,
			}

			got, err := s.CreateHandler(context.Background(), tc.give)

			if err != nil && tc.wantErr != nil {
				assert.Equal(t, tc.wantErr.Error(), err.Error())
				return
			}
			assert.Equal(t, tc.want, got)

		})
	}

}
