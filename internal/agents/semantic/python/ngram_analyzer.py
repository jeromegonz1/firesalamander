"""
N-gram analyzer for extracting keyword candidates from French text.
"""
import re
from typing import List, Set, Dict
from collections import Counter


class NgramAnalyzer:
    """Extract and analyze n-grams from text content"""
    
    def __init__(self, min_length: int = 2, max_length: int = 5, stopwords: Set[str] = None):
        self.min_length = min_length
        self.max_length = max_length
        self.stopwords = stopwords or set()
    
    def extract_ngrams(self, text: str) -> List[str]:
        """Extract n-grams from text"""
        # Clean and tokenize
        text = self._clean_text(text)
        tokens = self._tokenize(text)
        
        # Filter tokens
        filtered_tokens = [t for t in tokens if t not in self.stopwords and len(t) > 2]
        
        # Extract n-grams of different lengths using simple sliding window
        all_ngrams = []
        for n in range(self.min_length, self.max_length + 1):
            if len(filtered_tokens) >= n:
                # Simple n-gram extraction
                for i in range(len(filtered_tokens) - n + 1):
                    ngram = ' '.join(filtered_tokens[i:i+n])
                    all_ngrams.append(ngram)
        
        return all_ngrams
    
    def count_frequencies(self, texts: List[str]) -> Dict[str, int]:
        """Count n-gram frequencies across multiple texts"""
        all_ngrams = []
        for text in texts:
            all_ngrams.extend(self.extract_ngrams(text))
        
        return dict(Counter(all_ngrams))
    
    def _clean_text(self, text: str) -> str:
        """Clean text for n-gram extraction"""
        # Convert to lowercase
        text = text.lower()
        
        # Remove URLs
        text = re.sub(r'https?://\S+', '', text)
        
        # Remove email addresses
        text = re.sub(r'\S+@\S+', '', text)
        
        # Keep only letters, numbers, spaces, and basic punctuation
        text = re.sub(r'[^a-zàâäéèêëïîôùûüÿç0-9\s\'-]', ' ', text)
        
        # Normalize whitespace
        text = re.sub(r'\s+', ' ', text)
        
        return text.strip()
    
    def _tokenize(self, text: str) -> List[str]:
        """Tokenize text into words using simple regex"""
        # Handle contractions
        text = re.sub(r"([dlmts])'", r"\1' ", text)
        
        # Simple tokenization with regex
        tokens = re.findall(r'\b[a-zàâäéèêëïîôùûüÿç]+\b', text, re.IGNORECASE)
        
        # Filter out short tokens
        tokens = [t.lower() for t in tokens if len(t) > 1]
        
        return tokens
    
    def filter_by_frequency(self, ngram_freq: Dict[str, int], min_freq: int = 2) -> Dict[str, int]:
        """Filter n-grams by minimum frequency"""
        return {ng: freq for ng, freq in ngram_freq.items() if freq >= min_freq}