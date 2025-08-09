package semantic

import (
	"log"
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// NGramAnalyzer analyseur de n-grammes pour analyse sémantique locale
type NGramAnalyzer struct {
	// Configuration
	minNGramLength  int
	maxNGramLength  int
	minFrequency    int
	stopWords       map[string]bool
	
	// Regex pour nettoyage
	punctuationRegex *regexp.Regexp
	numberRegex      *regexp.Regexp
	whitespaceRegex  *regexp.Regexp
}

// NGramResult résultat de l'analyse des n-grammes
type NGramResult struct {
	NGrams          map[int][]NGram `json:"ngrams"`
	TotalTokens     int             `json:"total_tokens"`
	UniqueTokens    int             `json:"unique_tokens"`
	VocabularySize  int             `json:"vocabulary_size"`
	Entropy         float64         `json:"entropy"`
	Complexity      float64         `json:"complexity"`
}

// TokenFrequency fréquence d'un token
type TokenFrequency struct {
	Token     string
	Frequency int
	Score     float64
}

// NewNGramAnalyzer crée un nouvel analyseur de n-grammes
func NewNGramAnalyzer() *NGramAnalyzer {
	analyzer := &NGramAnalyzer{
		minNGramLength:   1,
		maxNGramLength:   3,
		minFrequency:     2,
		stopWords:        make(map[string]bool),
		punctuationRegex: regexp.MustCompile(`[^\p{L}\p{N}\s]+`),
		numberRegex:      regexp.MustCompile(`^\d+$`),
		whitespaceRegex:  regexp.MustCompile(`\s+`),
	}
	
	// Initialiser les mots vides
	analyzer.initStopWords()
	
	return analyzer
}

// Analyze effectue l'analyse complète des n-grammes
func (nga *NGramAnalyzer) Analyze(text string) map[int][]NGram {
	log.Printf("Début analyse n-grammes - TextLength:%d MaxNGram:%d", len(text), nga.maxNGramLength)

	// Prétraitement du texte
	cleanText := nga.preprocessText(text)
	tokens := nga.tokenize(cleanText)
	
	log.Printf("Tokenisation terminée - TokensCount:%d UniqueTokens:%d", len(tokens), nga.countUniqueTokens(tokens))

	// Génération des n-grammes
	results := make(map[int][]NGram)
	
	for n := nga.minNGramLength; n <= nga.maxNGramLength; n++ {
		ngrams := nga.generateNGrams(tokens, n)
		scoredNGrams := nga.scoreNGrams(ngrams, tokens)
		results[n] = nga.filterAndSortNGrams(scoredNGrams)
		
		log.Printf("N-grammes générés - N:%d Count:%d", n, len(results[n]))
	}

	return results
}

// preprocessText prétraite le texte pour l'analyse
func (nga *NGramAnalyzer) preprocessText(text string) string {
	// Convertir en minuscules
	text = strings.ToLower(text)
	
	// Supprimer la ponctuation excessive
	text = nga.punctuationRegex.ReplaceAllString(text, " ")
	
	// Normaliser les espaces
	text = nga.whitespaceRegex.ReplaceAllString(text, " ")
	
	// Supprimer les espaces en début/fin
	text = strings.TrimSpace(text)
	
	return text
}

// tokenize divise le texte en tokens
func (nga *NGramAnalyzer) tokenize(text string) []string {
	words := strings.Fields(text)
	var tokens []string
	
	for _, word := range words {
		// Filtrer les mots trop courts
		if len(word) < 2 {
			continue
		}
		
		// Filtrer les nombres purs
		if nga.numberRegex.MatchString(word) {
			continue
		}
		
		// Filtrer les mots non-alphabétiques
		if !nga.isValidToken(word) {
			continue
		}
		
		// Filtrer les mots vides (mais pas pour les n-grammes)
		tokens = append(tokens, word)
	}
	
	return tokens
}

// generateNGrams génère les n-grammes de taille n
func (nga *NGramAnalyzer) generateNGrams(tokens []string, n int) map[string]int {
	ngrams := make(map[string]int)
	
	if len(tokens) < n {
		return ngrams
	}
	
	for i := 0; i <= len(tokens)-n; i++ {
		// Créer le n-gramme
		ngramTokens := tokens[i : i+n]
		
		// Filtrer les n-grammes contenant trop de mots vides
		if nga.hasExcessiveStopWords(ngramTokens) {
			continue
		}
		
		ngramText := strings.Join(ngramTokens, " ")
		ngrams[ngramText]++
	}
	
	return ngrams
}

// scoreNGrams calcule les scores des n-grammes
func (nga *NGramAnalyzer) scoreNGrams(ngrams map[string]int, allTokens []string) []NGram {
	totalTokens := len(allTokens)
	var results []NGram
	
	// Calculer les fréquences des tokens individuels pour TF-IDF
	tokenFreq := make(map[string]int)
	for _, token := range allTokens {
		tokenFreq[token]++
	}
	
	for ngramText, frequency := range ngrams {
		if frequency < nga.minFrequency {
			continue
		}
		
		tokens := strings.Fields(ngramText)
		
		// Calculer différents scores
		tfScore := nga.calculateTF(frequency, totalTokens)
		positionScore := nga.calculatePositionScore(ngramText, allTokens)
		lengthScore := nga.calculateLengthScore(len(tokens))
		diversityScore := nga.calculateDiversityScore(tokens, tokenFreq)
		
		// Score composite
		compositeScore := (tfScore * 0.4) + (positionScore * 0.2) + (lengthScore * 0.2) + (diversityScore * 0.2)
		
		ngram := NGram{
			Text:  ngramText,
			Count: frequency,
			Score: compositeScore,
			Type:  nga.getNGramType(len(tokens)),
		}
		
		results = append(results, ngram)
	}
	
	return results
}

// calculateTF calcule le score TF (Term Frequency)
func (nga *NGramAnalyzer) calculateTF(frequency int, totalTokens int) float64 {
	if totalTokens == 0 {
		return 0
	}
	return float64(frequency) / float64(totalTokens)
}

// calculatePositionScore calcule un score basé sur la position dans le texte
func (nga *NGramAnalyzer) calculatePositionScore(ngram string, allTokens []string) float64 {
	// Les n-grammes apparaissant tôt dans le texte ont un score plus élevé
	allText := strings.Join(allTokens, " ")
	position := strings.Index(allText, ngram)
	
	if position == -1 {
		return 0
	}
	
	// Score inversement proportionnel à la position (normalisé)
	maxPosition := len(allText)
	if maxPosition == 0 {
		return 1
	}
	
	return 1.0 - (float64(position) / float64(maxPosition))
}

// calculateLengthScore calcule un score basé sur la longueur du n-gramme
func (nga *NGramAnalyzer) calculateLengthScore(length int) float64 {
	// Les bigrammes et trigrammes ont généralement plus de valeur sémantique
	switch length {
	case 1:
		return 0.5 // Unigrammes moins importants
	case 2:
		return 1.0 // Bigrammes très importants
	case 3:
		return 0.8 // Trigrammes importants mais plus rares
	default:
		return 0.3 // N-grammes plus longs moins fiables
	}
}

// calculateDiversityScore calcule un score de diversité lexicale
func (nga *NGramAnalyzer) calculateDiversityScore(tokens []string, tokenFreq map[string]int) float64 {
	if len(tokens) == 0 {
		return 0
	}
	
	// Score basé sur la rareté moyenne des tokens composant le n-gramme
	totalScore := 0.0
	for _, token := range tokens {
		// Plus un token est rare, plus il contribue au score
		freq := tokenFreq[token]
		if freq > 0 {
			// Score inversement proportionnel à la fréquence
			totalScore += 1.0 / math.Log(float64(freq)+1)
		}
	}
	
	return totalScore / float64(len(tokens))
}

// filterAndSortNGrams filtre et trie les n-grammes par pertinence
func (nga *NGramAnalyzer) filterAndSortNGrams(ngrams []NGram) []NGram {
	// Filtrer les n-grammes de faible qualité
	var filtered []NGram
	for _, ngram := range ngrams {
		if nga.isQualityNGram(ngram) {
			filtered = append(filtered, ngram)
		}
	}
	
	// Trier par score décroissant
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Score > filtered[j].Score
	})
	
	// Limiter le nombre de résultats
	maxResults := nga.getMaxResults(len(filtered))
	if len(filtered) > maxResults {
		filtered = filtered[:maxResults]
	}
	
	return filtered
}

