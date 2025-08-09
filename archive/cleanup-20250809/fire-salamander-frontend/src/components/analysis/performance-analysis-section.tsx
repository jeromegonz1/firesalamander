/**
 * Fire Salamander - Performance Analysis Section
 * Professional Core Web Vitals & Performance Analysis Component
 * Lead Tech implementation - Comparable to GTmetrix, PageSpeed Insights, WebPageTest
 */

"use client";

import React, { useState, useEffect, useMemo } from 'react';
import { 
  PerformanceAnalysis, 
  PagePerformanceMetrics,
  PerformanceRecommendation,
  CoreWebVitals,
  PerformanceGrade,
  DeviceType,
  RecommendationType,
  ImpactLevel,
  PerformanceSortBy,
  PerformanceFilterBy,
} from '@/types/performance-analysis';
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
  ArrowRight, Copy, RotateCcw, Users, Gauge, Smartphone,
  Monitor, Wifi, HardDrive, Minimize2, Maximize2, Database,
  Server, Network, Cpu, MemoryStick
} from 'lucide-react';

interface PerformanceAnalysisSectionProps {
  performanceData: PerformanceAnalysis;
  analysisId: string;
}

// ==================== MAIN COMPONENT ====================

export function PerformanceAnalysisSection({ 
  performanceData, 
  analysisId 
}: PerformanceAnalysisSectionProps) {
  const [activeTab, setActiveTab] = useState('overview');
  const [deviceView, setDeviceView] = useState<DeviceType>(DeviceType.MOBILE);
  const [recommendationFilter, setRecommendationFilter] = useState<PerformanceFilterBy>(PerformanceFilterBy.ALL);
  const [recommendationSort, setRecommendationSort] = useState<PerformanceSortBy>(PerformanceSortBy.SCORE);
  const [searchQuery, setSearchQuery] = useState('');
  const [expandedRecommendations, setExpandedRecommendations] = useState<Set<string>>(new Set());
  const [isLoading, setIsLoading] = useState(false);

  // Handle loading state
  if (!performanceData) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="text-center">
          <Activity className="h-8 w-8 animate-spin mx-auto mb-4 text-blue-500" />
          <p className="text-gray-600">Chargement de l'analyse performance...</p>
        </div>
      </div>
    );
  }

  // Handle empty data
  if (performanceData.summary.totalPages === 0) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="text-center">
          <Gauge className="h-12 w-12 mx-auto mb-4 text-gray-400" />
          <h3 className="text-lg font-semibold mb-2">Aucune données disponibles</h3>
          <p className="text-gray-600">Aucune page n'a été analysée pour les performances.</p>
        </div>
      </div>
    );
  }

  const currentPageMetrics = performanceData.pageMetrics[0];
  const currentDeviceMetrics = deviceView === DeviceType.MOBILE 
    ? currentPageMetrics?.mobile 
    : currentPageMetrics?.desktop;

  // Filter and sort recommendations
  const filteredRecommendations = useMemo(() => {
    let filtered = performanceData.recommendations;

    // Apply search filter
    if (searchQuery) {
      filtered = filtered.filter(rec => 
        rec.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        rec.description.toLowerCase().includes(searchQuery.toLowerCase())
      );
    }

    // Apply impact filter
    if (recommendationFilter !== PerformanceFilterBy.ALL) {
      const impactMap = {
        [PerformanceFilterBy.POOR]: ImpactLevel.HIGH,
        [PerformanceFilterBy.NEEDS_IMPROVEMENT]: ImpactLevel.MEDIUM,
        [PerformanceFilterBy.GOOD]: ImpactLevel.LOW,
      };
      const targetImpact = impactMap[recommendationFilter as keyof typeof impactMap];
      if (targetImpact) {
        filtered = filtered.filter(rec => rec.impact === targetImpact);
      }
    }

    // Apply sorting
    filtered.sort((a, b) => {
      switch (recommendationSort) {
        case PerformanceSortBy.SCORE:
          return b.estimatedGain.scoreImprovement - a.estimatedGain.scoreImprovement;
        case PerformanceSortBy.LOAD_TIME:
          return b.metrics.improvement - a.metrics.improvement;
        default:
          return 0;
      }
    });

    return filtered;
  }, [performanceData.recommendations, searchQuery, recommendationFilter, recommendationSort]);

  const toggleRecommendation = (id: string) => {
    const newExpanded = new Set(expandedRecommendations);
    if (newExpanded.has(id)) {
      newExpanded.delete(id);
    } else {
      newExpanded.add(id);
    }
    setExpandedRecommendations(newExpanded);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Analyse Performance</h1>
          <div className="flex items-center space-x-4 mt-2 text-gray-600">
            <span className="flex items-center space-x-1">
              <Globe className="h-4 w-4" />
              <span>{performanceData.summary.totalPages} page{performanceData.summary.totalPages > 1 ? 's' : ''} analysée{performanceData.summary.totalPages > 1 ? 's' : ''}</span>
            </span>
            <span className="flex items-center space-x-1">
              <Clock className="h-4 w-4" />
              <span>{new Date(performanceData.metadata.analysisDate).toLocaleString()}</span>
            </span>
            <span className="flex items-center space-x-1">
              <Activity className="h-4 w-4" />
              <span>{performanceData.summary.avgPerformanceScore}/100</span>
            </span>
          </div>
        </div>
        
        <div className="flex items-center space-x-3">
          <Button variant="outline" size="sm">
            <Download className="h-4 w-4 mr-2" />
            Exporter
          </Button>
          <Button variant="outline" size="sm">
            <RefreshCw className="h-4 w-4 mr-2" />
            Actualiser
          </Button>
        </div>
      </div>

      {/* Quick Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-600">Score Performance</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center space-x-2">
              <div className="text-2xl font-bold">
                {performanceData.summary.avgPerformanceScore}/100
              </div>
              <div className={`px-2 py-1 rounded text-xs font-medium ${
                performanceData.summary.avgPerformanceScore >= 90 ? 'bg-green-100 text-green-800' :
                performanceData.summary.avgPerformanceScore >= 75 ? 'bg-yellow-100 text-yellow-800' :
                'bg-red-100 text-red-800'
              }`}>
                {performanceData.summary.avgPerformanceScore >= 90 ? 'Excellent' :
                 performanceData.summary.avgPerformanceScore >= 75 ? 'Bon' : 'À améliorer'}
              </div>
            </div>
            <Progress 
              value={performanceData.summary.avgPerformanceScore} 
              className="mt-2"
            />
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-600">Temps de Chargement</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-1">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-1">
                  <Smartphone className="h-3 w-3" />
                  <span className="text-xs">Mobile</span>
                </div>
                <span className="text-sm font-medium">{(performanceData.summary.avgLoadTime.mobile / 1000).toFixed(1)}s</span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-1">
                  <Monitor className="h-3 w-3" />
                  <span className="text-xs">Desktop</span>
                </div>
                <span className="text-sm font-medium">{(performanceData.summary.avgLoadTime.desktop / 1000).toFixed(1)}s</span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-600">Core Web Vitals</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-1">
              <div className="flex items-center justify-between">
                <span className="text-xs">LCP</span>
                <span className={`text-xs px-1 rounded ${
                  performanceData.summary.coreWebVitals.mobile.lcp.passing >= 75 ? 'bg-green-100 text-green-800' :
                  performanceData.summary.coreWebVitals.mobile.lcp.passing >= 50 ? 'bg-yellow-100 text-yellow-800' :
                  'bg-red-100 text-red-800'
                }`}>
                  {performanceData.summary.coreWebVitals.mobile.lcp.passing}%
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-xs">FID</span>
                <span className={`text-xs px-1 rounded ${
                  performanceData.summary.coreWebVitals.mobile.fid.passing >= 75 ? 'bg-green-100 text-green-800' :
                  performanceData.summary.coreWebVitals.mobile.fid.passing >= 50 ? 'bg-yellow-100 text-yellow-800' :
                  'bg-red-100 text-red-800'
                }`}>
                  {performanceData.summary.coreWebVitals.mobile.fid.passing}%
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-xs">CLS</span>
                <span className={`text-xs px-1 rounded ${
                  performanceData.summary.coreWebVitals.mobile.cls.passing >= 75 ? 'bg-green-100 text-green-800' :
                  performanceData.summary.coreWebVitals.mobile.cls.passing >= 50 ? 'bg-yellow-100 text-yellow-800' :
                  'bg-red-100 text-red-800'
                }`}>
                  {performanceData.summary.coreWebVitals.mobile.cls.passing}%
                </span>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-600">Opportunités</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-1">
              <div className="flex items-center justify-between">
                <span className="text-xs">Impact élevé</span>
                <span className="text-sm font-bold text-orange-600">
                  {performanceData.summary.opportunities.highImpact}
                </span>
              </div>
              <div className="text-xs text-gray-600">
                Gain estimé: +{performanceData.summary.opportunities.estimatedGain.score} pts
              </div>
              <div className="text-xs text-gray-600">
                -{performanceData.summary.opportunities.estimatedGain.loadTime.toFixed(1)}s de chargement
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Vue d'ensemble</TabsTrigger>
          <TabsTrigger value="core-web-vitals">Core Web Vitals</TabsTrigger>
          <TabsTrigger value="recommendations">Recommandations</TabsTrigger>
          <TabsTrigger value="resources">Ressources</TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Score Distribution */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <BarChart3 className="h-5 w-5" />
                  <span>Distribution des Scores</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <div className="w-3 h-3 bg-green-500 rounded"></div>
                      <span className="text-sm">Excellent (90-100)</span>
                    </div>
                    <span className="font-medium">{performanceData.summary.scoreDistribution.excellent}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <div className="w-3 h-3 bg-yellow-500 rounded"></div>
                      <span className="text-sm">Bon (75-89)</span>
                    </div>
                    <span className="font-medium">{performanceData.summary.scoreDistribution.good}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <div className="w-3 h-3 bg-orange-500 rounded"></div>
                      <span className="text-sm">À améliorer (50-74)</span>
                    </div>
                    <span className="font-medium">{performanceData.summary.scoreDistribution.needsImprovement}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <div className="w-3 h-3 bg-red-500 rounded"></div>
                      <span className="text-sm">Pauvre (0-49)</span>
                    </div>
                    <span className="font-medium">{performanceData.summary.scoreDistribution.poor}</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Lighthouse Scores */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Gauge className="h-5 w-5" />
                  <span>Score global</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {currentPageMetrics && (
                    <>
                      <div className="flex items-center justify-between">
                        <span className="text-sm">Performance</span>
                        <div className="flex items-center space-x-2">
                          <Progress value={currentPageMetrics.scores.performance} className="w-20" />
                          <span className="font-medium">{currentPageMetrics.scores.performance}/100</span>
                        </div>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm">Accessibilité</span>
                        <div className="flex items-center space-x-2">
                          <Progress value={currentPageMetrics.scores.accessibility} className="w-20" />
                          <span className="font-medium">{currentPageMetrics.scores.accessibility}/100</span>
                        </div>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm">Bonnes pratiques</span>
                        <div className="flex items-center space-x-2">
                          <Progress value={currentPageMetrics.scores.bestPractices} className="w-20" />
                          <span className="font-medium">{currentPageMetrics.scores.bestPractices}/100</span>
                        </div>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm">SEO</span>
                        <div className="flex items-center space-x-2">
                          <Progress value={currentPageMetrics.scores.seo} className="w-20" />
                          <span className="font-medium">{currentPageMetrics.scores.seo}/100</span>
                        </div>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-sm">PWA</span>
                        <div className="flex items-center space-x-2">
                          <Progress value={currentPageMetrics.scores.pwa} className="w-20" />
                          <span className="font-medium">{currentPageMetrics.scores.pwa}/100</span>
                        </div>
                      </div>
                    </>
                  )}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Core Web Vitals Tab */}
        <TabsContent value="core-web-vitals" className="space-y-6">
          {/* Device Toggle */}
          <div className="flex items-center space-x-2">
            <Button
              variant={deviceView === DeviceType.MOBILE ? "default" : "outline"}
              size="sm"
              onClick={() => setDeviceView(DeviceType.MOBILE)}
            >
              <Smartphone className="h-4 w-4 mr-2" />
              Mobile
            </Button>
            <Button
              variant={deviceView === DeviceType.DESKTOP ? "default" : "outline"}
              size="sm"
              onClick={() => setDeviceView(DeviceType.DESKTOP)}
            >
              <Monitor className="h-4 w-4 mr-2" />
              Desktop
            </Button>
          </div>

          {currentDeviceMetrics && (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {/* LCP Card */}
              <CoreWebVitalCard
                title="Largest Contentful Paint"
                abbreviation="LCP"
                value={currentDeviceMetrics.lcp.value}
                unit="ms"
                grade={currentDeviceMetrics.lcp.grade}
                percentile={currentDeviceMetrics.lcp.percentile}
                description="Mesure le temps de chargement du plus grand élément visible"
                threshold="&lt; 2,5s"
                icon={<Eye className="h-5 w-5" />}
              />

              {/* FID Card */}
              <CoreWebVitalCard
                title="First Input Delay"
                abbreviation="FID"
                value={currentDeviceMetrics.fid.value}
                unit="ms"
                grade={currentDeviceMetrics.fid.grade}
                percentile={currentDeviceMetrics.fid.percentile}
                description="Mesure le délai avant la première interaction"
                threshold="&lt; 100ms"
                icon={<Zap className="h-5 w-5" />}
              />

              {/* CLS Card */}
              <CoreWebVitalCard
                title="Cumulative Layout Shift"
                abbreviation="CLS"
                value={currentDeviceMetrics.cls.value}
                unit=""
                grade={currentDeviceMetrics.cls.grade}
                percentile={currentDeviceMetrics.cls.percentile}
                description="Mesure la stabilité visuelle de la page"
                threshold="&lt; 0,1"
                icon={<Target className="h-5 w-5" />}
              />

              {/* FCP Card */}
              <CoreWebVitalCard
                title="First Contentful Paint"
                abbreviation="FCP"
                value={currentDeviceMetrics.fcp.value}
                unit="ms"
                grade={currentDeviceMetrics.fcp.grade}
                percentile={currentDeviceMetrics.fcp.percentile}
                description="Temps d'affichage du premier contenu"
                threshold="&lt; 1,8s"
                icon={<Activity className="h-5 w-5" />}
              />

              {/* TTI Card */}
              <CoreWebVitalCard
                title="Time to Interactive"
                abbreviation="TTI"
                value={currentDeviceMetrics.tti.value}
                unit="ms"
                grade={currentDeviceMetrics.tti.grade}
                percentile={currentDeviceMetrics.tti.percentile}
                description="Temps avant que la page soit interactive"
                threshold="&lt; 3,8s"
                icon={<Users className="h-5 w-5" />}
              />

              {/* Speed Index Card */}
              <CoreWebVitalCard
                title="Speed Index"
                abbreviation="SI"
                value={currentDeviceMetrics.speedIndex.value}
                unit="ms"
                grade={currentDeviceMetrics.speedIndex.grade}
                percentile={currentDeviceMetrics.speedIndex.percentile}
                description="Vitesse d'affichage du contenu visible"
                threshold="&lt; 3,4s"
                icon={<TrendingUp className="h-5 w-5" />}
              />
            </div>
          )}
        </TabsContent>

        {/* Recommendations Tab */}
        <TabsContent value="recommendations" className="space-y-6">
          {/* Filters and Search */}
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="relative">
                <Search className="h-4 w-4 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
                <Input
                  placeholder="Rechercher une recommandation..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 w-64"
                />
              </div>
              
              <select
                value={recommendationFilter}
                onChange={(e) => setRecommendationFilter(e.target.value as PerformanceFilterBy)}
                className="border rounded px-3 py-2 text-sm"
              >
                <option value={PerformanceFilterBy.ALL}>Tous les impacts</option>
                <option value={PerformanceFilterBy.POOR}>Impact élevé</option>
                <option value={PerformanceFilterBy.NEEDS_IMPROVEMENT}>Impact moyen</option>
                <option value={PerformanceFilterBy.GOOD}>Impact faible</option>
              </select>

              <select
                value={recommendationSort}
                onChange={(e) => setRecommendationSort(e.target.value as PerformanceSortBy)}
                className="border rounded px-3 py-2 text-sm"
              >
                <option value={PerformanceSortBy.SCORE}>Trier par impact</option>
                <option value={PerformanceSortBy.LOAD_TIME}>Par gain estimé</option>
              </select>
            </div>

            <div className="text-sm text-gray-600">
              {filteredRecommendations.length} recommandation{filteredRecommendations.length > 1 ? 's' : ''}
            </div>
          </div>

          {/* Recommendations List */}
          <div className="space-y-4">
            {filteredRecommendations.map((recommendation) => (
              <RecommendationCard
                key={recommendation.id}
                recommendation={recommendation}
                isExpanded={expandedRecommendations.has(recommendation.id)}
                onToggleExpanded={() => toggleRecommendation(recommendation.id)}
              />
            ))}
          </div>
        </TabsContent>

        {/* Resources Tab */}
        <TabsContent value="resources" className="space-y-6">
          {currentPageMetrics && (
            <>
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center space-x-2">
                    <HardDrive className="h-5 w-5" />
                    <span>Répartition des ressources</span>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                    <ResourceCard
                      title="Images"
                      count={currentPageMetrics.performance.resources.images.count}
                      size={currentPageMetrics.performance.resources.images.size}
                      icon={<Image className="h-5 w-5" />}
                      color="bg-blue-500"
                      details={[
                        `${currentPageMetrics.performance.resources.images.unoptimized} non optimisées`,
                        `${currentPageMetrics.performance.resources.images.withoutAlt} sans alt`,
                        `Formats: JPG (${currentPageMetrics.performance.resources.images.formats.jpg}), PNG (${currentPageMetrics.performance.resources.images.formats.png}), WebP (${currentPageMetrics.performance.resources.images.formats.webp})`
                      ]}
                    />
                    
                    <ResourceCard
                      title="CSS"
                      count={currentPageMetrics.performance.resources.css.count}
                      size={currentPageMetrics.performance.resources.css.size}
                      icon={<FileText className="h-5 w-5" />}
                      color="bg-green-500"
                      details={[
                        `${currentPageMetrics.performance.resources.css.unused}% inutilisé`,
                        `${currentPageMetrics.performance.resources.css.blocking} bloquant`,
                        `Minifié: ${currentPageMetrics.performance.resources.css.minified ? 'Oui' : 'Non'}`
                      ]}
                    />
                    
                    <ResourceCard
                      title="JavaScript"
                      count={currentPageMetrics.performance.resources.js.count}
                      size={currentPageMetrics.performance.resources.js.size}
                      icon={<Code className="h-5 w-5" />}
                      color="bg-yellow-500"
                      details={[
                        `${currentPageMetrics.performance.resources.js.unused}% inutilisé`,
                        `${currentPageMetrics.performance.resources.js.async} async`,
                        `${currentPageMetrics.performance.resources.js.defer} defer`
                      ]}
                    />
                    
                    <ResourceCard
                      title="Polices"
                      count={currentPageMetrics.performance.resources.fonts.count}
                      size={currentPageMetrics.performance.resources.fonts.size}
                      icon={<Type className="h-5 w-5" />}
                      color="bg-purple-500"
                      details={[
                        `${currentPageMetrics.performance.resources.fonts.preloaded} préchargée`,
                        `Fallbacks: ${currentPageMetrics.performance.resources.fonts.fallbacks ? 'Oui' : 'Non'}`,
                        `WOFF2: ${currentPageMetrics.performance.resources.fonts.formats.woff2}`
                      ]}
                    />
                  </div>
                </CardContent>
              </Card>

              {/* Optimization Status */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center space-x-2">
                    <Settings className="h-5 w-5" />
                    <span>État des Optimisations</span>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    <OptimizationCard
                      title="Compression"
                      enabled={currentPageMetrics.performance.optimization.compression.enabled}
                      details={[
                        `Algorithme: ${currentPageMetrics.performance.optimization.compression.algorithm || 'Non défini'}`,
                        `Ratio: ${currentPageMetrics.performance.optimization.compression.ratio}%`
                      ]}
                      icon={<Minimize2 className="h-4 w-4" />}
                    />

                    <OptimizationCard
                      title="CDN"
                      enabled={currentPageMetrics.performance.optimization.cdn.enabled}
                      details={[
                        `Fournisseur: ${currentPageMetrics.performance.optimization.cdn.provider || 'Non défini'}`,
                        `Couverture: ${currentPageMetrics.performance.optimization.cdn.coverage}%`
                      ]}
                      icon={<Network className="h-4 w-4" />}
                    />

                    <OptimizationCard
                      title="Minification CSS"
                      enabled={currentPageMetrics.performance.optimization.minification.css.enabled}
                      details={[
                        `Économies: ${formatBytes(currentPageMetrics.performance.optimization.minification.css.savings)}`,
                        `Ratio: ${currentPageMetrics.performance.optimization.minification.css.ratio}%`
                      ]}
                      icon={<FileText className="h-4 w-4" />}
                    />

                    <OptimizationCard
                      title="Minification JS"
                      enabled={currentPageMetrics.performance.optimization.minification.js.enabled}
                      details={[
                        `Économies: ${formatBytes(currentPageMetrics.performance.optimization.minification.js.savings)}`,
                        `Ratio: ${currentPageMetrics.performance.optimization.minification.js.ratio}%`
                      ]}
                      icon={<Code className="h-4 w-4" />}
                    />

                    <OptimizationCard
                      title="Cache navigateur"
                      enabled={currentPageMetrics.performance.optimization.caching.etags}
                      details={[
                        `ETags: ${currentPageMetrics.performance.optimization.caching.etags ? 'Activé' : 'Désactivé'}`,
                        `Last-Modified: ${currentPageMetrics.performance.optimization.caching.lastModified ? 'Activé' : 'Désactivé'}`
                      ]}
                      icon={<Database className="h-4 w-4" />}
                    />

                    <OptimizationCard
                      title="HTTP/2"
                      enabled={currentPageMetrics.performance.optimization.modernOptimizations.http2}
                      details={[
                        `Code splitting: ${currentPageMetrics.performance.optimization.modernOptimizations.codeSplitting ? 'Oui' : 'Non'}`,
                        `Lazy loading: ${currentPageMetrics.performance.optimization.modernOptimizations.lazyLoading ? 'Oui' : 'Non'}`
                      ]}
                      icon={<Wifi className="h-4 w-4" />}
                    />
                  </div>
                </CardContent>
              </Card>
            </>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
}

// ==================== SUB COMPONENTS ====================

interface CoreWebVitalCardProps {
  title: string;
  abbreviation: string;
  value: number;
  unit: string;
  grade: PerformanceGrade;
  percentile: number;
  description: string;
  threshold: string;
  icon: React.ReactNode;
}

function CoreWebVitalCard({ 
  title, 
  abbreviation, 
  value, 
  unit, 
  grade, 
  percentile, 
  description, 
  threshold, 
  icon 
}: CoreWebVitalCardProps) {
  const getGradeColor = (grade: PerformanceGrade) => {
    switch (grade) {
      case PerformanceGrade.EXCELLENT:
        return 'text-green-600 bg-green-100';
      case PerformanceGrade.GOOD:
        return 'text-blue-600 bg-blue-100';
      case PerformanceGrade.NEEDS_IMPROVEMENT:
        return 'text-yellow-600 bg-yellow-100';
      case PerformanceGrade.POOR:
        return 'text-red-600 bg-red-100';
      default:
        return 'text-gray-600 bg-gray-100';
    }
  };

  const getGradeText = (grade: PerformanceGrade) => {
    switch (grade) {
      case PerformanceGrade.EXCELLENT:
        return 'Excellent';
      case PerformanceGrade.GOOD:
        return 'Bon';
      case PerformanceGrade.NEEDS_IMPROVEMENT:
        return 'À améliorer';
      case PerformanceGrade.POOR:
        return 'Pauvre';
      default:
        return 'Non évalué';
    }
  };

  const formatValue = (val: number, unit: string) => {
    if (unit === 'ms') {
      return val >= 1000 ? `${(val / 1000).toFixed(1)}s` : `${val.toLocaleString()}ms`;
    }
    return `${val.toLocaleString()}${unit}`;
  };

  return (
    <Card>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            {icon}
            <div>
              <CardTitle className="text-lg">{abbreviation}</CardTitle>
              <p className="text-xs text-gray-600 font-normal">{title}</p>
            </div>
          </div>
          <Badge className={getGradeColor(grade)}>
            {getGradeText(grade)}
          </Badge>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        <div>
          <div className="text-2xl font-bold">{formatValue(value, unit)}</div>
          <div className="text-sm text-gray-600">
            {percentile}e percentile
          </div>
        </div>
        
        <div className="text-xs text-gray-600 leading-relaxed">
          {description}
        </div>
        
        <div className="pt-2 border-t">
          <div className="text-xs text-gray-500">
            Seuil recommandé: <span dangerouslySetInnerHTML={{ __html: threshold }} />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

interface RecommendationCardProps {
  recommendation: PerformanceRecommendation;
  isExpanded: boolean;
  onToggleExpanded: () => void;
}

function RecommendationCard({ recommendation, isExpanded, onToggleExpanded }: RecommendationCardProps) {
  const getImpactColor = (impact: ImpactLevel) => {
    switch (impact) {
      case ImpactLevel.HIGH:
        return 'bg-red-100 text-red-800';
      case ImpactLevel.MEDIUM:
        return 'bg-yellow-100 text-yellow-800';
      case ImpactLevel.LOW:
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const getImpactText = (impact: ImpactLevel) => {
    switch (impact) {
      case ImpactLevel.HIGH:
        return 'Impact élevé';
      case ImpactLevel.MEDIUM:
        return 'Impact moyen';
      case ImpactLevel.LOW:
        return 'Impact faible';
      default:
        return 'Impact non défini';
    }
  };

  const getRecommendationIcon = (type: RecommendationType) => {
    switch (type) {
      case RecommendationType.IMAGES:
        return <Image className="h-5 w-5" />;
      case RecommendationType.CSS:
        return <FileText className="h-5 w-5" />;
      case RecommendationType.JAVASCRIPT:
        return <Code className="h-5 w-5" />;
      case RecommendationType.FONTS:
        return <Type className="h-5 w-5" />;
      case RecommendationType.CACHING:
        return <Database className="h-5 w-5" />;
      case RecommendationType.COMPRESSION:
        return <Minimize2 className="h-5 w-5" />;
      case RecommendationType.CDN:
        return <Network className="h-5 w-5" />;
      default:
        return <Settings className="h-5 w-5" />;
    }
  };

  return (
    <Card className="transition-all duration-200 hover:shadow-md">
      <CardHeader className="cursor-pointer" onClick={onToggleExpanded}>
        <div className="flex items-start justify-between">
          <div className="flex items-start space-x-3">
            <div className="p-2 rounded-lg bg-blue-50 text-blue-600">
              {getRecommendationIcon(recommendation.type)}
            </div>
            <div className="flex-1">
              <div className="flex items-center space-x-2 mb-1">
                <CardTitle className="text-lg">{recommendation.title}</CardTitle>
                <Badge className={getImpactColor(recommendation.impact)}>
                  {getImpactText(recommendation.impact)}
                </Badge>
              </div>
              <p className="text-sm text-gray-600">{recommendation.description}</p>
              <div className="flex items-center space-x-4 mt-2 text-sm text-gray-500">
                <span>+{recommendation.estimatedGain.scoreImprovement} points</span>
                <span>{recommendation.estimatedGain.loadTime}</span>
                <span>{recommendation.estimatedGain.bandwidth}</span>
              </div>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            <div className="text-right">
              <div className="text-sm font-medium">
                {recommendation.metrics.improvement}% d'amélioration
              </div>
              <div className="text-xs text-gray-500">
                {recommendation.solution.estimatedTime}
              </div>
            </div>
            {isExpanded ? <ChevronDown className="h-5 w-5" /> : <ChevronRight className="h-5 w-5" />}
          </div>
        </div>
      </CardHeader>
      
      {isExpanded && (
        <CardContent className="pt-0">
          <div className="border-t pt-4 space-y-4">
            <div>
              <h4 className="font-medium mb-2">Solution recommandée</h4>
              <p className="text-sm text-gray-600 mb-2">{recommendation.solution.description}</p>
              <p className="text-sm text-gray-600">{recommendation.solution.implementation}</p>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <span className="text-xs font-medium text-gray-500">DIFFICULTÉ</span>
                <div className="capitalize">{recommendation.solution.difficulty}</div>
              </div>
              <div>
                <span className="text-xs font-medium text-gray-500">TEMPS ESTIMÉ</span>
                <div>{recommendation.solution.estimatedTime}</div>
              </div>
              <div>
                <span className="text-xs font-medium text-gray-500">GAIN UX</span>
                <div className="text-sm">{recommendation.estimatedGain.userExperience}</div>
              </div>
            </div>

            <div>
              <h4 className="font-medium mb-2">Impact détaillé</h4>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-3 text-sm">
                <div className="p-2 bg-gray-50 rounded">
                  <div className="text-xs text-gray-500">AVANT</div>
                  <div className="font-medium">{formatBytes(recommendation.metrics.before)}</div>
                </div>
                <div className="p-2 bg-gray-50 rounded">
                  <div className="text-xs text-gray-500">APRÈS</div>
                  <div className="font-medium">{formatBytes(recommendation.metrics.after)}</div>
                </div>
                <div className="p-2 bg-green-50 rounded">
                  <div className="text-xs text-gray-500">ÉCONOMIES</div>
                  <div className="font-medium text-green-600">
                    {formatBytes(recommendation.metrics.before - recommendation.metrics.after)}
                  </div>
                </div>
                <div className="p-2 bg-blue-50 rounded">
                  <div className="text-xs text-gray-500">AMÉLIORATION</div>
                  <div className="font-medium text-blue-600">{recommendation.metrics.improvement}%</div>
                </div>
              </div>
            </div>

            <div className="flex justify-between items-center pt-2">
              <div className="text-xs text-gray-500">
                Pages affectées: {recommendation.pages.length}
              </div>
              <Button size="sm" variant="outline">
                <ExternalLink className="h-4 w-4 mr-2" />
                Voir la documentation
              </Button>
            </div>
          </div>
        </CardContent>
      )}
    </Card>
  );
}

interface ResourceCardProps {
  title: string;
  count: number;
  size: number;
  icon: React.ReactNode;
  color: string;
  details: string[];
}

function ResourceCard({ title, count, size, icon, color, details }: ResourceCardProps) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <div className="flex items-center space-x-2">
          <div className={`p-2 rounded ${color} text-white`}>
            {icon}
          </div>
          <CardTitle className="text-base">{title}</CardTitle>
        </div>
      </CardHeader>
      <CardContent>
        <div className="space-y-2">
          <div>
            <div className="text-lg font-bold">{count} fichiers</div>
            <div className="text-sm text-gray-600">{formatBytes(size)}</div>
          </div>
          <div className="space-y-1">
            {details.map((detail, index) => (
              <div key={index} className="text-xs text-gray-500">
                {detail}
              </div>
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

interface OptimizationCardProps {
  title: string;
  enabled: boolean;
  details: string[];
  icon: React.ReactNode;
}

function OptimizationCard({ title, enabled, details, icon }: OptimizationCardProps) {
  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex items-start space-x-3">
          <div className={`p-2 rounded-lg ${enabled ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-400'}`}>
            {icon}
          </div>
          <div className="flex-1">
            <div className="flex items-center space-x-2 mb-1">
              <h3 className="font-medium">{title}</h3>
              {enabled ? (
                <CheckCircle2 className="h-4 w-4 text-green-500" />
              ) : (
                <XCircle className="h-4 w-4 text-red-500" />
              )}
            </div>
            <div className="space-y-1">
              {details.map((detail, index) => (
                <div key={index} className="text-xs text-gray-600">
                  {detail}
                </div>
              ))}
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

// ==================== UTILITY FUNCTIONS ====================

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Import missing components
function Code({ className, ...props }: any) {
  return <FileText className={className} {...props} />;
}

function Type({ className, ...props }: any) {
  return <FileText className={className} {...props} />;
}