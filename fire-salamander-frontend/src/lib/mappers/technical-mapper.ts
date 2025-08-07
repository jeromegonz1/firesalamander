/**
 * Fire Salamander - Technical Analysis Mapper
 * Maps backend data to TechnicalAnalysis interface with strict validation
 * Lead Tech quality implementation
 */

import {
  TechnicalAnalysis,
  PageAnalysis,
  PageSEOElement,
  PageHeadings,
  PageImage,
  PageLinks,
  SchemaMarkup,
  DuplicateContent,
  BrokenLink,
  RedirectChain,
  GlobalIssues,
  RobotsTxtAnalysis,
  SitemapAnalysis,
  CrawlBudget,
  IssueType,
  PageStatus,
  HeadingStructure,
  SimilarityLevel,
  PageAnalysisTableRow
} from '@/types/technical-analysis';

/**
 * Maps backend analysis data to TechnicalAnalysis interface
 * @param analysis Backend analysis data
 * @returns Complete TechnicalAnalysis object
 */
export function mapBackendToTechnicalAnalysis(analysis: any): TechnicalAnalysis {
  try {
    const resultData = JSON.parse(analysis.result_data || '{}');
    const techData = resultData.technical_analysis || {};
    
    return {
      pageAnalysis: mapPageAnalysis(techData.pages || []),
      globalIssues: mapGlobalIssues(techData.global_issues || {}),
      crawlability: mapCrawlability(techData.crawl_info || {}),
      metrics: mapMetrics(techData.metrics || {}, techData.pages || []),
      config: mapAnalysisConfig(techData.analysis_config || {}),
      status: mapAnalysisStatus(analysis, techData.status || {})
    };
  } catch (error) {
    console.error('Error mapping technical analysis:', error);
    return createEmptyTechnicalAnalysis();
  }
}

/**
 * Maps page analysis data with strict validation
 */
function mapPageAnalysis(pages: any[]): PageAnalysis[] {
  return pages.map((page, index) => {
    try {
      return {
        url: validateUrl(page.url),
        statusCode: validateStatusCode(page.status_code),
        loadTime: validatePositiveNumber(page.load_time, 0),
        size: validatePositiveNumber(page.page_size, 0),
        lastCrawled: validateISODate(page.last_crawled),
        depth: validatePositiveNumber(page.depth, 0),
        
        // SEO Elements
        title: mapSEOElement(page.seo_elements?.title),
        metaDescription: mapSEOElement(page.seo_elements?.meta_description),
        headings: mapHeadings(page.seo_elements?.headings),
        
        // Technical Elements
        canonical: page.technical_elements?.canonical || '',
        robots: page.technical_elements?.robots || '',
        metaRobots: page.technical_elements?.meta_robots,
        lang: page.technical_elements?.lang,
        hreflang: mapHreflang(page.technical_elements?.hreflang),
        
        // Structured Data
        schema: mapSchemaMarkup(page.schema || []),
        openGraph: page.open_graph || {},
        twitterCard: page.twitter_card,
        
        // Media & Resources
        images: mapImages(page.images || []),
        videos: page.videos,
        
        // Links
        links: mapLinks(page.links || {}),
        
        // Performance
        performance: mapPerformance(page.performance),
        
        // Mobile
        mobile: mapMobile(page.mobile)
      };
    } catch (error) {
      console.error(`Error mapping page ${index}:`, error);
      return createEmptyPageAnalysis(page.url || `page-${index}`);
    }
  });
}

/**
 * Maps SEO elements with validation
 */
function mapSEOElement(element: any): PageSEOElement {
  if (!element) {
    return {
      content: '',
      length: 0,
      hasKeyword: false,
      issues: [IssueType.TITLE_MISSING],
      recommendations: ['Ajouter un élément SEO manquant']
    };
  }
  
  const content = String(element.content || '');
  const issues = validateIssueTypes(element.issues || []);
  
  return {
    content,
    length: content.length,
    hasKeyword: Boolean(element.has_keyword),
    issues,
    recommendations: generateSEORecommendations(content, issues)
  };
}

