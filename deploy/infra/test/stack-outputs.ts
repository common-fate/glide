import { StackOutputs } from "../lib/helpers/outputs";
import { writeFileSync } from "fs";

// testOutputs will have a type error when new fields are added to stack outputs
// a new entry should be added here as this is used in a test to ensure consistency with the corresponding go type
// in pkg/deploy/output.go
const testOutputs: StackOutputs = {
  CognitoClientID: "abcdefg",
  CloudFrontDomain: "abcdefg",
  FrontendDomainOutput: "abcdefg",
  CloudFrontDistributionID: "abcdefg",
  S3BucketName: "abcdefg",
  UserPoolID: "abcdefg",
  UserPoolDomain: "abcdefg",
  APIURL: "abcdefg",
  WebhookURL: "abcdefg",
  APILogGroupName: "abcdefg",
  WebhookLogGroupName: "abcdefg",
  IDPSyncLogGroupName: "abcdefg",
  AccessHandlerLogGroupName: "abcdefg",
  EventBusLogGroupName: "abcdefg",
  EventsHandlerLogGroupName: "abcdefg",
  GranterLogGroupName: "abcdefg",
  SlackNotifierLogGroupName: "abcdefg",
  DynamoDBTable: "abcdefg",
  GranterStateMachineArn: "abcdefg",
  EventBusArn: "abcdefg",
  EventBusSource: "abcdefg",
  IdpSyncFunctionName: "abcdefg",
  Region: "abcdefg",
  PaginationKMSKeyARN: "abcdefg",
  AccessHandlerExecutionRoleARN: "abcdefg",
};

// Write the json object to ./testOutputs.json so that it can be parsed by a go test in pkg/deploy.output_test.go
writeFileSync("./testOutputs.json", JSON.stringify(testOutputs));
