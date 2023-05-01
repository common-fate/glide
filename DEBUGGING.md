# Running the API with VS-Code debugger

1. Create the file `.vscode/launch.json` and copy this launch configuration in.

```
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Run CommonFate ",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "cmd/server/main.go",
      "envFile": "${workspaceFolder}/.env"

    }
  ]
}
```

2. Export AWS credentials to a .env file in the root of the repo
3. Run the debugger for this profile
