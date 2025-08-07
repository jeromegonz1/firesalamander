/**
 * Fire Salamander - TDD Tests pour /compare/[id1]/[id2]
 * Lead Tech Standards: Tests AVANT implémentation
 * 
 * Focus: Comparaison analyses SEO side-by-side
 * Tests obligatoires: Cross-browser, Mobile, Screenshots, Accessibility > 95%, SEPTEO colors
 */

const { test, expect } = require('@playwright/test');

test.describe('Analysis Compare Page - TDD Suite', () => {
  const analysisId1 = '10'; // Marina-plage (score 13.283)
  const analysisId2 = '14'; // Example.com (score 13.22875)  
  const compareUrl = `/compare/${analysisId1}/${analysisId2}`;

  test.beforeEach(async ({ page }) => {
    await page.goto('/dashboard');
    await expect(page).toHaveTitle(/Fire Salamander/);
  });

  test.describe('Page Structure & Loading', () => {
    test('should load comparison page with correct title', async ({ page }) => {
      await page.goto(compareUrl);
      
      await expect(page).toHaveTitle(/Fire Salamander/);
      await expect(page.locator('h1')).toContainText('Comparaison d\'Analyses');
      
      await page.screenshot({ path: 'test-results/compare-page-loaded.png', fullPage: true });
    });

    test('should display both analysis headers side by side', async ({ page }) => {
      await page.goto(compareUrl);
      
      // Headers des deux analyses
      await expect(page.locator('[data-testid=\"analysis-1-header\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"analysis-2-header\"]')).toBeVisible();
      
      // URLs des analyses
      await expect(page.locator('[data-testid=\"analysis-1-url\"]')).toContainText('marina-plage');
      await expect(page.locator('[data-testid=\"analysis-2-url\"]')).toContainText('example.com');
      
      await page.screenshot({ path: 'test-results/compare-headers-side-by-side.png' });
    });

    test('should show loading state for comparison data', async ({ page }) => {
      await page.route('**/api/v1/analysis/10', async route => {
        await page.waitForTimeout(1000);
        route.continue();
      });
      
      await page.goto(compareUrl);
      
      await expect(page.locator('[data-testid=\"comparison-loading\"]')).toBeVisible();
      await page.screenshot({ path: 'test-results/compare-loading-state.png' });
      
      await expect(page.locator('[data-testid=\"comparison-loading\"]')).toBeHidden({ timeout: 15000 });
    });
  });

  test.describe('Score Comparison Components', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(compareUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should display overall scores comparison', async ({ page }) => {
      // Section scores globaux
      await expect(page.locator('[data-testid=\"overall-scores-comparison\"]')).toBeVisible();
      
      // Scores côte à côte avec indicateur différence
      await expect(page.locator('[data-testid=\"score-1-value\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"score-2-value\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"score-difference\"]')).toBeVisible();
      
      // Winner indicator
      await expect(page.locator('[data-testid=\"score-winner\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/overall-scores-comparison.png' });
    });

    test('should show detailed metrics comparison', async ({ page }) => {
      // Section métriques détaillées
      await expect(page.locator('[data-testid=\"detailed-metrics-comparison\"]')).toBeVisible();
      
      // Métriques SEO, Technical, Performance, Content
      const metricRows = page.locator('[data-testid=\"metric-comparison-row\"]');
      await expect(metricRows).toHaveCount(4, { timeout: 10000 }); // 4 catégories principales
      
      // Chaque row a 3 colonnes: Métrique | Site 1 | Site 2
      const firstRow = metricRows.first();
      await expect(firstRow.locator('[data-testid=\"metric-name\"]')).toBeVisible();
      await expect(firstRow.locator('[data-testid=\"metric-value-1\"]')).toBeVisible();
      await expect(firstRow.locator('[data-testid=\"metric-value-2\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/detailed-metrics-comparison.png' });
    });

    test('should display visual chart comparison', async ({ page }) => {
      // Graphique comparatif (radar/bar chart)
      await expect(page.locator('[data-testid=\"comparison-chart\"]')).toBeVisible();
      
      // Chart with both datasets
      await expect(page.locator('[data-testid=\"chart-legend\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"dataset-1-legend\"]')).toContainText('marina-plage');
      await expect(page.locator('[data-testid=\"dataset-2-legend\"]')).toContainText('example.com');
      
      await page.screenshot({ path: 'test-results/visual-chart-comparison.png' });
    });
  });

  test.describe('Side-by-Side Analysis Sections', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(compareUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should compare SEO recommendations side by side', async ({ page }) => {
      // Section recommandations SEO
      await expect(page.locator('[data-testid=\"seo-recommendations-comparison\"]')).toBeVisible();
      
      // Colonnes recommandations
      await expect(page.locator('[data-testid=\"recommendations-site-1\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"recommendations-site-2\"]')).toBeVisible();
      
      // Nombre de recommandations par site
      await expect(page.locator('[data-testid=\"recommendations-count-1\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"recommendations-count-2\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/seo-recommendations-comparison.png' });
    });

    test('should compare technical issues', async ({ page }) => {
      // Section problèmes techniques
      await expect(page.locator('[data-testid=\"technical-issues-comparison\"]')).toBeVisible();
      
      // Issues par sévérité pour chaque site
      await expect(page.locator('[data-testid=\"critical-issues-1\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"critical-issues-2\"]')).toBeVisible();
      
      // Indicateurs visuels différence
      await expect(page.locator('[data-testid=\"issues-difference-indicator\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/technical-issues-comparison.png' });
    });

    test('should compare keywords performance', async ({ page }) => {
      // Section mots-clés
      await expect(page.locator('[data-testid=\"keywords-comparison\"]')).toBeVisible();
      
      // Top keywords per site
      await expect(page.locator('[data-testid=\"top-keywords-1\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"top-keywords-2\"]')).toBeVisible();
      
      // Keywords overlap analysis
      await expect(page.locator('[data-testid=\"keywords-overlap\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/keywords-comparison.png' });
    });

    test('should compare performance metrics', async ({ page }) => {
      // Section performance
      await expect(page.locator('[data-testid=\"performance-comparison\"]')).toBeVisible();
      
      // Core Web Vitals comparison
      await expect(page.locator('[data-testid=\"cwv-comparison\"]')).toBeVisible();
      
      // Performance scores with visual indicators
      const perfScores = page.locator('[data-testid=\"performance-score\"]');
      await expect(perfScores).toHaveCount(2);
      
      await page.screenshot({ path: 'test-results/performance-comparison.png' });
    });
  });

  test.describe('Interactive Comparison Features', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(compareUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should allow toggling comparison views', async ({ page }) => {
      // Toggle buttons: Overview | Detailed | Side-by-Side
      await expect(page.locator('[data-testid=\"view-toggle\"]')).toBeVisible();
      
      // Switch to detailed view
      await page.locator('[data-testid=\"toggle-detailed\"]').click();
      await expect(page.locator('[data-testid=\"detailed-comparison-view\"]')).toBeVisible();
      
      // Switch to side-by-side
      await page.locator('[data-testid=\"toggle-sidebyside\"]').click();
      await expect(page.locator('[data-testid=\"sidebyside-comparison-view\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/comparison-view-toggle.png' });
    });

    test('should highlight differences on hover', async ({ page }) => {
      // Hover sur une métrique
      const firstMetricRow = page.locator('[data-testid=\"metric-comparison-row\"]').first();
      await firstMetricRow.hover();
      
      // Vérifier highlight différence
      await expect(firstMetricRow).toHaveClass(/highlighted/);
      await expect(page.locator('[data-testid=\"difference-tooltip\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/difference-highlight-hover.png' });
    });

    test('should filter comparison by categories', async ({ page }) => {
      // Filtres catégories
      await expect(page.locator('[data-testid=\"category-filters\"]')).toBeVisible();
      
      // Filter SEO only
      await page.locator('[data-testid=\"filter-seo\"]').click();
      
      // Vérifier que seules les métriques SEO sont affichées
      await expect(page.locator('[data-testid=\"seo-metrics-only\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"technical-metrics\"]')).toBeHidden();
      
      await page.screenshot({ path: 'test-results/comparison-filtered-seo.png' });
    });

    test('should export comparison report', async ({ page }) => {
      // Bouton export comparaison
      await expect(page.locator('[data-testid=\"export-comparison\"]')).toBeVisible();
      
      const downloadPromise = page.waitForEvent('download');
      await page.locator('[data-testid=\"export-comparison\"]').click();
      
      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/comparison.*\\.pdf/);
      
      await page.screenshot({ path: 'test-results/comparison-export-action.png' });
    });
  });

  test.describe('Winner/Loser Indicators', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(compareUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should show clear winner indicators', async ({ page }) => {
      // Winner badge sur le score global
      await expect(page.locator('[data-testid=\"overall-winner-badge\"]')).toBeVisible();
      
      // Winner indicators par catégorie
      const categoryWinners = page.locator('[data-testid=\"category-winner\"]');
      await expect(categoryWinners.count()).toBeGreaterThan(0);
      
      // Vérifier couleurs winner/loser
      const winnerElement = categoryWinners.first();
      await expect(winnerElement).toHaveClass(/winner/);
      
      await page.screenshot({ path: 'test-results/winner-indicators.png' });
    });

    test('should display improvement opportunities', async ({ page }) => {
      // Section opportunités d'amélioration
      await expect(page.locator('[data-testid=\"improvement-opportunities\"]')).toBeVisible();
      
      // Suggestions basées sur le site qui performe mieux
      await expect(page.locator('[data-testid=\"improvement-suggestion\"]')).toHaveCount(3, { timeout: 10000 });
      
      // Chaque suggestion a un impact score
      const suggestions = page.locator('[data-testid=\"improvement-suggestion\"]');
      await expect(suggestions.first().locator('[data-testid=\"impact-score\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/improvement-opportunities.png' });
    });
  });

  test.describe('SEPTEO Design Standards', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(compareUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should use SEPTEO orange for primary actions', async ({ page }) => {
      // Bouton export
      const exportButton = page.locator('[data-testid=\"export-comparison\"]');
      await expect(exportButton).toHaveCSS('background-color', 'rgb(255, 97, 54)');
      
      // Winner badges
      const winnerBadges = page.locator('[data-testid=\"winner-badge\"]');
      if (await winnerBadges.count() > 0) {
        await expect(winnerBadges.first()).toHaveCSS('color', 'rgb(255, 97, 54)');
      }
      
      await page.screenshot({ path: 'test-results/septeo-orange-comparison.png' });
    });

    test('should maintain consistent comparison layout', async ({ page }) => {
      // Vérifier layout colonnes égales
      const column1 = page.locator('[data-testid=\"comparison-column-1\"]');
      const column2 = page.locator('[data-testid=\"comparison-column-2\"]');
      
      const col1Box = await column1.boundingBox();
      const col2Box = await column2.boundingBox();
      
      // Colonnes de largeur égale (± 10px)
      expect(Math.abs(col1Box.width - col2Box.width)).toBeLessThan(10);
      
      await page.screenshot({ path: 'test-results/comparison-layout-consistency.png' });
    });
  });

  test.describe('Responsive Design', () => {
    test('should stack columns on mobile', async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
      await page.goto(compareUrl);
      
      // Sur mobile, colonnes empilées verticalement
      await expect(page.locator('[data-testid=\"mobile-stacked-view\"]')).toBeVisible();
      
      // Headers des analyses restent visibles
      await expect(page.locator('[data-testid=\"analysis-1-header\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"analysis-2-header\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/comparison-mobile-stacked.png', fullPage: true });
    });

    test('should adapt comparison table on tablet', async ({ page }) => {
      await page.setViewportSize({ width: 768, height: 1024 });
      await page.goto(compareUrl);
      
      // Table responsive sur tablet
      await expect(page.locator('[data-testid=\"detailed-metrics-comparison\"]')).toBeVisible();
      
      // Colonnes adaptées
      const metricTable = page.locator('[data-testid=\"comparison-table\"]');
      const tableBox = await metricTable.boundingBox();
      expect(tableBox.width).toBeLessThan(768);
      
      await page.screenshot({ path: 'test-results/comparison-tablet-responsive.png', fullPage: true });
    });
  });

  test.describe('Accessibility Standards (> 95%)', () => {
    test('should support keyboard navigation', async ({ page }) => {
      await page.goto(compareUrl);
      await page.waitForLoadState('networkidle');
      
      // Tab navigation
      await page.keyboard.press('Tab');
      await expect(page.locator(':focus')).toBeVisible();
      
      // Navigation entre comparaisons
      await page.keyboard.press('ArrowRight');
      await page.keyboard.press('ArrowLeft');
      
      await page.screenshot({ path: 'test-results/comparison-keyboard-nav.png' });
    });

    test('should have proper ARIA labels for comparisons', async ({ page }) => {
      await page.goto(compareUrl);
      
      // ARIA labels sur éléments de comparaison
      await expect(page.locator('[data-testid=\"comparison-table\"]')).toHaveAttribute('aria-label');
      await expect(page.locator('[data-testid=\"overall-scores-comparison\"]')).toHaveAttribute('aria-label');
      
      // Role tables appropriés
      await expect(page.locator('[data-testid=\"comparison-table\"]')).toHaveAttribute('role', 'table');
      
      await page.screenshot({ path: 'test-results/comparison-aria-labels.png' });
    });

    test('should announce comparison results to screen readers', async ({ page }) => {
      await page.goto(compareUrl);
      
      // Live regions pour annonces dynamiques
      await expect(page.locator('[data-testid=\"comparison-results\"]')).toHaveAttribute('aria-live', 'polite');
      
      // Labels descriptifs sur données comparatives
      const scoreElements = page.locator('[data-testid=\"score-value\"]');
      if (await scoreElements.count() > 0) {
        await expect(scoreElements.first()).toHaveAttribute('aria-label');
      }
      
      await page.screenshot({ path: 'test-results/screen-reader-comparison.png' });
    });
  });

  test.describe('Error Handling & Edge Cases', () => {
    test('should handle invalid analysis IDs', async ({ page }) => {
      await page.goto('/compare/999999/888888'); // IDs inexistants
      
      // Message d'erreur approprié
      await expect(page.locator('[data-testid=\"comparison-error\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"invalid-analyses-message\"]')).toContainText('Analyses non trouvées');
      
      // Bouton retour dashboard
      await expect(page.locator('[data-testid=\"back-to-dashboard\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/invalid-analysis-ids-error.png' });
    });

    test('should handle same analysis comparison', async ({ page }) => {
      await page.goto('/compare/10/10'); // Même analyse
      
      // Message d'avertissement
      await expect(page.locator('[data-testid=\"same-analysis-warning\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"same-analysis-warning\"]')).toContainText('même analyse');
      
      // Suggestion autre analyse
      await expect(page.locator('[data-testid=\"suggest-other-analysis\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/same-analysis-comparison.png' });
    });

    test('should handle API failures gracefully', async ({ page }) => {
      await page.route('**/api/v1/analysis/**', route => {
        route.fulfill({
          status: 500,
          body: JSON.stringify({ success: false, error: 'Server Error' })
        });
      });
      
      await page.goto(compareUrl);
      
      // Message d'erreur API
      await expect(page.locator('[data-testid=\"api-error-message\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"retry-comparison\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/comparison-api-error.png' });
    });
  });
});

/**
 * Cross-Browser Tests - OBLIGATOIRES Lead Tech Standards
 */
['chromium', 'firefox', 'webkit'].forEach(browserName => {
  test.describe(`Comparison Page - ${browserName}`, () => {
    test(`should render comparison correctly on ${browserName}`, async ({ page }) => {
      await page.goto('/compare/10/14');
      await page.waitForLoadState('networkidle');
      
      // Tests basiques cross-browser
      await expect(page.locator('h1')).toContainText('Comparaison d\'Analyses');
      await expect(page.locator('[data-testid=\"overall-scores-comparison\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"detailed-metrics-comparison\"]')).toBeVisible();
      
      await page.screenshot({ 
        path: `test-results/comparison-${browserName}-compatibility.png`,
        fullPage: true 
      });
    });
  });
});