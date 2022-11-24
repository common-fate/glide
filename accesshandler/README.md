# access-handler

Common Fate Access Handler.

## Getting started

Copy the .env file:

```bash
cp .env.template .env
```

Run the server:

```bash
go run cmd/server/main.go
```

By default the server will run on http://localhost:9092.

## Testing

Run tests:

```
go test ./...
```

To check coverage:

```bash
go test ./... -coverprofile cover.out
go tool cover -html=cover.out -o cover.html
# open cover.html in your browser
```

## Editing API endpoints

To add a new endpoint, follow the below steps:

1. Edit `openapi.yaml` in this repository.

2. Run `make generate` to update the generated handler code. The code is generated into types.gen.go, and the function signatures can be found on the ServerInterface interface.

3. You'll get a compilation error because API no longer meets the ServerInterface interface. The missing function will be your new endpoint. Implement the function on API, ensuring that the function signature matches the ServerInterface interface.

## Running the access handler

1. Assume the cf-dev role by running:

   ```
   assume cf-dev --env
   ```

   To put the access credentials into your env file automatically

2. The access handler can be run using the `run access handler` vscode debug option or by running
   ```
   go run cmd/server/main.go
   ```
3. Run dev access commands
   Currently we just have the grant create command. Which can be executed by running
   ```
   go run cmd/cli/main.go grants create
   ```
   or the vscode debug profile `run access handler create grant`
   While the server is running
   - Note: Make sure you have the `COMMON_FATE_ACCESS_HANDLER_RUNTIME` environment variable set to `lambda` if you want to run the access handler against the live version of the access handler step functions.
   - It will otherwise just run locally and do nothing.
