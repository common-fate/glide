import { test, expect } from "@playwright/test";

test("test login gets to granted page", async ({ page }) => {
  await page.goto("https://" + process.env.TESTING_DOMAIN ?? "");
  await expect(page).toHaveTitle(/Granted/);
});


