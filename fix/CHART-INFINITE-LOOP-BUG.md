# 🚨 PROBLÈME CRITIQUE: Boucle Infinie Chart.js

**Date:** 04/08/2025
**Priorité:** 🔴 CRITIQUE
**Status:** En cours de résolution

## 📋 Description du Problème

Fire Salamander subit une **boucle infinie de redimensionnement** causée par Chart.js, résultant en une hauteur de page qui croît indéfiniment (observé jusqu'à **90+ millions de pixels**).

### 🔍 Symptômes Observés

- Page height qui grandit continuellement (519136.col1644+ pixels)
- Interface complètement inutilisable sur mobile et desktop
- Navigateur devient très lent/non-responsive
- DevTools Chrome montre une hauteur anormalement élevée

### 🏷️ Root Cause Identifiée

**Chart.js avec `maintainAspectRatio: false`** dans `/web/static/app.js`:

```javascript
// PROBLÉMATIQUE - Cause la boucle infinie
this.charts.scoresChart = new Chart(ctx, {
    options: {
        responsive: true,
        maintainAspectRatio: false,  // ⚠️ BOUCLE INFINIE
        // ...
    }
});
```

## 🧪 Tests Effectués

### ✅ Tests de Diagnostic
1. **Redémarrage serveur** - ❌ Inefficace (fichiers embedded)
2. **Modifications CSS** - ❌ Inefficace (problème JS)
3. **Désactivation charts** - ✅ Confirmé comme cause root
4. **Monkey-patching Chart.js** - ⚠️ Partiellement efficace

### 📊 Résultats des Tests
- **Sans charts:** Hauteur normale (~2000px)
- **Avec Chart.js maintainAspectRatio: false:** Boucle infinie (90M+ pixels)
- **Avec fix CSS seul:** Inefficace

## 🔧 Solutions Testées

### 1. CSS Constraints (❌ Échec)
```css
.chart-container {
    height: 350px;
    max-height: 350px;
    overflow: hidden;
}
```
**Résultat:** Inefficace - la boucle persiste au niveau JS

### 2. JavaScript Monkey-Patch (⚠️ Partiel)
```javascript
// Override Chart constructor
window.Chart = function(ctx, config) {
    if (config.options.maintainAspectRatio === false) {
        config.options.maintainAspectRatio = true;
        config.options.aspectRatio = 2;
    }
    return new originalChart(ctx, config);
};
```
**Résultat:** Fonctionnel mais nécessite injection dynamique

### 3. Désactivation Temporaire (✅ Succès)
```javascript
// Désactiver charts pour debug
// this.updateDashboardCharts();
console.log('⚠️ Charts désactivés pour debug de la boucle infinie');
```
**Résultat:** Hauteur normale confirmée

## 🎯 Solution Recommandée

### Option A: Remplacement CSS-Only Charts
Remplacer Chart.js par des graphiques CSS purs pour éviter complètement le problème.

### Option B: Recompilation Fire Salamander
Modifier les fichiers embedded et recompiler le binaire.

### Option C: Charts CSS + SVG Personnalisés
Créer des visualisations légères sans dépendances externes.

## 📁 Fichiers Impactés

```
/web/static/app.js              - Code Chart.js problématique  
/web/static/styles.css          - Tentatives CSS fixes
/web/static/index.html          - Canvas elements
/fix/chart-infinite-loop.html   - Interface de test
```

## 🚀 Plan de Résolution

1. **Phase 1:** Créer charts CSS-only de remplacement
2. **Phase 2:** Tests incrémentaux avec logging
3. **Phase 3:** Déploiement avec fallback
4. **Phase 4:** Monitoring pour récidive

## 📝 Notes Techniques

- Fire Salamander a les fichiers statiques **embedded dans le binaire**
- Modifications `/web/static/` ne sont PAS prises en compte sans recompilation
- Le serveur sert les fichiers depuis le binaire, pas le filesystem
- Problème reproductible sur Chrome, Safari, Firefox

## ⚠️ Impact Business

- **Interface inutilisable** - Bloque tous les utilisateurs
- **Tests UX/UI impactés** - BackstopJS détecte le problème
- **Productivité équipe** - Développement bloqué

---
**Créé par:** Claude Code Assistant  
**Dernière MAJ:** 04/08/2025 11:54