#!/bin/bash
# scripts/mass-replace-hardcoded.sh
# MASS REPLACEMENT AUTOMATION ENGINE

set -e

FILE="$1"
if [ -z "$FILE" ]; then
    echo "Usage: $0 <file_to_fix>"
    exit 1
fi

echo "🤖 MASS REPLACEMENT ENGINE - Processing: $FILE"

# Backup original
cp "$FILE" "$FILE.backup"

# Replace common hardcoded strings with constants
sed -i '' 's/"schema"/'constants.IssueTypeSchema'/g' "$FILE"
sed -i '' 's/"data"/'constants.IssueTypeData'/g' "$FILE"  
sed -i '' 's/"integrity"/'constants.IssueTypeIntegrity'/g' "$FILE"
sed -i '' 's/"quality"/'constants.IssueTypeQuality'/g' "$FILE"
sed -i '' 's/"performance"/'constants.IssueTypePerformance'/g' "$FILE"

# Replace status strings
sed -i '' 's/"unknown"/'constants.StatusUnknown'/g' "$FILE"
sed -i '' 's/"passed"/'constants.StatusPassed'/g' "$FILE" 
sed -i '' 's/"failed"/'constants.StatusFailed'/g' "$FILE"
sed -i '' 's/"error"/'constants.StatusError'/g' "$FILE"
sed -i '' 's/"warning"/'constants.StatusWarning'/g' "$FILE"
sed -i '' 's/"success"/'constants.StatusSuccess'/g' "$FILE"
sed -i '' 's/"pass"/'constants.StatusPassed'/g' "$FILE"
sed -i '' 's/"fail"/'constants.StatusFailed'/g' "$FILE"

# Replace severity levels
sed -i '' 's/"low"/'constants.SeverityLow'/g' "$FILE"
sed -i '' 's/"medium"/'constants.SeverityMedium'/g' "$FILE"
sed -i '' 's/"high"/'constants.SeverityHigh'/g' "$FILE"  
sed -i '' 's/"critical"/'constants.SeverityCritical'/g' "$FILE"

# Replace test categories
sed -i '' 's/"schema_validation"/'constants.TestCategorySchemaValidation'/g' "$FILE"
sed -i '' 's/"data_consistency"/'constants.TestCategoryDataConsistency'/g' "$FILE"
sed -i '' 's/"referential_integrity"/'constants.TestCategoryReferentialIntegrity'/g' "$FILE"
sed -i '' 's/"data_quality"/'constants.TestCategoryDataQuality'/g' "$FILE"
sed -i '' 's/"performance_checks"/'constants.TestCategoryPerformanceChecks'/g' "$FILE"

# Replace test names
sed -i '' 's/"Schema Validation"/'constants.TestSchemaValidation'/g' "$FILE"
sed -i '' 's/"Constraints Check"/'constants.TestConstraintsCheck'/g' "$FILE"
sed -i '' 's/"Data Consistency"/'constants.TestDataConsistency'/g' "$FILE"
sed -i '' 's/"Referential Integrity"/'constants.TestReferentialCheck'/g' "$FILE"
sed -i '' 's/"Data Quality"/'constants.TestDataQuality'/g' "$FILE"
sed -i '' 's/"Performance Check"/'constants.TestPerformanceCheck'/g' "$FILE"

# Replace common messages
sed -i '' 's/"Missing required table"/'constants.MsgMissingRequiredTable'/g' "$FILE"
sed -i '' 's/"Application functionality may be impaired"/'constants.MsgApplicationImpaired'/g' "$FILE"
sed -i '' 's/"Data inconsistency detected"/'constants.MsgDataInconsistency'/g' "$FILE"
sed -i '' 's/"Referential integrity violation"/'constants.MsgReferentialIntegrityFail'/g' "$FILE"
sed -i '' 's/"Performance degradation detected"/'constants.MsgPerformanceDegradation'/g' "$FILE"

# Replace database related
sed -i '' 's/"sqlite3"/'constants.SQLite3Driver'/g' "$FILE"
sed -i '' 's/"fire_salamander_dev.db"/'constants.DefaultDatabasePath'/g' "$FILE"
sed -i '' 's/"tests\/reports\/data"/'constants.DefaultReportPath'/g' "$FILE"

# Replace table names
sed -i '' 's/"crawl_sessions"/'constants.TableCrawlSessions'/g' "$FILE"
sed -i '' 's/"pages"/'constants.TablePages'/g' "$FILE"
sed -i '' 's/"seo_metrics"/'constants.TableSEOMetrics'/g' "$FILE"

# Check if file compiles
echo "🔍 Validating Go compilation..."
if go build -o /dev/null "$FILE" 2>/dev/null; then
    echo "✅ File compiles successfully"
    rm "$FILE.backup"
else
    echo "❌ Compilation failed, restoring backup"
    mv "$FILE.backup" "$FILE"
    exit 1
fi

# Count remaining violations
REMAINING=$(grep -c '"[A-Za-z][A-Za-z0-9 ]\{4,\}"' "$FILE" | grep -v "const\|var\|//\|fmt\." || echo "0")
echo "📊 Remaining violations: $REMAINING"

echo "🎉 MASS REPLACEMENT COMPLETED for $FILE"