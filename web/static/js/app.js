/**
 * Fire Salamander Web Interface
 * Application JavaScript principale
 */

class FireSalamanderApp {
    constructor() {
        this.currentAuditId = null;
        this.progressInterval = null;
        this.apiBaseUrl = '/api/v1';
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.loadAuditHistory();
        this.setupTabs();
    }

    setupEventListeners() {
        // Formulaire d'audit
        const auditForm = document.getElementById('audit-form');
        if (auditForm) {
            auditForm.addEventListener('submit', (e) => this.handleAuditSubmit(e));
        }

        // Boutons d'export
        const exportJsonBtn = document.getElementById('export-json-btn');
        const exportHtmlBtn = document.getElementById('export-html-btn');
        
        if (exportJsonBtn) {
            exportJsonBtn.addEventListener('click', () => this.exportResults('json'));
        }
        
        if (exportHtmlBtn) {
            exportHtmlBtn.addEventListener('click', () => this.exportResults('html'));
        }
    }

    setupTabs() {
        const tabBtns = document.querySelectorAll('.tab-btn');
        
        tabBtns.forEach(btn => {
            btn.addEventListener('click', () => {
                const targetTab = btn.dataset.tab;
                this.switchTab(targetTab);
            });
        });
    }

    switchTab(targetTab) {
        // Désactiver tous les onglets et panneaux
        document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
        document.querySelectorAll('.tab-pane').forEach(pane => pane.classList.remove('active'));
        
        // Activer l'onglet et le panneau ciblé
        document.querySelector(`[data-tab="${targetTab}"]`).classList.add('active');
        document.getElementById(`${targetTab}-tab`).classList.add('active');
    }

    async handleAuditSubmit(e) {
        e.preventDefault();
        
        const form = e.target;
        const formData = new FormData(form);
        
        const auditRequest = {
            siteUrl: formData.get('siteUrl'),
            auditType: formData.get('auditType'),
            maxPages: parseInt(formData.get('maxPages')),
            timestamp: new Date().toISOString()
        };

        try {
            this.setFormLoading(true);
            
            // Simuler l'appel API pour le MVP
            // En production, ceci ferait appel à l'orchestrateur V2
            const response = await this.simulateAuditStart(auditRequest);
            
            this.currentAuditId = response.auditId;
            this.showProgressSection(auditRequest);
            this.startProgressMonitoring();
            
        } catch (error) {
            console.error('Erreur lors du démarrage de l\'audit:', error);
            this.showError('Erreur lors du démarrage de l\'audit. Veuillez réessayer.');
            this.setFormLoading(false);
        }
    }

    setFormLoading(loading) {
        const btn = document.getElementById('start-audit-btn');
        const btnText = btn.querySelector('.btn-text');
        const btnLoading = btn.querySelector('.btn-loading');
        
        if (loading) {
            btn.disabled = true;
            btnText.style.display = 'none';
            btnLoading.style.display = 'flex';
        } else {
            btn.disabled = false;
            btnText.style.display = 'block';
            btnLoading.style.display = 'none';
        }
    }

    showProgressSection(auditRequest) {
        // Mettre à jour les informations d'audit
        document.getElementById('current-audit-id').textContent = this.currentAuditId;
        document.getElementById('current-audit-url').textContent = auditRequest.siteUrl;
        
        // Afficher la section de progression
        const progressSection = document.getElementById('progress-section');
        progressSection.style.display = 'block';
        progressSection.classList.add('fade-in');
        
        // Réinitialiser les statuts des agents
        this.resetAgentStatuses();
        
        // Faire défiler vers la section
        progressSection.scrollIntoView({ behavior: 'smooth' });
    }

    resetAgentStatuses() {
        const agentCards = document.querySelectorAll('.agent-card');
        agentCards.forEach(card => {
            const status = card.querySelector('.agent-status');
            status.textContent = 'En attente';
            status.setAttribute('data-status', 'pending');
        });
    }

