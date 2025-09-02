# Pipeline d'intégration Fire Salamander

## Flux de données complet
1. Input utilisateur → Orchestrator (génère audit_id: FS-XXX)
2. Orchestrator → Crawler (async)
3. Crawler → crawl_index.json → Queue
4. Queue → Technical Analyzer (parallel)
5. Queue → Semantic Analyzer (après crawl)
6. Analyzers → results.json → Report Engine
7. Report Engine → PDF/HTML → User

## Matrice de dépendances
| Agent | Dépend de | Input | Output | Bloque | Timeout |
|-------|-----------|-------|--------|---------|---------|
| Orchestrator | User input | URL + options | audit_id | Tous | 10s |
| Crawler | Orchestrator | audit_request.json | crawl_index.json | Technical, Semantic | 5min |
| Technical | Crawler | crawl_index.json | tech_result.json | Report | 10min |
| Semantic | Crawler | crawl_index.json | semantic_result.json | Report | 5min |
| Report | Tech + Semantic | *_result.json | PDF/HTML | - | 30s |

## Communication inter-agents (JSON-RPC 2.0)

### Message: Orchestrator → Crawler
```json
{
  "jsonrpc": "2.0",
  "method": "start_crawl",
  "params": {
    "audit_id": "FS-001",
    "seed_url": "https://camping-test.fr",
    "max_urls": 300,
    "max_depth": 3,
    "config_file": "config/crawler.yaml"
  },
  "id": "orch-crawler-001"
}
```

### Réponse: Crawler → Orchestrator
```json
{
  "jsonrpc": "2.0",
  "result": {
    "audit_id": "FS-001",
    "status": "complete",
    "pages_crawled": 47,
    "duration_ms": 85420,
    "output_file": "/audits/FS-001/crawl_index.json"
  },
  "id": "orch-crawler-001"
}
```

### Notification d'erreur
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32001,
    "message": "Crawler timeout",
    "data": {
      "audit_id": "FS-001",
      "pages_crawled": 23,
      "partial_output": "/audits/FS-001/crawl_partial.json"
    }
  },
  "id": "orch-crawler-001"
}
```

## États d'audit
```
pending → crawling → analyzing_tech → analyzing_sem → reporting → complete
            ↓            ↓              ↓            ↓
          failed      partial       skipped      error
```