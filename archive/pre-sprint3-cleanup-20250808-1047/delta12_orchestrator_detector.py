#!/usr/bin/env python3
"""
Delta-12 Orchestrator Hardcoding Pattern Detector
Detects orchestration-specific hardcoding violations in Go files.

Focus Areas:
1. Service endpoint URLs
2. Timing and timeout values
3. Retry configurations
4. State management strings
5. Event names and types
6. Error handling messages
7. Workflow step names
8. Integration point definitions
"""

import re
import sys
import json
from typing import List, Dict, Any, Tuple
from dataclasses import dataclass, asdict
from pathlib import Path

@dataclass
class OrchestratorViolation:
    """Represents an orchestration hardcoding violation"""
    line_number: int
    line_content: str
    violation_type: str
    hardcoded_value: str
    severity: str
    category: str
    suggestion: str
    architecture_impact: str

class OrchestratorDetector:
    """Detects orchestration-specific hardcoding patterns"""
    
    def __init__(self):
        self.violations: List[OrchestratorViolation] = []
        self.patterns = self._initialize_patterns()
    
    def _initialize_patterns(self) -> Dict[str, Dict]:
        """Initialize orchestration-specific detection patterns"""
        return {
            'service_endpoints': {
                'patterns': [
                    r'"https?://[^"]+(?:api|service|endpoint)[^"]*"',
                    r'"localhost:\d+"',
                    r'"127\.0\.0\.1:\d+"',
                    r':\d+[^"]*/api[^"]*',
                    r'NewStorageManager\("([^"]+)"\)',
                    r'taskQueue:\s*make\(.*?,\s*(\d+)\)',
                    r'resultsChan:\s*make\(.*?,\s*(\d+)\)',
                ],
                'severity': 'HIGH',
                'category': 'Service Discovery'
            },
            'timing_values': {
                'patterns': [
                    r'time\.Duration\(\d+\)',
                    r'\d+\s*\*\s*time\.(Second|Minute|Hour|Millisecond|Microsecond)',
                    r'time\.Sleep\(\d+',
                    r'context\.WithTimeout\(.*?,\s*\d+',
                    r'Timeout:\s*(\d+)',
                    r'maxDepth\s*:=\s*(\d+)',
                    r'workers:\s*(\d+)',
                    r'\.Workers,?\s*$',
                    r'RetryAttempts:\s*(\d+)',
                    r'RetryDelay:\s*constants\.',
                ],
                'severity': 'HIGH',
                'category': 'Timing Configuration'
            },
            'retry_patterns': {
                'patterns': [
                    r'MaxRetries?\s*[:=]\s*(\d+)',
                    r'RetryAttempts?\s*[:=]\s*(\d+)',
                    r'for\s+i\s*:=\s*0;\s*i\s*<\s*(\d+)',
                    r'retry\.Times\(\d+\)',
                    r'attempts\s*<=\s*(\d+)',
                    r'RetryDelay:\s*constants\.',
                ],
                'severity': 'MEDIUM',
                'category': 'Retry Configuration'
            },
            'state_strings': {
                'patterns': [
                    r'Status:\s*"([^"]+)"',
                    r'TaskStatus(Pending|Running|Completed|Failed|Cancelled)',
                    r'AnalysisStatus(Success|Partial|Failed)',
                    r'Phase:\s*"([^"]+)"',
                    r'Type:\s*"([^"]+)"',
                    r'"(pending|running|completed|failed|cancelled)"',
                    r'constants\.OrchestratorStatus\w+',
                ],
                'severity': 'MEDIUM',
                'category': 'State Management'
            },
            'event_names': {
                'patterns': [
                    r'Event:\s*"([^"]+)"',
                    r'EventType:\s*"([^"]+)"',
                    r'log\.Printf\("([^"]+)"',
                    r'fmt\.Sprintf\("([^"]+)"',
                    r'constants\.Orchestrator\w+Message',
                    r'constants\.Orchestrator\w+Error',
                ],
                'severity': 'LOW',
                'category': 'Event Management'
            },
            'error_messages': {
                'patterns': [
                    r'fmt\.Errorf\("([^"]+)"',
                    r'errors\.New\("([^"]+)"',
                    r'return.*error.*"([^"]+)"',
                    r'panic\("([^"]+)"\)',
                    r'log\.Fatal\("([^"]+)"\)',
                    r'constants\.OrchestratorErr\w+',
                ],
                'severity': 'MEDIUM', 
                'category': 'Error Handling'
            },
            'workflow_steps': {
                'patterns': [
                    r'Step:\s*"([^"]+)"',
                    r'Phase:\s*"([^"]+)"',
                    r'Stage:\s*"([^"]+)"',
                    r'wg\.Add\((\d+)\)',
                    r'go func\(\) \{',
                    r'defer wg\.Done\(\)',
                    r'// \d+\. ([A-Z][^(]+)',
                ],
                'severity': 'LOW',
                'category': 'Workflow Management'
            },
            'integration_points': {
                'patterns': [
                    r'crawler\.\w+\(',
                    r'semanticAnalyzer\.\w+\(',
                    r'seoAnalyzer\.\w+\(',
                    r'storage\.\w+\(',
                    r'o\.\w+\(ctx,\s*task\)',
                    r'constants\.\w*Agent\w*',
                    r'Module:\s*constants\.',
                ],
                'severity': 'MEDIUM',
                'category': 'Integration Points'
            },
            'concurrency_patterns': {
                'patterns': [
                    r'workers\s*:=\s*(\d+)',
                    r'make\(.*?,\s*(\d+)\)',
                    r'pool\.Add\((\d+)\)',
                    r'runtime\.GOMAXPROCS\((\d+)\)',
                    r'MaxWorkers:\s*(\d+)',
                    r'sync\.WaitGroup',
                    r'chan\s+\w+.*?(\d+)',
                ],
                'severity': 'HIGH',
                'category': 'Concurrency Configuration'
            },
            'database_hardcoding': {
                'patterns': [
                    r'"fire_salamander\.db"',
                    r'NewStorageManager\("([^"]+)"\)',
                    r'sql\.Open\("([^"]+)"',
                    r'CREATE TABLE IF NOT EXISTS (\w+)',
                    r'SELECT .* FROM (\w+)',
                    r'INSERT INTO (\w+)',
                ],
                'severity': 'HIGH',
                'category': 'Database Configuration'
            }
        }
    
    def detect_violations(self, file_path: str) -> List[OrchestratorViolation]:
        """Detect orchestration hardcoding violations in a Go file"""
        self.violations.clear()
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                lines = f.readlines()
            
            for line_num, line in enumerate(lines, 1):
                line_content = line.strip()
                if not line_content or line_content.startswith('//'):
                    continue
                
                self._check_line_for_violations(line_num, line_content)
            
        except FileNotFoundError:
            print(f"Error: File '{file_path}' not found")
            return []
        except Exception as e:
            print(f"Error reading file '{file_path}': {e}")
            return []
        
        return self.violations
    
    def _check_line_for_violations(self, line_num: int, line_content: str):
        """Check a single line for orchestration violations"""
        for category, config in self.patterns.items():
            for pattern in config['patterns']:
                matches = re.finditer(pattern, line_content)
                for match in matches:
                    hardcoded_value = match.group(1) if match.groups() else match.group(0)
                    
                    # Skip if it's using constants properly
                    if self._is_using_constants_properly(line_content, hardcoded_value):
                        continue
                    
                    # Skip common false positives
                    if self._is_false_positive(hardcoded_value, category):
                        continue
                    
                    violation = self._create_violation(
                        line_num, line_content, category, hardcoded_value,
                        config['severity'], config['category']
                    )
                    
                    self.violations.append(violation)
    
    def _is_using_constants_properly(self, line: str, value: str) -> bool:
        """Check if the line is properly using constants"""
        return (
            'constants.' in line or
            'cfg.' in line or
            'config.' in line or
            'os.Getenv(' in line or
            line.startswith('const ')
        )
    
    def _is_false_positive(self, value: str, category: str) -> bool:
        """Filter out common false positives"""
        false_positives = {
            'timing_values': ['0', '1', '2'],
            'retry_patterns': ['0', '1'],
            'state_strings': ['', 'ok', 'OK'],
            'service_endpoints': [],
            'workflow_steps': ['1', '2', '3', '4', '5'],
            'concurrency_patterns': ['0', '1'],
        }
        
        return value in false_positives.get(category, [])
    
    def _create_violation(self, line_num: int, line_content: str, 
                         violation_type: str, hardcoded_value: str,
                         severity: str, category: str) -> OrchestratorViolation:
        """Create a violation object with detailed information"""
        suggestions = self._get_suggestions(violation_type, hardcoded_value)
        architecture_impact = self._get_architecture_impact(violation_type)
        
        return OrchestratorViolation(
            line_number=line_num,
            line_content=line_content.strip(),
            violation_type=violation_type,
            hardcoded_value=hardcoded_value,
            severity=severity,
            category=category,
            suggestion=suggestions,
            architecture_impact=architecture_impact
        )
    
    def _get_suggestions(self, violation_type: str, value: str) -> str:
        """Get specific suggestions for each violation type"""
        suggestions = {
            'service_endpoints': f'Move "{value}" to configuration file or environment variable',
            'timing_values': f'Move timeout/duration "{value}" to orchestrator configuration',
            'retry_patterns': f'Define retry policy "{value}" in orchestration config',
            'state_strings': f'Use orchestrator state constants for "{value}"',
            'event_names': f'Define event type "{value}" in orchestrator event registry',
            'error_messages': f'Use structured error handling for "{value}"',
            'workflow_steps': f'Define workflow step "{value}" in orchestration plan',
            'integration_points': f'Use dependency injection for service integration',
            'concurrency_patterns': f'Move concurrency setting "{value}" to orchestrator config',
            'database_hardcoding': f'Externalize database configuration "{value}"'
        }
        return suggestions.get(violation_type, f'Consider externalizing "{value}" to configuration')
    
    def _get_architecture_impact(self, violation_type: str) -> str:
        """Get architecture impact description for each violation type"""
        impacts = {
            'service_endpoints': 'Breaks service discovery patterns and deployment flexibility',
            'timing_values': 'Prevents adaptive timeout management and environment-specific tuning',
            'retry_patterns': 'Limits retry policy customization and fault tolerance configuration',
            'state_strings': 'Reduces state machine flexibility and workflow adaptability',
            'event_names': 'Complicates event-driven architecture and monitoring integration',
            'error_messages': 'Hinders structured error handling and error recovery mechanisms',
            'workflow_steps': 'Reduces workflow orchestration flexibility and step reusability',
            'integration_points': 'Creates tight coupling and reduces microservice independence',
            'concurrency_patterns': 'Prevents performance tuning and resource optimization',
            'database_hardcoding': 'Breaks database abstraction and multi-environment support'
        }
        return impacts.get(violation_type, 'May impact system flexibility and maintainability')
    
    def generate_report(self, violations: List[OrchestratorViolation], 
                       file_path: str) -> Dict[str, Any]:
        """Generate a comprehensive orchestration analysis report"""
        if not violations:
            return {
                'file_path': file_path,
                'total_violations': 0,
                'violations_by_category': {},
                'violations_by_severity': {},
                'violations': [],
                'summary': 'No orchestration hardcoding violations detected'
            }
        
        # Group violations
        by_category = {}
        by_severity = {}
        
        for violation in violations:
            # By category
            if violation.category not in by_category:
                by_category[violation.category] = []
            by_category[violation.category].append(asdict(violation))
            
            # By severity
            if violation.severity not in by_severity:
                by_severity[violation.severity] = 0
            by_severity[violation.severity] += 1
        
        return {
            'file_path': file_path,
            'total_violations': len(violations),
            'violations_by_category': by_category,
            'violations_by_severity': by_severity,
            'violations': [asdict(v) for v in violations],
            'summary': self._generate_summary(violations),
            'recommendations': self._generate_recommendations(by_category)
        }
    
    def _generate_summary(self, violations: List[OrchestratorViolation]) -> str:
        """Generate a summary of the violations found"""
        total = len(violations)
        categories = set(v.category for v in violations)
        high_severity = sum(1 for v in violations if v.severity == 'HIGH')
        
        return (
            f"Found {total} orchestration hardcoding violations across "
            f"{len(categories)} categories. {high_severity} high-severity issues "
            f"require immediate attention for proper microservice orchestration."
        )
    
    def _generate_recommendations(self, by_category: Dict[str, List]) -> List[str]:
        """Generate architectural recommendations"""
        recommendations = []
        
        if 'Service Discovery' in by_category:
            recommendations.append(
                "Implement service discovery pattern with registry (Consul, etcd, or Kubernetes DNS)"
            )
        
        if 'Timing Configuration' in by_category:
            recommendations.append(
                "Create centralized timeout management with environment-specific configurations"
            )
        
        if 'Retry Configuration' in by_category:
            recommendations.append(
                "Implement configurable retry policies with exponential backoff"
            )
        
        if 'State Management' in by_category:
            recommendations.append(
                "Use state machine pattern with externalized state definitions"
            )
        
        if 'Integration Points' in by_category:
            recommendations.append(
                "Apply dependency injection pattern for loose coupling between services"
            )
        
        if 'Concurrency Configuration' in by_category:
            recommendations.append(
                "Implement adaptive concurrency management based on system resources"
            )
        
        return recommendations

