import * as cdk from "aws-cdk-lib";

import { Construct } from "constructs";
import { AccessHandler } from "./constructs/access-handler";
import { AppBackend } from "./constructs/app-backend";
import { AppFrontend } from "./constructs/app-frontend";
import { WebUserPool } from "./constructs/app-user-pool";

import { CfnParameter } from "aws-cdk-lib";
import { EventBus } from "./constructs/events";
import { ProductionFrontendDeployer } from "./constructs/production-frontend-deployer";

interface Props extends cdk.StackProps {
  productionReleasesBucket: string;
  productionFrontendAssetObjectPrefix: string;
}
export class CustomerGrantedStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id, props);

    const cognitoDomainPrefix = new CfnParameter(this, "CognitoDomainPrefix", {
      type: "String",
      description:
        "CognitoDomainPrefix is a globally unique cognito domain prefix.",
      minLength: 1,
    });

    const idpType = new CfnParameter(this, "IdentityProviderType", {
      type: "String",
      description:
        "Configure your identity provider, okta requires SamlSSOMetadataURL to be provided",
      default: "COGNITO",
      allowedValues: ["COGNITO", "OKTA", "GOOGLE"],
    });

    const samlMetadataUrl = new CfnParameter(this, "SamlSSOMetadataURL", {
      type: "String",
      description:
        "Add your metadata url here to enable SSO, optionally leave this empty and provide the full metadata xml via SamlSSOMetadata",
      default: "",
    });
    const samlMetadata = new CfnParameter(this, "SamlSSOMetadata", {
      type: "String",
      description:
        "Add your metadata here to enable SSO, optionally, leave this empty and provide a metadata url SamlSSOMetadataURL",
      default: "",
    });

    const grantedAdminGroupId = new CfnParameter(this, "AdministratorGroupID", {
      type: "String",
      description:
        "Required, if you are not using cognito for your users you will need to provide a group id from your IDP which will control who has access to the administrator functions.",
      default: "granted_administrators",
    });

    const suffix = new CfnParameter(this, "DeploymentSuffix", {
      type: "String",
      description:
        "An optional suffix to be added to deployed resources (useful for testing scenarios where multiple stacks are deployed to a single AWS account)",
      default: "",
    });

    const frontendDomain = new CfnParameter(this, "FrontendDomain", {
      type: "String",
      description:
        "An optional custom domain name for the Granted web application. If not provided, an auto-generated CloudFront URL will be used.",
      default: "",
    });

    const frontendCertificate = new CfnParameter(
      this,
      "FrontendCertificateARN",
      {
        type: "String",
        description:
          "The ARN of an ACM certificate in us-east-1 for the frontend URL. Must be set if 'FrontendDomain' is defined.",
        default: "",
      }
    );

    const providerConfig = new CfnParameter(this, "ProviderConfiguration", {
      type: "String",
      description: "The Access Provider configuration in JSON format",
      default: "",
    });
    const slackConfig = new CfnParameter(this, "SlackConfiguration", {
      type: "String",
      description: "The Slack notifications configuration in JSON format",
      default: "",
    });
    const identityConfig = new CfnParameter(this, "IdentityConfiguration", {
      type: "String",
      description: "The Identity Provider Sync configuration in JSON format",
      default: "",
    });

    const appName = this.stackName + suffix.valueAsString;

    const appFrontend = new AppFrontend(this, "Frontend", {
      appName,
      // this is the same for all prod synthesis, it means that you can only deploy this once per account in production mode event with the suffix.
      // because the suffix cannot be appended to a logical id as it is a token.
      // the logical id must remain static to avoid issues with updates
      stableName: this.stackName,
    }).withProdCDN({
      frontendDomain: frontendDomain.valueAsString,
      frontendCertificateArn: frontendCertificate.valueAsString,
    });

    const webUserPool = new WebUserPool(this, "WebUserPool", {
      appName,
      domainPrefix: cognitoDomainPrefix.valueAsString,
      callbackUrls: appFrontend.getProdCallbackUrls(),
      idpType: idpType.valueAsString,
      samlMetadataUrl: samlMetadataUrl.valueAsString,
      samlMetadata: samlMetadata.valueAsString,
      devConfig: null,
    });

    const events = new EventBus(this, "EventBus", {
      appName,
    });

    const accessHandler = new AccessHandler(this, "AccessHandler", {
      appName,
      eventBus: events.getEventBus(),
      eventBusSourceName: events.getEventBusSourceName(),
      providerConfig: providerConfig.valueAsString,
    });
    const appBackend = new AppBackend(this, "API", {
      appName,
      userPool: webUserPool,
      frontendUrl: "https://" + appFrontend.getDomainName(),
      accessHandlerApi: accessHandler.getApiGateway(),
      eventBus: events.getEventBus(),
      eventBusSourceName: events.getEventBusSourceName(),
      adminGroupId: grantedAdminGroupId.valueAsString,
      identityProviderSyncConfiguration: identityConfig.valueAsString,
      slackConfiguration: slackConfig.valueAsString,
    });

    new ProductionFrontendDeployer(this, "FrontendDeployer", {
      apiUrl: appBackend.getApprovalsApiURL(),
      cloudfrontDistributionId: appFrontend.getDistributionId(),
      frontendDomain: appFrontend.getDomainName(),
      frontendBucket: appFrontend.getBucket(),
      cognitoClientId: webUserPool.getUserPoolClientId(),
      cognitoUserPoolId: webUserPool.getUserPoolId(),
      userPoolDomain: webUserPool.getUserPoolLoginFQDN(),
      cfReleaseBucket: props.productionReleasesBucket,
      cfReleaseBucketFrontendAssetObjectPrefix:
        props.productionFrontendAssetObjectPrefix,
    });

    /* Outputs */

    new cdk.CfnOutput(this, "CognitoClientID", {
      value: webUserPool.getUserPoolClientId(),
    });

    new cdk.CfnOutput(this, "CloudFrontDomain", {
      value: appFrontend.getCloudFrontDomain(),
    });

    new cdk.CfnOutput(this, "FrontendDomainOutput", {
      value: appFrontend.getDomainName(),
    });

    new cdk.CfnOutput(this, "CloudFrontDistributionID", {
      value: appFrontend.getDistributionId(),
    });

    new cdk.CfnOutput(this, "S3BucketName", {
      value: appFrontend.getBucketName(),
    });

    new cdk.CfnOutput(this, "UserPoolID", {
      value: webUserPool.getUserPoolId(),
    });

    new cdk.CfnOutput(this, "UserPoolDomain", {
      value: webUserPool.getUserPoolLoginFQDN(),
    }).node.addDependency(webUserPool);

    new cdk.CfnOutput(this, "APIURL", {
      value: appBackend.getApprovalsApiURL(),
    });

    new cdk.CfnOutput(this, "APILogGroupName", {
      value: appBackend.getLogGroupName(),
    });
    new cdk.CfnOutput(this, "IDPSyncLogGroupName", {
      value: appBackend.getIdpSync().getLogGroupName(),
    });
    new cdk.CfnOutput(this, "AccessHandlerLogGroupName", {
      value: accessHandler.getLogGroupName(),
    });

    new cdk.CfnOutput(this, "EventBusLogGroupName", {
      value: events.getLogGroupName(),
    });
    new cdk.CfnOutput(this, "EventsHandlerLogGroupName", {
      value: appBackend.getEventHandler().getLogGroupName(),
    });

    new cdk.CfnOutput(this, "GranterLogGroupName", {
      value: accessHandler.getGranter().getLogGroupName(),
    });

    new cdk.CfnOutput(this, "SlackNotifierLogGroupName", {
      value: appBackend.getNotifiers().getSlackLogGroupName(),
    });
    new cdk.CfnOutput(this, "DynamoDBTable", {
      value: appBackend.getDynamoTableName(),
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
    new cdk.CfnOutput(this, "Region", {
      value: this.region,
    });
  }
}
