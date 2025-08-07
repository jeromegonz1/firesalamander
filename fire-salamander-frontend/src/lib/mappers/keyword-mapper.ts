/**
 * Fire Salamander - Keyword Data Mapper
 * Maps backend data to KeywordAnalysis interface
 */

import { KeywordAnalysis } from '@/types/keyword-analysis';
import { AIKeywordService } from '@/lib/services/ai-keyword-service';

export async function mapBackendToKeywordAnalysis(
  analysis: any,
  includeAI: boolean = true
): Promise<KeywordAnalysis> {
  const resultData = JSON.parse(analysis.result_data || '{}');
  const seoAnalysis = resultData.seo_analysis || {};
  const content = seoAnalysis.content || '';
  
  const aiService = new AIKeywordService();

  // Extract keywords from various locations
  const foundKeywords = extractKeywordsFromAnalysis(seoAnalysis);
  
  // Generate n-grams
  const bigrams = aiService.analyzeNGrams(content, 2);
  const trigrams = aiService.analyzeNGrams(content, 3);
  
  // Get AI suggestions if enabled
  let aiSuggestions: KeywordAnalysis['aiSuggestions'] = [];
  if (includeAI) {
    try {
      aiSuggestions = await aiService.generateKeywordSuggestions({
        analysisId: analysis.id,
        context: content.substring(0, 1000), // First 1000 chars
        currentKeywords: foundKeywords.map(k => k.keyword).slice(0, 10),
        targetAudience: 'Professionnels du droit et entreprises',
        contentType: 'website'
      });
    } catch (error) {
      console.error('Failed to get AI suggestions:', error);
    }
  }

  // Semantic analysis
  const semanticAnalysis = performSemanticAnalysis(content, seoAnalysis);

  // Calculate metrics
  const metrics = calculateKeywordMetrics(foundKeywords, content);

  return {
    foundKeywords,
    ngrams: {
      bigrams,
      trigrams,
      topPhrases: [...bigrams.slice(0, 5), ...trigrams.slice(0, 5)]
        .sort((a, b) => b.count - a.count)
        .slice(0, 8)
        .map(item => ({
          phrase: item.phrase,
          count: item.count,
          type: item.phrase.split(' ').length === 2 ? 'bigram' as const : 'trigram' as const
        }))
    },
    aiSuggestions,
    semanticAnalysis,
    competitorGaps: generateMockCompetitorGaps(), // TODO: Implement real competitor analysis
    metrics,
    status: {
      lastUpdated: new Date().toISOString(),
      aiAnalysisEnabled: includeAI,
      competitorAnalysisEnabled: false,
      crawledPages: 1
    }
  };
}

/**
 * Extract keywords from various parts of the analysis
 */
