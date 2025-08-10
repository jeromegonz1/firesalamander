#!/bin/bash
set -e

echo "🔍 Fire Salamander - Complete Validation"
echo "======================================="

ERRORS=0
WARNINGS=0

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_section() {
    echo ""
    echo "📋 $1"
    echo "-----------------------------------"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
    ERRORS=$((ERRORS + 1))
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
    WARNINGS=$((WARNINGS + 1))
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

# 1. Check for duplicates
print_section "Checking for duplicate code"
if ./scripts/check-no-duplicates.sh; then
    print_success "No duplicates found"
else
    print_error "Duplicates detected"
fi

# 2. Check for hardcoding (if script exists)
print_section "Checking for hardcoding violations"
if [ -f "./scripts/detect-hardcoding.sh" ]; then
    if ./scripts/detect-hardcoding.sh 2>/dev/null; then
        print_success "No hardcoding violations"
    else
        print_warning "Some hardcoding violations found"
    fi
else
    print_warning "Hardcoding detection script not found"
fi

# 3. Run tests
print_section "Running all tests"
if go test ./...; then
    print_success "All tests pass"
else
    print_error "Some tests failed"
fi

# 4. Check test coverage
print_section "Checking test coverage"
COVERAGE_OUTPUT=$(go test ./... -cover 2>&1 | grep -E "coverage:" | tail -1)
if [ -n "$COVERAGE_OUTPUT" ]; then
    COVERAGE=$(echo "$COVERAGE_OUTPUT" | grep -oE '[0-9]+(\.[0-9]+)?%' | head -1 | sed 's/%//')
    
    if [ -n "$COVERAGE" ]; then
        # Use awk for floating point comparison (more portable than bc)
        if awk "BEGIN {exit !($COVERAGE >= 80)}"; then
            print_success "Coverage is sufficient: ${COVERAGE}%"
        else
            print_warning "Coverage is below 80%: ${COVERAGE}%"
        fi
    else
        print_warning "Could not determine coverage percentage"
    fi
else
    print_warning "No coverage information found"
fi

# 5. Check build
print_section "Building all packages"
if go build ./...; then
    print_success "Build successful"
else
    print_error "Build failed"
fi

# 6. Check for Go formatting
print_section "Checking Go formatting"
UNFORMATTED=$(gofmt -l . | grep -v vendor | grep -v archive)
if [ -z "$UNFORMATTED" ]; then
    print_success "All Go files are formatted correctly"
else
    print_warning "Some files need formatting:"
    echo "$UNFORMATTED"
fi

# 7. Check for Go vet issues
print_section "Running Go vet"
if go vet ./...; then
    print_success "Go vet passed"
else
    print_warning "Go vet found issues"
fi

# 8. Security check - look for panic usage
print_section "Checking for panic usage"
PANIC_COUNT=$(grep -r "panic(" --include="*.go" . | grep -v test | grep -v archive | wc -l)
if [ "$PANIC_COUNT" -eq 0 ]; then
    print_success "No panic usage in production code"
else
    print_warning "Found $PANIC_COUNT instances of panic() in production code"
    grep -r "panic(" --include="*.go" . | grep -v test | grep -v archive | head -5
fi

# 9. Check for TODO/FIXME comments
print_section "Checking for TODO/FIXME comments"
TODO_COUNT=$(grep -r -E "(TODO|FIXME|XXX)" --include="*.go" . | grep -v archive | wc -l)
if [ "$TODO_COUNT" -eq 0 ]; then
    print_success "No TODO/FIXME comments found"
else
    print_warning "Found $TODO_COUNT TODO/FIXME comments"
fi

# Final summary
echo ""
echo "🎯 VALIDATION SUMMARY"
echo "===================="

if [ $ERRORS -eq 0 ]; then
    print_success "ALL CRITICAL VALIDATIONS PASSED!"
    if [ $WARNINGS -eq 0 ]; then
        echo -e "${GREEN}🏆 PERFECT SCORE! No issues found.${NC}"
    else
        echo -e "${YELLOW}⚠️  $WARNINGS warnings found (non-blocking)${NC}"
    fi
    exit 0
else
    print_error "Found $ERRORS critical issues that MUST be fixed!"
    if [ $WARNINGS -gt 0 ]; then
        echo -e "${YELLOW}⚠️  Also found $WARNINGS warnings${NC}"
    fi
    echo ""
    echo "Please fix the critical issues before committing."
    exit 1
fi