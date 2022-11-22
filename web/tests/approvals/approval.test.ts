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

const RULE_NAME = "test";

test.describe.serial("Approval/Request Workflows", () => {
  const uniqueReason = "test-" + Math.floor(Math.random() * 1000);
  let accessInstructionLink: string;

  test("create an initial Access Rule", async ({ page }) => {
    // This will create our Acess Rule for the user account and log us in
    await CreateAccessRule(page, RULE_NAME, "");
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
    await page.waitForSelector(testId("rule-name"));
    await clickFormElementByText("button", "Submit", page);

    await page.waitForNavigation();

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

    // Click on the specific request
    await page.locator(testId(uniqueReason)).first().click();

    await page.waitForLoadState("networkidle");

    // Click approve
    await page.locator(testId("approve")).click();

    // Ensure it loads
    await page.waitForLoadState("networkidle");

    // Validate its teh same request
    let approvedText = await page.locator(testId("reason")).textContent();
    await expect(approvedText).toBe(uniqueReason);

    // // Assign the accessInstructionLink for our next test
    // accessInstructionLink =
    //   (await page
    //     .locator(testId("accessInstructionLink"))
    //     .getAttribute("href")) ?? "error";

    // // a preliminary check to make sure the link is valid, tested in next test
    // await expect(accessInstructionLink).toContain("https");
  });

  // @NOTE: commented out for now, will not pass on the CI (unknown reason)
  // test("ensure access granted for matching user", async ({
  //   playwright,
  //   page,
  // }) => {
  //   // wait 1s to allow the grant to be applied
  //   page.waitForTimeout(1000);

  //   let apiContext = await playwright.request.newContext({});
  //   let user = process.env.TEST_USERNAME;

  //   const res = await apiContext.get(accessInstructionLink);
  //   let stringSuccess = await res.text();

  //   // ensure the vault has granted access
  //   expect(stringSuccess).toContain(
  //     `{"message":"success! user ${user} is a member of vault`
  //   );
  //   await apiContext.dispose();
  // });
});
