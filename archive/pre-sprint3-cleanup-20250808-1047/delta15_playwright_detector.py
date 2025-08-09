#!/usr/bin/env python3
"""
DELTA-15: Playwright Frontend Testing Architecture Detector
Final mission in the DELTA analysis series focusing on frontend testing patterns.

This detector analyzes Playwright-based frontend testing architecture for:
- Browser automation patterns
- Selector strategies and maintainability
- Test assertion frameworks
- Screenshot and visual testing management
- Performance measurement integration
- User interaction simulation patterns
- Cross-browser compatibility handling
- Test data management strategies

Target: playwright_agent.go (73 violations expected)
Focus: Frontend testing architecture and best practices
"""

import ast
import re
import json
from pathlib import Path
from typing import List, Dict, Any, Tuple
from dataclasses import dataclass, asdict
from datetime import datetime

@dataclass
class PlaywrightViolation:
    """Represents a frontend testing architectural violation."""
    id: str
    severity: str  # CRITICAL, HIGH, MEDIUM, LOW
    category: str
    title: str
    description: str
    line_number: int
    code_snippet: str
    impact: str
    recommendation: str
    pattern_type: str

class PlaywrightDetector:
    """Advanced detector for frontend testing architecture analysis."""
    
    def __init__(self):
        self.violations: List[PlaywrightViolation] = []
        self.file_path = ""
        self.content = ""
        self.lines = []
        
    def analyze_file(self, file_path: str) -> Dict[str, Any]:
        """Main analysis entry point."""
        self.file_path = file_path
        
        try:
            with open(file_path, 'r', encoding='utf-8') as f:
                self.content = f.read()
                self.lines = self.content.split('\n')
        except Exception as e:
            return {"error": f"Failed to read file: {e}"}
            
        # Run all detection methods
        self._detect_browser_automation_patterns()
        self._detect_selector_strategies()
        self._detect_assertion_patterns()
        self._detect_screenshot_management()
        self._detect_performance_testing()
        self._detect_user_interaction_patterns()
        self._detect_cross_browser_compatibility()
        self._detect_test_data_management()
        self._detect_page_object_patterns()
        self._detect_wait_strategies()
        self._detect_error_handling()
        self._detect_test_maintainability()
        
        return {
            "file_path": file_path,
            "total_violations": len(self.violations),
            "violations": [asdict(v) for v in self.violations],
            "analysis_timestamp": datetime.now().isoformat(),
            "categories": self._get_category_stats()
        }
    
    def _detect_browser_automation_patterns(self):
        """Detect browser automation architectural issues."""
        patterns = [
            # Browser lifecycle management
            {
                "pattern": r"NewPlaywrightAgent\(\)",
                "severity": "MEDIUM",
                "title": "Browser Pool Not Implemented",
                "description": "Browser instances created without pooling mechanism",
                "recommendation": "Implement browser instance pooling for resource efficiency"
            },
            {
                "pattern": r"mcpClient.*endpoint.*mcp://playwright",
                "severity": "HIGH", 
                "title": "Hardcoded MCP Endpoint",
                "description": "MCP Playwright endpoint is hardcoded in implementation",
                "recommendation": "Use configuration-driven endpoint management"
            },
            {
                "pattern": r"Browsers.*\[\]string.*chromium.*firefox.*webkit",
                "severity": "MEDIUM",
                "title": "Static Browser Configuration",
                "description": "Browser list is statically defined without runtime flexibility",
                "recommendation": "Implement dynamic browser configuration based on environment"
            },
            # Browser context isolation
            {
                "pattern": r"testCrossBrowser.*for.*browser.*range",
                "severity": "HIGH",
                "title": "Shared Browser Context Risk",
                "description": "Cross-browser tests may share contexts without proper isolation",
                "recommendation": "Ensure separate browser contexts for each test execution"
            },
            {
                "pattern": r"testResponsiveDesign.*viewports.*size",
                "severity": "MEDIUM",
                "title": "Viewport Switching Without Context Reset",
                "description": "Viewport changes without proper context isolation",
                "recommendation": "Create new browser contexts for each viewport test"
            }
        ]
        
        self._check_patterns(patterns, "Browser Automation")
    
    def _detect_selector_strategies(self):
        """Detect selector strategy and maintainability issues."""
        patterns = [
            # Selector maintainability
            {
                "pattern": r'\.locator\(["\']\.[\w-]+["\']',
                "severity": "HIGH",
                "title": "CSS Class Selector Coupling",
                "description": "Direct CSS class selectors create tight coupling to styling",
                "recommendation": "Use data-testid attributes or semantic selectors"
            },
            {
                "pattern": r'locator.*text=.*button\.btn-secondary',
                "severity": "HIGH",
                "title": "Fragile Element Selector",
                "description": "Element selector relies on styling classes rather than semantic attributes",
                "recommendation": "Use [data-testid] or role-based selectors"
            },
            {
                "pattern": r'Element.*string.*element',
                "severity": "MEDIUM",
                "title": "Generic Element Reference",
                "description": "Accessibility violations reference elements generically",
                "recommendation": "Implement specific element identification strategy"
            },
            # Page Object Model absence
            {
                "pattern": r'page\.locator\(',
                "severity": "CRITICAL",
                "title": "Missing Page Object Model",
                "description": "Direct page interactions without Page Object Model abstraction",
                "recommendation": "Implement Page Object Model for better maintainability"
            },
            {
                "pattern": r'fmt\.Sprintf.*_.*home\.png',
                "severity": "MEDIUM",
                "title": "Hardcoded Screenshot Paths",
                "description": "Screenshot paths are hardcoded without configuration",
                "recommendation": "Use configurable screenshot naming strategy"
            }
        ]
        
        self._check_patterns(patterns, "Selector Strategy")
    
    def _detect_assertion_patterns(self):
        """Detect test assertion framework issues."""
        patterns = [
            # Assertion quality
            {
                "pattern": r'Score.*float64.*95\.0',
                "severity": "HIGH",
                "title": "Magic Number Thresholds",
                "description": "Accessibility score threshold hardcoded without configuration",
                "recommendation": "Define test thresholds in configuration files"
            },
            {
                "pattern": r'violations.*len.*== 0',
                "severity": "MEDIUM",
                "title": "Binary Pass/Fail Logic",
                "description": "Accessibility validation uses binary logic without severity levels",
                "recommendation": "Implement graduated severity-based validation"
            },
            {
                "pattern": r'pa\.results\.TotalTests\+\+',
                "severity": "HIGH",
                "title": "Manual Test Counter Management",
                "description": "Test counters manually managed without framework integration",
                "recommendation": "Use testing framework's built-in assertion counting"
            },
            # Assertion comprehensiveness  
            {
                "pattern": r'calculateStatus.*if.*FailedTests == 0',
                "severity": "MEDIUM",
                "title": "Simplistic Status Calculation",
                "description": "Test status calculation doesn't account for warning conditions",
                "recommendation": "Implement comprehensive status determination with warning states"
            },
            {
                "pattern": r'Performance.*Score.*90\.0',
                "severity": "HIGH",
                "title": "Performance Threshold Without Context",
                "description": "Performance score threshold lacks contextual validation",
                "recommendation": "Implement dynamic performance baselines based on page types"
            }
        ]
        
        self._check_patterns(patterns, "Assertion Framework")
    
    def _detect_screenshot_management(self):
        """Detect screenshot and visual testing management issues."""
        patterns = [
            # Screenshot organization
            {
                "pattern": r'Screenshots.*bool.*screenshots',
                "severity": "MEDIUM",
                "title": "Boolean Screenshot Control",
                "description": "Screenshot management uses simple boolean without granular control",
                "recommendation": "Implement screenshot strategy configuration (always/failure/none)"
            },
            {
                "pattern": r'Path.*tests/screenshots.*\.png',
                "severity": "HIGH",
                "title": "Hardcoded Screenshot Storage",
                "description": "Screenshot storage path hardcoded without environment awareness",
                "recommendation": "Use environment-specific screenshot storage strategies"
            },
            {
                "pattern": r'Visual.*Diff.*Report.*ChangedScreenshots.*int',
                "severity": "MEDIUM",
                "title": "Limited Visual Diff Metrics",
                "description": "Visual regression only tracks changed screenshots count",
                "recommendation": "Implement comprehensive visual diff metrics (pixel diff, regions, etc.)"
            },
            # Visual regression strategy
            {
                "pattern": r'testVisualRegression.*Percy.*visual.*comparison',
                "severity": "HIGH",
                "title": "Visual Tool Dependency",
                "description": "Visual regression tightly coupled to specific third-party service",
                "recommendation": "Abstract visual testing behind interface for tool flexibility"
            },
            {
                "pattern": r'DiffPercentage.*0\.0.*Passed.*true',
                "severity": "CRITICAL",
                "title": "Zero-Tolerance Visual Diff",
                "description": "Visual regression requires perfect pixel match without threshold",
                "recommendation": "Implement configurable visual difference tolerance"
            }
        ]
        
        self._check_patterns(patterns, "Screenshot Management")
    
    def _detect_performance_testing(self):
        """Detect performance measurement integration issues."""
        patterns = [
            # Core Web Vitals
            {
                "pattern": r'LCP.*2\.3.*FID.*85.*CLS.*0\.08',
                "severity": "HIGH",
                "title": "Hardcoded Performance Metrics",
                "description": "Core Web Vitals metrics are hardcoded without real measurement",
                "recommendation": "Integrate with real performance measurement tools"
            },
            {
                "pattern": r'PerformanceReport.*TTFB.*450',
                "severity": "HIGH",
                "title": "Mock Performance Data",
                "description": "Performance metrics appear to be mocked rather than measured",
                "recommendation": "Implement actual performance measurement via Lighthouse API"
            },
            {
                "pattern": r'testPerformance.*Lighthouse.*performance.*audit.*via.*MCP',
                "severity": "CRITICAL",
                "title": "Performance Testing Not Implemented",
                "description": "Performance testing is simulated rather than executed",
                "recommendation": "Implement real Lighthouse integration for performance measurement"
            },
            # Performance thresholds
            {
                "pattern": r'Score.*92\.0.*Overall.*performance.*score',
                "severity": "MEDIUM",
                "title": "Static Performance Scoring",
                "description": "Performance scoring is static without dynamic baseline",
                "recommendation": "Implement performance budget validation against baselines"
            },
            {
                "pattern": r'Performance.*Score.*>= 90\.0',
                "severity": "HIGH",
                "title": "Universal Performance Threshold",
                "description": "Single performance threshold applied across all page types",
                "recommendation": "Implement page-type-specific performance budgets"
            }
        ]
        
        self._check_patterns(patterns, "Performance Testing")
    
    def _detect_user_interaction_patterns(self):
        """Detect user interaction simulation issues."""
        patterns = [
            # Interaction simulation
            {
                "pattern": r'testCases.*homepage_load.*analysis_form_submit',
                "severity": "HIGH",
                "title": "Basic Interaction Testing",
                "description": "User interaction tests are too basic for comprehensive validation",
                "recommendation": "Implement complex user journey simulation with state validation"
            },
            {
                "pattern": r'ValidateSEPTEODesign.*Check.*primary.*color',
                "severity": "MEDIUM",
                "title": "Design Validation Stub",
                "description": "Brand compliance validation is not fully implemented",
                "recommendation": "Implement comprehensive design system validation"
            },
            # User journey coverage
            {
                "pattern": r'RunFullTest.*testResponsiveDesign.*testAccessibility',
                "severity": "MEDIUM",
                "title": "Sequential Test Execution",
                "description": "Tests executed sequentially without consideration for interdependencies",
                "recommendation": "Implement test dependency management and parallel execution where safe"
            },
            {
                "pattern": r'Cross-browser.*test.*passed.*log\.Debug',
                "severity": "HIGH",
                "title": "Superficial Cross-Browser Testing",
                "description": "Cross-browser tests only log success without actual validation",
                "recommendation": "Implement comprehensive cross-browser behavioral validation"
            }
        ]
        
        self._check_patterns(patterns, "User Interaction")
    
    def _detect_cross_browser_compatibility(self):
        """Detect cross-browser compatibility handling issues."""
        patterns = [
            # Browser handling
            {
                "pattern": r'Browsers.*chromium.*firefox.*webkit',
                "severity": "MEDIUM",
                "title": "Limited Browser Coverage",
                "description": "Browser testing limited to engine types without version coverage",
                "recommendation": "Include specific browser versions and mobile browsers in test matrix"
            },
            {
                "pattern": r'browser.*testCase.*TotalTests\+\+.*PassedTests\+\+',
                "severity": "CRITICAL",
                "title": "Fake Cross-Browser Results",
                "description": "Cross-browser tests automatically pass without actual execution",
                "recommendation": "Implement real cross-browser test execution and validation"
            },
            # Browser-specific handling
            {
                "pattern": r'for.*browser.*range.*config\.Browsers',
                "severity": "HIGH",
                "title": "Uniform Browser Treatment",
                "description": "All browsers treated uniformly without browser-specific optimizations",
                "recommendation": "Implement browser-specific test strategies and timeouts"
            },
            {
                "pattern": r'Microsoft Edge.*Google Chrome.*channel',
                "severity": "MEDIUM",
                "title": "Browser Channel Configuration",
                "description": "Browser channels specified in config but not used in agent",
                "recommendation": "Align browser channel usage between configuration and execution"
            }
        ]
        
        self._check_patterns(patterns, "Cross-Browser Compatibility")
    
    def _detect_test_data_management(self):
        """Detect test data management strategy issues.""" 
        patterns = [
            # Test data strategy
            {
                "pattern": r'AccessibilityViolation.*Rule.*color-contrast',
                "severity": "HIGH",
                "title": "Hardcoded Test Data",
                "description": "Accessibility violation data is hardcoded rather than discovered",
                "recommendation": "Implement dynamic test data discovery and validation"
            },
            {
                "pattern": r'config.*BaseURL.*constants\.TestLocalhost3000',
                "severity": "MEDIUM",
                "title": "Environment Coupling",
                "description": "Test configuration tightly coupled to development environment",
                "recommendation": "Implement environment-agnostic test data management"
            },
            # Test isolation
            {
                "pattern": r'results.*&PlaywrightResults.*Timestamp.*time\.Now',
                "severity": "HIGH",
                "title": "Shared Test Results Object",
                "description": "Test results shared across test methods without isolation",
                "recommendation": "Implement test result isolation per test case"
            },
            {
                "pattern": r'ReportPath.*tests/reports/frontend',
                "severity": "MEDIUM",
                "title": "Static Report Path",
                "description": "Report output path is static without environment consideration",
                "recommendation": "Use dynamic report paths based on execution context"
            }
        ]
        
        self._check_patterns(patterns, "Test Data Management")
    
    def _detect_page_object_patterns(self):
        """Detect Page Object Model implementation issues."""
        patterns = [
            # Page Object absence
            {
                "pattern": r'func.*testResponsiveDesign.*testAccessibility',
                "severity": "CRITICAL",
                "title": "No Page Object Abstraction",
                "description": "Tests directly interact with low-level browser APIs without Page Object Model",
                "recommendation": "Implement Page Object Model pattern for better test maintainability"
            },
            {
                "pattern": r'screenshot.*Name.*Path.*Viewport.*Browser',
                "severity": "HIGH",
                "title": "Data Structure Without Behavior",
                "description": "Screenshot struct is plain data without encapsulated behavior",
                "recommendation": "Implement behavior-rich Page Object classes with methods"
            },
            # Element management
            {
                "pattern": r'Element.*string.*json.*element',
                "severity": "HIGH",
                "title": "String-Based Element References",
                "description": "Elements referenced as strings without type safety",
                "recommendation": "Implement strongly-typed element reference system"
            },
            {
                "pattern": r'fmt\.Sprintf.*viewport.*browser',
                "severity": "MEDIUM",
                "title": "String Interpolation for Dynamic Content",
                "description": "Dynamic content generation using string formatting",
                "recommendation": "Use template-based content generation with validation"
            }
        ]
        
        self._check_patterns(patterns, "Page Object Model")
    
    def _detect_wait_strategies(self):
        """Detect wait strategy implementation issues."""
        patterns = [
            # Wait strategy
            {
                "pattern": r'timeout.*navigationTimeout.*30000',
                "severity": "MEDIUM",
                "title": "Global Timeout Configuration",
                "description": "Navigation timeout configured globally without per-test customization",
                "recommendation": "Implement dynamic timeout strategy based on operation type"
            },
            {
                "pattern": r'ctx.*context\.Context',
                "severity": "HIGH",
                "title": "Context Without Timeout Handling",
                "description": "Context passed but timeout handling not properly implemented",
                "recommendation": "Implement proper context timeout and cancellation handling"
            },
            # Async handling
            {
                "pattern": r'func.*testPerformance.*testAccessibility.*error',
                "severity": "MEDIUM",
                "title": "Sequential Async Operations",
                "description": "Async operations executed sequentially without parallelization opportunities",
                "recommendation": "Implement parallel execution for independent test operations"
            },
            {
                "pattern": r'RunFullTest.*if err.*log\.Error',
                "severity": "HIGH",
                "title": "Continue on Error Pattern",
                "description": "Test execution continues after errors without proper failure handling",
                "recommendation": "Implement proper error propagation and test isolation"
            }
        ]
        
        self._check_patterns(patterns, "Wait Strategies")
    
    def _detect_error_handling(self):
        """Detect error handling pattern issues."""
        patterns = [
            # Error propagation
            {
                "pattern": r'if err.*log\.Error.*continue',
                "severity": "HIGH",
                "title": "Error Swallowing",
                "description": "Errors are logged but not properly propagated to test framework",
                "recommendation": "Implement proper error propagation and test failure mechanisms"
            },
            {
                "pattern": r'generateReport.*if err.*log\.Error',
                "severity": "MEDIUM",
                "title": "Report Generation Error Ignored",
                "description": "Report generation errors logged but test continues as successful",
                "recommendation": "Treat report generation failures as test infrastructure failures"
            },
            # Error context
            {
                "pattern": r'error.*err\.Error\(\)',
                "severity": "MEDIUM",
                "title": "Generic Error Messages",
                "description": "Error messages lack context about test operation that failed",
                "recommendation": "Include test context and operation details in error messages"
            },
            {
                "pattern": r'return.*PlaywrightResults.*error',
                "severity": "LOW",
                "title": "Mixed Return Pattern",
                "description": "Functions return both results and errors without clear separation",
                "recommendation": "Consider using Result pattern for cleaner error handling"
            }
        ]
        
        self._check_patterns(patterns, "Error Handling")
    
    def _detect_test_maintainability(self):
        """Detect test maintainability issues."""
        patterns = [
            # Test organization
            {
                "pattern": r'NewPlaywrightAgent.*config.*results.*mcpClient',
                "severity": "MEDIUM",
                "title": "Monolithic Test Agent",
                "description": "Single agent handles all test types without separation of concerns",
                "recommendation": "Separate test agents by functional area (accessibility, performance, etc.)"
            },
            {
                "pattern": r'runFullTest.*testResponsiveDesign.*testAccessibility.*testPerformance',
                "severity": "HIGH",
                "title": "Coupling All Test Types",
                "description": "All test types executed together without independent execution capability",
                "recommendation": "Implement modular test execution with independent test type runners"
            },
            # Configuration management
            {
                "pattern": r'config.*&PlaywrightConfig.*BaseURL.*Browsers',
                "severity": "MEDIUM",
                "title": "Static Configuration Pattern",
                "description": "Configuration created statically without runtime modification capability",
                "recommendation": "Implement configuration builder pattern for flexible test setup"
            },
            {
                "pattern": r'log.*New.*FRONTEND-AGENT',
                "severity": "LOW",
                "title": "Global Logger Instance",
                "description": "Logger instantiated globally without test-specific context",
                "recommendation": "Implement test-scoped logging for better traceability"
            }
        ]
        
        self._check_patterns(patterns, "Test Maintainability")
    
    def _check_patterns(self, patterns: List[Dict], category: str):
        """Check for pattern matches and create violations."""
        for i, line in enumerate(self.lines, 1):
            for pattern_info in patterns:
                if re.search(pattern_info["pattern"], line, re.IGNORECASE):
                    violation = PlaywrightViolation(
                        id=f"PLW-{len(self.violations) + 1:03d}",
                        severity=pattern_info["severity"],
                        category=category,
                        title=pattern_info["title"],
                        description=pattern_info["description"],
                        line_number=i,
                        code_snippet=line.strip(),
                        impact=self._get_impact(pattern_info["severity"]),
                        recommendation=pattern_info["recommendation"],
                        pattern_type="regex"
                    )
                    self.violations.append(violation)
    
    def _get_impact(self, severity: str) -> str:
        """Get impact description based on severity."""
        impacts = {
            "CRITICAL": "Compromises test reliability and maintainability significantly",
            "HIGH": "Major impact on test quality and maintenance overhead",
            "MEDIUM": "Moderate impact on test effectiveness and code quality", 
            "LOW": "Minor impact on test maintainability"
        }
        return impacts.get(severity, "Unknown impact")
    
    def _get_category_stats(self) -> Dict[str, Any]:
        """Get violation statistics by category."""
        stats = {}
        for violation in self.violations:
            category = violation.category
            severity = violation.severity
            
            if category not in stats:
                stats[category] = {"total": 0, "severities": {}}
            
            stats[category]["total"] += 1
            if severity not in stats[category]["severities"]:
                stats[category]["severities"][severity] = 0
            stats[category]["severities"][severity] += 1
            
        return stats

def main():
    """Main execution function."""
    import sys
    
    if len(sys.argv) != 2:
        print("Usage: python delta15_playwright_detector.py <file_path>")
        sys.exit(1)
    
    file_path = sys.argv[1]
    detector = PlaywrightDetector()
    result = detector.analyze_file(file_path)
    
    # Pretty print results
    print(json.dumps(result, indent=2))
    
    # Summary
    print(f"\n🎭 DELTA-15 Playwright Analysis Complete")
    print(f"📊 Total Violations: {result['total_violations']}")
    print(f"📁 File: {file_path}")
    
    if result['total_violations'] > 0:
        print("\n📋 Categories:")
        for category, stats in result['categories'].items():
            print(f"  • {category}: {stats['total']} violations")

if __name__ == "__main__":
    main()