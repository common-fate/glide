## Getting Started

Common Fate uses [Playwright](https://playwright.dev/) to run full end to end tests. These tests run on Github actions for each branch and controls part of the CI/CD pipeline.
It's important to make sure these tests are passing while developing. These docs will serve as instructions to make sure you can run them locally to test against local changes.

First make sure the frontend dependencies are up to date by running

```bash
cd web
pnpm install
```

Then you can get your environment ready to run tests locally
Create a `.env` file within `common-fate/web` and copy in the following template

```
TEST_USERNAME=""
TEST_ADMIN_USERNAME=""
TEST_PASSWORD=""
USER_POOL_ID=""
COGNITO_CLIENT_ID=""
TESTING_DOMAIN=""
```

From here we will need to do some work to setup these variables. By this point you should have setup a dev environment using mage. If not checkout out [Deploying](./deploying.md).

- Grab `USER_POOL_ID` and `COGNITO_CLIENT_ID` from your `.env` from your dev deployment, or find the cognito user pool in the console to get these ID's
- Create a test user and test admin user in your cognito user pool (add the admin user to your admin group in cognito)
  _You might need to update their login passwords upon making the accounts, this will be emailed to your email if setup correctly_. - To make extra accounts you can use your own email address with a suffix to make infinite amount of accounts. Eg. `test+1@commonfate.io` where +1 is the suffix
- Once the accounts are created and you have reset the password. Add the usernames and password to the `TEST_USERNAME TEST_ADMIN_USERNAME TEST_PASSWORD` variables
- Set the `TESTING_DOMAIN` to be `http://localhost:3000`

## Running The Tests

To run the tests with live browser windows run:

```
pnpm e2e
```

To run the tests with headless run:

```
pnpm e2e:ci
```

A helpful debugging tool provided with Playwright is the Vscode [extension](https://marketplace.visualstudio.com/items?itemName=ms-playwright.playwright)
