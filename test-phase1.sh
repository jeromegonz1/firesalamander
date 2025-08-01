#!/bin/bash

# Fire Salamander - Phase 1 Testing Script
echo "🔥🧪 Fire Salamander - Test Phase 1"
echo "=================================="

# Couleurs pour l'output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Compteurs
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

function run_test() {
    local test_name="$1"
    local test_command="$2"
    local expected_result="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "\n${BLUE}[TEST $TOTAL_TESTS]${NC} $test_name"
    echo "Command: $test_command"
    
    if eval "$test_command" > /dev/null 2>&1; then
        if [ "$expected_result" = "success" ]; then
            echo -e "${GREEN}✅ PASSED${NC}"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        else
            echo -e "${RED}❌ FAILED${NC} (Expected failure but got success)"
            FAILED_TESTS=$((FAILED_TESTS + 1))
        fi
    else
        if [ "$expected_result" = "failure" ]; then
            echo -e "${GREEN}✅ PASSED${NC} (Expected failure)"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        else
            echo -e "${RED}❌ FAILED${NC}"
            FAILED_TESTS=$((FAILED_TESTS + 1))
        fi
    fi
}

function test_file_exists() {
    local file="$1"
    local description="$2"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "\n${BLUE}[TEST $TOTAL_TESTS]${NC} $description"
    echo "Checking: $file"
    
    if [ -f "$file" ] || [ -d "$file" ]; then
        echo -e "${GREEN}✅ PASSED${NC} - $file exists"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ FAILED${NC} - $file not found"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

function test_content_contains() {
    local file="$1"
    local content="$2"
    local description="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "\n${BLUE}[TEST $TOTAL_TESTS]${NC} $description"
    echo "Checking: $file contains '$content'"
    
    if [ -f "$file" ] && grep -q "$content" "$file"; then
        echo -e "${GREEN}✅ PASSED${NC} - Content found"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ FAILED${NC} - Content not found"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

echo -e "\n${YELLOW}Phase 1: Setup Initial - Testing${NC}"
echo "================================="

# Test 1: File Structure
echo -e "\n${YELLOW}🏗️  Testing File Structure${NC}"
test_file_exists "go.mod" "Go module file exists"
test_file_exists "main.go" "Main application file exists"
test_file_exists "README.md" "README file exists"
test_file_exists ".gitignore" "Git ignore file exists"
test_file_exists "docker-compose.yml" "Docker Compose file exists"
test_file_exists "config" "Config directory exists"
test_file_exists "config/config.go" "Config Go file exists"
test_file_exists "config/config.dev.yaml" "Dev config exists"
test_file_exists "config/config.prod.yaml" "Prod config exists"
test_file_exists "deploy" "Deploy directory exists"
test_file_exists "deploy/deploy.sh" "Deploy script exists"
test_file_exists "deploy/setup-infomaniak.sh" "Setup script exists"
test_file_exists "internal" "Internal directory exists"
test_file_exists "internal/logger" "Logger directory exists"
test_file_exists "internal/debug" "Debug directory exists"

# Test 2: Git Setup
echo -e "\n${YELLOW}📦 Testing Git Setup${NC}"
test_file_exists ".git" "Git repository initialized"
run_test "Git remote configured" "git remote get-url origin" "success"

# Test 3: Configuration Content
echo -e "\n${YELLOW}⚙️  Testing Configuration${NC}"
test_content_contains "config/config.dev.yaml" "Fire Salamander" "App name in dev config"
test_content_contains "config/config.dev.yaml" "#ff6136" "SEPTEO orange color in config"
test_content_contains "config/config.dev.yaml" "🦎" "Salamander icon in config"
test_content_contains "config/config.dev.yaml" "SEPTEO" "SEPTEO branding in config"

# Test 4: SEPTEO Branding
echo -e "\n${YELLOW}🎨 Testing SEPTEO Branding${NC}"
test_content_contains "main.go" "septeo.svg" "SEPTEO logo URL in main.go"
test_content_contains "main.go" "Branding.PrimaryColor" "SEPTEO orange color template in main.go"
test_content_contains "main.go" "App.Icon" "Fire Salamander icon template in main.go"
test_content_contains "main.go" "SEPTEO" "SEPTEO branding mentions"

# Test 5: Go Module
echo -e "\n${YELLOW}📋 Testing Go Module${NC}"
test_content_contains "go.mod" "github.com/jeromegonz1/firesalamander" "Correct module name"
test_content_contains "go.mod" "gopkg.in/yaml.v3" "YAML dependency"
run_test "Go mod tidy successful" "go mod tidy" "success"

# Test 6: Compilation
echo -e "\n${YELLOW}🔨 Testing Compilation${NC}"
run_test "Application compiles successfully" "go build -o firesalamander-test ." "success"

# Cleanup test binary
if [ -f "firesalamander-test" ]; then
    rm firesalamander-test
fi

# Test 7: Docker Configuration
echo -e "\n${YELLOW}🐳 Testing Docker Setup${NC}"
test_content_contains "docker-compose.yml" "app:" "App service defined"
test_content_contains "docker-compose.yml" "db:" "Database service defined"
test_content_contains "docker-compose.yml" "3000:3000" "Port mapping configured"
test_content_contains "docker-compose.yml" "mysql:8" "MySQL 8 image specified"

# Test 8: Deploy Scripts
echo -e "\n${YELLOW}🚀 Testing Deploy Scripts${NC}"
test_content_contains "deploy/deploy.sh" "GOOS=linux" "Linux build in deploy script"
test_content_contains "deploy/deploy.sh" "🔥" "Fire emoji in deploy script"
test_content_contains "deploy/setup-infomaniak.sh" "mysql" "MySQL setup in Infomaniak script"
run_test "Deploy script is executable" "[ -x deploy/deploy.sh ]" "success"
run_test "Setup script is executable" "[ -x deploy/setup-infomaniak.sh ]" "success"

# Test 9: Logger and Debug System
echo -e "\n${YELLOW}🐛 Testing Debug System${NC}"
test_file_exists "internal/logger/logger.go" "Logger system exists"
test_file_exists "internal/debug/checker.go" "Debug checker exists"
test_file_exists "internal/debug/phase_tests.go" "Phase tests exist"
test_content_contains "internal/logger/logger.go" "DEBUG" "Debug level logging"
test_content_contains "internal/debug/checker.go" "HealthCheck" "Health check struct"

# Résultats finaux
echo -e "\n${YELLOW}📊 Test Results Summary${NC}"
echo "======================="
echo -e "Total Tests: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed: ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 ALL TESTS PASSED! Phase 1 is complete and functional.${NC}"
    echo -e "${GREEN}🦎 Fire Salamander Phase 1 - Setup Initial: ✅ SUCCESS${NC}"
    exit 0
else
    echo -e "\n${RED}⚠️  Some tests failed. Please review the issues above.${NC}"
    echo -e "${RED}🦎 Fire Salamander Phase 1 - Setup Initial: ❌ NEEDS FIXES${NC}"
    exit 1
fi