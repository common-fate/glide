import { test, expect } from "@playwright/test";

import {
  clickFormElementByClass,
  clickFormElementByID,
  fillFormElement,
  clickFormElementByText,
  fillFormElementById,
  Logout,
  LoginUser,
  LoginAdmin,
} from "../utils/helpers";

//has to be admin to create access rule

//test user cannot create access rule
test("non admin cannot create access rule", async ({ page }) => {
  await Logout(page);
  await LoginUser(page);
  await page.goto("/");
  await expect(page).toHaveTitle(/Granted/);
  await page.goto("/admin/access-rules").then(() => expect(page.locator("#app")).toContainText(
    "Sorry, you  don't have access"
  ));

});

//test access rule create
test("admin can create access rule", async ({ page }) => {
  await Logout(page);
  await LoginAdmin(page);
  await page.goto("/");
  await clickFormElementByID("admin-button", page);
  await expect(page).toHaveTitle(/Granted/);
  await expect(
    page.locator(".chakra-container #new-access-rule-button")
  ).toHaveText("New Access Rule");

  //click new access rule
  await clickFormElementByID("new-access-rule-button", page);

  //enter a name for new rule
  await fillFormElement("input", "name", "test-rule", page);
  await fillFormElement(
    "textarea",
    "description",
    "test-rule description",
    page
  );
  await clickFormElementByID("form-step-next-button", page);

  //selec the test vault provider
  await clickFormElementByID("provider-selector", page);
  await fillFormElementById("provider-vault", "test", page);
  await clickFormElementByID("form-step-next-button", page);

  //select max duration for rule
  await fillFormElementById("rule-max-duration", "1", page);
  await clickFormElementByID("form-step-next-button", page);

  //click on group select
  await clickFormElementByID("group-select", page);
  await clickFormElementByID("react-select-2-listbox", page);

  //ensure granted_admins was added to selection box
  await clickFormElementByID("form-step-next-button", page);

  //add an approver
  await clickFormElementByClass("chakra-switch", page);

  //ensure granted_admins was added to selection box
  await clickFormElementByID("user-select", page);
  await page.keyboard.press("Enter");

  await clickFormElementByID("rule-create-button", page);

  //check to see if the rule was successfully created

  //check that we are redirected
  await expect(page).toHaveURL("/admin/access-rules");

  // await fillFormElement('input', "name", "test-rule", page)
  // await fillFormElement('input', "name", "test-rule", page)
  // await fillFormElement('input', "name", "test-rule", page)
});
