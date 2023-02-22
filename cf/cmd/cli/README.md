# Target Group Deployment Workflow

## Getting started (outline)

The Common Fate provider framework is spread out across some microservices. You'll need to clone the following to get a wholistic view over the codebase:

- [Common Fate Provider Core Framework](https://github.com/common-fate/commonfate-provider-core)
- [Common Fate Provider Registry](https://github.com/common-fate/provider-registry)
- [Common Fate Provider Core Types SDK](https://github.com/common-fate/provider-registry-sdk-go)
- [cli github repo](https://github.com/common-fate/cli)


We have some self managed in house implementation of some open source python providers, heres a list of some:
- [testvault provider github repo](https://github.com/common-fate/testvault-provider)
- [AWS Identity Centre (SSO) provider github repo](https://github.com/common-fate/commonfate-provider-aws-sso)
- [Azure Groups provider github repo](https://github.com/common-fate/azure-provider)
- [Okta Groups provider github repo](https://github.com/common-fate/okta-provider)
- [Test Groups provider github repo](https://github.com/common-fate/commonfate-provider-testgroups)


## Interacting with Providers v2 Overview
Providers v2 can be interacted with through the cli tooling.
- pdk-cli in the Provider Registry repo
    - The pdk-cli can be made by running `make pdk-cli` in the provider registry repo
    - Handles packaging providers to zip files and uploading them to s3
    - The pdk-cli also has testing commands to run test. See more here TODO
- cfcli which is in the Common Fate repo (current repo)
    - The cfcli can be made by running `make cfcli` in the common fate repo
    - Handles bootstrapping deployment s3 buckets 
    - Handles creating cloudformation based on an loaded provider from the pdk-cli upload.
    - Handles manually creating target groups and target group deployments



## Package and upload a local provider

```bash
# cd into provider-registry directory and assume into the target account
assume demo-sandbox1 --env
# run the registry server
go run cmd/local-dev-server/main.go
# make sure you have a local copy of any provider, for this example we are using the testvault provider
# This assumes you have the provider in a folder one step away from the current directory
go run cmd/cli/main.go package --path ../testvault-provider/provider --name cf --version v1.0.0 --publisher testvault
go run cmd/cli/main.go upload --publisher testvault --name cf --version v1.0.0                                      
```
- This will:
    - Create a provider in the provider registry dynamo database (to be interacted with in the frontend later down the line)


Authenticate, create a deployment and target group, then link them together

```bash
# cd into common-fate repo
# pre-liminary req: ensure you're running a recent deployment `mage deploy:dev`
# fetch the cloudfront URL from gdeploy status 
gdeploy status 
# will return -> CloudFrontDomain | d1eyyid6r8dmqf.cloudfront.net  <- copy me            
# now log in using the deployment URL
cf login d1eyyid6r8dmqf.cloudfront.net
# now we can run cf auth'd commands
go run cmd/devcli/main.go deploy --runtime=aws --publisher=jack --version=v0.1.4 --name=testvault --accountId=012345678912 --aws-region=ap-southeast-2 --suffix=jacktest6
```
