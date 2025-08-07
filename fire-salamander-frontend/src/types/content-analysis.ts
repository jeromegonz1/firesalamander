/**
 * Fire Salamander - Content Analysis Types
 * Professional Content Analysis with AI Insights & Advanced Visualizations
 * Lead Tech implementation - Niveau SEMrush/Ahrefs
 */

// Enums pour type safety stricte
export enum ContentQuality {
  EXCELLENT = 'excellent',
  GOOD = 'good', 
  AVERAGE = 'average',
  POOR = 'poor',
  CRITICAL = 'critical',
}

export enum ContentType {
  ARTICLE = 'article',
  PRODUCT = 'product',
  CATEGORY = 'category',
  HOMEPAGE = 'homepage',
  LANDING = 'landing',
  BLOG_POST = 'blog-post',
}

export enum ReadabilityLevel {
  VERY_EASY = 'very-easy',      // Score 90-100
  EASY = 'easy',                // Score 80-89
  FAIRLY_EASY = 'fairly-easy',  // Score 70-79
  STANDARD = 'standard',        // Score 60-69
  FAIRLY_DIFFICULT = 'fairly-difficult', // Score 50-59
  DIFFICULT = 'difficult',      // Score 30-49
  VERY_DIFFICULT = 'very-difficult',     // Score 0-29
}

export enum ContentIssueType {
  THIN_CONTENT = 'thin-content',
  DUPLICATE_CONTENT = 'duplicate-content',
  MISSING_KEYWORDS = 'missing-keywords',
  KEYWORD_STUFFING = 'keyword-stuffing',
  LOW_READABILITY = 'low-readability',
  MISSING_HEADINGS = 'missing-headings',
  POOR_STRUCTURE = 'poor-structure',
  NO_INTERNAL_LINKS = 'no-internal-links',
  MISSING_IMAGES = 'missing-images',
  NO_CTA = 'no-cta',
  OUTDATED_CONTENT = 'outdated-content',
  BROKEN_LINKS = 'broken-links',
  MISSING_SCHEMA = 'missing-schema',
  POOR_UX = 'poor-ux',
}

export enum TopicImportance {
  HIGH = 'high',
  MEDIUM = 'medium',
  LOW = 'low',
}

export enum ContentSeverity {
  CRITICAL = 'critical',
  WARNING = 'warning', 
  INFO = 'info',
}

export enum CompetitorPosition {
  ABOVE = 'above',
  EQUAL = 'equal',
  BELOW = 'below',
}

// ==================== INTERFACES PRINCIPALES ====================

export interface ContentPage {
  url: string;
  title: string;
  contentType: ContentType;
  publishDate?: string;
  lastModified?: string;
  
  // Métriques quantitatives
  metrics: {
    wordCount: number;
    uniqueWords: number;
    sentenceCount: number;
    paragraphCount: number;
    readingTime: number; // en minutes
    avgWordsPerSentence: number;
    avgSentencesPerParagraph: number;
  };
  
  // Scores de qualité (0-100)
  quality: {
    overallScore: number;
    originalityScore: number;
    topicRelevance: number;
    keywordIntegration: number;
    readabilityScore: number;
    engagementPotential: number;
    expertiseLevel: number; // E-A-T Score
    freshness: number; // Content freshness
  };
  
  // Analyse de lisibilité détaillée
  readability: {
    level: ReadabilityLevel;
    fleschKincaidGrade: number;
    fleschReadingEase: number;
    gunningFogIndex: number;
    smogIndex: number;
    colemanLiauIndex: number;
    automatedReadabilityIndex: number;
    averageGradeLevel: number;
  };
  
