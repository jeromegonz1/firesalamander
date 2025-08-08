package semantic

import (
	"log"
	"math"
	"regexp"
	"strings"
	
	"firesalamander/internal/constants"
)

// SEOScorer système de scoring SEO intelligent
type SEOScorer struct {
	// Seuils de scoring
	thresholds SEOThresholds
	
	// Regex pour validation
	titleRegex       *regexp.Regexp
	metaDescRegex    *regexp.Regexp
	keywordRegex     *regexp.Regexp
}

// SEOThresholds seuils pour le scoring SEO
type SEOThresholds struct {
	TitleMinLength        int     `json:"title_min_length"`
	TitleMaxLength        int     `json:"title_max_length"`
	MetaDescMinLength     int     `json:"meta_desc_min_length"`
	MetaDescMaxLength     int     `json:"meta_desc_max_length"`
	MinWordCount          int     `json:"min_word_count"`
	OptimalWordCount      int     `json:"optimal_word_count"`
	MinKeywordDensity     float64 `json:"min_keyword_density"`
	MaxKeywordDensity     float64 `json:"max_keyword_density"`
	MinInternalLinks      int     `json:"min_internal_links"`
	MinHeadings           int     `json:"min_headings"`
	MinReadabilityScore   float64 `json:"min_readability_score"`
}

// SEOFactorScore score d'un facteur SEO individuel
type SEOFactorScore struct {
	Name        string  `json:"name"`
	Score       float64 `json:"score"`
	Weight      float64 `json:"weight"`
	Status      string  `json:"status"` // excellent, good, warning, critical
	Message     string  `json:"message"`
	Suggestions []string `json:"suggestions"`
}

// NewSEOScorer crée un nouveau scorer SEO
func NewSEOScorer() *SEOScorer {
	return &SEOScorer{
		thresholds: SEOThresholds{
			TitleMinLength:        30,
			TitleMaxLength:        60,
			MetaDescMinLength:     120,
			MetaDescMaxLength:     160,
			MinWordCount:          300,
			OptimalWordCount:      1000,
			MinKeywordDensity:     0.5,
			MaxKeywordDensity:     3.0,
			MinInternalLinks:      3,
			MinHeadings:           2,
			MinReadabilityScore:   60.0,
		},
		titleRegex:    regexp.MustCompile(`^.{1,200}$`),
		metaDescRegex: regexp.MustCompile(`^.{1,300}$`),
		keywordRegex:  regexp.MustCompile(`\b\w+\b`),
	}
}

