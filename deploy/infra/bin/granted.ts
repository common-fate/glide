#!/usr/bin/env node
import { App, DefaultStackSynthesizer } from "aws-cdk-lib";
import "source-map-support/register";
import { DevGrantedStack } from "../lib/granted-approvals-stack";
import { CustomerGrantedStack } from "../lib/granted-approvals-stack-prod";
import {
  DevEnvironmentConfig,
  DevEnvironments,
} from "../lib/helpers/dev-accounts";
import { IdentityProviderRegistry } from "../lib/helpers/registry";

const app = new App();
const stage = app.node.tryGetContext("stage");
const cognitoDomainPrefix = app.node.tryGetContext("cognitoDomainPrefix");
const idpType = app.node.tryGetContext("idpType");
const samlMetadataUrl = app.node.tryGetContext("samlMetadataUrl");
const samlMetadata = app.node.tryGetContext("samlMetadata");
const adminGroupId = app.node.tryGetContext("adminGroupId");
const providerConfig = app.node.tryGetContext("providerConfiguration");
const identityConfig = app.node.tryGetContext("identityConfiguration");
const notificationsConfiguration = app.node.tryGetContext("notificationsConfiguration");
const productionReleasesBucket = app.node.tryGetContext(
  "productionReleasesBucket"
);
const productionReleaseBucketPrefix = app.node.tryGetContext(
  "productionReleaseBucketPrefix"
);

// https://github.com/aws/aws-cdk/issues/11625
// cdk processes both stacks event if you specify only one
// To prepare the prod stack only, set the env var to "prod"
const stackTarget = process.env.STACK_TARGET || "dev";

if (stackTarget === "dev") {
  // devEnvironment is used to set the environment for internal
  // development deployments of the Approvals stack.
  const devEnvironmentName = app.node.tryGetContext("devEnvironment");

  let devConfig: DevEnvironmentConfig | null = null;

  if (devEnvironmentName !== undefined) {
    const conf = DevEnvironments.get(devEnvironmentName);

    if (conf === undefined) {
      throw new Error(`invalid dev environment name: ${devEnvironmentName}`);
    }
    devConfig = conf;
  }

  new DevGrantedStack(app, "GrantedDev", {
    cognitoDomainPrefix,
    stage,
    providerConfig: providerConfig || "{}",
    // We have inadvertently propagated this "granted-approvals-" through our dev tooling, so if we want to change this then it needs to be changed everywhere
    stackName: "granted-approvals-" + stage,
    idpType: idpType || IdentityProviderRegistry.CognitoV1Key,
    samlMetadataUrl: samlMetadataUrl || "",
    devConfig,
    adminGroupId: adminGroupId || "granted_administrators",
    samlMetadata: samlMetadata || "",
    notificationsConfiguration: notificationsConfiguration || "[]",
    identityProviderSyncConfiguration: identityConfig || "[]",
  });
} else if (stackTarget === "prod") {
  new CustomerGrantedStack(app, "Granted", {
    productionReleasesBucket: productionReleasesBucket,
    productionFrontendAssetObjectPrefix:
      productionReleaseBucketPrefix + "/frontend-assets",
    synthesizer: new DefaultStackSynthesizer({
      generateBootstrapVersionRule: false,
      fileAssetsBucketName: productionReleasesBucket,
      bucketPrefix: productionReleaseBucketPrefix + "/",
      // This role ARN is critical, make sure that it is in the account that you are using to publish!
      // cdk-assets will return obscure error messages if you are using teh wrong account to publish.
      // If you want to push to a new bucket in a different account, remember to setup a publishing role in that account that can be assumed with the credentials that you are using.
      fileAssetPublishingRoleArn:
        "arn:${AWS::Partition}:iam::${AWS::AccountId}:role/granted-cdk-asset-publishing-role",
    }),
  });
}