/**
 * Maps headings structure with validation
 */
function mapHeadings(headings: any): PageHeadings {
  if (!headings) {
    return {
      h1: [],
      h2: [],
      h3: [],
      structure: HeadingStructure.BAD,
      issues: [IssueType.H1_MISSING],
      recommendations: ['Ajouter une structure de titres appropriée']
    };
  }
  
  const h1Array = Array.isArray(headings.h1) ? headings.h1 : [];
  const h2Array = Array.isArray(headings.h2) ? headings.h2 : [];
  const h3Array = Array.isArray(headings.h3) ? headings.h3 : [];
  
  const issues = validateIssueTypes(headings.issues || []);
  const structure = validateHeadingStructure(headings.structure);
  
  // Add automatic issue detection
  const detectedIssues = [...issues];
  if (h1Array.length === 0) {
    detectedIssues.push(IssueType.H1_MISSING);
  }
  if (h1Array.length > 1) {
    detectedIssues.push(IssueType.H1_MULTIPLE);
  }
  
  return {
    h1: h1Array,
    h2: h2Array,
    h3: h3Array,
    h4: headings.h4,
    h5: headings.h5,
    h6: headings.h6,
    structure,
    issues: [...new Set(detectedIssues)], // Remove duplicates
    recommendations: generateHeadingRecommendations(h1Array, h2Array, h3Array, detectedIssues)
  };
}

/**
 * Maps images with validation
 */
function mapImages(images: any[]): PageImage[] {
  return images.map((img, index) => {
    try {
      const issues = validateIssueTypes(img.issues || []);
      
      // Auto-detect alt text issues
      const detectedIssues = [...issues];
      if (!img.alt || img.alt.trim() === '') {
        detectedIssues.push(IssueType.IMAGE_NO_ALT);
      }
      
      // Auto-detect large image issues
      const size = validatePositiveNumber(img.size, 0);
      if (size > 1000000) { // > 1MB
        detectedIssues.push(IssueType.IMAGE_LARGE_SIZE);
      }
      
      return {
        src: img.src || '',
        alt: img.alt || '',
        title: img.title,
        size,
        dimensions: img.dimensions,
        format: img.format,
        loading: img.loading,
        issues: [...new Set(detectedIssues)],
        recommendations: generateImageRecommendations(img, detectedIssues)
      };
    } catch (error) {
      console.error(`Error mapping image ${index}:`, error);
      return {
        src: img?.src || '',
        alt: '',
        size: 0,
        issues: [IssueType.IMAGE_NO_ALT],
        recommendations: ['Corriger les informations de l\'image']
      };
    }
  });
}

/**
 * Maps links with validation
 */
function mapLinks(links: any): PageLinks {
  const brokenLinks = Array.isArray(links.broken) ? links.broken : [];
  
  return {
    internal: validatePositiveNumber(links.internal, 0),
    external: validatePositiveNumber(links.external, 0),
    broken: brokenLinks.map(mapBrokenLink),
    nofollow: validatePositiveNumber(links.nofollow, 0),
    totalLinks: validatePositiveNumber(links.total, 0)
  };
}

/**
 * Maps broken link with validation
 */
function mapBrokenLink(link: any): PageLinks['broken'][0] {
  return {
    url: link.url || '',
    status: validateStatusCode(link.status),
    anchorText: link.anchor_text || '',
    position: validateLinkPosition(link.position)
  };
}

/**
 * Maps schema markup with validation
 */
function mapSchemaMarkup(schemas: any[]): SchemaMarkup[] {
  return schemas.map((schema, index) => {
    try {
      return {
        type: schema.type || 'Unknown',
        valid: Boolean(schema.valid),
        errors: Array.isArray(schema.errors) ? schema.errors : [],
        properties: schema.properties || {},
        url: schema.url
      };
    } catch (error) {
      console.error(`Error mapping schema ${index}:`, error);
      return {
        type: 'Unknown',
        valid: false,
        errors: ['Erreur de traitement du schema'],
        properties: {}
      };
    }
  });
}

