"""
Topic modeling using French language models for semantic analysis.
"""
import numpy as np
from typing import List, Dict, Any
from collections import defaultdict, Counter
import logging
import re
from sklearn.cluster import DBSCAN
from sklearn.feature_extraction.text import TfidfVectorizer

logger = logging.getLogger(__name__)


class TopicModeler:
    """Extract semantic topics from content using embeddings and clustering"""
    
    def __init__(self, model_name: str = 'distilcambert-base', min_cluster_size: int = 3):
        self.model_name = model_name
        self.min_cluster_size = min_cluster_size
        self.model = None
        self._load_model()
    
    def _load_model(self):
        """Load the sentence transformer model"""
        try:
            # For now, we'll use TF-IDF based approach
            # TODO: Integrate actual CamemBERT/DistilCamemBERT when dependencies are available
            self.model = None  # Will use TF-IDF fallback
            logger.info(f"Using TF-IDF based topic modeling (lightweight mode)")
        except Exception as e:
            logger.warning(f"Could not load {self.model_name}, using fallback: {e}")
            self.model = None
    
    def extract_topics(self, texts: List[str], num_topics: int = 5) -> List[Dict[str, Any]]:
        """Extract topics from a list of texts"""
        if not texts:
            return []
        
        try:
            # Method 1: If we have embeddings model, use clustering
            if self.model:
                return self._extract_topics_clustering(texts, num_topics)
            else:
                # Method 2: Fallback to TF-IDF based topic extraction
                return self._extract_topics_tfidf(texts, num_topics)
        
        except Exception as e:
            logger.error(f"Topic extraction failed: {e}")
            return self._fallback_topics(texts)
    
    def _extract_topics_clustering(self, texts: List[str], num_topics: int) -> List[Dict[str, Any]]:
        """Extract topics using embeddings and clustering"""
        # Generate embeddings
        embeddings = self.model.encode(texts)
        
        # Cluster embeddings
        clustering = DBSCAN(eps=0.5, min_samples=2)
        cluster_labels = clustering.fit_predict(embeddings)
        
        # Extract topics from clusters
        topics = []
        clusters = defaultdict(list)
        
        for i, label in enumerate(cluster_labels):
            if label != -1:  # Not noise
                clusters[label].append((i, texts[i]))
        
        for cluster_id, cluster_texts in clusters.items():
            if len(cluster_texts) >= 2:
                # Extract key terms from cluster
                cluster_content = [text for _, text in cluster_texts]
                terms = self._extract_cluster_terms(cluster_content)
                
                topic = {
                    'id': f'topic_{cluster_id}',
                    'label': self._generate_topic_label(terms),
                    'terms': terms[:10]  # Top 10 terms
                }
                topics.append(topic)
        
        return topics[:num_topics]
    
    def _extract_topics_tfidf(self, texts: List[str], num_topics: int) -> List[Dict[str, Any]]:
        """Extract topics using TF-IDF (fallback method)"""
        vectorizer = TfidfVectorizer(
            max_features=1000,
            stop_words=None,  # We'll handle French stopwords separately
            ngram_range=(1, 3),
            min_df=2
        )
        
        try:
            tfidf_matrix = vectorizer.fit_transform(texts)
            feature_names = vectorizer.get_feature_names_out()
            
            # Get top terms across all documents
            mean_scores = np.mean(tfidf_matrix.toarray(), axis=0)
            top_indices = np.argsort(mean_scores)[-50:][::-1]  # Top 50 terms
            
            # Group into topics (simple clustering)
            topics = []
            terms_per_topic = len(top_indices) // num_topics
            
            for i in range(num_topics):
                start_idx = i * terms_per_topic
                end_idx = start_idx + terms_per_topic if i < num_topics - 1 else len(top_indices)
                
                topic_terms = [feature_names[idx] for idx in top_indices[start_idx:end_idx]]
                topic_terms = [term for term in topic_terms if len(term) > 2]
                
                if topic_terms:
                    topic = {
                        'id': f'tfidf_topic_{i}',
                        'label': self._generate_topic_label(topic_terms),
                        'terms': topic_terms[:10]
                    }
                    topics.append(topic)
            
            return topics
            
        except Exception as e:
            logger.error(f"TF-IDF topic extraction failed: {e}")
            return self._fallback_topics(texts)
    
    def _extract_cluster_terms(self, texts: List[str]) -> List[str]:
        """Extract representative terms from a cluster of texts"""
        # Simple frequency-based extraction
        all_words = []
        for text in texts:
            words = re.findall(r'\b[a-zàâäéèêëïîôùûüÿç]{3,}\b', text.lower())
            all_words.extend(words)
        
        word_freq = Counter(all_words)
        return [word for word, freq in word_freq.most_common(20) if freq >= 2]
    
    def _generate_topic_label(self, terms: List[str]) -> str:
        """Generate a human-readable label for a topic"""
        if not terms:
            return "Topic inconnu"
        
        # Take the most representative terms
        top_terms = terms[:3]
        return ' & '.join(top_terms).title()
    
    def _fallback_topics(self, texts: List[str]) -> List[Dict[str, Any]]:
        """Fallback topic extraction when models fail"""
        # Extract common words as a single topic
        all_text = ' '.join(texts).lower()
        words = re.findall(r'\b[a-zàâäéèêëïîôùûüÿç]{4,}\b', all_text)
        
        if not words:
            return []
        
        word_freq = Counter(words)
        top_words = [word for word, freq in word_freq.most_common(10) if freq >= 2]
        
        if top_words:
            return [{
                'id': 'fallback_topic_0',
                'label': f"Thème principal: {top_words[0]}",
                'terms': top_words
            }]
        
        return []


# Import Counter for the fallback method
from collections import Counter
import re