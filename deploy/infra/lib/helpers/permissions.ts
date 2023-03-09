import * as lambda from "aws-cdk-lib/aws-lambda";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";

/**
 * grants a lambda function permissions to assume the invocation role for a handler
 * @param _lambda
 */
export const grantAssumeHandlerRole = (_lambda: lambda.Function) => {
  _lambda.addToRolePolicy(
    new PolicyStatement({
      resources: ["*"],
      actions: ["sts:AssumeRole"],
      conditions: {
        StringEquals: {
          "iam:ResourceTag/common-fate-abac-role": "handler-invoke",
        },
      },
    })
  );
};

/**
 * grants a lambda function permissions to assume an aws identity sync role
 * Used when AWS SSO is used for SAML SSO
 * @param _lambda
 */
export const grantAssumeIdentitySyncRole = (_lambda: lambda.Function) => {
  _lambda.addToRolePolicy(
    new PolicyStatement({
      resources: ["*"],
      actions: ["sts:AssumeRole"],
      conditions: {
        StringEquals: {
          "iam:ResourceTag/common-fate-abac-role": "aws-sso-identity-provider",
        },
      },
    })
  );
};
