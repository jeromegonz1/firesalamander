#!/usr/bin/env python3
"""
🤖 SMART HARDCODE ELIMINATOR
Intelligent automation for massive hardcoding violations elimination
"""

import re
import os
import subprocess
import sys
from typing import Dict, List, Tuple

class SmartHardcodeEliminator:
    def __init__(self, file_path: str):
        self.file_path = file_path
        self.backup_path = f"{file_path}.backup"
        
        # Mapping of hardcoded strings to constants
        self.string_mappings = {
            # Status constants
            '"unknown"': 'constants.StatusUnknown',
            '"passed"': 'constants.StatusPassed',
            '"failed"': 'constants.StatusFailed',
            '"error"': 'constants.StatusError',
            '"warning"': 'constants.StatusWarning',
            '"success"': 'constants.StatusSuccess',
            '"pass"': 'constants.StatusPassed',
            '"fail"': 'constants.StatusFailed',
            '"pending"': 'constants.StatusPending',
            '"running"': 'constants.StatusRunning',
            '"completed"': 'constants.StatusCompleted',
            
            # Severity constants
            '"low"': 'constants.SeverityLow',
            '"medium"': 'constants.SeverityMedium', 
            '"high"': 'constants.SeverityHigh',
            '"critical"': 'constants.SeverityCritical',
            
            # Issue types
            '"schema"': 'constants.IssueTypeSchema',
            '"data"': 'constants.IssueTypeData',
            '"integrity"': 'constants.IssueTypeIntegrity',
            '"quality"': 'constants.IssueTypeQuality',
            '"performance"': 'constants.IssueTypePerformance',
            '"data_quality"': 'constants.IssueTypeQuality',
            '"data_consistency"': 'constants.IssueTypeData',
            '"referential_integrity"': 'constants.IssueTypeIntegrity',
            
            # Test categories
            '"schema_validation"': 'constants.TestCategorySchemaValidation',
            '"data_consistency"': 'constants.TestCategoryDataConsistency',
            '"referential_integrity"': 'constants.TestCategoryReferentialIntegrity',
            '"data_quality"': 'constants.TestCategoryDataQuality',
            '"performance_checks"': 'constants.TestCategoryPerformanceChecks',
            
            # Test names
            '"Schema Validation"': 'constants.TestSchemaValidation',
            '"Constraints Check"': 'constants.TestConstraintsCheck',
            '"Data Consistency"': 'constants.TestDataConsistency',
            '"Referential Integrity"': 'constants.TestReferentialCheck',
            '"Data Quality"': 'constants.TestDataQuality',
            '"Performance Check"': 'constants.TestPerformanceCheck',
            '"Data Constraints"': 'constants.TestDataConstraints',
            '"Timestamp Consistency"': 'constants.TestTimestampConsistency',
            '"NULL Values"': 'constants.TestNullValues',
            '"Unique Constraints"': 'constants.TestUniqueConstraints',
            '"URL Quality"': 'constants.TestURLQuality',
            '"HTTP Status Codes"': 'constants.TestHTTPStatusCodes',
            '"SEO Score Validity"': 'constants.TestSEOScoreValidity',
            '"Query Performance"': 'constants.TestQueryPerformance',
            '"Database Size"': 'constants.TestDatabaseSize',
            '"Numeric Consistency"': 'constants.TestNumericConsistency',
            
            # Database related
            '"sqlite3"': 'constants.SQLite3Driver',
            '"fire_salamander_dev.db"': 'constants.DefaultDatabasePath',
            '"tests/reports/data"': 'constants.DefaultReportPath',
            
            # Table names  
            '"crawl_sessions"': 'constants.TableCrawlSessions',
            '"pages"': 'constants.TablePages',
            '"seo_metrics"': 'constants.TableSEOMetrics',
            
            # Messages
            '"Missing required table"': 'constants.MsgMissingRequiredTable',
            '"Application functionality may be impaired"': 'constants.MsgApplicationImpaired',
            '"Data inconsistency detected"': 'constants.MsgDataInconsistency',
            '"Referential integrity violation"': 'constants.MsgReferentialIntegrityFail',
            '"Performance degradation detected"': 'constants.MsgPerformanceDegradation',
            '"NULL or empty URL values found"': 'constants.MsgNullOrEmptyURL',
            '"Crawl sessions cannot function without valid URLs"': 'constants.MsgCrawlSessionsNoURL',
            '"Duplicate page records found"': 'constants.MsgDuplicatePageRecords',
            '"May cause inconsistent crawl results"': 'constants.MsgInconsistentCrawlResults',
            '"No duplicate page records found"': 'constants.MsgNoDuplicateRecords',
            '"All timestamps are consistent"': 'constants.MsgAllTimestampsConsistent',
            '"Page count relationships are consistent"': 'constants.MsgPageCountConsistent',
            '"All URLs appear to be well-formed"': 'constants.MsgAllURLsWellFormed',
            '"All status codes are in valid range"': 'constants.MsgAllStatusCodesValid',
            '"All SEO scores are within valid range"': 'constants.MsgAllSEOScoresValid',
            '"No orphaned records found"': 'constants.MsgNoOrphanedRecords',
            '"Incorrect session duration calculations"': 'constants.MsgIncorrectSessionDurationCalculations',
            '"Data inconsistency and potential application errors"': 'constants.MsgDataInconsistencyAndPotentialErrors',
            '"Simple Count"': 'constants.QueryNameSimpleCount',
            '"Complex Join"': 'constants.QueryNameComplexJoin',
            '"excellent"': 'constants.StatusExcellent',
            '"acceptable"': 'constants.StatusAcceptable',
            '"CHECK"': 'constants.SQLKeywordCHECK',
            '"header"': 'constants.HTMLClassHeader',
            '"section"': 'constants.HTMLClassSection',
            
            # Debug and Phase Test Messages
            '"Configuration Loading"': 'constants.TestConfigurationLoading',
            '"File Structure"': 'constants.TestFileStructure',
            '"Git Repository"': 'constants.TestGitRepository',
            '"HTTP Server"': 'constants.TestHTTPServer',
            '"SEPTEO Branding"': 'constants.TestSEPTEOBranding',
            '"Docker Setup"': 'constants.TestDockerSetup',
            '"Deploy Scripts"': 'constants.TestDeployScripts',
            '"Phase 1 - Setup Initial"': 'constants.Phase1SetupInitial',
            '"Verify configuration is loaded correctly with all required fields"': 'constants.DescConfigurationLoading',
            '"Verify all required files and directories are present"': 'constants.DescFileStructure',
            '"Verify git repository is properly initialized"': 'constants.DescGitRepository',
            '"Test server endpoints are responding correctly"': 'constants.DescHTTPServer',
            '"Verify SEPTEO branding is properly integrated"': 'constants.DescSEPTEOBranding',
            '"Verify Docker Compose configuration exists and is valid"': 'constants.DescDockerSetup',
            '"Verify deployment scripts exist and are executable"': 'constants.DescDeployScripts',
            '"Configuration is nil"': 'constants.ErrConfigurationNil',
            '"Configuration is valid and complete"': 'constants.MsgConfigurationValid',
            '"All required files and directories exist"': 'constants.MsgAllFilesExist',
            '"Git repository properly initialized"': 'constants.MsgGitProperlyInitialized',
            '"All server endpoints responding correctly"': 'constants.MsgAllEndpointsResponding',
            '"SEPTEO branding properly integrated"': 'constants.MsgSEPTEOBrandingIntegrated',
            '"Docker Compose configuration is valid"': 'constants.MsgDockerComposeValid',
            '"All deployment scripts exist and are executable"': 'constants.MsgAllDeployScriptsReady',
            '"Endpoint responding correctly"': 'constants.MsgEndpointRespondingCorrect',
        }
        
        # Special patterns for more complex replacements
        self.pattern_mappings = [
            # SQL queries that can be replaced
            (r'"SELECT COUNT\(\*\) FROM sqlite_master WHERE type=\'table\' AND name=\?"', 
             'constants.QueryTableExists'),
        ]

    def create_backup(self) -> None:
        """Create backup of original file"""
        with open(self.file_path, 'r') as src:
            with open(self.backup_path, 'w') as dst:
                dst.write(src.read())
        print(f"📋 Backup created: {self.backup_path}")

    def count_violations(self, content: str) -> int:
        """Count hardcoded string violations in content"""
        # Pattern matches strings with 5+ characters, not in comments or imports
        pattern = r'"[A-Za-z][A-Za-z0-9 ]{4,}"'
        matches = re.findall(pattern, content)
        
        # Filter out imports, comments, etc.
        violations = 0
        for match in matches:
            if not any(exclusion in match.lower() for exclusion in ['import', '//', 'const', 'var', 'fmt.']):
                violations += 1
        
        return violations

    def eliminate_hardcoding(self) -> Tuple[int, int]:
        """Eliminate hardcoding violations and return before/after count"""
        
        with open(self.file_path, 'r') as f:
            content = f.read()
        
        original_violations = self.count_violations(content)
        print(f"📊 Original violations: {original_violations}")
        
        # Apply string mappings
        for hardcoded, constant in self.string_mappings.items():
            if hardcoded in content:
                content = content.replace(hardcoded, constant)
                print(f"🔄 Replaced: {hardcoded} → {constant}")
        
        # Apply pattern mappings
        for pattern, replacement in self.pattern_mappings:
            matches = re.findall(pattern, content)
            if matches:
                content = re.sub(pattern, replacement, content)
                print(f"🔄 Pattern replaced: {pattern[:50]}... → {replacement}")
        
        # Write modified content
        with open(self.file_path, 'w') as f:
            f.write(content)
        
        final_violations = self.count_violations(content)
        eliminated = original_violations - final_violations
        
        print(f"📊 Final violations: {final_violations}")
        print(f"🎯 Eliminated: {eliminated}")
        
        return original_violations, final_violations

    def validate_compilation(self) -> bool:
        """Validate that the file still compiles"""
        try:
            result = subprocess.run(['go', 'build', '-o', '/dev/null', self.file_path], 
                                  capture_output=True, text=True)
            if result.returncode == 0:
                print("✅ File compiles successfully")
                return True
            else:
                print(f"❌ Compilation failed: {result.stderr}")
                return False
        except Exception as e:
            print(f"❌ Compilation check failed: {e}")
            return False

    def restore_backup(self) -> None:
        """Restore from backup if compilation fails"""
        with open(self.backup_path, 'r') as src:
            with open(self.file_path, 'w') as dst:
                dst.write(src.read())
        print("🔄 Restored from backup")

    def cleanup_backup(self) -> None:
        """Remove backup file"""
        if os.path.exists(self.backup_path):
            os.remove(self.backup_path)
            print("🧹 Backup cleaned up")

    def process(self) -> Dict[str, int]:
        """Main processing method"""
        print(f"🤖 SMART HARDCODE ELIMINATOR - Processing: {self.file_path}")
        
        # Create backup
        self.create_backup()
        
        try:
            # Eliminate hardcoding
            original, final = self.eliminate_hardcoding()
            
            # Validate compilation
            if self.validate_compilation():
                self.cleanup_backup()
                print(f"🎉 SUCCESS: {original - final} violations eliminated!")
                return {
                    'status': 'success',
                    'original': original,
                    'final': final,
                    'eliminated': original - final
                }
            else:
                self.restore_backup()
                print("❌ FAILED: Compilation errors, restored backup")
                return {
                    'status': 'failed',
                    'original': original,
                    'final': original,
                    'eliminated': 0
                }
        
        except Exception as e:
            print(f"❌ ERROR: {e}")
            self.restore_backup()
            return {
                'status': 'error',
                'original': 0,
                'final': 0,
                'eliminated': 0
            }

def main():
    if len(sys.argv) != 2:
        print("Usage: python3 smart_hardcode_eliminator.py <go_file>")
        sys.exit(1)
    
    file_path = sys.argv[1]
    
    if not os.path.exists(file_path):
        print(f"❌ File not found: {file_path}")
        sys.exit(1)
    
    eliminator = SmartHardcodeEliminator(file_path)
    result = eliminator.process()
    
    print("\n📊 FINAL RESULT:")
    print(f"Status: {result['status']}")
    print(f"Original violations: {result['original']}")  
    print(f"Final violations: {result['final']}")
    print(f"Eliminated: {result['eliminated']}")
    
    # Return appropriate exit code
    sys.exit(0 if result['status'] == 'success' else 1)

if __name__ == "__main__":
    main()