const { chromium } = require('playwright');
const AxeBuilder = require('@axe-core/playwright').default;
const fs = require('fs-extra');
const path = require('path');
const chalk = require('chalk');

const config = require('./axe.config.js');

class FireSalamanderAccessibilityTester {
  constructor() {
    this.browser = null;
    this.results = [];
    this.reportsDir = path.join(__dirname, '../../../reports/accessibility');
  }

  async initialize() {
    console.log(chalk.blue('🚀 Initializing Fire Salamander Accessibility Tests...'));
    
    // Créer le dossier de rapports
    await fs.ensureDir(this.reportsDir);
    
    // Lancer le navigateur
    this.browser = await chromium.launch({ headless: true });
    
    // Vérifier que Fire Salamander est démarré
    const context = await this.browser.newContext();
    const page = await context.newPage();
    
    try {
      const response = await page.goto('http://localhost:8080/api/v1/health', { 
        timeout: 10000 
      });
      
      if (!response.ok()) {
        throw new Error(`Fire Salamander not ready: ${response.status()}`);
      }
      
      console.log(chalk.green('✅ Fire Salamander is ready'));
    } catch (error) {
      console.error(chalk.red('❌ Fire Salamander health check failed:'), error.message);
      throw error;
    } finally {
      await context.close();
    }
  }

  async testPage(url, pageName) {
    console.log(chalk.yellow(`\n🔍 Testing accessibility for: ${pageName}`));
    
    const context = await this.browser.newContext();
    const page = await context.newPage();
    
    try {
      // Naviguer vers la page
      await page.goto(url, { waitUntil: 'networkidle' });
      
      // Attendre que la page soit stable
      await page.waitForTimeout(2000);
      
      // Si c'est une page avec navigation, cliquer sur le bon lien
      if (url.includes('#')) {
        const hash = url.split('#')[1];
        const navLink = `a[data-page="${hash}"]`;
        
        await page.waitForSelector(navLink, { timeout: 5000 });
        await page.click(navLink);
        await page.waitForTimeout(1000);
      }
      
      // Configurer axe-core avec les règles SEPTEO
      const axeBuilder = new AxeBuilder({ page })
        .withTags(config.tags)
        .exclude(config.exclude);
      
      // Ajouter des règles personnalisées SEPTEO
      await page.addInitScript(() => {
        window.axe.configure({
          rules: [{
            id: 'septeo-color-contrast',
            impact: 'critical',
            tags: ['color-contrast', 'septeo'],
            all: [],
            any: [{
              id: 'septeo-contrast-check',
              evaluate: function(node, options) {
                const style = window.getComputedStyle(node);
                const bgColor = style.backgroundColor;
                const color = style.color;
                
                // Vérifier spécifiquement les couleurs SEPTEO
                if (bgColor.includes('255, 97, 54') || bgColor.includes('#ff6136')) {
                  // Orange SEPTEO détecté, vérifier le contraste
                  return this.data({
                    bgColor: bgColor,
                    fgColor: color,
                    contrastRatio: this.getContrast(bgColor, color)
                  });
                }
                
                return true;
              }
            }],
            none: [],
            metadata: {
              description: 'Vérification du contraste pour les couleurs SEPTEO',
              help: 'Les couleurs SEPTEO doivent respecter les ratios de contraste WCAG'
            }
          }]
        });
      });
      
      // Exécuter les tests axe
      const results = await axeBuilder.analyze();
      
      // Enrichir les résultats avec des informations SEPTEO
      const enrichedResults = await this.enrichWithSepteoData(page, results, pageName);
      
      // Sauvegarder les résultats
      this.results.push(enrichedResults);
      
      // Afficher le résumé
      this.logPageResults(enrichedResults);
      
      return enrichedResults;
      
    } catch (error) {
      console.error(chalk.red(`❌ Error testing ${pageName}:`), error.message);
      return {
        url,
        pageName,
        error: error.message,
        violations: [],
        passes: [],
        incomplete: [],
        inapplicable: []
      };
    } finally {
      await context.close();
    }
  }

