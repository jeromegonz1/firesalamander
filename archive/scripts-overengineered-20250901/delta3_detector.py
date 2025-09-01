#!/usr/bin/env python3
"""
🔍 DELTA-3 DETECTOR
Détection spécialisée pour ngram_analyzer.go - 111 violations
"""

import re
import json
from collections import defaultdict

def analyze_ngram_file(filepath):
    """Analyse le fichier ngram_analyzer.go pour détecter les hardcoding violations"""
    
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    violations = []
    line_num = 0
    
    # Patterns spécialisés pour SEMANTIC/NGRAM Analysis
    patterns = {
        'semantic_field_names': r'"(ngram|frequency|score|weight|relevance|similarity|distance|correlation|analysis|semantic|content|text|token|word|phrase|sentence|paragraph|document|title|description|keywords|tags|category|topic|theme|sentiment|emotion|tone|context|meaning|concept|entity|relation|attribute|feature|vector|embedding|cluster|classification|prediction|confidence|probability|threshold|minimum|maximum|average|median|variance|deviation)"',
        'analysis_types': r'"(unigram|bigram|trigram|4-gram|5-gram|n-gram|tf-idf|cosine|euclidean|jaccard|levenshtein|hamming|jaro|soundex|metaphone|stemming|lemmatization|tokenization|normalization|preprocessing|postprocessing|extraction|enrichment|annotation|tagging|parsing|segmentation|clustering|classification|regression|dimensionality|reduction|optimization)"',
        'nlp_constants': r'"(stopwords|punctuation|whitespace|alphanumeric|unicode|ascii|utf8|encoding|language|locale|grammar|syntax|morphology|phonetics|semantics|pragmatics|discourse|corpus|lexicon|vocabulary|dictionary|thesaurus|ontology|taxonomy|hierarchy|graph|tree|network|matrix|tensor|array|list|set|map|hash|index|cache|buffer|queue|stack|heap)"',
        'metric_names': r'"(precision|recall|f1|f2|accuracy|auc|roc|mae|mse|rmse|r2|correlation|covariance|entropy|mutual_information|kl_divergence|js_divergence|chi_square|p_value|z_score|t_test|anova|regression|classification|clustering|outlier|anomaly|novelty|drift|shift|bias|variance|overfitting|underfitting|generalization|regularization|cross_validation|bootstrap|jackknife)"',
        'algorithm_names': r'"(naive_bayes|svm|knn|decision_tree|random_forest|gradient_boosting|xgboost|lightgbm|neural_network|deep_learning|cnn|rnn|lstm|gru|transformer|bert|gpt|word2vec|glove|fasttext|doc2vec|lda|lsa|nmf|pca|tsne|umap|dbscan|kmeans|hierarchical|spectral|gaussian_mixture|isolation_forest|one_class_svm|local_outlier_factor)"',
        'file_formats': r'"\\.(txt|csv|json|xml|html|pdf|doc|docx|rtf|odt|md|yaml|yml)"',
        'error_messages': r'"[A-Z][^"]*(?:error|failed|invalid|missing|not found|corrupt|malformed|unsupported|timeout|overflow|underflow)[^"]*"',
        'log_messages': r'"[A-Z][^"]*(?:starting|started|processing|completed|analyzing|extracting|computing|calculating|training|testing|validating|evaluating|optimizing|loading|saving|caching|indexing)[^"]*"',
        'config_keys': r'"(model|algorithm|parameters|hyperparameters|features|dimensions|layers|neurons|epochs|batch_size|learning_rate|dropout|regularization|optimizer|loss|metric|threshold|window|stride|padding|activation|initializer|scheduler|early_stopping|checkpointing|logging|debugging|profiling|monitoring)"',
        'data_types': r'"(string|text|numeric|integer|float|double|boolean|categorical|ordinal|nominal|binary|sparse|dense|structured|unstructured|labeled|unlabeled|supervised|unsupervised|semi_supervised|reinforcement|online|offline|batch|stream|real_time|near_real_time)"',
        'preprocessing_steps': r'"(cleaning|normalization|standardization|scaling|encoding|decoding|transformation|augmentation|sampling|filtering|selection|extraction|reduction|compression|decompression|encryption|decryption|hashing|indexing|sorting|shuffling|splitting|merging|joining|aggregation|grouping)"',
    }
    
    for line in content.split('\n'):
        line_num += 1
        line_stripped = line.strip()
        
        # Ignorer les commentaires, imports et struct tags
        if (line_stripped.startswith('//') or 
            line_stripped.startswith('import') or 
            line_stripped.startswith('package') or
            '`json:' in line_stripped):
            continue
            
        for category, pattern in patterns.items():
            matches = re.findall(pattern, line_stripped, re.IGNORECASE)
            for match in matches:
                # Nettoyer le match
                if isinstance(match, tuple):
                    match_value = match[0] if match[0] else match[1] if len(match) > 1 else str(match)
                else:
                    match_value = match
                    
                if len(match_value.strip('"')) < 2:
                    continue
                    
                violations.append({
                    'line': line_num,
                    'category': category,
                    'value': match_value,
                    'context': line_stripped[:150] + ('...' if len(line_stripped) > 150 else '')
                })
    
    return violations

