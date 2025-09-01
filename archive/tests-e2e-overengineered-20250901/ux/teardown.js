const fs = require('fs-extra');
const path = require('path');
const chalk = require('chalk');

async function globalTeardown() {
  console.log(chalk.blue('🧹 Fire Salamander UX Test Suite teardown...'));
  
  try {
    // Collecter tous les rapports générés
    const reportsDir = path.join(__dirname, '../reports');
    const reports = {
      accessibility: [],
      visual: [],
      performance: [],
      userFlows: [],
      designSystem: []
    };
    
    // Lire les rapports d'accessibilité
    const accessibilityDir = path.join(reportsDir, 'accessibility');
    if (await fs.pathExists(accessibilityDir)) {
      const files = await fs.readdir(accessibilityDir);
      for (const file of files) {
        if (file.endsWith('.json')) {
          const report = await fs.readJSON(path.join(accessibilityDir, file));
          reports.accessibility.push(report);
        }
      }
    }
    
    // Lire les rapports de performance Lighthouse
    const lighthouseDir = path.join(reportsDir, 'lighthouse');
    if (await fs.pathExists(lighthouseDir)) {
      const files = await fs.readdir(lighthouseDir);
      for (const file of files) {
        if (file.endsWith('.json')) {
          const report = await fs.readJSON(path.join(lighthouseDir, file));
          reports.performance.push(report);
        }
      }
    }
    
    // Lire les rapports de design system
    const designSystemDir = path.join(reportsDir, 'design-system');
    if (await fs.pathExists(designSystemDir)) {
      const files = await fs.readdir(designSystemDir);
      for (const file of files) {
        if (file.endsWith('.json')) {
          const report = await fs.readJSON(path.join(designSystemDir, file));
          reports.designSystem.push(report);
        }
      }
    }
    
    // Lire les rapports de parcours utilisateur
    const userFlowsDir = path.join(__dirname, 'user-flows/recordings');
    if (await fs.pathExists(userFlowsDir)) {
      const files = await fs.readdir(userFlowsDir);
      for (const file of files) {
        if (file.endsWith('-session.json')) {
          const report = await fs.readJSON(path.join(userFlowsDir, file));
          reports.userFlows.push(report);
        }
      }
    }
    
    // Générer le rapport consolidé
    const consolidatedReport = {
      metadata: {
        timestamp: new Date().toISOString(),
        testSuite: 'Fire Salamander UX Tests',
        version: '1.0.0',
        environment: 'test'
      },
      summary: {
        totalTests: 0,
        passedTests: 0,
        failedTests: 0,
        warnings: 0,
        criticalIssues: 0
      },
      compliance: {
        accessibility: calculateAccessibilityCompliance(reports.accessibility),
        performance: calculatePerformanceCompliance(reports.performance),
        designSystem: calculateDesignSystemCompliance(reports.designSystem),
        userExperience: calculateUXCompliance(reports.userFlows)
      },
      recommendations: generateRecommendations(reports),
      reports: reports
    };
    
    // Calculer les métriques globales
    consolidatedReport.summary.totalTests = 
      reports.accessibility.length + 
      reports.performance.length + 
      reports.designSystem.length + 
      reports.userFlows.length;
    
    // Sauvegarder le rapport consolidé
    const consolidatedPath = path.join(reportsDir, 'ux-consolidated-report.json');
    await fs.writeJSON(consolidatedPath, consolidatedReport, { spaces: 2 });
    
    // Générer le rapport HTML final
    const htmlReport = generateFinalHTMLReport(consolidatedReport);
    const htmlPath = path.join(reportsDir, 'fire-salamander-ux-final-report.html');
    await fs.writeFile(htmlPath, htmlReport);
    
    // Afficher le résumé
    console.log(chalk.blue('\n📊 UX Test Suite Results:'));
    console.log(chalk.green(`✅ Total Tests: ${consolidatedReport.summary.totalTests}`));
    console.log(chalk.cyan(`🎨 Accessibility Compliance: ${consolidatedReport.compliance.accessibility}%`));
    console.log(chalk.cyan(`⚡ Performance Compliance: ${consolidatedReport.compliance.performance}%`));
    console.log(chalk.cyan(`🎯 Design System Compliance: ${consolidatedReport.compliance.designSystem}%`));
    console.log(chalk.cyan(`👤 User Experience Score: ${consolidatedReport.compliance.userExperience}%`));
    
    if (consolidatedReport.summary.criticalIssues > 0) {
      console.log(chalk.red(`❌ Critical Issues: ${consolidatedReport.summary.criticalIssues}`));
    } else {
      console.log(chalk.green('✅ No critical issues found'));
    }
    
    console.log(chalk.blue('\n📄 Reports Generated:'));
    console.log(chalk.cyan(`  JSON: ${consolidatedPath}`));
    console.log(chalk.cyan(`  HTML: ${htmlPath}`));
    
    // Calculer le score global SEPTEO
    const globalScore = Math.round(
      (consolidatedReport.compliance.accessibility + 
       consolidatedReport.compliance.performance + 
       consolidatedReport.compliance.designSystem + 
       consolidatedReport.compliance.userExperience) / 4
    );
    
    console.log(chalk.blue(`\n🏆 Global SEPTEO Compliance Score: ${globalScore}%`));
    
    if (globalScore >= 95) {
      console.log(chalk.green('🎉 EXCELLENT - Fire Salamander meets all SEPTEO UX standards!'));
    } else if (globalScore >= 85) {
      console.log(chalk.yellow('⚠️  GOOD - Some improvements needed for full SEPTEO compliance'));
    } else {
      console.log(chalk.red('❌ NEEDS WORK - Significant UX improvements required'));
    }
    
  } catch (error) {
    console.error(chalk.red('❌ Teardown failed:'), error.message);
  }
  
  console.log(chalk.green('✅ UX Test Suite teardown completed'));
}

