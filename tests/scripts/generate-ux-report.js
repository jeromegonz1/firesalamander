const fs = require('fs-extra');
const path = require('path');
const chalk = require('chalk');

// Générateur de rapport UX consolidé pour Fire Salamander
class UXReportGenerator {
  constructor() {
    this.reportsDir = path.join(__dirname, '../reports');
    this.outputDir = path.join(this.reportsDir, 'consolidated');
  }

  async generateWeeklyReport() {
    console.log(chalk.blue('📊 Generating weekly UX report for Fire Salamander...'));
    
    await fs.ensureDir(this.outputDir);
    
    const report = {
      metadata: {
        type: 'weekly',
        generatedAt: new Date().toISOString(),
        period: this.getWeekPeriod(),
        version: '1.0.0'
      },
      summary: {
        totalTests: 0,
        passedTests: 0,
        failedTests: 0,
        avgScore: 0,
        trends: {}
      },
      categories: {
        accessibility: await this.analyzeAccessibilityReports(),
        performance: await this.analyzePerformanceReports(),
        designSystem: await this.analyzeDesignSystemReports(),
        userExperience: await this.analyzeUserExperienceReports(),
        visualRegression: await this.analyzeVisualRegressionReports()
      },
      recommendations: [],
      alerts: [],
      septeoCompliance: {}
    };

    // Calculer les métriques globales
    this.calculateGlobalMetrics(report);
    
    // Générer les recommandations
    this.generateRecommendations(report);
    
    // Évaluer la conformité SEPTEO
    this.evaluateSepteoCompliance(report);
    
    // Sauvegarder le rapport JSON
    const jsonPath = path.join(this.outputDir, `weekly-ux-report-${this.getDateString()}.json`);
    await fs.writeJSON(jsonPath, report, { spaces: 2 });
    
    // Générer le rapport HTML
    const htmlReport = this.generateHTMLReport(report);
    const htmlPath = path.join(this.outputDir, `weekly-ux-report-${this.getDateString()}.html`);
    await fs.writeFile(htmlPath, htmlReport);
    
    // Générer le rapport exécutif PDF-ready
    const executiveReport = this.generateExecutiveReport(report);
    const execPath = path.join(this.outputDir, `executive-ux-summary-${this.getDateString()}.html`);
    await fs.writeFile(execPath, executiveReport);
    
    console.log(chalk.green('✅ Weekly UX report generated:'));
    console.log(chalk.cyan(`   JSON: ${jsonPath}`));
    console.log(chalk.cyan(`   HTML: ${htmlPath}`));
    console.log(chalk.cyan(`   Executive: ${execPath}`));
    
    return report;
  }

  async analyzeAccessibilityReports() {
    const category = {
      totalReports: 0,
      avgScore: 0,
      totalViolations: 0,
      trends: [],
      criticalIssues: []
    };

    try {
      const accessibilityDir = path.join(this.reportsDir, 'accessibility');
      if (await fs.pathExists(accessibilityDir)) {
        const files = await fs.readdir(accessibilityDir);
        const reportFiles = files.filter(f => f.endsWith('.json') && f.includes('report'));
        
        let totalScore = 0;
        let totalViolations = 0;
        
        for (const file of reportFiles) {
          const report = await fs.readJSON(path.join(accessibilityDir, file));
          category.totalReports++;
          
          if (report.summary) {
            const score = report.summary.septeoCompliance?.overall || 0;
            totalScore += score;
            totalViolations += report.summary.totalViolations || 0;
            
            // Collecter les violations critiques
            if (report.summary.totalViolations > 0) {
              category.criticalIssues.push({
                file,
                violations: report.summary.totalViolations,
                score: score,
                timestamp: report.summary.timestamp
              });
            }
          }
        }
        
        if (category.totalReports > 0) {
          category.avgScore = Math.round(totalScore / category.totalReports);
          category.totalViolations = totalViolations;
        }
      }
    } catch (error) {
      console.warn('Error analyzing accessibility reports:', error.message);
    }

    return category;
  }

