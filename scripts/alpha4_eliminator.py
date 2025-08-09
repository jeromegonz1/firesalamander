#!/usr/bin/env python3
"""
🤖 ALPHA-4 SMART ELIMINATOR
Élimination automatisée des violations hardcoding dans checker.go
Mission: 70 violations → 0 violations
"""

import re
import shutil
import subprocess
import sys

def create_string_mappings():
    """Créer les mappings de remplacement pour ALPHA-4"""
    
    mappings = {
        # Status Values
        '"healthy"': 'constants.HealthStatusHealthy',
        '"degraded"': 'constants.HealthStatusDegraded',
        '"ok"': 'constants.CheckStatusOK',
        '"error"': 'constants.CheckStatusError',
        '"warning"': 'constants.CheckStatusWarning',
        '"passed"': 'constants.TestStatusPassed',
        '"failed"': 'constants.TestStatusFailed',
        
        # Messages
        '"Check failed"': 'constants.MsgCheckFailed',
        '"Check passed"': 'constants.MsgCheckPassed',
        '"Health check completed"': 'constants.MsgHealthCheckCompleted',
        '"Debug endpoint called"': 'constants.MsgDebugEndpointCalled',
        '"Running Phase 1 tests"': 'constants.MsgRunningPhase1Tests',
        '"Creating new health checker"': 'constants.MsgCreatingHealthChecker',
        '"Running all health checks"': 'constants.MsgRunningAllHealthChecks',
        '"Running check"': 'constants.MsgRunningCheck',
        '"Failed to encode debug response"': 'constants.MsgDebugResponseEncodeFailed',
        '"Failed to encode phase tests response"': 'constants.MsgPhaseTestsEncodeFailed',
        '"Internal server error"': 'constants.MsgInternalServerError',
        
        # Configuration Messages
        '"Configuration is nil"': 'constants.MsgConfigIsNil',
        '"Configuration validation failed"': 'constants.MsgConfigValidationFailed',
        '"Configuration is valid"': 'constants.MsgConfigIsValid',
        '"app.name is empty"': 'constants.MsgAppNameEmpty',
        '"server.port is invalid"': 'constants.MsgServerPortInvalid',
        '"database.path is empty"': 'constants.MsgDBPathEmpty',
        
        # Database Messages
        '"SQLite path is empty"': 'constants.MsgSQLitePathEmpty',
        '"SQLite configuration is valid"': 'constants.MsgSQLiteConfigValid',
        '"SQLite directory doesn\'t exist yet"': 'constants.MsgSQLiteDirNotExist',
        '"MySQL configuration incomplete"': 'constants.MsgMySQLConfigIncomplete',
        '"MySQL configuration is valid"': 'constants.MsgMySQLConfigValid',
        '"Unknown database type"': 'constants.MsgUnknownDatabaseType',
        
        # Filesystem Messages
        '"Required files/directories missing"': 'constants.MsgRequiredFilesMissing',
        '"All required files and directories present"': 'constants.MsgAllFilesPresent',
        
        # Network Messages
        '"Network configuration valid"': 'constants.MsgNetworkConfigValid',
        
        # AI Messages
        '"AI is disabled"': 'constants.MsgAIDisabled',
        '"AI is in mock mode"': 'constants.MsgAIMockMode',
        '"AI configuration is valid"': 'constants.MsgAIConfigValid',
        
        # Error Codes
        '"config_nil"': 'constants.ErrorCodeConfigNil',
        '"sqlite_path_missing"': 'constants.ErrorCodeSQLitePathMissing',
        '"mysql_config_incomplete"': 'constants.ErrorCodeMySQLConfigIncomplete',
        '"unknown_db_type"': 'constants.ErrorCodeUnknownDBType',
        
        # Database Types
        '"sqlite"': 'constants.DatabaseTypeSQLite',
        '"mysql"': 'constants.DatabaseTypeMySQL',
        '"firesalamander"': 'constants.DefaultDatabaseName',
        '"localhost"': 'constants.DefaultHost',
        
        # Files and Directories
        '"config"': 'constants.ConfigDir',
        '"deploy"': 'constants.DeployDir',
        '"go.mod"': 'constants.GoModFile',
        '"main.go"': 'constants.MainGoFile',
        '"/firesalamander.db"': 'constants.DefaultDatabaseFile',
        
        # HTTP Constants
        '"Content-Type"': 'constants.HeaderContentType',
        '"application/json"': 'constants.ContentTypeJSON',
        
        # Environments
        '"development"': 'constants.EnvDevelopment',
        '"test"': 'constants.EnvTest',
        
        # Debug Field Names  
        '"check"': 'constants.DebugFieldCheck',
        '"status"': 'constants.DebugFieldStatus',
        '"message"': 'constants.DebugFieldMessage',
        '"error"': 'constants.DebugFieldError',
        '"method"': 'constants.DebugFieldMethod',
        '"path"': 'constants.DebugFieldPath',
        '"remote"': 'constants.DebugFieldRemote',
        '"checks"': 'constants.DebugFieldChecks',
        '"enabled"': 'constants.DebugFieldEnabled',
        '"mock_mode"': 'constants.DebugFieldMockMode',
        '"api_key"': 'constants.DebugFieldAPIKey',
        '"host"': 'constants.DebugFieldHost',
        '"name"': 'constants.DebugFieldName',
        '"type"': 'constants.DebugFieldType',
        '"port"': 'constants.DebugFieldPort',
        '"addr"': 'constants.DebugFieldAddr',
        '"app_name"': 'constants.DebugFieldAppName',
        '"server_port"': 'constants.DebugFieldServerPort',
        '"database_type"': 'constants.DebugFieldDBType',
        
        # Prefixes and Suffixes
        '"directory: "': 'constants.DebugDirPrefix',
        '"file: "': 'constants.DebugFilePrefix',
        '"***"': 'constants.DebugMaskSuffix',
        '"test"': 'constants.DebugQueryParam',
        '"phase1"': 'constants.DebugPhase1Value',
        '"ENV"': 'constants.DebugEnvVariable',
        
        # Numeric Constants
        '1024': 'constants.MemoryDivisor1024',
        '4:': 'constants.APIKeySuffixLength:',
    }
    
    return mappings

