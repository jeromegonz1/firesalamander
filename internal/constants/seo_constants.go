package constants

// SEO Scoring Constants
// Constantes pour le système de scoring SEO

// SEO Meta Names
const (
	SEOMetaNameDescription = "description"
	SEOMetaNameViewport    = "viewport"
	SEOMetaNameCanonical   = "canonical"
	SEOMetaNameOGTitle     = "og:title"
	SEOMetaNameOGDescription = "og:description"
	SEOMetaNameRobots      = "robots"
)

// SEO Categories
const (
	SEOCategoryPerformance = "performance"
	SEOCategoryTechnical   = "technical"
)

// SEO Recommendation Levels
const (
	SEORecommendationNeedsImprovement = "needs-improvement"
)

// SEO Factor Names
const (
	SEOFactorTitle                = "title"
	SEOFactorMetaDescription      = "meta_description"
	SEOFactorContentQuality       = "content_quality"
	SEOFactorKeywordOptimization  = "keyword_optimization"
	SEOFactorContentStructure     = "content_structure"
	SEOFactorReadability          = "readability"
	SEOFactorLinkOptimization     = "link_optimization"
	SEOFactorImageOptimization    = "image_optimization"
	SEOFactorAIEnrichment         = "ai_enrichment"
)

// SEO Detailed Factor Names  
const (
	SEOFactorTitleOptimization = "title_optimization"
	SEOFactorMetaDesc          = "meta_description"
	SEOFactorContentQual       = "content_quality"
	SEOFactorKeywordOpt        = "keyword_optimization"
	SEOFactorContentStruct     = "content_structure"
	SEOFactorReadabilityScore  = "readability"
	SEOFactorLinkOpt           = "link_optimization"
	SEOFactorImageOpt          = "image_optimization"
)

// SEO Status Values (additional to main constants)
const (
	SEOStatusExcellent = "excellent"
	SEOStatusGood      = "good"
	SEOStatusWarning   = "warning"
	SEOStatusCritical  = "critical"
)

// SEO Messages - Title Optimization
const (
	MsgTitleMissing                = "Titre manquant"
	MsgTitleTooShort              = "Titre trop court"
	MsgTitleTooLong               = "Titre trop long"
	MsgTitleOptimalLength         = "Longueur du titre optimale"
)

// SEO Messages - Meta Description
const (
	MsgMetaDescMissing            = "Meta description manquante"
	MsgMetaDescTooShort          = "Meta description trop courte"
	MsgMetaDescTooLong           = "Meta description trop longue"
	MsgMetaDescOptimalLength     = "Longueur de meta description optimale"
)

// SEO Messages - Content Quality
const (
	MsgContentTooShort           = "Contenu trop court"
	MsgContentOptimalLength      = "Longueur de contenu optimale"
	MsgContentCorrectLength      = "Longueur de contenu correcte"
)

// SEO Messages - Keyword Optimization
const (
	MsgNoKeywords                = "Aucun mot-clé identifié"
	MsgFewKeywords              = "Peu de mots-clés pertinents"
	MsgGoodKeywordCoverage      = "Bonne couverture de mots-clés"
	MsgCorrectKeywordCoverage   = "Couverture de mots-clés correcte"
)

// SEO Messages - Content Structure
const (
	MsgNoSectionTitles          = "Aucun titre de section"
	MsgFewSectionTitles         = "Peu de titres de section"
	MsgCorrectContentStructure  = "Structure de contenu correcte"
)

// SEO Messages - Readability
const (
	MsgVeryLowReadability       = "Lisibilité très faible"
	MsgLowReadability           = "Lisibilité faible"
	MsgCorrectReadability       = "Lisibilité correcte"
	MsgExcellentReadability     = "Excellente lisibilité"
)

// SEO Messages - Link Optimization
const (
	MsgNoInternalLinks          = "Aucun lien interne"
	MsgFewInternalLinks         = "Peu de liens internes"
	MsgGoodInternalLinks        = "Bon maillage interne"
)

// SEO Messages - Image Optimization
const (
	MsgNoImagesToOptimize       = "Pas d'images à optimiser"
	MsgNoImageAltText           = "Aucune image n'a de texte alternatif"
	MsgFewImageAltText          = "Peu d'images ont un texte alternatif"
	MsgMostImagesHaveAlt        = "La plupart des images ont un texte alternatif"
	MsgAllImagesHaveAlt         = "Toutes les images ont un texte alternatif"
)

// SEO Suggestions - Title
const (
	SuggAddDescriptiveTitle     = "Ajouter un titre H1 descriptif"
	SuggExtendTitle             = "Allonger le titre (30-60 caractères optimal)"
	SuggShortenTitle            = "Raccourcir le titre (risque de troncature)"
	SuggIncludeKeywordsTitle    = "Inclure des mots-clés pertinents dans le titre"
	SuggAvoidTitleOverOpt       = "Éviter la sur-optimisation du titre"
)

// SEO Suggestions - Meta Description
const (
	SuggAddMetaDescription      = "Ajouter une meta description attrayante"
	SuggExpandMetaDesc          = "Étoffer la meta description (120-160 caractères)"
	SuggShortenMetaDesc         = "Raccourcir la meta description"
	SuggIncludeKeywordsMetaDesc = "Inclure des mots-clés dans la meta description"
	SuggAddCallToAction         = "Ajouter un appel à l'action dans la meta description"
)

