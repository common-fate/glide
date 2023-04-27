import { Duration } from "aws-cdk-lib";
import { Table } from "aws-cdk-lib/aws-dynamodb";
import { EventBus, Rule } from "aws-cdk-lib/aws-events";
import { LambdaFunction } from "aws-cdk-lib/aws-events-targets";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import { TargetGroupGranter } from "./targetgroup-granter";
import * as path from "path";

interface Props {
  eventBusSourceName: string;
  eventBus: EventBus;
  dynamoTable: Table;
  targetGroupGranter: TargetGroupGranter;
}
export class EventHandler extends Construct {
  private _lambda: lambda.Function;
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);
    const code = lambda.Code.fromAsset(
      path.join(__dirname, "..", "..", "..", "..", "bin", "event-handler.zip")
    );
    this._lambda = new lambda.Function(this, "Function", {
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

    new Rule(this, "EventBusRule", {
      eventBus: props.eventBus,
      eventPattern: { source: [props.eventBusSourceName] },
      targets: [
        new LambdaFunction(this._lambda, {
          retryAttempts: 2,
        }),
      ],
    });
    props.dynamoTable.grantReadWriteData(this._lambda);
  }
  getLogGroupName(): string {
    return this._lambda.logGroup.logGroupName;
  }
}
