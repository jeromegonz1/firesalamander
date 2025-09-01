package constants

// ========================================
// BRAVO-2 TAG ANALYZER CONSTANTS
// ========================================

// Tag JSON Field Names - Analysis Structure
const (
	TagJSONFieldTitle           = "title"
	TagJSONFieldMetaDescription = "meta_description"
	TagJSONFieldHeadings        = "headings"
	TagJSONFieldMetaTags        = "meta_tags"
	TagJSONFieldImages          = "images"
	TagJSONFieldLinks           = "links"
	TagJSONFieldMicrodata       = "microdata"
)

// Tag JSON Field Names - Common Fields
const (
	TagJSONFieldPresent         = "present"
	TagJSONFieldContent         = "content"
	TagJSONFieldLength          = "length"
	TagJSONFieldOptimalLength   = "optimal_length"
	TagJSONFieldIssues          = "issues"
	TagJSONFieldRecommendations = "recommendations"
)

// Tag JSON Field Names - Title Analysis
const (
	TagJSONFieldHasKeywords    = "has_keywords"
	TagJSONFieldDuplicateWords = "duplicate_words"
)

// Tag JSON Field Names - Meta Description Analysis
const (
	TagJSONFieldHasCallToAction = "has_call_to_action"
)

// Tag JSON Field Names - Heading Analysis
const (
	TagJSONFieldH1Count         = "h1_count"
	TagJSONFieldH1Content       = "h1_content"
	TagJSONFieldHeadingStructure = "heading_structure"
	TagJSONFieldHasHierarchy    = "has_hierarchy"
	TagJSONFieldMissingLevels   = "missing_levels"
)

// Tag JSON Field Names - Meta Tags Analysis
const (
	TagJSONFieldHasRobots      = "has_robots"
	TagJSONFieldRobotsContent  = "robots_content"
	TagJSONFieldHasCanonical   = "has_canonical"
	TagJSONFieldCanonicalURL   = "canonical_url"
	TagJSONFieldHasOGTags      = "has_og_tags"
	TagJSONFieldOGTags         = "og_tags"
	TagJSONFieldHasTwitterCard = "has_twitter_card"
	TagJSONFieldTwitterCard    = "twitter_card"
	TagJSONFieldHasViewport    = "has_viewport"
	TagJSONFieldViewportContent = "viewport_content"
)

// Tag JSON Field Names - Image Analysis
const (
	TagJSONFieldTotalImages      = "total_images"
	TagJSONFieldImagesWithAlt    = "images_with_alt"
	TagJSONFieldAltTextCoverage  = "alt_text_coverage"
	TagJSONFieldMissingAltImages = "missing_alt_images"
	TagJSONFieldOptimizedFormats = "optimized_formats"
	TagJSONFieldLazyLoading      = "lazy_loading"
)

// Tag JSON Field Names - Link Analysis
const (
	TagJSONFieldInternalLinks      = "internal_links"
	TagJSONFieldExternalLinks      = "external_links"
	TagJSONFieldNoFollowLinks      = "nofollow_links"
	TagJSONFieldBrokenLinks        = "broken_links"
	TagJSONFieldAnchorOptimization = "anchor_optimization"
)

// Tag JSON Field Names - Microdata Analysis
const (
	TagJSONFieldHasJSONLD       = "has_json_ld"
	TagJSONFieldJSONLDTypes     = "json_ld_types"
	TagJSONFieldHasMicrodata    = "has_microdata"
	TagJSONFieldMicrodataTypes  = "microdata_types"
	TagJSONFieldHasRDFa         = "has_rdfa"
)

// Tag JSON Field Names - Object Fields
const (
	TagJSONFieldProperty = "property"
	TagJSONFieldName     = "name"
)

// ========================================
// TAG HTML META NAMES
// ========================================

// Meta Tag Names
const (
	TagMetaNameDescription = "description"
	TagMetaNameRobots      = "robots"
	TagMetaNameViewport    = "viewport"
)