/**
 * Maps global issues with validation
 */
function mapGlobalIssues(globalIssues: any): GlobalIssues {
  return {
    duplicateContent: mapDuplicateContent(globalIssues.duplicate_content || []),
    duplicateTitles: mapDuplicateTitles(globalIssues.duplicate_titles || []),
    duplicateMeta: mapDuplicateMeta(globalIssues.duplicate_meta || []),
    missingTitles: Array.isArray(globalIssues.missing_titles) ? globalIssues.missing_titles : [],
    missingMeta: Array.isArray(globalIssues.missing_meta) ? globalIssues.missing_meta : [],
    brokenLinks: mapGlobalBrokenLinks(globalIssues.broken_links || []),
    orphanPages: Array.isArray(globalIssues.orphan_pages) ? globalIssues.orphan_pages : [],
    redirectChains: mapRedirectChains(globalIssues.redirect_chains || []),
    largePages: mapLargePages(globalIssues.large_pages || []),
    slowPages: mapSlowPages(globalIssues.slow_pages || [])
  };
}

/**
 * Maps duplicate content with validation
 */
function mapDuplicateContent(duplicates: any[]): DuplicateContent[] {
  return duplicates.map((dup, index) => {
    try {
      const similarity = validatePercentage(dup.similarity);
      
      return {
        pages: Array.isArray(dup.pages) ? dup.pages : [],
        similarity,
        level: determineSimilarityLevel(similarity),
        contentHash: dup.content_hash,
        affectedElements: Array.isArray(dup.affected_elements) ? dup.affected_elements : []
      };
    } catch (error) {
      console.error(`Error mapping duplicate content ${index}:`, error);
      return {
        pages: [],
        similarity: 0,
        level: SimilarityLevel.LOW,
        affectedElements: []
      };
    }
  });
}

/**
 * Maps crawlability information
 */
function mapCrawlability(crawlInfo: any): TechnicalAnalysis['crawlability'] {
  return {
    robotsTxt: mapRobotsTxt(crawlInfo.robots_txt || {}),
    sitemap: mapSitemap(crawlInfo.sitemap || {}),
    crawlBudget: mapCrawlBudget(crawlInfo.crawl_budget || {})
  };
}

/**
 * Maps robots.txt analysis
 */
function mapRobotsTxt(robotsData: any): RobotsTxtAnalysis {
  return {
    exists: Boolean(robotsData.exists),
    url: robotsData.url,
    valid: Boolean(robotsData.valid),
    userAgents: Array.isArray(robotsData.user_agents) ? robotsData.user_agents : [],
    disallowedPaths: Array.isArray(robotsData.disallowed_paths) ? robotsData.disallowed_paths : [],
    allowedPaths: Array.isArray(robotsData.allowed_paths) ? robotsData.allowed_paths : [],
    sitemapUrls: Array.isArray(robotsData.sitemap_urls) ? robotsData.sitemap_urls : [],
    crawlDelay: robotsData.crawl_delay,
    issues: Array.isArray(robotsData.issues) ? robotsData.issues : [],
    recommendations: generateRobotsRecommendations(robotsData)
  };
}

/**
 * Maps sitemap analysis
 */
function mapSitemap(sitemapData: any): SitemapAnalysis {
  return {
    exists: Boolean(sitemapData.exists),
    url: sitemapData.url || '',
    valid: Boolean(sitemapData.valid),
    format: sitemapData.format || 'unknown',
    pagesInSitemap: validatePositiveNumber(sitemapData.pages_in_sitemap, 0),
    lastModified: sitemapData.last_modified,
    images: validatePositiveNumber(sitemapData.images, 0),
    videos: validatePositiveNumber(sitemapData.videos, 0),
    issues: Array.isArray(sitemapData.issues) ? sitemapData.issues : [],
    recommendations: generateSitemapRecommendations(sitemapData),
    indexSitemaps: sitemapData.index_sitemaps
  };
}

