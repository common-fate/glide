# Contributing

## Requirements

- [Mage](https://magefile.org/): install via `brew install mage`
- Go 1.19.0
- NodeJS
- [pnpm](https://pnpm.io/): install via `brew install pnpm`

Note: you can alternatively run `mage` natively with Go with `go run mage.go`.

## Getting started

1. Install NodeJS dependencies..

   ```bash
   pnpm install
   ```

1. Assume an AWS role with permission to create resources.

   ```bash
   assume cf-dev
   ```

1. Build the `gdeploy` CLI.

   ```bash
   make gdeploy
   ```

1. Deploy the development stack.

   ```bash
   mage -v deploy:dev
   ```

1. Setup some users in Cognito.

   ```bash
   gdeploy users create --admin -u you@gmail.com
   gdeploy users create -u you+1@gmail.com
   ```

1. Open the web dashboard.

   ```bash
   gdeploy dashboard open
   ```

## Tearing down

To clean up a deployment run:

```
mage -v destroy
```

## Cloud Development

You can deploy to AWS with hot-reloading as follows. Deployments are very fast (currently around 15-30s):

```
mage -v watch
```

## Documentation

- [Backend](./docs/backend/)
- [Frontend](./docs/frontend/)
- [Infrastructure](./docs/infrastructure/)
- [OpenAPI standards](./docs/openapi-standards/)
