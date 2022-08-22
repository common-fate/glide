import * as cdk from "aws-cdk-lib";

import { Construct } from "constructs";
import { AccessHandler } from "./constructs/access-handler";
import { AppBackend } from "./constructs/app-backend";
import { AppFrontend } from "./constructs/app-frontend";
import { WebUserPool } from "./constructs/app-user-pool";

import { EventBus } from "./constructs/events";
import { DevEnvironmentConfig } from "./helpers/dev-accounts";
import { generateOutputs } from "./helpers/outputs";
import { IdentityProviderTypes } from "./helpers/registry";
interface Props extends cdk.StackProps {
  stage: string;
  cognitoDomainPrefix: string;
  idpType: IdentityProviderTypes;
  providerConfig: string;
  samlMetadataUrl: string;
  samlMetadata: string;
  devConfig: DevEnvironmentConfig | null;
  notificationsConfiguration: string;
  identityProviderSyncConfiguration: string;
  adminGroupId: string;
}
export class DevGrantedStack extends cdk.Stack {
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
    } = props;
    const appName = `granted-approvals-${stage}`;

    const cdn = new AppFrontend(this, "Frontend", {
      appName,
      // this will be unique for dev deployments
      stableName: appName,
    }).withDevCDN(stage, devConfig);

    const webUserPool = new WebUserPool(this, "WebUserPool", {
      appName: appName,
      domainPrefix: cognitoDomainPrefix,
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
    });

    const approvals = new AppBackend(this, "API", {
      appName: appName,
      userPool: webUserPool,
      frontendUrl: "https://" + cdn.getDomainName(),
      accessHandlerApi: accessHandler.getApiGateway(),
      eventBus: events.getEventBus(),
      eventBusSourceName: events.getEventBusSourceName(),
      adminGroupId,
      identityProviderSyncConfiguration: identityProviderSyncConfiguration,
      notificationsConfiguration: notificationsConfiguration,
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
      APIURL: approvals.getApprovalsApiURL(),
      APILogGroupName: approvals.getLogGroupName(),
      IDPSyncLogGroupName: approvals.getIdpSync().getLogGroupName(),
      AccessHandlerLogGroupName: accessHandler.getLogGroupName(),
      EventBusLogGroupName: events.getLogGroupName(),
      EventsHandlerLogGroupName: approvals.getEventHandler().getLogGroupName(),
      GranterLogGroupName: accessHandler.getGranter().getLogGroupName(),
      SlackNotifierLogGroupName: approvals
        .getNotifiers()
        .getSlackLogGroupName(),
      DynamoDBTable: approvals.getDynamoTableName(),
      GranterStateMachineArn: accessHandler.getGranter().getStateMachineARN(),
      EventBusArn: events.getEventBus().eventBusArn,
      EventBusSource: events.getEventBusSourceName(),
      IdpSyncFunctionName: approvals.getIdpSync().getFunctionName(),
      Region: this.region,
      KmsKey: approvals.getKmsKeyArn()
    });
  }
}
