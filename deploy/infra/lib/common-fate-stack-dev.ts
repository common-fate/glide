import * as cdk from "aws-cdk-lib";

import { Construct } from "constructs";
import { AccessHandler } from "./constructs/access-handler";
import { AppBackend } from "./constructs/app-backend";
import { AppFrontend } from "./constructs/app-frontend";
import { WebUserPool } from "./constructs/app-user-pool";
import { Database } from "./constructs/database";

import { EventBus } from "./constructs/events";
import { DevEnvironmentConfig } from "./helpers/dev-accounts";
import { generateOutputs } from "./helpers/outputs";
import { IdentityProviderTypes } from "./helpers/registry";
import { access } from "fs";

interface Props extends cdk.StackProps {
  stage: string;
  cognitoDomainPrefix: string;
  idpType: IdentityProviderTypes;
  providerConfig: string;
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

    const webUserPool = new WebUserPool(this, "WebUserPool", {
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

    const accessHandler = new AccessHandler(this, "AccessHandler", {
      appName: appName,
      eventBus: events.getEventBus(),
      eventBusSourceName: events.getEventBusSourceName(),
      providerConfig: props.providerConfig,
      remoteConfigUrl,
      remoteConfigHeaders,
    });

    const appBackend = new AppBackend(this, "API", {
      appName: appName,
      userPool: webUserPool,
      frontendUrl: "https://" + cdn.getDomainName(),
      accessHandler: accessHandler,
      eventBus: events.getEventBus(),
      eventBusSourceName: events.getEventBusSourceName(),
      adminGroupId,
      providerConfig: props.providerConfig,
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
    });
    /* Outputs */
    generateOutputs(this, {
      CognitoClientID: webUserPool.getUserPoolClientId(),
      CloudFrontDomain: cdn.getCloudFrontDomain(),
      FrontendDomainOutput: cdn.getDomainName(),
      CloudFrontDistributionID: cdn.getDistributionId(),
      S3BucketName: cdn.getBucketName(),
      UserPoolID: webUserPool.getUserPoolId(),
      UserPoolDomain: webUserPool.getUserPoolLoginFQDN(),
      APIURL: appBackend.getRestApiURL(),
      WebhookURL: appBackend.getWebhookApiURL(),
      GovernanceURL: appBackend.getGovernanceApiURL(),
      APILogGroupName: appBackend.getLogGroupName(),
      WebhookLogGroupName: appBackend.getWebhookLogGroupName(),
      IDPSyncLogGroupName: appBackend.getIdpSync().getLogGroupName(),
      AccessHandlerLogGroupName: accessHandler.getLogGroupName(),
      EventBusLogGroupName: events.getLogGroupName(),
      EventsHandlerLogGroupName: appBackend.getEventHandler().getLogGroupName(),
      GranterLogGroupName: accessHandler.getGranter().getLogGroupName(),
      SlackNotifierLogGroupName: appBackend
        .getNotifiers()
        .getSlackLogGroupName(),
      DynamoDBTable: appBackend.getDynamoTableName(),
      GranterStateMachineArn: accessHandler.getGranter().getStateMachineARN(),
      EventBusArn: events.getEventBus().eventBusArn,
      EventBusSource: events.getEventBusSourceName(),
      IdpSyncFunctionName: appBackend.getIdpSync().getFunctionName(),
      SAMLIdentityProviderName:
        webUserPool.getSamlUserPoolClient()?.getUserPoolName() || "",
      Region: this.region,
      PaginationKMSKeyARN: appBackend.getKmsKeyArn(),
      AccessHandlerExecutionRoleARN:
        accessHandler.getAccessHandlerExecutionRoleArn(),
      CacheSyncLogGroupName: appBackend.getCacheSync().getLogGroupName(),
      IDPSyncExecutionRoleARN: appBackend.getIdpSync().getExecutionRoleArn(),
      RestAPIExecutionRoleARN: appBackend.getExecutionRoleArn(),
      CacheSyncFunctionName: appBackend.getCacheSync().getFunctionName(),
    });
  }
}
