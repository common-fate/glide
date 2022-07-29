import { test, expect } from "@playwright/test";
import {
  clickFormElementByText,
  CreateAccessRule,
  fillFormElementById,
  LoginAdmin,
  LoginUser,
  Logout,
} from "../utils/helpers";

test("test approval workflow", async ({ page }) => {
  // This will create our Acess Rule for the user account and log us in
  await CreateAccessRule(page);

  // This will log us out of the admin account
  await Logout(page);

  // This will log us in as a user
  await LoginUser(page);

  //   We can now validate that the rule is there
  await page.goto("/");

  //   Click on the first rule
  await page.click("#r_0");

  await page.waitForNavigation();

  let uniqueReason = "test-" + Math.random().toString();

  await fillFormElementById("reasonField", uniqueReason, page);

  await page.waitForNavigation();

  await clickFormElementByText("button", "Submit", page);
});
