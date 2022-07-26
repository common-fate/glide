## API

The backend is a Go HTTP server, which implements a generic interface on top of a [provider framework](providers.md).

We use a code generation library to generate all the go structs and HTTP endpoint stubs for us, based on our OpenAPI spec in `openapi.yml` [oapi-codegen](https://github.com/deepmap/oapi-codegen).

### Generating code / Updating the API spec

Refer to the [main backend docs](../backend/backend.md) as we use the same processes for the approvals backend and the access handler api.