    updateAgentStatus(agentName, status, message = null) {
        const agentCard = document.querySelector(`[data-agent="${agentName}"]`);
        if (agentCard) {
            const statusElement = agentCard.querySelector('.agent-status');
            
            const statusMessages = {
                'pending': 'En attente',
                'running': 'En cours...',
                'completed': 'Terminé',
                'failed': 'Échec'
            };
            
            statusElement.textContent = message || statusMessages[status];
            statusElement.setAttribute('data-status', status);
        }
    }

    updateProgress(percentage, step) {
        const progressFill = document.getElementById('progress-fill');
        const progressPercentage = document.getElementById('progress-percentage');
        const progressStep = document.getElementById('progress-step');
        
        if (progressFill) {
            progressFill.style.width = `${percentage}%`;
        }
        
        if (progressPercentage) {
            progressPercentage.textContent = `${Math.round(percentage)}%`;
        }
        
        if (progressStep && step) {
            progressStep.textContent = step;
        }
    }

    async startProgressMonitoring() {
        // Simulation du monitoring de progression pour le MVP
        let progress = 0;
        const agents = ['crawler', 'keyword_extractor', 'technical_auditor', 'linking_mapper', 'broken_links_detector'];
        let currentAgentIndex = 0;
        
        this.progressInterval = setInterval(async () => {
            progress += Math.random() * 15 + 5; // Progression simulée
            
            if (progress >= 100) {
                progress = 100;
                this.updateProgress(progress, 'Audit terminé');
                this.completeAudit();
                clearInterval(this.progressInterval);
                return;
            }
            
            // Simuler le progrès des agents
            const currentAgent = agents[currentAgentIndex];
            if (currentAgent && progress > (currentAgentIndex + 1) * 20) {
                this.updateAgentStatus(currentAgent, 'completed');
                currentAgentIndex++;
                if (currentAgentIndex < agents.length) {
                    this.updateAgentStatus(agents[currentAgentIndex], 'running');
                }
            }
            
            this.updateProgress(progress, `Analyse en cours...`);
            
        }, 1000);
        
        // Démarrer le premier agent
        if (agents.length > 0) {
            this.updateAgentStatus(agents[0], 'running');
        }
    }

    async completeAudit() {
        try {
            // Simuler la récupération des résultats
            const results = await this.simulateAuditResults();
            
            this.displayResults(results);
            this.saveAuditToHistory(results);
            this.setFormLoading(false);
            
        } catch (error) {
            console.error('Erreur lors de la récupération des résultats:', error);
            this.showError('Erreur lors de la récupération des résultats.');
            this.setFormLoading(false);
        }
    }

    displayResults(results) {
        // Mettre à jour le résumé
        document.getElementById('total-pages').textContent = results.summary.totalPages;
        document.getElementById('total-keywords').textContent = results.summary.totalKeywords;
        document.getElementById('broken-links').textContent = results.summary.brokenLinks;
        document.getElementById('seo-score').textContent = results.summary.averageSeoScore + '/100';
        
        // Remplir les onglets de contenu
        this.populateOverviewTab(results);
        this.populateKeywordsTab(results.keywords);
        this.populateTechnicalTab(results.technical);
        this.populateLinksTab(results.links);
        
        // Afficher la section des résultats
        const resultsSection = document.getElementById('results-section');
        resultsSection.style.display = 'block';
        resultsSection.classList.add('slide-up');
        
        // Faire défiler vers les résultats
        resultsSection.scrollIntoView({ behavior: 'smooth' });
    }

    populateOverviewTab(results) {
        const overviewContent = document.getElementById('overview-content');
        
        overviewContent.innerHTML = `
            <div class="overview-grid">
                <div class="overview-section">
                    <h4>📊 Aperçu général</h4>
                    <ul class="overview-list">
                        <li><strong>Pages analysées:</strong> ${results.summary.totalPages}</li>
                        <li><strong>Durée de l'audit:</strong> ${results.summary.duration}</li>
                        <li><strong>Score SEO moyen:</strong> ${results.summary.averageSeoScore}/100</li>
                        <li><strong>Problèmes détectés:</strong> ${results.summary.totalIssues}</li>
                    </ul>
                </div>
                
                <div class="overview-section">
                    <h4>🎯 Recommandations principales</h4>
                    <ul class="recommendations-list">
                        ${results.recommendations.slice(0, 5).map(rec => `<li>${rec}</li>`).join('')}
                    </ul>
                </div>
            </div>
        `;
    }

