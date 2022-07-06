package rulesvc

import (
	"context"
	"net/http"
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/common-fate/ddb/ddbmock"
	ahTypes "github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types/ahmocks"
	"github.com/common-fate/granted-approvals/pkg/identity"
	"github.com/common-fate/granted-approvals/pkg/rule"
	"github.com/common-fate/granted-approvals/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccessRule(t *testing.T) {
	type testcase struct {
		name                 string
		givenUserID          identity.User
		give                 types.CreateAccessRuleRequest
		wantErr              error
		withProviderResponse ahTypes.Provider
		want                 *rule.AccessRule
	}

	in := types.CreateAccessRuleRequest{}

	ruleID := "override"
	versionID := "overrideVersion"
	userID := "user1"
	clk := clock.NewMock()
	now := clk.Now()

	mockRule := rule.AccessRule{
		ID:          ruleID,
		Version:     versionID,
		Approval:    rule.Approval(in.Approval),
		Status:      rule.ACTIVE,
		Description: in.Description,
		Name:        in.Name,
		Groups:      in.Groups,
		Metadata: rule.AccessRuleMetadata{
			CreatedAt: now,
			CreatedBy: userID,
			UpdatedAt: now,
			UpdatedBy: userID,
		},
		Target: rule.Target{
			ProviderID:   in.Target.ProviderId,
			ProviderType: "okta",
			With:         in.Target.With.AdditionalProperties,
		},
		TimeConstraints: in.TimeConstraints,
		Current:         true,
	}

	/**
	There are two test cases here:
	- Create a valid rule
	*/
	testcases := []testcase{
		{
			name:        "ok",
			givenUserID: identity.User{ID: userID},
			give:        in,
			want:        &mockRule,
			withProviderResponse: ahTypes.Provider{
				Id:   in.Target.ProviderId,
				Type: "okta",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {

			dbc := ddbmock.Client{
				PutErr: tc.wantErr,
			}

			clk := clock.NewMock()
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			m := ahmocks.NewMockClientWithResponsesInterface(ctrl)
			m.EXPECT().GetProviderWithResponse(gomock.Any(), gomock.Eq(tc.give.Target.ProviderId)).Return(&ahTypes.GetProviderResponse{
				JSON200: &tc.withProviderResponse,
				HTTPResponse: &http.Response{
					StatusCode: http.StatusOK,
				},
			}, nil)
			s := Service{
				Clock:    clk,
				DB:       &dbc,
				AHClient: m,
			}

			got, err := s.CreateAccessRule(context.Background(), &tc.givenUserID, tc.give)

			// This is the only thing from service layer that we can't mock yet, hence the override
			if err == nil {
				got.ID = ruleID
				got.Version = versionID
			}

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.want, got)

		})
	}

}
