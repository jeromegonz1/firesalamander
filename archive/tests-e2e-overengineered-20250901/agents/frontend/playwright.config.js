// @ts-check
const { defineConfig, devices } = require('@playwright/test');

/**
 * Fire Salamander - Playwright Configuration
 * Tests E2E pour l'interface utilisateur
 */
module.exports = defineConfig({
  testDir: './tests',
  outputDir: 'test-results/',
  
  /* Run tests in files in parallel */
  fullyParallel: true,
  
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  
  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,
  
  /* Opt out of parallel tests on CI. */
  workers: process.env.CI ? 1 : undefined,
  
  /* Reporter to use */
  reporter: [
    ['html', { outputFolder: '../../../tests/reports/frontend/html' }],
    ['json', { outputFile: '../../../tests/reports/frontend/results.json' }],
    ['junit', { outputFile: '../../../tests/reports/frontend/results.xml' }]
  ],
  
  /* Shared settings for all the projects below */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: process.env.BASE_URL || 'http://localhost:3000',
    
    /* Collect trace when retrying the failed test */
    trace: 'on-first-retry',
    
    /* Take screenshot on failure */
    screenshot: 'only-on-failure',
    
    /* Record video on failure */
    video: 'retain-on-failure',
    
    /* Navigation timeout */
    navigationTimeout: 30000,
    
    /* Action timeout */
    actionTimeout: 10000
  },

  /* Configure projects for major browsers */
  projects: [
    {
      name: 'chromium',
      use: { 
        ...devices['Desktop Chrome'],
        viewport: { width: 1280, height: 720 }
      },
    },

    {
      name: 'firefox',
      use: { 
        ...devices['Desktop Firefox'],
        viewport: { width: 1280, height: 720 }
      },
    },

    {
      name: 'webkit',
      use: { 
        ...devices['Desktop Safari'],
        viewport: { width: 1280, height: 720 }
      },
    },

    /* Test against mobile viewports. */
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] },
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone 12'] },
    },

    /* Test against branded browsers. */
    {
      name: 'Microsoft Edge',
      use: { ...devices['Desktop Edge'], channel: 'msedge' },
    },
    {
      name: 'Google Chrome',
      use: { ...devices['Desktop Chrome'], channel: 'chrome' },
    },
  ],

  /* Run your local dev server before starting the tests */
  // webServer: [
  //   {
  //     command: 'cd ../../../ && ./fire-salamander',
  //     url: 'http://localhost:8080/api/v1/health',
  //     reuseExistingServer: !process.env.CI,
  //     timeout: 30000,
  //   },
  //   {
  //     command: 'cd ../../../fire-salamander-frontend && npm run dev',
  //     url: 'http://localhost:3000',
  //     reuseExistingServer: !process.env.CI,
  //     timeout: 30000,
  //   }
  // ],

  /* Global setup and teardown */
  // globalSetup: require.resolve('./global-setup.js'),
  // globalTeardown: require.resolve('./global-teardown.js'),

  /* Test timeout */
  timeout: 30000,
  
  /* Expect timeout */
  expect: {
    timeout: 5000
  }
});