# Philosophie des Rôles de Développement - Fire Salamander

## Principe Fondamental : La Séparation des Responsabilités

Ce document explique pourquoi Fire Salamander exige le respect strict de différents rôles durant le développement, même en travail solo.

## Le Danger du Développement "Bulldozer"

### Sans rôles séparés
Développeur écrit 500 lignes → "Ça compile" → Code suivant →
→ 1 mois plus tard : 5000 lignes non testées, 148 hardcodings, build cassé

### Avec rôles séparés
[Dev] 100 lignes → [QA] "Tests?" → [Dev] Tests → [Reviewer] "Hardcoding!" →
[Architecte] "Interface?" → 100 lignes SOLIDES

## Les 5 Rôles Obligatoires

### 1. 👨‍💻 Développeur
**Focus** : Faire fonctionner le code
**Produit** : Code qui compile
**Piège évité** : Code qui marche mais inmaintenable

### 2. 🧪 QA/Testeur  
**Focus** : Casser le code
**Produit** : Tests exhaustifs, coverage >85%
**Piège évité** : "Ça marche sur ma machine"

### 3. 🔍 Reviewer
**Focus** : Qualité et standards
**Produit** : Code clean, pas de dette technique
**Piège évité** : Accumulation silencieuse de problèmes

### 4. 🏗️ Architecte
**Focus** : Vision long terme
**Produit** : Interfaces SOLID, modularité
**Piège évité** : Architecture spaghetti après 6 mois

### 5. 📝 Documentaliste
**Focus** : Clarté et traçabilité
**Produit** : Code compréhensible par d'autres
**Piège évité** : "Seul Jean comprend ce code"

## Points de Contrôle Obligatoires

Chaque rôle crée un "checkpoint" qu'on ne peut pas ignorer :

| Étape | Rôle | Validation requise | Impossible de continuer si... |
|-------|------|-------------------|-------------------------------|
| 1 | Dev | Code compile | Erreurs de compilation |
| 2 | QA | Tests passent | Coverage < 85% |
| 3 | Reviewer | Standards respectés | Hardcoding détecté |
| 4 | Architecte | SOLID respecté | Interfaces manquantes |
| 5 | Doc | Documentation à jour | README obsolète |

## Exemple Concret : Sprint de Correction

**Sans les rôles**, on n'aurait jamais découvert :
- 148 valeurs hardcodées
- Interfaces SOLID manquantes  
- Coverage à 62% seulement
- Build cassé avec panics

**Avec les rôles**, Fire Salamander est passé de 4/10 à 8.5/10.

## Pourquoi c'est CRITIQUE pour l'IA

Claude Code peut générer 5000 lignes en 10 minutes. Sans ces rôles :
- Accumulation massive de dette technique
- Code non testé et non validé
- Impossible à maintenir après 1 semaine

Les rôles forcent des pauses réflexives et des validations.

## Application Pratique

### Workflow obligatoire pour chaque feature
```bash
# 1. Développeur
git checkout -b feature/X
# Écrire code minimal

# 2. QA
# STOP - Écrire tests d'abord
# Vérifier coverage

# 3. Reviewer  
# STOP - Vérifier standards
# Éliminer hardcoding

# 4. Architecte
# STOP - Vérifier architecture
# Ajouter interfaces si nécessaire

# 5. Documentaliste
# STOP - Mettre à jour docs
# Committer seulement après
```

## Règles Non-Négociables

1. Jamais plus de 100 lignes sans tests
2. Jamais de merge sans review (même en solo)
3. Jamais de feature sans documentation
4. Jamais de hardcoding (tout en config)
5. Jamais ignorer l'architecture pour aller vite

## Métriques de Succès

- Temps de développement initial : +20%
- Temps de maintenance : -80%
- Bugs en production : -90%
- Dette technique : Proche de zéro
- Note de code : >8/10 constant

## Conclusion

Les rôles ne sont pas une contrainte bureaucratique mais une protection contre notre propre enthousiasme. Ils garantissent que Fire Salamander reste maintenable, testable et professionnel.

Cette philosophie est obligatoire, non négociable, et doit être respectée même sous pression de deadline.