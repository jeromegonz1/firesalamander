/**
 * Fire Salamander - AI Keyword Service
 * Integration with ChatGPT 3.5 for keyword suggestions
 */

import { KeywordAnalysis, AIKeywordGenerationRequest } from '@/types/keyword-analysis';

const OPENAI_API_URL = 'https://api.openai.com/v1/chat/completions';

export class AIKeywordService {
  private apiKey: string;

  constructor(apiKey?: string) {
    // Get from environment or use server-side proxy
    this.apiKey = apiKey || process.env.NEXT_PUBLIC_OPENAI_API_KEY || '';
  }

  /**
   * Generate keyword suggestions using ChatGPT 3.5
   */
  async generateKeywordSuggestions(
    request: AIKeywordGenerationRequest
  ): Promise<KeywordAnalysis['aiSuggestions']> {
    try {
      // Call backend endpoint that handles OpenAI integration
      const response = await fetch(`/api/v1/analysis/${request.analysisId}/keywords/ai-suggestions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          context: request.context,
          currentKeywords: request.currentKeywords,
          targetAudience: request.targetAudience,
          contentType: request.contentType,
          limit: request.limit || 10
        })
      });

      if (!response.ok) {
        throw new Error('Failed to generate AI suggestions');
      }

      const data = await response.json();
      return data.suggestions;
    } catch (error) {
      console.error('AI Keyword generation error:', error);
      
      // Fallback to mock data in development
      return this.getMockSuggestions(request);
    }
  }

  /**
   * Analyze content with n-grams
   */
  analyzeNGrams(content: string, n: number = 2): Array<{ phrase: string; count: number; density: number }> {
    // Clean and tokenize content
    const words = content.toLowerCase()
      .replace(/[^\w\s]/g, ' ')
      .split(/\s+/)
      .filter(word => word.length > 2);

    const ngrams = new Map<string, number>();
    
    // Generate n-grams
    for (let i = 0; i <= words.length - n; i++) {
      const ngram = words.slice(i, i + n).join(' ');
      ngrams.set(ngram, (ngrams.get(ngram) || 0) + 1);
    }

    // Convert to array and calculate density
    const totalNgrams = words.length - n + 1;
    const results = Array.from(ngrams.entries())
      .map(([phrase, count]) => ({
        phrase,
        count,
        density: (count / totalNgrams) * 100
      }))
      .filter(item => item.count > 1) // Filter out single occurrences
      .sort((a, b) => b.count - a.count)
      .slice(0, 20); // Top 20

    return results;
  }

  /**
   * Extract entities from content
   */
  extractEntities(content: string): KeywordAnalysis['semanticAnalysis']['entities'] {
    // Simple entity extraction (in production, use NLP service)
    const entities: KeywordAnalysis['semanticAnalysis']['entities'] = [];
    
    // Common patterns for entities
    const patterns = {
      organization: /\b([A-Z][a-z]+(?:\s+[A-Z][a-z]+)*)\s+(?:SA|SAS|SARL|Inc|LLC|Ltd|GmbH|AG)\b/g,
      person: /\b(?:M\.|Mme|Dr\.|Prof\.)\s+([A-Z][a-z]+(?:\s+[A-Z][a-z]+)*)\b/g,
      location: /\b(?:à|de|en)\s+([A-Z][a-z]+(?:\s+[A-Z][a-z]+)*)\b/g,
    };

    // Extract organizations
    let match;
    while ((match = patterns.organization.exec(content)) !== null) {
      entities.push({
        name: match[1],
        type: 'Organization',
        salience: 0.8,
        mentions: 1
      });
    }

    return entities;
  }

  /**
   * Calculate readability metrics
   */
  calculateReadability(content: string): KeywordAnalysis['semanticAnalysis']['readability'] {
    const sentences = content.split(/[.!?]+/).filter(s => s.trim().length > 0);
    const words = content.split(/\s+/).filter(w => w.length > 0);
    const syllables = words.reduce((total, word) => total + this.countSyllables(word), 0);
    
    // Flesch Reading Ease Score
    const avgSentenceLength = words.length / sentences.length;
    const avgSyllablesPerWord = syllables / words.length;
    const fleschScore = 206.835 - 1.015 * avgSentenceLength - 84.6 * avgSyllablesPerWord;
    
    // Determine level
    let level: KeywordAnalysis['semanticAnalysis']['readability']['level'];
    if (fleschScore >= 90) level = 'very-easy';
    else if (fleschScore >= 80) level = 'easy';
    else if (fleschScore >= 70) level = 'fairly-easy';
    else if (fleschScore >= 60) level = 'standard';
    else if (fleschScore >= 50) level = 'fairly-difficult';
    else if (fleschScore >= 30) level = 'difficult';
    else level = 'very-difficult';

    return {
      score: Math.max(0, Math.min(100, fleschScore)),
      level,
      avgSentenceLength,
      avgWordLength: words.reduce((sum, word) => sum + word.length, 0) / words.length,
      avgWordsPerSentence: avgSentenceLength,
      complexWords: (words.filter(w => this.countSyllables(w) > 2).length / words.length) * 100,
      paragraphs: content.split(/\n\n+/).length
    };
  }

  /**
   * Count syllables in a word (French approximation)
   */
  private countSyllables(word: string): number {
    word = word.toLowerCase();
    let count = 0;
    let previousWasVowel = false;
    
    for (let i = 0; i < word.length; i++) {
      const isVowel = /[aeiouàâäéèêëïîôùûü]/.test(word[i]);
      if (isVowel && !previousWasVowel) {
        count++;
      }
      previousWasVowel = isVowel;
    }
    
    // French specific rules
    if (word.endsWith('e') && count > 1) count--;
    if (word.endsWith('es') && count > 1) count--;
    
    return Math.max(1, count);
  }

  /**
   * Mock data for development
   */
  private getMockSuggestions(request: AIKeywordGenerationRequest): KeywordAnalysis['aiSuggestions'] {
    return [
      {
        keyword: "solutions logicielles métiers",
        searchVolume: 2400,
        difficulty: 45,
        cpc: 3.2,
        intent: 'commercial',
        reason: "Forte pertinence avec votre activité principale et volume de recherche élevé",
        contentIdeas: [
          "Guide complet des solutions logicielles pour entreprises",
          "Comparatif des meilleures solutions métiers 2025",
          "Comment choisir sa solution logicielle métier"
        ],
        trendData: {
          direction: 'up',
          changePercent: 15
        },
        relatedKeywords: ["logiciel métier", "software professionnel", "ERP métier"]
      },
      {
        keyword: "transformation digitale juridique",
        searchVolume: 890,
        difficulty: 35,
        cpc: 4.5,
        intent: 'informational',
        reason: "Niche spécifique avec faible concurrence et CPC élevé",
        contentIdeas: [
          "La transformation digitale des cabinets d'avocats en 2025",
          "Étude de cas : digitalisation réussie d'un cabinet juridique",
          "Les outils essentiels pour la transformation digitale juridique"
        ],
        trendData: {
          direction: 'up',
          changePercent: 32
        }
      },
      {
        keyword: "logiciel gestion notaire",
        searchVolume: 1600,
        difficulty: 28,
        cpc: 5.8,
        intent: 'transactional',
        reason: "Intention d'achat élevée avec excellent CPC pour votre marché cible",
        contentIdeas: [
          "Démonstration interactive de votre solution notariale",
          "ROI d'un logiciel de gestion pour notaire",
          "Migration vers une solution moderne : guide étape par étape"
        ],
        trendData: {
          direction: 'stable',
          changePercent: 2
        }
      }
    ];
  }
}