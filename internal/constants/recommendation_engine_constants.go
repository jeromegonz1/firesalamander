package constants

// ========================================
// CHARLIE-1 RECOMMENDATION ENGINE CONSTANTS
// ========================================

// ========================================
// RECOMMENDATION JSON FIELD NAMES
// ========================================

// Recommendation Structure Fields
const (
	RecommendationJSONFieldID           = "id"
	RecommendationJSONFieldTitle        = "title"
	RecommendationJSONFieldDescription  = "description"
	RecommendationJSONFieldCategory     = "category"
	RecommendationJSONFieldPriority     = "priority"
	RecommendationJSONFieldImpact       = "impact"
	RecommendationJSONFieldEffort       = "effort"
	RecommendationJSONFieldActions      = "actions"
	RecommendationJSONFieldResources    = "resources"
	RecommendationJSONFieldMetrics      = "metrics"
	RecommendationJSONFieldTags         = "tags"
)

// Action Item Fields
const (
	RecommendationJSONFieldTask         = "task"
	RecommendationJSONFieldTechnical    = "technical"
	RecommendationJSONFieldEstimatedTime = "estimated_time"
)

// ========================================
// RECOMMENDATION PRIORITY LEVELS
// ========================================

// Priority Constants (string values for JSON)
const (
	RecommendationPriorityCritical = "critical"
	RecommendationPriorityHigh     = "high"
	RecommendationPriorityMedium   = "medium"
	RecommendationPriorityLow      = "low"
)

// ========================================
// RECOMMENDATION IMPACT LEVELS
// ========================================

// Impact Constants (string values for JSON)
const (
	RecommendationImpactHigh   = "high"
	RecommendationImpactMedium = "medium"
	RecommendationImpactLow    = "low"
)

// ========================================
// RECOMMENDATION EFFORT LEVELS
// ========================================

// Effort Constants (string values for JSON)
const (
	RecommendationEffortLow    = "low"
	RecommendationEffortMedium = "medium"
	RecommendationEffortHigh   = "high"
)

// ========================================
// RECOMMENDATION CATEGORIES
// ========================================

// Category Names
const (
	RecommendationCategoryTags        = "tags"
	RecommendationCategoryPerformance = "performance"
	RecommendationCategorySecurity    = "security"
	RecommendationCategoryGeneral     = "general"
	RecommendationCategoryTechnical   = "technical"
	RecommendationCategoryContent     = "content"
	RecommendationCategoryMobile      = "mobile"
	RecommendationCategoryStructure   = "structure"
	RecommendationCategoryAccessibility = "accessibility"
)

// ========================================
// RECOMMENDATION TEMPLATE IDS
// ========================================

// Tag-related Templates
const (
	RecommendationTemplateIDMissingTitle     = "missing-title"
	RecommendationTemplateIDMissingMetaDesc  = "missing-meta-desc"
	RecommendationTemplateIDTitleLength      = "title-length"
	RecommendationTemplateIDMetaDescLength   = "meta-desc-length"
	RecommendationTemplateIDMissingH1        = "missing-h1"
	RecommendationTemplateIDMultipleH1       = "multiple-h1"
	RecommendationTemplateIDHeadingHierarchy = "heading-hierarchy"
	RecommendationTemplateIDMissingAltTags   = "missing-alt-tags"
	RecommendationTemplateIDMissingViewport  = "missing-viewport"
	RecommendationTemplateIDMissingCanonical = "missing-canonical"
	RecommendationTemplateIDMissingOGTags    = "missing-og-tags"
)

// Performance-related Templates
const (
	RecommendationTemplateIDImproveLCP       = "improve-lcp"
	RecommendationTemplateIDImproveFID       = "improve-fid"
	RecommendationTemplateIDImproveCLS       = "improve-cls"
	RecommendationTemplateIDEnableCompression = "enable-compression"
	RecommendationTemplateIDConfigureCaching = "configure-caching"
	RecommendationTemplateIDOptimizeImages   = "optimize-images"
	RecommendationTemplateIDMinifyResources  = "minify-resources"
	RecommendationTemplateIDReducePageSize   = "reduce-page-size"
)

// Security-related Templates
const (
	RecommendationTemplateIDMigrateHTTPS    = "migrate-https"
	RecommendationTemplateIDFixMixedContent = "fix-mixed-content"
	RecommendationTemplateIDAddHSTS         = "add-hsts"
)

