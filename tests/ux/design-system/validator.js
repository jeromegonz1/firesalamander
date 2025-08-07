const fs = require('fs-extra');
const path = require('path');
const postcss = require('postcss');
const { parse } = require('css-tree');
const chalk = require('chalk');
const { chromium } = require('playwright');

class SepteoDesignSystemValidator {
  constructor() {
    this.errors = [];
    this.warnings = [];
    this.passes = [];
    
    // Standards SEPTEO
    this.septeoStandards = {
      colors: {
        allowed: {
          primary: '#ff6136',
          primaryDark: '#e55a2e',
          secondary: '#2c3e50',
          success: '#27ae60',
          warning: '#f39c12',
          danger: '#e74c3c',
          info: '#3498db',
          white: '#ffffff',
          black: '#000000',
          gray100: '#f8f9fa',
          gray200: '#e9ecef',
          gray300: '#dee2e6',
          gray400: '#ced4da',
          gray500: '#adb5bd',
          gray600: '#6c757d',
          gray700: '#495057',
          gray800: '#343a40',
          gray900: '#212529'
        },
        forbidden: [
          '#ff0000', // Rouge pur interdit
          '#00ff00', // Vert pur interdit
          '#0000ff', // Bleu pur interdit
          '#ffff00', // Jaune pur interdit
          '#ff00ff', // Magenta pur interdit
          '#00ffff'  // Cyan pur interdit
        ]
      },
      
      spacing: {
        // Grille 8px obligatoire
        baseUnit: 8,
        allowedValues: [0, 4, 8, 12, 16, 20, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128],
        properties: ['margin', 'padding', 'gap', 'top', 'right', 'bottom', 'left']
      },
      
      typography: {
        allowedFamilies: [
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'Roboto',
          'sans-serif',
          'system-ui'
        ],
        forbiddenFamilies: [
          'Times New Roman',
          'serif',
          'Comic Sans MS',
          'Courier New',
          'monospace'
        ],
        allowedSizes: [
          '12px', '14px', '16px', '18px', '20px', '24px', 
          '28px', '32px', '36px', '40px', '48px', '56px', '64px'
        ]
      },
      
      components: {
        requiredClasses: [
          'btn', 'btn-primary', 'btn-secondary',
          'card', 'form-input', 'form-select',
          'navbar', 'main-content', 'section'
        ],
        forbiddenInlineStyles: true,
        requireBEM: false // Pour l'instant
      }
    };
    
    this.reportsDir = path.join(__dirname, '../../../reports/design-system');
  }

  async initialize() {
    console.log(chalk.blue('🎨 Initializing SEPTEO Design System Validation...'));
    await fs.ensureDir(this.reportsDir);
  }

  async validateCSS() {
    console.log(chalk.yellow('\n🎨 Validating CSS against SEPTEO standards...'));
    
    const cssPath = path.join(__dirname, '../../../web/static/styles.css');
    
    if (!await fs.pathExists(cssPath)) {
      this.errors.push({
        type: 'file-missing',
        message: 'CSS file not found',
        file: cssPath
      });
      return;
    }
    
    const cssContent = await fs.readFile(cssPath, 'utf8');
    
    try {
      // Parser le CSS
      const ast = parse(cssContent);
      
      // Valider les couleurs
      await this.validateColors(ast);
      
      // Valider l'espacement
      await this.validateSpacing(ast);
      
      // Valider la typographie
      await this.validateTypography(ast);
      
      // Valider la structure des composants
      await this.validateComponents(ast);
      
    } catch (error) {
      this.errors.push({
        type: 'css-parse-error',
        message: `Error parsing CSS: ${error.message}`,
        file: cssPath
      });
    }
  }

