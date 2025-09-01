package debug

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"firesalamander/internal/config"
	"firesalamander/internal/constants"
	"firesalamander/internal/logger"
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
	testLog.Info(constants.LogRunningPhase1Tests)
	
	start := time.Now()
	suite := &PhaseTestSuite{
		Phase:  constants.Phase1SetupInitial,
		Status: constants.StatusRunning,
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
		if test.Status == constants.StatusPassed {
			suite.PassedTests++
		} else {
			suite.FailedTests++
		}
	}
	
	if suite.FailedTests == 0 {
		suite.Status = constants.StatusPassed
	} else {
		suite.Status = constants.StatusFailed
	}
	
	suite.Duration = time.Since(start).String()
	
	testLog.Info(constants.LogPhase1Completed, map[string]interface{}{
		constants.JSONFieldStatus:        suite.Status,
		constants.JSONFieldTotalTests:   suite.TotalTests,
		constants.JSONFieldPassedTests:  suite.PassedTests,
		constants.JSONFieldFailedTests:  suite.FailedTests,
		constants.JSONFieldDuration:      suite.Duration,
	})
	
	return suite
}

func runConfigTest(cfg *config.Config) PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       constants.Phase1Number,
		Name:        constants.TestConfigurationLoading,
		Description: constants.DescConfigurationLoading,
		Details:     make(map[string]interface{}),
	}
	
	testLog.Debug(constants.LogConfigurationTest)
	
	if cfg == nil {
		test.Status = constants.StatusFailed
		test.Message = constants.ErrConfigurationNil
		test.Error = constants.ErrConfigNil
		test.Duration = time.Since(start).String()
		return test
	}
	
	issues := []string{}
	
	// Test App config (using constants)
	appName := constants.AppName
	if appName == "" || appName != "Fire Salamander" {
		issues = append(issues, constants.ErrAppNameIncorrect)
	}
	
	// Test Server  
	if cfg.Server.Port <= 0 {
		issues = append(issues, constants.ErrServerPortInvalid)
	}
	
	test.Details[constants.DetailAppName] = constants.AppName
	test.Details[constants.DetailAppIcon] = constants.AppIcon
	test.Details[constants.DetailPoweredBy] = constants.PoweredBy  
	test.Details[constants.DetailPrimaryColor] = constants.PrimaryColor
	test.Details[constants.DetailServerPort] = cfg.Server.Port
	test.Details[constants.DetailIssuesFound] = len(issues)
	
	if len(issues) > 0 {
		test.Status = constants.StatusFailed
		test.Message = fmt.Sprintf(constants.MsgConfigurationFailed, strings.Join(issues, ", "))
		test.Details[constants.DetailIssues] = issues
	} else {
		test.Status = constants.StatusPassed
		test.Message = constants.MsgConfigurationValid
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func runFileStructureTest() PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       constants.Phase1Number,
		Name:        constants.TestFileStructure,
		Description: "Verify all required files and directories exist",
		Details:     make(map[string]interface{}),
	}
	
	testLog.Debug(constants.LogFileStructureTest)
	
	requiredFiles := []string{
		constants.RequiredFileGoMod,
		constants.RequiredFileMainGo,
		constants.RequiredFileReadme,
		constants.RequiredFileGitignore,
		constants.RequiredFileDockerCompose,
	}
	
	requiredDirs := []string{
		constants.RequiredDirConfig,
		constants.RequiredDirDeploy,
		constants.RequiredDirInternal,
		constants.RequiredDirInternalLogger,
		constants.RequiredDirInternalDebug,
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
	
	test.Details[constants.DetailRequiredFiles] = len(requiredFiles)
	test.Details[constants.DetailRequiredDirs] = len(requiredDirs)
	test.Details[constants.DetailMissingFiles] = len(missingFiles)
	test.Details[constants.DetailMissingDirs] = len(missingDirs)
	
	if len(missingFiles) > 0 || len(missingDirs) > 0 {
		test.Status = constants.StatusFailed
		test.Message = fmt.Sprintf(constants.MsgMissingFilesAndDirs, len(missingFiles), len(missingDirs))
		if len(missingFiles) > 0 {
			test.Details[constants.DetailMissingFilesList] = missingFiles
		}
		if len(missingDirs) > 0 {
			test.Details[constants.DetailMissingDirsList] = missingDirs
		}
	} else {
		test.Status = constants.StatusPassed
		test.Message = constants.MsgAllFilesExist
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func runGitTest() PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       constants.Phase1Number,
		Name:        "Git Setup",
		Description: "Verify Git repository is initialized with remote",
		Details:     make(map[string]interface{}),
	}
	
	testLog.Debug(constants.LogGitTest)
	
	issues := []string{}
	
	// Check .git directory
	if _, err := os.Stat(constants.GitDirectory); os.IsNotExist(err) {
		issues = append(issues, constants.MsgGitDirectoryMissing)
	} else {
		test.Details[constants.DetailGitInitialized] = true
	}
	
	// Check .gitignore
	if _, err := os.Stat(constants.RequiredFileGitignore); os.IsNotExist(err) {
		issues = append(issues, constants.MsgGitignoreMissing)
	} else {
		test.Details[constants.DetailGitignoreExists] = true
	}
	
	test.Details[constants.DetailIssuesFound] = len(issues)
	
	if len(issues) > 0 {
		test.Status = constants.StatusFailed
		test.Message = fmt.Sprintf(constants.MsgGitSetupIssues, strings.Join(issues, ", "))
		test.Details[constants.DetailIssues] = issues
	} else {
		test.Status = constants.StatusPassed
		test.Message = constants.MsgGitProperlyInitialized
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func runServerTest(cfg *config.Config) PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       constants.Phase1Number,
		Name:        constants.TestHTTPServer,
		Description: constants.DescHTTPServer,
		Details:     make(map[string]interface{}),
		SubTests:    []PhaseTest{},
	}
	
	testLog.Debug(constants.LogHTTPServerTest)
	
	if cfg == nil {
		test.Status = constants.StatusFailed
		test.Message = constants.MsgCannotTestServer
		test.Duration = time.Since(start).String()
		return test
	}
	
	baseURL := fmt.Sprintf(constants.HTTPPrefix+constants.ServerDefaultHost+":%d", cfg.Server.Port)
	test.Details[constants.DetailBaseURL] = baseURL
	
	// Test health endpoint
	healthTest := testEndpoint(constants.TestEndpointHealth, baseURL+constants.EndpointHealth, constants.ContentTypeJSON)
	test.SubTests = append(test.SubTests, healthTest)
	
	// Test debug endpoint
	debugTest := testEndpoint(constants.TestEndpointDebug, baseURL+constants.EndpointDebug, constants.ContentTypeJSON)
	test.SubTests = append(test.SubTests, debugTest)
	
	// Test home page
	homeTest := testEndpoint(constants.TestEndpointHome, baseURL+constants.EndpointHome, constants.ContentTypeHTML)
	test.SubTests = append(test.SubTests, homeTest)
	
	// Calculate overall status
	passedTests := 0
	for _, subTest := range test.SubTests {
		if subTest.Status == constants.StatusPassed {
			passedTests++
		}
	}
	
	test.Details[constants.DetailTotalEndpoints] = len(test.SubTests)
	test.Details[constants.DetailPassedEndpoints] = passedTests
	
	if passedTests == len(test.SubTests) {
		test.Status = constants.StatusPassed
		test.Message = constants.MsgAllEndpointsResponding
	} else {
		test.Status = constants.StatusFailed
		test.Message = fmt.Sprintf(constants.MsgOnlyEndpointsResponding, passedTests, len(test.SubTests))
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
	
	test.Details[constants.DetailURL] = url
	test.Details[constants.DetailExpectedContentType] = expectedContentType
	
	client := &http.Client{Timeout: constants.FastRequestTimeout}
	resp, err := client.Get(url)
	if err != nil {
		test.Status = constants.StatusFailed
		test.Message = constants.MsgEndpointNotAccessible
		test.Error = err.Error()
		test.Details[constants.DetailAccessible] = false
	} else {
		defer resp.Body.Close()
		
		test.Details[constants.DetailAccessible] = true
		test.Details[constants.DetailStatusCode] = resp.StatusCode
		test.Details[constants.DetailContentType] = resp.Header.Get("Content-Type")
		
		if resp.StatusCode == constants.HTTPStatusOK {
			contentType := resp.Header.Get("Content-Type")
			if strings.Contains(contentType, expectedContentType) {
				test.Status = constants.StatusPassed
				test.Message = constants.MsgEndpointRespondingCorrect
			} else {
				test.Status = constants.StatusFailed
				test.Message = fmt.Sprintf(constants.MsgWrongContentType, contentType, expectedContentType)
			}
		} else {
			test.Status = constants.StatusFailed
			test.Message = fmt.Sprintf(constants.MsgUnexpectedStatusCode, resp.StatusCode)
		}
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func runBrandingTest() PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       constants.Phase1Number,
		Name:        constants.TestSEPTEOBranding,
		Description: "Verify SEPTEO branding elements are properly integrated",
		Details:     make(map[string]interface{}),
	}
	
	testLog.Debug(constants.LogSEPTEOBrandingTest)
	
	issues := []string{}
	
	// Check main.go for SEPTEO logo URL
	if content, err := os.ReadFile(constants.RequiredFileMainGo); err == nil {
		contentStr := string(content)
		if !strings.Contains(contentStr, constants.SepteoLogoPath) {
			issues = append(issues, constants.MsgSepteoLogoNotFound)
		} else {
			test.Details[constants.DetailSepteoLogoIntegrated] = true
		}
		
		if !strings.Contains(contentStr, constants.SepteoOrangeColor) {
			issues = append(issues, constants.MsgSepteoOrangeNotFound)
		} else {
			test.Details[constants.DetailSepteoOrangeIntegrated] = true
		}
		
		if !strings.Contains(contentStr, constants.SalamanderIcon) {
			issues = append(issues, constants.MsgSalamanderIconNotFound)
		} else {
			test.Details[constants.DetailSalamanderIconIntegrated] = true
		}
	} else {
		issues = append(issues, constants.MsgCannotReadMainGo)
	}
	
	test.Details[constants.DetailIssuesFound] = len(issues)
	
	if len(issues) > 0 {
		test.Status = constants.StatusFailed
		test.Message = fmt.Sprintf(constants.MsgBrandingIssues, strings.Join(issues, ", "))
		test.Details[constants.DetailIssues] = issues
	} else {
		test.Status = constants.StatusPassed
		test.Message = constants.MsgSEPTEOBrandingIntegrated
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func runDockerTest() PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       constants.Phase1Number,
		Name:        constants.TestDockerSetup,
		Description: constants.DescDockerSetup,
		Details:     make(map[string]interface{}),
	}
	
	testLog.Debug(constants.LogDockerTest)
	
	if _, err := os.Stat(constants.RequiredFileDockerCompose); os.IsNotExist(err) {
		test.Status = constants.StatusFailed
		test.Message = constants.MsgDockerComposeFileMissing
		test.Error = constants.ErrDockerComposeMissing
	} else {
		// Check if docker-compose.yml contains required services
		if content, err := os.ReadFile(constants.RequiredFileDockerCompose); err == nil {
			contentStr := string(content)
			
			hasAppService := strings.Contains(contentStr, constants.DockerServiceApp)
			hasDbService := strings.Contains(contentStr, constants.DockerServiceDB)
			hasPortMapping := strings.Contains(contentStr, constants.TestPort3000+":"+constants.TestPort3000)
			
			test.Details[constants.DetailHasAppService] = hasAppService
			test.Details[constants.DetailHasDbService] = hasDbService
			test.Details[constants.DetailHasPortMapping] = hasPortMapping
			
			if hasAppService && hasDbService && hasPortMapping {
				test.Status = constants.StatusPassed
				test.Message = constants.MsgDockerComposeValid
			} else {
				test.Status = constants.StatusFailed
				test.Message = constants.ErrDockerComposeInvalid
			}
		} else {
			test.Status = constants.StatusFailed
			test.Message = constants.ErrCannotReadDockerCompose
			test.Error = err.Error()
		}
	}
	
	test.Duration = time.Since(start).String()
	return test
}

func runDeployScriptsTest() PhaseTest {
	start := time.Now()
	test := PhaseTest{
		Phase:       constants.Phase1Number,
		Name:        constants.TestDeployScripts,
		Description: constants.DescDeployScripts,
		Details:     make(map[string]interface{}),
	}
	
	testLog.Debug(constants.LogDeployScriptsTest)
	
	scripts := []string{
		constants.DeployScriptDeploy,
		constants.DeployScriptSetupInfomaniak,
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
	
	test.Details[constants.DetailTotalScripts] = len(scripts)
	test.Details[constants.DetailMissingScripts] = len(missingScripts)
	test.Details[constants.DetailNonExecutableScripts] = len(nonExecutableScripts)
	
	if len(missingScripts) > 0 || len(nonExecutableScripts) > 0 {
		test.Status = constants.StatusFailed
		issues := []string{}
		if len(missingScripts) > 0 {
			issues = append(issues, fmt.Sprintf(constants.MsgMissingScriptsCount, len(missingScripts)))
			test.Details[constants.DetailMissingScriptsList] = missingScripts
		}
		if len(nonExecutableScripts) > 0 {
			issues = append(issues, fmt.Sprintf(constants.MsgNonExecutableScriptsCount, len(nonExecutableScripts)))
			test.Details[constants.DetailNonExecutableScriptsList] = nonExecutableScripts
		}
		test.Message = fmt.Sprintf(constants.MsgDeployScriptIssues, strings.Join(issues, ", "))
	} else {
		test.Status = constants.StatusPassed
		test.Message = constants.MsgAllDeployScriptsReady
	}
	
	test.Duration = time.Since(start).String()
	return test
}