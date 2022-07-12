import type { PlaywrightTestConfig } from "@playwright/test";

const config: PlaywrightTestConfig = {
  forbidOnly: !!process.env.CI,
  testDir: "test",
  testMatch: "integration/**/*.test.ts",
  retries: process.env.CI ? 2 : 0,
  globalSetup: "./globalSetup",
  use: {
    // Tell all tests to load signed-in state from 'storageState.json'.
    storageState: "storageState.json",
    trace: "on",
    baseURL: "https://develop.sk9932ksdi.amplifyapp.com/",
  },
  globalTimeout: process.env.CI ? 60 * 60 * 1000 : undefined,
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
    // {
    //   name: 'firefox',
    //   use: { ...devices['Desktop Firefox'] },
    // },
    // {
    //   name: 'webkit',
    //   use: { ...devices['Desktop Safari'] },
    // },
  ],
};
export default config;