def backup_file(filepath):
    """Créer une sauvegarde du fichier original"""
    backup_path = f"{filepath}.alpha4_backup"
    shutil.copy2(filepath, backup_path)
    return backup_path

def restore_file(filepath):
    """Restaurer le fichier depuis la sauvegarde"""
    backup_path = f"{filepath}.alpha4_backup"
    if os.path.exists(backup_path):
        shutil.copy2(backup_path, filepath)
        return True
    return False

def apply_replacements(filepath, mappings):
    """Appliquer les remplacements de hardcoding"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    replacements_made = 0
    
    # Appliquer chaque remplacement
    for old_string, new_string in mappings.items():
        if old_string in content:
            content = content.replace(old_string, new_string)
            replacements_made += 1
            print(f"✅ {old_string} → {new_string}")
    
    # Sauvegarder le fichier modifié
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    
    return replacements_made, len(content) != len(original_content)

def test_compilation():
    """Tester la compilation du code Go"""
    try:
        result = subprocess.run(['go', 'build', './...'], 
                              capture_output=True, text=True, cwd='/Users/jeromegonzalez/claude-code/fire-salamander')
        
        if result.returncode == 0:
            print("✅ Compilation successful!")
            return True
        else:
            print(f"❌ Compilation failed:\n{result.stderr}")
            return False
    except Exception as e:
        print(f"❌ Error testing compilation: {e}")
        return False

def main():
    """Fonction principale ALPHA-4 Eliminator"""
    
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/debug/checker.go'
    
    print("🤖 ALPHA-4 SMART ELIMINATOR")
    print("=" * 50)
    print(f"Target: {filepath}")
    print("Mission: Élimination complète des hardcoding violations")
    
    # 1. Backup
    print("\n📦 Creating backup...")
    backup_path = backup_file(filepath)
    print(f"✅ Backup created: {backup_path}")
    
    # 2. Load mappings
    print("\n🔧 Loading string mappings...")
    mappings = create_string_mappings()
    print(f"✅ {len(mappings)} mappings loaded")
    
    # 3. Apply replacements
    print("\n🔄 Applying hardcode eliminations...")
    replacements_made, file_modified = apply_replacements(filepath, mappings)
    
    print(f"\n📊 RÉSULTATS:")
    print(f"   - Replacements applied: {replacements_made}")
    print(f"   - File modified: {file_modified}")
    
    # 4. Test compilation
    print(f"\n🔨 Testing compilation...")
    if test_compilation():
        print(f"\n🎯 ALPHA-4 MISSION ACCOMPLISHED!")
        print(f"   - Target file: checker.go")
        print(f"   - Violations eliminated: {replacements_made}")
        print(f"   - Status: 100% SUCCESS")
        return True
    else:
        print(f"\n⚠️  Compilation failed - restoring backup...")
        restore_file(filepath)
        print(f"   - File restored from backup")
        print(f"   - Status: NEEDS MANUAL INTERVENTION")
        return False

if __name__ == "__main__":
    import os
    main()