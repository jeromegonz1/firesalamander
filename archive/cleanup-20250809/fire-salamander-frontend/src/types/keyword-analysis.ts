/**
 * Fire Salamander - Keyword Analysis Types
 * Professional keyword analysis data structure (SEMrush/Ahrefs style)
 */

export interface KeywordAnalysis {
  // Mots-clés trouvés sur la page
  foundKeywords: Array<{
    keyword: string;
    density: number; // Pourcentage
    occurrences: number;
    prominence: number; // Score de 0-100 basé sur la position
    locations: Array<'title' | 'h1' | 'h2' | 'h3' | 'meta' | 'content' | 'alt' | 'url'>;
    variations?: string[]; // Variantes trouvées
  }>;
  
  // Analyse n-grammes
  ngrams: {
    bigrams: Array<{ 
      phrase: string;
      count: number;
      density: number;
    }>;
    trigrams: Array<{ 
      phrase: string;
      count: number;
      density: number;
    }>;
    topPhrases?: Array<{
      phrase: string;
      count: number;
      type: 'bigram' | 'trigram';
    }>;
  };
  
  // Suggestions IA (ChatGPT 3.5)
  aiSuggestions: Array<{
    keyword: string;
    searchVolume: number;
    difficulty: number; // 0-100
    cpc: number; // Coût par clic en €
    intent: 'informational' | 'navigational' | 'transactional' | 'commercial';
    reason: string; // Explication de la suggestion
    contentIdeas: string[]; // Idées de contenu
    trendData?: {
      direction: 'up' | 'down' | 'stable';
      changePercent: number;
    };
    relatedKeywords?: string[];
  }>;
  
  // Analyse sémantique
  semanticAnalysis: {
    mainTopics: string[];
    entities: Array<{ 
      name: string;
      type: 'Person' | 'Organization' | 'Location' | 'Product' | 'Event' | 'Other';
      salience: number; // Importance 0-1
      mentions: number;
    }>;
    sentiment: 'positive' | 'neutral' | 'negative' | 'mixed';
    sentimentScore?: number; // -1 to 1
    readability: {
      score: number; // Flesch Reading Ease (0-100)
      level: 'very-easy' | 'easy' | 'fairly-easy' | 'standard' | 'fairly-difficult' | 'difficult' | 'very-difficult';
      avgSentenceLength: number;
      avgWordLength: number;
      avgWordsPerSentence: number;
      complexWords: number; // Pourcentage
      paragraphs: number;
    };
    language: string; // Code ISO (fr, en, etc.)
    wordCount: number;
  };
  
  // Comparaison concurrence (optionnel)
  competitorGaps?: Array<{
    keyword: string;
    competitorRank: number;
    ourRank: number | null;
    searchVolume: number;
    difficulty: number;
    opportunity: 'high' | 'medium' | 'low';
    estimatedTraffic: number;
    competitors: Array<{
      domain: string;
      rank: number;
    }>;
  }>;

  // Métriques globales
  metrics: {
    totalKeywords: number;
    uniqueKeywords: number;
    avgKeywordDensity: number;
    keywordDiversity: number; // Score 0-100
    focusKeywords: string[]; // Top 3-5 mots-clés principaux
    missingImportantKeywords?: string[]; // Mots-clés importants manquants
  };

  // Statut de l'analyse
  status: {
    lastUpdated: string;
    aiAnalysisEnabled: boolean;
    competitorAnalysisEnabled: boolean;
    crawledPages: number;
  };
}

// Types pour les composants UI
export interface KeywordItemProps {
  keyword: string;
  density: number;
  occurrences: number;
  prominence: number;
  locations: string[];
  isExpanded?: boolean;
  onToggle?: () => void;
}

export interface NGramChartProps {
  data: Array<{
    phrase: string;
    count: number;
    density: number;
  }>;
  type: 'bigram' | 'trigram';
  limit?: number;
}

export interface AIKeywordSuggestionProps {
  suggestion: KeywordAnalysis['aiSuggestions'][0];
  onSelect?: (keyword: string) => void;
  onGenerateContent?: (keyword: string) => void;
}

export interface ReadabilityGaugeProps {
  score: number;
  level: string;
  metrics: KeywordAnalysis['semanticAnalysis']['readability'];
}

// Enums pour filtrage et tri
export enum KeywordSortBy {
  DENSITY = 'density',
  OCCURRENCES = 'occurrences',
  PROMINENCE = 'prominence',
  ALPHABETICAL = 'alphabetical'
}

export enum KeywordFilterBy {
  ALL = 'all',
  TITLE = 'title',
  HEADINGS = 'headings',
  CONTENT = 'content',
  HIGH_DENSITY = 'high_density',
  LOW_DENSITY = 'low_density'
}

// Types pour les requêtes API
export interface KeywordAnalysisRequest {
  analysisId: number;
  includeAI?: boolean;
  includeCompetitors?: boolean;
  targetKeywords?: string[];
  language?: string;
}

export interface AIKeywordGenerationRequest {
  analysisId: number;
  context: string;
  currentKeywords: string[];
  targetAudience?: string;
  contentType?: string;
  limit?: number;
}

// Utility types
export type KeywordLocation = KeywordAnalysis['foundKeywords'][0]['locations'][0];
export type SearchIntent = KeywordAnalysis['aiSuggestions'][0]['intent'];
export type EntityType = KeywordAnalysis['semanticAnalysis']['entities'][0]['type'];