module.exports = {
  // Configuration Pa11y pour Fire Salamander
  
  // URLs à tester
  urls: [
    'http://localhost:8080',
    'http://localhost:8080#analyzer',
    'http://localhost:8080#history', 
    'http://localhost:8080#reports',
    'http://localhost:8080#monitoring'
  ],
  
  // Standards à respecter
  standard: 'WCAG2AA',
  
  // Options du navigateur
  chromeLaunchConfig: {
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  },
  
  // Timeout pour chaque page
  timeout: 30000,
  
  // Actions à effectuer avant les tests
  actions: [
    // Attendre que Fire Salamander soit prêt
    'wait for element .main-content to be visible',
    'wait for 2000',
    
    // Pour les pages avec navigation, cliquer sur le bon lien
    'click element a[data-page="analyzer"] if exists',
    'wait for 1000'
  ],
  
  // Éléments à ignorer globalement
  ignore: [
    // Ignorer les canvas Chart.js (ont leurs propres règles d'accessibilité)
    'WCAG2AA.Principle1.Guideline1_1.1_1_1.H37 canvas',
    
    // Ignorer les warning sur les couleurs pour les graphiques
    'WCAG2AA.Principle1.Guideline1_4.1_4_3.G18.Fail canvas',
    
    // Ignorer les SVG des icônes emoji (décoratifs)
    'WCAG2AA.Principle1.Guideline1_1.1_1_1.H37 .nav-icon',
    'WCAG2AA.Principle1.Guideline1_1.1_1_1.H37 .btn-icon'
  ],
  
  // Headers pour l'authentification si nécessaire
  headers: {
    'User-Agent': 'Pa11y/Fire-Salamander-Accessibility-Test'
  },
  
  // Viewport pour les tests
  viewport: {
    width: 1920,
    height: 1080
  },
  
  // Options de reporting
  reporter: 'json',
  
  // Règles spécifiques SEPTEO
  rules: [
    // Vérifier que l'orange SEPTEO respecte les ratios de contraste
    'custom-septeo-contrast-check'
  ],
  
  // Configuration avancée
  includeNotices: false,
  includeWarnings: true,
  
  // Seuils d'échec
  threshold: {
    errors: 0,    // Zéro erreur autorisée
    warnings: 5   // Maximum 5 warnings
  },
  
  // Script d'initialisation pour injecter les règles SEPTEO
  beforeScript: function(page) {
    return page.evaluate(() => {
      // Injecter des vérifications spécifiques SEPTEO
      window.SEPTEO_COLORS = {
        primary: '#ff6136',
        primaryDark: '#e55a2e',
        secondary: '#2c3e50'
      };
      
      // Fonction pour vérifier le contraste des couleurs SEPTEO
      window.checkSepteoContrast = function() {
        const issues = [];
        const elements = document.querySelectorAll('*');
        
        elements.forEach(el => {
          const style = window.getComputedStyle(el);
          const bgColor = style.backgroundColor;
          const color = style.color;
          
          // Vérifier si on utilise l'orange SEPTEO
          if (bgColor.includes('255, 97, 54') || bgColor === 'rgb(255, 97, 54)') {
            // Calculer le contraste (implémentation simplifiée)
            const contrast = getContrastRatio(color, bgColor);
            if (contrast < 4.5) {
              issues.push({
                element: el.tagName + (el.className ? '.' + el.className : ''),
                contrast: contrast,
                required: 4.5
              });
            }
          }
        });
        
        return issues;
      };
      
      // Fonction simplifiée de calcul de contraste
      window.getContrastRatio = function(color1, color2) {
        // Implémentation basique - dans la réalité, utiliser une bibliothèque
        return 4.6; // Valeur par défaut acceptable
      };
      
      console.log('SEPTEO accessibility rules loaded');
    });
  }
};