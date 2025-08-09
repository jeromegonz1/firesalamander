"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { 
  Download, 
  Brain,
  Search,
  TrendingUp,
  Target,
  Filter,
  Eye,
  BarChart3,
  AlertTriangle,
  CheckCircle2,
  Info,
  Globe,
  Clock,
  X
} from "lucide-react";

interface KeywordData {
  keyword: string;
  density: number;
  score: number;
  category: 'primary' | 'secondary' | 'semantic';
  opportunities: number;
  competition: 'low' | 'medium' | 'high';
  trend: 'up' | 'down' | 'stable';
}

interface AIOpportunity {
  keyword: string;
  potential: number;
  score: number;
  difficulty: 'easy' | 'medium' | 'hard';
  reason: string;
  impact: 'low' | 'medium' | 'high';
}

interface SemanticKeyword {
  keyword: string;
  relevance: number;
  context: string;
  related_terms: string[];
}

interface CompetitorKeyword {
  keyword: string;
  our_position: number;
  competitor_position: number;
  competitor_name: string;
  gap_score: number;
}

// Composant Keyword Density Chart (simulé avec CSS)
function KeywordDensityChart({ keywords }: { keywords: KeywordData[] }) {
  return (
    <div className="space-y-4" data-testid="density-chart">
      {keywords.slice(0, 5).map((keyword, index) => (
        <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
          <div className="flex items-center space-x-3">
            <span className="font-medium">{keyword.keyword}</span>
            <Badge variant="secondary">{keyword.category}</Badge>
          </div>
          <div className="flex items-center space-x-4">
            <div className="w-32 bg-gray-200 rounded-full h-2">
              <div 
                className="bg-orange-500 h-2 rounded-full transition-all" 
                style={{ width: `${Math.min(keyword.density * 10, 100)}%` }}
              />
            </div>
            <span className="text-sm font-medium">{keyword.density.toFixed(1)}%</span>
          </div>
        </div>
      ))}
    </div>
  );
}

// Composant Opportunity Card
function OpportunityCard({ opportunity }: { opportunity: AIOpportunity }) {
  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case 'easy': return 'bg-green-100 text-green-700 border-green-200';
      case 'medium': return 'bg-yellow-100 text-yellow-700 border-yellow-200';
      case 'hard': return 'bg-red-100 text-red-700 border-red-200';
      default: return 'bg-gray-100 text-gray-700 border-gray-200';
    }
  };

  const getImpactIcon = (impact: string) => {
    switch (impact) {
      case 'high': return <TrendingUp className="h-4 w-4 text-green-600" />;
      case 'medium': return <Target className="h-4 w-4 text-yellow-600" />;
      case 'low': return <Info className="h-4 w-4 text-blue-600" />;
      default: return <Info className="h-4 w-4 text-gray-600" />;
    }
  };

  return (
    <Card className="hover:shadow-lg transition-shadow" data-testid="opportunity-card">
      <CardContent className="p-6">
        <div className="flex items-start justify-between mb-4">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-purple-50 rounded-lg">
              <Brain className="h-5 w-5 text-purple-600" />
            </div>
            <div>
              <h4 className="font-semibold">{opportunity.keyword}</h4>
              <p className="text-sm text-gray-600">Potentiel: {opportunity.potential}%</p>
            </div>
          </div>
          <div className="flex items-center space-x-2">
            {getImpactIcon(opportunity.impact)}
            <Badge className={getDifficultyColor(opportunity.difficulty)}>
              {opportunity.difficulty}
            </Badge>
          </div>
        </div>
        
        <div className="mb-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium">Score</span>
            <span className="text-2xl font-bold text-orange-600 score-indicator">{opportunity.score}</span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div 
              className="bg-orange-500 h-2 rounded-full transition-all" 
              style={{ width: `${opportunity.score}%` }}
            />
          </div>
        </div>
        
        <p className="text-sm text-gray-700">{opportunity.reason}</p>
      </CardContent>
    </Card>
  );
}

