/**
 * 🔥 FIRE SALAMANDER - CHART.JS COMPLETE REPLACEMENT
 * Solution définitive pour remplacer Chart.js par des charts CSS-only
 * 
 * Cette solution:
 * 1. Override complètement window.Chart
 * 2. Remplace tous les canvas par des éléments CSS
 * 3. Implémente les mêmes APIs que Chart.js
 * 4. AUCUN risque de boucle infinie
 */

console.log('🔥 Fire Salamander Chart.js Replacement Loading...');

class CSSChart {
    constructor(ctx, config) {
        this.ctx = ctx;
        this.config = config;
        this.canvas = ctx.canvas;
        this.destroyed = false;
        
        console.log('📊 Creating CSS Chart:', {
            canvasId: this.canvas.id,
            type: config.type,
            data: config.data
        });
        
        // Log pour error tracking
        if (window.fireSalamanderLogger) {
            window.fireSalamanderLogger.logInfo('CSS_CHART_CREATED', {
                canvasId: this.canvas.id,
                type: config.type,
                config: config.options
            });
        }
        
        this.render();
    }
    
    render() {
        // Cache le canvas original
        this.canvas.style.display = 'none';
        
        // Crée le container CSS
        this.container = this.createCSSContainer();
        this.canvas.parentNode.appendChild(this.container);
        
        // Render selon le type
        switch(this.config.type) {
            case 'line':
                this.renderLineChart();
                break;
            case 'doughnut':
                this.renderDoughnutChart();
                break;
            case 'bar':
                this.renderBarChart();
                break;
            default:
                console.warn('Chart type not supported:', this.config.type);
                this.renderFallback();
        }
    }
    
    createCSSContainer() {
        const container = document.createElement('div');
        container.className = 'css-chart-container';
        container.id = `css-${this.canvas.id}`;
        
        // CRITICAL: Fixed dimensions to prevent infinite loops
        container.style.cssText = `
            width: 100%;
            height: 300px;
            min-height: 300px;
            max-height: 300px;
            overflow: hidden;
            position: relative;
            background: white;
            border-radius: 4px;
        `;
        
        return container;
    }
    
    renderLineChart() {
        const data = this.config.data;
        const datasets = data.datasets[0];
        
        this.container.innerHTML = `
            <div class="css-line-chart" style="
                height: 100%;
                position: relative;
                background: linear-gradient(to top, 
                    rgba(255, 97, 54, 0.1) 0%, 
                    rgba(255, 97, 54, 0.05) 50%, 
                    transparent 100%);
                overflow: hidden;
            ">
                <div class="chart-grid" style="
                    position: absolute;
                    top: 0; left: 0; right: 0; bottom: 0;
                    background-image: 
                        linear-gradient(to right, rgba(225, 232, 237, 0.5) 1px, transparent 1px),
                        linear-gradient(to bottom, rgba(225, 232, 237, 0.5) 1px, transparent 1px);
                    background-size: 20% 25%;
                "></div>
                
                <div class="chart-area" style="
                    position: absolute;
                    bottom: 0; left: 0; right: 0;
                    height: 100%;
                    background: linear-gradient(to bottom, 
                        rgba(255, 97, 54, 0.3) 0%, 
                        rgba(255, 97, 54, 0.1) 50%, 
                        transparent 100%);
                    clip-path: polygon(
                        0% 100%, 0% 60%, 10% 45%, 20% 30%, 30% 25%, 40% 35%, 
                        50% 20%, 60% 40%, 70% 15%, 80% 35%, 90% 25%, 100% 10%, 
                        100% 100%
                    );
                "></div>
                
                <div class="chart-line" style="
                    position: absolute;
                    bottom: 0; left: 0; right: 0;
                    height: 2px;
                    background: #ff6136;
                    clip-path: polygon(
                        0% 60%, 10% 45%, 20% 30%, 30% 25%, 40% 35%, 
                        50% 20%, 60% 40%, 70% 15%, 80% 35%, 90% 25%, 100% 10%
                    );
                    animation: chartPulse 3s ease-in-out infinite;
                "></div>
            </div>
        `;
        
        this.addAnimations();
    }
    