  // Structure du contenu
  structure: {
    headings: {
      h1: number;
      h2: number;
      h3: number;
      h4: number;
      h5: number;
      h6: number;
    };
    lists: {
      ordered: number;
      unordered: number;
      totalItems: number;
    };
    images: {
      total: number;
      withAlt: number;
      withCaption: number;
      decorative: number;
    };
    links: {
      internal: number;
      external: number;
      outbound: number;
      anchor: number;
    };
    multimedia: {
      videos: number;
      audios: number;
      embeds: number;
      infographics: number;
    };
  };
  
  // Analyse des mots-clés
  keywords: {
    primary: Array<{
      keyword: string;
      density: number;
      frequency: number;
      positions: number[];
      inTitle: boolean;
      inHeadings: boolean;
      inMeta: boolean;
      prominence: number; // 0-100
    }>;
    secondary: Array<{
      keyword: string;
      density: number;
      frequency: number;
      relevanceScore: number;
    }>;
    lsi: Array<{
      term: string;
      relevance: number;
      frequency: number;
    }>; // Latent Semantic Indexing
    entities: Array<{
      entity: string;
      type: 'person' | 'place' | 'organization' | 'product' | 'concept';
      confidence: number;
      mentions: number;
    }>;
  };
  
  // Problèmes identifiés
  issues: Array<{
    type: ContentIssueType;
    severity: ContentSeverity;
    description: string;
    affectedElements: string[];
    impact: string;
    howToFix: HowToFix;
  }>;
  
  // Suggestions d'amélioration
  suggestions: {
    contentLength: {
      current: number;
      recommended: number;
      reason: string;
    };
    missingTopics: string[];
    keywordsToAdd: Array<{
      keyword: string;
      searchVolume: number;
      difficulty: number;
      opportunity: number;
    }>;
    structureImprovements: string[];
    readabilityEnhancements: string[];
  };
}

// Interface pour les actions correctives détaillées
export interface HowToFix {
  steps: string[];
  estimatedTime: string;
  difficulty: 'easy' | 'medium' | 'hard';
  codeExample?: string;
  resources: Array<{
    title: string;
    url: string;
    type: 'documentation' | 'tutorial' | 'tool' | 'guide';
  }>;
  beforeAfter?: {
    before: string;
    after: string;
  };
  priority: number; // 1-10
  potentialGain: {
    traffic: string;
    performance: string;
    ranking: string;
    engagement: string;
  };
}

// ==================== IA CONTENT ANALYSIS ====================

export interface AIContentInsights {
  // Analyse sémantique avancée
  topicsCovered: Array<{
    topic: string;
    coverage: number; // 0-100%
    depth: number; // 0-100%
    authority: number; // 0-100%
    searchIntent: 'informational' | 'navigational' | 'transactional' | 'commercial';
  }>;
  
  topicsMissing: Array<{
    topic: string;
    importance: TopicImportance;
    searchVolume: number;
    difficulty: number;
    opportunity: number;
    relatedKeywords: string[];
  }>;
  
  // Gaps de contenu identifiés par l'IA
  contentGaps: Array<{
    gap: string;
    category: 'knowledge' | 'depth' | 'format' | 'perspective';
    importance: TopicImportance;
    suggestedContent: string;
    estimatedLength: number;
    competitorsCovering: number;
    potentialTraffic: number;
  }>;
  
  // Analyse concurrentielle IA
  competitorComparison: {
    avgContentLength: number;
    ourAvgLength: number;
    lengthRecommendation: string;
    topicCoverageGap: number; // %
    qualityGap: number; // %
    topCompetitors: Array<{
      url: string;
      domain: string;
      contentLength: number;
      qualityScore: number;
      topicsCovered: string[];
      strengthsOverUs: string[];
    }>;
  };
  
  // Recommandations IA avancées
  aiRecommendations: Array<{
    type: 'content_expansion' | 'restructuring' | 'keyword_optimization' | 'format_change';
    priority: number;
    impact: 'high' | 'medium' | 'low';
    description: string;
    implementation: string;
    expectedResults: {
      trafficIncrease: string;
      rankingImprovement: string;
      engagementBoost: string;
    };
  }>;
  