  async analyzePerformanceReports() {
    const category = {
      totalReports: 0,
      avgScore: 0,
      avgLoadTime: 0,
      coreWebVitals: {
        lcp: 0,
        fid: 0,
        cls: 0
      },
      trends: [],
      issues: []
    };

    try {
      const performanceDir = path.join(this.reportsDir, 'lighthouse');
      if (await fs.pathExists(performanceDir)) {
        const files = await fs.readdir(performanceDir);
        const reportFiles = files.filter(f => f.endsWith('.json'));
        
        let totalScore = 0;
        let totalLoadTime = 0;
        let totalLCP = 0, totalFID = 0, totalCLS = 0;
        
        for (const file of reportFiles) {
          try {
            const report = await fs.readJSON(path.join(performanceDir, file));
            category.totalReports++;
            
            if (report.lhr && report.lhr.categories && report.lhr.categories.performance) {
              const score = report.lhr.categories.performance.score * 100;
              totalScore += score;
              
              // Extraire les métriques Core Web Vitals
              const audits = report.lhr.audits;
              if (audits['largest-contentful-paint']) {
                totalLCP += audits['largest-contentful-paint'].numericValue || 0;
              }
              if (audits['first-input-delay']) {
                totalFID += audits['first-input-delay'].numericValue || 0;
              }
              if (audits['cumulative-layout-shift']) {
                totalCLS += audits['cumulative-layout-shift'].numericValue || 0;
              }
              
              // Identifier les problèmes de performance
              if (score < 75) {
                category.issues.push({
                  file,
                  score: Math.round(score),
                  url: report.lhr.finalUrl,
                  timestamp: report.lhr.fetchTime
                });
              }
            }
          } catch (e) {
            console.warn(`Error reading performance report ${file}:`, e.message);
          }
        }
        
        if (category.totalReports > 0) {
          category.avgScore = Math.round(totalScore / category.totalReports);
          category.avgLoadTime = Math.round(totalLoadTime / category.totalReports);
          category.coreWebVitals.lcp = Math.round(totalLCP / category.totalReports);
          category.coreWebVitals.fid = Math.round(totalFID / category.totalReports);
          category.coreWebVitals.cls = Math.round((totalCLS / category.totalReports) * 100) / 100;
        }
      }
    } catch (error) {
      console.warn('Error analyzing performance reports:', error.message);
    }

    return category;
  }

  async analyzeDesignSystemReports() {
    const category = {
      totalReports: 0,
      complianceScore: 0,
      colorViolations: 0,
      spacingViolations: 0,
      typographyViolations: 0,
      issues: []
    };

    try {
      const designDir = path.join(this.reportsDir, 'design-system');
      if (await fs.pathExists(designDir)) {
        const files = await fs.readdir(designDir);
        const reportFiles = files.filter(f => f.endsWith('.json') && f.includes('report'));
        
        let totalCompliance = 0;
        let totalColorViolations = 0;
        let totalSpacingViolations = 0;
        let totalTypographyViolations = 0;
        
        for (const file of reportFiles) {
          const report = await fs.readJSON(path.join(designDir, file));
          category.totalReports++;
          
          if (report.summary) {
            totalCompliance += report.summary.septeoCompliance || 0;
            
            // Compter les violations par type
            if (report.errors) {
              const colorErrors = report.errors.filter(e => e.type.includes('color')).length;
              const spacingErrors = report.errors.filter(e => e.type.includes('spacing')).length;
              const typographyErrors = report.errors.filter(e => e.type.includes('font')).length;
              
              totalColorViolations += colorErrors;
              totalSpacingViolations += spacingErrors; 
              totalTypographyViolations += typographyErrors;
              
              if (report.errors.length > 0) {
                category.issues.push({
                  file,
                  totalErrors: report.errors.length,
                  colorErrors,
                  spacingErrors,
                  typographyErrors,
                  timestamp: report.summary.timestamp
                });
              }
            }
          }
        }
        
        if (category.totalReports > 0) {
          category.complianceScore = Math.round(totalCompliance / category.totalReports);
          category.colorViolations = totalColorViolations;
          category.spacingViolations = totalSpacingViolations;
          category.typographyViolations = totalTypographyViolations;
        }
      }
    } catch (error) {
      console.warn('Error analyzing design system reports:', error.message);
    }

    return category;
  }

