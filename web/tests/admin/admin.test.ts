import { expect, test } from "@playwright/test";
import { LoginAdmin } from "../utils/helpers";

test("test admin login gets to Common Fate page with admin nav", async ({
  page,
}) => {
  await LoginAdmin(page);
  await page.goto("/");
  await expect(page).toHaveTitle(/Common Fate/);
  await expect(page.locator("#admin-button >> visible=true")).toHaveText(
    "Switch To Admin"
  );
});
