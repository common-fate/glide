## Runtimes

The access handler has support for different runtimes, this may allow the access handle to run in different Clouds eventually.
For now, we have a local runtime and a lambda runtime.

### Local

Local is used in local development, it is a No-Op runtime in which no calls are actually made to the provider to grant or revoke access, it simply logs to the terminal, uses time.Sleep and responds with success responses to the api.

### Lambda

The lambda runtime is built for AWS Lambda with AWS Step Functions.
