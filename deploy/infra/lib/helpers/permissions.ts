import { Duration } from "aws-cdk-lib";
import { Table } from "aws-cdk-lib/aws-dynamodb";
import * as events from "aws-cdk-lib/aws-events";
import * as targets from "aws-cdk-lib/aws-events-targets";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Construct } from "constructs";
import * as path from "path";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";

export const grantInvokeCommunityProviders = (_lambda: lambda.Function) => {
  _lambda.addToRolePolicy(
    new PolicyStatement({
      resources: ["arn:aws:lambda:#:#:#"],
      actions: ["lambda:InvokeFunction"],
      conditions: {
        StringEquals: {
          "iam:ResourceTag/common-fate-abac-role": "access-provider",
        },
      },
    })
  );
};
