import * as lambda from "aws-cdk-lib/aws-lambda";
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