// Technical-related Templates
const (
	RecommendationTemplateIDMakeResponsive         = "make-responsive"
	RecommendationTemplateIDAddSitemap             = "add-sitemap"
	RecommendationTemplateIDAddRobotsTxt           = "add-robots-txt"
	RecommendationTemplateIDRemoveNoIndex          = "remove-noindex"
	RecommendationTemplateIDFixDuplicateContent    = "fix-duplicate-content"
	RecommendationTemplateIDImproveAccessibility  = "improve-accessibility"
	RecommendationTemplateIDFixBrokenLinks         = "fix-broken-links"
	RecommendationTemplateIDImproveInternalLinking = "improve-internal-linking"
)

// Score-based Templates
const (
	RecommendationTemplateIDOverallSEOAudit = "overall-seo-audit"
)

// ========================================
// RECOMMENDATION SCORE COMPARISONS
// ========================================

// Performance Score Values
const (
	RecommendationScorePoor             = "poor"
	RecommendationScoreGood             = "good"
	RecommendationScoreNeedsImprovement = "needs_improvement"
)

// ========================================
// RECOMMENDATION TIME ESTIMATES
// ========================================

// Time Estimates for Different Effort Levels
const (
	RecommendationTimeLow      = "1-2 heures"
	RecommendationTimeMedium   = "4-8 heures"
	RecommendationTimeHigh     = "1-2 jours"
	RecommendationTimeVariable = "Variable"
)

// ========================================
// RECOMMENDATION CONTENT RANGES
// ========================================

// Optimal Content Length Ranges
const (
	RecommendationTitleRange    = "30-60 caractères"
	RecommendationMetaDescRange = "120-160 caractères"
)

// ========================================
// RECOMMENDATION TARGET VALUES
// ========================================

// Core Web Vitals Targets
const (
	RecommendationTargetLCP      = "≤ 2.5s"
	RecommendationTargetFID      = "≤ 100ms"
	RecommendationTargetCLS      = "≤ 0.1"
	RecommendationTargetPageSize = "< 2 MB"
)

// SEO Score Targets
const (
	RecommendationTargetScore       = "≥ 70%"
	RecommendationTargetScoreHigh   = "≥ 85%"
	RecommendationTargetLinks       = "≥ 3"
)

// ========================================
// RECOMMENDATION RESOURCE TYPES
// ========================================

// Resource Types
const (
	RecommendationResourceTypeDoc     = "documentation"
	RecommendationResourceTypeGuide   = "guide"
	RecommendationResourceTypeTool    = "tool"
	RecommendationResourceTypeExample = "example"
)

// ========================================
// RECOMMENDATION NUMERIC THRESHOLDS
// ========================================

// Scoring Thresholds
const (
	RecommendationMaxRecommendations = 20
	RecommendationScoreThresholdLow  = 50.0
	RecommendationCategoryThresholdLow = 0.5
	RecommendationAccessibilityThreshold = 0.7
	RecommendationAltTextThreshold = 1.0
	RecommendationMinInternalLinks = 3
	RecommendationMaxPageSizeBytes = 2 * 1024 * 1024 // 2MB
)

// Priority Weights for Sorting
const (
	RecommendationPriorityWeightCritical = 4
	RecommendationPriorityWeightHigh     = 3
	RecommendationPriorityWeightMedium   = 2
	RecommendationPriorityWeightLow      = 1
	RecommendationPriorityWeightDefault  = 0
)

// Impact Weights for Sorting
const (
	RecommendationImpactWeightHigh    = 3
	RecommendationImpactWeightMedium  = 2
	RecommendationImpactWeightLow     = 1
	RecommendationImpactWeightDefault = 0
)

// ========================================
// RECOMMENDATION PRIORITY RULE SCORES
// ========================================

// Priority Rule Scores (out of 100)
const (
	RecommendationRuleMissingTitle    = 90
	RecommendationRuleMissingMetaDesc = 85 // Using existing HighQualityScore
	RecommendationRuleMigrateHTTPS    = 95
	RecommendationRuleMissingViewport = 85
	RecommendationRuleImproveLCP      = 75
	RecommendationRuleMissingH1       = 70
)

// ========================================
// RECOMMENDATION DEFAULT MESSAGES
// ========================================

// Default Recommendation Content
const (
	RecommendationDefaultTitle       = "Recommandation SEO"
	RecommendationDefaultDescription = "Amélioration SEO recommandée"
	RecommendationDefaultCategory    = "general"
)

