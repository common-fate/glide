import { expect, test } from "@playwright/test";
import {
  CreateAccessRule,
  fillFormElementById,
  fillFormElementByTestId,
  LoginAdmin,
  LoginUser,
  selectOptionByID,
  testId,
} from "../utils/helpers";

import { randomBytes } from "crypto";

var id = randomBytes(20).toString("hex");
const ruleName = "test-rule-" + id;
const groupName = "test-group-" + id;

const groupNameUpdated = groupName + "-updated";
const groupDescription = groupName + " description";
const uniqueReason = "test-reason-" + id;

const username = process.env.TEST_USERNAME ?? "";
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

    await fillFormElementById("name", groupName, page);
    await fillFormElementById("description", groupDescription, page);

    await selectOptionByID("user-select-input", username, page);

    await page.click(testId("save-group-button"));

    await page.click(testId(groupName));

    await page.waitForLoadState("networkidle");

    //check we are on the newly created group page
    await page.waitForSelector(testId("group-source"));

    const locator1 = page.locator(testId("group-source"));
    await expect(locator1).toContainText("Internal");

    const locator2 = page.locator(testId("group-name"));
    await expect(locator2).toContainText(groupName);
  });

  test("test update internal group details", async ({ page }) => {
    // This will log us in as an admin
    await LoginAdmin(page);

    await page.waitForLoadState("networkidle");

    //   go to the admin page
    await page.goto("/admin/groups");

    await page.waitForLoadState("networkidle");

    //   find the group we want to update

    await page.click(testId(groupName));

    // await page.waitForNavigation();

    await page.waitForLoadState("networkidle");

    //check we are on the newly created group page
    await page.waitForSelector(testId("group-source"));

    const locator1 = page.locator(testId("group-source"));
    await expect(locator1).toContainText("Internal");

    const locator2 = page.locator(testId("group-name"));
    await expect(locator2).toContainText(groupName);

    const locator = page.locator(testId("group-description"));
    await expect(locator).toContainText(groupDescription);

    await page.click(testId("edit-group"));
    await page.waitForLoadState("networkidle");

    await fillFormElementById("name", groupNameUpdated, page);
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
      groupDescription,
      page
    );

    await page.click(testId("save-group"));
    await page.click(testId("save-group"));
    await page.waitForLoadState("networkidle");

    const locator3 = page.locator(testId("group-source"));
    await expect(locator3).toContainText("Internal");

    const locator4 = page.locator(testId("group-name"));
    await expect(locator4).toContainText(groupNameUpdated);

    const locator5 = page.locator(testId("group-description"));
    await expect(locator5).toContainText(groupDescription);
  });

  test("test create access rule using internal group and request access", async ({
    page,
  }) => {
    // This will log us in as an admin and create an access rule
    await CreateAccessRule(page, ruleName, groupNameUpdated);
    await LoginUser(page);

    //check that we have access to the new rule
    await page.goto("/");
    await page.click(`text=${ruleName}`);
    await page.waitForLoadState("networkidle");

    //add reason
    await fillFormElementById("reasonField", uniqueReason, page);

    await page.click(testId("request-submit-button"));

    await page.waitForLoadState("networkidle");
    // const locator5 = page.locator(testId("req_")).first();
    // await expect(locator5).toBeVisible();
  });

  test("test review access rule", async ({ page }) => {
    //check that we are redirected

    await LoginAdmin(page);

    await page.waitForLoadState("networkidle");
    await page.goto("/reviews?status=pending");
    await page.waitForLoadState("networkidle");
    await page.click(testId(uniqueReason), { force: true });

    await page.waitForLoadState("networkidle");

    await page.locator(testId("approve")).click();

    // Ensure it loads
    await page.waitForLoadState("networkidle");

    // Validate its teh same request
    let approvedText = await page.locator(testId("reason")).textContent();
    await expect(approvedText).toBe(uniqueReason);
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
    await page.click(testId("select-input-deselect-all"));
    await selectOptionByID("group-select", groupNameUpdated, page);
    await page.click(testId("save-group-submit"));

    const locator5 = page.locator(testId(groupNameUpdated)).first();
    await expect(locator5).toBeVisible();
  });
});
