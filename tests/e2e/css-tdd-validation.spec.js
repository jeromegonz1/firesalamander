// 🧪 TDD TESTS CSS - MISSION PM Phase 1
// Playwright E2E tests pour validation CSS SEPTEO
const { test, expect } = require('@playwright/test');

test.describe('🔥🦎 FIRE SALAMANDER CSS TDD - Mission PM', () => {

  test('CSS File Loading - FileServer Test', async ({ page }) => {
    // GIVEN - CSS file should be served by FileServer
    const cssResponse = await page.request.get('/static/css/fire-salamander.css');
    
    // THEN - CSS file should be accessible
    expect(cssResponse.status()).toBe(200);
    expect(cssResponse.headers()['content-type']).toMatch(/text\/(css|plain)/);
    
    const cssContent = await cssResponse.text();
    
    // THEN - CSS should contain SEPTEO variables
    expect(cssContent).toContain('--septeo-orange: #ff6136');
    expect(cssContent).toContain('--septeo-blue: #1e3a8a');
    expect(cssContent).toContain('--septeo-gray-50: #f8fafc');
    
    // THEN - CSS should contain utility classes
    expect(cssContent).toContain('bg-septeo-orange');
    expect(cssContent).toContain('text-septeo-orange');
    expect(cssContent).toContain('bg-septeo-gray-50');
  });

  test('CSS Application - Visual Validation', async ({ page }) => {
    // GIVEN - Visit homepage
    await page.goto('/');
    
    // THEN - Check that CSS is linked in HTML
    const cssLink = page.locator('link[href="/static/css/fire-salamander.css"]');
    await expect(cssLink).toHaveAttribute('rel', 'stylesheet');
    
    // THEN - Verify SEPTEO orange color is applied to Fire icon
    const fireIcon = page.locator('i.fa-fire.text-septeo-orange');
    await expect(fireIcon).toBeVisible();
    
    // THEN - Verify computed styles for SEPTEO orange
    const fireIconColor = await fireIcon.evaluate(el => {
      return window.getComputedStyle(el).getPropertyValue('color');
    });
    
    // Convert SEPTEO orange #ff6136 to RGB for comparison
    // #ff6136 = rgb(255, 97, 54)
    expect(fireIconColor).toBe('rgb(255, 97, 54)');
  });

  test('CSS SEPTEO Background Colors Applied', async ({ page }) => {
    await page.goto('/');
    
    // THEN - Check SEPTEO blue background on brand element
    const brandElement = page.locator('.bg-septeo-blue').first();
    await expect(brandElement).toBeVisible();
    
    const backgroundColor = await brandElement.evaluate(el => {
      return window.getComputedStyle(el).getPropertyValue('background-color');
    });
    
    // #1e3a8a = rgb(30, 58, 138)
    expect(backgroundColor).toBe('rgb(30, 58, 138)');
  });

  test('CSS No CDN - Local Only Policy', async ({ page }) => {
    // GIVEN - Start monitoring network requests
    const cdnRequests = [];
    page.on('request', request => {
      if (request.url().includes('cdn.tailwindcss.com') || 
          request.url().includes('tailwind.com') ||
          request.url().includes('unpkg.com/tailwindcss')) {
        cdnRequests.push(request.url());
      }
    });
    
    // WHEN - Load the homepage
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    
    // THEN - No Tailwind CDN requests should be made
    expect(cdnRequests).toHaveLength(0);
    
    // THEN - Verify styles are still working (from local CSS)
    const fireIcon = page.locator('i.fa-fire.text-septeo-orange');
    await expect(fireIcon).toBeVisible();
  });

  test('CSS Loading - Zero CORS Errors', async ({ page }) => {
    const consoleErrors = [];
    page.on('console', msg => {
      if (msg.type() === 'error' && 
          (msg.text().includes('CORS') || 
           msg.text().includes('Access-Control-Allow-Origin') ||
           msg.text().includes('Cross-Origin'))) {
        consoleErrors.push(msg.text());
      }
    });

    // WHEN - Load all main pages
    await page.goto('/');
    await page.goto('/analyze?url=https://example.com');
    await page.goto('/results?url=https://example.com');
    
    // THEN - Zero CORS errors should occur
    expect(consoleErrors).toHaveLength(0);
  });

  test('CSS Visual Regression - SEPTEO Styling Screenshot', async ({ page }) => {
    await page.goto('/');
    
    // Wait for fonts and styles to load
    await page.waitForLoadState('networkidle');
    await page.waitForTimeout(1000);
    
    // Take screenshot for visual comparison
    await page.screenshot({ 
      path: 'tests/screenshots/css-tdd-septeo-styles.png', 
      fullPage: true 
    });
    
    // Verify key visual elements are styled correctly
    const header = page.locator('header');
    await expect(header).toBeVisible();
    
    const brandIcon = page.locator('.bg-septeo-blue');
    await expect(brandIcon).toBeVisible();
    
    const orangeElements = await page.locator('.text-septeo-orange').count();
    expect(orangeElements).toBeGreaterThan(0);
  });

});