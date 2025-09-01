#!/usr/bin/env python3
"""
ALPHA-2 Smart Hardcode Eliminator
Phase Tests Hardcoding Violation Eliminator
"""

import os
import re
import shutil
import subprocess
from typing import Dict, List, Tuple

class Alpha2HardcodeEliminator:
    def __init__(self, file_path: str):
        self.file_path = file_path
        self.backup_path = f"{file_path}.alpha2.backup"
        
        # Mapping complet des 74 violations critiques à éliminer
        self.string_mappings = {
            # Phase Numbers
            '"1"': 'constants.Phase1Number',
            
            # Test Names (sans les JSON tags)
            'testLog.Info("🧪 Running Phase 1 Tests - Setup Initial")': 'testLog.Info(constants.LogRunningPhase1Tests)',
            'testLog.Debug("Running configuration test")': 'testLog.Debug(constants.LogConfigurationTest)',
            'testLog.Debug("Running file structure test")': 'testLog.Debug(constants.LogFileStructureTest)',
            'testLog.Debug("Running Git setup test")': 'testLog.Debug(constants.LogGitTest)',
            'testLog.Debug("Running HTTP server test")': 'testLog.Debug(constants.LogHTTPServerTest)',
            'testLog.Debug("Running branding test")': 'testLog.Debug(constants.LogSEPTEOBrandingTest)',
            'testLog.Debug("Running Docker setup test")': 'testLog.Debug(constants.LogDockerTest)',
            'testLog.Debug("Running deploy scripts test")': 'testLog.Debug(constants.LogDeployScriptsTest)',
            'testLog.Info("🧪 Phase 1 Tests completed"': 'testLog.Info(constants.LogPhase1Completed',
            
            # Configuration Keys
            '"status"': 'constants.JSONFieldStatus',
            '"total_tests"': 'constants.JSONFieldTotalTests',
            '"passed_tests"': 'constants.JSONFieldPassedTests',
            '"failed_tests"': 'constants.JSONFieldFailedTests',
            '"duration"': 'constants.JSONFieldDuration',
            
            # Error Messages
            '"config_nil"': 'constants.ErrConfigNil',
            '"app.name incorrect"': 'constants.ErrAppNameIncorrect',
            '"server.port invalid"': 'constants.ErrServerPortInvalid',
            
            # Detail Keys
            '"app_name"': 'constants.DetailAppName',
            '"app_icon"': 'constants.DetailAppIcon',
            '"powered_by"': 'constants.DetailPoweredBy',
            '"primary_color"': 'constants.DetailPrimaryColor',
            '"server_port"': 'constants.DetailServerPort',
            '"issues_found"': 'constants.DetailIssuesFound',
            '"issues"': 'constants.DetailIssues',
            '"required_files"': 'constants.DetailRequiredFiles',
            '"required_dirs"': 'constants.DetailRequiredDirs',
            '"missing_files"': 'constants.DetailMissingFiles',
            '"missing_dirs"': 'constants.DetailMissingDirs',
            '"missing_files_list"': 'constants.DetailMissingFilesList',
            '"missing_dirs_list"': 'constants.DetailMissingDirsList',
            '"git_initialized"': 'constants.DetailGitInitialized',
            '"gitignore_exists"': 'constants.DetailGitignoreExists',
            '"base_url"': 'constants.DetailBaseURL',
            '"total_endpoints"': 'constants.DetailTotalEndpoints',
            '"passed_endpoints"': 'constants.DetailPassedEndpoints',
            '"url"': 'constants.DetailURL',
            '"expected_content_type"': 'constants.DetailExpectedContentType',
            '"accessible"': 'constants.DetailAccessible',
            '"status_code"': 'constants.DetailStatusCode',
            '"content_type"': 'constants.DetailContentType',
            '"septeo_logo_integrated"': 'constants.DetailSepteoLogoIntegrated',
            '"septeo_orange_integrated"': 'constants.DetailSepteoOrangeIntegrated',
            '"salamander_icon_integrated"': 'constants.DetailSalamanderIconIntegrated',
            '"has_app_service"': 'constants.DetailHasAppService',
            '"has_db_service"': 'constants.DetailHasDbService',
            '"has_port_mapping"': 'constants.DetailHasPortMapping',
            '"total_scripts"': 'constants.DetailTotalScripts',
            '"missing_scripts"': 'constants.DetailMissingScripts',
            '"non_executable_scripts"': 'constants.DetailNonExecutableScripts',
            '"missing_scripts_list"': 'constants.DetailMissingScriptsList',
            '"non_executable_scripts_list"': 'constants.DetailNonExecutableScriptsList',
            
            # File Paths
            '"go.mod"': 'constants.RequiredFileGoMod',
            '"main.go"': 'constants.RequiredFileMainGo',
            '"README.md"': 'constants.RequiredFileReadme',
            '".gitignore"': 'constants.RequiredFileGitignore',
            '"docker-compose.yml"': 'constants.RequiredFileDockerCompose',
            '"config"': 'constants.RequiredDirConfig',
            '"deploy"': 'constants.RequiredDirDeploy',
            '"internal"': 'constants.RequiredDirInternal',
            '"internal/logger"': 'constants.RequiredDirInternalLogger',
            '"internal/debug"': 'constants.RequiredDirInternalDebug',
            '".git"': 'constants.GitDirectory',
            '"deploy/deploy.sh"': 'constants.DeployScriptDeploy',
            '"deploy/setup-infomaniak.sh"': 'constants.DeployScriptSetupInfomaniak',
            
            # Docker Services
            '"app:"': 'constants.DockerServiceApp',
            '"db:"': 'constants.DockerServiceDB',
            
            # Branding Constants
            '"septeo.svg"': 'constants.SepteoLogoPath',
            '"#ff6136"': 'constants.SepteoOrangeColor',
            '"🦎"': 'constants.SalamanderIcon',
            
            # HTTP Endpoints
            '"/health"': 'constants.EndpointHealth',
            '"/debug"': 'constants.EndpointDebug',
            '"/"': 'constants.EndpointHome',
            '"application/json"': 'constants.ContentTypeJSON',
            '"text/html"': 'constants.ContentTypeHTML',
            
            # Test Endpoint Names
            '"Health Endpoint"': 'constants.TestEndpointHealth',
            '"Debug Endpoint"': 'constants.TestEndpointDebug',
            '"Home Page"': 'constants.TestEndpointHome',
            
            # Messages
            '"Cannot test server - configuration is nil"': 'constants.MsgCannotTestServer',
            '"Endpoint not accessible"': 'constants.MsgEndpointNotAccessible',
            '"SEPTEO logo URL not found in main.go"': 'constants.MsgSepteoLogoNotFound',
            '"SEPTEO orange color not found"': 'constants.MsgSepteoOrangeNotFound',
            '"Fire Salamander icon not found"': 'constants.MsgSalamanderIconNotFound',
            '"Cannot read main.go file"': 'constants.MsgCannotReadMainGo',
            '"docker-compose.yml file missing"': 'constants.MsgDockerComposeFileMissing',
            '"docker_compose_missing"': 'constants.ErrDockerComposeMissing',
            '"Docker Compose configuration is incomplete"': 'constants.ErrDockerComposeInvalid',
            '"Cannot read docker-compose.yml"': 'constants.ErrCannotReadDockerCompose',
            '".git directory missing"': 'constants.MsgGitDirectoryMissing',
            '".gitignore missing"': 'constants.MsgGitignoreMissing',
        }
        
        # Mappings pour les messages avec sprintf
        self.sprintf_mappings = {
            '"Configuration validation failed: %s"': 'constants.MsgConfigurationFailed',
            '"Missing %d files and %d directories"': 'constants.MsgMissingFilesAndDirs',
            '"Git setup issues: %s"': 'constants.MsgGitSetupIssues',
            '"Only %d/%d endpoints responding correctly"': 'constants.MsgOnlyEndpointsResponding',
            '"Wrong content type: got %s, expected %s"': 'constants.MsgWrongContentType',
            '"Unexpected status code: %d"': 'constants.MsgUnexpectedStatusCode',
            '"Branding issues: %s"': 'constants.MsgBrandingIssues',
            '"Deploy script issues: %s"': 'constants.MsgDeployScriptIssues',
            '"%d missing scripts"': 'constants.MsgMissingScriptsCount',
            '"%d non-executable scripts"': 'constants.MsgNonExecutableScriptsCount',
        }

    def create_backup(self) -> bool:
        """Créer une sauvegarde du fichier original"""
        try:
            shutil.copy2(self.file_path, self.backup_path)
            print(f"✅ Backup créé: {self.backup_path}")
            return True
        except Exception as e:
            print(f"❌ Erreur backup: {e}")
            return False

    def restore_backup(self) -> bool:
        """Restaurer depuis la sauvegarde"""
        try:
            shutil.copy2(self.backup_path, self.file_path)
            print(f"🔄 Fichier restauré depuis backup")
            return True
        except Exception as e:
            print(f"❌ Erreur restore: {e}")
            return False

    def test_compilation(self) -> bool:
        """Tester la compilation Go"""
        try:
            result = subprocess.run(['go', 'build', './...'], 
                                  cwd=os.path.dirname(self.file_path), 
                                  capture_output=True, text=True)
            if result.returncode == 0:
                print("✅ Compilation réussie")
                return True
            else:
                print(f"❌ Erreur compilation: {result.stderr}")
                return False
        except Exception as e:
            print(f"❌ Erreur test compilation: {e}")
            return False

    def count_violations_before_after(self, content_before: str, content_after: str) -> Tuple[int, int]:
        """Compter les violations avant et après"""
        violations_before = 0
        violations_after = 0
        
        for pattern in self.string_mappings.keys():
            violations_before += len(re.findall(re.escape(pattern), content_before))
            violations_after += len(re.findall(re.escape(pattern), content_after))
            
        for pattern in self.sprintf_mappings.keys():
            violations_before += len(re.findall(re.escape(pattern), content_before))
            violations_after += len(re.findall(re.escape(pattern), content_after))
            
        return violations_before, violations_after

    def eliminate_hardcoding(self) -> bool:
        """Éliminer le hardcoding avec validation"""
        print(f"🚀 Début élimination hardcoding: {self.file_path}")
        
        # Backup
        if not self.create_backup():
            return False
            
        try:
            # Lire le fichier
            with open(self.file_path, 'r', encoding='utf-8') as f:
                content_original = f.read()
            
            content = content_original
            replacements_made = 0
            
            # Appliquer les remplacements simples
            for old_string, new_string in self.string_mappings.items():
                if old_string in content:
                    content = content.replace(old_string, new_string)
                    replacements_made += 1
                    print(f"✅ Remplacé: {old_string} -> {new_string}")
            
            # Appliquer les remplacements sprintf
            for old_string, new_string in self.sprintf_mappings.items():
                if old_string in content:
                    content = content.replace(old_string, new_string)
                    replacements_made += 1
                    print(f"✅ Remplacé: {old_string} -> {new_string}")
            
            # Compter les violations
            violations_before, violations_after = self.count_violations_before_after(content_original, content)
            reduction = violations_before - violations_after
            percentage = (reduction / violations_before * 100) if violations_before > 0 else 0
            
            print(f"\n📊 RÉSULTATS:")
            print(f"🔢 Violations avant: {violations_before}")
            print(f"🔢 Violations après: {violations_after}")
            print(f"📉 Réduction: {reduction} ({percentage:.1f}%)")
            print(f"🔄 Remplacements: {replacements_made}")
            
            # Écrire le fichier modifié
            with open(self.file_path, 'w', encoding='utf-8') as f:
                f.write(content)
            
            # Test de compilation
            if not self.test_compilation():
                print("❌ Échec compilation, restoration backup...")
                self.restore_backup()
                return False
            
            print(f"🎉 SUCCESS: {reduction} violations éliminées ({percentage:.1f}% réduction)")
            return True
            
        except Exception as e:
            print(f"❌ Erreur élimination: {e}")
            self.restore_backup()
            return False

if __name__ == "__main__":
    file_path = "/Users/jeromegonzalez/claude-code/fire-salamander/internal/debug/phase_tests.go"
    
    eliminator = Alpha2HardcodeEliminator(file_path)
    success = eliminator.eliminate_hardcoding()
    
    if success:
        print("\n🏆 ALPHA-2 HARDCODE ELIMINATION: SUCCESS")
    else:
        print("\n💥 ALPHA-2 HARDCODE ELIMINATION: FAILED")