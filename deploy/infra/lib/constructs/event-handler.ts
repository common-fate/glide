import { Duration } from "aws-cdk-lib";
import { Table } from "aws-cdk-lib/aws-dynamodb";
import { EventBus, Rule } from "aws-cdk-lib/aws-events";
import { LambdaFunction } from "aws-cdk-lib/aws-events-targets";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import * as path from "path";
import {BaseLambdaFunction} from "../helpers/base-lambda";

interface Props {
  eventBusSourceName: string;
  eventBus: EventBus;
  dynamoTable: Table;
  vpcConfig: any;
}
export class EventHandler extends Construct {
  private _lambda: lambda.Function;
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);
    const code = lambda.Code.fromAsset(
      path.join(__dirname, "..", "..", "..", "..", "bin", "event-handler.zip")
    );
    this._lambda = new BaseLambdaFunction(this, "Function", {
      functionProps: {
        code,
        timeout: Duration.seconds(20),
        environment: {
          COMMONFATE_TABLE_NAME: props.dynamoTable.tableName,
        },
        runtime: lambda.Runtime.GO_1_X,
        handler: "event-handler",
      },
      vpcConfig: props.vpcConfig,
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
