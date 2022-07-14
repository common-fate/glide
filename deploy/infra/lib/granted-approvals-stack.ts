import * as cdk from "aws-cdk-lib";

import { Construct } from "constructs";
import { AccessHandler } from "./constructs/access-handler";
import { AppBackend } from "./constructs/app-backend";
import { AppFrontend } from "./constructs/app-frontend";
import { WebUserPool } from "./constructs/app-user-pool";

import { EventBus } from "./constructs/events";
import { DevEnvironmentConfig } from "./helpers/dev-accounts";

interface Props extends cdk.StackProps {
  stage: string;
  cognitoDomainPrefix: string;
  idpType: string;
  providerConfig: string;
  samlMetadataUrl: string;
  samlMetadata: string;
  devConfig: DevEnvironmentConfig | null;
  slackConfiguration: string;
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
      slackConfiguration,
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
      slackConfiguration: slackConfiguration,
    });
    /* Outputs */

    new cdk.CfnOutput(this, "CognitoClientID", {
      value: webUserPool.getUserPoolClientId(),
    });

    new cdk.CfnOutput(this, "CloudFrontDomain", {
      value: cdn.getCloudFrontDomain(),
    });

    new cdk.CfnOutput(this, "FrontendDomain", {
      value: cdn.getDomainName(),
    });

    new cdk.CfnOutput(this, "CloudFrontDistributionID", {
      value: cdn.getDistributionId(),
    });

    new cdk.CfnOutput(this, "S3BucketName", {
      value: cdn.getBucketName(),
    });

    new cdk.CfnOutput(this, "UserPoolID", {
      value: webUserPool.getUserPoolId(),
    });

    new cdk.CfnOutput(this, "UserPoolDomain", {
      value: webUserPool.getUserPoolLoginFQDN(),
    }).node.addDependency(webUserPool);

    new cdk.CfnOutput(this, "APIURL", {
      value: approvals.getApprovalsApiURL(),
    });

    new cdk.CfnOutput(this, "APILogGroupName", {
      value: approvals.getLogGroupName(),
    });
    new cdk.CfnOutput(this, "IDPSyncLogGroupName", {
      value: approvals.getIdpSync().getLogGroupName(),
    });
    new cdk.CfnOutput(this, "AccessHandlerLogGroupName", {
      value: accessHandler.getLogGroupName(),
    });

    new cdk.CfnOutput(this, "EventBusLogGroupName", {
      value: events.getLogGroupName(),
    });
    new cdk.CfnOutput(this, "EventsHandlerLogGroupName", {
      value: approvals.getEventHandler().getLogGroupName(),
    });

    new cdk.CfnOutput(this, "GranterLogGroupName", {
      value: accessHandler.getGranter().getLogGroupName(),
    });

    new cdk.CfnOutput(this, "SlackNotifierLogGroupName", {
      value: approvals.getNotifiers().getSlackLogGroupName(),
    });

    new cdk.CfnOutput(this, "DynamoDBTable", {
      value: approvals.getDynamoTableName(),
    });

    new cdk.CfnOutput(this, "GranterStateMachineArn", {
      value: accessHandler.getGranter().getStateMachineARN(),
    });
    new cdk.CfnOutput(this, "EventBusArn", {
      value: events.getEventBus().eventBusArn,
    });
    new cdk.CfnOutput(this, "EventBusSource", {
      value: events.getEventBusSourceName(),
    });
    new cdk.CfnOutput(this, "IdpSyncFunctionName", {
      value: approvals.getIdpSync().getFunctionName(),
    });
    new cdk.CfnOutput(this, "Region", {
      value: this.region,
    });
  }
}
