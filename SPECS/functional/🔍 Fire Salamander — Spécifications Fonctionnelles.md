🔍 Fire Salamander — Spécifications Fonctionnelles Détaillées (V1.1)

Ce document structure en détail les spécifications fonctionnelles du projet Fire Salamander. Il s'appuie sur le CDC V4.1 et tient compte des retours (phases, performance, dataset, IA optionnelle). Son objectif est de guider Claude Code dans la mise en œuvre des agents techniques, des interfaces, des flux de données et du cycle de vie d'un audit SEO automatisé.

1. Parcours utilisateur (SEPTEO consultant SEO)
1.1 Cas d’usage type

Objectif : réaliser un audit SEO automatisé pour un site client SEPTEO (camping, hôtel, etc.)

🧭 Étapes UX principales :

Connexion à l’interface web interne SEPTEO (aucune authentification nécessaire pour le MVP).

Saisie d’un domaine à auditer : exemple https://www.monsite.com.

Choix des options :

Langue : FR (pré‑sélectionné).

Profondeur max : 3.

Mode analyse sémantique : Rapide (rules-only) ou Approfondie (rules + embeddings + IA).

Nombre de suggestions souhaitées : 25, 50, 100.

Lancement de l’audit (bouton « Analyser ce site ») → déclenche un appel RPC vers le backend (Agent Orchestrateur).

Suivi en temps réel : affichage en streaming des étapes (Crawler en cours … | Audit technique OK | Analyse sémantique OK | Génération du rapport PDF).

Affichage des résultats dans l’UI (onglets) :

Vue synthétique : score global, nombre d’erreurs, suggestions clés.

Vue technique : tableau Lighthouse (perf/SEO/accessibilité/best practices), liste d’erreurs et de warnings, analyse du maillage.

Vue sémantique : expressions proposées, raisons et evidences.

Téléchargement du PDF ou consultation via l’interface SEPTEO.

Feedback optionnel : 👍 / 👎 sur les suggestions sémantiques, avec motif (trop générique, hors sujet…).

Audit sauvegardé : historique par domaine (audit_id, date, options utilisées, scores, suggestions).

2. Agents fonctionnels à développer
2.1 Agent Crawler

Rôle : visiter intelligemment le site web donné, en respectant robots.txt, sitemap.xml et les liens internes, pour constituer la base de l’audit.

Exploration : algorithme BFS ou priorisation par profondeur sémantique (Home → Catégorie → Produit/Article).

Domaines : ne pas sortir du domaine principal ni suivre les liens externes.

Langue : detection sur chaque page (seules les pages FR sont envoyées au module sémantique).

Paramètres de performance : concurrent_requests=5, request_timeout=10s, retry_attempts=2, cache_ttl=3600s (configurable via YAML).

Sampling (pour sites volumineux) : possibilité de réduire le crawl à un sous‑ensemble représentatif (Home, catégories, top produits, articles récents).

Sortie : JSON par page : URL, lang, title, h1–h3, ancres, texte, canonical, indexation, profondeur, liens entrants/sortants, etc.

2.2 Agent Audit Technique

Rôle : exécuter un audit technique détaillé et générer des scores et des recommandations.

Scores Lighthouse : performance, SEO, accessibilité, best practices (mode desktop).

Analyse SEO : détection des balises manquantes, meta dupliquées ou trop longues, h1 multiples, images sans alt, absence de HTTPS, redirections 302/301 abusives, etc.

Maillage interne : statistiques de profondeur, pages orphelines, ancres pauvres (ex. « cliquez ici »), distribution des ancres par page.

Sortie : tableau de scores + liste d’erreurs et de warnings + structure du maillage (mesh.json).

2.3 Agent Analyse Sémantique (hybride rules + ML)

Rôle : comprendre le contenu métier du site, identifier les thèmes et intentions, et extraire des expressions longue traîne (≥ 2 mots) en français.

Pipeline :

Pré‑traitement : normalisation, lemmatisation, extraction n‑grammes (2–5 mots), filtrage des expressions trop génériques
noiise.com
.

Compréhension métier : embeddings (CamemBERT/DistilCamemBERT), clustering thématique (HDBSCAN/BERTOPIC), intent classifier (info/comm/transac/nav), type de page (via DOM + heuristiques).

Génération candidates : n‑grammes + TF‑IDF + TextRank/YAKE + ancres internes + variantes morphosyntaxiques.

Ranking : score pondéré (similarité topic, intent cible, preuves maillage, lisibilité FR, historique feedback).

Sortie : top‑N suggestions (≥ 2 mots), avec keyword, reason, confidence, evidence_urls.

Modes :

Rapide : pipeline rules + heuristiques (n‑grammes, TF‑IDF, TextRank).

Approfondie : pipeline complet (rules + embeddings + clustering + ML) avec option IA.

IA optionnelle : appel à un modèle de langage (GPT‑3.5 ou Mistral 7B local) pour générer des variantes créatives (mode « deep »). Ce composant est facultatif et activé pour les clients premium. Les prompts sont construits à partir des topics et n‑grammes existants et les résultats sont filtrés pour respecter la langue et la pertinence.

Dataset & feedback : constitution d’un jeu de 100 sites annotés (phase 0) pour mesurer Precision@K et nDCG@K et calibrer le ranker. Feedback structuré pour alimenter l’apprentissage (voir §7).

2.4 Agent Reporting

Rôle : générer un rapport complet et facile à lire (PDF + HTML).

Gabarits : design responsive conforme à la charte SEPTEO.