// Score calcule le score SEO global
func (seo *SEOScorer) Score(content *ExtractedContent, localAnalysis LocalAnalysis, aiAnalysis *AIAnalysis) SEOScore {
	log.Printf(constants.LogSEOScoringStart, content.Title, content.WordCount, aiAnalysis != nil)

	factors := make(map[string]float64)
	var issues []string
	var recommendations []string
	
	// 1. Score du titre
	titleScore := seo.scoreTitleOptimization(content, localAnalysis)
	factors[constants.SEOFactorTitle] = titleScore.Score
	if titleScore.Status != constants.SEOStatusExcellent && titleScore.Status != constants.SEOStatusGood {
		issues = append(issues, titleScore.Message)
		recommendations = append(recommendations, titleScore.Suggestions...)
	}

	// 2. Score de la meta description
	metaScore := seo.scoreMetaDescription(content, localAnalysis)
	factors[constants.SEOFactorMetaDescription] = metaScore.Score
	if metaScore.Status != constants.SEOStatusExcellent && metaScore.Status != constants.SEOStatusGood {
		issues = append(issues, metaScore.Message)
		recommendations = append(recommendations, metaScore.Suggestions...)
	}

	// 3. Score du contenu
	contentScore := seo.scoreContentQuality(content, localAnalysis)
	factors[constants.SEOFactorContentQuality] = contentScore.Score
	if contentScore.Status != constants.SEOStatusExcellent && contentScore.Status != constants.SEOStatusGood {
		issues = append(issues, contentScore.Message)
		recommendations = append(recommendations, contentScore.Suggestions...)
	}

	// 4. Score des mots-clés
	keywordScore := seo.scoreKeywordOptimization(content, localAnalysis)
	factors[constants.SEOFactorKeywordOptimization] = keywordScore.Score
	if keywordScore.Status != constants.SEOStatusExcellent && keywordScore.Status != constants.SEOStatusGood {
		issues = append(issues, keywordScore.Message)
		recommendations = append(recommendations, keywordScore.Suggestions...)
	}

	// 5. Score de la structure
	structureScore := seo.scoreContentStructure(content, localAnalysis)
	factors[constants.SEOFactorContentStructure] = structureScore.Score
	if structureScore.Status != constants.SEOStatusExcellent && structureScore.Status != constants.SEOStatusGood {
		issues = append(issues, structureScore.Message)
		recommendations = append(recommendations, structureScore.Suggestions...)
	}

	// 6. Score de la lisibilité
	readabilityScore := seo.scoreReadability(content, localAnalysis)
	factors[constants.SEOFactorReadability] = readabilityScore.Score
	if readabilityScore.Status != constants.SEOStatusExcellent && readabilityScore.Status != constants.SEOStatusGood {
		issues = append(issues, readabilityScore.Message)
		recommendations = append(recommendations, readabilityScore.Suggestions...)
	}

	// 7. Score des liens
	linkScore := seo.scoreLinkOptimization(content, localAnalysis)
	factors[constants.SEOFactorLinkOptimization] = linkScore.Score
	if linkScore.Status != constants.SEOStatusExcellent && linkScore.Status != constants.SEOStatusGood {
		issues = append(issues, linkScore.Message)
		recommendations = append(recommendations, linkScore.Suggestions...)
	}

	// 8. Score des images
	imageScore := seo.scoreImageOptimization(content, localAnalysis)
	factors[constants.SEOFactorImageOptimization] = imageScore.Score
	if imageScore.Status != constants.SEOStatusExcellent && imageScore.Status != constants.SEOStatusGood {
		issues = append(issues, imageScore.Message)
		recommendations = append(recommendations, imageScore.Suggestions...)
	}

	// 9. Bonus IA (si disponible)
	aiBonus := 0.0
	if aiAnalysis != nil {
		aiBonus = seo.calculateAIBonus(aiAnalysis)
		factors[constants.SEOFactorAIEnrichment] = aiBonus
		
		// Ajouter les recommandations IA
		recommendations = append(recommendations, aiAnalysis.Recommendations...)
	}

	// Calcul du score global pondéré
	overallScore := seo.calculateWeightedScore(factors, aiBonus)

	log.Printf(constants.LogSEOScoringComplete, overallScore, len(factors), len(issues), len(recommendations))

	return SEOScore{
		Overall:         overallScore,
		Factors:         factors,
		Recommendations: recommendations,
		Issues:          issues,
	}
}

