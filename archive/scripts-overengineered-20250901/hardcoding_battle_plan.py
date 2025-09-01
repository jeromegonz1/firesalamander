#!/usr/bin/env python3
"""
🚨 FIRE SALAMANDER - HARDCODING BATTLE PLAN ANALYZER 🚨
Mission: Analyser les 1,862 violations pour plan d'attaque optimal
"""

import os
import re
import json
from collections import defaultdict
import datetime

class HardcodingBattlePlan:
    def __init__(self):
        self.patterns = {
            'http_codes': r'\b(200|201|202|204|301|302|304|400|401|403|404|405|409|422|429|500|502|503|504)\b',
            'strings': r'"[A-Za-z][^"]{4,}"',  # Strings de 5+ caractères
            'urls': r'(http://|https://|localhost|127\.0\.0\.1)',
            'ports': r'\b(8080|3000|5432|6379|27017|80|443)\b',
            'paths': r'"/[^"]+/"',
            'messages': r'"(Error|Warning|Info|Success|Failed|Invalid|Required|Missing|Found|Loading|Complete)[^"]*"',
            'config_keys': r'"[a-z][a-z_]+[a-z]"',  # Config-style strings
            'css_classes': r'"[a-z-]+(?:\s+[a-z-]+)*"',  # CSS classes
            'json_fields': r'`json:"[^"]+"`',
        }
        
        # Patterns à ignorer (légitimes)
        self.safe_patterns = [
            r'json:',  # JSON struct tags
            r'yaml:',  # YAML struct tags  
            r'import\s+',  # Import statements
            r'fmt\.Sprintf',  # Format strings
            r'_test\.go',  # Test files
            r'//.*',  # Comments
            r'package\s+',  # Package declarations
        ]
    
    def is_safe_violation(self, line, match):
        """Check if violation is in a safe context"""
        line_lower = line.lower()
        
        # Check safe patterns
        for pattern in self.safe_patterns:
            if re.search(pattern, line):
                return True
        
        # JSON struct tags are acceptable
        if 'json:' in line or 'yaml:' in line:
            return True
            
        # Import statements are safe
        if line.strip().startswith('import') or 'package ' in line:
            return True
            
        return False
    
    def analyze_file(self, filepath):
        """Analyze a single Go file for hardcoding violations"""
        violations = defaultdict(list)
        
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                lines = f.readlines()
                
            for line_num, line in enumerate(lines, 1):
                # Skip if line is safe
                if any(re.search(pattern, line) for pattern in self.safe_patterns):
                    continue
                    
                for vtype, pattern in self.patterns.items():
                    matches = re.findall(pattern, line)
                    if matches:
                        for match in matches:
                            if not self.is_safe_violation(line, match):
                                violations[vtype].append({
                                    'line': line_num,
                                    'content': line.strip(),
                                    'match': match,
                                    'severity': self.get_severity(vtype, match)
                                })
        except Exception as e:
            print(f"❌ Error processing {filepath}: {e}")
            
        return violations
    
    def get_severity(self, vtype, match):
        """Determine severity of violation"""
        if vtype in ['http_codes', 'urls', 'ports']:
            return 'CRITICAL'
        elif vtype in ['messages', 'paths']:
            return 'HIGH'
        elif vtype in ['strings', 'config_keys']:
            return 'MEDIUM'
        else:
            return 'LOW'
    
    def analyze_codebase(self):
        """Analyze entire codebase"""
        print("🔍 FIRE SALAMANDER - Analyzing codebase for hardcoding violations...")
        
        all_violations = defaultdict(lambda: defaultdict(list))
        file_stats = {}
        severity_stats = defaultdict(int)
        
        # Walk through Go files
        for root, dirs, files in os.walk('.'):
            # Skip certain directories
            if any(skip in root for skip in ['archive', 'vendor', '.git', 'node_modules']):
                continue
                
            for file in files:
                if file.endswith('.go'):
                    filepath = os.path.join(root, file)
                    file_violations = self.analyze_file(filepath)
                    
                    if file_violations:
                        total_violations = sum(len(violations) for violations in file_violations.values())
                        file_stats[filepath] = {
                            'total': total_violations,
                            'by_type': {vtype: len(violations) for vtype, violations in file_violations.items()},
                            'violations': file_violations
                        }
                        
                        # Add to global stats
                        for vtype, violations in file_violations.items():
                            all_violations[filepath][vtype] = violations
                            for violation in violations:
                                severity_stats[violation['severity']] += 1
        
        return self.generate_battle_plan(all_violations, file_stats, severity_stats)
    
    def generate_battle_plan(self, all_violations, file_stats, severity_stats):
        """Generate strategic battle plan"""
        
        # Sort files by violation count
        top_files = sorted(file_stats.items(), key=lambda x: x[1]['total'], reverse=True)
        
        # Calculate totals by type
        type_totals = defaultdict(int)
        for filepath, file_data in file_stats.items():
            for vtype, count in file_data['by_type'].items():
                type_totals[vtype] += count
        
        # Create battle plan
        battle_plan = {
            'timestamp': datetime.datetime.now().isoformat(),
            'summary': {
                'total_files': len(file_stats),
                'total_violations': sum(severity_stats.values()),
                'by_severity': dict(severity_stats),
                'by_type': dict(type_totals)
            },
            'battle_phases': [],
            'top_files': top_files[:20],  # Top 20 worst files
            'detailed_violations': all_violations
        }
        
        # Define attack phases
        phases = [
            {
                'name': 'BATCH 1 - HTTP Status Codes',
                'types': ['http_codes'],
                'priority': 'CRITICAL',
                'estimated_violations': type_totals.get('http_codes', 0),
                'estimated_time': '20min'
            },
            {
                'name': 'BATCH 2 - Error Messages',
                'types': ['messages'],
                'priority': 'HIGH',
                'estimated_violations': type_totals.get('messages', 0),
                'estimated_time': '30min'
            },
            {
                'name': 'BATCH 3 - URLs and Ports',
                'types': ['urls', 'ports'],
                'priority': 'CRITICAL',
                'estimated_violations': type_totals.get('urls', 0) + type_totals.get('ports', 0),
                'estimated_time': '25min'
            },
            {
                'name': 'BATCH 4 - Paths and Routes',
                'types': ['paths'],
                'priority': 'HIGH', 
                'estimated_violations': type_totals.get('paths', 0),
                'estimated_time': '20min'
            },
            {
                'name': 'BATCH 5 - Generic Strings',
                'types': ['strings', 'config_keys'],
                'priority': 'MEDIUM',
                'estimated_violations': type_totals.get('strings', 0) + type_totals.get('config_keys', 0),
                'estimated_time': '40min'
            }
        ]
        
        battle_plan['battle_phases'] = phases
        
        return battle_plan
    
    def print_battle_report(self, battle_plan):
        """Print detailed battle report"""
        print("\n" + "="*60)
        print("🚨 FIRE SALAMANDER - HARDCODING BATTLE PLAN 🚨")
        print("="*60)
        
        summary = battle_plan['summary']
        print(f"📊 SITUATION ANALYSIS:")
        print(f"   Total files with violations: {summary['total_files']}")
        print(f"   Total violations detected: {summary['total_violations']}")
        
        print(f"\n🎯 VIOLATIONS BY SEVERITY:")
        for severity, count in summary['by_severity'].items():
            print(f"   {severity}: {count}")
        
        print(f"\n📈 VIOLATIONS BY TYPE:")
        for vtype, count in sorted(summary['by_type'].items(), key=lambda x: x[1], reverse=True):
            print(f"   {vtype}: {count}")
        
        print(f"\n🎪 TOP 10 FILES TO ATTACK FIRST:")
        for i, (filepath, stats) in enumerate(battle_plan['top_files'][:10], 1):
            print(f"   {i:2}. {stats['total']:3} violations: {filepath}")
            for vtype, count in sorted(stats['by_type'].items(), key=lambda x: x[1], reverse=True):
                print(f"        {vtype}: {count}")
        
        print(f"\n🗡️  BATTLE PHASES:")
        total_estimated = 0
        for phase in battle_plan['battle_phases']:
            total_estimated += phase['estimated_violations']
            print(f"   {phase['name']} ({phase['priority']})")
            print(f"      Violations: {phase['estimated_violations']}")
            print(f"      Time: {phase['estimated_time']}")
            print()
        
        remaining = summary['total_violations'] - total_estimated
        reduction_percent = (total_estimated / summary['total_violations']) * 100
        
        print(f"📊 BATTLE IMPACT PROJECTION:")
        print(f"   Total violations to eliminate: {total_estimated}")
        print(f"   Estimated remaining: {remaining}")
        print(f"   Projected reduction: {reduction_percent:.1f}%")
        
        return battle_plan

def main():
    analyzer = HardcodingBattlePlan()
    battle_plan = analyzer.analyze_codebase()
    
    # Print report
    analyzer.print_battle_report(battle_plan)
    
    # Save detailed plan
    with open('violation_battle_plan.json', 'w') as f:
        json.dump(battle_plan, f, indent=2)
    
    print(f"\n💾 Detailed battle plan saved to: violation_battle_plan.json")
    print(f"🚀 Ready to commence Operation Zero Hardcoding!")

if __name__ == '__main__':
    main()