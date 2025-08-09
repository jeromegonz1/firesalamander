# 🔍 ANALYSE PATTERNS HARDCODING - FIRE SALAMANDER

## 📊 CLASSIFICATION DES 54 VIOLATIONS RESTANTES

### 🚨 HIGH PRIORITY (25 violations)

#### 1. URLs de Test (8 violations)
```go
// AVANT:
BaseURL: "http://localhost:3000"
VALUES ('https://example.com', 'completed', 10, 8)

// APRÈS:
BaseURL: constants.TestLocalhost3000
VALUES (constants.TestExampleURL, 'completed', 10, 8)
```

#### 2. Scores Hardcodés "80" (12 violations)
```go
// AVANT:
MinCoverage: 80.0
case score >= 80:
"target": "≥ 80%"

// APRÈS:
MinCoverage: constants.MinCoverageThreshold
case score >= constants.HighQualityScore:
"target": "≥ " + strconv.Itoa(constants.HighQualityScore) + "%"
```

#### 3. Magic Number "3000" (5 violations)
```go
// AVANT:
} else if value <= 3000 {
hasPortMapping := strings.Contains(contentStr, "3000:3000")

// APRÈS:
} else if value <= constants.TestValue3000 {
hasPortMapping := strings.Contains(contentStr, constants.TestPortMapping)
```

### 🟡 MEDIUM PRIORITY (19 violations)

#### 4. Timeouts Tests (8 violations)
```go
// AVANT:
Duration: 2 * time.Minute
RampUpTime: 30 * time.Second
ResponseTime: 200 * time.Millisecond

// APRÈS:  
Duration: constants.TestDuration2Min
RampUpTime: constants.TestRampUpTime
ResponseTime: constants.FastResponseTime
```

#### 5. URLs Documentation (4 violations)
```go
// AVANT:
"https://developers.google.com/search/docs/appearance/title-link"
"https://web.dev/lcp/"

// APRÈS:
constants.GoogleTitleLinkDocs
constants.WebDevLCPDocs
```

#### 6. Protocoles HTTP (7 violations)
```go
// AVANT:
strings.Contains(htmlContent, `http://`)
strings.HasPrefix(link, "http://")
"HEAD", "https://"+host

// APRÈS:
strings.Contains(htmlContent, constants.HTTPPrefix)
strings.HasPrefix(link, constants.HTTPPrefix)  
"HEAD", constants.HTTPSPrefix+host
```

### 🟢 LOW PRIORITY (10 violations)

#### 7. JSON Tags Struct (6 violations)
```go
// AVANT:
Domain string `json:"domain"`
Count int `json:"count"`

// Strategy: Keep as-is (JSON tags are API contract)
```

#### 8. Imports System (4 violations)
```go
// AVANT:
"context"
"syscall"

// Strategy: Keep as-is (Go standard library)
```

---

## 🎯 STRATÉGIE DE CORRECTION PAR BATCH

### Batch 1 (25 violations) - HIGH PRIORITY
**Focus**: URLs de test + Scores 80 + Magic number 3000  
**ETA**: 30 minutes  
**Impact**: Critique pour les tests

### Batch 2 (19 violations) - MEDIUM PRIORITY  
**Focus**: Timeouts tests + URLs doc + Protocoles HTTP  
**ETA**: 25 minutes  
**Impact**: Amélioration qualité code

### Batch 3 (10 violations) - LOW PRIORITY
**Focus**: JSON tags + Imports (évaluation si vraiment nécessaire)  
**ETA**: 15 minutes  
**Impact**: Perfectionnisme

---

## 📈 MÉTRIQUES PATTERNS

| Pattern Type | Count | % Total | Complexity | Priority |
|--------------|-------|---------|------------|----------|
| URLs Test | 8 | 15% | Simple | 🚨 HIGH |
| Score "80" | 12 | 22% | Simple | 🚨 HIGH |
| Magic "3000" | 5 | 9% | Simple | 🚨 HIGH |
| Timeouts | 8 | 15% | Medium | 🟡 MED |
| URLs Doc | 4 | 7% | Simple | 🟡 MED |
| Protocoles | 7 | 13% | Medium | 🟡 MED |
| JSON Tags | 6 | 11% | Complex | 🟢 LOW |
| Imports | 4 | 8% | N/A | 🟢 LOW |

---

## 🎯 RECOMMANDATIONS SPRINT SUIVANT

### Objectif Réaliste: 35 violations en 2 sprints
- **Sprint 0.3**: Éliminer 25 violations HIGH PRIORITY → 29 restantes
- **Sprint 0.4**: Éliminer 19 violations MEDIUM PRIORITY → 10 restantes
- **Sprint 0.5**: Évaluation finale des 10 LOW PRIORITY

### Seuil de Tolérance:
- **ACCEPTABLE**: < 15 violations (GREEN)
- **EXCELLENT**: < 5 violations (GOLD)  
- **PERFECTION**: 0 violations (LEGEND)

---

*Analyse générée par Code Quality Inspector*