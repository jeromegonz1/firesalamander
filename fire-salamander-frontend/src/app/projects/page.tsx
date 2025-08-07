"use client";

import { useState, useEffect } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { 
  Plus, 
  Search, 
  Filter,
  Globe,
  BarChart3,
  Calendar,
  TrendingUp,
  TrendingDown,
  Minus,
  ArrowRight,
  Eye
} from "lucide-react";
import Link from "next/link";

interface ProjectAnalysis {
  id: string;
  url: string;
  domain: string;
  analyzed_at: string;
  overall_score: number;
  status: "completed" | "running" | "failed";
  processing_time?: number;
  scores?: {
    seo: number;
    technical: number; 
    performance: number;
    content: number;
  };
}

interface Project {
  domain: string;
  url: string;
  analyses: ProjectAnalysis[];
  latest_score: number;
  trend: "up" | "down" | "stable";
  last_analysis: string;
  total_analyses: number;
}

function ProjectCard({ project }: { project: Project }) {
  const getTrendIcon = () => {
    switch (project.trend) {
      case "up": return <TrendingUp className="h-4 w-4 text-green-600" />;
      case "down": return <TrendingDown className="h-4 w-4 text-red-600" />;
      default: return <Minus className="h-4 w-4 text-gray-600" />;
    }
  };

  const getScoreColor = (score: number) => {
    if (score >= 80) return "text-green-600 bg-green-50 border-green-200";
    if (score >= 60) return "text-orange-600 bg-orange-50 border-orange-200";
    return "text-red-600 bg-red-50 border-red-200";
  };

  const latestAnalysis = project.analyses[0];

  return (
    <Card className="hover:shadow-lg transition-shadow cursor-pointer">
      <CardContent className="p-6">
        <div className="flex items-start justify-between mb-4">
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-blue-50 rounded-lg">
              <Globe className="h-5 w-5 text-blue-600" />
            </div>
            <div>
              <h3 className="font-semibold text-lg">{project.domain}</h3>
              <p className="text-sm text-gray-600">{project.url}</p>
            </div>
          </div>
          
          <div className="flex items-center space-x-2">
            {getTrendIcon()}
            <Badge className={getScoreColor(project.latest_score)}>
              {project.latest_score}/100
            </Badge>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4 mb-4">
          <div>
            <p className="text-sm text-gray-600">Dernière analyse</p>
            <p className="font-medium">{new Date(project.last_analysis).toLocaleDateString()}</p>
          </div>
          <div>
            <p className="text-sm text-gray-600">Total analyses</p>
            <p className="font-medium">{project.total_analyses}</p>
          </div>
        </div>

        {latestAnalysis?.scores && (
          <div className="grid grid-cols-4 gap-2 mb-4">
            <div className="text-center p-2 bg-blue-50 rounded">
              <div className="text-sm font-semibold text-blue-600">{latestAnalysis.scores.seo}</div>
              <div className="text-xs text-gray-600">SEO</div>
            </div>
            <div className="text-center p-2 bg-purple-50 rounded">
              <div className="text-sm font-semibold text-purple-600">{latestAnalysis.scores.technical}</div>
              <div className="text-xs text-gray-600">Tech</div>
            </div>
            <div className="text-center p-2 bg-green-50 rounded">
              <div className="text-sm font-semibold text-green-600">{latestAnalysis.scores.performance}</div>
              <div className="text-xs text-gray-600">Perf</div>
            </div>
            <div className="text-center p-2 bg-orange-50 rounded">
              <div className="text-sm font-semibold text-orange-600">{latestAnalysis.scores.content}</div>
              <div className="text-xs text-gray-600">Contenu</div>
            </div>
          </div>
        )}

        <div className="flex justify-between items-center">
          <div className="flex space-x-2">
            <Link href={`/analysis/${latestAnalysis?.id}/report`}>
              <Button variant="outline" size="sm">
                <Eye className="h-4 w-4 mr-1" />
                Voir
              </Button>
            </Link>
            <Link href="/analysis/new">
              <Button variant="outline" size="sm">
                <BarChart3 className="h-4 w-4 mr-1" />
                Analyser
              </Button>
            </Link>
          </div>
          
          <Link href={`/analysis/${latestAnalysis?.id}/report`}>
            <Button variant="ghost" size="sm">
              <ArrowRight className="h-4 w-4" />
            </Button>
          </Link>
        </div>
      </CardContent>
    </Card>
  );
}

