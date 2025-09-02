package keyword

import (
	"strings"
)

// LanguageDetector détecte la langue d'un contenu textuel
type LanguageDetector struct {
	frenchWords map[string]int
	englishWords map[string]int
}

// NewLanguageDetector crée un nouveau détecteur de langue
func NewLanguageDetector() *LanguageDetector {
	frenchWords := map[string]int{
		"le": 10, "la": 10, "les": 10, "de": 10, "du": 8, "des": 8,
		"et": 9, "est": 9, "une": 8, "un": 8, "ce": 7, "cette": 7,
		"pour": 6, "avec": 6, "sur": 6, "dans": 6, "par": 5,
		"mais": 5, "ou": 5, "où": 5, "que": 9, "qui": 8,
		"son": 4, "sa": 4, "ses": 4, "leur": 4, "leurs": 4,
		"nous": 5, "vous": 5, "ils": 5, "elles": 5,
		"avoir": 6, "être": 8, "faire": 5, "aller": 4,
	}

	englishWords := map[string]int{
		"the": 10, "and": 10, "of": 9, "to": 9, "in": 8, "is": 8,
		"you": 7, "that": 7, "it": 7, "he": 6, "was": 6, "for": 6,
		"on": 6, "are": 6, "as": 5, "with": 5, "his": 5, "they": 5,
		"i": 8, "at": 5, "be": 5, "this": 5, "have": 6, "from": 5,
		"or": 4, "one": 4, "had": 4, "by": 4, "word": 3, "but": 4,
		"not": 6, "what": 5, "all": 4, "were": 4, "we": 5,
	}

	return &LanguageDetector{
		frenchWords:  frenchWords,
		englishWords: englishWords,
	}
}

// DetectLanguage détecte la langue prédominante dans le contenu
func (ld *LanguageDetector) DetectLanguage(content string) string {
	if content == "" {
		return "unknown"
	}

	words := strings.Fields(strings.ToLower(content))
	if len(words) < 5 {
		return "unknown"
	}

	frenchScore := 0
	englishScore := 0
	totalWords := 0

	// Analyse seulement les 100 premiers mots pour la performance
	maxWords := len(words)
	if maxWords > 100 {
		maxWords = 100
	}

	for i := 0; i < maxWords; i++ {
		word := words[i]
		totalWords++

		// Vérification des mots français
		if score, exists := ld.frenchWords[word]; exists {
			frenchScore += score
		}

		// Vérification des mots anglais
		if score, exists := ld.englishWords[word]; exists {
			englishScore += score
		}
	}

	// Si pas assez d'indicateurs linguistiques
	if frenchScore == 0 && englishScore == 0 {
		return ld.detectByCharacteristics(content)
	}

	// Détermine la langue dominante
	if frenchScore > englishScore {
		return "fr"
	} else if englishScore > frenchScore {
		return "en"
	}

	// En cas d'égalité, utilise les caractéristiques
	return ld.detectByCharacteristics(content)
}

// detectByCharacteristics détecte la langue par les caractéristiques spécifiques
func (ld *LanguageDetector) detectByCharacteristics(content string) string {
	content = strings.ToLower(content)

	// Caractères spécifiques au français
	frenchChars := []string{"à", "é", "è", "ê", "ë", "î", "ï", "ô", "ù", "û", "ü", "ÿ", "ç"}
	frenchCharCount := 0

	for _, char := range frenchChars {
		frenchCharCount += strings.Count(content, char)
	}

	// Structures typiquement françaises
	frenchStructures := []string{"qu'", "d'", "l'", "c'est", "il y a", "n'est", "s'est"}
	frenchStructureCount := 0

	for _, structure := range frenchStructures {
		frenchStructureCount += strings.Count(content, structure)
	}

	// Structures typiquement anglaises
	englishStructures := []string{"'s ", " 's", "n't", "'re", "'ve", "'ll", "'d", "ing ", "tion ", "ly "}
	englishStructureCount := 0

	for _, structure := range englishStructures {
		englishStructureCount += strings.Count(content, structure)
	}

	totalFrenchIndicators := frenchCharCount + frenchStructureCount
	totalEnglishIndicators := englishStructureCount

	if totalFrenchIndicators > totalEnglishIndicators {
		return "fr"
	} else if totalEnglishIndicators > totalFrenchIndicators {
		return "en"
	}

	// Par défaut, considère comme anglais si indéterminé
	return "en"
}

// GetLanguageName retourne le nom complet de la langue
func (ld *LanguageDetector) GetLanguageName(code string) string {
	switch code {
	case "fr":
		return "French"
	case "en":
		return "English"
	default:
		return "Unknown"
	}
}

// IsSupported vérifie si une langue est supportée
func (ld *LanguageDetector) IsSupported(code string) bool {
	return code == "fr" || code == "en"
}