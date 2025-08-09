/**
 * Fire Salamander - Navigation TypeScript Interfaces
 * Lead Tech quality navigation system
 */

export interface NavigationItem {
  id: string;
  label: string;
  href: string;
  icon: string;
  badge?: string;
  isActive?: boolean;
  isDisabled?: boolean;
  children?: NavigationItem[];
}

export interface AnalysisNavigationItem extends NavigationItem {
  moduleType: AnalysisModuleType;
  score?: number;
  grade?: string;
  status: AnalysisModuleStatus;
  lastUpdated?: string;
}

export enum AnalysisModuleType {
  OVERVIEW = 'overview',
  TECHNICAL = 'technical',
  PERFORMANCE = 'performance',
  CONTENT = 'content',
  SECURITY = 'security',
  BACKLINKS = 'backlinks',
}

export enum AnalysisModuleStatus {
  COMPLETED = 'completed',
  IN_PROGRESS = 'in_progress',
  PENDING = 'pending',
  ERROR = 'error',
  NOT_STARTED = 'not_started',
}

export interface BreadcrumbItem {
  label: string;
  href?: string;
  isActive?: boolean;
}

export interface NavigationState {
  currentAnalysisId?: string;
  currentModule?: AnalysisModuleType;
  breadcrumbs: BreadcrumbItem[];
  sidebarCollapsed: boolean;
}

export interface AnalysisProgress {
  overall: number;
  modules: Record<AnalysisModuleType, {
    progress: number;
    status: AnalysisModuleStatus;
    score?: number;
    grade?: string;
  }>;
}