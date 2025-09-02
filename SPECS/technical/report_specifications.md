# Spécifications Report Engine - Sprint 1.5

## Formats de sortie supportés

### 1. HTML Report (Primaire)
```html
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Audit SEO - {site_url} - Fire Salamander</title>
    <style>/* CSS SEPTEO branding */</style>
</head>
<body>
    <header class="septeo-header">
        <img src="data:image/svg+xml;base64,{septeo_logo}" alt="SEPTEO">
        <h1>Fire Salamander - Audit SEO</h1>
    </header>
    
    <main>
        <section class="audit-summary">
            <h2>Résumé Audit {audit_id}</h2>
            <div class="metrics">
                <div class="score-card">
                    <span class="score">{overall_score}</span>
                    <span class="label">Score Global</span>
                </div>
                <div class="meta">
                    <p>Site: {site_url}</p>
                    <p>Date: {started_at}</p>
                    <p>Durée: {duration}</p>
                    <p>Pages: {total_pages}</p>
                </div>
            </div>
        </section>
        
        <section class="technical-findings">
            <h2>Analyse Technique</h2>
            {{range .TechResults.Findings}}
            <div class="finding severity-{{.Severity}}">
                <h3>{{.Message}}</h3>
                <ul>{{range .Evidence}}<li>{{.}}</li>{{end}}</ul>
            </div>
            {{end}}
        </section>
        
        <section class="semantic-keywords">
            <h2>Analyse Sémantique</h2>
            <div class="keywords">
                {{range .SemanticResults.Keywords}}
                <span class="keyword" data-score="{{.Score}}">{{.Keyword}}</span>
                {{end}}
            </div>
            <div class="suggestions">
                {{range .SemanticResults.Suggestions}}
                <div class="suggestion">
                    <strong>{{.Keyword}}</strong>
                    <p>{{.Reason}} ({{.Confidence}}%)</p>
                </div>
                {{end}}
            </div>
        </section>
        
        <section class="crawl-data">
            <h2>Données d'exploration</h2>
            <table>
                <thead>
                    <tr><th>URL</th><th>Titre</th><th>Status</th><th>Taille</th></tr>
                </thead>
                <tbody>
                    {{range .CrawlData.Pages}}
                    <tr>
                        <td><a href="{{.URL}}">{{.URL}}</a></td>
                        <td>{{.Title}}</td>
                        <td class="status-{{.StatusCode}}">{{.StatusCode}}</td>
                        <td>{{len .Content}} chars</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </section>
    </main>
    
    <footer class="septeo-footer">
        <p>Généré par Fire Salamander v{version} - © SEPTEO {year}</p>
        <p>Audit ID: {audit_id} | {timestamp}</p>
    </footer>
</body>
</html>
```

### 2. JSON Report (API/Intégration)
```json
{
  "audit_id": "FS-PROD-001",
  "site_url": "https://camping-bretagne.fr",
  "started_at": "2025-09-02T14:30:00Z",
  "completed_at": "2025-09-02T14:35:42Z",
  "duration": "5m42s",
  "status": "completed",
  "total_pages": 47,
  "overall_score": 76.5,
  "scores": {
    "technical": 82.0,
    "semantic": 71.0,
    "performance": 78.5,
    "mobile": 85.0
  },
  "crawl_data": {
    "pages_discovered": 47,
    "pages_crawled": 47,
    "errors": 0,
    "redirects": 5,
    "robots_txt": "allowed",
    "sitemap_found": true
  },
  "technical_analysis": {
    "findings_count": 12,
    "findings": [
      {
        "id": "missing-meta-description",
        "severity": "medium",
        "message": "Meta description manquante",
        "evidence": ["https://camping-bretagne.fr/contact"],
        "impact": "seo_visibility",
        "recommendation": "Ajouter meta description 150-160 chars"
      }
    ],
    "categories": {
      "titles": {"score": 90, "issues": 2},
      "headings": {"score": 85, "issues": 1},
      "meta": {"score": 70, "issues": 5},
      "links": {"score": 95, "issues": 0}
    }
  },
  "semantic_analysis": {
    "model_version": "camembert-base",
    "language": "fr",
    "keywords": [
      {
        "keyword": "camping familial bretagne",
        "score": 0.92,
        "frequency": 15,
        "contexts": ["title", "h1", "content"]
      }
    ],
    "topics": [
      {
        "id": "hebergement_outdoor",
        "confidence": 0.87,
        "keywords": ["camping", "mobil-home", "emplacement"]
      }
    ],
    "suggestions": [
      {
        "keyword": "location mobil-home mer",
        "reason": "Fort potentiel commercial détecté",
        "confidence": 0.78,
        "evidence": ["Page services", "Page tarifs"]
      }
    ]
  },
  "recommendations": {
    "priority_high": [
      {
        "id": "meta_descriptions",
        "title": "Compléter meta descriptions",
        "impact": "high",
        "effort": "medium",
        "estimated_time": "2-4h"
      }
    ],
    "priority_medium": [],
    "priority_low": []
  },
  "metadata": {
    "version": "1.0.0",
    "user_agent": "Fire Salamander/1.0",
    "schema_version": "1.0",
    "export_format": "json",
    "generated_at": "2025-09-02T14:35:42Z"
  }
}
```