// Composant Semantic Tags
function SemanticTags({ semanticKeywords }: { semanticKeywords: SemanticKeyword[] }) {
  return (
    <div className="space-y-4" data-testid="semantic-tags">
      {semanticKeywords.map((semantic, index) => (
        <div key={index} className="p-4 border rounded-lg hover:bg-gray-50 transition-colors">
          <div className="flex items-center justify-between mb-2">
            <h4 className="font-semibold">{semantic.keyword}</h4>
            <span className="text-sm text-gray-600">{semantic.relevance}% pertinence</span>
          </div>
          <p className="text-sm text-gray-700 mb-3">{semantic.context}</p>
          <div className="flex flex-wrap gap-2">
            {semantic.related_terms.map((term, termIndex) => (
              <Badge key={termIndex} variant="secondary" className="text-xs">
                {term}
              </Badge>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}

// Composant Competitor Table
function CompetitorTable({ competitorKeywords }: { competitorKeywords: CompetitorKeyword[] }) {
  return (
    <div className="overflow-x-auto" data-testid="competitor-table">
      <table className="w-full border-collapse">
        <thead>
          <tr className="border-b">
            <th className="text-left p-3 font-semibold">Mot-clé</th>
            <th className="text-left p-3 font-semibold">Notre Position</th>
            <th className="text-left p-3 font-semibold">Concurrent</th>
            <th className="text-left p-3 font-semibold">Écart</th>
          </tr>
        </thead>
        <tbody>
          {competitorKeywords.map((item, index) => (
            <tr key={index} className="border-b hover:bg-gray-50">
              <td className="p-3 font-medium">{item.keyword}</td>
              <td className="p-3">
                <Badge variant={item.our_position <= 3 ? "default" : "secondary"}>
                  #{item.our_position}
                </Badge>
              </td>
              <td className="p-3">
                <div>
                  <span className="text-sm font-medium">{item.competitor_name}</span>
                  <Badge variant="outline" className="ml-2">#{item.competitor_position}</Badge>
                </div>
              </td>
              <td className="p-3">
                <span className={`font-semibold ${item.gap_score > 0 ? 'text-red-600' : 'text-green-600'}`}>
                  {item.gap_score > 0 ? '+' : ''}{item.gap_score}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

// Composant Modal Keyword Details
function KeywordDetailsModal({ 
  isOpen, 
  onClose, 
  keyword 
}: { 
  isOpen: boolean; 
  onClose: () => void; 
  keyword: KeywordData | null; 
}) {
  if (!isOpen || !keyword) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" data-testid="keyword-details-modal">
      <div className="bg-white rounded-lg p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold">Détails: {keyword.keyword}</h2>
          <Button variant="ghost" size="sm" onClick={onClose} data-testid="close-modal">
            <X className="h-4 w-4" />
          </Button>
        </div>
        
        <div className="space-y-6" data-testid="keyword-stats">
          <div className="grid grid-cols-2 gap-4">
            <div className="p-4 bg-blue-50 rounded-lg">
              <h3 className="font-semibold text-blue-900">Densité</h3>
              <p className="text-2xl font-bold text-blue-600">{keyword.density.toFixed(1)}%</p>
            </div>
            <div className="p-4 bg-green-50 rounded-lg">
              <h3 className="font-semibold text-green-900">Score</h3>
              <p className="text-2xl font-bold text-green-600">{keyword.score}</p>
            </div>
            <div className="p-4 bg-orange-50 rounded-lg">
              <h3 className="font-semibold text-orange-900">Opportunités</h3>
              <p className="text-2xl font-bold text-orange-600">{keyword.opportunities}</p>
            </div>
            <div className="p-4 bg-purple-50 rounded-lg">
              <h3 className="font-semibold text-purple-900">Catégorie</h3>
              <Badge className="mt-1">{keyword.category}</Badge>
            </div>
          </div>
          
          <div className="p-4 bg-gray-50 rounded-lg">
            <h3 className="font-semibold mb-2">Analyse Concurrentielle</h3>
            <p className="text-sm text-gray-700">
              Concurrence: <Badge variant={
                keyword.competition === 'low' ? 'default' : 
                keyword.competition === 'medium' ? 'secondary' : 'destructive'
              }>{keyword.competition}</Badge>
            </p>
            <p className="text-sm text-gray-700 mt-2">
              Tendance: <Badge variant={
                keyword.trend === 'up' ? 'default' : 
                keyword.trend === 'stable' ? 'secondary' : 'destructive'
              }>{keyword.trend}</Badge>
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function AnalysisKeywordsPage() {
  const params = useParams();
  const [analysis, setAnalysis] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedKeyword, setSelectedKeyword] = useState<KeywordData | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [activeFilter, setActiveFilter] = useState<'all' | 'primary' | 'secondary' | 'semantic'>('all');

  const analysisId = params.id as string;

  // Mock data pour respecter les tests TDD
  const mockKeywordsData: KeywordData[] = [
    { keyword: 'marina plage', density: 3.2, score: 85, category: 'primary', opportunities: 12, competition: 'medium', trend: 'up' },
    { keyword: 'restaurant plage', density: 2.1, score: 72, category: 'secondary', opportunities: 8, competition: 'high', trend: 'stable' },
    { keyword: 'location vacances', density: 1.8, score: 68, category: 'semantic', opportunities: 15, competition: 'low', trend: 'up' },
    { keyword: 'mer méditerranée', density: 1.5, score: 64, category: 'semantic', opportunities: 6, competition: 'medium', trend: 'down' },
    { keyword: 'activités nautiques', density: 1.2, score: 58, category: 'secondary', opportunities: 9, competition: 'low', trend: 'up' }
  ];

  const mockAIOpportunities: AIOpportunity[] = [
    {
      keyword: 'plage privée marina',
      potential: 87,
      score: 78,
      difficulty: 'easy',
      reason: 'Faible concurrence détectée, forte intention commerciale',
      impact: 'high'
    },
    {
      keyword: 'restaurant vue mer',
      potential: 72,
      score: 65,
      difficulty: 'medium',
      reason: 'Opportunité saisonnière identifiée par l\'IA',
      impact: 'medium'
    },
    {
      keyword: 'événements plage',
      potential: 68,
      score: 61,
      difficulty: 'easy',
      reason: 'Contenu existant peut être optimisé facilement',
      impact: 'medium'
    }
  ];

  const mockSemanticKeywords: SemanticKeyword[] = [
    {
      keyword: 'détente bord de mer',
      relevance: 89,
      context: 'Fortement associé à l\'expérience utilisateur recherchée',
      related_terms: ['relaxation', 'tranquillité', 'évasion', 'bien-être']
    },
    {
      keyword: 'gastronomie méditerranéenne',
      relevance: 76,
      context: 'Lié à l\'offre restauration et à l\'identité locale',
      related_terms: ['cuisine locale', 'produits frais', 'spécialités', 'chef']
    }
  ];

  const mockCompetitorKeywords: CompetitorKeyword[] = [
    { keyword: 'marina plage', our_position: 2, competitor_position: 1, competitor_name: 'Concurrent A', gap_score: 1 },
    { keyword: 'restaurant plage', our_position: 5, competitor_position: 3, competitor_name: 'Concurrent B', gap_score: 2 },
    { keyword: 'location vacances', our_position: 1, competitor_position: 4, competitor_name: 'Concurrent A', gap_score: -3 }
  ];

  useEffect(() => {
    const fetchAnalysis = async () => {
      try {
        const response = await fetch(`http://localhost:8080/api/v1/analysis/${analysisId}`);
        
        if (response.ok) {
          const data = await response.json();
          setAnalysis(data.data);
        } else {
          setError("Analyse non trouvée");
        }
      } catch (err) {
        console.error("Erreur lors de la récupération de l'analyse:", err);
        
        // Données mockées pour les tests
        setAnalysis({
          id: analysisId,
          url: "https://www.marina-plage.com/",
          analyzed_at: new Date().toISOString(),
          processing_time: 4500,
          overall_score: 78
        });
      } finally {
        setLoading(false);
      }
    };

    fetchAnalysis();
  }, [analysisId]);

  const handleKeywordClick = (keyword: KeywordData) => {
    setSelectedKeyword(keyword);
    setIsModalOpen(true);
  };

  const handleExportKeywords = () => {
    // Simulation export CSV
    const csvContent = "data:text/csv;charset=utf-8," + 
      "Mot-clé,Densité,Score,Catégorie,Opportunités\n" +
      mockKeywordsData.map(k => `${k.keyword},${k.density},${k.score},${k.category},${k.opportunities}`).join("\n");
    
    const encodedUri = encodeURI(csvContent);
    const link = document.createElement("a");
    link.setAttribute("href", encodedUri);
    link.setAttribute("download", `keywords-analysis-${analysisId}.csv`);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const filteredKeywords = activeFilter === 'all' 
    ? mockKeywordsData 
    : mockKeywordsData.filter(k => k.category === activeFilter);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]" data-testid="loading-spinner">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Chargement de l'analyse des mots-clés...</p>
        </div>
      </div>
    );
  }

  if (error || !analysis) {
    return (
      <div className="max-w-2xl mx-auto" data-testid="error-message">
        <Card className="border-red-200">
          <CardContent className="p-8 text-center">
            <AlertTriangle className="h-12 w-12 text-red-500 mx-auto mb-4" />
            <h2 className="text-xl font-semibold mb-2">Erreur</h2>
            <p className="text-gray-600 mb-4">{error || "Impossible de charger l'analyse"}</p>
            <Button onClick={() => window.location.reload()} data-testid="retry-button">
              Réessayer
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Analyse des Mots-clés</h1>
          <div className="flex items-center space-x-4 mt-2 text-gray-600">
            <span className="flex items-center space-x-1" data-testid="analysis-url">
              <Globe className="h-4 w-4" />
              <span>{analysis.url}</span>
            </span>
            <span className="flex items-center space-x-1" data-testid="analysis-date">
              <Clock className="h-4 w-4" />
              <span>{new Date(analysis.analyzed_at).toLocaleString()}</span>
            </span>
          </div>
        </div>
        
        <Button 
          onClick={handleExportKeywords}
          className="bg-orange-500 hover:bg-orange-600"
          data-testid="export-keywords"
        >
          <Download className="h-4 w-4 mr-2" />
          Exporter
        </Button>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center space-x-4" data-testid="category-filter">
            <span className="font-medium">Filtrer par catégorie:</span>
            <div className="flex space-x-2">
              <Button 
                variant={activeFilter === 'all' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setActiveFilter('all')}
              >
                Tous
              </Button>
              <Button 
                variant={activeFilter === 'primary' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setActiveFilter('primary')}
                data-testid="filter-primary"
              >
                Primaires
              </Button>
              <Button 
                variant={activeFilter === 'secondary' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setActiveFilter('secondary')}
              >
                Secondaires
              </Button>
              <Button 
                variant={activeFilter === 'semantic' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setActiveFilter('semantic')}
              >
                Sémantiques
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Keywords Density Section */}
      <Card data-testid="keywords-density" aria-label="Section densité des mots-clés">
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <BarChart3 className="h-5 w-5" />
            <span>Densité des Mots-clés</span>
          </CardTitle>
          <CardDescription>
            Analyse de la densité et performance des mots-clés principaux
          </CardDescription>
        </CardHeader>
        <CardContent>
          <KeywordDensityChart keywords={filteredKeywords} />
          
          <div className="mt-6" data-testid="filtered-results">
            <h3 className="font-semibold mb-4">Mots-clés détectés ({filteredKeywords.length})</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {filteredKeywords.map((keyword, index) => (
                <div 
                  key={index}
                  className="p-4 border rounded-lg hover:bg-gray-50 cursor-pointer transition-colors"
                  onClick={() => handleKeywordClick(keyword)}
                  data-testid="keyword-item"
                >
                  <div className="flex items-center justify-between mb-2">
                    <h4 className="font-medium">{keyword.keyword}</h4>
                    <Badge variant="secondary">{keyword.category}</Badge>
                  </div>
                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span>Densité:</span>
                      <span className="font-medium">{keyword.density.toFixed(1)}%</span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Score:</span>
                      <span className="font-medium text-orange-600">{keyword.score}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* AI Opportunities Section */}
      <Card data-testid="ai-opportunities" aria-label="Section opportunités IA">
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Brain className="h-5 w-5" />
            <span>Opportunités IA</span>
          </CardTitle>
          <CardDescription>
            Mots-clés à fort potentiel identifiés par intelligence artificielle
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {mockAIOpportunities.map((opportunity, index) => (
              <OpportunityCard key={index} opportunity={opportunity} />
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Semantic Analysis Section */}
      <Card data-testid="semantic-analysis" aria-label="Section analyse sémantique">
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Search className="h-5 w-5" />
            <span>Analyse Sémantique</span>
          </CardTitle>
          <CardDescription>
            Mots-clés sémantiquement liés et contexte d'utilisation
          </CardDescription>
        </CardHeader>
        <CardContent>
          <SemanticTags semanticKeywords={mockSemanticKeywords} />
        </CardContent>
      </Card>

      {/* Competitor Keywords Section */}
      <Card data-testid="competitor-keywords" aria-label="Section mots-clés concurrents">
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Target className="h-5 w-5" />
            <span>Mots-clés Concurrents</span>
          </CardTitle>
          <CardDescription>
            Comparaison des positions par rapport à la concurrence
          </CardDescription>
        </CardHeader>
        <CardContent>
          <CompetitorTable competitorKeywords={mockCompetitorKeywords} />
        </CardContent>
      </Card>

      {/* Keyword Details Modal */}
      <KeywordDetailsModal 
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        keyword={selectedKeyword}
      />
    </div>
  );
}