import { defineConfig, devices } from '@playwright/test';

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,
  /* Opt out of parallel tests on CI. */
  workers: process.env.CI ? 1 : undefined,
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  reporter: 'html',

  /* Run your local dev server before starting the tests */
  webServer: [
    {
      command: 'npm --prefix ../../frontend run dev',
      url: 'http://localhost:5173',
      reuseExistingServer: !process.env.CI,
    },
    {
      command: 'npm --prefix ../../backoffice run dev',
      url: 'http://localhost:5174',
      reuseExistingServer: !process.env.CI,
    }
  ],

  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',
  },

  /* Configure projects for major browsers */
  projects: [
    // === Frontend Projects ===
    {
      name: 'frontend:chromium',
      use: {
        ...devices['Desktop Chrome'],
        baseURL: 'http://localhost:5173',
      },
      testDir: './tests/frontend',
    },
    {
      name: 'frontend:firefox',
      use: {
        ...devices['Desktop Firefox'],
        baseURL: 'http://localhost:5173',
      },
      testDir: './tests/frontend',
    },
    {
      name: 'frontend:webkit',
      use: {
        ...devices['Desktop Safari'],
        baseURL: 'http://localhost:5173',
      },
      testDir: './tests/frontend',
    },

    // === Backoffice Projects ===
    {
      name: 'backoffice:chromium',
      use: {
        ...devices['Desktop Chrome'],
        baseURL: 'http://localhost:5174',
      },
      testDir: './tests/backoffice',
    },
    {
      name: 'backoffice:firefox',
      use: {
        ...devices['Desktop Firefox'],
        baseURL: 'http://localhost:5174',
      },
      testDir: './tests/backoffice',
    },
    {
      name: 'backoffice:webkit',
      use: {
        ...devices['Desktop Safari'],
        baseURL: 'http://localhost:5174',
      },
      testDir: './tests/backoffice',
    },
  ],
});
