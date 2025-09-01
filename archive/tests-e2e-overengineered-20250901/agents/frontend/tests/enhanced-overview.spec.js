/**
 * Fire Salamander - Enhanced Overview Section Tests
 * TDD Suite for Professional SEO Report Overview
 * 
 * Requirements: Display comprehensive data like SEMrush/Ahrefs
 */

const { test, expect } = require('@playwright/test');

test.describe('Enhanced Overview Section - TDD Suite', () => {
  const reportUrl = '/analysis/26/report';

  test.beforeEach(async ({ page }) => {
    await page.goto(reportUrl);
    await expect(page).toHaveTitle(/Rapport d'Analyse SEO/);
  });

  test.describe('Score Details Display', () => {
    test('should display global score with all sub-scores', async ({ page }) => {
      const overviewSection = page.locator('[data-testid="enhanced-overview"]');
      await expect(overviewSection).toBeVisible();
      
      // Global Score
      await expect(overviewSection.locator('[data-testid="global-score"]')).toBeVisible();
      await expect(overviewSection.locator('[data-testid="global-score-value"]')).toContainText('49');
      
      // SEO Score with details
      const seoDetails = overviewSection.locator('[data-testid="seo-score-details"]');
      await expect(seoDetails).toBeVisible();
      await expect(seoDetails.locator('[data-testid="title-optimization"]')).toContainText('60');
      await expect(seoDetails.locator('[data-testid="meta-optimization"]')).toContainText('90');
      await expect(seoDetails.locator('[data-testid="heading-structure"]')).toContainText('0');
      await expect(seoDetails.locator('[data-testid="keyword-usage"]')).toContainText('80');
      await expect(seoDetails.locator('[data-testid="content-quality"]')).toContainText('0');
    });

    test('should display technical score breakdown', async ({ page }) => {
      const techDetails = page.locator('[data-testid="technical-score-details"]');
      await expect(techDetails).toBeVisible();
      
      await expect(techDetails.locator('[data-testid="crawlability"]')).toContainText('100');
      await expect(techDetails.locator('[data-testid="indexability"]')).toContainText('100');
      await expect(techDetails.locator('[data-testid="site-speed"]')).toContainText('80');
      await expect(techDetails.locator('[data-testid="mobile"]')).toContainText('100');
      await expect(techDetails.locator('[data-testid="security"]')).toContainText('90');
    });

    test('should display performance metrics with Core Web Vitals', async ({ page }) => {
      const perfDetails = page.locator('[data-testid="performance-details"]');
      await expect(perfDetails).toBeVisible();
      
      // Core Web Vitals
      await expect(perfDetails.locator('[data-testid="lcp-value"]')).toContainText('0.52s');
      await expect(perfDetails.locator('[data-testid="lcp-status"]')).toHaveClass(/good/);
      
      await expect(perfDetails.locator('[data-testid="fid-value"]')).toContainText('100ms');
      await expect(perfDetails.locator('[data-testid="fid-status"]')).toHaveClass(/good/);
      
      await expect(perfDetails.locator('[data-testid="cls-value"]')).toContainText('0.15');
      await expect(perfDetails.locator('[data-testid="cls-status"]')).toHaveClass(/needs-improvement/);
      
      await expect(perfDetails.locator('[data-testid="ttfb-value"]')).toContainText('431ms');
      await expect(perfDetails.locator('[data-testid="ttfb-status"]')).toHaveClass(/poor/);
    });
  });

  test.describe('Global Metrics Display', () => {
    test('should display issue summary', async ({ page }) => {
      const metrics = page.locator('[data-testid="global-metrics"]');
      await expect(metrics).toBeVisible();
      
      // Issue counts
      await expect(metrics.locator('[data-testid="total-issues"]')).toContainText('6');
      await expect(metrics.locator('[data-testid="critical-issues"]')).toContainText('2');
      await expect(metrics.locator('[data-testid="warnings"]')).toContainText('3');
      await expect(metrics.locator('[data-testid="passed-checks"]')).toContainText('8');
    });

    test('should display resource metrics', async ({ page }) => {
      const metrics = page.locator('[data-testid="global-metrics"]');
      
      // Pages
      await expect(metrics.locator('[data-testid="pages-analyzed"]')).toContainText('1');
      
      // Images
      await expect(metrics.locator('[data-testid="total-images"]')).toContainText('29');
      
      // Links
      await expect(metrics.locator('[data-testid="internal-links"]')).toContainText('34');
      await expect(metrics.locator('[data-testid="external-links"]')).toContainText('2');
      await expect(metrics.locator('[data-testid="broken-links"]')).toContainText('0');
      
      // Load time
      await expect(metrics.locator('[data-testid="avg-load-time"]')).toContainText('431ms');
    });

    test('should display visual indicators for scores', async ({ page }) => {
      const overview = page.locator('[data-testid="enhanced-overview"]');
      
      // Progress bars or circular indicators
      await expect(overview.locator('[data-testid="global-score-indicator"]')).toBeVisible();
      await expect(overview.locator('[data-testid="seo-score-indicator"]')).toBeVisible();
      await expect(overview.locator('[data-testid="technical-score-indicator"]')).toBeVisible();
      await expect(overview.locator('[data-testid="performance-score-indicator"]')).toBeVisible();
    });
  });

  test.describe('Professional Layout', () => {
    test('should have SEMrush-style grid layout', async ({ page }) => {
      const overview = page.locator('[data-testid="enhanced-overview"]');
      
      // Check grid structure
      await expect(overview).toHaveCSS('display', 'grid');
      
      // Score cards in top row
      const scoreCards = overview.locator('[data-testid="score-card"]');
      expect(await scoreCards.count()).toBeGreaterThanOrEqual(4);
      
      // Detailed metrics below
      await expect(overview.locator('[data-testid="detailed-metrics"]')).toBeVisible();
    });

    test('should have color-coded status indicators', async ({ page }) => {
      // Good = green, Needs Improvement = yellow, Poor = red
      await expect(page.locator('[data-testid="status-good"]')).toHaveCSS('background-color', /rgb\(.*green.*/);
      await expect(page.locator('[data-testid="status-warning"]')).toHaveCSS('background-color', /rgb\(.*yellow.*/);
      await expect(page.locator('[data-testid="status-critical"]')).toHaveCSS('background-color', /rgb\(.*red.*/);
    });
  });

  test.describe('Data Accuracy', () => {
    test('should calculate scores from real backend data', async ({ page }) => {
      // Verify data matches backend response
      const response = await page.evaluate(async () => {
        const res = await fetch('http://localhost:8080/api/v1/analysis/26');
        return res.json();
      });
      
      const resultData = JSON.parse(response.data.result_data);
      const overview = page.locator('[data-testid="enhanced-overview"]');
      
      // Verify global score calculation
      const expectedGlobalScore = Math.round(response.data.overall_score * 100);
      await expect(overview.locator('[data-testid="global-score-value"]')).toContainText(expectedGlobalScore.toString());
    });
  });

  test.describe('Interactive Elements', () => {
    test('should have expandable details sections', async ({ page }) => {
      const seoDetails = page.locator('[data-testid="seo-expand-button"]');
      
      // Click to expand
      await seoDetails.click();
      await expect(page.locator('[data-testid="seo-expanded-content"]')).toBeVisible();
      
      // Click to collapse
      await seoDetails.click();
      await expect(page.locator('[data-testid="seo-expanded-content"]')).not.toBeVisible();
    });

    test('should have tooltips for metrics', async ({ page }) => {
      const metricWithTooltip = page.locator('[data-testid="lcp-metric"]');
      
      await metricWithTooltip.hover();
      await expect(page.locator('[data-testid="tooltip"]')).toBeVisible();
      await expect(page.locator('[data-testid="tooltip"]')).toContainText('Largest Contentful Paint');
    });
  });

  test.describe('Mobile Responsiveness', () => {
    test('should adapt layout on mobile', async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
      
      const overview = page.locator('[data-testid="enhanced-overview"]');
      await expect(overview).toBeVisible();
      
      // Cards should stack vertically
      const scoreCards = overview.locator('[data-testid="score-card"]');
      const firstCard = await scoreCards.first().boundingBox();
      const secondCard = await scoreCards.nth(1).boundingBox();
      
      expect(secondCard.y).toBeGreaterThan(firstCard.y + firstCard.height);
    });
  });
});

/**
 * Integration Tests
 */
test.describe('Integration with Backend', () => {
  test('should handle missing data gracefully', async ({ page }) => {
    // Mock response with missing data
    await page.route('**/api/v1/analysis/**', async route => {
      const response = await route.fetch();
      const json = await response.json();
      
      // Remove some data
      delete json.data.result_data;
      
      route.fulfill({ json });
    });
    
    await page.goto(reportUrl);
    
    // Should show fallback values
    await expect(page.locator('[data-testid="no-data-message"]')).toBeVisible();
  });

  test('should refresh data on demand', async ({ page }) => {
    await page.goto(reportUrl);
    
    const refreshButton = page.locator('[data-testid="refresh-analysis"]');
    await refreshButton.click();
    
    // Should show loading state
    await expect(page.locator('[data-testid="loading-overlay"]')).toBeVisible();
    
    // Should update data
    await expect(page.locator('[data-testid="loading-overlay"]')).not.toBeVisible({ timeout: 5000 });
  });
});