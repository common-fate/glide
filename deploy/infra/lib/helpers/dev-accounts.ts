export interface DevEnvironmentConfig {
  certificateArn: string;
  domainName: string;
  sharedCognitoUserPoolId: string;
  sharedCognitoUserPoolDomain: string;
}

/**
 * DevEnvironments are Common Fate test accounts that our CI pipelines
 * create branch deployments in.
 */
export const DevEnvironments = new Map<string, DevEnvironmentConfig>([
  [
    "test",
    {
      certificateArn:
        "arn:aws:acm:us-east-1:963589028267:certificate/55db5633-5e2e-4b26-b37f-d07248997bfd",
      domainName: "test.granted.run",
      sharedCognitoUserPoolDomain: "granted-test",
      sharedCognitoUserPoolId: "us-east-1_Q0RLaRTOY",
    },
  ],
]);
