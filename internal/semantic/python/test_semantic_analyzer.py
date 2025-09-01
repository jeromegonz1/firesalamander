import pytest
import json
import yaml
from pathlib import Path
from unittest.mock import Mock, patch
import numpy as np

# Import our modules (to be created)
from semantic_analyzer import SemanticAnalyzer
from ngram_analyzer import NgramAnalyzer
from topic_modeler import TopicModeler
from keyword_ranker import KeywordRanker


class TestSemanticAnalyzer:
    """Test suite for semantic analyzer following TDD approach"""
    
    @pytest.fixture
    def config(self):
        """Load test configuration"""
        return {
            'semantic': {
                'language': {
                    'target': 'fr',
                    'confidence_threshold': 0.8
                },
                'preprocessing': {
                    'min_ngram_length': 2,
                    'max_ngram_length': 5,
                    'min_word_count': 2,
                    'stopwords_file': 'config/stopwords_fr.txt'
                },
                'ranker': {
                    'weights': {
                        'thematic_score': 0.3,
                        'intent_score': 0.25,
                        'mesh_evidence': 0.2,
                        'readability': 0.15,
                        'historical_feedback': 0.1
                    }
                },
                'clustering': {
                    'algorithm': 'HDBSCAN',
                    'min_cluster_size': 3,
                    'min_samples': 2
                },
                'models': {
                    'embeddings': 'distilcambert-base',
                    'intent_classifier': 'logistic_regression'
                },
                'output': {
                    'max_candidates': 500,
                    'top_suggestions': 20,
                    'diversity_threshold': 0.8
                }
            },
            'filters': {
                'banlist': ['cliquez ici', 'en savoir plus'],
                'brand_terms': ['septeo', 'fire salamander'],
                'generic_terms': ['accueil', 'contact']
            }
        }
    
    @pytest.fixture
    def sample_crawl_data(self):
        """Sample crawl data for testing"""
        return {
            "pages": [
                {
                    "url": "https://example.com/",
                    "lang": "fr",
                    "title": "Logiciel de gestion pour cabinets d'avocats",
                    "h1": "Solutions numériques pour professionnels du droit",
                    "h2": ["Gestion de dossiers", "Facturation clients", "Agenda partagé"],
                    "content": "Notre logiciel de gestion permet aux cabinets d'avocats de gérer efficacement leurs dossiers clients, la facturation et les échéances. Solution cloud sécurisée pour avocats et juristes."
                },
                {
                    "url": "https://example.com/fonctionnalites",
                    "lang": "fr",
                    "title": "Fonctionnalités - Logiciel juridique",
                    "h1": "Fonctionnalités complètes pour votre cabinet",
                    "h2": ["Gestion documentaire", "Comptabilité trust", "Signature électronique"],
                    "content": "Découvrez toutes les fonctionnalités de notre solution : gestion électronique des documents, comptabilité spécialisée, signature numérique, et bien plus."
                }
            ]
        }
    
    def test_analyzer_initialization(self, config):
        """Test that analyzer initializes correctly with config"""
        analyzer = SemanticAnalyzer(config)
        assert analyzer.config == config
        assert analyzer.language == 'fr'
        assert analyzer.model_version == 'sem-v1.0'
    
    def test_analyze_crawl_data(self, config, sample_crawl_data):
        """Test complete analysis pipeline"""
        analyzer = SemanticAnalyzer(config)
        result = analyzer.analyze("test_audit_123", sample_crawl_data)
        
        # Validate result structure matches JSON schema
        assert 'audit_id' in result
        assert result['audit_id'] == "test_audit_123"
        assert 'model_version' in result
        assert 'topics' in result
        assert 'suggestions' in result
        assert 'metadata' in result
        
        # Check topics
        assert isinstance(result['topics'], list)
        if result['topics']:
            topic = result['topics'][0]
            assert 'id' in topic
            assert 'label' in topic
            assert 'terms' in topic
            assert isinstance(topic['terms'], list)
        
        # Check suggestions
        assert isinstance(result['suggestions'], list)
        assert len(result['suggestions']) <= config['semantic']['output']['top_suggestions']
        
        if result['suggestions']:
            suggestion = result['suggestions'][0]
            assert 'keyword' in suggestion
            assert 'reason' in suggestion
            assert 'confidence' in suggestion
            assert 0 <= suggestion['confidence'] <= 1
            assert 'evidence' in suggestion
            assert isinstance(suggestion['evidence'], list)
    
    def test_ngram_extraction(self):
        """Test n-gram extraction from content"""
        stopwords = {'de', 'pour', 'd'}  # Add common stopwords
        analyzer = NgramAnalyzer(min_length=2, max_length=3, stopwords=stopwords)
        text = "logiciel de gestion pour cabinets d'avocats"
        ngrams = analyzer.extract_ngrams(text)
        
        # With stopwords, we should get meaningful combinations
        assert "logiciel gestion" in ngrams
        assert "gestion cabinets" in ngrams
        assert "cabinets avocats" in ngrams
    
    def test_topic_modeling(self, sample_crawl_data):
        """Test topic extraction from pages"""
        modeler = TopicModeler(model_name='distilcambert-base')
        
        # Extract all content
        texts = []
        for page in sample_crawl_data['pages']:
            text = f"{page['title']} {page['h1']} {' '.join(page['h2'])} {page['content']}"
            texts.append(text)
        
        topics = modeler.extract_topics(texts)
        
        assert isinstance(topics, list)
        assert len(topics) > 0
        assert all('id' in t and 'label' in t and 'terms' in t for t in topics)
    
    def test_keyword_ranking(self, config):
        """Test keyword ranking with multiple signals"""
        ranker = KeywordRanker(config['semantic']['ranker']['weights'])
        
        candidates = [
            {
                'keyword': 'logiciel avocat',
                'thematic_score': 0.9,
                'intent_score': 0.8,
                'mesh_evidence': 0.7,
                'readability': 0.85,
                'historical_feedback': 0.0
            },
            {
                'keyword': 'gestion cabinet juridique',
                'thematic_score': 0.85,
                'intent_score': 0.9,
                'mesh_evidence': 0.6,
                'readability': 0.8,
                'historical_feedback': 0.0
            }
        ]
        
        ranked = ranker.rank_keywords(candidates)
        
        assert len(ranked) == 2
        assert all('score' in kw for kw in ranked)
        assert ranked[0]['score'] >= ranked[1]['score']  # Descending order
    
    def test_filter_banned_keywords(self, config):
        """Test filtering of banned/generic keywords"""
        analyzer = SemanticAnalyzer(config)
        
        suggestions = [
            {'keyword': 'logiciel avocat', 'confidence': 0.9},
            {'keyword': 'cliquez ici', 'confidence': 0.8},  # Should be filtered
            {'keyword': 'accueil', 'confidence': 0.7},  # Should be filtered
            {'keyword': 'gestion juridique', 'confidence': 0.85}
        ]
        
        filtered = analyzer._filter_suggestions(suggestions)
        
        assert len(filtered) == 2
        assert all(s['keyword'] not in config['filters']['banlist'] for s in filtered)
        assert all(s['keyword'] not in config['filters']['generic_terms'] for s in filtered)
    
    def test_diversity_filtering(self, config):
        """Test diversity filtering to avoid redundant suggestions"""
        analyzer = SemanticAnalyzer(config)
        
        suggestions = [
            {'keyword': 'logiciel avocat', 'confidence': 0.9},
            {'keyword': 'avocat logiciel', 'confidence': 0.85},  # Same words, different order - high similarity
            {'keyword': 'gestion juridique', 'confidence': 0.8},
        ]
        
        # Test with similarity threshold that should filter similar keywords
        diverse = analyzer._apply_diversity_filter(suggestions, threshold=0.5)
        
        assert len(diverse) == 2  # Similar one should be filtered
        assert diverse[0]['keyword'] == 'logiciel avocat'  # Higher confidence kept
        assert diverse[1]['keyword'] == 'gestion juridique'
    
    def test_language_detection(self, config):
        """Test that non-French content is handled correctly"""
        analyzer = SemanticAnalyzer(config)
        
        english_data = {
            "pages": [{
                "url": "https://example.com/en",
                "lang": "en",
                "title": "Legal software for law firms",
                "content": "Our software helps law firms manage cases efficiently."
            }]
        }
        
        result = analyzer.analyze("test_audit_en", english_data)
        
        # Should return empty or minimal suggestions for non-French content
        assert len(result['suggestions']) == 0 or all(
            s['confidence'] < config['semantic']['language']['confidence_threshold'] 
            for s in result['suggestions']
        )
    
    @patch('semantic_analyzer.time.time')
    def test_execution_time_tracking(self, mock_time, config, sample_crawl_data):
        """Test that execution time is properly tracked"""
        mock_time.side_effect = [1000.0, 1000.5]  # 500ms execution
        
        analyzer = SemanticAnalyzer(config)
        result = analyzer.analyze("test_audit_123", sample_crawl_data)
        
        assert result['metadata']['execution_time_ms'] == 500
    
    def test_error_handling(self, config):
        """Test graceful error handling"""
        analyzer = SemanticAnalyzer(config)
        
        # Test with empty data
        result = analyzer.analyze("test_audit_empty", {"pages": []})
        assert result['suggestions'] == []
        assert result['topics'] == []
        
        # Test with malformed data - should return error result, not raise exception
        result = analyzer.analyze("test_audit_bad", {"invalid": "data"})
        assert 'error' in result['metadata']
        assert result['suggestions'] == []
        assert result['topics'] == []


