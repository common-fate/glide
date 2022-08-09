import { test, expect, Page } from "@playwright/test";
import config from "../../playwright.config";

import {
  clickFormElementByClass,
  clickFormElementByID,
  fillFormElement,
  clickFormElementByText,
  fillFormElementById,
  Logout,
  LoginUser,
  LoginAdmin,
  CreateAccessRule,
} from "../utils/helpers";

//has to be admin to create access rule

//test user cannot create access rule
test("non admin cannot create access rule", async ({ page }) => {
  await Logout(page);
  await LoginUser(page);
  await page.waitForLoadState("networkidle");
  await expect(page).toHaveTitle(/Granted/);
  await page
    .goto("/admin/access-rules")
    .then(() =>
      expect(page.locator("#app")).toContainText(
        "Sorry, you  don't have access"
      )
    );
});

//test access rule create
test("admin can create access rule", async ({ page }) => {
  await CreateAccessRule(page);
});