// scoreTitleOptimization évalue l'optimisation du titre
func (seo *SEOScorer) scoreTitleOptimization(content *ExtractedContent, analysis LocalAnalysis) SEOFactorScore {
	score := 0.0
	status := constants.SEOStatusCritical
	message := ""
	suggestions := []string{}

	if content.Title == "" {
		message = constants.MsgTitleMissing
		suggestions = append(suggestions, constants.SuggAddDescriptiveTitle)
	} else {
		titleLength := len(content.Title)
		
		// Score basé sur la longueur
		if titleLength < seo.thresholds.TitleMinLength {
			score = 0.3
			status = constants.SEOStatusWarning
			message = constants.MsgTitleTooShort
			suggestions = append(suggestions, constants.SuggExtendTitle)
		} else if titleLength > seo.thresholds.TitleMaxLength {
			score = 0.6
			status = constants.SEOStatusWarning
			message = constants.MsgTitleTooLong
			suggestions = append(suggestions, constants.SuggShortenTitle)
		} else {
			score = 0.8
			status = constants.SEOStatusGood
			message = constants.MsgTitleOptimalLength
		}

		// Bonus pour présence de mots-clés dans le titre
		keywordsInTitle := 0
		for _, keyword := range analysis.Keywords {
			if keyword.InTitle && keyword.Relevance > 0.3 {
				keywordsInTitle++
				score += 0.05
			}
		}

		if keywordsInTitle == 0 {
			suggestions = append(suggestions, constants.SuggIncludeKeywordsTitle)
		} else if keywordsInTitle >= 2 {
			score += 0.1
			if status == constants.SEOStatusGood {
				status = constants.SEOStatusExcellent
			}
		}

		// Pénalité pour sur-optimisation
		if keywordsInTitle > 4 {
			score -= 0.2
			suggestions = append(suggestions, constants.SuggAvoidTitleOverOpt)
		}
	}

	return SEOFactorScore{
		Name:        constants.SEOFactorTitleOptimization,
		Score:       math.Min(1.0, score),
		Weight:      0.20, // 20% du score total
		Status:      status,
		Message:     message,
		Suggestions: suggestions,
	}
}

// scoreMetaDescription évalue la meta description
func (seo *SEOScorer) scoreMetaDescription(content *ExtractedContent, analysis LocalAnalysis) SEOFactorScore {
	score := 0.0
	status := constants.SEOStatusCritical
	message := ""
	suggestions := []string{}

	if content.MetaDescription == "" {
		message = constants.MsgMetaDescMissing
		suggestions = append(suggestions, constants.SuggAddMetaDescription)
	} else {
		descLength := len(content.MetaDescription)
		
		// Score basé sur la longueur
		if descLength < seo.thresholds.MetaDescMinLength {
			score = 0.4
			status = constants.SEOStatusWarning
			message = constants.MsgMetaDescTooShort
			suggestions = append(suggestions, constants.SuggExpandMetaDesc)
		} else if descLength > seo.thresholds.MetaDescMaxLength {
			score = 0.6
			status = constants.SEOStatusWarning
			message = constants.MsgMetaDescTooLong
			suggestions = append(suggestions, constants.SuggShortenMetaDesc)
		} else {
			score = 0.8
			status = constants.SEOStatusGood
			message = constants.MsgMetaDescOptimalLength
		}

		// Bonus pour présence de mots-clés
		keywordsInMeta := 0
		for _, keyword := range analysis.Keywords {
			if keyword.InMeta && keyword.Relevance > 0.3 {
				keywordsInMeta++
				score += 0.05
			}
		}

		if keywordsInMeta == 0 {
			suggestions = append(suggestions, constants.SuggIncludeKeywordsMetaDesc)
		} else if keywordsInMeta >= 2 {
			score += 0.1
			if status == constants.SEOStatusGood {
				status = constants.SEOStatusExcellent
			}
		}

		// Bonus pour appel à l'action
		if seo.hasCallToAction(content.MetaDescription) {
			score += 0.05
		} else {
			suggestions = append(suggestions, constants.SuggAddCallToAction)
		}
	}

	return SEOFactorScore{
		Name:        constants.SEOFactorMetaDescription,
		Score:       math.Min(1.0, score),
		Weight:      0.15, // 15% du score total
		Status:      status,
		Message:     message,
		Suggestions: suggestions,
	}
}

