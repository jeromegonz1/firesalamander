#!/usr/bin/env python3
"""
🤖 BRAVO-2 SMART ELIMINATOR
Élimination automatisée des violations hardcoding dans tag_analyzer.go
Mission: 251 violations → 0 violations
"""

import re
import shutil
import subprocess
import sys

def create_string_mappings():
    """Créer les mappings de remplacement pour BRAVO-2"""
    
    mappings = {
        # JSON Field Names - Analysis Structure
        '"title"': 'constants.TagJSONFieldTitle',
        '"meta_description"': 'constants.TagJSONFieldMetaDescription',
        '"headings"': 'constants.TagJSONFieldHeadings',
        '"meta_tags"': 'constants.TagJSONFieldMetaTags',
        '"images"': 'constants.TagJSONFieldImages',
        '"links"': 'constants.TagJSONFieldLinks',
        '"microdata"': 'constants.TagJSONFieldMicrodata',
        
        # JSON Field Names - Common Fields
        '"present"': 'constants.TagJSONFieldPresent',
        '"content"': 'constants.TagJSONFieldContent',
        '"length"': 'constants.TagJSONFieldLength',
        '"optimal_length"': 'constants.TagJSONFieldOptimalLength',
        '"issues"': 'constants.TagJSONFieldIssues',
        '"recommendations"': 'constants.TagJSONFieldRecommendations',
        
        # JSON Field Names - Specific Analysis
        '"has_keywords"': 'constants.TagJSONFieldHasKeywords',
        '"duplicate_words"': 'constants.TagJSONFieldDuplicateWords',
        '"has_call_to_action"': 'constants.TagJSONFieldHasCallToAction',
        '"h1_count"': 'constants.TagJSONFieldH1Count',
        '"h1_content"': 'constants.TagJSONFieldH1Content',
        '"heading_structure"': 'constants.TagJSONFieldHeadingStructure',
        '"has_hierarchy"': 'constants.TagJSONFieldHasHierarchy',
        '"missing_levels"': 'constants.TagJSONFieldMissingLevels',
        '"has_robots"': 'constants.TagJSONFieldHasRobots',
        '"robots_content"': 'constants.TagJSONFieldRobotsContent',
        '"has_canonical"': 'constants.TagJSONFieldHasCanonical',
        '"canonical_url"': 'constants.TagJSONFieldCanonicalURL',
        '"has_og_tags"': 'constants.TagJSONFieldHasOGTags',
        '"og_tags"': 'constants.TagJSONFieldOGTags',
        '"has_twitter_card"': 'constants.TagJSONFieldHasTwitterCard',
        '"twitter_card"': 'constants.TagJSONFieldTwitterCard',
        '"has_viewport"': 'constants.TagJSONFieldHasViewport',
        '"viewport_content"': 'constants.TagJSONFieldViewportContent',
        '"total_images"': 'constants.TagJSONFieldTotalImages',
        '"images_with_alt"': 'constants.TagJSONFieldImagesWithAlt',
        '"alt_text_coverage"': 'constants.TagJSONFieldAltTextCoverage',
        '"missing_alt_images"': 'constants.TagJSONFieldMissingAltImages',
        '"optimized_formats"': 'constants.TagJSONFieldOptimizedFormats',
        '"lazy_loading"': 'constants.TagJSONFieldLazyLoading',
        '"internal_links"': 'constants.TagJSONFieldInternalLinks',
        '"external_links"': 'constants.TagJSONFieldExternalLinks',
        '"nofollow_links"': 'constants.TagJSONFieldNoFollowLinks',
        '"broken_links"': 'constants.TagJSONFieldBrokenLinks',
        '"anchor_optimization"': 'constants.TagJSONFieldAnchorOptimization',
        '"has_json_ld"': 'constants.TagJSONFieldHasJSONLD',
        '"json_ld_types"': 'constants.TagJSONFieldJSONLDTypes',
        '"has_microdata"': 'constants.TagJSONFieldHasMicrodata',
        '"microdata_types"': 'constants.TagJSONFieldMicrodataTypes',
        '"has_rdfa"': 'constants.TagJSONFieldHasRDFa',
        '"property"': 'constants.TagJSONFieldProperty',
        '"name"': 'constants.TagJSONFieldName',
        
        # Meta Names
        '"description"': 'constants.TagMetaNameDescription',
        '"robots"': 'constants.TagMetaNameRobots',
        '"viewport"': 'constants.TagMetaNameViewport',
        
        # HTML Attributes
        '"src"': 'constants.TagAttrSrc',
        '"alt"': 'constants.TagAttrAlt',
        '"loading"': 'constants.TagAttrLoading',
        '"href"': 'constants.TagAttrHref',
        '"rel"': 'constants.TagAttrRel',
        '"type"': 'constants.TagAttrType',
        
        # HTML Values
        '"lazy"': 'constants.TagValueLazy',
        '"canonical"': 'constants.TagValueCanonical',
        '"nofollow"': 'constants.TagValueNoFollow',
        '"application/ld+json"': 'constants.TagValueJSONLD',
        '"og:"': 'constants.TagPrefixOG',
        '"twitter:"': 'constants.TagPrefixTwitter',
        '"itemscope"': 'constants.TagValueItemScope',
        '"typeof"': 'constants.TagValueTypeOf',
        '"vocab"': 'constants.TagValueVocab',
        
        # Heading Levels
        '"h1"': 'constants.TagHeadingH1',
        '"h2"': 'constants.TagHeadingH2',
        '"h3"': 'constants.TagHeadingH3',
        '"h4"': 'constants.TagHeadingH4',
        '"h5"': 'constants.TagHeadingH5',
        '"h6"': 'constants.TagHeadingH6',
        
        # Error Messages
        '"Balise title manquante"': 'constants.TagMsgTitleMissing',
        '"Titre trop court"': 'constants.TagMsgTitleTooShort',
        '"Titre trop long"': 'constants.TagMsgTitleTooLong',
        '"Mots dupliqués dans le titre"': 'constants.TagMsgTitleDuplicates',
        '"Meta description manquante"': 'constants.TagMsgMetaDescMissing',
        '"Meta description trop courte"': 'constants.TagMsgMetaDescTooShort',
        '"Meta description trop longue"': 'constants.TagMsgMetaDescTooLong',
        '"Aucun titre H1"': 'constants.TagMsgNoH1',
        '"Plusieurs titres H1"': 'constants.TagMsgMultipleH1',
        '"Hiérarchie des titres incorrecte"': 'constants.TagMsgBadHierarchy',
        '"Meta viewport manquante"': 'constants.TagMsgViewportMissing',
        '"Textes d\'ancre peu optimisés"': 'constants.TagMsgBadAnchors',
        
        # Recommendations
        '"Ajouter une balise <title> descriptive"': 'constants.TagRecommendAddTitle',
        '"Étendre le titre (30-60 caractères optimal)"': 'constants.TagRecommendExtendTitle',
        '"Raccourcir le titre (risque de troncature)"': 'constants.TagRecommendShortenTitle',
        '"Éviter la répétition de mots dans le titre"': 'constants.TagRecommendAvoidDuplicates',
        '"Ajouter une meta description attrayante"': 'constants.TagRecommendAddMetaDesc',
        '"Étendre la meta description (120-160 caractères)"': 'constants.TagRecommendExtendMetaDesc',
        '"Raccourcir la meta description"': 'constants.TagRecommendShortenMetaDesc',
        '"Ajouter un appel à l\'action dans la meta description"': 'constants.TagRecommendAddCTA',
        '"Ajouter un titre H1 principal"': 'constants.TagRecommendAddH1',
        '"Utiliser un seul H1 par page"': 'constants.TagRecommendSingleH1',
        '"Respecter la hiérarchie H1 > H2 > H3..."': 'constants.TagRecommendRespectHierarchy',
        '"Ajouter une URL canonique"': 'constants.TagRecommendAddCanonical',
        '"Ajouter meta viewport pour le responsive"': 'constants.TagRecommendAddViewport',
        '"Ajouter les balises Open Graph"': 'constants.TagRecommendAddOGTags',
        '"Ajouter les balises Twitter Card"': 'constants.TagRecommendAddTwitter',
        '"Ajouter des textes alternatifs à toutes les images"': 'constants.TagRecommendAddAltText',
        '"Implémenter le lazy loading pour les images"': 'constants.TagRecommendLazyLoading',
        '"Améliorer le maillage interne"': 'constants.TagRecommendImproveInternal',
        '"Optimiser les textes d\'ancres des liens"': 'constants.TagRecommendOptimizeAnchors',
        '"Ajouter des données structurées (JSON-LD recommandé)"': 'constants.TagRecommendAddStructuredData',
        
        # CTA Words - French
        '"découvrir"': 'constants.TagCTADecouvrir',
        '"en savoir plus"': 'constants.TagCTASavoirPlus',
        '"contacter"': 'constants.TagCTAContacter',
        '"commander"': 'constants.TagCTACommander',
        '"acheter"': 'constants.TagCTAAcheter',
        '"télécharger"': 'constants.TagCTATelecharger',
        '"s\'inscrire"': 'constants.TagCTAInscrire',
        '"essayer"': 'constants.TagCTAEssayer',
        '"commencer"': 'constants.TagCTACommencer',
        '"cliquer"': 'constants.TagCTACliquer',
        
        # CTA Words - English
        '"discover"': 'constants.TagCTADiscover',
        '"learn more"': 'constants.TagCTALearnMore',
        '"contact"': 'constants.TagCTAContact',
        '"order"': 'constants.TagCTAOrder',
        '"buy"': 'constants.TagCTABuy',
        '"download"': 'constants.TagCTADownload',
        '"sign up"': 'constants.TagCTASignUp',
        '"try"': 'constants.TagCTATry',
        '"start"': 'constants.TagCTAStart',
        '"click"': 'constants.TagCTAClick',
        
        # Bad Anchors
        '"cliquez ici"': 'constants.TagBadAnchorCliquezIci',
        '"click here"': 'constants.TagBadAnchorClickHere',
        '"lire la suite"': 'constants.TagBadAnchorLireSuite',
        '"read more"': 'constants.TagBadAnchorReadMore',
        '"ici"': 'constants.TagBadAnchorIci',
        '"here"': 'constants.TagBadAnchorHere',
        
        # Remove regex patterns - they need special handling
        
        # Detected value
        '"detected"': 'constants.TagDetectedValue',
        
        # Numeric thresholds
        '30': 'constants.TagTitleMinLength',
        '60': 'constants.TagTitleMaxLength',
        '120': 'constants.TagMetaDescMinLength',
        '160': 'constants.TagMetaDescMaxLength',
        '200': 'constants.TagTitleRegexMax',
        '300': 'constants.TagMetaDescRegexMax',
        '3': 'constants.TagMinWordLength',
        '0.7': 'constants.TagMinAnchorOptim',
        '1.0': 'constants.TagFullCoverage',
    }
    
    return mappings

