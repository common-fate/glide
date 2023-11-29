import * as cdk from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";

export type VpcConfig = {
  // The intended data type for the SubnetIds and SecurityGroupIds fields is a list of strings.
  // However, in the production stack, the type of these fields is ICfnRuleConditionExpression, while in the development stack, it is string[].
  // To accommodate this difference, we are currently using 'any' data type to represent these fields.
  SubnetIds: any;
  SecurityGroupIds: any;
};

export interface BaseLambdaFunProps {
  functionProps: lambda.FunctionProps;
  vpcConfig: VpcConfig;
}

export class BaseLambdaFunction extends lambda.Function {
  constructor(scope: any, id: string, props: BaseLambdaFunProps) {
    super(scope, id, props.functionProps);
    this.addToRolePolicy(
      new PolicyStatement({
        actions: [
          "ec2:DescribeNetworkInterfaces",
          "ec2:CreateNetworkInterface",
          "ec2:DeleteNetworkInterface",
          "ec2:DescribeInstances",
          "ec2:AttachNetworkInterface",
        ],
        resources: ["*"],
      })
    );
    const cfnDist = this.node.defaultChild as cdk.CfnResource;
    cfnDist.addPropertyOverride("VpcConfig", props.vpcConfig);
  }
}
