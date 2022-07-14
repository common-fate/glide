import { test, expect } from "@playwright/test";

import {LoginUser, clickFormElementByID, fillFormElement, clickFormElementByText} from '../utils/helpers'

//has to be admin to create access rule

//test user cannot create access rule

test("non admin cannot create access rule", async ({ browser }) => {
    const userContext = await browser.newContext({ storageState: 'userAuthCookies.json' });
    const page = await userContext.newPage();
    await page.goto("/admin/access-rules");
    await expect(page).toHaveTitle(/Granted/);
    await expect(page.locator("#app")).toContainText("Sorry, you  don't have access");


});

//test access rule create

test("admin can create access rule", async ({ browser }) => {
    const adminContext = await browser.newContext({ storageState: 'adminAuthCookies.json' });
    const page = await adminContext.newPage();
    await page.goto("/admin/access-rules");
    await expect(page).toHaveTitle(/Granted/);
    await expect(page.locator(".chakra-container #new-access-rule-button")).toHaveText("New Access Rule");

    //click new access rule
    await clickFormElementByID("new-access-rule-button", page)

    await fillFormElement('input', "name", "test-rule", page)
    await fillFormElement('textarea', "description", "test-rule description", page)
    await clickFormElementByID("form-step-next-button", page)

    // await fillFormElement('input', "name", "test-rule", page)
    // await fillFormElement('input', "name", "test-rule", page)
    // await fillFormElement('input', "name", "test-rule", page)


});