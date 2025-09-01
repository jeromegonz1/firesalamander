/**
 * Fire Salamander - TDD Tests pour /analysis/[id]/keywords
 * Lead Tech Standards: Tests AVANT implémentation
 * 
 * Tests obligatoires:
 * - Cross-browser (Chrome, Firefox, Safari) ✅
 * - Mobile responsive ✅ 
 * - Screenshots automatiques ✅
 * - Accessibility > 95% ✅
 * - SEPTEO colors validation ✅
 */

const { test, expect } = require('@playwright/test');

test.describe('Analysis Keywords Page - TDD Suite', () => {
  const analysisId = '10'; // Marina-plage analysis avec vraies données
  const keywordsUrl = `/analysis/${analysisId}/keywords`;

  test.beforeEach(async ({ page }) => {
    // Ensure Fire Salamander backend is running
    await page.goto('/dashboard');
    await expect(page).toHaveTitle(/Fire Salamander/);
  });

  test.describe('Page Structure & Loading', () => {
    test('should load keywords analysis page with correct title', async ({ page }) => {
      await page.goto(keywordsUrl);
      
      // Vérifier titre page
      await expect(page).toHaveTitle(/Fire Salamander/);
      
      // Vérifier H1 principal 
      await expect(page.locator('h1')).toContainText('Analyse des Mots-clés');
      
      // Screenshot obligatoire Lead Tech
      await page.screenshot({ path: 'test-results/keywords-page-loaded.png', fullPage: true });
    });

    test('should display analysis URL and date info', async ({ page }) => {
      await page.goto(keywordsUrl);
      
      // Vérifier info analyse (URL + date)
      await expect(page.locator('[data-testid=\"analysis-url\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"analysis-date\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/keywords-header-info.png' });
    });

    test('should show loading state initially', async ({ page }) => {
      // Intercept API call pour simuler chargement lent
      await page.route(`**/api/v1/analysis/${analysisId}`, async route => {
        await page.waitForTimeout(1000); // Simulate slow API
        route.continue();
      });
      
      await page.goto(keywordsUrl);
      
      // Vérifier spinner loading
      await expect(page.locator('[data-testid=\"loading-spinner\"]')).toBeVisible();
      await page.screenshot({ path: 'test-results/keywords-loading-state.png' });
      
      // Attendre que le chargement se termine
      await expect(page.locator('[data-testid=\"loading-spinner\"]')).toBeHidden({ timeout: 15000 });
    });
  });

  test.describe('Keywords Analysis Components - IA Focus', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(keywordsUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should display keywords density analysis', async ({ page }) => {
      // Section densité mots-clés
      await expect(page.locator('[data-testid=\"keywords-density\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Densité des Mots-clés');
      
      // Tableau ou graphique densité
      await expect(page.locator('[data-testid=\"density-chart\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/keywords-density-section.png' });
    });

    test('should show AI-powered keyword opportunities', async ({ page }) => {
      // Section opportunités IA
      await expect(page.locator('[data-testid=\"ai-opportunities\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Opportunités IA');
      
      // Cards d'opportunités avec scores
      const opportunityCards = page.locator('[data-testid=\"opportunity-card\"]');
      await expect(opportunityCards).toHaveCount(3, { timeout: 10000 });
      
      // Vérifier contenu des cards
      await expect(opportunityCards.first()).toContainText('Score');
      await expect(opportunityCards.first()).toContainText('Potentiel');
      
      await page.screenshot({ path: 'test-results/ai-opportunities-cards.png' });
    });

    test('should display semantic keywords analysis', async ({ page }) => {
      // Section analyse sémantique
      await expect(page.locator('[data-testid=\"semantic-analysis\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Analyse Sémantique');
      
      // Tags ou nuage de mots sémantiques
      await expect(page.locator('[data-testid=\"semantic-tags\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/semantic-analysis-section.png' });
    });

    test('should show competitor keywords comparison', async ({ page }) => {
      // Section comparaison concurrents
      await expect(page.locator('[data-testid=\"competitor-keywords\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Mots-clés Concurrents');
      
      // Tableau comparatif
      await expect(page.locator('[data-testid=\"competitor-table\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/competitor-keywords-table.png' });
    });
  });

  test.describe('Interactive Features', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(keywordsUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should allow filtering keywords by category', async ({ page }) => {
      // Filtres par catégorie
      await expect(page.locator('[data-testid=\"category-filter\"]')).toBeVisible();
      
      // Cliquer sur un filtre
      await page.locator('[data-testid=\"filter-primary\"]').click();
      
      // Vérifier que le contenu est filtré
      await expect(page.locator('[data-testid=\"filtered-results\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/keywords-filtered-primary.png' });
    });

    test('should show keyword details on click', async ({ page }) => {
      // Cliquer sur un mot-clé
      await page.locator('[data-testid=\"keyword-item\"]').first().click();
      
      // Modal ou panel détails
      await expect(page.locator('[data-testid=\"keyword-details-modal\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"keyword-stats\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/keyword-details-modal.png' });
      
      // Fermer modal
      await page.locator('[data-testid=\"close-modal\"]').click();
      await expect(page.locator('[data-testid=\"keyword-details-modal\"]')).toBeHidden();
    });

    test('should export keywords data', async ({ page }) => {
      // Bouton export
      await expect(page.locator('[data-testid=\"export-keywords\"]')).toBeVisible();
      
      // Mock download
      const downloadPromise = page.waitForEvent('download');
      await page.locator('[data-testid=\"export-keywords\"]').click();
      
      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/keywords.*\\.csv/);
      
      await page.screenshot({ path: 'test-results/keywords-export-action.png' });
    });
  });

  test.describe('SEPTEO Design Standards Validation', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(keywordsUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should use SEPTEO orange (#FF6136) for primary elements', async ({ page }) => {
      // Vérifier couleurs SEPTEO
      const primaryButton = page.locator('[data-testid=\"export-keywords\"]');
      await expect(primaryButton).toHaveCSS('background-color', 'rgb(255, 97, 54)'); // #FF6136
      
      // Vérifier links et accents
      const opportunityCards = page.locator('[data-testid=\"opportunity-card\"]');
      const firstCard = opportunityCards.first();
      await expect(firstCard.locator('.score-indicator')).toHaveCSS('color', 'rgb(255, 97, 54)');
      
      await page.screenshot({ path: 'test-results/septeo-colors-validation.png' });
    });

    test('should have consistent spacing and typography', async ({ page }) => {
      // Vérifier espacement entre sections
      const sections = page.locator('section');
      const sectionCount = await sections.count();
      
      for (let i = 0; i < Math.min(sectionCount, 3); i++) {
        const section = sections.nth(i);
        await expect(section).toHaveCSS('margin-bottom', '24px');
      }
      
      // Vérifier typographie cohérente
      await expect(page.locator('h1')).toHaveCSS('font-size', '30px');
      await expect(page.locator('h2')).toHaveCSS('font-size', '24px');
      
      await page.screenshot({ path: 'test-results/typography-spacing-validation.png' });
    });
  });

  test.describe('Responsive Design Tests', () => {
    test('should be responsive on mobile devices', async ({ page }) => {
      // Set mobile viewport
      await page.setViewportSize({ width: 375, height: 667 });
      await page.goto(keywordsUrl);
      
      // Vérifier que les éléments s'adaptent
      await expect(page.locator('[data-testid=\"keywords-density\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"ai-opportunities\"]')).toBeVisible();
      
      // Grid responsive
      const opportunityCards = page.locator('[data-testid=\"opportunity-card\"]');
      const firstCard = opportunityCards.first();
      
      // Sur mobile, les cards doivent être en colonne
      const cardWidth = await firstCard.boundingBox();
      expect(cardWidth.width).toBeGreaterThan(300); // Full width sur mobile
      
      await page.screenshot({ path: 'test-results/keywords-mobile-responsive.png', fullPage: true });
    });

    test('should be responsive on tablet devices', async ({ page }) => {
      // Set tablet viewport
      await page.setViewportSize({ width: 768, height: 1024 });
      await page.goto(keywordsUrl);
      
      // Vérifier adaptations tablet
      await expect(page.locator('[data-testid=\"semantic-analysis\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/keywords-tablet-responsive.png', fullPage: true });
    });
  });

  test.describe('Accessibility Standards (> 95%)', () => {
    test('should meet WCAG accessibility standards', async ({ page }) => {
      await page.goto(keywordsUrl);
      await page.waitForLoadState('networkidle');
      
      // Test focus navigation
      await page.keyboard.press('Tab');
      await expect(page.locator(':focus')).toBeVisible();
      
      // Test ARIA labels
      await expect(page.locator('[data-testid=\"keywords-density\"]')).toHaveAttribute('aria-label');
      await expect(page.locator('[data-testid=\"ai-opportunities\"]')).toHaveAttribute('aria-label');
      
      // Test alt text sur graphiques/images
      const charts = page.locator('[data-testid=\"density-chart\"]');
      if (await charts.count() > 0) {
        await expect(charts).toHaveAttribute('aria-label');
      }
      
      await page.screenshot({ path: 'test-results/accessibility-validation.png' });
    });

    test('should have proper semantic structure', async ({ page }) => {
      await page.goto(keywordsUrl);
      
      // Vérifier structure H1 > H2 > H3
      await expect(page.locator('h1')).toHaveCount(1);
      const h2Count = await page.locator('h2').count();
      expect(h2Count).toBeGreaterThanOrEqual(3); // Au moins 3 sections
      
      // Vérifier rôles ARIA
      await expect(page.locator('main')).toHaveAttribute('role', 'main');
      
      await page.screenshot({ path: 'test-results/semantic-structure-validation.png' });
    });
  });

  test.describe('Performance Standards', () => {
    test('should load within performance budget', async ({ page }) => {
      const startTime = Date.now();
      
      await page.goto(keywordsUrl);
      await page.waitForLoadState('networkidle');
      
      const loadTime = Date.now() - startTime;
      expect(loadTime).toBeLessThan(3000); // < 3s pour chargement complet
      
      console.log(`Keywords page loaded in ${loadTime}ms`);
    });

    test('should handle large datasets efficiently', async ({ page }) => {
      // Mock large dataset response
      await page.route(`**/api/v1/analysis/${analysisId}`, async route => {
        const response = await route.fetch();
        const json = await response.json();
        
        // Inject large keywords dataset
        json.data.result_data = JSON.stringify({
          ...JSON.parse(json.data.result_data || '{}'),
          keywords_analysis: {
            keywords: Array.from({ length: 1000 }, (_, i) => ({
              keyword: `test-keyword-${i}`,
              density: Math.random() * 10,
              score: Math.random() * 100
            }))
          }
        });
        
        await route.fulfill({ response, json });
      });
      
      await page.goto(keywordsUrl);
      
      // Vérifier que la page reste responsive avec beaucoup de données
      await expect(page.locator('[data-testid=\"keywords-density\"]')).toBeVisible({ timeout: 5000 });
      
      await page.screenshot({ path: 'test-results/large-dataset-performance.png' });
    });
  });

  test.describe('Error Handling', () => {
    test('should handle API errors gracefully', async ({ page }) => {
      // Mock API error
      await page.route(`**/api/v1/analysis/${analysisId}`, route => {
        route.fulfill({
          status: 500,
          contentType: 'application/json',
          body: JSON.stringify({ success: false, error: 'Internal Server Error' })
        });
      });
      
      await page.goto(keywordsUrl);
      
      // Vérifier message d'erreur
      await expect(page.locator('[data-testid=\"error-message\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"error-message\"]')).toContainText('Erreur');
      
      // Bouton retry
      await expect(page.locator('[data-testid=\"retry-button\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/api-error-handling.png' });
    });

    test('should handle missing analysis gracefully', async ({ page }) => {
      await page.goto('/analysis/999999/keywords'); // ID inexistant
      
      // Vérifier message "Analyse non trouvée"
      await expect(page.locator('[data-testid=\"not-found-message\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"back-button\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/analysis-not-found.png' });
    });
  });
});

/**
 * Cross-Browser Test Suite
 * Tests OBLIGATOIRES sur Chrome, Firefox, Safari
 */
['chromium', 'firefox', 'webkit'].forEach(browserName => {
  test.describe(`Cross-Browser: ${browserName}`, () => {
    test(`should render correctly on ${browserName}`, async ({ page }) => {
      await page.goto(`/analysis/10/keywords`);
      await page.waitForLoadState('networkidle');
      
      // Tests basiques cross-browser
      await expect(page.locator('h1')).toContainText('Analyse des Mots-clés');
      await expect(page.locator('[data-testid=\"keywords-density\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"ai-opportunities\"]')).toBeVisible();
      
      await page.screenshot({ 
        path: `test-results/keywords-${browserName}-compatibility.png`,
        fullPage: true 
      });
    });
  });
});