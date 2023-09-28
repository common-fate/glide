## Getting Started

The backend is a Go HTTP server which communicates with a DynamoDB table for persisting data.

We use a code generation library to generate all the go structs and HTTP endpoint stubs for us, based on our OpenAPI spec in `openapi.yml` [oapi-codegen](https://github.com/deepmap/oapi-codegen).

### Generating code / Updating the API spec

If you need to make changes to the API spec (creating a new response type or creating a new endpoint). This will need to be done in the openapi.yml file first.

Once you have made your changes you will want to run `make generate` from the root folder. This will generate the types, client, services noted above.

### Implementing endpoint stub

If you create a new endpoint and generate the code for it, you will need to make a stub for it in the `api` package as you'll get a compilation error because API no longer meets the ServerInterface interface. The missing function will be your new endpoint. Implement the function on API, ensuring that the function signature matches the ServerInterface interface.

- You can copy the function definition from the `ServerInterface` in `service.gen.go`. eg:

```
// (GET /api/v1/access-rules)
GetApiV1Rules(w http.ResponseWriter, r *http.Request)
```

### Implementing List API responses

When implementing a List API, first you will create a a listResponse type following our openapi-standards.
Then in the api stub, it is critical to follow this pattern for the response.
In Go, a nil slice will marshal to JSON as null rather than an empty array.
This breaks our frontend typing and causes errors.

To avoid this, always initialise the list response using make([]types.ResponseType, len(responseElements)) as below
This will make a non nil empty slice if there are no elements which marshals to [] in JSON.

```go
res := types.ListGroupsResponse{
    Groups: make([]types.Group, len(q.Result)),
}
for i, g := range q.Result {
    res.Groups[i] = g.ToAPI()
}
```

## Webhook Handler

When deployed, the backend API runs in AWS Lambda, behind AWS API Gateway. The API Gateway handles Cognito authentication.

We additionally have a webhook API for third party integrations like Slack. This API is defined in [cmd/lambda/webhook/handler.go](../../cmd/lambda/webhook/handler.go) and does not use Cognito authentication. Our CDK code [app-backend.ts](../../deploy/infra/lib/constructs/app-backend.ts) defines the API Gateway routing for this. Our routing rules are:

| Path                   | Handler     | Authentication |
| ---------------------- | ----------- | -------------- |
| `/api/v1/{proxy+}`     | Glide API   | Cognito        |
| `/webhook/v1/{proxy+}` | Webhook API | -              |

_Note: `{proxy+}` refers to the [API Gateway Lambda Proxy integration](https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html), where all subpaths still point to the same Lambda. So `/api/v1/grants/gra_123` will still be handled by the Glide API._

## Environment Variables

Convention for environment variables is any variable directly related to the common fate application are prefixed with COMMONFATE\_

Variables used by packages relating to global services at commonfate such as analytics will be prefixed with CF\_

Variables that are directly derived from the runtime environment should not be prefixed, e.g AWS_REGION
