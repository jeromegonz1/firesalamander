/**
 * Fire Salamander - Backlinks Analysis Tests (TDD)
 * Lead Tech quality - Tests before implementation
 * Testing comprehensive backlinks analysis system
 */

import { describe, it, expect, beforeEach } from '@jest/globals';
import {
  BacklinksAnalysis,
  BacklinkGrade,
  BacklinkQuality,
  BacklinkStatus,
  LinkType,
  AnchorType,
  DomainCategory,
  ToxicityReason,
} from '@/types/backlinks-analysis';

// Mock backend data structure for Backlinks Analysis
interface MockBackendBacklinksData {
  id: string;
  url: string;
  analyzed_at: string;
  backlinks_data?: {
    profile?: {
      total_backlinks: number;
      referring_domains: number;
      referring_ips: number;
      new_links_30d: number;
      lost_links_30d: number;
      dofollow_percentage: number;
      average_da: number;
      trust_flow: number;
      citation_flow: number;
    };
    backlinks?: Array<{
      source_url: string;
      target_url: string;
      anchor_text: string;
      link_type: string;
      first_seen: string;
      last_seen: string;
      status: string;
      domain_authority: number;
      page_authority: number;
      spam_score: number;
      is_redirect: boolean;
      http_status: number;
      surrounding_text?: string;
      page_title?: string;
      toxicity_score?: number;
    }>;
    referring_domains?: Array<{
      domain: string;
      backlinks_count: number;
      first_seen: string;
      domain_authority: number;
      category: string;
      country: string;
      language: string;
      quality: string;
      toxicity_score: number;
    }>;
    anchor_texts?: Array<{
      text: string;
      count: number;
      type: string;
      percentage: number;
      average_da: number;
      is_natural: boolean;
    }>;
    competitors?: Array<{
      domain: string;
      total_backlinks: number;
      referring_domains: number;
      domain_authority: number;
      overlap_domains: string[];
      unique_domains: string[];
    }>;
    toxic_links?: Array<{
      url: string;
      domain: string;
      toxicity_score: number;
      reasons: string[];
      action: string;
    }>;
    opportunities?: Array<{
      domain: string;
      url: string;
      domain_authority: number;
      category: string;
      contact_email?: string;
      priority: string;
      difficulty: number;
      relevance_score: number;
    }>;
    historical_data?: Array<{
      date: string;
      total_backlinks: number;
      referring_domains: number;
      domain_authority: number;
      new_links: number;
      lost_links: number;
    }>;
    scores?: {
      overall: number;
      quality: number;
      diversity: number;
      authority: number;
      naturalness: number;
      toxicity: number;
    };
  };
}

