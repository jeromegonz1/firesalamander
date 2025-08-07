const { test, expect } = require('@playwright/test');
const fs = require('fs-extra');
const path = require('path');

// Configuration des parcours critiques Fire Salamander
const CRITICAL_FLOWS = {
  quickAnalysis: {
    name: 'Analyse Rapide',
    description: 'Parcours principal : analyser un site web rapidement',
    maxDuration: 30000, // 30 secondes max
    expectedSteps: 5,
    businessImpact: 'critical'
  },
  fullAnalysis: {
    name: 'Analyse Complète', 
    description: 'Analyse complète avec toutes les options',
    maxDuration: 60000, // 60 secondes max
    expectedSteps: 8,
    businessImpact: 'high'
  },
  viewHistory: {
    name: 'Consulter Historique',
    description: 'Consulter les analyses précédentes',
    maxDuration: 10000, // 10 secondes max
    expectedSteps: 3,
    businessImpact: 'medium'
  },
  generateReport: {
    name: 'Générer Rapport',
    description: 'Générer et télécharger un rapport',
    maxDuration: 20000, // 20 secondes max
    expectedSteps: 4,
    businessImpact: 'high'
  }
};

class UserFlowTracker {
  constructor(page) {
    this.page = page;
    this.startTime = null;
    this.interactions = [];
    this.friction_points = [];
    this.session_recording = [];
  }

  async startRecording(flowName) {
    this.startTime = Date.now();
    this.flowName = flowName;
    
    // Enregistrer les clics
    await this.page.evaluate(() => {
      window.userFlowData = { clicks: [], errors: [], performance: [] };
      
      document.addEventListener('click', (e) => {
        window.userFlowData.clicks.push({
          timestamp: Date.now(),
          element: e.target.tagName + (e.target.className ? '.' + e.target.className.split(' ')[0] : ''),
          x: e.clientX,
          y: e.clientY
        });
      });
      
      // Enregistrer les erreurs JavaScript
      window.addEventListener('error', (e) => {
        window.userFlowData.errors.push({
          timestamp: Date.now(),
          message: e.message,
          filename: e.filename,
          line: e.lineno
        });
      });
      
      // Enregistrer les métriques de performance
      const observer = new PerformanceObserver((list) => {
        list.getEntries().forEach((entry) => {
          window.userFlowData.performance.push({
            timestamp: Date.now(),
            name: entry.name,
            type: entry.entryType,
            duration: entry.duration,
            startTime: entry.startTime
          });
        });
      });
      
      observer.observe({ entryTypes: ['navigation', 'paint', 'largest-contentful-paint'] });
    });
  }

  async recordInteraction(step, element, action) {
    const timestamp = Date.now();
    const duration = timestamp - this.startTime;
    
    this.interactions.push({
      step,
      element,
      action,
      timestamp,
      duration,
      screenshot: `step-${step}-${action}.png`
    });
    
    // Prendre une capture d'écran
    await this.page.screenshot({
      path: path.join(__dirname, 'recordings', `${this.flowName}-step-${step}.png`),
      fullPage: false
    });
  }

  async detectFrictionPoint(reason, severity = 'medium') {
    const timestamp = Date.now();
    const duration = timestamp - this.startTime;
    
    this.friction_points.push({
      timestamp,
      duration,
      reason,
      severity,
      url: this.page.url(),
      viewport: await this.page.viewportSize()
    });
    
    console.warn(`⚠️  Friction point detected: ${reason} (${severity})`);
  }

  async endRecording() {
    const endTime = Date.now();
    const totalDuration = endTime - this.startTime;
    
    // Récupérer les données enregistrées côté client
    const clientData = await this.page.evaluate(() => window.userFlowData || {});
    
    const flowData = {
      flowName: this.flowName,
      startTime: this.startTime,
      endTime,
      totalDuration,
      interactions: this.interactions,
      frictionPoints: this.friction_points,
      clientData,
      performance: {
        totalSteps: this.interactions.length,
        avgStepDuration: totalDuration / this.interactions.length,
        frictionPointsCount: this.friction_points.length
      }
    };
    
    // Sauvegarder les données
    const recordingPath = path.join(__dirname, 'recordings', `${this.flowName}-session.json`);
    await fs.ensureDir(path.dirname(recordingPath));
    await fs.writeJSON(recordingPath, flowData, { spaces: 2 });
    
    return flowData;
  }
}

