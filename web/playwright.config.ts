import { type PlaywrightTestConfig, devices } from "@playwright/test";
// Read from default ".env" file.
const config: PlaywrightTestConfig = {
  testDir: "./tests",
  forbidOnly: !!process.env.CI,

  retries: 0,
  globalSetup: "./globalSetup.ts",
  use: {
    trace: "on",
    baseURL: process.env.TESTING_DOMAIN,
    // headless: false,
  },
  globalTimeout: process.env.CI ? 60 * 60 * 1000 : undefined,
  projects: [
    //     {
    //   name: 'firefox',
    //   use: { ...devices['Desktop Firefox'] },
    // },
    // {
    //   name: 'webkit',
    //   use: { ...devices['Desktop Safari'] },
    // },
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],
};
export default config;
