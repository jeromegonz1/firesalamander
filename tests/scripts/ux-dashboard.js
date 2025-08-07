const express = require('express');
const fs = require('fs-extra');
const path = require('path');
const { chromium } = require('playwright');
const chalk = require('chalk');

// Dashboard UX temps réel pour Fire Salamander
class FireSalamanderUXDashboard {
  constructor() {
    this.app = express();
    this.port = 9002;
    this.reportsDir = path.join(__dirname, '../reports');
    this.metricsHistory = [];
    this.alertsConfig = {
      accessibility: { threshold: 95, enabled: true },
      performance: { threshold: 90, enabled: true },
      designSystem: { threshold: 90, enabled: true },
      userExperience: { threshold: 90, enabled: true }
    };
    
    this.setupMiddleware();
    this.setupRoutes();
    this.startMetricsCollection();
  }

  setupMiddleware() {
    this.app.use(express.json());
    this.app.use(express.static(path.join(__dirname, '../dashboard')));
    
    // CORS pour les tests en dev
    this.app.use((req, res, next) => {
      res.header('Access-Control-Allow-Origin', '*');
      res.header('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept');
      next();
    });
  }

  setupRoutes() {
    // Page principale du dashboard
    this.app.get('/', (req, res) => {
      res.send(this.generateDashboardHTML());
    });

    // API pour les métriques en temps réel
    this.app.get('/api/metrics', async (req, res) => {
      try {
        const metrics = await this.collectCurrentMetrics();
        res.json({
          success: true,
          data: metrics,
          timestamp: new Date().toISOString()
        });
      } catch (error) {
        res.status(500).json({
          success: false,
          error: error.message
        });
      }
    });

    // API pour l'historique des métriques
    this.app.get('/api/metrics/history', (req, res) => {
      const limit = parseInt(req.query.limit) || 50;
      const history = this.metricsHistory.slice(-limit);
      
      res.json({
        success: true,
        data: history,
        count: history.length
      });
    });

    // API pour lancer des tests à la demande
    this.app.post('/api/tests/run', async (req, res) => {
      const { testType } = req.body;
      
      try {
        const result = await this.runOnDemandTest(testType);
        res.json({
          success: true,
          data: result
        });
      } catch (error) {
        res.status(500).json({
          success: false,
          error: error.message
        });
      }
    });

    // API pour les alertes
    this.app.get('/api/alerts', (req, res) => {
      const alerts = this.checkAlerts();
      res.json({
        success: true,
        data: alerts
      });
    });

    // API pour la configuration des alertes
    this.app.post('/api/alerts/config', (req, res) => {
      const { category, threshold, enabled } = req.body;
      
      if (this.alertsConfig[category]) {
        this.alertsConfig[category] = { threshold, enabled };
        res.json({ success: true });
      } else {
        res.status(400).json({ success: false, error: 'Invalid category' });
      }
    });

    // API pour les rapports consolidés
    this.app.get('/api/reports/latest', async (req, res) => {
      try {
        const reports = await this.getLatestReports();
        res.json({
          success: true,
          data: reports
        });
      } catch (error) {
        res.status(500).json({
          success: false,
          error: error.message
        });
      }
    });

    // WebSocket pour les mises à jour en temps réel
    this.app.get('/api/stream', (req, res) => {
      res.writeHead(200, {
        'Content-Type': 'text/event-stream',
        'Cache-Control': 'no-cache',
        'Connection': 'keep-alive'
      });

      const sendUpdate = (data) => {
        res.write(`data: ${JSON.stringify(data)}\n\n`);
      };

      // Envoyer des mises à jour toutes les 10 secondes
      const interval = setInterval(async () => {
        try {
          const metrics = await this.collectCurrentMetrics();
          sendUpdate({ type: 'metrics', data: metrics });
          
          const alerts = this.checkAlerts();
          if (alerts.length > 0) {
            sendUpdate({ type: 'alerts', data: alerts });
          }
        } catch (error) {
          sendUpdate({ type: 'error', data: error.message });
        }
      }, 10000);

      req.on('close', () => {
        clearInterval(interval);
      });
    });
  }