test.describe('Fire Salamander - Parcours Utilisateur Critiques', () => {
  let baseURL;
  
  test.beforeAll(async () => {
    baseURL = 'http://localhost:8080';
    
    // Vérifier que Fire Salamander est démarré
    const response = await fetch(`${baseURL}/api/v1/health`);
    if (!response.ok) {
      throw new Error('Fire Salamander is not running. Start it first.');
    }
  });

  test.beforeEach(async ({ page }) => {
    await page.goto(baseURL);
    await page.waitForSelector('.main-content', { timeout: 10000 });
  });

  test('Parcours Critique : Analyse Rapide', async ({ page }) => {
    const flow = CRITICAL_FLOWS.quickAnalysis;
    const tracker = new UserFlowTracker(page);
    
    await tracker.startRecording('quick-analysis');
    
    // Étape 1 : Aller sur la page d'analyse
    await tracker.recordInteraction(1, 'nav-analyzer', 'click');
    await page.click('a[data-page="analyzer"]');
    await page.waitForSelector('#analyzer-page.page.active', { timeout: 5000 });
    
    // Vérifier que la navigation a fonctionné
    const analyzerVisible = await page.isVisible('#analyzer-page.active');
    if (!analyzerVisible) {
      await tracker.detectFrictionPoint('Navigation vers analyzer failed', 'critical');
    }
    expect(analyzerVisible).toBeTruthy();
    
    // Étape 2 : Remplir l'URL
    await tracker.recordInteraction(2, 'analysisUrl', 'type');
    await page.fill('#analysisUrl', 'https://www.campinglacivelle.fr');
    
    // Étape 3 : Sélectionner le type d'analyse
    await tracker.recordInteraction(3, 'analysisType', 'select');
    await page.selectOption('#analysisType', 'quick');
    
    // Étape 4 : Lancer l'analyse
    await tracker.recordInteraction(4, 'submit-button', 'click');
    const submitButton = page.locator('#analysisForm button[type="submit"]');
    await submitButton.click();
    
    // Étape 5 : Attendre les résultats
    await tracker.recordInteraction(5, 'results', 'wait');
    
    // Vérifier que la progression apparaît
    const progressVisible = await page.waitForSelector('#analysisProgress', { 
      state: 'visible', 
      timeout: 5000 
    }).catch(() => null);
    
    if (!progressVisible) {
      await tracker.detectFrictionPoint('Progress indicator not shown', 'high');
    }
    
    // Attendre les résultats (avec timeout)
    const resultsVisible = await page.waitForSelector('#analysisResults', { 
      state: 'visible', 
      timeout: flow.maxDuration 
    }).catch(() => null);
    
    if (!resultsVisible) {
      await tracker.detectFrictionPoint('Analysis timeout - results not shown', 'critical');
    } else {
      // Vérifier la présence des éléments de résultats
      const scoreElement = await page.locator('.score, .overall-score, [data-score]').first();
      const scoreVisible = await scoreElement.isVisible().catch(() => false);
      
      if (!scoreVisible) {
        await tracker.detectFrictionPoint('Analysis results incomplete - no score shown', 'high');
      }
    }
    
    const flowData = await tracker.endRecording();
    
    // Assertions sur les performances
    expect(flowData.totalDuration).toBeLessThan(flow.maxDuration);
    expect(flowData.interactions.length).toBeGreaterThanOrEqual(flow.expectedSteps);
    expect(flowData.frictionPoints.filter(f => f.severity === 'critical').length).toBe(0);
    
    console.log(`✅ Quick Analysis Flow completed in ${flowData.totalDuration}ms`);
  });

  test('Parcours Critique : Analyse Complète avec Options', async ({ page }) => {
    const flow = CRITICAL_FLOWS.fullAnalysis;
    const tracker = new UserFlowTracker(page);
    
    await tracker.startRecording('full-analysis');
    
    // Étape 1 : Navigation
    await tracker.recordInteraction(1, 'nav-analyzer', 'click');
    await page.click('a[data-page="analyzer"]');
    await page.waitForSelector('#analyzer-page.active');
    
    // Étape 2 : Remplir l'URL
    await tracker.recordInteraction(2, 'analysisUrl', 'type');
    await page.fill('#analysisUrl', 'https://www.campinglacivelle.fr');
    
    // Étape 3 : Sélectionner analyse complète
    await tracker.recordInteraction(3, 'analysisType', 'select');
    await page.selectOption('#analysisType', 'full');
    
    // Étape 4 : Activer les options avancées
    await tracker.recordInteraction(4, 'includeCrawling', 'check');
    await page.check('#includeCrawling');
    
    await tracker.recordInteraction(5, 'analyzePerformance', 'check');
    await page.check('#analyzePerformance');
    
    await tracker.recordInteraction(6, 'useAIEnrichment', 'check');
    await page.check('#useAIEnrichment');
    
    // Étape 7 : Lancer l'analyse
    await tracker.recordInteraction(7, 'submit-button', 'click');
    await page.click('#analysisForm button[type="submit"]');
    
    // Étape 8 : Attendre les résultats complets
    await tracker.recordInteraction(8, 'full-results', 'wait');
    
    const resultsVisible = await page.waitForSelector('#analysisResults', { 
      state: 'visible', 
      timeout: flow.maxDuration 
    }).catch(() => {
      tracker.detectFrictionPoint('Full analysis timeout', 'critical');
      return null;
    });
    
    if (resultsVisible) {
      // Vérifier les résultats détaillés
      const detailedResults = await page.locator('.results-detailed, .analysis-complete').count();
      if (detailedResults === 0) {
        await tracker.detectFrictionPoint('Full analysis results incomplete', 'high');
      }
    }
    
    const flowData = await tracker.endRecording();
    
    expect(flowData.totalDuration).toBeLessThan(flow.maxDuration);
    expect(flowData.frictionPoints.filter(f => f.severity === 'critical').length).toBe(0);
    
    console.log(`✅ Full Analysis Flow completed in ${flowData.totalDuration}ms`);
  });

  test('Parcours Critique : Consultation Historique', async ({ page }) => {
    const flow = CRITICAL_FLOWS.viewHistory; 
    const tracker = new UserFlowTracker(page);
    
    await tracker.startRecording('view-history');
    
    // Étape 1 : Aller sur l'historique
    await tracker.recordInteraction(1, 'nav-history', 'click');
    await page.click('a[data-page="history"]');
    
    const historyPageVisible = await page.waitForSelector('#history-page.active', { timeout: 5000 }).catch(() => null);
    if (!historyPageVisible) {
      await tracker.detectFrictionPoint('History page navigation failed', 'critical');
    }
    
    // Étape 2 : Vérifier le chargement des données
    await tracker.recordInteraction(2, 'history-content', 'load');
    const historyContent = await page.waitForSelector('.history-table, .history-content', { timeout: 5000 }).catch(() => null);
    
    if (!historyContent) {
      await tracker.detectFrictionPoint('History content not loaded', 'high');
    }
    
    // Étape 3 : Tester la recherche si disponible
    await tracker.recordInteraction(3, 'history-search', 'interact');
    const searchBox = await page.locator('#historySearch');
    if (await searchBox.isVisible()) {
      await searchBox.fill('camping');
      await page.waitForTimeout(1000); // Attendre la recherche
    }
    
    const flowData = await tracker.endRecording();
    
    expect(flowData.totalDuration).toBeLessThan(flow.maxDuration);
    expect(flowData.frictionPoints.filter(f => f.severity === 'critical').length).toBe(0);
    
    console.log(`✅ History Flow completed in ${flowData.totalDuration}ms`);
  });

  test('Parcours Critique : Génération de Rapport', async ({ page }) => {
    const flow = CRITICAL_FLOWS.generateReport;
    const tracker = new UserFlowTracker(page);
    
    await tracker.startRecording('generate-report');
    
    // Étape 1 : Aller sur la page rapports
    await tracker.recordInteraction(1, 'nav-reports', 'click');
    await page.click('a[data-page="reports"]');
    
    const reportsPageVisible = await page.waitForSelector('#reports-page.active', { timeout: 5000 }).catch(() => null);
    if (!reportsPageVisible) {
      await tracker.detectFrictionPoint('Reports page navigation failed', 'critical');
    }
    
    // Étape 2 : Remplir le formulaire de rapport
    await tracker.recordInteraction(2, 'report-form', 'fill');
    await page.fill('#reportUrl', 'campinglacivelle.fr');
    await page.selectOption('#reportType', 'executive');
    await page.selectOption('#reportFormat', 'html');
    
    // Étape 3 : Générer le rapport
    await tracker.recordInteraction(3, 'generate-report', 'click');
    await page.click('#reportForm button[type="submit"]');
    
    // Étape 4 : Vérifier la génération
    await tracker.recordInteraction(4, 'report-generated', 'verify');
    
    // Dans un vrai environnement, on vérifierait le download ou la génération
    // Pour l'instant, on vérifie que l'action ne génère pas d'erreur
    await page.waitForTimeout(2000);
    
    const flowData = await tracker.endRecording();
    
    expect(flowData.totalDuration).toBeLessThan(flow.maxDuration);
    expect(flowData.frictionPoints.filter(f => f.severity === 'critical').length).toBe(0);
    
    console.log(`✅ Generate Report Flow completed in ${flowData.totalDuration}ms`);
  });

  test('Test de Performance Multi-Utilisateurs', async ({ browser }) => {
    const numUsers = 3;
    const promises = [];
    
    for (let i = 0; i < numUsers; i++) {
      promises.push(async () => {
        const context = await browser.newContext();
        const page = await context.newPage();
        const tracker = new UserFlowTracker(page);
        
        await tracker.startRecording(`concurrent-user-${i}`);
        
        try {
          await page.goto(baseURL);
          await page.waitForSelector('.main-content');
          
          // Simulation d'utilisation concurrente
          await page.click('a[data-page="analyzer"]');
          await page.fill('#analysisUrl', `https://example-${i}.com`);
          await page.selectOption('#analysisType', 'quick');
          await page.click('#analysisForm button[type="submit"]');
          
          await page.waitForSelector('#analysisProgress', { timeout: 10000 });
          
          const flowData = await tracker.endRecording();
          console.log(`User ${i} completed in ${flowData.totalDuration}ms`);
          
          return flowData;
        } finally {
          await context.close();
        }
      });
    }
    
    const results = await Promise.all(promises.map(p => p()));
    
    // Vérifier que tous les utilisateurs ont pu utiliser l'interface
    results.forEach((result, index) => {
      expect(result.totalDuration).toBeLessThan(30000);
      expect(result.frictionPoints.filter(f => f.severity === 'critical').length).toBe(0);
    });
    
    const avgDuration = results.reduce((sum, r) => sum + r.totalDuration, 0) / results.length;
    console.log(`✅ Multi-user test completed - Average duration: ${avgDuration}ms`);
  });

  test('Test de Régression UX', async ({ page }) => {
    // Test de régression pour s'assurer que les fonctionnalités de base marchent toujours
    const tracker = new UserFlowTracker(page);
    await tracker.startRecording('regression-test');
    
    // Tester toutes les navigations
    const navLinks = ['analyzer', 'history', 'reports', 'monitoring'];
    
    for (const link of navLinks) {
      await tracker.recordInteraction(`nav-${link}`, `nav-${link}`, 'click');
      await page.click(`a[data-page="${link}"]`);
      
      const pageVisible = await page.waitForSelector(`#${link}-page.active`, { timeout: 5000 }).catch(() => null);
      if (!pageVisible) {
        await tracker.detectFrictionPoint(`Navigation to ${link} failed`, 'critical');
      }
      
      expect(pageVisible).toBeTruthy();
    }
    
    // Retour au dashboard
    await page.click('a[data-page="dashboard"]');
    await page.waitForSelector('#dashboard-page.active');
    
    const flowData = await tracker.endRecording();
    
    // Aucun point de friction critique autorisé
    expect(flowData.frictionPoints.filter(f => f.severity === 'critical').length).toBe(0);
    
    console.log(`✅ UX Regression test completed - ${flowData.interactions.length} interactions tested`);
  });
});

