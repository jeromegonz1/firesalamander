/**
 * Fire Salamander - Backlinks Analysis Section Component
 * Enterprise-grade backlinks analysis with Lead Tech quality
 * Professional SEO backlinks dashboard comparable to industry leaders
 */

'use client';

import React, { useState, useMemo, useCallback } from 'react';
import { 
  BacklinksAnalysis, 
  BacklinkData,
  BacklinkStatus,
  BacklinkQuality,
  LinkType,
  AnchorType,
  BacklinkGrade,
  DomainCategory,
  ToxicityReason,
  ReferringDomain,
  BacklinkOpportunity,
} from '@/types/backlinks-analysis';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Progress } from '@/components/ui/progress';
import {
  ChevronDown,
  ChevronUp,
  ExternalLink,
  TrendingUp,
  TrendingDown,
  Minus,
  AlertTriangle,
  Shield,
  Globe,
  Link2,
  Target,
  Award,
  Search,
  Filter,
  Download,
  Eye,
  EyeOff,
  RefreshCw,
  BarChart3,
  PieChart,
  Activity,
  Clock,
  MapPin,
  Flag,
  Zap,
  Users,
  CheckCircle,
  XCircle,
  AlertCircle,
  Info,
  Star,
} from 'lucide-react';

interface BacklinksAnalysisSectionProps {
  backlinksData: BacklinksAnalysis;
  analysisId: string;
}

