import { Duration, Stack } from "aws-cdk-lib";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import { EventBus } from "aws-cdk-lib/aws-events";
import * as iam from "aws-cdk-lib/aws-iam";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import * as path from "path";
import { Granter } from "./granter";
import { BaseLambdaFunction } from "../helpers/base-lambda";
interface Props {
  appName: string;
  eventBusSourceName: string;
  eventBus: EventBus;
  remoteConfigUrl: string;
  remoteConfigHeaders: string;
  /** A JSON payload of the access provider configuration. */
  providerConfig: string;
  vpcConfig: any;
}

export class AccessHandler extends Construct {
  private _lambda: lambda.Function;
  private _apigateway: apigateway.RestApi;
  private readonly _granter: Granter;
  private readonly _restApiName: string;
  private _executionRole: iam.Role;
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);
    this._restApiName = props.appName + "-access-handler";

    // Create the access handler role with the lambda execution role ARNs as principals
    this._executionRole = new iam.Role(this, "ExecutionRole", {
      assumedBy: new iam.ServicePrincipal("lambda.amazonaws.com"),
      description:
        "This role is assumed by the Common Fate access handler lambdas to grant and revoke access. It has permissions to assume any role depending on the equirements on individual providers.",
      // https://docs.aws.amazon.com/cdk/api/v2/docs/aws-cdk-lib.aws_lambda.Function.html#role
      managedPolicies: [
        iam.ManagedPolicy.fromAwsManagedPolicyName(
          "service-role/AWSLambdaVPCAccessExecutionRole"
        ),
        iam.ManagedPolicy.fromAwsManagedPolicyName(
          "service-role/AWSLambdaBasicExecutionRole"
        ),
      ],

      inlinePolicies: {
        AccessHandlerPolicy: new iam.PolicyDocument({
          statements: [
            new iam.PolicyStatement({
              actions: ["ssm:GetParameter"],
              resources: [
                `arn:aws:ssm:${Stack.of(this).region}:${
                  Stack.of(this).account
                }:parameter/granted/providers/*`,
              ],
            }),
            new PolicyStatement({
              actions: ["states:StopExecution"],
              resources: ["*"],
            }),
            new PolicyStatement({
              actions: ["sts:AssumeRole"],
              resources: ["*"],
            }),
            new PolicyStatement({
              actions: [
                "states:StartExecution",
                "states:DescribeExecution",
                "states:GetExecutionHistory",
                "states:StopExecution",
              ],
              resources: ["*"],
            }),
          ],
        }),
      },
    });

    props.eventBus.grantPutEventsTo(this._executionRole);

    this._granter = new Granter(this, "Granter", {
      eventBus: props.eventBus,
      eventBusSourceName: props.eventBusSourceName,
      providerConfig: props.providerConfig,
      executionRole: this._executionRole,
      remoteConfigUrl: props.remoteConfigUrl,
      remoteConfigHeaders: props.remoteConfigHeaders,
      vpcConfig: props.vpcConfig,
    });

    const code = lambda.Code.fromAsset(
      path.join(__dirname, "..", "..", "..", "..", "bin", "access-handler.zip")
    );

    this._lambda = new BaseLambdaFunction(this, "RestAPIHandlerFunction", {
      functionProps: {
        code,
        timeout: Duration.seconds(60),
        environment: {
          COMMONFATE_ACCESS_HANDLER_RUNTIME: "lambda",
          COMMONFATE_STATE_MACHINE_ARN: this._granter.getStateMachineARN(),
          COMMONFATE_EVENT_BUS_ARN: props.eventBus.eventBusArn,
          COMMONFATE_EVENT_BUS_SOURCE: props.eventBusSourceName,
          COMMONFATE_PROVIDER_CONFIG: props.providerConfig,
          COMMONFATE_ACCESS_REMOTE_CONFIG_URL: props.remoteConfigUrl,
          COMMONFATE_REMOTE_CONFIG_HEADERS: props.remoteConfigHeaders,
        },
        runtime: lambda.Runtime.GO_1_X,
        handler: "access-handler",
        role: this._executionRole,
      },
      vpcConfig: props.vpcConfig,
    });

    this._apigateway = new apigateway.RestApi(this, "RestAPI", {
      restApiName: this._restApiName,
    });

    const lambdaProxy = this._apigateway.root.addResource("{proxy+}");
    lambdaProxy.addMethod(
      "ANY",
      new apigateway.LambdaIntegration(this._lambda, {
        allowTestInvoke: false,
      }),
      { authorizationType: apigateway.AuthorizationType.IAM }
    );
  }
  getGranter(): Granter {
    return this._granter;
  }
  getApiUrl(): string {
    return this._apigateway.url;
  }
  getApiGateway(): apigateway.RestApi {
    return this._apigateway;
  }
  getLogGroupName(): string {
    return this._lambda.logGroup.logGroupName;
  }
  getAccessHandlerARN(): string {
    return this._lambda.functionArn;
  }
  getAccessHandlerExecutionRoleArn(): string {
    return this._executionRole.roleArn;
  }
}