  async collectCurrentMetrics() {
    const metrics = {
      timestamp: new Date().toISOString(),
      firesalamander: {
        status: 'unknown',
        uptime: 0,
        version: '1.0.0'
      },
      accessibility: {
        score: 0,
        violations: 0,
        lastTest: null
      },
      performance: {
        lighthouse: 0,
        loadTime: 0,
        lcp: 0,
        fid: 0,
        cls: 0,
        lastTest: null
      },
      designSystem: {
        compliance: 0,
        colorViolations: 0,
        spacingViolations: 0,
        lastTest: null
      },
      userExperience: {
        score: 0,
        criticalFriction: 0,
        avgFlowTime: 0,
        lastTest: null
      }
    };

    try {
      // Vérifier le statut de Fire Salamander
      const browser = await chromium.launch({ headless: true });
      const context = await browser.newContext();
      const page = await context.newPage();

      try {
        // Test de santé
        const healthResponse = await page.goto('http://localhost:8080/api/v1/health', { timeout: 5000 });
        if (healthResponse.ok()) {
          const healthData = await healthResponse.json();
          metrics.firesalamander.status = healthData.data.status;
          metrics.firesalamander.uptime = healthData.data.uptime;
        }

        // Métriques de performance en temps réel
        await page.goto('http://localhost:8080', { waitUntil: 'networkidle' });
        
        const performanceMetrics = await page.evaluate(() => {
          const navigation = performance.getEntriesByType('navigation')[0];
          const paint = performance.getEntriesByType('paint');
          
          return {
            loadTime: navigation ? navigation.loadEventEnd - navigation.fetchStart : 0,
            fcp: paint.find(p => p.name === 'first-contentful-paint')?.startTime || 0,
            domContentLoaded: navigation ? navigation.domContentLoadedEventEnd - navigation.fetchStart : 0
          };
        });

        metrics.performance.loadTime = Math.round(performanceMetrics.loadTime);
        metrics.performance.lcp = Math.round(performanceMetrics.fcp); // Approximation

      } catch (error) {
        metrics.firesamander.status = 'down';
      } finally {
        await browser.close();
      }

      // Lire les derniers rapports s'ils existent
      await this.enrichMetricsFromReports(metrics);

    } catch (error) {
      console.error('Error collecting metrics:', error);
    }

    // Ajouter à l'historique
    this.metricsHistory.push(metrics);
    if (this.metricsHistory.length > 1000) {
      this.metricsHistory = this.metricsHistory.slice(-500); // Garder les 500 derniers
    }

    return metrics;
  }

  async enrichMetricsFromReports(metrics) {
    // Rapport d'accessibilité
    const accessibilityReport = path.join(this.reportsDir, 'accessibility/accessibility-report.json');
    if (await fs.pathExists(accessibilityReport)) {
      try {
        const report = await fs.readJSON(accessibilityReport);
        metrics.accessibility.score = Math.round(report.summary?.septeoCompliance?.overall || 0);
        metrics.accessibility.violations = report.summary?.totalViolations || 0;
        metrics.accessibility.lastTest = report.summary?.timestamp;
      } catch (e) { /* ignore */ }
    }

    // Rapport de design system
    const designReport = path.join(this.reportsDir, 'design-system/design-system-report.json');
    if (await fs.pathExists(designReport)) {
      try {
        const report = await fs.readJSON(designReport);
        metrics.designSystem.compliance = report.summary?.septeoCompliance || 0;
        metrics.designSystem.colorViolations = report.errors?.filter(e => e.type.includes('color')).length || 0;
        metrics.designSystem.spacingViolations = report.errors?.filter(e => e.type.includes('spacing')).length || 0;
        metrics.designSystem.lastTest = report.summary?.timestamp;
      } catch (e) { /* ignore */ }
    }

    // Rapports de user flows
    const userFlowsDir = path.join(__dirname, '../ux/user-flows/recordings');
    if (await fs.pathExists(userFlowsDir)) {
      try {
        const files = await fs.readdir(userFlowsDir);
        const sessionFiles = files.filter(f => f.endsWith('-session.json'));
        
        if (sessionFiles.length > 0) {
          let totalFriction = 0;
          let totalTime = 0;
          let criticalFriction = 0;

          for (const file of sessionFiles.slice(-5)) { // 5 dernières sessions
            const session = await fs.readJSON(path.join(userFlowsDir, file));
            totalTime += session.totalDuration || 0;
            totalFriction += session.frictionPoints?.length || 0;
            criticalFriction += session.frictionPoints?.filter(f => f.severity === 'critical').length || 0;
          }

          metrics.userExperience.avgFlowTime = Math.round(totalTime / sessionFiles.slice(-5).length);
          metrics.userExperience.criticalFriction = criticalFriction;
          metrics.userExperience.score = criticalFriction === 0 ? 100 : Math.max(0, 100 - (criticalFriction * 20));
        }
      } catch (e) { /* ignore */ }
    }
  }

