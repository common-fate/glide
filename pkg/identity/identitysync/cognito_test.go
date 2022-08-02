package identitysync

// // TestListUsers is a smoke test to check that we don't get an
// // error when reading from Cognito.
// func TestListUsers(t *testing.T) {
// 	c := newTestingCognito(t)
// 	ctx := context.Background()
// 	_, err := c.ListUsers(ctx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// // newTestingCognito creates a new testing Cognito.
// // It skips the tests if the APPROVALS_COGNITO_USER_POOL_ID env var isn't set.
// func newTestingCognito(t *testing.T) *Cognito {
// 	ctx := context.Background()
// 	_ = godotenv.Load("../../../.env")
// 	poolID := os.Getenv("APPROVALS_COGNITO_USER_POOL_ID")
// 	if poolID == "" {
// 		t.Skip("APPROVALS_COGNITO_USER_POOL_ID is not set")
// 	}
// 	s, err := NewCognito(ctx, Opts{
// 		UserPoolID: poolID,
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	return s
// }
