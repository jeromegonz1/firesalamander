"use client";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { 
  BarChart3, 
  Search, 
  Clock, 
  FileText, 
  FolderOpen,
  TrendingUp,
  Wrench,
  GitCompare
} from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

const navigationItems = [
  {
    name: "Dashboard",
    href: "/dashboard",
    icon: BarChart3,
  },
  {
    name: "Nouvelle Analyse",
    href: "/analysis/new",
    icon: Search,
  },
  {
    name: "Projets",
    href: "/projects",
    icon: FolderOpen,
  },
  {
    name: "Mots-clés",
    href: "/analysis/keywords",
    icon: TrendingUp,
    badge: "IA",
  },
  {
    name: "Technique",
    href: "/analysis/technical", 
    icon: Wrench,
  },
  {
    name: "Comparaison",
    href: "/compare",
    icon: GitCompare,
  },
];

export function Navigation() {
  const pathname = usePathname();

  return (
    <nav className="border-b bg-white">
      <div className="container mx-auto px-4">
        <div className="flex space-x-1 overflow-x-auto">
          {navigationItems.map((item) => {
            const isActive = pathname === item.href || pathname.startsWith(item.href + "/");
            
            return (
              <Link key={item.name} href={item.href}>
                <Button
                  variant={isActive ? "default" : "ghost"}
                  className={cn(
                    "flex items-center space-x-2 whitespace-nowrap",
                    isActive && "bg-orange-500 hover:bg-orange-600"
                  )}
                  size="sm"
                >
                  <item.icon className="h-4 w-4" />
                  <span>{item.name}</span>
                  {item.badge && (
                    <span className="ml-1 rounded-full bg-purple-100 px-2 py-0.5 text-xs font-medium text-purple-700">
                      {item.badge}
                    </span>
                  )}
                </Button>
              </Link>
            );
          })}
        </div>
      </div>
    </nav>
  );
}