import { test, expect } from "@playwright/test";

import { Logout, LoginUser } from "../utils/helpers";

// test("test loging through form works and gets to granted page", async ({
//   page,
// }) => {
//   await Logout(page);
//   await LoginUser(page);

//   //verify login
//   //verify we are on the granted homepage
//   await expect(page).toHaveTitle(/Granted/);
// });

// test("test login bypass works gets to granted page", async ({ page }) => {
//   await Logout(page)
//     await LoginUser(page)
//   await page.goto("/");
//   await expect(page).toHaveTitle(/Granted/);
// });