/**
 * Maps crawl budget analysis
 */
function mapCrawlBudget(budgetData: any): CrawlBudget {
  const totalPages = validatePositiveNumber(budgetData.total_pages, 0);
  const crawlablePages = validatePositiveNumber(budgetData.crawlable_pages, 0);
  const blockedPages = validatePositiveNumber(budgetData.blocked_pages, 0);
  
  return {
    totalPages,
    crawlablePages,
    blockedPages,
    indexedPages: budgetData.indexed_pages,
    pagesPerLevel: budgetData.pages_per_level || {},
    averageCrawlTime: validatePositiveNumber(budgetData.average_crawl_time, 0),
    crawlEfficiency: totalPages > 0 ? (crawlablePages / totalPages) * 100 : 0
  };
}

// ==================== UTILITY FUNCTIONS ====================

/**
 * Validates URL format
 */
function validateUrl(url: string): string {
  if (!url) return '';
  
  try {
    new URL(url);
    return url;
  } catch {
    return url; // Return as-is if invalid, let the UI handle display
  }
}

/**
 * Validates HTTP status code
 */
function validateStatusCode(code: any): PageStatus {
  const numCode = Number(code);
  if (isNaN(numCode) || numCode < 100 || numCode > 599) {
    return PageStatus.SERVER_ERROR;
  }
  return numCode as PageStatus;
}

/**
 * Validates positive numbers with fallback
 */
function validatePositiveNumber(value: any, fallback: number): number {
  const num = Number(value);
  return isNaN(num) || num < 0 ? fallback : num;
}

/**
 * Validates percentage (0-100)
 */
function validatePercentage(value: any): number {
  const num = Number(value);
  if (isNaN(num)) return 0;
  return Math.max(0, Math.min(100, num));
}

/**
 * Validates ISO date string
 */
function validateISODate(dateStr: any): string {
  if (!dateStr) return new Date().toISOString();
  
  try {
    const date = new Date(dateStr);
    return date.toISOString();
  } catch {
    return new Date().toISOString();
  }
}

/**
 * Validates issue types array
 */
function validateIssueTypes(issues: any[]): IssueType[] {
  if (!Array.isArray(issues)) return [];
  
  const validIssues = Object.values(IssueType);
  return issues.filter(issue => validIssues.includes(issue as IssueType));
}

/**
 * Validates heading structure
 */
function validateHeadingStructure(structure: any): HeadingStructure {
  return structure === HeadingStructure.GOOD ? HeadingStructure.GOOD : HeadingStructure.BAD;
}

/**
 * Validates link position
 */
function validateLinkPosition(position: any): 'header' | 'footer' | 'content' | 'navigation' {
  const validPositions = ['header', 'footer', 'content', 'navigation'];
  return validPositions.includes(position) ? position : 'content';
}

/**
 * Determines similarity level from percentage
 */
function determineSimilarityLevel(similarity: number): SimilarityLevel {
  if (similarity >= 80) return SimilarityLevel.HIGH;
  if (similarity >= 50) return SimilarityLevel.MEDIUM;
  return SimilarityLevel.LOW;
}

// ==================== RECOMMENDATION GENERATORS ====================

/**
 * Generates SEO element recommendations
 */
function generateSEORecommendations(content: string, issues: IssueType[]): string[] {
  const recommendations: string[] = [];
  
  if (issues.includes(IssueType.TITLE_MISSING)) {
    recommendations.push('Ajouter un titre à la page');
  }
  if (issues.includes(IssueType.TITLE_TOO_SHORT)) {
    recommendations.push('Rallonger le titre (30-60 caractères recommandés)');
  }
  if (issues.includes(IssueType.TITLE_TOO_LONG)) {
    recommendations.push('Raccourcir le titre (60 caractères maximum)');
  }
  if (issues.includes(IssueType.META_MISSING)) {
    recommendations.push('Ajouter une meta description');
  }
  if (issues.includes(IssueType.META_TOO_SHORT)) {
    recommendations.push('Rallonger la meta description (120-160 caractères)');
  }
  if (issues.includes(IssueType.META_TOO_LONG)) {
    recommendations.push('Raccourcir la meta description (160 caractères maximum)');
  }
  
  return recommendations;
}