// SEO Suggestions - Content
const (
	SuggExpandContent           = "Étoffer le contenu (minimum 300 mots)"
	SuggImproveVocabulary       = "Améliorer la diversité du vocabulaire"
	SuggAvoidDuplicateContent   = "Éviter le contenu dupliqué"
	SuggAddRelevantKeywords     = "Ajouter des mots-clés pertinents au contenu"
	SuggEnrichWithKeywords      = "Enrichir le contenu avec plus de mots-clés"
	SuggIncreaseKeywordDensity  = "Augmenter la densité des mots-clés principaux"
	SuggAvoidKeywordOverOpt     = "Éviter la sur-optimisation (densité trop élevée)"
	SuggStrategicKeywordPlace   = "Placer des mots-clés dans les positions stratégiques"
)

// SEO Suggestions - Structure
const (
	SuggStructureWithHeadings   = "Structurer le contenu avec des titres H2, H3"
	SuggImproveWithSubtitles    = "Améliorer la structure avec plus de sous-titres"
	SuggRespectHeadingHierarchy = "Respecter la hiérarchie H1 > H2 > H3"
	SuggUseLists                = "Utiliser des listes pour améliorer la lisibilité"
	SuggShortenParagraphs       = "Raccourcir les paragraphes pour améliorer la lisibilité"
	SuggDevelopParagraphs       = "Développer davantage les paragraphes"
)

// SEO Suggestions - Readability
const (
	SuggSimplifySentences       = "Simplifier les phrases et le vocabulaire"
	SuggImproveReadability      = "Améliorer la lisibilité du contenu"
	SuggShortenSentences        = "Raccourcir les phrases (max 20-25 mots)"
	SuggVarySentenceLength      = "Varier la longueur des phrases"
)

// SEO Suggestions - Links
const (
	SuggAddInternalLinks        = "Ajouter des liens internes vers d'autres pages"
	SuggIncreaseInternalLinks   = "Augmenter le maillage interne"
	SuggLimitExternalLinks      = "Limiter le nombre de liens externes"
	SuggOptimizeAnchorTexts     = "Optimiser les textes d'ancres des liens"
)

// SEO Suggestions - Images
const (
	SuggAddAltTexts             = "Ajouter des textes alternatifs à toutes les images"
	SuggCompleteMissingAlt      = "Compléter les textes alternatifs manquants"
	SuggCompleteLastAltTexts    = "Compléter les derniers textes alternatifs"
)

// SEO Weights for scoring factors
const (
	WeightTitle                 = 0.20
	WeightMetaDescription       = 0.15
	WeightContentQuality        = 0.15
	WeightKeywordOptimization   = 0.15
	WeightContentStructure      = 0.10
	WeightReadability           = 0.10
	WeightLinkOptimization      = 0.10
	WeightImageOptimization     = 0.05
)

// Call-to-Action Terms (French & English)
const (
	CTADiscover    = "découvrir"
	CTALearnMore   = "en savoir plus"
	CTAContact     = "contacter"
	CTAOrder       = "commander"
	CTABuy         = "acheter"
	CTADownload    = "télécharger"
	CTASignUp      = "s'inscrire"
	CTATry         = "essayer"
	CTAStart       = "commencer"
	CTAClick       = "cliquer"
	CTADiscoverEN  = "discover"
	CTALearnMoreEN = "learn more"
	CTAContactEN   = "contact"
	CTAOrderEN     = "order"
	CTABuyEN       = "buy"
	CTADownloadEN  = "download"
	CTASignUpEN    = "sign up"
	CTATryEN       = "try"
	CTAStartEN     = "start"
	CTAClickEN     = "click"
)

// Bad Anchor Texts (to avoid)
const (
	BadAnchorClickHere   = "cliquez ici"
	BadAnchorClickHereEN = "click here"
	BadAnchorReadMore    = "lire la suite"
	BadAnchorReadMoreEN  = "read more"
	BadAnchorHere        = "ici"
	BadAnchorHereEN      = "here"
)

// Heading Types
const (
	HeadingH1 = "h1"
	HeadingH2 = "h2"
	HeadingH3 = "h3"
)

// Log Messages for SEO Scoring
const (
	LogSEOScoringStart    = "Début scoring SEO - Title:%s WordCount:%d HasAI:%t"
	LogSEOScoringComplete = "Scoring SEO terminé - OverallScore:%.1f FactorsCount:%d IssuesCount:%d RecommendationsCount:%d"
)

// JSON Field Names for SEO Configuration
const (
	JSONFieldTitleMinLength      = "title_min_length"
	JSONFieldTitleMaxLength      = "title_max_length"
	JSONFieldMetaDescMinLength   = "meta_desc_min_length"
	JSONFieldMetaDescMaxLength   = "meta_desc_max_length"
	JSONFieldMinWordCount        = "min_word_count"
	JSONFieldOptimalWordCount    = "optimal_word_count"
	JSONFieldMinKeywordDensity   = "min_keyword_density"
	JSONFieldMaxKeywordDensity   = "max_keyword_density"
	JSONFieldMinInternalLinks    = "min_internal_links"
	JSONFieldMinHeadings         = "min_headings"
	JSONFieldMinReadabilityScore = "min_readability_score"
	JSONFieldNameSEO             = "name"
	JSONFieldScore               = "score"
	JSONFieldWeight              = "weight"
	JSONFieldMessage             = "message"
	JSONFieldSuggestions         = "suggestions"
)