Sections : Résumé (scores globaux), Technique (détails Lighthouse, erreurs, maillage), Sémantique (topics, suggestions, justifications), Recommandations & Next Steps, Feedback consultant.

Export : PDF via Puppeteer/wkhtmltopdf, avec numérotation unique et horodatage (FSAUDIT‑nnn).

Historique : rapport disponible dans /audits/{audit_id}/report.pdf.

2.5 Agent Orchestrateur

Rôle : contrôler le flux de l’audit et coordonner les agents.

Génère un audit_id et crée un dossier /audits/{audit_id}.

Lance les agents en parallèle ou en pipeline selon la configuration (Crawler → Audit Technique & Sémantique).

Publie des événements streaming (JSON‑RPC) à l’UI pour le suivi en temps réel.

Expose des endpoints API : GET /api/audits/{audit_id}/status, POST /api/audits/{audit_id}/cancel.

Enregistre les erreurs et les durées d’exécution dans un log.

3. Schémas d’échange JSON (contrats)
3.1 semantic_request.schema.json
{
  "audit_id": "string",
  "lang": "fr",
  "max_candidates": 500,
  "top_n": 50
}

3.2 semantic_response.schema.json
{
  "audit_id": "string",
  "topics": [
    { "id": "t1", "label": "camping nature", "terms": ["camping familial", "mobil-home nature"] }
  ],
  "suggestions": [
    {
      "keyword": "camping familial bord de mer",
      "reason": "topic:t1 + anchor keyword + intent=transactional",
      "confidence": 0.87,
      "evidence": ["https://site/page1", "/offres/early-booking"]
    }
  ]
}

3.3 feedback.schema.json

Exemple de structure de feedback envoyé par un consultant :

{
  "feedback_type": "keyword_relevance",
  "keyword": "camping familial",
  "rating": -1,
  "reason": "too_generic",
  "context": {
    "page_url": "https://www.monsite.com/offres",
    "intent": "commercial"
  }
}

4. Limites et performances

Max_urls : 300 par audit (configurable).

Concurrence : 5 requêtes simultanées dans le crawler.

Time-out : 10 s par requête, deux retries en cas d’échec.

Cache TTL : 3600 s (1 h) pour les pages déjà visitées.

Mode sampling : possibilité de réduire l’exploration aux pages prioritaires (home, catégories, top produits) pour les sites e‑commerce.

IA Deep : appels à un LLM local (Mistral 7B) pour limiter les coûts et la latence ; cache les suggestions générées pour réutilisation.

5. CI/CD et tests

TDD obligatoire pour chaque agent, avec tests unitaires et tests de contrat.

Contrats JSON validés via jsonschema et versionnés.

GitHub Actions : pipeline lint → test → build → coverage (≥ 85 %) → security (scan secrets).

Data anonymisation : purge des artefacts d’audit après 30 jours ; pas de données personnelles dans les prompts IA.

6. Roadmap de développement

Phase 0 (Semaine 1) : mise en place de l’infrastructure (CI/CD, dossiers audits/, contrats JSON) et constitution d’un dataset annoté (100 sites FR) pour la sémantique.

Phase 1 (Semaines 2–3) : développement du Crawler et de l’Audit Technique (Lighthouse) + génération d’un rapport PDF simple.

Phase 2 (Semaines 4–5) : implémentation de la sémantique basique (n‑grammes, TF‑IDF, règles POS) + collecte de feedback utilisateur.

Phase 3 (Semaines 6–7) : ajout des embeddings, clustering, intent classifier, ranker ML + possibilité d’utiliser un LLM local pour suggestions longues.

Phase 4+ : optimisation continue (A/B testing, cache intelligent, cross‑audit learning), gestion d’un LLM premium (GPT‑3.5) pour certains clients, ajout de nouveaux connecteurs (Memory Bank, GitHub MCP) si besoin.

7. Dataset d’évaluation et métriques

Dataset : 100 sites FR annotés manuellement (titre, thème, intent, suggestions acceptées/rejetées) pour évaluer la précision des suggestions.

Métriques : Precision@10 ≥ 0.6, nDCG@10 ≥ 0.7, taux d’acceptation consultant ≥ 60 %, ratio expressions ≥ 2 mots ≥ 90 %.

KPIs métier : distribution des scores Lighthouse, top erreurs techniques, temps moyen d’audit, taux d’intégration des suggestions dans les sites, progression des pages indexées.

8. Annexes à venir dans le dépôt

/specs/crawler.md : algorithmes de crawl détaillés, gestion des sitemaps, sampling.

/specs/semantic.md : diagramme complet du pipeline, description des modèles, heuristiques POS, exemples de prompts pour LLM.

/specs/reporting.md : structure du PDF, détails des gabarits, style guide SEPTEO.

/specs/agents.md : protocole JSON‑RPC et orchestration, diagrammes de séquence.

/contracts/*.schema.json : schémas versionnés (audit_request, tech_result, semantic_response, feedback).

/ccpm/context.md, /memory.md, /epics/ : fichiers de mémoire et de gestion de projet (CCPM).

Fire Salamander reste fidèle à sa vision : un outil SEO modulaire, rapide et intelligent, capable de comprendre un site en français, d’en extraire les enjeux métiers et d’en tirer des recommandations concrètes, mesurables et évolutives. Cette V1.1 intègre un découpage progressif, des optimisations de performance, une attention à la gestion des coûts IA et une phase de constitution de dataset pour valider le moteur sémantique.