/**
 * Generates heading recommendations
 */
function generateHeadingRecommendations(
  h1: string[], 
  h2: string[], 
  h3: string[], 
  issues: IssueType[]
): string[] {
  const recommendations: string[] = [];
  
  if (issues.includes(IssueType.H1_MISSING)) {
    recommendations.push('Ajouter un titre H1 unique à la page');
  }
  if (issues.includes(IssueType.H1_MULTIPLE)) {
    recommendations.push('Utiliser un seul titre H1 par page');
  }
  if (issues.includes(IssueType.HEADING_SKIP)) {
    recommendations.push('Respecter la hiérarchie des titres (H1 > H2 > H3)');
  }
  
  if (h2.length === 0 && h3.length > 0) {
    recommendations.push('Ajouter des titres H2 avant les H3');
  }
  
  return recommendations;
}

/**
 * Generates image recommendations
 */
function generateImageRecommendations(img: any, issues: IssueType[]): string[] {
  const recommendations: string[] = [];
  
  if (issues.includes(IssueType.IMAGE_NO_ALT)) {
    recommendations.push('Ajouter un texte alternatif descriptif');
  }
  if (issues.includes(IssueType.IMAGE_LARGE_SIZE)) {
    recommendations.push('Optimiser la taille de l\'image (< 1MB)');
  }
  
  if (img.format && !['webp', 'avif'].includes(img.format.toLowerCase())) {
    recommendations.push('Utiliser des formats modernes (WebP, AVIF)');
  }
  
  return recommendations;
}

/**
 * Generates robots.txt recommendations
 */
function generateRobotsRecommendations(robotsData: any): string[] {
  const recommendations: string[] = [];
  
  if (!robotsData.exists) {
    recommendations.push('Créer un fichier robots.txt');
  }
  if (!robotsData.valid) {
    recommendations.push('Corriger la syntax du fichier robots.txt');
  }
  if (!robotsData.sitemap_urls || robotsData.sitemap_urls.length === 0) {
    recommendations.push('Ajouter l\'URL du sitemap dans robots.txt');
  }
  
  return recommendations;
}

/**
 * Generates sitemap recommendations
 */
function generateSitemapRecommendations(sitemapData: any): string[] {
  const recommendations: string[] = [];
  
  if (!sitemapData.exists) {
    recommendations.push('Créer un sitemap XML');
  }
  if (!sitemapData.valid) {
    recommendations.push('Corriger les erreurs du sitemap');
  }
  if (sitemapData.pages_in_sitemap === 0) {
    recommendations.push('Ajouter des URLs au sitemap');
  }
  
  return recommendations;
}

// ==================== HELPER MAPPERS ====================

function mapHreflang(hreflang: any[]): PageAnalysis['hreflang'] {
  if (!Array.isArray(hreflang)) return undefined;
  
  return hreflang.map(item => ({
    lang: item.lang || '',
    url: item.url || ''
  }));
}

function mapPerformance(performance: any): PageAnalysis['performance'] {
  if (!performance) return undefined;
  
  return {
    fcp: validatePositiveNumber(performance.fcp, 0),
    lcp: validatePositiveNumber(performance.lcp, 0),
    cls: validatePositiveNumber(performance.cls, 0),
    ttfb: validatePositiveNumber(performance.ttfb, 0)
  };
}

function mapMobile(mobile: any): PageAnalysis['mobile'] {
  if (!mobile) return undefined;
  
  return {
    viewport: Boolean(mobile.viewport),
    responsive: Boolean(mobile.responsive),
    touchTargets: validatePositiveNumber(mobile.touch_targets, 0),
    fontSizes: Array.isArray(mobile.font_sizes) ? mobile.font_sizes : []
  };
}