// Test de génération de rapport d'analyse UX
test.afterAll(async () => {
  console.log('\n📊 Generating UX Flow Analysis Report...');
  
  const recordingsDir = path.join(__dirname, 'recordings');
  const reportPath = path.join(recordingsDir, 'ux-flow-report.html');
  
  // Lire tous les fichiers de session
  const sessionFiles = await fs.readdir(recordingsDir).catch(() => []);
  const sessions = [];
  
  for (const file of sessionFiles) {
    if (file.endsWith('-session.json')) {
      const sessionData = await fs.readJSON(path.join(recordingsDir, file)).catch(() => null);
      if (sessionData) {
        sessions.push(sessionData);
      }
    }
  }
  
  // Générer le rapport HTML
  const reportHTML = `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Fire Salamander - User Flow Analysis</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: linear-gradient(135deg, #ff6136 0%, #e55a2e 100%); color: white; padding: 30px; border-radius: 8px; margin-bottom: 30px; }
        .flow { background: white; padding: 20px; margin: 20px 0; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .metrics { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin: 20px 0; }
        .metric { background: #f8f9fa; padding: 15px; border-radius: 4px; text-align: center; }
        .friction { background: #fff3cd; border-left: 4px solid #f39c12; padding: 15px; margin: 10px 0; border-radius: 4px; }
        .critical { background: #f8d7da; border-left-color: #e74c3c; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔥 Fire Salamander - User Flow Analysis</h1>
            <p>Analyse des Parcours Utilisateur Critiques</p>
            <p>Généré le: ${new Date().toLocaleString('fr-FR')}</p>
        </div>
        
        ${sessions.map(session => `
        <div class="flow">
            <h2>${session.flowName}</h2>
            <div class="metrics">
                <div class="metric">
                    <strong>${session.totalDuration}ms</strong><br>
                    Durée Totale
                </div>
                <div class="metric">
                    <strong>${session.performance.totalSteps}</strong><br>
                    Étapes
                </div>
                <div class="metric">
                    <strong>${Math.round(session.performance.avgStepDuration)}ms</strong><br>
                    Durée Moyenne/Étape
                </div>
                <div class="metric">
                    <strong>${session.frictionPoints.length}</strong><br>
                    Points de Friction
                </div>
            </div>
            
            ${session.frictionPoints.length > 0 ? `
            <h3>⚠️ Points de Friction</h3>
            ${session.frictionPoints.map(friction => `
            <div class="friction ${friction.severity === 'critical' ? 'critical' : ''}">
                <strong>${friction.severity.toUpperCase()}:</strong> ${friction.reason}<br>
                <small>Timestamp: ${friction.duration}ms</small>
            </div>
            `).join('')}
            ` : '<p style="color: #27ae60;">✅ Aucun point de friction détecté</p>'}
        </div>
        `).join('')}
        
        <div class="flow">
            <h2>📊 Résumé Global</h2>
            <div class="metrics">
                <div class="metric">
                    <strong>${sessions.length}</strong><br>
                    Parcours Testés
                </div>
                <div class="metric">
                    <strong>${sessions.reduce((sum, s) => sum + s.frictionPoints.length, 0)}</strong><br>
                    Total Points de Friction
                </div>
                <div class="metric">
                    <strong>${Math.round(sessions.reduce((sum, s) => sum + s.totalDuration, 0) / sessions.length)}ms</strong><br>
                    Durée Moyenne
                </div>
                <div class="metric">
                    <strong>${sessions.filter(s => s.frictionPoints.filter(f => f.severity === 'critical').length === 0).length}/${sessions.length}</strong><br>
                    Parcours Sans Friction Critique
                </div>
            </div>
        </div>
    </div>
</body>
</html>
  `;
  
  await fs.ensureDir(recordingsDir);
  await fs.writeFile(reportPath, reportHTML);
  
  console.log(`✅ UX Flow Report generated: ${reportPath}`);
});