def backup_file(filepath):
    """Créer une sauvegarde du fichier original"""
    backup_path = f"{filepath}.bravo2_backup"
    shutil.copy2(filepath, backup_path)
    return backup_path

def restore_file(filepath):
    """Restaurer le fichier depuis la sauvegarde"""
    backup_path = f"{filepath}.bravo2_backup"
    try:
        shutil.copy2(backup_path, filepath)
        return True
    except:
        return False

def apply_targeted_replacements(filepath, mappings):
    """Appliquer les remplacements de hardcoding de manière ciblée"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    replacements_made = 0
    
    # Remplacements spéciaux avec contexte pour éviter les conflits
    special_replacements = [
        # String arrays for CTA words
        (r'ctas := \[\]string\{[^}]+\}', lambda m: replace_cta_array(m.group(0))),
        
        # String arrays for bad anchors
        (r'badAnchors := \[\]string\{[^}]+\}', lambda m: replace_bad_anchors_array(m.group(0))),
        
        # Heading levels array
        (r'levelNames := \[\]string\{[^}]+\}', lambda m: replace_level_names_array(m.group(0))),
        
        # Numeric threshold comparisons context-sensitive
        (r'analysis\.Length >= 30 && analysis\.Length <= 60', 
         'analysis.Length >= constants.TagTitleMinLength && analysis.Length <= constants.TagTitleMaxLength'),
        (r'analysis\.Length < 30', 'analysis.Length < constants.TagTitleMinLength'),
        (r'analysis\.Length >= 120 && analysis\.Length <= 160',
         'analysis.Length >= constants.TagMetaDescMinLength && analysis.Length <= constants.TagMetaDescMaxLength'),
        (r'analysis\.Length < 120', 'analysis.Length < constants.TagMetaDescMinLength'),
        
        # String formatting with images count
        (r'fmt\.Sprintf\("%d images sans texte alternatif", len\(analysis\.MissingAltImages\)\)',
         'fmt.Sprintf(constants.TagMsgImagesNoAlt, len(analysis.MissingAltImages))'),
    ]
    
    def replace_cta_array(match_text):
        """Replace CTA array with constants"""
        # Build array with constants
        cta_constants = [
            'constants.TagCTADecouvrir', 'constants.TagCTASavoirPlus', 'constants.TagCTAContacter',
            'constants.TagCTACommander', 'constants.TagCTAAcheter', 'constants.TagCTATelecharger',
            'constants.TagCTAInscrire', 'constants.TagCTAEssayer', 'constants.TagCTACommencer', 'constants.TagCTACliquer',
            'constants.TagCTADiscover', 'constants.TagCTALearnMore', 'constants.TagCTAContact',
            'constants.TagCTAOrder', 'constants.TagCTABuy', 'constants.TagCTADownload',
            'constants.TagCTASignUp', 'constants.TagCTATry', 'constants.TagCTAStart', 'constants.TagCTAClick'
        ]
        return f'ctas := []string{{\n\t\t{", ".join(cta_constants)},\n\t}}'
    
    def replace_bad_anchors_array(match_text):
        """Replace bad anchors array with constants"""
        bad_anchor_constants = [
            'constants.TagBadAnchorCliquezIci', 'constants.TagBadAnchorClickHere',
            'constants.TagBadAnchorLireSuite', 'constants.TagBadAnchorReadMore',
            'constants.TagBadAnchorIci', 'constants.TagBadAnchorHere'
        ]
        return f'badAnchors := []string{{{", ".join(bad_anchor_constants)}}}'
    
    def replace_level_names_array(match_text):
        """Replace heading level names array with constants"""
        level_constants = [
            'constants.TagHeadingH1', 'constants.TagHeadingH2', 'constants.TagHeadingH3',
            'constants.TagHeadingH4', 'constants.TagHeadingH5', 'constants.TagHeadingH6'
        ]
        return f'levelNames := []string{{{", ".join(level_constants)}}}'
    
    # Appliquer les remplacements spéciaux d'abord
    for pattern, replacement in special_replacements:
        if isinstance(replacement, str):
            if re.search(pattern, content):
                content = re.sub(pattern, replacement, content)
                replacements_made += 1
                print(f"✅ Special: {pattern[:50]}... → {replacement[:50]}...")
        else:
            # C'est une fonction lambda
            matches = re.findall(pattern, content, re.DOTALL)
            if matches:
                content = re.sub(pattern, replacement, content, flags=re.DOTALL)
                replacements_made += len(matches)
                print(f"✅ Special function: {pattern[:50]}... → {len(matches)} replacements")
    
    # Appliquer les remplacements standards
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
        result = subprocess.run(['go', 'build', './internal/seo'], 
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
    """Fonction principale BRAVO-2 Eliminator"""
    
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/seo/tag_analyzer.go'
    
    print("🤖 BRAVO-2 SMART ELIMINATOR")
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
    replacements_made, file_modified = apply_targeted_replacements(filepath, mappings)
    
    print(f"\n📊 RÉSULTATS:")
    print(f"   - Replacements applied: {replacements_made}")
    print(f"   - File modified: {file_modified}")
    
    # 4. Test compilation
    print(f"\n🔨 Testing compilation...")
    if test_compilation():
        print(f"\n🎯 BRAVO-2 MISSION ACCOMPLISHED!")
        print(f"   - Target file: tag_analyzer.go")
        print(f"   - Violations eliminated: {replacements_made}")
        print(f"   - Status: 100% SUCCESS")
        return True
    else:
        print(f"\n⚠️  Compilation failed - restoring backup...")
        if restore_file(filepath):
            print(f"   - File restored from backup")
        print(f"   - Status: NEEDS MANUAL INTERVENTION")
        return False

if __name__ == "__main__":
    import os
    main()