  checkAlerts() {
    const alerts = [];
    const latestMetrics = this.metricsHistory[this.metricsHistory.length - 1];
    
    if (!latestMetrics) return alerts;

    // Vérifier les seuils configurés
    Object.entries(this.alertsConfig).forEach(([category, config]) => {
      if (!config.enabled) return;

      let currentScore = 0;
      let description = '';

      switch (category) {
        case 'accessibility':
          currentScore = latestMetrics.accessibility.score;
          description = `Score d'accessibilité: ${currentScore}%`;
          break;
        case 'performance':
          currentScore = latestMetrics.performance.lighthouse || 
                        (latestMetrics.performance.loadTime < 3000 ? 90 : 60);
          description = `Performance: ${currentScore}% (Load: ${latestMetrics.performance.loadTime}ms)`;
          break;
        case 'designSystem':
          currentScore = latestMetrics.designSystem.compliance;
          description = `Conformité design: ${currentScore}%`;
          break;
        case 'userExperience':
          currentScore = latestMetrics.userExperience.score;
          description = `UX Score: ${currentScore}% (${latestMetrics.userExperience.criticalFriction} friction critique)`;
          break;
      }

      if (currentScore < config.threshold) {
        alerts.push({
          category,
          severity: currentScore < config.threshold * 0.8 ? 'critical' : 'warning',
          currentScore,
          threshold: config.threshold,
          description,
          timestamp: latestMetrics.timestamp
        });
      }
    });

    // Alerte spéciale si Fire Salamander est down
    if (latestMetrics.firesamander?.status === 'down') {
      alerts.push({
        category: 'system',
        severity: 'critical',
        description: 'Fire Salamander is not responding',
        timestamp: latestMetrics.timestamp
      });
    }

    return alerts;
  }

  async runOnDemandTest(testType) {
    const { spawn } = require('child_process');
    
    return new Promise((resolve, reject) => {
      let command, args;
      
      switch (testType) {
        case 'accessibility':
          command = 'npm';
          args = ['run', 'test:axe'];
          break;
        case 'performance':
          command = 'npm';
          args = ['run', 'test:lighthouse'];
          break;
        case 'design-system':
          command = 'npm';
          args = ['run', 'test:design-system'];
          break;
        case 'visual':
          command = 'npm';
          args = ['run', 'test:visual'];
          break;
        default:
          return reject(new Error('Unknown test type'));
      }

      const process = spawn(command, args, {
        cwd: path.join(__dirname, '..'),
        stdio: 'pipe'
      });

      let output = '';
      process.stdout.on('data', (data) => {
        output += data.toString();
      });

      process.stderr.on('data', (data) => {
        output += data.toString();
      });

      process.on('close', (code) => {
        resolve({
          testType,
          success: code === 0,
          output,
          timestamp: new Date().toISOString()
        });
      });

      process.on('error', reject);
    });
  }

