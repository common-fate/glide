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

  test("create an initial Access Rule", async ({ page }) => {
    // This will create our Acess Rule for the user account and log us in
    await CreateAccessRule(page);
    // This will log us out of the admin account
    await Logout(page);
  });

  test("test request workflow", async ({ page }) => {
    // This will log us in as an admin
    await LoginUser(page);

    await page.waitForLoadState("networkidle");

    //   We can now validate that the rule is there
    await page.goto("/");

    await page.waitForLoadState("networkidle");

    //   Click on the first rule
    await page.click(testId("r_0"));

    await fillFormElementById("reasonField", uniqueReason, page);

    await clickFormElementByText("button", "Submit", page);

    await page.waitForLoadState("networkidle");

    // expect to see the reason in the list
    await expect(page).toHaveURL("/requests");

    await page.waitForLoadState("networkidle");

    await page.goto("/requests");

    await page.waitForLoadState("networkidle");

    // Click on the first request
    await page.click(testId("req_" + uniqueReason), { force: true });

    await page.waitForLoadState("networkidle");

    const locator = page.locator(testId("reason"));
    await expect(locator).toContainText(uniqueReason);
  });

  test("test approval workflow", async ({ page }) => {
    // This will log us in as an admin
    await LoginAdmin(page);

    await page.waitForLoadState("networkidle");

    await page.goto("/reviews");

    await page.waitForLoadState("networkidle");

    // Click on the first review
    await page.locator(testId(uniqueReason)).first().click();

    await page.waitForLoadState("networkidle");

    await page.locator(testId("approve")).click();
  });

  test("ensure access granted for matching user", async ({
    page,
    playwright,
  }) => {
    // wait 5s to give the granter time to approve the request
    await page.waitForTimeout(5000);

    let apiContext = await playwright.request.newContext({});
    let user = process.env.TEST_USERNAME ?? "jordi@commonfate.io";
    let vault = process.env.VAULT_ID ?? "2CBsuomHFRE3mrpLGWFaxbyKXG6_5";
    const res = await apiContext.get(
      `https://prod.testvault.granted.run/vaults/${vault}/members/${user}`
    );
    let stringSuccess = await res.text();
    expect(stringSuccess).toBe(
      `{"message":"success! user ${user} is a member of vault ${vault}"}`
    );
    await apiContext.dispose();
  });
});