### 3. CSV Export (Données)
```csv
audit_id,url,title,status_code,title_length,h1_count,meta_description,semantic_keywords,technical_score,semantic_score
FS-PROD-001,https://camping-bretagne.fr,Camping Familial Bretagne,200,25,1,yes,"camping familial|location mobil-home",85,92
FS-PROD-001,https://camping-bretagne.fr/services,Services Camping,200,16,1,no,"services camping|animations",70,78
FS-PROD-001,https://camping-bretagne.fr/tarifs,Tarifs Saison 2025,200,18,1,yes,"tarifs camping|prix location",90,82
```

## Templates et branding

### Template HTML responsive
- CSS Grid/Flexbox pour responsive
- Couleurs SEPTEO: #0066CC (bleu), #FF6600 (orange)
- Police: Roboto/Open Sans
- Logo SEPTEO en base64 intégré
- Print-friendly avec @media print

### Template variables
```go
type TemplateData struct {
    AuditID      string          `json:"audit_id"`
    SiteURL      string          `json:"site_url"`
    StartedAt    string          `json:"started_at"`
    Duration     string          `json:"duration"`
    TotalPages   int             `json:"total_pages"`
    OverallScore float64         `json:"overall_score"`
    TechResults  TechResults     `json:"tech_results"`
    SemanticResults SemanticResults `json:"semantic_results"`
    CrawlData    CrawlData       `json:"crawl_data"`
    Version      string          `json:"version"`
    GeneratedAt  string          `json:"generated_at"`
}
```

## Règles de génération

### Scoring consolidé
```go
// OverallScore calculation
func CalculateOverallScore(techScore, semanticScore, crawlScore float64) float64 {
    weights := map[string]float64{
        "technical": 0.4,  // 40% - SEO technique primordial
        "semantic":  0.4,  // 40% - Contenu et mots-clés
        "crawl":     0.2,  // 20% - Accessibilité et structure
    }
    
    return weights["technical"]*techScore + 
           weights["semantic"]*semanticScore + 
           weights["crawl"]*crawlScore
}
```

### Seuils de couleur
- 🔴 0-40: Score critique (rouge)
- 🟠 41-70: Score moyen (orange)  
- 🟡 71-85: Score bon (jaune)
- 🟢 86-100: Score excellent (vert)

### Fallbacks par erreur
- **Crawler failed**: Rapport "Site inaccessible"
- **Technical failed**: Analyse basique (titre/meta uniquement)
- **Semantic failed**: Mots-clés extraits par regex simple
- **Template error**: Rapport JSON en fallback

## Formats d'export spécialisés

### Executive Summary (1 page)
- Score global et tendance
- Top 5 recommandations
- Graphique radar des catégories
- Pas de détails techniques

### Technical Deep Dive
- Tous les findings avec evidence
- Code source des erreurs
- Recommandations de correction
- Checklist d'actions

### Competitive Analysis (futur)
- Comparaison multi-sites
- Benchmarks sectoriels
- Écarts concurrentiels

## Livrables Sprint 1.5

### INT-004: Tests et implémentation
- ✅ Template HTML responsive avec branding SEPTEO
- ✅ Export JSON structuré avec schema
- ✅ Export CSV pour analyse données
- ✅ Scoring consolidé avec poids configurables
- ✅ Fallbacks gracieux pour chaque type d'erreur
- ✅ Tests TDD pour chaque format
- ✅ Validation schema JSON Schema
- ✅ Performance: génération < 500ms pour 100 pages