  validateColors(ast) {
    const colorRegex = /#[0-9a-fA-F]{6}|#[0-9a-fA-F]{3}|rgb\([^)]+\)|rgba\([^)]+\)|hsl\([^)]+\)|hsla\([^)]+\)/g;
    
    ast.children.forEach(node => {
      if (node.type === 'Rule') {
        node.block.children.forEach(declaration => {
          if (declaration.type === 'Declaration') {
            const value = declaration.value.children.first();
            if (value && value.type === 'Raw') {
              const colors = value.value.match(colorRegex);
              
              if (colors) {
                colors.forEach(color => {
                  const normalizedColor = this.normalizeColor(color);
                  
                  // Vérifier si la couleur est autorisée
                  const isAllowed = Object.values(this.septeoStandards.colors.allowed)
                    .some(allowedColor => this.colorsMatch(normalizedColor, allowedColor));
                  
                  // Vérifier si la couleur est interdite
                  const isForbidden = this.septeoStandards.colors.forbidden
                    .some(forbiddenColor => this.colorsMatch(normalizedColor, forbiddenColor));
                  
                  if (isForbidden) {
                    this.errors.push({
                      type: 'forbidden-color',
                      message: `Forbidden color used: ${color}`,
                      property: declaration.property,
                      selector: this.getSelector(node),
                      color: color
                    });
                  } else if (!isAllowed && !this.isGrayScale(normalizedColor)) {
                    this.warnings.push({
                      type: 'non-standard-color',
                      message: `Non-standard color used: ${color}`,
                      property: declaration.property,
                      selector: this.getSelector(node),
                      color: color,
                      suggestion: this.suggestSepteoColor(normalizedColor)
                    });
                  } else {
                    this.passes.push({
                      type: 'color-approved',
                      message: `SEPTEO color correctly used: ${color}`,
                      property: declaration.property,
                      selector: this.getSelector(node)
                    });
                  }
                });
              }
            }
          }
        });
      }
    });
  }

  validateSpacing(ast) {
    const spacingProps = this.septeoStandards.spacing.properties;
    
    ast.children.forEach(node => {
      if (node.type === 'Rule') {
        node.block.children.forEach(declaration => {
          if (declaration.type === 'Declaration' && spacingProps.includes(declaration.property)) {
            const value = this.getDeclarationValue(declaration);
            const pxValues = this.extractPixelValues(value);
            
            pxValues.forEach(pxValue => {
              const numValue = parseInt(pxValue);
              
              if (!this.septeoStandards.spacing.allowedValues.includes(numValue)) {
                // Vérifier si c'est un multiple de 8
                if (numValue % this.septeoStandards.spacing.baseUnit !== 0) {
                  this.errors.push({
                    type: 'invalid-spacing',
                    message: `Spacing not on 8px grid: ${pxValue}px`,
                    property: declaration.property,
                    selector: this.getSelector(node),
                    value: pxValue,
                    suggestion: this.suggestValidSpacing(numValue)
                  });
                } else {
                  this.warnings.push({
                    type: 'non-standard-spacing',
                    message: `Spacing on grid but not in standard values: ${pxValue}px`,
                    property: declaration.property,
                    selector: this.getSelector(node),
                    value: pxValue
                  });
                }
              } else {
                this.passes.push({
                  type: 'spacing-approved',
                  message: `Standard spacing used: ${pxValue}px`,
                  property: declaration.property,
                  selector: this.getSelector(node)
                });
              }
            });
          }
        });
      }
    });
  }