function extractKeywordsFromAnalysis(seoAnalysis: any): KeywordAnalysis['foundKeywords'] {
  const keywords = new Map<string, {
    occurrences: number;
    locations: Set<KeywordAnalysis['foundKeywords'][0]['locations'][0]>;
  }>();

  // Helper to add keywords
  const addKeywords = (
    text: string, 
    location: KeywordAnalysis['foundKeywords'][0]['locations'][0],
    weight: number = 1
  ) => {
    if (!text) return;
    
    // Extract words (3+ characters, no stopwords)
    const words = text.toLowerCase()
      .replace(/[^\w\sàâäéèêëïîôùûü]/g, ' ')
      .split(/\s+/)
      .filter(word => word.length > 3 && !isStopWord(word));
    
    words.forEach(word => {
      const existing = keywords.get(word) || { occurrences: 0, locations: new Set() };
      existing.occurrences += weight;
      existing.locations.add(location);
      keywords.set(word, existing);
    });
  };

  // Extract from title (high weight)
  if (seoAnalysis.tag_analysis?.title?.content) {
    addKeywords(seoAnalysis.tag_analysis.title.content, 'title', 3);
  }

  // Extract from meta description
  if (seoAnalysis.tag_analysis?.meta_description?.content) {
    addKeywords(seoAnalysis.tag_analysis.meta_description.content, 'meta', 2);
  }

  // Extract from headings
  if (seoAnalysis.tag_analysis?.headings?.h1_content) {
    seoAnalysis.tag_analysis.headings.h1_content.forEach((h1: string) => {
      addKeywords(h1, 'h1', 3);
    });
  }

  // Extract from content (if available)
  if (seoAnalysis.content) {
    addKeywords(seoAnalysis.content, 'content', 1);
  }

  // Extract from image alt texts
  if (seoAnalysis.tag_analysis?.images?.alt_texts) {
    seoAnalysis.tag_analysis.images.alt_texts.forEach((alt: string) => {
      addKeywords(alt, 'alt', 1);
    });
  }

  // Calculate total words for density
  const totalWords = seoAnalysis.content?.split(/\s+/).length || 100;

  // Convert to array and calculate metrics
  const keywordArray = Array.from(keywords.entries())
    .map(([keyword, data]) => {
      const density = (data.occurrences / totalWords) * 100;
      const prominence = calculateProminence(keyword, data.locations, seoAnalysis);
      
      return {
        keyword,
        density: parseFloat(density.toFixed(2)),
        occurrences: data.occurrences,
        prominence,
        locations: Array.from(data.locations)
      };
    })
    .filter(k => k.occurrences > 1) // Filter single occurrences
    .sort((a, b) => b.prominence - a.prominence)
    .slice(0, 50); // Top 50 keywords

  return keywordArray;
}

/**
 * Calculate keyword prominence based on location
 */
function calculateProminence(
  keyword: string, 
  locations: Set<string>, 
  seoAnalysis: any
): number {
  let score = 0;
  
  // Location weights
  const weights = {
    title: 40,
    h1: 30,
    h2: 20,
    h3: 15,
    meta: 25,
    content: 10,
    alt: 5,
    url: 35
  };

  locations.forEach(location => {
    score += weights[location as keyof typeof weights] || 0;
  });

  // Bonus for appearing in multiple locations
  if (locations.size > 1) {
    score += locations.size * 5;
  }

  return Math.min(100, score);
}

/**
 * Perform semantic analysis on content
 */
function performSemanticAnalysis(content: string, seoAnalysis: any): KeywordAnalysis['semanticAnalysis'] {
  const aiService = new AIKeywordService();
  
  // Extract main topics (simplified - in production use NLP)
  const mainTopics = extractMainTopics(content);
  
  // Extract entities
  const entities = aiService.extractEntities(content);
  
  // Calculate readability
  const readability = aiService.calculateReadability(content);
  
  // Determine sentiment (simplified)
  const sentiment = analyzeSentiment(content);
  
  return {
    mainTopics,
    entities,
    sentiment: sentiment.type,
    sentimentScore: sentiment.score,
    readability,
    language: 'fr', // TODO: Detect language
    wordCount: content.split(/\s+/).filter(w => w.length > 0).length
  };
}

/**
 * Extract main topics from content
 */
function extractMainTopics(content: string): string[] {
  // Simplified topic extraction
  const topics = new Set<string>();
  
  // Common topic patterns for Septeo
  const topicPatterns = [
    { pattern: /solutions?\s+logicielles?/gi, topic: 'Solutions logicielles' },
    { pattern: /transformation\s+digitale/gi, topic: 'Transformation digitale' },
    { pattern: /notaires?/gi, topic: 'Notaires' },
    { pattern: /avocats?/gi, topic: 'Avocats' },
    { pattern: /experts?\s+comptables?/gi, topic: 'Experts comptables' },
    { pattern: /immobilier/gi, topic: 'Immobilier' },
    { pattern: /gestion/gi, topic: 'Gestion' },
    { pattern: /juridique/gi, topic: 'Juridique' },
    { pattern: /cloud/gi, topic: 'Cloud' },
    { pattern: /sécurité/gi, topic: 'Sécurité' }
  ];

  topicPatterns.forEach(({ pattern, topic }) => {
    if (pattern.test(content)) {
      topics.add(topic);
    }
  });

  return Array.from(topics).slice(0, 8);
}

