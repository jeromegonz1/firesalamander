"""
Keyword ranking system with multiple signals for semantic analysis.
"""
import numpy as np
from typing import List, Dict, Any
import re
import logging

logger = logging.getLogger(__name__)


class KeywordRanker:
    """Rank keyword candidates using multiple scoring signals"""
    
    def __init__(self, weights: Dict[str, float]):
        self.weights = weights
        self._validate_weights()
    
    def _validate_weights(self):
        """Ensure weights sum to approximately 1.0"""
        total = sum(self.weights.values())
        if abs(total - 1.0) > 0.01:
            logger.warning(f"Weights sum to {total}, expected ~1.0")
    
    def rank_keywords(self, candidates: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
        """Rank keyword candidates by weighted score"""
        scored_candidates = []
        
        for candidate in candidates:
            score = self.calculate_score(candidate)
            candidate_copy = candidate.copy()
            candidate_copy['score'] = score
            scored_candidates.append(candidate_copy)
        
        # Sort by score descending
        return sorted(scored_candidates, key=lambda x: x['score'], reverse=True)
    
    def calculate_score(self, candidate: Dict[str, Any]) -> float:
        """Calculate weighted score for a keyword candidate"""
        total_score = 0.0
        
        for signal, weight in self.weights.items():
            signal_score = candidate.get(signal, 0.0)
            total_score += signal_score * weight
        
        return total_score
    
    def calculate_thematic_score(self, keyword: str, content_texts: List[str]) -> float:
        """Calculate how well keyword fits the thematic content"""
        if not content_texts or not keyword:
            return 0.0
        
        keyword_lower = keyword.lower()
        keyword_words = set(keyword_lower.split())
        
        # Count relevance across all texts
        relevance_scores = []
        for text in content_texts:
            text_lower = text.lower()
            text_words = set(re.findall(r'\b[a-zàâäéèêëïîôùûüÿç]+\b', text_lower))
            
            # Exact match bonus
            exact_match = keyword_lower in text_lower
            
            # Word overlap score
            word_overlap = len(keyword_words & text_words) / len(keyword_words) if keyword_words else 0
            
            # Context relevance (words appearing near keyword)
            context_score = self._calculate_context_relevance(keyword_lower, text_lower)
            
            text_score = (0.5 * (1.0 if exact_match else 0.0) + 
                         0.3 * word_overlap + 
                         0.2 * context_score)
            relevance_scores.append(text_score)
        
        return np.mean(relevance_scores) if relevance_scores else 0.0
    
    def calculate_intent_score(self, keyword: str) -> float:
        """Calculate search intent score (commercial, informational, etc.)"""
        keyword_lower = keyword.lower()
        
        # Commercial intent indicators
        commercial_terms = ['logiciel', 'solution', 'service', 'prix', 'tarif', 'gratuit', 'acheter']
        commercial_score = sum(1 for term in commercial_terms if term in keyword_lower)
        
        # Informational intent indicators  
        info_terms = ['comment', 'pourquoi', 'qu\'est-ce', 'guide', 'aide', 'conseil']
        info_score = sum(1 for term in info_terms if term in keyword_lower)
        
        # Navigational intent indicators
        nav_terms = ['connexion', 'login', 'accueil', 'contact', 'à propos']
        nav_score = sum(1 for term in nav_terms if term in keyword_lower)
        
        # Weight commercial intent higher for SEO
        total_signals = commercial_score * 1.0 + info_score * 0.7 + nav_score * 0.3
        
        # Normalize to 0-1 range
        return min(1.0, total_signals / 3.0)
    
    def calculate_mesh_evidence(self, keyword: str, page_data: List[Dict]) -> float:
        """Calculate how well keyword is supported by internal linking mesh"""
        if not page_data or not keyword:
            return 0.0
        
        keyword_lower = keyword.lower()
        evidence_count = 0
        total_pages = len(page_data)
        
        for page in page_data:
            # Check if keyword appears in anchors
            anchors = page.get('anchors', [])
            for anchor in anchors:
                anchor_text = anchor.get('text', '').lower()
                if keyword_lower in anchor_text:
                    evidence_count += 1
                    break  # Count once per page
            
            # Check if keyword appears in headings
            h1 = page.get('h1', '').lower()
            h2_list = page.get('h2', [])
            h3_list = page.get('h3', [])
            
            all_headings = [h1] + h2_list + h3_list
            for heading in all_headings:
                if keyword_lower in heading.lower():
                    evidence_count += 1
                    break  # Count once per page
        
        # Return ratio of pages with evidence
        return evidence_count / total_pages if total_pages > 0 else 0.0

    def calculate_readability_score(self, keyword: str) -> float:
        """Calculate readability/usability score for keyword"""
        if not keyword:
            return 0.0
        
        # Factors that improve readability
        score = 0.0
        
        # Length penalty (very long keywords are harder to read)
        length_penalty = max(0, 1.0 - (len(keyword) - 30) / 50) if len(keyword) > 30 else 1.0
        score += 0.3 * length_penalty
        
        # Word count (2-4 words is optimal)
        word_count = len(keyword.split())
        if 2 <= word_count <= 4:
            score += 0.3
        elif word_count == 1:
            score += 0.2  # Single words are less specific
        else:
            score += 0.1  # 5+ words get penalty
        
        # Natural language flow
        has_natural_flow = self._has_natural_flow(keyword)
        score += 0.2 if has_natural_flow else 0.0
        
        # No special characters (except accents and hyphens)
        clean_chars = re.match(r'^[a-zàâäéèêëïîôùûüÿç\s\'-]+$', keyword.lower()) is not None
        score += 0.2 if clean_chars else 0.0
        
        return min(1.0, score)
    
    def _has_natural_flow(self, keyword: str) -> bool:
        """Check if keyword has natural French language flow"""
        words = keyword.lower().split()
        
        if len(words) < 2:
            return True
        
        # Check for common French patterns
        common_patterns = [
            # Article + noun patterns
            (r'^(le|la|les|un|une|des)\s+\w+', True),
            # Preposition patterns  
            (r'\b(de|du|des|pour|avec|sans|sur|dans)\s+\w+', True),
            # Adjective + noun patterns
            (r'\w+\s+(logiciel|solution|service|gestion)', True),
        ]
        
        keyword_lower = keyword.lower()
        for pattern, is_good in common_patterns:
            if re.search(pattern, keyword_lower):
                return is_good
        
        return True  # Default to natural if no specific pattern detected
    
    def _calculate_context_relevance(self, keyword: str, text: str) -> float:
        """Calculate how relevant keyword is in its context"""
        if keyword not in text:
            return 0.0
        
        # Find all occurrences of keyword
        keyword_positions = []
        start = 0
        while True:
            pos = text.find(keyword, start)
            if pos == -1:
                break
            keyword_positions.append(pos)
            start = pos + 1
        
        if not keyword_positions:
            return 0.0
        
        # Analyze context around each occurrence
        context_scores = []
        for pos in keyword_positions:
            # Extract context window (±50 characters)
            context_start = max(0, pos - 50)
            context_end = min(len(text), pos + len(keyword) + 50)
            context = text[context_start:context_end]
            
            # Score based on context quality
            context_score = self._score_context(context, keyword)
            context_scores.append(context_score)
        
        return np.mean(context_scores) if context_scores else 0.0
    
    def _score_context(self, context: str, keyword: str) -> float:
        """Score the quality of context around a keyword"""
        context_lower = context.lower()
        
        # Positive indicators
        positive_terms = [
            'logiciel', 'solution', 'service', 'professionnel', 'efficace',
            'gestion', 'cabinet', 'avocat', 'juridique', 'droit'
        ]
        
        positive_count = sum(1 for term in positive_terms if term in context_lower)
        
        # Negative indicators (spam-like context)
        negative_terms = ['cliquez', 'gratuit', 'promo', 'offre spéciale']
        negative_count = sum(1 for term in negative_terms if term in context_lower)
        
        # Calculate score
        score = (positive_count * 0.2) - (negative_count * 0.3)
        return max(0.0, min(1.0, score))