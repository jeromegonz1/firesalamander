"use client";

import { useState, useEffect } from 'react';
import { KeywordAnalysis, KeywordSortBy, KeywordFilterBy } from '@/types/keyword-analysis';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Progress } from '@/components/ui/progress';
import { 
  Search, Brain, TrendingUp, Target, BarChart3, Download, 
  Filter, RefreshCw, ChevronDown, ChevronRight, Star,
  Tag, Lightbulb, BookOpen, Activity, Globe, Zap
} from 'lucide-react';

interface KeywordAnalysisSectionProps {
  keywordData: KeywordAnalysis;
  analysisId: number;
}

// Composant pour afficher un mot-clé
function KeywordItem({ 
  keyword, 
  isExpanded, 
  onToggle 
}: { 
  keyword: KeywordAnalysis['foundKeywords'][0]; 
  isExpanded: boolean;
  onToggle: () => void;
}) {
  const getDensityColor = (density: number) => {
    if (density > 3) return 'text-red-600 bg-red-50';
    if (density > 1.5) return 'text-yellow-600 bg-yellow-50';
    return 'text-green-600 bg-green-50';
  };

  const getProminenceColor = (prominence: number) => {
    if (prominence > 70) return 'text-green-600';
    if (prominence > 40) return 'text-yellow-600';
    return 'text-red-600';
  };

  return (
    <div className="border rounded-lg p-4 hover:shadow-md transition-shadow" data-testid="keyword-item">
      <div className="flex items-center justify-between cursor-pointer" onClick={onToggle}>
        <div className="flex items-center space-x-3 flex-1">
          <div className="flex items-center space-x-1">
            {isExpanded ? <ChevronDown className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
            <Tag className="h-4 w-4 text-blue-500" />
          </div>
          
          <div className="flex-1">
            <div className="font-semibold text-lg" data-testid="keyword-text">
              {keyword.keyword}
            </div>
            
            <div className="flex items-center space-x-4 mt-1 text-sm text-gray-600">
              <span className={`px-2 py-1 rounded ${getDensityColor(keyword.density)}`} data-testid="keyword-density">
                {keyword.density}% densité
              </span>
              <span data-testid="keyword-occurrences">{keyword.occurrences} occurrences</span>
              <span className={getProminenceColor(keyword.prominence)} data-testid="keyword-prominence">
                {keyword.prominence}/100 prominence
              </span>
            </div>
          </div>
        </div>

        {/* Locations badges */}
        <div className="flex flex-wrap gap-1" data-testid="keyword-locations">
          {keyword.locations.map((location, idx) => (
            <Badge 
              key={idx} 
              variant="outline" 
              className="text-xs"
              data-testid="location-badge"
            >
              {location}
            </Badge>
          ))}
        </div>
      </div>

      {isExpanded && (
        <div className="mt-4 pt-4 border-t">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <h5 className="font-medium mb-2">Analyse détaillée</h5>
              <div className="space-y-2 text-sm">
                <div>Densité optimale : 1-3%</div>
                <div>Position : {keyword.prominence > 50 ? 'Excellente' : 'À améliorer'}</div>
                <div>Répartition : {keyword.locations.length} emplacements</div>
              </div>
            </div>
            <div>
              <h5 className="font-medium mb-2">Recommandations</h5>
              <div className="space-y-1 text-sm">
                {keyword.density > 3 && (
                  <div className="text-red-600">• Réduire la densité (sur-optimisation)</div>
                )}
                {!keyword.locations.includes('title') && (
                  <div className="text-yellow-600">• Ajouter dans le titre</div>
                )}
                {!keyword.locations.includes('h1') && (
                  <div className="text-yellow-600">• Ajouter dans H1</div>
                )}
                {keyword.prominence < 50 && (
                  <div className="text-blue-600">• Améliorer la position stratégique</div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

// Composant pour les suggestions IA
function AIKeywordSuggestion({ 
  suggestion 
}: { 
  suggestion: KeywordAnalysis['aiSuggestions'][0] 
}) {
  const [isExpanded, setIsExpanded] = useState(false);
  
  const getDifficultyColor = (difficulty: number) => {
    if (difficulty < 30) return 'text-green-600';
    if (difficulty < 60) return 'text-yellow-600';
    return 'text-red-600';
  };

  const getIntentBadge = (intent: string) => {
    const colors = {
      informational: 'bg-blue-100 text-blue-800',
      navigational: 'bg-gray-100 text-gray-800',
      transactional: 'bg-green-100 text-green-800',
      commercial: 'bg-purple-100 text-purple-800'
    };
    return colors[intent as keyof typeof colors] || 'bg-gray-100 text-gray-800';
  };

  return (
    <div className="border rounded-lg p-4 bg-gradient-to-r from-purple-50 to-blue-50" data-testid="ai-suggestion-item">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center space-x-2 mb-2">
            <Brain className="h-4 w-4 text-purple-600" />
            <h4 className="font-semibold text-lg" data-testid="suggestion-keyword">
              {suggestion.keyword}
            </h4>
            <Badge className={getIntentBadge(suggestion.intent)} data-testid="suggestion-intent">
              {suggestion.intent}
            </Badge>
          </div>
          
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-3 text-sm">
            <div>
              <div className="text-gray-600">Volume</div>
              <div className="font-semibold" data-testid="suggestion-volume">
                {suggestion.searchVolume.toLocaleString()}
              </div>
            </div>
            <div>
              <div className="text-gray-600">Difficulté</div>
              <div className={`font-semibold ${getDifficultyColor(suggestion.difficulty)}`} data-testid="suggestion-difficulty">
                {suggestion.difficulty}/100
              </div>
            </div>
            <div>
              <div className="text-gray-600">CPC</div>
              <div className="font-semibold" data-testid="suggestion-cpc">
                {suggestion.cpc.toFixed(2)}€
              </div>
            </div>
            <div>
              <div className="text-gray-600">Tendance</div>
              <div className="flex items-center space-x-1">
                {suggestion.trendData?.direction === 'up' ? 
                  <TrendingUp className="h-3 w-3 text-green-500" /> :
                  <Activity className="h-3 w-3 text-gray-500" />
                }
                <span className="text-xs">
                  {suggestion.trendData?.changePercent}%
                </span>
              </div>
            </div>
          </div>
          
          <p className="text-gray-700 text-sm mb-3" data-testid="suggestion-reason">
            <strong>Pourquoi :</strong> {suggestion.reason}
          </p>
          
          <div className="flex items-center space-x-2">
            <Button 
              size="sm" 
              variant="outline"
              onClick={() => setIsExpanded(!isExpanded)}
              data-testid="expand-suggestion"
            >
              {isExpanded ? 'Masquer' : 'Voir idées contenu'}
              {isExpanded ? <ChevronDown className="h-3 w-3 ml-1" /> : <ChevronRight className="h-3 w-3 ml-1" />}
            </Button>
            <Button size="sm" className="bg-purple-600 hover:bg-purple-700">
              Ajouter au planning
            </Button>
          </div>
        </div>
      </div>
      
      {isExpanded && (
        <div className="mt-4 pt-4 border-t" data-testid="content-ideas-list">
          <h5 className="font-medium mb-2 flex items-center space-x-1">
            <Lightbulb className="h-4 w-4 text-yellow-500" />
            <span>Idées de contenu</span>
          </h5>
          <div className="space-y-2">
            {suggestion.contentIdeas.map((idea, idx) => (
              <div key={idx} className="flex items-center space-x-2 p-2 bg-white rounded" data-testid="content-idea">
                <BookOpen className="h-3 w-3 text-blue-500" />
                <span className="text-sm">{idea}</span>
              </div>
            ))}
          </div>
          
          {suggestion.relatedKeywords && (
            <div className="mt-3">
              <h6 className="font-medium text-sm mb-1">Mots-clés associés :</h6>
              <div className="flex flex-wrap gap-1">
                {suggestion.relatedKeywords.map((kw, idx) => (
                  <Badge key={idx} variant="secondary" className="text-xs">
                    {kw}
                  </Badge>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

// Composant principal
export function KeywordAnalysisSection({ keywordData, analysisId }: KeywordAnalysisSectionProps) {
  const [expandedKeywords, setExpandedKeywords] = useState<Set<string>>(new Set());
  const [sortBy, setSortBy] = useState<KeywordSortBy>(KeywordSortBy.PROMINENCE);
  const [filterBy, setFilterBy] = useState<KeywordFilterBy>(KeywordFilterBy.ALL);
  const [searchQuery, setSearchQuery] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const toggleKeyword = (keyword: string) => {
    const newExpanded = new Set(expandedKeywords);
    if (newExpanded.has(keyword)) {
      newExpanded.delete(keyword);
    } else {
      newExpanded.add(keyword);
    }
    setExpandedKeywords(newExpanded);
  };

  // Filter and sort keywords
  const filteredKeywords = keywordData.foundKeywords
    .filter(kw => {
      if (searchQuery && !kw.keyword.toLowerCase().includes(searchQuery.toLowerCase())) {
        return false;
      }
      
      switch (filterBy) {
        case KeywordFilterBy.TITLE:
          return kw.locations.includes('title');
        case KeywordFilterBy.HEADINGS:
          return kw.locations.some(loc => ['h1', 'h2', 'h3'].includes(loc));
        case KeywordFilterBy.HIGH_DENSITY:
          return kw.density > 2;
        case KeywordFilterBy.LOW_DENSITY:
          return kw.density < 1;
        default:
          return true;
      }
    })
    .sort((a, b) => {
      switch (sortBy) {
        case KeywordSortBy.DENSITY:
          return b.density - a.density;
        case KeywordSortBy.OCCURRENCES:
          return b.occurrences - a.occurrences;
        case KeywordSortBy.ALPHABETICAL:
          return a.keyword.localeCompare(b.keyword);
        default:
          return b.prominence - a.prominence;
      }
    });

  const refreshAISuggestions = async () => {
    setIsLoading(true);
    // TODO: Implement refresh logic
    setTimeout(() => setIsLoading(false), 2000);
  };

  return (
    <div className="space-y-6" data-testid="keyword-analysis-section">
      {/* Header with metrics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card data-testid="analysis-card">
          <CardContent className="p-4">
            <div className="flex items-center space-x-2">
              <Tag className="h-5 w-5 text-blue-500" />
              <div>
                <div className="text-2xl font-bold">{keywordData.metrics.totalKeywords}</div>
                <div className="text-sm text-gray-600">Mots-clés trouvés</div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card data-testid="analysis-card">
          <CardContent className="p-4">
            <div className="flex items-center space-x-2">
              <BarChart3 className="h-5 w-5 text-green-500" />
              <div>
                <div className="text-2xl font-bold">{keywordData.metrics.keywordDiversity}</div>
                <div className="text-sm text-gray-600">Diversité (%)</div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card data-testid="analysis-card">
          <CardContent className="p-4">
            <div className="flex items-center space-x-2">
              <Target className="h-5 w-5 text-orange-500" />
              <div>
                <div className="text-2xl font-bold">{keywordData.aiSuggestions.length}</div>
                <div className="text-sm text-gray-600">Suggestions IA</div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card data-testid="analysis-card">
          <CardContent className="p-4">
            <div className="flex items-center space-x-2">
              <Activity className="h-5 w-5 text-purple-500" />
              <div>
                <div className="text-2xl font-bold">{keywordData.semanticAnalysis.wordCount}</div>
                <div className="text-sm text-gray-600">Mots analysés</div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Tabs pour différentes sections */}
      <Tabs defaultValue="found-keywords" data-testid="tab-keywords">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="found-keywords">Mots-clés trouvés</TabsTrigger>
          <TabsTrigger value="ngrams">N-grammes</TabsTrigger>
          <TabsTrigger value="ai-suggestions">Suggestions IA</TabsTrigger>
          <TabsTrigger value="semantic">Analyse sémantique</TabsTrigger>
        </TabsList>

        {/* Mots-clés trouvés */}
        <TabsContent value="found-keywords" className="space-y-4">
          {/* Contrôles */}
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div className="flex items-center space-x-2">
              <Input 
                placeholder="Rechercher un mot-clé..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-64"
                data-testid="filter-keywords"
              />
              <select 
                className="border rounded px-3 py-2"
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value as KeywordSortBy)}
                data-testid="sort-keywords"
              >
                <option value={KeywordSortBy.PROMINENCE}>Par prominence</option>
                <option value={KeywordSortBy.DENSITY}>Par densité</option>
                <option value={KeywordSortBy.OCCURRENCES}>Par occurrences</option>
                <option value={KeywordSortBy.ALPHABETICAL}>Alphabétique</option>
              </select>
            </div>
            
            <Button variant="outline" data-testid="export-keywords">
              <Download className="h-4 w-4 mr-2" />
              Exporter
            </Button>
          </div>

          {/* Liste des mots-clés */}
          <div className="space-y-3" data-testid="found-keywords-list">
            {filteredKeywords.map((keyword) => (
              <KeywordItem
                key={keyword.keyword}
                keyword={keyword}
                isExpanded={expandedKeywords.has(keyword.keyword)}
                onToggle={() => toggleKeyword(keyword.keyword)}
              />
            ))}
          </div>
        </TabsContent>

        {/* N-grammes */}
        <TabsContent value="ngrams" className="space-y-4" data-testid="ngrams-analysis">
          <Tabs defaultValue="bigrams">
            <TabsList>
              <TabsTrigger value="bigrams" data-testid="bigrams-tab">Bigrammes</TabsTrigger>
              <TabsTrigger value="trigrams" data-testid="trigrams-tab">Trigrammes</TabsTrigger>
            </TabsList>

            <TabsContent value="bigrams">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Top Bigrammes</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3" data-testid="bigrams-list">
                      {keywordData.ngrams.bigrams.slice(0, 10).map((bigram, idx) => (
                        <div key={idx} className="flex justify-between items-center" data-testid="ngram-item">
                          <span className="font-medium" data-testid="ngram-phrase">{bigram.phrase}</span>
                          <div className="text-right">
                            <div className="font-semibold" data-testid="ngram-count">{bigram.count}</div>
                            <div className="text-xs text-gray-600">{bigram.density.toFixed(2)}%</div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>

                {/* Graphique placeholder */}
                <Card>
                  <CardHeader>
                    <CardTitle>Visualisation</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="h-64 bg-gray-100 rounded flex items-center justify-center" data-testid="ngrams-chart">
                      <div className="text-gray-500">Graphique des n-grammes</div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>

            <TabsContent value="trigrams">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Top Trigrammes</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-3" data-testid="trigrams-list">
                      {keywordData.ngrams.trigrams.slice(0, 10).map((trigram, idx) => (
                        <div key={idx} className="flex justify-between items-center" data-testid="ngram-item">
                          <span className="font-medium" data-testid="ngram-phrase">{trigram.phrase}</span>
                          <div className="text-right">
                            <div className="font-semibold" data-testid="ngram-count">{trigram.count}</div>
                            <div className="text-xs text-gray-600">{trigram.density.toFixed(2)}%</div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </div>
            </TabsContent>
          </Tabs>
        </TabsContent>

        {/* Suggestions IA */}
        <TabsContent value="ai-suggestions" className="space-y-4" data-testid="ai-suggestions">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Brain className="h-5 w-5 text-purple-600" />
              <h3 className="text-lg font-semibold">Suggestions IA</h3>
              <Badge className="bg-purple-100 text-purple-700" data-testid="ai-badge">
                GPT-3.5
              </Badge>
            </div>
            
            <Button 
              variant="outline" 
              onClick={refreshAISuggestions}
              disabled={isLoading}
              data-testid="refresh-ai-suggestions"
            >
              <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
              {isLoading ? 'Génération...' : 'Actualiser'}
            </Button>
          </div>

          {isLoading ? (
            <div className="text-center py-8" data-testid="ai-loading">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-purple-600 mx-auto mb-4"></div>
              <p>Génération de nouvelles suggestions IA...</p>
            </div>
          ) : (
            <div className="space-y-4" data-testid="ai-suggestions-list">
              {keywordData.aiSuggestions.map((suggestion, idx) => (
                <AIKeywordSuggestion key={idx} suggestion={suggestion} />
              ))}
            </div>
          )}
        </TabsContent>

        {/* Analyse sémantique */}
        <TabsContent value="semantic" className="space-y-4" data-testid="semantic-analysis">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Topics et entités */}
            <Card>
              <CardHeader>
                <CardTitle>Topics & Entités</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div>
                    <h4 className="font-medium mb-2">Sujets principaux</h4>
                    <div className="flex flex-wrap gap-2" data-testid="main-topics">
                      {keywordData.semanticAnalysis.mainTopics.map((topic, idx) => (
                        <Badge key={idx} className="bg-blue-100 text-blue-800" data-testid="topic-tag">
                          {topic}
                        </Badge>
                      ))}
                    </div>
                  </div>
                  
                  <div>
                    <h4 className="font-medium mb-2">Entités détectées</h4>
                    <div className="space-y-2" data-testid="entities-list">
                      {keywordData.semanticAnalysis.entities.slice(0, 8).map((entity, idx) => (
                        <div key={idx} className="flex justify-between items-center" data-testid="entity-item">
                          <span>{entity.name}</span>
                          <Badge variant="outline">{entity.type}</Badge>
                        </div>
                      ))}
                    </div>
                  </div>
                  
                  <div>
                    <h4 className="font-medium mb-2">Sentiment</h4>
                    <Badge 
                      className={`${
                        keywordData.semanticAnalysis.sentiment === 'positive' ? 'bg-green-100 text-green-800' :
                        keywordData.semanticAnalysis.sentiment === 'negative' ? 'bg-red-100 text-red-800' :
                        'bg-gray-100 text-gray-800'
                      }`}
                      data-testid="sentiment-badge"
                    >
                      {keywordData.semanticAnalysis.sentiment === 'positive' ? 'Positive' :
                       keywordData.semanticAnalysis.sentiment === 'negative' ? 'Negative' : 'Neutral'}
                    </Badge>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Lisibilité */}
            <Card>
              <CardHeader>
                <CardTitle>Lisibilité</CardTitle>
              </CardHeader>
              <CardContent data-testid="readability-metrics">
                <div className="space-y-4">
                  <div className="text-center">
                    <div className="text-3xl font-bold mb-2" data-testid="readability-score">
                      {Math.round(keywordData.semanticAnalysis.readability.score)}
                    </div>
                    <div className="text-sm text-gray-600" data-testid="readability-level">
                      Niveau : {keywordData.semanticAnalysis.readability.level}
                    </div>
                    <Progress 
                      value={keywordData.semanticAnalysis.readability.score} 
                      className="mt-2"
                      data-testid="readability-gauge"
                    />
                  </div>
                  
                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <div className="text-gray-600">Mots par phrase</div>
                      <div className="font-semibold" data-testid="avg-sentence-length">
                        {keywordData.semanticAnalysis.readability.avgSentenceLength.toFixed(1)}
                      </div>
                    </div>
                    <div>
                      <div className="text-gray-600">Taille mots</div>
                      <div className="font-semibold" data-testid="avg-word-length">
                        {keywordData.semanticAnalysis.readability.avgWordLength.toFixed(1)}
                      </div>
                    </div>
                    <div>
                      <div className="text-gray-600">Mots complexes</div>
                      <div className="font-semibold">
                        {keywordData.semanticAnalysis.readability.complexWords.toFixed(1)}%
                      </div>
                    </div>
                    <div>
                      <div className="text-gray-600">Paragraphes</div>
                      <div className="font-semibold">
                        {keywordData.semanticAnalysis.readability.paragraphs}
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}