    populateKeywordsTab(keywords) {
        const keywordsContent = document.getElementById('keywords-content');
        
        if (!keywords || keywords.length === 0) {
            keywordsContent.innerHTML = '<p>Aucun mot-clé détecté.</p>';
            return;
        }
        
        const keywordsHtml = `
            <div class="keywords-section">
                <h4>🔍 Mots-clés les plus pertinents</h4>
                <div class="keywords-grid">
                    ${keywords.slice(0, 20).map(keyword => `
                        <div class="keyword-card">
                            <div class="keyword-term">${keyword.term}</div>
                            <div class="keyword-metrics">
                                <span class="keyword-count">${keyword.count} occurrences</span>
                                <span class="keyword-density">${keyword.density.toFixed(1)}% densité</span>
                            </div>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
        
        keywordsContent.innerHTML = keywordsHtml;
    }

    populateTechnicalTab(technical) {
        const technicalContent = document.getElementById('technical-content');
        
        if (!technical) {
            technicalContent.innerHTML = '<p>Aucune donnée technique disponible.</p>';
            return;
        }
        
        const technicalHtml = `
            <div class="technical-section">
                <div class="scores-grid">
                    <div class="score-card">
                        <div class="score-value">${technical.performance || 0}/100</div>
                        <div class="score-label">Performance</div>
                    </div>
                    <div class="score-card">
                        <div class="score-value">${technical.accessibility || 0}/100</div>
                        <div class="score-label">Accessibilité</div>
                    </div>
                    <div class="score-card">
                        <div class="score-value">${technical.seo || 0}/100</div>
                        <div class="score-label">SEO technique</div>
                    </div>
                </div>
                
                <div class="issues-section">
                    <h4>⚠️ Problèmes détectés</h4>
                    <div class="issues-list">
                        ${technical.issues ? technical.issues.map(issue => `
                            <div class="issue-item ${issue.severity}">
                                <strong>${issue.type}:</strong> ${issue.description}
                            </div>
                        `).join('') : '<p>Aucun problème détecté.</p>'}
                    </div>
                </div>
            </div>
        `;
        
        technicalContent.innerHTML = technicalHtml;
    }

    populateLinksTab(links) {
        const linksContent = document.getElementById('links-content');
        
        if (!links) {
            linksContent.innerHTML = '<p>Aucune donnée de liens disponible.</p>';
            return;
        }
        
        const linksHtml = `
            <div class="links-section">
                <div class="links-stats">
                    <div class="stat-card">
                        <div class="stat-value">${links.totalLinks || 0}</div>
                        <div class="stat-label">Total liens</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-value">${links.internalLinks || 0}</div>
                        <div class="stat-label">Liens internes</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-value">${links.externalLinks || 0}</div>
                        <div class="stat-label">Liens externes</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-value">${links.brokenLinks || 0}</div>
                        <div class="stat-label">Liens brisés</div>
                    </div>
                </div>
                
                ${links.brokenLinks > 0 ? `
                    <div class="broken-links-section">
                        <h4>🚫 Liens brisés détectés</h4>
                        <div class="broken-links-list">
                            ${links.brokenLinksDetails ? links.brokenLinksDetails.slice(0, 10).map(link => `
                                <div class="broken-link-item">
                                    <div class="broken-url">${link.url}</div>
                                    <div class="broken-error">Status: ${link.statusCode} - ${link.error}</div>
                                </div>
                            `).join('') : ''}
                        </div>
                    </div>
                ` : ''}
            </div>
        `;
        
        linksContent.innerHTML = linksHtml;
    }

