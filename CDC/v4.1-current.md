# Fire Salamander — Cahier des Charges Fonctionnel & Technique (V4.1)

## 1. Contexte & vision produit

SEPTEO Digital Services souhaite disposer d’un outil interne d’audit SEO automatisé afin de remplacer des solutions payantes (Screaming Frog, etc.) et d’accélérer les audits. **Fire Salamander** se doit d’être :

- **Rapide et exhaustif** : capable de crawler un site complet en respectant le robots.txt, d’identifier les problèmes techniques et sémantiques.
- **Intelligent** : comprendre le métier du client et suggérer des expressions SEO pertinentes.
- **Actionnable** : générer des rapports PDF/web avec des recommandations claires.

## 2. Objectifs fonctionnels (MVP)

1. **Crawl intelligent** : sitemap + maillage interne + règles anti-boucle.
2. **Audit technique** : erreurs SEO classiques + score Lighthouse + analyse du maillage.
3. **Analyse sémantique** (FR uniquement) : compréhension métier + suggestions ≥ 2 mots.
4. **Rapport automatisé** : PDF + interface web claire.
5. **Historique et versioning** : suivi des audits et évolution des suggestions.

## 3. Architecture & découpage par agents

```
 Crawler → Queue → Technical Analyzer
                     ↘ Semantic Analyzer
                               ↓
                        Report Engine
```

### Technologies
- **Back-end** : Go (ou Rust si requis pour perf réseau)
- **Frontend** : HTML léger (templates) + CSS (custom SEPTEO)
- **ML/NLP** : Python (CamemBERT, spaCy, scikit-learn)
- **API inter-agents** : JSON-RPC streaming
- **CI/CD** : GitHub Actions, tests contractuels

## 4. Analyse sémantique (hybride rules + ML)

### 4.1 Pipeline général

1. **Pré-traitement** : lemmatisation FR, nettoyage, détection langue
2. **Extraction candidats** : n-grammes (2–5 mots), TF-IDF, ancres internes
3. **Compréhension métier** : clustering via embeddings (CamemBERT/DistilCamemBERT), intent classifier
4. **Ranking** : par pertinence, intent, maillage interne, lisibilité FR
5. **Retour consultant** : 👍 / 👎 avec raison (trop générique, duplicata, etc.)
6. **Apprentissage** : régression logistique ou XGBoost entraîné périodiquement (versionné)

### 4.2 Suggestions enrichies par IA (LLM optionnel)

- Utilisation possible de GPT-3.5 ou Mistral 7B pour générer des expressions métier longue traîne
- Prompts auto-générés à partir des topics et n-grammes détectés
- Cache des suggestions par hash contenu + fallback local
- Limité à certains cas premium ou rapide via mode "deep"

### 4.3 Recommandations performances

```yaml
crawler:
  concurrent_requests: 5
  request_timeout: 10s
  retry_attempts: 2
  cache_ttl: 3600s
  respect_crawl_delay: true
```

## 5. Feedback structuré (JSON)

```json
{
  "feedback_type": "keyword_relevance",
  "keyword": "camping familial",
  "rating": -1,
  "reason": "too_generic",
  "context": {
    "page_url": "...",
    "intent": "commercial"
  }
}
```

## 6. Métriques d’évaluation sémantique

- Precision@10 ≥ 0.6
- nDCG@10 ≥ 0.7
- Ratio expressions ≥ 2 mots ≥ 80 %
- Taux d’intégration des suggestions dans les audits suivants
- Temps moyen génération audit < 2 min

## 7. Roadmap réaliste

### Phase 0 (Semaine 1) — Fondations
- CI/CD GitHub Actions
- Contrats JSON-RPC
- Dataset d’évaluation (100 sites FR annotés)

### Phase 1 (Semaines 2–3) — Cœur technique
- Crawler (sitemap + interne) avec limites réglables
- Analyse Lighthouse (accessibilité, SEO, perf)
- Génération rapport PDF simple

### Phase 2 (Semaines 4–5) — Sémantique baseline
- N-grammes + TF-IDF + règles POS
- Intent classifier simple + clustering light
- Feedback utilisateur intégré

### Phase 3 (Semaines 6–7) — ML & IA enrichie
- Embeddings + ranker ML
- Cache + IA LLM locale (Mistral 7B)
- Suggestions IA + mode "fast" / "deep"

### Phase 4+ — Optimisation continue
- A/B testing des suggestions
- Banlist secteur, analyse cross-audit
- Visualisation score + monitoring métier

## 8. Monitoring & KPIs clés

- % suggestions intégrées dans le site
- Top erreurs SEO techniques récurrentes
- Score Lighthouse moyen par type de page
- Taux d’adoption suggestions
- Temps moyen d’audit

## 9. Sécurité & RGPD

- Pas de données personnelles collectées
- Tous les textes sont publics (pages web)
- Suggestions LLM anonymisées et non stockées

## 10. Annexes

- `contracts/semantic_request.schema.json`
- `contracts/semantic_response.schema.json`
- `config/semantic.yaml`
- `tests/semantic/test_ranker.py`
- `ccpm/context.md`, `memory.md`, `epics/semantic-engine.md`

---

**Fire Salamander** est conçu comme une boîte à outils SEO modulaire, rapide et intelligente, capable de comprendre un site en français, d’en extraire les enjeux métiers et d’en tirer des recommandations concrètes, mesurables, et surtout **améliorables dans le temps**.