  // Analyse de sentiment et ton
  contentTone: {
    sentiment: 'positive' | 'neutral' | 'negative';
    formality: number; // 0-100 (0 = très informel, 100 = très formel)
    complexity: number; // 0-100
    enthusiasm: number; // 0-100
    trustworthiness: number; // 0-100
    expertise: number; // 0-100 (E-A-T score)
    targetAudience: string[];
    toneConsistency: number; // 0-100
  };
}

// ==================== VISUALISATIONS AVANCÉES ====================

export interface ContentVisualizationData {
  // Données pour graphiques
  charts: {
    // Distribution des scores de contenu
    qualityDistribution: Array<{
      range: string; // '0-20', '21-40', etc.
      count: number;
      percentage: number;
    }>;
    
    // Timeline de performance du contenu
    contentPerformanceTimeline: Array<{
      date: string;
      avgQualityScore: number;
      totalWords: number;
      newContent: number;
      updatedContent: number;
    }>;
    
    // Répartition des types de contenu
    contentTypeDistribution: Array<{
      type: ContentType;
      count: number;
      avgQuality: number;
      avgLength: number;
    }>;
    
    // Performance par longueur de contenu
    lengthVsPerformance: Array<{
      lengthRange: string; // '0-500', '501-1000', etc.
      avgQualityScore: number;
      avgEngagement: number;
      count: number;
    }>;
  };
  
  // Données pour heatmaps
  heatmaps: {
    // Heatmap de densité des mots-clés
    keywordDensityMap: Array<{
      page: string;
      keyword: string;
      density: number;
      position: { x: number; y: number }; // Position dans la page
      performance: number; // Impact sur le ranking
    }>;
    
    // Heatmap de couverture des sujets
    topicCoverageMap: Array<{
      topic: string;
      pages: Array<{
        url: string;
        coverage: number; // 0-100
        depth: number; // 0-100
      }>;
    }>;
  };
  
  // Données pour graphes de réseau
  networkGraphs: {
    // Graphe du maillage interne
    internalLinkingGraph: {
      nodes: Array<{
        id: string;
        url: string;
        title: string;
        contentType: ContentType;
        qualityScore: number;
        inLinks: number;
        outLinks: number;
        pageRank: number;
      }>;
      edges: Array<{
        source: string;
        target: string;
        anchorText: string;
        strength: number; // Force du lien
        relevance: number; // Pertinence sémantique
      }>;
    };
    
    // Graphe des relations sémantiques
    semanticGraph: {
      nodes: Array<{
        id: string;
        term: string;
        type: 'keyword' | 'entity' | 'topic';
        importance: number;
        frequency: number;
      }>;
      edges: Array<{
        source: string;
        target: string;
        relation: 'synonym' | 'related' | 'broader' | 'narrower';
        strength: number;
      }>;
    };
  };
}

// ==================== INTERFACE PRINCIPALE ====================

export interface ContentAnalysis {
  // Pages analysées
  pages: ContentPage[];
  
  // Insights IA avancés
  aiContentAnalysis: AIContentInsights;
  
  // Données de visualisation
  visualizationData: ContentVisualizationData;
  
  // Métriques globales
  globalMetrics: {
    totalPages: number;
    totalWords: number;
    avgWordsPerPage: number;
    avgQualityScore: number;
    uniqueTopicsCovered: number;
    contentFreshness: number; // % de contenu récent
    duplicateContentPages: number;
    thinContentPages: number;
    highQualityPages: number;
    
    // Benchmarking concurrentiel
    competitorBenchmark: {
      position: CompetitorPosition;
      gapAnalysis: {
        contentVolume: number; // % de différence
        contentQuality: number;
        topicCoverage: number;
        keywordCoverage: number;
      };
    };
  };
  
