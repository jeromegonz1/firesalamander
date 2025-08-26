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
	Actions      []ActionItem                `json:"actions"`
	Resources    []RecommendationResource    `json:"resources"`
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

// RecommendationResource represents a helpful resource for implementing recommendations
type RecommendationResource struct {
	URL  string `json:"url"`
	Type string `json:"type"`
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

	// TODO: Fix field mapping after duplicates elimination
	// 1. Recommandations basées sur le titre
	// recommendations = append(recommendations, re.generateTitleRecommendations(&analysis.Title)...)

	// 2. Recommandations basées sur la performance  
	// recommendations = append(recommendations, re.generatePerformanceRecommendations(&analysis.Performance)...)

	// 3. Recommandations basées sur le contenu
	// recommendations = append(recommendations, re.generateContentRecommendations(&analysis.Content)...)

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
	if analysis.TotalScore < 50 {
		recs = append(recs, re.createRecommendation("overall-seo-audit", map[string]interface{}{
			"current_score": fmt.Sprintf("%d", analysis.TotalScore),
			"target":        "70+",
		}))
	}

	// TODO: Re-enable CategoryScores after field mapping
	// CategoryScores field doesn't exist in RealPageAnalysis yet
	// for category, score := range analysis.CategoryScores {
	//     if score < 0.5 {
	//         recs = append(recs, re.createRecommendation(...))
	//     }
	// }

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
		resource := RecommendationResource{
			URL:  resourceTemplate,
			Type: "documentation",
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
	re.priorityRules["missing-title"] = 90
	re.priorityRules["missing-meta-desc"] = 80
	re.priorityRules["migrate-https"] = 85
	re.priorityRules["missing-viewport"] = 70
	re.priorityRules["improve-lcp"] = 75
	re.priorityRules["missing-h1"] = 80
	// More rules can be added here
}

// initRecommendationTemplates initialise les templates de recommandations
func (re *RecommendationEngine) initRecommendationTemplates() {
	// Template: Titre manquant
	re.templates["missing-title"] = RecommendationTemplate{
		ID:          "missing-title",
		Title:       "Add Page Title",
		Description: "This page is missing a title tag. Add one to improve SEO.",
		Category:    "Tags",
		Priority:    PriorityCritical,
		Impact:      ImpactHigh,
		Effort:      EffortLow,
		Actions: []string{
			"Add a <title> tag to the HTML head",
			"Include relevant keywords in the title",
			"Keep title between 30-60 characters",
		},
		Resources: []string{
			"https://developers.google.com/search/docs/beginner/seo-starter-guide#create-good-titles",
		},
		Metrics: []string{"CTR", "SERP Position"},
		Tags:    []string{"critical", "tags", "on-page"},
	}

	// Template: Meta description manquante
	re.templates["missing-meta-desc"] = RecommendationTemplate{
		ID:          "missing-meta-desc",
		Title:       "Add Meta Description",
		Description: "This page is missing a meta description. Add one to improve click-through rates.",
		Category:    "Tags",
		Priority:    PriorityHigh,
		Impact:      ImpactHigh,
		Effort:      EffortLow,
		Actions: []string{
			"Add a meta description tag to HTML head",
			"Include a compelling call-to-action",
			"Keep description between 120-160 characters",
		},
		Resources: []string{
			"https://developers.google.com/search/docs/advanced/appearance/snippet",
		},
		Metrics: []string{"CTR", "Impressions"},
		Tags:    []string{"meta", "tags", "ctr"},
	}

	// Template: Migration HTTPS
	re.templates["migrate-https"] = RecommendationTemplate{
		ID:          "migrate-https",
		Title:       "Migrate to HTTPS",
		Description: "This site is not using HTTPS. Migrate to HTTPS for better security and SEO.",
		Category:    "Security",
		Priority:    PriorityCritical,
		Impact:      ImpactHigh,
		Effort:      EffortHigh,
		Actions: []string{
			"Obtain SSL certificate from certificate authority",
			"Configure HTTPS on web server",
			"Set up HTTP to HTTPS redirects",
			"Update all internal links to use HTTPS",
		},
		Resources: []string{
			"https://developers.google.com/web/fundamentals/security/encrypt-in-transit/why-https",
		},
		Metrics: []string{"Trust Signals", "Ranking Boost"},
		Tags:    []string{"critical", "security", "technical"},
	}

	// Template: Core Web Vitals - LCP
	re.templates["improve-lcp"] = RecommendationTemplate{
		ID:          "improve-lcp",
		Title:       "Improve Largest Contentful Paint",
		Description: "The Largest Contentful Paint (LCP) is slower than recommended. Optimize for better Core Web Vitals.",
		Category:    "Performance",
		Priority:    PriorityHigh,
		Impact:      ImpactHigh,
		Effort:      EffortMedium,
		Actions: []string{
			"Optimize and compress images",
			"Improve server response times",
			"Preload important resources",
			"Use a content delivery network (CDN)",
		},
		Resources: []string{
			"https://web.dev/lcp/",
		},
		Metrics: []string{"LCP", "Page Experience", "Core Web Vitals"},
		Tags:    []string{"performance", "core-web-vitals", "lcp"},
	}

	// More templates can be added here
}