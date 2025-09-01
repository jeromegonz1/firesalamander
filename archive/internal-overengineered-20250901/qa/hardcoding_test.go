package qa

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestNoHardcoding - Test critique qui échoue si du hardcoding est détecté
func TestNoHardcoding(t *testing.T) {
	// Patterns de hardcoding à détecter
	patterns := []struct {
		name    string
		regex   *regexp.Regexp
		exclude []string
	}{
		{
			name:  "Durées hardcodées",
			regex: regexp.MustCompile(`\d+\s*\*\s*time\.(Millisecond|Second|Minute|Hour)`),
			exclude: []string{"_test.go"},
		},
		{
			name:  "Pourcentages hardcodés",
			regex: regexp.MustCompile(`\b(0|[1-9]\d?|100)\s*%|:\s*(0|[1-9]\d?|100)\s*,`),
			exclude: []string{"_test.go", ".md"},
		},
		{
			name:  "Nombres magiques",
			regex: regexp.MustCompile(`[\+\-\*/=]\s*\d+\.?\d*\s*[;\),]|return\s+\d+`),
			exclude: []string{"_test.go", "config.go"},
		},
		{
			name:  "Messages hardcodés",
			regex: regexp.MustCompile(`"[A-Z][^"]{10,}"|'[A-Z][^']{10,}'`),
			exclude: []string{"_test.go", ".md"},
		},
		{
			name:  "Chemins hardcodés",
			regex: regexp.MustCompile(`"\.{0,2}/[^"]+"|'\.{0,2}/[^']+'`),
			exclude: []string{"_test.go"},
		},
		{
			name:  "URLs hardcodées",
			regex: regexp.MustCompile(`https?://[^"'\s]+`),
			exclude: []string{"_test.go", ".md", "go.mod"},
		},
	}

	// Dossiers à scanner
	dirs := []string{
		"../../cmd",
		"../../internal/api",
		"../../internal/web",
		"../../templates",
	}

	violations := []string{}

	for _, dir := range dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Ignorer les dossiers
			if info.IsDir() {
				return nil
			}

			// Vérifier seulement les fichiers Go et templates
			if !strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, ".html") {
				return nil
			}

			// Scanner le fichier
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			lineNum := 0

			for scanner.Scan() {
				lineNum++
				line := scanner.Text()

				// Ignorer les commentaires
				if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "//") {
					continue
				}

				// Vérifier chaque pattern
				for _, p := range patterns {
					// Vérifier les exclusions
					excluded := false
					for _, exc := range p.exclude {
						if strings.Contains(path, exc) {
							excluded = true
							break
						}
					}
					if excluded {
						continue
					}

					if p.regex.MatchString(line) {
						relPath, _ := filepath.Rel("../..", path)
						violation := fmt.Sprintf("%s:%d - %s: %s", relPath, lineNum, p.name, strings.TrimSpace(line))
						violations = append(violations, violation)
					}
				}
			}

			return nil
		})

		if err != nil {
			t.Logf("Erreur lors du scan de %s: %v", dir, err)
		}
	}

	// Rapport des violations
	if len(violations) > 0 {
		t.Errorf("\n🚨 VIOLATIONS DE LA POLITIQUE NO HARDCODING DÉTECTÉES: %d\n", len(violations))
		for _, v := range violations {
			t.Errorf("  ❌ %s", v)
		}
		t.Error("\n⚠️  VALIDATION REFUSÉE - Corrigez toutes les violations avant de continuer")
	}
}

// TestConfigurationComplete - Vérifie que toutes les valeurs sont externalisées
func TestConfigurationComplete(t *testing.T) {
	requiredEnvVars := []string{
		"PORT",
		"HOST",
		"ENV",
		"LOG_LEVEL",
		"MAX_PAGES_CRAWL",
		"TIMEOUT_SECONDS",
		// Ajouter toutes les variables requises
	}

	// Charger le fichier .env.example
	envExample := "../../.env.example"
	if _, err := os.Stat(envExample); os.IsNotExist(err) {
		t.Errorf("Fichier .env.example manquant")
		return
	}

	file, err := os.Open(envExample)
	if err != nil {
		t.Errorf("Erreur lecture .env.example: %v", err)
		return
	}
	defer file.Close()

	envVars := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
			parts := strings.Split(line, "=")
			envVars[parts[0]] = true
		}
	}

	// Vérifier que toutes les variables requises sont documentées
	for _, v := range requiredEnvVars {
		if !envVars[v] {
			t.Errorf("Variable d'environnement manquante dans .env.example: %s", v)
		}
	}
}

// TestNoHardcodedTestData - Vérifie que les données de test sont externalisées
func TestNoHardcodedTestData(t *testing.T) {
	testDataFiles := []string{
		"../../config/test-data.yaml",
		"../../config/simulation.yaml",
		"../../config/messages.yaml",
	}

	for _, file := range testDataFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Fichier de configuration manquant: %s", file)
		}
	}
}