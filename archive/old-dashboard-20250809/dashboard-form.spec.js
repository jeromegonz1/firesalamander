/**
 * Fire Salamander - TDD Tests pour Dashboard avec formulaire URL
 * Lead Tech Standards: Tests AVANT implémentation selon Image #1
 * 
 * Focus: Formulaire d'analyse SEO intégré au dashboard
 * Tests obligatoires: URL input, bouton analyze, options scan, UX conditionnelle
 */

const { test, expect } = require('@playwright/test');

test.describe('Dashboard Form - TDD Suite (Image #1)', () => {
  const dashboardUrl = '/dashboard';

  test.beforeEach(async ({ page }) => {
    await page.goto(dashboardUrl);
    await expect(page).toHaveTitle(/Fire Salamander/);
  });

  test.describe('URL Analysis Form - Selon Image #1', () => {
    test('should display main analysis form when no recent analyses', async ({ page }) => {
      // Clear any existing analyses for clean test
      await page.evaluate(() => localStorage.clear());
      await page.reload();
      
      // Titre principal selon Image #1
      await expect(page.locator('h1')).toContainText('New SEO Analysis');
      
      // Sous-titre descriptif
      await expect(page.locator('[data-testid="analysis-subtitle"]'))
        .toContainText('Enter a website URL to start a comprehensive SEO audit');
      
      await page.screenshot({ path: 'test-results/dashboard-form-main-view.png' });
    });

    test('should have URL input field with correct placeholder', async ({ page }) => {
      // Champ URL principal selon Image #1
      const urlInput = page.locator('[data-testid="url-input"]');
      await expect(urlInput).toBeVisible();
      await expect(urlInput).toHaveAttribute('placeholder', 'Enter website URL (e.g. example.com)');
      await expect(urlInput).toHaveAttribute('type', 'url');
      
      await page.screenshot({ path: 'test-results/dashboard-url-input-field.png' });
    });

    test('should display analyze button with SEPTEO orange styling', async ({ page }) => {
      // Bouton Analyze Now selon Image #1
      const analyzeButton = page.locator('[data-testid="analyze-button"]');
      await expect(analyzeButton).toBeVisible();
      await expect(analyzeButton).toContainText('Analyze Now');
      
      // Vérifier style SEPTEO orange
      await expect(analyzeButton).toHaveCSS('background-color', 'rgb(255, 97, 54)');
      
      await page.screenshot({ path: 'test-results/dashboard-analyze-button.png' });
    });

    test('should display analysis options checkboxes', async ({ page }) => {
      // Options selon Image #1
      await expect(page.locator('[data-testid="advanced-options"]')).toBeVisible();
      await expect(page.locator('[data-testid="competitor-analysis"]')).toBeVisible();
      await expect(page.locator('[data-testid="ai-insights"]')).toBeVisible();
      
      // AI Insights doit être coché par défaut selon Image #1
      await expect(page.locator('[data-testid="ai-insights"]')).toBeChecked();
      
      await page.screenshot({ path: 'test-results/dashboard-analysis-options.png' });
    });

    test('should display security notice', async ({ page }) => {
      // Message sécurité selon Image #1
      await expect(page.locator('[data-testid="security-notice"]'))
        .toContainText('All analyses are private and secure');
      
      // Icône shield visible
      await expect(page.locator('[data-testid="security-icon"]')).toBeVisible();
    });
  });

  test.describe('Scan Type Selection - 3 Cards selon Image #1', () => {
    test('should display three scan type cards', async ({ page }) => {
      // 3 cartes de scan selon Image #1
      await expect(page.locator('[data-testid="comprehensive-scan"]')).toBeVisible();
      await expect(page.locator('[data-testid="quick-scan"]')).toBeVisible();
      await expect(page.locator('[data-testid="custom-scan"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/dashboard-scan-type-cards.png' });
    });

    test('should have comprehensive scan selected by default', async ({ page }) => {
      // Comprehensive sélectionné par défaut selon Image #1
      const comprehensiveCard = page.locator('[data-testid="comprehensive-scan"]');
      await expect(comprehensiveCard).toHaveClass(/selected|active|border-orange/);
      
      // Icône de sélection visible
      await expect(page.locator('[data-testid="comprehensive-selected-icon"]')).toBeVisible();
    });

    test('should display comprehensive scan features', async ({ page }) => {
      const comprehensiveCard = page.locator('[data-testid="comprehensive-scan"]');
      
      // Features selon Image #1
      await expect(comprehensiveCard).toContainText('Full website crawl (up to 500 pages)');
      await expect(comprehensiveCard).toContainText('Technical SEO analysis');
      await expect(comprehensiveCard).toContainText('Content quality assessment');
      await expect(comprehensiveCard).toContainText('Backlink profile analysis');
      await expect(comprehensiveCard).toContainText('AI-powered recommendations');
      await expect(comprehensiveCard).toContainText('Estimated time: 4-6 minutes');
    });

    test('should display quick scan features', async ({ page }) => {
      const quickCard = page.locator('[data-testid="quick-scan"]');
      
      // Features selon Image #1
      await expect(quickCard).toContainText('Home page analysis');
      await expect(quickCard).toContainText('Basic technical checks');
      await expect(quickCard).toContainText('Core Web Vitals');
      await expect(quickCard).toContainText('Top 5 issues identified');
      await expect(quickCard).toContainText('Estimated time: 1-2 minutes');
      
      // AI recommendations non disponible (texte barré)
      await expect(page.locator('[data-testid="quick-no-ai"]')).toBeVisible();
    });

    test('should display custom scan features', async ({ page }) => {
      const customCard = page.locator('[data-testid="custom-scan"]');
      
      // Features selon Image #1
      await expect(customCard).toContainText('Select specific pages to scan');
      await expect(customCard).toContainText('Choose analysis modules');
      await expect(customCard).toContainText('Set crawl depth and limits');
      await expect(customCard).toContainText('Configure scan frequency');
      await expect(customCard).toContainText('Optional AI insights');
      await expect(customCard).toContainText('Time varies based on settings');
    });

    test('should allow switching between scan types', async ({ page }) => {
      // Sélectionner Quick Scan
      await page.locator('[data-testid="quick-scan"]').click();
      
      // Vérifier changement de sélection
      await expect(page.locator('[data-testid="quick-scan"]')).toHaveClass(/selected|active|border-orange/);
      await expect(page.locator('[data-testid="comprehensive-scan"]')).not.toHaveClass(/selected|active|border-orange/);
      
      // Sélectionner Custom Scan
      await page.locator('[data-testid="custom-scan"]').click();
      await expect(page.locator('[data-testid="custom-scan"]')).toHaveClass(/selected|active|border-orange/);
      
      await page.screenshot({ path: 'test-results/dashboard-scan-type-selection.png' });
    });
  });

  test.describe('Form Validation & Submission', () => {
    test('should validate URL format', async ({ page }) => {
      const urlInput = page.locator('[data-testid="url-input"]');
      const analyzeButton = page.locator('[data-testid="analyze-button"]');
      
      // URL invalide
      await urlInput.fill('invalid-url');
      await analyzeButton.click();
      
      await expect(page.locator('[data-testid="url-validation-error"]')).toBeVisible();
      await expect(page.locator('[data-testid="url-validation-error"]'))
        .toContainText('Please enter a valid URL');
    });

    test('should accept valid URLs', async ({ page }) => {
      const urlInput = page.locator('[data-testid="url-input"]');
      
      // URLs valides
      const validUrls = [
        'https://example.com',
        'http://test-site.fr',
        'example.com',
        'subdomain.example.org'
      ];
      
      for (const url of validUrls) {
        await urlInput.fill(url);
        await expect(page.locator('[data-testid="url-validation-error"]')).not.toBeVisible();
      }
    });

    test('should start analysis with valid URL', async ({ page }) => {
      const urlInput = page.locator('[data-testid="url-input"]');
      const analyzeButton = page.locator('[data-testid="analyze-button"]');
      
      // Saisir URL valide
      await urlInput.fill('https://test-website.com');
      
      // Cliquer analyser
      await analyzeButton.click();
      
      // Vérifier redirection vers page de progression
      await expect(page).toHaveURL(/\/analysis\/\d+\/progress/);
      
      await page.screenshot({ path: 'test-results/dashboard-analysis-started.png' });
    });
  });

  test.describe('Conditional Display Logic', () => {
    test('should show form when no analyses exist', async ({ page }) => {
      // Mock API pour retourner 0 analyses
      await page.route('**/api/v1/analyses', async route => {
        route.fulfill({
          json: { success: true, data: [] }
        });
      });
      
      await page.reload();
      
      // Formulaire visible
      await expect(page.locator('[data-testid="analysis-form"]')).toBeVisible();
      
      // KPIs cachés
      await expect(page.locator('[data-testid="kpi-cards"]')).not.toBeVisible();
    });

    test('should show metrics when analyses exist', async ({ page }) => {
      // Mock API avec analyses existantes
      await page.route('**/api/v1/analyses', async route => {
        route.fulfill({
          json: { 
            success: true, 
            data: [
              { id: 1, url: 'test.com', score: 78, status: 'success' },
              { id: 2, url: 'example.com', score: 65, status: 'success' }
            ]
          }
        });
      });
      
      await page.reload();
      
      // KPIs visibles
      await expect(page.locator('[data-testid="kpi-cards"]')).toBeVisible();
      
      // Formulaire toujours visible mais moins prominent
      await expect(page.locator('[data-testid="analysis-form"]')).toBeVisible();
    });
  });

  test.describe('SEPTEO Design Compliance', () => {
    test('should use SEPTEO orange for primary elements', async ({ page }) => {
      // Bouton principal
      const analyzeButton = page.locator('[data-testid="analyze-button"]');
      await expect(analyzeButton).toHaveCSS('background-color', 'rgb(255, 97, 54)');
      
      // Carte sélectionnée
      const selectedCard = page.locator('[data-testid="comprehensive-scan"]');
      await expect(selectedCard).toHaveCSS('border-color', 'rgb(255, 97, 54)');
    });

    test('should be responsive on mobile', async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
      
      // Cards en stack vertical
      const cards = page.locator('[data-testid*="scan"]');
      const cardCount = await cards.count();
      
      for (let i = 0; i < cardCount; i++) {
        await expect(cards.nth(i)).toBeVisible();
      }
      
      // URL input full width
      const urlInput = page.locator('[data-testid="url-input"]');
      const inputBox = await urlInput.boundingBox();
      expect(inputBox.width).toBeGreaterThan(300);
      
      await page.screenshot({ path: 'test-results/dashboard-mobile-responsive.png', fullPage: true });
    });
  });

  test.describe('Accessibility Standards', () => {
    test('should have proper form labels', async ({ page }) => {
      // Label pour URL input
      await expect(page.locator('label[for="url-input"]')).toBeVisible();
      
      // Labels pour checkboxes
      await expect(page.locator('label[for="advanced-options"]')).toBeVisible();
      await expect(page.locator('label[for="competitor-analysis"]')).toBeVisible();
      await expect(page.locator('label[for="ai-insights"]')).toBeVisible();
    });

    test('should support keyboard navigation', async ({ page }) => {
      // Navigation par Tab
      await page.keyboard.press('Tab'); // URL input
      await expect(page.locator('[data-testid="url-input"]')).toBeFocused();
      
      await page.keyboard.press('Tab'); // Advanced options
      await expect(page.locator('[data-testid="advanced-options"]')).toBeFocused();
      
      await page.keyboard.press('Tab'); // Competitor analysis
      await expect(page.locator('[data-testid="competitor-analysis"]')).toBeFocused();
    });

    test('should announce form state changes', async ({ page }) => {
      // Région live pour les messages
      await expect(page.locator('[data-testid="form-feedback"]')).toHaveAttribute('aria-live', 'polite');
    });
  });
});

/**
 * Cross-Browser Tests - OBLIGATOIRES Lead Tech Standards
 */
['chromium', 'firefox', 'webkit'].forEach(browserName => {
  test.describe(`Dashboard Form - ${browserName}`, () => {
    test(`should render form correctly on ${browserName}`, async ({ page }) => {
      await page.goto('/dashboard');
      await page.waitForLoadState('networkidle');
      
      // Tests basiques cross-browser
      await expect(page.locator('[data-testid="url-input"]')).toBeVisible();
      await expect(page.locator('[data-testid="analyze-button"]')).toBeVisible();
      await expect(page.locator('[data-testid="comprehensive-scan"]')).toBeVisible();
      
      await page.screenshot({ 
        path: `test-results/dashboard-form-${browserName}-compatibility.png`,
        fullPage: true 
      });
    });
  });
});