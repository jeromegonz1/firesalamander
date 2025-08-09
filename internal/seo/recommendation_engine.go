package seo

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"firesalamander/internal/constants"
)

// RecommendationEngine moteur de recommandations SEO intelligent
type RecommendationEngine struct {
	// Règles de priorité
	priorityRules map[string]int
	// Templates de recommandations
	templates map[string]RecommendationTemplate
}

// SEORecommendation recommandation SEO structurée
type SEORecommendation struct {
	ID           string                 `json:"id"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	Category     string                 `json:"category"`
	Priority     Priority               `json:"priority"`
	Impact       Impact                 `json:"impact"`
	Effort       Effort                 `json:"effort"`
	Actions      []ActionItem           `json:"actions"`
	Resources    []Resource             `json:"resources"`
	Metrics      []string               `json:"metrics"`
	Tags         []string               `json:"tags"`
}

// RecommendationTemplate template pour générer des recommandations
type RecommendationTemplate struct {
	ID          string
	Title       string
	Description string
	Category    string
	Priority    Priority
	Impact      Impact
	Effort      Effort
	Actions     []string
	Resources   []string
	Metrics     []string
	Tags        []string
	Conditions  []string // Conditions pour déclencher cette recommandation
}

// ActionItem action à effectuer
type ActionItem struct {
	Task        string `json:"task"`
	Description string `json:"description"`
	Technical   bool   `json:"technical"`
	Estimated   string `json:"estimated_time"`
}

// Priority niveau de priorité
type Priority string

const (
	PriorityCritical Priority = "critical"
	PriorityHigh     Priority = "high"
	PriorityMedium   Priority = "medium"
	PriorityLow      Priority = "low"
)

// Impact impact sur le SEO
type Impact string

const (
	ImpactHigh   Impact = "high"
	ImpactMedium Impact = "medium"
	ImpactLow    Impact = "low"
)

// Effort effort requis
type Effort string

const (
	EffortLow    Effort = "low"
	EffortMedium Effort = "medium"
	EffortHigh   Effort = "high"
)

// NewRecommendationEngine crée un nouveau moteur de recommandations
func NewRecommendationEngine() *RecommendationEngine {
	engine := &RecommendationEngine{
		priorityRules: make(map[string]int),
		templates:     make(map[string]RecommendationTemplate),
	}

	// Initialiser les règles de priorité
	engine.initPriorityRules()
	
	// Initialiser les templates
	engine.initRecommendationTemplates()

	return engine
}

// GenerateRecommendations génère les recommandations basées sur l'analyse
func (re *RecommendationEngine) GenerateRecommendations(analysis *RealPageAnalysis) []SEORecommendation {
	var recommendations []SEORecommendation

	// 1. Recommandations basées sur les balises
	recommendations = append(recommendations, re.generateTagRecommendations(&analysis.TagAnalysis)...)

	// 2. Recommandations basées sur les performances
	recommendations = append(recommendations, re.generatePerformanceRecommendations(&analysis.PerformanceMetrics)...)

	// 3. Recommandations basées sur l'audit technique
	recommendations = append(recommendations, re.generateTechnicalRecommendations(&analysis.TechnicalAudit)...)

	// 4. Recommandations générales basées sur les scores
	recommendations = append(recommendations, re.generateScoreBasedRecommendations(analysis)...)

	// 5. Trier par priorité et impact
	re.sortRecommendations(recommendations)

	// 6. Déduplication et optimisation
	recommendations = re.deduplicateRecommendations(recommendations)

	// 7. Limiter le nombre de recommandations (top 20)
	if len(recommendations) > constants.RecommendationMaxRecommendations {
		recommendations = recommendations[:20]
	}

	return recommendations
}

// generateTagRecommendations génère les recommandations pour les balises
func (re *RecommendationEngine) generateTagRecommendations(tagAnalysis *TagAnalysisResult) []SEORecommendation {
	var recs []SEORecommendation

	// Titre
	if !tagAnalysis.Title.Present {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDMissingTitle, map[string]interface{}{
			"issue": constants.RecommendationIssueTitleMissing,
		}))
	} else if !tagAnalysis.Title.OptimalLength {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDTitleLength, map[string]interface{}{
			"current_length": tagAnalysis.Title.Length,
			"optimal_range":  constants.RecommendationTitleRange,
		}))
	}

	// Meta description
	if !tagAnalysis.MetaDescription.Present {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDMissingMetaDesc, map[string]interface{}{}))
	} else if !tagAnalysis.MetaDescription.OptimalLength {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDMetaDescLength, map[string]interface{}{
			"current_length": tagAnalysis.MetaDescription.Length,
			"optimal_range":  constants.RecommendationMetaDescRange,
		}))
	}

	// Structure des headings
	if tagAnalysis.Headings.H1Count == 0 {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDMissingH1, map[string]interface{}{}))
	} else if tagAnalysis.Headings.H1Count > 1 {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDMultipleH1, map[string]interface{}{
			"count": tagAnalysis.Headings.H1Count,
		}))
	}

	if !tagAnalysis.Headings.HasHierarchy {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDHeadingHierarchy, map[string]interface{}{}))
	}

	// Images
	if tagAnalysis.Images.AltTextCoverage < constants.RecommendationAltTextThreshold {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDMissingAltTags, map[string]interface{}{
			"coverage": fmt.Sprintf("%.1f%%", tagAnalysis.Images.AltTextCoverage*100),
			"missing":  len(tagAnalysis.Images.MissingAltImages),
		}))
	}

	// Meta tags essentiels
	if !tagAnalysis.MetaTags.HasViewport {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDMissingViewport, map[string]interface{}{}))
	}

	if !tagAnalysis.MetaTags.HasCanonical {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDMissingCanonical, map[string]interface{}{}))
	}

	if !tagAnalysis.MetaTags.HasOGTags {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDMissingOGTags, map[string]interface{}{}))
	}

	return recs
}

// generatePerformanceRecommendations génère les recommandations de performance
func (re *RecommendationEngine) generatePerformanceRecommendations(perfMetrics *PerformanceMetricsResult) []SEORecommendation {
	var recs []SEORecommendation

	// Temps de chargement
	if perfMetrics.CoreWebVitals.LCP.Score == constants.RecommendationScorePoor {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDImproveLCP, map[string]interface{}{
			"current_value": fmt.Sprintf("%.1fms", perfMetrics.CoreWebVitals.LCP.Value),
			"target":        constants.RecommendationTargetLCP,
		}))
	}

	if perfMetrics.CoreWebVitals.FID.Score == constants.RecommendationScorePoor {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDImproveFID, map[string]interface{}{
			"current_value": fmt.Sprintf("%.1fms", perfMetrics.CoreWebVitals.FID.Value),
			"target":        constants.RecommendationTargetFID,
		}))
	}

	if perfMetrics.CoreWebVitals.CLS.Score == constants.RecommendationScorePoor {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDImproveCLS, map[string]interface{}{
			"current_value": fmt.Sprintf("%.3f", perfMetrics.CoreWebVitals.CLS.Value),
			"target":        constants.RecommendationTargetCLS,
		}))
	}

	// Compression
	if !perfMetrics.HasCompression {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDEnableCompression, map[string]interface{}{}))
	}

	// Cache
	if !perfMetrics.HasCaching {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDConfigureCaching, map[string]interface{}{}))
	}

	// Images
	if !perfMetrics.OptimizedImages && perfMetrics.ResourceCounts.Images > 0 {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDOptimizeImages, map[string]interface{}{
			"image_count": perfMetrics.ResourceCounts.Images,
		}))
	}

	// Minification
	if !perfMetrics.MinifiedResources {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDMinifyResources, map[string]interface{}{}))
	}

	// Taille de page
	if perfMetrics.PageSize > constants.RecommendationMaxPageSizeBytes { // > 2MB
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDReducePageSize, map[string]interface{}{
			"current_size": fmt.Sprintf("%.1f MB", float64(perfMetrics.PageSize)/(1024*1024)),
			"target":       constants.RecommendationTargetPageSize,
		}))
	}

	return recs
}

// generateTechnicalRecommendations génère les recommandations techniques
func (re *RecommendationEngine) generateTechnicalRecommendations(techAudit *TechnicalAuditResult) []SEORecommendation {
	var recs []SEORecommendation

	// Sécurité
	if !techAudit.Security.HasHTTPS {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDMigrateHTTPS, map[string]interface{}{}))
	}

	if techAudit.Security.MixedContent {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDFixMixedContent, map[string]interface{}{
			"insecure_links": len(techAudit.Security.InsecureLinks),
		}))
	}

	if !techAudit.Security.HasHSTS {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDAddHSTS, map[string]interface{}{}))
	}

	// Mobile
	if !techAudit.Mobile.IsResponsive {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDMakeResponsive, map[string]interface{}{}))
	}

	// Structure
	if !techAudit.Structure.HasSitemap {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDAddSitemap, map[string]interface{}{}))
	}

	if !techAudit.Structure.HasRobotsTxt {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDAddRobotsTxt, map[string]interface{}{}))
	}

	// Indexabilité
	if techAudit.Indexability.HasNoIndex {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDRemoveNoIndex, map[string]interface{}{}))
	}

	if techAudit.Indexability.DuplicateContent {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDFixDuplicateContent, map[string]interface{}{}))
	}

	// Accessibilité
	if techAudit.Accessibility.Score < constants.RecommendationAccessibilityThreshold {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDImproveAccessibility, map[string]interface{}{
			"current_score": fmt.Sprintf("%.1f%%", techAudit.Accessibility.Score*100),
			"target":        "≥ " + strconv.Itoa(constants.RecommendationRuleMissingMetaDesc) + "%",
		}))
	}

	// Crawlabilité
	if len(techAudit.Crawlability.BrokenLinks) > 0 {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDFixBrokenLinks, map[string]interface{}{
			"broken_count": len(techAudit.Crawlability.BrokenLinks),
		}))
	}

	if techAudit.Crawlability.InternalLinks < constants.RecommendationMinInternalLinks {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDImproveInternalLinking, map[string]interface{}{
			"current_links": techAudit.Crawlability.InternalLinks,
			"target":        constants.RecommendationTargetLinks,
		}))
	}

	return recs
}

// generateScoreBasedRecommendations génère des recommandations basées sur les scores globaux
func (re *RecommendationEngine) generateScoreBasedRecommendations(analysis *RealPageAnalysis) []SEORecommendation {
	var recs []SEORecommendation

	// Score global faible
	if analysis.OverallScore < constants.RecommendationScoreThresholdLow {
		recs = append(recs, re.createRecommendation(constants.RecommendationTemplateIDOverallSEOAudit, map[string]interface{}{
			"current_score": fmt.Sprintf("%.1f", analysis.OverallScore),
			"target":        constants.RecommendationTargetScore,
		}))
	}

	// Scores par catégorie
	for category, score := range analysis.CategoryScores {
		if score < constants.RecommendationCategoryThresholdLow { // Score < 50%
			recs = append(recs, re.createRecommendation(fmt.Sprintf("improve-%s", category), map[string]interface{}{
				"category":      category,
				"current_score": fmt.Sprintf("%.1f%%", score*100),
				"target":        constants.RecommendationTargetScore,
			}))
		}
	}

	return recs
}

// createRecommendation crée une recommandation à partir d'un template
func (re *RecommendationEngine) createRecommendation(templateID string, params map[string]interface{}) SEORecommendation {
	template, exists := re.templates[templateID]
	if !exists {
		// Template par défaut
		return SEORecommendation{
			ID:          templateID,
			Title:       constants.RecommendationDefaultTitle,
			Description: constants.RecommendationDefaultDescription,
			Category:    constants.RecommendationCategoryGeneral,
			Priority:    PriorityMedium,
			Impact:      ImpactMedium,
			Effort:      EffortMedium,
		}
	}

	rec := SEORecommendation{
		ID:          template.ID,
		Title:       template.Title,
		Description: template.Description,
		Category:    template.Category,
		Priority:    template.Priority,
		Impact:      template.Impact,
		Effort:      template.Effort,
		Metrics:     template.Metrics,
		Tags:        template.Tags,
	}

	// Personnaliser avec les paramètres
	rec.Description = re.interpolateTemplate(template.Description, params)

	// Créer les actions
	for _, actionTemplate := range template.Actions {
		action := ActionItem{
			Task:        re.interpolateTemplate(actionTemplate, params),
			Description: re.interpolateTemplate(actionTemplate, params),
			Technical:   strings.Contains(actionTemplate, constants.RecommendationTagTechnical) || strings.Contains(actionTemplate, "code"),
			Estimated:   re.estimateTime(template.Effort),
		}
		rec.Actions = append(rec.Actions, action)
	}

	// Créer les ressources
	for _, resourceTemplate := range template.Resources {
		resource := Resource{
			URL:  resourceTemplate,
			Type: constants.RecommendationResourceTypeDoc,
		}
		rec.Resources = append(rec.Resources, resource)
	}

	return rec
}

// interpolateTemplate remplace les placeholders dans les templates
func (re *RecommendationEngine) interpolateTemplate(template string, params map[string]interface{}) string {
	result := template
	for key, value := range params {
		placeholder := fmt.Sprintf(constants.RecommendationPlaceholderPattern, key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

// estimateTime estime le temps requis
func (re *RecommendationEngine) estimateTime(effort Effort) string {
	switch effort {
	case EffortLow:
		return constants.RecommendationTimeLow
	case EffortMedium:
		return constants.RecommendationTimeMedium
	case EffortHigh:
		return constants.RecommendationTimeHigh
	default:
		return constants.RecommendationTimeVariable
	}
}

// sortRecommendations trie les recommandations par priorité et impact
func (re *RecommendationEngine) sortRecommendations(recommendations []SEORecommendation) {
	sort.Slice(recommendations, func(i, j int) bool {
		// Priorité d'abord
		iPrio := re.getPriorityWeight(recommendations[i].Priority)
		jPrio := re.getPriorityWeight(recommendations[j].Priority)
		
		if iPrio != jPrio {
			return iPrio > jPrio
		}
		
		// Puis impact
		iImpact := re.getImpactWeight(recommendations[i].Impact)
		jImpact := re.getImpactWeight(recommendations[j].Impact)
		
		return iImpact > jImpact
	})
}

// getPriorityWeight retourne le poids numérique de la priorité
func (re *RecommendationEngine) getPriorityWeight(priority Priority) int {
	switch priority {
	case PriorityCritical:
		return 4
	case PriorityHigh:
		return 3
	case PriorityMedium:
		return 2
	case PriorityLow:
		return 1
	default:
		return 0
	}
}

// getImpactWeight retourne le poids numérique de l'impact
func (re *RecommendationEngine) getImpactWeight(impact Impact) int {
	switch impact {
	case ImpactHigh:
		return 3
	case ImpactMedium:
		return 2
	case ImpactLow:
		return 1
	default:
		return 0
	}
}

// deduplicateRecommendations supprime les doublons
func (re *RecommendationEngine) deduplicateRecommendations(recommendations []SEORecommendation) []SEORecommendation {
	seen := make(map[string]bool)
	var unique []SEORecommendation

	for _, rec := range recommendations {
		key := rec.ID + "|" + rec.Category
		if !seen[key] {
			seen[key] = true
			unique = append(unique, rec)
		}
	}

	return unique
}

// initPriorityRules initialise les règles de priorité
func (re *RecommendationEngine) initPriorityRules() {
	re.priorityRules[constants.RecommendationTemplateIDMissingTitle] = constants.RecommendationRuleMissingTitle
	re.priorityRules[constants.RecommendationTemplateIDMissingMetaDesc] = constants.RecommendationRuleMissingMetaDesc
	re.priorityRules[constants.RecommendationTemplateIDMigrateHTTPS] = constants.RecommendationRuleMigrateHTTPS
	re.priorityRules[constants.RecommendationTemplateIDMissingViewport] = constants.RecommendationRuleMissingMetaDesc
	re.priorityRules[constants.RecommendationTemplateIDImproveLCP] = constants.RecommendationRuleImproveLCP
	re.priorityRules[constants.RecommendationTemplateIDMissingH1] = constants.RecommendationRuleMissingH1
	// ... plus de règles
}

// initRecommendationTemplates initialise les templates de recommandations
func (re *RecommendationEngine) initRecommendationTemplates() {
	// Template: Titre manquant
	re.templates[constants.RecommendationTemplateIDMissingTitle] = RecommendationTemplate{
		ID:          constants.RecommendationTemplateIDMissingTitle,
		Title:       constants.RecommendationTitleAddPageTitle,
		Description: constants.RecommendationDescMissingTitle,
		Category:    constants.RecommendationCategoryTags,
		Priority:    PriorityCritical,
		Impact:      ImpactHigh,
		Effort:      EffortLow,
		Actions: []string{
			constants.RecommendationActionAddTitleTag,
			constants.RecommendationActionIncludeKeywords,
			constants.RecommendationActionRespectLength,
		},
		Resources: []string{
			constants.GoogleTitleLinkDocs,
		},
		Metrics: []string{constants.RecommendationMetricCTR, constants.RecommendationMetricSERPPosition},
		Tags:    []string{constants.RecommendationTagCritical, constants.RecommendationTagTags, constants.RecommendationTagOnPage},
	}

	// Template: Meta description manquante
	re.templates[constants.RecommendationTemplateIDMissingMetaDesc] = RecommendationTemplate{
		ID:          constants.RecommendationTemplateIDMissingMetaDesc,
		Title:       constants.RecommendationTitleAddMetaDesc,
		Description: constants.RecommendationDescMissingMetaDesc,
		Category:    constants.RecommendationCategoryTags,
		Priority:    PriorityHigh,
		Impact:      ImpactHigh,
		Effort:      EffortLow,
		Actions: []string{
			constants.RecommendationActionAddMetaDesc,
			constants.RecommendationActionIncludeCTA,
			constants.RecommendationActionRespectLength,
		},
		Resources: []string{
			constants.GoogleSnippetDocs,
		},
		Metrics: []string{constants.RecommendationMetricCTR, constants.RecommendationMetricImpressions},
		Tags:    []string{constants.RecommendationTagMeta, constants.RecommendationTagTags, constants.RecommendationTagCTR},
	}

	// Template: Migration HTTPS
	re.templates[constants.RecommendationTemplateIDMigrateHTTPS] = RecommendationTemplate{
		ID:          constants.RecommendationTemplateIDMigrateHTTPS,
		Title:       constants.RecommendationTitleMigrateHTTPS,
		Description: constants.RecommendationDescMigrateHTTPS,
		Category:    constants.RecommendationCategorySecurity,
		Priority:    PriorityCritical,
		Impact:      ImpactHigh,
		Effort:      EffortHigh,
		Actions: []string{
			constants.RecommendationActionGetSSLCert,
			constants.RecommendationActionConfigureHTTPS,
			constants.RecommendationActionRedirectHTTPS,
			constants.RecommendationActionUpdateLinks,
		},
		Resources: []string{
			constants.GoogleHTTPSDocs,
		},
		Metrics: []string{constants.RecommendationMetricTrustSignals, constants.RecommendationMetricRankingBoost},
		Tags:    []string{constants.RecommendationTagCritical, constants.RecommendationTagSecurity, constants.RecommendationTagTechnical},
	}

	// Template: Core Web Vitals - LCP
	re.templates[constants.RecommendationTemplateIDImproveLCP] = RecommendationTemplate{
		ID:          constants.RecommendationTemplateIDImproveLCP,
		Title:       constants.RecommendationTitleImproveLCP,
		Description: constants.RecommendationDescImproveLCP,
		Category:    constants.RecommendationCategoryPerformance,
		Priority:    PriorityHigh,
		Impact:      ImpactHigh,
		Effort:      EffortMedium,
		Actions: []string{
			constants.RecommendationActionOptimizeImages,
			constants.RecommendationActionImproveServer,
			constants.RecommendationActionPreloadResources,
			constants.RecommendationActionUseCDN,
		},
		Resources: []string{
			constants.WebDevLCPDocs,
		},
		Metrics: []string{constants.RecommendationMetricLCP, constants.RecommendationMetricPageExperience, constants.RecommendationMetricCoreWebVitals},
		Tags:    []string{constants.RecommendationCategoryPerformance, constants.RecommendationTagCoreWebVitals, constants.RecommendationTagLCP},
	}

	// ... Ajouter plus de templates selon les besoins
}