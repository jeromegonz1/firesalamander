#!/usr/bin/env python3
"""
Script de correction spécifique pour checker.go
Remplace toutes les références à l'ancienne structure Config
"""

import re

def fix_checker_config():
    file_path = "./internal/debug/checker.go"
    
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Mappings de remplacement
    replacements = [
        # App references
        ('cfg.App.Name', 'constants.AppName'),
        ('cfg.App.Version', 'constants.AppVersion'),
        ('cfg.App.Icon', 'constants.AppIcon'),
        ('cfg.App.PoweredBy', 'constants.PoweredBy'),
        
        # Server references
        ('cfg.Server.Port', 'cfg.Port'),
        
        # Database references - remplacer par des valeurs appropriées
        ('cfg.Database.Type', '"sqlite"'),
        ('cfg.Database.Path', 'cfg.DBPath'),
        ('cfg.Database.Host', '"localhost"'),
        ('cfg.Database.Name', '"firesalamander"'),
        
        # AI references
        ('cfg.AI.Enabled', '(cfg.OpenAIAPIKey != "")'),
        ('cfg.AI.MockMode', 'false'),
        ('cfg.AI.APIKey', 'cfg.OpenAIAPIKey'),
        
        # Checks simplifiés
        ('if cfg.App.Name == ""', 'if constants.AppName == ""'),
        ('if cfg.Server.Port <= 0', 'if cfg.Port <= 0'),
        ('if cfg.Database.Type == ""', 'if cfg.DBPath == ""'),  # Check du path au lieu du type
        
        # Switch statement sur Database.Type
        ('switch cfg.Database.Type {', 'switch "sqlite" {'),
        
        # Conditions spécifiques
        ('if cfg.Database.Path == ""', 'if cfg.DBPath == ""'),
        ('if cfg.Database.Host == "" || cfg.Database.Name == ""', 'if false'), # toujours false pour sqlite
        
        # Dir extraction fix
        ('dir := cfg.Database.Path[:len(cfg.Database.Path)-len("/firesalamander.db")]', 
         'dir := cfg.DBPath[:len(cfg.DBPath)-len("/firesalamander.db")]'),
        
        # Data mappings
        ('"app_name":      cfg.App.Name,', '"app_name":      constants.AppName,'),
        ('"server_port":   cfg.Server.Port,', '"server_port":   cfg.Port,'),
        ('"database_type": cfg.Database.Type,', '"database_type": "sqlite",'),
        ('"path": cfg.Database.Path', '"path": cfg.DBPath'),
        ('"host": cfg.Database.Host,', '"host": "localhost",'),
        ('"name": cfg.Database.Name,', '"name": "firesalamander",'),
        ('"type": cfg.Database.Type', '"type": "sqlite"'),
        ('"port": cfg.Server.Port', '"port": cfg.Port'),
        
        # AI checks
        ('if !cfg.AI.Enabled', 'if cfg.OpenAIAPIKey == ""'),
        ('if cfg.AI.MockMode', 'if false'),
        ('if cfg.AI.APIKey == ""', 'if cfg.OpenAIAPIKey == ""'),
        
        # Address format
        ('fmt.Sprintf(":%d", cfg.Server.Port)', 'fmt.Sprintf(":%d", cfg.Port)'),
        
        # API Key masking
        ('"api_key":   "***" + cfg.AI.APIKey[len(cfg.AI.APIKey)-4:],', 
         '"api_key":   "***" + cfg.OpenAIAPIKey[len(cfg.OpenAIAPIKey)-4:],'),
    ]
    
    # Appliquer tous les remplacements
    for old, new in replacements:
        content = content.replace(old, new)
    
    # Écrire le fichier modifié
    with open(file_path, 'w') as f:
        f.write(content)
    
    print(f"✅ Fixed all config references in {file_path}")

if __name__ == "__main__":
    fix_checker_config()