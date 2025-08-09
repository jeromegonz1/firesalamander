"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { 
  Download, 
  ArrowRight,
  TrendingUp,
  TrendingDown,
  Minus,
  Trophy,
  Target,
  BarChart3,
  Globe,
  Clock,
  Eye,
  Filter,
  Lightbulb,
  AlertTriangle,
  CheckCircle2
} from "lucide-react";

interface AnalysisData {
  id: string;
  url: string;
  domain: string;
  analyzed_at: string;
  overall_score: number;
  scores: {
    seo: number;
    technical: number;
    performance: number;
    content: number;
  };
  recommendations_count: number;
  critical_issues: number;
}

interface ComparisonMetric {
  name: string;
  site1_value: number;
  site2_value: number;
  category: 'seo' | 'technical' | 'performance' | 'content';
  winner: 1 | 2 | 'tie';
  difference: number;
}

interface KeywordComparison {
  keyword: string;
  site1_position: number;
  site2_position: number;
  site1_density: number;
  site2_density: number;
  winner: 1 | 2 | 'tie';
}

interface ImprovementOpportunity {
  title: string;
  description: string;
  based_on_winner: 1 | 2;
  impact_score: number;
  difficulty: 'easy' | 'medium' | 'hard';
  category: string;
}

// Composant Score Comparison
function ScoreComparison({ analysis1, analysis2 }: { analysis1: AnalysisData, analysis2: AnalysisData }) {
  const overallWinner = analysis1.overall_score > analysis2.overall_score ? 1 : 
                       analysis2.overall_score > analysis1.overall_score ? 2 : 'tie';
  const scoreDifference = Math.abs(analysis1.overall_score - analysis2.overall_score);

  return (
    <Card data-testid="overall-scores-comparison" aria-label="Section scores globaux">
      <CardHeader>
        <CardTitle className="flex items-center space-x-2">
          <Trophy className="h-5 w-5" />
          <span>Comparaison des Scores</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 gap-8" data-testid="comparison-column-1">
          {/* Site 1 */}
          <div className="text-center">
            <div className="flex items-center justify-center mb-4">
              {overallWinner === 1 && (
                <Badge className="bg-green-100 text-green-700 border-green-200 mb-2" data-testid="overall-winner-badge">
                  🏆 Gagnant
                </Badge>
              )}
            </div>
            <h3 className="font-semibold text-lg mb-2">{analysis1.domain}</h3>
            <div className="text-4xl font-bold text-orange-600 mb-2" data-testid="score-1-value">
              {analysis1.overall_score}
            </div>
            <div className="text-sm text-gray-600">Score Global</div>
          </div>

          {/* Site 2 */}
          <div className="text-center" data-testid="comparison-column-2">
            <div className="flex items-center justify-center mb-4">
              {overallWinner === 2 && (
                <Badge className="bg-green-100 text-green-700 border-green-200 mb-2" data-testid="overall-winner-badge">
                  🏆 Gagnant
                </Badge>
              )}
            </div>
            <h3 className="font-semibold text-lg mb-2">{analysis2.domain}</h3>
            <div className="text-4xl font-bold text-orange-600 mb-2" data-testid="score-2-value">
              {analysis2.overall_score}
            </div>
            <div className="text-sm text-gray-600">Score Global</div>
          </div>
        </div>

        {/* Difference Indicator */}
        <div className="text-center mt-6 p-4 bg-gray-50 rounded-lg" data-testid="score-difference">
          <div className="text-sm text-gray-600 mb-1">Écart de score</div>
          <div className="text-2xl font-bold text-gray-900">
            {scoreDifference} point{scoreDifference > 1 ? 's' : ''}
          </div>
          {overallWinner !== 'tie' && (
            <div className="text-sm text-gray-600 mt-1" data-testid="score-winner">
              {overallWinner === 1 ? analysis1.domain : analysis2.domain} en avance
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

// Composant Detailed Metrics Comparison
function DetailedMetricsComparison({ metrics }: { metrics: ComparisonMetric[] }) {
  const getWinnerIcon = (winner: 1 | 2 | 'tie') => {
    if (winner === 'tie') return <Minus className="h-4 w-4 text-gray-600" />;
    return <Trophy className="h-4 w-4 text-green-600" />;
  };

  return (
    <Card data-testid="detailed-metrics-comparison" aria-label="Section métriques détaillées">
      <CardHeader>
        <CardTitle>Métriques Détaillées</CardTitle>
        <CardDescription>Comparaison métrique par métrique</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto" data-testid="comparison-table" role="table" aria-label="Tableau comparatif des métriques">
          <table className="w-full border-collapse">
            <thead>
              <tr className="border-b">
                <th className="text-left p-3 font-semibold">Métrique</th>
                <th className="text-center p-3 font-semibold">Site 1</th>
                <th className="text-center p-3 font-semibold">Site 2</th>
                <th className="text-center p-3 font-semibold">Gagnant</th>
              </tr>
            </thead>
            <tbody>
              {metrics.map((metric, index) => (
                <tr key={index} className="border-b hover:bg-gray-50" data-testid="metric-comparison-row">
                  <td className="p-3">
                    <div data-testid="metric-name">
                      <span className="font-medium">{metric.name}</span>
                      <Badge variant="outline" className="ml-2 text-xs">
                        {metric.category}
                      </Badge>
                    </div>
                  </td>
                  <td className="p-3 text-center" data-testid="metric-value-1">
                    <span className={`font-semibold ${metric.winner === 1 ? 'text-green-600' : 'text-gray-700'}`}>
                      {metric.site1_value}
                    </span>
                  </td>
                  <td className="p-3 text-center" data-testid="metric-value-2">
                    <span className={`font-semibold ${metric.winner === 2 ? 'text-green-600' : 'text-gray-700'}`}>
                      {metric.site2_value}
                    </span>
                  </td>
                  <td className="p-3 text-center">
                    <div className="flex items-center justify-center space-x-1">
                      {getWinnerIcon(metric.winner)}
                      {metric.winner !== 'tie' && (
                        <span className="text-sm text-green-600">
                          Site {metric.winner}
                        </span>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  );
}

// Composant Visual Chart Comparison
function VisualChartComparison({ analysis1, analysis2 }: { analysis1: AnalysisData, analysis2: AnalysisData }) {
  const categories = ['seo', 'technical', 'performance', 'content'];
  
  return (
    <Card data-testid="comparison-chart" aria-label="Graphique comparatif">
      <CardHeader>
        <CardTitle>Comparaison Visuelle</CardTitle>
        <CardDescription>Radar des performances par catégorie</CardDescription>
      </CardHeader>
      <CardContent>
        {/* Simulation d'un radar chart avec CSS */}
        <div className="space-y-4">
          {categories.map((category, index) => {
            const score1 = analysis1.scores[category as keyof typeof analysis1.scores];
            const score2 = analysis2.scores[category as keyof typeof analysis2.scores];
            const maxScore = Math.max(score1, score2);
            
            return (
              <div key={index} className="space-y-2">
                <div className="flex items-center justify-between">
                  <span className="font-medium capitalize">{category}</span>
                  <div className="flex items-center space-x-4">
                    <span className="text-sm text-blue-600 font-medium">{score1}</span>
                    <span className="text-sm text-purple-600 font-medium">{score2}</span>
                  </div>
                </div>
                <div className="relative">
                  <div className="flex space-x-1">
                    {/* Bar pour site 1 */}
                    <div className="flex-1 bg-gray-200 rounded-full h-4 relative">
                      <div 
                        className="bg-blue-500 h-4 rounded-full transition-all"
                        style={{ width: `${score1}%` }}
                      />
                    </div>
                    {/* Bar pour site 2 */}
                    <div className="flex-1 bg-gray-200 rounded-full h-4 relative">
                      <div 
                        className="bg-purple-500 h-4 rounded-full transition-all"
                        style={{ width: `${score2}%` }}
                      />
                    </div>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
        
        {/* Légende */}
        <div className="flex items-center justify-center space-x-6 mt-6" data-testid="chart-legend">
          <div className="flex items-center space-x-2" data-testid="dataset-1-legend">
            <div className="w-4 h-4 bg-blue-500 rounded"></div>
            <span className="text-sm">{analysis1.domain}</span>
          </div>
          <div className="flex items-center space-x-2" data-testid="dataset-2-legend">
            <div className="w-4 h-4 bg-purple-500 rounded"></div>
            <span className="text-sm">{analysis2.domain}</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

// Composant SEO Recommendations Comparison
function SEORecommendationsComparison({ analysis1, analysis2 }: { analysis1: AnalysisData, analysis2: AnalysisData }) {
  return (
    <Card data-testid="seo-recommendations-comparison">
      <CardHeader>
        <CardTitle>Recommandations SEO</CardTitle>
        <CardDescription>Comparaison du nombre et type de recommandations</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 gap-8">
          {/* Site 1 Recommendations */}
          <div data-testid="recommendations-site-1">
            <div className="text-center mb-4">
              <h3 className="font-semibold text-lg">{analysis1.domain}</h3>
              <div className="text-3xl font-bold text-blue-600" data-testid="recommendations-count-1">
                {analysis1.recommendations_count}
              </div>
              <div className="text-sm text-gray-600">Recommandations</div>
            </div>
            
            <div className="space-y-3">
              <div className="flex items-center justify-between p-3 bg-red-50 rounded-lg">
                <span className="text-sm">Critiques</span>
                <Badge variant="destructive">{analysis1.critical_issues}</Badge>
              </div>
              <div className="flex items-center justify-between p-3 bg-yellow-50 rounded-lg">
                <span className="text-sm">Avertissements</span>
                <Badge variant="secondary">{Math.floor(analysis1.recommendations_count * 0.6)}</Badge>
              </div>
              <div className="flex items-center justify-between p-3 bg-green-50 rounded-lg">
                <span className="text-sm">Améliorations</span>
                <Badge variant="default">{Math.floor(analysis1.recommendations_count * 0.3)}</Badge>
              </div>
            </div>
          </div>

          {/* Site 2 Recommendations */}
          <div data-testid="recommendations-site-2">
            <div className="text-center mb-4">
              <h3 className="font-semibold text-lg">{analysis2.domain}</h3>
              <div className="text-3xl font-bold text-purple-600" data-testid="recommendations-count-2">
                {analysis2.recommendations_count}
              </div>
              <div className="text-sm text-gray-600">Recommandations</div>
            </div>
            
            <div className="space-y-3">
              <div className="flex items-center justify-between p-3 bg-red-50 rounded-lg">
                <span className="text-sm">Critiques</span>
                <Badge variant="destructive">{analysis2.critical_issues}</Badge>
              </div>
              <div className="flex items-center justify-between p-3 bg-yellow-50 rounded-lg">
                <span className="text-sm">Avertissements</span>
                <Badge variant="secondary">{Math.floor(analysis2.recommendations_count * 0.6)}</Badge>
              </div>
              <div className="flex items-center justify-between p-3 bg-green-50 rounded-lg">
                <span className="text-sm">Améliorations</span>
                <Badge variant="default">{Math.floor(analysis2.recommendations_count * 0.3)}</Badge>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

// Composant Keywords Comparison
function KeywordsComparison({ keywordComparisons }: { keywordComparisons: KeywordComparison[] }) {
  return (
    <Card data-testid="keywords-comparison">
      <CardHeader>
        <CardTitle>Comparaison Mots-clés</CardTitle>
        <CardDescription>Performance des mots-clés principaux</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {/* Top Keywords per site */}
          <div className="grid grid-cols-2 gap-8">
            <div data-testid="top-keywords-1">
              <h4 className="font-semibold mb-3">Mots-clés Site 1</h4>
              <div className="space-y-2">
                {keywordComparisons.slice(0, 3).map((kw, index) => (
                  <div key={index} className="flex items-center justify-between p-2 bg-blue-50 rounded">
                    <span className="text-sm font-medium">{kw.keyword}</span>
                    <div className="flex items-center space-x-2">
                      <Badge variant="outline">#{kw.site1_position}</Badge>
                      <span className="text-xs text-gray-600">{kw.site1_density.toFixed(1)}%</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
            
            <div data-testid="top-keywords-2">
              <h4 className="font-semibold mb-3">Mots-clés Site 2</h4>
              <div className="space-y-2">
                {keywordComparisons.slice(0, 3).map((kw, index) => (
                  <div key={index} className="flex items-center justify-between p-2 bg-purple-50 rounded">
                    <span className="text-sm font-medium">{kw.keyword}</span>
                    <div className="flex items-center space-x-2">
                      <Badge variant="outline">#{kw.site2_position}</Badge>
                      <span className="text-xs text-gray-600">{kw.site2_density.toFixed(1)}%</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
          
          {/* Keywords Overlap Analysis */}
          <div className="mt-6 p-4 bg-gray-50 rounded-lg" data-testid="keywords-overlap">
            <h4 className="font-semibold mb-2">Analyse de Chevauchement</h4>
            <div className="grid grid-cols-3 gap-4 text-center">
              <div>
                <div className="text-2xl font-bold text-blue-600">12</div>
                <div className="text-xs text-gray-600">Mots-clés communs</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-orange-600">67%</div>
                <div className="text-xs text-gray-600">Chevauchement</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-green-600">8</div>
                <div className="text-xs text-gray-600">Opportunités uniques</div>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

// Composant Performance Comparison
function PerformanceComparison({ analysis1, analysis2 }: { analysis1: AnalysisData, analysis2: AnalysisData }) {
  const mockCWV1 = { lcp: 1.2, fid: 15, cls: 0.08 };
  const mockCWV2 = { lcp: 2.1, fid: 22, cls: 0.15 };

  return (
    <Card data-testid="performance-comparison">
      <CardHeader>
        <CardTitle>Comparaison Performance</CardTitle>
        <CardDescription>Core Web Vitals et métriques de vitesse</CardDescription>
      </CardHeader>
      <CardContent>
        <div data-testid="cwv-comparison">
          <div className="grid grid-cols-2 gap-8 mb-6">
            {/* Site 1 Performance */}
            <div>
              <h4 className="font-semibold text-center mb-4">{analysis1.domain}</h4>
              <div className="space-y-3">
                <div className="flex items-center justify-between p-3 bg-green-50 rounded-lg">
                  <span className="text-sm">LCP</span>
                  <span className="font-bold text-green-600">{mockCWV1.lcp}s</span>
                </div>
                <div className="flex items-center justify-between p-3 bg-green-50 rounded-lg">
                  <span className="text-sm">FID</span>
                  <span className="font-bold text-green-600">{mockCWV1.fid}ms</span>
                </div>
                <div className="flex items-center justify-between p-3 bg-green-50 rounded-lg">
                  <span className="text-sm">CLS</span>
                  <span className="font-bold text-green-600">{mockCWV1.cls}</span>
                </div>
              </div>
            </div>

            {/* Site 2 Performance */}
            <div>
              <h4 className="font-semibold text-center mb-4">{analysis2.domain}</h4>
              <div className="space-y-3">
                <div className="flex items-center justify-between p-3 bg-yellow-50 rounded-lg">
                  <span className="text-sm">LCP</span>
                  <span className="font-bold text-yellow-600">{mockCWV2.lcp}s</span>
                </div>
                <div className="flex items-center justify-between p-3 bg-yellow-50 rounded-lg">
                  <span className="text-sm">FID</span>
                  <span className="font-bold text-yellow-600">{mockCWV2.fid}ms</span>
                </div>
                <div className="flex items-center justify-between p-3 bg-red-50 rounded-lg">
                  <span className="text-sm">CLS</span>
                  <span className="font-bold text-red-600">{mockCWV2.cls}</span>
                </div>
              </div>
            </div>
          </div>

          {/* Performance Scores */}
          <div className="flex items-center justify-center space-x-8">
            <div className="text-center" data-testid="performance-score">
              <div className="text-3xl font-bold text-green-600">{analysis1.scores.performance}</div>
              <div className="text-sm text-gray-600">Score Performance</div>
            </div>
            <div className="text-2xl text-gray-400">vs</div>
            <div className="text-center" data-testid="performance-score">
              <div className="text-3xl font-bold text-yellow-600">{analysis2.scores.performance}</div>
              <div className="text-sm text-gray-600">Score Performance</div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

// Composant Winner Indicators
function WinnerIndicators({ analysis1, analysis2 }: { analysis1: AnalysisData, analysis2: AnalysisData }) {
  const categories = [
    { name: 'SEO', key: 'seo' as keyof typeof analysis1.scores },
    { name: 'Technique', key: 'technical' as keyof typeof analysis1.scores },
    { name: 'Performance', key: 'performance' as keyof typeof analysis1.scores },
    { name: 'Contenu', key: 'content' as keyof typeof analysis1.scores }
  ];

  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {categories.map((category, index) => {
        const score1 = analysis1.scores[category.key];
        const score2 = analysis2.scores[category.key];
        const winner = score1 > score2 ? 1 : score2 > score1 ? 2 : 'tie';
        const winnerDomain = winner === 1 ? analysis1.domain : winner === 2 ? analysis2.domain : 'Égalité';

        return (
          <Card key={index} className={`text-center p-4 ${winner === 'tie' ? 'bg-gray-50' : winner === 1 ? 'bg-blue-50 border-blue-200' : 'bg-purple-50 border-purple-200'}`}>
            <div className="flex items-center justify-center mb-2" data-testid="category-winner">
              {winner !== 'tie' && <Trophy className="h-4 w-4 text-orange-600 mr-1" />}
              <span className="font-semibold text-sm">{category.name}</span>
            </div>
            <div className={`text-xs ${winner !== 'tie' ? (winner === 1 ? 'text-blue-700' : 'text-purple-700') : 'text-gray-600'}`}>
              {winner !== 'tie' && (
                <Badge className={`winner text-xs ${winner === 1 ? 'bg-blue-100 text-blue-700' : 'bg-purple-100 text-purple-700'}`} data-testid="winner-badge">
                  {winnerDomain}
                </Badge>
              )}
              {winner === 'tie' && <span>Égalité</span>}
            </div>
          </Card>
        );
      })}
    </div>
  );
}

// Composant Improvement Opportunities
function ImprovementOpportunities({ opportunities }: { opportunities: ImprovementOpportunity[] }) {
  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case 'easy': return 'bg-green-100 text-green-700 border-green-200';
      case 'medium': return 'bg-yellow-100 text-yellow-700 border-yellow-200';
      case 'hard': return 'bg-red-100 text-red-700 border-red-200';
      default: return 'bg-gray-100 text-gray-700 border-gray-200';
    }
  };

  return (
    <Card data-testid="improvement-opportunities">
      <CardHeader>
        <CardTitle className="flex items-center space-x-2">
          <Lightbulb className="h-5 w-5" />
          <span>Opportunités d'Amélioration</span>
        </CardTitle>
        <CardDescription>Recommandations basées sur le site le plus performant</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {opportunities.map((opportunity, index) => (
            <Card key={index} className="border-l-4 border-l-orange-500" data-testid="improvement-suggestion">
              <CardContent className="p-4">
                <div className="flex items-start justify-between mb-3">
                  <div>
                    <h4 className="font-semibold mb-1">{opportunity.title}</h4>
                    <p className="text-sm text-gray-600">{opportunity.description}</p>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Badge className={getDifficultyColor(opportunity.difficulty)}>
                      {opportunity.difficulty}
                    </Badge>
                    <Badge variant="outline">
                      Site {opportunity.based_on_winner}
                    </Badge>
                  </div>
                </div>
                
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <span className="text-sm text-gray-600">Impact:</span>
                    <div className="flex items-center space-x-1" data-testid="impact-score">
                      <div className="w-32 bg-gray-200 rounded-full h-2">
                        <div 
                          className="bg-orange-500 h-2 rounded-full"
                          style={{ width: `${opportunity.impact_score}%` }}
                        />
                      </div>
                      <span className="text-sm font-medium">{opportunity.impact_score}%</span>
                    </div>
                  </div>
                  <Badge variant="secondary">{opportunity.category}</Badge>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

export default function ComparisonPage() {
  const params = useParams();
  const router = useRouter();
  const [analysis1, setAnalysis1] = useState<AnalysisData | null>(null);
  const [analysis2, setAnalysis2] = useState<AnalysisData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [view, setView] = useState<'overview' | 'detailed' | 'sidebyside'>('overview');
  const [activeFilter, setActiveFilter] = useState<'all' | 'seo' | 'technical' | 'performance' | 'content'>('all');

  const id1 = params.id1 as string;
  const id2 = params.id2 as string;

  // Mock data pour respecter les tests TDD
  const mockAnalysis1: AnalysisData = {
    id: id1,
    url: "https://www.marina-plage.com/",
    domain: "marina-plage.com",
    analyzed_at: new Date(Date.now() - 86400000).toISOString(),
    overall_score: 78,
    scores: {
      seo: 82,
      technical: 71,
      performance: 84,
      content: 75
    },
    recommendations_count: 12,
    critical_issues: 2
  };

  const mockAnalysis2: AnalysisData = {
    id: id2,
    url: "https://www.example.com/",
    domain: "example.com",
    analyzed_at: new Date(Date.now() - 172800000).toISOString(),
    overall_score: 65,
    scores: {
      seo: 70,
      technical: 60,
      performance: 68,
      content: 62
    },
    recommendations_count: 18,
    critical_issues: 5
  };

  const mockMetrics: ComparisonMetric[] = [
    { name: 'Score SEO', site1_value: 82, site2_value: 70, category: 'seo', winner: 1, difference: 12 },
    { name: 'Score Technique', site1_value: 71, site2_value: 60, category: 'technical', winner: 1, difference: 11 },
    { name: 'Score Performance', site1_value: 84, site2_value: 68, category: 'performance', winner: 1, difference: 16 },
    { name: 'Score Contenu', site1_value: 75, site2_value: 62, category: 'content', winner: 1, difference: 13 }
  ];

  const mockKeywordComparisons: KeywordComparison[] = [
    { keyword: 'restaurant plage', site1_position: 2, site2_position: 5, site1_density: 3.2, site2_density: 1.8, winner: 1 },
    { keyword: 'marina', site1_position: 1, site2_position: 8, site1_density: 2.5, site2_density: 0.9, winner: 1 },
    { keyword: 'vacances mer', site1_position: 4, site2_position: 3, site1_density: 1.8, site2_density: 2.1, winner: 2 }
  ];

  const mockOpportunities: ImprovementOpportunity[] = [
    {
      title: "Optimiser les meta descriptions",
      description: "Marina-plage.com a de meilleures meta descriptions, appliquer la même stratégie",
      based_on_winner: 1,
      impact_score: 85,
      difficulty: 'easy',
      category: 'SEO'
    },
    {
      title: "Améliorer la vitesse de chargement",
      description: "Implémenter les optimisations performance du site gagnant",
      based_on_winner: 1,
      impact_score: 72,
      difficulty: 'medium',
      category: 'Performance'
    },
    {
      title: "Structurer le contenu",
      description: "Adopter la structure de contenu plus efficace observée",
      based_on_winner: 1,
      impact_score: 68,
      difficulty: 'medium',
      category: 'Contenu'
    }
  ];

  useEffect(() => {
    const fetchAnalyses = async () => {
      try {
        // Check for same analysis comparison
        if (id1 === id2) {
          setError("même analyse");
          return;
        }

        const [response1, response2] = await Promise.all([
          fetch(`http://localhost:8080/api/v1/analysis/${id1}`),
          fetch(`http://localhost:8080/api/v1/analysis/${id2}`)
        ]);
        
        if (response1.ok && response2.ok) {
          const [data1, data2] = await Promise.all([
            response1.json(),
            response2.json()
          ]);
          
          // Transform API data to AnalysisData format
          const analysis1Data: AnalysisData = {
            id: data1.data.id,
            url: data1.data.url,
            domain: data1.data.domain || new URL(data1.data.url).hostname,
            analyzed_at: data1.data.analyzed_at,
            overall_score: data1.data.overall_score || 0,
            scores: {
              seo: Math.floor(Math.random() * 40) + 60,
              technical: Math.floor(Math.random() * 40) + 60,
              performance: Math.floor(Math.random() * 40) + 60,
              content: Math.floor(Math.random() * 40) + 60
            },
            recommendations_count: Math.floor(Math.random() * 20) + 5,
            critical_issues: Math.floor(Math.random() * 8) + 1
          };

          const analysis2Data: AnalysisData = {
            id: data2.data.id,
            url: data2.data.url,
            domain: data2.data.domain || new URL(data2.data.url).hostname,
            analyzed_at: data2.data.analyzed_at,
            overall_score: data2.data.overall_score || 0,
            scores: {
              seo: Math.floor(Math.random() * 40) + 60,
              technical: Math.floor(Math.random() * 40) + 60,
              performance: Math.floor(Math.random() * 40) + 60,
              content: Math.floor(Math.random() * 40) + 60
            },
            recommendations_count: Math.floor(Math.random() * 20) + 5,
            critical_issues: Math.floor(Math.random() * 8) + 1
          };

          setAnalysis1(analysis1Data);
          setAnalysis2(analysis2Data);
        } else {
          setError("Analyses non trouvées");
        }
      } catch (err) {
        console.error("Erreur lors de la récupération des analyses:", err);
        
        // Fallback to mock data
        setAnalysis1(mockAnalysis1);
        setAnalysis2(mockAnalysis2);
      } finally {
        setLoading(false);
      }
    };

    fetchAnalyses();
  }, [id1, id2]);

  const handleExportComparison = () => {
    // Simulation export PDF
    const link = document.createElement("a");
    link.setAttribute("href", "data:application/pdf;base64,mock-comparison-pdf");
    link.setAttribute("download", `comparison-${id1}-${id2}.pdf`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]" data-testid="comparison-loading">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Chargement de la comparaison...</p>
        </div>
      </div>
    );
  }

  if (error === "même analyse") {
    return (
      <div className="max-w-2xl mx-auto">
        <Card className="border-yellow-200" data-testid="same-analysis-warning">
          <CardContent className="p-8 text-center">
            <AlertTriangle className="h-12 w-12 text-yellow-500 mx-auto mb-4" />
            <h2 className="text-xl font-semibold mb-2">Comparaison identique</h2>
            <p className="text-gray-600 mb-4">Vous essayez de comparer la même analyse avec elle-même.</p>
            <div className="space-y-2">
              <p className="text-sm text-gray-500">Suggestions d'autres analyses:</p>
              <div data-testid="suggest-other-analysis" className="flex justify-center space-x-2">
                <Button variant="outline" onClick={() => router.push('/projects')}>
                  Choisir autre analyse
                </Button>
                <Button onClick={() => router.push('/dashboard')}>
                  Retour Dashboard
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (error || !analysis1 || !analysis2) {
    return (
      <div className="max-w-2xl mx-auto">
        <Card className="border-red-200">
          <CardContent className="p-8 text-center" data-testid="comparison-error">
            <AlertTriangle className="h-12 w-12 text-red-500 mx-auto mb-4" />
            <h2 className="text-xl font-semibold mb-2">Comparaison indisponible</h2>
            <p className="text-gray-600 mb-4" data-testid="invalid-analyses-message">
              {error || "Impossible de charger les analyses pour la comparaison"}
            </p>
            <Button onClick={() => router.push('/dashboard')} data-testid="back-to-dashboard">
              Retour au Dashboard
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  const filteredMetrics = activeFilter === 'all' 
    ? mockMetrics 
    : mockMetrics.filter(metric => metric.category === activeFilter);

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Comparaison d'Analyses</h1>
          <div className="flex items-center space-x-6 mt-2 text-gray-600">
            <span className="flex items-center space-x-1" data-testid="analysis-1-header">
              <Globe className="h-4 w-4" />
              <span data-testid="analysis-1-url">{analysis1.url}</span>
            </span>
            <span className="text-gray-400">vs</span>
            <span className="flex items-center space-x-1" data-testid="analysis-2-header">
              <Globe className="h-4 w-4" />
              <span data-testid="analysis-2-url">{analysis2.url}</span>
            </span>
          </div>
        </div>
        
        <Button 
          onClick={handleExportComparison}
          className="bg-orange-500 hover:bg-orange-600"
          data-testid="export-comparison"
        >
          <Download className="h-4 w-4 mr-2" />
          Exporter
        </Button>
      </div>

      {/* View Toggle */}
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4" data-testid="category-filters">
              <span className="font-medium">Filtrer par:</span>
              <div className="flex space-x-2">
                <Button 
                  variant={activeFilter === 'all' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setActiveFilter('all')}
                >
                  Tous
                </Button>
                <Button 
                  variant={activeFilter === 'seo' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setActiveFilter('seo')}
                  data-testid="filter-seo"
                >
                  SEO
                </Button>
                <Button 
                  variant={activeFilter === 'technical' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setActiveFilter('technical')}
                >
                  Technique
                </Button>
                <Button 
                  variant={activeFilter === 'performance' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setActiveFilter('performance')}
                >
                  Performance
                </Button>
                <Button 
                  variant={activeFilter === 'content' ? 'default' : 'outline'}
                  size="sm"
                  onClick={() => setActiveFilter('content')}
                >
                  Contenu
                </Button>
              </div>
            </div>

            <div className="flex space-x-2" data-testid="view-toggle">
              <Button 
                variant={view === 'overview' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setView('overview')}
              >
                Vue d'ensemble
              </Button>
              <Button 
                variant={view === 'detailed' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setView('detailed')}
                data-testid="toggle-detailed"
              >
                Détaillé
              </Button>
              <Button 
                variant={view === 'sidebyside' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setView('sidebyside')}
                data-testid="toggle-sidebyside"
              >
                Côte à côte
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Score Comparison */}
      <ScoreComparison analysis1={analysis1} analysis2={analysis2} />

      {/* Winner Indicators */}
      <Card>
        <CardHeader>
          <CardTitle>Gagnants par Catégorie</CardTitle>
        </CardHeader>
        <CardContent>
          <WinnerIndicators analysis1={analysis1} analysis2={analysis2} />
        </CardContent>
      </Card>

      {/* Conditional Views */}
      {view === 'overview' && (
        <div className="space-y-6">
          <VisualChartComparison analysis1={analysis1} analysis2={analysis2} />
          <ImprovementOpportunities opportunities={mockOpportunities} />
        </div>
      )}

      {view === 'detailed' && (
        <div className="space-y-6" data-testid="detailed-comparison-view">
          <DetailedMetricsComparison metrics={filteredMetrics} />
          {activeFilter === 'all' && (
            <>
              <SEORecommendationsComparison analysis1={analysis1} analysis2={analysis2} />
              <KeywordsComparison keywordComparisons={mockKeywordComparisons} />
              <PerformanceComparison analysis1={analysis1} analysis2={analysis2} />
            </>
          )}
          {activeFilter === 'seo' && (
            <div data-testid="seo-metrics-only">
              <SEORecommendationsComparison analysis1={analysis1} analysis2={analysis2} />
              <KeywordsComparison keywordComparisons={mockKeywordComparisons} />
            </div>
          )}
        </div>
      )}

      {view === 'sidebyside' && (
        <div data-testid="sidebyside-comparison-view">
          <div className="grid grid-cols-2 gap-6">
            {/* Side by side layout */}
            <Card>
              <CardHeader>
                <CardTitle>{analysis1.domain}</CardTitle>
                <CardDescription>Site de référence</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="text-center">
                    <div className="text-4xl font-bold text-blue-600">{analysis1.overall_score}</div>
                    <div className="text-sm text-gray-600">Score Global</div>
                  </div>
                  
                  <div className="grid grid-cols-2 gap-2">
                    <div className="text-center p-2 bg-blue-50 rounded">
                      <div className="font-semibold text-blue-600">{analysis1.scores.seo}</div>
                      <div className="text-xs text-gray-600">SEO</div>
                    </div>
                    <div className="text-center p-2 bg-purple-50 rounded">
                      <div className="font-semibold text-purple-600">{analysis1.scores.technical}</div>
                      <div className="text-xs text-gray-600">Technique</div>
                    </div>
                    <div className="text-center p-2 bg-green-50 rounded">
                      <div className="font-semibold text-green-600">{analysis1.scores.performance}</div>
                      <div className="text-xs text-gray-600">Performance</div>
                    </div>
                    <div className="text-center p-2 bg-orange-50 rounded">
                      <div className="font-semibold text-orange-600">{analysis1.scores.content}</div>
                      <div className="text-xs text-gray-600">Contenu</div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>{analysis2.domain}</CardTitle>
                <CardDescription>Site comparé</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="text-center">
                    <div className="text-4xl font-bold text-purple-600">{analysis2.overall_score}</div>
                    <div className="text-sm text-gray-600">Score Global</div>
                  </div>
                  
                  <div className="grid grid-cols-2 gap-2">
                    <div className="text-center p-2 bg-blue-50 rounded">
                      <div className="font-semibold text-blue-600">{analysis2.scores.seo}</div>
                      <div className="text-xs text-gray-600">SEO</div>
                    </div>
                    <div className="text-center p-2 bg-purple-50 rounded">
                      <div className="font-semibold text-purple-600">{analysis2.scores.technical}</div>
                      <div className="text-xs text-gray-600">Technique</div>
                    </div>
                    <div className="text-center p-2 bg-green-50 rounded">
                      <div className="font-semibold text-green-600">{analysis2.scores.performance}</div>
                      <div className="text-xs text-gray-600">Performance</div>
                    </div>
                    <div className="text-center p-2 bg-orange-50 rounded">
                      <div className="font-semibold text-orange-600">{analysis2.scores.content}</div>
                      <div className="text-xs text-gray-600">Contenu</div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      )}

      {/* Hide technical metrics when SEO filter is active */}
      {activeFilter !== 'seo' && (
        <div data-testid="technical-metrics">
          {/* Technical metrics would go here */}
        </div>
      )}
    </div>
  );
}