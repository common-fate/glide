import { test, expect, APIRequestContext } from "@playwright/test";
import axios from "axios";
import { LoginAdmin, Logout } from "../utils/helpers";

// Request context is reused by all tests in the file.
let apiContext: APIRequestContext;

test("test request workflow", async ({ page, context, playwright }) => {
  //   await LoginAdmin(page);

  //   const res = await axios.request({
  //     url:
  //       "https://prod.testvault.granted.run/vaults/testvault/members/usr_2CPrGPm7oB35fzT6tmxAwXNrMCf",
  //     method: "GET",
  //   });

  //   await expect(res).toBe("test");
  apiContext = await playwright.request.newContext({
    // All requests we send go to this API endpoint.
    // baseURL: "https://prod.testvault.granted.run/",
    extraHTTPHeaders: {
      // We set this header per GitHub guidelines.
      //   Accept: "application/vnd.github.v3+json",
      // Add authorization token to all requests.
      // Assuming personal access token available in the environment.
      //   Authorization: `token ${process.env.API_TOKEN}`,
    },
  });

  //   https://prod.testvault.granted.run/vaults/2CBsuomHFRE3mrpLGWFaxbyKXG6/members/usr_2CBst12mpb4o9NzXqhqJ5gJvNH3%3B2D
  //   const res = await apiContext.post(
  //     "/vaults/2CBsuomHFRE3mrpLGWFaxbyKXG6/members/usr_2CPrGPm7oB35fzT6tmxAwXNrMCf"
  //   );
  const res = await apiContext.get(
    "https://prod.testvault.granted.run/vaults/2CBsuomHFRE3mrpLGWFaxbyKXG6/members/usr_2CBst12mpb4o9NzXqhqJ5gJvNH3%3B2D"
  );
  //   const res = await apiContext.get("https://api.github.com");
  const ctx = await playwright.request.newContext();

  // const res = await page.request.get("https://api.github.com");

  const body = res.body();

  let json = await res.json();

  console.log(json);

  expect(res.text()).toBe("hello");
  expect(res.json()).toBe("hello");

  expect(res.ok()).toBe(true);

  await apiContext.dispose();
});
