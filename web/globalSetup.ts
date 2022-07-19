import {
  chromium,
  FullConfig,
  PlaywrightTestConfig,
  PlaywrightTestOptions,
} from "@playwright/test";
import { Auth } from "@aws-amplify/auth";
import dotenv from "dotenv";

// Read from default ".env" file.
dotenv.config();

async function globalSetup(config: FullConfig) {
  //normal user
  const username = process.env.TEST_USERNAME ?? "";
  const password = process.env.TEST_PASSWORD;

  //admin user
  const adminUsername = process.env.TEST_ADMIN_USERNAME ?? "";
  const adminPassword = process.env.TEST_PASSWORD;

  //setup aws config
  const userPoolId = process.env.USER_POOL_ID;
  const clientId = process.env.COGNITO_CLIENT_ID;
  const awsconfig = {
    aws_user_pools_id: userPoolId,
    aws_user_pools_web_client_id: clientId,
  };
  Auth.configure(awsconfig);

  //sign in using user creds and save the login creds to a json containing cookies to be inserted into the tests
  //These functions will create the json including login creds
  //user
  await LoginUser(username, password)
   
  // // ... log in

  // //admin
  await LoginAdmin(adminUsername, adminPassword)

 

  
}
export default globalSetup;


const LoginUser = async (username, password, ) => {
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
          expires: -1,
          httpOnly: false,
          sameSite: "Strict",
          secure: false,
        },
        {
          name: makeKey("clockDrift"),
          value: "0",
          domain: process.env.TESTING_DOMAIN ?? "",
          path: "/",
          expires: -1,
          httpOnly: false,
          sameSite: "Strict",
          secure: false,
        },
        {
          name: "amplify-signin-with-hostedUI",
          value: "false",
          domain: process.env.TESTING_DOMAIN ?? "",
          path: "/",
          expires: -1,
          httpOnly: false,
          sameSite: "Strict",
          secure: false,
        },
        {
          name: makeKey("accessToken"),
          value: cognitoUser.signInUserSession.accessToken.jwtToken,
          domain: process.env.TESTING_DOMAIN ?? "",
          path: "/",
          expires: -1,
          httpOnly: false,
          sameSite: "Strict",
          secure: false,
        },
        {
          name: `CognitoIdentityServiceProvider.${cognitoUser.pool.clientId}.LastAuthUser`,
          value: cognitoUser.username,
          domain: process.env.TESTING_DOMAIN ?? "",
          path: "/",
          expires: -1,
          httpOnly: false,
          sameSite: "Strict",
          secure: false,
        },
      ],
      origins: [
        {
          origin: "http://" + process.env.TESTING_DOMAIN ?? "",
          localStorage: [],
        },
      ],
    };

    const data = JSON.stringify(amplifyData);
    const fs = require("fs");
    fs.writeFile("./userAuthCookies.json", data, (err) => {
      if (err) {
        throw err;
      }
      console.log(
        `AWS Cognito login information stored in storageState.json for: ${username}`
      );
    });
  });
}


const LoginAdmin = async (username, password, ) => {
  await Auth.signIn(username, password).then(async (cognitoUser) => {
    const makeKey = (name) =>
      `CognitoIdentityServiceProvider.${cognitoUser.pool.clientId}.${cognitoUser.username}.${name}`;

    let adminData: PlaywrightTestOptions["storageState"] = {
      cookies: [
        {
          name: makeKey("idToken"),
          value: cognitoUser.signInUserSession.idToken.jwtToken,
          domain: process.env.TESTING_DOMAIN ?? "",
          path: "/",
          expires: -1,
          httpOnly: false,
          sameSite: "Strict",
          secure: true,
        },
        {
          name: makeKey("clockDrift"),
          value: "0",
          domain: process.env.TESTING_DOMAIN ?? "",
          path: "/",
          expires: -1,
          httpOnly: false,
          sameSite: "Strict",
          secure: true,
        },
        {
          name: "amplify-signin-with-hostedUI",
          value: "false",
          domain: process.env.TESTING_DOMAIN ?? "",
          path: "/",
          expires: -1,
          httpOnly: false,
          sameSite: "Strict",
          secure: true,
        },
        {
          name: makeKey("accessToken"),
          value: cognitoUser.signInUserSession.accessToken.jwtToken,
          domain: process.env.TESTING_DOMAIN ?? "",
          path: "/",
          expires: -1,
          httpOnly: false,
          sameSite: "Strict",
          secure: true,
        },
        {
          name: `CognitoIdentityServiceProvider.${cognitoUser.pool.clientId}.LastAuthUser`,
          value: cognitoUser.username,
          domain: process.env.TESTING_DOMAIN ?? "",
          path: "/",
          expires: -1,
          httpOnly: false,
          sameSite: "Strict",
          secure: true,
        },
      ],
      origins: [
        {
          origin: "https://" + process.env.TESTING_DOMAIN ?? "",
          localStorage: [],
        },
      ],
    };

    const data = JSON.stringify(adminData);
    const fs = require("fs");
    fs.writeFile("./adminAuthCookies.json", data, (err) => {
      if (err) {
        throw err;
      }
      console.log(
        `AWS Cognito login information stored in storageState.json for: ${username}`
      );
    });
  }).catch((e) => {console.log(`error: ${e}`)});
}