  async analyzeUserExperienceReports() {
    const category = {
      totalFlows: 0,
      avgCompletionTime: 0,
      successRate: 0,
      frictionPoints: 0,
      criticalFriction: 0,
      flowAnalysis: []
    };

    try {
      const userFlowsDir = path.join(__dirname, '../ux/user-flows/recordings');
      if (await fs.pathExists(userFlowsDir)) {
        const files = await fs.readdir(userFlowsDir);
        const sessionFiles = files.filter(f => f.endsWith('-session.json'));
        
        let totalTime = 0;
        let successfulFlows = 0;
        let totalFriction = 0;
        let totalCriticalFriction = 0;
        
        for (const file of sessionFiles) {
          const session = await fs.readJSON(path.join(userFlowsDir, file));
          category.totalFlows++;
          
          totalTime += session.totalDuration || 0;
          
          const frictionCount = session.frictionPoints?.length || 0;
          const criticalCount = session.frictionPoints?.filter(f => f.severity === 'critical').length || 0;
          
          totalFriction += frictionCount;
          totalCriticalFriction += criticalCount;
          
          // Considérer comme succès si pas de friction critique
          if (criticalCount === 0) {
            successfulFlows++;
          }
          
          category.flowAnalysis.push({
            flowName: session.flowName,
            duration: session.totalDuration,
            frictionPoints: frictionCount,
            criticalFriction: criticalCount,
            success: criticalCount === 0,
            timestamp: session.startTime
          });
        }
        
        if (category.totalFlows > 0) {
          category.avgCompletionTime = Math.round(totalTime / category.totalFlows);
          category.successRate = Math.round((successfulFlows / category.totalFlows) * 100);
          category.frictionPoints = totalFriction;
          category.criticalFriction = totalCriticalFriction;
        }
      }
    } catch (error) {
      console.warn('Error analyzing user experience reports:', error.message);
    }

    return category;
  }

  async analyzeVisualRegressionReports() {
    const category = {
      totalTests: 0,
      passedTests: 0,
      failedTests: 0,
      avgPixelDiff: 0,
      regressions: []
    };

    try {
      const visualDir = path.join(__dirname, '../ux/visual-regression/backstop_data');
      const reportPath = path.join(visualDir, 'html_report/config.js');
      
      if (await fs.pathExists(reportPath)) {
        // BackstopJS génère un fichier de configuration avec les résultats
        const configContent = await fs.readFile(reportPath, 'utf8');
        
        // Parser les résultats (simplifié)
        const passMatches = configContent.match(/"passed":\s*(\d+)/);
        const failMatches = configContent.match(/"failed":\s*(\d+)/);
        
        if (passMatches) category.passedTests = parseInt(passMatches[1]);
        if (failMatches) category.failedTests = parseInt(failMatches[1]);
        category.totalTests = category.passedTests + category.failedTests;
      }
    } catch (error) {
      console.warn('Error analyzing visual regression reports:', error.message);
    }

    return category;
  }

  calculateGlobalMetrics(report) {
    const categories = report.categories;
    
    // Calculer les totaux
    report.summary.totalTests = 
      categories.accessibility.totalReports +
      categories.performance.totalReports +
      categories.designSystem.totalReports +
      categories.userExperience.totalFlows +
      categories.visualRegression.totalTests;
    
    report.summary.passedTests = 
      categories.accessibility.totalReports - categories.accessibility.criticalIssues.length +
      categories.performance.totalReports - categories.performance.issues.length +
      categories.designSystem.totalReports - categories.designSystem.issues.length +
      Math.round(categories.userExperience.totalFlows * (categories.userExperience.successRate / 100)) +
      categories.visualRegression.passedTests;
    
    report.summary.failedTests = report.summary.totalTests - report.summary.passedTests;
    
    // Score global pondéré
    const scores = [];
    if (categories.accessibility.avgScore > 0) scores.push(categories.accessibility.avgScore);
    if (categories.performance.avgScore > 0) scores.push(categories.performance.avgScore);
    if (categories.designSystem.complianceScore > 0) scores.push(categories.designSystem.complianceScore);
    if (categories.userExperience.successRate > 0) scores.push(categories.userExperience.successRate);
    
    report.summary.avgScore = scores.length > 0 ? 
      Math.round(scores.reduce((a, b) => a + b, 0) / scores.length) : 0;
  }

