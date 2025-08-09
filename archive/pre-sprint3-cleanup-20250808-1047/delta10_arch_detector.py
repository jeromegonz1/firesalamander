#!/usr/bin/env python3
"""
Delta 10 Architecture Detector for Message Hardcoding Violations

This script analyzes messages.go for hardcoded strings and provides
architectural recommendations for message management systems.
"""

import re
import json
from typing import Dict, List, Tuple, Set
from pathlib import Path


class MessageArchitectureAnalyzer:
    """Analyzes hardcoded messages and provides architectural recommendations."""
    
    def __init__(self, file_path: str):
        self.file_path = file_path
        self.violations = []
        self.message_categories = {
            'error': [],
            'info': [],
            'warning': [],
            'success': [],
            'ui': [],
            'log': [],
            'help': [],
            'time': [],
            'phase': [],
            'mode': [],
            'seo': [],
            'ai': []
        }
        
    def extract_messages(self) -> List[Dict]:
        """Extract all hardcoded message strings from the Go file."""
        try:
            with open(self.file_path, 'r', encoding='utf-8') as file:
                content = file.read()
        except FileNotFoundError:
            print(f"Error: File {self.file_path} not found")
            return []
        
        # Pattern to match constant definitions with string values
        pattern = r'(\w+)\s*=\s*"([^"]*)"'
        matches = re.findall(pattern, content)
        
        violations = []
        for const_name, message_text in matches:
            violation = {
                'constant_name': const_name,
                'message_text': message_text,
                'line_number': self._find_line_number(content, const_name),
                'category': self._categorize_message(const_name, message_text),
                'i18n_potential': self._assess_i18n_potential(message_text),
                'user_facing': self._is_user_facing(const_name, message_text),
                'technical_level': self._assess_technical_level(const_name, message_text)
            }
            violations.append(violation)
            
        return violations
    
    def _find_line_number(self, content: str, const_name: str) -> int:
        """Find the line number where a constant is defined."""
        lines = content.split('\n')
        for i, line in enumerate(lines, 1):
            if const_name in line and '=' in line:
                return i
        return -1
    
    def _categorize_message(self, const_name: str, message_text: str) -> str:
        """Categorize message based on constant name and content."""
        const_lower = const_name.lower()
        message_lower = message_text.lower()
        
        # Error messages
        if 'err' in const_lower or 'error' in message_lower or 'failed' in message_lower:
            return 'error'
        
        # UI messages
        if const_lower.startswith('ui'):
            return 'ui'
        
        # Log messages  
        if const_lower.startswith('log'):
            return 'log'
            
        # Help messages
        if const_lower.startswith('help'):
            return 'help'
            
        # Time estimates
        if 'time' in const_lower or 'seconds' in message_lower or 'minutes' in message_lower:
            return 'time'
            
        # Phase messages
        if 'phase' in const_lower or 'progress' in message_lower or 'analyzing' in message_lower:
            return 'phase'
            
        # Mode messages
        if 'mode' in const_lower or 'interface' in message_lower or 'api' in message_lower:
            return 'mode'
            
        # SEO messages
        if const_lower.startswith('seo'):
            return 'seo'
            
        # AI messages
        if const_lower.startswith('ai'):
            return 'ai'
            
        # Server/success messages
        if 'start' in message_lower or 'complete' in message_lower or 'available' in message_lower:
            return 'success'
            
        return 'info'
    
    def _assess_i18n_potential(self, message_text: str) -> str:
        """Assess internationalization potential of the message."""
        # High potential: user-facing messages with natural language
        if any(word in message_text.lower() for word in ['please', 'try again', 'error', 'complete', 'analyzing']):
            return 'high'
        
        # Medium potential: status messages and notifications
        if any(word in message_text.lower() for word in ['started', 'stopped', 'progress', 'available']):
            return 'medium'
        
        # Low potential: technical identifiers or URLs
        if 'https://' in message_text or message_text.isupper() or len(message_text.split()) == 1:
            return 'low'
            
        return 'medium'
    
    def _is_user_facing(self, const_name: str, message_text: str) -> bool:
        """Determine if message is user-facing."""
        const_lower = const_name.lower()
        
        # UI messages are definitely user-facing
        if const_lower.startswith('ui'):
            return True
            
        # Log messages are typically not user-facing
        if const_lower.startswith('log'):
            return False
            
        # Error messages shown to users
        if const_lower.startswith('err') and not const_lower.startswith('log'):
            return True
            
        # Help messages are user-facing
        if const_lower.startswith('help'):
            return True
            
        # Phase and time messages are often user-facing
        if any(prefix in const_lower for prefix in ['phase', 'time', 'analysis']):
            return True
            
        return False
    
    def _assess_technical_level(self, const_name: str, message_text: str) -> str:
        """Assess the technical level of the message."""
        message_lower = message_text.lower()
        
        # High technical: developer/system messages
        if any(term in message_lower for term in ['json', 'http', 'server', 'database', 'template', 'configuration']):
            return 'technical'
            
        # Medium technical: process/status messages
        if any(term in message_lower for term in ['analysis', 'api', 'timeout', 'connection']):
            return 'semi-technical'
            
        # Low technical: user-friendly messages
        return 'user-friendly'
    
    def analyze(self) -> Dict:
        """Perform complete analysis of messages."""
        self.violations = self.extract_messages()
        
        # Categorize violations
        for violation in self.violations:
            category = violation['category']
            if category in self.message_categories:
                self.message_categories[category].append(violation)
        
        return self._generate_analysis_report()
    
    def _generate_analysis_report(self) -> Dict:
        """Generate comprehensive analysis report."""
        total_violations = len(self.violations)
        
        # Category statistics
        category_stats = {}
        for category, messages in self.message_categories.items():
            category_stats[category] = {
                'count': len(messages),
                'percentage': (len(messages) / total_violations * 100) if total_violations > 0 else 0
            }
        
        # I18n potential analysis
        i18n_analysis = {
            'high': len([v for v in self.violations if v['i18n_potential'] == 'high']),
            'medium': len([v for v in self.violations if v['i18n_potential'] == 'medium']),
            'low': len([v for v in self.violations if v['i18n_potential'] == 'low'])
        }
        
        # User-facing analysis
        user_facing_count = len([v for v in self.violations if v['user_facing']])
        
        return {
            'summary': {
                'total_violations': total_violations,
                'file_analyzed': self.file_path,
                'user_facing_messages': user_facing_count,
                'internal_messages': total_violations - user_facing_count
            },
            'category_breakdown': category_stats,
            'i18n_potential': i18n_analysis,
            'violations': self.violations,
            'architectural_recommendations': self._generate_recommendations()
        }
    
    def _generate_recommendations(self) -> Dict:
        """Generate architectural recommendations for message management."""
        return {
            'immediate_actions': [
                "Implement a centralized message management system",
                "Create message constants with semantic naming conventions",
                "Separate user-facing messages from technical/log messages",
                "Implement message templates for dynamic content"
            ],
            'architectural_patterns': {
                'message_manager': {
                    'description': "Centralized message management with category-based organization",
                    'benefits': ["Single source of truth", "Easy maintenance", "Consistent formatting"],
                    'implementation': "Create MessageManager interface with category-specific implementations"
                },
                'i18n_ready_structure': {
                    'description': "Prepare architecture for internationalization",
                    'benefits': ["Future-proof for multiple languages", "Consistent message keys", "Easy localization"],
                    'implementation': "Use message keys instead of hardcoded strings, implement locale-based message resolution"
                },
                'message_templates': {
                    'description': "Template-based messages for dynamic content",
                    'benefits': ["Reusable message patterns", "Consistent formatting", "Easy parameterization"],
                    'implementation': "Implement template engine for messages with placeholder support"
                }
            },
            'refactoring_strategy': {
                'phase_1': "Extract all hardcoded messages to constants (current state - partially done)",
                'phase_2': "Categorize messages by type and audience",
                'phase_3': "Implement message manager with category-based organization",
                'phase_4': "Add template support for dynamic messages",
                'phase_5': "Prepare i18n structure with message keys"
            },
            'code_quality_improvements': [
                "Use consistent naming conventions for message constants",
                "Group related messages together",
                "Add documentation for message categories",
                "Implement message validation for required parameters"
            ]
        }


def main():
    """Main execution function."""
    file_path = "/Users/jeromegonzalez/claude-code/fire-salamander/internal/messages/messages.go"
    
    analyzer = MessageArchitectureAnalyzer(file_path)
    analysis_result = analyzer.analyze()
    
    # Save analysis to JSON file
    output_file = "/Users/jeromegonzalez/claude-code/fire-salamander/delta10_arch_analysis.json"
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(analysis_result, f, indent=2, ensure_ascii=False)
    
    print(f"Architecture analysis completed. Results saved to: {output_file}")
    print(f"\nSummary:")
    print(f"- Total violations detected: {analysis_result['summary']['total_violations']}")
    print(f"- User-facing messages: {analysis_result['summary']['user_facing_messages']}")
    print(f"- Internal messages: {analysis_result['summary']['internal_messages']}")
    
    print(f"\nCategory breakdown:")
    for category, stats in analysis_result['category_breakdown'].items():
        if stats['count'] > 0:
            print(f"- {category}: {stats['count']} messages ({stats['percentage']:.1f}%)")
    
    print(f"\nI18n potential:")
    for level, count in analysis_result['i18n_potential'].items():
        print(f"- {level}: {count} messages")


if __name__ == "__main__":
    main()