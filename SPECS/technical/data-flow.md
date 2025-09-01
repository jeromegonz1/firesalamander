# Fire Salamander - Flux de données

## Pipeline principal

```
[URL Input] → [Orchestrateur] → [Crawler] → [crawl_index.json]
                    ↓
            [Audit Technique] → [tech_result.json]
                    ↓
            [Analyse Sémantique] → [semantic_result.json]
                    ↓
            [Agent Reporting] → [PDF + Web Report]
```

## Formats de données

### 1. crawl_index.json
Sortie du crawler, entrée pour les analyses :
```json
{
  "pages": [
    {
      "url": "https://example.com/",
      "lang": "fr",
      "title": "Accueil - Mon Site",
      "h1": "Bienvenue",
      "content": "Texte principal...",
      "depth": 0,
      "anchors": [{"text": "Services", "href": "/services"}]
    }
  ],
  "metadata": {
    "total_pages": 150,
    "duration_ms": 45000
  }
}
```

### 2. tech_result.json
Résultats de l'audit technique :
```json
{
  "audit_id": "audit_123",
  "scores": {
    "performance": 0.78,
    "seo": 0.91
  },
  "findings": [
    {
      "id": "missing-title",
      "severity": "critical",
      "message": "Titre manquant",
      "evidence": ["https://example.com/contact"]
    }
  ],
  "mesh": {
    "orphans": [],
    "weak_anchors": ["cliquez ici"]
  }
}
```

### 3. semantic_result.json
Suggestions sémantiques :
```json
{
  "audit_id": "audit_123",
  "topics": [
    {
      "id": "t1",
      "label": "camping nature",
      "terms": ["mobil-home", "emplacement"]
    }
  ],
  "suggestions": [
    {
      "keyword": "camping familial nature",
      "confidence": 0.87,
      "reason": "topic + intent",
      "evidence": ["/services", "/emplacements"]
    }
  ]
}
```

## Communication inter-agents

### JSON-RPC 2.0
Tous les agents communiquent via JSON-RPC :

```json
{
  "jsonrpc": "2.0",
  "method": "crawl",
  "id": "audit_123",
  "params": {
    "seed_url": "https://example.com",
    "max_urls": 300
  }
}
```

### Streaming
Le crawler peut envoyer des mises à jour en temps réel :
```json
{
  "jsonrpc": "2.0",
  "id": "audit_123",
  "result": {
    "status": "page_found",
    "progress": "15/300",
    "current_url": "https://example.com/services"
  }
}
```

## Persistance

### SQLite (phase MVP)
- `audits` - Métadonnées des audits
- `feedback` - Évaluations des suggestions
- `cache` - Cache des embeddings

### Structure audit
```sql
CREATE TABLE audits (
  id TEXT PRIMARY KEY,
  url TEXT NOT NULL,
  status TEXT DEFAULT 'pending',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  completed_at DATETIME,
  config_snapshot TEXT -- JSON des paramètres utilisés
);
```

## Gestion d'erreurs

### Codes d'erreur standardisés
- `CRAWLER_001` - Robots.txt inaccessible
- `CRAWLER_002` - Trop de redirections
- `TECH_001` - Lighthouse timeout
- `SEMANTIC_001` - Langue non supportée

### Fallback
- Si un agent échoue, l'audit continue avec les données partielles
- Les erreurs sont reportées dans le rapport final
- Cache des résultats partiels pour reprise