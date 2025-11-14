import { defineConfig, devices } from '@playwright/test';

const onlyChromium = process.env.PLAYWRIGHT_ONLY_CHROMIUM === '1';
const includeWebkit = process.env.PLAYWRIGHT_INCLUDE_WEBKIT === '1';

const browserMatrix = (() => {
  if (onlyChromium) return ['chromium'];
  const defaults = ['chromium', 'firefox'];
  return includeWebkit ? [...defaults, 'webkit'] : defaults;
})();

const browserDevice = {
  chromium: devices['Desktop Chrome'],
  firefox: devices['Desktop Firefox'],
  webkit: devices['Desktop Safari'],
} as const;

const createProjects = (
  prefix: 'frontend' | 'backoffice',
  baseURL: string,
  testDir: string,
) =>
  browserMatrix.map((browser) => ({
    name: `${prefix}:${browser}`,
    use: {
      ...browserDevice[browser as keyof typeof browserDevice],
      baseURL,
    },
    testDir,
  }));

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
    ...createProjects('frontend', 'http://localhost:5173', './tests/frontend'),
    ...createProjects('backoffice', 'http://localhost:5174', './tests/backoffice'),
  ],
});
