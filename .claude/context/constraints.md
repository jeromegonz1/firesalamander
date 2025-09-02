# Contraintes Techniques Fire Salamander

## Limites Crawler
- Max URLs : 300 par audit
- Max profondeur : 3 niveaux
- Timeout requête : 10 secondes
- Retry attempts : 2 maximum
- Cache TTL : 3600 secondes
- Concurrent requests : 5

## Limites Analyse Sémantique
- Langue : FR uniquement (confiance ≥ 0.8)
- N-grammes : 2-5 mots
- Max candidats : 500
- Top suggestions : 50 par défaut

## Contraintes Performance
- Crawl : < 2s par page (p95)
- Analyse sémantique : < 15s pour 50 pages
- Génération PDF : < 15s
- Mémoire max : 2GB par audit

## Contraintes Qualité
- Coverage tests : ≥ 85%
- Pas d'API externe payante
- Logs structurés JSON
- Validation schémas obligatoire

## Sécurité
- Respect strict robots.txt
- Pas de données personnelles
- Purge après 30 jours
- Secrets en .env uniquement