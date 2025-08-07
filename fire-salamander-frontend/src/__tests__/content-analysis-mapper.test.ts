/**
 * Fire Salamander - Content Analysis Mapper Tests (TDD)
 * Lead Tech quality - Tests before mapper implementation
 */

import { describe, it, expect, beforeEach } from '@jest/globals';
import {
  ContentAnalysis,
  ContentType,
  ContentQuality,
  ReadabilityLevel,
  ContentIssueType,
  ContentSeverity,
  TopicImportance,
} from '@/types/content-analysis';

// Mock backend data structure pour Content Analysis
interface MockBackendContentData {
  id: string;
  url: string;
  analyzed_at: string;
  content_analysis: {
    pages: Array<{
      url: string;
      title: string;
      word_count: number;
      unique_words: number;
      reading_time: number;
      quality_score: number;
      readability_score: number;
      keyword_density: number;
      headings: {
        h1: number;
        h2: number;
        h3: number;
      };
      issues: Array<{
        type: string;
        severity: string;
        description: string;
      }>;
    }>;
    ai_insights: {
      topics_covered: Array<{
        topic: string;
        coverage: number;
      }>;
      missing_topics: Array<{
        topic: string;
        importance: string;
        search_volume: number;
      }>;
      competitor_analysis: {
        avg_content_length: number;
        quality_gap: number;
      };
    };
    global_metrics: {
      total_pages: number;
      avg_quality_score: number;
      thin_content_pages: number;
    };
  };
}

