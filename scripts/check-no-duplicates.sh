#!/bin/bash
set -e

echo "🔍 Fire Salamander - Duplicate Check"
echo "====================================="

ERRORS=0

# Check for duplicate patterns
echo "Checking for duplicate files..."

if find . -path ./archive -prune -o -name "real_*.go" -print | grep -v test | grep -q .; then
    echo "❌ Found 'real_' prefixed files (excluding tests):"
    find . -path ./archive -prune -o -name "real_*.go" -print | grep -v test
    ERRORS=$((ERRORS + 1))
fi

if find . -path ./archive -prune -o -name "*_v2.go" -print | grep -q .; then
    echo "❌ Found '_v2' versioned files:"
    find . -path ./archive -prune -o -name "*_v2.go" -print
    ERRORS=$((ERRORS + 1))
fi

if find . -path ./archive -prune -o -name "new_*.go" -print | grep -q .; then
    echo "❌ Found 'new_' prefixed files:"
    find . -path ./archive -prune -o -name "new_*.go" -print
    ERRORS=$((ERRORS + 1))
fi

# Check for multiple handlers
echo "Checking for duplicate handlers..."
HANDLERS=$(grep -r "func.*Handler" --include="*.go" . | grep -v test | grep -v archive | wc -l)
if [ "$HANDLERS" -gt 10 ]; then
    echo "⚠️  Warning: Many handlers found ($HANDLERS), verify no duplicates"
    grep -r "func.*Handler" --include="*.go" . | grep -v test | grep -v archive | cut -d: -f2 | sort
fi

# Check for specific duplicate patterns
echo "Checking for specific duplicate patterns..."
ORCHESTRATORS=$(find . -path ./archive -prune -o -name "*orchestrator*.go" -print | grep -v test | wc -l)
ANALYZERS=$(find . -path ./archive -prune -o -name "*analyzer*.go" -print | grep -v test | grep "internal/" | wc -l)

if [ "$ORCHESTRATORS" -gt 1 ]; then
    echo "⚠️  Multiple orchestrator files found:"
    find . -path ./archive -prune -o -name "*orchestrator*.go" -print | grep -v test
fi

if [ "$ANALYZERS" -gt 2 ]; then
    echo "⚠️  Multiple main analyzer files found:"
    find . -path ./archive -prune -o -name "*analyzer*.go" -print | grep -v test | grep "internal/"
fi

if [ $ERRORS -eq 0 ]; then
    echo "✅ No critical duplicates found!"
    echo "📊 Summary:"
    echo "   - Orchestrators: $ORCHESTRATORS"
    echo "   - Main Analyzers: $ANALYZERS"
    echo "   - Total Handlers: $HANDLERS"
    exit 0
else
    echo "❌ Found $ERRORS critical duplicate issues!"
    echo "Please fix before committing."
    exit 1
fi