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
} from "../utils/helpers";

test.describe.serial("Running test sequentially", () => {
  let accessRuleId = "";

  test("admin can create multi target access rule", async ({ page }) => {
    await Logout(page);
    await LoginAdmin(page);

    await page.waitForLoadState("networkidle");
    await expect(page).toHaveTitle(/Granted/);
    await page.goto("/admin/access-rules");

    await clickFormElementByID("new-access-rule-button", page);

    await fillFormElement("input", "name", "new-access-rule", page);
    await fillFormElement("textarea", "description", "test description", page);

    await clickFormElementByText("button", "Next", page);
    await page.getByTestId("provider-selector-testgroups").click();

    await page.getByTestId("argumentField").click()
    await page.locator('text=fifth >> nth=1').click()
    await page.locator('internal:attr=[data-testid="argumentField"] >> text=Groups').click()

    await page.locator('internal:attr=[data-testid="arg-group-view"] >> div:has-text("Select...") >> nth=3').click()
    await page.getByText("alla category containing all groupsall").click();
    await clickFormElementByText("button", "Next", page);

    await page.locator("#increment >> nth=0").click();
    await clickFormElementByText("button", "Next", page);

    //select max duration for rule
    await fillFormElementById("hour-duration-input", "1", page);
    await clickFormElementByID("form-step-next-button", page);

    //click on group select, add both groups for approval
    await page.locator('#group-select div:has-text("Select...") >> nth=1').click()
    await page.locator('text=granted_administrators >> nth=1').click()
    await page.locator('text=Add or remove groups').click()

    await clickFormElementByID("form-step-next-button", page);

    await clickFormElementByClass("chakra-switch", page);
    await page.locator("#approval-group-select >> visible=true").click();
    await page.keyboard.press("Enter");
    await page.locator("#approval-group-select").click();
    await page.keyboard.press("Enter");

    await clickFormElementByID("rule-create-button", page);

    const response = await page.waitForResponse(response =>  response.url().includes("/api/v1/admin/access-rules") && response.status() === 201 )

    // console.log('the reponse method is', response.request().method() )
    // console.log('the response is', (await response.json()))

    accessRuleId = (await response.json()).id

    await expect(page.locator("#toast-access-rule-created")).toHaveText("Access rule created")
  });

  test("admin can update existing access rule", async({page }) => {
    await Logout(page);
    await LoginAdmin(page);

    await expect(page).toHaveTitle(/Granted/);
    await page.goto(`/admin/access-rules/${accessRuleId}`);

    await page.locator(`role=button[name="Edit"] >> nth=0`).click()
    await fillFormElement("input", "name", "new-access-rule-updated", page);

    await page.locator(`role=button[name="Update"]`).click()

    await expect(page.locator("#toast-access-rule-updated")).toHaveText("Access rule updated")
  })

  test("user can request access to multiple options in one slot", async ({
    page,
  }) => {
    await Logout(page);
    await LoginUser(page);

    await expect(page).toHaveTitle(/Granted/);

    // make sure the access rule has permission for the user
    await page.goto(`/access/request/${accessRuleId}`);

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
    await page.locator('role=button[name="Remove second"]').click()

    await page.locator("#increment >> nth=0").click();
    await fillFormElement("textarea", "reason", "need access", page);

    await page.locator('role=button[name="Submit"]').click()

    await expect(page.locator("#toast-user-request-created")).toHaveText("Request created")
  });

  test("user can request access to multiple request slots", async ({
    page,
  }) => {
    await Logout(page);
    await LoginUser(page);

    await expect(page).toHaveTitle(/Granted/);

    // make sure the access rule has permission for the user
    await page.goto(`/access/request/${accessRuleId}`);

    // add first request
    await clickFormElementByID("user-request-access", page);
    await page.locator("text=fifth >> nth=1").click();
    await page.getByText("Group").click();

    // add second request
    await page.locator('role=button[name="add"]').click();
    await page.locator('#subrequest-1').click();
    await page.locator("text=second >> nth=1").click();

    // add third request
    await page.locator('role=button[name="add"]').click();
    await page.locator('#subrequest-2').click();
    await page.locator("text=third >> nth=1").click();

    // remove second request 
    await page.locator('role=button[name="remove"] >> nth=1').click()

    await page.locator("#increment >> nth=0").click();
    await fillFormElement("textarea", "reason", "need access", page);

    await page.locator('role=button[name="Submit"]').click()

    await expect(page.locator("#toast-user-request-created")).toHaveText("Request created")
  });
});
