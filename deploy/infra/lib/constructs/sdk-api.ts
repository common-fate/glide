import * as cdk from "aws-cdk-lib";
import { CfnCondition, Duration, Stack } from "aws-cdk-lib";
import * as apigv2 from "aws-cdk-lib/aws-apigatewayv2";
import * as lambda from "aws-cdk-lib/aws-lambda";
import * as iam from "aws-cdk-lib/aws-iam";
import { Construct } from "constructs";
import * as path from "path";

interface Props {
  appName: string;
  restApiLambda: lambda.Function;
  audience: string;
  issuer: string;
}

export class SDKAPI extends Construct {
  private readonly _appName: string;
  private _apigateway: apigv2.CfnApi;
  private _createResourcesCondition: CfnCondition;

  constructor(scope: Construct, id: string, props: Props) {
    super(scope, id);

    this._appName = props.appName;

    /**
     * If the audience or issuer are not provided, then do not create these resources
     *
     */
    this._createResourcesCondition = new CfnCondition(
      this,
      "EnableSDKAPICondition",
      {
        expression: cdk.Fn.conditionNot(
          cdk.Fn.conditionOr(
            cdk.Fn.conditionEquals(props.audience, ""),
            cdk.Fn.conditionEquals(props.issuer, "")
          )
        ),
      }
    );
    /**
     * This defines the base API gateway, it is configured as a lambda proxy
     * The default route is defined here as /api/v1/{proxy+} which means that all requests to this base path get proxied to the lambda function handler
     * The route here cannot be changed without recreating the API, this would cause a change in the api URL
     */
    this._apigateway = new apigv2.CfnApi(this, "HTTPAPI", {
      name: this._appName,
      target: props.restApiLambda.functionArn,
      protocolType: "HTTP",
      routeKey: "ANY /api/v1/{proxy+}",
    });
    this._apigateway.cfnOptions.condition = this._createResourcesCondition;

    const lambdaPermission = new lambda.CfnPermission(
      this,
      "APIInvokeLambdaPermission",
      {
        action: "lambda:InvokeFunction",
        functionName: props.restApiLambda.functionName,
        principal: "apigateway.amazonaws.com",
        sourceArn: cdk.Fn.join("", [
          "arn:",
          Stack.of(this).partition,
          ":execute-api:",
          Stack.of(this).region,
          ":",
          Stack.of(this).account,
          ":",
          this._apigateway.ref,
          "/*/*",
        ]),
      }
    );

    lambdaPermission.cfnOptions.condition = this._createResourcesCondition;

    /**
     * Create JWT Authorizer. Here, the 'audience' & 'issuer' should be user provided.
     */
    const jwtAuthorizer = new apigv2.CfnAuthorizer(this, "JWTAuthorizer", {
      apiId: this._apigateway.ref,
      authorizerType: "JWT",
      identitySource: ["$request.header.Authorization"],
      name: "JWTAuthorizer",
      jwtConfiguration: {
        audience: [props.audience],
        issuer: props.issuer,
      },
    });
    jwtAuthorizer.cfnOptions.condition = this._createResourcesCondition;

    /**
     * Managed overrides allow customization of the resources deployed by the quickstart deployment metho that we use
     * The quickstarts mean that we don;t need to deploy everything manually, it also works properly where manually defining each resource has issues
     */
    const managedOverrides = new apigv2.CfnApiGatewayManagedOverrides(
      this,
      "ManagedAPIOverrides",
      {
        apiId: this._apigateway.ref,
        route: {
          authorizationType: "JWT",
          authorizerId: jwtAuthorizer.ref,
        },
        integration: {
          payloadFormatVersion: "1.0",
        },
      }
    );
    managedOverrides.cfnOptions.condition = this._createResourcesCondition;
  }

  getRestApiURL(): string {
    return cdk.Fn.conditionIf(
      this._createResourcesCondition.logicalId,
      cdk.Fn.join("", [
        "https://",
        this._apigateway.ref,
        ".execute-api.",
        Stack.of(this).region,
        ".amazonaws.com",
      ]),
      ""
    ).toString();
  }
}
