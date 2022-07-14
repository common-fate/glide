import { test, expect } from "@playwright/test";

import {fillFormElement, clickFormElement, LoginUser} from "../utils/helpers"

test("test loging through form works and gets to granted page", async ({ browser }) => {
   test.slow();
  const noUserContext = await browser.newContext({ storageState: undefined });
  const page = await noUserContext.newPage();
  await page.goto("/");
  await fillFormElement("input", "username", process.env.TEST_USERNAME ?? "", page)
  await fillFormElement("input", "password", process.env.TEST_PASSWORD ?? "", page)
  await clickFormElement("input", "Sign in", page)

  //verify login
  //verify we are on the granted homepage
    await expect(page).toHaveTitle(/Granted/);

});


test("test login bypass works gets to granted page", async ({ browser }) => {
  const userContext = await browser.newContext({ storageState: 'userAuthCookies.json' });
  const page = await userContext.newPage();
  await page.goto("/");
  await expect(page).toHaveTitle(/Granted/);
});


