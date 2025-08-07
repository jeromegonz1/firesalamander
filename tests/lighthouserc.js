module.exports = {
  ci: {
    collect: {
      // URLs à analyser
      url: [
        'http://localhost:8080',
        'http://localhost:8080/#analyzer',
        'http://localhost:8080/#history',
        'http://localhost:8080/#reports',
        'http://localhost:8080/#monitoring'
      ],
      
      // Paramètres de collecte
      numberOfRuns: 3,  // Moyenne sur 3 runs pour plus de fiabilité
      settings: {
        chromeFlags: '--no-sandbox --disable-setuid-sandbox --headless',
        preset: 'desktop',
        onlyCategories: ['performance', 'accessibility', 'best-practices', 'seo', 'pwa'],
        
        // Configuration spécifique pour Fire Salamander
        throttlingMethod: 'simulate',
        throttling: {
          rttMs: 40,      // RTT réseau simulé
          throughputKbps: 10240, // Bande passante simulée
          cpuSlowdownMultiplier: 1
        },
        
        // Désactiver certains audits non pertinents
        skipAudits: [
          'uses-http2',         // HTTP/2 pas forcément nécessaire
          'redirects-http',     // Redirections HTTP pas testables localement
          'is-on-https'         // HTTPS pas disponible en local
        ]
      }
    },
    
    // Seuils minimum SEPTEO (très stricts)
    assert: {
      assertions: {
        // Performance > 90
        'categories:performance': ['error', { minScore: 0.9 }],
        
        // Accessibilité > 95 (standard SEPTEO élevé)
        'categories:accessibility': ['error', { minScore: 0.95 }],
        
        // Bonnes pratiques > 90
        'categories:best-practices': ['error', { minScore: 0.9 }],
        
        // SEO > 85
        'categories:seo': ['error', { minScore: 0.85 }],
        
        // PWA (optionnel mais recommandé)
        'categories:pwa': ['warn', { minScore: 0.7 }],
        
        // Métriques Core Web Vitals (critiques pour UX)
        'largest-contentful-paint': ['error', { maxNumericValue: 2500 }],    // LCP < 2.5s
        'first-input-delay': ['error', { maxNumericValue: 100 }],           // FID < 100ms
        'cumulative-layout-shift': ['error', { maxNumericValue: 0.1 }],     // CLS < 0.1
        
        // Métriques de performance supplémentaires
        'first-contentful-paint': ['warn', { maxNumericValue: 1800 }],      // FCP < 1.8s
        'speed-index': ['warn', { maxNumericValue: 3000 }],                 // SI < 3s
        'interactive': ['warn', { maxNumericValue: 3800 }],                 // TTI < 3.8s
        
        // Audits spécifiques Fire Salamander
        'unused-css-rules': ['warn', { maxLength: 10 }],                    // CSS inutilisé
        'unused-javascript': ['warn', { maxLength: 10 }],                   // JS inutilisé
        'render-blocking-resources': ['warn', { maxLength: 3 }],            // Ressources bloquantes
        'offscreen-images': ['warn', { maxLength: 0 }],                     // Images hors écran
        
        // Accessibilité spécifique
        'color-contrast': ['error', { minScore: 1 }],                       // Contraste parfait requis
        'button-name': ['error', { minScore: 1 }],                          // Tous les boutons nommés
        'link-name': ['error', { minScore: 1 }],                            // Tous les liens nommés
        
        // SEO technique
        'meta-description': ['error', { minScore: 1 }],                     // Meta description obligatoire
        'document-title': ['error', { minScore: 1 }],                       // Titre obligatoire
        'html-has-lang': ['error', { minScore: 1 }]                         // Langue déclarée
      }
    },
    
    // Upload des résultats
    upload: {
      target: 'filesystem',
      outputDir: './reports/lighthouse',
      reportFilenamePattern: '%%PATHNAME%%-%%DATETIME%%-report.%%EXTENSION%%'
    },
    
    // Serveur pour les rapports
    server: {
      port: 9001,
      storage: {
        storageMethod: 'filesystem',
        storagePath: './reports/lighthouse-server'
      }
    }
  },
  
  // Configuration personnalisée SEPTEO
  septeoConfig: {
    brandColors: {
      primary: '#ff6136',
      secondary: '#2c3e50'
    },
    
    // Métriques business spécifiques
    businessMetrics: {
      // Temps max acceptable pour une analyse SEO
      maxAnalysisTime: 30000,  // 30 secondes
      
      // Score minimum acceptable pour les rapports
      minReportScore: 80,
      
      // Taille maximum des assets
      maxAssetSize: {
        css: 50000,    // 50KB max pour CSS
        js: 100000,    // 100KB max pour JS  
        images: 500000 // 500KB max pour images
      }
    },
    
    // Budgets de performance
    budgets: [
      {
        resourceType: 'total',
        budget: 500 // 500KB total max
      },
      {
        resourceType: 'script',
        budget: 150 // 150KB JS max
      },
      {
        resourceType: 'image',
        budget: 200 // 200KB images max
      },
      {
        resourceType: 'stylesheet',
        budget: 50  // 50KB CSS max
      }
    ]
  }
};