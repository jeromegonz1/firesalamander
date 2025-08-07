#!/usr/bin/env node

const lighthouse = require('lighthouse');
const chromeLauncher = require('chrome-launcher');
const fs = require('fs-extra');
const path = require('path');
const chalk = require('chalk');

// Configuration Lighthouse pour SEPTEO
const lighthouseConfig = {
  extends: 'lighthouse:default',
  settings: {
    onlyAudits: [
      'first-contentful-paint',
      'largest-contentful-paint',
      'first-meaningful-paint',
      'speed-index',
      'total-blocking-time',
      'cumulative-layout-shift',
      'interactive',
      'accessibility',
      'color-contrast',
      'tap-targets',
      'link-name',
      'button-name',
      'aria-*'
    ],
    output: ['html', 'json'],
    throttling: {
      rttMs: 150,
      throughputKbps: 1.6 * 1024,
      cpuSlowdownMultiplier: 4
    }
  }
};

async function runLighthouseAudit() {
  console.log(chalk.blue('⚡ Démarrage audit Lighthouse Fire Salamander...'));
  
  const chrome = await chromeLauncher.launch({
    chromeFlags: ['--headless', '--no-sandbox', '--disable-dev-shm-usage']
  });
  
  try {
    const runnerResult = await lighthouse(
      'http://localhost:8080',
      {
        port: chrome.port,
        output: ['html', 'json'],
        logLevel: 'info',
        throttling: 'simulated3G'
      }
    );

    // Sauvegarder les rapports
    const reportsDir = path.join(__dirname, '../reports/lighthouse');
    await fs.ensureDir(reportsDir);
    
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const htmlPath = path.join(reportsDir, `lighthouse-report-${timestamp}.html`);
    const jsonPath = path.join(reportsDir, `lighthouse-report-${timestamp}.json`);
    
    await fs.writeFile(htmlPath, runnerResult.report[0]);
    await fs.writeFile(jsonPath, runnerResult.report[1]);
    
    // Analyser les résultats selon standards SEPTEO
    const lhr = runnerResult.lhr;
    const scores = {
      performance: Math.round(lhr.categories.performance.score * 100),
      accessibility: Math.round(lhr.categories.accessibility.score * 100),
      bestPractices: Math.round(lhr.categories['best-practices'].score * 100),
      seo: Math.round(lhr.categories.seo.score * 100)
    };

    // Métriques Core Web Vitals
    const vitals = {
      lcp: lhr.audits['largest-contentful-paint'].numericValue,
      fid: lhr.audits['total-blocking-time'].numericValue, // Approximation
      cls: lhr.audits['cumulative-layout-shift'].numericValue,
      fcp: lhr.audits['first-contentful-paint'].numericValue,
      ttfb: lhr.audits['server-response-time']?.numericValue || 0
    };

    // Validation standards SEPTEO
    const septeoCompliance = {
      performance: scores.performance >= 90,
      accessibility: scores.accessibility >= 95,
      lcp: vitals.lcp <= 2500,
      cls: vitals.cls <= 0.1,
      fcp: vitals.fcp <= 1800
    };

    const overallCompliance = Object.values(septeoCompliance).filter(Boolean).length;
    const totalChecks = Object.keys(septeoCompliance).length;
    const compliancePercentage = Math.round((overallCompliance / totalChecks) * 100);

    console.log(chalk.green('\n📊 Résultats Lighthouse SEPTEO:'));
    console.log(chalk.cyan(`Performance: ${scores.performance}% ${septeoCompliance.performance ? '✅' : '❌'}`));
    console.log(chalk.cyan(`Accessibilité: ${scores.accessibility}% ${septeoCompliance.accessibility ? '✅' : '❌'}`));
    console.log(chalk.cyan(`Bonnes pratiques: ${scores.bestPractices}%`));
    console.log(chalk.cyan(`SEO: ${scores.seo}%`));
    
    console.log(chalk.green('\n🎯 Core Web Vitals:'));
    console.log(chalk.cyan(`LCP: ${Math.round(vitals.lcp)}ms ${septeoCompliance.lcp ? '✅' : '❌'}`));
    console.log(chalk.cyan(`FCP: ${Math.round(vitals.fcp)}ms ${septeoCompliance.fcp ? '✅' : '❌'}`));
    console.log(chalk.cyan(`CLS: ${vitals.cls.toFixed(3)} ${septeoCompliance.cls ? '✅' : '❌'}`));
    console.log(chalk.cyan(`TBT: ${Math.round(vitals.fid)}ms`));
    
    console.log(chalk.green(`\n🏆 Conformité SEPTEO: ${compliancePercentage}%`));
    console.log(chalk.blue(`📄 Rapports générés:`));
    console.log(chalk.blue(`  HTML: ${htmlPath}`));
    console.log(chalk.blue(`  JSON: ${jsonPath}`));

    // Générer un rapport consolidé
    const consolidatedReport = {
      timestamp: new Date().toISOString(),
      url: 'http://localhost:8080',
      scores,
      vitals,
      septeoCompliance: {
        individual: septeoCompliance,
        overall: compliancePercentage
      },
      reports: {
        html: htmlPath,
        json: jsonPath
      }
    };

    await fs.writeJSON(
      path.join(reportsDir, 'lighthouse-summary.json'),
      consolidatedReport,
      { spaces: 2 }
    );

    // Retourner le statut de succès selon SEPTEO
    const success = compliancePercentage >= 85;
    
    if (success) {
      console.log(chalk.green('✅ Audit Lighthouse réussi - Standards SEPTEO respectés'));
      process.exit(0);
    } else {
      console.log(chalk.red('❌ Audit Lighthouse échoué - Standards SEPTEO non respectés'));
      process.exit(1);
    }

  } catch (error) {
    console.error(chalk.red('❌ Erreur pendant l\'audit Lighthouse:'), error.message);
    process.exit(1);
  } finally {
    await chrome.kill();
  }
}

// Lancer l'audit si exécuté directement
if (require.main === module) {
  runLighthouseAudit().catch(console.error);
}

module.exports = { runLighthouseAudit };