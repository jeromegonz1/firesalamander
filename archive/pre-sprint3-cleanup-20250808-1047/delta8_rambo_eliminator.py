#!/usr/bin/env python3
"""
🔥💀 DELTA-8 RAMBO ELIMINATOR 💀🔥
"Nothing is over! Nothing!"
"They drew first blood, not me!"
"I could have killed them all!"

RAMBO MISSION: Hardcore elimination of web test violations
TARGET: web_test.go
WEAPON: Surgical precision replacements
STATUS: Armed and dangerous

This is WAR against hardcoded strings!
"""

import json
import re
import shutil
import subprocess
import sys
from pathlib import Path
from typing import Dict, List, Tuple, Optional

class RamboEliminator:
    """John Rambo's hardcoded string elimination unit"""
    
    def __init__(self):
        self.target_file = "/Users/jeromegonzalez/claude-code/fire-salamander/internal/web/web_test.go"
        self.constants_file = "/Users/jeromegonzalez/claude-code/fire-salamander/web_test_constants.go"
        self.analysis_file = "/Users/jeromegonzalez/claude-code/fire-salamander/delta8_rambo_analysis.json"
        self.backup_file = self.target_file + ".rambo_backup"
        self.kills = 0
        self.failed_ops = []
        self.rambo_quotes = [
            "Nothing is over! Nothing!",
            "They drew first blood, not me!",
            "I could have killed them all!",
            "Don't push it! Don't push it or I'll give you a war you won't believe!",
            "Live for nothing, or die for something!",
            "When you're pushed, killing's as easy as breathing!",
            "I'm your worst nightmare!"
        ]
        
    def rambo_print(self, message: str, level: str = "INFO"):
        """Print messages with Rambo intensity"""
        if level == "KILL":
            print(f"🎯 RAMBO KILL: {message}")
        elif level == "WARN":
            print(f"⚠️  RAMBO WARNING: {message}")
        elif level == "ERROR":
            print(f"💥 RAMBO FAILURE: {message}")
        elif level == "MISSION":
            print(f"🔥 RAMBO MISSION: {message}")
        else:
            print(f"📡 RAMBO: {message}")
    
    def load_analysis_data(self) -> Dict:
        """Load the analysis data with violation locations"""
        try:
            with open(self.analysis_file, 'r') as f:
                return json.load(f)
        except FileNotFoundError:
            self.rambo_print("Analysis file not found - going in blind!", "ERROR")
            return {"violations": [], "constants_mapping": {}}
    
    def create_backup(self):
        """Create backup before the assault begins"""
        self.rambo_print(f"Creating backup at {self.backup_file}", "MISSION")
        shutil.copy2(self.target_file, self.backup_file)
        self.rambo_print("Backup secured - ready for combat!", "INFO")
    
    def read_target_file(self) -> str:
        """Read the target file content"""
        with open(self.target_file, 'r') as f:
            return f.read()
    
    def write_target_file(self, content: str):
        """Write modified content back to target"""
        with open(self.target_file, 'w') as f:
            f.write(content)
    
    def test_compilation(self) -> bool:
        """Test if Go code compiles after changes"""
        try:
            result = subprocess.run(
                ["go", "build", "-o", "/dev/null", self.target_file],
                capture_output=True,
                text=True,
                cwd=Path(self.target_file).parent
            )
            return result.returncode == 0
        except Exception as e:
            self.rambo_print(f"Compilation test failed: {e}", "ERROR")
            return False
    
    def get_safe_replacements(self, analysis_data: Dict) -> List[Dict]:
        """Get safe replacements that won't break syntax"""
        safe_types = {
            "http_methods": True,
            "status_codes": True,  # Only replace http.StatusOK references
            "http_headers": True,
            "content_types": True,
            "port_numbers": True,  # Only simple port numbers
        }
        
        safe_violations = []
        for violation in analysis_data.get("violations", []):
            v_type = violation.get("violation_type")
            hardcoded_value = violation.get("hardcoded_value")
            
            # Skip unsafe replacements
            if v_type not in safe_types:
                continue
                
            # Special cases for safety
            if v_type == "status_codes" and hardcoded_value != "http.StatusOK":
                continue  # Skip non-StatusOK status codes for safety
                
            if v_type == "port_numbers" and not hardcoded_value.isdigit():
                continue  # Only replace pure numeric ports
                
            safe_violations.append(violation)
        
        return safe_violations
    
    def perform_surgical_strike(self, content: str, violation: Dict) -> Tuple[str, bool]:
        """Perform a surgical replacement on specific violation"""
        line_number = violation.get("line_number")
        hardcoded_value = violation.get("hardcoded_value")
        suggested_constant = violation.get("suggested_constant")
        v_type = violation.get("violation_type")
        
        # Define replacement strategies
        if v_type == "http_methods":
            # Replace "GET" with HTTP_METHOD_GET
            if hardcoded_value == "GET":
                old_pattern = r'"GET"'
                new_replacement = "HTTP_METHOD_GET"
                content = content.replace(old_pattern, new_replacement)
                return content, True
        
        elif v_type == "status_codes":
            # Replace http.StatusOK with HTTP_STATUS_OK
            if hardcoded_value == "http.StatusOK":
                # This is tricky - we need to import http and reference the constant
                old_pattern = r'http\.StatusOK'
                new_replacement = "http.StatusOK"  # Keep as is for now - would need import changes
                # Skip this for safety
                return content, False
        
        elif v_type == "http_headers":
            # Replace "Content-Type" with HEADER_CONTENT_TYPE
            if hardcoded_value in ["Content-Type", "Content-Disposition"]:
                old_pattern = f'"{hardcoded_value}"'
                new_replacement = f"HEADER_{hardcoded_value.upper().replace('-', '_')}"
                if old_pattern in content:
                    content = content.replace(old_pattern, new_replacement)
                    return content, True
        
        elif v_type == "content_types":
            # Replace content types
            if hardcoded_value in ["text/html", "application/json"]:
                old_pattern = f'"{hardcoded_value}"'
                if hardcoded_value == "text/html":
                    new_replacement = "CONTENT_TYPE_HTML"
                elif hardcoded_value == "application/json":
                    new_replacement = "CONTENT_TYPE_JSON"
                
                if old_pattern in content:
                    content = content.replace(old_pattern, new_replacement)
                    return content, True
        
        elif v_type == "port_numbers":
            # Replace port numbers in struct initialization
            if hardcoded_value in ["8080", "8083", "8084"]:
                # Look for Port: 8080 pattern
                old_pattern = f"Port: {hardcoded_value}"
                new_replacement = f"Port: TEST_PORT_{hardcoded_value}"
                if old_pattern in content:
                    content = content.replace(old_pattern, new_replacement)
                    return content, True
                
                # Also look for localhost:port patterns
                old_pattern = f"localhost:{hardcoded_value}"
                new_replacement = f"localhost:\" + strconv.Itoa(TEST_PORT_{hardcoded_value}) + \""
                # This is getting complex, skip for safety
                return content, False
        
        return content, False
    
    def update_constants_import(self, content: str) -> str:
        """Add import for constants if needed"""
        # Check if constants are already imported
        if "web_test_constants" not in content:
            # Add import after existing imports
            import_pattern = r'(import \(\n(?:[^)]*\n)*)'
            replacement = r'\1\t"firesalamander/internal/web/constants"\n'
            content = re.sub(import_pattern, replacement, content)
        
        return content
    
    def execute_rambo_mission(self):
        """Execute the main Rambo elimination mission"""
        self.rambo_print("🔥💀 DELTA-8 RAMBO ELIMINATOR DEPLOYED 💀🔥", "MISSION")
        self.rambo_print("FIRST BLOOD! Engaging hardcoded strings!", "MISSION")
        
        # Load analysis data
        analysis_data = self.load_analysis_data()
        
        # Create backup
        self.create_backup()
        
        # Read target file
        content = self.read_target_file()
        original_content = content
        
        # Get safe replacements only
        safe_violations = self.get_safe_replacements(analysis_data)
        self.rambo_print(f"Identified {len(safe_violations)} safe targets for elimination", "INFO")
        
        # Execute surgical strikes
        for violation in safe_violations:
            hardcoded_value = violation.get("hardcoded_value")
            v_type = violation.get("violation_type")
            
            self.rambo_print(f"Targeting {v_type}: '{hardcoded_value}'", "INFO")
            
            modified_content, success = self.perform_surgical_strike(content, violation)
            
            if success and modified_content != content:
                content = modified_content
                self.kills += 1
                self.rambo_print(f"TARGET ELIMINATED: {hardcoded_value}", "KILL")
            else:
                self.failed_ops.append(f"{v_type}: {hardcoded_value}")
                self.rambo_print(f"Target survived: {hardcoded_value}", "WARN")
        
        # Write back if changes were made
        if content != original_content:
            # Add constants import if constants are used
            if self.kills > 0:
                content = self.add_constants_declarations(content)
            
            self.write_target_file(content)
            self.rambo_print("Modified file written - testing compilation", "INFO")
            
            # Test compilation
            if self.test_compilation():
                self.rambo_print("Compilation successful - mission accomplished!", "MISSION")
            else:
                self.rambo_print("Compilation failed - reverting to backup!", "ERROR")
                # Restore backup
                shutil.copy2(self.backup_file, self.target_file)
                self.rambo_print("Backup restored", "INFO")
                self.kills = 0  # Reset kills since we reverted
        else:
            self.rambo_print("No changes made - target too heavily fortified", "WARN")
        
        # Generate mission report
        self.generate_rambo_report()
    
    def add_constants_declarations(self, content: str) -> str:
        """Add constant declarations to the file"""
        # Add constants at the top of the file after imports
        constants_block = '''
// 🔥💀 RAMBO CONSTANTS - Hardcoded strings eliminated! 💀🔥
const (
	HTTP_METHOD_GET           = "GET"
	HEADER_CONTENT_TYPE      = "Content-Type"
	HEADER_CONTENT_DISPOSITION = "Content-Disposition"
	CONTENT_TYPE_HTML        = "text/html"
	CONTENT_TYPE_JSON        = "application/json"
	TEST_PORT_8080           = 8080
	TEST_PORT_8083           = 8083
	TEST_PORT_8084           = 8084
)

'''
        
        # Find the position after imports
        import_end = content.find(")")
        if import_end != -1:
            # Find the next newline after imports
            next_line = content.find("\n", import_end)
            if next_line != -1:
                content = content[:next_line + 1] + constants_block + content[next_line + 1:]
        
        return content
    
    def generate_rambo_report(self):
        """Generate Rambo-style mission report"""
        self.rambo_print("", "INFO")
        self.rambo_print("=" * 60, "INFO")
        self.rambo_print("🔥💀 DELTA-8 RAMBO ELIMINATION REPORT 💀🔥", "MISSION")
        self.rambo_print("=" * 60, "INFO")
        self.rambo_print("", "INFO")
        
        # Mission quote
        import random
        quote = random.choice(self.rambo_quotes)
        self.rambo_print(f'"{quote}"', "MISSION")
        self.rambo_print("- John Rambo", "INFO")
        self.rambo_print("", "INFO")
        
        # Statistics
        self.rambo_print(f"CONFIRMED KILLS: {self.kills}", "KILL")
        self.rambo_print(f"FAILED OPERATIONS: {len(self.failed_ops)}", "WARN")
        self.rambo_print(f"MISSION SUCCESS RATE: {(self.kills / max(self.kills + len(self.failed_ops), 1)) * 100:.1f}%", "INFO")
        self.rambo_print("", "INFO")
        
        # Detailed results
        if self.kills > 0:
            self.rambo_print("SUCCESSFUL ELIMINATIONS:", "KILL")
            kill_types = {
                "HTTP Methods": ["GET"],
                "HTTP Headers": ["Content-Type", "Content-Disposition"],
                "Content Types": ["text/html", "application/json"],
                "Port Numbers": ["8080", "8083", "8084"]
            }
            
            for category, targets in kill_types.items():
                self.rambo_print(f"  {category}: Eliminated hardcoded strings", "KILL")
        
        if self.failed_ops:
            self.rambo_print("SURVIVORS (Too dangerous to eliminate):", "WARN")
            for survivor in self.failed_ops[:5]:  # Show first 5
                self.rambo_print(f"  - {survivor}", "WARN")
            if len(self.failed_ops) > 5:
                self.rambo_print(f"  ... and {len(self.failed_ops) - 5} more", "WARN")
        
        self.rambo_print("", "INFO")
        self.rambo_print("RAMBO TACTICAL ASSESSMENT:", "MISSION")
        if self.kills > 0:
            self.rambo_print("✅ Surgical strikes successful - constants deployed", "MISSION")
            self.rambo_print("✅ Code compilation maintained - no casualties", "MISSION")
            self.rambo_print("✅ Backup secured - retreat route available", "MISSION")
        else:
            self.rambo_print("🚫 No safe targets eliminated - enemy too fortified", "WARN")
            self.rambo_print("🚫 Complex strings require heavy artillery", "WARN")
        
        self.rambo_print("", "INFO")
        self.rambo_print("RAMBO RECOMMENDATION:", "MISSION")
        self.rambo_print("Manual cleanup required for complex violations", "INFO")
        self.rambo_print("Focus on test descriptions, config values, and error messages", "INFO")
        self.rambo_print("Use constants file: web_test_constants.go", "INFO")
        self.rambo_print("", "INFO")
        
        self.rambo_print("🎯 MISSION STATUS: COMBAT COMPLETE", "MISSION")
        self.rambo_print("They drew first blood, but Rambo drew last!", "MISSION")
        self.rambo_print("=" * 60, "INFO")

def main():
    """Main execution function"""
    eliminator = RamboEliminator()
    eliminator.execute_rambo_mission()

if __name__ == "__main__":
    main()