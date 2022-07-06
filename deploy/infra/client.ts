import {
  Account,
  ClientOptions,
  EventType,
  IAws,
  IPublishProgress,
  IPublishProgressListener,
} from "cdk-assets";
import { LogLevel } from "cdk-assets/bin/logging";
import { log } from "console";
import * as os from "os";

/**
 * The DefaultAwsClient class in cdk-assets is inflexible and
 * has a hard requirement that credentials are loaded from `~/.aws/credentials`.
 *
 * This class removes that requirement.
 *
 * See: https://github.com/aws/aws-cdk/issues/12696
 */
export class AwsClient implements IAws {
  private readonly AWS: typeof import("aws-sdk");
  private account?: Account;

  constructor(profile?: string) {
    if (profile) {
      process.env.AWS_PROFILE = profile;
    }

    // We need to set the environment before we load this library for the first time.
    // eslint-disable-next-line @typescript-eslint/no-require-imports
    this.AWS = require("aws-sdk");
  }

  public async s3Client(options: ClientOptions) {
    return new this.AWS.S3();
  }

  public async ecrClient(options: ClientOptions) {
    return new this.AWS.ECR();
  }

  public async secretsManagerClient(options: ClientOptions) {
    return new this.AWS.SecretsManager();
  }

  public async discoverPartition(): Promise<string> {
    return (await this.discoverCurrentAccount()).partition;
  }

  public async discoverDefaultRegion(): Promise<string> {
    return this.AWS.config.region || "us-east-1";
  }

  public async discoverCurrentAccount(): Promise<Account> {
    if (this.account === undefined) {
      const sts = new this.AWS.STS();
      const response = await sts.getCallerIdentity().promise();
      if (!response.Account || !response.Arn) {
        throw new Error(
          `Unrecognized reponse from STS: '${JSON.stringify(response)}'`
        );
      }
      this.account = {
        accountId: response.Account!,
        partition: response.Arn!.split(":")[1],
      };
    }

    return this.account;
  }

  public async discoverTargetAccount(options: ClientOptions): Promise<Account> {
    const sts = new this.AWS.STS();
    const response = await sts.getCallerIdentity().promise();
    if (!response.Account || !response.Arn) {
      throw new Error(
        `Unrecognized reponse from STS: '${JSON.stringify(response)}'`
      );
    }
    return {
      accountId: response.Account!,
      partition: response.Arn!.split(":")[1],
    };
  }

  private async awsOptions(options: ClientOptions) {
    let credentials;

    if (options.assumeRoleArn) {
      credentials = await this.assumeRole(
        options.region,
        options.assumeRoleArn,
        options.assumeRoleExternalId
      );
    }

    return {
      region: options.region,
      customUserAgent: "cdk-assets",
      credentials,
    };
  }

  /**
   * Explicit manual AssumeRole call
   *
   * Necessary since I can't seem to get the built-in support for ChainableTemporaryCredentials to work.
   *
   * It needs an explicit configuration of `masterCredentials`, we need to put
   * a `DefaultCredentialProverChain()` in there but that is not possible.
   */
  private async assumeRole(
    region: string | undefined,
    roleArn: string,
    externalId?: string
  ): Promise<AWS.Credentials> {
    return new this.AWS.ChainableTemporaryCredentials({
      params: {
        RoleArn: roleArn,
        ExternalId: externalId,
        RoleSessionName: `cdk-assets-${safeUsername()}`,
      },
      stsConfig: {
        region,
        customUserAgent: "cdk-assets",
      },
    });
  }
}

/**
 * Return the username with characters invalid for a RoleSessionName removed
 *
 * @see https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html#API_AssumeRole_RequestParameters
 */
function safeUsername() {
  try {
    return os.userInfo().username.replace(/[^\w+=,.@-]/g, "@");
  } catch (e) {
    return "noname";
  }
}

const EVENT_TO_LEVEL: Record<EventType, LogLevel> = {
  build: "verbose",
  cached: "verbose",
  check: "verbose",
  debug: "verbose",
  fail: "error",
  found: "verbose",
  start: "info",
  success: "info",
  upload: "verbose",
};

export class ConsoleProgress implements IPublishProgressListener {
  public onPublishEvent(type: EventType, event: IPublishProgress): void {
    log(
      EVENT_TO_LEVEL[type],
      `[${event.percentComplete}%] ${type}: ${event.message}`
    );
  }
}
