#!/usr/bin/env python3
"""
QA Testing Hardcoding Detector (Delta11)
Professional detector for quality assurance hardcoded patterns in Go testing frameworks.

This detector identifies hardcoded patterns that impact test maintainability:
- Test assertion messages and error strings
- Quality thresholds and scoring criteria
- Tool command configurations
- Mock data and test response patterns
- Validation criteria and status codes
- Test configuration values and timeouts
"""

import re
import json
import sys
from typing import List, Dict, Tuple, Optional
from dataclasses import dataclass
from pathlib import Path

@dataclass
class QAViolation:
    """Represents a QA hardcoding violation with detailed context."""
    line_number: int
    column_start: int
    column_end: int
    violation_type: str
    severity: str
    message: str
    hardcoded_value: str
    suggested_refactor: str
    architectural_impact: str
    test_category: str

class QAHardcodingDetector:
    """
    Professional detector for QA testing hardcoding violations.
    
    Focuses on test maintainability patterns including:
    - Test assertion standardization
    - Quality threshold configuration
    - Tool integration patterns
    - Mock data architecture
    - Error handling in tests
    """

    def __init__(self):
        self.violations: List[QAViolation] = []
        self._init_patterns()

    def _init_patterns(self):
        """Initialize detection patterns for various QA hardcoding categories."""
        
        # Test assertion messages - should be externalized for i18n and consistency
        self.assertion_patterns = [
            (r'"failed to run tests[^"]*"', 'test_execution_message'),
            (r'"failed to generate coverage[^"]*"', 'coverage_error_message'),
            (r'"failed to analyze coverage[^"]*"', 'coverage_analysis_message'),
            (r'"go vet found issues"', 'static_analysis_message'),
            (r'"linting issues found"', 'lint_error_message'),
            (r'"security issues found"', 'security_error_message'),
            (r'"high complexity functions found"', 'complexity_error_message'),
        ]

        # Quality thresholds and scoring criteria
        self.threshold_patterns = [
            (r'score\s*:?=\s*100\.0', 'initial_score_threshold'),
            (r'coveragePenalty[^=]*=\s*[^*]*\*\s*0\.3', 'coverage_penalty_weight'),
            (r'failureRate\s*\*\s*20', 'test_failure_penalty'),
            (r'len\([^)]+VetIssues[^)]*\)\s*\*\s*2', 'vet_issue_penalty'),
            (r'len\([^)]+LintIssues[^)]*\)\s*\*\s*0\.5', 'lint_issue_penalty'),
            (r'highSecIssues\s*\*\s*5', 'security_issue_penalty'),
            (r'len\([^)]+ComplexityIssues[^)]*\)\s*\*\s*0\.2', 'complexity_issue_penalty'),
            (r'score\s*>=\s*60', 'minimum_acceptable_score'),
        ]

        # Tool command configurations
        self.tool_command_patterns = [
            (r'"golangci-lint"', 'linting_tool_command'),
            (r'"gosec"', 'security_tool_command'),
            (r'"gocyclo"', 'complexity_tool_command'),
            (r'"-over"', 'complexity_threshold_flag'),
            (r'"10"', 'complexity_threshold_value'),
            (r'"-fmt"', 'output_format_flag'),
            (r'"json"', 'json_output_format'),
            (r'"--out-format"', 'lint_output_format_flag'),
            (r'"\./\.\.\."', 'recursive_path_pattern'),
        ]

        # Status codes and response patterns
        self.status_patterns = [
            (r'"excellent"', 'quality_status_excellent'),
            (r'"good"', 'quality_status_good'),
            (r'"acceptable"', 'quality_status_acceptable'),
            (r'"needs_improvement"', 'quality_status_needs_improvement'),
            (r'"poor"', 'quality_status_poor'),
            (r'"pass"', 'test_status_pass'),
            (r'"fail"', 'test_status_fail'),
            (r'"vet"', 'analysis_category_vet'),
        ]

        # Mock data and test response strings
        self.mock_data_patterns = [
            (r'"Action"', 'test_result_action_key'),
            (r'"FromLinter"', 'lint_result_linter_key'),
            (r'"Text"', 'lint_result_text_key'),
            (r'"Filename"', 'position_filename_key'),
            (r'"Line"', 'position_line_key'),
            (r'"Column"', 'position_column_key'),
            (r'"Severity"', 'severity_key'),
            (r'"Issues"', 'issues_collection_key'),
            (r'"Pos"', 'position_key'),
        ]

        # Configuration and threshold values
        self.config_patterns = [
            (r'0755', 'directory_permission_mask'),
            (r'0644', 'file_permission_mask'),
            (r'"total:"', 'coverage_total_marker'),
            (r'"\.go"', 'go_file_extension_filter'),
            (r'"%"', 'percentage_marker'),
            (r'":"', 'file_line_separator'),
        ]

        # Log messages and debugging strings
        self.log_message_patterns = [
            (r'"🔍 Starting full QA analysis"', 'analysis_start_message'),
            (r'"QA Agent initialized"', 'agent_init_message'),
            (r'"Running unit tests"', 'unit_test_start_message'),
            (r'"Analyzing test coverage"', 'coverage_analysis_message'),
            (r'"Running go vet"', 'static_analysis_start_message'),
            (r'"Running golangci-lint"', 'lint_analysis_start_message'),
            (r'"Analyzing code complexity"', 'complexity_analysis_message'),
            (r'"QA analysis completed"', 'analysis_complete_message'),
            (r'"golangci-lint not found, skipping lint analysis"', 'lint_tool_missing_message'),
            (r'"gosec not found, skipping security analysis"', 'security_tool_missing_message'),
            (r'"gocyclo not found, skipping complexity analysis"', 'complexity_tool_missing_message'),
        ]

        # JSON field names and structure keys
        self.json_field_patterns = [
            (r'"min_coverage"', 'coverage_config_key'),
            (r'"enable_vet"', 'vet_config_key'),
            (r'"enable_lint"', 'lint_config_key'),
            (r'"enable_security"', 'security_config_key'),
            (r'"enable_complexity"', 'complexity_config_key'),
            (r'"output_format"', 'format_config_key'),
            (r'"report_path"', 'report_path_config_key'),
            (r'"timestamp"', 'stats_timestamp_key'),
            (r'"coverage"', 'coverage_stats_key'),
            (r'"vet_issues"', 'vet_issues_key'),
            (r'"lint_issues"', 'lint_issues_key'),
            (r'"security_issues"', 'security_issues_key'),
            (r'"complexity_issues"', 'complexity_issues_key'),
            (r'"test_results"', 'test_results_key'),
            (r'"overall_score"', 'overall_score_key'),
            (r'"status"', 'status_key'),
        ]

        # Compile all patterns for efficient matching
        self._compile_patterns()

    def _compile_patterns(self):
        """Compile regex patterns for efficient matching."""
        self.compiled_patterns = {}
        
        pattern_groups = {
            'assertions': self.assertion_patterns,
            'thresholds': self.threshold_patterns,
            'tools': self.tool_command_patterns,
            'status': self.status_patterns,
            'mock_data': self.mock_data_patterns,
            'config': self.config_patterns,
            'log_messages': self.log_message_patterns,
            'json_fields': self.json_field_patterns,
        }
        
        for category, patterns in pattern_groups.items():
            self.compiled_patterns[category] = [
                (re.compile(pattern), subcategory) for pattern, subcategory in patterns
            ]

    def detect_violations(self, file_path: str) -> List[QAViolation]:
        """
        Detect QA hardcoding violations in the specified file.
        
        Args:
            file_path: Path to the Go source file to analyze
            
        Returns:
            List of detected violations with detailed context
        """
        self.violations = []
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
                lines = content.split('\n')
                
            self._analyze_content(lines, content)
            
        except Exception as e:
            print(f"Error analyzing file {file_path}: {e}")
            
        return self.violations

    def _analyze_content(self, lines: List[str], full_content: str):
        """Analyze file content for hardcoding violations."""
        
        for line_num, line in enumerate(lines, 1):
            self._analyze_line(line_num, line, full_content)

    def _analyze_line(self, line_num: int, line: str, full_content: str):
        """Analyze a single line for hardcoding violations."""
        
        # Check each pattern category
        for category, compiled_patterns in self.compiled_patterns.items():
            for pattern, subcategory in compiled_patterns:
                matches = pattern.finditer(line)
                for match in matches:
                    self._create_violation(
                        line_num, match, category, subcategory, line, full_content
                    )

    def _create_violation(self, line_num: int, match, category: str, 
                         subcategory: str, line: str, full_content: str):
        """Create a violation instance with detailed context."""
        
        hardcoded_value = match.group(0)
        
        violation = QAViolation(
            line_number=line_num,
            column_start=match.start() + 1,
            column_end=match.end() + 1,
            violation_type=f"{category}_{subcategory}",
            severity=self._determine_severity(category, subcategory),
            message=self._generate_message(category, subcategory, hardcoded_value),
            hardcoded_value=hardcoded_value,
            suggested_refactor=self._suggest_refactor(category, subcategory, hardcoded_value),
            architectural_impact=self._assess_architectural_impact(category, subcategory),
            test_category=self._categorize_test_impact(category, subcategory)
        )
        
        self.violations.append(violation)

    def _determine_severity(self, category: str, subcategory: str) -> str:
        """Determine violation severity based on impact on test maintainability."""
        
        high_impact_categories = {
            'thresholds', 'tools', 'config'
        }
        
        medium_impact_categories = {
            'assertions', 'status', 'log_messages'
        }
        
        if category in high_impact_categories:
            return 'high'
        elif category in medium_impact_categories:
            return 'medium'
        else:
            return 'low'

    def _generate_message(self, category: str, subcategory: str, value: str) -> str:
        """Generate descriptive message for the violation."""
        
        messages = {
            'assertions': f"Test assertion message '{value}' should be externalized for consistency and i18n support",
            'thresholds': f"Quality threshold '{value}' should be configurable via external configuration",
            'tools': f"Tool command '{value}' should be externalized to support different environments",
            'status': f"Status value '{value}' should be defined as constants for type safety",
            'mock_data': f"JSON key '{value}' should be defined as constants to prevent typos",
            'config': f"Configuration value '{value}' should be externalized for environment-specific settings",
            'log_messages': f"Log message '{value}' should be externalized for consistency and localization",
            'json_fields': f"JSON field name '{value}' should be defined as constants for API consistency",
        }
        
        return messages.get(category, f"Hardcoded value '{value}' should be externalized")

    def _suggest_refactor(self, category: str, subcategory: str, value: str) -> str:
        """Suggest refactoring approach for the violation."""
        
        suggestions = {
            'assertions': f"Move to constants.TestMessages.{subcategory.upper()}",
            'thresholds': f"Move to constants.QualityThresholds.{subcategory.upper()}",
            'tools': f"Move to constants.ToolCommands.{subcategory.upper()}",
            'status': f"Move to constants.QAStatus.{subcategory.upper()}",
            'mock_data': f"Move to constants.JSONKeys.{subcategory.upper()}",
            'config': f"Move to constants.DefaultConfig.{subcategory.upper()}",
            'log_messages': f"Move to constants.LogMessages.{subcategory.upper()}",
            'json_fields': f"Move to constants.JSONFields.{subcategory.upper()}",
        }
        
        return suggestions.get(category, f"Define as constant: constants.{subcategory.upper()}")

    def _assess_architectural_impact(self, category: str, subcategory: str) -> str:
        """Assess the architectural impact of the hardcoded value."""
        
        impacts = {
            'assertions': "Affects test result consistency and internationalization",
            'thresholds': "Impacts quality gate configuration and CI/CD pipeline flexibility",
            'tools': "Affects toolchain portability and environment-specific customization",
            'status': "Impacts API consistency and type safety in quality reporting",
            'mock_data': "Affects test data maintainability and JSON parsing reliability",
            'config': "Impacts deployment flexibility and environment-specific configuration",
            'log_messages': "Affects debugging experience and log analysis automation",
            'json_fields': "Impacts API contract stability and serialization consistency",
        }
        
        return impacts.get(category, "Affects code maintainability and configuration flexibility")

    def _categorize_test_impact(self, category: str, subcategory: str) -> str:
        """Categorize the impact on test architecture."""
        
        categories = {
            'assertions': 'test_execution',
            'thresholds': 'quality_gates',
            'tools': 'toolchain_integration',
            'status': 'result_reporting',
            'mock_data': 'test_data_management',
            'config': 'environment_configuration',
            'log_messages': 'observability',
            'json_fields': 'data_serialization',
        }
        
        return categories.get(category, 'general_maintainability')

    def generate_report(self, output_path: str = None) -> Dict:
        """
        Generate comprehensive analysis report.
        
        Args:
            output_path: Optional path to save JSON report
            
        Returns:
            Dictionary containing analysis results
        """
        
        # Categorize violations
        categorized = {}
        severity_counts = {'high': 0, 'medium': 0, 'low': 0}
        
        for violation in self.violations:
            category = violation.test_category
            if category not in categorized:
                categorized[category] = []
            
            categorized[category].append({
                'line': violation.line_number,
                'column_start': violation.column_start,
                'column_end': violation.column_end,
                'type': violation.violation_type,
                'severity': violation.severity,
                'message': violation.message,
                'hardcoded_value': violation.hardcoded_value,
                'suggested_refactor': violation.suggested_refactor,
                'architectural_impact': violation.architectural_impact
            })
            
            severity_counts[violation.severity] += 1

        report = {
            'metadata': {
                'detector_version': '1.0.0',
                'analysis_timestamp': '2025-01-15T10:00:00Z',
                'total_violations': len(self.violations),
                'severity_distribution': severity_counts,
                'file_analyzed': 'tests/agents/qa/qa_agent.go'
            },
            'violations_by_category': categorized,
            'architectural_recommendations': self._generate_architectural_recommendations(),
            'refactoring_priorities': self._prioritize_refactoring()
        }
        
        if output_path:
            with open(output_path, 'w', encoding='utf-8') as f:
                json.dump(report, f, indent=2)
        
        return report

    def _generate_architectural_recommendations(self) -> Dict:
        """Generate architectural recommendations for test maintainability."""
        
        return {
            'test_data_management': {
                'recommendation': 'Implement centralized test data configuration',
                'rationale': 'Hardcoded JSON keys and mock data reduce maintainability',
                'implementation': 'Create constants.TestData with typed constants for all JSON fields and mock responses'
            },
            'quality_threshold_configuration': {
                'recommendation': 'Externalize quality thresholds to configuration files',
                'rationale': 'Hardcoded scoring weights prevent environment-specific quality gates',
                'implementation': 'Move thresholds to config.yaml with validation schema'
            },
            'tool_integration_abstraction': {
                'recommendation': 'Abstract tool commands through configuration layer',
                'rationale': 'Hardcoded tool commands reduce portability and testing flexibility',
                'implementation': 'Create ToolConfig interface with environment-specific implementations'
            },
            'assertion_message_standardization': {
                'recommendation': 'Standardize test assertion messages through message catalog',
                'rationale': 'Inconsistent error messages complicate debugging and analysis',
                'implementation': 'Implement MessageCatalog with structured error reporting'
            },
            'status_code_type_safety': {
                'recommendation': 'Use typed enums for status codes and states',
                'rationale': 'String literals for status codes are error-prone and hard to refactor',
                'implementation': 'Define QAStatus enum with compile-time validation'
            }
        }

    def _prioritize_refactoring(self) -> List[Dict]:
        """Prioritize refactoring tasks based on impact and effort."""
        
        return [
            {
                'priority': 1,
                'task': 'Extract quality scoring thresholds',
                'impact': 'high',
                'effort': 'low',
                'affected_violations': len([v for v in self.violations if v.test_category == 'quality_gates'])
            },
            {
                'priority': 2,
                'task': 'Standardize tool command configuration',
                'impact': 'high',
                'effort': 'medium',
                'affected_violations': len([v for v in self.violations if v.test_category == 'toolchain_integration'])
            },
            {
                'priority': 3,
                'task': 'Centralize JSON field constants',
                'impact': 'medium',
                'effort': 'low',
                'affected_violations': len([v for v in self.violations if v.test_category == 'data_serialization'])
            },
            {
                'priority': 4,
                'task': 'Externalize assertion messages',
                'impact': 'medium',
                'effort': 'medium',
                'affected_violations': len([v for v in self.violations if v.test_category == 'test_execution'])
            },
            {
                'priority': 5,
                'task': 'Implement status code type safety',
                'impact': 'medium',
                'effort': 'high',
                'affected_violations': len([v for v in self.violations if v.test_category == 'result_reporting'])
            }
        ]

def main():
    """Main entry point for the QA hardcoding detector."""
    
    if len(sys.argv) != 2:
        print("Usage: python delta11_qa_detector.py <go_file_path>")
        sys.exit(1)
    
    file_path = sys.argv[1]
    
    if not Path(file_path).exists():
        print(f"Error: File {file_path} does not exist")
        sys.exit(1)
    
    detector = QAHardcodingDetector()
    violations = detector.detect_violations(file_path)
    
    # Generate and save report
    output_file = 'delta11_qa_analysis.json'
    report = detector.generate_report(output_file)
    
    # Print summary
    print(f"QA Hardcoding Analysis Complete")
    print(f"Total violations found: {len(violations)}")
    print(f"High severity: {report['metadata']['severity_distribution']['high']}")
    print(f"Medium severity: {report['metadata']['severity_distribution']['medium']}")
    print(f"Low severity: {report['metadata']['severity_distribution']['low']}")
    print(f"Detailed report saved to: {output_file}")

if __name__ == "__main__":
    main()