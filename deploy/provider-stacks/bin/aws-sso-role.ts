import { App, CfnOutput, Stack } from "aws-cdk-lib";
import { Construct } from "constructs";
import * as iam from "aws-cdk-lib/aws-iam";
class AWSSSORoleStack extends Stack {
  constructor(scope: Construct, id: string, props: any) {
    super(scope, id, props);
    const role = new iam.Role(this, "GrantedAccessHandlerSSO", {
      assumedBy: new iam.ArnPrincipal(
        "arn:aws:iam::12345678912:role/granted-approvals"
      ),
      description:
        "This role grants management access to AWS SSO for the Granted Access Handler.",
      inlinePolicies: {
        AccessHandlerSSOPolicy: new iam.PolicyDocument({
          statements: [
            new iam.PolicyStatement({
              actions: [
                "sso:CreateAccountAssignment",
                "sso:DescribeAccountAssignmentDeletionStatus",
                "sso:DescribeAccountAssignmentCreationStatus",
                "sso:DescribePermissionSet",
                "sso:DeleteAccountAssignment",
                "sso:ListPermissionSets",
                "sso:ListTagsForResource",
                "sso:ListAccountAssignments",
                "organizations:ListAccounts",
                "organizations:DescribeAccount",
                "organizations:DescribeOrganization",
                "iam:GetSAMLProvider",
                "iam:GetRole",
                "iam:ListAttachedRolePolicies",
                "iam:ListRolePolicies",
                "identitystore:ListUsers",
              ],
              resources: ["*"],
            }),
          ],
        }),
      },
    });
    new CfnOutput(this, "RoleARN", {
      value: role.roleArn,
    });
  }
}

const app = new App();
new AWSSSORoleStack(app, "AWSSSORoleStack", {});