  generateRecommendations(report) {
    const recommendations = [];
    
    // Recommandations d'accessibilité
    if (report.categories.accessibility.avgScore < 95) {
      recommendations.push({
        priority: report.categories.accessibility.avgScore < 85 ? 'critical' : 'high',
        category: 'Accessibilité',
        title: 'Améliorer la conformité WCAG',
        description: `Score actuel: ${report.categories.accessibility.avgScore}% (objectif: 95%+)`,
        action: 'Corriger les violations d\'accessibilité identifiées',
        impact: 'Utilisabilité pour tous les utilisateurs',
        effort: 'Moyen',
        deadline: '2 semaines'
      });
    }
    
    // Recommandations de performance
    if (report.categories.performance.avgScore < 90) {
      recommendations.push({
        priority: report.categories.performance.avgScore < 75 ? 'critical' : 'high',
        category: 'Performance',
        title: 'Optimiser les performances',
        description: `Score Lighthouse: ${report.categories.performance.avgScore}% (objectif: 90%+)`,
        action: 'Optimiser les Core Web Vitals et temps de chargement',
        impact: 'Expérience utilisateur et SEO',
        effort: 'Élevé',
        deadline: '1 mois'
      });
    }
    
    // Recommandations design system
    if (report.categories.designSystem.complianceScore < 90) {
      recommendations.push({
        priority: 'high',
        category: 'Design System',
        title: 'Respecter les standards SEPTEO',
        description: `Conformité: ${report.categories.designSystem.complianceScore}% (objectif: 90%+)`,
        action: 'Corriger les violations de couleurs, espacement et typographie',
        impact: 'Cohérence de marque SEPTEO',
        effort: 'Moyen',
        deadline: '2 semaines'
      });
    }
    
    // Recommandations UX
    if (report.categories.userExperience.criticalFriction > 0) {
      recommendations.push({
        priority: 'critical',
        category: 'Expérience Utilisateur',
        title: 'Éliminer les points de friction critiques',
        description: `${report.categories.userExperience.criticalFriction} points de friction critiques détectés`,
        action: 'Résoudre les blocages dans les parcours utilisateur',
        impact: 'Conversion et satisfaction utilisateur',
        effort: 'Élevé',
        deadline: '1 semaine'
      });
    }
    
    report.recommendations = recommendations;
  }

  evaluateSepteoCompliance(report) {
    const compliance = {
      overall: 0,
      accessibility: {
        score: report.categories.accessibility.avgScore,
        status: report.categories.accessibility.avgScore >= 95 ? 'excellent' : 
                report.categories.accessibility.avgScore >= 85 ? 'good' : 'needs_improvement',
        requirements: ['WCAG 2.1 AA', 'Contraste couleurs', 'Navigation clavier']
      },
      performance: {
        score: report.categories.performance.avgScore,
        status: report.categories.performance.avgScore >= 90 ? 'excellent' :
                report.categories.performance.avgScore >= 75 ? 'good' : 'needs_improvement',
        requirements: ['LCP < 2.5s', 'FID < 100ms', 'CLS < 0.1']
      },
      designSystem: {
        score: report.categories.designSystem.complianceScore,
        status: report.categories.designSystem.complianceScore >= 90 ? 'excellent' :
                report.categories.designSystem.complianceScore >= 75 ? 'good' : 'needs_improvement',
        requirements: ['Couleurs SEPTEO', 'Grille 8px', 'Typographie système']
      },
      userExperience: {
        score: report.categories.userExperience.successRate,
        status: report.categories.userExperience.successRate >= 95 ? 'excellent' :
                report.categories.userExperience.successRate >= 85 ? 'good' : 'needs_improvement',
        requirements: ['Parcours fluides', 'Zéro friction critique', 'Temps optimaux']
      }
    };
    
    // Score global SEPTEO
    compliance.overall = Math.round(
      (compliance.accessibility.score + 
       compliance.performance.score + 
       compliance.designSystem.score + 
       compliance.userExperience.score) / 4
    );
    
    report.septeoCompliance = compliance;
  }

