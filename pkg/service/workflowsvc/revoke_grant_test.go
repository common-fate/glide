package workflowsvc

// func TestRevokeGrant(t *testing.T) {
// 	type testcase struct {
// 		name                          string
// 		withRevokeGrantResponseErr    error
// 		withUser                      *identity.User
// 		getRule                       rule.AccessRule
// 		giveRequest                   requests.Requestv2
// 		wantUserErr                   error
// 		want                          *requests.Requestv2
// 		revokerID                     string
// 		requestReviewers              []access.Reviewer
// 		withGetRuleVersionResponseErr error
// 	}
// 	clk := clock.NewMock()

// 	testcases := []testcase{
// 		{
// 			name: "ok",
// 			getRule: rule.AccessRule{ID: "rule1",
// 				Status: rule.ACTIVE,

// 				Description: "string",
// 				Name:        "string",
// 				Groups:      []string{"string"},
// 				Target: rule.Target{
// 					With:          map[string]string{},
// 					TargetGroupID: "123",
// 				}},
// 			giveRequest: requests.Requestv2{
// 				RequestedBy: "user1",
// 			},
// 			withRevokeGrantResponseErr: nil,
// 			withUser:                   &identity.User{Groups: []string{"testAdmin"}},
// 			wantUserErr:                nil,
// 			want:                       nil,
// 			revokerID:                  "test",
// 			requestReviewers:           []access.Reviewer{{ReviewerID: "123"}},
// 		},
// 		{
// 			name: "no grant",
// 			getRule: rule.AccessRule{ID: "rule1",
// 				Status: rule.ACTIVE,

// 				Description: "string",
// 				Name:        "string",
// 				Groups:      []string{"string"},
// 				Target: rule.Target{
// 					With:          map[string]string{},
// 					TargetGroupID: "123",
// 				}},
// 			giveRequest: requests.Requestv2{
// 				RequestedBy: "user1",
// 				Grant: &requests.Grantv2{
// 					Status: types.GrantStatus(types.GrantStatusACTIVE),
// 					End:    time.Now().Add(time.Hour),
// 				},
// 			},
// 			withRevokeGrantResponseErr: ErrNoGrant,
// 			withUser:                   &identity.User{Groups: []string{"testAdmin"}},
// 			wantUserErr:                nil,
// 			want:                       nil,
// 			revokerID:                  "test",
// 			requestReviewers:           []access.Reviewer{{ReviewerID: "123"}},
// 		},
// 		{
// 			name: "trying to revoke inactive grant",
// 			getRule: rule.AccessRule{ID: "rule1",
// 				Status: rule.ACTIVE,

// 				Description: "string",
// 				Name:        "string",
// 				Groups:      []string{"string"},
// 				Target: rule.Target{
// 					With:          map[string]string{},
// 					TargetGroupID: "123",
// 				}},
// 			giveRequest: requests.Requestv2{
// 				RequestedBy: "user1",
// 				Grant: &requests.Grantv2{
// 					Status: types.GrantStatusEXPIRED,
// 				},
// 			},
// 			withRevokeGrantResponseErr: ErrGrantInactive,
// 			withUser:                   &identity.User{Groups: []string{"testAdmin"}},
// 			wantUserErr:                nil,
// 			want:                       nil,
// 			revokerID:                  "test",
// 		},
// 		{
// 			name: "access rule version not found",
// 			getRule: rule.AccessRule{ID: "rule1",
// 				Status: rule.ACTIVE,

// 				Description: "string",
// 				Name:        "string",
// 				Groups:      []string{"string"},
// 				Target: rule.Target{
// 					With:          map[string]string{},
// 					TargetGroupID: "123",
// 				}},
// 			giveRequest: requests.Requestv2{
// 				RequestedBy: "user1",
// 				Grant: &requests.Grantv2{
// 					Status: types.GrantStatus(types.GrantStatusACTIVE),
// 					End:    time.Now().Add(time.Hour),
// 				},
// 			},
// 			withRevokeGrantResponseErr:    ddb.ErrNoItems,
// 			withUser:                      &identity.User{Groups: []string{"testAdmin"}},
// 			wantUserErr:                   nil,
// 			want:                          nil,
// 			revokerID:                     "test",
// 			requestReviewers:              []access.Reviewer{{ReviewerID: "123"}},
// 			withGetRuleVersionResponseErr: ddb.ErrNoItems,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			runtime := mocks.NewMockRuntime(ctrl)
// 			runtime.EXPECT().Revoke(gomock.Any(), gomock.Any()).Return(tc.withRevokeGrantResponseErr).AnyTimes()

// 			eventbus := mocks.NewMockEventPutter(ctrl)
// 			eventbus.EXPECT().Put(gomock.Any(), gomock.Any()).Return(tc.withRevokeGrantResponseErr).AnyTimes()

// 			c := ddbmock.New(t)
// 			c.MockQueryWithErr(&storage.GetUser{Result: tc.withUser}, tc.wantUserErr)
// 			c.MockQueryWithErr(&storage.GetAccessRuleCurrent{Result: &tc.getRule}, tc.withGetRuleVersionResponseErr)
// 			c.MockQueryWithErr(&storage.ListRequestReviewers{Result: tc.requestReviewers}, tc.wantUserErr)

// 			s := Service{
// 				Runtime:  runtime,
// 				DB:       c,
// 				Clk:      clk,
// 				Eventbus: eventbus,
// 			}

// 			gotRequest, err := s.Revoke(context.Background(), tc.giveRequest, tc.revokerID, tc.withUser.Email)
// 			assert.Equal(t, tc.withRevokeGrantResponseErr, err)

// 			assert.Equal(t, tc.want, gotRequest)
// 		})
// 	}
// }
