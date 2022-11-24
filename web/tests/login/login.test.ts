import { test, expect } from "@playwright/test";

import { LoginUser } from "../utils/helpers";

test("test loging through form works and gets to Common Fate page", async ({
  page,
}) => {
  await LoginUser(page);

  //verify login
  //verify we are on the Common Fate homepage
  await expect(page).toHaveTitle(/Common Fate/);
});

test("test login bypass works gets to granted page", async ({ page }) => {
  await LoginUser(page);
  await page.goto("/");
  await expect(page).toHaveTitle(/Common Fate/);
});
