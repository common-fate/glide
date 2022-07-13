import { test, expect } from "@playwright/test";

import {fillFormElement, clickFormElement} from "../utils/helpers"


test("test loging through form works and gets to granted page", async ({ page }) => {
  await page.goto("http://" + process.env.TESTING_DOMAIN ?? "");
  await fillFormElement("input", "username", process.env.TEST_USERNAME ?? "", page)
  await fillFormElement("input", "password", process.env.TEST_PASSWORD ?? "", page)
  // await page.click('input[value="Sign in"]')

  await clickFormElement("input", "Sign in", page)

  //verify login

  //verify we are on the granted homepage
    await expect(page).toHaveTitle(/Granted/);

});


test.use({ storageState: 'userAuthCookies.json' });

test("test login bypass works gets to granted page", async ({ page }) => {
  await page.goto("http://" + process.env.TESTING_DOMAIN ?? "");
  await expect(page).toHaveTitle(/Granted/);
});


