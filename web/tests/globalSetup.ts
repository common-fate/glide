import {
  chromium,
  FullConfig,
  PlaywrightTestConfig,
  PlaywrightTestOptions,
} from "@playwright/test";
import { Auth } from "aws-amplify";
import fs from "fs";
import dotenv from "dotenv";
import { OriginURL } from "./consts";

// Read from default ".env" file.
dotenv.config();

async function globalSetup(config: FullConfig) {
  const username = process.env.TEST_USERNAME ?? "";
  const password = process.env.TEST_PASSWORD;
  // get the userPoolId from Amazon Cognito
  const userPoolId = "ap-southeast-2_XaL0vYUTJ";

  // get the clientId from Amazon Cognito :
  // teamclientb2h5644q_userpool_b2h5644q-dev > App Integration, App Client List > bentlyb2h5644q_app_clientWeb
  const clientId = "7afe1lvpncn20dfm50rjh7nov";
  const awsconfig = {
    aws_user_pools_id: userPoolId,
    aws_user_pools_web_client_id: clientId,
  };
  Auth.configure(awsconfig);
  await Auth.signIn(username, password).then(async (cognitoUser) => {
    const makeKey = (name) =>
      `CognitoIdentityServiceProvider.${cognitoUser.pool.clientId}.${cognitoUser.username}.${name}`;
    let amplifyData: PlaywrightTestOptions["storageState"] = {
      cookies: [],
      origins: [
        {
          origin:
            "https://granted-login-cd-dev-test.auth.ap-southeast-2.amazoncognito.com/login",
          localStorage: [],
        },
      ],
    };
    amplifyData.cookies = [
      {
        name: makeKey("idToken"),
        value: cognitoUser.signInUserSession.idToken.jwtToken,
      },
      {
        name: makeKey("clockDrift"),
        value: "0",
      },
      {
        name: "amplify-signin-with-hostedUI",
        value: "false",
      },
      {
        name: makeKey("accessToken"),
        value: cognitoUser.signInUserSession.accessToken.jwtToken,
      },
      {
        name: `CognitoIdentityServiceProvider.${cognitoUser.pool.clientId}.LastAuthUser`,
        value: cognitoUser.username,
      },
    ];
    const data = JSON.stringify(amplifyData);
    fs.writeFile("./tests/storageState.json", data, (err) => {
      if (err) {
        throw err;
      }
      console.log(
        `AWS Cognito login information stored in storageState.json for: ${username}`
      );
    });
  });
}
export default globalSetup;
