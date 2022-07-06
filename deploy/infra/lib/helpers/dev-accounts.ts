export interface DevEnvironmentConfig {
  certificateArn: string;
  domainName: string;
  sharedCognitoUserPoolId: string;
  sharedCognitoUserPoolDomain: string;
  account: string;
}