  // Issues globales avec actions
  globalIssues: Array<{
    id: string;
    type: ContentIssueType;
    severity: ContentSeverity;
    affectedPages: string[];
    impact: string;
    description: string;
    howToFix: HowToFix;
    priority: number;
    estimatedROI: {
      trafficGain: number; // %
      rankingImprovement: number; // positions
      implementationCost: 'low' | 'medium' | 'high';
    };
  }>;
  
  // Recommandations stratégiques
  strategicRecommendations: Array<{
    id: string;
    category: 'content_strategy' | 'technical_seo' | 'user_experience' | 'competition';
    title: string;
    description: string;
    implementation: {
      phases: Array<{
        phase: string;
        duration: string;
        tasks: string[];
        resources: string[];
      }>;
      totalEffort: string;
      priority: number;
    };
    expectedResults: {
      shortTerm: string; // 1-3 mois
      mediumTerm: string; // 3-6 mois
      longTerm: string; // 6+ mois
    };
  }>;
  
  // Configuration et métadonnées
  config: {
    analysisDepth: 'basic' | 'advanced' | 'comprehensive';
    includeCompetitorAnalysis: boolean;
    targetKeywords: string[];
    competitorDomains: string[];
    contentLanguage: string;
    targetAudience: string[];
  };
  
  metadata: {
    analysisDate: string;
    analysisId: string;
    processingTime: number; // en secondes
    aiModelUsed: string;
    confidenceScore: number; // 0-100
    version: string;
    tool: string;
  };
}

// ==================== TYPES UTILITAIRES ====================

export interface ContentAnalysisTableRow extends ContentPage {
  id: string;
  selected?: boolean;
}

export enum ContentSortBy {
  QUALITY_SCORE = 'quality-score',
  WORD_COUNT = 'word-count',
  READABILITY = 'readability',
  KEYWORD_DENSITY = 'keyword-density',
  LAST_MODIFIED = 'last-modified',
  ISSUES_COUNT = 'issues-count',
  URL = 'url',
}

export enum ContentFilterBy {
  ALL = 'all',
  HIGH_QUALITY = 'high-quality',
  LOW_QUALITY = 'low-quality',
  THIN_CONTENT = 'thin-content',
  DUPLICATE_CONTENT = 'duplicate-content',
  MISSING_KEYWORDS = 'missing-keywords',
  CONTENT_TYPE = 'content-type',
}

// ==================== CONSTANTES ====================

export const CONTENT_QUALITY_THRESHOLDS = {
  excellent: { min: 90, max: 100 },
  good: { min: 75, max: 89 },
  average: { min: 60, max: 74 },
  poor: { min: 40, max: 59 },
  critical: { min: 0, max: 39 },
};

export const READABILITY_THRESHOLDS = {
  [ReadabilityLevel.VERY_EASY]: { min: 90, max: 100 },
  [ReadabilityLevel.EASY]: { min: 80, max: 89 },
  [ReadabilityLevel.FAIRLY_EASY]: { min: 70, max: 79 },
  [ReadabilityLevel.STANDARD]: { min: 60, max: 69 },
  [ReadabilityLevel.FAIRLY_DIFFICULT]: { min: 50, max: 59 },
  [ReadabilityLevel.DIFFICULT]: { min: 30, max: 49 },
  [ReadabilityLevel.VERY_DIFFICULT]: { min: 0, max: 29 },
};

export const CONTENT_LENGTH_RECOMMENDATIONS = {
  [ContentType.ARTICLE]: { min: 1500, ideal: 2500, max: 4000 },
  [ContentType.BLOG_POST]: { min: 800, ideal: 1500, max: 3000 },
  [ContentType.PRODUCT]: { min: 300, ideal: 500, max: 1000 },
  [ContentType.CATEGORY]: { min: 200, ideal: 400, max: 800 },
  [ContentType.HOMEPAGE]: { min: 500, ideal: 800, max: 1200 },
  [ContentType.LANDING]: { min: 600, ideal: 1000, max: 1500 },
};