  validateTypography(ast) {
    ast.children.forEach(node => {
      if (node.type === 'Rule') {
        node.block.children.forEach(declaration => {
          if (declaration.type === 'Declaration') {
            // Valider font-family
            if (declaration.property === 'font-family') {
              const fontFamily = this.getDeclarationValue(declaration);
              const fonts = fontFamily.split(',').map(f => f.trim().replace(/['"]/g, ''));
              
              const hasForbiddenFont = fonts.some(font => 
                this.septeoStandards.typography.forbiddenFamilies.includes(font)
              );
              
              const hasAllowedFont = fonts.some(font => 
                this.septeoStandards.typography.allowedFamilies.includes(font)
              );
              
              if (hasForbiddenFont) {
                this.errors.push({
                  type: 'forbidden-font',
                  message: `Forbidden font family used: ${fontFamily}`,
                  selector: this.getSelector(node),
                  value: fontFamily
                });
              } else if (!hasAllowedFont) {
                this.warnings.push({
                  type: 'non-standard-font',
                  message: `Non-standard font family: ${fontFamily}`,
                  selector: this.getSelector(node),
                  value: fontFamily
                });
              } else {
                this.passes.push({
                  type: 'font-approved',
                  message: `Standard font family used: ${fontFamily}`,
                  selector: this.getSelector(node)
                });
              }
            }
            
            // Valider font-size
            if (declaration.property === 'font-size') {
              const fontSize = this.getDeclarationValue(declaration);
              
              if (!this.septeoStandards.typography.allowedSizes.includes(fontSize)) {
                this.warnings.push({
                  type: 'non-standard-font-size',
                  message: `Non-standard font size: ${fontSize}`,
                  selector: this.getSelector(node),
                  value: fontSize,
                  suggestion: this.suggestFontSize(fontSize)
                });
              } else {
                this.passes.push({
                  type: 'font-size-approved',
                  message: `Standard font size used: ${fontSize}`,
                  selector: this.getSelector(node)
                });
              }
            }
          }
        });
      }
    });
  }

  validateComponents(ast) {
    // Vérifier la présence des classes de composants obligatoires
    const foundClasses = new Set();
    
    ast.children.forEach(node => {
      if (node.type === 'Rule') {
        const selector = this.getSelector(node);
        
        // Extraire les classes du sélecteur
        const classMatches = selector.match(/\.[a-zA-Z_-]+/g);
        if (classMatches) {
          classMatches.forEach(match => {
            const className = match.substring(1); // Enlever le point
            foundClasses.add(className);
          });
        }
      }
    });
    
    // Vérifier les classes obligatoires
    this.septeoStandards.components.requiredClasses.forEach(requiredClass => {
      if (foundClasses.has(requiredClass)) {
        this.passes.push({
          type: 'component-class-found',
          message: `Required component class found: .${requiredClass}`,
          className: requiredClass
        });
      } else {
        this.errors.push({
          type: 'missing-component-class',
          message: `Missing required component class: .${requiredClass}`,
          className: requiredClass
        });
      }
    });
  }

  async validateLiveInterface() {
    console.log(chalk.yellow('\n🌐 Validating live interface...'));
    
    const browser = await chromium.launch({ headless: true });
    const context = await browser.newContext();
    const page = await context.newPage();
    
    try {
      await page.goto('http://localhost:8080', { waitUntil: 'networkidle' });
      
      // Vérifier les styles inline interdits
      const inlineStyles = await page.evaluate(() => {
        const elements = document.querySelectorAll('[style]');
        return Array.from(elements).map(el => ({
          tag: el.tagName,
          className: el.className,
          style: el.getAttribute('style')
        }));
      });
      
      if (this.septeoStandards.components.forbiddenInlineStyles && inlineStyles.length > 0) {
        inlineStyles.forEach(element => {
          this.warnings.push({
            type: 'inline-styles-detected',
            message: `Inline styles detected on ${element.tag}`,
            element: element,
            suggestion: 'Use CSS classes instead of inline styles'
          });
        });
      }
      
      // Vérifier la cohérence des couleurs affichées
      const computedColors = await page.evaluate(() => {
        const colorData = [];
        const elements = document.querySelectorAll('*');
        
        elements.forEach(el => {
          const style = window.getComputedStyle(el);
          const bgColor = style.backgroundColor;
          const color = style.color;
          
          if (bgColor !== 'rgba(0, 0, 0, 0)' && bgColor !== 'transparent') {
            colorData.push({
              element: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : ''),
              property: 'background-color',
              value: bgColor
            });
          }
          
          if (color !== 'rgb(0, 0, 0)') {
            colorData.push({
              element: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : ''),
              property: 'color',
              value: color
            });
          }
        });
        
        return colorData;
      });
      
      // Analyser les couleurs calculées
      computedColors.forEach(colorInfo => {
        const normalizedColor = this.normalizeColor(colorInfo.value);
        const isStandardColor = Object.values(this.septeoStandards.colors.allowed)
          .some(standardColor => this.colorsMatch(normalizedColor, standardColor));
        
        if (!isStandardColor && !this.isGrayScale(normalizedColor)) {
          this.warnings.push({
            type: 'non-standard-computed-color',
            message: `Non-standard computed color: ${colorInfo.value}`,
            element: colorInfo.element,
            property: colorInfo.property,
            value: colorInfo.value
          });
        }
      });
      
    } catch (error) {
      this.errors.push({
        type: 'live-validation-error',
        message: `Error validating live interface: ${error.message}`
      });
    } finally {
      await browser.close();
    }
  }

  // Fonctions utilitaires
  normalizeColor(color) {
    // Convertir différents formats de couleur vers hex
    if (color.startsWith('rgb')) {
      const matches = color.match(/\d+/g);
      if (matches && matches.length >= 3) {
        const r = parseInt(matches[0]);
        const g = parseInt(matches[1]);
        const b = parseInt(matches[2]);
        return '#' + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
      }
    }
    return color.toLowerCase();
  }

  colorsMatch(color1, color2) {
    return this.normalizeColor(color1) === this.normalizeColor(color2);
  }

  isGrayScale(color) {
    // Vérifier si c'est une couleur grise (R=G=B)
    if (color.startsWith('#')) {
      const r = parseInt(color.slice(1, 3), 16);
      const g = parseInt(color.slice(3, 5), 16);
      const b = parseInt(color.slice(5, 7), 16);
      return r === g && g === b;
    }
    return false;
  }

  suggestSepteoColor(color) {
    // Suggérer la couleur SEPTEO la plus proche
    const colors = this.septeoStandards.colors.allowed;
    
    // Pour simplifier, suggérer selon la teinte
    if (color.includes('ff') && color.includes('36')) return colors.primary;
    if (color.includes('2c') && color.includes('3e')) return colors.secondary;
    return colors.primary; // Par défaut
  }

  suggestValidSpacing(value) {
    const allowed = this.septeoStandards.spacing.allowedValues;
    const closest = allowed.reduce((prev, curr) => 
      Math.abs(curr - value) < Math.abs(prev - value) ? curr : prev
    );
    return `${closest}px`;
  }

  suggestFontSize(fontSize) {
    const allowed = this.septeoStandards.typography.allowedSizes;
    const numValue = parseInt(fontSize);
    const closest = allowed.reduce((prev, curr) => {
      const prevNum = parseInt(prev);
      const currNum = parseInt(curr);
      return Math.abs(currNum - numValue) < Math.abs(prevNum - numValue) ? curr : prev;
    });
    return closest;
  }

  getSelector(ruleNode) {
    return ruleNode.prelude.children.first().name || 'unknown';
  }

  getDeclarationValue(declarationNode) {
    return declarationNode.value.children.first().value || '';
  }

  extractPixelValues(value) {
    const matches = value.match(/(\d+)px/g);
    return matches ? matches.map(match => parseInt(match)) : [];
  }

  async generateReport() {
    console.log(chalk.blue('\n📄 Generating design system report...'));
    
    const summary = {
      timestamp: new Date().toISOString(),
      totalErrors: this.errors.length,
      totalWarnings: this.warnings.length,
      totalPasses: this.passes.length,
      septeoCompliance: this.calculateCompliance()
    };
    
    const report = {
      summary,
      errors: this.errors,
      warnings: this.warnings,
      passes: this.passes,
      standards: this.septeoStandards
    };
    
    // Sauvegarder JSON
    const jsonPath = path.join(this.reportsDir, 'design-system-report.json');
    await fs.writeJSON(jsonPath, report, { spaces: 2 });
    
    // Générer HTML
    const htmlReport = this.generateHTMLReport(report);
    const htmlPath = path.join(this.reportsDir, 'design-system-report.html');
    await fs.writeFile(htmlPath, htmlReport);
    
    console.log(chalk.green(`✅ Reports generated:`));
    console.log(chalk.cyan(`  JSON: ${jsonPath}`));
    console.log(chalk.cyan(`  HTML: ${htmlPath}`));
    
    return summary;
  }

  calculateCompliance() {
    const total = this.errors.length + this.warnings.length + this.passes.length;
    if (total === 0) return 100;
    
    const compliance = (this.passes.length / total) * 100;
    return Math.round(compliance);
  }

  generateHTMLReport(report) {
    return `
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Fire Salamander - Design System Validation</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: linear-gradient(135deg, #ff6136 0%, #e55a2e 100%); color: white; padding: 30px; border-radius: 8px; margin-bottom: 30px; }
        .summary { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .metric { background: white; padding: 20px; border-radius: 8px; text-align: center; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .metric-value { font-size: 2rem; font-weight: bold; color: #ff6136; }
        .section { background: white; padding: 20px; margin-bottom: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .error { background: #fee; border-left: 4px solid #e74c3c; padding: 15px; margin: 10px 0; border-radius: 4px; }
        .warning { background: #fff8e1; border-left: 4px solid #f39c12; padding: 15px; margin: 10px 0; border-radius: 4px; }
        .pass { background: #efe; border-left: 4px solid #27ae60; padding: 15px; margin: 10px 0; border-radius: 4px; }
        .color-swatch { display: inline-block; width: 20px; height: 20px; border-radius: 3px; margin-right: 10px; vertical-align: middle; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎨 Fire Salamander - Design System Validation</h1>
            <p>Conformité aux Standards SEPTEO</p>
            <p>Généré le: ${new Date(report.summary.timestamp).toLocaleString('fr-FR')}</p>
        </div>
        
        <div class="summary">
            <div class="metric">
                <div class="metric-value">${report.summary.septeoCompliance}%</div>
                <div class="metric-label">Conformité SEPTEO</div>
            </div>
            <div class="metric">
                <div class="metric-value">${report.summary.totalErrors}</div>
                <div class="metric-label">Erreurs</div>
            </div>
            <div class="metric">
                <div class="metric-value">${report.summary.totalWarnings}</div>
                <div class="metric-label">Avertissements</div>
            </div>
            <div class="metric">
                <div class="metric-value">${report.summary.totalPasses}</div>
                <div class="metric-label">Validations</div>
            </div>
        </div>
        
        ${report.errors.length > 0 ? `
        <div class="section">
            <h2>❌ Erreurs (${report.errors.length})</h2>
            ${report.errors.map(error => `
            <div class="error">
                <h4>${error.type}: ${error.message}</h4>
                ${error.selector ? `<p><strong>Sélecteur:</strong> ${error.selector}</p>` : ''}
                ${error.property ? `<p><strong>Propriété:</strong> ${error.property}</p>` : ''}
                ${error.value ? `<p><strong>Valeur:</strong> ${error.value}</p>` : ''}
                ${error.suggestion ? `<p><strong>Suggestion:</strong> ${error.suggestion}</p>` : ''}
            </div>
            `).join('')}
        </div>
        ` : ''}
        
        ${report.warnings.length > 0 ? `
        <div class="section">
            <h2>⚠️ Avertissements (${report.warnings.length})</h2>
            ${report.warnings.map(warning => `
            <div class="warning">
                <h4>${warning.type}: ${warning.message}</h4>
                ${warning.selector ? `<p><strong>Sélecteur:</strong> ${warning.selector}</p>` : ''}
                ${warning.suggestion ? `<p><strong>Suggestion:</strong> ${warning.suggestion}</p>` : ''}
            </div>
            `).join('')}
        </div>
        ` : ''}
        
        <div class="section">
            <h2>✅ Validations Réussies (${report.summary.totalPasses})</h2>
            <p>Éléments conformes aux standards SEPTEO.</p>
        </div>
        
        <div class="section">
            <h2>🎨 Palette SEPTEO Autorisée</h2>
            <div>
                ${Object.entries(report.standards.colors.allowed).map(([name, color]) => `
                <div style="margin: 10px 0;">
                    <span class="color-swatch" style="background-color: ${color}"></span>
                    <strong>${name}:</strong> ${color}
                </div>
                `).join('')}
            </div>
        </div>
    </div>
</body>
</html>
    `;
  }

  async run() {
    try {
      await this.initialize();
      
      console.log(chalk.blue('🎨 Running SEPTEO Design System Validation...'));
      
      // Valider le CSS
      await this.validateCSS();
      
      // Valider l'interface live
      await this.validateLiveInterface();
      
      // Générer le rapport
      const summary = await this.generateReport();
      
      // Afficher le résumé
      console.log(chalk.blue('\n📊 Design System Validation Summary:'));
      console.log(chalk.green(`✅ ${summary.totalPasses} validations passed`));
      
      if (summary.totalWarnings > 0) {
        console.log(chalk.yellow(`⚠️  ${summary.totalWarnings} warnings`));
      }
      
      if (summary.totalErrors > 0) {
        console.log(chalk.red(`❌ ${summary.totalErrors} errors`));
      }
      
      console.log(chalk.cyan(`🎨 SEPTEO Compliance: ${summary.septeoCompliance}%`));
      
      // Critères de réussite
      const meetsCriteria = summary.totalErrors === 0 && summary.septeoCompliance >= 90;
      
      if (meetsCriteria) {
        console.log(chalk.green('\n🏆 Design system validation passed!'));
        process.exit(0);
      } else {
        console.log(chalk.red('\n💥 Design system validation failed'));
        process.exit(1);
      }
      
    } catch (error) {
      console.error(chalk.red('💥 Design system validation error:'), error);
      process.exit(1);
    }
  }
}

// Exécuter si appelé directement
if (require.main === module) {
  const validator = new SepteoDesignSystemValidator();
  validator.run();
}

module.exports = SepteoDesignSystemValidator;