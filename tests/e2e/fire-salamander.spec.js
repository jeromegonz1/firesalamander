// Fire Salamander E2E Tests
const { test, expect } = require('@playwright/test');

test.describe('Fire Salamander - US-1.2 Interface Visuelle', () => {
  
  test('Page d\'accueil - Formulaire d\'analyse', async ({ page }) => {
    await page.goto('/');
    
    // Vérifier le titre et branding SEPTEO
    await expect(page).toHaveTitle(/Fire Salamander/);
    await expect(page.locator('text=Fire Salamander')).toBeVisible();
    
    // Vérifier le formulaire
    const urlInput = page.locator('input[type="url"]');
    await expect(urlInput).toBeVisible();
    await expect(urlInput).toHaveAttribute('placeholder', 'https://votre-site.com');
    
    // Vérifier le bouton d'analyse
    const analyzeButton = page.locator('button:has-text("ANALYSER MON SITE")');
    await expect(analyzeButton).toBeVisible();
    
    // Vérifier les couleurs SEPTEO (orange #ff6136)
    const fireIcon = page.locator('i.fa-fire');
    await expect(fireIcon).toHaveClass(/text-septeo-orange/);
    
    // Screenshot pour validation visuelle
    await page.screenshot({ path: 'tests/screenshots/home.png', fullPage: true });
  });

  test('Navigation vers page d\'analyse', async ({ page }) => {
    await page.goto('/');
    
    // Remplir le formulaire
    await page.fill('input[type="url"]', 'https://example.com');
    
    // Cliquer sur analyser
    await page.click('button:has-text("ANALYSER MON SITE")');
    
    // Vérifier la redirection
    await expect(page).toHaveURL(/.*analyze\?url=https%3A%2F%2Fexample\.com/);
    
    // Vérifier le contenu de la page d'analyse
    await expect(page.locator('text=Analyse en cours')).toBeVisible();
    await expect(page.locator('text=https://example.com')).toBeVisible();
    
    // Vérifier la barre de progression
    const progressBar = page.locator('.bg-septeo-orange');
    await expect(progressBar).toBeVisible();
    
    // Screenshot
    await page.screenshot({ path: 'tests/screenshots/analyzing.png', fullPage: true });
  });

  test('Page de résultats SEO', async ({ page }) => {
    await page.goto('/results?url=https://example.com');
    
    // Vérifier le titre et score
    await expect(page.locator('text=Score Global SEO')).toBeVisible();
    await expect(page.locator('text=example.com')).toBeVisible();
    
    // Vérifier le score (85 dans les données de test)
    await expect(page.locator('text=85')).toBeVisible();
    await expect(page.locator('text=sur 100')).toBeVisible();
    
    // Vérifier les sections principales
    await expect(page.locator('text=Problèmes Critiques')).toBeVisible();
    await expect(page.locator('text=Avertissements')).toBeVisible();
    await expect(page.locator('text=Suggestions IA')).toBeVisible();
    
    // Vérifier le bouton PDF
    const pdfButton = page.locator('button:has-text("Exporter PDF")');
    await expect(pdfButton).toBeVisible();
    await expect(pdfButton).toHaveClass(/bg-septeo-orange/);
    
    // Screenshot
    await page.screenshot({ path: 'tests/screenshots/results.png', fullPage: true });
  });

  test('Design responsive mobile', async ({ page }) => {
    // Simuler un mobile
    await page.setViewportSize({ width: 375, height: 667 });
    await page.goto('/');
    
    // Vérifier que le design est responsive
    await expect(page.locator('text=Fire Salamander')).toBeVisible();
    const form = page.locator('form');
    await expect(form).toBeVisible();
    
    // Screenshot mobile
    await page.screenshot({ path: 'tests/screenshots/mobile-home.png', fullPage: true });
  });

  test('Validation des couleurs SEPTEO', async ({ page }) => {
    await page.goto('/');
    
    // Vérifier que les couleurs SEPTEO sont appliquées
    const orangeElements = await page.locator('.text-septeo-orange').count();
    const blueElements = await page.locator('.bg-septeo-blue').count();
    
    expect(orangeElements).toBeGreaterThan(0);
    expect(blueElements).toBeGreaterThan(0);
  });

  test('Performance - Chargement < 2s', async ({ page }) => {
    const startTime = Date.now();
    await page.goto('/');
    
    // Attendre que la page soit complètement chargée
    await page.waitForLoadState('networkidle');
    
    const loadTime = Date.now() - startTime;
    expect(loadTime).toBeLessThan(2000); // < 2 secondes
  });
});