// scoreContentQuality évalue la qualité du contenu
func (seo *SEOScorer) scoreContentQuality(content *ExtractedContent, analysis LocalAnalysis) SEOFactorScore {
	score := 0.0
	status := constants.SEOStatusCritical
	message := ""
	suggestions := []string{}

	wordCount := content.WordCount

	// Score basé sur la longueur du contenu
	if wordCount < seo.thresholds.MinWordCount {
		score = 0.3
		status = constants.SEOStatusWarning
		message = constants.MsgContentTooShort
		suggestions = append(suggestions, constants.SuggExpandContent)
	} else if wordCount >= seo.thresholds.OptimalWordCount {
		score = 1.0
		status = constants.SEOStatusExcellent
		message = constants.MsgContentOptimalLength
	} else {
		// Score proportionnel entre min et optimal
		ratio := float64(wordCount-seo.thresholds.MinWordCount) / float64(seo.thresholds.OptimalWordCount-seo.thresholds.MinWordCount)
		score = 0.5 + (ratio * 0.5)
		status = constants.SEOStatusGood
		message = constants.MsgContentCorrectLength
	}

	// Bonus pour diversité lexicale
	if analysis.Statistics.LexicalDiversity > 0.5 {
		score += 0.1
		if status != constants.SEOStatusExcellent {
			status = constants.SEOStatusGood
		}
	} else if analysis.Statistics.LexicalDiversity < 0.3 {
		suggestions = append(suggestions, constants.SuggImproveVocabulary)
	}

	// Pénalité pour contenu dupliqué détecté (heuristique simple)
	if seo.detectDuplicateContent(content) {
		score -= 0.2
		suggestions = append(suggestions, constants.SuggAvoidDuplicateContent)
	}

	return SEOFactorScore{
		Name:        constants.SEOFactorContentQuality,
		Score:       math.Min(1.0, math.Max(0.0, score)),
		Weight:      0.15, // 15% du score total
		Status:      status,
		Message:     message,
		Suggestions: suggestions,
	}
}

// scoreKeywordOptimization évalue l'optimisation des mots-clés
func (seo *SEOScorer) scoreKeywordOptimization(content *ExtractedContent, analysis LocalAnalysis) SEOFactorScore {
	score := 0.0
	status := constants.SEOStatusCritical
	message := ""
	suggestions := []string{}

	if len(analysis.Keywords) == 0 {
		message = constants.MsgNoKeywords
		suggestions = append(suggestions, constants.SuggAddRelevantKeywords)
	} else {
		// Score basé sur le nombre de mots-clés pertinents
		relevantKeywords := 0
		totalDensity := 0.0

		for _, keyword := range analysis.Keywords {
			if keyword.Relevance > 0.3 {
				relevantKeywords++
				totalDensity += keyword.Density
			}
		}

		if relevantKeywords < 5 {
			score = 0.4
			status = constants.SEOStatusWarning
			message = constants.MsgFewKeywords
			suggestions = append(suggestions, constants.SuggEnrichWithKeywords)
		} else if relevantKeywords >= 10 {
			score = 0.8
			status = constants.SEOStatusGood
			message = constants.MsgGoodKeywordCoverage
		} else {
			score = 0.6
			status = constants.SEOStatusGood
			message = constants.MsgCorrectKeywordCoverage
		}

		// Vérifier la densité moyenne
		avgDensity := totalDensity / float64(relevantKeywords)
		if avgDensity < seo.thresholds.MinKeywordDensity {
			suggestions = append(suggestions, constants.SuggIncreaseKeywordDensity)
		} else if avgDensity > seo.thresholds.MaxKeywordDensity {
			score -= 0.2
			suggestions = append(suggestions, constants.SuggAvoidKeywordOverOpt)
		} else {
			score += 0.1
			if status == constants.SEOStatusGood && score >= 0.9 {
				status = constants.SEOStatusExcellent
			}
		}

		// Bonus pour mots-clés en position stratégique
		strategicKeywords := 0
		for _, keyword := range analysis.Keywords[:min(5, len(analysis.Keywords))] {
			if keyword.InTitle || keyword.InMeta || keyword.InHeadings {
				strategicKeywords++
			}
		}

		if strategicKeywords >= 3 {
			score += 0.1
		} else {
			suggestions = append(suggestions, constants.SuggStrategicKeywordPlace)
		}
	}

	return SEOFactorScore{
		Name:        constants.SEOFactorKeywordOptimization,
		Score:       math.Min(1.0, math.Max(0.0, score)),
		Weight:      0.15, // 15% du score total
		Status:      status,
		Message:     message,
		Suggestions: suggestions,
	}
}

