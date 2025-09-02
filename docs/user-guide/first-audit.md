# Lancer son premier audit

## Interface web

### 1. Accéder à l'interface
```bash
make server
# Ouvrir http://localhost:8080
```

### 2. Saisir l'URL
- Entrer l'URL complète (https://monsite.fr)
- Vérifier que le site est accessible
- Cliquer "Démarrer l'audit"

### 3. Suivre la progression
L'audit progresse en 4 étapes :
1. 🕷️ **Crawling** (1-5 min) - Exploration des pages
2. 🔧 **Analyse technique** (1-2 min) - SEO et performance
3. 🧠 **Analyse sémantique** (30s-2 min) - Mots-clés français
4. 📊 **Génération rapport** (10-30s) - PDF final

### 4. Télécharger le rapport
- Cliquer "Télécharger PDF" quand prêt
- Le rapport contient scores, recommandations et détails

## Via API REST

### Démarrer un audit
```bash
curl -X POST http://localhost:8080/api/audit \
  -H "Content-Type: application/json" \
  -d '{
    "seed_url": "https://camping-bretagne.fr",
    "options": {
      "max_urls": 100,
      "max_depth": 3
    }
  }'
```

Réponse :
```json
{
  "audit_id": "FS-PROD-001",
  "status": "started",
  "estimated_duration": "3-5 minutes"
}
```

### Suivre la progression
```bash
curl http://localhost:8080/api/audit/FS-PROD-001/status
```

Réponse :
```json
{
  "audit_id": "FS-PROD-001",
  "status": "analyzing_technical",
  "progress": 65.0,
  "current_step": "lighthouse_analysis",
  "estimated_completion": "2 minutes"
}
```

### Récupérer le rapport
```bash
# Rapport HTML
curl http://localhost:8080/api/audit/FS-PROD-001/report/html > rapport.html

# Rapport JSON  
curl http://localhost:8080/api/audit/FS-PROD-001/report/json > rapport.json
```

## Options d'audit

### Configuration basique
```json
{
  "seed_url": "https://monsite.fr",
  "options": {
    "max_urls": 300,
    "max_depth": 3,
    "respect_robots": true,
    "include_sitemap": true
  }
}
```

### Configuration avancée
```json
{
  "seed_url": "https://monsite.fr",
  "options": {
    "max_urls": 500,
    "max_depth": 4,
    "crawl_delay": 2000,
    "user_agent": "Fire Salamander Custom",
    "semantic_mode": "deep",
    "lighthouse_categories": ["seo", "performance", "accessibility"],
    "output_formats": ["html", "json", "csv"]
  }
}
```

## Troubleshooting

### Audit bloqué en "crawling"
- Site peut bloquer le crawler
- Vérifier robots.txt : `curl https://monsite.fr/robots.txt`
- Réduire max_urls à 50

### Pas de suggestions sémantiques
- Site doit être majoritairement en français
- Contenu minimum : 1000 mots
- Vérifier que Python semantic service est actif

### Rapport vide ou incomplet
- Vérifier logs : `make logs`
- Audit peut avoir échoué sur certains agents
- Mode dégradé génère quand même un rapport partiel