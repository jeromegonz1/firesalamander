# ADR-003: Architecture Machine Learning

## Status
Accepté

## Context
L'analyse sémantique nécessite des modèles NLP français et potentiellement de l'IA générative pour les suggestions.

## Decision
Architecture **hybride** : règles heuristiques + ML léger + IA optionnelle.

## Reasoning
- **Baseline robuste** : TF-IDF + n-grammes fonctionnent sans dépendances lourdes
- **ML ciblé** : CamemBERT/DistilCamemBERT pour embeddings français
- **IA premium** : Mistral 7B local ou GPT-3.5 en option
- **Évolutivité** : possibilité d'améliorer le ML sans réécriture

## Architecture retenue

```
Texte → [Règles] → [ML Embeddings] → [Ranker] → [IA optionnelle] → Suggestions
```

### Composants
1. **Règles heuristiques** : n-grammes, TF-IDF, filtres qualité
2. **ML** : embeddings pour clustering thématique et similarité
3. **Ranker** : modèle léger (régression logistique) entrainable
4. **IA** : LLM pour génération créative (mode premium)

## Consequences
✅ Performance dégradée gracieusement si ML indisponible
✅ Amélioration continue via feedback
✅ Coûts maîtrisés (IA optionnelle)
✅ Résultats français optimisés
⚠️ Complexité technique accrue
⚠️ Dépendances Python pour le ML