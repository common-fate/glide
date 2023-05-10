import * as cdk from "aws-cdk-lib";

import { Construct } from "constructs";
import { AppBackend } from "./constructs/app-backend";
import { AppFrontend } from "./constructs/app-frontend";
import { WebUserPool } from "./constructs/app-user-pool";
import { Database } from "./constructs/database";
import * as kms from "aws-cdk-lib/aws-kms";

import { EventBus } from "./constructs/events";
import { DevEnvironmentConfig } from "./helpers/dev-accounts";
import { generateOutputs } from "./helpers/outputs";
import { IdentityProviderTypes } from "./helpers/registry";
import { Governance } from "./constructs/governance";
import { TargetGroupGranter } from "./constructs/targetgroup-granter";

interface Props extends cdk.StackProps {
  stage: string;
  cognitoDomainPrefix: string;
  idpType: IdentityProviderTypes;
  samlMetadataUrl: string;
  samlMetadata: string;
  remoteConfigUrl: string;
  remoteConfigHeaders: string;
  devConfig: DevEnvironmentConfig | null;
  notificationsConfiguration: string;
  identityProviderSyncConfiguration: string;
  adminGroupId: string;
  cloudfrontWafAclArn: string;
  apiGatewayWafAclArn: string;
  analyticsDisabled: string;
  analyticsUrl: string;
  analyticsLogLevel: string;
  analyticsDeploymentStage: string;
  shouldRunCronHealthCheckCacheSync: boolean;
  identityGroupFilter: string;
  idpSyncTimeoutSeconds: number;
  idpSyncSchedule: string;
  idpSyncMemory: number;
}

export class CommonFateStackDev extends cdk.Stack {
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id, props);

    const {
      stage,
      cognitoDomainPrefix,
      idpType,
      samlMetadataUrl,
      samlMetadata,
      devConfig,
      adminGroupId,
      notificationsConfiguration,
      identityProviderSyncConfiguration,
      remoteConfigUrl,
      remoteConfigHeaders,
      cloudfrontWafAclArn,
      apiGatewayWafAclArn,
      analyticsDisabled,
      analyticsUrl,
      analyticsLogLevel,
      analyticsDeploymentStage,
      identityGroupFilter,
      idpSyncTimeoutSeconds,
      idpSyncSchedule,
      idpSyncMemory,
    } = props;
    const appName = `common-fate-${stage}`;

    const db = new Database(this, "Database", {
      appName,
    });

    const cdn = new AppFrontend(this, "Frontend", {
      appName,
      // this will be unique for dev deployments
      stableName: appName,
    }).withDevCDN(stage, devConfig, cloudfrontWafAclArn);

    const userPool = new WebUserPool(this, "WebUserPool", {
      appName: appName,
      domainPrefix: cognitoDomainPrefix,
      frontendUrl: "https://" + cdn.getDomainName(),
      callbackUrls: cdn.getDevCallbackUrls(),
      idpType: idpType,
      samlMetadataUrl: samlMetadataUrl,
      samlMetadata: samlMetadata,
      devConfig,
    });

    const events = new EventBus(this, "EventBus", {
      appName: appName,
    });

    //KMS key is used in governance api as well as appBackend - both for tokinization for ddb use
    const kmsKey = new kms.Key(this, "PaginationKMSKey", {
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      pendingWindow: cdk.Duration.days(7),
      enableKeyRotation: true,
      description:
        "Used for encrypting and decrypting pagination tokens for Common Fate",
    });

    const governance = new Governance(this, "Governance", {
      appName: appName,
      kmsKey: kmsKey,

      providerConfig: "",

      dynamoTable: db.getTable(),
    });
    const targetGroupGranter = new TargetGroupGranter(
      this,
      "TargetGroupGranter",
      {
        eventBus: events.getEventBus(),
        eventBusSourceName: events.getEventBusSourceName(),
        dynamoTable: db.getTable(),
      }
    );
    const appBackend = new AppBackend(this, "API", {
      appName: appName,
      userPool: userPool,
      frontendUrl: "https://" + cdn.getDomainName(),
      governanceHandler: governance,
      eventBus: events.getEventBus(),
      eventBusSourceName: events.getEventBusSourceName(),
      adminGroupId,
      identityProviderSyncConfiguration: identityProviderSyncConfiguration,
      notificationsConfiguration: notificationsConfiguration,
      deploymentSuffix: stage,
      dynamoTable: db.getTable(),
      remoteConfigUrl,
      remoteConfigHeaders,
      apiGatewayWafAclArn,
      analyticsDisabled,
      analyticsUrl,
      analyticsLogLevel,
      analyticsDeploymentStage,
      kmsKey: kmsKey,
      idpSyncMemory: idpSyncMemory,
      idpSyncSchedule: idpSyncSchedule,
      idpSyncTimeoutSeconds: idpSyncTimeoutSeconds,
      shouldRunCronHealthCheckCacheSync:
        props.shouldRunCronHealthCheckCacheSync || false,
      targetGroupGranter: targetGroupGranter,
      identityGroupFilter,
    });

    /* Outputs */
    generateOutputs(this, {
      APILogGroupName: appBackend.getLogGroupName(),
      APIURL: appBackend.getRestApiURL(),
      CacheSyncFunctionName: appBackend.getCacheSync().getFunctionName(),
      CacheSyncLogGroupName: appBackend.getCacheSync().getLogGroupName(),
      CLIAppClientID: userPool.getCLIAppClient().userPoolClientId,
      CloudFrontDistributionID: cdn.getDistributionId(),
      CloudFrontDomain: cdn.getCloudFrontDomain(),
      CognitoClientID: userPool.getUserPoolClientId(),
      DynamoDBTable: appBackend.getDynamoTableName(),
      EventBusArn: events.getEventBus().eventBusArn,
      EventBusLogGroupName: events.getLogGroupName(),
      EventBusSource: events.getEventBusSourceName(),
      EventsHandlerConcurrentLogGroupName: appBackend
        .getEventHandler()
        .getConcurrentLogGroupName(),
      EventsHandlerSequentialLogGroupName: appBackend
        .getEventHandler()
        .getSequentialLogGroupName(),
      FrontendDomainOutput: cdn.getDomainName(),
      GovernanceURL: governance.getGovernanceApiURL(),
      GranterLogGroupName: targetGroupGranter.getLogGroupName(),
      GranterV2StateMachineArn: targetGroupGranter.getStateMachineARN(),
      HealthcheckFunctionName: appBackend.getHealthChecker().getFunctionName(),
      HealthcheckLogGroupName: appBackend.getHealthChecker().getLogGroupName(),
      IDPSyncExecutionRoleARN: appBackend.getIdpSync().getExecutionRoleArn(),
      IDPSyncFunctionName: appBackend.getIdpSync().getFunctionName(),
      IDPSyncLogGroupName: appBackend.getIdpSync().getLogGroupName(),
      PaginationKMSKeyARN: appBackend.getKmsKeyArn(),
      Region: this.region,
      RestAPIExecutionRoleARN: appBackend.getExecutionRoleArn(),
      S3BucketName: cdn.getBucketName(),
      SAMLIdentityProviderName:
        userPool.getSamlUserPoolClient()?.getUserPoolName() || "",
      SlackNotifierLogGroupName: appBackend
        .getNotifiers()
        .getSlackLogGroupName(),
      UserPoolDomain: userPool.getUserPoolLoginFQDN(),
      UserPoolID: userPool.getUserPoolId(),
      WebhookLogGroupName: appBackend.getWebhookLogGroupName(),
      WebhookURL: appBackend.getWebhookApiURL(),
    });
  }
}
