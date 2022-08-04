import { test, expect, APIRequestContext } from "@playwright/test";

import { LoginAdmin, Logout } from "../utils/helpers";

// Request context is reused by all tests in the file.
let apiContext: APIRequestContext;

test("verify access is provisioned", async ({ page, context, playwright }) => {
  apiContext = await playwright.request.newContext({});

  const res = await apiContext.get(
    "https://prod.testvault.granted.run/vaults/false_vault_id/members/false_member_id"
  );
  let stringErr = await res.text();
  expect(stringErr).toBe('{"error":"user is not a member of this vault"}');

  // let user = process.env.TEST_USERNAME ?? "jordi@commonfate.io";
  // let vault = "2CBsuomHFRE3mrpLGWFaxbyKXG6_5";

  // const res2 = await apiContext.get(
  //   `https://prod.testvault.granted.run/vaults/${vault}/members/${user}`
  // );

  // let stringSuccess = await res2.text();
  // expect(stringSuccess).toBe(
  //   `{"message":"success! user ${user} is a member of vault ${vault}"}`
  // );

  await apiContext.dispose();
});
