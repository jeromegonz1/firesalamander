/**
 * Fire Salamander - Overview Data Mapper
 * Maps backend data to AnalysisOverview interface
 */

import { AnalysisOverview } from '@/types/analysis-overview';

export function mapBackendToOverview(analysis: any): AnalysisOverview {
  const resultData = JSON.parse(analysis.result_data || '{}');
  const seoAnalysis = resultData.seo_analysis || {};
  const categoryScores = resultData.category_scores || {};
  const unifiedMetrics = resultData.unified_metrics || {};
  const tagAnalysis = seoAnalysis.tag_analysis || {};
  const technicalAudit = seoAnalysis.technical_audit || {};
  const performanceMetrics = seoAnalysis.performance_metrics || {};
  const coreWebVitals = performanceMetrics.core_web_vitals || {};

  // Helper functions
  const formatTime = (ms: number) => {
    if (ms < 1000) return `${ms}ms`;
    return `${(ms / 1000).toFixed(1)}s`;
  };

  const formatBytes = (bytes: number) => {
    if (bytes < 1024) return `${bytes}B`;
    if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)}KB`;
    return `${(bytes / 1024 / 1024).toFixed(1)}MB`;
  };

  const getAnalysisType = (type: string): 'Quick Scan' | 'AI Boost Scan' | 'Custom Scan' => {
    switch (type) {
      case 'quick': return 'Quick Scan';
      case 'full': return 'AI Boost Scan';
      case 'custom': return 'Custom Scan';
      default: return 'Quick Scan';
    }
  };

  const formatDate = (dateStr: string): string => {
    const date = new Date(dateStr);
    const options: Intl.DateTimeFormatOptions = {
      day: 'numeric',
      month: 'long',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    };
    return date.toLocaleDateString('fr-FR', options);
  };

  // Calculate SEO sub-scores
  const calculateTitleScore = () => {
    const title = tagAnalysis.title || {};
    let score = 100;
    if (!title.present) return 0;
    if (!title.optimal_length) score -= 40;
    if (!title.has_keywords) score -= 20;
    return Math.max(0, score);
  };

  const calculateMetaScore = () => {
    const meta = tagAnalysis.meta_description || {};
    let score = 100;
    if (!meta.present) return 0;
    if (!meta.optimal_length) score -= 30;
    if (!meta.has_call_to_action) score -= 10;
    return Math.max(0, score);
  };

  const calculateHeadingScore = () => {
    const headings = tagAnalysis.headings || {};
    if (headings.h1_count === 0) return 0;
    if (headings.h1_count > 1) return 50;
    if (!headings.has_hierarchy) return 70;
    return 100;
  };

  // Count issues
  const countIssues = () => {
    let critical = 0;
    let warnings = 0;
    let notices = 0;

    // Critical issues
    if (!tagAnalysis.headings?.h1_count) critical++;
    if (technicalAudit.security?.mixed_content) critical++;
    if (coreWebVitals.cls?.score === 'poor') critical++;
    if (coreWebVitals.lcp?.score === 'poor') critical++;

    // Warnings
    if (!performanceMetrics.has_caching) warnings++;
    if (!performanceMetrics.optimized_images) warnings++;
    if (coreWebVitals.cls?.score === 'needs-improvement') warnings++;
    if (!tagAnalysis.title?.optimal_length) warnings++;

    // Notices
    if (!tagAnalysis.meta_description?.has_call_to_action) notices++;
    if (tagAnalysis.images?.alt_text_coverage < 1) notices++;

    return { critical, warnings, notices };
  };

  const issues = countIssues();

  // Map Core Web Vitals with proper units
  const mapCoreWebVitals = () => {
    return {
      lcp: {
        value: coreWebVitals.lcp?.value ? coreWebVitals.lcp.value / 1000 : 0,
        unit: 's',
        score: coreWebVitals.lcp?.score || 'poor'
      },
      fid: {
        value: coreWebVitals.fid?.value || 0,
        unit: 'ms',
        score: coreWebVitals.fid?.score || 'poor'
      },
      cls: {
        value: parseFloat((coreWebVitals.cls?.value || 0).toFixed(2)),
        unit: '',
        score: coreWebVitals.cls?.score || 'poor'
      },
      ttfb: {
        value: coreWebVitals.ttfb ? Math.round(coreWebVitals.ttfb / 1000000) : 0,
        unit: 'ms',
        score: coreWebVitals.ttfb < 800000000 ? 'good' : coreWebVitals.ttfb < 1800000000 ? 'needs-improvement' : 'poor'
      },
      inp: {
        value: null,
        unit: 'ms',
        score: 'n/a' as const
      }
    };
  };

  // Build top issues
  const buildTopIssues = () => {
    const topIssues = [];

    if (!tagAnalysis.headings?.h1_count) {
      topIssues.push({
        id: 'missing-h1',
        title: 'Balise H1 manquante',
        category: 'SEO' as const,
        severity: 'critical' as const,
        pagesAffected: 1,
        description: 'La page n\'a pas de balise H1, essentielle pour le SEO'
      });
    }

    if (tagAnalysis.images?.alt_text_coverage < 1) {
      const missingAlt = tagAnalysis.images.total_images - tagAnalysis.images.images_with_alt;
      topIssues.push({
        id: 'missing-alt-text',
        title: 'Images sans texte alternatif',
        category: 'Accessibility' as const,
        severity: 'warning' as const,
        pagesAffected: 1,
        description: `${missingAlt} images sur ${tagAnalysis.images.total_images} n'ont pas de texte alternatif`
      });
    }

    if (!performanceMetrics.has_caching) {
      topIssues.push({
        id: 'no-caching',
        title: 'Cache HTTP non configuré',
        category: 'Performance' as const,
        severity: 'warning' as const,
        pagesAffected: 1,
        description: 'Aucune stratégie de cache n\'est configurée, impactant les performances'
      });
    }

    return topIssues.slice(0, 5); // Top 5 issues
  };

  return {
    metadata: {
      url: analysis.url,
      analyzedAt: formatDate(analysis.created_at),
      analysisType: getAnalysisType(analysis.analysis_type),
      processingTime: formatTime(analysis.processing_time || 0),
      domain: analysis.domain,
      protocol: analysis.url.startsWith('https') ? 'https' : 'http'
    },

    scores: {
      global: Math.round((analysis.overall_score || 0) * 100),
      seo: {
        score: Math.round(categoryScores.basics || 0),
        trend: 'stable', // TODO: Calculate from historical data
        details: {
          titleOptimization: calculateTitleScore(),
          metaOptimization: calculateMetaScore(),
          headingStructure: calculateHeadingScore(),
          keywordUsage: tagAnalysis.title?.has_keywords ? 80 : 0,
          contentQuality: Math.round((unifiedMetrics.content_quality_score || 0) * 100)
        }
      },
      technical: {
        score: Math.round(categoryScores.technical || 0),
        trend: 'stable',
        details: {
          crawlability: Math.round((technicalAudit.crawlability?.crawlability_score || 0) * 100),
          indexability: Math.round((technicalAudit.indexability?.indexability_score || 0) * 100),
          siteSpeed: Math.round(categoryScores.performance || 0),
          mobile: Math.round((technicalAudit.mobile?.mobile_score || 0) * 100),
          security: Math.round((technicalAudit.security?.security_score || 0) * 100)
        }
      },
      performance: {
        score: Math.round(categoryScores.performance || 0),
        trend: 'stable',
        coreWebVitals: mapCoreWebVitals()
      },
      accessibility: {
        score: Math.round((technicalAudit.accessibility?.score || 0) * 100),
        wcagLevel: technicalAudit.accessibility?.score > 0.9 ? 'AA' : 
                   technicalAudit.accessibility?.score > 0.7 ? 'A' : 'Fail'
      }
    },

    metrics: {
      totalPages: 1, // TODO: Get from crawl data
      pagesAnalyzed: 1,
      pagesWithErrors: issues.critical > 0 ? 1 : 0,
      
      issues: {
        total: issues.critical + issues.warnings + issues.notices,
        critical: issues.critical,
        warnings: issues.warnings,
        notices: issues.notices,
        passedChecks: 10 // TODO: Calculate from all checks
      },
      
      performance: {
        avgLoadTime: Math.round((performanceMetrics.load_time || 0) / 1000000),
        pageSize: performanceMetrics.page_size || 0,
        requests: (performanceMetrics.resource_counts?.images || 0) + 
                  (performanceMetrics.resource_counts?.scripts || 0) +
                  (performanceMetrics.resource_counts?.stylesheets || 0),
        avgTimeToFirstByte: Math.round((coreWebVitals.ttfb || 0) / 1000000)
      },
      
      resources: {
        images: {
          total: tagAnalysis.images?.total_images || 0,
          optimized: tagAnalysis.images?.optimized_formats || 0,
          missingAlt: (tagAnalysis.images?.total_images || 0) - (tagAnalysis.images?.images_with_alt || 0),
          oversized: 0 // TODO: Calculate from image analysis
        },
        scripts: {
          total: performanceMetrics.resource_counts?.scripts || 0,
          minified: performanceMetrics.minified_resources ? performanceMetrics.resource_counts?.scripts : 0,
          external: 0, // TODO: Calculate
          renderBlocking: 0 // TODO: Calculate
        },
        stylesheets: {
          total: performanceMetrics.resource_counts?.stylesheets || 0,
          minified: performanceMetrics.minified_resources ? performanceMetrics.resource_counts?.stylesheets : 0,
          external: 0,
          renderBlocking: 0
        },
        fonts: performanceMetrics.resource_counts?.fonts || 0,
        videos: performanceMetrics.resource_counts?.videos || 0
      },
      
      links: {
        internal: {
          total: tagAnalysis.links?.internal_links || 0,
          unique: tagAnalysis.links?.internal_links || 0, // TODO: Calculate unique
          broken: 0
        },
        external: {
          total: tagAnalysis.links?.external_links || 0,
          unique: tagAnalysis.links?.external_links || 0,
          broken: tagAnalysis.links?.broken_links || 0,
          nofollow: tagAnalysis.links?.nofollow_links || 0
        }
      },
      
      seo: {
        metaTags: {
          title: tagAnalysis.title?.present || false,
          description: tagAnalysis.meta_description?.present || false,
          keywords: false, // Not in current data
          ogTags: tagAnalysis.meta_tags?.has_og_tags || false,
          twitterCards: tagAnalysis.meta_tags?.has_twitter_card || false
        },
        headings: {
          h1Count: tagAnalysis.headings?.h1_count || 0,
          h2Count: 0, // TODO: Extract from heading_structure
          h3Count: 0,
          missingH1: !tagAnalysis.headings?.h1_count,
          multipleH1: (tagAnalysis.headings?.h1_count || 0) > 1
        },
        schema: {
          present: tagAnalysis.microdata?.has_json_ld || false,
          types: tagAnalysis.microdata?.json_ld_types || []
        }
      }
    },

    topIssues: buildTopIssues()
  };
}