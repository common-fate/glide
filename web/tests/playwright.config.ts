import { type PlaywrightTestConfig, devices } from '@playwright/test';
import { OriginURL } from "./consts";
import storageState from './storageState.json'

const config: PlaywrightTestConfig = {
  forbidOnly: !!process.env.CI,
  
  retries: process.env.CI ? 2 : 0,
  globalSetup: "./globalSetup.ts",
  use: {
    // Tell all tests to load signed-in state from 'storageState.json'.
    storageState,
    trace: "on-first-retry",
    baseURL: OriginURL,
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
