import { Duration, Stack } from "aws-cdk-lib";
import * as apigateway from "aws-cdk-lib/aws-apigateway";
import { EventBus } from "aws-cdk-lib/aws-events";
import * as iam from "aws-cdk-lib/aws-iam";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import * as path from "path";
import { Granter } from "./granter";
interface Props {
  appName: string;
  eventBusSourceName: string;
  eventBus: EventBus;
  /** A JSON payload of the access provider configuration. */
  providerConfig: string;
}

export class AccessHandler extends Construct {
  private _lambda: lambda.Function;
  private _apigateway: apigateway.RestApi;
  private _accessHandlerRole: iam.Role;
  private readonly _granter: Granter;
  private readonly _restApiName: string;
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);
    this._restApiName = props.appName + "-access-handler";

    const accessHandlerRolePolicy = new iam.PolicyDocument({
      statements: [
        new iam.PolicyStatement({
          actions: ["ssm:GetParameter"],
          resources: [
            `arn:aws:ssm:${Stack.of(this).region}:${
              Stack.of(this).account
            }:parameter/granted/providers/*`,
          ],
        }),
        new iam.PolicyStatement({
          actions: [
            "sso:ListPermissionSets",
            "sso:ListTagsForResource",
            "sso:DescribePermissionSet",
            "organizations:ListAccounts",
            "sso:DeleteAccountAssignment",
            "sso:ListAccountAssignments",
            "identitystore:ListUsers",
            "organizations:DescribeAccount",
          ],
          resources: ["*"],
        }),
        new PolicyStatement({
          actions: ["states:StopExecution"],
          resources: ["*"],
        }),
        new PolicyStatement({
          actions: ["sts:AssumeRole"],
          resources: ["*"],
        }),
      ],
    });

    this._accessHandlerRole = new iam.Role(this, "ExecutionRole", {
      assumedBy: new iam.CompositePrincipal(),
      description:
        "This role is assumed by the Granted Approvals access handler lambdas to grant and revoke access. It has permissions to assume any role depending on the equirements on individual providers.",
      inlinePolicies: {
        AccessHandlerPolicy: accessHandlerRolePolicy,
      },
    });

    this._granter = new Granter(this, "Granter", {
      eventBus: props.eventBus,
      eventBusSourceName: props.eventBusSourceName,
      providerConfig: props.providerConfig,
      assumeExecutionRoleArn: this._accessHandlerRole.roleArn,
    });

    const code = lambda.Code.fromAsset(
      path.join(__dirname, "..", "..", "..", "..", "bin", "access-handler.zip")
    );

    this._lambda = new lambda.Function(this, "RestAPIHandlerFunction", {
      code,
      timeout: Duration.seconds(60),
      environment: {
        GRANTED_RUNTIME: "lambda",
        STATE_MACHINE_ARN: this._granter.getStateMachineARN(),
        EVENT_BUS_ARN: props.eventBus.eventBusArn,
        EVENT_BUS_SOURCE: props.eventBusSourceName,
        PROVIDER_CONFIG: props.providerConfig,
        ASSUME_EXECUTION_ROLE_ARN: this._accessHandlerRole.roleArn,
      },
      runtime: lambda.Runtime.GO_1_X,
      handler: "access-handler",
    });

    this._accessHandlerRole.grantAssumeRole(
      new iam.ArnPrincipal(this._lambda.role?.roleArn ?? "")
    );
    this._accessHandlerRole.grantAssumeRole(
      new iam.ArnPrincipal(this._granter.getGranterLambdaExecutionRoleARN())
    );

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

    this._granter.getStateMachine().grantStartExecution(this._lambda);

    this._granter.getStateMachine().grantRead(this._lambda);

    props.eventBus.grantPutEventsTo(this._lambda);
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
  getAccessHandlerAssumeRoleArn(): string {
    return this._accessHandlerRole.roleArn;
  }
}
