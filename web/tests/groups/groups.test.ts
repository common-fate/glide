import { test, expect } from "@playwright/test";
import {
  clickFormElementByText,
  CreateAccessRule,
  fillFormElementById,
  fillFormElementByTestId,
  LoginAdmin,
  LoginUser,
  testId,
  uniqueReason,
} from "../utils/helpers";

const RULE_NAME = "test";

test.describe.serial("Internal Groups Workflows", () => {
  test("test create internal group workflow", async ({ page }) => {
    // This will log us in as an admin

    await LoginAdmin(page);

    await page.waitForLoadState("networkidle");

    //   go to the admin page
    await page.goto("/admin/groups");

    await page.waitForLoadState("networkidle");

    //   Click on the add internal group button
    await page.click(testId("create-group-button"));

    await fillFormElementById("name", "group_1_internal", page);
    await fillFormElementById(
      "description",
      "group_1_internal description",
      page
    );
    await fillFormElementById("user-select-input", "jack@commonfate.io", page);
    page.keyboard.press("Enter");

    await page.click(testId("save-group-button"));

    await page.click(testId("group_1_internal"));

    await page.waitForLoadState("networkidle");

    //check we are on the newly created group page
    await page.waitForSelector(testId("group-source"));

    const locator1 = page.locator(testId("group-source"));
    await expect(locator1).toContainText("Internal");

    const locator2 = page.locator(testId("group-name"));
    await expect(locator2).toContainText("group_1_internal");
  });

  test("test update internal group details", async ({ page }) => {
    // This will log us in as an admin
    await LoginAdmin(page);

    await page.waitForLoadState("networkidle");

    //   go to the admin page
    await page.goto("/admin/groups");

    await page.waitForLoadState("networkidle");

    //   find the group we want to update

    await page.click(testId("group_1_internal"));

    // await page.waitForNavigation();

    await page.waitForLoadState("networkidle");

    //check we are on the newly created group page
    await page.waitForSelector(testId("group-source"));

    const locator1 = page.locator(testId("group-source"));
    await expect(locator1).toContainText("Internal");

    const locator2 = page.locator(testId("group-name"));
    await expect(locator2).toContainText("group_1_internal");

    const locator = page.locator(testId("group-description"));
    await expect(locator).toContainText("group_1_internal description");

    await page.click(testId("edit-group"));
    await page.waitForLoadState("networkidle");

    await fillFormElementById("name", "group_1_internal_updated", page);
    await page.waitForLoadState("networkidle");

    await fillFormElementById(
      "user-select-input",
      process.env.TEST_ADMIN_USERNAME ?? "",
      page
    );
    page.keyboard.press("Enter");

    await page.waitForLoadState("networkidle");

    //need to focus the input first by clicking
    await page.click(testId("description"));

    await fillFormElementByTestId(
      testId("description"),
      "group_1_internal updated description",
      page
    );

    await page.click(testId("save-group"));
    await page.click(testId("save-group"));
    await page.waitForLoadState("networkidle");

    const locator3 = page.locator(testId("group-source"));
    await expect(locator3).toContainText("Internal");

    const locator4 = page.locator(testId("group-name"));
    await expect(locator4).toContainText("group_1_internal_updated");

    const locator5 = page.locator(testId("group-description"));
    await expect(locator5).toContainText(
      "group_1_internal updated description"
    );
  });

  test("test create access rule using internal group and request access", async ({
    page,
  }) => {
    // This will log us in as an admin and create an access rule
    await CreateAccessRule(page, RULE_NAME, "group_1_internal_updated");

    await LoginUser(page);

    //check that we have access to the new rule
    await page.goto("/");
    await page.click(testId("r_0"));
    await page.waitForLoadState("networkidle");

    //add reason
    await fillFormElementById("reasonField", uniqueReason, page);

    await page.click(testId("request-submit-button"));

    await page.waitForLoadState("networkidle");
    // const locator5 = page.locator(testId("req_")).first();
    // await expect(locator5).toBeVisible();
  });

  test("test review access rule", async ({ page }) => {
    // This will log us in as an admin
    const tempUniqueReason = uniqueReason;

    //create a new request
    await LoginUser(page);

    //check that we have access to the new rule
    await page.goto("/");
    await page.click(testId("r_0"));
    await page.waitForLoadState("networkidle");

    //add reason
    await fillFormElementById("reasonField", tempUniqueReason, page);

    await page.click(testId("request-submit-button"));

    await page.waitForLoadState("networkidle");
    await expect(page).toHaveURL("/requests");

    page.waitForNavigation();

    //check that we are redirected

    await LoginAdmin(page);

    await page.waitForLoadState("networkidle");
    await page.goto("/reviews?status=pending");
    await page.waitForLoadState("networkidle");
    await page.click(testId(tempUniqueReason), { force: true });

    await page.waitForLoadState("networkidle");

    await page.locator(testId("approve")).click();

    // Ensure it loads
    await page.waitForLoadState("networkidle");

    // Validate its teh same request
    let approvedText = await page.locator(testId("reason")).textContent();
    await expect(approvedText).toBe(tempUniqueReason);
  });

  test("test admin can update user groups", async ({ page }) => {
    // This will log us in as an admin
    await LoginAdmin(page);

    await page.waitForLoadState("networkidle");

    //   go to the admin page
    await page.goto("/admin/users");

    await page.waitForLoadState("networkidle");

    //   find the group we want to update

    await page.click(testId(process.env.TEST_USERNAME ?? ""));

    // await page.waitForNavigation();

    await page.waitForLoadState("networkidle");
    await page.click(testId("edit-groups-icon"));
    await page.waitForLoadState("networkidle");

    await fillFormElementById(
      "group-select-input",
      "group_1_internal_updated",
      page
    );
    //wait for the dropdown to fillout
    await new Promise((r) => setTimeout(r, 2000));

    await page.keyboard.press("Enter");
    await page.keyboard.press("Enter");

    // await page.click(testId("save-group-submit"));

    const locator5 = page.locator(testId("group_1_internal_updated")).first();
    await expect(locator5).toBeVisible();
  });
});
