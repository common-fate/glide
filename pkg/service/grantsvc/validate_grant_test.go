package grantsvc

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/benbjohnson/clock"
	"github.com/common-fate/iso8601"

	"github.com/common-fate/ddb/ddbmock"
	ahTypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types/ahmocks"
	"github.com/common-fate/granted-approvals/pkg/access"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/storage"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestValidateGrant(t *testing.T) {
	type testcase struct {
		name                       string
		withCreateGrantResponseErr error
		withUser                   identity.User
		give                       CreateGrantOpts

		wantPostGrantsWithRequestBody  ahTypes.PostGrantsJSONRequestBody
		wantValidateGrantsResponseBody *ahTypes.ValidateGrantResponse

		wantErr error
	}
	clk := clock.NewMock()
	now := clk.Now()
	testcases := []testcase{
		{
			name: "validation passes",
			give: CreateGrantOpts{
				Request: access.Request{
					Status: access.APPROVED,
					RequestedTiming: access.Timing{
						Duration:  time.Minute,
						StartTime: &now,
					},
				},
			},

			withUser: identity.User{
				Email: "test@test.com",
			},
			wantPostGrantsWithRequestBody: ahTypes.PostGrantsJSONRequestBody{
				Start:   iso8601.New(now),
				End:     iso8601.New(now.Add(time.Minute)),
				Subject: "test@test.com",
				With: ahTypes.CreateGrant_With{
					AdditionalProperties: make(map[string]string),
				},
			},
			wantValidateGrantsResponseBody: &ahTypes.ValidateGrantResponse{HTTPResponse: &http.Response{StatusCode: http.StatusOK}},
		},
		{
			name: "validation has a known error (400)",
			give: CreateGrantOpts{
				Request: access.Request{
					Status: access.APPROVED,
					RequestedTiming: access.Timing{
						Duration:  time.Minute,
						StartTime: &now,
					},
				},
			},

			withUser: identity.User{
				Email: "test@test.com",
			},
			wantPostGrantsWithRequestBody: ahTypes.PostGrantsJSONRequestBody{
				Start:   iso8601.New(now),
				End:     iso8601.New(now.Add(time.Minute)),
				Subject: "test@test.com",
				With: ahTypes.CreateGrant_With{
					AdditionalProperties: make(map[string]string),
				},
			},
			wantValidateGrantsResponseBody: &ahTypes.ValidateGrantResponse{HTTPResponse: &http.Response{StatusCode: http.StatusBadRequest}, JSON400: &struct {
				Error *string "json:\"error,omitempty\""
			}{Error: aws.String("400 error")}},
			wantErr: &GrantValidationError{ValidationFailureMsg: "400 error"},
		},
		{
			name: "validation passes",
			give: CreateGrantOpts{
				Request: access.Request{
					Status: access.APPROVED,
					RequestedTiming: access.Timing{
						Duration:  time.Minute,
						StartTime: &now,
					},
				},
			},

			withUser: identity.User{
				Email: "test@test.com",
			},
			wantPostGrantsWithRequestBody: ahTypes.PostGrantsJSONRequestBody{
				Start:   iso8601.New(now),
				End:     iso8601.New(now.Add(time.Minute)),
				Subject: "test@test.com",
				With: ahTypes.CreateGrant_With{
					AdditionalProperties: make(map[string]string),
				},
			},
			wantValidateGrantsResponseBody: &ahTypes.ValidateGrantResponse{HTTPResponse: &http.Response{StatusCode: http.StatusInternalServerError}, JSON500: &struct {
				Error *string "json:\"error,omitempty\""
			}{Error: aws.String("500 error")}},
			wantErr: fmt.Errorf("access handler returned internal server error while validating grant. error: 500 error"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			g := ahmocks.NewMockClientWithResponsesInterface(ctrl)
			g.EXPECT().ValidateGrantWithResponse(gomock.Any(), gomock.Eq(tc.wantPostGrantsWithRequestBody)).Return(tc.wantValidateGrantsResponseBody, tc.withCreateGrantResponseErr).AnyTimes()

			c := ddbmock.New(t)
			c.MockQuery(&storage.GetUser{Result: &tc.withUser})

			s := Granter{
				AHClient: g,
				DB:       c,
				Clock:    clk,
			}

			err := s.ValidateGrant(context.Background(), tc.give)
			assert.Equal(t, tc.wantErr, err)

		})
	}
}
