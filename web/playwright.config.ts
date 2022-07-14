import { type PlaywrightTestConfig, devices } from '@playwright/test';

const config: PlaywrightTestConfig = {
  forbidOnly: !!process.env.CI,
  
 
  retries: 0,
  globalSetup: "./globalSetup.ts",
  use: {
    
    trace: "on",
    baseURL: "http://" + process.env.TESTING_DOMAIN,
  },
  globalTimeout: process.env.CI ? 60 * 60 * 1000 : undefined,
  projects: [
        {
      name: 'firefox',
      use: { ...devices['Desktop Firefox'] },
    },
    {
      name: 'webkit',
      use: { ...devices['Desktop Safari'] },
    },
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },

  ],
};
export default config;
