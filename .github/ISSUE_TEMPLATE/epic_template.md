---
name: 🎯 Epic Template
about: Créer un nouvel epic pour Fire Salamander
title: '[EPIC] '
labels: 'type:epic'
assignees: ''
---

## 🎯 Vision de l'Epic
Description de l'objectif global et de la valeur métier.

## 🏗️ Contexte technique
Agent(s) concerné(s) et architecture impactée.

## 📋 User Stories
Liste des user stories qui composent cet epic:

### US-001: [Titre]
**En tant que** [utilisateur]
**Je veux** [fonctionnalité]  
**Afin de** [bénéfice]

**Critères BDD:**
- **Given** [contexte]
- **When** [action]
- **Then** [résultat]

**Estimation**: [points] points

### US-002: [Titre]
[Répéter pour chaque user story]

## 🧪 Stratégie de test
- **Tests unitaires**: Couverture cible et mocks
- **Tests d'intégration**: Scénarios end-to-end
- **Tests contractuels**: Validation JSON Schema
- **Tests de performance**: Métriques cibles

## 📊 Estimation globale
- **Story Points totaux**: [X] points
- **Durée estimée**: [Y] jours
- **Complexité**: [Simple/Moyen/Complexe]
- **Risques**: [Technique/Fonctionnel/Délais]

## 🔗 Dépendances
- **Prérequis**: Epics/features qui doivent être terminés avant
- **Bloquants**: Issues connues qui pourraient ralentir
- **Librairies**: Nouvelles dépendances Go/Python nécessaires

## 📈 Métriques de succès
- [ ] Tous les tests passent (coverage ≥ 85%)
- [ ] Performance dans contraintes (.claude/context/constraints.md)
- [ ] Documentation complète (ADR + README)
- [ ] Validation utilisateur réussie

## 🎛️ Configuration
Nouveaux paramètres à ajouter dans `config/*.yaml`:
```yaml
# Exemple de nouvelle configuration
agent_name:
  feature_x: true
  threshold: 0.85
```

## 📋 Definition of Done Epic
- [ ] Toutes les user stories terminées
- [ ] Tests d'intégration passent
- [ ] Documentation mise à jour
- [ ] Performance validée
- [ ] Déployement sans régression