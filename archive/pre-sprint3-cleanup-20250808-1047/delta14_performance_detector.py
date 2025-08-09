#!/usr/bin/env python3
"""
DELTA-14 Performance Analyzer Detector
Analyzes performance monitoring patterns and architectural designs
Fire Salamander - Professional Code Analysis
"""

import json
import re
import ast
import sys
from typing import Dict, List, Any, Set, Tuple
from pathlib import Path

class PerformanceAnalyzerDetector:
    """Advanced detector for performance analysis patterns and agent architecture"""
    
    def __init__(self):
        self.patterns = {
            'performance_interfaces': [],
            'metric_collection': [],
            'threshold_management': [],
            'optimization_strategies': [],
            'monitoring_patterns': [],
            'caching_mechanisms': [],
            'communication_protocols': [],
            'architectural_patterns': []
        }
        
        self.violations = []
        self.performance_insights = {}
        
    def analyze_go_file(self, file_path: str) -> Dict[str, Any]:
        """Analyze Go file for performance analyzer patterns"""
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                content = f.read()
            
            results = {
                'file_path': file_path,
                'performance_architecture': self.extract_performance_architecture(content),
                'metric_collection_patterns': self.analyze_metric_collection(content),
                'threshold_patterns': self.analyze_threshold_management(content),
                'optimization_analysis': self.analyze_optimization_strategies(content),
                'monitoring_mechanisms': self.analyze_monitoring_patterns(content),
                'caching_strategies': self.analyze_caching_mechanisms(content),
                'communication_patterns': self.analyze_communication_protocols(content),
                'web_vitals_analysis': self.analyze_core_web_vitals(content),
                'http_analysis': self.analyze_http_patterns(content),
                'architectural_violations': self.find_violations(content),
                'reusable_components': self.identify_reusable_components(content),
                'performance_metrics': self.extract_performance_metrics(content)
            }
            
            return results
            
        except Exception as e:
            return {'error': f"Failed to analyze file: {str(e)}"}
    
    def extract_performance_architecture(self, content: str) -> Dict[str, Any]:
        """Extract performance analyzer architecture patterns"""
        architecture = {
            'analyzer_type': 'PerformanceAnalyzer',
            'primary_struct': None,
            'result_struct': None,
            'metric_structs': [],
            'analysis_methods': [],
            'initialization_pattern': None,
            'client_integration': {}
        }
        
        # Find primary analyzer struct
        analyzer_match = re.search(r'type\s+(PerformanceAnalyzer)\s+struct\s*{([^}]+)}', content, re.MULTILINE | re.DOTALL)
        if analyzer_match:
            architecture['primary_struct'] = analyzer_match.group(1)
            struct_content = analyzer_match.group(2)
            
            # Extract struct fields
            fields = re.findall(r'(\w+)\s+\*?([^`\n]+)', struct_content)
            architecture['struct_composition'] = fields
            
            # Check for HTTP client integration
            if 'http.Client' in struct_content:
                architecture['client_integration']['http_client'] = True
            if 'regexp.Regexp' in struct_content:
                architecture['client_integration']['regex_engine'] = True
        
        # Find result struct
        result_match = re.search(r'type\s+(PerformanceMetricsResult)\s+struct\s*{([^}]+)}', content, re.MULTILINE | re.DOTALL)
        if result_match:
            architecture['result_struct'] = result_match.group(1)
            result_content = result_match.group(2)
            
            # Extract JSON-tagged fields
            json_fields = re.findall(r'(\w+)\s+[^`]+`json:"([^"]+)"', result_content)
            architecture['result_fields'] = json_fields
        
        # Find metric structures
        metric_structs = re.findall(r'type\s+(\w+(?:Metric|Analysis|Vitals|Headers|Count)s?)\s+struct', content)
        architecture['metric_structs'] = list(set(metric_structs))
        
        # Extract analysis methods
        analysis_methods = re.findall(r'func\s+\([^)]+\)\s+(analyze\w+|estimate\w+|score\w+)\([^)]*\)', content)
        architecture['analysis_methods'] = analysis_methods
        
        # Find constructor
        constructor_match = re.search(r'func\s+(New\w+Analyzer)\([^)]*\)\s+\*\w+', content)
        if constructor_match:
            architecture['initialization_pattern'] = constructor_match.group(1)
        
        return architecture
    
    def analyze_metric_collection(self, content: str) -> Dict[str, Any]:
        """Analyze metric collection strategies and patterns"""
        metrics = {
            'collection_strategies': [],
            'metric_types': [],
            'measurement_methods': [],
            'aggregation_patterns': [],
            'temporal_patterns': []
        }
        
        # Find metric collection patterns
        metric_patterns = [
            r'LoadTime\s*=\s*time\.Since\(',
            r'PageSize\s*=\s*int64\(',
            r'CompressionRatio\s*=',
            r'ResourceCounts\.\w+\s*=',
            r'CoreWebVitals\.\w+\s*='
        ]
        
        for pattern in metric_patterns:
            matches = re.findall(pattern, content)
            if matches:
                metrics['collection_strategies'].append({
                    'pattern': pattern,
                    'occurrences': len(matches)
                })
        
        # Extract metric types
        metric_types = re.findall(r'(\w+)\s+time\.Duration\s+`json:"([^"]+)"', content)
        metrics['metric_types'].extend([{'name': name, 'json_key': key} for name, key in metric_types])
        
        # Find measurement methods
        measurement_methods = [
            'time.Now()',
            'time.Since(',
            'resp.Header.Get(',
            'strings.Count(',
            'regexp.MustCompile('
        ]
        
        for method in measurement_methods:
            count = len(re.findall(re.escape(method), content))
            if count > 0:
                metrics['measurement_methods'].append({
                    'method': method,
                    'usage_count': count
                })
        
        # Find aggregation patterns
        aggregation_patterns = [
            r'ResourceCounts\s*{[^}]+}',
            r'CoreWebVitals\s*{[^}]+}',
            r'HTTPHeaders\s*{[^}]+}'
        ]
        
        for pattern in aggregation_patterns:
            if re.search(pattern, content, re.MULTILINE | re.DOTALL):
                metrics['aggregation_patterns'].append(pattern)
        
        return metrics
    
    def analyze_threshold_management(self, content: str) -> Dict[str, Any]:
        """Analyze performance threshold patterns and management"""
        thresholds = {
            'threshold_definitions': [],
            'scoring_systems': [],
            'performance_boundaries': [],
            'threshold_constants': []
        }
        
        # Find threshold constants
        threshold_constants = re.findall(r'constants\.(\w*(?:Time|Size|Score|Limit|Threshold)\w*)', content)
        thresholds['threshold_constants'] = list(set(threshold_constants))
        
        # Find scoring systems
        scoring_patterns = [
            r'func\s+\([^)]+\)\s+(score\w+)\([^)]*\)\s+string',
            r'if\s+value\s*<=\s*(\d+)\s*{[^}]*return\s+"(\w+)"',
            r'else\s+if\s+value\s*<=\s*(\d+)\s*{[^}]*return\s+"(\w+)"'
        ]
        
        for pattern in scoring_patterns:
            matches = re.findall(pattern, content)
            if matches:
                thresholds['scoring_systems'].extend(matches)
        
        # Find performance boundaries
        boundary_patterns = [
            r'(\d+)\s*\*\s*1024\s*\*\s*1024',  # MB calculations
            r'value\s*<=\s*(\d+)',  # Threshold comparisons
            r'duration\s*>\s*constants\.(\w+)',  # Duration thresholds
            r'size\s*>\s*(\d+)'  # Size thresholds
        ]
        
        for pattern in boundary_patterns:
            matches = re.findall(pattern, content)
            if matches:
                thresholds['performance_boundaries'].extend(matches)
        
        # Extract Core Web Vitals thresholds
        vitals_thresholds = re.findall(r'"(\w+)\s*≤\s*([\d.]+)s?\s*\([^)]+\)', content)
        thresholds['core_web_vitals_thresholds'] = vitals_thresholds
        
        return thresholds
    
    def analyze_optimization_strategies(self, content: str) -> Dict[str, Any]:
        """Analyze optimization strategies and recommendations"""
        optimization = {
            'optimization_checks': [],
            'recommendation_patterns': [],
            'performance_flags': [],
            'resource_optimization': []
        }
        
        # Find optimization checks
        opt_checks = [
            r'HasCompression',
            r'HasCaching',
            r'OptimizedImages',
            r'MinifiedResources',
            r'HasHTTP2',
            r'HasPreload'
        ]
        
        for check in opt_checks:
            if check in content:
                optimization['optimization_checks'].append(check)
        
        # Extract recommendation patterns
        recommendation_patterns = re.findall(r'result\.Recommendations\s*=\s*append\([^,]+,\s*"([^"]+)"', content)
        optimization['recommendation_patterns'] = recommendation_patterns
        
        # Find performance flags
        flag_patterns = [
            r'\.Has\w+\s*=\s*(true|false)',
            r'\.Enabled\s*=\s*(true|false)',
            r'\.Optimized\w*\s*=\s*\w+'
        ]
        
        for pattern in flag_patterns:
            matches = re.findall(pattern, content)
            if matches:
                optimization['performance_flags'].extend(matches)
        
        # Resource optimization analysis
        resource_patterns = [
            'checkMinification',
            'webp|avif',
            'gzip|deflate|br',
            'keep-alive',
            'max-age'
        ]
        
        for pattern in resource_patterns:
            if pattern in content.lower():
                optimization['resource_optimization'].append(pattern)
        
        return optimization
    
    def analyze_monitoring_patterns(self, content: str) -> Dict[str, Any]:
        """Analyze monitoring and instrumentation patterns"""
        monitoring = {
            'instrumentation_points': [],
            'timing_mechanisms': [],
            'data_collection': [],
            'reporting_patterns': []
        }
        
        # Find instrumentation points
        instrumentation_patterns = [
            r'start\s*:=\s*time\.Now\(\)',
            r'duration\s*:=\s*time\.Since\(',
            r'result\.\w+\s*=\s*[^;]+',
            r'resp\.Header\.Get\('
        ]
        
        for pattern in instrumentation_patterns:
            matches = len(re.findall(pattern, content))
            if matches > 0:
                monitoring['instrumentation_points'].append({
                    'pattern': pattern,
                    'count': matches
                })
        
        # Find timing mechanisms
        timing_patterns = [
            'time.Now()',
            'time.Since(',
            'time.Duration',
            'Timeout:'
        ]
        
        for pattern in timing_patterns:
            count = content.count(pattern)
            if count > 0:
                monitoring['timing_mechanisms'].append({
                    'mechanism': pattern,
                    'usage_count': count
                })
        
        # Data collection patterns
        collection_patterns = [
            r'buf\s*:=\s*make\(\[\]byte',
            r'resp\.Body\.Read\(',
            r'strings\.Count\(',
            r'regexp\.\w+\('
        ]
        
        for pattern in collection_patterns:
            if re.search(pattern, content):
                monitoring['data_collection'].append(pattern)
        
        return monitoring
    
    def analyze_caching_mechanisms(self, content: str) -> Dict[str, Any]:
        """Analyze caching strategies and patterns"""
        caching = {
            'cache_headers': [],
            'cache_strategies': [],
            'cache_validation': [],
            'cache_patterns': []
        }
        
        # Find cache headers
        cache_headers = [
            'Cache-Control',
            'ETag',
            'Last-Modified',
            'Expires',
            'max-age'
        ]
        
        for header in cache_headers:
            if header in content:
                caching['cache_headers'].append(header)
        
        # Find cache strategies
        cache_strategy_patterns = [
            r'HasCacheControl',
            r'HasETag',
            r'HasLastModified',
            r'MaxAge\s*time\.Duration',
            r'cacheRegex'
        ]
        
        for pattern in cache_strategy_patterns:
            if re.search(pattern, content):
                caching['cache_strategies'].append(pattern)
        
        # Cache validation patterns
        validation_patterns = [
            r'resp\.Header\.Get\("Cache-Control"\)',
            r'resp\.Header\.Get\("ETag"\)',
            r'resp\.Header\.Get\("Last-Modified"\)'
        ]
        
        for pattern in validation_patterns:
            if re.search(pattern, content):
                caching['cache_validation'].append(pattern)
        
        return caching
    
    def analyze_communication_protocols(self, content: str) -> Dict[str, Any]:
        """Analyze communication protocol patterns"""
        protocols = {
            'http_patterns': [],
            'request_optimization': [],
            'response_handling': [],
            'protocol_features': []
        }
        
        # HTTP patterns
        http_patterns = [
            r'http\.NewRequestWithContext\(',
            r'http\.Client\{',
            r'req\.Header\.Set\(',
            r'resp\.Header\.Get\(',
            r'resp\.Body\.Close\(\)'
        ]
        
        for pattern in http_patterns:
            count = len(re.findall(pattern, content))
            if count > 0:
                protocols['http_patterns'].append({
                    'pattern': pattern,
                    'count': count
                })
        
        # Request optimization
        opt_patterns = [
            r'Accept-Encoding',
            r'User-Agent',
            r'Accept',
            r'Timeout:'
        ]
        
        for pattern in opt_patterns:
            if pattern in content:
                protocols['request_optimization'].append(pattern)
        
        # Protocol features
        feature_patterns = [
            r'resp\.ProtoMajor\s*==\s*2',  # HTTP/2 detection
            r'Connection.*keep-alive',      # Keep-alive
            r'gzip.*deflate.*br'           # Compression
        ]
        
        for pattern in feature_patterns:
            if re.search(pattern, content):
                protocols['protocol_features'].append(pattern)
        
        return protocols
    
    def analyze_core_web_vitals(self, content: str) -> Dict[str, Any]:
        """Analyze Core Web Vitals implementation patterns"""
        vitals = {
            'vitals_metrics': [],
            'scoring_functions': [],
            'estimation_algorithms': [],
            'threshold_definitions': []
        }
        
        # Find Core Web Vitals metrics
        vital_metrics = [
            'LCP',  # Largest Contentful Paint
            'FID',  # First Input Delay
            'CLS',  # Cumulative Layout Shift
            'TTFB', # Time To First Byte
            'FCP',  # First Contentful Paint
            'SpeedIndex'
        ]
        
        for metric in vital_metrics:
            if f'result.CoreWebVitals.{metric}' in content:
                vitals['vitals_metrics'].append(metric)
        
        # Find scoring functions
        scoring_functions = re.findall(r'func\s+\([^)]+\)\s+(score\w+)\([^)]*\)\s+string', content)
        vitals['scoring_functions'] = scoring_functions
        
        # Estimation algorithms
        estimation_patterns = [
            r'lcpValue\s*\*=\s*[\d.]+',
            r'fidValue\s*:=\s*float64\([^)]+\)',
            r'clsValue\s*\+=\s*[\d.]+',
            r'fcpValue\s*:=\s*[^;]+'
        ]
        
        for pattern in estimation_patterns:
            if re.search(pattern, content):
                vitals['estimation_algorithms'].append(pattern)
        
        # Threshold definitions
        threshold_patterns = re.findall(r'(\w+)\s*≤\s*([\d.]+)s?\s*\(([^)]+)\)', content)
        vitals['threshold_definitions'] = threshold_patterns
        
        return vitals
    
    def analyze_http_patterns(self, content: str) -> Dict[str, Any]:
        """Analyze HTTP-specific patterns and optimizations"""
        http = {
            'header_analysis': [],
            'security_headers': [],
            'performance_headers': [],
            'compression_analysis': []
        }
        
        # Header analysis patterns
        header_patterns = [
            r'resp\.Header\.Get\("([^"]+)"\)',
            r'req\.Header\.Set\("([^"]+)"',
            r'headers\.(\w+)\.(\w+)\s*='
        ]
        
        for pattern in header_patterns:
            matches = re.findall(pattern, content)
            if matches:
                http['header_analysis'].extend(matches)
        
        # Security headers
        security_headers = [
            'Strict-Transport-Security',
            'Content-Security-Policy',
            'X-Frame-Options',
            'X-Content-Type-Options'
        ]
        
        for header in security_headers:
            if header in content:
                http['security_headers'].append(header)
        
        # Performance headers
        perf_headers = [
            'keep-alive',
            'preload',
            'prefetch',
            'Content-Encoding'
        ]
        
        for header in perf_headers:
            if header in content:
                http['performance_headers'].append(header)
        
        # Compression analysis
        compression_patterns = [
            r'compressionRegex',
            r'gzip|deflate|br',
            r'Content-Encoding',
            r'CompressionRatio'
        ]
        
        for pattern in compression_patterns:
            if pattern in content:
                http['compression_analysis'].append(pattern)
        
        return http
    
    def find_violations(self, content: str) -> List[Dict[str, Any]]:
        """Find architectural violations and performance anti-patterns"""
        violations = []
        
        # Check for blocking operations without context
        blocking_patterns = [
            r'\.Do\(req\)',
            r'\.Read\(',
            r'time\.Sleep\('
        ]
        
        for pattern in blocking_patterns:
            matches = re.findall(pattern, content)
            for match in matches:
                violations.append({
                    'type': 'BLOCKING_OPERATION',
                    'severity': 'MEDIUM',
                    'description': 'Potentially blocking operation without proper timeout handling',
                    'pattern': match,
                    'recommendation': 'Ensure proper context timeout and error handling'
                })
        
        # Check for hardcoded performance values
        hardcoded_patterns = [
            r'\d+\s*\*\s*1024\s*\*\s*1024',  # Hardcoded MB values
            r'4\s*\*\s*1024\s*\*\s*1024',   # 4MB buffer
            r'make\(\[\]byte,\s*\d+'          # Fixed buffer sizes
        ]
        
        for pattern in hardcoded_patterns:
            matches = re.findall(pattern, content)
            for match in matches:
                violations.append({
                    'type': 'HARDCODED_VALUES',
                    'severity': 'LOW',
                    'description': 'Hardcoded performance values should be configurable',
                    'value': match,
                    'recommendation': 'Use configuration constants'
                })
        
        # Check for missing error handling in performance critical paths
        critical_calls = re.findall(r'(resp\.Body\.Read\([^)]+\))[^;]*;[^i]*$', content, re.MULTILINE)
        for call in critical_calls:
            violations.append({
                'type': 'ERROR_HANDLING',
                'severity': 'MEDIUM',
                'description': 'Missing comprehensive error handling in performance-critical operation',
                'pattern': call,
                'recommendation': 'Add proper error handling and resource cleanup'
            })
        
        return violations
    
    def identify_reusable_components(self, content: str) -> Dict[str, Any]:
        """Identify reusable performance analysis components"""
        reusable = {
            'utility_functions': [],
            'metric_calculators': [],
            'analyzer_interfaces': [],
            'common_patterns': []
        }
        
        # Find utility functions
        utility_patterns = [
            r'func\s+(score\w+)\(',
            r'func\s+(analyze\w+)\(',
            r'func\s+(estimate\w+)\(',
            r'func\s+(check\w+)\('
        ]
        
        for pattern in utility_patterns:
            matches = re.findall(pattern, content)
            reusable['utility_functions'].extend(matches)
        
        # Metric calculators
        calculator_patterns = [
            r'func\s+\([^)]+\)\s+(calculate\w+)\(',
            r'func\s+\([^)]+\)\s+(measure\w+)\(',
            r'func\s+\([^)]+\)\s+(compute\w+)\('
        ]
        
        for pattern in calculator_patterns:
            matches = re.findall(pattern, content)
            reusable['metric_calculators'].extend(matches)
        
        # Interface candidates
        interface_methods = {
            'PerformanceAnalyzer': ['Analyze', 'analyzePageLoad', 'analyzeHTTPHeaders'],
            'MetricCollector': ['estimateCoreWebVitals', 'extract', 'measure'],
            'ThresholdManager': ['score', 'evaluate', 'check'],
            'ReportGenerator': ['generate', 'format', 'export']
        }
        
        for interface, methods in interface_methods.items():
            found_methods = []
            for method in methods:
                if any(method.lower() in func.lower() for func in reusable['utility_functions']):
                    found_methods.append(method)
            
            if len(found_methods) >= 2:
                reusable['analyzer_interfaces'].append({
                    'interface': interface,
                    'methods': found_methods
                })
        
        return reusable
    
    def extract_performance_metrics(self, content: str) -> Dict[str, Any]:
        """Extract performance metrics and measurement patterns"""
        metrics = {
            'timing_metrics': [],
            'size_metrics': [],
            'quality_metrics': [],
            'optimization_metrics': []
        }
        
        # Timing metrics
        timing_patterns = [
            r'LoadTime',
            r'TTFB',
            r'duration',
            r'time\.Duration'
        ]
        
        for pattern in timing_patterns:
            count = len(re.findall(pattern, content))
            if count > 0:
                metrics['timing_metrics'].append({
                    'metric': pattern,
                    'occurrences': count
                })
        
        # Size metrics
        size_patterns = [
            r'PageSize',
            r'CompressedSize',
            r'CompressionRatio',
            r'int64'
        ]
        
        for pattern in size_patterns:
            count = len(re.findall(pattern, content))
            if count > 0:
                metrics['size_metrics'].append({
                    'metric': pattern,
                    'occurrences': count
                })
        
        # Quality metrics
        quality_patterns = [
            r'OptimizedImages',
            r'MinifiedResources',
            r'HasCompression',
            r'HasCaching'
        ]
        
        for pattern in quality_patterns:
            if pattern in content:
                metrics['quality_metrics'].append(pattern)
        
        return metrics
    
    def generate_analysis_report(self, analysis_data: Dict[str, Any]) -> Dict[str, Any]:
        """Generate comprehensive performance analysis report"""
        report = {
            'analysis_metadata': {
                'analyzer': 'DELTA-14 Performance Analyzer Detector',
                'target_file': analysis_data.get('file_path', 'unknown'),
                'analysis_timestamp': '2025-08-07T00:00:00Z',
                'analyzer_type': 'PerformanceAnalyzer'
            },
            'performance_architecture': analysis_data.get('performance_architecture', {}),
            'metric_collection_analysis': analysis_data.get('metric_collection_patterns', {}),
            'threshold_analysis': analysis_data.get('threshold_patterns', {}),
            'optimization_analysis': analysis_data.get('optimization_analysis', {}),
            'monitoring_analysis': analysis_data.get('monitoring_mechanisms', {}),
            'caching_analysis': analysis_data.get('caching_strategies', {}),
            'communication_analysis': analysis_data.get('communication_patterns', {}),
            'core_web_vitals_analysis': analysis_data.get('web_vitals_analysis', {}),
            'http_analysis': analysis_data.get('http_analysis', {}),
            'violations_found': analysis_data.get('architectural_violations', []),
            'reusability_analysis': analysis_data.get('reusable_components', {}),
            'performance_metrics_analysis': analysis_data.get('performance_metrics', {}),
            'architectural_insights': {
                'performance_agent_pattern': 'HTTP-based performance analyzer with metric aggregation',
                'metric_collection_strategy': 'Multi-faceted metrics with Core Web Vitals estimation',
                'threshold_management': 'Score-based evaluation with configurable thresholds',
                'optimization_detection': 'Boolean flags with recommendation generation',
                'communication_protocol': 'HTTP client with header analysis and timeout handling'
            },
            'recommendations': {
                'architectural_improvements': [
                    'Extract common performance interface for multiple analyzers',
                    'Implement configurable metric collection strategies',
                    'Create reusable threshold management system',
                    'Standardize optimization recommendation patterns'
                ],
                'performance_improvements': [
                    'Implement connection pooling for HTTP clients',
                    'Add request/response caching mechanisms',
                    'Optimize regex compilation and reuse',
                    'Implement parallel resource analysis'
                ],
                'monitoring_improvements': [
                    'Add distributed tracing integration',
                    'Implement real-time metric streaming',
                    'Create performance baseline tracking',
                    'Add anomaly detection for performance regressions'
                ]
            }
        }
        
        return report

