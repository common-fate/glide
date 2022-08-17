import * as iam from "aws-cdk-lib/aws-iam";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as sfn from "aws-cdk-lib/aws-stepfunctions";
import { Construct } from "constructs";
import { Duration, Stack } from "aws-cdk-lib";
import * as path from "path";
import { EventBus } from "aws-cdk-lib/aws-events";

interface Props {
  eventBusSourceName: string;
  eventBus: EventBus;
  providerConfig: string;
}
export class Granter extends Construct {
  private _stateMachine: sfn.StateMachine;
  private _lambda: lambda.Function;
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);
    const code = lambda.Code.fromAsset(
      path.join(__dirname, "..", "..", "..", "..", "bin", "granter.zip")
    );

    this._lambda = new lambda.Function(this, "StepHandlerFunction", {
      code,
      timeout: Duration.minutes(5),
      environment: {
        EVENT_BUS_ARN: props.eventBus.eventBusArn,
        EVENT_BUS_SOURCE: props.eventBusSourceName,
        PROVIDER_CONFIG: props.providerConfig,
      },
      runtime: lambda.Runtime.GO_1_X,
      handler: "granter",
    });

    const definition = {
      StartAt: "Validate End is in the Future",
      States: {
        "Validate End is in the Future": {
          Type: "Choice",
          Choices: [
            {
              Variable: "$.grant.end",
              TimestampGreaterThanPath: "$$.State.EnteredTime",
              Next: "Wait for Grant Start Time",
            },
          ],
          Default: "Fail",
          Comment: "Do not provision any access if the end time is in the past",
        },
        "Wait for Grant Start Time": {
          Type: "Wait",
          TimestampPath: "$.grant.start",
          Next: "Activate Access",
        },
        "Activate Access": {
          Type: "Task",
          Resource: "arn:aws:states:::lambda:invoke",
          Parameters: {
            FunctionName: this._lambda.functionArn,
            Payload: {
              "action": "ACTIVATE",
              "grant.$": "$.grant",
            },
          },
          Retry: [
            {
              ErrorEquals: [
                "Lambda.ServiceException",
                "Lambda.AWSLambdaException",
                "Lambda.SdkClientException",
              ],
              IntervalSeconds: 2,
              MaxAttempts: 6,
              BackoffRate: 2,
            },
          ],
          Next: "Wait for Window End",
          ResultPath: "$",
          OutputPath: "$.Payload",
        },
        "Wait for Window End": {
          Type: "Wait",
          TimestampPath: "$.grant.end",
          Next: "Expire Access",
        },
        "Expire Access": {
          Type: "Task",
          Resource: "arn:aws:states:::lambda:invoke",
          OutputPath: "$.Payload",
          Parameters: {
            FunctionName: this._lambda.functionArn,
            Payload: {
              "action": "DEACTIVATE",
              "grant.$": "$.grant",
            },
          },
          Retry: [
            {
              ErrorEquals: [
                "Lambda.ServiceException",
                "Lambda.AWSLambdaException",
                "Lambda.SdkClientException",
              ],
              IntervalSeconds: 2,
              MaxAttempts: 6,
              BackoffRate: 2,
            },
          ],
          ResultPath: "$",
          End: true,
        },
        "Fail": {
          Type: "Fail",
        },
      },
      Comment: "Granted Access Handler State Machine",
    };

    this._stateMachine = new sfn.StateMachine(this, "StateMachine", {
      definition: new sfn.Pass(this, "StartState"),
    });

    const cfnStatemachine = this._stateMachine.node
      .defaultChild as sfn.CfnStateMachine;

    cfnStatemachine.definitionString = JSON.stringify(definition);

    const smRole = iam.Role.fromRoleArn(
      this,
      "StateMachineRole",
      cfnStatemachine.roleArn
    );
    this._lambda.grantInvoke(smRole);

    //grant read from ssm
    this._lambda.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ["ssm:GetParameter"],
        resources: [
          `arn:aws:ssm:${Stack.of(this).region}:${
            Stack.of(this).account
          }:parameter/granted/providers/*`,
        ],
      })
    );

    // add permissions for the AWS SSO provider
    this._lambda.addToRolePolicy(
      new iam.PolicyStatement({
        actions: [
          "sso:CreateAccountAssignment",
          "sso:DeleteAccountAssignment",
          "sso:ListAccountAssignments",
          "identitystore:ListUsers",
          "organizations:DescribeAccount",
        ],
        resources: ["*"],
      })
    );

    //permissions for the eks access handler when it is added
    //Ideally these should get set dynamically when the provider is added
    this._lambda.addToRolePolicy(
      new iam.PolicyStatement({
        actions: [
          "sso:CreatePermissionSet",
          "sso:PutInlinePolicyToPermissionSet",
          "sso:CreateAccountAssignment",
          "sso:ListPermissionSets",
          "sso:DescribePermissionSet",
          "sso:DeletePermissionSet",
          "iam:ListRoles",
          "sts:AssumeRole"
        ],
        resources: ["*"],

      })
    )
    this._lambda.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ["sts:AssumeRole"],
        resources: ["*"],
      })
    );

    props.eventBus.grantPutEventsTo(this._lambda);
  }
  getStateMachineARN(): string {
    return this._stateMachine.stateMachineArn;
  }
  getStateMachine(): sfn.StateMachine {
    return this._stateMachine;
  }
  getLogGroupName(): string {
    return this._lambda.logGroup.logGroupName;
  }
  getGranterARN(): string {
    return this._lambda.functionArn;
  }
  getGranterLambdaExecutionRoleARN(): string {
    return this._lambda.role?.roleArn ?? "";
  }
}