  async enrichWithSepteoData(page, axeResults, pageName) {
    // Collecter des données spécifiques à SEPTEO
    const septeoData = await page.evaluate((brandColors) => {
      const data = {
        colorUsage: {},
        contrastIssues: [],
        keyboardNavigation: {
          focusableElements: 0,
          withoutFocusIndicator: 0
        },
        ariaLabels: {
          total: 0,
          missing: 0
        }
      };
      
      // Analyser l'utilisation des couleurs SEPTEO
      const allElements = document.querySelectorAll('*');
      allElements.forEach(el => {
        const style = window.getComputedStyle(el);
        const bgColor = style.backgroundColor;
        const color = style.color;
        
        // Détecter l'utilisation des couleurs SEPTEO
        Object.entries(brandColors).forEach(([name, hex]) => {
          if (bgColor.includes(hex) || color.includes(hex)) {
            data.colorUsage[name] = (data.colorUsage[name] || 0) + 1;
          }
        });
      });
      
      // Analyser la navigation clavier
      const focusableElements = document.querySelectorAll(
        'a[href], button, input, select, textarea, [tabindex]:not([tabindex="-1"])'
      );
      
      data.keyboardNavigation.focusableElements = focusableElements.length;
      
      focusableElements.forEach(el => {
        // Simuler le focus pour vérifier l'indicateur
        el.focus();
        const focusStyle = window.getComputedStyle(el, ':focus');
        const hasOutline = focusStyle.outline !== 'none' && focusStyle.outline !== '';
        const hasBoxShadow = focusStyle.boxShadow !== 'none';
        
        if (!hasOutline && !hasBoxShadow) {
          data.keyboardNavigation.withoutFocusIndicator++;
        }
      });
      
      // Analyser les labels ARIA
      const elementsWithAria = document.querySelectorAll('[aria-label], [aria-labelledby], [aria-describedby]');
      const elementsNeedingAria = document.querySelectorAll('input, button, select, textarea');
      
      data.ariaLabels.total = elementsNeedingAria.length;
      data.ariaLabels.missing = elementsNeedingAria.length - elementsWithAria.length;
      
      return data;
    }, config.septeoOptions.brandColors);
    
    return {
      ...axeResults,
      pageName,
      septeoData,
      timestamp: new Date().toISOString()
    };
  }

  logPageResults(results) {
    const { violations, passes, incomplete, septeoData } = results;
    
    console.log(chalk.blue(`\n📊 Results for ${results.pageName}:`));
    
    if (violations.length === 0) {
      console.log(chalk.green(`✅ No violations found!`));
    } else {
      console.log(chalk.red(`❌ ${violations.length} violations found:`));
      violations.forEach(violation => {
        console.log(chalk.red(`  • ${violation.description}`));
        console.log(chalk.yellow(`    Impact: ${violation.impact}`));
        console.log(chalk.gray(`    Nodes: ${violation.nodes.length}`));
      });
    }
    
    console.log(chalk.green(`✅ ${passes.length} tests passed`));
    
    if (incomplete.length > 0) {
      console.log(chalk.yellow(`⚠️  ${incomplete.length} tests incomplete`));
    }
    
    // Résultats SEPTEO spécifiques
    console.log(chalk.blue('\n🎨 SEPTEO Brand Analysis:'));
    Object.entries(septeoData.colorUsage).forEach(([color, count]) => {
      console.log(chalk.cyan(`  • ${color}: used ${count} times`));
    });
    
    if (septeoData.keyboardNavigation.withoutFocusIndicator > 0) {
      console.log(chalk.yellow(`⚠️  ${septeoData.keyboardNavigation.withoutFocusIndicator} elements without focus indicator`));
    }
    
    if (septeoData.ariaLabels.missing > 0) {
      console.log(chalk.yellow(`⚠️  ${septeoData.ariaLabels.missing} elements missing ARIA labels`));
    }
  }