function mapDuplicateTitles(duplicates: any[]): GlobalIssues['duplicateTitles'] {
  return duplicates.map(dup => ({
    title: dup.title || '',
    pages: Array.isArray(dup.pages) ? dup.pages : [],
    count: validatePositiveNumber(dup.count, 0)
  }));
}

function mapDuplicateMeta(duplicates: any[]): GlobalIssues['duplicateMeta'] {
  return duplicates.map(dup => ({
    meta: dup.meta || '',
    pages: Array.isArray(dup.pages) ? dup.pages : [],
    count: validatePositiveNumber(dup.count, 0)
  }));
}

function mapGlobalBrokenLinks(links: any[]): BrokenLink[] {
  return links.map(link => ({
    from: link.from || '',
    to: link.to || '',
    status: validateStatusCode(link.status),
    anchorText: link.anchor_text || '',
    linkType: link.link_type === 'external' ? 'external' : 'internal',
    position: validateLinkPosition(link.position),
    lastChecked: validateISODate(link.last_checked)
  }));
}

function mapRedirectChains(chains: any[]): RedirectChain[] {
  return chains.map(chain => ({
    chain: Array.isArray(chain.chain) ? chain.chain : [],
    finalUrl: chain.final_url || '',
    totalHops: validatePositiveNumber(chain.total_hops, 0),
    totalTime: validatePositiveNumber(chain.total_time, 0),
    statusCodes: Array.isArray(chain.status_codes) ? chain.status_codes : []
  }));
}

function mapLargePages(pages: any[]): GlobalIssues['largePages'] {
  return pages.map(page => ({
    url: page.url || '',
    size: validatePositiveNumber(page.size, 0)
  }));
}

function mapSlowPages(pages: any[]): GlobalIssues['slowPages'] {
  return pages.map(page => ({
    url: page.url || '',
    loadTime: validatePositiveNumber(page.load_time, 0)
  }));
}

function mapMetrics(metricsData: any, pages: any[]): TechnicalAnalysis['metrics'] {
  const totalIssues = pages.reduce((total, page) => {
    const pageIssues = [
      ...(page.seo_elements?.title?.issues || []),
      ...(page.seo_elements?.meta_description?.issues || []),
      ...(page.seo_elements?.headings?.issues || []),
      ...(page.images?.flatMap((img: any) => img.issues || []) || [])
    ];
    return total + pageIssues.length;
  }, 0);
  
  const issuesByType = pages.reduce((acc, page) => {
    const pageIssues = [
      ...(page.seo_elements?.title?.issues || []),
      ...(page.seo_elements?.meta_description?.issues || []),
      ...(page.seo_elements?.headings?.issues || []),
      ...(page.images?.flatMap((img: any) => img.issues || []) || [])
    ];
    
    pageIssues.forEach((issue: IssueType) => {
      acc[issue] = (acc[issue] || 0) + 1;
    });
    
    return acc;
  }, {} as Record<IssueType, number>);
  
  return {
    totalPages: validatePositiveNumber(metricsData.total_pages, pages.length),
    pagesWithIssues: validatePositiveNumber(metricsData.pages_with_issues, 0),
    avgLoadTime: validatePositiveNumber(metricsData.avg_load_time, 0),
    avgPageSize: validatePositiveNumber(metricsData.avg_page_size, 0),
    totalIssues: validatePositiveNumber(metricsData.total_issues, totalIssues),
    issuesByType,
    healthScore: validatePercentage(metricsData.health_score)
  };
}

function mapAnalysisConfig(configData: any): TechnicalAnalysis['config'] {
  return {
    maxDepth: validatePositiveNumber(configData.max_depth, 3),
    respectRobots: Boolean(configData.respect_robots ?? true),
    userAgent: configData.user_agent || 'Fire-Salamander-Bot/1.0',
    crawlDelay: validatePositiveNumber(configData.crawl_delay, 1000),
    maxPages: configData.max_pages,
    includedPaths: configData.included_paths,
    excludedPaths: configData.excluded_paths
  };
}

