module.exports = async (page, scenario, vp) => {
  console.log('ENABLING DARK MODE > ' + scenario.label);
  
  // Injecter le mode sombre pour tester la compatibilité
  await page.evaluate(() => {
    // Créer un thème sombre compatible SEPTEO
    const darkTheme = `
      :root {
        --primary-color: #ff6136 !important;
        --primary-dark: #e55a2e !important;
        --secondary-color: #ffffff !important;
        --bg-primary: #1a1a1a !important;
        --bg-secondary: #2d2d2d !important;
        --bg-card: #333333 !important;
        --text-primary: #ffffff !important;
        --text-secondary: #cccccc !important;
        --text-light: #999999 !important;
      }
      
      body {
        background-color: var(--bg-secondary) !important;
        color: var(--text-primary) !important;
      }
      
      .navbar {
        background: linear-gradient(135deg, #ff6136 0%, #e55a2e 100%) !important;
      }
      
      .card, .section, .main-content, .page {
        background-color: var(--bg-card) !important;
        color: var(--text-primary) !important;
      }
      
      .form-input, .form-select {
        background-color: var(--bg-primary) !important;
        color: var(--text-primary) !important;
        border-color: #555555 !important;
      }
      
      .stat-card {
        background-color: var(--bg-card) !important;
        border: 1px solid #555555 !important;
      }
      
      /* Conserver les couleurs SEPTEO pour les éléments importants */
      .btn-primary {
        background: linear-gradient(135deg, #ff6136 0%, #e55a2e 100%) !important;
        color: white !important;
      }
      
      .status-dot.status-healthy {
        background-color: #27ae60 !important;
      }
      
      .chart-container {
        background-color: var(--bg-card) !important;
      }
    `;
    
    const styleElement = document.createElement('style');
    styleElement.textContent = darkTheme;
    document.head.appendChild(styleElement);
    
    // Ajouter un indicateur de mode sombre
    const darkModeIndicator = document.createElement('div');
    darkModeIndicator.style.cssText = `
      position: fixed;
      top: 10px;
      left: 10px;
      background: #ff6136;
      color: white;
      padding: 5px 10px;
      border-radius: 4px;
      font-size: 12px;
      z-index: 9999;
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    `;
    darkModeIndicator.textContent = '🌙 Dark Mode';
    document.body.appendChild(darkModeIndicator);
  });
  
  // Attendre que les styles s'appliquent
  await page.waitForTimeout(1000);
  
  console.log('✅ Dark mode enabled with SEPTEO branding preserved');
};