function AnalysisItem({ analysis }: { analysis: ProjectAnalysis }) {
  const getStatusColor = (status: string) => {
    switch (status) {
      case "completed": return "bg-green-100 text-green-700";
      case "running": return "bg-blue-100 text-blue-700";
      case "failed": return "bg-red-100 text-red-700";
      default: return "bg-gray-100 text-gray-700";
    }
  };

  const getScoreColor = (score: number) => {
    if (score >= 80) return "text-green-600";
    if (score >= 60) return "text-orange-600";
    return "text-red-600";
  };

  return (
    <div className="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50 transition-colors">
      <div className="flex items-center space-x-4">
        <div className="text-center">
          <div className={`text-2xl font-bold ${getScoreColor(analysis.overall_score)}`}>
            {analysis.overall_score}
          </div>
          <div className="text-xs text-gray-500">Score</div>
        </div>
        
        <div>
          <p className="font-medium">{analysis.url}</p>
          <div className="flex items-center space-x-4 text-sm text-gray-600">
            <span className="flex items-center space-x-1">
              <Calendar className="h-3 w-3" />
              <span>{new Date(analysis.analyzed_at).toLocaleDateString()}</span>
            </span>
            {analysis.processing_time && (
              <span>Durée: {(analysis.processing_time / 1000).toFixed(1)}s</span>
            )}
          </div>
        </div>
      </div>
      
      <div className="flex items-center space-x-3">
        <Badge className={getStatusColor(analysis.status)}>
          {analysis.status === "completed" && "Terminé"}
          {analysis.status === "running" && "En cours"}
          {analysis.status === "failed" && "Échoué"}
        </Badge>
        
        <Link href={`/analysis/${analysis.id}/report`}>
          <Button variant="ghost" size="sm">
            <ArrowRight className="h-4 w-4" />
          </Button>
        </Link>
      </div>
    </div>
  );
}

