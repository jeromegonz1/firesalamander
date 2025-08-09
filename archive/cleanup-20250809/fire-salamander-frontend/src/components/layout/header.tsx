"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Bell, Settings, User } from "lucide-react";

export function Header() {
  return (
    <header className="border-b bg-white shadow-sm">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        {/* Logo & Brand */}
        <div className="flex items-center space-x-4">
          <div className="flex items-center space-x-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-orange-500 text-white font-bold">
              🔥
            </div>
            <div className="flex flex-col">
              <span className="text-sm font-medium text-gray-600">SEPTEO</span>
              <span className="text-lg font-bold text-gray-900">Fire Salamander</span>
            </div>
          </div>
          <Badge variant="secondary" className="ml-2">
            v1.0.0
          </Badge>
        </div>

        {/* Actions */}
        <div className="flex items-center space-x-2">
          <div className="flex items-center space-x-1 text-sm text-gray-600">
            <div className="h-2 w-2 rounded-full bg-green-500"></div>
            <span>En ligne</span>
          </div>
          
          <Button variant="ghost" size="sm">
            <Bell className="h-4 w-4" />
          </Button>
          
          <Button variant="ghost" size="sm">
            <Settings className="h-4 w-4" />
          </Button>
          
          <Button variant="ghost" size="sm">
            <User className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </header>
  );
}