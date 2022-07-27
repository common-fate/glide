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
Here we just run `assume` using [Granted](https://granted.dev/)
```
assume rolename
```



```
mage deploy:dev
```
The command wil