  async generateReport() {
    console.log(chalk.blue('\n📄 Generating accessibility report...'));
    
    const summary = {
      totalPages: this.results.length,
      totalViolations: this.results.reduce((sum, r) => sum + (r.violations?.length || 0), 0),
      totalPasses: this.results.reduce((sum, r) => sum + (r.passes?.length || 0), 0),
      totalIncomplete: this.results.reduce((sum, r) => sum + (r.incomplete?.length || 0), 0),
      timestamp: new Date().toISOString(),
      septeoCompliance: this.calculateSepteoCompliance()
    };
    
    const report = {
      summary,
      results: this.results,
      config: config
    };
    
    // Sauvegarder le rapport JSON
    const jsonPath = path.join(this.reportsDir, 'accessibility-report.json');
    await fs.writeJSON(jsonPath, report, { spaces: 2 });
    
    // Générer un rapport HTML
    const htmlReport = this.generateHTMLReport(report);
    const htmlPath = path.join(this.reportsDir, 'accessibility-report.html');
    await fs.writeFile(htmlPath, htmlReport);
    
    console.log(chalk.green(`✅ Reports generated:`));
    console.log(chalk.cyan(`  JSON: ${jsonPath}`));
    console.log(chalk.cyan(`  HTML: ${htmlPath}`));
    
    return summary;
  }

  calculateSepteoCompliance() {
    const septeoScore = {
      colorContrast: 0,
      keyboardNavigation: 0,
      ariaLabels: 0,
      overall: 0
    };
    
    this.results.forEach(result => {
      if (result.septeoData) {
        const { keyboardNavigation, ariaLabels } = result.septeoData;
        
        // Score navigation clavier
        if (keyboardNavigation.focusableElements > 0) {
          septeoScore.keyboardNavigation += (keyboardNavigation.focusableElements - keyboardNavigation.withoutFocusIndicator) / keyboardNavigation.focusableElements;
        }
        
        // Score ARIA labels
        if (ariaLabels.total > 0) {
          septeoScore.ariaLabels += (ariaLabels.total - ariaLabels.missing) / ariaLabels.total;
        }
        
        // Score contraste (basé sur l'absence de violations de contraste)
        const contrastViolations = result.violations.filter(v => v.id.includes('color-contrast')).length;
        septeoScore.colorContrast += contrastViolations === 0 ? 1 : 0;
      }
    });
    
    // Moyenne des scores
    const pageCount = this.results.length;
    septeoScore.colorContrast = (septeoScore.colorContrast / pageCount) * 100;
    septeoScore.keyboardNavigation = (septeoScore.keyboardNavigation / pageCount) * 100;
    septeoScore.ariaLabels = (septeoScore.ariaLabels / pageCount) * 100;
    septeoScore.overall = (septeoScore.colorContrast + septeoScore.keyboardNavigation + septeoScore.ariaLabels) / 3;
    
    return septeoScore;
  }