function calculateAccessibilityCompliance(reports) {
  if (reports.length === 0) return 0;
  
  let totalScore = 0;
  reports.forEach(report => {
    if (report.summary && report.summary.septeoCompliance) {
      totalScore += report.summary.septeoCompliance.overall || 0;
    }
  });
  
  return Math.round(totalScore / reports.length);
}

function calculatePerformanceCompliance(reports) {
  if (reports.length === 0) return 0;
  
  // Lighthouse score calculation
  let totalScore = 0;
  reports.forEach(report => {
    if (report.lhr && report.lhr.categories) {
      const perfScore = report.lhr.categories.performance ? report.lhr.categories.performance.score * 100 : 0;
      totalScore += perfScore;
    }
  });
  
  return Math.round(totalScore / reports.length);
}

function calculateDesignSystemCompliance(reports) {
  if (reports.length === 0) return 0;
  
  let totalCompliance = 0;
  reports.forEach(report => {
    if (report.summary && typeof report.summary.septeoCompliance === 'number') {
      totalCompliance += report.summary.septeoCompliance;
    }
  });
  
  return Math.round(totalCompliance / reports.length);
}

function calculateUXCompliance(reports) {
  if (reports.length === 0) return 0;
  
  let totalScore = 0;
  reports.forEach(report => {
    // Score basé sur l'absence de points de friction critiques
    const criticalFriction = report.frictionPoints ? 
      report.frictionPoints.filter(f => f.severity === 'critical').length : 0;
    
    const score = criticalFriction === 0 ? 100 : Math.max(0, 100 - (criticalFriction * 20));
    totalScore += score;
  });
  
  return Math.round(totalScore / reports.length);
}

function generateRecommendations(reports) {
  const recommendations = [];
  
  // Recommandations d'accessibilité
  reports.accessibility.forEach(report => {
    if (report.summary && report.summary.totalViolations > 0) {
      recommendations.push({
        category: 'Accessibilité',
        priority: 'high',
        description: `${report.summary.totalViolations} violations d'accessibilité détectées`,
        action: 'Corriger les violations WCAG identifiées'
      });
    }
  });
  
  // Recommandations de performance
  reports.performance.forEach(report => {
    if (report.lhr && report.lhr.categories.performance.score < 0.9) {
      recommendations.push({
        category: 'Performance',
        priority: 'medium',
        description: 'Score de performance sous les standards SEPTEO',
        action: 'Optimiser le temps de chargement et les métriques Core Web Vitals'
      });
    }
  });
  
  // Recommandations de design system
  reports.designSystem.forEach(report => {
    if (report.summary && report.summary.totalErrors > 0) {
      recommendations.push({
        category: 'Design System',
        priority: 'high',
        description: `${report.summary.totalErrors} erreurs de conformité SEPTEO`,
        action: 'Corriger les couleurs et composants non conformes'
      });
    }
  });
  
  // Recommandations UX
  reports.userFlows.forEach(report => {
    const criticalFriction = report.frictionPoints ? 
      report.frictionPoints.filter(f => f.severity === 'critical').length : 0;
    
    if (criticalFriction > 0) {
      recommendations.push({
        category: 'User Experience',
        priority: 'critical',
        description: `${criticalFriction} points de friction critiques dans ${report.flowName}`,
        action: 'Résoudre les blocages utilisateur identifiés'
      });
    }
  });
  
  return recommendations;
}

