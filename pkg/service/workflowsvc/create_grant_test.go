package workflowsvc

// func TestCreateGrant(t *testing.T) {
// 	type testcase struct {
// 		name                       string
// 		withCreateGrantResponseErr error
// 		withUser                   *identity.User
// 		giveRule                   rule.AccessRule
// 		giveRequest                access.Request
// 		createGrant                types.CreateGrant
// 		wantErr                    error
// 		wantUserErr                error
// 		want                       *access.Grant
// 	}
// 	clk := clock.NewMock()

// 	testcases := []testcase{
// 		{
// 			name: "ok",
// 			createGrant: types.CreateGrant{
// 				Subject:  openapi_types.Email("test@commonfate.io"),
// 				Start:    iso8601.New(clk.Now().Add(time.Second * 2)),
// 				End:      iso8601.New(clk.Now().Add(time.Hour)),
// 				Provider: "test",
// 				Id:       types.NewGrantID(),
// 				With: types.CreateGrant_With{
// 					AdditionalProperties: map[string]string{
// 						"vault": "test",
// 					},
// 				}},
// 			giveRule: rule.AccessRule{ID: "rule1",
// 				Status: rule.ACTIVE,

// 				Description: "string",
// 				Name:        "string",
// 				Groups:      []string{"string"},
// 				Target: rule.Target{
// 					With:          map[string]string{},
// 					TargetGroupID: "123",
// 				}},
// 			giveRequest: access.Request{
// 				RequestedBy: "user1",
// 			},
// 			withCreateGrantResponseErr: nil,
// 			withUser:                   &identity.User{Groups: []string{"testAdmin"}},
// 			wantUserErr:                nil,
// 			want:                       &access.Grant{Provider: "string", Subject: "", With: types.Grant_With{AdditionalProperties: map[string]string{}}, Start: clk.Now(), End: clk.Now(), Status: "PENDING", CreatedAt: clk.Now(), UpdatedAt: clk.Now()},
// 		},
// 		{
// 			name: "user doesn't exist",
// 			createGrant: types.CreateGrant{
// 				Subject:  openapi_types.Email("test@commonfate.io"),
// 				Start:    iso8601.New(time.Now().Add(time.Second * 2)),
// 				End:      iso8601.New(time.Now().Add(time.Hour)),
// 				Provider: "test",
// 				Id:       types.NewGrantID(),
// 				With: types.CreateGrant_With{
// 					AdditionalProperties: map[string]string{
// 						"vault": "test",
// 					},
// 				}},
// 			giveRule: rule.AccessRule{ID: "rule1",
// 				Status: rule.ACTIVE,

// 				Description: "string",
// 				Name:        "string",
// 				Groups:      []string{"string"},
// 				Target: rule.Target{
// 					With:          map[string]string{},
// 					TargetGroupID: "123",
// 				}},
// 			giveRequest: access.Request{
// 				RequestedBy: "user1",
// 			},
// 			withCreateGrantResponseErr: nil,
// 			withUser:                   nil,
// 			want:                       nil,
// 			wantUserErr:                ddb.ErrNoItems,
// 			wantErr:                    ddb.ErrNoItems,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			runtime := mocks.NewMockRuntime(ctrl)
// 			runtime.EXPECT().Grant(gomock.Any(), gomock.Any()).Return(tc.withCreateGrantResponseErr).AnyTimes()

// 			eventbus := mocks.NewMockEventPutter(ctrl)
// 			eventbus.EXPECT().Put(gomock.Any(), gomock.Any()).Return(tc.withCreateGrantResponseErr).AnyTimes()

// 			c := ddbmock.New(t)
// 			c.MockQueryWithErr(&storage.GetUser{Result: tc.withUser}, tc.wantUserErr)

// 			s := Service{
// 				Runtime:  runtime,
// 				DB:       c,
// 				Clk:      clk,
// 				Eventbus: eventbus,
// 			}

// 			gotGrant, err := s.Grant(context.Background(), tc.giveRequest, tc.giveRule)
// 			assert.Equal(t, tc.wantErr, err)

// 			assert.Equal(t, tc.want, gotGrant)
// 		})
// 	}
// }
