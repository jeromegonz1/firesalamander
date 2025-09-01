module.exports = async (page, scenario, vp) => {
  console.log('SCENARIO > ' + scenario.label);
  
  // Attendre que Fire Salamander soit prêt
  await page.waitForTimeout(2000);
  
  // Désactiver les animations pour des tests visuels consistants
  await page.addStyleTag({
    content: `
      *, *::before, *::after {
        animation-duration: 0s !important;
        animation-delay: 0s !important;
        transition-duration: 0s !important;
        transition-delay: 0s !important;
      }
    `
  });
  
  // Vérifier que Fire Salamander est démarré
  try {
    const response = await page.goto('http://localhost:8080/api/v1/health', { 
      waitUntil: 'networkidle0',
      timeout: 10000 
    });
    
    if (!response.ok()) {
      throw new Error(`Fire Salamander not ready: ${response.status()}`);
    }
    
    console.log('✅ Fire Salamander is ready');
  } catch (error) {
    console.error('❌ Fire Salamander health check failed:', error.message);
    throw error;
  }
};