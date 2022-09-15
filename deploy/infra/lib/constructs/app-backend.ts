import * as cdk from "aws-cdk-lib";
import { Duration, Stack } from "aws-cdk-lib";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import * as kms from "aws-cdk-lib/aws-kms";
import { EventBus } from "aws-cdk-lib/aws-events";
import * as iam from "aws-cdk-lib/aws-iam";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import * as path from "path";
import { WebUserPool } from "./app-user-pool";
import { EventHandler } from "./event-handler";
import { IdpSync } from "./idp-sync";
import { Notifiers } from "./notifiers";
import { AccessHandler } from "./access-handler";

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
  dynamoTable: dynamodb.Table;
}

export class AppBackend extends Construct {
  private readonly _appName: string;
  private _dynamoTable: dynamodb.Table;
  private _lambda: lambda.Function;
  private _apigateway: apigateway.LambdaRestApi;
  private _notifiers: Notifiers;
  private _eventHandler: EventHandler;
  private _idpSync: IdpSync;
  private _KMSkey: cdk.aws_kms.Key;
  private _webhook: apigateway.Resource;
  private _webhookLambda: lambda.Function;

  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);

    this._appName = props.appName;

    this._dynamoTable = props.dynamoTable;

    this._KMSkey = new kms.Key(this, "PaginationKMSKey", {
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      pendingWindow: cdk.Duration.days(7),
      description:
        "used for encrypting and decrypting pagination tokens for granted approvals",
    });

    // used to handle webhook events from third party integrations such as Slack
    this._webhookLambda = new lambda.Function(this, "WebhookHandlerFunction", {
      code: lambda.Code.fromAsset(
        path.join(__dirname, "..", "..", "..", "..", "bin", "webhook.zip")
      ),
      timeout: Duration.seconds(20),
      runtime: lambda.Runtime.GO_1_X,
      handler: "webhook",
      environment: {
        APPROVALS_TABLE_NAME: this._dynamoTable.tableName,
      },
    });

    this._dynamoTable.grantReadWriteData(this._webhookLambda);

    this._apigateway = new apigateway.RestApi(this, "RestAPI", {
      restApiName: this._appName,
    });

    const webhook = this._apigateway.root.addResource("webhook");
    const webhookv1 = webhook.addResource("v1");

    const webhookProxy = webhookv1.addResource("{proxy+}");
    webhookProxy.addMethod(
      "ANY",
      new apigateway.LambdaIntegration(this._webhookLambda, {
        allowTestInvoke: false,
      })
    );

    this._webhook = webhookv1;

    const code = lambda.Code.fromAsset(
      path.join(__dirname, "..", "..", "..", "..", "bin", "approvals.zip")
    );

    this._lambda = new lambda.Function(this, "RestAPIHandlerFunction", {
      code,
      timeout: Duration.seconds(60),
      environment: {
        APPROVALS_TABLE_NAME: this._dynamoTable.tableName,
        APPROVALS_FRONTEND_URL: props.frontendUrl,
        APPROVALS_COGNITO_USER_POOL_ID: props.userPool.getUserPoolId(),
        IDENTITY_PROVIDER: props.userPool.getIdpType(),
        APPROVALS_ADMIN_GROUP: props.adminGroupId,
        MOCK_ACCESS_HANDLER: "false",
        ACCESS_HANDLER_URL: props.accessHandler.getApiGateway().url,
        PROVIDER_CONFIG: props.providerConfig,
        // SENTRY_DSN: can be added here
        EVENT_BUS_ARN: props.eventBus.eventBusArn,
        EVENT_BUS_SOURCE: props.eventBusSourceName,
        IDENTITY_SETTINGS: props.identityProviderSyncConfiguration,
        PAGINATION_KMS_KEY_ARN: this._KMSkey.keyArn,
        ACCESS_HANDLER_EXECUTION_ROLE_ARN: props.accessHandler.getAccessHandlerExecutionRoleArn(),
        DEPLOYMENT_SUFFIX: props.deploymentSuffix,
        REMOTE_CONFIG_URL: props.remoteConfigUrl,
      },
      runtime: lambda.Runtime.GO_1_X,
      handler: "approvals",
    });

    this._KMSkey.grantEncryptDecrypt(this._lambda);

    this._lambda.addToRolePolicy(
      new PolicyStatement({
        resources: [props.userPool.getUserPool().userPoolArn],
        actions: [
          "cognito-idp:AdminListGroupsForUser",
          "cognito-idp:ListUsers",
          "cognito-idp:ListGroups",
          "cognito-idp:ListUsersInGroup",
          "cognito-idp:AdminGetUser",
          "cognito-idp:AdminListUserAuthEvents",
          "cognito-idp:AdminUserGlobalSignOut",
          "cognito-idp:DescribeUserPool",
          "cognito-idp:AdminAddUserToGroup",
          "cognito-idp:AdminCreateUser",
          "cognito-idp:CreateGroup",
          "cognito-idp:AdminRemoveUserFromGroup",
        ],
      })
    );
    this._lambda.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ["ssm:GetParameter", "ssm:PutParameter"],
        resources: [
          `arn:aws:ssm:${Stack.of(this).region}:${
            Stack.of(this).account
          }:parameter/granted/secrets/identity/*`,
        ],
      })
    );

    // allow the Approvals API to write SSM parameters as part of the guided setup workflow.
    this._lambda.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ["ssm:PutParameter"],
        resources: [
          `arn:aws:ssm:${Stack.of(this).region}:${
            Stack.of(this).account
          }:parameter/granted/providers/*`,
        ],
      })
    );

    const api = this._apigateway.root.addResource("api");
    const apiv1 = api.addResource("v1");

    const lambdaProxy = apiv1.addResource("{proxy+}");
    lambdaProxy.addMethod(
      "ANY",
      new apigateway.LambdaIntegration(this._lambda, {
        allowTestInvoke: false,
      }),
      {
        authorizationType: apigateway.AuthorizationType.COGNITO,
        authorizer: new apigateway.CognitoUserPoolsAuthorizer(
          this,
          "Authorizer",
          {
            cognitoUserPools: [props.userPool.getUserPool()],
          }
        ),
      }
    );

    const ALLOWED_HEADERS = [
      "Content-Type",
      "X-Amz-Date",
      "X-Amz-Security-Token",
      "Authorization",
      "X-Api-Key",
      "X-Requested-With",
      "Accept",
      "Access-Control-Allow-Methods",
      "Access-Control-Allow-Origin",
      "Access-Control-Allow-Headers",
    ];

    const standardCorsMockIntegration = new apigateway.MockIntegration({
      integrationResponses: [
        {
          statusCode: "200",
          responseParameters: {
            "method.response.header.Access-Control-Allow-Headers": `'${ALLOWED_HEADERS.join(
              ","
            )}'`,
            "method.response.header.Access-Control-Allow-Origin": "'*'",
            "method.response.header.Access-Control-Allow-Credentials":
              "'false'",
            "method.response.header.Access-Control-Allow-Methods":
              "'OPTIONS,GET,PUT,POST,DELETE'",
          },
        },
      ],
      passthroughBehavior: apigateway.PassthroughBehavior.NEVER,
      requestTemplates: {
        "application/json": '{"statusCode": 200}',
      },
    });

    const optionsMethodResponse = {
      statusCode: "200",
      responseModels: {
        "application/json": apigateway.Model.EMPTY_MODEL,
      },
      responseParameters: {
        "method.response.header.Access-Control-Allow-Headers": true,
        "method.response.header.Access-Control-Allow-Methods": true,
        "method.response.header.Access-Control-Allow-Credentials": true,
        "method.response.header.Access-Control-Allow-Origin": true,
      },
    };

    lambdaProxy.addMethod("OPTIONS", standardCorsMockIntegration, {
      authorizationType: apigateway.AuthorizationType.NONE,
      methodResponses: [optionsMethodResponse],
    });

    this._dynamoTable.grantReadWriteData(this._lambda);

    // Grant the approvals app access to invoke the access handler api
    this._lambda.addToRolePolicy(
      new PolicyStatement({
        resources: [props.accessHandler.getApiGateway().arnForExecuteApi()],
        actions: ["execute-api:Invoke"],
      })
    );
    props.eventBus.grantPutEventsTo(this._lambda);

    this._eventHandler = new EventHandler(this, "EventHandler", {
      dynamoTable: this._dynamoTable,
      eventBus: props.eventBus,
      eventBusSourceName: props.eventBusSourceName,
    });
    this._notifiers = new Notifiers(this, "Notifiers", {
      dynamoTable: this._dynamoTable,
      eventBus: props.eventBus,
      eventBusSourceName: props.eventBusSourceName,
      frontendUrl: props.frontendUrl,
      userPool: props.userPool,
      notificationsConfig: props.notificationsConfiguration,
    });

    this._idpSync = new IdpSync(this, "IdpSync", {
      dynamoTable: this._dynamoTable,
      userPool: props.userPool,
      identityProviderSyncConfiguration:
        props.identityProviderSyncConfiguration,
    });
  }

  getApprovalsApiURL(): string {
    return this._apigateway.url;
  }

  getWebhookApiURL(): string {
    // both prepend and append a / so we have to remove one out
    return (
      this._apigateway.url +
      this._webhook.path.substring(1, this._webhook.path.length)
    );
  }

  getDynamoTableName(): string {
    return this._dynamoTable.tableName;
  }
  getDynamoTable(): dynamodb.Table {
    return this._dynamoTable;
  }
  getWebhookLogGroupName(): string {
    return this._webhookLambda.logGroup.logGroupName;
  }
  getLogGroupName(): string {
    return this._lambda.logGroup.logGroupName;
  }
  getEventHandler(): EventHandler {
    return this._eventHandler;
  }
  getNotifiers(): Notifiers {
    return this._notifiers;
  }
  getIdpSync(): IdpSync {
    return this._idpSync;
  }

  getKmsKeyArn(): string {
    return this._KMSkey.keyArn;
  }
}
