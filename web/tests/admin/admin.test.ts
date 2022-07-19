import { test, expect } from "@playwright/test";
import { LoginAdmin, Logout } from "../utils/helpers";

test("test admin login gets to granted page with admin nav", async ({
  page,
}) => {
  await Logout(page);
  await LoginAdmin(page);
  await page.goto("/");
  await expect(page).toHaveTitle(/Granted/);
  await expect(page.locator("#admin-button")).toHaveText("Admin");
});
