import * as cdk from "aws-cdk-lib";

import { Construct } from "constructs";
import { AccessHandler } from "./constructs/access-handler";
import { AppBackend } from "./constructs/app-backend";
import { AppFrontend } from "./constructs/app-frontend";
import { WebUserPool } from "./constructs/app-user-pool";

import { CfnParameter } from "aws-cdk-lib";
import { EventBus } from "./constructs/events";
import { ProductionFrontendDeployer } from "./constructs/production-frontend-deployer";
import { generateOutputs } from "./helpers/outputs";
import { IdentityProviderRegistry, IdentityProviderTypes } from "./helpers/registry";

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
      default: IdentityProviderRegistry.Cognito,
      allowedValues: Object.values(IdentityProviderRegistry),
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
      default: "{}",
    });
    const notificationsConfiguration = new CfnParameter(this, "NotificationsConfiguration", {
      type: "String",
      description: "The Notifications configuration in JSON format",
      default: "{}",
    });
    const identityConfig = new CfnParameter(this, "IdentityConfiguration", {
      type: "String",
      description: "The Identity Provider Sync configuration in JSON format",
      default: "{}",
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
      idpType: idpType.valueAsString as IdentityProviderTypes,
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
      notificationsConfiguration: notificationsConfiguration.valueAsString,
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
    generateOutputs(this, {
      CognitoClientID: webUserPool.getUserPoolClientId(),
      CloudFrontDomain: appFrontend.getCloudFrontDomain(),
      FrontendDomainOutput: appFrontend.getDomainName(),
      CloudFrontDistributionID: appFrontend.getDistributionId(),
      S3BucketName: appFrontend.getBucketName(),
      UserPoolID: webUserPool.getUserPoolId(),
      UserPoolDomain: webUserPool.getUserPoolLoginFQDN(),
      APIURL: appBackend.getApprovalsApiURL(),
      APILogGroupName: appBackend.getLogGroupName(),
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
      Region: this.region,
      AccessHandlerRestAPILambdaExecutionRoleARN: accessHandler.getAccessHandlerRestAPILambdaExecutionRoleARN(),
      GranterLambdaExecutionRoleARN: accessHandler.getGranter().getGranterLambdaExecutionRoleARN(),
      PaginationKMSKeyARN: appBackend.getKmsKeyArn()
    });
  }
}