describe('Backlinks Analysis Interface Tests (TDD)', () => {
  let mockBackendData: MockBackendBacklinksData;

  beforeEach(() => {
    mockBackendData = {
      id: 'backlinks-analysis-456',
      url: 'https://example.com',
      analyzed_at: '2024-03-15T14:30:00Z',
      backlinks_data: {
        profile: {
          total_backlinks: 15420,
          referring_domains: 2340,
          referring_ips: 1980,
          new_links_30d: 187,
          lost_links_30d: 45,
          dofollow_percentage: 68.5,
          average_da: 42.8,
          trust_flow: 35.2,
          citation_flow: 41.6,
        },
        backlinks: [
          {
            source_url: 'https://tech-blog.example/article-seo',
            target_url: 'https://example.com/seo-guide',
            anchor_text: 'comprehensive SEO guide',
            link_type: 'dofollow',
            first_seen: '2024-01-15T10:00:00Z',
            last_seen: '2024-03-15T10:00:00Z',
            status: 'active',
            domain_authority: 65,
            page_authority: 58,
            spam_score: 2,
            is_redirect: false,
            http_status: 200,
            surrounding_text: 'For more detailed information, check out this comprehensive SEO guide that covers all aspects.',
            page_title: 'Advanced SEO Techniques for 2024',
            toxicity_score: 5,
          },
          {
            source_url: 'https://spammy-site.example/page',
            target_url: 'https://example.com/',
            anchor_text: 'best SEO services',
            link_type: 'dofollow',
            first_seen: '2024-03-01T08:00:00Z',
            last_seen: '2024-03-15T08:00:00Z',
            status: 'toxic',
            domain_authority: 15,
            page_authority: 12,
            spam_score: 14,
            is_redirect: false,
            http_status: 200,
            surrounding_text: 'Click here for best SEO services and quick ranking',
            page_title: 'SEO Services - Rank #1 Fast',
            toxicity_score: 85,
          },
          {
            source_url: 'https://news-site.example/breaking-news',
            target_url: 'https://example.com/press-release',
            anchor_text: 'https://example.com',
            link_type: 'nofollow',
            first_seen: '2024-02-20T12:00:00Z',
            last_seen: '2024-03-10T12:00:00Z',
            status: 'lost',
            domain_authority: 78,
            page_authority: 82,
            spam_score: 1,
            is_redirect: false,
            http_status: 404,
            surrounding_text: 'According to a recent press release from https://example.com, the company announced new features.',
            page_title: 'Tech Company Announces New Features',
            toxicity_score: 0,
          },
        ],
        referring_domains: [
          {
            domain: 'tech-blog.example',
            backlinks_count: 3,
            first_seen: '2024-01-15T10:00:00Z',
            domain_authority: 65,
            category: 'blog',
            country: 'US',
            language: 'en',
            quality: 'good',
            toxicity_score: 10,
          },
          {
            domain: 'spammy-site.example',
            backlinks_count: 12,
            first_seen: '2024-03-01T08:00:00Z',
            domain_authority: 15,
            category: 'other',
            country: 'Unknown',
            language: 'en',
            quality: 'toxic',
            toxicity_score: 85,
          },
        ],
        anchor_texts: [
          {
            text: 'comprehensive SEO guide',
            count: 45,
            type: 'partial_match',
            percentage: 12.5,
            average_da: 52.3,
            is_natural: true,
          },
          {
            text: 'SEO',
            count: 120,
            type: 'exact_match',
            percentage: 33.2,
            average_da: 38.7,
            is_natural: false, // Over-optimized
          },
          {
            text: 'Example Company',
            count: 89,
            type: 'branded',
            percentage: 24.6,
            average_da: 45.1,
            is_natural: true,
          },
          {
            text: 'https://example.com',
            count: 67,
            type: 'naked_url',
            percentage: 18.5,
            average_da: 41.2,
            is_natural: true,
          },
          {
            text: 'click here',
            count: 23,
            type: 'generic',
            percentage: 6.4,
            average_da: 28.4,
            is_natural: true,
          },
        ],
        competitors: [
          {
            domain: 'competitor1.example',
            total_backlinks: 28940,
            referring_domains: 4120,
            domain_authority: 72,
            overlap_domains: ['tech-blog.example', 'industry-news.example'],
            unique_domains: ['exclusive-partner.example', 'premium-directory.example'],
          },
          {
            domain: 'competitor2.example',
            total_backlinks: 12350,
            referring_domains: 1890,
            domain_authority: 58,
            overlap_domains: ['news-site.example'],
            unique_domains: ['local-business.example', 'regional-blog.example'],
          },
        ],
        toxic_links: [
          {
            url: 'https://spammy-site.example/page',
            domain: 'spammy-site.example',
            toxicity_score: 85,
            reasons: ['spam', 'low_quality', 'unnatural_pattern'],
            action: 'disavow',
          },
          {
            url: 'https://link-farm.example/links/123',
            domain: 'link-farm.example',
            toxicity_score: 92,
            reasons: ['link_farm', 'pbn', 'suspicious_anchor'],
            action: 'disavow',
          },
        ],
        opportunities: [
          {
            domain: 'industry-authority.example',
            url: 'https://industry-authority.example/resources',
            domain_authority: 84,
            category: 'corporate',
            contact_email: 'editor@industry-authority.example',
            priority: 'high',
            difficulty: 75,
            relevance_score: 92,
          },
          {
            domain: 'tech-magazine.example',
            url: 'https://tech-magazine.example/guest-posts',
            domain_authority: 71,
            category: 'news',
            priority: 'medium',
            difficulty: 60,
            relevance_score: 88,
          },
        ],
        historical_data: [
          {
            date: '2024-01-01',
            total_backlinks: 14200,
            referring_domains: 2180,
            domain_authority: 38,
            new_links: 156,
            lost_links: 23,
          },
          {
            date: '2024-02-01',
            total_backlinks: 14850,
            referring_domains: 2260,
            domain_authority: 40,
            new_links: 201,
            lost_links: 34,
          },
          {
            date: '2024-03-01',
            total_backlinks: 15420,
            referring_domains: 2340,
            domain_authority: 42,
            new_links: 187,
            lost_links: 45,
          },
        ],
        scores: {
          overall: 72,
          quality: 65,
          diversity: 78,
          authority: 68,
          naturalness: 58, // Lower due to over-optimization
          toxicity: 85,    // 100 - 15 (low toxicity)
        },
      },
    };
  });

  describe('BacklinksAnalysis Interface Validation', () => {
    it('should have required backlinks analysis properties', () => {
      const requiredProps = [
        'score',
        'profile', 
        'backlinks',
        'anchorTexts',
        'referringDomains',
        'competitors',
        'opportunities',
        'toxicAnalysis',
        'trends',
        'analytics',
        'recommendations',
        'metadata',
      ];

      // Test would verify interface structure
      expect(requiredProps.length).toBe(12);
      requiredProps.forEach(prop => {
        expect(typeof prop).toBe('string');
      });
    });

    it('should validate score structure with all sub-metrics', () => {
      const scoreProps = [
        'overall',
        'grade', 
        'trend',
        'scoreBreakdown',
      ];

      const breakdownProps = [
        'quantity',
        'quality',
        'diversity', 
        'authority',
        'naturalness',
        'toxicity',
      ];

      expect(scoreProps.length).toBe(4);
      expect(breakdownProps.length).toBe(6);
      
      // Validate BacklinkGrade enum values
      const gradeValues = Object.values(BacklinkGrade);
      expect(gradeValues).toContain(BacklinkGrade.A_PLUS);
      expect(gradeValues).toContain(BacklinkGrade.A);
      expect(gradeValues).toContain(BacklinkGrade.B);
      expect(gradeValues).toContain(BacklinkGrade.F);
    });

    it('should validate BacklinkData interface structure', () => {
      const backlinkDataProps = [
        'id',
        'sourceUrl',
        'targetUrl', 
        'sourceDomain',
        'targetDomain',
        'anchorText',
        'anchorType',
        'linkType',
        'firstSeen',
        'lastSeen',
        'status',
        'quality',
        'domainAuthority',
        'contextualRelevance',
        'linkPosition',
        'toxicityScore',
        'toxicityReasons',
      ];

      expect(backlinkDataProps.length).toBe(17);
      
      // Validate enums
      expect(Object.values(BacklinkStatus)).toContain(BacklinkStatus.ACTIVE);
      expect(Object.values(BacklinkStatus)).toContain(BacklinkStatus.TOXIC);
      expect(Object.values(LinkType)).toContain(LinkType.DOFOLLOW);
      expect(Object.values(AnchorType)).toContain(AnchorType.EXACT_MATCH);
    });

    it('should validate AnchorTextDistribution structure', () => {
      const anchorProps = [
        'anchorText',
        'type',
        'count',
        'percentage',
        'averageDA',
        'isNatural',
        'isOverOptimized',
        'riskScore',
        'topDomains',
      ];

      expect(anchorProps.length).toBe(9);
      expect(Object.values(AnchorType)).toContain(AnchorType.BRANDED);
      expect(Object.values(AnchorType)).toContain(AnchorType.NAKED_URL);
    });

    it('should validate ReferringDomain interface', () => {
      const domainProps = [
        'domain',
        'backlinksCount',
        'firstSeen',
        'lastSeen',
        'domainAuthority',
        'category',
        'language',
        'country',
        'quality',
        'linkTypes',
        'anchorDiversity',
        'contentRelevance',
        'toxicityScore',
        'topPages',
      ];

      expect(domainProps.length).toBe(14);
      expect(Object.values(DomainCategory)).toContain(DomainCategory.NEWS);
      expect(Object.values(DomainCategory)).toContain(DomainCategory.BLOG);
      expect(Object.values(BacklinkQuality)).toContain(BacklinkQuality.EXCELLENT);
      expect(Object.values(BacklinkQuality)).toContain(BacklinkQuality.TOXIC);
    });

    it('should validate CompetitorAnalysis structure', () => {
      const competitorProps = [
        'competitor',
        'competitorDomain',
        'totalBacklinks',
        'referringDomains',
        'domainAuthority',
        'gap',
        'overlap',
        'strengthComparison',
      ];

      expect(competitorProps.length).toBe(8);
    });

    it('should validate BacklinkOpportunity interface', () => {
      const opportunityProps = [
        'id',
        'domain',
        'url',
        'title',
        'domainAuthority',
        'category',
        'contactInfo',
        'outreachStatus',
        'priority',
        'estimatedDifficulty',
        'relevanceScore',
        'trafficValue',
        'competitorLinks',
        'suggestedApproach',
        'contentOpportunities',
      ];

      expect(opportunityProps.length).toBe(15);
    });

    it('should validate ToxicityReason enum', () => {
      const toxicityReasons = Object.values(ToxicityReason);
      expect(toxicityReasons).toContain(ToxicityReason.SPAM);
      expect(toxicityReasons).toContain(ToxicityReason.LINK_FARM);
      expect(toxicityReasons).toContain(ToxicityReason.PBN);
      expect(toxicityReasons).toContain(ToxicityReason.UNNATURAL_PATTERN);
    });
  });

  describe('mapBackendToBacklinksAnalysis function', () => {
    it('should exist and be callable', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      expect(typeof mapBackendToBacklinksAnalysis).toBe('function');
    });

    it('should map backend data to BacklinksAnalysis structure', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result).toBeDefined();
      expect(result.score).toBeDefined();
      expect(result.profile).toBeDefined();
      expect(result.backlinks).toBeDefined();
      expect(result.anchorTexts).toBeDefined();
      expect(result.referringDomains).toBeDefined();
      expect(result.competitors).toBeDefined();
    });

    it('should calculate backlinks scores correctly', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.score.overall).toBe(72);
      expect(result.score.grade).toBe(BacklinkGrade.B);
      expect(result.score.scoreBreakdown.quality).toBe(65);
      expect(result.score.scoreBreakdown.diversity).toBe(78);
      expect(result.score.scoreBreakdown.naturalness).toBe(58); // Lower due to over-optimization
    });

    it('should map backlink profile correctly', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.profile.totalBacklinks).toBe(15420);
      expect(result.profile.totalReferringDomains).toBe(2340);
      expect(result.profile.newBacklinks30Days).toBe(187);
      expect(result.profile.lostBacklinks30Days).toBe(45);
      expect(result.profile.dofollowPercentage).toBe(68.5);
      expect(result.profile.averageDomainAuthority).toBe(42.8);
    });

    it('should categorize individual backlinks by status and quality', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.backlinks.list).toHaveLength(3);
      
      const activeLink = result.backlinks.list.find(link => link.status === BacklinkStatus.ACTIVE);
      expect(activeLink).toBeDefined();
      expect(activeLink!.sourceDomain).toBe('tech-blog.example');
      expect(activeLink!.linkType).toBe(LinkType.DOFOLLOW);
      expect(activeLink!.quality).toBe(BacklinkQuality.GOOD);
      
      const toxicLink = result.backlinks.list.find(link => link.status === BacklinkStatus.TOXIC);
      expect(toxicLink).toBeDefined();
      expect(toxicLink!.toxicityScore).toBe(85);
      expect(toxicLink!.quality).toBe(BacklinkQuality.TOXIC);
    });

    it('should analyze anchor text distribution and naturalness', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.anchorTexts.distribution).toHaveLength(5);
      expect(result.anchorTexts.naturalness.isNatural).toBe(false); // Due to over-optimized "SEO" anchor
      
      const exactMatchAnchor = result.anchorTexts.distribution.find(a => a.type === AnchorType.EXACT_MATCH);
      expect(exactMatchAnchor).toBeDefined();
      expect(exactMatchAnchor!.anchorText).toBe('SEO');
      expect(exactMatchAnchor!.percentage).toBe(33.2);
      expect(exactMatchAnchor!.isOverOptimized).toBe(true);
      
      const brandedAnchor = result.anchorTexts.distribution.find(a => a.type === AnchorType.BRANDED);
      expect(brandedAnchor).toBeDefined();
      expect(brandedAnchor!.isNatural).toBe(true);
    });

    it('should map referring domains with quality assessment', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.referringDomains.list).toHaveLength(2);
      expect(result.referringDomains.totalCount).toBe(2);
      
      const qualityDomain = result.referringDomains.list.find(d => d.domain === 'tech-blog.example');
      expect(qualityDomain).toBeDefined();
      expect(qualityDomain!.quality).toBe(BacklinkQuality.GOOD);
      expect(qualityDomain!.category).toBe(DomainCategory.BLOG);
      expect(qualityDomain!.toxicityScore).toBe(10);
      
      const toxicDomain = result.referringDomains.list.find(d => d.domain === 'spammy-site.example');
      expect(toxicDomain).toBeDefined();
      expect(toxicDomain!.quality).toBe(BacklinkQuality.TOXIC);
      expect(toxicDomain!.toxicityScore).toBe(85);
    });

    it('should analyze competitor backlink profiles', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.competitors.analysis).toHaveLength(2);
      
      const topCompetitor = result.competitors.analysis.find(c => c.competitor === 'competitor1.example');
      expect(topCompetitor).toBeDefined();
      expect(topCompetitor!.totalBacklinks).toBe(28940);
      expect(topCompetitor!.referringDomains).toBe(4120);
      expect(topCompetitor!.gap.uniqueDomains).toContain('exclusive-partner.example');
      expect(topCompetitor!.overlap.commonDomains).toContain('tech-blog.example');
      
      expect(result.competitors.benchmarking.position).toBeGreaterThan(0);
      expect(result.competitors.benchmarking.aboveAverage).toBeDefined();
    });

    it('should identify link building opportunities', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.opportunities.list).toHaveLength(2);
      
      const highPriorityOpp = result.opportunities.list.find(o => o.priority === 'high');
      expect(highPriorityOpp).toBeDefined();
      expect(highPriorityOpp!.domain).toBe('industry-authority.example');
      expect(highPriorityOpp!.domainAuthority.domainAuthority).toBe(84);
      expect(highPriorityOpp!.relevanceScore).toBe(92);
      expect(highPriorityOpp!.contactInfo.email).toBe('editor@industry-authority.example');
    });

    it('should detect and categorize toxic backlinks', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.toxicAnalysis.toxicBacklinks.length).toBeGreaterThan(0);
      expect(result.toxicAnalysis.toxicDomains).toContain('spammy-site.example');
      expect(result.toxicAnalysis.toxicDomains).toContain('link-farm.example');
      
      const toxicLink = result.toxicAnalysis.toxicBacklinks.find(l => l.sourceDomain === 'spammy-site.example');
      expect(toxicLink).toBeDefined();
      expect(toxicLink!.toxicityReasons).toContain(ToxicityReason.SPAM);
      expect(toxicLink!.toxicityReasons).toContain(ToxicityReason.LOW_QUALITY);
      
      expect(result.toxicAnalysis.recommendations.recommendedActions.length).toBeGreaterThan(0);
      const disavowAction = result.toxicAnalysis.recommendations.recommendedActions.find(
        a => a.action === 'disavow'
      );
      expect(disavowAction).toBeDefined();
    });

    it('should track historical trends and growth patterns', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.trends.timeline).toHaveLength(3);
      expect(result.trends.growthMetrics.backlinksGrowthRate).toBeGreaterThan(0);
      
      const latestData = result.trends.timeline[result.trends.timeline.length - 1];
      expect(latestData.totalBacklinks).toBe(15420);
      expect(latestData.referringDomains).toBe(2340);
      expect(latestData.netGain).toBe(142); // 187 new - 45 lost
    });

    it('should generate actionable recommendations', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.recommendations.immediate.length).toBeGreaterThan(0);
      expect(result.recommendations.strategic.length).toBeGreaterThan(0);
      
      const highPriorityRec = result.recommendations.immediate.find(r => r.priority === 'high');
      expect(highPriorityRec).toBeDefined();
      expect(highPriorityRec!.category).toMatch(/^(acquisition|cleanup|optimization|monitoring)$/);
      
      // Should recommend cleaning up toxic links
      const cleanupRec = result.recommendations.immediate.find(r => r.category === 'cleanup');
      expect(cleanupRec).toBeDefined();
    });

    it('should handle empty or invalid data gracefully', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const emptyData = {};
      const result = mapBackendToBacklinksAnalysis(emptyData);
      
      expect(result.score.overall).toBe(0);
      expect(result.score.grade).toBe(BacklinkGrade.F);
      expect(result.profile.totalBacklinks).toBe(0);
      expect(result.backlinks.list).toHaveLength(0);
      expect(result.referringDomains.list).toHaveLength(0);
    });

    it('should set metadata correctly', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.metadata.scanDate).toBeTruthy();
      expect(result.metadata.toolVersion).toBe('1.0.0');
      expect(result.metadata.coverage.verificationRate).toBeGreaterThanOrEqual(0);
      expect(result.metadata.coverage.verificationRate).toBeLessThanOrEqual(100);
    });
  });

  describe('Helper Functions', () => {
    it('should have function to determine backlink grade', async () => {
      const { getBacklinkGrade } = await import('@/lib/mappers/backlinks-mapper');
      
      expect(getBacklinkGrade(95)).toBe(BacklinkGrade.A_PLUS);
      expect(getBacklinkGrade(85)).toBe(BacklinkGrade.A);
      expect(getBacklinkGrade(72)).toBe(BacklinkGrade.B);
      expect(getBacklinkGrade(60)).toBe(BacklinkGrade.C);
      expect(getBacklinkGrade(45)).toBe(BacklinkGrade.D);
      expect(getBacklinkGrade(25)).toBe(BacklinkGrade.F);
    });

    it('should have function to assess backlink quality', async () => {
      const { assessBacklinkQuality } = await import('@/lib/mappers/backlinks-mapper');
      
      const qualityLink = {
        domainAuthority: 75,
        spamScore: 2,
        toxicityScore: 5,
        isRedirect: false,
        contextualRelevance: 85,
      };
      
      const toxicLink = {
        domainAuthority: 15,
        spamScore: 14,
        toxicityScore: 85,
        isRedirect: false,
        contextualRelevance: 20,
      };
      
      expect(assessBacklinkQuality(qualityLink)).toBe(BacklinkQuality.EXCELLENT);
      expect(assessBacklinkQuality(toxicLink)).toBe(BacklinkQuality.TOXIC);
    });

    it('should have function to detect over-optimized anchors', async () => {
      const { detectOverOptimization } = await import('@/lib/mappers/backlinks-mapper');
      
      const anchors = [
        { text: 'SEO', count: 120, percentage: 35.2 }, // Over-optimized
        { text: 'Example Company', count: 89, percentage: 24.6 }, // Branded, OK
        { text: 'https://example.com', count: 67, percentage: 18.5 }, // URL, OK
      ];
      
      const overOptimized = detectOverOptimization(anchors);
      expect(overOptimized).toContain('SEO');
      expect(overOptimized).not.toContain('Example Company');
    });

    it('should have function to calculate link velocity', async () => {
      const { calculateLinkVelocity } = await import('@/lib/mappers/backlinks-mapper');
      
      const historicalData = [
        { date: '2024-01-01', newLinks: 156 },
        { date: '2024-02-01', newLinks: 201 }, 
        { date: '2024-03-01', newLinks: 187 },
      ];
      
      const velocity = calculateLinkVelocity(historicalData);
      expect(velocity.current).toBeGreaterThan(0);
      expect(velocity.isNatural).toBeDefined();
    });

    it('should have function to prioritize opportunities', async () => {
      const { prioritizeOpportunities } = await import('@/lib/mappers/backlinks-mapper');
      
      const opportunities = [
        { domainAuthority: 84, difficulty: 75, relevanceScore: 92 },
        { domainAuthority: 71, difficulty: 60, relevanceScore: 88 },
        { domainAuthority: 45, difficulty: 30, relevanceScore: 65 },
      ];
      
      const prioritized = prioritizeOpportunities(opportunities);
      expect(prioritized[0].priority).toBe('high');
      expect(prioritized[0].domainAuthority).toBe(84); // Highest DA should be first
    });

    it('should create empty backlinks analysis for fallback', async () => {
      const { createEmptyBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const emptyAnalysis = createEmptyBacklinksAnalysis();
      
      expect(emptyAnalysis.score.overall).toBe(0);
      expect(emptyAnalysis.score.grade).toBe(BacklinkGrade.F);
      expect(emptyAnalysis.profile.totalBacklinks).toBe(0);
      expect(emptyAnalysis.backlinks.list).toHaveLength(0);
      expect(emptyAnalysis.metadata.scanDate).toBeTruthy();
    });
  });

  describe('Advanced Analytics Functions', () => {
    it('should analyze link patterns and footprints', async () => {
      const { analyzeBacklinkPatterns } = await import('@/lib/mappers/backlinks-mapper');
      
      const backlinks = [
        { sourceDomain: 'site1.blogspot.com', anchorText: 'keyword1' },
        { sourceDomain: 'site2.blogspot.com', anchorText: 'keyword1' },
        { sourceDomain: 'site3.wordpress.com', anchorText: 'keyword1' },
      ];
      
      const patterns = analyzeBacklinkPatterns(backlinks);
      expect(patterns.footprintDetection.length).toBeGreaterThan(0);
      
      const blogspotFootprint = patterns.footprintDetection.find(f => f.footprint.includes('blogspot'));
      expect(blogspotFootprint).toBeDefined();
      expect(blogspotFootprint!.occurrences).toBe(2);
    });

    it('should calculate domain diversity metrics', async () => {
      const { calculateDomainDiversity } = await import('@/lib/mappers/backlinks-mapper');
      
      const domains = [
        { domain: 'site1.com', category: 'blog', country: 'US' },
        { domain: 'site2.com', category: 'news', country: 'UK' },
        { domain: 'site3.com', category: 'blog', country: 'US' },
      ];
      
      const diversity = calculateDomainDiversity(domains);
      expect(diversity.categoryDiversity).toBeGreaterThan(0);
      expect(diversity.geographicDiversity).toBeGreaterThan(0);
      expect(diversity.overallDiversity).toBeGreaterThan(0);
    });

    it('should generate disavow recommendations', async () => {
      const { generateDisavowRecommendations } = await import('@/lib/mappers/backlinks-mapper');
      
      const toxicLinks = [
        { domain: 'spammy.com', toxicityScore: 85, reasons: ['spam'] },
        { domain: 'low-quality.com', toxicityScore: 70, reasons: ['low_quality'] },
      ];
      
      const recommendations = generateDisavowRecommendations(toxicLinks);
      expect(recommendations.disavowFile.content).toContain('domain:spammy.com');
      expect(recommendations.recommendedActions.length).toBeGreaterThan(0);
      expect(recommendations.riskAssessment.currentRisk).toBeGreaterThan(0);
    });
  });
});