# Fire Salamander - Contexte Projet

## Vision
Outil d'audit SEO automatisé pour SEPTEO remplaçant Screaming Frog
- Focus marché français
- Pipeline modulaire 5 agents
- Analyse sémantique avec ML

## Repository
https://github.com/jeromegonz1/firesalamander

## Architecture
- Crawler (Go) → collecte pages
- Technical Analyzer → Lighthouse + SEO rules
- Semantic Analyzer → CamemBERT + n-grammes FR
- Report Engine → PDF/HTML
- Orchestrator → JSON-RPC streaming

## Stack
- Backend: Go
- ML: Python (spaCy, CamemBERT)
- LLM: Mistral 7B local (optionnel)
- Frontend: Templates HTML/CSS existants

## Documentation
- CDC: CDC/v4.1-current.md
- Specs: SPECS/functional/full-specifications.md

## Philosophie de Développement

Fire Salamander suit une approche stricte de séparation des rôles :
- Voir `.claude/standards/development-roles-philosophy.md`
- Chaque feature passe par 5 validations (Quality Gates)
- Aucun code n'est mergé sans passer tous les checkpoints
- Cette approche a permis de passer de 4/10 à 8.5/10 en qualité