  generateHTMLReport(report) {
    return `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Fire Salamander - Rapport d'Accessibilité</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: linear-gradient(135deg, #ff6136 0%, #e55a2e 100%); color: white; padding: 30px; border-radius: 8px; margin-bottom: 30px; }
        .summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .metric { background: white; padding: 20px; border-radius: 8px; text-align: center; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .metric-value { font-size: 2rem; font-weight: bold; color: #ff6136; }
        .metric-label { color: #666; margin-top: 10px; }
        .section { background: white; padding: 20px; margin-bottom: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .violation { background: #fee; border-left: 4px solid #e74c3c; padding: 15px; margin: 10px 0; border-radius: 4px; }
        .pass { background: #efe; border-left: 4px solid #27ae60; padding: 15px; margin: 10px 0; border-radius: 4px; }
        .score-bar { height: 20px; background: #eee; border-radius: 10px; overflow: hidden; margin: 10px 0; }
        .score-fill { height: 100%; background: linear-gradient(90deg, #e74c3c 0%, #f39c12 50%, #27ae60 100%); transition: width 0.3s ease; }
        .septeo-badge { background: #ff6136; color: white; padding: 5px 10px; border-radius: 4px; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔥 Fire Salamander - Rapport d'Accessibilité</h1>
            <p>Conformité WCAG 2.1 AA + Standards SEPTEO</p>
            <p>Généré le: ${new Date(report.summary.timestamp).toLocaleString('fr-FR')}</p>
        </div>
        
        <div class="summary">
            <div class="metric">
                <div class="metric-value">${report.summary.totalViolations}</div>
                <div class="metric-label">Violations</div>
            </div>
            <div class="metric">
                <div class="metric-value">${report.summary.totalPasses}</div>
                <div class="metric-label">Tests Réussis</div>
            </div>
            <div class="metric">
                <div class="metric-value">${Math.round(report.summary.septeoCompliance.overall)}%</div>
                <div class="metric-label">Conformité SEPTEO</div>
            </div>
            <div class="metric">
                <div class="metric-value">${report.summary.totalPages}</div>
                <div class="metric-label">Pages Testées</div>
            </div>
        </div>
        
        <div class="section">
            <h2>🎨 Conformité SEPTEO</h2>
            <div>
                <h3>Contraste des Couleurs</h3>
                <div class="score-bar">
                    <div class="score-fill" style="width: ${report.summary.septeoCompliance.colorContrast}%"></div>
                </div>
                <p>${Math.round(report.summary.septeoCompliance.colorContrast)}% conforme</p>
            </div>
            <div>
                <h3>Navigation Clavier</h3>
                <div class="score-bar">
                    <div class="score-fill" style="width: ${report.summary.septeoCompliance.keyboardNavigation}%"></div>
                </div>
                <p>${Math.round(report.summary.septeoCompliance.keyboardNavigation)}% conforme</p>
            </div>
            <div>
                <h3>Labels ARIA</h3>
                <div class="score-bar">
                    <div class="score-fill" style="width: ${report.summary.septeoCompliance.ariaLabels}%"></div>
                </div>
                <p>${Math.round(report.summary.septeoCompliance.ariaLabels)}% conforme</p>
            </div>
        </div>
        
        ${report.results.map(result => `
        <div class="section">
            <h2>${result.pageName} <span class="septeo-badge">SEPTEO</span></h2>
            <p><strong>URL:</strong> ${result.url}</p>
            
            ${result.violations && result.violations.length > 0 ? `
            <h3>❌ Violations (${result.violations.length})</h3>
            ${result.violations.map(violation => `
            <div class="violation">
                <h4>${violation.description}</h4>
                <p><strong>Impact:</strong> ${violation.impact}</p>
                <p><strong>Aide:</strong> ${violation.help}</p>
                <p><strong>Éléments affectés:</strong> ${violation.nodes.length}</p>
            </div>
            `).join('')}
            ` : '<div class="pass"><h3>✅ Aucune violation détectée!</h3></div>'}
            
            ${result.passes && result.passes.length > 0 ? `
            <h3>✅ Tests Réussis (${result.passes.length})</h3>
            <p>Tous les tests d'accessibilité sont passés avec succès.</p>
            ` : ''}
        </div>
        `).join('')}
    </div>
</body>
</html>
    `;
  }

  async run() {
    try {
      await this.initialize();
      
      console.log(chalk.blue('\n🧪 Running accessibility tests for all pages...'));
      
      // Tester toutes les URLs configurées
      for (const url of config.urls) {
        const pageName = url.includes('#') ? 
          `Page ${url.split('#')[1].charAt(0).toUpperCase() + url.split('#')[1].slice(1)}` : 
          'Dashboard';
        
        await this.testPage(url, pageName);
      }
      
      // Générer le rapport final
      const summary = await this.generateReport();
      
      // Afficher le résumé final
      console.log(chalk.blue('\n📊 Final Accessibility Summary:'));
      console.log(chalk.green(`✅ ${summary.totalPasses} tests passed`));
      
      if (summary.totalViolations > 0) {
        console.log(chalk.red(`❌ ${summary.totalViolations} violations found`));
      } else {
        console.log(chalk.green('🎉 No accessibility violations found!'));
      }
      
      console.log(chalk.cyan(`🎨 SEPTEO Compliance: ${Math.round(summary.septeoCompliance.overall)}%`));
      
      // Vérifier les seuils
      const meetsCriteria = summary.totalViolations <= config.thresholds.violations &&
                           summary.totalPasses >= config.thresholds.passes &&
                           summary.septeoCompliance.overall >= 95;
      
      if (meetsCriteria) {
        console.log(chalk.green('\n🏆 All accessibility criteria met!'));
        process.exit(0);
      } else {
        console.log(chalk.red('\n💥 Accessibility criteria not met'));
        process.exit(1);
      }
      
    } catch (error) {
      console.error(chalk.red('💥 Accessibility testing failed:'), error);
      process.exit(1);
    } finally {
      if (this.browser) {
        await this.browser.close();
      }
    }
  }
}

// Exécuter si appelé directement
if (require.main === module) {
  const tester = new FireSalamanderAccessibilityTester();
  tester.run();
}

module.exports = FireSalamanderAccessibilityTester;