// ========================================
// TAG HTML ATTRIBUTES
// ========================================

// HTML Attributes
const (
	TagAttrSrc      = "src"
	TagAttrAlt      = "alt"
	TagAttrLoading  = "loading"
	TagAttrHref     = "href"
	TagAttrRel      = "rel"
	TagAttrType     = "type"
	TagAttrProperty = "property"
	TagAttrContent  = "content"
	TagAttrName     = "name"
)

// ========================================
// TAG HTML VALUES
// ========================================

// HTML Attribute Values
const (
	TagValueLazy       = "lazy"
	TagValueCanonical  = "canonical"
	TagValueNoFollow   = "nofollow"
	TagValueJSONLD     = "application/ld+json"
	TagValueItemScope  = "itemscope"
	TagValueTypeOf     = "typeof"
	TagValueVocab      = "vocab"
)

// HTML Prefixes
const (
	TagPrefixOG      = "og:"
	TagPrefixTwitter = "twitter:"
)

// ========================================
// TAG HEADING LEVELS
// ========================================

// Heading Level Names
const (
	TagHeadingH1 = "h1"
	TagHeadingH2 = "h2"
	TagHeadingH3 = "h3"
	TagHeadingH4 = "h4"
	TagHeadingH5 = "h5"
	TagHeadingH6 = "h6"
)

// ========================================
// TAG SEO THRESHOLDS
// ========================================

// Title Length Thresholds
const (
	TagTitleMinLength = 30
	TagTitleMaxLength = 60
	TagTitleRegexMax  = 200
)

// Meta Description Length Thresholds
const (
	TagMetaDescMinLength = 120
	TagMetaDescMaxLength = 160
	TagMetaDescRegexMax  = 300
)

// Link and Content Thresholds
const (
	TagMinWordLength     = 3
	TagMinInternalLinks  = 3
	TagMinAnchorOptim    = 0.7
	TagFullCoverage      = 1.0
)

// ========================================
// TAG ERROR MESSAGES
// ========================================

// Title Error Messages
const (
	TagMsgTitleMissing   = "Balise title manquante"
	TagMsgTitleTooShort  = "Titre trop court"
	TagMsgTitleTooLong   = "Titre trop long"
	TagMsgTitleDuplicates = "Mots dupliqués dans le titre"
)

// Meta Description Error Messages
const (
	TagMsgMetaDescMissing  = "Meta description manquante"
	TagMsgMetaDescTooShort = "Meta description trop courte"
	TagMsgMetaDescTooLong  = "Meta description trop longue"
)

// Heading Error Messages
const (
	TagMsgNoH1           = "Aucun titre H1"
	TagMsgMultipleH1     = "Plusieurs titres H1"
	TagMsgBadHierarchy   = "Hiérarchie des titres incorrecte"
)

// Meta Tags Error Messages
const (
	TagMsgViewportMissing = "Meta viewport manquante"
)

// Image Error Messages  
const (
	TagMsgImagesNoAlt = "%d images sans texte alternatif"
)

// Link Error Messages
const (
	TagMsgBadAnchors = "Textes d'ancre peu optimisés"
)

// ========================================
// TAG RECOMMENDATIONS
// ========================================

// Title Recommendations
const (
	TagRecommendAddTitle       = "Ajouter une balise <title> descriptive"
	TagRecommendExtendTitle    = "Étendre le titre (30-60 caractères optimal)"
	TagRecommendShortenTitle   = "Raccourcir le titre (risque de troncature)"
	TagRecommendAvoidDuplicates = "Éviter la répétition de mots dans le titre"
)

// Meta Description Recommendations
const (
	TagRecommendAddMetaDesc     = "Ajouter une meta description attrayante"
	TagRecommendExtendMetaDesc  = "Étendre la meta description (120-160 caractères)"
	TagRecommendShortenMetaDesc = "Raccourcir la meta description"
	TagRecommendAddCTA          = "Ajouter un appel à l'action dans la meta description"
)

