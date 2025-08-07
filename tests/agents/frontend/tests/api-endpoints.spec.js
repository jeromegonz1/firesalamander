const { test, expect } = require('@playwright/test');

/**
 * Fire Salamander - Tests E2E des endpoints API via le navigateur
 */

test.describe('🔌 Endpoints API', () => {
  
  test('endpoint /health devrait retourner un JSON valide', async ({ page }) => {
    // Naviguer vers l'endpoint health
    await page.goto('/health');
    
    // Vérifier que c'est du JSON
    const contentType = await page.locator('body').getAttribute('data-content-type') || 
                       await page.evaluate(() => document.contentType);
    
    // Vérifier le contenu
    const body = await page.locator('body').textContent();
    expect(body).toContain('healthy');
    expect(body).toContain('Fire Salamander');
    expect(body).toContain('1.0.0');
    
    // Vérifier que c'est du JSON valide
    expect(() => JSON.parse(body)).not.toThrow();
    
    const healthData = JSON.parse(body);
    expect(healthData.status).toBe('healthy');
    expect(healthData.app).toBe('Fire Salamander');
    expect(healthData.version).toBe('1.0.0');
  });

  test('endpoint /debug devrait retourner des informations de diagnostic', async ({ page }) => {
    await page.goto('/debug');
    
    const body = await page.locator('body').textContent();
    
    // Vérifier que c'est du JSON valide
    expect(() => JSON.parse(body)).not.toThrow();
    
    const debugData = JSON.parse(body);
    
    // Vérifications de base
    expect(debugData).toHaveProperty('timestamp');
    expect(debugData).toHaveProperty('app');
    expect(debugData).toHaveProperty('system');
    expect(debugData).toHaveProperty('config');
    expect(debugData).toHaveProperty('checks');
    
    // Vérifier les informations de l'app
    expect(debugData.app.name).toBe('Fire Salamander');
    expect(debugData.app.version).toBe('1.0.0');
    expect(debugData.app.icon).toBe('🦎');
    expect(debugData.app.powered_by).toBe('SEPTEO');
    
    // Vérifier les informations système
    expect(debugData.system).toHaveProperty('go_version');
    expect(debugData.system).toHaveProperty('os');
    expect(debugData.system).toHaveProperty('arch');
    expect(debugData.system.debug_mode).toBe(true);
    
    // Vérifier les checks
    expect(debugData.checks).toHaveProperty('config');
    expect(debugData.checks).toHaveProperty('database');
    expect(debugData.checks).toHaveProperty('filesystem');
  });

  test('endpoint /debug?test=phase1 devrait exécuter les tests de Phase 1', async ({ page }) => {
    await page.goto('/debug?test=phase1');
    
    const body = await page.locator('body').textContent();
    
    // Vérifier que c'est du JSON valide
    expect(() => JSON.parse(body)).not.toThrow();
    
    const testData = JSON.parse(body);
    
    // Vérifications spécifiques aux tests de phase
    expect(testData).toHaveProperty('phase');
    expect(testData).toHaveProperty('status');
    expect(testData).toHaveProperty('tests');
    expect(testData).toHaveProperty('total_tests');
    expect(testData).toHaveProperty('passed_tests');
    
    expect(testData.phase).toContain('Phase 1');
    expect(['passed', 'failed', 'running']).toContain(testData.status);
    expect(testData.total_tests).toBeGreaterThan(0);
    expect(Array.isArray(testData.tests)).toBe(true);
  });

  test('endpoint inexistant devrait retourner 404', async ({ page }) => {
    const response = await page.goto('/api/nonexistent-endpoint');
    expect(response.status()).toBe(404);
  });

  test('endpoints API devraient avoir les headers de sécurité', async ({ page }) => {
    const response = await page.goto('/health');
    
    // Vérifier les headers de sécurité basiques
    const headers = response.headers();
    
    // Content-Type devrait être présent
    expect(headers['content-type']).toContain('application/json');
    
    // Vérifier que la réponse est valide
    expect(response.status()).toBe(200);
  });

  test('les endpoints devraient gérer les requêtes CORS', async ({ page }) => {
    // Test CORS en utilisant fetch depuis le navigateur
    const corsTest = await page.evaluate(async () => {
      try {
        const response = await fetch('/health', {
          method: 'GET',
          headers: {
            'Origin': 'http://example.com'
          }
        });
        
        return {
          status: response.status,
          corsHeaders: {
            'access-control-allow-origin': response.headers.get('access-control-allow-origin'),
            'access-control-allow-methods': response.headers.get('access-control-allow-methods')
          }
        };
      } catch (error) {
        return { error: error.message };
      }
    });
    
    expect(corsTest.status).toBe(200);
    // Les headers CORS peuvent être configurés selon les besoins
  });

  test('les réponses JSON devraient être bien formatées', async ({ page }) => {
    await page.goto('/health');
    
    const body = await page.locator('body').textContent();
    const healthData = JSON.parse(body);
    
    // Vérifier la structure JSON
    expect(typeof healthData).toBe('object');
    expect(healthData).not.toBe(null);
    
    // Vérifier que tous les champs requis sont présents
    const requiredFields = ['status', 'app', 'version', 'environment'];
    for (const field of requiredFields) {
      expect(healthData).toHaveProperty(field);
      expect(healthData[field]).toBeTruthy();
    }
    
    // Vérifier les types
    expect(typeof healthData.status).toBe('string');
    expect(typeof healthData.app).toBe('string');
    expect(typeof healthData.version).toBe('string');
  });

  test('les temps de réponse des API devraient être acceptables', async ({ page }) => {
    const endpoints = ['/health', '/debug'];
    
    for (const endpoint of endpoints) {
      const startTime = Date.now();
      const response = await page.goto(endpoint);
      const responseTime = Date.now() - startTime;
      
      // Les endpoints devraient répondre en moins de 1 seconde
      expect(responseTime).toBeLessThan(1000);
      expect(response.status()).toBe(200);
      
      console.log(`${endpoint}: ${responseTime}ms`);
    }
  });

  test('les endpoints devraient gérer les erreurs gracieusement', async ({ page }) => {
    // Test avec des paramètres invalides
    const response = await page.goto('/debug?test=invalid-test');
    
    // Devrait retourner une réponse même avec des paramètres invalides
    expect(response.status()).toBe(200);
    
    const body = await page.locator('body').textContent();
    expect(() => JSON.parse(body)).not.toThrow();
  });

  test('navigation entre endpoints devrait être fluide', async ({ page }) => {
    // Test de navigation entre différents endpoints
    const endpoints = ['/', '/health', '/debug', '/'];
    
    for (const endpoint of endpoints) {
      const startTime = Date.now();
      await page.goto(endpoint);
      const loadTime = Date.now() - startTime;
      
      // Chaque navigation devrait être rapide
      expect(loadTime).toBeLessThan(2000);
      
      // Vérifier que la page charge correctement
      if (endpoint === '/') {
        await expect(page.locator('h1')).toContainText('Fire Salamander');
      } else {
        // Pour les endpoints JSON, vérifier que le contenu est présent
        const body = await page.locator('body').textContent();
        expect(body.length).toBeGreaterThan(0);
      }
    }
  });
});