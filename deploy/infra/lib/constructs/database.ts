import { Construct } from "constructs";
import * as dynamodb from "aws-cdk-lib/aws-dynamodb";
import * as cdk from "aws-cdk-lib";

interface Props {
  appName: string;
}

export class Database extends Construct {
  private readonly _appName: string;
  private _dynamoTable: dynamodb.Table;

  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);

    this._appName = props.appName;

    const dbTable = new dynamodb.Table(this, "DBTable", {
      tableName: this._appName,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      partitionKey: { name: "PK", type: dynamodb.AttributeType.STRING },
      sortKey: { name: "SK", type: dynamodb.AttributeType.STRING },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      pointInTimeRecovery: true,
    });

    const gsi1: dynamodb.GlobalSecondaryIndexProps = {
      indexName: "GSI1",
      partitionKey: {
        name: "GSI1PK",
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: "GSI1SK",
        type: dynamodb.AttributeType.STRING,
      },
    };
    const gsi2: dynamodb.GlobalSecondaryIndexProps = {
      indexName: "GSI2",
      partitionKey: {
        name: "GSI2PK",
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: "GSI2SK",
        type: dynamodb.AttributeType.STRING,
      },
    };
    const gsi3: dynamodb.GlobalSecondaryIndexProps = {
      indexName: "GSI3",
      partitionKey: {
        name: "GSI3PK",
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: "GSI3SK",
        type: dynamodb.AttributeType.STRING,
      },
    };

    const gsi4: dynamodb.GlobalSecondaryIndexProps = {
      indexName: "GSI4",
      partitionKey: {
        name: "GSI4PK",
        type: dynamodb.AttributeType.STRING,
      },
      sortKey: {
        name: "GSI4SK",
        type: dynamodb.AttributeType.STRING,
      },
    };

    dbTable.addGlobalSecondaryIndex(gsi1);
    dbTable.addGlobalSecondaryIndex(gsi2);
    dbTable.addGlobalSecondaryIndex(gsi3);
    dbTable.addGlobalSecondaryIndex(gsi4);

    const cfnTable = dbTable.node.defaultChild as dynamodb.CfnTable;

    // this table was previously part of the 'app-backend' construct.
    // to avoid recreating the table as part of a deployment, we force
    // the logical ID to be the same as the previous one.
    cfnTable.overrideLogicalId("APIDBTableA8FD77F9");

    this._dynamoTable = dbTable;
  }

  getTable(): dynamodb.Table {
    return this._dynamoTable;
  }
}
