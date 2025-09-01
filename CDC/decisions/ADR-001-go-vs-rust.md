# ADR-001: Choix du langage Go vs Rust

## Status
Accepté

## Context
Fire Salamander nécessite un langage performant pour le crawling concurrent et l'analyse de données. Les options principales étaient Go et Rust.

## Decision
Utilisation de **Go** comme langage principal.

## Reasoning
- **Simplicité** : syntaxe accessible pour l'équipe SEPTEO
- **Concurrence native** : goroutines idéales pour le crawling parallèle
- **Écosystème** : bibliothèques matures pour HTTP, HTML parsing, JSON
- **Productivité** : développement plus rapide vs Rust
- **Maintenance** : courbe d'apprentissage plus douce

## Consequences
✅ Développement rapide et maintenable
✅ Concurrence native pour le crawler
✅ Intégration facile avec les templates HTML
⚠️ Performance légèrement inférieure à Rust pour le parsing intensif