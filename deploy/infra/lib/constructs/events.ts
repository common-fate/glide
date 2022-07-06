import { EventBus as EB, Rule } from "aws-cdk-lib/aws-events";
import { CloudWatchLogGroup } from "aws-cdk-lib/aws-events-targets";
import { LogGroup } from "aws-cdk-lib/aws-logs";
import { Construct } from "constructs";
interface Props {
  appName: string;
}

export class EventBus extends Construct {
  private _eventBus: EB;
  private _sourceName: string;
  private _logGroup: LogGroup;
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);
    this._sourceName = "commonfate.io/granted";

    this._eventBus = new EB(this, "EventBus", {
      eventBusName: props.appName,
    });

    // write all events to CloudWatch for debugging.
    this._logGroup = new LogGroup(this, "EventBusLog");

    new Rule(this, "EventBusCloudWatchRule", {
      eventBus: this._eventBus,
      eventPattern: { source: [this._sourceName] },
      targets: [new CloudWatchLogGroup(this._logGroup)],
    });
  }
  getEventBus(): EB {
    return this._eventBus;
  }

  getEventBusSourceName(): string {
    return this._sourceName;
  }

  getLogGroupName(): string {
    return this._logGroup.logGroupName;
  }
}