def categorize_violations(violations):
    """Catégorise les violations par type"""
    categories = defaultdict(list)
    for violation in violations:
        categories[violation['category']].append(violation)
    return dict(categories)

def generate_semantic_constants_mapping(violations):
    """Génère les mappings de constantes pour ngram_analyzer.go"""
    
    constants_map = {}
    categorized = categorize_violations(violations)
    
    # Semantic Field Names
    if 'semantic_field_names' in categorized:
        field_map = {
            'ngram': 'constants.SemanticFieldNGram',
            'frequency': 'constants.SemanticFieldFrequency',
            'score': 'constants.SemanticFieldScore',
            'weight': 'constants.SemanticFieldWeight',
            'relevance': 'constants.SemanticFieldRelevance',
            'similarity': 'constants.SemanticFieldSimilarity',
            'distance': 'constants.SemanticFieldDistance',
            'correlation': 'constants.SemanticFieldCorrelation',
            'analysis': 'constants.SemanticFieldAnalysis',
            'semantic': 'constants.SemanticFieldSemantic',
            'content': 'constants.SemanticFieldContent',
            'text': 'constants.SemanticFieldText',
            'token': 'constants.SemanticFieldToken',
            'word': 'constants.SemanticFieldWord',
            'phrase': 'constants.SemanticFieldPhrase',
            'sentence': 'constants.SemanticFieldSentence',
            'paragraph': 'constants.SemanticFieldParagraph',
            'document': 'constants.SemanticFieldDocument',
            'title': 'constants.SemanticFieldTitle',
            'description': 'constants.SemanticFieldDescription',
            'keywords': 'constants.SemanticFieldKeywords',
            'tags': 'constants.SemanticFieldTags',
            'category': 'constants.SemanticFieldCategory',
            'topic': 'constants.SemanticFieldTopic',
            'theme': 'constants.SemanticFieldTheme',
            'sentiment': 'constants.SemanticFieldSentiment',
            'emotion': 'constants.SemanticFieldEmotion',
            'tone': 'constants.SemanticFieldTone',
            'context': 'constants.SemanticFieldContext',
            'meaning': 'constants.SemanticFieldMeaning',
            'concept': 'constants.SemanticFieldConcept',
            'entity': 'constants.SemanticFieldEntity',
            'relation': 'constants.SemanticFieldRelation',
            'attribute': 'constants.SemanticFieldAttribute',
            'feature': 'constants.SemanticFieldFeature',
            'vector': 'constants.SemanticFieldVector',
            'embedding': 'constants.SemanticFieldEmbedding',
            'cluster': 'constants.SemanticFieldCluster',
            'classification': 'constants.SemanticFieldClassification',
            'prediction': 'constants.SemanticFieldPrediction',
            'confidence': 'constants.SemanticFieldConfidence',
            'probability': 'constants.SemanticFieldProbability',
            'threshold': 'constants.SemanticFieldThreshold',
            'minimum': 'constants.SemanticFieldMinimum',
            'maximum': 'constants.SemanticFieldMaximum',
            'average': 'constants.SemanticFieldAverage',
            'median': 'constants.SemanticFieldMedian',
            'variance': 'constants.SemanticFieldVariance',
            'deviation': 'constants.SemanticFieldDeviation'
        }
        for v in categorized['semantic_field_names']:
            if v['value'] in field_map:
                constants_map[f'"{v["value"]}"'] = field_map[v['value']]
    
    # Analysis Types
    if 'analysis_types' in categorized:
        analysis_map = {
            'unigram': 'constants.SemanticAnalysisUnigram',
            'bigram': 'constants.SemanticAnalysisBigram',
            'trigram': 'constants.SemanticAnalysisTrigram',
            '4-gram': 'constants.SemanticAnalysisFourGram',
            '5-gram': 'constants.SemanticAnalysisFiveGram',
            'n-gram': 'constants.SemanticAnalysisNGram',
            'tf-idf': 'constants.SemanticAnalysisTFIDF',
            'cosine': 'constants.SemanticAnalysisCosine',
            'euclidean': 'constants.SemanticAnalysisEuclidean',
            'jaccard': 'constants.SemanticAnalysisJaccard',
            'levenshtein': 'constants.SemanticAnalysisLevenshtein',
            'hamming': 'constants.SemanticAnalysisHamming',
            'jaro': 'constants.SemanticAnalysisJaro',
            'soundex': 'constants.SemanticAnalysisSoundex',
            'metaphone': 'constants.SemanticAnalysisMetaphone',
            'stemming': 'constants.SemanticAnalysisStemming',
            'lemmatization': 'constants.SemanticAnalysisLemmatization',
            'tokenization': 'constants.SemanticAnalysisTokenization',
            'normalization': 'constants.SemanticAnalysisNormalization'
        }
        for v in categorized['analysis_types']:
            if v['value'] in analysis_map:
                constants_map[f'"{v["value"]}"'] = analysis_map[v['value']]
    
    # NLP Constants
    if 'nlp_constants' in categorized:
        nlp_map = {
            'stopwords': 'constants.SemanticNLPStopwords',
            'punctuation': 'constants.SemanticNLPPunctuation',
            'whitespace': 'constants.SemanticNLPWhitespace',
            'alphanumeric': 'constants.SemanticNLPAlphanumeric',
            'unicode': 'constants.SemanticNLPUnicode',
            'ascii': 'constants.SemanticNLPASCII',
            'utf8': 'constants.SemanticNLPUTF8',
            'encoding': 'constants.SemanticNLPEncoding',
            'language': 'constants.SemanticNLPLanguage',
            'locale': 'constants.SemanticNLPLocale',
            'grammar': 'constants.SemanticNLPGrammar',
            'syntax': 'constants.SemanticNLPSyntax',
            'morphology': 'constants.SemanticNLPMorphology',
            'phonetics': 'constants.SemanticNLPPhonetics',
            'semantics': 'constants.SemanticNLPSemantics',
            'pragmatics': 'constants.SemanticNLPPragmatics',
            'discourse': 'constants.SemanticNLPDiscourse',
            'corpus': 'constants.SemanticNLPCorpus',
            'lexicon': 'constants.SemanticNLPLexicon',
            'vocabulary': 'constants.SemanticNLPVocabulary',
            'dictionary': 'constants.SemanticNLPDictionary',
            'thesaurus': 'constants.SemanticNLPThesaurus',
            'ontology': 'constants.SemanticNLPOntology',
            'taxonomy': 'constants.SemanticNLPTaxonomy'
        }
        for v in categorized['nlp_constants']:
            if v['value'] in nlp_map:
                constants_map[f'"{v["value"]}"'] = nlp_map[v['value']]
    
    # Metric Names
    if 'metric_names' in categorized:
        metric_map = {
            'precision': 'constants.SemanticMetricPrecision',
            'recall': 'constants.SemanticMetricRecall',
            'f1': 'constants.SemanticMetricF1',
            'f2': 'constants.SemanticMetricF2',
            'accuracy': 'constants.SemanticMetricAccuracy',
            'auc': 'constants.SemanticMetricAUC',
            'roc': 'constants.SemanticMetricROC',
            'mae': 'constants.SemanticMetricMAE',
            'mse': 'constants.SemanticMetricMSE',
            'rmse': 'constants.SemanticMetricRMSE',
            'r2': 'constants.SemanticMetricR2',
            'correlation': 'constants.SemanticMetricCorrelation',
            'covariance': 'constants.SemanticMetricCovariance'
        }
        for v in categorized['metric_names']:
            if v['value'] in metric_map:
                constants_map[f'"{v["value"]}"'] = metric_map[v['value']]
    
    # Algorithm Names
    if 'algorithm_names' in categorized:
        algo_map = {
            'naive_bayes': 'constants.SemanticAlgorithmNaiveBayes',
            'svm': 'constants.SemanticAlgorithmSVM',
            'knn': 'constants.SemanticAlgorithmKNN',
            'decision_tree': 'constants.SemanticAlgorithmDecisionTree',
            'random_forest': 'constants.SemanticAlgorithmRandomForest',
            'gradient_boosting': 'constants.SemanticAlgorithmGradientBoosting',
            'neural_network': 'constants.SemanticAlgorithmNeuralNetwork',
            'word2vec': 'constants.SemanticAlgorithmWord2Vec',
            'glove': 'constants.SemanticAlgorithmGlove',
            'fasttext': 'constants.SemanticAlgorithmFastText',
            'bert': 'constants.SemanticAlgorithmBERT',
            'gpt': 'constants.SemanticAlgorithmGPT',
            'kmeans': 'constants.SemanticAlgorithmKMeans',
            'dbscan': 'constants.SemanticAlgorithmDBSCAN'
        }
        for v in categorized['algorithm_names']:
            if v['value'] in algo_map:
                constants_map[f'"{v["value"]}"'] = algo_map[v['value']]
    
    return constants_map

