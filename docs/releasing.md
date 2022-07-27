# Granted Approvals releases

Our versioning follows [SemVer numbering](https://semver.org/). Currently both `gdeploy` and our CloudFormation templates are versioned and released together. In future, this may change if it's more convenient to use separate release cycles for each one.

Prior to a release being made, we create a 'Release Candidate' (RC) version. The tagging format of an RC is `<NEXT_VERSION>-rc<RC_NUMBER>`, where `NEXT_VERSION` is the version we are planning to release (e.g. `v0.2.0`) and `RC_NUMBER` is a number starting at 1 which is incremented for each subsequent RC.

An example of an RC tag is `v0.2.0-rc2`.

The RC is released to the `granted-test-releases-us-east-1` S3 bucket. The CloudFormation template can be found at:

```
s3://granted-test-releases-us-east-1/rc/<TAG>/Granted.template.json
```

Where `<TAG>` is the RC tag, e.g. `v0.2.0-rc2`.


# Running Granted Approvals Locally
The goals with Granted Approvals was to keep local development environments as similar to deployments as possible. 
This makes spinning up a dev environment super simple.

Start by exporting some AWS credentials, using whatever method you like.
Here we just run `assume` using [Granted](https://granted.dev/).
```
assume rolename
```


Then initiate the dev deployment by running: 
```
mage deploy:dev
```
The command will ask some prompts for naming your dev deployment and which region to deploy some of the resources into.

It will create a changeset and ask to continue with the provisioning, input yes

Once its completed successfully you will recieve the following success message in your console:
```
 ✅  GrantedDev (granted-approvals-dev-deployment)

✨  Deployment time: 149.51s
```
**Note: This can take up to 5 minutes to complete, it is highly recommended to go make a coffee during this process ☕**

The mage scripts will export all necessary variables from the CDK outputs and set them in your `.env` file

From here you should have successfully deployed a local dev environment and can run the server by running:
```
go run cmd/server/main.go
```
And the frontend locally by running:
```
cd web
pnpm install
pnpm dev
```