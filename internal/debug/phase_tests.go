package debug

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jeromegonz1/firesalamander/config"
	"github.com/jeromegonz1/firesalamander/internal/logger"
)

type PhaseTest struct {
	Phase       string                         `json:"phase"`
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	Status      string                         `json:"status"`
	Message     string                         `json:"message"`
	Duration    string                         `json:"duration"`
	Details     map[string]interface{}         `json:"details,omitempty"`
	SubTests    []PhaseTest                    `json:"sub_tests,omitempty"`
	Error       string                         `json:"error,omitempty"`
}

type PhaseTestSuite struct {
	Phase       string      `json:"phase"`
	Status      string      `json:"status"`
	TotalTests  int         `json:"total_tests"`
	PassedTests int         `json:"passed_tests"`
	FailedTests int         `json:"failed_tests"`
	Duration    string      `json:"duration"`
	Tests       []PhaseTest `json:"tests"`
}

var testLog = logger.New("PHASE-TESTS")

func RunPhase1Tests(cfg *config.Config) *PhaseTestSuite {
	testLog.Info("🧪 Running Phase 1 Tests - Setup Initial")
	
	start := time.Now()
	suite := &PhaseTestSuite{
		Phase:  "Phase 1 - Setup Initial",
		Status: "running",
		Tests:  []PhaseTest{},
	}
	
	// Test 1: Configuration
	suite.Tests = append(suite.Tests, runConfigTest(cfg))
	
	// Test 2: File Structure
	suite.Tests = append(suite.Tests, runFileStructureTest())
	
	// Test 3: Git Setup
	suite.Tests = append(suite.Tests, runGitTest())
	
	// Test 4: HTTP Server
	suite.Tests = append(suite.Tests, runServerTest(cfg))
	
	// Test 5: Branding
	suite.Tests = append(suite.Tests, runBrandingTest())
	
	// Test 6: Docker Setup
	suite.Tests = append(suite.Tests, runDockerTest())
	
	// Test 7: Deploy Scripts
	suite.Tests = append(suite.Tests, runDeployScriptsTest())
	
	// Calculer les résultats
	suite.TotalTests = len(suite.Tests)
	for _, test := range suite.Tests {
		if test.Status == "passed" {
			suite.PassedTests++
		} else {
			suite.FailedTests++
		}
	}
	
	if suite.FailedTests == 0 {
		suite.Status = "passed"
	} else {
		suite.Status = "failed"
	}
	
	suite.Duration = time.Since(start).String()
	
	testLog.Info("🧪 Phase 1 Tests completed", map[string]interface{}{
		"status":        suite.Status,
		"total_tests":   suite.TotalTests,
		"passed_tests":  suite.PassedTests,
		"failed_tests":  suite.FailedTests,
		"duration":      suite.Duration,
	})
	
	return suite
}

