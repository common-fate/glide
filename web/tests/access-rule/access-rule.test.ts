import { expect, test } from "@playwright/test";

import { CreateAccessRule, LoginUser } from "../utils/helpers";

//has to be admin to create access rule

//test user cannot create access rule
test("non admin cannot create access rule", async ({ page }) => {
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
import { randomBytes } from "crypto";

var id = randomBytes(20).toString("hex");
const ruleName = "test-rule-" + id;
//test access rule create
test("admin can create access rule", async ({ page }) => {
  await CreateAccessRule(page, ruleName, "");
});
