import { Duration, Stack } from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import * as path from "path";
import { Effect, PolicyStatement } from "aws-cdk-lib/aws-iam";
import * as events from "aws-cdk-lib/aws-events";
import * as targets from "aws-cdk-lib/aws-events-targets";
import { Table } from "aws-cdk-lib/aws-dynamodb";
import { WebUserPool } from "./app-user-pool";

interface Props {
  dynamoTable: Table;
  userPool: WebUserPool;
  identityProviderSyncConfiguration: string;
  analyticsDisabled: string;
  analyticsUrl: string;
  analyticsLogLevel: string;
  analyticsDeploymentStage: string;
}

export class IdpSync extends Construct {
  private _lambda: lambda.Function;
  private eventRule: events.Rule;

  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);
    const code = lambda.Code.fromAsset(
      path.join(__dirname, "..", "..", "..", "..", "bin", "syncer.zip")
    );

    this._lambda = new lambda.Function(this, "HandlerFunction", {
      code,
      timeout: Duration.seconds(20),
      environment: {
        APPROVALS_TABLE_NAME: props.dynamoTable.tableName,
        IDENTITY_PROVIDER: props.userPool.getIdpType(),
        APPROVALS_COGNITO_USER_POOL_ID: props.userPool.getUserPoolId(),
        IDENTITY_SETTINGS: props.identityProviderSyncConfiguration,
        CF_ANALYTICS_DISABLED: props.analyticsDisabled,
        CF_ANALYTICS_URL: props.analyticsUrl,
        CF_ANALYTICS_LOG_LEVEL: props.analyticsLogLevel,
        CF_ANALYTICS_DEPLOYMENT_STAGE: props.analyticsDeploymentStage,
      },
      runtime: lambda.Runtime.GO_1_X,
      handler: "syncer",
    });

    props.dynamoTable.grantReadWriteData(this._lambda);

    //add event bridge trigger to lambda
    this.eventRule = new events.Rule(this, "EventBridgeCronRule", {
      schedule: events.Schedule.cron({ minute: "0/5" }),
    });

    // add the Lambda function as a target for the Event Rule
    this.eventRule.addTarget(new targets.LambdaFunction(this._lambda));

    // allow the Event Rule to invoke the Lambda function
    targets.addLambdaPermission(this.eventRule, this._lambda);

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
        ],
      })
    );
    this._lambda.addToRolePolicy(
      new PolicyStatement({
        actions: ["ssm:GetParameter"],
        resources: [
          `arn:aws:ssm:${Stack.of(this).region}:${
            Stack.of(this).account
          }:parameter/granted/secrets/identity/*`,
        ],
      })
    );

    this._lambda.addToRolePolicy(
      new PolicyStatement({
        actions: ["sts:AssumeRole"],
        resources: ["*"],
        conditions: {
          StringEquals: {
            "iam:ResourceTag/common-fate-abac-role":
              "aws-sso-identity-provider",
          },
        },
      })
    );
    //allow the lambda to write to the table
    props.dynamoTable.grantWriteData(this._lambda);
  }
  getLogGroupName(): string {
    return this._lambda.logGroup.logGroupName;
  }
  getFunctionName(): string {
    return this._lambda.functionName;
  }
  getExecutionRoleArn(): string {
    return this._lambda.role?.roleArn || "";
  }
}
