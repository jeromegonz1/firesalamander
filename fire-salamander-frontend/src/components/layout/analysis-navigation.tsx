/**
 * Fire Salamander - Analysis Navigation Component
 * Lead Tech quality navigation system for analysis modules
 */

'use client';

import React, { useState, useMemo } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { 
  AnalysisNavigationItem, 
  AnalysisModuleType, 
  AnalysisModuleStatus, 
  BreadcrumbItem 
} from '@/types/navigation';
import { Card, CardContent } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { Button } from '@/components/ui/button';
import {
  BarChart3,
  Wrench,
  Zap,
  FileText,
  Shield,
  Link2,
  ChevronLeft,
  ChevronRight,
  Home,
  ChevronDown,
  CheckCircle,
  Clock,
  AlertCircle,
  XCircle,
  Minus,
} from 'lucide-react';

interface AnalysisNavigationProps {
  analysisId: string;
  analysisUrl?: string;
  collapsed?: boolean;
  onToggleCollapsed?: () => void;
}

export function AnalysisNavigation({ 
  analysisId, 
  analysisUrl = 'https://example.com',
  collapsed = false,
  onToggleCollapsed 
}: AnalysisNavigationProps) {
  const pathname = usePathname();
  const [expandedSections, setExpandedSections] = useState<Record<string, boolean>>({});

  // Mock analysis progress - in real app, this would come from props/context
  const analysisProgress = useMemo(() => ({
    overall: 78,
    modules: {
      [AnalysisModuleType.OVERVIEW]: {
        progress: 100,
        status: AnalysisModuleStatus.COMPLETED,
        score: 78,
        grade: 'B',
      },
      [AnalysisModuleType.TECHNICAL]: {
        progress: 100,
        status: AnalysisModuleStatus.COMPLETED,
        score: 85,
        grade: 'A',
      },
      [AnalysisModuleType.PERFORMANCE]: {
        progress: 100,
        status: AnalysisModuleStatus.COMPLETED,
        score: 72,
        grade: 'B',
      },
      [AnalysisModuleType.CONTENT]: {
        progress: 100,
        status: AnalysisModuleStatus.COMPLETED,
        score: 69,
        grade: 'C+',
      },
      [AnalysisModuleType.SECURITY]: {
        progress: 100,
        status: AnalysisModuleStatus.COMPLETED,
        score: 78,
        grade: 'B',
      },
      [AnalysisModuleType.BACKLINKS]: {
        progress: 100,
        status: AnalysisModuleStatus.COMPLETED,
        score: 72,
        grade: 'B',
      },
    },
  }), []);

  const navigationItems: AnalysisNavigationItem[] = useMemo(() => [
    {
      id: 'overview',
      label: 'Vue d\'ensemble',
      href: `/analysis/${analysisId}`,
      icon: 'BarChart3',
      moduleType: AnalysisModuleType.OVERVIEW,
      status: analysisProgress.modules[AnalysisModuleType.OVERVIEW].status,
      score: analysisProgress.modules[AnalysisModuleType.OVERVIEW].score,
      grade: analysisProgress.modules[AnalysisModuleType.OVERVIEW].grade,
    },
    {
      id: 'technical',
      label: 'Technique',
      href: `/analysis/${analysisId}/technical`,
      icon: 'Wrench',
      moduleType: AnalysisModuleType.TECHNICAL,
      status: analysisProgress.modules[AnalysisModuleType.TECHNICAL].status,
      score: analysisProgress.modules[AnalysisModuleType.TECHNICAL].score,
      grade: analysisProgress.modules[AnalysisModuleType.TECHNICAL].grade,
    },
    {
      id: 'performance',
      label: 'Performance',
      href: `/analysis/${analysisId}/performance`,
      icon: 'Zap',
      moduleType: AnalysisModuleType.PERFORMANCE,
      status: analysisProgress.modules[AnalysisModuleType.PERFORMANCE].status,
      score: analysisProgress.modules[AnalysisModuleType.PERFORMANCE].score,
      grade: analysisProgress.modules[AnalysisModuleType.PERFORMANCE].grade,
    },
    {
      id: 'content',
      label: 'Contenu',
      href: `/analysis/${analysisId}/content`,
      icon: 'FileText',
      moduleType: AnalysisModuleType.CONTENT,
      status: analysisProgress.modules[AnalysisModuleType.CONTENT].status,
      score: analysisProgress.modules[AnalysisModuleType.CONTENT].score,
      grade: analysisProgress.modules[AnalysisModuleType.CONTENT].grade,
    },
    {
      id: 'security',
      label: 'Sécurité',
      href: `/analysis/${analysisId}/security`,
      icon: 'Shield',
      moduleType: AnalysisModuleType.SECURITY,
      status: analysisProgress.modules[AnalysisModuleType.SECURITY].status,
      score: analysisProgress.modules[AnalysisModuleType.SECURITY].score,
      grade: analysisProgress.modules[AnalysisModuleType.SECURITY].grade,
    },
    {
      id: 'backlinks',
      label: 'Backlinks',
      href: `/analysis/${analysisId}/backlinks`,
      icon: 'Link2',
      moduleType: AnalysisModuleType.BACKLINKS,
      status: analysisProgress.modules[AnalysisModuleType.BACKLINKS].status,
      score: analysisProgress.modules[AnalysisModuleType.BACKLINKS].score,
      grade: analysisProgress.modules[AnalysisModuleType.BACKLINKS].grade,
    },
  ], [analysisId, analysisProgress]);

  const getIcon = (iconName: string) => {
    const icons = {
      BarChart3,
      Wrench,
      Zap,
      FileText,
      Shield,
      Link2,
    };
    const IconComponent = icons[iconName as keyof typeof icons] || BarChart3;
    return <IconComponent className="h-5 w-5" />;
  };

  const getStatusIcon = (status: AnalysisModuleStatus) => {
    switch (status) {
      case AnalysisModuleStatus.COMPLETED:
        return <CheckCircle className="h-4 w-4 text-green-500" />;
      case AnalysisModuleStatus.IN_PROGRESS:
        return <Clock className="h-4 w-4 text-blue-500 animate-pulse" />;
      case AnalysisModuleStatus.ERROR:
        return <XCircle className="h-4 w-4 text-red-500" />;
      case AnalysisModuleStatus.PENDING:
        return <AlertCircle className="h-4 w-4 text-orange-500" />;
      default:
        return <Minus className="h-4 w-4 text-gray-400" />;
    }
  };

  const getGradeColor = (grade?: string) => {
    if (!grade) return 'text-gray-500';
    const letter = grade.charAt(0);
    switch (letter) {
      case 'A': return 'text-green-600';
      case 'B': return 'text-blue-600';
      case 'C': return 'text-yellow-600';
      case 'D': return 'text-orange-600';
      case 'F': return 'text-red-600';
      default: return 'text-gray-600';
    }
  };

  return (
    <div className={`bg-white border-r border-gray-200 h-full flex flex-col transition-all duration-300 ${
      collapsed ? 'w-16' : 'w-80'
    }`}>
      {/* Header */}
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center justify-between">
          {!collapsed && (
            <div className="flex-1 min-w-0">
              <h2 className="text-lg font-semibold text-gray-900 truncate">
                Fire Salamander
              </h2>
              <div className="text-sm text-gray-600 truncate" title={analysisUrl}>
                {analysisUrl}
              </div>
            </div>
          )}
          <Button
            variant="ghost"
            size="sm"
            onClick={onToggleCollapsed}
            className="shrink-0"
          >
            {collapsed ? <ChevronRight className="h-4 w-4" /> : <ChevronLeft className="h-4 w-4" />}
          </Button>
        </div>
      </div>

      {/* Overall Progress */}
      {!collapsed && (
        <div className="p-4 border-b border-gray-200">
          <Card>
            <CardContent className="p-4">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium">Progression globale</span>
                <span className="text-lg font-bold text-blue-600">{analysisProgress.overall}%</span>
              </div>
              <Progress value={analysisProgress.overall} className="h-2" />
              <div className="text-xs text-gray-600 mt-2">
                6/6 modules analysés
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Navigation Items */}
      <div className="flex-1 overflow-y-auto">
        <nav className="p-2">
          <div className="space-y-1">
            {/* Home Link */}
            <Link href="/dashboard">
              <div className={`flex items-center space-x-3 px-3 py-2 rounded-lg hover:bg-gray-50 transition-colors ${
                pathname === '/dashboard' ? 'bg-blue-50 text-blue-700' : 'text-gray-700'
              }`}>
                <Home className="h-5 w-5 shrink-0" />
                {!collapsed && <span className="font-medium">Dashboard</span>}
              </div>
            </Link>

            {/* Analysis Modules */}
            {navigationItems.map((item) => {
              const isActive = pathname === item.href;
              const progress = analysisProgress.modules[item.moduleType];
              
              return (
                <Link key={item.id} href={item.href}>
                  <div className={`flex items-center justify-between px-3 py-2 rounded-lg hover:bg-gray-50 transition-colors ${
                    isActive ? 'bg-blue-50 text-blue-700' : 'text-gray-700'
                  }`}>
                    <div className="flex items-center space-x-3 flex-1 min-w-0">
                      <div className="shrink-0">
                        {getIcon(item.icon)}
                      </div>
                      {!collapsed && (
                        <>
                          <span className="font-medium truncate">{item.label}</span>
                          <div className="flex items-center space-x-1 ml-auto">
                            {getStatusIcon(item.status)}
                          </div>
                        </>
                      )}
                    </div>
                    
                    {!collapsed && item.grade && (
                      <div className={`text-sm font-bold ${getGradeColor(item.grade)}`}>
                        {item.grade}
                      </div>
                    )}
                  </div>
                </Link>
              );
            })}
          </div>
        </nav>
      </div>

      {/* Footer Actions */}
      {!collapsed && (
        <div className="p-4 border-t border-gray-200">
          <div className="space-y-2">
            <Link href="/analysis/new">
              <Button className="w-full" size="sm">
                Nouvelle analyse
              </Button>
            </Link>
            <Link href={`/analysis/${analysisId}/report`}>
              <Button variant="outline" className="w-full" size="sm">
                Rapport complet
              </Button>
            </Link>
          </div>
        </div>
      )}
    </div>
  );
}