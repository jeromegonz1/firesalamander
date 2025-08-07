const { defineConfig, devices } = require('@playwright/test');

module.exports = defineConfig({
  // Répertoire des tests
  testDir: './ux/user-flows',
  
  // Timeout global
  timeout: 120000, // 2 minutes par test
  expect: {
    timeout: 10000  // 10 secondes pour les assertions
  },
  
  // Configuration de la parallélisation
  fullyParallel: false, // Tests séquentiels pour éviter les conflits de ressources
  workers: 1,           // Un seul worker pour Fire Salamander
  
  // Configuration des rapports
  reporter: [
    ['html', { outputFolder: '../reports/playwright' }],
    ['json', { outputFile: '../reports/playwright/results.json' }],
    ['junit', { outputFile: '../reports/playwright/junit.xml' }],
    ['list'] // Console output
  ],
  
  // Dossier des artefacts
  outputDir: '../reports/playwright/artifacts',
  
  // Configuration des projets (navigateurs et appareils)
  projects: [
    {
      name: 'Desktop Chrome',
      use: { 
        ...devices['Desktop Chrome'],
        viewport: { width: 1920, height: 1080 },
        // Enregistrer les traces pour debug
        trace: 'retain-on-failure',
        screenshot: 'only-on-failure',
        video: 'retain-on-failure'
      },
    },
    {
      name: 'Desktop Firefox',
      use: { 
        ...devices['Desktop Firefox'],
        viewport: { width: 1920, height: 1080 },
        trace: 'retain-on-failure',
        screenshot: 'only-on-failure'
      },
    },
    {
      name: 'Desktop Safari',
      use: { 
        ...devices['Desktop Safari'],
        viewport: { width: 1920, height: 1080 },
        trace: 'retain-on-failure',
        screenshot: 'only-on-failure'
      },
    },
    {
      name: 'Mobile Chrome',
      use: { 
        ...devices['Pixel 5'],
        trace: 'retain-on-failure',
        screenshot: 'only-on-failure'
      },
    },
    {
      name: 'Mobile Safari',
      use: { 
        ...devices['iPhone 12'],
        trace: 'retain-on-failure',
        screenshot: 'only-on-failure'
      },
    },
    {
      name: 'Tablet',
      use: { 
        ...devices['iPad Pro'],
        trace: 'retain-on-failure',
        screenshot: 'only-on-failure'
      },
    }
  ],
  
  // Configuration globale
  use: {
    // URL de base
    baseURL: 'http://localhost:8080',
    
    // Attendre les requêtes réseau
    waitForNetworkIdle: true,
    
    // Configuration des captures
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    trace: 'retain-on-failure',
    
    // Permissions du navigateur
    permissions: ['clipboard-read', 'clipboard-write'],
    
    // Accepter les certificats auto-signés
    ignoreHTTPSErrors: true,
    
    // User agent personnalisé pour identifier les tests
    userAgent: 'Fire-Salamander-UX-Tests/1.0 Playwright',
    
    // Locale française
    locale: 'fr-FR',
    timezoneId: 'Europe/Paris',
    
    // Configuration des timeouts
    navigationTimeout: 30000,
    actionTimeout: 10000
  },
  
  // Serveur de développement (optionnel)
  webServer: {
    command: '../fire-salamander --config ../config.yaml',
    port: 8080,
    timeout: 30000,
    reuseExistingServer: true,
    cwd: '..'
  },
  
  // Configuration spécifique aux tests UX
  metadata: {
    project: 'Fire Salamander',
    version: '1.0.0',
    environment: 'test',
    septeoCompliance: true
  },
  
  // Hooks globaux
  globalSetup: require.resolve('./ux/setup.js'),
  globalTeardown: require.resolve('./ux/teardown.js')
});