// isQualityNGram vérifie si un n-gramme est de qualité
func (nga *NGramAnalyzer) isQualityNGram(ngram NGram) bool {
	tokens := strings.Fields(ngram.Text)
	
	// Rejeter si trop de mots vides
	if nga.hasExcessiveStopWords(tokens) {
		return false
	}
	
	// Rejeter si score trop faible
	if ngram.Score < 0.01 {
		return false
	}
	
	// Rejeter si contient des caractères indésirables
	if strings.ContainsAny(ngram.Text, "0123456789@#$%^&*()") {
		return false
	}
	
	return true
}

// hasExcessiveStopWords vérifie si le n-gramme contient trop de mots vides
func (nga *NGramAnalyzer) hasExcessiveStopWords(tokens []string) bool {
	if len(tokens) == 0 {
		return true
	}
	
	stopWordCount := 0
	for _, token := range tokens {
		if nga.stopWords[token] {
			stopWordCount++
		}
	}
	
	// Seuils par longueur de n-gramme
	switch len(tokens) {
	case 1:
		return stopWordCount > 0 // Pas de mots vides pour les unigrammes
	case 2:
		return stopWordCount > 1 // Maximum 1 mot vide pour les bigrammes
	case 3:
		return stopWordCount > 1 // Maximum 1 mot vide pour les trigrammes
	default:
		return float64(stopWordCount)/float64(len(tokens)) > 0.5 // Maximum 50%
	}
}

