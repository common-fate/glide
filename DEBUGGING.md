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
# Launch Profiles

```json
{
    "version": "0.2.0",
    "configurations": [
      {
        "name": "Run CommonFate",
        "type": "go",
        "request": "launch",
        "mode": "auto",
        "program": "cmd/server/main.go",
        "envFile": "${workspaceFolder}/.env"
      }
    ]
  }

```


# Force Typescript Linting to Run For all Files in Web Directory

Copy this snippet into the file `.vscode/tasks.json`

```
{
  "version": "2.0.0",
  "tasks": [
    {
      "type": "typescript",
      "tsconfig": "web/tsconfig.json",
      "option": "watch",
      "problemMatcher": ["$tsc-watch"],
      "group": "build",
      "label": "tsc: watch - web/tsconfig.json"
    }
  ]
}

```
- ctrl + shift + P 
- Search for 'Tasks: run task'
- click on `tsc: watch - web/tsconfig.json`


# Workflows

## Connect a Local PDK Provider Without the Registry

1. in the provider-registry repo, pull the latest changes, the `make pdk`
2. open your provider repo
3. make sure you have updated to the latest version of the pdk package 
4. `source .venv/bin/activate`
5. `pip install commonfate-provider`
6. `pdk configure`
7. `pdk package`
8. assume aws credentials for the account where you want to deploy the handler lambda
9. use the dev deployment process to deploy the lambda handler `pdk devhandler deploy --id <Choose an ID> --confirm true` 
10. switch back to teh common fate repo
11. run `go run cmd/devcli/main.go targetgroup create --path=../cf-provider-testvault --id=<TargetGroupId> --kind=Vault --provider <Can be anything>`
12. Register the newly deployed handler with `cf handler register`
    a. The cfm stack id will be the same name as the id you used in step 9
13. Link the target group with the handler with `cf targetgroup link --handler-id <handler_id> --target-group-id <target_grouo id> --kind Vault`
14. Manually run the heathcheck go run `go run cmd/devcli/main.go healthcheck local`
15. Then the cache sync using `go run cmd/devcli/main.go cache sync`
