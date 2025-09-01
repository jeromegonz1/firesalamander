# 🚨 VIOLATIONS CRITIQUES - NO HARDCODING POLICY

## DATE: 2025-08-07
## ARCHITECTE PRINCIPAL: Claude Code
## STATUT: ❌ NON CONFORME

---

## ⚠️ AVERTISSEMENT À TOUS LES AGENTS

**RÈGLE ABSOLUE : AUCUN HARDCODING TOLÉRÉ**

Toute validation de code est **SUSPENDUE** jusqu'à correction complète des violations détectées.

### VIOLATIONS CRITIQUES DÉTECTÉES

#### 1. `/internal/api/simulator.go` - ❌ CRITIQUE
- Durées hardcodées : 500ms, 800ms, 600ms, 300ms
- Pourcentages hardcodés : 0-30%, 30-70%, 70-95%, 95-100%
- Facteurs hardcodés : 1.5, 40.0, 0.1
- Messages hardcodés : "Découverte des pages...", "Analyse SEO..."
- Plages de pages hardcodées : 50-150, 100-300, 200-700

#### 2. `/internal/api/handlers.go` - ❌ CRITIQUE
- Messages d'erreur hardcodés
- Données de test hardcodées (score: 72, pages: 47)
- URLs de test hardcodées

#### 3. `/cmd/server/main.go` - ⚠️ MAJEUR
- Chemins hardcodés : "./templates"
- Données de simulation hardcodées
- Valeurs par défaut hardcodées

### ACTIONS IMMÉDIATES REQUISES

1. **STOP** - Aucune nouvelle fonctionnalité avant correction
2. **REVIEW** - Audit complet de tout le code
3. **REFACTOR** - Externalisation de TOUTES les valeurs
4. **TEST** - Ajout de tests anti-hardcoding

---

**SANCTION : Code rejeté en review si hardcoding détecté**