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
  // await CreateAccessRule(page);

  // This will log us out of the admin account
  // await Logout(page);

  // This will log us in as a user
  await LoginUser(page);

  //   We can now validate that the rule is there
  await page.goto("/");

  //   Click on the first rule
  await page.click("#r_0");

  let uniqueReason = "test-" + Math.floor(Math.random() * 100);

  await fillFormElementById("reasonField", uniqueReason, page);

  await clickFormElementByText("button", "Submit", page);

  await page.waitForLoadState("domcontentloaded");

  // const [response] = await Promise.all([
  //   // Clicking the link will indirectly cause a navigation.
  //   clickFormElementByText("button", "Submit", page),
  //   // expect to redirect to home page
  //   page.waitForNavigation({ url: "/", waitUntil: "domcontentloaded" }),
  // ]);

  // expect to see the reason in the list
  await expect(page).toHaveURL("/requests");

  // Click on the first request
  await page.click("#req_" + uniqueReason);

  await page.innerText("#reason").then(async (text) => {
    await expect(text).toBe(uniqueReason);
  });

  // expect(page).

  // page.innerText("#reason").then(async (text) => {

  // let the page hang so we can inspect it
  // await page.waitForTimeout(30000);
});
