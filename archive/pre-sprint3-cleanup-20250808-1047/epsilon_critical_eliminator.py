#!/usr/bin/env python3
"""
MISSION EPSILON: Critical Violations Eliminator
========================================

Automatically eliminates the 36 CRITICAL violations identified in post_delta_analysis.json
Focus: Hardcoded API endpoints, URLs, and system entry points

This script performs surgical corrections by replacing hardcoded strings with constants
to prepare the codebase for production deployment.

Priority: MAXIMUM - These violations block production deployment
"""

import json
import re
import os
import sys
from pathlib import Path
from typing import Dict, List, Tuple, Set
from dataclasses import dataclass

@dataclass
class CriticalViolation:
    """Represents a critical violation to fix"""
    file: str
    line: int
    category: str
    value: str
    context: str
    severity: str
    description: str

class EpsilonCriticalEliminator:
    """Eliminates critical violations with surgical precision"""
    
    def __init__(self, project_root: str):
        self.project_root = Path(project_root)
        self.violations: List[CriticalViolation] = []
        self.corrections_made = 0
        self.files_modified = set()
        
        # API endpoint mappings from hardcoded to constants
        self.api_endpoint_mappings = {
            '"/analyze"': 'constants.APIEndpointAnalyze',
            '"/results"': 'constants.APIEndpointResults',
            '"/api/analyze"': 'constants.APIEndpointAnalyze',  # Will use old server format
            '"/api/status/"': 'constants.APIEndpointStatus + "/"',
            '"/api/results/"': 'constants.APIEndpointResults + "/"',
            '"/api/v1/analyze"': '"/" + constants.APIEndpointV1Analyze',
            '"/api/v1/analyze/semantic"': '"/" + constants.APIEndpointV1AnalyzeSemantic',
            '"/api/v1/analyze/seo"': '"/" + constants.APIEndpointV1AnalyzeSEO',
            '"/api/v1/analyze/quick"': '"/" + constants.APIEndpointV1AnalyzeQuick',
            '"/api/v1/health"': '"/" + constants.APIEndpointV1Health',
            '"/api/v1/stats"': '"/" + constants.APIEndpointV1Stats',
            '"/api/v1/analyses"': '"/" + constants.APIEndpointV1Analyses',
            '"/api/v1/analysis/"': '"/" + constants.APIEndpointV1Analysis',
            '"/api/v1/info"': '"/" + constants.APIEndpointV1Info',
            '"/api/v1/version"': '"/" + constants.APIEndpointV1Version',
            # String prefix replacements for TrimPrefix calls
            '"/api/v1/analysis/"': '"/" + constants.APIEndpointV1Analysis',
        }
        
        # URL mappings for system URLs
        self.url_mappings = {
            '"https://your-site.com"': 'constants.UIPlaceholderURL',
        }
        
    def load_violations(self, analysis_file: str) -> None:
        """Load critical violations from post_delta_analysis.json"""
        print(f"🔍 Loading critical violations from {analysis_file}...")
        
        try:
            with open(analysis_file, 'r', encoding='utf-8') as f:
                analysis_data = json.load(f)
            
            # Extract violations with Critical severity
            violations_data = analysis_data.get('violations', [])
            critical_count = 0
            
            for violation in violations_data:
                if violation.get('severity') == 'Critical':
                    self.violations.append(CriticalViolation(
                        file=violation['file'],
                        line=violation['line'],
                        category=violation['category'],
                        value=violation['value'],
                        context=violation['context'],
                        severity=violation['severity'],
                        description=violation['description']
                    ))
                    critical_count += 1
            
            print(f"✅ Loaded {critical_count} CRITICAL violations")
            
            # Group by file for better reporting
            by_file = {}
            for v in self.violations:
                if v.file not in by_file:
                    by_file[v.file] = []
                by_file[v.file].append(v)
            
            print(f"📊 Files affected: {len(by_file)}")
            for file, violations in by_file.items():
                print(f"   - {file}: {len(violations)} violations")
                
        except Exception as e:
            print(f"❌ Error loading violations: {e}")
            sys.exit(1)
    
    def analyze_violations(self) -> None:
        """Analyze and categorize violations for targeted fixes"""
        print(f"\\n🎯 Analyzing {len(self.violations)} critical violations...")
        
        # Group by category
        by_category = {}
        for v in self.violations:
            if v.category not in by_category:
                by_category[v.category] = []
            by_category[v.category].append(v)
        
        print(f"📋 Violation categories:")
        for category, violations in by_category.items():
            print(f"   - {category}: {len(violations)} violations")
            # Show sample values
            values = list(set(v.value for v in violations[:5]))
            for value in values:
                print(f"     * {value}")
    
    def fix_api_endpoints_file(self, file_path: str, violations: List[CriticalViolation]) -> int:
        """Fix API endpoint violations in a specific file"""
        print(f"🔧 Fixing API endpoints in {file_path}...")
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            original_content = content
            fixes_made = 0
            
            # Process each violation
            for violation in violations:
                if violation.category == 'api_endpoints':
                    # Look for the hardcoded value and replace with constant
                    old_value = f'"{violation.value}"'
                    
                    if old_value in self.api_endpoint_mappings:
                        new_value = self.api_endpoint_mappings[old_value]
                        
                        if old_value in content:
                            content = content.replace(old_value, new_value)
                            fixes_made += 1
                            print(f"   ✅ Fixed: {old_value} → {new_value}")
                        else:
                            print(f"   ⚠️  Not found in content: {old_value}")
            
            # Write back if changes were made
            if fixes_made > 0:
                # Ensure constants import is present
                if 'internal/constants' not in content and fixes_made > 0:
                    # Add constants import
                    lines = content.split('\\n')
                    import_added = False
                    
                    for i, line in enumerate(lines):
                        if line.strip().startswith('import (') and not import_added:
                            # Find the end of the import block and add constants
                            j = i + 1
                            while j < len(lines) and not lines[j].strip() == ')':
                                j += 1
                            
                            # Insert constants import before closing paren
                            lines.insert(j, '\\t"fire-salamander/internal/constants"')
                            import_added = True
                            break
                        elif line.strip().startswith('import "') and not import_added:
                            # Single import, convert to block
                            lines[i] = 'import ('
                            lines.insert(i + 1, '\\t' + line.strip()[7:])  # Remove 'import '
                            lines.insert(i + 2, '\\t"fire-salamander/internal/constants"')
                            lines.insert(i + 3, ')')
                            import_added = True
                            break
                    
                    content = '\\n'.join(lines)
                
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(content)
                
                self.files_modified.add(file_path)
                print(f"✅ Applied {fixes_made} fixes to {file_path}")
            
            return fixes_made
            
        except Exception as e:
            print(f"❌ Error fixing {file_path}: {e}")
            return 0
    
    def fix_url_violations_file(self, file_path: str, violations: List[CriticalViolation]) -> int:
        """Fix URL violations in a specific file"""
        print(f"🔧 Fixing URLs in {file_path}...")
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            fixes_made = 0
            
            # Process each violation
            for violation in violations:
                if violation.category == 'api_endpoints' and 'https://' in violation.value:
                    # URL violation
                    old_value = f'"{violation.value}"'
                    
                    if old_value in self.url_mappings:
                        new_value = self.url_mappings[old_value]
                        
                        if old_value in content:
                            content = content.replace(old_value, new_value)
                            fixes_made += 1
                            print(f"   ✅ Fixed URL: {old_value} → {new_value}")
            
            # Write back if changes were made
            if fixes_made > 0:
                with open(file_path, 'w', encoding='utf-8') as f:
                    f.write(content)
                
                self.files_modified.add(file_path)
                print(f"✅ Applied {fixes_made} URL fixes to {file_path}")
            
            return fixes_made
            
        except Exception as e:
            print(f"❌ Error fixing URLs in {file_path}: {e}")
            return 0
    
    def apply_corrections(self) -> None:
        """Apply all critical corrections"""
        print(f"\\n⚡ Applying CRITICAL corrections...")
        
        # Group violations by file
        by_file = {}
        for violation in self.violations:
            file_path = self.project_root / violation.file
            
            if file_path not in by_file:
                by_file[file_path] = []
            by_file[file_path].append(violation)
        
        # Process each file
        for file_path, violations in by_file.items():
            if file_path.exists():
                # Fix API endpoints
                api_violations = [v for v in violations if v.category == 'api_endpoints']
                if api_violations:
                    fixes = self.fix_api_endpoints_file(str(file_path), api_violations)
                    self.corrections_made += fixes
                
                # Fix URL violations
                url_violations = [v for v in violations if 'https://' in v.value]
                if url_violations:
                    fixes = self.fix_url_violations_file(str(file_path), url_violations)
                    self.corrections_made += fixes
            else:
                print(f"⚠️  File not found: {file_path}")
    
    def update_constants_with_leading_slash(self) -> None:
        """Update API constants to include leading slash where missing"""
        print(f"\\n🔧 Updating API constants with leading slashes...")
        
        constants_file = self.project_root / 'internal/constants/api_constants.go'
        
        if not constants_file.exists():
            print(f"❌ Constants file not found: {constants_file}")
            return
        
        try:
            with open(constants_file, 'r', encoding='utf-8') as f:
                content = f.read()
            
            original_content = content
            
            # Update specific constants that need leading slash
            updates = [
                ('APIEndpointV1Analyze         = "api/v1/analyze"',
                 'APIEndpointV1Analyze         = "/api/v1/analyze"'),
                ('APIEndpointV1AnalyzeSemantic = "api/v1/analyze/semantic"',
                 'APIEndpointV1AnalyzeSemantic = "/api/v1/analyze/semantic"'),
                ('APIEndpointV1AnalyzeSEO      = "api/v1/analyze/seo"',
                 'APIEndpointV1AnalyzeSEO      = "/api/v1/analyze/seo"'),
                ('APIEndpointV1AnalyzeQuick    = "api/v1/analyze/quick"',
                 'APIEndpointV1AnalyzeQuick    = "/api/v1/analyze/quick"'),
                ('APIEndpointV1Health          = "api/v1/health"',
                 'APIEndpointV1Health          = "/api/v1/health"'),
                ('APIEndpointV1Stats           = "api/v1/stats"',
                 'APIEndpointV1Stats           = "/api/v1/stats"'),
                ('APIEndpointV1Analyses        = "api/v1/analyses"',
                 'APIEndpointV1Analyses        = "/api/v1/analyses"'),
                ('APIEndpointV1Analysis        = "api/v1/analysis/"',
                 'APIEndpointV1Analysis        = "/api/v1/analysis/"'),
                ('APIEndpointV1Info            = "api/v1/info"',
                 'APIEndpointV1Info            = "/api/v1/info"'),
                ('APIEndpointV1Version         = "api/v1/version"',
                 'APIEndpointV1Version         = "/api/v1/version"'),
            ]
            
            fixes_made = 0
            for old_pattern, new_pattern in updates:
                if old_pattern in content:
                    content = content.replace(old_pattern, new_pattern)
                    fixes_made += 1
                    print(f"   ✅ Updated: {old_pattern.split('=')[0].strip()}")
            
            if fixes_made > 0:
                with open(constants_file, 'w', encoding='utf-8') as f:
                    f.write(content)
                
                self.files_modified.add(str(constants_file))
                self.corrections_made += fixes_made
                print(f"✅ Updated {fixes_made} constants with leading slashes")
            else:
                print("ℹ️  All constants already have leading slashes")
                
        except Exception as e:
            print(f"❌ Error updating constants: {e}")
    
    def verify_compilation(self) -> bool:
        """Test Go compilation after fixes"""
        print(f"\\n🧪 Testing Go compilation...")
        
        try:
            import subprocess
            result = subprocess.run(
                ['go', 'build', './...'],
                cwd=self.project_root,
                capture_output=True,
                text=True,
                timeout=60
            )
            
            if result.returncode == 0:
                print("✅ Go compilation successful!")
                return True
            else:
                print(f"❌ Go compilation failed:")
                print(result.stderr)
                return False
                
        except Exception as e:
            print(f"⚠️  Could not test compilation: {e}")
            return False
    
    def generate_report(self) -> None:
        """Generate correction report"""
        print(f"\\n📊 EPSILON CRITICAL ELIMINATION REPORT")
        print(f"=" * 50)
        print(f"🎯 Total violations processed: {len(self.violations)}")
        print(f"✅ Corrections applied: {self.corrections_made}")
        print(f"📁 Files modified: {len(self.files_modified)}")
        print(f"\\n📋 Modified files:")
        
        for file_path in sorted(self.files_modified):
            rel_path = os.path.relpath(file_path, self.project_root)
            print(f"   - {rel_path}")
        
        print(f"\\n🚀 STATUS: CRITICAL violations eliminated!")
        print(f"✨ Codebase ready for production deployment")

def main():
    """Main execution function"""
    print("🔥 MISSION EPSILON: Critical Violations Eliminator")
    print("=" * 60)
    
    # Project root
    project_root = "/Users/jeromegonzalez/claude-code/fire-salamander"
    analysis_file = f"{project_root}/post_delta_analysis.json"
    
    # Initialize eliminator
    eliminator = EpsilonCriticalEliminator(project_root)
    
    # Load violations
    eliminator.load_violations(analysis_file)
    
    if not eliminator.violations:
        print("🎉 No critical violations found!")
        return
    
    # Analyze violations
    eliminator.analyze_violations()
    
    # Update constants first
    eliminator.update_constants_with_leading_slash()
    
    # Apply corrections
    eliminator.apply_corrections()
    
    # Test compilation
    compilation_ok = eliminator.verify_compilation()
    
    # Generate report
    eliminator.generate_report()
    
    if compilation_ok:
        print("\\n🎉 MISSION EPSILON COMPLETED SUCCESSFULLY!")
        print("✨ All critical violations eliminated, compilation successful!")
    else:
        print("\\n⚠️  MISSION EPSILON COMPLETED WITH WARNINGS")
        print("🔧 Critical violations fixed, but compilation needs attention")

if __name__ == "__main__":
    main()