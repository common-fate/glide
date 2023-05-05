import { CfnOutput } from "aws-cdk-lib";
import { Construct } from "constructs";

export type StackOutputs = {
  APILogGroupName: string;
  APIURL: string;
  CacheSyncFunctionName: string;
  CacheSyncLogGroupName: string;
  CLIAppClientID: string;
  CloudFrontDistributionID: string;
  CloudFrontDomain: string;
  CognitoClientID: string;
  DynamoDBTable: string;
  EventBusArn: string;
  EventBusLogGroupName: string;
  EventBusSource: string;
  EventsHandlerConcurrentLogGroupName: string;
  EventsHandlerSequentialLogGroupName: string;
  FrontendDomainOutput: string;
  GovernanceURL: string;
  GranterV2StateMachineArn: string;
  GranterLogGroupName: string;
  HealthcheckFunctionName: string;
  HealthcheckLogGroupName: string;
  IDPSyncExecutionRoleARN: string;
  IDPSyncFunctionName: string;
  IDPSyncLogGroupName: string;
  PaginationKMSKeyARN: string;
  Region: string;
  RestAPIExecutionRoleARN: string;
  S3BucketName: string;
  SAMLIdentityProviderName: string;
  SlackNotifierLogGroupName: string;
  UserPoolDomain: string;
  UserPoolID: string;
  WebhookLogGroupName: string;
  WebhookURL: string;
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