// scoreContentStructure évalue la structure du contenu
func (seo *SEOScorer) scoreContentStructure(content *ExtractedContent, analysis LocalAnalysis) SEOFactorScore {
	score := 0.0
	status := constants.SEOStatusCritical
	message := ""
	suggestions := []string{}

	headingCount := len(content.Headings)

	// Score basé sur le nombre de headings
	if headingCount == 0 {
		message = constants.MsgNoSectionTitles
		suggestions = append(suggestions, constants.SuggStructureWithHeadings)
	} else if headingCount < seo.thresholds.MinHeadings {
		score = 0.4
		status = constants.SEOStatusWarning
		message = constants.MsgFewSectionTitles
		suggestions = append(suggestions, constants.SuggImproveWithSubtitles)
	} else {
		score = 0.7
		status = constants.SEOStatusGood
		message = constants.MsgCorrectContentStructure
	}

	// Vérifier la hiérarchie des headings
	if seo.hasGoodHeadingHierarchy(content.HeadingStructure) {
		score += 0.2
		if status == constants.SEOStatusGood {
			status = constants.SEOStatusExcellent
		}
	} else {
		suggestions = append(suggestions, constants.SuggRespectHeadingHierarchy)
	}

	// Bonus pour présence de listes
	if len(content.Lists) > 0 {
		score += 0.1
	} else {
		suggestions = append(suggestions, constants.SuggUseLists)
	}

	// Vérifier la longueur des paragraphes
	if analysis.Statistics.AvgSentencesPerParagraph > 10 {
		suggestions = append(suggestions, constants.SuggShortenParagraphs)
	} else if analysis.Statistics.AvgSentencesPerParagraph < 2 {
		suggestions = append(suggestions, constants.SuggDevelopParagraphs)
	} else {
		score += 0.05
	}

	return SEOFactorScore{
		Name:        constants.SEOFactorContentStructure,
		Score:       math.Min(1.0, math.Max(0.0, score)),
		Weight:      0.10, // 10% du score total
		Status:      status,
		Message:     message,
		Suggestions: suggestions,
	}
}

// scoreReadability évalue la lisibilité
func (seo *SEOScorer) scoreReadability(content *ExtractedContent, analysis LocalAnalysis) SEOFactorScore {
	score := analysis.ReadabilityScore / 100.0
	status := constants.SEOStatusCritical
	message := ""
	suggestions := []string{}

	if analysis.ReadabilityScore < 30 {
		status = constants.SEOStatusCritical
		message = constants.MsgVeryLowReadability
		suggestions = append(suggestions, constants.SuggSimplifySentences)
	} else if analysis.ReadabilityScore < seo.thresholds.MinReadabilityScore {
		status = constants.SEOStatusWarning
		message = constants.MsgLowReadability
		suggestions = append(suggestions, constants.SuggImproveReadability)
	} else if analysis.ReadabilityScore < constants.HighQualityScore {
		status = constants.SEOStatusGood
		message = constants.MsgCorrectReadability
	} else {
		status = constants.SEOStatusExcellent
		message = constants.MsgExcellentReadability
	}

	// Vérifier la longueur des phrases
	if analysis.Statistics.AvgWordsPerSentence > 25 {
		score -= 0.1
		suggestions = append(suggestions, constants.SuggShortenSentences)
	} else if analysis.Statistics.AvgWordsPerSentence < 10 {
		suggestions = append(suggestions, constants.SuggVarySentenceLength)
	}

	return SEOFactorScore{
		Name:        constants.SEOFactorReadability,
		Score:       math.Min(1.0, math.Max(0.0, score)),
		Weight:      0.10, // 10% du score total
		Status:      status,
		Message:     message,
		Suggestions: suggestions,
	}
}

