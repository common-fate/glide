## Testing

The access handler has generic testing framework for providers, it consists of two steps.

### Fixture generation

There is a fixture [generation CLI](../../accesshandler/cmd/gdk/main.go)
Checkout the Okta provider for an example of how the integration tests work with fixture generation.

### Test

Checkout the [Okta tests](../../accesshandler/pkg/providers/okta/okta_test.go) for an example of hour the provider integration tests work.

[integration.RunTests](../../accesshandler/pkg/providertest/integration/integration.go)
