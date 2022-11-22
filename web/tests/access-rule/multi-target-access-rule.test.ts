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
  let accessRuleId = ""

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
    await page.getByText("first").click();
    await page.locator('internal:attr=[data-testid="argumentField"] >> text=Groups').click()

    await page.getByTestId("argumentGroupMultiSelect").click()
    await page.getByText("all").click();
    await clickFormElementByText("button", "Next", page);
  
    await page.locator("#increment >> nth=0").click();
    await clickFormElementByText("button", "Next", page);
  
    //select max duration for rule
    await fillFormElementById("hour-duration-input", "1", page);
    await clickFormElementByID("form-step-next-button", page);
  
    //click on group select, add both groups for approval
    await page.locator('#group-select div:has-text("Select...") >> nth=1').click()  
    await page.locator('p:has-text("granted_administrators")').click()
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
})