function mapAnalysisStatus(analysis: any, statusData: any): TechnicalAnalysis['status'] {
  return {
    analysisDate: validateISODate(analysis.analyzed_at),
    crawlDuration: validatePositiveNumber(statusData.crawl_duration, 0),
    crawlStatus: statusData.crawl_status || 'completed',
    lastUpdate: validateISODate(statusData.last_update || analysis.analyzed_at),
    version: statusData.version || '1.0'
  };
}

// ==================== EMPTY OBJECT CREATORS ====================

function createEmptyTechnicalAnalysis(): TechnicalAnalysis {
  return {
    pageAnalysis: [],
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
        exists: false,
        valid: false,
        userAgents: [],
        disallowedPaths: [],
        allowedPaths: [],
        sitemapUrls: [],
        issues: ['Données indisponibles']
      },
      sitemap: {
        exists: false,
        url: '',
        valid: false,
        format: 'unknown',
        pagesInSitemap: 0,
        images: 0,
        videos: 0,
        issues: ['Données indisponibles']
      },
      crawlBudget: {
        totalPages: 0,
        crawlablePages: 0,
        blockedPages: 0,
        pagesPerLevel: {},
        averageCrawlTime: 0,
        crawlEfficiency: 0
      }
    },
    metrics: {
      totalPages: 0,
      pagesWithIssues: 0,
      avgLoadTime: 0,
      avgPageSize: 0,
      totalIssues: 0,
      issuesByType: {},
      healthScore: 0
    },
    config: {
      maxDepth: 3,
      respectRobots: true,
      userAgent: 'Fire-Salamander-Bot/1.0',
      crawlDelay: 1000
    },
    status: {
      analysisDate: new Date().toISOString(),
      crawlDuration: 0,
      crawlStatus: 'failed',
      lastUpdate: new Date().toISOString(),
      version: '1.0'
    }
  };
}

function createEmptyPageAnalysis(url: string): PageAnalysis {
  return {
    url,
    statusCode: PageStatus.SERVER_ERROR,
    loadTime: 0,
    size: 0,
    lastCrawled: new Date().toISOString(),
    depth: 0,
    title: {
      content: '',
      length: 0,
      hasKeyword: false,
      issues: [IssueType.TITLE_MISSING],
      recommendations: ['Données de page indisponibles']
    },
    metaDescription: {
      content: '',
      length: 0,
      hasKeyword: false,
      issues: [IssueType.META_MISSING]
    },
    headings: {
      h1: [],
      h2: [],
      h3: [],
      structure: HeadingStructure.BAD,
      issues: [IssueType.H1_MISSING]
    },
    canonical: '',
    robots: '',
    schema: [],
    openGraph: {},
    images: [],
    links: {
      internal: 0,
      external: 0,
      broken: [],
      nofollow: 0,
      totalLinks: 0
    }
  };
}

/**
 * Transforms PageAnalysis to table row format
 */
export function mapPageAnalysisToTableRow(page: PageAnalysis): PageAnalysisTableRow {
  const allIssues = [
    ...page.title.issues,
    ...page.metaDescription.issues,
    ...page.headings.issues,
    ...page.images.flatMap(img => img.issues),
    ...(page.links.broken.length > 0 ? [IssueType.LINK_BROKEN] : [])
  ];
  
  const uniqueIssues = [...new Set(allIssues)];
  const issueCount = uniqueIssues.length;
  
  let overallHealth: 'good' | 'warning' | 'error' = 'good';
  if (issueCount > 5) {
    overallHealth = 'error';
  } else if (issueCount > 0) {
    overallHealth = 'warning';
  }
  
  return {
    ...page,
    issueCount,
    issueTypes: uniqueIssues,
    overallHealth
  };
}