  generateHTMLReport(report) {
    return `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Fire Salamander - Rapport UX Hebdomadaire</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 20px; background: #f8f9fa; }
        .container { max-width: 1400px; margin: 0 auto; }
        .header { background: linear-gradient(135deg, #ff6136 0%, #e55a2e 100%); color: white; padding: 40px; border-radius: 12px; margin-bottom: 30px; text-align: center; }
        .summary-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin: 30px 0; }
        .summary-card { background: white; padding: 25px; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); text-align: center; }
        .score { font-size: 3rem; font-weight: bold; margin: 15px 0; }
        .score.excellent { color: #27ae60; }
        .score.good { color: #f39c12; }
        .score.poor { color: #e74c3c; }
        .section { background: white; padding: 25px; margin: 20px 0; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .category-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin: 20px 0; }
        .category-card { background: #f8f9fa; padding: 20px; border-radius: 6px; }
        .recommendations { background: #fff3cd; border-left: 4px solid #f39c12; padding: 20px; margin: 15px 0; border-radius: 4px; }
        .recommendation { margin: 15px 0; padding: 15px; background: white; border-radius: 4px; }
        .priority-critical { border-left: 4px solid #e74c3c; }
        .priority-high { border-left: 4px solid #f39c12; }
        .priority-medium { border-left: 4px solid #3498db; }
        .badge { padding: 4px 8px; border-radius: 4px; font-size: 12px; font-weight: bold; color: white; }
        .badge-critical { background: #e74c3c; }
        .badge-high { background: #f39c12; }
        .badge-medium { background: #3498db; }
        .badge-excellent { background: #27ae60; }
        .badge-good { background: #f39c12; }
        .badge-poor { background: #e74c3c; }
        .chart { height: 200px; background: #f8f9fa; border-radius: 4px; display: flex; align-items: center; justify-content: center; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔥 Fire Salamander</h1>
            <h2>Rapport UX/UI Hebdomadaire</h2>
            <p>Standards SEPTEO & Expérience Utilisateur</p>
            <p><strong>Période:</strong> ${report.metadata.period}</p>
            <p>Généré le: ${new Date(report.metadata.generatedAt).toLocaleString('fr-FR')}</p>
        </div>
        
        <div class="summary-grid">
            <div class="summary-card">
                <h3>Score Global SEPTEO</h3>
                <div class="score ${this.getScoreClass(report.septeoCompliance.overall)}">${report.septeoCompliance.overall}%</div>
                <p>Conformité générale</p>
            </div>
            <div class="summary-card">
                <h3>Tests Exécutés</h3>
                <div class="score">${report.summary.totalTests}</div>
                <p>Total cette semaine</p>
            </div>
            <div class="summary-card">
                <h3>Taux de Réussite</h3>
                <div class="score ${this.getScoreClass(Math.round(report.summary.passedTests / report.summary.totalTests * 100))}">${Math.round(report.summary.passedTests / report.summary.totalTests * 100)}%</div>
                <p>${report.summary.passedTests}/${report.summary.totalTests} tests</p>
            </div>
            <div class="summary-card">
                <h3>Recommandations</h3>
                <div class="score">${report.recommendations.length}</div>
                <p>Actions à entreprendre</p>
            </div>
        </div>
        
        <div class="section">
            <h2>🎯 Conformité SEPTEO par Catégorie</h2>
            <div class="category-grid">
                <div class="category-card">
                    <h3>🎨 Accessibilité</h3>
                    <div class="score ${this.getScoreClass(report.septeoCompliance.accessibility.score)}">${report.septeoCompliance.accessibility.score}%</div>
                    <span class="badge badge-${report.septeoCompliance.accessibility.status}">${report.septeoCompliance.accessibility.status.toUpperCase()}</span>
                    <p><strong>Violations:</strong> ${report.categories.accessibility.totalViolations}</p>
                    <p><strong>Tests:</strong> ${report.categories.accessibility.totalReports}</p>
                </div>
                
                <div class="category-card">
                    <h3>⚡ Performance</h3>
                    <div class="score ${this.getScoreClass(report.septeoCompliance.performance.score)}">${report.septeoCompliance.performance.score}%</div>
                    <span class="badge badge-${report.septeoCompliance.performance.status}">${report.septeoCompliance.performance.status.toUpperCase()}</span>
                    <p><strong>LCP:</strong> ${report.categories.performance.coreWebVitals.lcp}ms</p>
                    <p><strong>CLS:</strong> ${report.categories.performance.coreWebVitals.cls}</p>
                </div>
                
                <div class="category-card">
                    <h3>🎯 Design System</h3>
                    <div class="score ${this.getScoreClass(report.septeoCompliance.designSystem.score)}">${report.septeoCompliance.designSystem.score}%</div>
                    <span class="badge badge-${report.septeoCompliance.designSystem.status}">${report.septeoCompliance.designSystem.status.toUpperCase()}</span>
                    <p><strong>Violations couleur:</strong> ${report.categories.designSystem.colorViolations}</p>
                    <p><strong>Violations espacement:</strong> ${report.categories.designSystem.spacingViolations}</p>
                </div>
                
                <div class="category-card">
                    <h3>👤 Expérience Utilisateur</h3>
                    <div class="score ${this.getScoreClass(report.septeoCompliance.userExperience.score)}">${report.septeoCompliance.userExperience.score}%</div>
                    <span class="badge badge-${report.septeoCompliance.userExperience.status}">${report.septeoCompliance.userExperience.status.toUpperCase()}</span>
                    <p><strong>Friction critique:</strong> ${report.categories.userExperience.criticalFriction}</p>
                    <p><strong>Temps moyen:</strong> ${Math.round(report.categories.userExperience.avgCompletionTime / 1000)}s</p>
                </div>
            </div>
        </div>
        
        ${report.recommendations.length > 0 ? `
        <div class="section">
            <h2>💡 Recommandations Prioritaires</h2>
            ${report.recommendations.map(rec => `
            <div class="recommendation priority-${rec.priority}">
                <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
                    <h3>${rec.title}</h3>
                    <span class="badge badge-${rec.priority}">${rec.priority.toUpperCase()}</span>
                </div>
                <p><strong>Catégorie:</strong> ${rec.category}</p>
                <p><strong>Description:</strong> ${rec.description}</p>
                <p><strong>Action:</strong> ${rec.action}</p>
                <p><strong>Impact:</strong> ${rec.impact}</p>
                <p><strong>Effort:</strong> ${rec.effort} | <strong>Délai:</strong> ${rec.deadline}</p>
            </div>
            `).join('')}
        </div>
        ` : ''}
        
        <div class="section">
            <h2>📈 Analyse Détaillée</h2>
            
            <h3>Parcours Utilisateur</h3>
            ${report.categories.userExperience.flowAnalysis.length > 0 ? `
            <table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
                <thead>
                    <tr style="background: #f8f9fa;">
                        <th style="padding: 12px; text-align: left; border: 1px solid #dee2e6;">Parcours</th>
                        <th style="padding: 12px; text-align: left; border: 1px solid #dee2e6;">Durée</th>
                        <th style="padding: 12px; text-align: left; border: 1px solid #dee2e6;">Friction</th>
                        <th style="padding: 12px; text-align: left; border: 1px solid #dee2e6;">Statut</th>
                    </tr>
                </thead>
                <tbody>
                    ${report.categories.userExperience.flowAnalysis.map(flow => `
                    <tr>
                        <td style="padding: 12px; border: 1px solid #dee2e6;">${flow.flowName}</td>
                        <td style="padding: 12px; border: 1px solid #dee2e6;">${Math.round(flow.duration / 1000)}s</td>
                        <td style="padding: 12px; border: 1px solid #dee2e6;">${flow.frictionPoints} (${flow.criticalFriction} critique)</td>
                        <td style="padding: 12px; border: 1px solid #dee2e6;">
                            <span class="badge badge-${flow.success ? 'excellent' : 'poor'}">${flow.success ? 'SUCCÈS' : 'ÉCHEC'}</span>
                        </td>
                    </tr>
                    `).join('')}
                </tbody>
            </table>
            ` : '<p>Aucune donnée de parcours utilisateur disponible.</p>'}
        </div>
        
        <div class="section">
            <h2>🔍 Méthodologie</h2>
            <p>Ce rapport consolide les résultats des tests automatisés suivants :</p>
            <ul>
                <li><strong>Tests d'Accessibilité:</strong> Axe-core + Pa11y (WCAG 2.1 AA)</li>
                <li><strong>Tests de Performance:</strong> Lighthouse CI (Core Web Vitals)</li>
                <li><strong>Validation Design System:</strong> Conformité couleurs, espacement, typographie SEPTEO</li>
                <li><strong>Tests Parcours Utilisateur:</strong> Playwright + enregistrement des interactions</li>
                <li><strong>Tests Visuels:</strong> BackstopJS (régression visuelle)</li>
            </ul>
            
            <h3>Seuils SEPTEO</h3>
            <ul>
                <li>Accessibilité: ≥95% (zéro violation critique)</li>
                <li>Performance: ≥90% Lighthouse (LCP <2.5s, FID <100ms, CLS <0.1)</li>
                <li>Design System: ≥90% conformité (couleurs + grille 8px)</li>
                <li>UX: ≥95% succès parcours (zéro friction critique)</li>
            </ul>
        </div>
    </div>
</body>
</html>
    `;
  }

