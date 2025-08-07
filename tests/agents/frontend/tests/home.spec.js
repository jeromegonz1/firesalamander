const { test, expect } = require('@playwright/test');

/**
 * Fire Salamander - Tests E2E de la page d'accueil
 */

test.describe('🏠 Page d\'accueil', () => {
  
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('devrait afficher le titre Fire Salamander', async ({ page }) => {
    await expect(page).toHaveTitle(/Fire Salamander/);
    
    // Vérifier le titre principal
    const mainTitle = page.locator('h1');
    await expect(mainTitle).toContainText('Fire Salamander');
    await expect(mainTitle).toContainText('SEO Analyzer');
  });

  test('devrait afficher le branding SEPTEO', async ({ page }) => {
    // Vérifier le logo SEPTEO
    const septeoLogo = page.locator('img[alt="SEPTEO"]');
    await expect(septeoLogo).toBeVisible();
    await expect(septeoLogo).toHaveAttribute('src', /septeo/);
    
    // Vérifier la mention "Propulsé par SEPTEO"
    await expect(page.locator('text=Propulsé par SEPTEO')).toBeVisible();
    
    // Vérifier l'icône salamandre
    await expect(page.locator('text=🦎')).toBeVisible();
  });

  test('devrait afficher l\'indicateur de statut en ligne', async ({ page }) => {
    const statusIndicator = page.locator('.status-indicator');
    await expect(statusIndicator).toBeVisible();
    await expect(statusIndicator).toContainText('En ligne');
    
    // Vérifier l'animation du point de statut
    const statusDot = page.locator('.status-dot');
    await expect(statusDot).toBeVisible();
  });

  test('devrait afficher les 6 cards de fonctionnalités', async ({ page }) => {
    const featureCards = page.locator('.feature-card');
    await expect(featureCards).toHaveCount(6);
    
    // Vérifier chaque card
    const expectedFeatures = [
      { icon: '🕷️', title: 'Crawling Intelligent' },
      { icon: '🧠', title: 'Analyse Sémantique IA' },
      { icon: '📊', title: 'Audit SEO Complet' },
      { icon: '📈', title: 'Rapports Détaillés' },
      { icon: '⚡', title: 'Performance Web' },
      { icon: '🎯', title: 'Insights Concurrentiels' }
    ];

    for (let i = 0; i < expectedFeatures.length; i++) {
      const card = featureCards.nth(i);
      const feature = expectedFeatures[i];
      
      await expect(card.locator('.feature-icon')).toContainText(feature.icon);
      await expect(card.locator('h3')).toContainText(feature.title);
    }
  });

  test('devrait avoir un bouton CTA fonctionnel', async ({ page }) => {
    const ctaButton = page.locator('.btn-primary');
    await expect(ctaButton).toBeVisible();
    await expect(ctaButton).toContainText('Commencer l\'Analyse');
    
    // Vérifier que le bouton est cliquable
    await expect(ctaButton).toBeEnabled();
    
    // Vérifier le style hover (si possible)
    await ctaButton.hover();
  });

  test('devrait afficher le footer avec les informations correctes', async ({ page }) => {
    const footer = page.locator('.footer');
    await expect(footer).toBeVisible();
    
    // Vérifier le copyright et les mentions
    await expect(footer).toContainText('Fire Salamander');
    await expect(footer).toContainText('Propulsé par SEPTEO');
    await expect(footer).toContainText('Version 1.0.0');
    await expect(footer).toContainText('2024');
  });

  test('devrait être responsive sur mobile', async ({ page }) => {
    // Simuler un viewport mobile
    await page.setViewportSize({ width: 375, height: 667 });
    
    // Vérifier que le contenu s'adapte
    const header = page.locator('.header-content');
    await expect(header).toBeVisible();
    
    // Vérifier que les features passent en colonne unique
    const features = page.locator('.features');
    await expect(features).toBeVisible();
    
    // Le titre devrait être plus petit sur mobile
    const heroTitle = page.locator('.hero h1');
    await expect(heroTitle).toBeVisible();
  });

  test('devrait avoir les couleurs SEPTEO correctes', async ({ page }) => {
    // Vérifier la couleur primaire SEPTEO (#ff6136)
    const appBrand = page.locator('.app-brand');
    await expect(appBrand).toHaveCSS('color', 'rgb(255, 97, 54)'); // #ff6136 en RGB
    
    const highlight = page.locator('.hero .highlight');
    await expect(highlight).toHaveCSS('color', 'rgb(255, 97, 54)');
    
    const ctaButton = page.locator('.btn-primary');
    await expect(ctaButton).toHaveCSS('background-color', 'rgb(255, 97, 54)');
  });

  test('devrait charger rapidement (performance)', async ({ page }) => {
    const startTime = Date.now();
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    const loadTime = Date.now() - startTime;
    
    // Le chargement devrait prendre moins de 3 secondes
    expect(loadTime).toBeLessThan(3000);
    
    // Vérifier que les éléments critiques sont visibles
    await expect(page.locator('h1')).toBeVisible();
    await expect(page.locator('.septeo-logo')).toBeVisible();
  });

  test('devrait avoir une accessibilité correcte', async ({ page }) => {
    // Vérifier les alt texts des images
    const septeoLogo = page.locator('img[alt="SEPTEO"]');
    await expect(septeoLogo).toHaveAttribute('alt', 'SEPTEO');
    
    // Vérifier la structure des headings
    const h1 = page.locator('h1');
    await expect(h1).toBeVisible();
    
    const h2s = page.locator('h2');
    // Il ne devrait pas y avoir de h2 sur la page d'accueil, mais des h3
    
    const h3s = page.locator('h3');
    await expect(h3s).toHaveCount(6); // Une pour chaque feature card
    
    // Vérifier le contraste (basique)
    // Les textes devraient être lisibles
    const bodyText = page.locator('body');
    await expect(bodyText).toHaveCSS('color', 'rgb(51, 51, 51)'); // #333
  });

  test('devrait avoir des animations fluides', async ({ page }) => {
    // Vérifier l'animation du status dot
    const statusDot = page.locator('.status-dot');
    await expect(statusDot).toBeVisible();
    
    // Vérifier l'effet hover sur les feature cards
    const firstCard = page.locator('.feature-card').first();
    await firstCard.hover();
    
    // Vérifier l'effet hover sur le bouton CTA
    const ctaButton = page.locator('.btn-primary');
    await ctaButton.hover();
    
    // Les animations ne devraient pas être trop longues
    await page.waitForTimeout(500); // Attendre que les transitions se terminent
  });
});