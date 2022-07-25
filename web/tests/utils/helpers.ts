import { Page } from "@playwright/test";

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

export const LoginAdmin = async (page) => {
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
  await page.locator(`#${id} >> visible=true`).click();
};

export const clickFormElementByClass = async (id: string, page: Page) => {
  await page.locator(`.${id} >> visible=true`).click();
};
