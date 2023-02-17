

# Target Group Deployment Workflow

## Getting started (outline)

To begin with you'll need to download: 
- [provider regsitry github repo](https://github.com/common-fate/provider-registry)
- [testvault provider github repo](https://github.com/common-fate/testvault-provider)
- [cli github repo](https://github.com/common-fate/cli)

Package and upload testvault

```bash
# cd into provider-registry directory
assume demo-sandbox1 --env

# run the registry
go run cmd/local-dev-server/main.go

# ensure you have ../testvault-provider
go run cmd/cli/main.go package --path ../testvault-provider/provider --name cf --version v1.0.0 --publisher testvault
go run cmd/cli/main.go upload --publisher testvault --name cf --version v1.0.0                                      
```


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