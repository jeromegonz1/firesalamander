/**
 * Fire Salamander - Overview Analysis Section Component
 * Lead Tech quality overview dashboard with comprehensive SEO metrics
 */

'use client';

import React, { useState, useMemo } from 'react';
import { OverviewAnalysis, ModuleSummary, OverviewGrade, IssueItem, RecommendationItem } from '@/types/overview';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Progress } from '@/components/ui/progress';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  TrendingUp,
  TrendingDown,
  Minus,
  AlertTriangle,
  CheckCircle,
  Clock,
  Target,
  BarChart3,
  Wrench,
  Zap,
  FileText,
  Shield,
  Link2,
  Award,
  Users,
  Activity,
  Calendar,
  ExternalLink,
  RefreshCw,
  Download,
  ArrowUp,
  ArrowDown,
  Eye,
} from 'lucide-react';
import Link from 'next/link';

interface OverviewSectionProps {
  overviewData: OverviewAnalysis;
  analysisId: string;
}

export function OverviewSection({ overviewData, analysisId }: OverviewSectionProps) {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [refreshing, setRefreshing] = useState(false);

  const getTrendIcon = (trend: 'up' | 'down' | 'stable' | 'improving' | 'declining') => {
    switch (trend) {
      case 'up':
      case 'improving':
        return <TrendingUp className="h-4 w-4 text-green-500" />;
      case 'down':
      case 'declining':
        return <TrendingDown className="h-4 w-4 text-red-500" />;
      default:
        return <Minus className="h-4 w-4 text-gray-500" />;
    }
  };

  const getGradeColor = (grade: OverviewGrade) => {
    switch (grade) {
      case OverviewGrade.A_PLUS:
      case OverviewGrade.A:
        return 'text-green-600';
      case OverviewGrade.B:
        return 'text-blue-600';
      case OverviewGrade.C:
        return 'text-yellow-600';
      case OverviewGrade.D:
        return 'text-orange-600';
      case OverviewGrade.F:
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  };

  const getModuleIcon = (module: string) => {
    const icons = {
      technical: Wrench,
      performance: Zap,
      content: FileText,
      security: Shield,
      backlinks: Link2,
    };
    const IconComponent = icons[module as keyof typeof icons] || BarChart3;
    return <IconComponent className="h-5 w-5" />;
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'text-red-600 bg-red-50';
      case 'high':
        return 'text-orange-600 bg-orange-50';
      case 'medium':
        return 'text-yellow-600 bg-yellow-50';
      case 'low':
        return 'text-blue-600 bg-blue-50';
      default:
        return 'text-gray-600 bg-gray-50';
    }
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    // Simulate API call
    setTimeout(() => {
      setRefreshing(false);
    }, 2000);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Vue d'ensemble SEO</h1>
          <div className="flex items-center space-x-4 mt-2">
            <div className="flex items-center space-x-2">
              <div className={`text-4xl font-bold ${getGradeColor(overviewData.score.grade)}`}>
                {overviewData.score.grade}
              </div>
              <div className="text-2xl font-semibold text-gray-900">
                {overviewData.score.overall}/100
              </div>
              {getTrendIcon(overviewData.score.trend)}
            </div>
            <div className="text-sm text-gray-600">
              {overviewData.siteInfo.url}
            </div>
          </div>
        </div>
        
        <div className="flex items-center space-x-2">
          <Button onClick={handleRefresh} disabled={refreshing} variant="outline" size="sm">
            <RefreshCw className={`h-4 w-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
            Actualiser
          </Button>
          <Button variant="outline" size="sm">
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
          <Link href={`/analysis/${analysisId}/report`}>
            <Button size="sm">
              <Eye className="h-4 w-4 mr-2" />
              Rapport complet
            </Button>
          </Link>
        </div>
      </div>

      {/* Main Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="dashboard">Dashboard</TabsTrigger>
          <TabsTrigger value="modules">Modules</TabsTrigger>
          <TabsTrigger value="issues">Problèmes</TabsTrigger>
          <TabsTrigger value="recommendations">Recommandations</TabsTrigger>
        </TabsList>

        {/* Dashboard Tab */}
        <TabsContent value="dashboard" className="space-y-6">
          {/* KPIs Cards */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600">Trafic Organique</p>
                    <div className="flex items-center space-x-2">
                      <p className="text-2xl font-bold">{overviewData.kpis.organicTraffic.current.toLocaleString()}</p>
                      {getTrendIcon(overviewData.kpis.organicTraffic.trend)}
                    </div>
                    <p className="text-xs text-gray-500">
                      {overviewData.kpis.organicTraffic.change > 0 ? '+' : ''}
                      {overviewData.kpis.organicTraffic.change}% ce mois
                    </p>
                  </div>
                  <Activity className="h-8 w-8 text-blue-500" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600">Mots-clés Top 10</p>
                    <div className="flex items-center space-x-2">
                      <p className="text-2xl font-bold">{overviewData.kpis.keywordRankings.top10}</p>
                      <div className="text-green-500">
                        <ArrowUp className="h-4 w-4" />
                      </div>
                    </div>
                    <p className="text-xs text-gray-500">
                      {overviewData.kpis.keywordRankings.total} mots-clés total
                    </p>
                  </div>
                  <Target className="h-8 w-8 text-green-500" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600">Domain Authority</p>
                    <div className="flex items-center space-x-2">
                      <p className="text-2xl font-bold">{overviewData.kpis.domainAuthority.current}</p>
                      {getTrendIcon(overviewData.kpis.domainAuthority.trend)}
                    </div>
                    <p className="text-xs text-gray-500">
                      {overviewData.kpis.domainAuthority.change > 0 ? '+' : ''}
                      {overviewData.kpis.domainAuthority.change} points
                    </p>
                  </div>
                  <Award className="h-8 w-8 text-purple-500" />
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600">Page Speed</p>
                    <div className="flex items-center space-x-2">
                      <p className="text-2xl font-bold">{overviewData.kpis.pageSpeed.mobile}</p>
                      <span className="text-sm text-gray-500">/{overviewData.kpis.pageSpeed.desktop}</span>
                    </div>
                    <p className="text-xs text-gray-500">Mobile / Desktop</p>
                  </div>
                  <Zap className="h-8 w-8 text-yellow-500" />
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Progress & Competition */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <BarChart3 className="h-5 w-5" />
                  <span>Progression Globale</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-sm font-medium">Tâches complétées</span>
                    <span className="text-sm font-semibold">
                      {overviewData.progress.completedTasks}/{overviewData.progress.totalTasks}
                    </span>
                  </div>
                  <Progress 
                    value={(overviewData.progress.completedTasks / overviewData.progress.totalTasks) * 100} 
                    className="h-3" 
                  />
                </div>
                
                <div className="space-y-2 text-sm">
                  <div className="flex items-center justify-between">
                    <span className="text-gray-600">Dernière activité</span>
                    <span>{overviewData.progress.lastActivity}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-gray-600">Prochain objectif</span>
                    <span className="font-medium">{overviewData.progress.nextMilestone}</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Users className="h-5 w-5" />
                  <span>Position Concurrentielle</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="text-center">
                  <div className="text-3xl font-bold text-blue-600">
                    #{overviewData.competition.position}
                  </div>
                  <div className="text-sm text-gray-600">
                    sur {overviewData.competition.totalCompetitors} concurrents
                  </div>
                </div>
                
                <div className="space-y-3">
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-gray-600">Part de marché</span>
                    <span className="font-semibold">{overviewData.competition.marketShare}%</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-sm text-gray-600">Performance</span>
                    <span className={`font-semibold ${
                      overviewData.competition.aboveAverage ? 'text-green-600' : 'text-red-600'
                    }`}>
                      {overviewData.competition.aboveAverage ? 'Au-dessus' : 'En-dessous'} de la moyenne
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Site Information */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <ExternalLink className="h-5 w-5" />
                <span>Informations du Site</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div>
                  <div className="text-sm text-gray-600">URL Analysée</div>
                  <div className="font-medium">{overviewData.siteInfo.url}</div>
                </div>
                <div>
                  <div className="text-sm text-gray-600">Pages Analysées</div>
                  <div className="font-medium">{overviewData.siteInfo.pagesAnalyzed.toLocaleString()}</div>
                </div>
                <div>
                  <div className="text-sm text-gray-600">Dernière Analyse</div>
                  <div className="font-medium">
                    {new Date(overviewData.siteInfo.lastCrawled).toLocaleDateString('fr-FR')}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Modules Tab */}
        <TabsContent value="modules" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {Object.entries(overviewData.modules).map(([moduleKey, module]) => (
              <Card key={moduleKey}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center space-x-3">
                      {getModuleIcon(moduleKey)}
                      <div>
                        <div className="font-semibold capitalize">{moduleKey}</div>
                        <div className="text-sm text-gray-600">
                          Mis à jour {new Date(module.lastUpdated).toLocaleDateString('fr-FR')}
                        </div>
                      </div>
                    </div>
                    <div className="text-right">
                      <div className={`text-2xl font-bold ${getGradeColor(module.grade)}`}>
                        {module.grade}
                      </div>
                      <div className="text-sm text-gray-600">{module.score}/100</div>
                    </div>
                  </div>

                  <div className="space-y-3">
                    <Progress value={module.score} className="h-2" />
                    
                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div className="flex items-center justify-between">
                        <span className="text-gray-600">Problèmes critiques</span>
                        <span className="font-semibold text-red-600">{module.criticalIssues}</span>
                      </div>
                      <div className="flex items-center justify-between">
                        <span className="text-gray-600">Améliorations</span>
                        <span className="font-semibold text-green-600">{module.improvements}</span>
                      </div>
                    </div>

                    <div className="flex justify-end">
                      <Link href={`/analysis/${analysisId}/${moduleKey}`}>
                        <Button variant="outline" size="sm">
                          Voir détails
                          <ExternalLink className="h-4 w-4 ml-2" />
                        </Button>
                      </Link>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        {/* Issues Tab */}
        <TabsContent value="issues" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Critical Issues */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2 text-red-600">
                  <AlertTriangle className="h-5 w-5" />
                  <span>Problèmes Critiques ({overviewData.issues.critical.length})</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {overviewData.issues.critical.slice(0, 5).map((issue, index) => (
                  <div key={issue.id} className="p-3 border border-red-200 rounded-lg">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="font-medium text-red-900">{issue.title}</div>
                        <div className="text-sm text-red-700 mt-1">{issue.description}</div>
                        <div className="flex items-center space-x-4 mt-2 text-xs">
                          <span className="flex items-center space-x-1">
                            {getModuleIcon(issue.module)}
                            <span className="capitalize">{issue.module}</span>
                          </span>
                          <span>{issue.affectedPages} pages</span>
                          <span>{issue.estimatedFix}</span>
                        </div>
                      </div>
                      <span className={`px-2 py-1 text-xs rounded-full font-medium ${getSeverityColor(issue.severity)}`}>
                        {issue.severity}
                      </span>
                    </div>
                  </div>
                ))}
              </CardContent>
            </Card>

            {/* Warnings */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2 text-orange-600">
                  <Clock className="h-5 w-5" />
                  <span>Avertissements ({overviewData.issues.warnings.length})</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {overviewData.issues.warnings.slice(0, 5).map((issue, index) => (
                  <div key={issue.id} className="p-3 border border-orange-200 rounded-lg">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="font-medium text-orange-900">{issue.title}</div>
                        <div className="text-sm text-orange-700 mt-1">{issue.description}</div>
                        <div className="flex items-center space-x-4 mt-2 text-xs">
                          <span className="flex items-center space-x-1">
                            {getModuleIcon(issue.module)}
                            <span className="capitalize">{issue.module}</span>
                          </span>
                          <span>{issue.affectedPages} pages</span>
                          <span>{issue.estimatedFix}</span>
                        </div>
                      </div>
                      <span className={`px-2 py-1 text-xs rounded-full font-medium ${getSeverityColor(issue.severity)}`}>
                        {issue.severity}
                      </span>
                    </div>
                  </div>
                ))}
              </CardContent>
            </Card>
          </div>

          {/* Issues Summary */}
          <Card>
            <CardHeader>
              <CardTitle>Résumé des Problèmes</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="text-center">
                  <div className="text-2xl font-bold text-red-600">{overviewData.issues.totalIssues}</div>
                  <div className="text-sm text-gray-600">Total des problèmes</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-green-600">{overviewData.issues.resolvedIssues}</div>
                  <div className="text-sm text-gray-600">Problèmes résolus</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-blue-600">
                    {Math.round((overviewData.issues.resolvedIssues / overviewData.issues.totalIssues) * 100)}%
                  </div>
                  <div className="text-sm text-gray-600">Taux de résolution</div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Recommendations Tab */}
        <TabsContent value="recommendations" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Immediate Actions */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <AlertTriangle className="h-5 w-5 text-red-500" />
                  <span>Actions Immédiates</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {overviewData.recommendations.immediate.slice(0, 5).map((rec, index) => (
                  <div key={rec.id} className="p-3 border border-gray-200 rounded-lg">
                    <div className="flex items-start space-x-3">
                      <div className="w-6 h-6 rounded-full bg-red-500 text-white flex items-center justify-center text-xs font-bold">
                        {rec.priority}
                      </div>
                      <div className="flex-1">
                        <div className="font-medium">{rec.title}</div>
                        <div className="text-sm text-gray-600 mt-1">{rec.description}</div>
                        <div className="flex items-center space-x-4 mt-2 text-xs">
                          <span className="flex items-center space-x-1">
                            {getModuleIcon(rec.module)}
                            <span className="capitalize">{rec.module}</span>
                          </span>
                          <span>Impact: {rec.expectedImpact}</span>
                          <span>Effort: {rec.effort}</span>
                          <span>Délai: {rec.timeline}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </CardContent>
            </Card>

            {/* Quick Wins */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <CheckCircle className="h-5 w-5 text-green-500" />
                  <span>Victoires Rapides</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {overviewData.recommendations.quickWins.slice(0, 5).map((rec, index) => (
                  <div key={rec.id} className="p-3 border border-gray-200 rounded-lg">
                    <div className="flex items-start space-x-3">
                      <div className="w-6 h-6 rounded-full bg-green-500 text-white flex items-center justify-center text-xs font-bold">
                        {rec.priority}
                      </div>
                      <div className="flex-1">
                        <div className="font-medium">{rec.title}</div>
                        <div className="text-sm text-gray-600 mt-1">{rec.description}</div>
                        <div className="flex items-center space-x-4 mt-2 text-xs">
                          <span className="flex items-center space-x-1">
                            {getModuleIcon(rec.module)}
                            <span className="capitalize">{rec.module}</span>
                          </span>
                          <span>Impact: {rec.expectedImpact}</span>
                          <span>Effort: {rec.effort}</span>
                          <span>Délai: {rec.timeline}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}