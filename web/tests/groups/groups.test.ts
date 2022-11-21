import { test, expect } from "@playwright/test";
import {
  clickFormElementByText,
  CreateAccessRule,
  fillFormElementById,
  fillFormElementByTestId,
  LoginAdmin,
  LoginUser,
  Logout,
  testId,
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
    await fillFormElementById(
      "react-select-3-input",
      "jack@commonfate.io",
      page
    );
    page.keyboard.press("Enter");

    await page.click(testId("save-group-button"));

    await page.click(testId("group_1_internal"));

    await page.waitForNavigation();

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
      "react-select-3-input",
      "jack@commonfate.io",
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

  test("test create access rule using internal group", async ({ page }) => {
    // This will log us in as an admin
    await CreateAccessRule(page, RULE_NAME, "group_1_internal_updated");

    //check that we have access to the new rule
    await page.goto("/");
  });
});
