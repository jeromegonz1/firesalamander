#!/usr/bin/env python3
"""
ALPHA-3 Smart Hardcode Eliminator
SEO Scorer Hardcoding Violation Eliminator
"""

import os
import re
import shutil
import subprocess
from typing import Dict, List, Tuple

class Alpha3SEOEliminator:
    def __init__(self, file_path: str):
        self.file_path = file_path
        self.backup_path = f"{file_path}.alpha3.backup"
        
        # Mapping complet des 103 violations critiques SEO à éliminer
        self.string_mappings = {
            # SEO Factor Names
            '"title"': 'constants.SEOFactorTitle',
            '"meta_description"': 'constants.SEOFactorMetaDescription',
            '"content_quality"': 'constants.SEOFactorContentQuality',
            '"keyword_optimization"': 'constants.SEOFactorKeywordOptimization',
            '"content_structure"': 'constants.SEOFactorContentStructure',
            '"readability"': 'constants.SEOFactorReadability',
            '"link_optimization"': 'constants.SEOFactorLinkOptimization',
            '"image_optimization"': 'constants.SEOFactorImageOptimization',
            '"ai_enrichment"': 'constants.SEOFactorAIEnrichment',
            
            # Detailed Factor Names
            '"title_optimization"': 'constants.SEOFactorTitleOptimization',
            
            # Status Values
            '"excellent"': 'constants.SEOStatusExcellent',
            '"good"': 'constants.SEOStatusGood',
            '"warning"': 'constants.SEOStatusWarning',
            '"critical"': 'constants.SEOStatusCritical',
            
            # Title Messages
            '"Titre manquant"': 'constants.MsgTitleMissing',
            '"Titre trop court"': 'constants.MsgTitleTooShort',
            '"Titre trop long"': 'constants.MsgTitleTooLong',
            '"Longueur du titre optimale"': 'constants.MsgTitleOptimalLength',
            
            # Meta Description Messages
            '"Meta description manquante"': 'constants.MsgMetaDescMissing',
            '"Meta description trop courte"': 'constants.MsgMetaDescTooShort',
            '"Meta description trop longue"': 'constants.MsgMetaDescTooLong',
            '"Longueur de meta description optimale"': 'constants.MsgMetaDescOptimalLength',
            
            # Content Quality Messages
            '"Contenu trop court"': 'constants.MsgContentTooShort',
            '"Longueur de contenu optimale"': 'constants.MsgContentOptimalLength',
            '"Longueur de contenu correcte"': 'constants.MsgContentCorrectLength',
            
            # Keyword Messages
            '"Aucun mot-clé identifié"': 'constants.MsgNoKeywords',
            '"Peu de mots-clés pertinents"': 'constants.MsgFewKeywords',
            '"Bonne couverture de mots-clés"': 'constants.MsgGoodKeywordCoverage',
            '"Couverture de mots-clés correcte"': 'constants.MsgCorrectKeywordCoverage',
            
            # Structure Messages
            '"Aucun titre de section"': 'constants.MsgNoSectionTitles',
            '"Peu de titres de section"': 'constants.MsgFewSectionTitles',
            '"Structure de contenu correcte"': 'constants.MsgCorrectContentStructure',
            
            # Readability Messages
            '"Lisibilité très faible"': 'constants.MsgVeryLowReadability',
            '"Lisibilité faible"': 'constants.MsgLowReadability',
            '"Lisibilité correcte"': 'constants.MsgCorrectReadability',
            '"Excellente lisibilité"': 'constants.MsgExcellentReadability',
            
            # Link Messages
            '"Aucun lien interne"': 'constants.MsgNoInternalLinks',
            '"Peu de liens internes"': 'constants.MsgFewInternalLinks',
            '"Bon maillage interne"': 'constants.MsgGoodInternalLinks',
            
            # Image Messages
            '"Pas d\'images à optimiser"': 'constants.MsgNoImagesToOptimize',
            '"Aucune image n\'a de texte alternatif"': 'constants.MsgNoImageAltText',
            '"Peu d\'images ont un texte alternatif"': 'constants.MsgFewImageAltText',
            '"La plupart des images ont un texte alternatif"': 'constants.MsgMostImagesHaveAlt',
            '"Toutes les images ont un texte alternatif"': 'constants.MsgAllImagesHaveAlt',
            
            # Suggestions - Title
            '"Ajouter un titre H1 descriptif"': 'constants.SuggAddDescriptiveTitle',
            '"Allonger le titre (30-60 caractères optimal)"': 'constants.SuggExtendTitle',
            '"Raccourcir le titre (risque de troncature)"': 'constants.SuggShortenTitle',
            '"Inclure des mots-clés pertinents dans le titre"': 'constants.SuggIncludeKeywordsTitle',
            '"Éviter la sur-optimisation du titre"': 'constants.SuggAvoidTitleOverOpt',
            
            # Suggestions - Meta Description
            '"Ajouter une meta description attrayante"': 'constants.SuggAddMetaDescription',
            '"Étoffer la meta description (120-160 caractères)"': 'constants.SuggExpandMetaDesc',
            '"Raccourcir la meta description"': 'constants.SuggShortenMetaDesc',
            '"Inclure des mots-clés dans la meta description"': 'constants.SuggIncludeKeywordsMetaDesc',
            '"Ajouter un appel à l\'action dans la meta description"': 'constants.SuggAddCallToAction',
            
            # Suggestions - Content
            '"Étoffer le contenu (minimum 300 mots)"': 'constants.SuggExpandContent',
            '"Améliorer la diversité du vocabulaire"': 'constants.SuggImproveVocabulary',
            '"Éviter le contenu dupliqué"': 'constants.SuggAvoidDuplicateContent',
            '"Ajouter des mots-clés pertinents au contenu"': 'constants.SuggAddRelevantKeywords',
            '"Enrichir le contenu avec plus de mots-clés"': 'constants.SuggEnrichWithKeywords',
            '"Augmenter la densité des mots-clés principaux"': 'constants.SuggIncreaseKeywordDensity',
            '"Éviter la sur-optimisation (densité trop élevée)"': 'constants.SuggAvoidKeywordOverOpt',
            '"Placer des mots-clés dans les positions stratégiques"': 'constants.SuggStrategicKeywordPlace',
            
            # Suggestions - Structure
            '"Structurer le contenu avec des titres H2, H3"': 'constants.SuggStructureWithHeadings',
            '"Améliorer la structure avec plus de sous-titres"': 'constants.SuggImproveWithSubtitles',
            '"Respecter la hiérarchie H1 > H2 > H3"': 'constants.SuggRespectHeadingHierarchy',
            '"Utiliser des listes pour améliorer la lisibilité"': 'constants.SuggUseLists',
            '"Raccourcir les paragraphes pour améliorer la lisibilité"': 'constants.SuggShortenParagraphs',
            '"Développer davantage les paragraphes"': 'constants.SuggDevelopParagraphs',
            
            # Suggestions - Readability
            '"Simplifier les phrases et le vocabulaire"': 'constants.SuggSimplifySentences',
            '"Améliorer la lisibilité du contenu"': 'constants.SuggImproveReadability',
            '"Raccourcir les phrases (max 20-25 mots)"': 'constants.SuggShortenSentences',
            '"Varier la longueur des phrases"': 'constants.SuggVarySentenceLength',
            
            # Suggestions - Links
            '"Ajouter des liens internes vers d\'autres pages"': 'constants.SuggAddInternalLinks',
            '"Augmenter le maillage interne"': 'constants.SuggIncreaseInternalLinks',
            '"Limiter le nombre de liens externes"': 'constants.SuggLimitExternalLinks',
            '"Optimiser les textes d\'ancres des liens"': 'constants.SuggOptimizeAnchorTexts',
            
            # Suggestions - Images
            '"Ajouter des textes alternatifs à toutes les images"': 'constants.SuggAddAltTexts',
            '"Compléter les textes alternatifs manquants"': 'constants.SuggCompleteMissingAlt',
            '"Compléter les derniers textes alternatifs"': 'constants.SuggCompleteLastAltTexts',
            
            # Call-to-Action Terms
            '"découvrir"': 'constants.CTADiscover',
            '"en savoir plus"': 'constants.CTALearnMore',
            '"contacter"': 'constants.CTAContact',
            '"commander"': 'constants.CTAOrder',
            '"acheter"': 'constants.CTABuy',
            '"télécharger"': 'constants.CTADownload',
            '"s\'inscrire"': 'constants.CTASignUp',
            '"essayer"': 'constants.CTATry',
            '"commencer"': 'constants.CTAStart',
            '"cliquer"': 'constants.CTAClick',
            '"discover"': 'constants.CTADiscoverEN',
            '"learn more"': 'constants.CTALearnMoreEN',
            '"contact"': 'constants.CTAContactEN',
            '"order"': 'constants.CTAOrderEN',
            '"buy"': 'constants.CTABuyEN',
            '"download"': 'constants.CTADownloadEN',
            '"sign up"': 'constants.CTASignUpEN',
            '"try"': 'constants.CTATryEN',
            '"start"': 'constants.CTAStartEN',
            '"click"': 'constants.CTAClickEN',
            
            # Bad Anchor Texts
            '"cliquez ici"': 'constants.BadAnchorClickHere',
            '"click here"': 'constants.BadAnchorClickHereEN',
            '"lire la suite"': 'constants.BadAnchorReadMore',
            '"read more"': 'constants.BadAnchorReadMoreEN',
            '"ici"': 'constants.BadAnchorHere',
            '"here"': 'constants.BadAnchorHereEN',
            
            # Heading Types
            '"h1"': 'constants.HeadingH1',
            '"h2"': 'constants.HeadingH2',
            '"h3"': 'constants.HeadingH3',
        }
        
        # Mappings pour les printf
        self.printf_mappings = {
            '"Début scoring SEO - Title:%s WordCount:%d HasAI:%t"': 'constants.LogSEOScoringStart',
            '"Scoring SEO terminé - OverallScore:%.1f FactorsCount:%d IssuesCount:%d RecommendationsCount:%d"': 'constants.LogSEOScoringComplete',
        }
        
        # Mappings pour les valeurs numériques  
        self.numeric_mappings = {
            '0.20': 'constants.WeightTitle',
            '0.15': 'constants.WeightMetaDescription', # first occurrence
            # Note: autres 0.15, 0.10, 0.05 seront gérées contextuellement
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
            
        for pattern in self.printf_mappings.keys():
            violations_before += len(re.findall(re.escape(pattern), content_before))
            violations_after += len(re.findall(re.escape(pattern), content_after))
            
        return violations_before, violations_after

    def eliminate_hardcoding(self) -> bool:
        """Éliminer le hardcoding avec validation"""
        print(f"🚀 Début élimination hardcoding SEO: {self.file_path}")
        
        # Backup
        if not self.create_backup():
            return False
            
        try:
            # Lire le fichier
            with open(self.file_path, 'r', encoding='utf-8') as f:
                content_original = f.read()
            
            content = content_original
            replacements_made = 0
            
            # Appliquer les remplacements strings
            for old_string, new_string in self.string_mappings.items():
                if old_string in content:
                    content = content.replace(old_string, new_string)
                    replacements_made += 1
                    print(f"✅ Remplacé: {old_string} -> {new_string}")
            
            # Appliquer les remplacements printf  
            for old_string, new_string in self.printf_mappings.items():
                if old_string in content:
                    content = content.replace(old_string, new_string)
                    replacements_made += 1
                    print(f"✅ Remplacé: {old_string} -> {new_string}")
            
            # Remplacements contextuels pour les weights
            weight_replacements = [
                ('		"meta_description":     0.15,', f'		{self.string_mappings['"meta_description"']}:     constants.WeightMetaDescription,'),
                ('		"content_quality":      0.15,', f'		{self.string_mappings['"content_quality"']}:      constants.WeightContentQuality,'),
                ('		"keyword_optimization": 0.15,', f'		{self.string_mappings['"keyword_optimization"']}: constants.WeightKeywordOptimization,'),
                ('		"content_structure":    0.10,', f'		{self.string_mappings['"content_structure"']}:    constants.WeightContentStructure,'),
                ('		"readability":          0.10,', f'		{self.string_mappings['"readability"']}:          constants.WeightReadability,'),
                ('		"link_optimization":    0.10,', f'		{self.string_mappings['"link_optimization"']}:    constants.WeightLinkOptimization,'),
                ('		"image_optimization":   0.05,', f'		{self.string_mappings['"image_optimization"']}:   constants.WeightImageOptimization,'),
            ]
            
            for old_line, new_line in weight_replacements:
                if old_line in content:
                    content = content.replace(old_line, new_line)
                    replacements_made += 1
                    print(f"✅ Weight remplacé: {old_line.strip()} -> {new_line.strip()}")
            
            # Compter les violations
            violations_before, violations_after = self.count_violations_before_after(content_original, content)
            reduction = violations_before - violations_after
            percentage = (reduction / violations_before * 100) if violations_before > 0 else 0
            
            print(f"\n📊 RÉSULTATS SEO:")
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
            
            print(f"🎉 SUCCESS SEO: {reduction} violations éliminées ({percentage:.1f}% réduction)")
            return True
            
        except Exception as e:
            print(f"❌ Erreur élimination SEO: {e}")
            self.restore_backup()
            return False

if __name__ == "__main__":
    file_path = "/Users/jeromegonzalez/claude-code/fire-salamander/internal/semantic/seo_scorer.go"
    
    eliminator = Alpha3SEOEliminator(file_path)
    success = eliminator.eliminate_hardcoding()
    
    if success:
        print("\n🏆 ALPHA-3 SEO HARDCODE ELIMINATION: SUCCESS")
    else:
        print("\n💥 ALPHA-3 SEO HARDCODE ELIMINATION: FAILED")