    saveAuditToHistory(results) {
        let history = JSON.parse(localStorage.getItem('auditHistory') || '[]');
        
        const auditRecord = {
            id: this.currentAuditId,
            url: results.siteUrl,
            date: new Date().toISOString(),
            status: 'completed',
            summary: results.summary,
            results: results
        };
        
        history.unshift(auditRecord);
        
        // Garder seulement les 10 derniers audits
        if (history.length > 10) {
            history = history.slice(0, 10);
        }
        
        localStorage.setItem('auditHistory', JSON.stringify(history));
        this.loadAuditHistory();
    }

    loadAuditHistory() {
        const history = JSON.parse(localStorage.getItem('auditHistory') || '[]');
        const historyContainer = document.getElementById('audit-history');
        
        if (history.length === 0) {
            historyContainer.innerHTML = '<p class="no-history">Aucun audit récent.</p>';
            return;
        }
        
        const historyHtml = history.map(audit => `
            <div class="history-item" data-audit-id="${audit.id}">
                <div class="history-meta">
                    <div class="history-url">${audit.url}</div>
                    <div class="history-date">${new Date(audit.date).toLocaleString()}</div>
                </div>
                <div class="history-status ${audit.status}">${audit.status}</div>
            </div>
        `).join('');
        
        historyContainer.innerHTML = historyHtml;
        
        // Ajouter les listeners pour charger les résultats précédents
        historyContainer.querySelectorAll('.history-item').forEach(item => {
            item.addEventListener('click', () => {
                const auditId = item.dataset.auditId;
                this.loadPreviousAudit(auditId);
            });
        });
    }

    loadPreviousAudit(auditId) {
        const history = JSON.parse(localStorage.getItem('auditHistory') || '[]');
        const audit = history.find(a => a.id === auditId);
        
        if (audit && audit.results) {
            this.displayResults(audit.results);
        }
    }

    async exportResults(format) {
        if (!this.currentAuditId) {
            this.showError('Aucun audit à exporter.');
            return;
        }
        
        try {
            // En production, ceci ferait appel à l'API d'export
            const results = await this.getCurrentResults();
            
            if (format === 'json') {
                this.downloadJSON(results);
            } else if (format === 'html') {
                this.downloadHTML(results);
            }
            
        } catch (error) {
            console.error('Erreur lors de l\'export:', error);
            this.showError('Erreur lors de l\'export des résultats.');
        }
    }

    downloadJSON(data) {
        const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `fire-salamander-audit-${this.currentAuditId}.json`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }

