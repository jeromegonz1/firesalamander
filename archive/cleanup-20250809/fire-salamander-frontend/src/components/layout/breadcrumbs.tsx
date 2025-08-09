/**
 * Fire Salamander - Breadcrumbs Component
 * Lead Tech quality breadcrumb navigation
 */

'use client';

import React from 'react';
import Link from 'next/link';
import { BreadcrumbItem } from '@/types/navigation';
import { ChevronRight, Home } from 'lucide-react';

interface BreadcrumbsProps {
  items: BreadcrumbItem[];
  className?: string;
}

export function Breadcrumbs({ items, className = '' }: BreadcrumbsProps) {
  if (items.length === 0) return null;

  return (
    <nav className={`flex items-center space-x-1 text-sm ${className}`} aria-label="Breadcrumb">
      <Link 
        href="/dashboard" 
        className="flex items-center text-gray-500 hover:text-gray-700 transition-colors"
      >
        <Home className="h-4 w-4" />
      </Link>
      
      {items.map((item, index) => (
        <React.Fragment key={index}>
          <ChevronRight className="h-4 w-4 text-gray-400" />
          {item.href && !item.isActive ? (
            <Link 
              href={item.href}
              className="text-gray-500 hover:text-gray-700 transition-colors"
            >
              {item.label}
            </Link>
          ) : (
            <span className={`${item.isActive ? 'text-gray-900 font-medium' : 'text-gray-500'}`}>
              {item.label}
            </span>
          )}
        </React.Fragment>
      ))}
    </nav>
  );
}

// Utility function to generate breadcrumbs from pathname
export function generateBreadcrumbs(pathname: string, analysisId?: string, analysisUrl?: string): BreadcrumbItem[] {
  const segments = pathname.split('/').filter(Boolean);
  const breadcrumbs: BreadcrumbItem[] = [];

  if (segments[0] === 'analysis' && segments[1] && analysisId) {
    breadcrumbs.push({
      label: analysisUrl ? `Analyse: ${analysisUrl}` : `Analyse ${analysisId}`,
      href: `/analysis/${analysisId}`,
    });

    if (segments[2]) {
      const moduleLabels: Record<string, string> = {
        technical: 'Technique',
        performance: 'Performance', 
        content: 'Contenu',
        security: 'Sécurité',
        backlinks: 'Backlinks',
        report: 'Rapport',
        progress: 'Progression',
      };

      breadcrumbs.push({
        label: moduleLabels[segments[2]] || segments[2],
        isActive: segments.length === 3,
      });
    }
  } else if (segments[0] === 'dashboard') {
    breadcrumbs.push({
      label: 'Dashboard',
      isActive: true,
    });
  } else if (segments[0] === 'projects') {
    breadcrumbs.push({
      label: 'Projets',
      isActive: true,
    });
  }

  return breadcrumbs;
}