/**
 * Analyze sentiment of content
 */
function analyzeSentiment(content: string): { type: KeywordAnalysis['semanticAnalysis']['sentiment']; score: number } {
  // Simplified sentiment analysis
  const positiveWords = /excellent|optimal|performant|innovant|leader|réussite|efficace|moderne|sécurisé/gi;
  const negativeWords = /problème|erreur|lent|obsolète|complexe|difficile|risque|échec/gi;
  
  const positiveMatches = (content.match(positiveWords) || []).length;
  const negativeMatches = (content.match(negativeWords) || []).length;
  
  const total = positiveMatches + negativeMatches;
  if (total === 0) return { type: 'neutral', score: 0 };
  
  const score = (positiveMatches - negativeMatches) / total;
  
  if (score > 0.3) return { type: 'positive', score };
  if (score < -0.3) return { type: 'negative', score };
  if (positiveMatches > 0 && negativeMatches > 0) return { type: 'mixed', score };
  return { type: 'neutral', score };
}

/**
 * Calculate keyword metrics
 */
function calculateKeywordMetrics(
  keywords: KeywordAnalysis['foundKeywords'],
  content: string
): KeywordAnalysis['metrics'] {
  const uniqueKeywords = new Set(keywords.map(k => k.keyword));
  const avgDensity = keywords.reduce((sum, k) => sum + k.density, 0) / keywords.length || 0;
  
  // Keyword diversity score (unique/total ratio * 100)
  const totalOccurrences = keywords.reduce((sum, k) => sum + k.occurrences, 0);
  const diversityScore = (uniqueKeywords.size / totalOccurrences) * 100;
  
  // Focus keywords (top 5 by prominence)
  const focusKeywords = keywords
    .sort((a, b) => b.prominence - a.prominence)
    .slice(0, 5)
    .map(k => k.keyword);

  return {
    totalKeywords: keywords.length,
    uniqueKeywords: uniqueKeywords.size,
    avgKeywordDensity: parseFloat(avgDensity.toFixed(2)),
    keywordDiversity: Math.min(100, Math.round(diversityScore)),
    focusKeywords,
    missingImportantKeywords: [] // TODO: Compare with industry standards
  };
}

/**
 * Check if word is a stop word
 */
function isStopWord(word: string): boolean {
  const stopWords = new Set([
    'le', 'la', 'les', 'un', 'une', 'des', 'de', 'du', 'et', 'ou', 'mais',
    'pour', 'dans', 'sur', 'avec', 'sans', 'sous', 'vers', 'chez', 'par',
    'ce', 'ces', 'cette', 'cet', 'celui', 'celle', 'ceux', 'celles',
    'qui', 'que', 'quoi', 'dont', 'où', 'si', 'ne', 'pas', 'plus', 'très',
    'tout', 'toute', 'tous', 'toutes', 'être', 'avoir', 'faire'
  ]);
  
  return stopWords.has(word.toLowerCase());
}

/**
 * Generate mock competitor gaps (TODO: Implement real analysis)
 */
function generateMockCompetitorGaps(): KeywordAnalysis['competitorGaps'] {
  return [
    {
      keyword: "logiciel juridique cloud",
      competitorRank: 3,
      ourRank: null,
      searchVolume: 1900,
      difficulty: 42,
      opportunity: 'high',
      estimatedTraffic: 285,
      competitors: [
        { domain: 'competitor1.com', rank: 1 },
        { domain: 'competitor2.com', rank: 2 },
        { domain: 'competitor3.com', rank: 3 }
      ]
    },
    {
      keyword: "gestion cabinet notarial",
      competitorRank: 5,
      ourRank: 15,
      searchVolume: 720,
      difficulty: 35,
      opportunity: 'medium',
      estimatedTraffic: 108,
      competitors: [
        { domain: 'leader-market.com', rank: 1 },
        { domain: 'solution-pro.fr', rank: 2 }
      ]
    }
  ];
}