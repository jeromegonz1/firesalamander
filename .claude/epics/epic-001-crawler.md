# Epic 001 - Agent Crawler

## Vue d'ensemble
Crawler intelligent pour l'audit SEO avec respect des robots.txt et gestion des erreurs robuste.

## Statut
**IMPLEMENTÉ** ✅ (Tests: 6/6 passants)

## Fonctionnalités Core
- [x] Crawling respectueux des robots.txt
- [x] Gestion de la profondeur et limites
- [x] Extraction du contenu HTML structuré
- [x] Détection automatique de la langue
- [x] Normalisation et déduplication des URLs
- [x] Gestion des timeouts et retry logic

## Architecture Technique
- **Package**: `internal/crawler`
- **Point d'entrée**: `crawler.go:NewCrawler()`
- **Configuration**: `crawler.yaml`
- **Tests**: `crawler_test.go` (6 tests)
- **Dépendances**: net/http, golang.org/x/net/html

## Contrats API
```json
{
  "crawl_request": {
    "audit_id": "string",
    "base_url": "string", 
    "max_depth": "number",
    "max_pages": "number"
  },
  "crawl_response": {
    "pages": "array",
    "statistics": "object",
    "errors": "array"
  }
}
```

## Points de Performance
- Crawling concurrent avec pools de workers
- Respect du crawl-delay des robots.txt
- Cache intelligent des pages déjà visitées
- Limitation mémoire avec streaming

## 🧪 Critères BDD
**Given** un site web avec robots.txt et sitemap.xml
**When** le crawler analyse le site
**Then** il respecte les règles robots.txt et explore selon sitemap

## 📊 Estimation
- **Story Points**: 21 pts (Sprint 1)
- **Complexité**: Moyenne
- **Risques**: Gestion timeouts réseau

## 🎯 User Stories
- **US-001**: Parser robots.txt complet (5 pts)
- **US-002**: Support sitemap.xml (8 pts)  
- **US-003**: Gestion redirections 3xx (5 pts)
- **US-004**: Optimisation performance (3 pts)

## Issues Connues
- ⚠️ Gestion des redirections 3xx à améliorer
- ⚠️ Support des sitemaps XML partiel

## Métriques Qualité
- Coverage: 85%+
- Performance: <2s pour 50 pages
- Mémoire: <100MB pour 1000 pages