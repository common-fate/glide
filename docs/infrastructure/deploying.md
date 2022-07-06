# Deploying the stack

## Prerequisites

- Make sure you have [Granted OSS](https://github.com/common-fate/granted) installed to be able to assume roles:
- Have the cf-dev profile in your` ~/.aws` config

Bootstrap CDK into your account and region:

Skip this step if deploying to cf dev account.
This is a one time operation, if you want to remove the bootstrapping while testing, open the console and delete the "CDKToolkit" stack (change the account number to your desired account).

```
pnpm run cdk bootstrap aws://12345678900/us-east-1
```

## Watch mode

CDK can run in [watch mode](https://aws.amazon.com/blogs/developer/increasing-development-speed-with-cdk-watch/) which is super fast for development.

```bash
mage -v watch
```

CDK will continuously build and redeploy changes when it detects a change to files.

If CDK watches any build artifacts it will cause a redeploy loop to occur, where building changes cause another deployment over and over again. A sign of this occurring is in the CDK logs, such as:

```
Detected change to '../staticredirect/build/index.js' (type: change) while 'cdk deploy' is still running. Will queue for another deployment after this one finishes
```

If this happens, add the path to the "exclude" section in `deploy/infra/cdk.json`.

## Frontend

The fronted loads configuration variables at runtime to determine things like the API URL and the Cognito user pool to use.

In production the deployment process writes a file `aws-exports.json` into `next/public`. This file is uploaded to AWS S3 and contains the production configuration for the frontend.

The deployment process also writes a `aws-exports.js` file to be used in local development into `next/utils`.

To access the Admin routes in the app you will have to add the user to the admin group in Cognito.
