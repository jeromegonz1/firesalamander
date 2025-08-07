"use client";

import React, { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { TechnicalAnalysisSection } from "@/components/analysis/technical-analysis-section";
import { mapBackendToTechnicalAnalysis } from "@/lib/mappers/technical-mapper";
import { TechnicalAnalysis, PageStatus, HeadingStructure, IssueType } from "@/types/technical-analysis";

export default function TechnicalDebugPage() {
  const params = useParams();
  const [technicalData, setTechnicalData] = useState<TechnicalAnalysis | null>(null);
  const [loading, setLoading] = useState(true);
  const analysisId = params.id as string;

  const createMinimalMockData = (): TechnicalAnalysis => {
    return {
      pageAnalysis: [
        {
          url: 'https://example.com',
          statusCode: 200 as PageStatus,
          loadTime: 1200,
          size: 45678,
          lastCrawled: new Date().toISOString(),
          depth: 0,
          title: {
            content: 'Example Page Title',
            length: 18,
            hasKeyword: true,
            issues: [],
            recommendations: []
          },
          metaDescription: {
            content: 'This is a good meta description',
            length: 31,
            hasKeyword: false,
            issues: [IssueType.META_TOO_SHORT],
            recommendations: ['Rallonger la meta description']
          },
          headings: {
            h1: ['Main Heading'],
            h2: ['Section 1', 'Section 2'],
            h3: ['Subsection'],
            structure: HeadingStructure.GOOD,
            issues: [],
            recommendations: []
          },
          canonical: 'https://example.com',
          robots: 'index,follow',
          schema: [],
          openGraph: {},
          images: [],
          links: {
            internal: 15,
            external: 5,
            broken: [],
            nofollow: 2,
            totalLinks: 20
          }
        }
      ],
      globalIssues: {
        duplicateContent: [],
        duplicateTitles: [],
        duplicateMeta: [],
        missingTitles: [],
        missingMeta: [],
        brokenLinks: [],
        orphanPages: [],
        redirectChains: [],
        largePages: [],
        slowPages: []
      },
      crawlability: {
        robotsTxt: {
          exists: true,
          valid: true,
          userAgents: ['*'],
          disallowedPaths: ['/admin'],
          allowedPaths: [],
          sitemapUrls: ['https://example.com/sitemap.xml'],
          issues: []
        },
        sitemap: {
          exists: true,
          url: 'https://example.com/sitemap.xml',
          valid: true,
          format: 'xml',
          pagesInSitemap: 25,
          images: 0,
          videos: 0,
          issues: []
        },
        crawlBudget: {
          totalPages: 30,
          crawlablePages: 25,
          blockedPages: 5,
          pagesPerLevel: { 0: 1, 1: 10, 2: 14 },
          averageCrawlTime: 1.2,
          crawlEfficiency: 83.3
        }
      },
      metrics: {
        totalPages: 1,
        pagesWithIssues: 0,
        avgLoadTime: 1200,
        avgPageSize: 45678,
        totalIssues: 1,
        issuesByType: {
          [IssueType.META_TOO_SHORT]: 1
        },
        healthScore: 85
      },
      config: {
        maxDepth: 3,
        respectRobots: true,
        userAgent: 'Fire-Salamander-Bot/1.0',
        crawlDelay: 1000
      },
      status: {
        analysisDate: new Date().toISOString(),
        crawlDuration: 5,
        crawlStatus: 'completed',
        lastUpdate: new Date().toISOString(),
        version: '1.0'
      }
    };
  };

  useEffect(() => {
    // Always use mock data for debugging
    const mockData = createMinimalMockData();
    setTechnicalData(mockData);
    setLoading(false);
    console.log('DEBUG: Mock data loaded', mockData);
  }, []);

  if (loading) {
    return <div className="p-4">Loading debug page...</div>;
  }

  if (!technicalData) {
    return <div className="p-4">No technical data available</div>;
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Technical Analysis Debug</h1>
      
      <div className="mb-4 p-4 bg-blue-50 rounded">
        <h2 className="font-semibold">Debug Info:</h2>
        <p>Analysis ID: {analysisId}</p>
        <p>Pages: {technicalData.pageAnalysis.length}</p>
        <p>Health Score: {technicalData.metrics.healthScore}</p>
        <p>Total Issues: {technicalData.metrics.totalIssues}</p>
      </div>

      <TechnicalAnalysisSection 
        technicalData={technicalData} 
        analysisId={analysisId} 
      />
    </div>
  );
}