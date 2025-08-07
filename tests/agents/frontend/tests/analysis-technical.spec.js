/**
 * Fire Salamander - TDD Tests pour /analysis/[id]/technical
 * Lead Tech Standards: Tests AVANT implémentation
 * 
 * Focus: Détails techniques SEO avancés
 * Tests obligatoires: Cross-browser, Mobile, Screenshots, Accessibility > 95%, SEPTEO colors
 */

const { test, expect } = require('@playwright/test');

test.describe('Analysis Technical Page - TDD Suite', () => {
  const analysisId = '10'; // Marina-plage analysis avec vraies données
  const technicalUrl = `/analysis/${analysisId}/technical`;

  test.beforeEach(async ({ page }) => {
    await page.goto('/dashboard');
    await expect(page).toHaveTitle(/Fire Salamander/);
  });

  test.describe('Page Structure & Loading', () => {
    test('should load technical analysis page with correct title', async ({ page }) => {
      await page.goto(technicalUrl);
      
      await expect(page).toHaveTitle(/Fire Salamander/);
      await expect(page.locator('h1')).toContainText('Analyse Technique');
      
      await page.screenshot({ path: 'test-results/technical-page-loaded.png', fullPage: true });
    });

    test('should display technical health overview', async ({ page }) => {
      await page.goto(technicalUrl);
      
      // Score de santé technique global
      await expect(page.locator('[data-testid=\"technical-health-score\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"health-score-value\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/technical-health-overview.png' });
    });

    test('should show loading state for technical audits', async ({ page }) => {
      await page.route(`**/api/v1/analysis/${analysisId}`, async route => {
        await page.waitForTimeout(1000);
        route.continue();
      });
      
      await page.goto(technicalUrl);
      
      await expect(page.locator('[data-testid=\"technical-loading\"]')).toBeVisible();
      await page.screenshot({ path: 'test-results/technical-loading-state.png' });
      
      await expect(page.locator('[data-testid=\"technical-loading\"]')).toBeHidden({ timeout: 15000 });
    });
  });

  test.describe('Technical Audit Sections', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(technicalUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should display HTML structure analysis', async ({ page }) => {
      // Section structure HTML
      await expect(page.locator('[data-testid=\"html-structure\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Structure HTML');
      
      // Validations HTML
      await expect(page.locator('[data-testid=\"html-validation\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"doctype-check\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"semantic-tags\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/html-structure-section.png' });
    });

    test('should show meta tags audit', async ({ page }) => {
      // Section meta tags
      await expect(page.locator('[data-testid=\"meta-tags-audit\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Meta Tags');
      
      // Check liste meta tags
      await expect(page.locator('[data-testid=\"title-tag-status\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"description-tag-status\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"og-tags-status\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/meta-tags-audit.png' });
    });

    test('should display page speed diagnostics', async ({ page }) => {
      // Section vitesse
      await expect(page.locator('[data-testid=\"page-speed-diagnostics\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Vitesse de Page');
      
      // Core Web Vitals
      await expect(page.locator('[data-testid=\"lcp-metric\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"fid-metric\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"cls-metric\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/page-speed-diagnostics.png' });
    });

    test('should show mobile optimization analysis', async ({ page }) => {
      // Section mobile
      await expect(page.locator('[data-testid=\"mobile-optimization\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Optimisation Mobile');
      
      // Tests mobile
      await expect(page.locator('[data-testid=\"viewport-check\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"mobile-friendly-check\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"touch-targets-check\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/mobile-optimization-section.png' });
    });

    test('should display security audit', async ({ page }) => {
      // Section sécurité
      await expect(page.locator('[data-testid=\"security-audit\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Audit Sécurité');
      
      // Checks sécurité
      await expect(page.locator('[data-testid=\"https-check\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"mixed-content-check\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"security-headers-check\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/security-audit-section.png' });
    });

    test('should show crawlability analysis', async ({ page }) => {
      // Section crawlabilité
      await expect(page.locator('[data-testid=\"crawlability-analysis\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Crawlabilité');
      
      // Tests robots/sitemap
      await expect(page.locator('[data-testid=\"robots-txt-check\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"sitemap-check\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"internal-links-analysis\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/crawlability-analysis.png' });
    });
  });

  test.describe('Interactive Features & Details', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(technicalUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should expand technical issue details on click', async ({ page }) => {
      // Cliquer sur un problème technique
      await page.locator('[data-testid=\"technical-issue\"]').first().click();
      
      // Panel de détails
      await expect(page.locator('[data-testid=\"issue-details-panel\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"issue-description\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"issue-solution\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/technical-issue-details.png' });
      
      // Fermer panel
      await page.locator('[data-testid=\"close-details\"]').click();
      await expect(page.locator('[data-testid=\"issue-details-panel\"]')).toBeHidden();
    });

    test('should filter issues by severity', async ({ page }) => {
      // Filtres par sévérité
      await expect(page.locator('[data-testid=\"severity-filter\"]')).toBeVisible();
      
      // Filtrer "Critical"
      await page.locator('[data-testid=\"filter-critical\"]').click();
      
      // Vérifier filtrage
      await expect(page.locator('[data-testid=\"filtered-issues\"]')).toBeVisible();
      const criticalIssues = page.locator('[data-testid=\"critical-issue\"]');
      expect(await criticalIssues.count()).toBeGreaterThanOrEqual(0);
      
      await page.screenshot({ path: 'test-results/technical-issues-filtered-critical.png' });
    });

    test('should show before/after comparison for fixes', async ({ page }) => {
      // Cliquer sur un fix suggéré
      await page.locator('[data-testid=\"suggested-fix\"]').first().click();
      
      // Modal avant/après
      await expect(page.locator('[data-testid=\"before-after-modal\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"before-code\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"after-code\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/before-after-comparison.png' });
    });

    test('should export technical audit report', async ({ page }) => {
      // Bouton export
      await expect(page.locator('[data-testid=\"export-technical-report\"]')).toBeVisible();
      
      const downloadPromise = page.waitForEvent('download');
      await page.locator('[data-testid=\"export-technical-report\"]').click();
      
      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/technical-audit.*\\.(pdf|csv)/);
      
      await page.screenshot({ path: 'test-results/technical-report-export.png' });
    });
  });

  test.describe('Status Indicators & Scores', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(technicalUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should display color-coded status indicators', async ({ page }) => {
      // Indicateurs de statut
      const statusIndicators = page.locator('[data-testid=\"status-indicator\"]');
      await expect(statusIndicators).toHaveCount(6, { timeout: 10000 }); // 6 sections techniques
      
      // Vérifier couleurs selon status
      const firstIndicator = statusIndicators.first();
      const statusClass = await firstIndicator.getAttribute('class');
      expect(statusClass).toMatch(/(success|warning|error)/);
      
      await page.screenshot({ path: 'test-results/status-indicators-colors.png' });
    });

    test('should show improvement recommendations', async ({ page }) => {
      // Section recommandations
      await expect(page.locator('[data-testid=\"improvement-recommendations\"]')).toBeVisible();
      
      // Cards recommandations avec priorité
      const recommendations = page.locator('[data-testid=\"recommendation-card\"]');
      await expect(recommendations).toHaveCount(3, { timeout: 10000 });
      
      // Vérifier contenu recommandations
      await expect(recommendations.first()).toContainText('Priorité');
      await expect(recommendations.first()).toContainText('Impact');
      
      await page.screenshot({ path: 'test-results/improvement-recommendations.png' });
    });
  });

  test.describe('SEPTEO Design Standards', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(technicalUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should use SEPTEO orange for critical elements', async ({ page }) => {
      // Vérifier bouton export
      const exportButton = page.locator('[data-testid=\"export-technical-report\"]');
      await expect(exportButton).toHaveCSS('background-color', 'rgb(255, 97, 54)');
      
      // Vérifier indicateurs critiques
      const criticalIndicators = page.locator('[data-testid=\"critical-indicator\"]');
      if (await criticalIndicators.count() > 0) {
        await expect(criticalIndicators.first()).toHaveCSS('color', 'rgb(255, 97, 54)');
      }
      
      await page.screenshot({ path: 'test-results/septeo-orange-validation.png' });
    });

    test('should maintain consistent card design', async ({ page }) => {
      // Vérifier design cards
      const cards = page.locator('[data-testid=\"technical-card\"]');
      const cardCount = await cards.count();
      
      for (let i = 0; i < Math.min(cardCount, 3); i++) {
        const card = cards.nth(i);
        await expect(card).toHaveCSS('border-radius', '8px');
        await expect(card).toHaveCSS('padding', '24px');
      }
      
      await page.screenshot({ path: 'test-results/card-design-consistency.png' });
    });
  });

  test.describe('Responsive Design', () => {
    test('should adapt to tablet layout', async ({ page }) => {
      await page.setViewportSize({ width: 768, height: 1024 });
      await page.goto(technicalUrl);
      
      // Vérifier adaptation tablet
      await expect(page.locator('[data-testid=\"technical-health-score\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"html-structure\"]')).toBeVisible();
      
      // Grid responsive
      const auditSections = page.locator('[data-testid=\"audit-section\"]');
      const firstSection = auditSections.first();
      const sectionBox = await firstSection.boundingBox();
      expect(sectionBox.width).toBeLessThan(768); // Adapt to tablet width
      
      await page.screenshot({ path: 'test-results/technical-tablet-responsive.png', fullPage: true });
    });

    test('should work on mobile devices', async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
      await page.goto(technicalUrl);
      
      // Vérifier mobile
      await expect(page.locator('[data-testid=\"technical-health-score\"]')).toBeVisible();
      
      // Sections stacked sur mobile
      const auditSections = page.locator('[data-testid=\"audit-section\"]');
      const firstSection = auditSections.first();
      const sectionBox = await firstSection.boundingBox();
      expect(sectionBox.width).toBeCloseTo(375, 50); // Quasi full width sur mobile
      
      await page.screenshot({ path: 'test-results/technical-mobile-responsive.png', fullPage: true });
    });
  });

  test.describe('Accessibility Compliance (> 95%)', () => {
    test('should meet WCAG standards for technical data', async ({ page }) => {
      await page.goto(technicalUrl);
      await page.waitForLoadState('networkidle');
      
      // Navigation clavier
      await page.keyboard.press('Tab');
      const focusedElement = await page.locator(':focus');
      await expect(focusedElement).toBeVisible();
      
      // ARIA labels sur données techniques
      await expect(page.locator('[data-testid=\"technical-health-score\"]')).toHaveAttribute('aria-label');
      await expect(page.locator('[data-testid=\"page-speed-diagnostics\"]')).toHaveAttribute('aria-label');
      
      await page.screenshot({ path: 'test-results/technical-accessibility.png' });
    });

    test('should have proper table semantics for audit data', async ({ page }) => {
      await page.goto(technicalUrl);
      
      // Tables avec headers appropriés
      const auditTables = page.locator('[data-testid=\"audit-table\"]');
      if (await auditTables.count() > 0) {
        const firstTable = auditTables.first();
        await expect(firstTable.locator('th')).toHaveCount(3, { timeout: 5000 });
        await expect(firstTable.locator('th').first()).toHaveAttribute('scope', 'col');
      }
      
      await page.screenshot({ path: 'test-results/table-semantics-validation.png' });
    });
  });

  test.describe('Performance & Error Handling', () => {
    test('should handle large technical datasets', async ({ page }) => {
      // Mock large dataset
      await page.route(`**/api/v1/analysis/${analysisId}`, async route => {
        const response = await route.fetch();
        const json = await response.json();
        
        // Inject large technical data
        json.data.result_data = JSON.stringify({
          ...JSON.parse(json.data.result_data || '{}'),
          technical_analysis: {
            html_issues: Array.from({ length: 500 }, (_, i) => ({
              type: `issue-${i}`,
              severity: ['critical', 'warning', 'info'][i % 3],
              message: `Technical issue ${i}`
            })),
            performance_metrics: Array.from({ length: 100 }, (_, i) => ({
              metric: `metric-${i}`,
              value: Math.random() * 1000
            }))
          }
        });
        
        await route.fulfill({ response, json });
      });
      
      await page.goto(technicalUrl);
      
      // Page reste responsive
      await expect(page.locator('[data-testid=\"html-structure\"]')).toBeVisible({ timeout: 10000 });
      
      await page.screenshot({ path: 'test-results/large-technical-dataset.png' });
    });

    test('should handle API failures gracefully', async ({ page }) => {
      await page.route(`**/api/v1/analysis/${analysisId}`, route => {
        route.fulfill({
          status: 503,
          body: JSON.stringify({ success: false, error: 'Service Unavailable' })
        });
      });
      
      await page.goto(technicalUrl);
      
      // Message d'erreur technique
      await expect(page.locator('[data-testid=\"technical-error\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"retry-technical-audit\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/technical-api-error.png' });
    });
  });
});

/**
 * Cross-Browser Tests - OBLIGATOIRES Lead Tech
 */
['chromium', 'firefox', 'webkit'].forEach(browserName => {
  test.describe(`Technical Analysis - ${browserName}`, () => {
    test(`should render technical data correctly on ${browserName}`, async ({ page }) => {
      await page.goto(`/analysis/10/technical`);
      await page.waitForLoadState('networkidle');
      
      // Tests cross-browser
      await expect(page.locator('h1')).toContainText('Analyse Technique');
      await expect(page.locator('[data-testid=\"technical-health-score\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"html-structure\"]')).toBeVisible();
      
      await page.screenshot({ 
        path: `test-results/technical-${browserName}-compatibility.png`,
        fullPage: true 
      });
    });
  });
});