class TestNgramAnalyzer:
    """Test n-gram extraction functionality"""
    
    def test_basic_ngram_extraction(self):
        stopwords = {'des', 'le', 'la', 'et'}
        analyzer = NgramAnalyzer(min_length=2, max_length=3, stopwords=stopwords)
        text = "analyse sémantique des mots clés"
        ngrams = analyzer.extract_ngrams(text)
        
        # With stopwords filtered, 'des' should be removed
        expected_present = ['analyse sémantique', 'sémantique mots', 'mots clés']
        expected_absent = ['sémantique des', 'des mots']
        
        for expected in expected_present:
            assert expected in ngrams, f"Expected '{expected}' in ngrams"
        
        for unexpected in expected_absent:
            assert unexpected not in ngrams, f"Unexpected '{unexpected}' in ngrams"
    
    def test_stopword_filtering(self):
        stopwords = {'des', 'le', 'la', 'et'}
        analyzer = NgramAnalyzer(min_length=2, max_length=2, stopwords=stopwords)
        text = "analyse des mots et concepts"
        ngrams = analyzer.extract_ngrams(text)
        
        # 'des' and 'et' should be filtered
        assert 'analyse mots' in ngrams
        assert 'mots concepts' in ngrams
        assert 'analyse des' not in ngrams
        assert 'et concepts' not in ngrams
    
    def test_frequency_counting(self):
        analyzer = NgramAnalyzer(min_length=2, max_length=2)
        texts = [
            "logiciel gestion avocat",
            "logiciel gestion cabinet",
            "gestion cabinet avocat"
        ]
        
        freq_map = analyzer.count_frequencies(texts)
        
        assert freq_map['logiciel gestion'] == 2
        assert freq_map['gestion cabinet'] == 2
        assert freq_map['gestion avocat'] == 1
        assert freq_map['cabinet avocat'] == 1


