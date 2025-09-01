module.exports = {
  // Configuration Axe-core pour Fire Salamander
  rules: {
    // Règles critiques obligatoires
    'color-contrast': { enabled: true },
    'keyboard-navigation': { enabled: true },
    'focus-order-semantics': { enabled: true },
    'aria-labels': { enabled: true },
    
    // Règles SEPTEO spécifiques
    'custom-septeo-contrast': {
      enabled: true,
      options: {
        primaryColor: '#ff6136',
        minimumRatio: 4.5 // WCAG AA
      }
    }
  },
  
  tags: [
    'wcag2a',
    'wcag2aa',
    'wcag21aa',
    'best-practice'
  ],
  
  // URLs à tester
  urls: [
    'http://localhost:8080',
    'http://localhost:8080#analyzer',
    'http://localhost:8080#history',
    'http://localhost:8080#reports',
    'http://localhost:8080#monitoring'
  ],
  
  // Seuils minimum
  thresholds: {
    violations: 0,      // Zéro violation autorisée
    incomplete: 5,      // Maximum 5 tests incomplets
    passes: 50          // Minimum 50 tests passés
  },
  
  // Options spécifiques SEPTEO
  septeoOptions: {
    brandColors: {
      primary: '#ff6136',
      primaryDark: '#e55a2e',
      secondary: '#2c3e50',
      success: '#27ae60',
      warning: '#f39c12',
      danger: '#e74c3c'
    },
    contrastRequirements: {
      normal: 4.5,      // Texte normal
      large: 3.0,       // Texte large (18pt+ ou 14pt+ bold)
      decorative: 1.0   // Éléments décoratifs
    }
  },
  
  // Éléments à ignorer (false positives connus)
  exclude: [
    // Exclure les éléments Chart.js qui ont leurs propres règles d'accessibilité
    '#scoresChart',
    '#categoriesChart',
    '#metricsChart'
  ],
  
  // Configuration du reporter
  reporter: {
    format: 'html',
    outputDir: 'reports/accessibility/',
    fileName: 'axe-report.html'
  }
};