  async getLatestReports() {
    const reports = {};
    
    // Liste des rapports à collecter
    const reportPaths = [
      { key: 'accessibility', path: 'accessibility/accessibility-report.json' },
      { key: 'designSystem', path: 'design-system/design-system-report.json' },
      { key: 'lighthouse', path: 'lighthouse' }, // Dossier avec plusieurs rapports
      { key: 'userFlows', path: '../ux/user-flows/recordings' }
    ];

    for (const { key, path: reportPath } of reportPaths) {
      const fullPath = path.join(this.reportsDir, reportPath);
      
      try {
        if (await fs.pathExists(fullPath)) {
          const stat = await fs.stat(fullPath);
          
          if (stat.isFile()) {
            reports[key] = await fs.readJSON(fullPath);
          } else if (stat.isDirectory()) {
            // Prendre le fichier le plus récent
            const files = await fs.readdir(fullPath);
            const jsonFiles = files.filter(f => f.endsWith('.json'));
            
            if (jsonFiles.length > 0) {
              const latest = jsonFiles.sort((a, b) => {
                const statA = fs.statSync(path.join(fullPath, a));
                const statB = fs.statSync(path.join(fullPath, b));
                return statB.mtime - statA.mtime;
              })[0];
              
              reports[key] = await fs.readJSON(path.join(fullPath, latest));
            }
          }
        }
      } catch (error) {
        console.warn(`Could not load report ${key}:`, error.message);
      }
    }

    return reports;
  }

