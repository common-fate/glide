package grantsvc

import (
	"context"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/iso8601"

	ahTypes "github.com/common-fate/common-fate/accesshandler/pkg/types"
	"github.com/common-fate/common-fate/accesshandler/pkg/types/ahmocks"
	"github.com/common-fate/common-fate/pkg/access"
	"github.com/common-fate/common-fate/pkg/identity"
	"github.com/common-fate/common-fate/pkg/storage"
	"github.com/common-fate/ddb/ddbmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// testAccessTokenChecker is a mock implementation of NeedsAccessToken
type testAccessTokenChecker struct {
	NeedsToken bool
	Err        error
}

func (t testAccessTokenChecker) NeedsAccessToken(ctx context.Context, providerID string) (bool, error) {
	return t.NeedsToken, t.Err
}

func TestCreateGrant(t *testing.T) {
	type testcase struct {
		name                           string
		withCreateGrantResponse        *ahTypes.PostGrantsResponse
		withCreateGrantResponseErr     error
		withUser                       identity.User
		give                           CreateGrantOpts
		subject                        string
		wantPostGrantsWithResponseBody ahTypes.PostGrantsJSONRequestBody
		wantValidateRequestToProvider  ahTypes.ValidateGrantJSONRequestBody
		wantValidateRequestResponse    *ahTypes.ValidateGrantResponse
		needsAccessToken               bool
		needsAccessTokenErr            error

		wantRequest *access.Request
		wantErr     error
	}
	clk := clock.NewMock()
	now := clk.Now()
	overrideStart := now.Add(time.Hour)
	grantId := "abcd"
	testcases := []testcase{
		{
			name: "created success",
			give: CreateGrantOpts{
				Request: access.Request{
					Status: access.APPROVED,
					RequestedTiming: access.Timing{
						Duration:  time.Minute,
						StartTime: &now,
					},
				},
			},
			withCreateGrantResponse: &ahTypes.PostGrantsResponse{
				JSON201: &struct {
					Grant ahTypes.Grant "json:\"grant\""
				}{
					Grant: ahTypes.Grant{
						ID:      grantId,
						Start:   iso8601.New(now),
						End:     iso8601.New(now.Add(time.Minute)),
						Subject: "test@test.com",
					},
				},
			},
			withUser: identity.User{
				Email: "test@test.com",
			},
			wantPostGrantsWithResponseBody: ahTypes.PostGrantsJSONRequestBody{
				Start:   iso8601.New(now),
				End:     iso8601.New(now.Add(time.Minute)),
				Subject: "test@test.com",
				With: ahTypes.CreateGrant_With{
					AdditionalProperties: make(map[string]string),
				},
			},
			wantValidateRequestToProvider: ahTypes.ValidateGrantJSONRequestBody{
				Id:       "123",
				Provider: "OKTA",
				With: ahTypes.CreateGrant_With{
					AdditionalProperties: make(map[string]string),
				},
				Subject: "test@test.com",
				Start:   iso8601.New(overrideStart),
				End:     iso8601.New(overrideStart.Add(time.Minute * 2)),
			},
			wantValidateRequestResponse: &ahTypes.ValidateGrantResponse{},
			wantRequest: &access.Request{
				Status: access.APPROVED,
				RequestedTiming: access.Timing{
					Duration:  time.Minute,
					StartTime: &now,
				},
				Grant: &access.Grant{
					CreatedAt: clk.Now(),
					UpdatedAt: clk.Now(),
					Start:     iso8601.New(now).Time,
					End:       iso8601.New(now.Add(time.Minute)).Time,
					Subject:   "test@test.com",
				},
			},
			subject: "test@test.com",
		},
		{
			name: "created success with override timing",
			give: CreateGrantOpts{
				Request: access.Request{
					Status: access.APPROVED,
					RequestedTiming: access.Timing{
						Duration:  time.Minute,
						StartTime: &now,
					},
					OverrideTiming: &access.Timing{
						Duration:  time.Minute * 2,
						StartTime: &overrideStart,
					},
				},
			},
			subject: "test@test.com",

			withCreateGrantResponse: &ahTypes.PostGrantsResponse{
				JSON201: &struct {
					Grant ahTypes.Grant "json:\"grant\""
				}{
					Grant: ahTypes.Grant{
						ID:      grantId,
						Start:   iso8601.New(overrideStart),
						End:     iso8601.New(overrideStart.Add(time.Minute * 2)),
						Subject: "test@test.com",
					},
				},
			},
			withUser: identity.User{
				Email: "test@test.com",
			},
			wantPostGrantsWithResponseBody: ahTypes.PostGrantsJSONRequestBody{
				Start:   iso8601.New(overrideStart),
				End:     iso8601.New(overrideStart.Add(time.Minute * 2)),
				Subject: "test@test.com",
				With: ahTypes.CreateGrant_With{
					AdditionalProperties: make(map[string]string),
				},
			},
			wantValidateRequestToProvider: ahTypes.ValidateGrantJSONRequestBody{
				Id:       "123",
				Provider: "OKTA",
				With: ahTypes.CreateGrant_With{
					AdditionalProperties: make(map[string]string),
				},
				Subject: "test@test.com",
				Start:   iso8601.New(overrideStart),
				End:     iso8601.New(overrideStart.Add(time.Minute * 2)),
			},
			wantValidateRequestResponse: &ahTypes.ValidateGrantResponse{},

			wantRequest: &access.Request{
				Status: access.APPROVED,
				RequestedTiming: access.Timing{
					Duration:  time.Minute,
					StartTime: &now,
				},
				OverrideTiming: &access.Timing{
					Duration:  time.Minute * 2,
					StartTime: &overrideStart,
				},
				Grant: &access.Grant{
					CreatedAt: clk.Now(),
					UpdatedAt: clk.Now(),
					Start:     overrideStart,
					End:       overrideStart.Add(time.Minute * 2),
					Subject:   "test@test.com",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			g := ahmocks.NewMockClientWithResponsesInterface(ctrl)
			g.EXPECT().ValidateGrantWithResponse(gomock.Any(), gomock.Any(), gomock.Eq(tc.wantPostGrantsWithResponseBody)).Return(tc.wantValidateRequestResponse, tc.withCreateGrantResponseErr).AnyTimes()

			g.EXPECT().PostGrantsWithResponse(gomock.Any(), gomock.Eq(tc.wantPostGrantsWithResponseBody)).Return(tc.withCreateGrantResponse, tc.withCreateGrantResponseErr).AnyTimes()
			c := ddbmock.New(t)
			c.MockQuery(&storage.GetUser{Result: &tc.withUser})

			s := Granter{
				AHClient: g,
				DB:       c,
				Clock:    clk,
				accessTokenChecker: testAccessTokenChecker{
					NeedsToken: tc.needsAccessToken,
					Err:        tc.needsAccessTokenErr,
				},
			}

			gotRequest, err := s.CreateGrant(context.Background(), tc.give)
			assert.Equal(t, tc.wantErr, err)

			assert.Equal(t, tc.wantRequest, gotRequest)
		})
	}
}
