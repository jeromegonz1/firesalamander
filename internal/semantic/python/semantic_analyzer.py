"""
Main semantic analyzer that orchestrates the entire semantic analysis pipeline.
"""
import json
import time
import yaml
from pathlib import Path
from typing import Dict, List, Any
import logging
import numpy as np
from sklearn.metrics.pairwise import cosine_similarity

from ngram_analyzer import NgramAnalyzer
from topic_modeler import TopicModeler
from keyword_ranker import KeywordRanker

logger = logging.getLogger(__name__)


class SemanticAnalyzer:
    """Main semantic analysis agent for Fire Salamander"""
    
    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.semantic_config = config['semantic']
        self.filters = config['filters']
        
        # Basic properties
        self.language = self.semantic_config['language']['target']
        self.model_version = 'sem-v1.0'
        
        # Initialize components
        self.ngram_analyzer = NgramAnalyzer(
            min_length=self.semantic_config['preprocessing']['min_ngram_length'],
            max_length=self.semantic_config['preprocessing']['max_ngram_length'],
            stopwords=self._load_stopwords()
        )
        
        self.topic_modeler = TopicModeler(
            model_name=self.semantic_config['models']['embeddings'],
            min_cluster_size=self.semantic_config['clustering']['min_cluster_size']
        )
        
        self.keyword_ranker = KeywordRanker(
            weights=self.semantic_config['ranker']['weights']
        )
        
        logger.info(f"SemanticAnalyzer initialized for language: {self.language}")
    
    def analyze(self, audit_id: str, crawl_data: Dict[str, Any]) -> Dict[str, Any]:
        """
        Main analysis method that processes crawl data and returns semantic insights
        """
        start_time = time.time()
        
        try:
            # Validate input
            if 'pages' not in crawl_data:
                raise ValueError("Crawl data must contain 'pages' key")
            
            pages = crawl_data['pages']
            
            if not pages:
                return self._empty_result(audit_id, start_time)
            
            # Step 1: Extract topics from page content
            topics = self._extract_topics(pages)
            
            # Step 2: Generate keyword candidates using n-grams
            candidates = self._generate_keyword_candidates(pages)
            
            # Step 3: Score and rank candidates
            scored_candidates = self._score_candidates(candidates, pages)
            
            # Step 4: Apply filters and diversity
            filtered_suggestions = self._filter_suggestions(scored_candidates)
            diverse_suggestions = self._apply_diversity_filter(filtered_suggestions)
            
            # Step 5: Limit to top suggestions
            top_suggestions = diverse_suggestions[:self.semantic_config['output']['top_suggestions']]
            
            # Build result
            execution_time = int((time.time() - start_time) * 1000)
            
            return {
                'audit_id': audit_id,
                'model_version': self.model_version,
                'topics': topics,
                'suggestions': top_suggestions,
                'metadata': {
                    'schema_version': '1.0',
                    'weights_version': '1.0',
                    'execution_time_ms': execution_time,
                    'lang': self.language
                }
            }
            
        except Exception as e:
            logger.error(f"Semantic analysis failed for audit {audit_id}: {e}")
            return self._error_result(audit_id, str(e), start_time)
    
    def _extract_topics(self, pages: List[Dict]) -> List[Dict[str, Any]]:
        """Extract semantic topics from all page content"""
        # Combine all text content
        texts = []
        for page in pages:
            # Only process French content
            if page.get('lang', '').lower() not in ['fr', 'fr-fr']:
                continue
                
            combined_text = ' '.join([
                page.get('title', ''),
                page.get('h1', ''),
                ' '.join(page.get('h2', [])),
                ' '.join(page.get('h3', [])),
                page.get('content', '')
            ]).strip()
            
            if combined_text:
                texts.append(combined_text)
        
        if not texts:
            return []
        
        return self.topic_modeler.extract_topics(texts)
    
    def _generate_keyword_candidates(self, pages: List[Dict]) -> List[Dict[str, Any]]:
        """Generate keyword candidates from page content using n-grams"""
        candidates = []
        
        # Collect all texts for frequency analysis
        all_texts = []
        for page in pages:
            if page.get('lang', '').lower() not in ['fr', 'fr-fr']:
                continue
                
            # Combine different content types with different weights
            title_text = page.get('title', '') * 3  # Titles are more important
            heading_text = ' '.join([page.get('h1', '')] + page.get('h2', []) + page.get('h3', [])) * 2
            content_text = page.get('content', '')
            
            combined = f"{title_text} {heading_text} {content_text}".strip()
            if combined:
                all_texts.append(combined)
        
        if not all_texts:
            return []
        
        # Extract n-gram frequencies
        ngram_frequencies = self.ngram_analyzer.count_frequencies(all_texts)
        
        # Filter by minimum frequency
        min_freq = self.semantic_config['preprocessing']['min_word_count']
        filtered_ngrams = self.ngram_analyzer.filter_by_frequency(ngram_frequencies, min_freq)
        
        # Convert to candidate format
        for ngram, frequency in filtered_ngrams.items():
            if len(ngram.strip()) >= 3:  # Minimum length filter
                candidate = {
                    'keyword': ngram,
                    'frequency': frequency,
                    'thematic_score': 0.0,  # Will be calculated in scoring
                    'intent_score': 0.0,
                    'mesh_evidence': 0.0,
                    'readability': 0.0,
                    'historical_feedback': 0.0
                }
                candidates.append(candidate)
        
        return candidates
    
    def _score_candidates(self, candidates: List[Dict], pages: List[Dict]) -> List[Dict[str, Any]]:
        """Score all keyword candidates using multiple signals"""
        scored_candidates = []
        
        # Prepare content texts for thematic scoring
        content_texts = []
        for page in pages:
            if page.get('lang', '').lower() in ['fr', 'fr-fr']:
                text = f"{page.get('title', '')} {page.get('content', '')}"
                content_texts.append(text)
        
        for candidate in candidates:
            keyword = candidate['keyword']
            
            # Calculate all scoring signals
            candidate['thematic_score'] = self.keyword_ranker.calculate_thematic_score(keyword, content_texts)
            candidate['intent_score'] = self.keyword_ranker.calculate_intent_score(keyword)
            candidate['mesh_evidence'] = self.keyword_ranker.calculate_mesh_evidence(keyword, pages)
            candidate['readability'] = self.keyword_ranker.calculate_readability_score(keyword)
            # historical_feedback remains 0.0 for new keywords
            
            # Calculate confidence (overall score)
            confidence = self.keyword_ranker.calculate_score(candidate)
            candidate['confidence'] = confidence
            
            # Add reason and evidence
            candidate['reason'] = self._generate_reason(candidate)
            candidate['evidence'] = self._generate_evidence(candidate, pages)
            
            scored_candidates.append(candidate)
        
        return scored_candidates
    
    def _filter_suggestions(self, suggestions: List[Dict]) -> List[Dict]:
        """Filter out banned and generic keywords"""
        filtered = []
        
        for suggestion in suggestions:
            keyword = suggestion['keyword'].lower()
            
            # Check banlist
            if any(banned in keyword for banned in self.filters['banlist']):
                continue
            
            # Check generic terms
            if any(generic in keyword for generic in self.filters['generic_terms']):
                continue
            
            # Check brand terms (might want to keep or filter depending on strategy)
            # For now, we keep them but could add logic here
            
            filtered.append(suggestion)
        
        return filtered
    
    def _apply_diversity_filter(self, suggestions: List[Dict], threshold: float = None) -> List[Dict]:
        """Apply diversity filtering to avoid too similar suggestions"""
        if threshold is None:
            threshold = self.semantic_config['output']['diversity_threshold']
        
        if not suggestions:
            return []
        
        # Sort by confidence first
        suggestions = sorted(suggestions, key=lambda x: x['confidence'], reverse=True)
        
        diverse_suggestions = [suggestions[0]]  # Always keep the top one
        
        for suggestion in suggestions[1:]:
            # Check similarity with already selected suggestions
            is_diverse = True
            
            for selected in diverse_suggestions:
                similarity = self._calculate_keyword_similarity(
                    suggestion['keyword'], 
                    selected['keyword']
                )
                
                if similarity > threshold:
                    is_diverse = False
                    break
            
            if is_diverse:
                diverse_suggestions.append(suggestion)
        
        return diverse_suggestions
    
    def _calculate_keyword_similarity(self, kw1: str, kw2: str) -> float:
        """Calculate similarity between two keywords"""
        # Simple word overlap similarity
        words1 = set(kw1.lower().split())
        words2 = set(kw2.lower().split())
        
        if not words1 or not words2:
            return 0.0
        
        intersection = len(words1 & words2)
        union = len(words1 | words2)
        
        return intersection / union if union > 0 else 0.0
    
    def _generate_reason(self, candidate: Dict) -> str:
        """Generate human-readable reason for keyword suggestion"""
        scores = {
            'thématique': candidate['thematic_score'],
            'intention': candidate['intent_score'],
            'maillage': candidate['mesh_evidence'],
            'lisibilité': candidate['readability']
        }
        
        # Find strongest signal
        best_signal = max(scores.keys(), key=lambda k: scores[k])
        best_score = scores[best_signal]
        
        if best_score > 0.8:
            return f"Très pertinent pour le {best_signal} (score: {best_score:.2f})"
        elif best_score > 0.6:
            return f"Pertinent pour le {best_signal} (score: {best_score:.2f})"
        else:
            return f"Candidat basé sur la fréquence ({candidate.get('frequency', 0)} occurrences)"
    
    def _generate_evidence(self, candidate: Dict, pages: List[Dict]) -> List[str]:
        """Generate evidence URLs where keyword appears"""
        evidence = []
        keyword = candidate['keyword'].lower()
        
        for page in pages:
            page_url = page.get('url', '')
            
            # Check title
            if keyword in page.get('title', '').lower():
                evidence.append(f"{page_url} (titre)")
            
            # Check headings
            elif keyword in page.get('h1', '').lower():
                evidence.append(f"{page_url} (H1)")
            
            elif any(keyword in h.lower() for h in page.get('h2', [])):
                evidence.append(f"{page_url} (H2)")
            
            # Check content
            elif keyword in page.get('content', '').lower():
                evidence.append(f"{page_url} (contenu)")
            
            # Limit evidence to avoid too long lists
            if len(evidence) >= 5:
                break
        
        return evidence
    
    def _load_stopwords(self) -> set:
        """Load French stopwords from configuration"""
        stopwords_file = self.semantic_config['preprocessing']['stopwords_file']
        
        try:
            with open(stopwords_file, 'r', encoding='utf-8') as f:
                stopwords = set(line.strip().lower() for line in f if line.strip())
                logger.info(f"Loaded {len(stopwords)} stopwords from {stopwords_file}")
                return stopwords
        except FileNotFoundError:
            logger.warning(f"Stopwords file not found: {stopwords_file}, using default set")
            return self._default_french_stopwords()
    
    def _default_french_stopwords(self) -> set:
        """Default French stopwords if file is not available"""
        return {
            'le', 'de', 'et', 'à', 'un', 'il', 'être', 'et', 'en', 'avoir', 'que', 'pour',
            'dans', 'ce', 'son', 'une', 'avec', 'ne', 'se', 'pas', 'tout', 'plus', 'par',
            'sur', 'du', 'des', 'les', 'la', 'ou', 'qui', 'nous', 'vous', 'ils', 'elle',
            'aux', 'ses', 'ces', 'leur', 'leurs', 'mais', 'si', 'bien', 'comme', 'très',
            'donc', 'sans', 'chez', 'vers', 'sous', 'depuis', 'pendant', 'selon'
        }
    
    def _empty_result(self, audit_id: str, start_time: float) -> Dict[str, Any]:
        """Return empty result for cases with no data"""
        execution_time = int((time.time() - start_time) * 1000)
        
        return {
            'audit_id': audit_id,
            'model_version': self.model_version,
            'topics': [],
            'suggestions': [],
            'metadata': {
                'schema_version': '1.0',
                'weights_version': '1.0',
                'execution_time_ms': execution_time,
                'lang': self.language
            }
        }
    
    def _error_result(self, audit_id: str, error_msg: str, start_time: float) -> Dict[str, Any]:
        """Return error result"""
        execution_time = int((time.time() - start_time) * 1000)
        
        return {
            'audit_id': audit_id,
            'model_version': self.model_version,
            'topics': [],
            'suggestions': [],
            'metadata': {
                'schema_version': '1.0',
                'weights_version': '1.0',
                'execution_time_ms': execution_time,
                'lang': self.language,
                'error': error_msg
            }
        }


def load_config(config_path: str = 'config/semantic.yaml') -> Dict[str, Any]:
    """Load semantic analysis configuration"""
    try:
        with open(config_path, 'r', encoding='utf-8') as f:
            return yaml.safe_load(f)
    except FileNotFoundError:
        logger.error(f"Configuration file not found: {config_path}")
        raise
    except yaml.YAMLError as e:
        logger.error(f"Invalid YAML in configuration: {e}")
        raise


if __name__ == "__main__":
    # CLI entry point for standalone testing
    import sys
    
    if len(sys.argv) != 3:
        print("Usage: python semantic_analyzer.py <audit_id> <crawl_data.json>")
        sys.exit(1)
    
    audit_id = sys.argv[1]
    crawl_data_file = sys.argv[2]
    
    # Load configuration
    config = load_config()
    
    # Load crawl data
    with open(crawl_data_file, 'r', encoding='utf-8') as f:
        crawl_data = json.load(f)
    
    # Analyze
    analyzer = SemanticAnalyzer(config)
    result = analyzer.analyze(audit_id, crawl_data)
    
    # Output result
    print(json.dumps(result, indent=2, ensure_ascii=False))