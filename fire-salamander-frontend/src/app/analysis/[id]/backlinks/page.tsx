"use client";

import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { BacklinksAnalysisSection } from "@/components/analysis/backlinks-analysis-section";
import { mapBackendToBacklinksAnalysis } from "@/lib/mappers/backlinks-mapper";
import { 
  BacklinksAnalysis, 
  BacklinkGrade, 
  BacklinkQuality, 
  BacklinkStatus,
  LinkType,
  AnchorType,
  DomainCategory,
  ToxicityReason,
} from "@/types/backlinks-analysis";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertTriangle, Activity, Link2 } from 'lucide-react';

interface AnalysisData {
  id: string;
  url: string;
  analyzed_at: string;
  backlinks_data?: any;
}

export default function AnalysisBacklinksPage() {
  const params = useParams();
  const [backlinksData, setBacklinksData] = useState<BacklinksAnalysis | null>(null);
  const [analysis, setAnalysis] = useState<AnalysisData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const analysisId = params.id as string;

  useEffect(() => {
    const fetchBacklinksAnalysis = async () => {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/api/v1/analysis/${analysisId}/backlinks`);
        
        if (response.ok) {
          const data = await response.json();
          setAnalysis(data.data);
          
          // Map backend data to our BacklinksAnalysis interface
          const mappedData = mapBackendToBacklinksAnalysis(data.data);
          setBacklinksData(mappedData);
        } else {
          setError("Analyse de backlinks non trouvée");
        }
      } catch (err) {
        console.error("Erreur lors de la récupération de l'analyse de backlinks:", err);
        setError("Erreur de connexion");
        
        // Fallback to mock data for development
        const mockData = createMockBacklinksData();
        setBacklinksData(mockData);
        console.log('Using mock backlinks data:', mockData);
      } finally {
        setLoading(false);
      }
    };

    fetchBacklinksAnalysis();
  }, [analysisId]);

  const createMockBacklinksData = (): BacklinksAnalysis => {
    return {
      score: {
        overall: 72,
        grade: BacklinkGrade.B,
        trend: 'improving',
        previousScore: 68,
        scoreBreakdown: {
          quantity: 75,
          quality: 68,
          diversity: 80,
          authority: 71,
          naturalness: 65,
          toxicity: 82,
        },
      },

      profile: {
        totalBacklinks: 15420,
        totalReferringDomains: 2340,
        totalReferringIPs: 1980,
        totalReferringSubnets: 1650,
        newBacklinks30Days: 187,
        lostBacklinks30Days: 45,
        brokenBacklinks: 23,
        dofollowPercentage: 68.5,
        nofollowPercentage: 31.5,
        averageDomainAuthority: 42.8,
        averagePageAuthority: 38.6,
        trustFlow: 35.2,
        citationFlow: 41.6,
        majesticTrustRatio: 0.85,
      },

      backlinks: {
        list: [
          {
            id: 'backlink-1',
            sourceUrl: 'https://tech-authority.example/seo-guide-2024',
            targetUrl: 'https://example.com/services',
            sourceDomain: 'tech-authority.example',
            targetDomain: 'example.com',
            anchorText: 'comprehensive SEO services',
            anchorType: AnchorType.PARTIAL_MATCH,
            linkType: LinkType.DOFOLLOW,
            firstSeen: '2024-01-15T10:00:00Z',
            lastSeen: '2024-03-15T10:00:00Z',
            status: BacklinkStatus.ACTIVE,
            quality: BacklinkQuality.EXCELLENT,
            domainAuthority: {
              domainAuthority: 84,
              pageAuthority: 79,
              domainRating: 82,
              urlRating: 77,
              trustFlow: 68,
              citationFlow: 75,
              spamScore: 1,
              organicTraffic: 250000,
              organicKeywords: 15000,
              lastUpdated: new Date().toISOString(),
            },
            contextualRelevance: 92,
            linkPosition: {
              isMainContent: true,
              isNavigation: false,
              isFooter: false,
              isSidebar: false,
              position: 1,
            },
            suroundingText: 'For businesses looking for comprehensive SEO services, this platform offers advanced tools and strategies.',
            pageTitle: 'Ultimate SEO Guide 2024 - Best Practices and Tools',
            pageCategory: DomainCategory.BLOG,
            language: 'en',
            country: 'US',
            isRedirect: false,
            redirectChain: [],
            httpStatus: 200,
            toxicityScore: 5,
            toxicityReasons: [],
          },
          {
            id: 'backlink-2',
            sourceUrl: 'https://industry-news.example/tech-trends',
            targetUrl: 'https://example.com/about',
            sourceDomain: 'industry-news.example',
            targetDomain: 'example.com',
            anchorText: 'Example Company',
            anchorType: AnchorType.BRANDED,
            linkType: LinkType.DOFOLLOW,
            firstSeen: '2024-02-01T14:00:00Z',
            lastSeen: '2024-03-15T14:00:00Z',
            status: BacklinkStatus.ACTIVE,
            quality: BacklinkQuality.GOOD,
            domainAuthority: {
              domainAuthority: 71,
              pageAuthority: 68,
              domainRating: 69,
              urlRating: 65,
              trustFlow: 58,
              citationFlow: 67,
              spamScore: 2,
              organicTraffic: 180000,
              organicKeywords: 12000,
              lastUpdated: new Date().toISOString(),
            },
            contextualRelevance: 88,
            linkPosition: {
              isMainContent: true,
              isNavigation: false,
              isFooter: false,
              isSidebar: false,
              position: 2,
            },
            suroundingText: 'Example Company has been recognized as a leader in innovative SEO technology solutions.',
            pageTitle: 'Tech Industry Leaders 2024 - Innovation Report',
            pageCategory: DomainCategory.NEWS,
            language: 'en',
            country: 'US',
            isRedirect: false,
            redirectChain: [],
            httpStatus: 200,
            toxicityScore: 8,
            toxicityReasons: [],
          },
          {
            id: 'backlink-3',
            sourceUrl: 'https://spam-network.example/quick-links',
            targetUrl: 'https://example.com/',
            sourceDomain: 'spam-network.example',
            targetDomain: 'example.com',
            anchorText: 'buy cheap SEO services fast ranking guaranteed',
            anchorType: AnchorType.EXACT_MATCH,
            linkType: LinkType.DOFOLLOW,
            firstSeen: '2024-03-05T08:00:00Z',
            lastSeen: '2024-03-15T08:00:00Z',
            status: BacklinkStatus.TOXIC,
            quality: BacklinkQuality.TOXIC,
            domainAuthority: {
              domainAuthority: 18,
              pageAuthority: 15,
              domainRating: 16,
              urlRating: 12,
              trustFlow: 5,
              citationFlow: 22,
              spamScore: 14,
              organicTraffic: 500,
              organicKeywords: 50,
              lastUpdated: new Date().toISOString(),
            },
            contextualRelevance: 25,
            linkPosition: {
              isMainContent: false,
              isNavigation: false,
              isFooter: true,
              isSidebar: false,
              position: 15,
            },
            suroundingText: 'Get instant results! buy cheap SEO services fast ranking guaranteed with money back guarantee!',
            pageTitle: 'Cheap SEO Services - Fast Results Guaranteed',
            pageCategory: DomainCategory.OTHER,
            language: 'en',
            country: 'Unknown',
            isRedirect: false,
            redirectChain: [],
            httpStatus: 200,
            toxicityScore: 89,
            toxicityReasons: [ToxicityReason.SPAM, ToxicityReason.LOW_QUALITY, ToxicityReason.SUSPICIOUS_ANCHOR, ToxicityReason.UNNATURAL_PATTERN],
          },
        ],
        totalCount: 15420,
        pagination: {
          currentPage: 1,
          totalPages: 771,
          pageSize: 20,
        },
        filtering: {
          availableFilters: {
            quality: Object.values(BacklinkQuality),
            status: Object.values(BacklinkStatus),
            linkType: Object.values(LinkType),
            domainCategories: Object.values(DomainCategory),
            countries: ['US', 'UK', 'CA', 'AU', 'DE', 'FR', 'Unknown'],
            languages: ['en', 'de', 'fr', 'es', 'it'],
          },
        },
      },

      anchorTexts: {
        distribution: [
          {
            anchorText: 'Example Company',
            type: AnchorType.BRANDED,
            count: 3420,
            percentage: 22.2,
            averageDA: 58.3,
            isNatural: true,
            isOverOptimized: false,
            riskScore: 5,
            topDomains: ['tech-authority.example', 'industry-news.example', 'business-review.example'],
          },
          {
            anchorText: 'https://example.com',
            type: AnchorType.NAKED_URL,
            count: 2850,
            percentage: 18.5,
            averageDA: 45.7,
            isNatural: true,
            isOverOptimized: false,
            riskScore: 0,
            topDomains: ['directory.example', 'forum.example', 'social-media.example'],
          },
          {
            anchorText: 'SEO services',
            type: AnchorType.EXACT_MATCH,
            count: 1890,
            percentage: 12.3,
            averageDA: 41.2,
            isNatural: false,
            isOverOptimized: true,
            riskScore: 75,
            topDomains: ['blog1.example', 'blog2.example', 'article-site.example'],
          },
          {
            anchorText: 'comprehensive SEO platform',
            type: AnchorType.PARTIAL_MATCH,
            count: 1560,
            percentage: 10.1,
            averageDA: 52.8,
            isNatural: true,
            isOverOptimized: false,
            riskScore: 15,
            topDomains: ['tech-blog.example', 'marketing-site.example'],
          },
        ],
        naturalness: {
          score: 65,
          isNatural: false,
          overOptimizedAnchors: ['SEO services', 'best SEO company', 'cheap SEO'],
          recommendations: [
            'Reduce exact match anchor text percentage below 10%',
            'Increase branded anchor text variations',
            'Focus on natural contextual anchor text',
            'Diversify anchor text across different domains',
          ],
        },
        diversity: {
          uniqueAnchors: 892,
          diversityScore: 78,
          topAnchors: [
            { text: 'Example Company', count: 3420, percentage: 22.2 },
            { text: 'https://example.com', count: 2850, percentage: 18.5 },
            { text: 'SEO services', count: 1890, percentage: 12.3 },
          ],
        },
      },

      referringDomains: {
        list: [
          {
            domain: 'tech-authority.example',
            backlinksCount: 15,
            firstSeen: '2024-01-15T10:00:00Z',
            lastSeen: '2024-03-15T10:00:00Z',
            domainAuthority: {
              domainAuthority: 84,
              pageAuthority: 79,
              domainRating: 82,
              urlRating: 77,
              trustFlow: 68,
              citationFlow: 75,
              spamScore: 1,
              organicTraffic: 250000,
              organicKeywords: 15000,
              lastUpdated: new Date().toISOString(),
            },
            category: DomainCategory.BLOG,
            language: 'en',
            country: 'US',
            isActive: true,
            quality: BacklinkQuality.EXCELLENT,
            linkTypes: {
              dofollow: 12,
              nofollow: 3,
              sponsored: 0,
              ugc: 0,
            },
            anchorDiversity: 85,
            contentRelevance: 92,
            toxicityScore: 5,
            topPages: [
              {
                url: 'https://tech-authority.example/seo-guide-2024',
                title: 'Ultimate SEO Guide 2024',
                backlinksCount: 8,
                pageAuthority: 79,
              },
            ],
          },
        ],
        totalCount: 2340,
        qualityDistribution: {
          [BacklinkQuality.EXCELLENT]: 345,
          [BacklinkQuality.GOOD]: 892,
          [BacklinkQuality.AVERAGE]: 756,
          [BacklinkQuality.POOR]: 289,
          [BacklinkQuality.TOXIC]: 58,
        },
        authorityDistribution: {
          high: 423,  // 80+
          medium: 1156, // 50-79
          low: 761, // 0-49
        },
        categories: {
          [DomainCategory.BLOG]: 645,
          [DomainCategory.NEWS]: 423,
          [DomainCategory.CORPORATE]: 356,
          [DomainCategory.DIRECTORY]: 289,
          [DomainCategory.FORUM]: 234,
          [DomainCategory.EDUCATION]: 178,
          [DomainCategory.GOVERNMENT]: 89,
          [DomainCategory.SOCIAL_MEDIA]: 67,
          [DomainCategory.ECOMMERCE]: 45,
          [DomainCategory.OTHER]: 14,
        },
        geographicDistribution: {
          'US': 1245,
          'UK': 378,
          'CA': 234,
          'AU': 156,
          'DE': 123,
          'FR': 98,
          'Unknown': 106,
        },
        topDomains: [],
      },

      competitors: {
        analysis: [
          {
            competitor: 'Competitor A',
            competitorDomain: 'competitor-a.example',
            totalBacklinks: 28940,
            referringDomains: 4120,
            domainAuthority: {
              domainAuthority: 78,
              pageAuthority: 74,
              domainRating: 76,
              urlRating: 71,
              trustFlow: 65,
              citationFlow: 72,
              spamScore: 2,
              organicTraffic: 450000,
              organicKeywords: 28000,
              lastUpdated: new Date().toISOString(),
            },
            gap: {
              uniqueBacklinks: 12500,
              uniqueDomains: 1890,
              opportunities: [],
            },
            overlap: {
              commonBacklinks: 890,
              commonDomains: 234,
              sharedDomains: ['tech-authority.example', 'industry-news.example'],
            },
            strengthComparison: {
              stronger: true,
              areas: ['Domain Authority', 'Backlink Volume', 'Content Depth'],
              advantages: ['Established authority', 'Diverse link profile', 'Strong technical SEO'],
              weaknesses: ['Lower content freshness', 'Limited social engagement'],
            },
          },
        ],
        benchmarking: {
          position: 3,
          totalCompetitors: 8,
          aboveAverage: true,
          strengths: ['Quality backlink profile', 'Strong brand mentions', 'Good anchor text diversity'],
          opportunities: ['Industry partnerships', 'Guest posting', 'Resource page listings'],
          threats: ['Competitor link acquisition', 'Industry consolidation'],
        },
        gapAnalysis: {
          totalOpportunities: 145,
          highPriorityOpportunities: [],
          estimatedValue: 285000,
          timeToAcquire: '4-8 months',
        },
      },

      opportunities: {
        list: [
          {
            id: 'opportunity-1',
            domain: 'industry-authority.example',
            url: 'https://industry-authority.example/resources',
            title: 'Industry Resources Directory',
            domainAuthority: {
              domainAuthority: 87,
              pageAuthority: 82,
              domainRating: 85,
              urlRating: 80,
              trustFlow: 72,
              citationFlow: 79,
              spamScore: 1,
              organicTraffic: 380000,
              organicKeywords: 22000,
              lastUpdated: new Date().toISOString(),
            },
            category: DomainCategory.CORPORATE,
            contactInfo: {
              email: 'editor@industry-authority.example',
              contactPage: 'https://industry-authority.example/contact',
              socialMedia: {
                twitter: '@industryauth',
                linkedin: 'company/industry-authority',
              },
            },
            outreachStatus: 'not_contacted',
            priority: 'high',
            estimatedDifficulty: 75,
            relevanceScore: 94,
            trafficValue: 15000,
            competitorLinks: ['competitor-a.example', 'competitor-b.example'],
            suggestedApproach: 'Resource page outreach with valuable industry data',
            contentOpportunities: ['Industry report', 'Case study collaboration'],
          },
        ],
        totalCount: 145,
        priorityDistribution: {
          high: 23,
          medium: 67,
          low: 55,
        },
        estimatedTotalValue: 285000,
        insights: {
          topPerformingContent: [
            {
              url: 'https://example.com/seo-guide',
              title: 'Complete SEO Guide 2024',
              backlinksCount: 234,
              referringDomains: 89,
              socialShares: 1250,
              contentType: 'guide',
              publishDate: '2024-01-15T00:00:00Z',
              topics: ['SEO', 'Digital Marketing', 'Content Strategy'],
            },
          ],
          contentGaps: [
            {
              topic: 'Technical SEO Audits',
              competitorCoverage: 85,
              opportunity: 92,
              suggestedContentType: 'comprehensive guide',
              estimatedBacklinkPotential: 45,
            },
          ],
          linkBaitOpportunities: [
            {
              type: 'study',
              topic: 'SEO Industry Benchmarks 2024',
              difficulty: 70,
              potential: 88,
              competitorExamples: ['competitor-study.example'],
              resources: ['Survey platform', 'Data analysis tools'],
            },
          ],
          industryBenchmarks: {
            averageBacklinks: 12500,
            averageReferringDomains: 1890,
            averageDomainAuthority: 58,
            topPerformers: ['leader1.example', 'leader2.example'],
            industryTrends: [
              {
                trend: 'AI-powered content optimization',
                growth: 45,
                timeframe: 'Last 6 months',
              },
            ],
          },
        },
      },

      toxicAnalysis: {
        toxicBacklinks: [
          {
            id: 'backlink-3',
            sourceUrl: 'https://spam-network.example/quick-links',
            targetUrl: 'https://example.com/',
            sourceDomain: 'spam-network.example',
            targetDomain: 'example.com',
            anchorText: 'buy cheap SEO services fast ranking guaranteed',
            anchorType: AnchorType.EXACT_MATCH,
            linkType: LinkType.DOFOLLOW,
            firstSeen: '2024-03-05T08:00:00Z',
            lastSeen: '2024-03-15T08:00:00Z',
            status: BacklinkStatus.TOXIC,
            quality: BacklinkQuality.TOXIC,
            domainAuthority: {
              domainAuthority: 18,
              pageAuthority: 15,
              domainRating: 16,
              urlRating: 12,
              trustFlow: 5,
              citationFlow: 22,
              spamScore: 14,
              organicTraffic: 500,
              organicKeywords: 50,
              lastUpdated: new Date().toISOString(),
            },
            contextualRelevance: 25,
            linkPosition: {
              isMainContent: false,
              isNavigation: false,
              isFooter: true,
              isSidebar: false,
              position: 15,
            },
            suroundingText: 'Get instant results! buy cheap SEO services fast ranking guaranteed with money back guarantee!',
            pageTitle: 'Cheap SEO Services - Fast Results Guaranteed',
            pageCategory: DomainCategory.OTHER,
            language: 'en',
            country: 'Unknown',
            isRedirect: false,
            redirectChain: [],
            httpStatus: 200,
            toxicityScore: 89,
            toxicityReasons: [ToxicityReason.SPAM, ToxicityReason.LOW_QUALITY, ToxicityReason.SUSPICIOUS_ANCHOR],
          },
        ],
        toxicDomains: ['spam-network.example', 'link-farm-site.example', 'pbn-domain.example'],
        overallToxicity: 22,
        recommendations: {
          recommendedActions: [
            {
              backlink: {
                id: 'backlink-3',
                sourceUrl: 'https://spam-network.example/quick-links',
                targetUrl: 'https://example.com/',
                sourceDomain: 'spam-network.example',
                targetDomain: 'example.com',
                anchorText: 'buy cheap SEO services fast ranking guaranteed',
                anchorType: AnchorType.EXACT_MATCH,
                linkType: LinkType.DOFOLLOW,
                firstSeen: '2024-03-05T08:00:00Z',
                lastSeen: '2024-03-15T08:00:00Z',
                status: BacklinkStatus.TOXIC,
                quality: BacklinkQuality.TOXIC,
                domainAuthority: {
                  domainAuthority: 18,
                  pageAuthority: 15,
                  domainRating: 16,
                  urlRating: 12,
                  trustFlow: 5,
                  citationFlow: 22,
                  spamScore: 14,
                  organicTraffic: 500,
                  organicKeywords: 50,
                  lastUpdated: new Date().toISOString(),
                },
                contextualRelevance: 25,
                linkPosition: {
                  isMainContent: false,
                  isNavigation: false,
                  isFooter: true,
                  isSidebar: false,
                  position: 15,
                },
                suroundingText: 'Get instant results!',
                pageTitle: 'Cheap SEO Services',
                pageCategory: DomainCategory.OTHER,
                language: 'en',
                country: 'Unknown',
                isRedirect: false,
                redirectChain: [],
                httpStatus: 200,
                toxicityScore: 89,
                toxicityReasons: [ToxicityReason.SPAM],
              },
              action: 'disavow',
              priority: 'high',
              reason: 'Extremely high toxicity score (89%) with spammy anchor text and low-quality domain',
              potentialImpact: 'High risk of algorithmic penalty - immediate disavowal required',
            },
          ],
          disavowFile: {
            content: '# Fire Salamander Disavow File\n# Generated on 2024-03-15\ndomain:spam-network.example\ndomain:link-farm-site.example\ndomain:pbn-domain.example',
            lastUpdated: new Date().toISOString(),
            domainsCount: 3,
            urlsCount: 0,
          },
          riskAssessment: {
            currentRisk: 22,
            riskAfterDisavow: 8,
            estimatedRecoveryTime: '2-4 months',
            confidenceLevel: 87,
          },
        },
        riskFactors: [
          {
            factor: 'Over-optimized exact match anchors',
            severity: 'medium',
            affectedLinks: 1890,
            description: 'High percentage of exact match anchor text indicating potential over-optimization',
            impact: 'May reduce ranking effectiveness and trigger algorithmic review',
          },
          {
            factor: 'Toxic domain cluster',
            severity: 'high',
            affectedLinks: 47,
            description: 'Multiple toxic domains with similar spam patterns detected',
            impact: 'High risk of manual penalty if not addressed promptly',
          },
        ],
      },

      trends: {
        timeline: [
          {
            date: '2024-01-01',
            totalBacklinks: 14200,
            referringDomains: 2180,
            domainAuthority: 38,
            newBacklinks: 156,
            lostBacklinks: 23,
            netGain: 133,
          },
          {
            date: '2024-02-01',
            totalBacklinks: 14850,
            referringDomains: 2260,
            domainAuthority: 40,
            newBacklinks: 201,
            lostBacklinks: 34,
            netGain: 167,
          },
          {
            date: '2024-03-01',
            totalBacklinks: 15420,
            referringDomains: 2340,
            domainAuthority: 42,
            newBacklinks: 187,
            lostBacklinks: 45,
            netGain: 142,
          },
        ],
        growthMetrics: {
          backlinksGrowthRate: 4.2,
          domainsGrowthRate: 3.6,
          authorityGrowthRate: 5.3,
          velocityScore: 78,
        },
        seasonalTrends: [
          {
            period: 'Q1 2024',
            avgBacklinks: 14823,
            avgDomains: 2260,
            pattern: 'increasing',
          },
        ],
        milestones: [
          {
            date: '2024-02-15',
            event: 'Major industry publication feature',
            impact: 234,
            description: 'Featured in TechCrunch article resulting in 234 new backlinks',
          },
        ],
      },

      analytics: {
        linkVelocity: {
          current: 181,
          natural: 150,
          isNatural: false,
          warning: 'Link velocity is 20% above natural baseline - monitor for continued growth patterns',
        },
        linkPatterns: {
          patternAnalysis: [
            {
              pattern: 'High concentration of exact match anchors',
              confidence: 85,
              risk: 'medium',
              description: '12.3% of links use exact match anchor text "SEO services"',
            },
          ],
          footprintDetection: [
            {
              footprint: 'Multiple links from .blogspot.com domains',
              occurrences: 23,
              risk: 45,
              affectedLinks: [
                'https://blog1.blogspot.com/seo-post',
                'https://blog2.blogspot.com/marketing-tips',
              ],
            },
          ],
        },
        contentAnalysis: {
          topLinkingPages: [
            {
              url: 'https://example.com/seo-guide',
              title: 'Complete SEO Guide 2024',
              linksCount: 234,
              contentType: 'guide',
              topics: ['SEO', 'Digital Marketing'],
            },
          ],
          linkingContext: [
            {
              context: 'Editorial mentions in articles',
              frequency: 45,
              naturalness: 92,
            },
            {
              context: 'Resource page listings',
              frequency: 38,
              naturalness: 88,
            },
          ],
        },
      },

      recommendations: {
        immediate: [
          {
            id: 'rec-1',
            title: 'Disavow toxic backlinks immediately',
            priority: 'high',
            category: 'cleanup',
            description: 'Remove 47 toxic backlinks from spam domains to prevent algorithmic penalties',
            expectedImpact: 'Reduce penalty risk by 70% and improve domain health',
            effort: 'low',
            timeline: '1-2 weeks',
            resources: ['Google Disavow Tool', 'Link analysis team'],
          },
          {
            id: 'rec-2',
            title: 'Diversify anchor text profile',
            priority: 'high',
            category: 'optimization',
            description: 'Reduce exact match anchor percentage from 12.3% to below 8%',
            expectedImpact: 'Improve anchor text naturalness and reduce over-optimization risk',
            effort: 'medium',
            timeline: '2-3 months',
            resources: ['Content team', 'Outreach specialists'],
          },
          {
            id: 'rec-3',
            title: 'Pursue high-DA opportunities',
            priority: 'medium',
            category: 'acquisition',
            description: 'Target 23 high-priority domains with DA 80+ for quality link acquisition',
            expectedImpact: 'Increase average domain authority and organic traffic by 15%',
            effort: 'high',
            timeline: '3-6 months',
            resources: ['Outreach team', 'Content creation budget'],
          },
        ],
        strategic: [
          {
            id: 'strategy-1',
            title: 'Comprehensive link building campaign',
            goal: 'Increase referring domains by 30% while maintaining quality standards',
            tactics: [
              'Resource page outreach to industry authorities',
              'Guest posting on relevant high-DA publications',
              'Digital PR for thought leadership content',
              'Partnership development with complementary brands',
            ],
            kpis: [
              'New referring domains per month: 80+',
              'Average DA of new links: 55+',
              'Anchor text diversity score: 85+',
              'Link acquisition cost: <$150 per link',
            ],
            timeline: '6-12 months',
            budget: '$25,000-$50,000',
          },
        ],
      },

      metadata: {
        scanDate: new Date().toISOString(),
        scanDuration: 420,
        toolVersion: '1.0.0',
        dataFreshness: 6,
        coverage: {
          crawledDomains: 2340,
          totalFoundBacklinks: 15420,
          verifiedBacklinks: 14892,
          verificationRate: 97,
        },
        limitations: [
          'Some backlinks may not be discoverable due to robots.txt restrictions',
          'Historical data limited to 12 months based on crawl frequency',
          'Competitor analysis based on publicly available link data',
        ],
        dataSources: [
          'Common Crawl Archives',
          'Search Engine APIs', 
          'Majestic SEO Database',
          'Direct Website Crawling',
        ],
      },
    };
  };

  console.log('State check:', { loading, error, hasBacklinksData: !!backlinksData });

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <Activity className="h-8 w-8 animate-spin mx-auto mb-4 text-blue-500" />
          <p className="text-gray-600">Chargement de l'analyse des backlinks...</p>
        </div>
      </div>
    );
  }

  if (error && !backlinksData) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <Card className="w-full max-w-md">
          <CardHeader>
            <CardTitle className="flex items-center space-x-2 text-red-600">
              <AlertTriangle className="h-5 w-5" />
              <span>Erreur</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-gray-600 mb-4">{error}</p>
            <Button 
              onClick={() => window.location.reload()} 
              className="w-full"
            >
              Réessayer
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!backlinksData) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <Link2 className="h-12 w-12 mx-auto mb-4 text-gray-400" />
          <h3 className="text-lg font-semibold mb-2">Aucune données de backlinks</h3>
          <p className="text-gray-600">L'analyse des backlinks n'est pas encore disponible pour cette page.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Debug Info - Development only */}
      {process.env.NODE_ENV === 'development' && (
        <div className="mb-4 p-4 bg-gray-100 rounded">
          <p>Backlinks Data Status: {backlinksData ? 'Loaded' : 'Missing'}</p>
          <p>Overall Score: {backlinksData?.score.overall || 0}/100</p>
          <p>Grade: {backlinksData?.score.grade || 'N/A'}</p>
          <p>Total Backlinks: {backlinksData?.profile.totalBacklinks.toLocaleString() || 0}</p>
          <p>Referring Domains: {backlinksData?.profile.totalReferringDomains.toLocaleString() || 0}</p>
          <p>Analysis ID: {analysisId}</p>
          <p>Toxic Links: {backlinksData?.toxicAnalysis.toxicBacklinks.length || 0}</p>
          <p>Quality Score: {backlinksData?.score.scoreBreakdown.quality || 0}</p>
        </div>
      )}

      {/* Use the BacklinksAnalysisSection component */}
      <BacklinksAnalysisSection 
        backlinksData={backlinksData} 
        analysisId={analysisId} 
      />
    </div>
  );
}