  generateDashboardHTML() {
    return `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Fire Salamander - Dashboard UX Temps Réel</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/3.9.1/chart.min.js"></script>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; 
            background: #f8f9fa; 
            color: #2c3e50;
        }
        .header { 
            background: linear-gradient(135deg, #ff6136 0%, #e55a2e 100%); 
            color: white; 
            padding: 20px; 
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .container { max-width: 1400px; margin: 0 auto; padding: 20px; }
        .metrics-grid { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); 
            gap: 20px; 
            margin: 20px 0; 
        }
        .metric-card { 
            background: white; 
            padding: 20px; 
            border-radius: 12px; 
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            transition: transform 0.2s ease;
        }
        .metric-card:hover { transform: translateY(-2px); }
        .metric-value { 
            font-size: 2.5rem; 
            font-weight: bold; 
            margin: 10px 0; 
        }
        .metric-label { 
            color: #7f8c8d; 
            font-size: 0.9rem; 
            text-transform: uppercase; 
            letter-spacing: 0.5px;
        }
        .status-indicator { 
            display: inline-block; 
            width: 12px; 
            height: 12px; 
            border-radius: 50%; 
            margin-right: 8px; 
        }
        .status-online { background: #27ae60; }
        .status-warning { background: #f39c12; }
        .status-offline { background: #e74c3c; }
        .excellent { color: #27ae60; }
        .good { color: #f39c12; }
        .poor { color: #e74c3c; }
        .chart-container { 
            background: white; 
            padding: 20px; 
            border-radius: 12px; 
            margin: 20px 0; 
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .alerts { 
            background: white; 
            padding: 20px; 
            border-radius: 12px; 
            margin: 20px 0; 
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .alert { 
            padding: 12px; 
            margin: 10px 0; 
            border-radius: 6px; 
            border-left: 4px solid;
        }
        .alert-critical { background: #fee; border-left-color: #e74c3c; }
        .alert-warning { background: #fff8e1; border-left-color: #f39c12; }
        .controls { 
            background: white; 
            padding: 20px; 
            border-radius: 12px; 
            margin: 20px 0; 
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
        }
        .btn { 
            background: #ff6136; 
            color: white; 
            border: none; 
            padding: 10px 20px; 
            border-radius: 6px; 
            cursor: pointer; 
            margin: 5px;
            font-size: 14px;
            transition: background 0.2s ease;
        }
        .btn:hover { background: #e55a2e; }
        .btn:disabled { background: #bdc3c7; cursor: not-allowed; }
        .timestamp { 
            color: #7f8c8d; 
            font-size: 0.8rem; 
            margin-top: 5px; 
        }
        .grid-2 { 
            display: grid; 
            grid-template-columns: 1fr 1fr; 
            gap: 20px; 
        }
        @media (max-width: 768px) {
            .grid-2 { grid-template-columns: 1fr; }
            .metrics-grid { grid-template-columns: 1fr; }
        }
        .loading { 
            display: inline-block; 
            width: 20px; 
            height: 20px; 
            border: 3px solid #f3f3f3;
            border-top: 3px solid #ff6136;
            border-radius: 50%;
            animation: spin 1s linear infinite;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="container">
            <h1>🔥 Fire Salamander - Dashboard UX Temps Réel</h1>
            <p>Monitoring SEPTEO Standards & User Experience</p>
            <p id="lastUpdate">Dernière mise à jour: Chargement...</p>
        </div>
    </div>

    <div class="container">
        <!-- Métriques principales -->
        <div class="metrics-grid">
            <div class="metric-card">
                <div class="metric-label">
                    <span id="systemStatus" class="status-indicator status-offline"></span>
                    Système Fire Salamander
                </div>
                <div id="systemValue" class="metric-value">--</div>
                <div class="timestamp" id="systemTimestamp">--</div>
            </div>
            
            <div class="metric-card">
                <div class="metric-label">🎨 Accessibilité WCAG</div>
                <div id="accessibilityValue" class="metric-value">--%</div>
                <div class="timestamp" id="accessibilityTimestamp">--</div>
            </div>
            
            <div class="metric-card">
                <div class="metric-label">⚡ Performance Lighthouse</div>
                <div id="performanceValue" class="metric-value">--%</div>
                <div class="timestamp" id="performanceTimestamp">--</div>
            </div>
            
            <div class="metric-card">
                <div class="metric-label">🎯 Design System SEPTEO</div>
                <div id="designSystemValue" class="metric-value">--%</div>
                <div class="timestamp" id="designSystemTimestamp">--</div>
            </div>
            
            <div class="metric-card">
                <div class="metric-label">👤 Expérience Utilisateur</div>
                <div id="userExperienceValue" class="metric-value">--%</div>
                <div class="timestamp" id="userExperienceTimestamp">--</div>
            </div>
            
            <div class="metric-card">
                <div class="metric-label">⏱️ Temps de Chargement</div>
                <div id="loadTimeValue" class="metric-value">--ms</div>
                <div class="timestamp" id="loadTimeTimestamp">--</div>
            </div>
        </div>

        <div class="grid-2">
            <!-- Graphique des métriques -->
            <div class="chart-container">
                <h3>📈 Évolution des Métriques (Temps Réel)</h3>
                <canvas id="metricsChart" width="400" height="200"></canvas>
            </div>

            <!-- Alertes -->
            <div class="alerts">
                <h3>🚨 Alertes UX</h3>
                <div id="alertsList">
                    <p>Chargement des alertes...</p>
                </div>
            </div>
        </div>

        <!-- Contrôles -->
        <div class="controls">
            <h3>🔧 Actions de Test</h3>
            <button class="btn" onclick="runTest('accessibility')">
                <span id="btn-accessibility">🎨 Test Accessibilité</span>
            </button>
            <button class="btn" onclick="runTest('performance')">
                <span id="btn-performance">⚡ Test Performance</span>
            </button>
            <button class="btn" onclick="runTest('design-system')">
                <span id="btn-design-system">🎯 Test Design System</span>
            </button>
            <button class="btn" onclick="runTest('visual')">
                <span id="btn-visual">👁️ Test Visuel</span>
            </button>
            <button class="btn" onclick="refreshMetrics()">
                🔄 Actualiser
            </button>
        </div>
    </div>

    <script>
        // Configuration du dashboard
        let metricsChart;
        let metricsHistory = [];
        
        // Initialiser le graphique
        function initChart() {
            const ctx = document.getElementById('metricsChart').getContext('2d');
            metricsChart = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: [],
                    datasets: [
                        {
                            label: 'Accessibilité',
                            data: [],
                            borderColor: '#3498db',
                            backgroundColor: 'rgba(52, 152, 219, 0.1)',
                            tension: 0.4
                        },
                        {
                            label: 'Performance', 
                            data: [],
                            borderColor: '#27ae60',
                            backgroundColor: 'rgba(39, 174, 96, 0.1)',
                            tension: 0.4
                        },
                        {
                            label: 'Design System',
                            data: [],
                            borderColor: '#ff6136',
                            backgroundColor: 'rgba(255, 97, 54, 0.1)',
                            tension: 0.4
                        },
                        {
                            label: 'UX',
                            data: [],
                            borderColor: '#f39c12',
                            backgroundColor: 'rgba(243, 156, 18, 0.1)',
                            tension: 0.4
                        }
                    ]
                },
                options: {
                    responsive: true,
                    scales: {
                        y: {
                            beginAtZero: true,
                            max: 100
                        }
                    },
                    plugins: {
                        legend: {
                            display: true,
                            position: 'top'
                        }
                    }
                }
            });
        }

        // Mettre à jour les métriques
        async function updateMetrics() {
            try {
                const response = await fetch('/api/metrics');
                const result = await response.json();
                
                if (result.success) {
                    const metrics = result.data;
                    updateUI(metrics);
                    updateChart(metrics);
                    
                    document.getElementById('lastUpdate').textContent = 
                        'Dernière mise à jour: ' + new Date().toLocaleTimeString('fr-FR');
                }
            } catch (error) {
                console.error('Error fetching metrics:', error);
            }
        }

        // Mettre à jour l'UI
        function updateUI(metrics) {
            // Système
            const systemEl = document.getElementById('systemValue');
            const statusEl = document.getElementById('systemStatus');
            
            if (metrics.firesamander?.status === 'healthy') {
                systemEl.textContent = 'En ligne';
                systemEl.className = 'metric-value excellent';
                statusEl.className = 'status-indicator status-online';
            } else {
                systemEl.textContent = 'Hors ligne';
                systemEl.className = 'metric-value poor';
                statusEl.className = 'status-indicator status-offline';
            }

            // Accessibilité
            updateMetricCard('accessibility', metrics.accessibility.score, '%');
            
            // Performance
            const perfScore = metrics.performance.lighthouse || 
                             (metrics.performance.loadTime < 3000 ? 90 : 60);
            updateMetricCard('performance', perfScore, '%');
            
            // Design System
            updateMetricCard('designSystem', metrics.designSystem.compliance, '%');
            
            // UX
            updateMetricCard('userExperience', metrics.userExperience.score, '%');
            
            // Temps de chargement
            document.getElementById('loadTimeValue').textContent = metrics.performance.loadTime + 'ms';
            document.getElementById('loadTimeValue').className = 
                'metric-value ' + (metrics.performance.loadTime < 2000 ? 'excellent' : 
                                  metrics.performance.loadTime < 5000 ? 'good' : 'poor');
        }

        function updateMetricCard(metric, value, unit = '') {
            const element = document.getElementById(metric + 'Value');
            element.textContent = value + unit;
            element.className = 'metric-value ' + getScoreClass(value);
        }

        function getScoreClass(score) {
            return score >= 90 ? 'excellent' : score >= 75 ? 'good' : 'poor';
        }

        // Mettre à jour le graphique
        function updateChart(metrics) {
            const now = new Date().toLocaleTimeString('fr-FR', { 
                hour: '2-digit', 
                minute: '2-digit' 
            });

            metricsHistory.push(metrics);
            if (metricsHistory.length > 20) {
                metricsHistory.shift();
            }

            const labels = metricsHistory.map((_, i) => 
                new Date(Date.now() - (metricsHistory.length - 1 - i) * 30000)
                    .toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })
            );

            metricsChart.data.labels = labels;
            metricsChart.data.datasets[0].data = metricsHistory.map(m => m.accessibility.score);
            metricsChart.data.datasets[1].data = metricsHistory.map(m => 
                m.performance.lighthouse || (m.performance.loadTime < 3000 ? 90 : 60)
            );
            metricsChart.data.datasets[2].data = metricsHistory.map(m => m.designSystem.compliance);
            metricsChart.data.datasets[3].data = metricsHistory.map(m => m.userExperience.score);
            
            metricsChart.update('none');
        }

        // Mettre à jour les alertes
        async function updateAlerts() {
            try {
                const response = await fetch('/api/alerts');
                const result = await response.json();
                
                if (result.success) {
                    const alertsList = document.getElementById('alertsList');
                    
                    if (result.data.length === 0) {
                        alertsList.innerHTML = '<p style="color: #27ae60;">✅ Aucune alerte active</p>';
                    } else {
                        alertsList.innerHTML = result.data.map(alert => 
                            '<div class="alert alert-' + alert.severity + '">' +
                            '<strong>' + alert.category.toUpperCase() + '</strong><br>' +
                            alert.description + '<br>' +
                            '<small>Seuil: ' + alert.threshold + '% | Actuel: ' + alert.currentScore + '%</small>' +
                            '</div>'
                        ).join('');
                    }
                }
            } catch (error) {
                console.error('Error fetching alerts:', error);
            }
        }

        // Lancer un test à la demande
        async function runTest(testType) {
            const button = document.getElementById('btn-' + testType);
            const originalText = button.innerHTML;
            
            button.innerHTML = '<div class="loading"></div> En cours...';
            button.parentElement.disabled = true;
            
            try {
                const response = await fetch('/api/tests/run', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ testType })
                });
                
                const result = await response.json();
                
                if (result.success) {
                    button.innerHTML = '✅ ' + originalText.split(' ').slice(1).join(' ');
                    setTimeout(() => {
                        button.innerHTML = originalText;
                        button.parentElement.disabled = false;
                        updateMetrics(); // Rafraîchir les métriques
                    }, 3000);
                } else {
                    button.innerHTML = '❌ Erreur';
                    setTimeout(() => {
                        button.innerHTML = originalText;
                        button.parentElement.disabled = false;
                    }, 3000);
                }
            } catch (error) {
                button.innerHTML = '❌ Erreur';
                setTimeout(() => {
                    button.innerHTML = originalText;
                    button.parentElement.disabled = false;
                }, 3000);
            }
        }

        // Actualiser manuellement
        function refreshMetrics() {
            updateMetrics();
            updateAlerts();
        }

        // Initialisation
        document.addEventListener('DOMContentLoaded', function() {
            initChart();
            updateMetrics();
            updateAlerts();
            
            // Mise à jour automatique toutes les 30 secondes
            setInterval(() => {
                updateMetrics();
                updateAlerts();
            }, 30000);
            
            // Connexion WebSocket pour les mises à jour temps réel
            const eventSource = new EventSource('/api/stream');
            eventSource.onmessage = function(event) {
                const data = JSON.parse(event.data);
                
                if (data.type === 'metrics') {
                    updateUI(data.data);
                    updateChart(data.data);
                } else if (data.type === 'alerts') {
                    updateAlerts();
                }
            };
        });
    </script>
</body>
</html>
    `;
  }

  startMetricsCollection() {
    // Collecter les métriques toutes les 30 secondes
    setInterval(async () => {
      try {
        await this.collectCurrentMetrics();
      } catch (error) {
        console.error('Error in metrics collection:', error);
      }
    }, 30000);
  }

  async start() {
    await fs.ensureDir(this.reportsDir);
    
    this.app.listen(this.port, () => {
      console.log(chalk.blue('🔥 Fire Salamander UX Dashboard started'));
      console.log(chalk.cyan(`   Dashboard: http://localhost:${this.port}`));
      console.log(chalk.cyan(`   API: http://localhost:${this.port}/api/metrics`));
      console.log(chalk.green('   Real-time monitoring enabled'));
    });
  }
}

// Démarrer le dashboard si exécuté directement
if (require.main === module) {
  const dashboard = new FireSalamanderUXDashboard();
  dashboard.start().catch(console.error);
}

module.exports = FireSalamanderUXDashboard;