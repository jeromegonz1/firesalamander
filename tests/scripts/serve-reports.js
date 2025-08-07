#!/usr/bin/env node

const express = require('express');
const path = require('path');
const fs = require('fs');
const chalk = require('chalk');

const app = express();
const PORT = 9003;

// Servir les fichiers statiques depuis le dossier tests
app.use('/tests', express.static(path.join(__dirname, '..')));

// Route pour les rapports BackstopJS
app.get('/backstop', (req, res) => {
  const reportPath = path.join(__dirname, '../backstop_data/html_report/index.html');
  if (fs.existsSync(reportPath)) {
    res.sendFile(reportPath);
  } else {
    res.status(404).send(`
      <h1>Rapport BackstopJS non trouvé</h1>
      <p>Le rapport n'a pas encore été généré.</p>
      <p>Lancez <code>npm run test:visual</code> pour créer le rapport.</p>
    `);
  }
});

// Route pour les rapports UX consolidés
app.get('/reports', (req, res) => {
  const reportsDir = path.join(__dirname, '../reports');
  const reports = {
    consolidated: [],
    accessibility: [],
    lighthouse: [],
    designSystem: []
  };

  // Scanner les rapports consolidés
  const consolidatedDir = path.join(reportsDir, 'consolidated');
  if (fs.existsSync(consolidatedDir)) {
    const files = fs.readdirSync(consolidatedDir);
    reports.consolidated = files.filter(f => f.endsWith('.html'));
  }

  // Scanner les rapports d'accessibilité  
  const accessibilityDir = path.join(reportsDir, 'accessibility');
  if (fs.existsSync(accessibilityDir)) {
    const files = fs.readdirSync(accessibilityDir);
    reports.accessibility = files.filter(f => f.endsWith('.html'));
  }

  // Scanner les rapports Lighthouse
  const lighthouseDir = path.join(reportsDir, 'lighthouse');
  if (fs.existsSync(lighthouseDir)) {
    const files = fs.readdirSync(lighthouseDir);
    reports.lighthouse = files.filter(f => f.endsWith('.html'));
  }

  // Scanner les rapports Design System
  const designSystemDir = path.join(reportsDir, 'design-system');
  if (fs.existsSync(designSystemDir)) {
    const files = fs.readdirSync(designSystemDir);
    reports.designSystem = files.filter(f => f.endsWith('.html'));
  }

  // Générer une page d'index
  const html = `
    <!DOCTYPE html>
    <html lang="fr">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Fire Salamander - Rapports UX/UI</title>
        <style>
            * { margin: 0; padding: 0; box-sizing: border-box; }
            body { 
                font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; 
                background: #f8f9fa; 
                padding: 20px;
            }
            .header { 
                background: linear-gradient(135deg, #ff6136 0%, #e55a2e 100%); 
                color: white; 
                padding: 20px; 
                border-radius: 12px;
                margin-bottom: 30px;
                text-align: center;
            }
            .section { 
                background: white; 
                padding: 20px; 
                border-radius: 12px; 
                margin-bottom: 20px;
                box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            }
            .section h2 { 
                color: #2c3e50; 
                margin-bottom: 15px;
                border-bottom: 2px solid #ff6136;
                padding-bottom: 5px;
            }
            .report-list { 
                list-style: none; 
            }
            .report-list li { 
                margin: 10px 0; 
            }
            .report-list a { 
                color: #ff6136; 
                text-decoration: none; 
                padding: 8px 12px;
                border: 1px solid #ff6136;
                border-radius: 6px;
                display: inline-block;
                transition: all 0.2s ease;
            }
            .report-list a:hover { 
                background: #ff6136; 
                color: white; 
            }
            .no-reports { 
                color: #7f8c8d; 
                font-style: italic; 
            }
            .quick-links {
                display: grid;
                grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
                gap: 15px;
                margin-top: 20px;
            }
            .quick-link {
                background: #ff6136;
                color: white;
                padding: 15px;
                text-align: center;
                text-decoration: none;
                border-radius: 8px;
                transition: background 0.2s ease;
            }
            .quick-link:hover {
                background: #e55a2e;
                color: white;
            }
        </style>
    </head>
    <body>
        <div class="header">
            <h1>🔥 Fire Salamander - Rapports UX/UI</h1>
            <p>Centre de reporting SEPTEO</p>
        </div>

        <div class="quick-links">
            <a href="/backstop" class="quick-link">📸 Rapport Visuel BackstopJS</a>
            <a href="http://localhost:9002" class="quick-link">📊 Dashboard Temps Réel</a>
            <a href="http://localhost:8080" class="quick-link">🔥 Fire Salamander</a>
        </div>

        <div class="section">
            <h2>📊 Rapports Consolidés</h2>
            ${reports.consolidated.length > 0 ? 
              `<ul class="report-list">
                ${reports.consolidated.map(file => 
                  `<li><a href="/tests/reports/consolidated/${file}">${file}</a></li>`
                ).join('')}
              </ul>` : 
              '<p class="no-reports">Aucun rapport consolidé disponible</p>'
            }
        </div>

        <div class="section">
            <h2>♿ Rapports d'Accessibilité</h2>
            ${reports.accessibility.length > 0 ? 
              `<ul class="report-list">
                ${reports.accessibility.map(file => 
                  `<li><a href="/tests/reports/accessibility/${file}">${file}</a></li>`
                ).join('')}
              </ul>` : 
              '<p class="no-reports">Aucun rapport d\'accessibilité disponible</p>'
            }
        </div>

        <div class="section">
            <h2>⚡ Rapports Lighthouse</h2>
            ${reports.lighthouse.length > 0 ? 
              `<ul class="report-list">
                ${reports.lighthouse.map(file => 
                  `<li><a href="/tests/reports/lighthouse/${file}">${file}</a></li>`
                ).join('')}
              </ul>` : 
              '<p class="no-reports">Aucun rapport Lighthouse disponible</p>'
            }
        </div>

        <div class="section">
            <h2>🎨 Rapports Design System</h2>
            ${reports.designSystem.length > 0 ? 
              `<ul class="report-list">
                ${reports.designSystem.map(file => 
                  `<li><a href="/tests/reports/design-system/${file}">${file}</a></li>`
                ).join('')}
              </ul>` : 
              '<p class="no-reports">Aucun rapport Design System disponible</p>'
            }
        </div>

        <div class="section">
            <h2>ℹ️ Information</h2>
            <p>Tous les rapports UX/UI de Fire Salamander sont centralisés ici.</p>
            <p><strong>Dashboard temps réel:</strong> <a href="http://localhost:9002">http://localhost:9002</a></p>
            <p><strong>Fire Salamander:</strong> <a href="http://localhost:8080">http://localhost:8080</a></p>
        </div>
    </body>
    </html>
  `;

  res.send(html);
});

// Route racine
app.get('/', (req, res) => {
  res.redirect('/reports');
});

app.listen(PORT, () => {
  console.log(chalk.blue('🔥 Fire Salamander Reports Server started'));
  console.log(chalk.cyan(`   Reports Hub: http://localhost:${PORT}`));
  console.log(chalk.cyan(`   BackstopJS: http://localhost:${PORT}/backstop`));
  console.log(chalk.cyan(`   All Reports: http://localhost:${PORT}/reports`));
  console.log(chalk.green('   Static files served from /tests'));
});

module.exports = app;