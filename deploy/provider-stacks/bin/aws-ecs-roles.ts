import { App, CfnOutput, Stack } from "aws-cdk-lib";
import { Construct } from "constructs";
import * as iam from "aws-cdk-lib/aws-iam";
class AWSSSOECSRoleStack extends Stack {
  constructor(scope: Construct, id: string, props: any) {
    super(scope, id, props);
    const role = new iam.Role(this, "GrantedAccessHandlerECSFlask", {
      assumedBy: new iam.ArnPrincipal(
        "arn:aws:iam::12345678912:role/granted-approvals"
      ),
      description:
        "This role grants management access to AWS SSO and read access to ECS for the Granted Access Handler.",
      inlinePolicies: {
        AccessHandlerSSOPolicy: new iam.PolicyDocument({
          statements: [
            new iam.PolicyStatement({
              sid: "ReadSSO",
              actions: [
                "sso:DescribeAccountAssignmentDeletionStatus",
                "sso:DescribeAccountAssignmentCreationStatus",
                "sso:DescribePermissionSet",
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
                "iam:ListRoles",
                "iam:ListUsers",
              ],
              resources: ["*"],
            }),
            new iam.PolicyStatement({
              sid: "AssignSSO",
              actions: [
                "sso:DeletePermissionSet",
                "sso:DeleteAccountAssignment",
                "sso:CreatePermissionSet",
                "sso:PutInlinePolicyToPermissionSet",
                "sso:CreateAccountAssignment",
              ],
              resources: ["*"],
            }),
            new iam.PolicyStatement({
              sid: "ReadECS",
              actions: [
                "ecs:ListTasks",
                "ecs:DescribeTasks",
                "cloudtrail:LookupEvents",
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
new AWSSSOECSRoleStack(app, "AWSECSFlaskRoleStack", {});
