import {
  chromium,
  FullConfig,
  PlaywrightTestConfig,
  PlaywrightTestOptions,
} from "@playwright/test";
import { Auth } from "aws-amplify";
import fs from "fs";
import dotenv from "dotenv";

// Read from default ".env" file.
dotenv.config();

async function globalSetup(config: FullConfig) {
  const username = process.env.TEST_USERNAME ?? "";
  const password = process.env.TEST_PASSWORD;
  // get the userPoolId from Amazon Cognito
  const userPoolId = process.env.USER_POOL_ID;

  // get the clientId from Amazon Cognito :
  const clientId = process.env.COGNITO_CLIENT_ID;
  const awsconfig = {
    aws_user_pools_id: userPoolId,
    aws_user_pools_web_client_id: clientId,
  };
  Auth.configure(awsconfig);
  await Auth.signIn(username, password).then(async (cognitoUser) => {
    const makeKey = (name) =>
      `CognitoIdentityServiceProvider.${cognitoUser.pool.clientId}.${cognitoUser.username}.${name}`;


    let amplifyData: PlaywrightTestOptions["storageState"] = {
      cookies: [
        {
        name: makeKey("idToken"),
        value: cognitoUser.signInUserSession.idToken.jwtToken,
        domain: process.env.TESTING_DOMAIN ?? "",
        path: "/",
      },
      {
        name: makeKey("clockDrift"),
        value: "0",
        domain: process.env.TESTING_DOMAIN ?? "",
        path: "/",
      },
      {
        name: "amplify-signin-with-hostedUI",
        value: "false",
        domain: process.env.TESTING_DOMAIN ?? "",
        path: "/",
      },
      {
        name: makeKey("accessToken"),
        value: cognitoUser.signInUserSession.accessToken.jwtToken,
        domain: process.env.TESTING_DOMAIN ?? "",
        path: "/",
      },
      {
        name: `CognitoIdentityServiceProvider.${cognitoUser.pool.clientId}.LastAuthUser`,
        value: cognitoUser.username,
        domain: process.env.TESTING_DOMAIN ?? "",
        path: "/",
      },
    ],
    origins: [
          {
            "origin": "https://" + process.env.TESTING_DOMAIN ?? "",
            "localStorage": []
          }
        ]
    };

    
    const data = JSON.stringify(amplifyData);
    fs.writeFile("./authCookies.json", data, (err) => {
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
