import { CfnOutput } from "aws-cdk-lib";
import { Construct } from "constructs";

export type StackOutputs = {
  CognitoClientID: string;
  CloudFrontDomain: string;
  FrontendDomainOutput: string;
  CloudFrontDistributionID: string;
  S3BucketName: string;
  UserPoolID: string;
  UserPoolDomain: string;
  APIURL: string;
  WebhookURL: string;
  APILogGroupName: string;
  WebhookLogGroupName: string;
  IDPSyncLogGroupName: string;
  AccessHandlerLogGroupName: string;
  EventBusLogGroupName: string;
  EventsHandlerLogGroupName: string;
  GranterLogGroupName: string;
  SlackNotifierLogGroupName: string;
  DynamoDBTable: string;
  GranterStateMachineArn: string;
  EventBusArn: string;
  EventBusSource: string;
  IdpSyncFunctionName: string;
  Region: string;
  PaginationKMSKeyARN: string;
  AccessHandlerExecutionRoleARN: string;
};
/**
 * generateOutputs creates a Cloudformation Output for each key-value pair in the type StackOutputs
 *
 */
export const generateOutputs = (scope: Construct, o: StackOutputs) => {
  Object.entries(o).forEach(
    ([k, v]) =>
      new CfnOutput(scope, k, {
        value: v,
      })
  );
};
