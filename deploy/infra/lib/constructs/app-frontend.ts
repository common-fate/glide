import * as cdk from "aws-cdk-lib";
import { CfnCondition } from "aws-cdk-lib";
import { Certificate, ICertificate } from "aws-cdk-lib/aws-certificatemanager";
import * as cloudfront from "aws-cdk-lib/aws-cloudfront";
import { CfnDistribution } from "aws-cdk-lib/aws-cloudfront";
import * as origins from "aws-cdk-lib/aws-cloudfront-origins";
import * as s3 from "aws-cdk-lib/aws-s3";
import { Bucket } from "aws-cdk-lib/aws-s3";
import { Construct } from "constructs";
import { DevEnvironmentConfig } from "../helpers/dev-accounts";

enum HttpStatus {
  OK = 200,
  Unauthorized = 403,
  NotFound = 404,
}

interface AppDomain {
  // e.g. `granted.devcommonfate.io`
  domainName: string;
  // the AWS ACM certificate associated with the domain. Must be in us-east-1.
  certificate: ICertificate;
}

interface Props {
  appName: string;
  // The stack name should be stable at synth time and not include token
  // for dev this should be app name, for prod this should be stack name
  stableName: string;
}

interface ProdCdnConfig {
  wafAclArn?: string;
  frontendDomain: string;
  frontendCertificateArn: string;
}

interface ProdDomain {
  customFrontendCondition: cdk.CfnCondition;
  domainName: string;
  certificateArn: string;
}

