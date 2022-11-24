## Runtimes

The access handler has support for different runtimes, this may allow the access handle to run in different clouds eventually. For now, we have a local runtime and a lambda runtime.

Switch between the runtimes by setting the `COMMON_FATE_ACCESS_HANDLER_RUNTIME` variable in your `.env` file. The options are `lambda` or `local`

### Local

Local is used in local development, it is a No-Op runtime in which no calls are actually made to the provider to grant or revoke access. It simply logs to the terminal, uses `time.Sleep` and responds with success responses to the API.

### Lambda

The lambda runtime is built for AWS Lambda with AWS Step Functions. Since our lambda functions are all written in Go, they can be run locally when running the access handler.
If you are running `mage deploy:dev` locally it will set this environment variable to `lambda` by default.