// Heading Recommendations
const (
	TagRecommendAddH1          = "Ajouter un titre H1 principal"
	TagRecommendSingleH1       = "Utiliser un seul H1 par page"
	TagRecommendRespectHierarchy = "Respecter la hiérarchie H1 > H2 > H3..."
)

// Meta Tags Recommendations
const (
	TagRecommendAddCanonical  = "Ajouter une URL canonique"
	TagRecommendAddViewport   = "Ajouter meta viewport pour le responsive"
	TagRecommendAddOGTags     = "Ajouter les balises Open Graph"
	TagRecommendAddTwitter    = "Ajouter les balises Twitter Card"
)

// Image Recommendations
const (
	TagRecommendAddAltText    = "Ajouter des textes alternatifs à toutes les images"
	TagRecommendLazyLoading   = "Implémenter le lazy loading pour les images"
)

// Link Recommendations
const (
	TagRecommendImproveInternal = "Améliorer le maillage interne"
	TagRecommendOptimizeAnchors = "Optimiser les textes d'ancres des liens"
)

// Microdata Recommendations
const (
	TagRecommendAddStructuredData = "Ajouter des données structurées (JSON-LD recommandé)"
)

// ========================================
// TAG CALL TO ACTION WORDS
// ========================================

// French CTA Words
const (
	TagCTADecouvrir   = "découvrir"
	TagCTASavoirPlus  = "en savoir plus"
	TagCTAContacter   = "contacter"
	TagCTACommander   = "commander"
	TagCTAAcheter     = "acheter"
	TagCTATelecharger = "télécharger"
	TagCTAInscrire    = "s'inscrire"
	TagCTAEssayer     = "essayer"
	TagCTACommencer   = "commencer"
	TagCTACliquer     = "cliquer"
)

// English CTA Words
const (
	TagCTADiscover   = "discover"
	TagCTALearnMore  = "learn more"
	TagCTAContact    = "contact"
	TagCTAOrder      = "order"
	TagCTABuy        = "buy"
	TagCTADownload   = "download"
	TagCTASignUp     = "sign up"
	TagCTATry        = "try"
	TagCTAStart      = "start"
	TagCTAClick      = "click"
)

// ========================================
// TAG BAD ANCHOR TEXTS
// ========================================

// French Bad Anchors
const (
	TagBadAnchorCliquezIci  = "cliquez ici"
	TagBadAnchorLireSuite   = "lire la suite"
	TagBadAnchorIci         = "ici"
)

// English Bad Anchors
const (
	TagBadAnchorClickHere = "click here"
	TagBadAnchorReadMore  = "read more"
	TagBadAnchorHere      = "here"
)

// ========================================
// TAG REGEX PATTERNS
// ========================================

// Regex Pattern Constants
const (
	TagRegexTitlePattern    = "^.{1,200}$"
	TagRegexMetaDescPattern = "^.{1,300}$"
	TagRegexURLPattern      = "^https?://[^\\s/$.?#].[^\\s]*$"
	TagRegexImageExtPattern = "\\.(webp|avif|jpg|jpeg|png|gif|svg)$"
)

// File Extension Pattern
const (
	TagImageExtensions = ".webp|avif|jpg|jpeg|png|gif|svg"
)

// ========================================
// TAG URL PROTOCOLS
// ========================================

// URL Protocols
const (
	TagProtocolHTTP  = "http"
	TagProtocolHTTPS = "https"
	TagProtocolSlash = "/"
)

// Protocol Patterns
const (
	TagPatternHTTPProtocol = "http"
	TagPatternURLScheme    = "://"
)

// ========================================
// TAG DETECTION VALUES
// ========================================

// Microdata Detection Values
const (
	TagDetectedValue = "detected"
)