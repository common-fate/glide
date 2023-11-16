import { CustomResource, Duration, Stack } from "aws-cdk-lib";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { Bucket } from "aws-cdk-lib/aws-s3";
import { Construct } from "constructs";
import * as path from "path";
import {BaseLambdaFunction} from "../helpers/base-lambda";

interface Props {
  cfReleaseBucket: string;
  cfReleaseBucketFrontendAssetObjectPrefix: string;
  frontendBucket: Bucket;
  cognitoUserPoolId: string;
  cognitoClientId: string;
  cliClientId: string;
  userPoolDomain: string;
  frontendDomain: string;
  cloudfrontDistributionId: string;
  apiUrl: string;
  vpcConfig: any;
}
export class ProductionFrontendDeployer extends Construct {
  private _lambda: lambda.Function;
  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);
    const code = lambda.Code.fromAsset(
      path.join(
        __dirname,
        "..",
        "..",
        "..",
        "..",
        "bin",
        "frontend-deployer.zip"
      )
    );
    this._lambda = new BaseLambdaFunction(this, "Function", {
      functionProps: {
        code,
        // The frontend deployer has a 7 minute timeout
        // internally, the deployer has a 5 minute retry backoff around the invalidation cloudfront method
        // at worst execution would take around 5 mins 30s
        timeout: Duration.seconds(60 * 7),
        environment: {
          CF_RELEASES_BUCKET: props.cfReleaseBucket,
          CF_RELEASES_FRONTEND_ASSET_OBJECT_PREFIX:
            props.cfReleaseBucketFrontendAssetObjectPrefix,
          COMMONFATE_FRONTEND_BUCKET: props.frontendBucket.bucketName,
          COMMONFATE_COGNITO_USER_POOL_ID: props.cognitoUserPoolId,
          COMMONFATE_COGNITO_CLIENT_ID: props.cognitoClientId,
          COMMONFATE_USER_POOL_DOMAIN: props.userPoolDomain,
          COMMONFATE_FRONTEND_DOMAIN: props.frontendDomain,
          COMMONFATE_CLOUDFRONT_DISTRIBUTION_ID: props.cloudfrontDistributionId,
          COMMONFATE_CLI_CLIENT_ID: props.cliClientId,
          COMMONFATE_API_URL: props.apiUrl,
        },
        runtime: lambda.Runtime.GO_1_X,
        handler: "frontend-deployer",
      },
      vpcConfig: props.vpcConfig
    });

    // Allow the deployer to deploy to the frontend bucket
    props.frontendBucket.grantReadWrite(this._lambda);

    // Allow the deployer access to read objects from the releases buckets
    this._lambda.addToRolePolicy(
      new PolicyStatement({
        actions: ["s3:ListBucket", "s3:GetObject", "s3:GetObjectVersion"],
        resources: [`arn:aws:s3:::${props.cfReleaseBucket}/*`],
      })
    );
    this._lambda.addToRolePolicy(
      new PolicyStatement({
        actions: ["s3:ListBucket"],
        resources: [`arn:aws:s3:::${props.cfReleaseBucket}`],
      })
    );

    // Allow the deployer to invalidation the distribution cache after updating the files
    this._lambda.addToRolePolicy(
      new PolicyStatement({
        actions: ["cloudfront:CreateInvalidation"],
        resources: ["*"],
      })
    );

    // custom resource will deploy the frontend from the public releases bucket
    new CustomResource(this, "CustomResource", {
      serviceToken: this._lambda.functionArn,

      // These properties will cause the custom resource to run during an update when the cognito client id changes or when the frontend assets path changes
      properties: {
        Release: props.cfReleaseBucketFrontendAssetObjectPrefix,
        CognitoClientID: props.cognitoClientId,
        FrontendDomain: props.frontendDomain,
      },
    });
  }
}
