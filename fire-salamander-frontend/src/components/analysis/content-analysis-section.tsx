/**
 * Fire Salamander - Content Analysis Section Component
 * Professional Content Analysis Display with AI Insights & Advanced Visualizations
 * Lead Tech implementation - Niveau SEMrush/Ahrefs
 */

'use client';

import { useState, useEffect, useMemo } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { Alert, AlertDescription } from '@/components/ui/alert';
import {
  ContentAnalysis,
  ContentPage,
  ContentQuality,
  ContentType,
  ContentIssueType,
  ContentSeverity,
  ReadabilityLevel,
  TopicImportance,
  ContentSortBy,
  ContentFilterBy,
  CONTENT_QUALITY_THRESHOLDS,
} from '@/types/content-analysis';
import {
  BarChart3,
  BookOpen,
  Brain,
  CheckCircle,
  Clock,
  FileText,
  Filter,
  Hash,
  Heart,
  HelpCircle,
  Image,
  Lightbulb,
  LineChart,
  Link,
  MessageSquare,
  Search,
  Sparkles,
  Target,
  TrendingUp,
  Users,
  Zap,
  ArrowUpDown,
  ChevronDown,
  ChevronRight,
  Download,
  ExternalLink,
  Eye,
  Settings,
  Star,
  Wand2,
  XCircle,
  AlertTriangle,
  Info,
} from 'lucide-react';

interface ContentAnalysisSectionProps {
  contentData: ContentAnalysis;
  analysisId: string;
}

interface SortConfig {
  key: ContentSortBy;
  direction: 'asc' | 'desc';
}

