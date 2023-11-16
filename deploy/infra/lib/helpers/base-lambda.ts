import * as cdk from "aws-cdk-lib";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";

export interface BaseLambdaFunProps {
  functionProps: lambda.FunctionProps;
  vpcConfig: any;
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
