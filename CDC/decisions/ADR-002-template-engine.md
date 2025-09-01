# ADR-002: Choix du moteur de templates

## Status
Accepté

## Context
Les templates HTML/CSS existent déjà dans le projet avec Alpine.js et des styles SEPTEO.

## Decision
Conserver et améliorer les **templates Go natifs** existants plutôt que d'adopter un framework lourd.

## Reasoning
- **Capitalisation** : travail déjà effectué sur `templates/home.html`, `analyzing.html`, `results.html`
- **Performance** : templates Go compilés plus rapides que les SPA
- **Simplicité** : pas de build step frontend complexe
- **Branding SEPTEO** : styles et couleurs déjà intégrés
- **Alpine.js** : réactivité suffisante pour l'interface

## Consequences
✅ Préservation du travail existant
✅ Performance optimale pour le rendu
✅ Maintenance simplifiée
✅ Conformité au branding SEPTEO
⚠️ Interactivité limitée vs frameworks modernes