def main():
    """Main function to run the orchestrator detector"""
    if len(sys.argv) != 2:
        print("Usage: python delta12_orchestrator_detector.py <go_file_path>")
        sys.exit(1)
    
    file_path = sys.argv[1]
    if not Path(file_path).exists():
        print(f"Error: File '{file_path}' does not exist")
        sys.exit(1)
    
    detector = OrchestratorDetector()
    violations = detector.detect_violations(file_path)
    report = detector.generate_report(violations, file_path)
    
    # Print summary
    print(f"\n🔍 Orchestrator Hardcoding Analysis Results")
    print(f"File: {file_path}")
    print(f"Total Violations: {report['total_violations']}")
    print(f"Summary: {report['summary']}\n")
    
    # Print violations by severity
    print("📊 Violations by Severity:")
    for severity, count in sorted(report['violations_by_severity'].items()):
        print(f"  {severity}: {count}")
    
    # Print violations by category
    print(f"\n📋 Violations by Category:")
    for category, violations in report['violations_by_category'].items():
        print(f"  {category}: {len(violations)}")
    
    # Print detailed violations
    print(f"\n🚨 Detailed Violations:")
    for i, violation_dict in enumerate(report['violations'], 1):
        print(f"\n{i}. Line {violation_dict['line_number']} [{violation_dict['severity']}] {violation_dict['category']}")
        print(f"   Code: {violation_dict['line_content']}")
        print(f"   Issue: {violation_dict['hardcoded_value']}")
        print(f"   Fix: {violation_dict['suggestion']}")
        print(f"   Impact: {violation_dict['architecture_impact']}")
    
    # Print recommendations
    if report['recommendations']:
        print(f"\n💡 Architectural Recommendations:")
        for i, rec in enumerate(report['recommendations'], 1):
            print(f"{i}. {rec}")
    
    # Save JSON report
    json_file = file_path.replace('.go', '_orchestrator_analysis.json')
    with open(json_file, 'w') as f:
        json.dump(report, f, indent=2)
    
    print(f"\n📄 Detailed JSON report saved to: {json_file}")

if __name__ == "__main__":
    main()