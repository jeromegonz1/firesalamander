/**
 * Fire Salamander - Technical Analysis Section
 * Professional page-by-page technical analysis component
 * Lead Tech implementation with comprehensive features
 */

"use client";

import React, { useState, useEffect, useMemo } from 'react';
import { 
  TechnicalAnalysis, 
  PageAnalysisTableRow,
  PageSortBy,
  PageFilterBy,
  IssueType,
  PageStatus,
  GlobalIssueType
} from '@/types/technical-analysis';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Progress } from '@/components/ui/progress';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { 
  Search, Download, Filter, RefreshCw, ChevronDown, ChevronRight,
  AlertTriangle, CheckCircle2, XCircle, Clock, Globe, Image,
  Link2, BarChart3, Target, Shield, Zap, Eye, Settings,
  TrendingUp, TrendingDown, Activity, FileText, ExternalLink,
  ArrowRight, Copy, RotateCcw, Users, Gauge
} from 'lucide-react';

interface TechnicalAnalysisSectionProps {
  technicalData: TechnicalAnalysis;
  analysisId: string;
}

// ==================== MAIN COMPONENT ====================

export function TechnicalAnalysisSection({ 
  technicalData, 
  analysisId 
}: TechnicalAnalysisSectionProps) {
  const [expandedPages, setExpandedPages] = useState<Set<string>>(new Set());
  const [selectedPages, setSelectedPages] = useState<Set<string>>(new Set());
  const [searchQuery, setSearchQuery] = useState('');
  const [sortBy, setSortBy] = useState<PageSortBy>(PageSortBy.ISSUES);
  const [filterBy, setFilterBy] = useState<PageFilterBy>(PageFilterBy.ALL);
  const [isLoading, setIsLoading] = useState(false);

  // Convert pages to table rows with calculated metrics
  const tableRows: PageAnalysisTableRow[] = useMemo(() => {
    return technicalData.pageAnalysis.map(page => {
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
      if (page.statusCode >= 400) {
        overallHealth = 'error';
      } else if (issueCount > 5) {
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
    });
  }, [technicalData.pageAnalysis]);

  // Filter and sort pages
  const filteredAndSortedPages = useMemo(() => {
    let filtered = tableRows.filter(page => {
      // Search filter
      if (searchQuery && !page.url.toLowerCase().includes(searchQuery.toLowerCase())) {
        return false;
      }
      
      // Status filter
      switch (filterBy) {
        case PageFilterBy.ERRORS_ONLY:
          return page.statusCode >= 400;
        case PageFilterBy.WARNINGS_ONLY:
          return page.overallHealth === 'warning';
        case PageFilterBy.TITLE_ISSUES:
          return page.title.issues.length > 0;
        case PageFilterBy.META_ISSUES:
          return page.metaDescription.issues.length > 0;
        case PageFilterBy.HEADING_ISSUES:
          return page.headings.issues.length > 0;
        case PageFilterBy.IMAGE_ISSUES:
          return page.images.some(img => img.issues.length > 0);
        case PageFilterBy.LINK_ISSUES:
          return page.links.broken.length > 0;
        case PageFilterBy.SLOW_PAGES:
          return page.loadTime > 3000;
        case PageFilterBy.LARGE_PAGES:
          return page.size > 1000000; // > 1MB
        default:
          return true;
      }
    });

    // Sort pages
    return filtered.sort((a, b) => {
      switch (sortBy) {
        case PageSortBy.URL:
          return a.url.localeCompare(b.url);
        case PageSortBy.STATUS_CODE:
          return b.statusCode - a.statusCode;
        case PageSortBy.LOAD_TIME:
          return b.loadTime - a.loadTime;
        case PageSortBy.SIZE:
          return b.size - a.size;
        case PageSortBy.LAST_CRAWLED:
          return new Date(b.lastCrawled).getTime() - new Date(a.lastCrawled).getTime();
        case PageSortBy.ISSUES:
        default:
          return b.issueCount - a.issueCount;
      }
    });
  }, [tableRows, searchQuery, sortBy, filterBy]);

  const togglePageExpansion = (url: string) => {
    const newExpanded = new Set(expandedPages);
    if (newExpanded.has(url)) {
      newExpanded.delete(url);
    } else {
      newExpanded.add(url);
    }
    setExpandedPages(newExpanded);
  };

  const togglePageSelection = (url: string) => {
    const newSelected = new Set(selectedPages);
    if (newSelected.has(url)) {
      newSelected.delete(url);
    } else {
      newSelected.add(url);
    }
    setSelectedPages(newSelected);
  };

  const handleRefresh = async () => {
    setIsLoading(true);
    // TODO: Implement refresh logic
    setTimeout(() => setIsLoading(false), 2000);
  };

  return (
    <div className="space-y-6" data-testid="technical-analysis-section">
      {/* Metrics Overview */}
      <MetricsOverview metrics={technicalData.metrics} />

      {/* Main Content Tabs */}
      <Tabs defaultValue="pages" className="w-full">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="pages">Analysis by Page</TabsTrigger>
          <TabsTrigger value="issues">Global Issues</TabsTrigger>
          <TabsTrigger value="crawlability">Crawlability</TabsTrigger>
          <TabsTrigger value="overview">Overview</TabsTrigger>
        </TabsList>

        {/* Pages Analysis Tab */}
        <TabsContent value="pages" className="space-y-4">
          <PageAnalysisSection 
            pages={filteredAndSortedPages}
            expandedPages={expandedPages}
            selectedPages={selectedPages}
            searchQuery={searchQuery}
            sortBy={sortBy}
            filterBy={filterBy}
            isLoading={isLoading}
            onToggleExpansion={togglePageExpansion}
            onToggleSelection={togglePageSelection}
            onSearchChange={setSearchQuery}
            onSortChange={setSortBy}
            onFilterChange={setFilterBy}
            onRefresh={handleRefresh}
          />
        </TabsContent>

        {/* Global Issues Tab */}
        <TabsContent value="issues" className="space-y-4">
          <GlobalIssuesSection globalIssues={technicalData.globalIssues} />
        </TabsContent>

        {/* Crawlability Tab */}
        <TabsContent value="crawlability" className="space-y-4">
          <CrawlabilitySection crawlability={technicalData.crawlability} />
        </TabsContent>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-4">
          <AnalysisOverviewSection 
            analysis={technicalData}
            totalPages={filteredAndSortedPages.length}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}

// ==================== METRICS OVERVIEW ====================

function MetricsOverview({ metrics }: { metrics: TechnicalAnalysis['metrics'] }) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-4 gap-4" data-testid="metrics-overview">
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-blue-100 rounded-lg">
              <Globe className="h-5 w-5 text-blue-600" />
            </div>
            <div>
              <div className="text-2xl font-bold" data-testid="total-pages-metric">
                {metrics.totalPages}
              </div>
              <div className="text-sm text-gray-600">Total Pages</div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-6">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-orange-100 rounded-lg">
              <AlertTriangle className="h-5 w-5 text-orange-600" />
            </div>
            <div>
              <div className="text-2xl font-bold" data-testid="pages-with-issues-metric">
                {metrics.pagesWithIssues}
              </div>
              <div className="text-sm text-gray-600">Pages with Issues</div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-6">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-red-100 rounded-lg">
              <XCircle className="h-5 w-5 text-red-600" />
            </div>
            <div>
              <div className="text-2xl font-bold" data-testid="total-issues-metric">
                {metrics.totalIssues}
              </div>
              <div className="text-sm text-gray-600">Total Issues</div>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="p-6">
          <div className="flex items-center space-x-3">
            <HealthScoreGauge score={metrics.healthScore} />
            <div>
              <div className="text-2xl font-bold" data-testid="health-score-metric">
                {metrics.healthScore}
              </div>
              <div className="text-sm text-gray-600">Health Score</div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function HealthScoreGauge({ score }: { score: number }) {
  const getColorClass = (score: number) => {
    if (score >= 80) return 'text-green-600';
    if (score >= 60) return 'text-orange-600';
    return 'text-red-600';
  };

  const getFillClass = (score: number) => {
    if (score >= 80) return 'stroke-green-600';
    if (score >= 60) return 'stroke-orange-600';
    return 'stroke-red-600';
  };

  return (
    <div className="relative h-12 w-12" data-testid="health-score-gauge">
      <svg className="h-12 w-12 transform -rotate-90" viewBox="0 0 100 100">
        <circle
          cx="50"
          cy="50"
          r="40"
          stroke="currentColor"
          strokeWidth="8"
          fill="none"
          className="opacity-20"
        />
        <circle
          cx="50"
          cy="50"
          r="40"
          stroke="currentColor"
          strokeWidth="8"
          fill="none"
          strokeDasharray={`${score * 2.51} 251`}
          className={`transition-all duration-500 ${getFillClass(score)}`}
          data-testid="gauge-fill"
        />
      </svg>
      <div className="absolute inset-0 flex items-center justify-center">
        <Gauge className={`h-4 w-4 ${getColorClass(score)}`} />
      </div>
    </div>
  );
}

// ==================== PAGE ANALYSIS SECTION ====================

interface PageAnalysisSectionProps {
  pages: PageAnalysisTableRow[];
  expandedPages: Set<string>;
  selectedPages: Set<string>;
  searchQuery: string;
  sortBy: PageSortBy;
  filterBy: PageFilterBy;
  isLoading: boolean;
  onToggleExpansion: (url: string) => void;
  onToggleSelection: (url: string) => void;
  onSearchChange: (query: string) => void;
  onSortChange: (sort: PageSortBy) => void;
  onFilterChange: (filter: PageFilterBy) => void;
  onRefresh: () => void;
}

function PageAnalysisSection({
  pages,
  expandedPages,
  selectedPages,
  searchQuery,
  sortBy,
  filterBy,
  isLoading,
  onToggleExpansion,
  onToggleSelection,
  onSearchChange,
  onSortChange,
  onFilterChange,
  onRefresh
}: PageAnalysisSectionProps) {
  return (
    <div className="space-y-4" data-testid="page-analysis-section">
      {/* Controls */}
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div className="flex items-center space-x-2">
          <Input
            placeholder="Search pages by URL..."
            value={searchQuery}
            onChange={(e) => onSearchChange(e.target.value)}
            className="w-64"
            data-testid="page-search"
          />
          
          <select 
            className="border rounded px-3 py-2"
            value={filterBy}
            onChange={(e) => onFilterChange(e.target.value as PageFilterBy)}
            data-testid="issue-filter"
          >
            <option value={PageFilterBy.ALL}>All Pages</option>
            <option value={PageFilterBy.ERRORS_ONLY}>Errors Only</option>
            <option value={PageFilterBy.WARNINGS_ONLY}>Warnings Only</option>
            <option value={PageFilterBy.TITLE_ISSUES}>Title Issues</option>
            <option value={PageFilterBy.META_ISSUES}>Meta Issues</option>
            <option value={PageFilterBy.HEADING_ISSUES}>Heading Issues</option>
            <option value={PageFilterBy.IMAGE_ISSUES}>Image Issues</option>
            <option value={PageFilterBy.LINK_ISSUES}>Link Issues</option>
            <option value={PageFilterBy.SLOW_PAGES}>Slow Pages</option>
            <option value={PageFilterBy.LARGE_PAGES}>Large Pages</option>
          </select>

          <select 
            className="border rounded px-3 py-2"
            value={sortBy}
            onChange={(e) => onSortChange(e.target.value as PageSortBy)}
            data-testid="sort-pages"
          >
            <option value={PageSortBy.ISSUES}>By Issues</option>
            <option value={PageSortBy.URL}>By URL</option>
            <option value={PageSortBy.STATUS_CODE}>By Status</option>
            <option value={PageSortBy.LOAD_TIME} data-testid="sort-load-time">By Load Time</option>
            <option value={PageSortBy.SIZE}>By Size</option>
            <option value={PageSortBy.LAST_CRAWLED}>By Last Crawled</option>
          </select>
        </div>

        <div className="flex items-center space-x-2">
          {selectedPages.size > 0 && (
            <div className="flex items-center space-x-2" data-testid="bulk-actions">
              <Badge variant="outline">
                {selectedPages.size} selected
              </Badge>
              <Button size="sm" variant="outline" data-testid="reanalyze-selected">
                <RefreshCw className="h-3 w-3 mr-1" />
                Re-analyze
              </Button>
              <Button size="sm" variant="outline" data-testid="export-selected">
                <Download className="h-3 w-3 mr-1" />
                Export
              </Button>
            </div>
          )}
          
          <Button 
            variant="outline" 
            onClick={onRefresh}
            disabled={isLoading}
            data-testid="refresh-analysis"
          >
            <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
          
          <Button variant="outline" data-testid="export-technical">
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
        </div>
      </div>

      {/* Results Summary */}
      <div className="text-sm text-gray-600">
        Showing {pages.length} pages
        {searchQuery && ` matching "${searchQuery}"`}
        {filterBy !== PageFilterBy.ALL && ` with ${filterBy.replace('-', ' ')}`}
      </div>

      {/* Pages Table */}
      <Card>
        <CardContent className="p-0">
          {pages.length > 0 ? (
            <PageAnalysisTable 
              pages={pages}
              expandedPages={expandedPages}
              selectedPages={selectedPages}
              onToggleExpansion={onToggleExpansion}
              onToggleSelection={onToggleSelection}
            />
          ) : (
            <div className="p-8 text-center" data-testid="no-pages-message">
              <Globe className="h-12 w-12 text-gray-300 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-gray-900 mb-2">No pages found</h3>
              <p className="text-gray-600 mb-4">
                No pages match your current filters. Try adjusting your search or filters.
              </p>
              <Button data-testid="start-analysis-button">
                <Target className="h-4 w-4 mr-2" />
                Start New Analysis
              </Button>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

// ==================== PAGE ANALYSIS TABLE ====================

interface PageAnalysisTableProps {
  pages: PageAnalysisTableRow[];
  expandedPages: Set<string>;
  selectedPages: Set<string>;
  onToggleExpansion: (url: string) => void;
  onToggleSelection: (url: string) => void;
}

function PageAnalysisTable({
  pages,
  expandedPages,
  selectedPages,
  onToggleExpansion,
  onToggleSelection
}: PageAnalysisTableProps) {
  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  };

  const formatLoadTime = (ms: number) => {
    if (ms === 0) return '-';
    return ms < 1000 ? `${ms}ms` : `${(ms / 1000).toFixed(1)}s`;
  };

  const getStatusBadge = (statusCode: number) => {
    let variant: "default" | "secondary" | "destructive" | "outline" = "default";
    let color = "text-green-600 bg-green-50";
    
    if (statusCode >= 400) {
      variant = "destructive";
      color = "text-red-600 bg-red-50";
    } else if (statusCode >= 300) {
      color = "text-yellow-600 bg-yellow-50";
    }

    return (
      <Badge variant={variant} className={color}>
        {statusCode}
      </Badge>
    );
  };

  const getHealthBadge = (health: 'good' | 'warning' | 'error') => {
    const configs = {
      good: { icon: CheckCircle2, color: 'text-green-600 bg-green-50', text: 'Good' },
      warning: { icon: AlertTriangle, color: 'text-yellow-600 bg-yellow-50', text: 'Warning' },
      error: { icon: XCircle, color: 'text-red-600 bg-red-50', text: 'Error' }
    };

    const config = configs[health];
    const Icon = config.icon;

    return (
      <Badge variant="outline" className={config.color}>
        <Icon className="h-3 w-3 mr-1" />
        {config.text}
      </Badge>
    );
  };

  return (
    <div className="overflow-x-auto" data-testid="page-analysis-table">
      <table className="w-full">
        <thead className="bg-gray-50 border-b">
          <tr>
            <th className="w-8 p-3">
              <input
                type="checkbox"
                onChange={(e) => {
                  if (e.target.checked) {
                    pages.forEach(page => onToggleSelection(page.url));
                  } else {
                    pages.forEach(page => {
                      if (selectedPages.has(page.url)) {
                        onToggleSelection(page.url);
                      }
                    });
                  }
                }}
                checked={pages.length > 0 && pages.every(page => selectedPages.has(page.url))}
              />
            </th>
            <th className="text-left p-3 font-semibold" data-testid="url-header">URL</th>
            <th className="text-left p-3 font-semibold" data-testid="status-header">Status</th>
            <th className="text-left p-3 font-semibold" data-testid="load-time-header">Load Time</th>
            <th className="text-left p-3 font-semibold" data-testid="size-header">Size</th>
            <th className="text-left p-3 font-semibold" data-testid="issues-header">Issues</th>
            <th className="text-left p-3 font-semibold" data-testid="health-header">Health</th>
            <th className="w-8 p-3"></th>
          </tr>
        </thead>
        <tbody>
          {pages.map((page, index) => (
            <React.Fragment key={page.url}>
              <tr 
                className="border-b hover:bg-gray-50 cursor-pointer"
                data-testid={`page-row-${index}`}
                onClick={() => onToggleExpansion(page.url)}
              >
                <td className="p-3" onClick={(e) => e.stopPropagation()}>
                  <input
                    type="checkbox"
                    checked={selectedPages.has(page.url)}
                    onChange={() => onToggleSelection(page.url)}
                    data-testid={`page-checkbox-${index}`}
                  />
                </td>
                
                <td className="p-3">
                  <div className="flex items-center space-x-2">
                    <Globe className="h-4 w-4 text-gray-400" />
                    <span 
                      className="font-medium truncate max-w-xs" 
                      title={page.url}
                      data-testid="page-url"
                    >
                      {page.url}
                    </span>
                  </div>
                  <div className="text-xs text-gray-500 mt-1">
                    Depth: {page.depth} • Last crawled: {new Date(page.lastCrawled).toLocaleDateString()}
                  </div>
                </td>
                
                <td className="p-3" data-testid="page-status">
                  {getStatusBadge(page.statusCode)}
                </td>
                
                <td className="p-3" data-testid="page-load-time">
                  {formatLoadTime(page.loadTime)}
                </td>
                
                <td className="p-3" data-testid="page-size">
                  {formatBytes(page.size)}
                </td>
                
                <td className="p-3" data-testid="page-issues">
                  <div className="flex items-center space-x-2">
                    <Badge variant="outline" className="text-red-600">
                      {page.issueCount}
                    </Badge>
                    {page.issueCount > 0 && (
                      <div className="flex space-x-1">
                        {page.issueTypes.slice(0, 3).map((issue, i) => (
                          <Badge key={i} variant="secondary" className="text-xs">
                            {issue.replace('-', ' ')}
                          </Badge>
                        ))}
                        {page.issueTypes.length > 3 && (
                          <Badge variant="secondary" className="text-xs">
                            +{page.issueTypes.length - 3}
                          </Badge>
                        )}
                      </div>
                    )}
                  </div>
                </td>
                
                <td className="p-3" data-testid="page-health">
                  {getHealthBadge(page.overallHealth)}
                </td>
                
                <td className="p-3">
                  {expandedPages.has(page.url) ? (
                    <ChevronDown className="h-4 w-4" />
                  ) : (
                    <ChevronRight className="h-4 w-4" />
                  )}
                </td>
              </tr>
              
              {/* Expanded Details */}
              {expandedPages.has(page.url) && (
                <tr>
                  <td colSpan={8} className="p-0">
                    <PageDetailsPanel page={page} />
                  </td>
                </tr>
              )}
            </React.Fragment>
          ))}
        </tbody>
      </table>
    </div>
  );
}

// ==================== PAGE DETAILS PANEL ====================

function PageDetailsPanel({ page }: { page: PageAnalysisTableRow }) {
  return (
    <div 
      className="bg-gray-50 border-t border-b p-6 space-y-6"
      data-testid={`page-details-${page.url.split('/').pop()}`}
    >
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* SEO Elements */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">SEO Elements</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4" data-testid="seo-elements-section">
            {/* Title */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <h4 className="font-medium">Title</h4>
                <Badge variant="outline" data-testid="title-length">
                  {page.title.length} chars
                </Badge>
              </div>
              <p className="text-sm" data-testid="title-content">
                {page.title.content || <em className="text-gray-500">No title</em>}
              </p>
              {page.title.issues.length > 0 && (
                <div className="mt-2" data-testid="title-issues">
                  {page.title.issues.map((issue, i) => (
                    <Badge key={i} variant="outline" className="text-red-600 mr-1 text-xs">
                      {issue.replace('-', ' ')}
                    </Badge>
                  ))}
                </div>
              )}
            </div>

            {/* Meta Description */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <h4 className="font-medium">Meta Description</h4>
                <Badge variant="outline">
                  {page.metaDescription.length} chars
                </Badge>
              </div>
              <p className="text-sm" data-testid="meta-content">
                {page.metaDescription.content || <em className="text-gray-500">No meta description</em>}
              </p>
              {page.metaDescription.issues.length > 0 && (
                <div className="mt-2" data-testid="meta-issues">
                  {page.metaDescription.issues.map((issue, i) => (
                    <Badge key={i} variant="outline" className="text-red-600 mr-1 text-xs">
                      {issue.replace('-', ' ')}
                    </Badge>
                  ))}
                </div>
              )}
            </div>

            {/* Headings */}
            <div>
              <div className="flex items-center justify-between mb-2">
                <h4 className="font-medium">Headings Structure</h4>
                <Badge 
                  variant="outline" 
                  className={page.headings.structure === 'good' ? 'text-green-600' : 'text-red-600'}
                  data-testid="heading-structure"
                >
                  {page.headings.structure}
                </Badge>
              </div>
              <div className="space-y-2">
                {page.headings.h1.length > 0 && (
                  <div>
                    <span className="text-xs font-semibold text-gray-600">H1:</span>
                    <div className="ml-2" data-testid="h1-list">
                      {page.headings.h1.map((h1, i) => (
                        <div key={i} className="text-sm">{h1}</div>
                      ))}
                    </div>
                  </div>
                )}
                {page.headings.h2.length > 0 && (
                  <div>
                    <span className="text-xs font-semibold text-gray-600">H2:</span>
                    <div className="ml-2" data-testid="h2-list">
                      {page.headings.h2.slice(0, 3).map((h2, i) => (
                        <div key={i} className="text-sm">{h2}</div>
                      ))}
                      {page.headings.h2.length > 3 && (
                        <div className="text-xs text-gray-500">
                          +{page.headings.h2.length - 3} more
                        </div>
                      )}
                    </div>
                  </div>
                )}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Technical Elements */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Technical Elements</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4" data-testid="technical-elements-section">
            <div className="grid grid-cols-1 gap-3 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600">Canonical:</span>
                <span className="text-right max-w-xs truncate" title={page.canonical}>
                  {page.canonical || 'Not set'}
                </span>
              </div>
              
              <div className="flex justify-between">
                <span className="text-gray-600">Robots:</span>
                <span>{page.robots || 'Not set'}</span>
              </div>
              
              {page.lang && (
                <div className="flex justify-between">
                  <span className="text-gray-600">Language:</span>
                  <span>{page.lang}</span>
                </div>
              )}
              
              {page.schema.length > 0 && (
                <div>
                  <span className="text-gray-600">Schema:</span>
                  <div className="mt-1 space-y-1">
                    {page.schema.map((schema, i) => (
                      <Badge 
                        key={i} 
                        variant="outline"
                        className={schema.valid ? 'text-green-600' : 'text-red-600'}
                      >
                        {schema.type}
                      </Badge>
                    ))}
                  </div>
                </div>
              )}
              
              {Object.keys(page.openGraph).length > 0 && (
                <div>
                  <span className="text-gray-600">Open Graph:</span>
                  <div className="mt-1 text-xs text-gray-500">
                    {Object.keys(page.openGraph).length} properties
                  </div>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Images */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base flex items-center space-x-2">
              <Image className="h-4 w-4" />
              <span>Images ({page.images.length})</span>
            </CardTitle>
          </CardHeader>
          <CardContent data-testid="images-section">
            {page.images.length > 0 ? (
              <div className="space-y-3 max-h-48 overflow-y-auto">
                {page.images.slice(0, 5).map((img, i) => (
                  <div key={i} className="border rounded p-2" data-testid={`image-${i}`}>
                    <div className="text-sm font-medium truncate" data-testid="image-src">
                      {img.src}
                    </div>
                    <div className="text-xs text-gray-600 mt-1" data-testid="image-alt">
                      Alt: {img.alt || <em className="text-red-500">Missing</em>}
                    </div>
                    <div className="flex justify-between items-center mt-1">
                      <span className="text-xs text-gray-500" data-testid="image-size">
                        {img.size ? `${(img.size / 1024).toFixed(1)} KB` : 'Unknown size'}
                      </span>
                      {img.issues.length > 0 && (
                        <div className="flex space-x-1" data-testid="image-issues">
                          {img.issues.map((issue, j) => (
                            <Badge key={j} variant="outline" className="text-red-600 text-xs">
                              {issue.replace('-', ' ')}
                            </Badge>
                          ))}
                        </div>
                      )}
                    </div>
                  </div>
                ))}
                {page.images.length > 5 && (
                  <div className="text-center text-xs text-gray-500">
                    +{page.images.length - 5} more images
                  </div>
                )}
              </div>
            ) : (
              <div className="text-sm text-gray-500">No images found</div>
            )}
          </CardContent>
        </Card>

        {/* Links */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base flex items-center space-x-2">
              <Link2 className="h-4 w-4" />
              <span>Links</span>
            </CardTitle>
          </CardHeader>
          <CardContent data-testid="links-section">
            <div className="space-y-3 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-600">Internal:</span>
                <span className="font-medium">{page.links.internal}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">External:</span>
                <span className="font-medium">{page.links.external}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Nofollow:</span>
                <span className="font-medium">{page.links.nofollow}</span>
              </div>
              {page.links.broken.length > 0 && (
                <div>
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600">Broken:</span>
                    <Badge variant="destructive">{page.links.broken.length}</Badge>
                  </div>
                  <div className="mt-2 max-h-24 overflow-y-auto">
                    {page.links.broken.slice(0, 3).map((broken, i) => (
                      <div key={i} className="text-xs text-red-600 truncate" title={broken.url}>
                        {broken.status}: {broken.url}
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Performance */}
        {page.performance && (
          <Card>
            <CardHeader>
              <CardTitle className="text-base flex items-center space-x-2">
                <Zap className="h-4 w-4" />
                <span>Performance</span>
              </CardTitle>
            </CardHeader>
            <CardContent data-testid="performance-section">
              <div className="space-y-3 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-600">FCP:</span>
                  <span className="font-medium">{page.performance.fcp}ms</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">LCP:</span>
                  <span className="font-medium">{page.performance.lcp}ms</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">CLS:</span>
                  <span className="font-medium">{page.performance.cls.toFixed(3)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">TTFB:</span>
                  <span className="font-medium">{page.performance.ttfb}ms</span>
                </div>
              </div>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}

// ==================== GLOBAL ISSUES SECTION ====================

function GlobalIssuesSection({ globalIssues }: { globalIssues: TechnicalAnalysis['globalIssues'] }) {
  return (
    <div className="space-y-6" data-testid="global-issues-section">
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Duplicate Content */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Copy className="h-5 w-5 text-orange-600" />
              <span>Duplicate Content</span>
              <Badge variant="outline">{globalIssues.duplicateContent.length}</Badge>
            </CardTitle>
          </CardHeader>
          <CardContent data-testid="duplicate-content-section">
            {globalIssues.duplicateContent.length > 0 ? (
              <div className="space-y-3">
                {globalIssues.duplicateContent.slice(0, 5).map((dup, i) => (
                  <div key={i} className="border rounded p-3" data-testid={`duplicate-item-${i}`}>
                    <div className="flex justify-between items-center mb-2">
                      <span className="font-medium">
                        {dup.pages.length} pages affected
                      </span>
                      <div className="flex items-center space-x-2">
                        <Badge 
                          variant={dup.level === 'high' ? 'destructive' : dup.level === 'medium' ? 'default' : 'secondary'}
                          data-testid="similarity-level"
                        >
                          {dup.level}
                        </Badge>
                        <span className="text-sm font-medium" data-testid="similarity">
                          {dup.similarity}%
                        </span>
                      </div>
                    </div>
                    <div className="text-sm text-gray-600 space-y-1">
                      {dup.pages.slice(0, 2).map((page, j) => (
                        <div key={j} className="truncate">{page}</div>
                      ))}
                      {dup.pages.length > 2 && (
                        <div className="text-xs">+{dup.pages.length - 2} more pages</div>
                      )}
                    </div>
                    {dup.affectedElements.length > 0 && (
                      <div className="mt-2">
                        <div className="text-xs text-gray-500">Affected elements:</div>
                        <div className="flex space-x-1 mt-1">
                          {dup.affectedElements.map((element, j) => (
                            <Badge key={j} variant="outline" className="text-xs">
                              {element}
                            </Badge>
                          ))}
                        </div>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-6 text-gray-500">
                <CheckCircle2 className="h-8 w-8 mx-auto mb-2 text-green-500" />
                No duplicate content found
              </div>
            )}
          </CardContent>
        </Card>

        {/* Broken Links */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <XCircle className="h-5 w-5 text-red-600" />
              <span>Broken Links</span>
              <Badge variant="outline">{globalIssues.brokenLinks.length}</Badge>
            </CardTitle>
          </CardHeader>
          <CardContent data-testid="broken-links-section">
            {globalIssues.brokenLinks.length > 0 ? (
              <div className="space-y-3 max-h-80 overflow-y-auto">
                {globalIssues.brokenLinks.slice(0, 10).map((link, i) => (
                  <div key={i} className="border rounded p-3" data-testid={`broken-link-${i}`}>
                    <div className="flex justify-between items-start mb-2">
                      <div className="flex-1 min-w-0">
                        <div className="text-sm font-medium truncate" data-testid="broken-url">
                          {link.to}
                        </div>
                        <div className="text-xs text-gray-600 truncate" data-testid="broken-from">
                          From: {link.from}
                        </div>
                      </div>
                      <Badge variant="destructive" className="ml-2" data-testid="broken-status">
                        {link.status}
                      </Badge>
                    </div>
                    <div className="text-xs text-gray-500">
                      Anchor: "{link.anchorText}" • {link.linkType} • {link.position}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-6 text-gray-500">
                <CheckCircle2 className="h-8 w-8 mx-auto mb-2 text-green-500" />
                No broken links found
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Redirect Chains */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <RotateCcw className="h-5 w-5 text-yellow-600" />
              <span>Redirect Chains</span>
              <Badge variant="outline">{globalIssues.redirectChains.length}</Badge>
            </CardTitle>
          </CardHeader>
          <CardContent data-testid="redirect-chains-section">
            {globalIssues.redirectChains.length > 0 ? (
              <div className="space-y-3 max-h-64 overflow-y-auto">
                {globalIssues.redirectChains.slice(0, 5).map((chain, i) => (
                  <div key={i} className="border rounded p-3" data-testid={`redirect-chain-${i}`}>
                    <div className="flex justify-between items-center mb-2">
                      <span className="text-sm font-medium" data-testid="chain-hops">
                        {chain.totalHops} hops
                      </span>
                      <span className="text-xs text-gray-500">
                        {chain.totalTime}ms
                      </span>
                    </div>
                    <div className="space-y-1">
                      {chain.chain.map((url, j) => (
                        <div key={j} className="flex items-center space-x-2" data-testid="chain-step">
                          <div className="text-xs truncate flex-1">{url}</div>
                          {j < chain.chain.length - 1 && (
                            <ArrowRight className="h-3 w-3 text-gray-400 flex-shrink-0" />
                          )}
                        </div>
                      ))}
                    </div>
                    <div className="text-xs text-green-600 mt-2" data-testid="final-url">
                      Final: {chain.finalUrl}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-6 text-gray-500">
                <CheckCircle2 className="h-6 w-6 mx-auto mb-2 text-green-500" />
                <div className="text-xs">No redirect chains</div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Orphan Pages */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Users className="h-5 w-5 text-purple-600" />
              <span>Orphan Pages</span>
              <Badge variant="outline">{globalIssues.orphanPages.length}</Badge>
            </CardTitle>
          </CardHeader>
          <CardContent data-testid="orphan-pages-section">
            {globalIssues.orphanPages.length > 0 ? (
              <div className="space-y-2 max-h-64 overflow-y-auto">
                {globalIssues.orphanPages.slice(0, 10).map((page, i) => (
                  <div key={i} className="text-sm truncate p-2 bg-gray-50 rounded" data-testid={`orphan-page-${i}`}>
                    {page}
                  </div>
                ))}
                {globalIssues.orphanPages.length > 10 && (
                  <div className="text-xs text-gray-500 text-center">
                    +{globalIssues.orphanPages.length - 10} more
                  </div>
                )}
              </div>
            ) : (
              <div className="text-center py-6 text-gray-500">
                <CheckCircle2 className="h-6 w-6 mx-auto mb-2 text-green-500" />
                <div className="text-xs">No orphan pages</div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Missing Elements Summary */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <AlertTriangle className="h-5 w-5 text-red-600" />
              <span>Missing Elements</span>
            </CardTitle>
          </CardHeader>
          <CardContent data-testid="missing-elements-issues">
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-600">Missing Titles:</span>
                <Badge variant={globalIssues.missingTitles.length > 0 ? "destructive" : "outline"}>
                  {globalIssues.missingTitles.length}
                </Badge>
              </div>
              
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-600">Missing Meta Descriptions:</span>
                <Badge variant={globalIssues.missingMeta.length > 0 ? "destructive" : "outline"}>
                  {globalIssues.missingMeta.length}
                </Badge>
              </div>
              
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-600">Duplicate Titles:</span>
                <Badge variant={globalIssues.duplicateTitles.length > 0 ? "default" : "outline"}>
                  {globalIssues.duplicateTitles.length}
                </Badge>
              </div>
              
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-600">Large Pages (&gt;1MB):</span>
                <Badge variant={globalIssues.largePages.length > 0 ? "default" : "outline"}>
                  {globalIssues.largePages.length}
                </Badge>
              </div>
              
              <div className="flex justify-between items-center">
                <span className="text-sm text-gray-600">Slow Pages (&gt;3s):</span>
                <Badge variant={globalIssues.slowPages.length > 0 ? "default" : "outline"}>
                  {globalIssues.slowPages.length}
                </Badge>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

// ==================== CRAWLABILITY SECTION ====================

function CrawlabilitySection({ crawlability }: { crawlability: TechnicalAnalysis['crawlability'] }) {
  return (
    <div className="space-y-6" data-testid="crawlability-section">
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Robots.txt Analysis */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <Shield className="h-5 w-5 text-blue-600" />
              <span>Robots.txt Analysis</span>
            </CardTitle>
          </CardHeader>
          <CardContent data-testid="robots-txt-section">
            <div className="space-y-4">
              {/* Status Indicators */}
              <div className="grid grid-cols-2 gap-4">
                <div className="flex items-center space-x-2">
                  {crawlability.robotsTxt.exists ? (
                    <CheckCircle2 className="h-4 w-4 text-green-600" />
                  ) : (
                    <XCircle className="h-4 w-4 text-red-600" />
                  )}
                  <span className="text-sm" data-testid="robots-exists">
                    {crawlability.robotsTxt.exists ? '✓ Exists' : '✗ Missing'}
                  </span>
                </div>
                
                <div className="flex items-center space-x-2">
                  {crawlability.robotsTxt.valid ? (
                    <CheckCircle2 className="h-4 w-4 text-green-600" />
                  ) : (
                    <XCircle className="h-4 w-4 text-red-600" />
                  )}
                  <span className="text-sm" data-testid="robots-valid">
                    {crawlability.robotsTxt.valid ? '✓ Valid' : '✗ Invalid'}
                  </span>
                </div>
              </div>

              {/* User Agents */}
              {crawlability.robotsTxt.userAgents.length > 0 && (
                <div>
                  <h4 className="font-medium text-sm mb-2">User Agents:</h4>
                  <div className="flex flex-wrap gap-1">
                    {crawlability.robotsTxt.userAgents.map((agent, i) => (
                      <Badge key={i} variant="outline" className="text-xs">
                        {agent}
                      </Badge>
                    ))}
                  </div>
                </div>
              )}

              {/* Disallowed Paths */}
              {crawlability.robotsTxt.disallowedPaths.length > 0 && (
                <div>
                  <h4 className="font-medium text-sm mb-2">Disallowed Paths:</h4>
                  <div className="space-y-1 max-h-32 overflow-y-auto" data-testid="disallowed-paths">
                    {crawlability.robotsTxt.disallowedPaths.map((path, i) => (
                      <div key={i} className="text-xs font-mono bg-gray-100 px-2 py-1 rounded">
                        {path}
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Sitemap URLs */}
              {crawlability.robotsTxt.sitemapUrls.length > 0 && (
                <div>
                  <h4 className="font-medium text-sm mb-2">Sitemap URLs:</h4>
                  <div className="space-y-1" data-testid="sitemap-urls">
                    {crawlability.robotsTxt.sitemapUrls.map((url, i) => (
                      <div key={i} className="text-xs truncate">{url}</div>
                    ))}
                  </div>
                </div>
              )}

              {/* Issues */}
              {crawlability.robotsTxt.issues.length > 0 && (
                <div>
                  <h4 className="font-medium text-sm mb-2">Issues:</h4>
                  <div className="space-y-1">
                    {crawlability.robotsTxt.issues.map((issue, i) => (
                      <div key={i} className="text-xs text-red-600 flex items-center space-x-1">
                        <AlertTriangle className="h-3 w-3" />
                        <span>{issue}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Sitemap Analysis */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <FileText className="h-5 w-5 text-green-600" />
              <span>Sitemap Analysis</span>
            </CardTitle>
          </CardHeader>
          <CardContent data-testid="sitemap-section">
            <div className="space-y-4">
              {/* Status Indicators */}
              <div className="grid grid-cols-2 gap-4">
                <div className="flex items-center space-x-2">
                  {crawlability.sitemap.exists ? (
                    <CheckCircle2 className="h-4 w-4 text-green-600" />
                  ) : (
                    <XCircle className="h-4 w-4 text-red-600" />
                  )}
                  <span className="text-sm" data-testid="sitemap-exists">
                    {crawlability.sitemap.exists ? '✓ Exists' : '✗ Missing'}
                  </span>
                </div>
                
                <div className="flex items-center space-x-2">
                  <Badge variant="outline" data-testid="sitemap-format">
                    {crawlability.sitemap.format.toUpperCase()}
                  </Badge>
                </div>
              </div>

              {/* Metrics */}
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Pages in Sitemap:</span>
                  <span className="font-semibold" data-testid="pages-in-sitemap">
                    {crawlability.sitemap.pagesInSitemap}
                  </span>
                </div>
                
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Images:</span>
                  <span className="font-semibold" data-testid="images-in-sitemap">
                    {crawlability.sitemap.images}
                  </span>
                </div>
                
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">Videos:</span>
                  <span className="font-semibold">
                    {crawlability.sitemap.videos}
                  </span>
                </div>
                
                {crawlability.sitemap.lastModified && (
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-gray-600">Last Modified:</span>
                    <span className="text-xs">
                      {new Date(crawlability.sitemap.lastModified).toLocaleDateString()}
                    </span>
                  </div>
                )}
              </div>

              {/* URL */}
              {crawlability.sitemap.url && (
                <div>
                  <h4 className="font-medium text-sm mb-2">Sitemap URL:</h4>
                  <div className="text-xs font-mono bg-gray-100 px-2 py-1 rounded truncate">
                    {crawlability.sitemap.url}
                  </div>
                </div>
              )}

              {/* Issues */}
              {crawlability.sitemap.issues.length > 0 && (
                <div>
                  <h4 className="font-medium text-sm mb-2">Issues:</h4>
                  <div className="space-y-1">
                    {crawlability.sitemap.issues.map((issue, i) => (
                      <div key={i} className="text-xs text-red-600 flex items-center space-x-1">
                        <AlertTriangle className="h-3 w-3" />
                        <span>{issue}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Crawl Budget Analysis */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Activity className="h-5 w-5 text-purple-600" />
            <span>Crawl Budget Analysis</span>
          </CardTitle>
        </CardHeader>
        <CardContent data-testid="crawl-budget-section">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            {/* Metrics Cards */}
            <div className="space-y-4">
              <div className="text-center">
                <div className="text-2xl font-bold text-blue-600" data-testid="total-pages">
                  {crawlability.crawlBudget.totalPages}
                </div>
                <div className="text-sm text-gray-600">Total Pages</div>
              </div>
              
              <div className="text-center">
                <div className="text-2xl font-bold text-green-600" data-testid="crawlable-pages">
                  {crawlability.crawlBudget.crawlablePages}
                </div>
                <div className="text-sm text-gray-600">Crawlable</div>
              </div>
            </div>

            <div className="space-y-4">
              <div className="text-center">
                <div className="text-2xl font-bold text-red-600" data-testid="blocked-pages">
                  {crawlability.crawlBudget.blockedPages}
                </div>
                <div className="text-sm text-gray-600">Blocked</div>
              </div>
              
              <div className="text-center">
                <div className="text-2xl font-bold text-purple-600" data-testid="crawl-efficiency">
                  {crawlability.crawlBudget.crawlEfficiency.toFixed(1)}%
                </div>
                <div className="text-sm text-gray-600">Efficiency</div>
              </div>
            </div>

            {/* Crawl Efficiency Progress */}
            <div className="col-span-2">
              <div className="space-y-4">
                <div>
                  <div className="flex justify-between text-sm mb-2">
                    <span>Crawl Efficiency</span>
                    <span>{crawlability.crawlBudget.crawlEfficiency.toFixed(1)}%</span>
                  </div>
                  <Progress 
                    value={crawlability.crawlBudget.crawlEfficiency} 
                    className="h-2"
                  />
                </div>
                
                <div>
                  <div className="text-sm text-gray-600 mb-2">Average Crawl Time:</div>
                  <div className="text-lg font-semibold">
                    {crawlability.crawlBudget.averageCrawlTime.toFixed(1)}s per page
                  </div>
                </div>
                
                {/* Depth Distribution */}
                {Object.keys(crawlability.crawlBudget.pagesPerLevel).length > 0 && (
                  <div data-testid="depth-distribution">
                    <div className="text-sm text-gray-600 mb-2">Pages by Depth:</div>
                    <div className="space-y-1">
                      {Object.entries(crawlability.crawlBudget.pagesPerLevel)
                        .sort(([a], [b]) => Number(a) - Number(b))
                        .map(([depth, count]) => (
                          <div key={depth} className="flex justify-between text-xs">
                            <span>Level {depth}:</span>
                            <span className="font-medium">{count} pages</span>
                          </div>
                        ))}
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

// ==================== ANALYSIS OVERVIEW SECTION ====================

function AnalysisOverviewSection({ 
  analysis, 
  totalPages 
}: { 
  analysis: TechnicalAnalysis; 
  totalPages: number;
}) {
  return (
    <div className="space-y-6">
      {/* Analysis Status */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <BarChart3 className="h-5 w-5" />
            <span>Analysis Overview</span>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div>
              <h4 className="font-medium mb-3">Analysis Status</h4>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-600">Status:</span>
                  <Badge 
                    variant={analysis.status.crawlStatus === 'completed' ? 'default' : 'outline'}
                    data-testid="crawl-status"
                  >
                    {analysis.status.crawlStatus}
                  </Badge>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Duration:</span>
                  <span>{analysis.status.crawlDuration} min</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Date:</span>
                  <span>{new Date(analysis.status.analysisDate).toLocaleDateString()}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Version:</span>
                  <span>{analysis.status.version}</span>
                </div>
              </div>
            </div>

            <div>
              <h4 className="font-medium mb-3">Crawl Configuration</h4>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-600">Max Depth:</span>
                  <span>{analysis.config.maxDepth}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Respect Robots:</span>
                  <span>{analysis.config.respectRobots ? '✓' : '✗'}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Crawl Delay:</span>
                  <span>{analysis.config.crawlDelay}ms</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">User Agent:</span>
                  <span className="truncate">{analysis.config.userAgent}</span>
                </div>
              </div>
            </div>

            <div>
              <h4 className="font-medium mb-3">Performance Summary</h4>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-600">Avg Load Time:</span>
                  <span>{(analysis.metrics.avgLoadTime / 1000).toFixed(1)}s</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Avg Page Size:</span>
                  <span>{(analysis.metrics.avgPageSize / 1024).toFixed(1)} KB</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Pages Analyzed:</span>
                  <span>{totalPages}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-600">Health Score:</span>
                  <Badge 
                    variant={analysis.metrics.healthScore >= 80 ? 'default' : 
                           analysis.metrics.healthScore >= 60 ? 'secondary' : 'destructive'}
                  >
                    {analysis.metrics.healthScore}/100
                  </Badge>
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Issue Types Breakdown */}
      <Card>
        <CardHeader>
          <CardTitle>Issues Breakdown</CardTitle>
        </CardHeader>
        <CardContent>
          {Object.keys(analysis.metrics.issuesByType).length > 0 ? (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {Object.entries(analysis.metrics.issuesByType)
                .sort(([,a], [,b]) => b - a)
                .slice(0, 8)
                .map(([issueType, count]) => (
                  <div key={issueType} className="text-center p-3 bg-gray-50 rounded">
                    <div className="text-xl font-bold text-red-600">{count}</div>
                    <div className="text-xs text-gray-600 capitalize">
                      {issueType.replace('-', ' ')}
                    </div>
                  </div>
                ))}
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500">
              <CheckCircle2 className="h-12 w-12 mx-auto mb-4 text-green-500" />
              <h3 className="text-lg font-semibold mb-2">No Issues Found</h3>
              <p>All pages passed the technical analysis successfully.</p>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}