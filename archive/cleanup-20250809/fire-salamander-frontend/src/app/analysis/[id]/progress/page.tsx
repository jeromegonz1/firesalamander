"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { 
  CheckCircle2, 
  Clock, 
  AlertCircle,
  ArrowRight,
  RefreshCw,
  BarChart3,
  Globe,
  Search,
  Zap
} from "lucide-react";

interface AnalysisStatus {
  task_id: string;
  url: string;
  status: "pending" | "running" | "completed" | "failed";
  progress: number;
  current_step: string;
  steps: {
    name: string;
    status: "pending" | "running" | "completed" | "failed";
    duration?: number;
  }[];
  start_time: string;
  estimated_completion?: string;
  error?: string;
}

const analysisSteps = [
  { key: "extraction", name: "Extraction du contenu", icon: Globe },
  { key: "crawling", name: "Crawling des pages", icon: Search },
  { key: "seo_analysis", name: "Analyse SEO", icon: BarChart3 },
  { key: "performance", name: "Test de performance", icon: Zap },
  { key: "ai_insights", name: "Génération des insights IA", icon: RefreshCw },
];

export default function AnalysisProgressPage() {
  const params = useParams();
  const router = useRouter();
  const [status, setStatus] = useState<AnalysisStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const taskId = params.id as string;

  useEffect(() => {
    const fetchStatus = async () => {
      try {
        // Appel vers l'API Fire Salamander pour récupérer le statut
        const response = await fetch(`http://localhost:8080/api/v1/analysis/${taskId}/status`);
        
        if (response.ok) {
          const data = await response.json();
          setStatus(data.data);
          
          // Si l'analyse est terminée, redirection vers le rapport
          if (data.data.status === "completed") {
            setTimeout(() => {
              router.push(`/analysis/${taskId}/report`);
            }, 2000);
          }
        } else {
          // Fallback: essayer de récupérer l'analyse terminée
          const analysisResponse = await fetch(`http://localhost:8080/api/v1/analysis/${taskId}`);
          if (analysisResponse.ok) {
            router.push(`/analysis/${taskId}/report`);
          } else {
            setError("Analyse non trouvée");
          }
        }
      } catch (err) {
        console.error("Erreur lors de la récupération du statut:", err);
        
        // Simulation de statut pour la demo
        setStatus({
          task_id: taskId,
          url: "https://www.marina-plage.com/",
          status: "running",
          progress: 75,
          current_step: "seo_analysis",
          steps: [
            { name: "Extraction du contenu", status: "completed", duration: 1200 },
            { name: "Crawling des pages", status: "completed", duration: 2300 },
            { name: "Analyse SEO", status: "running" },
            { name: "Test de performance", status: "pending" },
            { name: "Génération des insights IA", status: "pending" },
          ],
          start_time: new Date().toISOString(),
          estimated_completion: new Date(Date.now() + 30000).toISOString()
        });
      } finally {
        setLoading(false);
      }
    };

    fetchStatus();
    
    // Polling toutes les 2 secondes si l'analyse est en cours
    const interval = setInterval(() => {
      if (status?.status === "running") {
        fetchStatus();
      }
    }, 2000);

    return () => clearInterval(interval);
  }, [taskId, router, status?.status]);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Chargement du statut...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-2xl mx-auto">
        <Card className="border-red-200">
          <CardContent className="p-8 text-center">
            <AlertCircle className="h-12 w-12 text-red-500 mx-auto mb-4" />
            <h2 className="text-xl font-semibold mb-2">Erreur</h2>
            <p className="text-gray-600 mb-4">{error}</p>
            <Button onClick={() => router.push("/dashboard")}>
              Retour au Dashboard
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (!status) return null;

  const getStatusColor = (stepStatus: string) => {
    switch (stepStatus) {
      case "completed": return "text-green-600 bg-green-50 border-green-200";
      case "running": return "text-blue-600 bg-blue-50 border-blue-200";
      case "failed": return "text-red-600 bg-red-50 border-red-200";
      default: return "text-gray-600 bg-gray-50 border-gray-200";
    }
  };

  const getStatusIcon = (stepStatus: string) => {
    switch (stepStatus) {
      case "completed": return <CheckCircle2 className="h-5 w-5 text-green-600" />;
      case "running": return <RefreshCw className="h-5 w-5 text-blue-600 animate-spin" />;
      case "failed": return <AlertCircle className="h-5 w-5 text-red-600" />;
      default: return <Clock className="h-5 w-5 text-gray-400" />;
    }
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">
          Analyse en cours
        </h1>
        <p className="text-lg text-gray-600">
          {status.url}
        </p>
      </div>

      {/* Progress Overview */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="flex items-center space-x-2">
                <RefreshCw className={`h-5 w-5 ${status.status === "running" ? "animate-spin text-blue-600" : "text-gray-600"}`} />
                <span>Progression de l&apos;analyse</span>
              </CardTitle>
              <CardDescription>
                Task ID: {status.task_id}
              </CardDescription>
            </div>
            <Badge 
              variant={status.status === "completed" ? "default" : "secondary"}
              className={status.status === "running" ? "bg-blue-100 text-blue-700" : ""}
            >
              {status.status === "running" && "En cours"}
              {status.status === "completed" && "Terminé"}
              {status.status === "failed" && "Échoué"}
              {status.status === "pending" && "En attente"}
            </Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between text-sm text-gray-600">
              <span>Progression globale</span>
              <span>{status.progress}%</span>
            </div>
            <Progress value={status.progress} className="h-2" />
            
            {status.estimated_completion && (
              <div className="text-sm text-gray-500">
                Estimation de fin: {new Date(status.estimated_completion).toLocaleTimeString()}
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Steps Detail */}
      <Card>
        <CardHeader>
          <CardTitle>Étapes de l&apos;analyse</CardTitle>
          <CardDescription>
            Suivi détaillé de chaque étape du processus d&apos;analyse
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {status.steps.map((step, index) => (
              <div key={index} className={`p-4 rounded-lg border transition-all ${getStatusColor(step.status)}`}>
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    {getStatusIcon(step.status)}
                    <div>
                      <h4 className="font-medium">{step.name}</h4>
                      {step.duration && (
                        <p className="text-sm opacity-80">
                          Durée: {(step.duration / 1000).toFixed(1)}s
                        </p>
                      )}
                    </div>
                  </div>
                  
                  {step.status === "running" && (
                    <div className="flex items-center space-x-2">
                      <div className="animate-pulse h-2 w-2 rounded-full bg-current"></div>
                      <span className="text-sm">En cours...</span>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Actions */}
      {status.status === "completed" && (
        <div className="text-center">
          <Button 
            size="lg" 
            onClick={() => router.push(`/analysis/${taskId}/report`)}
            className="bg-orange-500 hover:bg-orange-600"
          >
            Voir le rapport complet
            <ArrowRight className="h-5 w-5 ml-2" />
          </Button>
        </div>
      )}
      
      {status.status === "failed" && (
        <Card className="border-red-200">
          <CardContent className="p-6 text-center">
            <AlertCircle className="h-8 w-8 text-red-500 mx-auto mb-4" />
            <h3 className="text-lg font-semibold mb-2">Analyse échouée</h3>
            <p className="text-gray-600 mb-4">
              {status.error || "Une erreur est survenue pendant l'analyse"}
            </p>
            <div className="flex justify-center space-x-4">
              <Button variant="outline" onClick={() => router.push("/analysis/new")}>
                Nouvelle analyse
              </Button>
              <Button onClick={() => router.push("/dashboard")}>
                Retour au Dashboard
              </Button>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}