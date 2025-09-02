# Scénarios BDD - Fire Salamander Sprint 1.5

## Format des audit_ids
- Pattern: `FS-{ENV}-{SEQ}` où ENV = PROD/TEST/DEV, SEQ = 001-999
- Exemples: `FS-PROD-001`, `FS-TEST-042`, `FS-DEV-123`

## Feature: Pipeline d'audit complet

### Scenario: Audit réussi avec toutes les étapes
```gherkin
Feature: Pipeline audit Fire Salamander
  En tant qu'analyste SEO
  Je veux exécuter un audit complet
  Pour obtenir un rapport d'analyse technique et sémantique

Scenario: Audit complet réussi
  Given que j'ai une URL valide "https://camping-bretagne.fr"
  And que le système est opérationnel
  When je lance l'audit avec audit_id "FS-PROD-001"
  Then le crawler explore le site en ≤ 5 minutes
  And l'analyse technique valide 5 critères SEO
  And l'analyse sémantique extrait ≥ 10 mots-clés
  And le rapport final est généré en HTML/JSON
  And le statut final est "completed"
  And la progression atteint 100%

Scenario: Audit avec erreur de crawling
  Given que j'ai une URL invalide "https://site-inexistant.test"
  When je lance l'audit avec audit_id "FS-TEST-001"
  Then le crawler échoue en ≤ 30 secondes
  And le statut devient "failed"
  And l'erreur est capturée dans les logs
  And aucun rapport n'est généré

Scenario: Audit partiel (crawler OK, semantic KO)
  Given que j'ai une URL valide "https://site-test.fr"
  And que le service sémantique est indisponible
  When je lance l'audit avec audit_id "FS-DEV-001"
  Then le crawler termine avec succès
  And l'analyse technique produit des résultats
  And l'analyse sémantique échoue gracieusement
  And un rapport partiel est généré
  And le statut final est "partial"
```

### Scenario: Gestion des timeouts
```gherkin
Scenario: Timeout du crawler
  Given que j'ai une URL lente "https://site-tres-lent.fr"
  And un timeout crawler de 30 secondes
  When je lance l'audit avec audit_id "FS-PROD-002"
  Then le crawler s'arrête au timeout
  And les pages déjà crawlées sont sauvegardées
  And le statut devient "timeout"
  And un rapport partiel est généré

Scenario: Timeout de l'analyse sémantique
  Given que j'ai des données crawl volumineuses (500 pages)
  And un timeout sémantique de 2 minutes
  When l'analyse sémantique dépasse le timeout
  Then l'analyse s'arrête proprement
  And les résultats partiels sont conservés
  And le rapport indique "analyse incomplète"
```

### Scenario: Audit_id tracking
```gherkin
Scenario: Traçabilité complète de l'audit_id
  Given que je lance un audit avec "FS-PROD-003"
  When l'audit progresse dans le pipeline
  Then chaque agent loggue l'audit_id
  And tous les fichiers contiennent l'audit_id
  And le rapport final référence l'audit_id
  And l'historique est traçable par audit_id

Scenario: Audit_id unique et collisions
  Given qu'un audit "FS-PROD-004" est en cours
  When je tente de lancer un autre audit "FS-PROD-004"
  Then le système rejette la requête
  And retourne une erreur "audit_id already exists"
  And propose un nouvel audit_id "FS-PROD-005"
```

## Feature: Intégration inter-agents

### Scenario: Communication JSON-RPC
```gherkin
Scenario: Messages JSON-RPC valides
  Given que l'orchestrator veut contacter le crawler
  When il envoie un message JSON-RPC 2.0
  Then le message contient: jsonrpc, method, params, id
  And le crawler répond avec: jsonrpc, result/error, id
  And l'id de réponse = id de requête
  And le format JSON est strictement respecté

Scenario: Gestion des erreurs JSON-RPC
  Given qu'un agent renvoie une erreur
  When l'orchestrator reçoit la réponse
  Then la réponse contient un objet "error"
  And l'erreur a un code, message, et data
  And l'audit_id est dans les données d'erreur
  And l'orchestrator loggue l'erreur complète
```

## Feature: Formats de données

### Scenario: Structure crawl_index.json
```gherkin
Scenario: Validation du format crawl_index.json
  Given qu'un crawl est terminé pour "FS-PROD-005"
  When je lis le fichier crawl_index.json
  Then il contient: audit_id, pages[], metadata
  And chaque page a: url, title, content, status_code
  And metadata contient: crawl_date, duration, total_pages
  And le JSON respecte le schema défini

Scenario: Structure tech_result.json
  Given qu'une analyse technique est terminée
  When je lis tech_result.json
  Then il contient: audit_id, findings[], summary
  And chaque finding a: id, severity, message, evidence
  And summary contient: total_issues, score, recommendations

Scenario: Structure semantic_result.json
  Given qu'une analyse sémantique est terminée
  When je lis semantic_result.json
  Then il contient: audit_id, keywords[], topics[], suggestions[]
  And chaque keyword a: text, score, frequency, contexts
  And le format respecte le schema sémantique
```

## Critères d'acceptation Sprint 1.5

### INT-002: Tous les scénarios BDD
- ✅ 12 scénarios définis avec audit_ids
- ✅ Coverage: pipeline, erreurs, timeouts, formats
- ✅ Given/When/Then structure respectée
- ✅ Critères d'acceptation mesurables

### Gherkin Keywords utilisés
- **Given**: État initial du système
- **When**: Action déclenchée 
- **Then**: Résultat attendu observable
- **And**: Conditions/résultats additionnels