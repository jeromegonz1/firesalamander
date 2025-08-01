package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jeromegonz1/firesalamander/config"
)

var cfg *config.Config

func main() {
	// Charger la configuration
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	var err error
	cfg, err = config.Load(env)
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la configuration: %v", err)
	}

	// Configuration des routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)

	// Démarrage du serveur
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🔥 Fire Salamander %s démarré sur le port %d", cfg.App.Icon, cfg.Server.Port)
	log.Printf("🌐 Serveur accessible sur: http://localhost%s", addr)
	
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Erreur lors du démarrage du serveur: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.App.Name}} - {{.App.PoweredBy}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
            min-height: 100vh;
            color: #333;
        }

        .header {
            background: white;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            padding: 1rem 0;
        }

        .header-content {
            max-width: 1200px;
            margin: 0 auto;
            padding: 0 2rem;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }

        .logo-section {
            display: flex;
            align-items: center;
            gap: 1rem;
        }

        .septeo-logo {
            height: 40px;
        }

        .app-brand {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            font-size: 1.5rem;
            font-weight: 700;
            color: {{.Branding.PrimaryColor}};
        }

        .salamander {
            font-size: 2rem;
            filter: hue-rotate(10deg) brightness(1.1);
        }

        .powered-by {
            font-size: 0.9rem;
            color: #666;
            font-weight: 400;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 4rem 2rem;
        }

        .hero {
            text-align: center;
            margin-bottom: 4rem;
        }

        .hero h1 {
            font-size: 3rem;
            font-weight: 800;
            color: #2d3748;
            margin-bottom: 1rem;
            line-height: 1.2;
        }

        .hero .highlight {
            color: {{.Branding.PrimaryColor}};
        }

        .hero p {
            font-size: 1.25rem;
            color: #4a5568;
            margin-bottom: 2rem;
            max-width: 600px;
            margin-left: auto;
            margin-right: auto;
        }

        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 2rem;
            margin-top: 3rem;
        }

        .feature-card {
            background: white;
            padding: 2rem;
            border-radius: 12px;
            box-shadow: 0 4px 6px rgba(0,0,0,0.1);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }

        .feature-card:hover {
            transform: translateY(-4px);
            box-shadow: 0 8px 20px rgba(0,0,0,0.15);
        }

        .feature-icon {
            font-size: 2.5rem;
            margin-bottom: 1rem;
        }

        .feature-card h3 {
            font-size: 1.25rem;
            font-weight: 600;
            color: #2d3748;
            margin-bottom: 0.5rem;
        }

        .feature-card p {
            color: #4a5568;
            line-height: 1.6;
        }

        .cta-section {
            text-align: center;
            margin-top: 4rem;
        }

        .btn-primary {
            background: {{.Branding.PrimaryColor}};
            color: white;
            padding: 1rem 2rem;
            border: none;
            border-radius: 8px;
            font-size: 1.1rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            text-decoration: none;
            display: inline-block;
        }

        .btn-primary:hover {
            background: #e55a2b;
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(255,97,54,0.3);
        }

        .footer {
            background: #2d3748;
            color: white;
            text-align: center;
            padding: 2rem 0;
            margin-top: 4rem;
        }

        .status-indicator {
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            background: #10b981;
            color: white;
            padding: 0.5rem 1rem;
            border-radius: 20px;
            font-size: 0.9rem;
            font-weight: 500;
        }

        .status-dot {
            width: 8px;
            height: 8px;
            background: white;
            border-radius: 50%;
            animation: pulse 2s infinite;
        }

        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }

        @media (max-width: 768px) {
            .header-content {
                flex-direction: column;
                gap: 1rem;
            }
            
            .hero h1 {
                font-size: 2rem;
            }
            
            .features {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <header class="header">
        <div class="header-content">
            <div class="logo-section">
                <img src="https://cdn.prod.website-files.com/62f0eecf4db2ed8ccaf46cea/641d757432ef8bbe6de56c80_logo-septeo.svg" 
                     alt="SEPTEO" class="septeo-logo">
                <div class="app-brand">
                    <span class="salamander">{{.App.Icon}}</span>
                    <span>{{.App.Name}}</span>
                </div>
            </div>
            <div class="status-indicator">
                <span class="status-dot"></span>
                En ligne
            </div>
        </div>
    </header>

    <main class="container">
        <section class="hero">
            <h1>
                <span class="highlight">Fire Salamander</span><br>
                SEO Analyzer
            </h1>
            <p>
                L'outil d'analyse SEO nouvelle génération, alliant intelligence artificielle 
                et expertise technique pour optimiser vos performances web.
            </p>
            <div class="powered-by">Propulsé par {{.App.PoweredBy}}</div>
        </section>

        <section class="features">
            <div class="feature-card">
                <div class="feature-icon">🕷️</div>
                <h3>Crawling Intelligent</h3>
                <p>Exploration complète de vos sites web avec respect des bonnes pratiques et analyse en temps réel de la structure.</p>
            </div>

            <div class="feature-card">
                <div class="feature-icon">🧠</div>
                <h3>Analyse Sémantique IA</h3>
                <p>Compréhension avancée du contenu grâce à l'intelligence artificielle pour des recommandations précises.</p>
            </div>

            <div class="feature-card">
                <div class="feature-icon">📊</div>
                <h3>Audit SEO Complet</h3>
                <p>Analyse technique approfondie : performance, mots-clés, structure, maillage interne et recommandations prioritaires.</p>
            </div>

            <div class="feature-card">
                <div class="feature-icon">📈</div>
                <h3>Rapports Détaillés</h3>
                <p>Rapports visuels et actionnables avec insights IA, scores de performance et roadmap d'optimisation.</p>
            </div>

            <div class="feature-card">
                <div class="feature-icon">⚡</div>
                <h3>Performance Web</h3>
                <p>Mesure des Core Web Vitals, analyse de la vitesse de chargement et optimisations mobile-first.</p>
            </div>

            <div class="feature-card">
                <div class="feature-icon">🎯</div>
                <h3>Insights Concurrentiels</h3>
                <p>Analyse comparative et identification d'opportunités pour surpasser la concurrence.</p>
            </div>
        </section>

        <section class="cta-section">
            <a href="#" class="btn-primary">
                Commencer l'Analyse
            </a>
        </section>
    </main>

    <footer class="footer">
        <p>&copy; 2024 {{.App.Name}} - Propulsé par {{.App.PoweredBy}} | Version {{.App.Version}}</p>
    </footer>
</body>
</html>`

	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, "Erreur de template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, cfg); err != nil {
		http.Error(w, "Erreur d'exécution du template", http.StatusInternalServerError)
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{
		"status": "healthy",
		"app": "%s",
		"version": "%s",
		"environment": "%s"
	}`, cfg.App.Name, cfg.App.Version, os.Getenv("ENV"))
}