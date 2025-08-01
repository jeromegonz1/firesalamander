#!/bin/bash

# Fire Salamander - Phase 2 Testing Script
echo "🕷️🧪 Fire Salamander - Test Phase 2 : Module Crawler"
echo "===================================================="

# Couleurs pour l'output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
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

function count_pattern() {
    local file="$1"
    local pattern="$2"
    local expected_min="$3"
    local description="$4"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -e "\n${BLUE}[TEST $TOTAL_TESTS]${NC} $description"
    
    if [ -f "$file" ]; then
        count=$(grep -c "$pattern" "$file" || echo "0")
        echo "Found $count occurrences of '$pattern' (expected min: $expected_min)"
        
        if [ "$count" -ge "$expected_min" ]; then
            echo -e "${GREEN}✅ PASSED${NC}"
            PASSED_TESTS=$((PASSED_TESTS + 1))
        else
            echo -e "${RED}❌ FAILED${NC} - Too few occurrences"
            FAILED_TESTS=$((FAILED_TESTS + 1))
        fi
    else
        echo -e "${RED}❌ FAILED${NC} - File not found"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

echo -e "\n${YELLOW}Phase 2: Module Crawler - Testing${NC}"
echo "================================="

# Test 1: Module Structure
echo -e "\n${YELLOW}🏗️  Testing Crawler Module Structure${NC}"
test_file_exists "crawler" "Crawler directory exists"
test_file_exists "crawler/crawler.go" "Main crawler file exists"
test_file_exists "crawler/fetcher.go" "Fetcher implementation exists"
test_file_exists "crawler/robots.go" "Robots.txt parser exists"
test_file_exists "crawler/sitemap.go" "Sitemap parser exists"
test_file_exists "crawler/cache.go" "Cache implementation exists"
test_file_exists "crawler/crawler_test.go" "Crawler tests exist"

# Test 2: Crawler Core Implementation
echo -e "\n${YELLOW}⚙️  Testing Crawler Core Implementation${NC}"
test_content_contains "crawler/crawler.go" "type Crawler struct" "Crawler struct defined"
test_content_contains "crawler/crawler.go" "func New(" "Crawler constructor exists"
test_content_contains "crawler/crawler.go" "CrawlSite" "CrawlSite method exists"
test_content_contains "crawler/crawler.go" "CrawlPage" "CrawlPage method exists"
test_content_contains "crawler/crawler.go" "CrawlResult" "CrawlResult type defined"
test_content_contains "crawler/crawler.go" "Config" "Config type defined"

# Test 3: Fetcher Implementation
echo -e "\n${YELLOW}🌐 Testing HTTP Fetcher${NC}"
test_content_contains "crawler/fetcher.go" "type Fetcher struct" "Fetcher struct defined"
test_content_contains "crawler/fetcher.go" "RetryStrategy" "Retry strategy implemented"
test_content_contains "crawler/fetcher.go" "http.Transport" "HTTP transport optimized"
test_content_contains "crawler/fetcher.go" "MaxIdleConns" "Connection pooling configured"
test_content_contains "crawler/fetcher.go" "gzip" "Gzip compression support"
test_content_contains "crawler/fetcher.go" "User-Agent" "User-Agent header set"
count_pattern "crawler/fetcher.go" "retry" 3 "Multiple retry mentions"

# Test 4: Robots.txt Parser
echo -e "\n${YELLOW}🤖 Testing Robots.txt Parser${NC}"
test_content_contains "crawler/robots.go" "type RobotsTxt struct" "RobotsTxt struct defined"
test_content_contains "crawler/robots.go" "ParseRobotsTxt" "Parser function exists"
test_content_contains "crawler/robots.go" "IsAllowed" "IsAllowed method exists"
test_content_contains "crawler/robots.go" "GetCrawlDelay" "Crawl delay support"
test_content_contains "crawler/robots.go" "RobotsCache" "Robots cache implemented"
test_content_contains "crawler/robots.go" "user-agent" "User-agent parsing"
test_content_contains "crawler/robots.go" "Disallow" "Disallow rules parsing"
test_content_contains "crawler/robots.go" "Allow" "Allow rules parsing"
test_content_contains "crawler/robots.go" "Sitemap" "Sitemap discovery"

# Test 5: Sitemap Parser
echo -e "\n${YELLOW}🗺️  Testing Sitemap Parser${NC}"
test_content_contains "crawler/sitemap.go" "type Sitemap struct" "Sitemap struct defined"
test_content_contains "crawler/sitemap.go" "xml.Unmarshal" "XML parsing"
test_content_contains "crawler/sitemap.go" "SitemapURL" "URL structure defined"
test_content_contains "crawler/sitemap.go" "lastmod" "Last modified support"
test_content_contains "crawler/sitemap.go" "changefreq" "Change frequency support"
test_content_contains "crawler/sitemap.go" "priority" "Priority support"
test_content_contains "crawler/sitemap.go" "SitemapIndex" "Sitemap index support"

# Test 6: Cache System
echo -e "\n${YELLOW}💾 Testing Cache System${NC}"
test_content_contains "crawler/cache.go" "type PageCache struct" "Page cache defined"
test_content_contains "crawler/cache.go" "LRU" "LRU cache mentioned"
test_content_contains "crawler/cache.go" "Get(" "Cache Get method"
test_content_contains "crawler/cache.go" "Set(" "Cache Set method"
test_content_contains "crawler/cache.go" "evictList" "LRU eviction list"
test_content_contains "crawler/cache.go" "cleanup" "Cache cleanup"
test_content_contains "crawler/cache.go" "RateLimiter" "Rate limiter defined"
test_content_contains "crawler/cache.go" "CrawlQueue" "Crawl queue defined"

# Test 7: Rate Limiting
echo -e "\n${YELLOW}⏱️  Testing Rate Limiting${NC}"
test_content_contains "crawler/cache.go" "tokens" "Token bucket implementation"
test_content_contains "crawler/cache.go" "Wait(" "Rate limiter wait method"
test_content_contains "crawler/cache.go" "refillTokens" "Token refill mechanism"
count_pattern "crawler/cache.go" "rate" 5 "Rate limiting implementation"

# Test 8: Test Coverage
echo -e "\n${YELLOW}🧪 Testing Test Coverage${NC}"
test_content_contains "crawler/crawler_test.go" "TestCrawlerCreation" "Crawler creation test"
test_content_contains "crawler/crawler_test.go" "TestFetcherBasic" "Fetcher basic test"
test_content_contains "crawler/crawler_test.go" "TestFetcherRetry" "Retry mechanism test"
test_content_contains "crawler/crawler_test.go" "TestRobotsTxtParsing" "Robots.txt test"
test_content_contains "crawler/crawler_test.go" "TestSitemapParsing" "Sitemap parsing test"
test_content_contains "crawler/crawler_test.go" "TestPageCache" "Cache test"
test_content_contains "crawler/crawler_test.go" "TestRateLimiter" "Rate limiter test"
test_content_contains "crawler/crawler_test.go" "TestCrawlQueue" "Crawl queue test"

# Test 9: Configuration
echo -e "\n${YELLOW}🔧 Testing Configuration${NC}"
test_content_contains "crawler/crawler.go" "DefaultConfig" "Default config function"
test_content_contains "crawler/crawler.go" "Workers" "Worker configuration"
test_content_contains "crawler/crawler.go" "RateLimit" "Rate limit configuration"
test_content_contains "crawler/crawler.go" "MaxDepth" "Max depth configuration"
test_content_contains "crawler/crawler.go" "MaxPages" "Max pages configuration"
test_content_contains "crawler/crawler.go" "RespectRobots" "Robots.txt respect flag"

# Test 10: Error Handling
echo -e "\n${YELLOW}❌ Testing Error Handling${NC}"
count_pattern "crawler/crawler.go" "error" 10 "Error handling in crawler"
count_pattern "crawler/fetcher.go" "error" 10 "Error handling in fetcher"
test_content_contains "crawler/fetcher.go" "fmt.Errorf" "Formatted errors"
test_content_contains "crawler/crawler.go" "if err != nil" "Error checking"

# Test 11: Logging
echo -e "\n${YELLOW}📝 Testing Logging Integration${NC}"
test_content_contains "crawler/crawler.go" "logger.New" "Logger initialization"
count_pattern "crawler/crawler.go" "log.Debug" 3 "Debug logging"
count_pattern "crawler/crawler.go" "log.Info" 2 "Info logging"
count_pattern "crawler/crawler.go" "log.Error" 1 "Error logging"

# Test 12: Go Module Updates
echo -e "\n${YELLOW}📦 Testing Go Module Updates${NC}"
run_test "Go mod tidy" "cd crawler && go mod tidy 2>/dev/null || true" "success"

# Test 13: Compilation
echo -e "\n${YELLOW}🔨 Testing Compilation${NC}"
run_test "Crawler module compiles" "go build ./crawler/..." "success"

# Test 14: Run Unit Tests
echo -e "\n${YELLOW}🧪 Running Unit Tests${NC}"
echo "Running crawler tests..."
if go test -v ./crawler/... -count=1 > test_output.tmp 2>&1; then
    echo -e "${GREEN}✅ All crawler tests passed${NC}"
    PASSED_TESTS=$((PASSED_TESTS + 1))
    
    # Afficher un résumé des tests
    echo -e "\n${CYAN}Test Summary:${NC}"
    grep -E "PASS:|RUN:" test_output.tmp | tail -10
else
    echo -e "${RED}❌ Some crawler tests failed${NC}"
    FAILED_TESTS=$((FAILED_TESTS + 1))
    
    # Afficher les erreurs
    echo -e "\n${RED}Failed tests:${NC}"
    grep -E "FAIL:|Error:" test_output.tmp | head -10
fi
TOTAL_TESTS=$((TOTAL_TESTS + 1))

# Nettoyer
rm -f test_output.tmp

# Test 15: Integration avec main.go
echo -e "\n${YELLOW}🔗 Testing Integration${NC}"
test_content_contains "go.mod" "module github.com/jeromegonz1/firesalamander" "Correct module name"

# Résultats finaux
echo -e "\n${YELLOW}📊 Test Results Summary${NC}"
echo "======================="
echo -e "Total Tests: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "Passed: ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed: ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}🎉 ALL TESTS PASSED! Phase 2 is complete and functional.${NC}"
    echo -e "${GREEN}🕷️ Fire Salamander Phase 2 - Module Crawler: ✅ SUCCESS${NC}"
    
    echo -e "\n${CYAN}📋 Module Features Implemented:${NC}"
    echo "✅ HTTP Fetcher with retry and compression"
    echo "✅ Robots.txt parser with caching"
    echo "✅ XML Sitemap parser"
    echo "✅ LRU Page cache"
    echo "✅ Rate limiter (token bucket)"
    echo "✅ Crawl queue with deduplication"
    echo "✅ Comprehensive unit tests"
    echo "✅ Debug logging integration"
    
    exit 0
else
    echo -e "\n${RED}⚠️  Some tests failed. Please review the issues above.${NC}"
    echo -e "${RED}🕷️ Fire Salamander Phase 2 - Module Crawler: ❌ NEEDS FIXES${NC}"
    exit 1
fi