// isValidToken vérifie si un token est valide
func (nga *NGramAnalyzer) isValidToken(token string) bool {
	// Vérifier que le token contient au moins une lettre
	hasLetter := false
	for _, r := range token {
		if unicode.IsLetter(r) {
			hasLetter = true
			break
		}
	}
	
	return hasLetter
}

// countUniqueTokens compte les tokens uniques
func (nga *NGramAnalyzer) countUniqueTokens(tokens []string) int {
	unique := make(map[string]bool)
	for _, token := range tokens {
		unique[token] = true
	}
	return len(unique)
}

// getNGramType retourne le type de n-gramme
func (nga *NGramAnalyzer) getNGramType(length int) string {
	switch length {
	case 1:
		return "unigram"
	case 2:
		return "bigram"
	case 3:
		return "trigram"
	default:
		return "n-gram"
	}
}

// getMaxResults retourne le nombre maximum de résultats à conserver
func (nga *NGramAnalyzer) getMaxResults(totalResults int) int {
	// Adaptation dynamique selon le nombre total
	switch {
	case totalResults <= 20:
		return totalResults
	case totalResults <= 100:
		return 30
	case totalResults <= 500:
		return 50
	default:
		return 100
	}
}

// initStopWords initialise la liste des mots vides
func (nga *NGramAnalyzer) initStopWords() {
	// Mots vides français
	frenchStopWords := []string{
		"le", "de", "et", "à", "un", "il", "être", "et", "en", "avoir", "que", "pour",
		"dans", "ce", "son", "une", "sur", "avec", "ne", "se", "pas", "tout", "plus",
		"par", "grand", "en", "être", "et", "en", "avoir", "que", "pour", "dans",
		"ce", "son", "une", "sur", "avec", "ne", "se", "pas", "tout", "plus", "par",
		"mais", "comme", "faire", "leur", "si", "deux", "peut", "ces", "dont", "très",
		"sans", "nous", "vous", "ils", "elle", "bien", "où", "quand", "comment",
		"pourquoi", "qui", "quoi", "ou", "donc", "car", "ni", "soit",
	}
	
	// Mots vides anglais
	englishStopWords := []string{
		"the", "be", "to", "of", "and", "a", "in", "that", "have", "i", "it", "for",
		"not", "on", "with", "he", "as", "you", "do", "at", "this", "but", "his", "by",
		"from", "they", "she", "or", "an", "will", "my", "one", "all", "would", "there",
		"their", "what", "so", "up", "out", "if", "about", "who", "get", "which", "go",
		"me", "when", "make", "can", "like", "time", "no", "just", "him", "know", "take",
		"people", "into", "year", "your", "good", "some", "could", "them", "see", "other",
		"than", "then", "now", "look", "only", "come", "its", "over", "think", "also",
		"back", "after", "use", "two", "how", "our", "work", "first", "well", "way",
		"even", "new", "want", "because", "any", "these", "give", "day", "most", "us",
	}
	
	// Ajouter tous les mots vides au map
	for _, word := range frenchStopWords {
		nga.stopWords[word] = true
	}
	for _, word := range englishStopWords {
		nga.stopWords[word] = true
	}
	
	// Ajouter des mots vides communs du web
	webStopWords := []string{
		"www", "http", "https", "com", "org", "net", "html", "php", "asp", "jsp",
		"page", "site", "web", "internet", "online", "click", "here", "more", "info",
		"contact", "about", "home", "menu", "navigation", "footer", "header",
	}
	
	for _, word := range webStopWords {
		nga.stopWords[word] = true
	}
}

// CalculateEntropy calcule l'entropie du texte (diversité du vocabulaire)
func (nga *NGramAnalyzer) CalculateEntropy(tokens []string) float64 {
	if len(tokens) == 0 {
		return 0
	}
	
	// Compter les fréquences
	frequencies := make(map[string]int)
	for _, token := range tokens {
		frequencies[token]++
	}
	
	// Calculer l'entropie de Shannon
	entropy := 0.0
	totalTokens := float64(len(tokens))
	
	for _, freq := range frequencies {
		probability := float64(freq) / totalTokens
		if probability > 0 {
			entropy -= probability * math.Log2(probability)
		}
	}
	
	return entropy
}

// CalculateComplexity calcule la complexité linguistique
func (nga *NGramAnalyzer) CalculateComplexity(ngrams map[int][]NGram) float64 {
	// Complexité basée sur la diversité des n-grammes
	totalNGrams := 0
	uniqueNGrams := 0
	
	for _, ngramList := range ngrams {
		totalNGrams += len(ngramList)
		uniqueSet := make(map[string]bool)
		for _, ngram := range ngramList {
			uniqueSet[ngram.Text] = true
		}
		uniqueNGrams += len(uniqueSet)
	}
	
	if totalNGrams == 0 {
		return 0
	}
	
	// Ratio de diversité (0 à 1)
	complexity := float64(uniqueNGrams) / float64(totalNGrams)
	
	return complexity
}