export function ContentAnalysisSection({ 
  contentData, 
  analysisId 
}: ContentAnalysisSectionProps) {
  // États pour la gestion de l'interface
  const [activeTab, setActiveTab] = useState<'overview' | 'pages' | 'ai-insights' | 'visualizations'>('overview');
  const [expandedPages, setExpandedPages] = useState<Set<string>>(new Set());
  const [selectedPages, setSelectedPages] = useState<Set<string>>(new Set());
  const [sortConfig, setSortConfig] = useState<SortConfig>({
    key: ContentSortBy.QUALITY_SCORE,
    direction: 'desc'
  });
  const [activeFilter, setActiveFilter] = useState<ContentFilterBy>(ContentFilterBy.ALL);
  const [searchQuery, setSearchQuery] = useState('');
  const [showAIRecommendations, setShowAIRecommendations] = useState(false);

  // Mémorisation des données triées et filtrées
  const filteredAndSortedPages = useMemo(() => {
    let filtered = contentData.pages;

    // Application des filtres
    switch (activeFilter) {
      case ContentFilterBy.HIGH_QUALITY:
        filtered = filtered.filter(p => p.quality.overallScore >= CONTENT_QUALITY_THRESHOLDS.good.min);
        break;
      case ContentFilterBy.LOW_QUALITY:
        filtered = filtered.filter(p => p.quality.overallScore < CONTENT_QUALITY_THRESHOLDS.average.min);
        break;
      case ContentFilterBy.THIN_CONTENT:
        filtered = filtered.filter(p => p.issues.some(i => i.type === ContentIssueType.THIN_CONTENT));
        break;
      case ContentFilterBy.DUPLICATE_CONTENT:
        filtered = filtered.filter(p => p.issues.some(i => i.type === ContentIssueType.DUPLICATE_CONTENT));
        break;
      case ContentFilterBy.MISSING_KEYWORDS:
        filtered = filtered.filter(p => p.issues.some(i => i.type === ContentIssueType.MISSING_KEYWORDS));
        break;
    }

    // Application de la recherche
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter(page => 
        page.title.toLowerCase().includes(query) ||
        page.url.toLowerCase().includes(query) ||
        page.contentType.toLowerCase().includes(query)
      );
    }

    // Application du tri
    const sorted = [...filtered].sort((a, b) => {
      let aValue: any, bValue: any;

      switch (sortConfig.key) {
        case ContentSortBy.QUALITY_SCORE:
          aValue = a.quality.overallScore;
          bValue = b.quality.overallScore;
          break;
        case ContentSortBy.WORD_COUNT:
          aValue = a.metrics.wordCount;
          bValue = b.metrics.wordCount;
          break;
        case ContentSortBy.READABILITY:
          aValue = a.readability.fleschReadingEase;
          bValue = b.readability.fleschReadingEase;
          break;
        case ContentSortBy.KEYWORD_DENSITY:
          aValue = a.keywords.primary[0]?.density || 0;
          bValue = b.keywords.primary[0]?.density || 0;
          break;
        case ContentSortBy.LAST_MODIFIED:
          aValue = new Date(a.lastModified || 0).getTime();
          bValue = new Date(b.lastModified || 0).getTime();
          break;
        case ContentSortBy.ISSUES_COUNT:
          aValue = a.issues.length;
          bValue = b.issues.length;
          break;
        case ContentSortBy.URL:
          aValue = a.url;
          bValue = b.url;
          break;
        default:
          aValue = a.quality.overallScore;
          bValue = b.quality.overallScore;
      }

      if (typeof aValue === 'string') {
        return sortConfig.direction === 'asc' 
          ? aValue.localeCompare(bValue)
          : bValue.localeCompare(aValue);
      }

      return sortConfig.direction === 'asc' ? aValue - bValue : bValue - aValue;
    });

    return sorted;
  }, [contentData.pages, activeFilter, searchQuery, sortConfig]);

  // Fonctions utilitaires
  const togglePageExpansion = (url: string) => {
    setExpandedPages(prev => {
      const newSet = new Set(prev);
      if (newSet.has(url)) {
        newSet.delete(url);
      } else {
        newSet.add(url);
      }
      return newSet;
    });
  };

  const togglePageSelection = (url: string) => {
    setSelectedPages(prev => {
      const newSet = new Set(prev);
      if (newSet.has(url)) {
        newSet.delete(url);
      } else {
        newSet.add(url);
      }
      return newSet;
    });
  };

  const handleSort = (key: ContentSortBy) => {
    setSortConfig(prev => ({
      key,
      direction: prev.key === key && prev.direction === 'desc' ? 'asc' : 'desc'
    }));
  };

  const getQualityBadgeVariant = (score: number) => {
    if (score >= CONTENT_QUALITY_THRESHOLDS.excellent.min) return 'default';
    if (score >= CONTENT_QUALITY_THRESHOLDS.good.min) return 'secondary';
    if (score >= CONTENT_QUALITY_THRESHOLDS.average.min) return 'outline';
    return 'destructive';
  };

  const getQualityLabel = (score: number) => {
    if (score >= CONTENT_QUALITY_THRESHOLDS.excellent.min) return 'Excellent';
    if (score >= CONTENT_QUALITY_THRESHOLDS.good.min) return 'Bon';
    if (score >= CONTENT_QUALITY_THRESHOLDS.average.min) return 'Moyen';
    if (score >= CONTENT_QUALITY_THRESHOLDS.poor.min) return 'Faible';
    return 'Critique';
  };

  const getReadabilityColor = (level: ReadabilityLevel) => {
    const colors = {
      [ReadabilityLevel.VERY_EASY]: 'text-green-600',
      [ReadabilityLevel.EASY]: 'text-green-500',
      [ReadabilityLevel.FAIRLY_EASY]: 'text-blue-500',
      [ReadabilityLevel.STANDARD]: 'text-yellow-600',
      [ReadabilityLevel.FAIRLY_DIFFICULT]: 'text-orange-500',
      [ReadabilityLevel.DIFFICULT]: 'text-red-500',
      [ReadabilityLevel.VERY_DIFFICULT]: 'text-red-600',
    };
    return colors[level] || 'text-gray-500';
  };

  const getSeverityIcon = (severity: ContentSeverity) => {
    switch (severity) {
      case ContentSeverity.CRITICAL: return <XCircle className="h-4 w-4 text-red-500" />;
      case ContentSeverity.WARNING: return <AlertTriangle className="h-4 w-4 text-yellow-500" />;
      case ContentSeverity.INFO: return <Info className="h-4 w-4 text-blue-500" />;
    }
  };

  const getContentTypeIcon = (type: ContentType) => {
    const icons = {
      [ContentType.ARTICLE]: <FileText className="h-4 w-4" />,
      [ContentType.BLOG_POST]: <BookOpen className="h-4 w-4" />,
      [ContentType.PRODUCT]: <Star className="h-4 w-4" />,
      [ContentType.CATEGORY]: <Filter className="h-4 w-4" />,
      [ContentType.HOMEPAGE]: <Heart className="h-4 w-4" />,
      [ContentType.LANDING]: <Target className="h-4 w-4" />,
    };
    return icons[type] || <FileText className="h-4 w-4" />;
  };

  // Rendu du composant principal
  return (
    <div className="space-y-6">
      {/* En-tête avec métriques globales */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Pages analysées</p>
                <p className="text-2xl font-bold">{contentData.globalMetrics.totalPages}</p>
              </div>
              <FileText className="h-8 w-8 text-blue-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Score moyen</p>
                <p className="text-2xl font-bold">{contentData.globalMetrics.avgQualityScore}</p>
              </div>
              <TrendingUp className="h-8 w-8 text-green-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Mots total</p>
                <p className="text-2xl font-bold">{contentData.globalMetrics.totalWords.toLocaleString()}</p>
              </div>
              <Hash className="h-8 w-8 text-purple-500" />
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-gray-600">Contenu mince</p>
                <p className="text-2xl font-bold text-red-500">{contentData.globalMetrics.thinContentPages}</p>
              </div>
              <AlertTriangle className="h-8 w-8 text-red-500" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Navigation par onglets */}
      <div className="border-b">
        <nav className="-mb-px flex space-x-8">
          {[
            { id: 'overview', label: 'Vue d\'ensemble', icon: BarChart3 },
            { id: 'pages', label: 'Analyse des pages', icon: FileText },
            { id: 'ai-insights', label: 'Insights IA', icon: Brain },
            { id: 'visualizations', label: 'Visualisations', icon: LineChart },
          ].map(tab => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id as any)}
                className={`flex items-center py-2 px-1 border-b-2 font-medium text-sm whitespace-nowrap ${
                  activeTab === tab.id
                    ? 'border-blue-500 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }`}
              >
                <Icon className="h-4 w-4 mr-2" />
                {tab.label}
              </button>
            );
          })}
        </nav>
      </div>

      {/* Contenu des onglets */}
      {activeTab === 'overview' && (
        <div className="space-y-6">
          {/* Métriques détaillées */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-lg">Distribution de qualité</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {contentData.visualizationData.charts.qualityDistribution.map((item, index) => (
                    <div key={index} className="flex items-center justify-between">
                      <span className="text-sm font-medium">{item.range}</span>
                      <div className="flex items-center gap-2">
                        <Progress value={item.percentage} className="w-20" />
                        <span className="text-sm text-gray-600">{item.count}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-lg">Types de contenu</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {contentData.visualizationData.charts.contentTypeDistribution.map((item, index) => (
                    <div key={index} className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        {getContentTypeIcon(item.type)}
                        <span className="text-sm font-medium">{item.type}</span>
                      </div>
                      <div className="text-right">
                        <p className="text-sm font-bold">{item.count}</p>
                        <p className="text-xs text-gray-600">Q: {item.avgQuality}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-lg">Benchmarking concurrent</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <span className="text-sm">Position</span>
                    <Badge variant={
                      contentData.globalMetrics.competitorBenchmark.position === 'above' ? 'default' :
                      contentData.globalMetrics.competitorBenchmark.position === 'equal' ? 'secondary' : 'destructive'
                    }>
                      {contentData.globalMetrics.competitorBenchmark.position === 'above' ? 'Au-dessus' :
                       contentData.globalMetrics.competitorBenchmark.position === 'equal' ? 'Égal' : 'En-dessous'}
                    </Badge>
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Volume de contenu</span>
                      <span className={contentData.globalMetrics.competitorBenchmark.gapAnalysis.contentVolume < 0 ? 'text-red-600' : 'text-green-600'}>
                        {contentData.globalMetrics.competitorBenchmark.gapAnalysis.contentVolume > 0 ? '+' : ''}{contentData.globalMetrics.competitorBenchmark.gapAnalysis.contentVolume}%
                      </span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Qualité</span>
                      <span className={contentData.globalMetrics.competitorBenchmark.gapAnalysis.contentQuality < 0 ? 'text-red-600' : 'text-green-600'}>
                        {contentData.globalMetrics.competitorBenchmark.gapAnalysis.contentQuality > 0 ? '+' : ''}{contentData.globalMetrics.competitorBenchmark.gapAnalysis.contentQuality} pts
                      </span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Couverture sujets</span>
                      <span className={contentData.globalMetrics.competitorBenchmark.gapAnalysis.topicCoverage < 0 ? 'text-red-600' : 'text-green-600'}>
                        {contentData.globalMetrics.competitorBenchmark.gapAnalysis.topicCoverage > 0 ? '+' : ''}{contentData.globalMetrics.competitorBenchmark.gapAnalysis.topicCoverage}%
                      </span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Issues globales */}
          {contentData.globalIssues.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <AlertTriangle className="h-5 w-5 text-yellow-500" />
                  Issues globales ({contentData.globalIssues.length})
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {contentData.globalIssues.slice(0, 5).map((issue, index) => (
                    <div key={index} className="border rounded-lg p-4">
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-2">
                            {getSeverityIcon(issue.severity)}
                            <h4 className="font-medium">{issue.type.replace('-', ' ')}</h4>
                            <Badge variant="outline">
                              Priorité {issue.priority}
                            </Badge>
                          </div>
                          <p className="text-sm text-gray-600 mb-2">{issue.description}</p>
                          <p className="text-sm font-medium">
                            Pages affectées: {issue.affectedPages.length}
                          </p>
                          <div className="mt-2 flex items-center gap-4 text-xs text-gray-500">
                            <span>ROI: +{issue.estimatedROI.trafficGain}% trafic</span>
                            <span>Amélioration: +{issue.estimatedROI.rankingImprovement} positions</span>
                            <span>Coût: {issue.estimatedROI.implementationCost}</span>
                          </div>
                        </div>
                        <Button
                          variant="outline"
                          size="sm"
                          className="ml-4"
                          onClick={() => console.log('Show fix details for:', issue.id)}
                        >
                          <Wand2 className="h-4 w-4 mr-1" />
                          Corriger
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      )}

      {activeTab === 'pages' && (
        <div className="space-y-4">
          {/* Contrôles de filtrage et tri */}
          <div className="flex flex-col sm:flex-row gap-4 items-center justify-between">
            <div className="flex gap-2 flex-wrap">
              <select
                value={activeFilter}
                onChange={(e) => setActiveFilter(e.target.value as ContentFilterBy)}
                className="px-3 py-2 border border-gray-300 rounded-md text-sm"
              >
                <option value={ContentFilterBy.ALL}>Toutes les pages</option>
                <option value={ContentFilterBy.HIGH_QUALITY}>Haute qualité</option>
                <option value={ContentFilterBy.LOW_QUALITY}>Faible qualité</option>
                <option value={ContentFilterBy.THIN_CONTENT}>Contenu mince</option>
                <option value={ContentFilterBy.DUPLICATE_CONTENT}>Contenu dupliqué</option>
                <option value={ContentFilterBy.MISSING_KEYWORDS}>Mots-clés manquants</option>
              </select>
              
              <input
                type="text"
                placeholder="Rechercher une page..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="px-3 py-2 border border-gray-300 rounded-md text-sm min-w-48"
              />
            </div>

            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setSelectedPages(new Set())}
                disabled={selectedPages.size === 0}
              >
                Désélectionner tout
              </Button>
              <span className="text-sm text-gray-600">
                {filteredAndSortedPages.length} pages
              </span>
            </div>
          </div>

          {/* En-têtes de colonnes pour le tri */}
          <Card>
            <CardContent className="p-4">
              <div className="grid grid-cols-12 gap-4 text-sm font-medium text-gray-600 border-b pb-2">
                <div className="col-span-1">
                  <input
                    type="checkbox"
                    onChange={(e) => {
                      if (e.target.checked) {
                        setSelectedPages(new Set(filteredAndSortedPages.map(p => p.url)));
                      } else {
                        setSelectedPages(new Set());
                      }
                    }}
                    checked={selectedPages.size === filteredAndSortedPages.length && filteredAndSortedPages.length > 0}
                  />
                </div>
                <div className="col-span-4">
                  <button 
                    onClick={() => handleSort(ContentSortBy.URL)}
                    className="flex items-center gap-1 hover:text-gray-900"
                  >
                    Page
                    <ArrowUpDown className="h-3 w-3" />
                  </button>
                </div>
                <div className="col-span-1">
                  <button 
                    onClick={() => handleSort(ContentSortBy.QUALITY_SCORE)}
                    className="flex items-center gap-1 hover:text-gray-900"
                  >
                    Qualité
                    <ArrowUpDown className="h-3 w-3" />
                  </button>
                </div>
                <div className="col-span-1">
                  <button 
                    onClick={() => handleSort(ContentSortBy.WORD_COUNT)}
                    className="flex items-center gap-1 hover:text-gray-900"
                  >
                    Mots
                    <ArrowUpDown className="h-3 w-3" />
                  </button>
                </div>
                <div className="col-span-1">
                  <button 
                    onClick={() => handleSort(ContentSortBy.READABILITY)}
                    className="flex items-center gap-1 hover:text-gray-900"
                  >
                    Lisibilité
                    <ArrowUpDown className="h-3 w-3" />
                  </button>
                </div>
                <div className="col-span-1">
                  <button 
                    onClick={() => handleSort(ContentSortBy.ISSUES_COUNT)}
                    className="flex items-center gap-1 hover:text-gray-900"
                  >
                    Issues
                    <ArrowUpDown className="h-3 w-3" />
                  </button>
                </div>
                <div className="col-span-2">
                  Mots-clés
                </div>
                <div className="col-span-1">
                  Actions
                </div>
              </div>

              {/* Liste des pages */}
              <div className="space-y-2 mt-4">
                {filteredAndSortedPages.map((page, index) => (
                  <div key={page.url} className="border rounded-lg">
                    {/* Ligne principale */}
                    <div className="grid grid-cols-12 gap-4 p-4 items-center">
                      <div className="col-span-1">
                        <input
                          type="checkbox"
                          checked={selectedPages.has(page.url)}
                          onChange={() => togglePageSelection(page.url)}
                        />
                      </div>
                      
                      <div className="col-span-4">
                        <div className="flex items-center gap-2">
                          {getContentTypeIcon(page.contentType)}
                          <div>
                            <p className="font-medium text-sm truncate max-w-64">{page.title}</p>
                            <p className="text-xs text-gray-500 truncate max-w-64">{page.url}</p>
                          </div>
                        </div>
                      </div>
                      
                      <div className="col-span-1">
                        <Badge variant={getQualityBadgeVariant(page.quality.overallScore)}>
                          {page.quality.overallScore}
                        </Badge>
                      </div>
                      
                      <div className="col-span-1">
                        <span className="text-sm">{page.metrics.wordCount.toLocaleString()}</span>
                      </div>
                      
                      <div className="col-span-1">
                        <span className={`text-sm ${getReadabilityColor(page.readability.level)}`}>
                          {Math.round(page.readability.fleschReadingEase)}
                        </span>
                      </div>
                      
                      <div className="col-span-1">
                        {page.issues.length > 0 ? (
                          <Badge variant="destructive" className="text-xs">
                            {page.issues.length}
                          </Badge>
                        ) : (
                          <CheckCircle className="h-4 w-4 text-green-500" />
                        )}
                      </div>
                      
                      <div className="col-span-2">
                        <div className="flex flex-wrap gap-1">
                          {page.keywords.primary.slice(0, 2).map((kw, kwIndex) => (
                            <Badge key={kwIndex} variant="outline" className="text-xs">
                              {kw.keyword} ({kw.density.toFixed(1)}%)
                            </Badge>
                          ))}
                          {page.keywords.primary.length > 2 && (
                            <span className="text-xs text-gray-500">+{page.keywords.primary.length - 2}</span>
                          )}
                        </div>
                      </div>
                      
                      <div className="col-span-1">
                        <div className="flex gap-1">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => togglePageExpansion(page.url)}
                          >
                            {expandedPages.has(page.url) ? 
                              <ChevronDown className="h-4 w-4" /> : 
                              <ChevronRight className="h-4 w-4" />
                            }
                          </Button>
                          <Button variant="ghost" size="sm">
                            <ExternalLink className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                    </div>

                    {/* Détails expandus */}
                    {expandedPages.has(page.url) && (
                      <div className="border-t bg-gray-50 p-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                          
                          {/* Métriques détaillées */}
                          <div>
                            <h4 className="font-medium text-sm mb-2 flex items-center gap-2">
                              <BarChart3 className="h-4 w-4" />
                              Métriques
                            </h4>
                            <div className="space-y-1 text-xs">
                              <div className="flex justify-between">
                                <span>Mots uniques:</span>
                                <span>{page.metrics.uniqueWords.toLocaleString()}</span>
                              </div>
                              <div className="flex justify-between">
                                <span>Phrases:</span>
                                <span>{page.metrics.sentenceCount}</span>
                              </div>
                              <div className="flex justify-between">
                                <span>Paragraphes:</span>
                                <span>{page.metrics.paragraphCount}</span>
                              </div>
                              <div className="flex justify-between">
                                <span>Temps lecture:</span>
                                <span>{page.metrics.readingTime} min</span>
                              </div>
                            </div>
                          </div>

                          {/* Structure */}
                          <div>
                            <h4 className="font-medium text-sm mb-2 flex items-center gap-2">
                              <FileText className="h-4 w-4" />
                              Structure
                            </h4>
                            <div className="space-y-1 text-xs">
                              <div className="flex justify-between">
                                <span>Titres:</span>
                                <span>H1:{page.structure.headings.h1} H2:{page.structure.headings.h2} H3:{page.structure.headings.h3}</span>
                              </div>
                              <div className="flex justify-between">
                                <span>Images:</span>
                                <span>{page.structure.images.total} ({page.structure.images.withAlt} avec alt)</span>
                              </div>
                              <div className="flex justify-between">
                                <span>Liens:</span>
                                <span>Int:{page.structure.links.internal} Ext:{page.structure.links.external}</span>
                              </div>
                            </div>
                          </div>

                          {/* Issues */}
                          {page.issues.length > 0 && (
                            <div>
                              <h4 className="font-medium text-sm mb-2 flex items-center gap-2">
                                <AlertTriangle className="h-4 w-4 text-red-500" />
                                Issues ({page.issues.length})
                              </h4>
                              <div className="space-y-2">
                                {page.issues.slice(0, 3).map((issue, issueIndex) => (
                                  <div key={issueIndex} className="flex items-start gap-2">
                                    {getSeverityIcon(issue.severity)}
                                    <div className="flex-1">
                                      <p className="text-xs font-medium">{issue.type.replace('-', ' ')}</p>
                                      <p className="text-xs text-gray-600">{issue.description}</p>
                                    </div>
                                  </div>
                                ))}
                                {page.issues.length > 3 && (
                                  <p className="text-xs text-gray-500">+{page.issues.length - 3} autres issues</p>
                                )}
                              </div>
                            </div>
                          )}

                        </div>

                        {/* Suggestions d'amélioration */}
                        <div className="mt-4 pt-4 border-t">
                          <h4 className="font-medium text-sm mb-2 flex items-center gap-2">
                            <Lightbulb className="h-4 w-4 text-yellow-500" />
                            Suggestions d'amélioration
                          </h4>
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            
                            {/* Longueur de contenu */}
                            <div className="text-xs">
                              <p className="font-medium">Longueur recommandée:</p>
                              <p className="text-gray-600">
                                {page.suggestions.contentLength.current} → {page.suggestions.contentLength.recommended} mots
                              </p>
                              <p className="text-gray-500 mt-1">{page.suggestions.contentLength.reason}</p>
                            </div>

                            {/* Sujets manquants */}
                            {page.suggestions.missingTopics.length > 0 && (
                              <div className="text-xs">
                                <p className="font-medium">Sujets à ajouter:</p>
                                <div className="flex flex-wrap gap-1 mt-1">
                                  {page.suggestions.missingTopics.slice(0, 3).map((topic, topicIndex) => (
                                    <Badge key={topicIndex} variant="outline" className="text-xs">
                                      {topic}
                                    </Badge>
                                  ))}
                                </div>
                              </div>
                            )}

                          </div>
                        </div>
                      </div>
                    )}
                  </div>
                ))}
              </div>

              {filteredAndSortedPages.length === 0 && (
                <div className="text-center py-8 text-gray-500">
                  <Search className="h-12 w-12 mx-auto mb-4 opacity-50" />
                  <p>Aucune page trouvée avec les filtres actuels</p>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      )}

      {activeTab === 'ai-insights' && (
        <div className="space-y-6">
          {/* Sujets couverts */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Brain className="h-5 w-5 text-purple-500" />
                Analyse sémantique IA
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                
                {/* Sujets couverts */}
                <div>
                  <h4 className="font-medium mb-4">Sujets couverts ({contentData.aiContentAnalysis.topicsCovered.length})</h4>
                  <div className="space-y-3">
                    {contentData.aiContentAnalysis.topicsCovered.map((topic, index) => (
                      <div key={index} className="flex items-center justify-between p-3 border rounded-lg">
                        <div className="flex-1">
                          <p className="font-medium text-sm">{topic.topic}</p>
                          <div className="flex items-center gap-4 mt-1">
                            <span className="text-xs text-gray-500">Couverture: {topic.coverage}%</span>
                            <span className="text-xs text-gray-500">Profondeur: {topic.depth}%</span>
                            <span className="text-xs text-gray-500">Autorité: {topic.authority}%</span>
                          </div>
                        </div>
                        <Badge variant="outline">{topic.searchIntent}</Badge>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Sujets manquants */}
                <div>
                  <h4 className="font-medium mb-4">Opportunités identifiées ({contentData.aiContentAnalysis.topicsMissing.length})</h4>
                  <div className="space-y-3">
                    {contentData.aiContentAnalysis.topicsMissing.map((topic, index) => (
                      <div key={index} className="p-3 border rounded-lg">
                        <div className="flex items-center justify-between mb-2">
                          <p className="font-medium text-sm">{topic.topic}</p>
                          <Badge variant={
                            topic.importance === TopicImportance.HIGH ? 'destructive' :
                            topic.importance === TopicImportance.MEDIUM ? 'default' : 'secondary'
                          }>
                            {topic.importance}
                          </Badge>
                        </div>
                        <div className="flex items-center gap-4 text-xs text-gray-600">
                          <span>Vol. recherche: {topic.searchVolume.toLocaleString()}</span>
                          <span>Difficulté: {topic.difficulty}/100</span>
                          <span>Opportunité: {topic.opportunity}%</span>
                        </div>
                        <div className="mt-2">
                          <p className="text-xs text-gray-500 mb-1">Mots-clés associés:</p>
                          <div className="flex flex-wrap gap-1">
                            {topic.relatedKeywords.slice(0, 3).map((kw, kwIndex) => (
                              <Badge key={kwIndex} variant="outline" className="text-xs">
                                {kw}
                              </Badge>
                            ))}
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

              </div>
            </CardContent>
          </Card>

          {/* Comparaison concurrentielle */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Users className="h-5 w-5 text-blue-500" />
                Analyse concurrentielle IA
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                
                {/* Métriques comparatives */}
                <div>
                  <h4 className="font-medium mb-4">Position competitive</h4>
                  <div className="space-y-4">
                    <Alert>
                      <AlertDescription>
                        <strong>Recommandation:</strong> {contentData.aiContentAnalysis.competitorComparison.lengthRecommendation}
                      </AlertDescription>
                    </Alert>
                    
                    <div className="grid grid-cols-2 gap-4">
                      <div className="text-center p-3 border rounded-lg">
                        <p className="text-2xl font-bold text-blue-600">
                          {contentData.aiContentAnalysis.competitorComparison.avgContentLength}
                        </p>
                        <p className="text-sm text-gray-600">Mots moy. concurrent</p>
                      </div>
                      <div className="text-center p-3 border rounded-lg">
                        <p className="text-2xl font-bold text-purple-600">
                          {contentData.aiContentAnalysis.competitorComparison.ourAvgLength}
                        </p>
                        <p className="text-sm text-gray-600">Notre moy.</p>
                      </div>
                    </div>
                    
                    <div className="space-y-2">
                      <div className="flex justify-between text-sm">
                        <span>Gap couverture sujets:</span>
                        <span className="text-red-600">-{contentData.aiContentAnalysis.competitorComparison.topicCoverageGap}%</span>
                      </div>
                      <div className="flex justify-between text-sm">
                        <span>Gap qualité:</span>
                        <span className="text-red-600">-{contentData.aiContentAnalysis.competitorComparison.qualityGap}%</span>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Top concurrents */}
                <div>
                  <h4 className="font-medium mb-4">Concurrents analysés</h4>
                  <div className="space-y-3">
                    {contentData.aiContentAnalysis.competitorComparison.topCompetitors.map((competitor, index) => (
                      <div key={index} className="p-3 border rounded-lg">
                        <div className="flex items-center justify-between mb-2">
                          <p className="font-medium text-sm">{competitor.domain}</p>
                          <Badge variant="secondary">{competitor.qualityScore}</Badge>
                        </div>
                        <div className="text-xs text-gray-600 mb-2">
                          <span>{competitor.contentLength} mots</span> • 
                          <span className="ml-1">{competitor.topicsCovered.length} sujets</span>
                        </div>
                        <div className="text-xs">
                          <p className="text-gray-500 mb-1">Forces vs nous:</p>
                          <ul className="list-disc list-inside text-gray-600 space-y-1">
                            {competitor.strengthsOverUs.slice(0, 2).map((strength, sIndex) => (
                              <li key={sIndex}>{strength}</li>
                            ))}
                          </ul>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

              </div>
            </CardContent>
          </Card>

          {/* Recommandations IA */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Sparkles className="h-5 w-5 text-yellow-500" />
                Recommandations IA stratégiques
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {contentData.aiContentAnalysis.aiRecommendations
                  .sort((a, b) => a.priority - b.priority)
                  .map((rec, index) => (
                  <div key={index} className="border rounded-lg p-4">
                    <div className="flex items-start justify-between mb-3">
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-2">
                          <Badge variant="outline">Priorité {rec.priority}</Badge>
                          <Badge variant={
                            rec.impact === 'high' ? 'destructive' :
                            rec.impact === 'medium' ? 'default' : 'secondary'
                          }>
                            Impact {rec.impact}
                          </Badge>
                          <Badge variant="outline">{rec.type.replace('_', ' ')}</Badge>
                        </div>
                        <h4 className="font-medium mb-2">{rec.description}</h4>
                        <p className="text-sm text-gray-600 mb-3">{rec.implementation}</p>
                      </div>
                    </div>
                    
                    <div className="bg-green-50 border border-green-200 rounded p-3">
                      <h5 className="font-medium text-sm text-green-800 mb-2">Résultats attendus:</h5>
                      <div className="grid grid-cols-1 md:grid-cols-3 gap-2 text-sm">
                        <div className="text-green-700">
                          <span className="font-medium">Trafic:</span> {rec.expectedResults.trafficIncrease}
                        </div>
                        <div className="text-green-700">
                          <span className="font-medium">Classement:</span> {rec.expectedResults.rankingImprovement}
                        </div>
                        <div className="text-green-700">
                          <span className="font-medium">Engagement:</span> {rec.expectedResults.engagementBoost}
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Analyse du ton */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <MessageSquare className="h-5 w-5 text-indigo-500" />
                Analyse du ton et style
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                
                <div>
                  <h4 className="font-medium mb-4">Métriques de ton</h4>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <span className="text-sm">Sentiment:</span>
                      <Badge variant={
                        contentData.aiContentAnalysis.contentTone.sentiment === 'positive' ? 'default' :
                        contentData.aiContentAnalysis.contentTone.sentiment === 'neutral' ? 'secondary' : 'destructive'
                      }>
                        {contentData.aiContentAnalysis.contentTone.sentiment}
                      </Badge>
                    </div>
                    
                    {[
                      { label: 'Formalité', value: contentData.aiContentAnalysis.contentTone.formality },
                      { label: 'Complexité', value: contentData.aiContentAnalysis.contentTone.complexity },
                      { label: 'Enthousiasme', value: contentData.aiContentAnalysis.contentTone.enthusiasm },
                      { label: 'Fiabilité', value: contentData.aiContentAnalysis.contentTone.trustworthiness },
                      { label: 'Expertise', value: contentData.aiContentAnalysis.contentTone.expertise },
                      { label: 'Cohérence', value: contentData.aiContentAnalysis.contentTone.toneConsistency },
                    ].map((metric, index) => (
                      <div key={index} className="space-y-1">
                        <div className="flex justify-between text-sm">
                          <span>{metric.label}:</span>
                          <span>{metric.value}/100</span>
                        </div>
                        <Progress value={metric.value} className="h-2" />
                      </div>
                    ))}
                  </div>
                </div>

                <div>
                  <h4 className="font-medium mb-4">Audience cible identifiée</h4>
                  <div className="space-y-2">
                    {contentData.aiContentAnalysis.contentTone.targetAudience.map((audience, index) => (
                      <Badge key={index} variant="outline" className="mr-2 mb-2">
                        {audience}
                      </Badge>
                    ))}
                  </div>
                  
                  <div className="mt-6 p-3 bg-blue-50 border border-blue-200 rounded">
                    <h5 className="font-medium text-blue-800 mb-2">Recommandations de ton</h5>
                    <ul className="text-sm text-blue-700 space-y-1">
                      {contentData.aiContentAnalysis.contentTone.expertise < 70 && (
                        <li>• Renforcer les signaux d'expertise (citations, données, cas d'étude)</li>
                      )}
                      {contentData.aiContentAnalysis.contentTone.trustworthiness < 75 && (
                        <li>• Améliorer la crédibilité (témoignages, preuves sociales)</li>
                      )}
                      {contentData.aiContentAnalysis.contentTone.toneConsistency < 80 && (
                        <li>• Harmoniser le ton entre les différentes sections</li>
                      )}
                    </ul>
                  </div>
                </div>

              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {activeTab === 'visualizations' && (
        <div className="space-y-6">
          <Alert>
            <LineChart className="h-4 w-4" />
            <AlertDescription>
              <strong>Visualisations avancées</strong> - Graphiques interactifs, heatmaps et analyses de réseau à venir.
              Cette section sera développée avec D3.js pour des visualisations riches des données de contenu.
            </AlertDescription>
          </Alert>

          {/* Placeholder pour les futures visualisations */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <BarChart3 className="h-5 w-5" />
                  Performance vs Longueur
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="h-64 flex items-center justify-center border-2 border-dashed border-gray-300 rounded-lg">
                  <div className="text-center text-gray-500">
                    <LineChart className="h-12 w-12 mx-auto mb-4 opacity-50" />
                    <p>Graphique interactif à venir</p>
                    <p className="text-sm">Corrélation longueur/performance</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Target className="h-5 w-5" />
                  Heatmap des mots-clés
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="h-64 flex items-center justify-center border-2 border-dashed border-gray-300 rounded-lg">
                  <div className="text-center text-gray-500">
                    <Target className="h-12 w-12 mx-auto mb-4 opacity-50" />
                    <p>Heatmap interactive à venir</p>
                    <p className="text-sm">Densité des mots-clés par page</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card className="md:col-span-2">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Link className="h-5 w-5" />
                  Graphe de maillage interne
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="h-80 flex items-center justify-center border-2 border-dashed border-gray-300 rounded-lg">
                  <div className="text-center text-gray-500">
                    <Link className="h-16 w-16 mx-auto mb-4 opacity-50" />
                    <p className="text-lg">Graphe de réseau à venir</p>
                    <p className="text-sm">Visualisation des liens internes avec D3.js</p>
                    <p className="text-sm text-gray-400 mt-2">Force-directed graph des relations entre pages</p>
                  </div>
                </div>
              </CardContent>
            </Card>

          </div>
        </div>
      )}

    </div>
  );
}