def main():
    """Main execution function"""
    if len(sys.argv) != 2:
        print("Usage: python delta14_performance_detector.py <go_file_path>")
        sys.exit(1)
    
    file_path = sys.argv[1]
    
    if not Path(file_path).exists():
        print(f"Error: File {file_path} does not exist")
        sys.exit(1)
    
    detector = PerformanceAnalyzerDetector()
    analysis_data = detector.analyze_go_file(file_path)
    
    if 'error' in analysis_data:
        print(f"Analysis failed: {analysis_data['error']}")
        sys.exit(1)
    
    report = detector.generate_analysis_report(analysis_data)
    
    # Output JSON report
    print(json.dumps(report, indent=2, ensure_ascii=False))
    
    # Print summary statistics
    print(f"\n=== DELTA-14 Analysis Summary ===", file=sys.stderr)
    print(f"Analyzer Type: {report['analysis_metadata']['analyzer_type']}", file=sys.stderr)
    print(f"Violations Found: {len(report['violations_found'])}", file=sys.stderr)
    print(f"Metric Collection Strategies: {len(report['metric_collection_analysis'].get('collection_strategies', []))}", file=sys.stderr)
    print(f"Core Web Vitals Metrics: {len(report['core_web_vitals_analysis'].get('vitals_metrics', []))}", file=sys.stderr)
    print(f"Reusable Components: {len(report['reusability_analysis'].get('utility_functions', []))}", file=sys.stderr)

if __name__ == "__main__":
    main()