  generateExecutiveReport(report) {
    const statusText = report.septeoCompliance.overall >= 95 ? 'EXCELLENT' :
                       report.septeoCompliance.overall >= 85 ? 'SATISFAISANT' : 'DOIT ÊTRE AMÉLIORÉ';
    
    return `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Fire Salamander - Résumé Exécutif UX</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 40px; background: white; color: #2c3e50; line-height: 1.6; }
        .container { max-width: 800px; margin: 0 auto; }
        .header { text-align: center; margin-bottom: 40px; border-bottom: 2px solid #ff6136; padding-bottom: 20px; }
        .logo { color: #ff6136; font-size: 2.5rem; font-weight: bold; }
        .status-banner { background: ${report.septeoCompliance.overall >= 95 ? '#d4edda' : report.septeoCompliance.overall >= 85 ? '#fff3cd' : '#f8d7da'}; 
                        color: ${report.septeoCompliance.overall >= 95 ? '#155724' : report.septeoCompliance.overall >= 85 ? '#856404' : '#721c24'};
                        padding: 20px; border-radius: 8px; text-align: center; margin: 20px 0; font-weight: bold; font-size: 1.2rem; }
        .metrics { display: grid; grid-template-columns: repeat(2, 1fr); gap: 20px; margin: 30px 0; }
        .metric { background: #f8f9fa; padding: 20px; border-radius: 8px; text-align: center; }
        .metric-value { font-size: 2rem; font-weight: bold; color: #ff6136; }
        .section { margin: 30px 0; }
        .priorities { background: #fff3cd; padding: 20px; border-radius: 8px; border-left: 4px solid #f39c12; }
        .priority-item { margin: 10px 0; padding: 10px; background: white; border-radius: 4px; }
        @media print { body { background: white; } .container { max-width: none; } }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <div class="logo">🔥 FIRE SALAMANDER</div>
            <h1>Résumé Exécutif UX/UI</h1>
            <p><strong>Période:</strong> ${report.metadata.period} | <strong>Date:</strong> ${new Date(report.metadata.generatedAt).toLocaleDateString('fr-FR')}</p>
        </div>
        
        <div class="status-banner">
            STATUT GLOBAL SEPTEO: ${statusText} (${report.septeoCompliance.overall}%)
        </div>
        
        <div class="section">
            <h2>📊 Vue d'Ensemble</h2>
            <div class="metrics">
                <div class="metric">
                    <div class="metric-value">${report.summary.totalTests}</div>
                    <div>Tests Exécutés</div>
                </div>
                <div class="metric">
                    <div class="metric-value">${Math.round(report.summary.passedTests / report.summary.totalTests * 100)}%</div>
                    <div>Taux de Réussite</div>
                </div>
                <div class="metric">
                    <div class="metric-value">${report.categories.accessibility.totalViolations}</div>
                    <div>Violations Accessibilité</div>
                </div>
                <div class="metric">
                    <div class="metric-value">${report.categories.userExperience.criticalFriction}</div>
                    <div>Friction Critique UX</div>
                </div>
            </div>
        </div>
        
        <div class="section">
            <h2>🎯 Conformité SEPTEO</h2>
            <ul>
                <li><strong>Accessibilité WCAG 2.1 AA:</strong> ${report.septeoCompliance.accessibility.score}% 
                    (${report.septeoCompliance.accessibility.status === 'excellent' ? '✅ Excellent' : 
                       report.septeoCompliance.accessibility.status === 'good' ? '⚠️ Satisfaisant' : '❌ À améliorer'})</li>
                <li><strong>Performance Lighthouse:</strong> ${report.septeoCompliance.performance.score}% 
                    (${report.septeoCompliance.performance.status === 'excellent' ? '✅ Excellent' : 
                       report.septeoCompliance.performance.status === 'good' ? '⚠️ Satisfaisant' : '❌ À améliorer'})</li>
                <li><strong>Design System SEPTEO:</strong> ${report.septeoCompliance.designSystem.score}% 
                    (${report.septeoCompliance.designSystem.status === 'excellent' ? '✅ Excellent' : 
                       report.septeoCompliance.designSystem.status === 'good' ? '⚠️ Satisfaisant' : '❌ À améliorer'})</li>
                <li><strong>Expérience Utilisateur:</strong> ${report.septeoCompliance.userExperience.score}% 
                    (${report.septeoCompliance.userExperience.status === 'excellent' ? '✅ Excellent' : 
                       report.septeoCompliance.userExperience.status === 'good' ? '⚠️ Satisfaisant' : '❌ À améliorer'})</li>
            </ul>
        </div>
        
        ${report.recommendations.length > 0 ? `
        <div class="section">
            <div class="priorities">
                <h2>⚡ Actions Prioritaires</h2>
                ${report.recommendations.slice(0, 3).map((rec, i) => `
                <div class="priority-item">
                    <strong>${i + 1}. ${rec.title}</strong><br>
                    <em>${rec.category}</em> - Priorité: ${rec.priority.toUpperCase()}<br>
                    ${rec.description}<br>
                    <small>Délai: ${rec.deadline} | Effort: ${rec.effort}</small>
                </div>
                `).join('')}
                ${report.recommendations.length > 3 ? `<p><em>... et ${report.recommendations.length - 3} autres recommandations (voir rapport détaillé)</em></p>` : ''}
            </div>
        </div>
        ` : ''}
        
        <div class="section">
            <h2>📈 Impact Business</h2>
            <p><strong>Utilisabilité:</strong> ${report.categories.accessibility.totalViolations === 0 ? 'Excellente' : 'À améliorer'} - 
               ${report.categories.accessibility.totalViolations} violations d'accessibilité détectées</p>
            <p><strong>Performance:</strong> Temps de chargement moyen ${report.categories.performance.avgLoadTime || 'N/A'}ms 
               (objectif: <2000ms)</p>
            <p><strong>Conversion:</strong> ${report.categories.userExperience.successRate}% des parcours utilisateur réussis 
               (objectif: 95%+)</p>
            <p><strong>Conformité Marque:</strong> ${report.septeoCompliance.designSystem.score}% de conformité aux standards SEPTEO 
               (objectif: 90%+)</p>
        </div>
        
        <div class="section">
            <h2>🎯 Recommandations Stratégiques</h2>
            ${report.septeoCompliance.overall >= 95 ? 
              '<p>✅ <strong>Fire Salamander respecte excellemment les standards SEPTEO.</strong> Maintenir ce niveau d\'excellence et surveiller les régressions.</p>' :
              report.septeoCompliance.overall >= 85 ?
              '<p>⚠️ <strong>Fire Salamander atteint un niveau satisfaisant.</strong> Quelques améliorations ciblées permettront d\'atteindre l\'excellence SEPTEO.</p>' :
              '<p>❌ <strong>Fire Salamander nécessite des améliorations significatives.</strong> Un plan d\'action immédiat est recommandé pour respecter les standards SEPTEO.</p>'
            }
        </div>
        
        <div class="section" style="margin-top: 40px; padding-top: 20px; border-top: 1px solid #dee2e6; font-size: 0.9rem; color: #6c757d;">
            <p><strong>Note:</strong> Ce résumé est généré automatiquement à partir des tests UX/UI continus. 
               Rapport détaillé disponible avec métriques complètes et recommandations techniques.</p>
            <p><strong>Prochaine évaluation:</strong> ${new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toLocaleDateString('fr-FR')}</p>
        </div>
    </div>
</body>
</html>
    `;
  }

