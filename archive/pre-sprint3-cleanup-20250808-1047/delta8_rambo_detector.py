#!/usr/bin/env python3
"""
🔥💀 DELTA-8 RAMBO ASSAULT DETECTOR 💀🔥
Mission: Hunt down hardcoded web violations in enemy territory

"Nothing is over! Nothing! You just don't turn it off!"
- John Rambo
"""

import re
import json
import os
from typing import List, Dict, Tuple, Set
from dataclasses import dataclass

@dataclass
class RamboViolation:
    """A single violation detected by RAMBO"""
    line_number: int
    line_content: str
    violation_type: str
    hardcoded_value: str
    suggested_constant: str
    rambo_comment: str

class Delta8RamboDetector:
    def __init__(self):
        """Initialize RAMBO detection systems 🎯"""
        self.violations: List[RamboViolation] = []
        self.constants_mapping: Dict[str, str] = {}
        
        # 💥 RAMBO's hardcoded detection patterns
        self.web_patterns = {
            # HTTP Methods
            'http_methods': r'"(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)"',
            
            # HTTP Status Codes  
            'status_codes': r'(http\.Status\w+|\b[1-5]\d{2}\b)',
            
            # Content-Type headers
            'content_types': r'"(text/html|application/json|text/plain|application/xml|multipart/form-data)[^"]*"',
            
            # HTTP Headers
            'http_headers': r'"(Content-Type|Content-Disposition|Authorization|Accept|User-Agent|Cache-Control)[^"]*"',
            
            # URLs and endpoints
            'urls_endpoints': r'"(/[^"]*)"',
            
            # Port numbers
            'port_numbers': r'\b(80[0-9]{2}|[1-9]\d{3,4})\b',
            
            # JSON field names
            'json_fields': r'`"([a-zA-Z_][a-zA-Z0-9_]*)"[^`]*`',
            
            # HTML content
            'html_content': r'"(<!DOCTYPE html>|<[^>]+>)[^"]*"',
            
            # Error messages
            'error_messages': r'"(Erreur|Error|Failed|Invalid)[^"]*"',
            
            # Test descriptions
            'test_descriptions': r'"([^"]*Test[^"]*|[^"]*test[^"]*)"',
            
            # Configuration values
            'config_values': r'"(Fire Salamander[^"]*|running|web_server)[^"]*"',
            
            # File extensions
            'file_extensions': r'"[^"]*\.(html|json|csv|xml|txt)"',
            
            # Protocol schemes
            'protocols': r'"(http://|https://)[^"]*"',
            
            # MIME types
            'mime_types': r'"(attachment; filename=)[^"]*"',
            
            # Route patterns
            'route_patterns': r'"(/web/[^"]*)"',
        }
        
        # 🎯 RAMBO's violation comments
        self.rambo_comments = {
            'http_methods': "FIRST BLOOD! HTTP method needs extraction 💀",
            'status_codes': "Target locked on status codes! 🎯",
            'content_types': "Content-Type headers in crosshairs! ⚔️",
            'http_headers': "HTTP headers spotted - engage! 💥",
            'urls_endpoints': "URL endpoint detected - neutralize! 🔥",
            'port_numbers': "Port number hardcoded - tactical strike needed! 💣",
            'json_fields': "JSON field exposed - take the shot! 🎯",
            'html_content': "HTML content hardcoded - eliminate! ⚡",
            'error_messages': "Error message needs extraction - go loud! 🚨",
            'test_descriptions': "Test string hardcoded - silent takedown! 🔫",
            'config_values': "Config value exposed - breach and clear! 🏴‍☠️",
            'file_extensions': "File extension hardcoded - explosive ordnance! 💥",
            'protocols': "Protocol scheme detected - air strike! ✈️",
            'mime_types': "MIME type in the open - sniper shot! 🎯",
            'route_patterns': "Route pattern exposed - tactical assault! ⚔️",
        }

    def scan_target(self, file_path: str) -> None:
        """🔥 RAMBO scanning operation begins 🔥"""
        print("🔥💀 FIRST BLOOD - Scanning web_test.go... 💀🔥")
        print("🎯 'They drew first blood, not me!' - John Rambo")
        print()
        
        if not os.path.exists(file_path):
            print(f"❌ Target not found: {file_path}")
            return
            
        with open(file_path, 'r', encoding='utf-8') as f:
            lines = f.readlines()
            
        total_kills = 0
        
        for line_num, line in enumerate(lines, 1):
            line_stripped = line.strip()
            
            for pattern_type, pattern in self.web_patterns.items():
                matches = re.finditer(pattern, line)
                
                for match in matches:
                    hardcoded_value = match.group(1) if match.groups() else match.group(0)
                    
                    # Skip empty matches or single characters
                    if len(hardcoded_value.strip()) <= 1:
                        continue
                        
                    # Generate constant name
                    constant_name = self._generate_constant_name(hardcoded_value, pattern_type)
                    
                    violation = RamboViolation(
                        line_number=line_num,
                        line_content=line_stripped,
                        violation_type=pattern_type,
                        hardcoded_value=hardcoded_value,
                        suggested_constant=constant_name,
                        rambo_comment=self.rambo_comments.get(pattern_type, "Target acquired! 🎯")
                    )
                    
                    self.violations.append(violation)
                    self.constants_mapping[hardcoded_value] = constant_name
                    total_kills += 1
        
        print(f"🎯 RAMBO SCAN COMPLETE - {total_kills} cibles neutralisées! 💥")
        print()

    def _generate_constant_name(self, value: str, pattern_type: str) -> str:
        """Generate appropriate constant names - RAMBO style"""
        
        # Clean the value
        clean_value = re.sub(r'[^\w\s\-\.]', '', value)
        clean_value = clean_value.strip()
        
        if pattern_type == 'http_methods':
            return f"HTTP_METHOD_{clean_value.upper()}"
        elif pattern_type == 'status_codes':
            if 'StatusOK' in value:
                return "HTTP_STATUS_OK"
            elif 'Status' in value:
                return f"HTTP_{value.replace('http.', '').upper()}"
            return f"HTTP_STATUS_{clean_value}"
        elif pattern_type == 'content_types':
            if 'text/html' in value:
                return "CONTENT_TYPE_HTML"
            elif 'application/json' in value:
                return "CONTENT_TYPE_JSON"
            elif 'text/plain' in value:
                return "CONTENT_TYPE_TEXT"
            return "CONTENT_TYPE_DEFAULT"
        elif pattern_type == 'http_headers':
            if 'Content-Type' in value:
                return "HEADER_CONTENT_TYPE"
            elif 'Content-Disposition' in value:
                return "HEADER_CONTENT_DISPOSITION"
            return f"HEADER_{clean_value.upper().replace('-', '_')}"
        elif pattern_type == 'urls_endpoints':
            if value == '/':
                return "ENDPOINT_ROOT"
            elif '/web/health' in value:
                return "ENDPOINT_HEALTH"
            elif '/web/download' in value:
                return "ENDPOINT_DOWNLOAD"
            return f"ENDPOINT_{clean_value.upper().replace('/', '_').replace('-', '_')}"
        elif pattern_type == 'port_numbers':
            return f"TEST_PORT_{clean_value}"
        elif pattern_type == 'error_messages':
            words = clean_value.split()[:3]  # First 3 words
            return f"ERROR_MSG_{'_'.join(word.upper() for word in words)}"
        elif pattern_type == 'config_values':
            if 'Fire Salamander' in value:
                return "APP_NAME_TEST"
            elif 'running' in value:
                return "STATUS_RUNNING"
            elif 'web_server' in value:
                return "SERVICE_WEB_SERVER"
            return "CONFIG_VALUE_DEFAULT"
        elif pattern_type == 'file_extensions':
            ext = clean_value.split('.')[-1] if '.' in clean_value else clean_value
            return f"FILE_FORMAT_{ext.upper()}"
        elif pattern_type == 'protocols':
            return "PROTOCOL_HTTP" if 'http://' in value else "PROTOCOL_HTTPS"
        else:
            # Generic constant name
            words = clean_value.split()[:2]
            return f"WEB_{'_'.join(word.upper() for word in words if word)}"

    def generate_constants_file(self) -> str:
        """Generate Go constants file content"""
        constants_content = """package web

// 🔥💀 DELTA-8 RAMBO CONSTANTS - Web Testing Edition 💀🔥
// "I could have killed 'em all, I could've killed you. In town you're the law, 
// out here it's me. Don't push it! Don't push it or I'll give you a war you won't believe."
// - John Rambo

const (
\t// 🎯 HTTP Methods
"""
        
        # Group constants by type
        method_constants = []
        status_constants = []
        content_constants = []
        header_constants = []
        endpoint_constants = []
        port_constants = []
        other_constants = []
        
        for value, constant in self.constants_mapping.items():
            if constant.startswith('HTTP_METHOD_'):
                method_constants.append(f'\t{constant} = "{value}"')
            elif constant.startswith('HTTP_STATUS_') or constant.startswith('HTTP_'):
                if not any(c in constant for c in ['METHOD', 'HEADER']):
                    status_constants.append(f'\t{constant} = {value}' if value.isdigit() else f'\t{constant} = "{value}"')
            elif constant.startswith('CONTENT_TYPE_'):
                content_constants.append(f'\t{constant} = "{value}"')
            elif constant.startswith('HEADER_'):
                header_constants.append(f'\t{constant} = "{value}"')
            elif constant.startswith('ENDPOINT_'):
                endpoint_constants.append(f'\t{constant} = "{value}"')
            elif constant.startswith('TEST_PORT_'):
                port_constants.append(f'\t{constant} = {value}')
            else:
                other_constants.append(f'\t{constant} = "{value}"')
        
        if method_constants:
            constants_content += '\n'.join(method_constants) + '\n\n'
        
        if status_constants:
            constants_content += '\t// ⚡ HTTP Status Codes\n'
            constants_content += '\n'.join(status_constants) + '\n\n'
        
        if content_constants:
            constants_content += '\t// 🎯 Content Types\n'
            constants_content += '\n'.join(content_constants) + '\n\n'
        
        if header_constants:
            constants_content += '\t// 💥 HTTP Headers\n'
            constants_content += '\n'.join(header_constants) + '\n\n'
        
        if endpoint_constants:
            constants_content += '\t// 🔥 API Endpoints\n'
            constants_content += '\n'.join(endpoint_constants) + '\n\n'
        
        if port_constants:
            constants_content += '\t// 🚀 Test Ports\n'
            constants_content += '\n'.join(port_constants) + '\n\n'
        
        if other_constants:
            constants_content += '\t// ⚔️ Other Constants\n'
            constants_content += '\n'.join(other_constants) + '\n\n'
        
        constants_content += ')\n'
        return constants_content

    def print_rambo_results(self) -> None:
        """Print results with RAMBO style 💀"""
        print("🔥💀⚔️ RAMBO RÉSULTATS ⚔️💀🔥")
        print("🎯 'Nothing is over! NOTHING!'")
        print()
        
        if not self.violations:
            print("🏆 Mission accomplished - No hardcoded violations found!")
            print("🎖️ 'You did good, kid. Real good.'")
            return
        
        print(f"💀 Total kills: {len(self.violations)}")
        print(f"🎯 Cibles neutralisées:")
        print()
        
        # Group violations by type
        violations_by_type = {}
        for violation in self.violations:
            if violation.violation_type not in violations_by_type:
                violations_by_type[violation.violation_type] = []
            violations_by_type[violation.violation_type].append(violation)
        
        for violation_type, violations in violations_by_type.items():
            print(f"🔥 {violation_type.upper().replace('_', ' ')}: {len(violations)} kills")
            
            for violation in violations[:5]:  # Show first 5 of each type
                print(f"   💀 Line {violation.line_number}: {violation.hardcoded_value}")
                print(f"      🎯 Suggested: {violation.suggested_constant}")
                print(f"      ⚔️ {violation.rambo_comment}")
            
            if len(violations) > 5:
                print(f"      ... and {len(violations) - 5} more kills! 💥")
            print()
        
        print("🎖️ 'We're all gonna die anyway. At least this way we choose how.'")
        print("🔥 RAMBO OUT! 💀")

    def save_analysis(self, output_file: str) -> None:
        """Save analysis results to JSON file"""
        results = {
            "mission": "DELTA-8 RAMBO ASSAULT",
            "target": "/Users/jeromegonzalez/claude-code/fire-salamander/internal/web/web_test.go",
            "total_kills": len(self.violations),
            "rambo_quote": "They drew first blood, not me!",
            "violations": [
                {
                    "line_number": v.line_number,
                    "line_content": v.line_content,
                    "violation_type": v.violation_type,
                    "hardcoded_value": v.hardcoded_value,
                    "suggested_constant": v.suggested_constant,
                    "rambo_comment": v.rambo_comment
                }
                for v in self.violations
            ],
            "constants_mapping": self.constants_mapping,
            "rambo_stats": {
                "violations_by_type": {}
            }
        }
        
        # Calculate stats by type
        for violation in self.violations:
            vtype = violation.violation_type
            if vtype not in results["rambo_stats"]["violations_by_type"]:
                results["rambo_stats"]["violations_by_type"][vtype] = 0
            results["rambo_stats"]["violations_by_type"][vtype] += 1
        
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump(results, f, indent=2, ensure_ascii=False)
        
        print(f"📊 Analysis saved to: {output_file}")

def main():
    """🔥 RAMBO MAIN ASSAULT OPERATION 🔥"""
    print("🔥💀⚔️ DELTA-8 RAMBO DETECTOR ACTIVATED ⚔️💀🔥")
    print("🎯 Mission: Hunt hardcoded web violations")
    print("💀 'Live for nothing, or die for something!'")
    print("=" * 60)
    print()
    
    # Initialize RAMBO
    rambo = Delta8RamboDetector()
    
    # Target file
    target_file = "/Users/jeromegonzalez/claude-code/fire-salamander/internal/web/web_test.go"
    
    # Execute scan
    rambo.scan_target(target_file)
    
    # Print results
    rambo.print_rambo_results()
    
    # Save analysis
    rambo.save_analysis("delta8_rambo_analysis.json")
    
    # Generate constants file content
    constants_content = rambo.generate_constants_file()
    with open("web_test_constants.go", 'w', encoding='utf-8') as f:
        f.write(constants_content)
    
    print()
    print("🎖️ RAMBO MISSION COMPLETE!")
    print("🔥 Constants file generated: web_test_constants.go")
    print("💀 'Until you stop running, you're never really free.'")

if __name__ == "__main__":
    main()