export default function ProjectsPage() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [allAnalyses, setAllAnalyses] = useState<ProjectAnalysis[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState("");
  const [view, setView] = useState<"projects" | "analyses">("projects");

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await fetch("http://localhost:8080/api/v1/analyses");
        
        if (response.ok) {
          const data = await response.json();
          const analyses: ProjectAnalysis[] = data.data || [];
          
          setAllAnalyses(analyses);
          
          // Grouper par domaine pour créer les projets
          const projectsMap = new Map<string, Project>();
          
          analyses.forEach(analysis => {
            const domain = analysis.domain || new URL(analysis.url).hostname;
            
            if (!projectsMap.has(domain)) {
              projectsMap.set(domain, {
                domain,
                url: analysis.url,
                analyses: [],
                latest_score: analysis.overall_score,
                trend: "stable",
                last_analysis: analysis.analyzed_at,
                total_analyses: 0
              });
            }
            
            const project = projectsMap.get(domain)!;
            project.analyses.push(analysis);
            project.total_analyses = project.analyses.length;
            
            // Calculer la tendance basée sur les 2 dernières analyses
            if (project.analyses.length >= 2) {
              const current = project.analyses[0].overall_score;
              const previous = project.analyses[1].overall_score;
              
              if (current > previous + 5) project.trend = "up";
              else if (current < previous - 5) project.trend = "down";
              else project.trend = "stable";
            }
          });
          
          setProjects(Array.from(projectsMap.values()));
        } else {
          // Données mockées pour la demo
          const mockProjects: Project[] = [
            {
              domain: "marina-plage.com",
              url: "https://www.marina-plage.com/",
              analyses: [
                {
                  id: "task_123",
                  url: "https://www.marina-plage.com/",
                  domain: "marina-plage.com",
                  analyzed_at: new Date(Date.now() - 86400000).toISOString(),
                  overall_score: 78,
                  status: "completed",
                  processing_time: 4500,
                  scores: { seo: 82, technical: 71, performance: 84, content: 90 }
                }
              ],
              latest_score: 78,
              trend: "up",
              last_analysis: new Date(Date.now() - 86400000).toISOString(),
              total_analyses: 5
            },
            {
              domain: "example-business.fr",
              url: "https://www.example-business.fr/",
              analyses: [
                {
                  id: "task_456",
                  url: "https://www.example-business.fr/",
                  domain: "example-business.fr",
                  analyzed_at: new Date(Date.now() - 172800000).toISOString(),
                  overall_score: 65,
                  status: "completed",
                  processing_time: 3200,
                  scores: { seo: 70, technical: 60, performance: 68, content: 62 }
                }
              ],
              latest_score: 65,
              trend: "down",
              last_analysis: new Date(Date.now() - 172800000).toISOString(),
              total_analyses: 3
            }
          ];
          
          setProjects(mockProjects);
          setAllAnalyses(mockProjects.flatMap(p => p.analyses));
        }
      } catch (error) {
        console.error("Erreur lors du chargement des projets:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const filteredProjects = projects.filter(project =>
    project.domain.toLowerCase().includes(searchTerm.toLowerCase()) ||
    project.url.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const filteredAnalyses = allAnalyses.filter(analysis =>
    analysis.url.toLowerCase().includes(searchTerm.toLowerCase()) ||
    analysis.domain?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Chargement des projets...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Mes Projets</h1>
          <p className="text-gray-600 mt-1">
            Gérez et suivez vos analyses SEO par projet
          </p>
        </div>
        
        <Link href="/analysis/new">
          <Button className="bg-orange-500 hover:bg-orange-600">
            <Plus className="h-4 w-4 mr-2" />
            Nouvelle Analyse
          </Button>
        </Link>
      </div>

      {/* Filters & Search */}
      <Card>
        <CardContent className="p-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="Rechercher un projet ou domaine..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10 w-80"
                />
              </div>
              
              <Button variant="outline" size="sm">
                <Filter className="h-4 w-4 mr-2" />
                Filtres
              </Button>
            </div>
            
            <div className="flex space-x-2">
              <Button 
                variant={view === "projects" ? "default" : "outline"}
                size="sm"
                onClick={() => setView("projects")}
              >
                Par projets
              </Button>
              <Button 
                variant={view === "analyses" ? "default" : "outline"}
                size="sm"
                onClick={() => setView("analyses")}
              >
                Toutes les analyses
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Content */}
      {view === "projects" ? (
        <div>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold">
              {filteredProjects.length} projet{filteredProjects.length > 1 ? "s" : ""}
            </h2>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredProjects.map((project, index) => (
              <ProjectCard key={index} project={project} />
            ))}
          </div>
          
          {filteredProjects.length === 0 && (
            <Card>
              <CardContent className="p-12 text-center">
                <Globe className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">Aucun projet trouvé</h3>
                <p className="text-gray-600 mb-4">
                  {searchTerm ? "Aucun projet ne correspond à votre recherche" : "Vous n'avez pas encore de projets"}
                </p>
                <Link href="/analysis/new">
                  <Button className="bg-orange-500 hover:bg-orange-600">
                    <Plus className="h-4 w-4 mr-2" />
                    Créer votre première analyse
                  </Button>
                </Link>
              </CardContent>
            </Card>
          )}
        </div>
      ) : (
        <div>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold">
              {filteredAnalyses.length} analyse{filteredAnalyses.length > 1 ? "s" : ""}
            </h2>
          </div>
          
          <div className="space-y-4">
            {filteredAnalyses.map((analysis) => (
              <AnalysisItem key={analysis.id} analysis={analysis} />
            ))}
          </div>
          
          {filteredAnalyses.length === 0 && (
            <Card>
              <CardContent className="p-12 text-center">
                <BarChart3 className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">Aucune analyse trouvée</h3>
                <p className="text-gray-600 mb-4">
                  {searchTerm ? "Aucune analyse ne correspond à votre recherche" : "Vous n'avez pas encore d'analyses"}
                </p>
                <Link href="/analysis/new">
                  <Button className="bg-orange-500 hover:bg-orange-600">
                    <Plus className="h-4 w-4 mr-2" />
                    Lancer votre première analyse
                  </Button>
                </Link>
              </CardContent>
            </Card>
          )}
        </div>
      )}
    </div>
  );
}