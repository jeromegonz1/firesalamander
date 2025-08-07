/**
 * Fire Salamander - Analysis Layout Component
 * Lead Tech quality layout wrapper for analysis pages
 */

'use client';

import React, { useState, useEffect } from 'react';
import { usePathname, useParams } from 'next/navigation';
import { AnalysisNavigation } from './analysis-navigation';
import { Breadcrumbs, generateBreadcrumbs } from './breadcrumbs';

interface AnalysisLayoutProps {
  children: React.ReactNode;
}

export function AnalysisLayout({ children }: AnalysisLayoutProps) {
  const pathname = usePathname();
  const params = useParams();
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [analysisUrl, setAnalysisUrl] = useState<string>();

  const analysisId = params.id as string;

  // Mock analysis URL fetch - in real app this would come from API
  useEffect(() => {
    if (analysisId) {
      // Simulate API call to get analysis details
      setAnalysisUrl('https://example.com');
    }
  }, [analysisId]);

  const breadcrumbs = generateBreadcrumbs(pathname, analysisId, analysisUrl);

  const toggleSidebar = () => {
    setSidebarCollapsed(!sidebarCollapsed);
  };

  return (
    <div className="min-h-screen bg-gray-50 flex">
      {/* Sidebar Navigation */}
      <div className="shrink-0">
        <AnalysisNavigation
          analysisId={analysisId}
          analysisUrl={analysisUrl}
          collapsed={sidebarCollapsed}
          onToggleCollapsed={toggleSidebar}
        />
      </div>

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <header className="bg-white border-b border-gray-200 px-6 py-4">
          <div className="flex items-center justify-between">
            <Breadcrumbs items={breadcrumbs} />
            <div className="flex items-center space-x-4">
              {/* Analysis Status Indicator */}
              {analysisId && (
                <div className="flex items-center space-x-2">
                  <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                  <span className="text-sm text-gray-600">Analyse complète</span>
                </div>
              )}
            </div>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-6">
          <div className="max-w-7xl mx-auto">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}