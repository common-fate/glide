import { Duration } from "aws-cdk-lib";
import { Table } from "aws-cdk-lib/aws-dynamodb";
import { EventBus, Rule } from "aws-cdk-lib/aws-events";
import { LambdaFunction } from "aws-cdk-lib/aws-events-targets";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import { TargetGroupGranter } from "./targetgroup-granter";
import * as path from "path";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";

interface Props {
  eventBusSourceName: string;
  eventBus: EventBus;
  dynamoTable: Table;
  targetGroupGranter: TargetGroupGranter;
}
export class EventHandler extends Construct {
  private _sequentialLambda: lambda.Function;
  private _concurrentLambda: lambda.Function;
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);
    const code = lambda.Code.fromAsset(
      path.join(__dirname, "..", "..", "..", "..", "bin", "event-handler.zip")
    );

    /**
     * There are 2 event handler lambdas, a concurrent one and a sequential one.
     * The idea here is that the (currently) 3 event types which can only be processed one at a time.
     * initiating revoke, initiating cancel, review request
     *
     */
    this._concurrentLambda = new lambda.Function(this, "ConcurrentFunction", {
      code,
      timeout: Duration.seconds(20),
      environment: {
        COMMONFATE_TABLE_NAME: props.dynamoTable.tableName,
        COMMONFATE_EVENT_BUS_ARN: props.eventBus.eventBusArn,
        COMMONFATE_GRANTER_V2_STATE_MACHINE_ARN:
          props.targetGroupGranter.getStateMachineARN(),
      },
      runtime: lambda.Runtime.GO_1_X,
      handler: "event-handler",
    });
    this._sequentialLambda = new lambda.Function(this, "SequentialFunction", {
      code,
      timeout: Duration.seconds(20),
      environment: {
        COMMONFATE_TABLE_NAME: props.dynamoTable.tableName,
        COMMONFATE_EVENT_BUS_ARN: props.eventBus.eventBusArn,
        COMMONFATE_GRANTER_V2_STATE_MACHINE_ARN:
          props.targetGroupGranter.getStateMachineARN(),
      },
      runtime: lambda.Runtime.GO_1_X,
      handler: "event-handler",
      // Concurrency set to 1
      reservedConcurrentExecutions: 1,
    });

    this._concurrentLambda.addToRolePolicy(
      new PolicyStatement({
        actions: [
          "states:StopExecution",
          "states:StartExecution",
          "states:DescribeExecution",
          "states:GetExecutionHistory",
        ],
        resources: [props.targetGroupGranter.getStateMachineARN()],
      })
    );

    props.eventBus.grantPutEventsTo(this._concurrentLambda);
    this._sequentialLambda.addToRolePolicy(
      new PolicyStatement({
        actions: [
          "states:StopExecution",
          "states:StartExecution",
          "states:DescribeExecution",
          "states:GetExecutionHistory",
        ],
        resources: [props.targetGroupGranter.getStateMachineARN()],
      })
    );

    props.eventBus.grantPutEventsTo(this._sequentialLambda);
    new Rule(this, "SequentialEventBusRule", {
      eventBus: props.eventBus,
      eventPattern: {
        source: [props.eventBusSourceName],
        detailType: [
          "request.revoke.initiated",
          "request.cancel.initiated",
          "accessGroup.review",
        ],
      },
      targets: [
        new LambdaFunction(this._sequentialLambda, {
          retryAttempts: 2,
        }),
      ],
    });
    props.dynamoTable.grantReadWriteData(this._sequentialLambda);

    new Rule(this, "ConcurrentEventBusRule", {
      eventBus: props.eventBus,
      eventPattern: {
        source: [props.eventBusSourceName],
        detail: {
          detailType: [
            {
              "anything-but": [
                "request.revoke.initiated",
                "request.cancel.initiated",
                "accessGroup.review",
              ],
            },
          ],
        },
      },
      targets: [
        new LambdaFunction(this._concurrentLambda, {
          retryAttempts: 2,
        }),
      ],
    });
    props.dynamoTable.grantReadWriteData(this._concurrentLambda);
  }
  getConcurrentLogGroupName(): string {
    return this._concurrentLambda.logGroup.logGroupName;
  }
  getSequentialLogGroupName(): string {
    return this._sequentialLambda.logGroup.logGroupName;
  }
}