class TestKeywordRanker:
    """Test keyword ranking logic"""
    
    def test_weighted_scoring(self):
        weights = {
            'thematic_score': 0.4,
            'intent_score': 0.3,
            'mesh_evidence': 0.2,
            'readability': 0.1,
            'historical_feedback': 0.0
        }
        ranker = KeywordRanker(weights)
        
        keyword = {
            'keyword': 'test keyword',
            'thematic_score': 0.8,
            'intent_score': 0.9,
            'mesh_evidence': 0.7,
            'readability': 0.6,
            'historical_feedback': 0.0
        }
        
        score = ranker.calculate_score(keyword)
        expected = (0.8 * 0.4) + (0.9 * 0.3) + (0.7 * 0.2) + (0.6 * 0.1)
        assert abs(score - expected) < 0.001
    
    def test_ranking_order(self):
        ranker = KeywordRanker({
            'thematic_score': 1.0,
            'intent_score': 0.0,
            'mesh_evidence': 0.0,
            'readability': 0.0,
            'historical_feedback': 0.0
        })
        
        keywords = [
            {'keyword': 'low', 'thematic_score': 0.3},
            {'keyword': 'high', 'thematic_score': 0.9},
            {'keyword': 'medium', 'thematic_score': 0.6}
        ]
        
        ranked = ranker.rank_keywords(keywords)
        
        assert ranked[0]['keyword'] == 'high'
        assert ranked[1]['keyword'] == 'medium'
        assert ranked[2]['keyword'] == 'low'