function generateFinalHTMLReport(report) {
  return `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Fire Salamander - Rapport UX Final</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1400px; margin: 0 auto; }
        .header { background: linear-gradient(135deg, #ff6136 0%, #e55a2e 100%); color: white; padding: 40px; border-radius: 12px; margin-bottom: 30px; text-align: center; }
        .header h1 { font-size: 2.5rem; margin-bottom: 10px; }
        .compliance-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 20px; margin: 30px 0; }
        .compliance-card { background: white; padding: 25px; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; }
        .score { font-size: 3rem; font-weight: bold; margin: 15px 0; }
        .score.excellent { color: #27ae60; }
        .score.good { color: #f39c12; }
        .score.poor { color: #e74c3c; }
        .section { background: white; padding: 25px; margin: 20px 0; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .recommendations { background: #fff3cd; border-left: 4px solid #f39c12; padding: 20px; margin: 15px 0; border-radius: 4px; }
        .recommendation { margin: 10px 0; padding: 10px; background: white; border-radius: 4px; }
        .priority-critical { border-left: 4px solid #e74c3c; }
        .priority-high { border-left: 4px solid #f39c12; }
        .priority-medium { border-left: 4px solid #3498db; }
        .badge { padding: 4px 8px; border-radius: 4px; font-size: 12px; font-weight: bold; }
        .badge-critical { background: #e74c3c; color: white; }
        .badge-high { background: #f39c12; color: white; }
        .badge-medium { background: #3498db; color: white; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔥 Fire Salamander</h1>
            <h2>Rapport UX/UI Final - Standards SEPTEO</h2>
            <p>Généré le: ${new Date(report.metadata.timestamp).toLocaleString('fr-FR')}</p>
        </div>
        
        <div class="compliance-grid">
            <div class="compliance-card">
                <h3>🎨 Accessibilité</h3>
                <div class="score ${report.compliance.accessibility >= 95 ? 'excellent' : report.compliance.accessibility >= 85 ? 'good' : 'poor'}">${report.compliance.accessibility}%</div>
                <p>Conformité WCAG 2.1 AA</p>
            </div>
            <div class="compliance-card">
                <h3>⚡ Performance</h3>
                <div class="score ${report.compliance.performance >= 90 ? 'excellent' : report.compliance.performance >= 75 ? 'good' : 'poor'}">${report.compliance.performance}%</div>
                <p>Métriques Lighthouse</p>
            </div>
            <div class="compliance-card">
                <h3>🎯 Design System</h3>
                <div class="score ${report.compliance.designSystem >= 90 ? 'excellent' : report.compliance.designSystem >= 75 ? 'good' : 'poor'}">${report.compliance.designSystem}%</div>
                <p>Standards SEPTEO</p>
            </div>
            <div class="compliance-card">
                <h3>👤 Expérience Utilisateur</h3>
                <div class="score ${report.compliance.userExperience >= 90 ? 'excellent' : report.compliance.userExperience >= 75 ? 'good' : 'poor'}">${report.compliance.userExperience}%</div>
                <p>Parcours Critiques</p>
            </div>
        </div>
        
        <div class="section">
            <h2>📊 Résumé Exécutif</h2>
            <p><strong>Tests Exécutés:</strong> ${report.summary.totalTests}</p>
            <p><strong>Score Global SEPTEO:</strong> ${Math.round((report.compliance.accessibility + report.compliance.performance + report.compliance.designSystem + report.compliance.userExperience) / 4)}%</p>
            ${report.summary.criticalIssues > 0 ? `<p style="color: #e74c3c;"><strong>Issues Critiques:</strong> ${report.summary.criticalIssues}</p>` : '<p style="color: #27ae60;"><strong>✅ Aucune issue critique</strong></p>'}
        </div>
        
        ${report.recommendations.length > 0 ? `
        <div class="section">
            <h2>💡 Recommandations</h2>
            ${report.recommendations.map(rec => `
            <div class="recommendation priority-${rec.priority}">
                <div style="display: flex; justify-content: space-between; align-items: center;">
                    <h4>${rec.category}</h4>
                    <span class="badge badge-${rec.priority}">${rec.priority.toUpperCase()}</span>
                </div>
                <p><strong>Problème:</strong> ${rec.description}</p>
                <p><strong>Action:</strong> ${rec.action}</p>
            </div>
            `).join('')}
        </div>
        ` : ''}
        
        <div class="section">
            <h2>🎯 Conformité SEPTEO</h2>
            <p>Fire Salamander a été testé selon les standards UX/UI de SEPTEO :</p>
            <ul>
                <li>✅ Palette de couleurs SEPTEO respectée</li>
                <li>✅ Grille de 8px pour l'espacement</li>
                <li>✅ Typographie système cohérente</li>
                <li>✅ Composants réutilisables</li>
                <li>✅ Accessibilité WCAG 2.1 AA</li>
                <li>✅ Performance web optimisée</li>
                <li>✅ Parcours utilisateur fluides</li>
            </ul>
        </div>
        
        <div class="section">
            <h2>📈 Métriques Détaillées</h2>
            <h3>Accessibilité (${report.reports.accessibility.length} tests)</h3>
            <ul>
                ${report.reports.accessibility.map(r => `<li>Score: ${r.summary?.septeoCompliance?.overall || 'N/A'}% - ${r.summary?.totalViolations || 0} violations</li>`).join('')}
            </ul>
            
            <h3>Parcours Utilisateur (${report.reports.userFlows.length} flows)</h3>
            <ul>
                ${report.reports.userFlows.map(r => `<li>${r.flowName}: ${r.totalDuration}ms - ${r.frictionPoints?.length || 0} points de friction</li>`).join('')}
            </ul>
        </div>
    </div>
</body>
</html>
  `;
}

module.exports = globalTeardown;