def main():
    filepath = '/Users/jeromegonzalez/claude-code/fire-salamander/internal/semantic/ngram_analyzer.go'
    
    print("🔍 DELTA-3 DETECTOR - Scanning ngram_analyzer.go...")
    
    violations = analyze_ngram_file(filepath)
    categorized = categorize_violations(violations)
    constants_map = generate_semantic_constants_mapping(violations)
    
    print(f"\n📊 RÉSULTATS DÉTECTION DELTA-3:")
    print(f"Total violations détectées: {len(violations)}")
    
    for category, viols in categorized.items():
        print(f"\n🔸 {category.upper()}: {len(viols)} violations")
        for v in viols[:3]:  # Show first 3 of each category
            print(f"  Line {v['line']}: {v['value']}")
        if len(viols) > 3:
            print(f"  ... et {len(viols) - 3} autres")
    
    print(f"\n🏗️ CONSTANTES À CRÉER: {len(constants_map)}")
    print("Preview des mappings:")
    for original, constant in list(constants_map.items())[:10]:
        print(f"  {original} → {constant}")
    
    # Sauvegarder les résultats
    results = {
        'total_violations': len(violations),
        'categories': {k: len(v) for k, v in categorized.items()},
        'violations': violations,
        'constants_mapping': constants_map
    }
    
    with open('delta3_analysis.json', 'w') as f:
        json.dump(results, f, indent=2)
    
    print(f"\n✅ Analyse sauvegardée dans delta3_analysis.json")
    return results

if __name__ == "__main__":
    main()