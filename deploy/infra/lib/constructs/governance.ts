import { Duration, Stack } from "aws-cdk-lib";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import { EventBus } from "aws-cdk-lib/aws-events";
import * as iam from "aws-cdk-lib/aws-iam";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import * as path from "path";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import { WebUserPool } from "./app-user-pool";
import * as cdk from "aws-cdk-lib";
import { AccessHandler } from "./access-handler";
import * as kms from "aws-cdk-lib/aws-kms";

interface Props {
  appName: string;
  userPool: WebUserPool;
  frontendUrl: string;
  accessHandler: AccessHandler;
  eventBusSourceName: string;
  eventBus: EventBus;
  adminGroupId: string;
  providerConfig: string;
  notificationsConfiguration: string;
  identityProviderSyncConfiguration: string;
  deploymentSuffix: string;
  remoteConfigUrl: string;
  remoteConfigHeaders: string;
  analyticsDisabled: string;
  analyticsUrl: string;
  analyticsLogLevel: string;
  analyticsDeploymentStage: string;
  dynamoTable: dynamodb.Table;
  apiGatewayWafAclArn: string;
}

export class Governance extends Construct {
  private _governanceLambda: lambda.Function;
  private _governanceApiGateway: apigateway.Resource;
  private _apigateway: apigateway.RestApi;

  private _dynamoTable: dynamodb.Table;
  private _KMSkey: cdk.aws_kms.Key;

  private readonly _restApiName: string;
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);

    this._dynamoTable = props.dynamoTable;

    //todo passthrough kmskey
    this._KMSkey = new kms.Key(this, "PaginationKMSKey", {
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      pendingWindow: cdk.Duration.days(7),
      enableKeyRotation: true,
      description:
        "Used for encrypting and decrypting pagination tokens for Common Fate",
    });

    this._restApiName = props.appName + "_governance";

    const code = lambda.Code.fromAsset(
      path.join(__dirname, "..", "..", "..", "..", "bin", "governance.zip")
    );

    this._governanceLambda = new lambda.Function(
      this,
      "GovernanceAPIHandlerFunction",
      {
        code,
        timeout: Duration.seconds(60),
        environment: {
          COMMONFATE_TABLE_NAME: this._dynamoTable.tableName,
          COMMONFATE_FRONTEND_URL: props.frontendUrl,
          COMMONFATE_COGNITO_USER_POOL_ID: props.userPool.getUserPoolId(),
          COMMONFATE_IDENTITY_PROVIDER: props.userPool.getIdpType(),
          COMMONFATE_ADMIN_GROUP: props.adminGroupId,
          COMMONFATE_MOCK_ACCESS_HANDLER: "false",
          COMMONFATE_ACCESS_HANDLER_URL: props.accessHandler.getApiUrl(),
          COMMONFATE_PROVIDER_CONFIG: props.providerConfig,
          // COMMONFATE_SENTRY_DSN: can be added here
          COMMONFATE_EVENT_BUS_ARN: props.eventBus.eventBusArn,
          COMMONFATE_EVENT_BUS_SOURCE: props.eventBusSourceName,
          COMMONFATE_IDENTITY_SETTINGS: props.identityProviderSyncConfiguration,
          COMMONFATE_PAGINATION_KMS_KEY_ARN: this._KMSkey.keyArn,
          COMMONFATE_ACCESS_HANDLER_EXECUTION_ROLE_ARN:
            props.accessHandler.getAccessHandlerExecutionRoleArn(),
          COMMONFATE_DEPLOYMENT_SUFFIX: props.deploymentSuffix,
          COMMONFATE_ACCESS_REMOTE_CONFIG_URL: props.remoteConfigUrl,
          COMMONFATE_REMOTE_CONFIG_HEADERS: props.remoteConfigHeaders,
          CF_ANALYTICS_DISABLED: props.analyticsDisabled,
          CF_ANALYTICS_URL: props.analyticsUrl,
          CF_ANALYTICS_LOG_LEVEL: props.analyticsLogLevel,
          CF_ANALYTICS_DEPLOYMENT_STAGE: props.analyticsDeploymentStage,
        },
        runtime: lambda.Runtime.GO_1_X,
        handler: "governance",
      }
    );
    this._dynamoTable.grantReadWriteData(this._governanceLambda);

    this._apigateway = new apigateway.RestApi(this, "RestAPI", {
      restApiName: this._restApiName,
    });

    const api = this._apigateway.root.addResource("gov");
    const governancev1 = api.addResource("v1");

    const lambdaProxy = governancev1.addResource("{proxy+}");
    lambdaProxy.addMethod(
      "ANY",
      new apigateway.LambdaIntegration(this._governanceLambda, {
        allowTestInvoke: false,
      }),
      { authorizationType: apigateway.AuthorizationType.IAM }
    );

    this._governanceApiGateway = governancev1;
  }

  getGovernanceApiURL(): string {
    // both prepend and append a / so we have to remove one out
    return (
      this._apigateway.url +
      this._governanceApiGateway.path.substring(
        1,
        this._governanceApiGateway.path.length
      )
    );
  }
}
