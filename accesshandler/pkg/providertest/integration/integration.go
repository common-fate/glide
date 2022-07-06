package integration

import (
	"context"
	"testing"
	"time"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/genv"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/sethvargo/go-retry"
	"github.com/stretchr/testify/assert"
)

// TestCase is a test case for running integration tests.
type TestCase struct {
	Name              string
	Subject           string
	Args              string
	WantValidationErr error
}

func WithProviderConfig(config []byte) func(*IntegrationTests) {
	return func(it *IntegrationTests) {
		it.providerConfig = config
	}
}

// RunTests runs standardised integration tests to check the behaviour of a Granted Provider.
// It tests validation, granting, and revoking of access.
//
// This should be used against the live version of any integration APIs - you shouldn't mock the API that you are
// trying to test access against.
//
// RunTests is the entrypoint to the integration testing package.
func RunTests(t *testing.T, ctx context.Context, providerName string, p providers.Accessor, testcases []TestCase, opts ...func(*IntegrationTests)) {
	it := new(providerName, p, testcases, opts...)
	it.run(t, ctx)
}

type IntegrationTests struct {
	testcases      []TestCase
	providerName   string
	p              providers.Accessor
	providerConfig []byte
}

// new creates a new IntegrationTests holder struct.
func new(providerName string, p providers.Accessor, testcases []TestCase, opts ...func(*IntegrationTests)) *IntegrationTests {
	it := &IntegrationTests{
		testcases:    testcases,
		providerName: providerName,
		p:            p,
	}

	for _, o := range opts {
		o(it)
	}

	return it
}

func (it *IntegrationTests) run(t *testing.T, ctx context.Context) {
	// configure the provider, if it supports it.
	if c, ok := it.p.(providers.Configer); ok {
		err := c.Config().Load(ctx, genv.JSONLoader{Data: it.providerConfig})
		if err != nil {
			t.Fatal(err)
		}
	}

	// initialise the provider, if it supports it.
	if c, ok := it.p.(providers.Initer); ok {
		err := c.Init(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, tc := range it.testcases {
		t.Run(tc.Name, func(t *testing.T) {

			t.Run("validate", func(t *testing.T) {
				v, ok := it.p.(providers.Validator)
				if !ok {
					t.Skip("Provider does not implement providers.Validator")
				} else {
					err := v.Validate(ctx, tc.Subject, []byte(tc.Args))
					if tc.WantValidationErr == nil {
						// we shouldn't get any validation errors.
						assert.NoError(t, err)
					} else {
						assert.EqualError(t, err, tc.WantValidationErr.Error())
					}
				}
			})

			t.Run("access", func(t *testing.T) {
				err := it.p.Grant(ctx, tc.Subject, []byte(tc.Args))
				AssertAccessError(t, tc.WantValidationErr, err, "granting access")

				if tc.WantValidationErr == nil {
					t.Run("check provisioned", func(t *testing.T) {
						checker, ok := it.p.(IsActiver)
						if !ok {
							t.Skip("Provider does not implement IsActiver")
						} else {
							err = CheckIsProvisioned(ctx, checker, true)
							if err != nil {
								t.Fatal(err)
							}
						}
					})
				}

				b := retry.NewFibonacci(time.Second)
				b = retry.WithMaxDuration(time.Second*30, b)
				err = retry.Do(ctx, b, func(ctx context.Context) error {
					return it.p.Revoke(ctx, tc.Subject, []byte(tc.Args))
				})
				AssertAccessError(t, tc.WantValidationErr, err, "revoking access")

				if tc.WantValidationErr == nil {
					t.Run("check revoked", func(t *testing.T) {
						checker, ok := it.p.(IsActiver)
						if !ok {
							t.Skip("Provider does not implement IsActiver")
						} else {
							err = CheckIsProvisioned(ctx, checker, false)
							if err != nil {
								t.Fatal(err)
							}
						}
					})
				}
			})
		})
	}
}