    downloadHTML(data) {
        const htmlContent = this.generateHTMLReport(data);
        const blob = new Blob([htmlContent], { type: 'text/html' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `fire-salamander-audit-${this.currentAuditId}.html`;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    }

    generateHTMLReport(data) {
        return `
            <!DOCTYPE html>
            <html>
            <head>
                <meta charset="UTF-8">
                <title>Rapport SEO - ${data.siteUrl}</title>
                <style>
                    body { font-family: Arial, sans-serif; margin: 20px; line-height: 1.6; }
                    .header { border-bottom: 2px solid #2563eb; padding-bottom: 20px; margin-bottom: 30px; }
                    .section { margin-bottom: 30px; }
                    .score { font-size: 24px; font-weight: bold; color: #2563eb; }
                    .issue { margin-bottom: 10px; padding: 10px; background: #f8f9fa; border-left: 4px solid #dc3545; }
                </style>
            </head>
            <body>
                <div class="header">
                    <h1>🦎 Rapport d'audit SEO Fire Salamander</h1>
                    <p><strong>Site analysé:</strong> ${data.siteUrl}</p>
                    <p><strong>Date:</strong> ${new Date(data.timestamp).toLocaleString()}</p>
                </div>
                
                <div class="section">
                    <h2>Résumé</h2>
                    <p><strong>Pages analysées:</strong> ${data.summary.totalPages}</p>
                    <p><strong>Score SEO moyen:</strong> <span class="score">${data.summary.averageSeoScore}/100</span></p>
                    <p><strong>Mots-clés trouvés:</strong> ${data.summary.totalKeywords}</p>
                    <p><strong>Liens brisés:</strong> ${data.summary.brokenLinks}</p>
                </div>
                
                <!-- Sections détaillées seraient ajoutées ici -->
                
                <div class="footer">
                    <p><em>Rapport généré par Fire Salamander - Plateforme d'audit SEO</em></p>
                </div>
            </body>
            </html>
        `;
    }

    showError(message) {
        // Créer une notification d'erreur simple
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error-notification';
        errorDiv.textContent = message;
        errorDiv.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            background: #dc2626;
            color: white;
            padding: 15px 20px;
            border-radius: 8px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            z-index: 1000;
            animation: slideIn 0.3s ease-out;
        `;
        
        document.body.appendChild(errorDiv);
        
        setTimeout(() => {
            errorDiv.remove();
        }, 5000);
    }

    // Méthodes de simulation pour le MVP
    async simulateAuditStart(auditRequest) {
        // Simuler un délai d'API
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        return {
            auditId: 'aud_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9),
            status: 'started',
            message: 'Audit démarré avec succès'
        };
    }

    async simulateAuditResults() {
        // Simuler des résultats d'audit réalistes
        return {
            siteUrl: document.getElementById('site-url').value,
            timestamp: new Date().toISOString(),
            summary: {
                totalPages: Math.floor(Math.random() * 20) + 5,
                totalKeywords: Math.floor(Math.random() * 100) + 20,
                brokenLinks: Math.floor(Math.random() * 5),
                averageSeoScore: Math.floor(Math.random() * 40) + 60,
                duration: '2m 34s',
                totalIssues: Math.floor(Math.random() * 10) + 2
            },
            keywords: this.generateSampleKeywords(),
            technical: this.generateSampleTechnical(),
            links: this.generateSampleLinks(),
            recommendations: [
                'Améliorer les titres des pages (trop courts)',
                'Ajouter des descriptions meta manquantes',
                'Optimiser les images sans attribut alt',
                'Corriger les liens internes brisés',
                'Améliorer la structure des headings H1-H6'
            ]
        };
    }

    generateSampleKeywords() {
        const sampleKeywords = [
            'seo', 'audit', 'analyse', 'référencement', 'optimisation',
            'contenu', 'marketing', 'digital', 'web', 'site',
            'performance', 'technique', 'liens', 'mots-clés', 'meta'
        ];
        
        return sampleKeywords.map(term => ({
            term,
            count: Math.floor(Math.random() * 20) + 1,
            density: Math.random() * 3 + 0.5,
            relevance: Math.random() * 2 + 1
        })).sort((a, b) => b.relevance - a.relevance);
    }

    generateSampleTechnical() {
        return {
            performance: Math.floor(Math.random() * 30) + 70,
            accessibility: Math.floor(Math.random() * 40) + 60,
            seo: Math.floor(Math.random() * 30) + 70,
            issues: [
                { type: 'SEO', severity: 'medium', description: 'Titre trop court sur la page d\'accueil' },
                { type: 'Accessibilité', severity: 'high', description: 'Images sans attribut alt détectées' },
                { type: 'Performance', severity: 'low', description: 'Resources non compressées trouvées' }
            ]
        };
    }

    generateSampleLinks() {
        const brokenCount = Math.floor(Math.random() * 3);
        return {
            totalLinks: Math.floor(Math.random() * 100) + 20,
            internalLinks: Math.floor(Math.random() * 50) + 10,
            externalLinks: Math.floor(Math.random() * 30) + 5,
            brokenLinks: brokenCount,
            brokenLinksDetails: brokenCount > 0 ? [
                { url: 'https://example.com/broken-page', statusCode: 404, error: 'Not Found' },
                { url: 'https://external-site.com/missing', statusCode: 500, error: 'Internal Server Error' }
            ].slice(0, brokenCount) : []
        };
    }

    async getCurrentResults() {
        // Récupérer les résultats actuels depuis l'historique local
        const history = JSON.parse(localStorage.getItem('auditHistory') || '[]');
        const current = history.find(audit => audit.id === this.currentAuditId);
        return current ? current.results : null;
    }
}

// Initialiser l'application quand le DOM est prêt
document.addEventListener('DOMContentLoaded', () => {
    new FireSalamanderApp();
});