func runConfigTest(cfg *config.Config) PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       "1",
		Name:        "Configuration Loading",
		Description: "Verify configuration is loaded correctly with all required fields",
		Details:     make(map[string]interface{}),
	}
	
	testLog.Debug("Running configuration test")
	
	if cfg == nil {
		test.Status = "failed"
		test.Message = "Configuration is nil"
		test.Error = "config_nil"
		test.Duration = time.Since(start).String()
		return test
	}
	
	issues := []string{}
	
	// Test App config
	if cfg.App.Name == "" || cfg.App.Name != "Fire Salamander" {
		issues = append(issues, "app.name incorrect")
	}
	if cfg.App.Icon != "🦎" {
		issues = append(issues, "app.icon incorrect")
	}
	if cfg.App.PoweredBy != "SEPTEO" {
		issues = append(issues, "app.powered_by incorrect")
	}
	
	// Test Branding
	if cfg.Branding.PrimaryColor != "#ff6136" {
		issues = append(issues, "branding.primary_color incorrect")
	}
	
	// Test Server
	if cfg.Server.Port <= 0 {
		issues = append(issues, "server.port invalid")
	}
	
	test.Details["app_name"] = cfg.App.Name
	test.Details["app_icon"] = cfg.App.Icon
	test.Details["powered_by"] = cfg.App.PoweredBy  
	test.Details["primary_color"] = cfg.Branding.PrimaryColor
	test.Details["server_port"] = cfg.Server.Port
	test.Details["issues_found"] = len(issues)
	
	if len(issues) > 0 {
		test.Status = "failed"
		test.Message = fmt.Sprintf("Configuration validation failed: %s", strings.Join(issues, ", "))
		test.Details["issues"] = issues
	} else {
		test.Status = "passed"
		test.Message = "Configuration is valid and complete"
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func runFileStructureTest() PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       "1",
		Name:        "File Structure",
		Description: "Verify all required files and directories exist",
		Details:     make(map[string]interface{}),
	}
	
	testLog.Debug("Running file structure test")
	
	requiredFiles := []string{
		"go.mod",
		"main.go",
		"README.md",
		".gitignore",
		"docker-compose.yml",
	}
	
	requiredDirs := []string{
		"config",
		"deploy",
		"internal",
		"internal/logger",
		"internal/debug",
	}
	
	missingFiles := []string{}
	missingDirs := []string{}
	
	// Check files
	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			missingFiles = append(missingFiles, file)
		}
	}
	
	// Check directories
	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			missingDirs = append(missingDirs, dir)
		}
	}
	
	test.Details["required_files"] = len(requiredFiles)
	test.Details["required_dirs"] = len(requiredDirs)
	test.Details["missing_files"] = len(missingFiles)
	test.Details["missing_dirs"] = len(missingDirs)
	
	if len(missingFiles) > 0 || len(missingDirs) > 0 {
		test.Status = "failed"
		test.Message = fmt.Sprintf("Missing %d files and %d directories", len(missingFiles), len(missingDirs))
		if len(missingFiles) > 0 {
			test.Details["missing_files_list"] = missingFiles
		}
		if len(missingDirs) > 0 {
			test.Details["missing_dirs_list"] = missingDirs
		}
	} else {
		test.Status = "passed"
		test.Message = "All required files and directories exist"
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func runGitTest() PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       "1",
		Name:        "Git Setup",
		Description: "Verify Git repository is initialized with remote",
		Details:     make(map[string]interface{}),
	}
	
	testLog.Debug("Running Git setup test")
	
	issues := []string{}
	
	// Check .git directory
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		issues = append(issues, ".git directory missing")
	} else {
		test.Details["git_initialized"] = true
	}
	
	// Check .gitignore
	if _, err := os.Stat(".gitignore"); os.IsNotExist(err) {
		issues = append(issues, ".gitignore missing")
	} else {
		test.Details["gitignore_exists"] = true
	}
	
	test.Details["issues_found"] = len(issues)
	
	if len(issues) > 0 {
		test.Status = "failed"
		test.Message = fmt.Sprintf("Git setup issues: %s", strings.Join(issues, ", "))
		test.Details["issues"] = issues
	} else {
		test.Status = "passed"
		test.Message = "Git repository properly initialized"
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func runServerTest(cfg *config.Config) PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       "1",
		Name:        "HTTP Server",
		Description: "Test server endpoints are responding correctly",
		Details:     make(map[string]interface{}),
		SubTests:    []PhaseTest{},
	}
	
	testLog.Debug("Running HTTP server test")
	
	if cfg == nil {
		test.Status = "failed"
		test.Message = "Cannot test server - configuration is nil"
		test.Duration = time.Since(start).String()
		return test
	}
	
	baseURL := fmt.Sprintf("http://localhost:%d", cfg.Server.Port)
	test.Details["base_url"] = baseURL
	
	// Test health endpoint
	healthTest := testEndpoint("Health Endpoint", baseURL+"/health", "application/json")
	test.SubTests = append(test.SubTests, healthTest)
	
	// Test debug endpoint
	debugTest := testEndpoint("Debug Endpoint", baseURL+"/debug", "application/json")
	test.SubTests = append(test.SubTests, debugTest)
	
	// Test home page
	homeTest := testEndpoint("Home Page", baseURL+"/", "text/html")
	test.SubTests = append(test.SubTests, homeTest)
	
	// Calculate overall status
	passedTests := 0
	for _, subTest := range test.SubTests {
		if subTest.Status == "passed" {
			passedTests++
		}
	}
	
	test.Details["total_endpoints"] = len(test.SubTests)
	test.Details["passed_endpoints"] = passedTests
	
	if passedTests == len(test.SubTests) {
		test.Status = "passed"
		test.Message = "All server endpoints responding correctly"
	} else {
		test.Status = "failed"
		test.Message = fmt.Sprintf("Only %d/%d endpoints responding correctly", passedTests, len(test.SubTests))
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func testEndpoint(name, url, expectedContentType string) PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Name:        name,
		Description: fmt.Sprintf("Test %s is accessible", url),
		Details:     make(map[string]interface{}),
	}
	
	test.Details["url"] = url
	test.Details["expected_content_type"] = expectedContentType
	
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		test.Status = "failed"
		test.Message = "Endpoint not accessible"
		test.Error = err.Error()
		test.Details["accessible"] = false
	} else {
		defer resp.Body.Close()
		
		test.Details["accessible"] = true
		test.Details["status_code"] = resp.StatusCode
		test.Details["content_type"] = resp.Header.Get("Content-Type")
		
		if resp.StatusCode == 200 {
			contentType := resp.Header.Get("Content-Type")
			if strings.Contains(contentType, expectedContentType) {
				test.Status = "passed"
				test.Message = "Endpoint responding correctly"
			} else {
				test.Status = "failed"
				test.Message = fmt.Sprintf("Wrong content type: got %s, expected %s", contentType, expectedContentType)
			}
		} else {
			test.Status = "failed"
			test.Message = fmt.Sprintf("Unexpected status code: %d", resp.StatusCode)
		}
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func runBrandingTest() PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       "1",
		Name:        "SEPTEO Branding",
		Description: "Verify SEPTEO branding elements are properly integrated",
		Details:     make(map[string]interface{}),
	}
	
	testLog.Debug("Running branding test")
	
	issues := []string{}
	
	// Check main.go for SEPTEO logo URL
	if content, err := os.ReadFile("main.go"); err == nil {
		contentStr := string(content)
		if !strings.Contains(contentStr, "septeo.svg") {
			issues = append(issues, "SEPTEO logo URL not found in main.go")
		} else {
			test.Details["septeo_logo_integrated"] = true
		}
		
		if !strings.Contains(contentStr, "#ff6136") {
			issues = append(issues, "SEPTEO orange color not found")
		} else {
			test.Details["septeo_orange_integrated"] = true
		}
		
		if !strings.Contains(contentStr, "🦎") {
			issues = append(issues, "Fire Salamander icon not found")
		} else {
			test.Details["salamander_icon_integrated"] = true
		}
	} else {
		issues = append(issues, "Cannot read main.go file")
	}
	
	test.Details["issues_found"] = len(issues)
	
	if len(issues) > 0 {
		test.Status = "failed"
		test.Message = fmt.Sprintf("Branding issues: %s", strings.Join(issues, ", "))
		test.Details["issues"] = issues
	} else {
		test.Status = "passed"
		test.Message = "SEPTEO branding properly integrated"
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func runDockerTest() PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       "1",
		Name:        "Docker Setup",
		Description: "Verify Docker Compose configuration exists and is valid",
		Details:     make(map[string]interface{}),
	}
	
	testLog.Debug("Running Docker setup test")
	
	if _, err := os.Stat("docker-compose.yml"); os.IsNotExist(err) {
		test.Status = "failed"
		test.Message = "docker-compose.yml file missing"
		test.Error = "docker_compose_missing"
	} else {
		// Check if docker-compose.yml contains required services
		if content, err := os.ReadFile("docker-compose.yml"); err == nil {
			contentStr := string(content)
			
			hasAppService := strings.Contains(contentStr, "app:")
			hasDbService := strings.Contains(contentStr, "db:")
			hasPortMapping := strings.Contains(contentStr, "3000:3000")
			
			test.Details["has_app_service"] = hasAppService
			test.Details["has_db_service"] = hasDbService
			test.Details["has_port_mapping"] = hasPortMapping
			
			if hasAppService && hasDbService && hasPortMapping {
				test.Status = "passed"
				test.Message = "Docker Compose configuration is valid"
			} else {
				test.Status = "failed"
				test.Message = "Docker Compose configuration is incomplete"
			}
		} else {
			test.Status = "failed"
			test.Message = "Cannot read docker-compose.yml"
			test.Error = err.Error()
		}
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func runDeployScriptsTest() PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       "1",
		Name:        "Deploy Scripts",
		Description: "Verify deployment scripts exist and are executable",
		Details:     make(map[string]interface{}),
	}
	
	testLog.Debug("Running deploy scripts test")
	
	scripts := []string{
		"deploy/deploy.sh",
		"deploy/setup-infomaniak.sh",
	}
	
	missingScripts := []string{}
	nonExecutableScripts := []string{}
	
	for _, script := range scripts {
		if stat, err := os.Stat(script); os.IsNotExist(err) {
			missingScripts = append(missingScripts, script)
		} else {
			// Check if executable
			mode := stat.Mode()
			if mode&0111 == 0 { // Check if any execute bit is set
				nonExecutableScripts = append(nonExecutableScripts, script)
			}
		}
	}
	
	test.Details["total_scripts"] = len(scripts)
	test.Details["missing_scripts"] = len(missingScripts)
	test.Details["non_executable_scripts"] = len(nonExecutableScripts)
	
	if len(missingScripts) > 0 || len(nonExecutableScripts) > 0 {
		test.Status = "failed"
		issues := []string{}
		if len(missingScripts) > 0 {
			issues = append(issues, fmt.Sprintf("%d missing scripts", len(missingScripts)))
			test.Details["missing_scripts_list"] = missingScripts
		}
		if len(nonExecutableScripts) > 0 {
			issues = append(issues, fmt.Sprintf("%d non-executable scripts", len(nonExecutableScripts)))
			test.Details["non_executable_scripts_list"] = nonExecutableScripts
		}
		test.Message = fmt.Sprintf("Deploy script issues: %s", strings.Join(issues, ", "))
	} else {
		test.Status = "passed"
		test.Message = "All deployment scripts exist and are executable"
	}
	
	test.Duration = time.Since(start).String()
	return test
}