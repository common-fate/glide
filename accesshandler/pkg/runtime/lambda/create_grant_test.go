package lambda

import (
	"testing"
)

func TestCreateGrant(t *testing.T) {
	t.Skip("this needs to be replaced with an integration test which actually triggers a Step Functions invocation. The current test fails because of an invalid ARN.")
	// ctx := context.Background()

	// // Skip test if credentials are not set
	// c, err := config.LoadDefaultConfig(ctx)
	// if err != nil {
	// 	t.Skip(err)
	// }
	// creds, err := c.Credentials.Retrieve(ctx)
	// if err != nil || !creds.HasKeys() {
	// 	t.Skip(err)
	// }

	// r := Runtime{}
	// os.Setenv("COMMONFATE_STATE_MACHINE_ARN", "test:arn")

	// err = r.Init(ctx)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// g := types.CreateGrant{
	// 	Provider: "test",
	// 	Subject:  "test@acme.com",
	// 	Start:    iso8601.New(time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC)),
	// 	End:      iso8601.New(time.Date(2022, 1, 1, 10, 10, 0, 0, time.UTC)),
	// }

	// // testing time is 1st Jan 2022, 10:00am UTC
	// now := time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC)

	// vcg, err := g.Validate(ctx, now)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// _, err = r.CreateGrant(ctx, *vcg)
	// if err != nil {
	// 	t.Fatal(err)
	// }
}
