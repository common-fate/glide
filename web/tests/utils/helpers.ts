import { expect, Page } from "@playwright/test";

//login helper functions

export const LoginUser = async (page: Page) => {
  //  const page = await context.newPage();
  await page.goto("/");
  await fillFormElement(
    "input",
    "username",
    process.env.TEST_USERNAME ?? "",
    page
  );
  await fillFormElement(
    "input",
    "password",
    process.env.TEST_PASSWORD ?? "",
    page
  );
  await clickFormElementByText("input", "Sign in", page);
  // wait for the cognito login to redirect to the app frontend
  await page.waitForNavigation();
  // when auth has been successful, the me api will be called, this means we can continue our test from this point
  await page.waitForRequest(/me/);
};

export const LoginAdmin = async (page: Page) => {
  //  const page = await context.newPage();
  await page.goto("/", { timeout: 10000 });
  await fillFormElement(
    "input",
    "username",
    process.env.TEST_ADMIN_USERNAME ?? "",
    page
  );
  await fillFormElement(
    "input",
    "password",
    process.env.TEST_PASSWORD ?? "",
    page
  );
  await clickFormElementByText("input", "Sign in", page);
  // wait for the cognito login to redirect to the app frontend
  await page.waitForNavigation();
  // when auth has been successful, the me api will be called, this means we can continue our test from this point
  await page.waitForRequest(/me/);
};

export const Logout = async (page: Page) => {
  // wait for redirects to stop
  await page.goto("/logout", { waitUntil: "networkidle" });
};

export const CreateAccessRule = async (
  page: Page,
  ruleName: string,
  group: string
) => {
  await Logout(page);
  await LoginAdmin(page);
  await page.waitForLoadState("networkidle");
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

  //select the test vault provider
  await page.locator(testId("provider-selector-testvault")).click();
  await fillFormElementById("provider-vault", ruleName, page);
  await clickFormElementByID("form-step-next-button", page);

  //select max duration for rule
  await fillFormElementById("hour-duration-input", "1", page);
  await clickFormElementByID("form-step-next-button", page);

  //click on group select, add both groups for approval
  if (group != "") {
    await clickFormElementByID("group-select", page);
    await fillFormElementById("react-select-3-input", group, page);
    await page.keyboard.press("Enter");
    await page.keyboard.press("Escape");
  } else {
    await clickFormElementByID("group-select", page);
    await fillFormElementById("react-select-3-input", "administrator", page);
    await page.keyboard.press("Enter");
    await page.keyboard.press("Escape");
  }
  // page.keyboard.press("Escape");

  //ensure granted_admins was added to selection box
  await clickFormElementByID("form-step-next-button", page);

  //add an approver
  await clickFormElementByClass("chakra-switch", page);

  //ensure granted_admins was added to selection box
  // await clickFormElementByID("approval-group-select", page);
  await page.locator("#user-select >> visible=true").click();
  await page.keyboard.insertText(process.env.TEST_ADMIN_USERNAME ?? "");
  await page.keyboard.press("Enter");
  await page.keyboard.press("Escape");

  // await clickFormElementByID("approval-group-select", page);
  // await page.locator("#approval-group-select").click();
  // await page.keyboard.press("Enter");

  await clickFormElementByID("rule-create-button", page);

  //check to see if the rule was successfully created
  await page.waitForLoadState("networkidle");

  //check that we are redirected
  await expect(page).toHaveURL("/admin/access-rules");
};

//helper functions to click elements that are visible
export const fillFormElement = async (
  inputType: "input" | "textarea",
  name: string,
  value: string,
  page: Page
) => {
  await page.locator(`${inputType}[name=${name}] >> visible=true`).fill(value);
};

export const fillFormElementById = async (
  name: string,
  value: string,
  page: Page
) => {
  await page.locator(`#${name} >> visible=true`).fill(value);
};

export const fillFormElementByTestId = async (
  name: string,
  value: string,
  page: Page
) => {
  await page.locator(`${name} >> visible=true`).fill(value);
};

export const clickFormElementByText = async (
  inputType: "input" | "textarea" | "button",
  name: string,
  page: Page
) => {
  await page
    .locator(`${inputType}:has-text("${name}") >> visible=true`)
    .click();
};

export const clickFormElementByID = async (id: string, page: Page) => {
  await page.locator(`#${id} >> visible=true`).first().click();
};

export const clickFormElementByClass = async (id: string, page: Page) => {
  await page.locator(`.${id} >> visible=true`).click();
};

export const testId = (id: string) => {
  return `[data-testid="${id}"]`;
};

export const uniqueReason = "test-" + Math.floor(Math.random() * 1000);