// scoreLinkOptimization évalue l'optimisation des liens
func (seo *SEOScorer) scoreLinkOptimization(content *ExtractedContent, analysis LocalAnalysis) SEOFactorScore {
	score := 0.0
	status := constants.SEOStatusWarning
	message := ""
	suggestions := []string{}

	internalLinks := 0
	externalLinks := 0

	for _, link := range content.Links {
		if link.IsInternal {
			internalLinks++
		} else if link.IsExternal {
			externalLinks++
		}
	}

	// Score basé sur les liens internes
	if internalLinks == 0 {
		message = constants.MsgNoInternalLinks
		suggestions = append(suggestions, constants.SuggAddInternalLinks)
	} else if internalLinks < seo.thresholds.MinInternalLinks {
		score = 0.4
		message = constants.MsgFewInternalLinks
		suggestions = append(suggestions, constants.SuggIncreaseInternalLinks)
	} else {
		score = 0.7
		status = constants.SEOStatusGood
		message = constants.MsgGoodInternalLinks
	}

	// Bonus pour liens externes de qualité
	if externalLinks > 0 && externalLinks <= 3 {
		score += 0.1
	} else if externalLinks > 5 {
		suggestions = append(suggestions, constants.SuggLimitExternalLinks)
	}

	// Vérifier les ancres de liens
	if seo.hasGoodLinkAnchors(content.Links) {
		score += 0.2
		if status == constants.SEOStatusGood {
			status = constants.SEOStatusExcellent
		}
	} else {
		suggestions = append(suggestions, constants.SuggOptimizeAnchorTexts)
	}

	return SEOFactorScore{
		Name:        constants.SEOFactorLinkOptimization,
		Score:       math.Min(1.0, math.Max(0.0, score)),
		Weight:      0.05, // 5% du score total
		Status:      status,
		Message:     message,
		Suggestions: suggestions,
	}
}

// scoreImageOptimization évalue l'optimisation des images
func (seo *SEOScorer) scoreImageOptimization(content *ExtractedContent, analysis LocalAnalysis) SEOFactorScore {
	score := 1.0 // Par défaut excellent si pas d'images
	status := constants.SEOStatusExcellent
	message := constants.MsgNoImagesToOptimize
	suggestions := []string{}

	if len(content.Images) > 0 {
		score = 0.0
		imagesWithAlt := 0

		for _, image := range content.Images {
			if image.Alt != "" {
				imagesWithAlt++
			}
		}

		// Score basé sur le pourcentage d'images avec alt
		altRatio := float64(imagesWithAlt) / float64(len(content.Images))
		score = altRatio

		if altRatio == 0 {
			status = constants.SEOStatusCritical
			message = constants.MsgNoImageAltText
			suggestions = append(suggestions, constants.SuggAddAltTexts)
		} else if altRatio < 0.5 {
			status = constants.SEOStatusWarning
			message = constants.MsgFewImageAltText
			suggestions = append(suggestions, constants.SuggCompleteMissingAlt)
		} else if altRatio < 1.0 {
			status = constants.SEOStatusGood
			message = constants.MsgMostImagesHaveAlt
			suggestions = append(suggestions, constants.SuggCompleteLastAltTexts)
		} else {
			status = constants.SEOStatusExcellent
			message = constants.MsgAllImagesHaveAlt
		}
	}

	return SEOFactorScore{
		Name:        constants.SEOFactorImageOptimization,
		Score:       score,
		Weight:      0.05, // 5% du score total
		Status:      status,
		Message:     message,
		Suggestions: suggestions,
	}
}

// calculateAIBonus calcule le bonus IA
func (seo *SEOScorer) calculateAIBonus(aiAnalysis *AIAnalysis) float64 {
	bonus := 0.0

	// Bonus pour intention claire
	if aiAnalysis.Intent != "" {
		bonus += 0.02
	}

	// Bonus pour topics bien définis
	if len(aiAnalysis.MainTopics) >= 3 {
		bonus += 0.03
	}

	// Bonus pour audience cible identifiée
	if aiAnalysis.TargetAudience != "" {
		bonus += 0.02
	}

	// Bonus pour recommandations actionables
	if len(aiAnalysis.Recommendations) >= 3 {
		bonus += 0.03
	}

	return bonus
}

