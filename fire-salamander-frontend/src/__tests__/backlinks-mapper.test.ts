/**
 * Fire Salamander - Backlinks Mapper Tests (TDD)
 * Lead Tech quality - Tests for backlinks mapper functions
 */

import { describe, it, expect, beforeEach } from '@jest/globals';
import {
  BacklinksAnalysis,
  BacklinkGrade,
  BacklinkQuality,
  BacklinkStatus,
  LinkType,
  AnchorType,
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

describe('Backlinks Mapper Tests (TDD)', () => {
  let mockBackendData: MockBackendBacklinksData;

  beforeEach(() => {
    mockBackendData = {
      id: 'backlinks-mapper-test-123',
      url: 'https://example.com',
      analyzed_at: '2024-03-15T14:30:00Z',
      backlinks_data: {
        profile: {
          total_backlinks: 1250,
          referring_domains: 340,
          referring_ips: 290,
          new_links_30d: 47,
          lost_links_30d: 12,
          dofollow_percentage: 73.2,
          average_da: 48.5,
          trust_flow: 38.7,
          citation_flow: 45.1,
        },
        backlinks: [
          {
            source_url: 'https://quality-blog.example/seo-article',
            target_url: 'https://example.com/guide',
            anchor_text: 'SEO best practices',
            link_type: 'dofollow',
            first_seen: '2024-01-10T09:00:00Z',
            last_seen: '2024-03-15T09:00:00Z',
            status: 'active',
            domain_authority: 72,
            page_authority: 68,
            spam_score: 1,
            is_redirect: false,
            http_status: 200,
            surrounding_text: 'For comprehensive information on SEO best practices, visit this detailed guide.',
            page_title: 'Complete Guide to SEO in 2024',
            toxicity_score: 3,
          },
          {
            source_url: 'https://spam-domain.example/link-page',
            target_url: 'https://example.com/',
            anchor_text: 'buy cheap SEO services now',
            link_type: 'dofollow',
            first_seen: '2024-03-05T15:00:00Z',
            last_seen: '2024-03-15T15:00:00Z',
            status: 'toxic',
            domain_authority: 18,
            page_authority: 15,
            spam_score: 13,
            is_redirect: false,
            http_status: 200,
            surrounding_text: 'Get quick results! buy cheap SEO services now for instant rankings!',
            page_title: 'Cheap SEO Services - Fast Results',
            toxicity_score: 87,
          },
        ],
        referring_domains: [
          {
            domain: 'quality-blog.example',
            backlinks_count: 5,
            first_seen: '2024-01-10T09:00:00Z',
            domain_authority: 72,
            category: 'blog',
            country: 'US',
            language: 'en',
            quality: 'excellent',
            toxicity_score: 5,
          },
          {
            domain: 'spam-domain.example',
            backlinks_count: 23,
            first_seen: '2024-03-05T15:00:00Z',
            domain_authority: 18,
            category: 'other',
            country: 'Unknown',
            language: 'en',
            quality: 'toxic',
            toxicity_score: 87,
          },
        ],
        scores: {
          overall: 65,
          quality: 58,
          diversity: 72,
          authority: 61,
          naturalness: 71,
          toxicity: 78, // 100 - 22 (average toxicity)
        },
      },
    };
  });

  describe('mapBackendToBacklinksAnalysis function', () => {
    it('should exist and be callable', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      expect(typeof mapBackendToBacklinksAnalysis).toBe('function');
    });

    it('should map backend data to BacklinksAnalysis interface', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result).toBeDefined();
      expect(result.score).toBeDefined();
      expect(result.profile).toBeDefined();
      expect(result.backlinks).toBeDefined();
      expect(result.anchorTexts).toBeDefined();
      expect(result.referringDomains).toBeDefined();
      expect(result.competitors).toBeDefined();
      expect(result.opportunities).toBeDefined();
      expect(result.toxicAnalysis).toBeDefined();
      expect(result.trends).toBeDefined();
      expect(result.analytics).toBeDefined();
      expect(result.recommendations).toBeDefined();
      expect(result.metadata).toBeDefined();
    });

    it('should calculate overall scores and grade correctly', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.score.overall).toBe(65);
      expect(result.score.grade).toBe(BacklinkGrade.C);
      expect(result.score.scoreBreakdown.quality).toBe(58);
      expect(result.score.scoreBreakdown.diversity).toBe(72);
      expect(result.score.scoreBreakdown.authority).toBe(61);
      expect(result.score.scoreBreakdown.naturalness).toBe(71);
      expect(result.score.scoreBreakdown.toxicity).toBe(78);
    });

    it('should map backlink profile data correctly', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.profile.totalBacklinks).toBe(1250);
      expect(result.profile.totalReferringDomains).toBe(340);
      expect(result.profile.totalReferringIPs).toBe(290);
      expect(result.profile.newBacklinks30Days).toBe(47);
      expect(result.profile.lostBacklinks30Days).toBe(12);
      expect(result.profile.dofollowPercentage).toBe(73.2);
      expect(result.profile.averageDomainAuthority).toBe(48.5);
      expect(result.profile.trustFlow).toBe(38.7);
      expect(result.profile.citationFlow).toBe(45.1);
    });

    it('should categorize backlinks by status and quality', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.backlinks.list).toHaveLength(2);
      expect(result.backlinks.totalCount).toBe(2);
      
      const activeLink = result.backlinks.list.find(link => link.status === BacklinkStatus.ACTIVE);
      expect(activeLink).toBeDefined();
      expect(activeLink!.sourceUrl).toBe('https://quality-blog.example/seo-article');
      expect(activeLink!.sourceDomain).toBe('quality-blog.example');
      expect(activeLink!.linkType).toBe(LinkType.DOFOLLOW);
      expect(activeLink!.quality).toBe(BacklinkQuality.EXCELLENT);
      expect(activeLink!.anchorType).toBe(AnchorType.PARTIAL_MATCH);
      expect(activeLink!.toxicityScore).toBe(3);
      
      const toxicLink = result.backlinks.list.find(link => link.status === BacklinkStatus.TOXIC);
      expect(toxicLink).toBeDefined();
      expect(toxicLink!.sourceDomain).toBe('spam-domain.example');
      expect(toxicLink!.quality).toBe(BacklinkQuality.TOXIC);
      expect(toxicLink!.toxicityScore).toBe(87);
      expect(toxicLink!.toxicityReasons).toContain(ToxicityReason.SPAM);
    });

    it('should map referring domains with quality assessment', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.referringDomains.list).toHaveLength(2);
      expect(result.referringDomains.totalCount).toBe(2);
      
      const qualityDomain = result.referringDomains.list.find(d => d.domain === 'quality-blog.example');
      expect(qualityDomain).toBeDefined();
      expect(qualityDomain!.quality).toBe(BacklinkQuality.EXCELLENT);
      expect(qualityDomain!.backlinksCount).toBe(5);
      expect(qualityDomain!.domainAuthority.domainAuthority).toBe(72);
      expect(qualityDomain!.toxicityScore).toBe(5);
      expect(qualityDomain!.isActive).toBe(true);
      
      const toxicDomain = result.referringDomains.list.find(d => d.domain === 'spam-domain.example');
      expect(toxicDomain).toBeDefined();
      expect(toxicDomain!.quality).toBe(BacklinkQuality.TOXIC);
      expect(toxicDomain!.toxicityScore).toBe(87);
      expect(toxicDomain!.isActive).toBe(false); // Should be marked as inactive due to toxicity
      
      // Check quality distribution
      expect(result.referringDomains.qualityDistribution[BacklinkQuality.EXCELLENT]).toBe(1);
      expect(result.referringDomains.qualityDistribution[BacklinkQuality.TOXIC]).toBe(1);
    });

    it('should generate anchor text distribution analysis', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.anchorTexts.distribution.length).toBeGreaterThan(0);
      
      const partialMatchAnchor = result.anchorTexts.distribution.find(a => 
        a.anchorText === 'SEO best practices'
      );
      expect(partialMatchAnchor).toBeDefined();
      expect(partialMatchAnchor!.type).toBe(AnchorType.PARTIAL_MATCH);
      expect(partialMatchAnchor!.isNatural).toBe(true);
      
      const spammyAnchor = result.anchorTexts.distribution.find(a => 
        a.anchorText === 'buy cheap SEO services now'
      );
      expect(spammyAnchor).toBeDefined();
      expect(spammyAnchor!.type).toBe(AnchorType.EXACT_MATCH);
      expect(spammyAnchor!.isOverOptimized).toBe(true);
      expect(spammyAnchor!.riskScore).toBeGreaterThan(70);
      
      expect(result.anchorTexts.naturalness.isNatural).toBeDefined();
      expect(result.anchorTexts.diversity.diversityScore).toBeGreaterThan(0);
    });

    it('should identify toxic links and generate recommendations', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.toxicAnalysis.toxicBacklinks.length).toBeGreaterThan(0);
      expect(result.toxicAnalysis.toxicDomains).toContain('spam-domain.example');
      expect(result.toxicAnalysis.overallToxicity).toBeGreaterThan(0);
      
      const toxicLink = result.toxicAnalysis.toxicBacklinks.find(l => 
        l.sourceDomain === 'spam-domain.example'
      );
      expect(toxicLink).toBeDefined();
      expect(toxicLink!.toxicityReasons).toContain(ToxicityReason.SPAM);
      expect(toxicLink!.toxicityReasons).toContain(ToxicityReason.LOW_QUALITY);
      
      expect(result.toxicAnalysis.recommendations.recommendedActions.length).toBeGreaterThan(0);
      const disavowAction = result.toxicAnalysis.recommendations.recommendedActions.find(
        a => a.action === 'disavow'
      );
      expect(disavowAction).toBeDefined();
      expect(disavowAction!.priority).toBe('high');
      
      expect(result.toxicAnalysis.recommendations.disavowFile.content).toContain('domain:spam-domain.example');
      expect(result.toxicAnalysis.recommendations.riskAssessment.currentRisk).toBeGreaterThan(0);
    });

    it('should generate actionable recommendations', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.recommendations.immediate.length).toBeGreaterThan(0);
      expect(result.recommendations.strategic.length).toBeGreaterThan(0);
      
      // Should recommend cleanup for toxic links
      const cleanupRec = result.recommendations.immediate.find(r => r.category === 'cleanup');
      expect(cleanupRec).toBeDefined();
      expect(cleanupRec!.priority).toBe('high');
      
      // Should recommend acquisition opportunities
      const acquisitionRec = result.recommendations.immediate.find(r => r.category === 'acquisition');
      expect(acquisitionRec).toBeDefined();
      
      // Strategic recommendations should have goals and tactics
      const strategicRec = result.recommendations.strategic[0];
      expect(strategicRec.goal).toBeTruthy();
      expect(strategicRec.tactics).toBeDefined();
      expect(strategicRec.kpis).toBeDefined();
      expect(strategicRec.timeline).toBeTruthy();
    });

    it('should handle empty or invalid data gracefully', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const emptyData = {};
      const result = mapBackendToBacklinksAnalysis(emptyData);
      
      expect(result.score.overall).toBe(0);
      expect(result.score.grade).toBe(BacklinkGrade.F);
      expect(result.profile.totalBacklinks).toBe(0);
      expect(result.profile.totalReferringDomains).toBe(0);
      expect(result.backlinks.list).toHaveLength(0);
      expect(result.referringDomains.list).toHaveLength(0);
      expect(result.anchorTexts.distribution).toHaveLength(0);
      expect(result.competitors.analysis).toHaveLength(0);
      expect(result.opportunities.list).toHaveLength(0);
    });

    it('should set metadata correctly', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.metadata.scanDate).toBeTruthy();
      expect(result.metadata.toolVersion).toBe('1.0.0');
      expect(result.metadata.coverage.crawledDomains).toBeGreaterThanOrEqual(0);
      expect(result.metadata.coverage.verificationRate).toBeGreaterThanOrEqual(0);
      expect(result.metadata.coverage.verificationRate).toBeLessThanOrEqual(100);
      expect(result.metadata.dataSources.length).toBeGreaterThan(0);
    });

    it('should calculate link velocity and naturalness', async () => {
      const { mapBackendToBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const result = mapBackendToBacklinksAnalysis(mockBackendData);
      
      expect(result.analytics.linkVelocity.current).toBeGreaterThan(0);
      expect(result.analytics.linkVelocity.natural).toBeGreaterThan(0);
      expect(result.analytics.linkVelocity.isNatural).toBeDefined();
      
      // With mock data, we should detect some patterns
      expect(result.analytics.linkPatterns.patternAnalysis.length).toBeGreaterThan(0);
      expect(result.analytics.linkPatterns.footprintDetection.length).toBeGreaterThan(0);
    });
  });

  describe('Helper functions', () => {
    it('should have a function to determine backlink grade', async () => {
      const { getBacklinkGrade } = await import('@/lib/mappers/backlinks-mapper');
      
      expect(getBacklinkGrade(95)).toBe(BacklinkGrade.A_PLUS);
      expect(getBacklinkGrade(90)).toBe(BacklinkGrade.A);
      expect(getBacklinkGrade(80)).toBe(BacklinkGrade.B);
      expect(getBacklinkGrade(65)).toBe(BacklinkGrade.C);
      expect(getBacklinkGrade(50)).toBe(BacklinkGrade.D);
      expect(getBacklinkGrade(30)).toBe(BacklinkGrade.F);
    });

    it('should have a function to assess backlink quality', async () => {
      const { assessBacklinkQuality } = await import('@/lib/mappers/backlinks-mapper');
      
      const qualityLink = {
        domainAuthority: 72,
        pageAuthority: 68,
        spamScore: 1,
        toxicityScore: 3,
        contextualRelevance: 85,
        isRedirect: false,
        httpStatus: 200,
      };
      
      const toxicLink = {
        domainAuthority: 18,
        pageAuthority: 15,
        spamScore: 13,
        toxicityScore: 87,
        contextualRelevance: 20,
        isRedirect: false,
        httpStatus: 200,
      };
      
      expect(assessBacklinkQuality(qualityLink)).toBe(BacklinkQuality.EXCELLENT);
      expect(assessBacklinkQuality(toxicLink)).toBe(BacklinkQuality.TOXIC);
    });

    it('should have a function to determine anchor type', async () => {
      const { determineAnchorType } = await import('@/lib/mappers/backlinks-mapper');
      
      expect(determineAnchorType('SEO', 'SEO guide')).toBe(AnchorType.EXACT_MATCH);
      expect(determineAnchorType('SEO best practices', 'SEO guide')).toBe(AnchorType.PARTIAL_MATCH);
      expect(determineAnchorType('Example Company', 'SEO guide')).toBe(AnchorType.BRANDED);
      expect(determineAnchorType('https://example.com', 'SEO guide')).toBe(AnchorType.NAKED_URL);
      expect(determineAnchorType('click here', 'SEO guide')).toBe(AnchorType.GENERIC);
      expect(determineAnchorType('[Image: SEO Guide]', 'SEO guide')).toBe(AnchorType.IMAGE);
    });

    it('should have a function to detect toxicity reasons', async () => {
      const { detectToxicityReasons } = await import('@/lib/mappers/backlinks-mapper');
      
      const spammyLink = {
        sourceDomain: 'spam-domain.example',
        anchorText: 'buy cheap SEO services now',
        domainAuthority: 18,
        spamScore: 13,
        toxicityScore: 87,
        surroundingText: 'Get quick results! buy cheap SEO services now for instant rankings!',
        pageTitle: 'Cheap SEO Services - Fast Results',
      };
      
      const reasons = detectToxicityReasons(spammyLink);
      
      expect(reasons).toContain(ToxicityReason.SPAM);
      expect(reasons).toContain(ToxicityReason.LOW_QUALITY);
      expect(reasons).toContain(ToxicityReason.SUSPICIOUS_ANCHOR);
    });

    it('should have a function to calculate domain authority metrics', async () => {
      const { calculateDomainAuthority } = await import('@/lib/mappers/backlinks-mapper');
      
      const domainData = {
        domain_authority: 72,
        page_authority: 68,
        spam_score: 1,
        organic_traffic: 15000,
        organic_keywords: 2500,
      };
      
      const authority = calculateDomainAuthority(domainData);
      
      expect(authority.domainAuthority).toBe(72);
      expect(authority.pageAuthority).toBe(68);
      expect(authority.spamScore).toBe(1);
      expect(authority.organicTraffic).toBe(15000);
      expect(authority.organicKeywords).toBe(2500);
      expect(authority.lastUpdated).toBeTruthy();
    });

    it('should have a function to prioritize opportunities', async () => {
      const { prioritizeOpportunities } = await import('@/lib/mappers/backlinks-mapper');
      
      const opportunities = [
        { 
          domainAuthority: { domainAuthority: 72 },
          estimatedDifficulty: 75,
          relevanceScore: 90,
          trafficValue: 5000,
        },
        { 
          domainAuthority: { domainAuthority: 85 },
          estimatedDifficulty: 90,
          relevanceScore: 85,
          trafficValue: 8000,
        },
        { 
          domainAuthority: { domainAuthority: 45 },
          estimatedDifficulty: 40,
          relevanceScore: 70,
          trafficValue: 1200,
        },
      ];
      
      const prioritized = prioritizeOpportunities(opportunities);
      
      expect(prioritized[0].priority).toBe('high');
      expect(prioritized[0].domainAuthority.domainAuthority).toBe(85); // Highest DA first
      expect(prioritized[2].priority).toBe('low'); // Lowest value last
    });

    it('should have a function to generate disavow file content', async () => {
      const { generateDisavowFile } = await import('@/lib/mappers/backlinks-mapper');
      
      const toxicLinks = [
        { url: 'https://spam1.example/page1', domain: 'spam1.example' },
        { url: 'https://spam1.example/page2', domain: 'spam1.example' },
        { url: 'https://spam2.example/page1', domain: 'spam2.example' },
      ];
      
      const disavowFile = generateDisavowFile(toxicLinks);
      
      expect(disavowFile.content).toContain('# Fire Salamander Disavow File');
      expect(disavowFile.content).toContain('domain:spam1.example');
      expect(disavowFile.content).toContain('domain:spam2.example');
      expect(disavowFile.domainsCount).toBe(2);
      expect(disavowFile.urlsCount).toBe(3);
      expect(disavowFile.lastUpdated).toBeTruthy();
    });

    it('should create empty backlinks analysis for fallback', async () => {
      const { createEmptyBacklinksAnalysis } = await import('@/lib/mappers/backlinks-mapper');
      
      const emptyAnalysis = createEmptyBacklinksAnalysis();
      
      expect(emptyAnalysis.score.overall).toBe(0);
      expect(emptyAnalysis.score.grade).toBe(BacklinkGrade.F);
      expect(emptyAnalysis.profile.totalBacklinks).toBe(0);
      expect(emptyAnalysis.profile.totalReferringDomains).toBe(0);
      expect(emptyAnalysis.backlinks.list).toHaveLength(0);
      expect(emptyAnalysis.anchorTexts.distribution).toHaveLength(0);
      expect(emptyAnalysis.referringDomains.list).toHaveLength(0);
      expect(emptyAnalysis.metadata.scanDate).toBeTruthy();
    });
  });

  describe('Advanced Analytics Functions', () => {
    it('should analyze link patterns and detect footprints', async () => {
      const { analyzeBacklinkPatterns } = await import('@/lib/mappers/backlinks-mapper');
      
      const backlinks = [
        { sourceDomain: 'site1.blogspot.com', anchorText: 'keyword1', pageTitle: 'Article about keyword1' },
        { sourceDomain: 'site2.blogspot.com', anchorText: 'keyword1', pageTitle: 'Post about keyword1' },
        { sourceDomain: 'site3.wordpress.com', anchorText: 'keyword1', pageTitle: 'Blog about keyword1' },
        { sourceDomain: 'quality-site.com', anchorText: 'brand name', pageTitle: 'Quality article' },
      ];
      
      const patterns = analyzeBacklinkPatterns(backlinks);
      
      expect(patterns.patternAnalysis.length).toBeGreaterThan(0);
      expect(patterns.footprintDetection.length).toBeGreaterThan(0);
      
      const blogspotFootprint = patterns.footprintDetection.find(f => 
        f.footprint.includes('blogspot')
      );
      expect(blogspotFootprint).toBeDefined();
      expect(blogspotFootprint!.occurrences).toBe(2);
      expect(blogspotFootprint!.risk).toBeGreaterThan(0);
    });

    it('should calculate link velocity and naturalness', async () => {
      const { calculateLinkVelocity } = await import('@/lib/mappers/backlinks-mapper');
      
      const historicalData = [
        { date: '2024-01-01', newLinks: 45, lostLinks: 8 },
        { date: '2024-02-01', newLinks: 52, lostLinks: 12 },
        { date: '2024-03-01', newLinks: 47, lostLinks: 9 },
      ];
      
      const velocity = calculateLinkVelocity(historicalData);
      
      expect(velocity.current).toBeGreaterThan(0);
      expect(velocity.natural).toBeGreaterThan(0);
      expect(velocity.isNatural).toBeDefined();
      
      if (!velocity.isNatural) {
        expect(velocity.warning).toBeTruthy();
      }
    });

    it('should calculate domain diversity metrics', async () => {
      const { calculateDomainDiversity } = await import('@/lib/mappers/backlinks-mapper');
      
      const domains = [
        { domain: 'blog1.com', category: 'blog', country: 'US', language: 'en' },
        { domain: 'news1.com', category: 'news', country: 'UK', language: 'en' },
        { domain: 'blog2.com', category: 'blog', country: 'CA', language: 'en' },
        { domain: 'corp1.com', category: 'corporate', country: 'US', language: 'en' },
        { domain: 'edu1.edu', category: 'education', country: 'US', language: 'en' },
      ];
      
      const diversity = calculateDomainDiversity(domains);
      
      expect(diversity.categoryDiversity).toBeGreaterThan(0);
      expect(diversity.categoryDiversity).toBeLessThanOrEqual(100);
      expect(diversity.geographicDiversity).toBeGreaterThan(0);
      expect(diversity.geographicDiversity).toBeLessThanOrEqual(100);
      expect(diversity.overallDiversity).toBeGreaterThan(0);
      expect(diversity.overallDiversity).toBeLessThanOrEqual(100);
    });

    it('should detect over-optimized anchor text patterns', async () => {
      const { detectOverOptimization } = await import('@/lib/mappers/backlinks-mapper');
      
      const anchors = [
        { text: 'SEO services', count: 85, percentage: 42.5 }, // Over-optimized
        { text: 'Company Name', count: 45, percentage: 22.5 }, // Branded, OK
        { text: 'https://example.com', count: 30, percentage: 15.0 }, // URL, OK
        { text: 'best SEO company', count: 25, percentage: 12.5 }, // Potentially over-optimized
        { text: 'click here', count: 15, percentage: 7.5 }, // Generic, OK
      ];
      
      const overOptimized = detectOverOptimization(anchors);
      
      expect(overOptimized).toContain('SEO services'); // 42.5% is way too high
      expect(overOptimized).not.toContain('Company Name'); // Branded is OK
      expect(overOptimized).not.toContain('https://example.com'); // URL is OK
      expect(overOptimized).not.toContain('click here'); // Generic is OK
      
      // 'best SEO company' at 12.5% might be flagged depending on threshold
    });

    it('should generate link building insights', async () => {
      const { generateLinkBuildingInsights } = await import('@/lib/mappers/backlinks-mapper');
      
      const analysisData = {
        backlinks: mockBackendData.backlinks_data?.backlinks || [],
        competitorData: [],
        industryData: { averageBacklinks: 2500, averageDA: 52 },
      };
      
      const insights = generateLinkBuildingInsights(analysisData);
      
      expect(insights.topPerformingContent).toBeDefined();
      expect(insights.contentGaps).toBeDefined();
      expect(insights.linkBaitOpportunities).toBeDefined();
      expect(insights.industryBenchmarks).toBeDefined();
      
      expect(insights.industryBenchmarks.averageBacklinks).toBe(2500);
      expect(insights.industryBenchmarks.averageDomainAuthority).toBe(52);
    });
  });
});