export function BacklinksAnalysisSection({ backlinksData, analysisId }: BacklinksAnalysisSectionProps) {
  const [activeTab, setActiveTab] = useState('overview');
  const [expandedSections, setExpandedSections] = useState<Record<string, boolean>>({});
  const [sortBy, setSortBy] = useState<string>('domainAuthority');
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc');
  const [filterStatus, setFilterStatus] = useState<BacklinkStatus | 'all'>('all');
  const [filterQuality, setFilterQuality] = useState<BacklinkQuality | 'all'>('all');
  const [searchTerm, setSearchTerm] = useState('');
  const [currentPage, setCurrentPage] = useState(1);
  const [showAdvancedFilters, setShowAdvancedFilters] = useState(false);
  const [refreshing, setRefreshing] = useState(false);

  const itemsPerPage = 20;

  // Memoized calculations for performance
  const filteredBacklinks = useMemo(() => {
    let filtered = [...backlinksData.backlinks.list];

    // Apply search filter
    if (searchTerm) {
      const term = searchTerm.toLowerCase();
      filtered = filtered.filter(backlink =>
        backlink.sourceDomain.toLowerCase().includes(term) ||
        backlink.anchorText.toLowerCase().includes(term) ||
        backlink.sourceUrl.toLowerCase().includes(term)
      );
    }

    // Apply status filter
    if (filterStatus !== 'all') {
      filtered = filtered.filter(backlink => backlink.status === filterStatus);
    }

    // Apply quality filter
    if (filterQuality !== 'all') {
      filtered = filtered.filter(backlink => backlink.quality === filterQuality);
    }

    return filtered;
  }, [backlinksData.backlinks.list, searchTerm, filterStatus, filterQuality]);

  const sortedBacklinks = useMemo(() => {
    const sorted = [...filteredBacklinks];
    
    sorted.sort((a, b) => {
      let aVal: any = a;
      let bVal: any = b;

      switch (sortBy) {
        case 'domainAuthority':
          aVal = a.domainAuthority.domainAuthority;
          bVal = b.domainAuthority.domainAuthority;
          break;
        case 'firstSeen':
          aVal = new Date(a.firstSeen).getTime();
          bVal = new Date(b.firstSeen).getTime();
          break;
        case 'sourceDomain':
          aVal = a.sourceDomain.toLowerCase();
          bVal = b.sourceDomain.toLowerCase();
          break;
        case 'toxicityScore':
          aVal = a.toxicityScore;
          bVal = b.toxicityScore;
          break;
        default:
          return 0;
      }

      if (typeof aVal === 'string') {
        return sortDirection === 'asc' ? aVal.localeCompare(bVal) : bVal.localeCompare(aVal);
      }

      return sortDirection === 'asc' ? aVal - bVal : bVal - aVal;
    });

    return sorted;
  }, [filteredBacklinks, sortBy, sortDirection]);

  const paginatedBacklinks = useMemo(() => {
    const startIndex = (currentPage - 1) * itemsPerPage;
    return sortedBacklinks.slice(startIndex, startIndex + itemsPerPage);
  }, [sortedBacklinks, currentPage]);

  const totalPages = Math.ceil(sortedBacklinks.length / itemsPerPage);

  const toggleSection = useCallback((sectionId: string) => {
    setExpandedSections(prev => ({
      ...prev,
      [sectionId]: !prev[sectionId]
    }));
  }, []);

  const handleSort = useCallback((field: string) => {
    if (sortBy === field) {
      setSortDirection(prev => prev === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(field);
      setSortDirection('desc');
    }
    setCurrentPage(1);
  }, [sortBy]);

  const handleRefresh = useCallback(async () => {
    setRefreshing(true);
    // Simulate API call
    setTimeout(() => {
      setRefreshing(false);
    }, 2000);
  }, []);

  const getStatusIcon = (status: BacklinkStatus) => {
    switch (status) {
      case BacklinkStatus.ACTIVE:
        return <CheckCircle className="h-4 w-4 text-green-500" />;
      case BacklinkStatus.TOXIC:
        return <XCircle className="h-4 w-4 text-red-500" />;
      case BacklinkStatus.LOST:
        return <AlertCircle className="h-4 w-4 text-orange-500" />;
      case BacklinkStatus.NEW:
        return <Star className="h-4 w-4 text-blue-500" />;
      case BacklinkStatus.BROKEN:
        return <AlertTriangle className="h-4 w-4 text-red-500" />;
      default:
        return <Info className="h-4 w-4 text-gray-500" />;
    }
  };

  const getQualityColor = (quality: BacklinkQuality) => {
    switch (quality) {
      case BacklinkQuality.EXCELLENT:
        return 'text-green-600 bg-green-50';
      case BacklinkQuality.GOOD:
        return 'text-blue-600 bg-blue-50';
      case BacklinkQuality.AVERAGE:
        return 'text-yellow-600 bg-yellow-50';
      case BacklinkQuality.POOR:
        return 'text-orange-600 bg-orange-50';
      case BacklinkQuality.TOXIC:
        return 'text-red-600 bg-red-50';
      default:
        return 'text-gray-600 bg-gray-50';
    }
  };

  const getGradeColor = (grade: BacklinkGrade) => {
    switch (grade) {
      case BacklinkGrade.A_PLUS:
      case BacklinkGrade.A:
        return 'text-green-600';
      case BacklinkGrade.B:
        return 'text-blue-600';
      case BacklinkGrade.C:
        return 'text-yellow-600';
      case BacklinkGrade.D:
        return 'text-orange-600';
      case BacklinkGrade.F:
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  };

  const getTrendIcon = (trend: string) => {
    switch (trend) {
      case 'improving':
        return <TrendingUp className="h-4 w-4 text-green-500" />;
      case 'declining':
        return <TrendingDown className="h-4 w-4 text-red-500" />;
      default:
        return <Minus className="h-4 w-4 text-gray-500" />;
    }
  };

  return (
    <div className="space-y-6">
      {/* Header with Score and Actions */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2">
            <div className={`text-3xl font-bold ${getGradeColor(backlinksData.score.grade)}`}>
              {backlinksData.score.grade}
            </div>
            <div className="text-2xl font-semibold text-gray-900">
              {backlinksData.score.overall}/100
            </div>
            {getTrendIcon(backlinksData.score.trend)}
          </div>
          <div className="text-sm text-gray-600">
            <div>Backlinks Analysis</div>
            <div className="flex items-center space-x-2">
              <Link2 className="h-4 w-4" />
              <span>{backlinksData.profile.totalBacklinks.toLocaleString()} backlinks</span>
              <Globe className="h-4 w-4 ml-2" />
              <span>{backlinksData.profile.totalReferringDomains.toLocaleString()} domains</span>
            </div>
          </div>
        </div>
        
        <div className="flex items-center space-x-2">
          <Button
            onClick={handleRefresh}
            disabled={refreshing}
            size="sm"
            variant="outline"
          >
            <RefreshCw className={`h-4 w-4 mr-2 ${refreshing ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
          <Button size="sm" variant="outline">
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
          <Button size="sm" variant="outline">
            <Eye className="h-4 w-4 mr-2" />
            Monitor
          </Button>
        </div>
      </div>

      {/* Main Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="grid w-full grid-cols-7">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="backlinks">Backlinks</TabsTrigger>
          <TabsTrigger value="domains">Domains</TabsTrigger>
          <TabsTrigger value="anchors">Anchors</TabsTrigger>
          <TabsTrigger value="competitors">Competitors</TabsTrigger>
          <TabsTrigger value="opportunities">Opportunities</TabsTrigger>
          <TabsTrigger value="toxic">Toxic Analysis</TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-6">
          {/* Score Breakdown */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <BarChart3 className="h-5 w-5" />
                <span>Score Breakdown</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 md:grid-cols-3 gap-6">
                <div className="space-y-2">
                  <div className="flex justify-between items-center">
                    <span className="text-sm font-medium">Quantity</span>
                    <span className="text-sm font-semibold">{backlinksData.score.scoreBreakdown.quantity}/100</span>
                  </div>
                  <Progress value={backlinksData.score.scoreBreakdown.quantity} className="h-2" />
                </div>
                <div className="space-y-2">
                  <div className="flex justify-between items-center">
                    <span className="text-sm font-medium">Quality</span>
                    <span className="text-sm font-semibold">{backlinksData.score.scoreBreakdown.quality}/100</span>
                  </div>
                  <Progress value={backlinksData.score.scoreBreakdown.quality} className="h-2" />
                </div>
                <div className="space-y-2">
                  <div className="flex justify-between items-center">
                    <span className="text-sm font-medium">Diversity</span>
                    <span className="text-sm font-semibold">{backlinksData.score.scoreBreakdown.diversity}/100</span>
                  </div>
                  <Progress value={backlinksData.score.scoreBreakdown.diversity} className="h-2" />
                </div>
                <div className="space-y-2">
                  <div className="flex justify-between items-center">
                    <span className="text-sm font-medium">Authority</span>
                    <span className="text-sm font-semibold">{backlinksData.score.scoreBreakdown.authority}/100</span>
                  </div>
                  <Progress value={backlinksData.score.scoreBreakdown.authority} className="h-2" />
                </div>
                <div className="space-y-2">
                  <div className="flex justify-between items-center">
                    <span className="text-sm font-medium">Naturalness</span>
                    <span className="text-sm font-semibold">{backlinksData.score.scoreBreakdown.naturalness}/100</span>
                  </div>
                  <Progress value={backlinksData.score.scoreBreakdown.naturalness} className="h-2" />
                </div>
                <div className="space-y-2">
                  <div className="flex justify-between items-center">
                    <span className="text-sm font-medium">Toxicity</span>
                    <span className="text-sm font-semibold">{backlinksData.score.scoreBreakdown.toxicity}/100</span>
                  </div>
                  <Progress value={backlinksData.score.scoreBreakdown.toxicity} className="h-2" />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Key Metrics Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <Card>
              <CardContent className="p-6">
                <div className="flex items-center space-x-2">
                  <Link2 className="h-8 w-8 text-blue-500" />
                  <div>
                    <div className="text-2xl font-bold">{backlinksData.profile.totalBacklinks.toLocaleString()}</div>
                    <div className="text-sm text-gray-600">Total Backlinks</div>
                    <div className="text-xs text-green-600 flex items-center mt-1">
                      <TrendingUp className="h-3 w-3 mr-1" />
                      +{backlinksData.profile.newBacklinks30Days} this month
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center space-x-2">
                  <Globe className="h-8 w-8 text-green-500" />
                  <div>
                    <div className="text-2xl font-bold">{backlinksData.profile.totalReferringDomains.toLocaleString()}</div>
                    <div className="text-sm text-gray-600">Referring Domains</div>
                    <div className="text-xs text-blue-600">
                      {backlinksData.profile.totalReferringIPs.toLocaleString()} unique IPs
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center space-x-2">
                  <Award className="h-8 w-8 text-orange-500" />
                  <div>
                    <div className="text-2xl font-bold">{Math.round(backlinksData.profile.averageDomainAuthority)}</div>
                    <div className="text-sm text-gray-600">Avg. Domain Authority</div>
                    <div className="text-xs text-gray-500">
                      TF: {backlinksData.profile.trustFlow} / CF: {backlinksData.profile.citationFlow}
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="p-6">
                <div className="flex items-center space-x-2">
                  <Target className="h-8 w-8 text-purple-500" />
                  <div>
                    <div className="text-2xl font-bold">{Math.round(backlinksData.profile.dofollowPercentage)}%</div>
                    <div className="text-sm text-gray-600">Dofollow Links</div>
                    <div className="text-xs text-gray-500">
                      {Math.round(backlinksData.profile.nofollowPercentage)}% nofollow
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Recent Activity & Quality Distribution */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Activity className="h-5 w-5" />
                  <span>Recent Activity (30 days)</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <CheckCircle className="h-4 w-4 text-green-500" />
                    <span className="text-sm">New Backlinks</span>
                  </div>
                  <span className="font-semibold text-green-600">
                    +{backlinksData.profile.newBacklinks30Days}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <XCircle className="h-4 w-4 text-red-500" />
                    <span className="text-sm">Lost Backlinks</span>
                  </div>
                  <span className="font-semibold text-red-600">
                    -{backlinksData.profile.lostBacklinks30Days}
                  </span>
                </div>
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <AlertTriangle className="h-4 w-4 text-orange-500" />
                    <span className="text-sm">Broken Links</span>
                  </div>
                  <span className="font-semibold text-orange-600">
                    {backlinksData.profile.brokenBacklinks}
                  </span>
                </div>
                <div className="pt-2 border-t">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium">Net Growth</span>
                    <span className={`font-bold ${
                      (backlinksData.profile.newBacklinks30Days - backlinksData.profile.lostBacklinks30Days) > 0 
                        ? 'text-green-600' : 'text-red-600'
                    }`}>
                      {backlinksData.profile.newBacklinks30Days - backlinksData.profile.lostBacklinks30Days > 0 ? '+' : ''}
                      {backlinksData.profile.newBacklinks30Days - backlinksData.profile.lostBacklinks30Days}
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <PieChart className="h-5 w-5" />
                  <span>Quality Distribution</span>
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                {Object.entries(backlinksData.referringDomains.qualityDistribution).map(([quality, count]) => (
                  <div key={quality} className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      <div className={`w-3 h-3 rounded-full ${getQualityColor(quality as BacklinkQuality).replace('text-', 'bg-').replace('bg-', 'bg-').split(' ')[0]}`}></div>
                      <span className="text-sm capitalize">{quality}</span>
                    </div>
                    <div className="text-right">
                      <div className="text-sm font-semibold">{count}</div>
                      <div className="text-xs text-gray-500">
                        {backlinksData.referringDomains.totalCount > 0 
                          ? `${Math.round((count / backlinksData.referringDomains.totalCount) * 100)}%`
                          : '0%'
                        }
                      </div>
                    </div>
                  </div>
                ))}
              </CardContent>
            </Card>
          </div>

          {/* Top Recommendations */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Zap className="h-5 w-5" />
                <span>Immediate Action Items</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {backlinksData.recommendations.immediate.slice(0, 3).map((rec, index) => (
                  <div key={rec.id} className="flex items-start space-x-3 p-3 border border-gray-200 rounded-lg">
                    <div className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold text-white ${
                      rec.priority === 'high' ? 'bg-red-500' :
                      rec.priority === 'medium' ? 'bg-orange-500' : 'bg-blue-500'
                    }`}>
                      {index + 1}
                    </div>
                    <div className="flex-1">
                      <div className="font-medium">{rec.title}</div>
                      <div className="text-sm text-gray-600 mt-1">{rec.description}</div>
                      <div className="flex items-center space-x-4 mt-2 text-xs">
                        <span className="flex items-center space-x-1">
                          <Clock className="h-3 w-3" />
                          <span>{rec.timeline}</span>
                        </span>
                        <span className="flex items-center space-x-1">
                          <Users className="h-3 w-3" />
                          <span>{rec.effort} effort</span>
                        </span>
                        <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                          rec.priority === 'high' ? 'bg-red-100 text-red-800' :
                          rec.priority === 'medium' ? 'bg-orange-100 text-orange-800' : 'bg-blue-100 text-blue-800'
                        }`}>
                          {rec.priority} priority
                        </span>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Backlinks Tab */}
        <TabsContent value="backlinks" className="space-y-6">
          {/* Filters and Search */}
          <Card>
            <CardContent className="p-4">
              <div className="flex flex-col lg:flex-row lg:items-center space-y-4 lg:space-y-0 lg:space-x-4">
                <div className="flex-1">
                  <div className="relative">
                    <Search className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
                    <input
                      type="text"
                      placeholder="Search domains, URLs, or anchor text..."
                      value={searchTerm}
                      onChange={(e) => setSearchTerm(e.target.value)}
                      className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    />
                  </div>
                </div>
                
                <div className="flex items-center space-x-3">
                  <select
                    value={filterStatus}
                    onChange={(e) => setFilterStatus(e.target.value as BacklinkStatus | 'all')}
                    className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="all">All Status</option>
                    <option value={BacklinkStatus.ACTIVE}>Active</option>
                    <option value={BacklinkStatus.NEW}>New</option>
                    <option value={BacklinkStatus.LOST}>Lost</option>
                    <option value={BacklinkStatus.TOXIC}>Toxic</option>
                    <option value={BacklinkStatus.BROKEN}>Broken</option>
                  </select>

                  <select
                    value={filterQuality}
                    onChange={(e) => setFilterQuality(e.target.value as BacklinkQuality | 'all')}
                    className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="all">All Quality</option>
                    <option value={BacklinkQuality.EXCELLENT}>Excellent</option>
                    <option value={BacklinkQuality.GOOD}>Good</option>
                    <option value={BacklinkQuality.AVERAGE}>Average</option>
                    <option value={BacklinkQuality.POOR}>Poor</option>
                    <option value={BacklinkQuality.TOXIC}>Toxic</option>
                  </select>

                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setShowAdvancedFilters(!showAdvancedFilters)}
                  >
                    <Filter className="h-4 w-4 mr-2" />
                    {showAdvancedFilters ? 'Hide' : 'Show'} Filters
                  </Button>
                </div>
              </div>

              {showAdvancedFilters && (
                <div className="mt-4 pt-4 border-t border-gray-200">
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Link Type</label>
                      <select className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
                        <option value="all">All Types</option>
                        <option value={LinkType.DOFOLLOW}>Dofollow</option>
                        <option value={LinkType.NOFOLLOW}>Nofollow</option>
                        <option value={LinkType.SPONSORED}>Sponsored</option>
                        <option value={LinkType.UGC}>UGC</option>
                      </select>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Domain Authority</label>
                      <select className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
                        <option value="all">All DA</option>
                        <option value="80+">80+ (Excellent)</option>
                        <option value="60-79">60-79 (Good)</option>
                        <option value="40-59">40-59 (Average)</option>
                        <option value="0-39">0-39 (Low)</option>
                      </select>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Date Added</label>
                      <select className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
                        <option value="all">All Time</option>
                        <option value="7d">Last 7 days</option>
                        <option value="30d">Last 30 days</option>
                        <option value="90d">Last 90 days</option>
                        <option value="1y">Last year</option>
                      </select>
                    </div>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Results Summary */}
          <div className="flex items-center justify-between">
            <div className="text-sm text-gray-600">
              Showing {((currentPage - 1) * itemsPerPage) + 1}-{Math.min(currentPage * itemsPerPage, sortedBacklinks.length)} of {sortedBacklinks.length} backlinks
              {filteredBacklinks.length !== backlinksData.backlinks.list.length && (
                <span className="ml-2 text-blue-600">({backlinksData.backlinks.list.length} total)</span>
              )}
            </div>
            <div className="flex items-center space-x-2">
              <span className="text-sm text-gray-600">Sort by:</span>
              <select
                value={sortBy}
                onChange={(e) => handleSort(e.target.value)}
                className="px-3 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
              >
                <option value="domainAuthority">Domain Authority</option>
                <option value="firstSeen">Date Added</option>
                <option value="sourceDomain">Source Domain</option>
                <option value="toxicityScore">Toxicity Score</option>
              </select>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')}
              >
                {sortDirection === 'asc' ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
              </Button>
            </div>
          </div>

          {/* Backlinks List */}
          <Card>
            <CardContent className="p-0">
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="bg-gray-50 border-b">
                    <tr>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Source</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Anchor Text</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Quality</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">DA/PA</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">First Seen</th>
                      <th className="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-200">
                    {paginatedBacklinks.map((backlink, index) => (
                      <tr key={backlink.id} className="hover:bg-gray-50">
                        <td className="px-4 py-3">
                          <div className="flex items-start space-x-2">
                            <div className="flex-1 min-w-0">
                              <div className="font-medium text-sm truncate">{backlink.sourceDomain}</div>
                              <div className="text-xs text-gray-500 truncate" title={backlink.sourceUrl}>
                                {backlink.sourceUrl.length > 50 ? backlink.sourceUrl.substring(0, 50) + '...' : backlink.sourceUrl}
                              </div>
                              {backlink.pageTitle && (
                                <div className="text-xs text-gray-400 truncate mt-1" title={backlink.pageTitle}>
                                  {backlink.pageTitle.length > 40 ? backlink.pageTitle.substring(0, 40) + '...' : backlink.pageTitle}
                                </div>
                              )}
                            </div>
                            <div className="flex items-center space-x-1">
                              <MapPin className="h-3 w-3 text-gray-400" />
                              <span className="text-xs text-gray-500">{backlink.country}</span>
                            </div>
                          </div>
                        </td>
                        <td className="px-4 py-3">
                          <div className="max-w-48">
                            <div className="text-sm font-medium truncate" title={backlink.anchorText}>
                              {backlink.anchorText || '[No anchor text]'}
                            </div>
                            <div className="flex items-center space-x-1 mt-1">
                              <span className={`px-1.5 py-0.5 text-xs rounded-full font-medium ${
                                backlink.anchorType === AnchorType.EXACT_MATCH ? 'bg-red-100 text-red-800' :
                                backlink.anchorType === AnchorType.PARTIAL_MATCH ? 'bg-orange-100 text-orange-800' :
                                backlink.anchorType === AnchorType.BRANDED ? 'bg-green-100 text-green-800' :
                                'bg-gray-100 text-gray-800'
                              }`}>
                                {backlink.anchorType.replace('_', ' ')}
                              </span>
                            </div>
                          </div>
                        </td>
                        <td className="px-4 py-3">
                          <div className="flex items-center space-x-2">
                            {getStatusIcon(backlink.status)}
                            <span className="text-sm capitalize">{backlink.status.replace('_', ' ')}</span>
                          </div>
                        </td>
                        <td className="px-4 py-3">
                          <span className={`px-2 py-1 text-xs rounded-full font-medium ${getQualityColor(backlink.quality)}`}>
                            {backlink.quality}
                          </span>
                          {backlink.toxicityScore > 50 && (
                            <div className="text-xs text-red-600 mt-1 flex items-center space-x-1">
                              <AlertTriangle className="h-3 w-3" />
                              <span>Toxicity: {backlink.toxicityScore}</span>
                            </div>
                          )}
                        </td>
                        <td className="px-4 py-3">
                          <div className="text-sm">
                            <div className="font-semibold">{backlink.domainAuthority.domainAuthority}</div>
                            <div className="text-xs text-gray-500">{backlink.domainAuthority.pageAuthority}</div>
                          </div>
                        </td>
                        <td className="px-4 py-3">
                          <div className="flex items-center space-x-1">
                            {backlink.linkType === LinkType.DOFOLLOW && (
                              <div className="w-2 h-2 rounded-full bg-green-500" title="Dofollow"></div>
                            )}
                            {backlink.linkType === LinkType.NOFOLLOW && (
                              <div className="w-2 h-2 rounded-full bg-gray-400" title="Nofollow"></div>
                            )}
                            {backlink.linkType === LinkType.SPONSORED && (
                              <div className="w-2 h-2 rounded-full bg-orange-500" title="Sponsored"></div>
                            )}
                            <span className="text-xs text-gray-600">{backlink.linkType}</span>
                          </div>
                        </td>
                        <td className="px-4 py-3 text-sm text-gray-600">
                          {new Date(backlink.firstSeen).toLocaleDateString()}
                        </td>
                        <td className="px-4 py-3">
                          <div className="flex items-center space-x-2">
                            <Button
                              size="sm"
                              variant="ghost"
                              onClick={() => window.open(backlink.sourceUrl, '_blank')}
                            >
                              <ExternalLink className="h-4 w-4" />
                            </Button>
                            <Button
                              size="sm"
                              variant="ghost"
                              onClick={() => toggleSection(`backlink-${index}`)}
                            >
                              {expandedSections[`backlink-${index}`] ? 
                                <EyeOff className="h-4 w-4" /> : 
                                <Eye className="h-4 w-4" />
                              }
                            </Button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {/* Pagination */}
              {totalPages > 1 && (
                <div className="px-4 py-3 border-t border-gray-200 flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                      disabled={currentPage === 1}
                    >
                      Previous
                    </Button>
                    <span className="text-sm text-gray-600">
                      Page {currentPage} of {totalPages}
                    </span>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                      disabled={currentPage === totalPages}
                    >
                      Next
                    </Button>
                  </div>
                  <div className="text-sm text-gray-600">
                    {sortedBacklinks.length} results
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Referring Domains Tab */}
        <TabsContent value="domains" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2">
              <Card>
                <CardHeader>
                  <CardTitle>Top Referring Domains</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {backlinksData.referringDomains.topDomains.slice(0, 10).map((domain, index) => (
                      <div key={domain.domain} className="flex items-center space-x-4 p-3 border border-gray-200 rounded-lg">
                        <div className="w-8 h-8 rounded-full bg-blue-500 text-white flex items-center justify-center text-sm font-bold">
                          {index + 1}
                        </div>
                        <div className="flex-1">
                          <div className="font-medium">{domain.domain}</div>
                          <div className="text-sm text-gray-600">
                            {domain.backlinksCount} backlinks • DA {domain.domainAuthority.domainAuthority}
                          </div>
                        </div>
                        <div className="text-right">
                          <span className={`px-2 py-1 text-xs rounded-full font-medium ${getQualityColor(domain.quality)}`}>
                            {domain.quality}
                          </span>
                          <div className="text-xs text-gray-500 mt-1">
                            {new Date(domain.firstSeen).toLocaleDateString()}
                          </div>
                        </div>
                        <Button size="sm" variant="ghost" onClick={() => window.open(`https://${domain.domain}`, '_blank')}>
                          <ExternalLink className="h-4 w-4" />
                        </Button>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>

            <div className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle>Domain Categories</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  {Object.entries(backlinksData.referringDomains.categories)
                    .filter(([, count]) => count > 0)
                    .sort(([, a], [, b]) => b - a)
                    .slice(0, 8)
                    .map(([category, count]) => (
                      <div key={category} className="flex items-center justify-between">
                        <span className="text-sm capitalize">{category.replace('_', ' ')}</span>
                        <div className="text-right">
                          <div className="text-sm font-semibold">{count}</div>
                          <div className="text-xs text-gray-500">
                            {backlinksData.referringDomains.totalCount > 0 
                              ? `${Math.round((count / backlinksData.referringDomains.totalCount) * 100)}%`
                              : '0%'
                            }
                          </div>
                        </div>
                      </div>
                    ))}
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Geographic Distribution</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                  {Object.entries(backlinksData.referringDomains.geographicDistribution)
                    .sort(([, a], [, b]) => b - a)
                    .slice(0, 8)
                    .map(([country, count]) => (
                      <div key={country} className="flex items-center justify-between">
                        <div className="flex items-center space-x-2">
                          <Flag className="h-4 w-4 text-gray-400" />
                          <span className="text-sm">{country}</span>
                        </div>
                        <div className="text-right">
                          <div className="text-sm font-semibold">{count}</div>
                          <div className="text-xs text-gray-500">
                            {backlinksData.referringDomains.totalCount > 0 
                              ? `${Math.round((count / backlinksData.referringDomains.totalCount) * 100)}%`
                              : '0%'
                            }
                          </div>
                        </div>
                      </div>
                    ))}
                </CardContent>
              </Card>
            </div>
          </div>
        </TabsContent>

        {/* Anchor Texts Tab */}
        <TabsContent value="anchors" className="space-y-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-2">
              <Card>
                <CardHeader>
                  <CardTitle>Anchor Text Distribution</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    {backlinksData.anchorTexts.distribution.slice(0, 15).map((anchor, index) => (
                      <div key={`${anchor.anchorText}-${index}`} className="flex items-center space-x-4 p-3 border border-gray-200 rounded-lg">
                        <div className="flex-1">
                          <div className="flex items-center space-x-2">
                            <div className="font-medium">{anchor.anchorText}</div>
                            <span className={`px-1.5 py-0.5 text-xs rounded-full font-medium ${
                              anchor.type === AnchorType.EXACT_MATCH ? 'bg-red-100 text-red-800' :
                              anchor.type === AnchorType.PARTIAL_MATCH ? 'bg-orange-100 text-orange-800' :
                              anchor.type === AnchorType.BRANDED ? 'bg-green-100 text-green-800' :
                              'bg-gray-100 text-gray-800'
                            }`}>
                              {anchor.type.replace('_', ' ')}
                            </span>
                            {anchor.isOverOptimized && (
                              <AlertTriangle className="h-4 w-4 text-red-500" title="Over-optimized" />
                            )}
                          </div>
                          <div className="text-sm text-gray-600 mt-1">
                            {anchor.count} links • Avg DA: {anchor.averageDA} • {anchor.topDomains.length} domains
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="text-lg font-bold">{anchor.percentage}%</div>
                          {anchor.isOverOptimized && (
                            <div className="text-xs text-red-600">High risk</div>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>

            <div className="space-y-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center space-x-2">
                    <Shield className="h-5 w-5" />
                    <span>Naturalness Score</span>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-center mb-4">
                    <div className={`text-3xl font-bold ${backlinksData.anchorTexts.naturalness.isNatural ? 'text-green-600' : 'text-red-600'}`}>
                      {backlinksData.anchorTexts.naturalness.score}/100
                    </div>
                    <div className={`text-sm font-medium ${backlinksData.anchorTexts.naturalness.isNatural ? 'text-green-600' : 'text-red-600'}`}>
                      {backlinksData.anchorTexts.naturalness.isNatural ? 'Natural' : 'Unnatural'}
                    </div>
                  </div>
                  <Progress value={backlinksData.anchorTexts.naturalness.score} className="h-3 mb-4" />
                  
                  {backlinksData.anchorTexts.naturalness.overOptimizedAnchors.length > 0 && (
                    <div className="space-y-2">
                      <div className="text-sm font-medium text-red-600">Over-optimized anchors:</div>
                      <div className="space-y-1">
                        {backlinksData.anchorTexts.naturalness.overOptimizedAnchors.slice(0, 5).map((anchor, index) => (
                          <div key={index} className="text-sm text-red-600 bg-red-50 px-2 py-1 rounded">
                            {anchor}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}

                  <div className="mt-4 space-y-2">
                    <div className="text-sm font-medium">Recommendations:</div>
                    <div className="space-y-1">
                      {backlinksData.anchorTexts.naturalness.recommendations.slice(0, 3).map((rec, index) => (
                        <div key={index} className="text-xs text-gray-600 p-2 bg-blue-50 rounded">
                          {rec}
                        </div>
                      ))}
                    </div>
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Anchor Diversity</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-center mb-4">
                    <div className="text-2xl font-bold text-blue-600">{backlinksData.anchorTexts.diversity.uniqueAnchors}</div>
                    <div className="text-sm text-gray-600">Unique anchors</div>
                  </div>
                  <div className="text-center mb-4">
                    <div className="text-lg font-semibold">{Math.round(backlinksData.anchorTexts.diversity.diversityScore)}/100</div>
                    <div className="text-sm text-gray-600">Diversity score</div>
                  </div>
                  <Progress value={backlinksData.anchorTexts.diversity.diversityScore} className="h-2" />
                </CardContent>
              </Card>
            </div>
          </div>
        </TabsContent>

        {/* Competitors Tab */}
        <TabsContent value="competitors" className="space-y-6">
          {backlinksData.competitors.analysis.length > 0 ? (
            <div className="space-y-6">
              {/* Benchmarking Overview */}
              <Card>
                <CardHeader>
                  <CardTitle>Competitive Position</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                    <div className="text-center">
                      <div className="text-2xl font-bold text-blue-600">#{backlinksData.competitors.benchmarking.position}</div>
                      <div className="text-sm text-gray-600">Your Position</div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold">{backlinksData.competitors.benchmarking.totalCompetitors}</div>
                      <div className="text-sm text-gray-600">Total Competitors</div>
                    </div>
                    <div className="text-center">
                      <div className={`text-2xl font-bold ${backlinksData.competitors.benchmarking.aboveAverage ? 'text-green-600' : 'text-red-600'}`}>
                        {backlinksData.competitors.benchmarking.aboveAverage ? 'Above' : 'Below'}
                      </div>
                      <div className="text-sm text-gray-600">Industry Average</div>
                    </div>
                    <div className="text-center">
                      <div className="text-2xl font-bold text-orange-600">{backlinksData.competitors.gapAnalysis.totalOpportunities}</div>
                      <div className="text-sm text-gray-600">Gap Opportunities</div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Competitor Analysis */}
              <Card>
                <CardHeader>
                  <CardTitle>Competitor Comparison</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-6">
                    {backlinksData.competitors.analysis.map((competitor, index) => (
                      <div key={competitor.competitorDomain} className="p-4 border border-gray-200 rounded-lg">
                        <div className="flex items-center justify-between mb-4">
                          <div>
                            <div className="font-semibold text-lg">{competitor.competitor}</div>
                            <div className="text-sm text-gray-600">{competitor.competitorDomain}</div>
                          </div>
                          <div className="text-right">
                            <div className="text-2xl font-bold">{competitor.domainAuthority.domainAuthority}</div>
                            <div className="text-sm text-gray-600">Domain Authority</div>
                          </div>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                          <div>
                            <div className="text-sm text-gray-600">Total Backlinks</div>
                            <div className="font-semibold">{competitor.totalBacklinks.toLocaleString()}</div>
                          </div>
                          <div>
                            <div className="text-sm text-gray-600">Referring Domains</div>
                            <div className="font-semibold">{competitor.referringDomains.toLocaleString()}</div>
                          </div>
                          <div>
                            <div className="text-sm text-gray-600">Unique Opportunities</div>
                            <div className="font-semibold">{competitor.gap.uniqueDomains}</div>
                          </div>
                        </div>

                        <div className="flex items-center space-x-4">
                          <Button size="sm" variant="outline">
                            View Gap Analysis
                          </Button>
                          <Button size="sm" variant="outline">
                            <ExternalLink className="h-4 w-4 mr-2" />
                            Visit Site
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          ) : (
            <Card>
              <CardContent className="text-center py-8">
                <div className="text-gray-500">No competitor data available</div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        {/* Opportunities Tab */}
        <TabsContent value="opportunities" className="space-y-6">
          {backlinksData.opportunities.list.length > 0 ? (
            <div className="space-y-6">
              {/* Opportunities Summary */}
              <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                <Card>
                  <CardContent className="p-6 text-center">
                    <div className="text-2xl font-bold text-blue-600">{backlinksData.opportunities.totalCount}</div>
                    <div className="text-sm text-gray-600">Total Opportunities</div>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="p-6 text-center">
                    <div className="text-2xl font-bold text-green-600">{backlinksData.opportunities.priorityDistribution.high}</div>
                    <div className="text-sm text-gray-600">High Priority</div>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="p-6 text-center">
                    <div className="text-2xl font-bold text-orange-600">{backlinksData.opportunities.priorityDistribution.medium}</div>
                    <div className="text-sm text-gray-600">Medium Priority</div>
                  </CardContent>
                </Card>
                <Card>
                  <CardContent className="p-6 text-center">
                    <div className="text-2xl font-bold text-purple-600">
                      ${Math.round(backlinksData.opportunities.estimatedTotalValue / 1000)}K
                    </div>
                    <div className="text-sm text-gray-600">Est. Traffic Value</div>
                  </CardContent>
                </Card>
              </div>

              {/* Opportunities List */}
              <Card>
                <CardHeader>
                  <CardTitle>Link Building Opportunities</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {backlinksData.opportunities.list.slice(0, 20).map((opportunity, index) => (
                      <div key={opportunity.id} className="p-4 border border-gray-200 rounded-lg hover:bg-gray-50">
                        <div className="flex items-start space-x-4">
                          <div className={`w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold text-white ${
                            opportunity.priority === 'high' ? 'bg-red-500' :
                            opportunity.priority === 'medium' ? 'bg-orange-500' : 'bg-blue-500'
                          }`}>
                            {index + 1}
                          </div>
                          
                          <div className="flex-1">
                            <div className="flex items-center justify-between mb-2">
                              <div>
                                <div className="font-semibold">{opportunity.domain}</div>
                                <div className="text-sm text-gray-600">{opportunity.title}</div>
                              </div>
                              <div className="text-right">
                                <div className="font-semibold">DA {opportunity.domainAuthority.domainAuthority}</div>
                                <div className="text-sm text-gray-600">
                                  ${Math.round(opportunity.trafficValue / 100) * 100} value
                                </div>
                              </div>
                            </div>

                            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-3">
                              <div>
                                <div className="text-xs text-gray-500">Difficulty</div>
                                <Progress value={opportunity.estimatedDifficulty} className="h-2 mt-1" />
                                <div className="text-xs text-gray-600 mt-1">{opportunity.estimatedDifficulty}%</div>
                              </div>
                              <div>
                                <div className="text-xs text-gray-500">Relevance</div>
                                <Progress value={opportunity.relevanceScore} className="h-2 mt-1" />
                                <div className="text-xs text-gray-600 mt-1">{opportunity.relevanceScore}%</div>
                              </div>
                              <div>
                                <div className="text-xs text-gray-500">Category</div>
                                <div className="text-xs font-medium mt-1 capitalize">
                                  {opportunity.category.replace('_', ' ')}
                                </div>
                              </div>
                            </div>

                            <div className="flex items-center justify-between">
                              <div className="flex items-center space-x-2">
                                <span className={`px-2 py-1 text-xs rounded-full font-medium ${
                                  opportunity.priority === 'high' ? 'bg-red-100 text-red-800' :
                                  opportunity.priority === 'medium' ? 'bg-orange-100 text-orange-800' : 'bg-blue-100 text-blue-800'
                                }`}>
                                  {opportunity.priority} priority
                                </span>
                                <span className="text-xs text-gray-600">{opportunity.suggestedApproach}</span>
                              </div>
                              <div className="flex items-center space-x-2">
                                {opportunity.contactInfo.email && (
                                  <Button size="sm" variant="outline">
                                    Contact
                                  </Button>
                                )}
                                <Button size="sm" variant="outline" onClick={() => window.open(opportunity.url, '_blank')}>
                                  <ExternalLink className="h-4 w-4" />
                                </Button>
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          ) : (
            <Card>
              <CardContent className="text-center py-8">
                <div className="text-gray-500">No link building opportunities identified</div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        {/* Toxic Analysis Tab */}
        <TabsContent value="toxic" className="space-y-6">
          {/* Toxic Overview */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <AlertTriangle className="h-5 w-5 text-red-500" />
                <span>Toxic Link Analysis</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                <div className="text-center">
                  <div className="text-2xl font-bold text-red-600">{backlinksData.toxicAnalysis.toxicBacklinks.length}</div>
                  <div className="text-sm text-gray-600">Toxic Backlinks</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-orange-600">{backlinksData.toxicAnalysis.toxicDomains.length}</div>
                  <div className="text-sm text-gray-600">Toxic Domains</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-purple-600">{backlinksData.toxicAnalysis.overallToxicity}%</div>
                  <div className="text-sm text-gray-600">Overall Toxicity</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-blue-600">{backlinksData.toxicAnalysis.recommendations.recommendedActions.length}</div>
                  <div className="text-sm text-gray-600">Actions Needed</div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Risk Assessment */}
          <Card>
            <CardHeader>
              <CardTitle>Risk Assessment</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <div className="mb-4">
                    <div className="flex justify-between items-center mb-2">
                      <span className="text-sm font-medium">Current Risk Level</span>
                      <span className="text-sm font-semibold">{backlinksData.toxicAnalysis.recommendations.riskAssessment.currentRisk}%</span>
                    </div>
                    <Progress value={backlinksData.toxicAnalysis.recommendations.riskAssessment.currentRisk} className="h-3" />
                  </div>
                  <div className="mb-4">
                    <div className="flex justify-between items-center mb-2">
                      <span className="text-sm font-medium">Risk After Disavow</span>
                      <span className="text-sm font-semibold">{backlinksData.toxicAnalysis.recommendations.riskAssessment.riskAfterDisavow}%</span>
                    </div>
                    <Progress value={backlinksData.toxicAnalysis.recommendations.riskAssessment.riskAfterDisavow} className="h-3" />
                  </div>
                </div>
                <div className="space-y-3">
                  <div>
                    <span className="text-sm font-medium">Estimated Recovery Time:</span>
                    <span className="ml-2 text-sm">{backlinksData.toxicAnalysis.recommendations.riskAssessment.estimatedRecoveryTime}</span>
                  </div>
                  <div>
                    <span className="text-sm font-medium">Confidence Level:</span>
                    <span className="ml-2 text-sm">{backlinksData.toxicAnalysis.recommendations.riskAssessment.confidenceLevel}%</span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Risk Factors */}
          <Card>
            <CardHeader>
              <CardTitle>Risk Factors</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {backlinksData.toxicAnalysis.riskFactors.map((factor, index) => (
                  <div key={index} className="p-4 border border-gray-200 rounded-lg">
                    <div className="flex items-start space-x-3">
                      <div className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold text-white ${
                        factor.severity === 'high' ? 'bg-red-500' :
                        factor.severity === 'medium' ? 'bg-orange-500' : 'bg-yellow-500'
                      }`}>
                        !
                      </div>
                      <div className="flex-1">
                        <div className="font-medium">{factor.factor}</div>
                        <div className="text-sm text-gray-600 mt-1">{factor.description}</div>
                        <div className="text-xs text-gray-500 mt-2">
                          Impact: {factor.impact} • Affected links: {factor.affectedLinks}
                        </div>
                      </div>
                      <span className={`px-2 py-1 text-xs rounded-full font-medium ${
                        factor.severity === 'high' ? 'bg-red-100 text-red-800' :
                        factor.severity === 'medium' ? 'bg-orange-100 text-orange-800' : 'bg-yellow-100 text-yellow-800'
                      }`}>
                        {factor.severity}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Recommended Actions */}
          <Card>
            <CardHeader>
              <CardTitle>Recommended Actions</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {backlinksData.toxicAnalysis.recommendations.recommendedActions.slice(0, 10).map((action, index) => (
                  <div key={index} className="p-4 border border-gray-200 rounded-lg">
                    <div className="flex items-start space-x-3">
                      <div className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold text-white ${
                        action.priority === 'high' ? 'bg-red-500' :
                        action.priority === 'medium' ? 'bg-orange-500' : 'bg-blue-500'
                      }`}>
                        {index + 1}
                      </div>
                      <div className="flex-1">
                        <div className="font-medium">{action.backlink.sourceDomain}</div>
                        <div className="text-sm text-gray-600 mt-1">{action.reason}</div>
                        <div className="text-xs text-gray-500 mt-2">{action.potentialImpact}</div>
                      </div>
                      <div className="text-right">
                        <span className={`px-2 py-1 text-xs rounded-full font-medium ${
                          action.action === 'disavow' ? 'bg-red-100 text-red-800' :
                          action.action === 'contact_removal' ? 'bg-orange-100 text-orange-800' : 'bg-blue-100 text-blue-800'
                        }`}>
                          {action.action.replace('_', ' ')}
                        </span>
                        <div className="text-xs text-gray-500 mt-1">{action.priority} priority</div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Disavow File */}
          <Card>
            <CardHeader>
              <CardTitle>Disavow File</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                <div className="text-center">
                  <div className="text-xl font-bold">{backlinksData.toxicAnalysis.recommendations.disavowFile.domainsCount}</div>
                  <div className="text-sm text-gray-600">Domains to disavow</div>
                </div>
                <div className="text-center">
                  <div className="text-xl font-bold">{backlinksData.toxicAnalysis.recommendations.disavowFile.urlsCount}</div>
                  <div className="text-sm text-gray-600">URLs to disavow</div>
                </div>
                <div className="text-center">
                  <div className="text-sm text-gray-600">Last updated</div>
                  <div className="text-sm font-medium">
                    {new Date(backlinksData.toxicAnalysis.recommendations.disavowFile.lastUpdated).toLocaleDateString()}
                  </div>
                </div>
              </div>
              
              <div className="flex items-center space-x-2">
                <Button>
                  <Download className="h-4 w-4 mr-2" />
                  Download Disavow File
                </Button>
                <Button variant="outline">
                  Preview File
                </Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}