import { APIRequestContext, expect, test } from "@playwright/test";

// Request context is reused by all tests in the file.
let apiContext: APIRequestContext;

test("verify access is provisioned", async ({ page, context, playwright }) => {
  apiContext = await playwright.request.newContext({});

  const res = await apiContext.get(
    "https://prod.testvault.granted.run/vaults/false_vault_id/members/false_member_id"
  );
  let stringErr = await res.text();
  expect(stringErr).toBe('{"error":"user is not a member of this vault"}');

  await apiContext.dispose();
});
