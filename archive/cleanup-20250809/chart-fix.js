// Fix urgent pour la boucle infinie Chart.js
// Ce script monkey-patch Chart.js pour empêcher maintainAspectRatio: false

console.log('🔧 Application du fix Chart.js...');

// Override Chart constructor pour forcer maintainAspectRatio: true
const originalChart = window.Chart;
window.Chart = function(ctx, config) {
    // Force maintainAspectRatio à true pour éviter la boucle infinie
    if (config && config.options) {
        if (config.options.maintainAspectRatio === false) {
            console.log('⚠️ Correction: maintainAspectRatio forcé à true pour éviter boucle infinie');
            config.options.maintainAspectRatio = true;
            
            // Ajouter un aspectRatio approprié selon le type de chart
            if (config.type === 'line' || config.type === 'bar') {
                config.options.aspectRatio = 2; // Ratio 2:1 pour les graphiques linéaires
            } else if (config.type === 'doughnut' || config.type === 'pie') {
                config.options.aspectRatio = 1; // Ratio 1:1 pour les graphiques circulaires
            }
        }
    }
    
    return new originalChart(ctx, config);
};

// Conserver les propriétés statiques de Chart
Object.setPrototypeOf(window.Chart, originalChart);
Object.keys(originalChart).forEach(key => {
    window.Chart[key] = originalChart[key];
});

console.log('✅ Fix Chart.js appliqué avec succès');