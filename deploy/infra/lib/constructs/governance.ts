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
  accessHandler: AccessHandler;
  providerConfig: string;
  dynamoTable: dynamodb.Table;
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
          COMMONFATE_MOCK_ACCESS_HANDLER: "false",
          COMMONFATE_ACCESS_HANDLER_URL: props.accessHandler.getApiUrl(),
          COMMONFATE_PROVIDER_CONFIG: props.providerConfig,

          COMMONFATE_PAGINATION_KMS_KEY_ARN: this._KMSkey.keyArn,
        },
        runtime: lambda.Runtime.GO_1_X,
        handler: "governance",
      }
    );
    this._dynamoTable.grantReadWriteData(this._governanceLambda);

    this._governanceLambda.addToRolePolicy(
      new PolicyStatement({
        resources: [props.accessHandler.getApiGateway().arnForExecuteApi()],
        actions: ["execute-api:Invoke"],
      })
    );

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
