// Test CSS Fire Salamander selon standards NO HARDCODING
const { test, expect } = require('@playwright/test');

test.describe('CSS Tailwind Validation', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('http://localhost:8080');
  });

  test('Tailwind CSS should load and apply styles', async ({ page }) => {
    // Vérifier que Tailwind est chargé
    const tailwindLoaded = await page.evaluate(() => {
      return window.tailwind !== undefined;
    });
    expect(tailwindLoaded).toBe(true);

    // Vérifier que les classes SEPTEO sont appliquées
    const bgColor = await page.evaluate(() => {
      const body = document.querySelector('body');
      return window.getComputedStyle(body).backgroundColor;
    });
    
    // bg-septeo-gray-50 = #f8fafc = rgb(248, 250, 252)
    expect(bgColor).toBe('rgb(248, 250, 252)');

    // Vérifier le header
    const header = await page.locator('header');
    const headerBg = await header.evaluate(el => 
      window.getComputedStyle(el).backgroundColor
    );
    
    // bg-white = rgb(255, 255, 255)
    expect(headerBg).toBe('rgb(255, 255, 255)');

    // Vérifier SEPTEO colors
    const logoDiv = await page.locator('.bg-septeo-blue');
    const logoBg = await logoDiv.evaluate(el => 
      window.getComputedStyle(el).backgroundColor
    );
    
    // bg-septeo-blue = #1e3a8a = rgb(30, 58, 138)
    expect(logoBg).toBe('rgb(30, 58, 138)');
  });

  test('Font Awesome icons should load', async ({ page }) => {
    // Vérifier que Font Awesome est chargé
    const faLoaded = await page.evaluate(() => {
      return window.FontAwesomeConfig !== undefined;
    });
    expect(faLoaded).toBe(true);

    // Vérifier l'icône fire
    const fireIcon = await page.locator('.fa-fire');
    await expect(fireIcon).toBeVisible();
  });

  test('Responsive design should work', async ({ page }) => {
    // Desktop
    await page.setViewportSize({ width: 1280, height: 720 });
    const container = await page.locator('.container');
    await expect(container).toBeVisible();

    // Mobile
    await page.setViewportSize({ width: 375, height: 667 });
    await expect(container).toBeVisible();
  });
});