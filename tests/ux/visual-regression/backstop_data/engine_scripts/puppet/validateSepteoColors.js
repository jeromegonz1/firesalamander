module.exports = async (page, scenario, vp) => {
  console.log('VALIDATING SEPTEO COLORS > ' + scenario.label);
  
  const SEPTEO_COLORS = {
    primary: '#ff6136',
    primaryDark: '#e55a2e',
    secondary: '#2c3e50',
    success: '#27ae60',
    warning: '#f39c12',
    danger: '#e74c3c'
  };
  
  // Valider les couleurs SEPTEO
  const colorValidation = await page.evaluate((colors) => {
    const issues = [];
    
    // Vérifier la couleur primaire
    const primaryElements = document.querySelectorAll('.btn-primary, .navbar, [style*="' + colors.primary + '"]');
    primaryElements.forEach((el, index) => {
      const style = window.getComputedStyle(el);
      const bgColor = style.backgroundColor;
      const color = style.color;
      
      // Convertir RGB vers HEX pour comparaison
      const rgbToHex = (rgb) => {
        const result = rgb.match(/\d+/g);
        if (!result) return null;
        return '#' + ((1 << 24) + (parseInt(result[0]) << 16) + (parseInt(result[1]) << 8) + parseInt(result[2])).toString(16).slice(1);
      };
      
      const hexBg = rgbToHex(bgColor);
      const hexColor = rgbToHex(color);
      
      console.log(`Element ${index}: bg=${hexBg}, color=${hexColor}`);
    });
    
    // Vérifier le contraste
    const checkContrast = (fg, bg) => {
      // Implémentation simplifiée du calcul de contraste WCAG
      const getLuminance = (hex) => {
        const rgb = parseInt(hex.slice(1), 16);
        const r = (rgb >> 16) & 0xff;
        const g = (rgb >> 8) & 0xff;
        const b = (rgb >> 0) & 0xff;
        
        const rsRGB = r / 255;
        const gsRGB = g / 255;
        const bsRGB = b / 255;
        
        const rLum = rsRGB <= 0.03928 ? rsRGB / 12.92 : Math.pow((rsRGB + 0.055) / 1.055, 2.4);
        const gLum = gsRGB <= 0.03928 ? gsRGB / 12.92 : Math.pow((gsRGB + 0.055) / 1.055, 2.4);
        const bLum = bsRGB <= 0.03928 ? bsRGB / 12.92 : Math.pow((bsRGB + 0.055) / 1.055, 2.4);
        
        return 0.2126 * rLum + 0.7152 * gLum + 0.0722 * bLum;
      };
      
      const l1 = getLuminance(fg);
      const l2 = getLuminance(bg);
      const ratio = (Math.max(l1, l2) + 0.05) / (Math.min(l1, l2) + 0.05);
      return ratio;
    };
    
    // Vérifier le contraste minimum (4.5:1 pour WCAG AA)
    const orangeOnWhite = checkContrast(colors.primary, '#ffffff');
    const whiteOnOrange = checkContrast('#ffffff', colors.primary);
    
    if (orangeOnWhite < 4.5) {
      issues.push(`Orange SEPTEO sur blanc: contraste insuffisant (${orangeOnWhite.toFixed(2)}:1)`);
    }
    
    if (whiteOnOrange < 4.5) {
      issues.push(`Blanc sur orange SEPTEO: contraste insuffisant (${whiteOnOrange.toFixed(2)}:1)`);
    }
    
    return {
      valid: issues.length === 0,
      issues: issues,
      contrasts: {
        orangeOnWhite: orangeOnWhite.toFixed(2),
        whiteOnOrange: whiteOnOrange.toFixed(2)
      }
    };
  }, SEPTEO_COLORS);
  
  console.log('Color validation result:', colorValidation);
  
  if (!colorValidation.valid) {
    console.warn('⚠️  SEPTEO Color issues detected:', colorValidation.issues);
  } else {
    console.log('✅ SEPTEO Colors validated successfully');
  }
  
  // Ajouter des marqueurs visuels pour les tests
  await page.evaluate(() => {
    const style = document.createElement('style');
    style.textContent = `
      .color-test-marker {
        position: fixed;
        top: 10px;
        right: 10px;
        background: #ff6136;
        color: white;
        padding: 5px 10px;
        border-radius: 4px;
        font-size: 12px;
        z-index: 9999;
      }
    `;
    document.head.appendChild(style);
    
    const marker = document.createElement('div');
    marker.className = 'color-test-marker';
    marker.textContent = 'SEPTEO ✓';
    document.body.appendChild(marker);
  });
};