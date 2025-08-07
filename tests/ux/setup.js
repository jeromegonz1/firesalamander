const { chromium } = require('playwright');
const fs = require('fs-extra');
const path = require('path');
const chalk = require('chalk');

async function globalSetup() {
  console.log(chalk.blue('🚀 Setting up Fire Salamander UX Test Suite...'));
  
  // Créer les dossiers nécessaires
  const dirs = [
    '../reports/ux',
    '../reports/accessibility', 
    '../reports/lighthouse',
    '../reports/design-system',
    '../reports/playwright',
    './user-flows/recordings'
  ];
  
  for (const dir of dirs) {
    await fs.ensureDir(path.join(__dirname, dir));
  }
  
  // Vérifier que Fire Salamander est démarré
  console.log(chalk.yellow('🔍 Checking Fire Salamander status...'));
  
  const browser = await chromium.launch();
  const context = await browser.newContext();
  const page = await context.newPage();
  
  try {
    // Test de santé
    const response = await page.goto('http://localhost:8080/api/v1/health', { 
      timeout: 30000 
    });
    
    if (!response.ok()) {
      throw new Error(`Fire Salamander health check failed: ${response.status()}`);
    }
    
    const healthData = await response.json();
    console.log(chalk.green(`✅ Fire Salamander is healthy - Status: ${healthData.data.status}`));
    
    // Vérifier l'interface web
    await page.goto('http://localhost:8080', { timeout: 10000 });
    await page.waitForSelector('.main-content', { timeout: 10000 });
    console.log(chalk.green('✅ Web interface is accessible'));
    
    // Prendre une capture d'écran de référence
    await page.screenshot({ 
      path: path.join(__dirname, '../reports/ux/baseline-screenshot.png'),
      fullPage: true 
    });
    
    // Collecter les métriques initiales
    const initialMetrics = await page.evaluate(() => {
      return {
        userAgent: navigator.userAgent,
        viewport: {
          width: window.innerWidth,
          height: window.innerHeight
        },
        performance: {
          navigation: performance.getEntriesByType('navigation')[0],
          paint: performance.getEntriesByType('paint')
        },
        timestamp: Date.now()
      };
    });
    
    // Sauvegarder les métriques initiales
    await fs.writeJSON(
      path.join(__dirname, '../reports/ux/initial-metrics.json'), 
      initialMetrics, 
      { spaces: 2 }
    );
    
    // Test rapide de l'API
    await page.goto('http://localhost:8080/api/v1/stats');
    const statsData = await page.textContent('pre');
    const stats = JSON.parse(statsData);
    
    console.log(chalk.cyan(`📊 Initial stats - Total tasks: ${stats.data.total_tasks}`));
    
  } catch (error) {
    console.error(chalk.red('❌ Setup failed:'), error.message);
    console.error(chalk.yellow('💡 Make sure Fire Salamander is running: ./fire-salamander --config config.yaml'));
    throw error;
  } finally {
    await browser.close();
  }
  
  // Créer un fichier de configuration de test
  const testConfig = {
    startTime: new Date().toISOString(),
    environment: {
      baseURL: 'http://localhost:8080',
      userAgent: 'Fire-Salamander-UX-Tests/1.0 Playwright'
    },
    thresholds: {
      performance: {
        maxLoadTime: 5000,      // 5 secondes max pour charger une page
        maxAnalysisTime: 30000, // 30 secondes max pour une analyse
        maxInteractionTime: 1000 // 1 seconde max pour une interaction
      },
      accessibility: {
        minScore: 95,           // Score minimum d'accessibilité
        maxViolations: 0        // Zéro violation autorisée
      },
      visual: {
        maxPixelDiff: 0.1,      // 0.1% de différence visuelle max
        threshold: 0.2          // Seuil de comparaison
      }
    },
    septeoStandards: {
      primaryColor: '#ff6136',
      requiredComponents: ['btn', 'navbar', 'main-content', 'form-input'],
      maxInlineStyles: 0
    }
  };
  
  await fs.writeJSON(
    path.join(__dirname, '../reports/ux/test-config.json'),
    testConfig,
    { spaces: 2 }
  );
  
  console.log(chalk.green('✅ UX Test Suite setup completed successfully'));
  
  return testConfig;
}

module.exports = globalSetup;