module.exports = async (page, scenario, vp) => {
  console.log('READY > ' + scenario.label);
  
  // Attendre que les stats se chargent
  await page.waitForFunction(() => {
    const statsElements = document.querySelectorAll('.stat-value');
    return statsElements.length > 0;
  }, { timeout: 5000 }).catch(() => {
    console.warn('Stats did not load within timeout');
  });
  
  // Masquer les éléments temporels qui changent constamment
  await page.evaluate(() => {
    // Masquer timestamps et dates
    const timeElements = document.querySelectorAll('[data-time], .timestamp, .last-updated');
    timeElements.forEach(el => el.style.visibility = 'hidden');
    
    // Remplacer les valeurs dynamiques par des valeurs fixes pour les tests
    const statValues = document.querySelectorAll('.stat-value');
    if (statValues.length > 0) {
      statValues[0].textContent = '42'; // Analyses totales
      if (statValues[1]) statValues[1].textContent = '38'; // Succès
      if (statValues[2]) statValues[2].textContent = '2.3s'; // Temps moyen
      if (statValues[3]) statValues[3].textContent = '85'; // Score moyen
    }
    
    // Fixer l'heure dans le status
    const statusTime = document.querySelector('.status-text');
    if (statusTime) {
      statusTime.textContent = 'En ligne';
    }
  });
  
  // Attendre que les polices se chargent
  await page.evaluateHandle('document.fonts.ready');
  
  // Attendre un peu plus pour la stabilité visuelle
  await page.waitForTimeout(1000);
};