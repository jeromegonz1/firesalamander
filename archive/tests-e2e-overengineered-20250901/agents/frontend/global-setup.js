/**
 * Fire Salamander - Playwright Global Setup
 * Configuration globale avant l'exécution des tests
 */

async function globalSetup(config) {
  console.log('🔥 Fire Salamander - Global Setup');
  
  // Vérifier que le serveur est démarré
  const baseURL = config.use?.baseURL || process.env.BASE_URL || 'http://localhost:3000';
  console.log(`Checking server at: ${baseURL}`);
  
  // Attendre que le serveur soit prêt
  let serverReady = false;
  let attempts = 0;
  const maxAttempts = 30; // 30 secondes max
  
  while (!serverReady && attempts < maxAttempts) {
    try {
      const fetch = require('node-fetch');
      const response = await fetch(`${baseURL}/health`);
      
      if (response.status === 200) {
        const data = await response.json();
        if (data.status === 'healthy') {
          serverReady = true;
          console.log('✅ Server is ready for testing');
        }
      }
    } catch (error) {
      // Server not ready yet
    }
    
    if (!serverReady) {
      attempts++;
      console.log(`⏳ Waiting for server... (attempt ${attempts}/${maxAttempts})`);
      await new Promise(resolve => setTimeout(resolve, 1000));
    }
  }
  
  if (!serverReady) {
    console.error('❌ Server is not responding after 30 seconds');
    console.log('💡 Make sure Fire Salamander is running with: go run main.go');
    process.exit(1);
  }
  
  // Créer les répertoires de rapport si nécessaire
  const fs = require('fs').promises;
  const path = require('path');
  
  const reportDirs = [
    'tests/reports/frontend',
    'tests/reports/frontend/html',
    'tests/reports/frontend/screenshots',
    'tests/reports/frontend/videos'
  ];
  
  for (const dir of reportDirs) {
    try {
      await fs.mkdir(dir, { recursive: true });
    } catch (error) {
      // Directory might already exist
    }
  }
  
  console.log('📁 Report directories created');
  
  // Log de configuration
  console.log('🧪 Test Configuration:');
  console.log(`  Base URL: ${baseURL}`);
  console.log(`  Workers: ${config.workers || 'default'}`);
  console.log(`  Retries: ${config.retries || 0}`);
  console.log(`  Timeout: ${config.timeout || 30000}ms`);
  
  return {
    baseURL,
    serverReady: true,
    setupTime: new Date().toISOString()
  };
}

module.exports = globalSetup;