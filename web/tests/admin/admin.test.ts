import { test, expect } from "@playwright/test";


test("test admin login gets to granted page with admin nav", async ({ browser }) => {
    
    const adminContext = await browser.newContext({ storageState: 'adminAuthCookies.json' });
    const page = await adminContext.newPage();
    await page.goto("/");
    await expect(page).toHaveTitle(/Granted/);
    await expect(page.locator("#admin-button")).toHaveText("Admin");


});