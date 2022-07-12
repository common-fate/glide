import { type PlaywrightTestConfig, devices } from '@playwright/test';

const config: PlaywrightTestConfig = {
  forbidOnly: !!process.env.CI,
  
 
  retries: 2,
  globalSetup: "./globalSetup.ts",
  use: {
    // Tell all tests to load signed-in state from 'authCookies.json'.
    storageState: "./authCookies.json",
    trace: "on",
    baseURL: process.env.TESTING_DOMAIN,
  },
  globalTimeout: process.env.CI ? 60 * 60 * 1000 : undefined,
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
    {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
  ],
};
export default config;
