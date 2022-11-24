# Common Fate releases

Our versioning follows [SemVer numbering](https://semver.org/). Currently both `gdeploy` and our CloudFormation templates are versioned and released together. In future, this may change if it's more convenient to use separate release cycles for each one.

Prior to a release being made, we create a 'Release Candidate' (RC) version. The tagging format of an RC is `<NEXT_VERSION>-rc<RC_NUMBER>`, where `NEXT_VERSION` is the version we are planning to release (e.g. `v0.2.0`) and `RC_NUMBER` is a number starting at 1 which is incremented for each subsequent RC.

An example of an RC tag is `v0.2.0-rc2`.

The RC is released to the `granted-test-releases-us-east-1` S3 bucket. The CloudFormation template can be found at:

```
s3://granted-test-releases-us-east-1/rc/<TAG>/Granted.template.json
```

Where `<TAG>` is the RC tag, e.g. `v0.2.0-rc2`.

## Production Releases

Once you have tagged a commit and pushed to the remote repo, it must not be changed. If there is an issue found after the tag has been pushed.

1. Create a PR to fix the issue.
2. Merge to main
3. Tag the new commit and increment the minor number `v0.1.1` -> `v0.1.2`
4. Run the release workflow for this new tag

Side effects from manipulating tags after they are created include

- go pkg.dev will refer to the original commit for the tag rather than a new commit if you delete the tag and re create it
- the cli build pipeline will refer to the first tag for the commit hash, so even if you create a new tag for the same commit, it will not work as expected

# Running Common Fate Locally

The goal with Common Fate was to keep local development environments as similar to deployments as possible. This makes spinning up a dev environment super simple.

Start by exporting some AWS credentials, using whatever method you like. Here we just run `assume` using [Granted](https://granted.dev/).

```
assume rolename
```

Then initiate the dev deployment by running:

```
mage deploy:dev
```

The command will ask some prompts for naming your dev deployment and which region to deploy some of the resources into.

It will create a CDK changeset and ask to continue with the provisioning, input _yes_.

Once it has completed successfully you will receive the following success message in your console:

```
 ✅  GrantedDev (common-fate-dev-deployment)

✨  Deployment time: 749.51s
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

If at any point you make updates that effect any of the deployed cloud resources you can run `mage deploy:dev` to update these.

**Never use `gdeploy create` or `gdeploy update` when working with a local deployment.**

## Extras

All additional commands to be run are all the same as if it was a live deployment using `gdeploy`.

### Setting up a new provider

To create a new provider run:

```
gdeploy provider add
```

- Docs for setting this up here: [Provider Docs](https://docs.commonfate.io/granted-approvals/providers/introduction)

### Setting up a SSO identity provider]

To create a new identity provider run:

```
gdeploy sso configure
```

- Docs for setting this up here: [SSO Docs](https://docs.commonfate.io/granted-approvals/sso/introduction)
