package rulesvc

// func TestArchiveAccessRule(t *testing.T) {
// 	type testcase struct {
// 		name      string
// 		givenRule rule.AccessRule
// 		wantErr   error
// 		want      *rule.AccessRule
// 	}

// 	clk := clock.NewMock()
// 	now := clk.Now()
// 	mockRule := rule.AccessRule{
// 		ID: "rule",

// 		Status: rule.ACTIVE,
// 		Metadata: rule.AccessRuleMetadata{
// 			CreatedAt: now.Add(time.Minute),
// 			CreatedBy: "hello",
// 			UpdatedAt: now.Add(time.Minute),
// 			UpdatedBy: "hello",
// 		},
// 	}
// 	want := mockRule
// 	want.Status = rule.ARCHIVED
// 	want.Metadata.UpdatedAt = now

// 	testcases := []testcase{
// 		{
// 			name:      "ok",
// 			givenRule: mockRule,
// 			want:      &want,
// 		},
// 		{
// 			name: "already archived",
// 			givenRule: rule.AccessRule{
// 				Status: rule.ACTIVE,
// 			},
// 			wantErr: ErrAccessRuleAlreadyArchived,
// 		},
// 	}

// 	for _, tc := range testcases {
// 		t.Run(tc.name, func(t *testing.T) {

// 			db := ddbmock.New(t)
// 			db.PutBatchErr = tc.wantErr
// 			db.MockQueryWithErrWithResult(&storage.ListRequestsForStatus{Status: access.PENDING, Result: []access.Request{}}, &ddb.QueryResult{}, nil)
// 			db.MockQueryWithErrWithResult(&storage.ListRequestReviewers{Result: []access.Reviewer{}}, &ddb.QueryResult{}, nil)

// 			s := Service{
// 				Clock: clk,
// 				DB:    db,
// 			}

// 			got, err := s.ArchiveAccessRule(context.Background(), "", tc.givenRule)

// 			assert.Equal(t, tc.wantErr, err)
// 			assert.Equal(t, tc.want, got)
// 		})
// 	}
// }