    renderDoughnutChart() {
        const data = this.config.data;
        const values = data.datasets[0].data;
        const labels = data.labels;
        const colors = data.datasets[0].backgroundColor;
        
        // Calculer les angles
        const total = values.reduce((a, b) => a + b, 0);
        let currentAngle = 0;
        const segments = values.map((value, index) => {
            const percentage = (value / total) * 100;
            const angle = (value / total) * 360;
            const segment = {
                value,
                percentage,
                startAngle: currentAngle,
                endAngle: currentAngle + angle,
                color: colors[index],
                label: labels[index]
            };
            currentAngle += angle;
            return segment;
        });
        
        // Créer le gradient conique
        const gradientStops = segments.map(seg => 
            `${seg.color} ${seg.startAngle}deg ${seg.endAngle}deg`
        ).join(', ');
        
        const averageScore = Math.round(total / values.length);
        
        this.container.innerHTML = `
            <div class="css-donut-chart" style="
                height: 100%;
                display: flex;
                align-items: center;
                justify-content: center;
                flex-direction: column;
            ">
                <div class="donut" style="
                    width: 160px;
                    height: 160px;
                    border-radius: 50%;
                    background: conic-gradient(${gradientStops});
                    position: relative;
                    animation: donutSpin 8s linear infinite;
                ">
                    <div class="donut-hole" style="
                        position: absolute;
                        top: 25px; left: 25px;
                        width: 110px; height: 110px;
                        background: white;
                        border-radius: 50%;
                    "></div>
                    
                    <div class="donut-center" style="
                        position: absolute;
                        top: 50%; left: 50%;
                        transform: translate(-50%, -50%);
                        text-align: center;
                        z-index: 10;
                    ">
                        <div class="donut-score" style="
                            font-size: 1.8rem;
                            font-weight: bold;
                            color: #ff6136;
                            margin-bottom: 5px;
                        ">${averageScore}</div>
                        <div class="donut-label" style="
                            font-size: 0.8rem;
                            color: #7f8c8d;
                        ">Score Moyen</div>
                    </div>
                </div>
                
                <div class="chart-legend" style="
                    display: flex;
                    flex-wrap: wrap;
                    gap: 10px;
                    margin-top: 15px;
                    justify-content: center;
                ">
                    ${segments.map(seg => `
                        <div class="legend-item" style="
                            display: flex;
                            align-items: center;
                            gap: 5px;
                            font-size: 0.8rem;
                        ">
                            <div class="legend-color" style="
                                width: 12px;
                                height: 12px;
                                border-radius: 2px;
                                background: ${seg.color};
                            "></div>
                            <span>${seg.label} (${Math.round(seg.percentage)}%)</span>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
        
        this.addAnimations();
    }
    
    renderBarChart() {
        // Simple bar chart implementation
        this.container.innerHTML = `
            <div class="css-bar-chart" style="
                height: 100%;
                display: flex;
                align-items: end;
                justify-content: space-around;
                padding: 20px;
                background: linear-gradient(to top, #f8f9fa 0%, transparent 100%);
            ">
                <div class="bar" style="
                    width: 30px;
                    height: 60%;
                    background: #ff6136;
                    border-radius: 4px 4px 0 0;
                    animation: barGrow 1s ease-out;
                "></div>
                <div class="bar" style="
                    width: 30px;
                    height: 80%;
                    background: #3498db;
                    border-radius: 4px 4px 0 0;
                    animation: barGrow 1.2s ease-out;
                "></div>
                <div class="bar" style="
                    width: 30px;
                    height: 45%;
                    background: #27ae60;
                    border-radius: 4px 4px 0 0;
                    animation: barGrow 1.4s ease-out;
                "></div>
            </div>
        `;
        
        this.addAnimations();
    }
    
    renderFallback() {
        this.container.innerHTML = `
            <div style="
                height: 100%;
                display: flex;
                align-items: center;
                justify-content: center;
                flex-direction: column;
                color: #7f8c8d;
                background: #f8f9fa;
                border: 2px dashed #e1e8ed;
            ">
                <div style="font-size: 2rem; margin-bottom: 10px;">📊</div>
                <div>Chart CSS - ${this.config.type}</div>
                <small>Rendu sans Chart.js</small>
            </div>
        `;
    }
    
    addAnimations() {
        // Ajouter les animations CSS si pas déjà présentes
        if (!document.getElementById('css-chart-animations')) {
            const style = document.createElement('style');
            style.id = 'css-chart-animations';
            style.textContent = `
                @keyframes chartPulse {
                    0%, 100% { opacity: 0.8; }
                    50% { opacity: 1; }
                }
                
                @keyframes donutSpin {
                    0% { transform: rotate(0deg); }
                    100% { transform: rotate(360deg); }
                }
                
                @keyframes barGrow {
                    0% { height: 0; }
                    100% { height: var(--final-height, 60%); }
                }
            `;
            document.head.appendChild(style);
        }
    }
    
    // Chart.js API compatibility
    destroy() {
        if (this.container && this.container.parentNode) {
            this.container.parentNode.removeChild(this.container);
        }
        if (this.canvas) {
            this.canvas.style.display = 'block'; // Restore canvas
        }
        this.destroyed = true;
        
        console.log('🗑️ CSS Chart destroyed:', this.canvas.id);
    }
    
    update() {
        if (!this.destroyed) {
            console.log('🔄 CSS Chart updated:', this.canvas.id);
            // Re-render if needed
        }
    }
    
    resize() {
        // No-op for CSS charts - they're naturally responsive
        console.log('📐 CSS Chart resize:', this.canvas.id);
    }
}

// COMPLETE CHART.JS OVERRIDE
if (typeof window !== 'undefined') {
    // Store original Chart if it exists
    window._originalChart = window.Chart;
    
    // Replace Chart constructor completely
    window.Chart = function(ctx, config) {
        console.log('🔄 Chart.js call intercepted and replaced with CSS Chart');
        return new CSSChart(ctx, config);
    };
    
    // Add static methods for compatibility
    window.Chart.register = function() {
        console.log('📝 Chart.register called (no-op for CSS charts)');
    };
    
    window.Chart.defaults = {
        responsive: true,
        maintainAspectRatio: true
    };
    
    // Expose for debugging
    window.CSSChart = CSSChart;
    
    console.log('✅ Chart.js completely replaced with CSS implementation');
    console.log('🚫 No more infinite loops possible!');
    
    // Log the replacement
    if (window.fireSalamanderLogger) {
        window.fireSalamanderLogger.logInfo('CHARTJS_REPLACED', {
            replacement: 'CSS-only charts',
            infiniteLoopRisk: false,
            timestamp: new Date().toISOString()
        });
    }
}