export class AppFrontend extends Construct {
  private _frontendS3Bucket: s3.Bucket;
  private _domain?: AppDomain;
  private _distribution: cloudfront.Distribution;
  private _prodDomain?: ProdDomain;

  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);
  }
  private withFrontendBucket(logBucket?: Bucket) {
    /* S3 bucket for react app CDN */
    this._frontendS3Bucket = new s3.Bucket(this, "WebAppBucket", {
      accessControl: s3.BucketAccessControl.PRIVATE,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
      encryption: s3.BucketEncryption.S3_MANAGED,
      enforceSSL: true,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      autoDeleteObjects: true,
      serverAccessLogsBucket: logBucket && logBucket,
      serverAccessLogsPrefix: logBucket && "webApp/",
    });
  }

  //distributionConfig requires that withFrontendBucket has run which should be called by withProdCDN or withDevCDN
  private distributionConfig(opts: {
    logBucket?: Bucket;
    wafAclArn?: string;
  }): cdk.aws_cloudfront.DistributionProps {
    const defaultErrorResponseTTLSeconds = 10;
    return {
      webAclId: opts.wafAclArn,
      domainNames: this._domain ? [this._domain.domainName] : undefined,
      certificate: this._domain ? this._domain.certificate : undefined,
      defaultBehavior: {
        origin: new origins.S3Origin(this._frontendS3Bucket),
        allowedMethods: cloudfront.AllowedMethods.ALLOW_GET_HEAD,
        viewerProtocolPolicy: cloudfront.ViewerProtocolPolicy.REDIRECT_TO_HTTPS,
        cachedMethods: cloudfront.CachedMethods.CACHE_GET_HEAD,
        compress: false,
        cachePolicy:
          cloudfront.CachePolicy.CACHING_OPTIMIZED_FOR_UNCOMPRESSED_OBJECTS,
      },
      httpVersion: cloudfront.HttpVersion.HTTP1_1,
      enableIpv6: false,
      defaultRootObject: "/index.html",
      priceClass: cloudfront.PriceClass.PRICE_CLASS_100,
      errorResponses: [
        {
          httpStatus: HttpStatus.NotFound,
          responseHttpStatus: HttpStatus.OK,
          responsePagePath: "/index.html",
          ttl: cdk.Duration.seconds(defaultErrorResponseTTLSeconds),
        },
        {
          httpStatus: HttpStatus.Unauthorized,
          responseHttpStatus: HttpStatus.OK,
          responsePagePath: "/index.html",
          ttl: cdk.Duration.seconds(defaultErrorResponseTTLSeconds),
        },
      ],
      minimumProtocolVersion: cloudfront.SecurityPolicyProtocol.TLS_V1_2_2021,
      logBucket: opts.logBucket,
      logFilePrefix: opts.logBucket && "distribution/",
    };
  }

  // withProdCDN defines a cdn with access logging and optionally predefined urls
  withProdCDN(config: ProdCdnConfig): AppFrontend {
    /* CDN */

    const accessLogBucket = new s3.Bucket(this, "AccessLogBucket", {
      accessControl: s3.BucketAccessControl.LOG_DELIVERY_WRITE,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
      encryption: s3.BucketEncryption.S3_MANAGED,
      enforceSSL: true,
      serverAccessLogsPrefix: "thisBucket/",
    });

    const customFrontendCondition = new CfnCondition(
      this,
      "CustomFrontendUrlCondition",
      {
        expression: cdk.Fn.conditionNot(
          cdk.Fn.conditionEquals(config.frontendDomain, "")
        ),
      }
    );

    this.withFrontendBucket(accessLogBucket);

    this._distribution = new cloudfront.Distribution(
      this,
      "CloudfrontDistribution",
      this.distributionConfig({
        logBucket: accessLogBucket,
        wafAclArn: config.wafAclArn,
      })
    );

    const cfnDist = this._distribution.node.defaultChild as CfnDistribution;
    // https://docs.aws.amazon.com/cdk/v2/guide/cfn_layer.html#cfn_layer_raw
    cfnDist.addPropertyOverride("DistributionConfig.Aliases", [
      cdk.Fn.conditionIf(
        customFrontendCondition.logicalId,
        config.frontendDomain,
        cdk.Fn.ref("AWS::NoValue")
      ),
    ]);
    cfnDist.addPropertyOverride(
      "DistributionConfig.ViewerCertificate",
      cdk.Fn.conditionIf(
        customFrontendCondition.logicalId,
        {
          AcmCertificateArn: config.frontendCertificateArn,
          SslSupportMethod: "sni-only",
          MinimumProtocolVersion:
            cloudfront.SecurityPolicyProtocol.TLS_V1_2_2021,
        },
        cdk.Fn.ref("AWS::NoValue")
      )
    );

    this._prodDomain = {
      certificateArn: config.frontendCertificateArn,
      domainName: config.frontendDomain,
      customFrontendCondition,
    };

    return this;
  }

  //withDevCDN does not use any existing domains
  withDevCDN(
    stage: string,
    devConfig: DevEnvironmentConfig | null,
    wafAclArn?: string
  ): AppFrontend {
    if (devConfig !== null) {
      const domainName = `${stage}.${devConfig.domainName}`;

      const certificate = Certificate.fromCertificateArn(
        this,
        "Certificate",
        devConfig.certificateArn
      );

      this._domain = {
        certificate,
        domainName,
      };
    }
    /* S3 bucket for react app CDN */
    this.withFrontendBucket();

    /* CDN */
    this._distribution = new cloudfront.Distribution(
      this,
      "CloudfrontDistribution",
      this.distributionConfig({ wafAclArn: wafAclArn })
    );

    return this;
  }

  getDevCallbackUrls(): string[] {
    /* Cognito web app client for Frontend */
    return [`https://${this.getDomainName()}`, "http://localhost:3000"];
  }

  getProdCallbackUrls(): string[] {
    /* Cognito web app client for Frontend */
    return [`https://${this.getDomainName()}`];
  }

  /**
   * returns the CloudFront domain URL (e.g. abcd2oaolmclh.cloudfront.net).
   * Always returns the CloudFront domain, even if a custom domain is set.
   */
  getCloudFrontDomain(): string {
    return this._distribution.domainName;
  }

  getDomainName(): string {
    // if the CFN condition is defined, use a deploy-time lookup rather than a synth-time one.
    if (this._prodDomain !== undefined) {
      return cdk.Fn.conditionIf(
        this._prodDomain.customFrontendCondition.logicalId,
        this._prodDomain.domainName,
        this._distribution.domainName
      ).toString();
    }

    // return the custom domain name if it's configured, rather than the default CloudFront one.
    if (this._domain !== undefined) {
      return this._domain.domainName;
    }

    return this._distribution.domainName;
  }

  getDistributionId(): string {
    return this._distribution.distributionId;
  }

  getBucketName(): string {
    return this._frontendS3Bucket.bucketName;
  }
  getBucket(): Bucket {
    return this._frontendS3Bucket;
  }
}
