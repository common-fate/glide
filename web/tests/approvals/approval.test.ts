import { test, expect } from "@playwright/test";
import {
  clickFormElementByText,
  CreateAccessRule,
  fillFormElementById,
  LoginAdmin,
  LoginUser,
  Logout,
  testId,
} from "../utils/helpers";

test.describe.serial("Approval/Request Workflows", () => {
  const uniqueReason = "test-" + Math.floor(Math.random() * 1000);

  test("test request workflow", async ({ page }) => {
    // This will create our Acess Rule for the user account and log us in
    await CreateAccessRule(page);
    // This will log us out of the admin account
    await Logout(page);

    // This will log us in as an admin
    await LoginAdmin(page);

    await page.waitForLoadState("networkidle");

    //   We can now validate that the rule is there
    await page.goto("/");

    await page.waitForLoadState("networkidle");

    //   Click on the first rule
    await page.click(testId("r_0"));

    await fillFormElementById("reasonField", uniqueReason, page);

    await clickFormElementByText("button", "Submit", page);

    await page.waitForLoadState("domcontentloaded");

    // expect to see the reason in the list
    await expect(page).toHaveURL("/requests");

    await page.waitForLoadState("networkidle");

    // Click on the first request
    await page.click(testId("req_" + uniqueReason));

    const locator = page.locator(testId("reason"));
    await expect(locator).toContainText(uniqueReason);
  });

  test("test approval workflow", async ({ page }) => {
    // This will log us out of the admin account
    await Logout(page);

    // This will log us in as an admin
    await LoginAdmin(page);

    await page.waitForLoadState("networkidle");

    await page.goto("/reviews");

    await page.waitForLoadState("networkidle");

    // Click on the first review
    await page.locator(testId("tablerow-0")).click();

    await page.locator(testId("approve")).click();
  });
});
