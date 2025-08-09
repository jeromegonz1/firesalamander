#!/usr/bin/env python3
"""
DELTA-13 Data Integrity Agent Detector
Analyzes data integrity patterns and agent communication architectures
Fire Salamander - Professional Code Analysis
"""

import json
import re
import ast
import sys
from typing import Dict, List, Any, Set, Tuple
from pathlib import Path

class DataIntegrityAgentDetector:
    """Advanced detector for data integrity agent patterns and architectural concerns"""
    
    def __init__(self):
        self.patterns = {
            'agent_interfaces': [],
            'data_validation_mechanisms': [],
            'performance_thresholds': [],
            'metric_collection_strategies': [],
            'error_handling_patterns': [],
            'configuration_management': [],
            'communication_patterns': [],
            'architectural_concerns': []
        }
        
        self.violations = []
        self.architectural_insights = {}
        
    def analyze_go_file(self, file_path: str) -> Dict[str, Any]:
        """Analyze Go file for data integrity agent patterns"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            results = {
                'file_path': file_path,
                'agent_architecture': self.extract_agent_architecture(content),
                'data_validation': self.analyze_data_validation(content),
                'performance_patterns': self.analyze_performance_patterns(content),
                'error_handling': self.analyze_error_handling(content),
                'configuration_management': self.analyze_configuration(content),
                'communication_interfaces': self.analyze_communication(content),
                'architectural_violations': self.find_violations(content),
                'reusable_components': self.identify_reusable_components(content),
                'metrics': self.extract_metrics(content)
            }
            
            return results
            
        except Exception as e:
            return {'error': f"Failed to analyze file: {str(e)}"}
    
    def extract_agent_architecture(self, content: str) -> Dict[str, Any]:
        """Extract core agent architecture patterns"""
        architecture = {
            'agent_type': 'DataIntegrityAgent',
            'primary_struct': None,
            'config_struct': None,
            'stats_struct': None,
            'main_interface_methods': [],
            'initialization_pattern': None,
            'lifecycle_methods': []
        }
        
        # Find primary agent struct
        agent_struct_match = re.search(r'type\s+(\w+Agent)\s+struct\s*{([^}]+)}', content, re.MULTILINE | re.DOTALL)
        if agent_struct_match:
            architecture['primary_struct'] = agent_struct_match.group(1)
            struct_content = agent_struct_match.group(2)
            
            # Extract struct fields for architectural analysis
            fields = re.findall(r'(\w+)\s+\*?(\w+)', struct_content)
            architecture['struct_composition'] = fields
        
        # Find configuration struct
        config_match = re.search(r'type\s+(\w+Config)\s+struct\s*{([^}]+)}', content, re.MULTILINE | re.DOTALL)
        if config_match:
            architecture['config_struct'] = config_match.group(1)
            config_content = config_match.group(2)
            architecture['config_fields'] = re.findall(r'(\w+)\s+(\w+)', config_content)
        
        # Find stats/metrics struct
        stats_match = re.search(r'type\s+(\w+Stats)\s+struct\s*{([^}]+)}', content, re.MULTILINE | re.DOTALL)
        if stats_match:
            architecture['stats_struct'] = stats_match.group(1)
            
        # Extract main interface methods
        method_patterns = [
            r'func\s+\([^)]+\)\s+(Run\w+)\([^)]*\)\s+error',
            r'func\s+\([^)]+\)\s+(Test\w+)\([^)]*\)',
            r'func\s+\([^)]+\)\s+(Validate\w+)\([^)]*\)',
            r'func\s+\([^)]+\)\s+(Analyze\w+)\([^)]*\)'
        ]
        
        for pattern in method_patterns:
            matches = re.findall(pattern, content)
            architecture['main_interface_methods'].extend(matches)
        
        # Find constructor pattern
        constructor_match = re.search(r'func\s+(New\w+Agent)\([^)]*\)\s+\*\w+Agent', content)
        if constructor_match:
            architecture['initialization_pattern'] = constructor_match.group(1)
        
        return architecture
    
    def analyze_data_validation(self, content: str) -> Dict[str, Any]:
        """Analyze data validation mechanisms and patterns"""
        validation = {
            'validation_categories': [],
            'validation_methods': [],
            'constraint_checks': [],
            'integrity_tests': [],
            'data_quality_checks': [],
            'validation_architecture': {}
        }
        
        # Extract test categories
        test_categories = re.findall(r'constants\.TestCategory(\w+)', content)
        validation['validation_categories'] = list(set(test_categories))
        
        # Find validation methods
        validation_methods = re.findall(r'func\s+\([^)]+\)\s+(test\w+|validate\w+|check\w+)\([^)]*\)', content)
        validation['validation_methods'] = validation_methods
        
        # Extract constraint checks
        constraint_patterns = [
            r'CHECK\s*\([^)]+\)',
            r'NOT\s+NULL',
            r'UNIQUE\s*\([^)]+\)',
            r'FOREIGN\s+KEY'
        ]
        
        for pattern in constraint_patterns:
            matches = re.findall(pattern, content, re.IGNORECASE)
            if matches:
                validation['constraint_checks'].extend(matches)
        
        # Analyze integrity test patterns
        integrity_tests = re.findall(r'(TestResult\s*{[^}]+})', content, re.MULTILINE | re.DOTALL)
        validation['integrity_tests'] = len(integrity_tests)
        
        # Find data quality patterns
        quality_patterns = [
            r'testURLQuality',
            r'testStatusCodeQuality', 
            r'testSEOScoreQuality',
            r'testNullValues',
            r'testUniqueConstraints'
        ]
        
        for pattern in quality_patterns:
            if pattern in content:
                validation['data_quality_checks'].append(pattern)
        
        return validation
    
    def analyze_performance_patterns(self, content: str) -> Dict[str, Any]:
        """Analyze performance monitoring and threshold patterns"""
        performance = {
            'performance_tests': [],
            'threshold_definitions': [],
            'timing_mechanisms': [],
            'performance_metrics': [],
            'database_optimization': []
        }
        
        # Find performance test methods
        perf_tests = re.findall(r'func\s+\([^)]+\)\s+(test\w*[Pp]erformance\w*)\([^)]*\)', content)
        performance['performance_tests'] = perf_tests
        
        # Extract threshold patterns
        threshold_patterns = [
            r'SlowResponseTime',
            r'AcceptableLoadTime',
            r'DefaultTimeout',
            r'ClientTimeout'
        ]
        
        for pattern in threshold_patterns:
            if pattern in content:
                performance['threshold_definitions'].append(pattern)
        
        # Find timing mechanisms
        timing_patterns = [
            r'time\.Now\(\)',
            r'time\.Since\(',
            r'duration\s*:=',
            r'time\.Duration'
        ]
        
        for pattern in timing_patterns:
            matches = len(re.findall(pattern, content))
            if matches > 0:
                performance['timing_mechanisms'].append({'pattern': pattern, 'count': matches})
        
        # Extract performance metrics
        metric_patterns = [
            r'LoadTime',
            r'PageSize',
            r'CompressionRatio',
            r'ResponseTime'
        ]
        
        for pattern in metric_patterns:
            if pattern in content:
                performance['performance_metrics'].append(pattern)
        
        return performance
    
    def analyze_error_handling(self, content: str) -> Dict[str, Any]:
        """Analyze error handling patterns and strategies"""
        error_handling = {
            'error_types': [],
            'error_handling_methods': [],
            'error_propagation': [],
            'error_logging': [],
            'recovery_mechanisms': []
        }
        
        # Find error handling patterns
        error_patterns = [
            r'if\s+err\s*!=\s*nil\s*{([^}]+)}',
            r'return\s+[^,]*,\s*err',
            r'return\s+[^,]*,\s*fmt\.Errorf\(',
            r'log\.Printf\([^)]*err[^)]*\)'
        ]
        
        for pattern in error_patterns:
            matches = re.findall(pattern, content, re.MULTILINE | re.DOTALL)
            if matches:
                error_handling['error_handling_methods'].append({
                    'pattern': pattern,
                    'count': len(matches)
                })
        
        # Find error types and constants
        error_constants = re.findall(r'Status(Error|Failed|Warning)', content)
        error_handling['error_types'] = list(set(error_constants))
        
        # Analyze error logging
        log_patterns = re.findall(r'log\.(Print|Fatal|Panic)\w*\([^)]+\)', content)
        error_handling['error_logging'] = log_patterns
        
        return error_handling
    
    def analyze_configuration(self, content: str) -> Dict[str, Any]:
        """Analyze configuration management patterns"""
        config = {
            'config_structures': [],
            'default_configs': [],
            'config_validation': [],
            'environment_handling': [],
            'config_sources': []
        }
        
        # Find configuration structures
        config_structs = re.findall(r'type\s+(\w*Config)\s+struct', content)
        config['config_structures'] = config_structs
        
        # Find default configuration methods
        default_methods = re.findall(r'func\s+(default\w+)\([^)]*\)', content)
        config['default_configs'] = default_methods
        
        # Extract configuration fields with JSON tags
        json_tag_pattern = re.findall(r'(\w+)\s+\w+\s+`json:"([^"]+)"', content)
        config['config_fields_with_json'] = json_tag_pattern
        
        # Find constants used in configuration
        const_patterns = re.findall(r'constants\.(\w+)', content)
        config['constants_used'] = list(set(const_patterns))
        
        return config
    
    def analyze_communication(self, content: str) -> Dict[str, Any]:
        """Analyze agent communication and interface patterns"""
        communication = {
            'database_interfaces': [],
            'external_apis': [],
            'inter_agent_communication': [],
            'data_exchange_formats': [],
            'protocol_patterns': []
        }
        
        # Find database communication patterns
        db_patterns = [
            r'sql\.Open\(',
            r'\.Query\(',
            r'\.QueryRow\(',
            r'\.Exec\(',
            r'\.Prepare\('
        ]
        
        for pattern in db_patterns:
            matches = len(re.findall(pattern, content))
            if matches > 0:
                communication['database_interfaces'].append({
                    'pattern': pattern,
                    'usage_count': matches
                })
        
        # Find data exchange formats
        format_patterns = [
            r'json\.Marshal',
            r'json\.Unmarshal',
            r'json\.MarshalIndent',
            r'ioutil\.WriteFile'
        ]
        
        for pattern in format_patterns:
            if pattern in content:
                communication['data_exchange_formats'].append(pattern)
        
        return communication
    
    def find_violations(self, content: str) -> List[Dict[str, Any]]:
        """Find architectural violations and anti-patterns"""
        violations = []
        
        # Check for direct SQL string construction (SQL injection risk)
        sql_concat_patterns = [
            r'fmt\.Sprintf\([^)]*SELECT[^)]*\)',
            r'fmt\.Sprintf\([^)]*INSERT[^)]*\)',
            r'fmt\.Sprintf\([^)]*UPDATE[^)]*\)',
            r'fmt\.Sprintf\([^)]*DELETE[^)]*\)'
        ]
        
        for pattern in sql_concat_patterns:
            matches = re.findall(pattern, content, re.IGNORECASE)
            for match in matches:
                violations.append({
                    'type': 'SQL_INJECTION_RISK',
                    'severity': 'HIGH',
                    'description': 'Direct SQL string construction detected',
                    'pattern': match,
                    'recommendation': 'Use prepared statements or parameterized queries'
                })
        
        # Check for missing error handling
        error_ignore_patterns = re.findall(r'_,\s*err\s*:=', content)
        for match in error_ignore_patterns:
            violations.append({
                'type': 'ERROR_HANDLING',
                'severity': 'MEDIUM',
                'description': 'Error potentially ignored',
                'pattern': match,
                'recommendation': 'Handle or log errors appropriately'
            })
        
        # Check for hardcoded values
        hardcoded_patterns = [
            r'http://\w+',
            r'https://\w+',
            r'/[a-zA-Z0-9_/]+\.db',
            r':\d{4,5}/'
        ]
        
        for pattern in hardcoded_patterns:
            matches = re.findall(pattern, content)
            for match in matches:
                violations.append({
                    'type': 'HARDCODED_VALUES',
                    'severity': 'MEDIUM',
                    'description': 'Hardcoded value detected',
                    'value': match,
                    'recommendation': 'Use configuration or constants'
                })
        
        return violations
    
    def identify_reusable_components(self, content: str) -> Dict[str, Any]:
        """Identify components that could be reused across agents"""
        reusable = {
            'utility_functions': [],
            'common_patterns': [],
            'shared_structures': [],
            'interface_candidates': []
        }
        
        # Find utility functions that could be shared
        utility_patterns = [
            r'func\s+(connect\w+)\(',
            r'func\s+(validate\w+)\(',
            r'func\s+(generate\w+Report)\(',
            r'func\s+(calculate\w+Score)\('
        ]
        
        for pattern in utility_patterns:
            matches = re.findall(pattern, content)
            reusable['utility_functions'].extend(matches)
        
        # Find common structures
        common_structs = re.findall(r'type\s+(\w+Result)\s+struct', content)
        reusable['shared_structures'] = common_structs
        
        # Identify interface candidates
        method_groups = {
            'Reporter': ['generateReport', 'generateHTMLReport'],
            'Validator': ['validate', 'test', 'check'],
            'Configurable': ['defaultConfig', 'NewAgent'],
            'Analyzer': ['analyze', 'Run']
        }
        
        for interface_name, methods in method_groups.items():
            found_methods = []
            for method in methods:
                if any(method.lower() in func.lower() for func in reusable['utility_functions']):
                    found_methods.append(method)
            
            if len(found_methods) >= 2:
                reusable['interface_candidates'].append({
                    'interface': interface_name,
                    'methods': found_methods
                })
        
        return reusable
    
    def extract_metrics(self, content: str) -> Dict[str, Any]:
        """Extract metrics and measurement patterns"""
        metrics = {
            'measurement_types': [],
            'scoring_systems': [],
            'threshold_definitions': [],
            'metric_aggregation': []
        }
        
        # Find measurement types
        measurement_patterns = [
            r'COUNT\(\*\)',
            r'time\.Duration',
            r'file\.Size\(\)',
            r'len\([^)]+\)'
        ]
        
        for pattern in measurement_patterns:
            matches = len(re.findall(pattern, content))
            if matches > 0:
                metrics['measurement_types'].append({
                    'type': pattern,
                    'count': matches
                })
        
        # Find scoring systems
        scoring_patterns = re.findall(r'totalScore\s*[-+*/]=\s*\d+', content)
        if scoring_patterns:
            metrics['scoring_systems'] = scoring_patterns
        
        # Extract threshold definitions
        threshold_vars = re.findall(r'(\w+)\s*[><=]\s*(\d+)', content)
        metrics['threshold_definitions'] = threshold_vars
        
        return metrics
    
    def generate_analysis_report(self, analysis_data: Dict[str, Any]) -> Dict[str, Any]:
        """Generate comprehensive analysis report"""
        report = {
            'analysis_metadata': {
                'analyzer': 'DELTA-13 Data Integrity Agent Detector',
                'target_file': analysis_data.get('file_path', 'unknown'),
                'analysis_timestamp': '2025-08-07T00:00:00Z',
                'agent_type': 'DataIntegrityAgent'
            },
            'architectural_analysis': analysis_data.get('agent_architecture', {}),
            'data_integrity_patterns': analysis_data.get('data_validation', {}),
            'performance_analysis': analysis_data.get('performance_patterns', {}),
            'error_handling_analysis': analysis_data.get('error_handling', {}),
            'configuration_analysis': analysis_data.get('configuration_management', {}),
            'communication_analysis': analysis_data.get('communication_interfaces', {}),
            'violations_found': analysis_data.get('architectural_violations', []),
            'reusability_analysis': analysis_data.get('reusable_components', {}),
            'metrics_analysis': analysis_data.get('metrics', {}),
            'architectural_insights': {
                'agent_interface_pattern': 'Struct-based agent with configuration and stats composition',
                'data_validation_architecture': 'Category-based validation with test results aggregation',
                'error_handling_strategy': 'Error propagation with logging and status tracking',
                'performance_monitoring': 'Time-based measurements with threshold comparisons',
                'configuration_pattern': 'JSON-tagged structs with default configuration factory'
            },
            'recommendations': {
                'architectural_improvements': [
                    'Extract common agent interface for reusability',
                    'Implement dependency injection for database connections',
                    'Create shared validation framework',
                    'Standardize error handling patterns'
                ],
                'security_improvements': [
                    'Use prepared statements for all SQL queries',
                    'Implement input validation for all external data',
                    'Add authentication for database connections'
                ],
                'performance_improvements': [
                    'Implement connection pooling',
                    'Add query result caching',
                    'Optimize database queries with indexes'
                ]
            }
        }
        
        return report

def main():
    """Main execution function"""
    if len(sys.argv) != 2:
        print("Usage: python delta13_data_integrity_detector.py <go_file_path>")
        sys.exit(1)
    
    file_path = sys.argv[1]
    
    if not Path(file_path).exists():
        print(f"Error: File {file_path} does not exist")
        sys.exit(1)
    
    detector = DataIntegrityAgentDetector()
    analysis_data = detector.analyze_go_file(file_path)
    
    if 'error' in analysis_data:
        print(f"Analysis failed: {analysis_data['error']}")
        sys.exit(1)
    
    report = detector.generate_analysis_report(analysis_data)
    
    # Output JSON report
    print(json.dumps(report, indent=2, ensure_ascii=False))
    
    # Print summary statistics
    print(f"\n=== DELTA-13 Analysis Summary ===", file=sys.stderr)
    print(f"Agent Type: {report['analysis_metadata']['agent_type']}", file=sys.stderr)
    print(f"Violations Found: {len(report['violations_found'])}", file=sys.stderr)
    print(f"Validation Categories: {len(report['data_integrity_patterns'].get('validation_categories', []))}", file=sys.stderr)
    print(f"Performance Tests: {len(report['performance_analysis'].get('performance_tests', []))}", file=sys.stderr)
    print(f"Reusable Components: {len(report['reusability_analysis'].get('utility_functions', []))}", file=sys.stderr)

if __name__ == "__main__":
    main()