  getScoreClass(score) {
    return score >= 90 ? 'excellent' : score >= 75 ? 'good' : 'poor';
  }

  getWeekPeriod() {
    const now = new Date();
    const weekStart = new Date(now.setDate(now.getDate() - now.getDay()));
    const weekEnd = new Date(now.setDate(now.getDate() - now.getDay() + 6));
    
    return `${weekStart.toLocaleDateString('fr-FR')} - ${weekEnd.toLocaleDateString('fr-FR')}`;
  }

  getDateString() {
    return new Date().toISOString().split('T')[0];
  }
}

// Exécuter si appelé directement
if (require.main === module) {
  const generator = new UXReportGenerator();
  generator.generateWeeklyReport()
    .then(report => {
      console.log(chalk.green('✅ Weekly UX report generated successfully'));
      console.log(chalk.cyan(`Global SEPTEO Score: ${report.septeoCompliance.overall}%`));
      
      if (report.recommendations.length > 0) {
        console.log(chalk.yellow(`⚠️  ${report.recommendations.length} recommendations to address`));
      } else {
        console.log(chalk.green('🎉 No critical recommendations - excellent UX compliance!'));
      }
    })
    .catch(error => {
      console.error(chalk.red('❌ Error generating UX report:'), error);
      process.exit(1);
    });
}

module.exports = UXReportGenerator;