describe('Content Analysis Mapper Tests (TDD)', () => {
  let mockBackendData: MockBackendContentData;

  beforeEach(() => {
    mockBackendData = {
      id: 'content-analysis-123',
      url: 'https://example.com',
      analyzed_at: '2024-03-15T10:30:00Z',
      content_analysis: {
        pages: [
          {
            url: 'https://example.com/article/seo-guide',
            title: 'Complete SEO Guide 2024',
            word_count: 2500,
            unique_words: 1200,
            reading_time: 10,
            quality_score: 87,
            readability_score: 75,
            keyword_density: 2.1,
            headings: {
              h1: 1,
              h2: 8,
              h3: 15,
            },
            issues: [
              {
                type: 'missing_keywords',
                severity: 'warning',
                description: 'Missing important keywords: technical SEO, local SEO',
              },
            ],
          },
          {
            url: 'https://example.com/thin-content',
            title: 'Short Article',
            word_count: 150,
            unique_words: 120,
            reading_time: 1,
            quality_score: 35,
            readability_score: 60,
            keyword_density: 0.5,
            headings: {
              h1: 1,
              h2: 1,
              h3: 0,
            },
            issues: [
              {
                type: 'thin_content',
                severity: 'critical',
                description: 'Content is too short with only 150 words',
              },
            ],
          },
        ],
        ai_insights: {
          topics_covered: [
            {
              topic: 'SEO basics',
              coverage: 85,
            },
            {
              topic: 'On-page SEO',
              coverage: 70,
            },
          ],
          missing_topics: [
            {
              topic: 'Technical SEO',
              importance: 'high',
              search_volume: 15000,
            },
            {
              topic: 'Local SEO',
              importance: 'medium',
              search_volume: 8500,
            },
          ],
          competitor_analysis: {
            avg_content_length: 3200,
            quality_gap: -12,
          },
        },
        global_metrics: {
          total_pages: 2,
          avg_quality_score: 61,
          thin_content_pages: 1,
        },
      },
    };
  });

  describe('mapBackendToContentAnalysis function', () => {
    it('should exist and be callable', async () => {
      const { mapBackendToContentAnalysis } = await import('@/lib/mappers/content-analysis-mapper');
      expect(typeof mapBackendToContentAnalysis).toBe('function');
    });

    it('should map backend data to ContentAnalysis structure', async () => {
      const { mapBackendToContentAnalysis } = await import('@/lib/mappers/content-analysis-mapper');
      
      const result = mapBackendToContentAnalysis(mockBackendData);
      
      expect(result).toBeDefined();
      expect(result.pages).toHaveLength(2);
      expect(result.pages[0].url).toBe('https://example.com/article/seo-guide');
      expect(result.metadata.analysisId).toBe('content-analysis-123');
    });

    it('should correctly map content page metrics', async () => {
      const { mapBackendToContentAnalysis } = await import('@/lib/mappers/content-analysis-mapper');
      
      const result = mapBackendToContentAnalysis(mockBackendData);
      const firstPage = result.pages[0];
      
      expect(firstPage.metrics.wordCount).toBe(2500);
      expect(firstPage.metrics.uniqueWords).toBe(1200);
      expect(firstPage.metrics.readingTime).toBe(10);
      expect(firstPage.quality.overallScore).toBe(87);
    });

    it('should detect and categorize content types correctly', async () => {
      const { mapBackendToContentAnalysis } = await import('@/lib/mappers/content-analysis-mapper');
      
      const result = mapBackendToContentAnalysis(mockBackendData);
      const firstPage = result.pages[0];
      
      // Should detect article from URL pattern and length
      expect(firstPage.contentType).toBe(ContentType.ARTICLE);
      expect(firstPage.url).toContain('/article/');
      expect(firstPage.metrics.wordCount).toBeGreaterThan(1500);
    });

    it('should map readability scores correctly', async () => {
      const { mapBackendToContentAnalysis } = await import('@/lib/mappers/content-analysis-mapper');
      
      const result = mapBackendToContentAnalysis(mockBackendData);
      const firstPage = result.pages[0];
      
      expect(firstPage.readability.fleschReadingEase).toBe(75);
      expect(firstPage.readability.level).toBe(ReadabilityLevel.FAIRLY_EASY);
      expect(firstPage.quality.readabilityScore).toBe(75);
    });

    it('should identify and classify issues correctly', async () => {
      const { mapBackendToContentAnalysis } = await import('@/lib/mappers/content-analysis-mapper');
      
      const result = mapBackendToContentAnalysis(mockBackendData);
      
      // First page should have missing keywords issue
      const firstPageIssues = result.pages[0].issues;
      expect(firstPageIssues).toHaveLength(1);
      expect(firstPageIssues[0].type).toBe(ContentIssueType.MISSING_KEYWORDS);
      expect(firstPageIssues[0].severity).toBe(ContentSeverity.WARNING);
      
      // Second page should have thin content issue
      const secondPageIssues = result.pages[1].issues;
      expect(secondPageIssues).toHaveLength(1);
      expect(secondPageIssues[0].type).toBe(ContentIssueType.THIN_CONTENT);
      expect(secondPageIssues[0].severity).toBe(ContentSeverity.CRITICAL);
    });

    it('should map AI insights correctly', async () => {
      const { mapBackendToContentAnalysis } = await import('@/lib/mappers/content-analysis-mapper');
      
      const result = mapBackendToContentAnalysis(mockBackendData);
      const aiInsights = result.aiContentAnalysis;
      
      expect(aiInsights.topicsCovered).toHaveLength(2);
      expect(aiInsights.topicsCovered[0].topic).toBe('SEO basics');
      expect(aiInsights.topicsCovered[0].coverage).toBe(85);
      
      expect(aiInsights.topicsMissing).toHaveLength(2);
      expect(aiInsights.topicsMissing[0].topic).toBe('Technical SEO');
      expect(aiInsights.topicsMissing[0].importance).toBe(TopicImportance.HIGH);
    });

    it('should calculate global metrics correctly', async () => {
      const { mapBackendToContentAnalysis } = await import('@/lib/mappers/content-analysis-mapper');
      
      const result = mapBackendToContentAnalysis(mockBackendData);
      const globalMetrics = result.globalMetrics;
      
      expect(globalMetrics.totalPages).toBe(2);
      expect(globalMetrics.avgQualityScore).toBe(61);
      expect(globalMetrics.thinContentPages).toBe(1);
      expect(globalMetrics.totalWords).toBe(2650); // 2500 + 150
      expect(globalMetrics.avgWordsPerPage).toBe(1325); // 2650 / 2
    });

    it('should handle empty or invalid data gracefully', async () => {
      const { mapBackendToContentAnalysis } = await import('@/lib/mappers/content-analysis-mapper');
      
      const emptyData = {};
      const result = mapBackendToContentAnalysis(emptyData);
      
      expect(result.pages).toHaveLength(0);
      expect(result.globalMetrics.totalPages).toBe(0);
      expect(result.metadata.analysisId).toBeTruthy();
    });

    it('should generate visualization data', async () => {
      const { mapBackendToContentAnalysis } = await import('@/lib/mappers/content-analysis-mapper');
      
      const result = mapBackendToContentAnalysis(mockBackendData);
      const vizData = result.visualizationData;
      
      expect(vizData.charts.qualityDistribution).toBeDefined();
      expect(vizData.charts.contentTypeDistribution).toBeDefined();
      expect(vizData.heatmaps.keywordDensityMap).toBeDefined();
    });

    it('should create proper how-to-fix instructions', async () => {
      const { mapBackendToContentAnalysis } = await import('@/lib/mappers/content-analysis-mapper');
      
      const result = mapBackendToContentAnalysis(mockBackendData);
      const thinContentIssue = result.pages[1].issues[0];
      
      expect(thinContentIssue.howToFix.steps.length).toBeGreaterThan(0);
      expect(thinContentIssue.howToFix.estimatedTime).toBeTruthy();
      expect(['easy', 'medium', 'hard']).toContain(thinContentIssue.howToFix.difficulty);
      expect(thinContentIssue.howToFix.priority).toBeGreaterThanOrEqual(1);
      expect(thinContentIssue.howToFix.priority).toBeLessThanOrEqual(10);
    });
  });

  describe('Helper functions', () => {
    it('should have a function to classify content type', async () => {
      const { classifyContentType } = await import('@/lib/mappers/content-analysis-mapper');
      
      expect(classifyContentType('/article/seo-guide', 2500)).toBe(ContentType.ARTICLE);
      expect(classifyContentType('/blog/quick-tips', 800)).toBe(ContentType.BLOG_POST);
      expect(classifyContentType('/product/widget', 300)).toBe(ContentType.PRODUCT);
      expect(classifyContentType('/', 600)).toBe(ContentType.HOMEPAGE);
    });

    it('should have a function to calculate readability level', async () => {
      const { calculateReadabilityLevel } = await import('@/lib/mappers/content-analysis-mapper');
      
      expect(calculateReadabilityLevel(95)).toBe(ReadabilityLevel.VERY_EASY);
      expect(calculateReadabilityLevel(75)).toBe(ReadabilityLevel.FAIRLY_EASY);
      expect(calculateReadabilityLevel(65)).toBe(ReadabilityLevel.STANDARD);
      expect(calculateReadabilityLevel(45)).toBe(ReadabilityLevel.DIFFICULT);
      expect(calculateReadabilityLevel(25)).toBe(ReadabilityLevel.VERY_DIFFICULT);
    });

    it('should have a function to map issue types', async () => {
      const { mapIssueType } = await import('@/lib/mappers/content-analysis-mapper');
      
      expect(mapIssueType('thin_content')).toBe(ContentIssueType.THIN_CONTENT);
      expect(mapIssueType('missing_keywords')).toBe(ContentIssueType.MISSING_KEYWORDS);
      expect(mapIssueType('duplicate_content')).toBe(ContentIssueType.DUPLICATE_CONTENT);
      expect(mapIssueType('unknown_issue')).toBe(ContentIssueType.POOR_STRUCTURE); // fallback
    });

    it('should have a function to map severity levels', async () => {
      const { mapSeverityLevel } = await import('@/lib/mappers/content-analysis-mapper');
      
      expect(mapSeverityLevel('critical')).toBe(ContentSeverity.CRITICAL);
      expect(mapSeverityLevel('warning')).toBe(ContentSeverity.WARNING);
      expect(mapSeverityLevel('info')).toBe(ContentSeverity.INFO);
      expect(mapSeverityLevel('unknown')).toBe(ContentSeverity.INFO); // fallback
    });

    it('should calculate content quality scores', async () => {
      const { calculateContentQuality } = await import('@/lib/mappers/content-analysis-mapper');
      
      const mockPage = {
        word_count: 2500,
        readability_score: 75,
        keyword_density: 2.1,
        headings: { h1: 1, h2: 8, h3: 15 },
        issues_count: 1,
      };
      
      const quality = calculateContentQuality(mockPage);
      
      expect(quality.overallScore).toBeGreaterThanOrEqual(0);
      expect(quality.overallScore).toBeLessThanOrEqual(100);
      expect(quality.readabilityScore).toBe(75);
    });

    it('should generate content suggestions', async () => {
      const { generateContentSuggestions } = await import('@/lib/mappers/content-analysis-mapper');
      
      const mockPage = {
        word_count: 800,
        content_type: ContentType.ARTICLE,
        missing_topics: ['technical SEO', 'local SEO'],
      };
      
      const suggestions = generateContentSuggestions(mockPage);
      
      expect(suggestions.contentLength.recommended).toBeGreaterThan(800);
      expect(suggestions.missingTopics).toContain('technical SEO');
      expect(suggestions.keywordsToAdd.length).toBeGreaterThan(0);
    });

    it('should create visualization data from pages', async () => {
      const { createVisualizationData } = await import('@/lib/mappers/content-analysis-mapper');
      
      const mockPages = [
        { quality: { overallScore: 85 }, contentType: ContentType.ARTICLE },
        { quality: { overallScore: 65 }, contentType: ContentType.BLOG_POST },
      ];
      
      const vizData = createVisualizationData(mockPages as any);
      
      expect(vizData.charts.qualityDistribution).toBeDefined();
      expect(vizData.charts.contentTypeDistribution).toBeDefined();
    });

    it('should create empty content analysis for fallback', async () => {
      const { createEmptyContentAnalysis } = await import('@/lib/mappers/content-analysis-mapper');
      
      const emptyAnalysis = createEmptyContentAnalysis();
      
      expect(emptyAnalysis.pages).toHaveLength(0);
      expect(emptyAnalysis.globalMetrics.totalPages).toBe(0);
      expect(emptyAnalysis.metadata.analysisId).toBeTruthy();
    });
  });

  describe('Content Analysis Calculations', () => {
    it('should calculate accurate word count metrics', async () => {
      const { calculateWordCountMetrics } = await import('@/lib/mappers/content-analysis-mapper');
      
      const content = 'This is a test content. It has multiple sentences. Some words repeat.';
      const metrics = calculateWordCountMetrics(content);
      
      expect(metrics.wordCount).toBe(13);
      expect(metrics.sentenceCount).toBe(3);
      expect(metrics.uniqueWords).toBeLessThanOrEqual(metrics.wordCount);
    });

    it('should assess content structure quality', async () => {
      const { assessContentStructure } = await import('@/lib/mappers/content-analysis-mapper');
      
      const structure = {
        headings: { h1: 1, h2: 5, h3: 10, h4: 2, h5: 0, h6: 0 },
        images: { total: 8, withAlt: 6, withCaption: 3, decorative: 2 },
        links: { internal: 15, external: 5, outbound: 3, anchor: 2 },
      };
      
      const assessment = assessContentStructure(structure as any);
      
      expect(assessment.structureScore).toBeGreaterThanOrEqual(0);
      expect(assessment.structureScore).toBeLessThanOrEqual(100);
      expect(assessment.hasGoodHeadingHierarchy).toBe(true); // H1:1, H2:5, H3:10
    });

    it('should identify content gaps correctly', async () => {
      const { identifyContentGaps } = await import('@/lib/mappers/content-analysis-mapper');
      
      const topicsCovered = ['SEO basics', 'keyword research'];
      const targetTopics = ['SEO basics', 'technical SEO', 'local SEO', 'link building'];
      
      const gaps = identifyContentGaps(topicsCovered, targetTopics);
      
      expect(gaps).toContain('technical SEO');
      expect(gaps).toContain('local SEO');
      expect(gaps).toContain('link building');
      expect(gaps).not.toContain('SEO basics');
    });
  });
});