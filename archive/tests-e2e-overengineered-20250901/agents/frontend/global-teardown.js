/**
 * Fire Salamander - Playwright Global Teardown
 * Nettoyage après l'exécution des tests
 */

async function globalTeardown(config) {
  console.log('🧹 Fire Salamander - Global Teardown');
  
  const fs = require('fs').promises;
  const path = require('path');
  
  try {
    // Générer un résumé des tests
    const summaryPath = 'tests/reports/frontend/test-summary.json';
    
    const summary = {
      timestamp: new Date().toISOString(),
      baseURL: config.use.baseURL,
      testCompleted: true,
      reportLocation: 'tests/reports/frontend/',
      nextSteps: [
        'Review HTML report: tests/reports/frontend/html/index.html',
        'Check screenshots: tests/reports/frontend/screenshots/',
        'View videos: tests/reports/frontend/videos/'
      ]
    };
    
    await fs.writeFile(summaryPath, JSON.stringify(summary, null, 2));
    console.log(`📊 Test summary saved: ${summaryPath}`);
    
    // Nettoyer les fichiers temporaires anciens (plus de 7 jours)
    const cleanupDirs = [
      'tests/reports/frontend/screenshots',
      'tests/reports/frontend/videos',
      'test-results'
    ];
    
    const sevenDaysAgo = Date.now() - (7 * 24 * 60 * 60 * 1000);
    
    for (const dir of cleanupDirs) {
      try {
        const files = await fs.readdir(dir);
        for (const file of files) {
          const filePath = path.join(dir, file);
          const stats = await fs.stat(filePath);
          
          if (stats.mtime.getTime() < sevenDaysAgo) {
            await fs.unlink(filePath);
            console.log(`🗑️  Cleaned old file: ${filePath}`);
          }
        }
      } catch (error) {
        // Directory might not exist or other error
      }
    }
    
    console.log('✅ Teardown completed successfully');
    
  } catch (error) {
    console.error('⚠️  Teardown error:', error.message);
  }
}

module.exports = globalTeardown;