package integration

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/common-fate/granted-approvals/accesshandler/pkg/diagnostics"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/providers"
	"github.com/common-fate/granted-approvals/accesshandler/pkg/types"
	"github.com/common-fate/granted-approvals/pkg/gconfig"
	"github.com/segmentio/ksuid"
	"github.com/sethvargo/go-retry"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

// TestCase is a test case for running integration tests.
type TestCase struct {
	Name                      string
	Subject                   string
	Args                      string
	WantValidationDiagnostics map[string]types.GrantValidation
	WantValidationErr         error
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
	if c, ok := it.p.(gconfig.Configer); ok {
		err := c.Config().Load(ctx, gconfig.JSONLoader{Data: it.providerConfig})
		if err != nil {
			t.Fatal(err)
		}
	}

	// initialise the provider, if it supports it.
	if c, ok := it.p.(gconfig.Initer); ok {
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
					// the provider implements validation, so try and validate the request
					res := v.ValidateGrant([]byte(tc.Args))

					validationRes := types.GrantValidationResponse{}
					var mu sync.Mutex
					handleResults := func(key string, value providers.GrantValidationStep, logs diagnostics.Logs) {
						mu.Lock()
						defer mu.Unlock()

						result := types.GrantValidation{
							Id: key,
						}

						if logs.HasSucceeded() {
							result.Status = types.GrantValidationStatusSUCCESS
						} else {
							result.Status = types.GrantValidationStatusERROR
						}

						for _, l := range logs {
							result.Logs = append(result.Logs, types.Log{
								Level: types.LogLevel(l.Level),
								Msg:   l.Msg,
							})
						}

						validationRes.Validation = append(validationRes.Validation, result)
					}

					g, gctx := errgroup.WithContext(ctx)

					for key, val := range res {
						k := key
						v := val
						g.Go(func() error {
							logs := v.Run(gctx, string(tc.Subject), []byte(tc.Args))
							handleResults(k, v, logs)
							return nil
						})
					}

					_ = g.Wait()

					for _, res := range validationRes.Validation {
						it := tc.WantValidationDiagnostics[res.Id]
						assert.Equal(t, res.Status, it.Status)

					}

				}
			})

			t.Run("access", func(t *testing.T) {
				testGrantID := ksuid.New().String()
				err := it.p.Grant(ctx, tc.Subject, []byte(tc.Args), testGrantID)
				AssertAccessError(t, tc.WantValidationErr, err, "granting access")

				if tc.WantValidationErr == nil {
					t.Run("check provisioned", func(t *testing.T) {
						checker, ok := it.p.(IsActiver)
						if !ok {
							t.Skip("Provider does not implement IsActiver")
						} else {
							err = CheckIsProvisioned(ctx, checker, tc.Subject, []byte(tc.Args), testGrantID, true)
							if err != nil {
								t.Fatal(err)
							}
						}
					})
				}

				b := retry.NewFibonacci(time.Second)
				b = retry.WithMaxDuration(time.Second*30, b)
				err = retry.Do(ctx, b, func(ctx context.Context) error {
					return it.p.Revoke(ctx, tc.Subject, []byte(tc.Args), testGrantID)
				})
				AssertAccessError(t, tc.WantValidationErr, err, "revoking access")

				if tc.WantValidationErr == nil {
					t.Run("check revoked", func(t *testing.T) {
						checker, ok := it.p.(IsActiver)
						if !ok {
							t.Skip("Provider does not implement IsActiver")
						} else {
							err = CheckIsProvisioned(ctx, checker, tc.Subject, []byte(tc.Args), testGrantID, false)
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