// calculateWeightedScore calcule le score pondéré final
func (seo *SEOScorer) calculateWeightedScore(factors map[string]float64, aiBonus float64) float64 {
	// Poids des facteurs (total = 1.0)
	weights := map[string]float64{
		constants.SEOFactorTitle:                0.20,
		constants.SEOFactorMetaDescription:     0.15,
		constants.SEOFactorContentQuality:      0.15,
		constants.SEOFactorKeywordOptimization: 0.15,
		constants.SEOFactorContentStructure:    0.10,
		constants.SEOFactorReadability:          0.10,
		constants.SEOFactorLinkOptimization:    0.10,
		constants.SEOFactorImageOptimization:   0.05,
	}

	totalScore := 0.0
	for factor, score := range factors {
		if weight, exists := weights[factor]; exists {
			totalScore += score * weight
		}
	}

	// Ajouter le bonus IA
	totalScore += aiBonus

	// Convertir en score sur 100
	return math.Min(100.0, totalScore*100.0)
}

// Fonctions utilitaires

func (seo *SEOScorer) hasCallToAction(text string) bool {
	ctas := []string{
		constants.CTADiscover, constants.CTALearnMore, constants.CTAContact, constants.CTAOrder, constants.CTABuy,
		constants.CTADownload, constants.CTASignUp, constants.CTATry, constants.CTAStart, constants.CTAClick,
		constants.CTADiscoverEN, constants.CTALearnMoreEN, constants.CTAContactEN, constants.CTAOrderEN, constants.CTABuyEN, constants.CTADownloadEN,
		constants.CTASignUpEN, constants.CTATryEN, constants.CTAStartEN, constants.CTAClickEN,
	}

	lowerText := strings.ToLower(text)
	for _, cta := range ctas {
		if strings.Contains(lowerText, cta) {
			return true
		}
	}
	return false
}

func (seo *SEOScorer) detectDuplicateContent(content *ExtractedContent) bool {
	// Heuristique simple : détecter les répétitions excessives
	words := strings.Fields(content.CleanText)
	if len(words) < 50 {
		return false
	}

	wordCount := make(map[string]int)
	for _, word := range words {
		wordCount[strings.ToLower(word)]++
	}

	// Si plus de 30% des mots sont répétés plus de 5 fois
	repeatedWords := 0
	for _, count := range wordCount {
		if count > 5 {
			repeatedWords += count
		}
	}

	return float64(repeatedWords)/float64(len(words)) > 0.3
}

func (seo *SEOScorer) hasGoodHeadingHierarchy(headings map[string][]string) bool {
	h1Count := len(headings[constants.HeadingH1])
	h2Count := len(headings[constants.HeadingH2])

	// Il devrait y avoir exactement un H1 et au moins un H2
	return h1Count == 1 && h2Count >= 1
}

func (seo *SEOScorer) hasGoodLinkAnchors(links []Link) bool {
	if len(links) == 0 {
		return true
	}

	goodAnchors := 0
	badAnchors := []string{constants.BadAnchorClickHere, constants.BadAnchorClickHereEN, constants.BadAnchorReadMore, constants.BadAnchorReadMoreEN, constants.BadAnchorHere, constants.BadAnchorHereEN}

	for _, link := range links {
		isGood := true
		lowerText := strings.ToLower(link.Text)
		
		for _, badAnchor := range badAnchors {
			if strings.Contains(lowerText, badAnchor) {
				isGood = false
				break
			}
		}
		
		if isGood && len(link.Text) > 3 {
			goodAnchors++
		}
	}

	return float64(goodAnchors)/float64(len(links)) > 0.7
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}