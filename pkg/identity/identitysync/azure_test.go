package identitysync

// // TestListUsers is a smoke test to check that we don't get an
// // error when reading from Cognito.
// func TestAzureListUsers(t *testing.T) {
// 	c := newTestingAzure(t)
// 	ctx := context.Background()
// 	_, err := c.ListUsers(ctx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestAzureListGroups(t *testing.T) {
// 	c := newTestingAzure(t)
// 	ctx := context.Background()
// 	_, err := c.ListGroups(ctx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// // newTestingCognito creates a new testing Cognito.
// // It skips the tests if the APPROVALS_COGNITO_USER_POOL_ID env var isn't set.
// func newTestingAzure(t *testing.T) *AzureSync {
// 	ctx := context.Background()
// 	_ = godotenv.Load("../../../.env")
// 	tenant := os.Getenv("AZURE_TENANT_ID")
// 	if tenant == "" {
// 		t.Skip("AZURE_TENANT_ID is not set")
// 	}
// 	clientID := os.Getenv("AZURE_CLIENT_ID")
// 	if clientID == "" {
// 		t.Skip("AZURE_CLIENT_ID is not set")
// 	}
// 	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
// 	if clientSecret == "" {
// 		t.Skip("AZURE_CLIENT_SECRET is not set")
// 	}
// 	s, err := NewAzure(ctx, deploy.Azure{TenantID: tenant, ClientID: clientID, ClientSecret: clientSecret})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	return s
// }
