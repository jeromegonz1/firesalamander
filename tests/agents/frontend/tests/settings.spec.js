/**
 * Fire Salamander - TDD Tests pour /settings
 * Lead Tech Standards: Tests AVANT implémentation
 * 
 * Focus: Configuration utilisateur et préférences
 * Tests obligatoires: Cross-browser, Mobile, Screenshots, Accessibility > 95%, SEPTEO colors
 */

const { test, expect } = require('@playwright/test');

test.describe('Settings Page - TDD Suite', () => {
  const settingsUrl = '/settings';

  test.beforeEach(async ({ page }) => {
    await page.goto('/dashboard');
    await expect(page).toHaveTitle(/Fire Salamander/);
  });

  test.describe('Page Structure & Loading', () => {
    test('should load settings page with correct title', async ({ page }) => {
      await page.goto(settingsUrl);
      
      await expect(page).toHaveTitle(/Fire Salamander/);
      await expect(page.locator('h1')).toContainText('Paramètres');
      
      await page.screenshot({ path: 'test-results/settings-page-loaded.png', fullPage: true });
    });

    test('should display settings navigation tabs', async ({ page }) => {
      await page.goto(settingsUrl);
      
      // Tabs navigation
      await expect(page.locator('[data-testid=\"settings-tabs\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"tab-general\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"tab-analysis\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"tab-notifications\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"tab-api\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/settings-navigation-tabs.png' });
    });

    test('should show loading state for settings data', async ({ page }) => {
      await page.route('**/api/v1/settings', async route => {
        await page.waitForTimeout(1000);
        route.continue();
      });
      
      await page.goto(settingsUrl);
      
      await expect(page.locator('[data-testid=\"settings-loading\"]')).toBeVisible();
      await page.screenshot({ path: 'test-results/settings-loading-state.png' });
      
      await expect(page.locator('[data-testid=\"settings-loading\"]')).toBeHidden({ timeout: 15000 });
    });
  });

  test.describe('General Settings Tab', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(settingsUrl);
      await page.waitForLoadState('networkidle');
      await page.locator('[data-testid=\"tab-general\"]').click();
    });

    test('should display profile settings', async ({ page }) => {
      // Section profil
      await expect(page.locator('[data-testid=\"profile-settings\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Profil');
      
      // Champs profil
      await expect(page.locator('[data-testid=\"profile-name-input\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"profile-email-input\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"profile-company-input\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/profile-settings-section.png' });
    });

    test('should allow updating profile information', async ({ page }) => {
      // Modifier nom
      await page.locator('[data-testid=\"profile-name-input\"]').fill('John Doe Updated');
      
      // Modifier entreprise
      await page.locator('[data-testid=\"profile-company-input\"]').fill('SEPTEO Updated');
      
      // Sauvegarder
      await page.locator('[data-testid=\"save-profile-button\"]').click();
      
      // Vérifier message succès
      await expect(page.locator('[data-testid=\"success-message\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"success-message\"]')).toContainText('Profil mis à jour');
      
      await page.screenshot({ path: 'test-results/profile-updated-success.png' });
    });

    test('should display language and timezone settings', async ({ page }) => {
      // Section langue/timezone
      await expect(page.locator('[data-testid=\"locale-settings\"]')).toBeVisible();
      
      // Sélecteur langue
      await expect(page.locator('[data-testid=\"language-select\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"timezone-select\"]')).toBeVisible();
      
      // Options disponibles
      await page.locator('[data-testid=\"language-select\"]').click();
      await expect(page.locator('[data-testid=\"lang-fr\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"lang-en\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/locale-settings-options.png' });
    });

    test('should handle theme preference toggle', async ({ page }) => {
      // Toggle thème sombre
      await expect(page.locator('[data-testid=\"theme-toggle\"]')).toBeVisible();
      
      // Activer mode sombre
      await page.locator('[data-testid=\"theme-toggle\"]').click();
      
      // Vérifier changement visuel
      await expect(page.locator('body')).toHaveClass(/dark/);
      
      await page.screenshot({ path: 'test-results/dark-theme-activated.png' });
      
      // Désactiver mode sombre
      await page.locator('[data-testid=\"theme-toggle\"]').click();
      await expect(page.locator('body')).not.toHaveClass(/dark/);
    });
  });

  test.describe('Analysis Settings Tab', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(settingsUrl);
      await page.waitForLoadState('networkidle');
      await page.locator('[data-testid=\"tab-analysis\"]').click();
    });

    test('should display default analysis configuration', async ({ page }) => {
      // Section config analyses par défaut
      await expect(page.locator('[data-testid=\"default-analysis-config\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Configuration des Analyses');
      
      // Type d'analyse par défaut
      await expect(page.locator('[data-testid=\"default-analysis-type\"]')).toBeVisible();
      
      // Checkbox options
      await expect(page.locator('[data-testid=\"include-images-check\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"deep-crawl-check\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"competitor-analysis-check\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/default-analysis-config.png' });
    });

    test('should allow configuring crawl settings', async ({ page }) => {
      // Section paramètres crawl
      await expect(page.locator('[data-testid=\"crawl-settings\"]')).toBeVisible();
      
      // Paramètres crawl
      await expect(page.locator('[data-testid=\"max-pages-input\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"crawl-delay-input\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"user-agent-input\"]')).toBeVisible();
      
      // Modifier paramètres
      await page.locator('[data-testid=\"max-pages-input\"]').fill('500');
      await page.locator('[data-testid=\"crawl-delay-input\"]').fill('2000');
      
      // Sauvegarder
      await page.locator('[data-testid=\"save-crawl-settings\"]').click();
      
      await expect(page.locator('[data-testid=\"crawl-settings-saved\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/crawl-settings-updated.png' });
    });

    test('should configure SEO analysis preferences', async ({ page }) => {
      // Section préférences SEO
      await expect(page.locator('[data-testid=\"seo-preferences\"]')).toBeVisible();
      
      // Options SEO
      await expect(page.locator('[data-testid=\"keyword-density-check\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"meta-analysis-check\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"schema-validation-check\"]')).toBeVisible();
      
      // Cocher options
      await page.locator('[data-testid=\"schema-validation-check\"]').check();
      
      // Paramètres keywords
      await expect(page.locator('[data-testid=\"target-keywords-input\"]')).toBeVisible();
      await page.locator('[data-testid=\"target-keywords-input\"]').fill('plage, marina, restaurant');
      
      await page.screenshot({ path: 'test-results/seo-preferences-configured.png' });
    });
  });

  test.describe('Notifications Settings Tab', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(settingsUrl);
      await page.waitForLoadState('networkidle');
      await page.locator('[data-testid=\"tab-notifications\"]').click();
    });

    test('should display email notification preferences', async ({ page }) => {
      // Section notifications email
      await expect(page.locator('[data-testid=\"email-notifications\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Notifications Email');
      
      // Options notifications
      await expect(page.locator('[data-testid=\"analysis-complete-email\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"weekly-summary-email\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"critical-issues-email\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/email-notifications-preferences.png' });
    });

    test('should configure notification frequency', async ({ page }) => {
      // Fréquence notifications
      await expect(page.locator('[data-testid=\"notification-frequency\"]')).toBeVisible();
      
      // Options fréquence
      await page.locator('[data-testid=\"frequency-select\"]').click();
      await expect(page.locator('[data-testid=\"freq-immediate\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"freq-daily\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"freq-weekly\"]')).toBeVisible();
      
      // Sélectionner daily
      await page.locator('[data-testid=\"freq-daily\"]').click();
      
      await page.screenshot({ path: 'test-results/notification-frequency-daily.png' });
    });

    test('should allow testing notification settings', async ({ page }) => {
      // Bouton test notification
      await expect(page.locator('[data-testid=\"test-notification-button\"]')).toBeVisible();
      
      await page.locator('[data-testid=\"test-notification-button\"]').click();
      
      // Confirmation envoi test
      await expect(page.locator('[data-testid=\"test-notification-sent\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"test-notification-sent\"]')).toContainText('Email de test envoyé');
      
      await page.screenshot({ path: 'test-results/test-notification-sent.png' });
    });
  });

  test.describe('API Settings Tab', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(settingsUrl);
      await page.waitForLoadState('networkidle');
      await page.locator('[data-testid=\"tab-api\"]').click();
    });

    test('should display API key management', async ({ page }) => {
      // Section clé API
      await expect(page.locator('[data-testid=\"api-key-management\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Clé API');
      
      // Clé API actuelle (masquée)
      await expect(page.locator('[data-testid=\"current-api-key\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"current-api-key\"]')).toContainText('••••••••');
      
      // Boutons gestion
      await expect(page.locator('[data-testid=\"show-api-key\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"regenerate-api-key\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"copy-api-key\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/api-key-management.png' });
    });

    test('should allow revealing and copying API key', async ({ page }) => {
      // Révéler clé API
      await page.locator('[data-testid=\"show-api-key\"]').click();
      
      // Vérifier clé visible
      const apiKeyVisible = page.locator('[data-testid=\"api-key-visible\"]');
      await expect(apiKeyVisible).toBeVisible();
      await expect(apiKeyVisible).not.toContainText('••••••••');
      
      // Copier clé
      await page.locator('[data-testid=\"copy-api-key\"]').click();
      
      // Message copié
      await expect(page.locator('[data-testid=\"api-key-copied\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/api-key-revealed-copied.png' });
    });

    test('should display API usage statistics', async ({ page }) => {
      // Section stats usage API
      await expect(page.locator('[data-testid=\"api-usage-stats\"]')).toBeVisible();
      await expect(page.locator('h2')).toContainText('Utilisation API');
      
      // Métriques usage
      await expect(page.locator('[data-testid=\"total-requests\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"monthly-limit\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"remaining-requests\"]')).toBeVisible();
      
      // Graphique usage
      await expect(page.locator('[data-testid=\"usage-chart\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/api-usage-statistics.png' });
    });

    test('should configure API rate limits', async ({ page }) => {
      // Section rate limiting
      await expect(page.locator('[data-testid=\"rate-limiting-config\"]')).toBeVisible();
      
      // Paramètres rate limit
      await expect(page.locator('[data-testid=\"requests-per-minute\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"burst-limit\"]')).toBeVisible();
      
      // Modifier limites
      await page.locator('[data-testid=\"requests-per-minute\"]').fill('100');
      await page.locator('[data-testid=\"burst-limit\"]').fill('200');
      
      // Sauvegarder
      await page.locator('[data-testid=\"save-rate-limits\"]').click();
      
      await expect(page.locator('[data-testid=\"rate-limits-saved\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/api-rate-limits-configured.png' });
    });
  });

  test.describe('Form Validation & Error Handling', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(settingsUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should validate email format in profile settings', async ({ page }) => {
      await page.locator('[data-testid=\"tab-general\"]').click();
      
      // Email invalide
      await page.locator('[data-testid=\"profile-email-input\"]').fill('invalid-email');
      await page.locator('[data-testid=\"save-profile-button\"]').click();
      
      // Message erreur validation
      await expect(page.locator('[data-testid=\"email-validation-error\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"email-validation-error\"]')).toContainText('Email invalide');
      
      await page.screenshot({ path: 'test-results/email-validation-error.png' });
    });

    test('should validate crawl settings ranges', async ({ page }) => {
      await page.locator('[data-testid=\"tab-analysis\"]').click();
      
      // Valeur trop élevée pour max pages
      await page.locator('[data-testid=\"max-pages-input\"]').fill('99999');
      await page.locator('[data-testid=\"save-crawl-settings\"]').click();
      
      // Erreur validation
      await expect(page.locator('[data-testid=\"max-pages-validation-error\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"max-pages-validation-error\"]')).toContainText('Maximum 10000');
      
      await page.screenshot({ path: 'test-results/crawl-settings-validation-error.png' });
    });

    test('should handle API key regeneration confirmation', async ({ page }) => {
      await page.locator('[data-testid=\"tab-api\"]').click();
      
      // Régénérer clé API
      await page.locator('[data-testid=\"regenerate-api-key\"]').click();
      
      // Modal confirmation
      await expect(page.locator('[data-testid=\"regenerate-confirmation-modal\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"regenerate-warning\"]')).toContainText('ancienne clé sera invalidée');
      
      // Confirmer
      await page.locator('[data-testid=\"confirm-regenerate\"]').click();
      
      // Nouvelle clé générée
      await expect(page.locator('[data-testid=\"new-api-key-generated\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/api-key-regeneration-flow.png' });
    });
  });

  test.describe('SEPTEO Design Standards', () => {
    test.beforeEach(async ({ page }) => {
      await page.goto(settingsUrl);
      await page.waitForLoadState('networkidle');
    });

    test('should use SEPTEO orange for primary actions', async ({ page }) => {
      // Boutons principaux
      const saveButtons = page.locator('[data-testid*=\"save-\"]');
      const firstSaveButton = saveButtons.first();
      await expect(firstSaveButton).toHaveCSS('background-color', 'rgb(255, 97, 54)');
      
      // Onglets actifs
      const activeTab = page.locator('[data-testid=\"tab-general\"][class*=\"active\"]');
      await expect(activeTab).toHaveCSS('border-bottom-color', 'rgb(255, 97, 54)');
      
      await page.screenshot({ path: 'test-results/settings-septeo-orange.png' });
    });

    test('should maintain consistent form styling', async ({ page }) => {
      await page.locator('[data-testid=\"tab-general\"]').click();
      
      // Tous les inputs ont le même style
      const inputs = page.locator('input[type=\"text\"], input[type=\"email\"]');
      const inputCount = await inputs.count();
      
      for (let i = 0; i < Math.min(inputCount, 3); i++) {
        const input = inputs.nth(i);
        await expect(input).toHaveCSS('border-radius', '6px');
        await expect(input).toHaveCSS('padding', '12px');
      }
      
      await page.screenshot({ path: 'test-results/consistent-form-styling.png' });
    });
  });

  test.describe('Responsive Design', () => {
    test('should stack settings on mobile', async ({ page }) => {
      await page.setViewportSize({ width: 375, height: 667 });
      await page.goto(settingsUrl);
      
      // Tabs en dropdown sur mobile
      await expect(page.locator('[data-testid=\"mobile-tabs-dropdown\"]')).toBeVisible();
      
      // Formulaires adaptés
      await expect(page.locator('[data-testid=\"profile-settings\"]')).toBeVisible();
      
      // Inputs full width
      const nameInput = page.locator('[data-testid=\"profile-name-input\"]');
      const inputBox = await nameInput.boundingBox();
      expect(inputBox.width).toBeGreaterThan(300); // Quasi full width
      
      await page.screenshot({ path: 'test-results/settings-mobile-responsive.png', fullPage: true });
    });

    test('should adapt tabs layout on tablet', async ({ page }) => {
      await page.setViewportSize({ width: 768, height: 1024 });
      await page.goto(settingsUrl);
      
      // Tabs horizontaux sur tablet
      await expect(page.locator('[data-testid=\"settings-tabs\"]')).toBeVisible();
      
      // Layout en 2 colonnes possible
      const settingsContainer = page.locator('[data-testid=\"settings-container\"]');
      const containerBox = await settingsContainer.boundingBox();
      expect(containerBox.width).toBeLessThan(768);
      
      await page.screenshot({ path: 'test-results/settings-tablet-responsive.png', fullPage: true });
    });
  });

  test.describe('Accessibility Standards (> 95%)', () => {
    test('should support keyboard navigation between tabs', async ({ page }) => {
      await page.goto(settingsUrl);
      await page.waitForLoadState('networkidle');
      
      // Navigation tabs par clavier
      await page.keyboard.press('Tab');
      await page.keyboard.press('ArrowRight');
      
      // Tab Analysis devrait être actif
      await expect(page.locator('[data-testid=\"tab-analysis\"][class*=\"active\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/settings-keyboard-navigation.png' });
    });

    test('should have proper form labels and descriptions', async ({ page }) => {
      await page.goto(settingsUrl);
      
      // Tous les champs ont des labels
      const inputs = page.locator('input');
      const inputCount = await inputs.count();
      
      for (let i = 0; i < Math.min(inputCount, 5); i++) {
        const input = inputs.nth(i);
        const inputId = await input.getAttribute('id');
        if (inputId) {
          await expect(page.locator(`label[for=\"${inputId}\"]`)).toBeVisible();
        }
      }
      
      await page.screenshot({ path: 'test-results/form-labels-accessibility.png' });
    });

    test('should announce setting changes to screen readers', async ({ page }) => {
      await page.goto(settingsUrl);
      
      // Live regions pour feedback
      await expect(page.locator('[data-testid=\"settings-feedback\"]')).toHaveAttribute('aria-live', 'polite');
      
      // Status messages descriptifs
      const successMessages = page.locator('[data-testid*=\"success-message\"]');
      if (await successMessages.count() > 0) {
        await expect(successMessages.first()).toHaveAttribute('role', 'status');
      }
      
      await page.screenshot({ path: 'test-results/screen-reader-announcements.png' });
    });
  });

  test.describe('Data Persistence & Sync', () => {
    test('should persist settings across page reloads', async ({ page }) => {
      await page.goto(settingsUrl);
      
      // Modifier un paramètre
      await page.locator('[data-testid=\"profile-name-input\"]').fill('Test Persistence');
      await page.locator('[data-testid=\"save-profile-button\"]').click();
      
      // Recharger page
      await page.reload();
      await page.waitForLoadState('networkidle');
      
      // Vérifier persistance
      await expect(page.locator('[data-testid=\"profile-name-input\"]')).toHaveValue('Test Persistence');
      
      await page.screenshot({ path: 'test-results/settings-persistence-check.png' });
    });

    test('should handle concurrent settings modifications', async ({ page, context }) => {
      // Ouvrir deux onglets
      const page2 = await context.newPage();
      
      await page.goto(settingsUrl);
      await page2.goto(settingsUrl);
      
      // Modification simultanée
      await page.locator('[data-testid=\"profile-name-input\"]').fill('User Tab 1');
      await page2.locator('[data-testid=\"profile-name-input\"]').fill('User Tab 2');
      
      // Sauvegarder sur tab 1
      await page.locator('[data-testid=\"save-profile-button\"]').click();
      
      // Sauvegarder sur tab 2 (conflict potentiel)
      await page2.locator('[data-testid=\"save-profile-button\"]').click();
      
      // Vérifier gestion conflit
      await expect(page2.locator('[data-testid=\"conflict-warning\"]')).toBeVisible();
      
      await page.screenshot({ path: 'test-results/concurrent-modifications-tab1.png' });
      await page2.screenshot({ path: 'test-results/concurrent-modifications-tab2.png' });
    });
  });
});

/**
 * Cross-Browser Tests - OBLIGATOIRES Lead Tech Standards
 */
['chromium', 'firefox', 'webkit'].forEach(browserName => {
  test.describe(`Settings Page - ${browserName}`, () => {
    test(`should render settings correctly on ${browserName}`, async ({ page }) => {
      await page.goto('/settings');
      await page.waitForLoadState('networkidle');
      
      // Tests basiques cross-browser
      await expect(page.locator('h1')).toContainText('Paramètres');
      await expect(page.locator('[data-testid=\"settings-tabs\"]')).toBeVisible();
      await expect(page.locator('[data-testid=\"profile-settings\"]')).toBeVisible();
      
      await page.screenshot({ 
        path: `test-results/settings-${browserName}-compatibility.png`,
        fullPage: true 
      });
    });
  });
});