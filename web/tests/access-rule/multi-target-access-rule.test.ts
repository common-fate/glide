import { test, expect, Locator } from "@playwright/test";
import {
  Logout,
  LoginAdmin,
  clickFormElementByID,
  clickFormElementByClass,
  fillFormElement,
  clickFormElementByText,
  testId,
  fillFormElementById,
  LoginUser,
  randomRuleName,
  randomDescription,
  selectOptionByID,
  fillFormElementByTestId,
} from "../utils/helpers";
const ruleName = randomRuleName();
const ruleDescription = randomDescription();
const ruleNameUpdated = ruleName + "-updated";
test.describe.serial("Running test sequentially", () => {
  let accessRuleId = "";

  test("admin can create multi target access rule", async ({ page }) => {
    await LoginAdmin(page);

    await page.waitForLoadState("networkidle");
    await expect(page).toHaveTitle(/Common Fate/);
    await page.goto("/admin/access-rules");

    await clickFormElementByID("new-access-rule-button", page);

    await fillFormElement("input", "name", ruleName, page);
    await fillFormElement("textarea", "description", ruleDescription, page);

    await clickFormElementByText("button", "Next", page);
    await page.getByTestId("provider-selector-testgroups").click();

    await selectOptionByID("providerArgumentField", "fifth", page);
    await selectOptionByID(
      "providerArgumentGroupField-group-category",
      "a category containing",
      page
    );
    await clickFormElementByText("button", "Next", page);

    await page.locator("#increment >> nth=0").click();
    await clickFormElementByText("button", "Next", page);

    //select max duration for rule
    await fillFormElementById("hour-duration-input", "1", page);
    await clickFormElementByID("form-step-next-button", page);

    //click on group select, add both groups for approval
    await page
      .locator('#group-select div:has-text("Select...") >> nth=1')
      .click();
    await page.locator("text=everyone >> nth=1").click();
    await page.locator("text=Add or remove groups").click();

    await clickFormElementByID("form-step-next-button", page);

    await clickFormElementByClass("chakra-switch", page);
    await page.locator("#approval-group-select >> visible=true").click();
    await page.keyboard.press("Enter");
    await page.locator("#approval-group-select").click();
    await page.keyboard.press("Enter");

    await clickFormElementByID("rule-create-button", page);

    const response = await page.waitForResponse(
      (response) =>
        response.url().includes("/api/v1/admin/access-rules") &&
        response.status() === 201
    );

    // console.log('the reponse method is', response.request().method() )
    // console.log('the response is', (await response.json()))

    accessRuleId = (await response.json()).id;

    await expect(page.locator("#toast-access-rule-created")).toHaveText(
      "Access rule created"
    );
  });

  test("admin can update existing access rule", async ({ page }) => {
    await LoginAdmin(page);
    await expect(page).toHaveTitle(/Common Fate/);
    await page.goto(`/admin/access-rules/${accessRuleId}`);
    await page.waitForLoadState("networkidle");
    await page.locator(`role=button[name="Edit"] >> nth=0`).click();
    await fillFormElement("input", "name", ruleNameUpdated, page);
    await page.locator(`role=button[name="Update"]`).click();
    const response = await page.waitForResponse(
      (response) =>
        response.url().includes("/api/v1/admin/access-rules") &&
        response.status() === 200
    );
    await expect(page.locator("#toast-access-rule-updated")).toHaveText(
      "Access rule updated"
    );
  });

  test("user can request access to multiple options in one slot", async ({
    page,
  }) => {
    await LoginUser(page);

    await expect(page).toHaveTitle(/Common Fate/);

    // make sure the access rule has permission for the user
    await page.goto(`/access/requests/${accessRuleId}`);

    await clickFormElementByID("user-request-access", page);
    await page.locator("text=fifth >> nth=1").click();
    await page.getByText("Group").click();

    await clickFormElementByID("user-request-access", page);
    await page.locator("text=first >> nth=1").click();
    await page.getByText("Group").click();

    await clickFormElementByID("user-request-access", page);
    await page.locator("text=second >> nth=1").click();
    await page.getByText("Group").click();

    // remove one select item
    await page.locator('role=button[name="Remove second"]').click();

    await page.locator("#increment >> nth=0").click();
    await fillFormElement("textarea", "reason", "need access", page);

    await page.locator('role=button[name="Submit"]').click();
    const response = await page.waitForResponse(
      (response) =>
        response.url().includes("/api/v1/requests") && response.status() === 200
    );
    await expect(page.locator("#toast-user-request-created")).toHaveText(
      "Request created"
    );
  });

  test("user can request access to multiple request slots", async ({
    page,
  }) => {
    await LoginUser(page);

    await expect(page).toHaveTitle(/Common Fate/);

    // make sure the access rule has permission for the user
    await page.goto(`/access/requests/${accessRuleId}`);

    // add first request
    await clickFormElementByID("user-request-access", page);
    await page.locator("text=fifth >> nth=1").click();
    await page.getByText("Group").click();

    // add second request
    await page.locator('role=button[name="add"]').click();
    await page.locator("#subrequest-1").click();
    await page.locator("text=second >> nth=1").click();

    // add third request
    await page.locator('role=button[name="add"]').click();
    await page.locator("#subrequest-2").click();
    await page.locator("text=third >> nth=1").click();

    // remove second request
    await page.locator('role=button[name="remove"] >> nth=1').click();

    await page.locator("#increment >> nth=0").click();
    await fillFormElement("textarea", "reason", "need access", page);

    await page.locator('role=button[name="Submit"]').click();
    const response = await page.waitForResponse(
      (response) =>
        response.url().includes("/api/v1/requests") && response.status() === 200
    );
    await expect(page.locator("#toast-user-request-created")).toHaveText(
      "Request created"
    );
  });

  test("user can favourite an access request", async ({ page }) => {
    await LoginUser(page);
    await expect(page).toHaveTitle(/Common Fate/);
    await page.goto(`/access/requests/${accessRuleId}`);
    await page.waitForLoadState("networkidle");

    const requiredOptions = ["first", "second", "fifth"];

    for (const option of requiredOptions) {
      await clickFormElementByID("user-request-access", page);
      await page.locator(`text=${option} >> nth=1`).click();
    }

    await page.getByRole("button", { name: "Favorite" }).click();
    await page.getByTestId("favourite-request-button").click();
    await page.getByTestId("favourite-request-button").fill("test-fav-one");
    await page.getByRole("button", { name: "Save" }).click();

    const response = await page.waitForResponse(
      (response) =>
        response.url().includes("/api/v1/favorites") &&
        response.status() === 201
    );

    await expect(page.locator("#toast-favourite-created")).toHaveText(
      "Favorite created"
    );
  });

  test("user can see access requests that are favourited", async ({ page }) => {
    await LoginUser(page);
    await expect(page).toHaveTitle(/Common Fate/);

    await page.getByRole("heading", { name: "Favorites" }).click();
    await page.getByTestId("fav-request-item-test-fav-one").first().click();

    // check if the saved request contains required values in the groups.
    await expect(page.locator("#user-request-access >> nth=0")).toHaveText(
      "fifth secondfirst"
    );
  });

  test("user can update a favourite request", async ({ page }) => {
    await LoginUser(page);
    await expect(page).toHaveTitle(/Common Fate/);

    // const requiredOptions = ["first", "second", "fifth"];
    await page.getByRole("heading", { name: "Favorites" }).click();
    await page.getByTestId("fav-request-item-test-fav-one").first().click();

    await page.getByTestId("fav-icon-btn").click();
    // FIXME: playwright is re-rending the input field and unable to update the text value to smt else.¯\_(ツ)_/¯
    await page.getByTestId("favourite-request-button").fill("test-fav-one");
    await page.getByRole("button", { name: "Update" }).click();

    const response = await page.waitForResponse(
      (response) =>
        response.url().includes("/api/v1/favorites") &&
        response.status() === 201
    );

    await expect(page.locator("#toast-favourite-updated")).toHaveText(
      "Favorite updated"
    );
  });

  test("user can delete a favourite request", async ({ page }) => {
    await LoginUser(page);
    await expect(page).toHaveTitle(/Common Fate/);

    // const requiredOptions = ["first", "second", "fifth"];
    await page.getByRole("heading", { name: "Favorites" }).click();
    await page.getByTestId("fav-request-item-test-fav-one").first().click();

    // check if the saved request contains required values in the groups.
    await page.getByTestId("fav-icon-btn").click();
    await page.getByTestId("favourite-request-button").click();

    await page.getByTestId("del-fav-btn").click();

    const response = await page.waitForResponse(
      (response) =>
        response.url().includes("/api/v1/favorites") &&
        response.status() === 200
    );

    await expect(page.locator("#toast-favourite-removed")).toHaveText(
      "Favorite removed"
    );
  });
});