// Issue Messages
const (
	RecommendationIssueTitleMissing = "Titre manquant"
	RecommendationIssueMetaDescMissing = "Meta description manquante"
	RecommendationIssueH1Missing = "H1 manquant"
	RecommendationIssueViewportMissing = "Meta viewport manquante"
)

// ========================================
// RECOMMENDATION TEMPLATE STRINGS
// ========================================

// Template Titles
const (
	RecommendationTitleAddPageTitle = "Ajouter un titre de page"
	RecommendationTitleAddMetaDesc  = "Ajouter une meta description"
	RecommendationTitleMigrateHTTPS = "Migrer vers HTTPS"
	RecommendationTitleImproveLCP   = "Améliorer le Largest Contentful Paint (LCP)"
)

// Template Descriptions
const (
	RecommendationDescMissingTitle = "La page n'a pas de balise <title>. C'est un élément fondamental pour le SEO."
	RecommendationDescMissingMetaDesc = "La page n'a pas de meta description. Elle influence le taux de clic dans les résultats de recherche."
	RecommendationDescMigrateHTTPS = "Le site utilise HTTP au lieu de HTTPS. Google favorise les sites sécurisés."
	RecommendationDescImproveLCP = "Le LCP actuel est de {current_value}, l'objectif est {target}. Optimisez le chargement du contenu principal."
)

// ========================================
// RECOMMENDATION ACTION ITEMS
// ========================================

// Common Action Templates
const (
	RecommendationActionAddTitleTag      = "Ajouter une balise <title> descriptive et unique"
	RecommendationActionIncludeKeywords  = "Inclure les mots-clés principaux"
	RecommendationActionRespectLength    = "Respecter la longueur optimale ({range})"
	RecommendationActionAddMetaDesc      = "Ajouter une meta description attrayante"
	RecommendationActionIncludeCTA       = "Inclure un appel à l'action"
	RecommendationActionGetSSLCert       = "Obtenir un certificat SSL/TLS"
	RecommendationActionConfigureHTTPS   = "Configurer le serveur pour HTTPS"
	RecommendationActionRedirectHTTPS    = "Rediriger tout le trafic HTTP vers HTTPS"
	RecommendationActionUpdateLinks      = "Mettre à jour les liens internes"
	RecommendationActionOptimizeImages   = "Optimiser les images de l'above-the-fold"
	RecommendationActionImproveServer    = "Améliorer le temps de réponse du serveur"
	RecommendationActionPreloadResources = "Précharger les ressources critiques"
	RecommendationActionUseCDN           = "Utiliser un CDN"
)

// ========================================
// RECOMMENDATION METRICS
// ========================================

// SEO Metrics
const (
	RecommendationMetricCTR          = "Taux de clic"
	RecommendationMetricSERPPosition = "Position dans les SERP"
	RecommendationMetricImpressions  = "Impressions"
	RecommendationMetricTrustSignals = "Trust signals"
	RecommendationMetricRankingBoost = "Ranking boost"
	RecommendationMetricPageExperience = "Page Experience"
	RecommendationMetricCoreWebVitals = "Core Web Vitals"
	RecommendationMetricLCP          = "LCP"
)

// ========================================
// RECOMMENDATION TAGS
// ========================================

// Tag Categories
const (
	RecommendationTagCritical      = "critique"
	RecommendationTagTags          = "balises"
	RecommendationTagOnPage        = "onpage"
	RecommendationTagMeta          = "meta"
	RecommendationTagCTR           = "ctr"
	RecommendationTagSecurity      = "sécurité"
	RecommendationTagTechnical     = "technique"
	RecommendationTagPerformance   = "performance"
	RecommendationTagCoreWebVitals = "core-web-vitals"
	RecommendationTagLCP           = "lcp"
)

// ========================================
// RECOMMENDATION OPERATORS
// ========================================

// Comparison Operators
const (
	RecommendationOperatorLessEqual    = "≤"
	RecommendationOperatorGreaterEqual = "≥"
	RecommendationOperatorLess         = "<"
	RecommendationOperatorGreater      = ">"
)

// ========================================
// RECOMMENDATION PLACEHOLDER PATTERNS
// ========================================

// Template Placeholder Pattern
const (
	RecommendationPlaceholderPattern = "{%s}"
)

// Common Placeholders
const (
	RecommendationPlaceholderCurrentValue = "{current_value}"
	RecommendationPlaceholderTarget       = "{target}"
	RecommendationPlaceholderRange        = "{range}"
	RecommendationPlaceholderCategory     = "{category}"
	RecommendationPlaceholderCount        = "{count}"
)