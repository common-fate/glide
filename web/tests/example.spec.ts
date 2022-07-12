import { test, expect } from "@playwright/test";
import { OriginURL } from "./consts";

test("basic test", async ({ page }) => {
  await page.goto(OriginURL);
  const title = page.locator(".navbar__inner .navbar__title");

  // Check that there is a 'Sign in' button
  // await expect(title).toHaveText("Playwright");

  await page.waitForTimeout(1000 * 30);
});
