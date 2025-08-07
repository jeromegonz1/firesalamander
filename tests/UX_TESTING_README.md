# 🔥 Fire Salamander - Suite de Tests UX/UI

## Vue d'Ensemble

Cette suite de tests UX/UI garantit que Fire Salamander respecte les **standards SEPTEO** et offre une expérience utilisateur optimale. Elle inclut 5 agents automatisés qui valident tous les aspects de l'interface utilisateur.

## 🎯 Objectifs SEPTEO

- **Accessibilité:** ≥95% conformité WCAG 2.1 AA
- **Performance:** ≥90% score Lighthouse 
- **Design System:** ≥90% conformité couleurs/espacement SEPTEO
- **UX:** ≥95% succès parcours utilisateur (zéro friction critique)
- **Régression Visuelle:** 0 régression non approuvée

## 🤖 Agents UX Automatisés

### 1. **Visual Regression Agent** (BackstopJS)
```bash
npm run test:visual
```
- Détection automatique des changements visuels
- Validation du design SEPTEO sur tous les breakpoints
- Captures de référence et comparaisons pixel-perfect
- Tests responsive (mobile/tablet/desktop)

**Configuration:** `tests/ux/visual-regression/backstop.config.js`

### 2. **Accessibility Agent** (Axe-core + Pa11y)  
```bash
npm run test:accessibility
```
- Contraste couleur orange SEPTEO (#ff6136) validé
- Navigation clavier complète
- Labels ARIA obligatoires
- Conformité WCAG 2.1 niveau AA minimum
- Rapport d'accessibilité automatique

**Configuration:** `tests/ux/accessibility/axe.config.js`

### 3. **UX Metrics Agent** (Lighthouse CI)
```bash
npm run test:lighthouse
```
- Core Web Vitals: LCP < 2.5s, FID < 100ms, CLS < 0.1
- Score performance > 90, Accessibilité > 95
- Audit SEO technique automatique
- Budgets de performance configurés

**Configuration:** `tests/lighthouserc.js`

### 4. **Design System Validator Agent**
```bash
npm run test:design-system
```
- Palette SEPTEO exclusive (#ff6136, #e55a2e, #2c3e50...)
- Grille 8px obligatoire pour espacement
- Typographie système uniquement
- Composants réutilisables validés
- Détection styles inline interdits

**Configuration:** `tests/ux/design-system/validator.js`

### 5. **User Flow Testing Agent** (Playwright)
```bash
npm run test:user-flows
```
- Parcours utilisateur critiques (analyse rapide/complète)
- Mesure temps de completion
- Détection points de friction automatique
- Recording sessions problématiques
- Tests multi-navigateurs (Chrome/Firefox/Safari)

**Configuration:** `tests/playwright.config.js`

## 🚀 Installation et Configuration

### Prérequis
```bash
# Node.js 18+
node --version

# Fire Salamander démarré
./fire-salamander --config config.yaml
```

### Installation
```bash
cd tests
npm install
npx playwright install
```

### Variables d'Environnement
```bash
export SEPTEO_STANDARDS=enabled
export FIRE_SALAMANDER_URL=http://localhost:8080
```

## 📊 Commandes Principales

### Tests Individuels
```bash
# Tests visuels avec références
npm run test:visual:reference  # Première fois
npm run test:visual            # Tests de régression

# Accessibilité SEPTEO
npm run test:axe              # Tests axe-core
npm run test:pa11y            # Tests Pa11y

# Performance Lighthouse  
npm run test:lighthouse       # Audit complet

# Design System SEPTEO
npm run test:design-system    # Validation conformité

# Parcours utilisateur
npm run test:user-flows       # Tests Playwright
```

### Test Complet
```bash
# Tous les agents UX
npm run test:all

# Génération rapport consolidé
npm run report
```

### Dashboard Temps Réel
```bash
# Monitoring UX continu
npm run dashboard
# → http://localhost:9002
```

## 📁 Structure des Tests

```
tests/
├── ux/
│   ├── visual-regression/          # Agent 1: Régression visuelle
│   │   ├── backstop.config.js
│   │   ├── scenarios/
│   │   └── backstop_data/
│   ├── accessibility/              # Agent 2: Accessibilité
│   │   ├── axe.config.js
│   │   ├── pa11y.config.js
│   │   └── axe-runner.js
│   ├── design-system/              # Agent 4: Design System
│   │   └── validator.js
│   ├── user-flows/                 # Agent 5: Parcours utilisateur
│   │   ├── critical-paths.spec.js
│   │   └── recordings/
│   ├── setup.js                    # Configuration globale
│   └── teardown.js                 # Consolidation rapports
├── scripts/
│   ├── ux-dashboard.js             # Dashboard temps réel
│   └── generate-ux-report.js       # Rapport hebdomadaire
├── lighthouserc.js                 # Agent 3: Performance
├── playwright.config.js            # Configuration Playwright
└── package.json                    # Dépendances UX
```

## 🎨 Standards SEPTEO Validés

### Couleurs Autorisées
```css
--primary-color: #ff6136;        /* Orange SEPTEO */
--primary-dark: #e55a2e;         /* Orange sombre */
--secondary-color: #2c3e50;      /* Bleu-gris */
--success-color: #27ae60;        /* Vert */
--warning-color: #f39c12;        /* Orange clair */
--danger-color: #e74c3c;         /* Rouge */
```

### Espacement (Grille 8px)
```css
/* Valeurs autorisées */
margin: 0, 8px, 16px, 24px, 32px, 40px, 48px...
padding: 0, 8px, 16px, 24px, 32px, 40px, 48px...
```

### Typographie Système
```css
font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
```

## 🚨 Alertes et Monitoring

### Dashboard Temps Réel
- **URL:** http://localhost:9002
- **Métriques:** Mise à jour toutes les 30 secondes
- **Alertes:** Notification si seuils non respectés
- **Actions:** Tests à la demande via interface

### Seuils d'Alerte
```javascript
{
  accessibility: { threshold: 95 },    // Score minimum
  performance: { threshold: 90 },      // Lighthouse score
  designSystem: { threshold: 90 },     // Conformité SEPTEO
  userExperience: { threshold: 90 }    // Succès parcours
}
```

## 📈 Rapports Automatiques

### Rapport Hebdomadaire
```bash
node scripts/generate-ux-report.js
```
Génère:
- `weekly-ux-report-YYYY-MM-DD.html` - Rapport détaillé
- `executive-ux-summary-YYYY-MM-DD.html` - Résumé exécutif
- `ux-consolidated-report.json` - Données consolidées

### Intégration CI/CD
GitHub Actions configurées pour:
- Tests automatiques sur chaque PR
- Commentaires automatiques avec scores
- Blocage si standards SEPTEO non respectés
- Rapports d'artefacts conservés 7 jours

## 🔧 Configuration Avancée

### Personnalisation BackstopJS
```javascript
// backstop.config.js
{
  "misMatchThreshold": 0.1,      // Seuil différence visuelle
  "requireSameDimensions": true,  // Même dimensions requises
  "resembleOutputOptions": {
    "ignoreAntialiasing": true
  }
}
```

### Personnalisation Lighthouse
```javascript
// lighthouserc.js
{
  assert: {
    "categories:performance": ["error", { minScore: 0.9 }],
    "largest-contentful-paint": ["error", { maxNumericValue: 2500 }]
  }
}
```

### Personnalisation Playwright
```javascript
// playwright.config.js
{
  timeout: 120000,              // 2 minutes par test
  retries: 2,                   // 2 tentatives en cas d'échec
  workers: 1,                   // Séquentiel pour Fire Salamander
  reporter: [['html'], ['json']]
}
```

## 🛠️ Debugging et Troubleshooting

### Problèmes Courants

#### 1. Tests visuels échouent
```bash
# Mettre à jour les références
npm run test:visual:approve

# Vider le cache
rm -rf tests/ux/visual-regression/backstop_data/bitmaps_test
```

#### 2. Fire Salamander non accessible
```bash
# Vérifier le statut
curl http://localhost:8080/api/v1/health

# Redémarrer
./fire-salamander --config config.yaml
```

#### 3. Tests Playwright timeout
```bash
# Augmenter le timeout
export PLAYWRIGHT_TIMEOUT=180000

# Mode debug
npx playwright test --debug
```

#### 4. Erreurs de conformité SEPTEO
```bash
# Rapport détaillé design system
npm run test:design-system
cat tests/reports/design-system/design-system-report.html
```

### Logs et Debug
```bash
# Logs dashboard UX
tail -f tests/ux-dashboard.log

# Logs tests Playwright
ls tests/reports/playwright/

# Screenshots échecs
ls tests/reports/playwright/artifacts/
```

## 📋 Checklist Qualité

### Avant Déploiement
- [ ] `npm run test:all` - Tous tests passent
- [ ] Score global SEPTEO ≥ 85%
- [ ] Zéro violation accessibilité critique
- [ ] Zéro friction critique parcours utilisateur
- [ ] Design system 100% conforme
- [ ] Performance Lighthouse ≥ 90%

### Validation Continue
- [ ] Dashboard UX monitore en temps réel
- [ ] Rapport hebdomadaire généré
- [ ] GitHub Actions intégrées
- [ ] Alertes configurées
- [ ] Équipe notifiée des régressions

## 🎯 Métriques de Succès

### KPIs UX Principaux
1. **Score Global SEPTEO:** ≥95% (excellence)
2. **Temps Moyen Parcours:** <30s analyse rapide
3. **Taux Succès Parcours:** ≥95% sans friction
4. **Score Accessibilité:** ≥95% WCAG AA
5. **Performance Web:** ≥90% Lighthouse

### Objectifs Business
- **Adoption:** Interface utilisable par tous
- **Conversion:** Parcours optimisés
- **Satisfaction:** UX fluide et cohérente
- **Conformité:** Standards SEPTEO respectés
- **Maintenance:** Détection automatique régressions

---

## 🆘 Support

Pour toute question sur la suite UX:

1. **Documentation:** Ce README + commentaires code
2. **Logs:** `tests/reports/` + `tests/ux-dashboard.log`  
3. **Dashboard:** http://localhost:9002 (monitoring)
4. **Rapports:** `tests/reports/consolidated/`

**La suite UX garantit que Fire Salamander offre une expérience utilisateur exceptionnelle conforme aux standards SEPTEO** 🔥