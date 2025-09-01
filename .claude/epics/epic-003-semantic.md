# Epic 003 - Agent Analyse Sémantique

## Vue d'ensemble
Moteur d'analyse sémantique avancé avec NLP français et suggestions de mots-clés intelligentes.

## Statut
**IMPLEMENTÉ** ✅ (Tests: 15/15 passants)

## Fonctionnalités Core
- [x] Extraction automatique de n-grammes pertinents
- [x] Topic modeling avec clustering K-means
- [x] Scoring sémantique multi-critères (TF-IDF, proximité, densité)
- [x] Suggestions de mots-clés contextualisées
- [x] Analyse de la cohérence thématique
- [x] API REST Flask pour intégration

## Architecture Technique
- **Langage**: Python 3.9+
- **Framework**: Flask (API REST)
- **NLP**: spaCy français, scikit-learn
- **Package**: `internal/semantic/python/`
- **Point d'entrée**: `semantic_analyzer.py:SemanticAnalyzer`
- **Configuration**: `semantic.yaml`
- **Tests**: `test_semantic_analyzer.py` (15 tests)

## Contrats API
```json
{
  "semantic_analysis_request": {
    "audit_id": "string",
    "crawl_data": "object",
    "language": "string"
  },
  "semantic_analysis_response": {
    "keyword_suggestions": "array",
    "topic_analysis": "object",
    "semantic_score": "number",
    "recommendations": "array"
  }
}
```

## Algorithmes Implémentés
- **N-gram extraction**: 1-4 grammes avec filtrage stopwords
- **Topic modeling**: K-means clustering avec 3-8 clusters
- **Scoring**: Combinaison TF-IDF + position + densité + cohérence
- **Diversity filter**: Évite la sur-représentation thématique
- **French NLP**: Tokenisation, lemmatisation, POS tagging

## Points de Performance
- API asynchrone avec Flask
- Cache des modèles NLP
- Traitement par batch pour gros volumes
- Timeout configurable

## Issues Connues
- ⚠️ Modèles CamemBERT non intégrés (phase 1)
- ⚠️ Analyse sentiment basique

## Métriques Qualité
- Coverage: 95%+
- Précision suggestions: >80% (évaluation humaine)
- Performance: <3s pour 100 pages