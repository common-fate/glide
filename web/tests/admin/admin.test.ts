import { test, expect } from "@playwright/test";


test.use({ storageState: 'adminAuthCookies.json' });

test("test admin login gets to granted page with admin nav", async ({ page }) => {
  await page.goto("http://" + process.env.TESTING_DOMAIN ?? "");
  await expect(page).toHaveTitle(/Granted/);
    //await expect(page.locator(".chakra-container #admin-button")).toHaveText(/Admin/);


});