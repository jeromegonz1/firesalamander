"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { 
  Save,
  User,
  Globe,
  Bell,
  Key,
  Moon,
  Sun,
  Copy,
  RefreshCw,
  Check,
  AlertTriangle,
  Settings as SettingsIcon,
  Mail,
  Smartphone,
  Shield,
  BarChart3
} from "lucide-react";

interface UserProfile {
  name: string;
  email: string;
  company: string;
  language: string;
  timezone: string;
}

interface AnalysisSettings {
  default_analysis_type: string;
  include_images: boolean;
  deep_crawl: boolean;
  competitor_analysis: boolean;
  max_pages: number;
  crawl_delay: number;
  user_agent: string;
  target_keywords: string;
  keyword_density: boolean;
  meta_analysis: boolean;
  schema_validation: boolean;
}

interface NotificationSettings {
  analysis_complete_email: boolean;
  weekly_summary_email: boolean;
  critical_issues_email: boolean;
  frequency: 'immediate' | 'daily' | 'weekly';
}

interface APISettings {
  api_key: string;
  api_key_visible: boolean;
  total_requests: number;
  monthly_limit: number;
  requests_per_minute: number;
  burst_limit: number;
}

// Composant Settings Navigation (Mobile)
function MobileTabsDropdown({ activeTab, onTabChange }: { activeTab: string, onTabChange: (tab: string) => void }) {
  const [isOpen, setIsOpen] = useState(false);
  
  const tabs = [
    { id: 'general', label: 'Général', icon: User },
    { id: 'analysis', label: 'Analyses', icon: BarChart3 },
    { id: 'notifications', label: 'Notifications', icon: Bell },
    { id: 'api', label: 'API', icon: Key }
  ];

  const activeTabData = tabs.find(tab => tab.id === activeTab);

  return (
    <div className="md:hidden mb-6" data-testid="mobile-tabs-dropdown">
      <Button 
        variant="outline" 
        className="w-full justify-between"
        onClick={() => setIsOpen(!isOpen)}
      >
        <div className="flex items-center space-x-2">
          {activeTabData && <activeTabData.icon className="h-4 w-4" />}
          <span>{activeTabData?.label}</span>
        </div>
        <RefreshCw className={`h-4 w-4 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
      </Button>
      
      {isOpen && (
        <div className="mt-2 border rounded-lg bg-white shadow-lg">
          {tabs.map((tab) => (
            <button
              key={tab.id}
              className={`w-full text-left px-4 py-3 flex items-center space-x-2 hover:bg-gray-50 ${
                activeTab === tab.id ? 'bg-orange-50 text-orange-600' : ''
              }`}
              onClick={() => {
                onTabChange(tab.id);
                setIsOpen(false);
              }}
            >
              <tab.icon className="h-4 w-4" />
              <span>{tab.label}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}

// Composant Success Message
function SuccessMessage({ message, show }: { message: string, show: boolean }) {
  if (!show) return null;

  return (
    <div className="flex items-center space-x-2 p-3 bg-green-50 border border-green-200 rounded-lg" data-testid="success-message">
      <Check className="h-4 w-4 text-green-600" />
      <span className="text-sm text-green-800">{message}</span>
    </div>
  );
}

// Composant Validation Error
function ValidationError({ message, show }: { message: string, show: boolean }) {
  if (!show) return null;

  return (
    <div className="flex items-center space-x-2 p-3 bg-red-50 border border-red-200 rounded-lg" data-testid="email-validation-error">
      <AlertTriangle className="h-4 w-4 text-red-600" />
      <span className="text-sm text-red-800">{message}</span>
    </div>
  );
}

// Composant Theme Toggle
function ThemeToggle({ isDark, onToggle }: { isDark: boolean, onToggle: () => void }) {
  return (
    <div className="flex items-center justify-between p-4 border rounded-lg">
      <div className="flex items-center space-x-3">
        {isDark ? <Moon className="h-5 w-5" /> : <Sun className="h-5 w-5" />}
        <div>
          <Label className="font-medium">Mode sombre</Label>
          <p className="text-sm text-gray-600">Basculer vers le thème sombre</p>
        </div>
      </div>
      <Button
        variant="outline"
        size="sm"
        onClick={onToggle}
        data-testid="theme-toggle"
      >
        {isDark ? 'Désactiver' : 'Activer'}
      </Button>
    </div>
  );
}

// Composant API Key Management
function APIKeyManagement({ apiSettings, onUpdate }: { apiSettings: APISettings, onUpdate: (settings: Partial<APISettings>) => void }) {
  const [showRegenerateModal, setShowRegenerateModal] = useState(false);
  const [copied, setCopied] = useState(false);

  const handleShowKey = () => {
    onUpdate({ api_key_visible: !apiSettings.api_key_visible });
  };

  const handleCopyKey = () => {
    navigator.clipboard.writeText(apiSettings.api_key);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleRegenerateKey = () => {
    const newKey = 'fs_' + Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
    onUpdate({ api_key: newKey, api_key_visible: true });
    setShowRegenerateModal(false);
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Clé API</CardTitle>
          <CardDescription>Gérez votre clé d'accès à l'API Fire Salamander</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <Label>Clé API actuelle</Label>
            <div className="flex items-center space-x-2 mt-2">
              <div className="flex-1 p-3 bg-gray-50 rounded-lg font-mono text-sm" data-testid="current-api-key">
                {apiSettings.api_key_visible ? apiSettings.api_key : '••••••••••••••••••••••••••••'}
              </div>
              <Button variant="outline" size="sm" onClick={handleShowKey} data-testid="show-api-key">
                <Eye className="h-4 w-4" />
              </Button>
              <Button variant="outline" size="sm" onClick={handleCopyKey} data-testid="copy-api-key">
                <Copy className="h-4 w-4" />
              </Button>
            </div>
            {apiSettings.api_key_visible && (
              <div className="mt-2 p-3 bg-blue-50 rounded-lg" data-testid="api-key-visible">
                <p className="text-sm text-blue-800">Clé API visible - Copiez-la maintenant</p>
              </div>
            )}
            {copied && (
              <div className="mt-2 p-2 bg-green-50 rounded-lg" data-testid="api-key-copied">
                <p className="text-sm text-green-800">Clé API copiée dans le presse-papiers</p>
              </div>
            )}
          </div>

          <div className="flex space-x-2">
            <Button variant="outline" onClick={() => setShowRegenerateModal(true)} data-testid="regenerate-api-key">
              <RefreshCw className="h-4 w-4 mr-2" />
              Régénérer
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Regenerate Confirmation Modal */}
      {showRegenerateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" data-testid="regenerate-confirmation-modal">
          <Card className="max-w-md w-full mx-4">
            <CardHeader>
              <CardTitle>Confirmer la régénération</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <p className="text-sm text-gray-700" data-testid="regenerate-warning">
                  Attention : l'ancienne clé sera invalidée immédiatement et ne pourra plus être utilisée.
                </p>
                <div className="flex space-x-2">
                  <Button variant="outline" onClick={() => setShowRegenerateModal(false)}>
                    Annuler
                  </Button>
                  <Button onClick={handleRegenerateKey} data-testid="confirm-regenerate">
                    Confirmer
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* New API Key Generated Message */}
      {apiSettings.api_key.startsWith('fs_') && apiSettings.api_key_visible && (
        <div className="p-4 bg-green-50 border border-green-200 rounded-lg" data-testid="new-api-key-generated">
          <h4 className="font-semibold text-green-900 mb-2">Nouvelle clé API générée</h4>
          <p className="text-sm text-green-800">Votre nouvelle clé API a été générée avec succès. Assurez-vous de la copier maintenant.</p>
        </div>
      )}
    </div>
  );
}

export default function SettingsPage() {
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('general');
  const [isDarkMode, setIsDarkMode] = useState(false);
  const [showSuccess, setShowSuccess] = useState('');
  const [showError, setShowError] = useState('');

  // States for different settings
  const [profile, setProfile] = useState<UserProfile>({
    name: 'John Doe',
    email: 'john.doe@example.com',
    company: 'SEPTEO',
    language: 'fr',
    timezone: 'Europe/Paris'
  });

  const [analysisSettings, setAnalysisSettings] = useState<AnalysisSettings>({
    default_analysis_type: 'comprehensive',
    include_images: true,
    deep_crawl: false,
    competitor_analysis: true,
    max_pages: 100,
    crawl_delay: 1000,
    user_agent: 'Fire Salamander Bot 1.0',
    target_keywords: '',
    keyword_density: true,
    meta_analysis: true,
    schema_validation: false
  });

  const [notificationSettings, setNotificationSettings] = useState<NotificationSettings>({
    analysis_complete_email: true,
    weekly_summary_email: false,
    critical_issues_email: true,
    frequency: 'immediate'
  });

  const [apiSettings, setApiSettings] = useState<APISettings>({
    api_key: 'fs_example_key_1234567890abcdef',
    api_key_visible: false,
    total_requests: 1247,
    monthly_limit: 10000,
    requests_per_minute: 60,
    burst_limit: 100
  });

  useEffect(() => {
    // Simulate loading settings
    setTimeout(() => {
      setLoading(false);
    }, 1000);
  }, []);

  const showSuccessMessage = (message: string) => {
    setShowSuccess(message);
    setTimeout(() => setShowSuccess(''), 3000);
  };

  const showErrorMessage = (message: string) => {
    setShowError(message);
    setTimeout(() => setShowError(''), 5000);
  };

  const handleProfileSave = () => {
    // Validate email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(profile.email)) {
      showErrorMessage('Email invalide');
      return;
    }

    // Simulate API call
    setTimeout(() => {
      showSuccessMessage('Profil mis à jour avec succès');
    }, 500);
  };

  const handleAnalysisSettingsSave = () => {
    // Validate max pages
    if (analysisSettings.max_pages > 10000) {
      showErrorMessage('Maximum 10000 pages autorisées');
      return;
    }

    setTimeout(() => {
      showSuccessMessage('Paramètres d\'analyse sauvegardés');
    }, 500);
  };

  const handleNotificationTest = () => {
    setTimeout(() => {
      showSuccessMessage('Email de test envoyé avec succès');
    }, 1000);
  };

  const handleThemeToggle = () => {
    setIsDarkMode(!isDarkMode);
    // Apply theme to body
    if (!isDarkMode) {
      document.body.classList.add('dark');
    } else {
      document.body.classList.remove('dark');
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]" data-testid="settings-loading">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-orange-500 mx-auto mb-4"></div>
          <p className="text-gray-600">Chargement des paramètres...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-6" data-testid="settings-container">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Paramètres</h1>
        <p className="text-gray-600 mt-1">
          Gérez vos préférences et configurations Fire Salamander
        </p>
      </div>

      {/* Success/Error Messages */}
      <div className="space-y-2" data-testid="settings-feedback" aria-live="polite">
        <SuccessMessage message={showSuccess} show={!!showSuccess} />
        <ValidationError message={showError} show={!!showError} />
      </div>

      {/* Mobile Tabs Dropdown */}
      <MobileTabsDropdown activeTab={activeTab} onTabChange={setActiveTab} />

      {/* Desktop Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="hidden md:grid w-full grid-cols-4" data-testid="settings-tabs">
          <TabsTrigger value="general" data-testid="tab-general">
            <User className="h-4 w-4 mr-2" />
            Général
          </TabsTrigger>
          <TabsTrigger value="analysis" data-testid="tab-analysis">
            <BarChart3 className="h-4 w-4 mr-2" />
            Analyses
          </TabsTrigger>
          <TabsTrigger value="notifications" data-testid="tab-notifications">
            <Bell className="h-4 w-4 mr-2" />
            Notifications
          </TabsTrigger>
          <TabsTrigger value="api" data-testid="tab-api">
            <Key className="h-4 w-4 mr-2" />
            API
          </TabsTrigger>
        </TabsList>

        {/* General Settings Tab */}
        <TabsContent value="general" className="space-y-6">
          <Card data-testid="profile-settings">
            <CardHeader>
              <CardTitle>Profil</CardTitle>
              <CardDescription>Informations personnelles et préférences</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="profile-name">Nom complet</Label>
                  <Input
                    id="profile-name"
                    value={profile.name}
                    onChange={(e) => setProfile({ ...profile, name: e.target.value })}
                    data-testid="profile-name-input"
                  />
                </div>
                <div>
                  <Label htmlFor="profile-email">Email</Label>
                  <Input
                    id="profile-email"
                    type="email"
                    value={profile.email}
                    onChange={(e) => setProfile({ ...profile, email: e.target.value })}
                    data-testid="profile-email-input"
                  />
                </div>
                <div>
                  <Label htmlFor="profile-company">Entreprise</Label>
                  <Input
                    id="profile-company"
                    value={profile.company}
                    onChange={(e) => setProfile({ ...profile, company: e.target.value })}
                    data-testid="profile-company-input"
                  />
                </div>
              </div>
              
              <Button onClick={handleProfileSave} data-testid="save-profile-button" className="bg-orange-500 hover:bg-orange-600">
                <Save className="h-4 w-4 mr-2" />
                Sauvegarder Profil
              </Button>
            </CardContent>
          </Card>

          <Card data-testid="locale-settings">
            <CardHeader>
              <CardTitle>Langue et Région</CardTitle>
              <CardDescription>Préférences linguistiques et de fuseau horaire</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="language">Langue</Label>
                  <select
                    id="language"
                    className="w-full p-2 border rounded-lg"
                    value={profile.language}
                    onChange={(e) => setProfile({ ...profile, language: e.target.value })}
                    data-testid="language-select"
                  >
                    <option value="fr" data-testid="lang-fr">Français</option>
                    <option value="en" data-testid="lang-en">English</option>
                    <option value="es">Español</option>
                  </select>
                </div>
                <div>
                  <Label htmlFor="timezone">Fuseau horaire</Label>
                  <select
                    id="timezone"
                    className="w-full p-2 border rounded-lg"
                    value={profile.timezone}
                    onChange={(e) => setProfile({ ...profile, timezone: e.target.value })}
                    data-testid="timezone-select"
                  >
                    <option value="Europe/Paris">Europe/Paris</option>
                    <option value="America/New_York">America/New_York</option>
                    <option value="Asia/Tokyo">Asia/Tokyo</option>
                  </select>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Apparence</CardTitle>
              <CardDescription>Personnalisez l'apparence de l'interface</CardDescription>
            </CardHeader>
            <CardContent>
              <ThemeToggle isDark={isDarkMode} onToggle={handleThemeToggle} />
            </CardContent>
          </Card>
        </TabsContent>

        {/* Analysis Settings Tab */}
        <TabsContent value="analysis" className="space-y-6">
          <Card data-testid="default-analysis-config">
            <CardHeader>
              <CardTitle>Configuration des Analyses</CardTitle>
              <CardDescription>Paramètres par défaut pour vos analyses SEO</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <Label htmlFor="default-analysis-type">Type d'analyse par défaut</Label>
                <select
                  id="default-analysis-type"
                  className="w-full p-2 border rounded-lg mt-1"
                  value={analysisSettings.default_analysis_type}
                  onChange={(e) => setAnalysisSettings({ ...analysisSettings, default_analysis_type: e.target.value })}
                  data-testid="default-analysis-type"
                >
                  <option value="comprehensive">Analyse Complète</option>
                  <option value="quick">Analyse Rapide</option>
                  <option value="custom">Analyse Personnalisée</option>
                </select>
              </div>

              <div className="space-y-3">
                <h4 className="font-semibold">Options d'analyse</h4>
                <div className="space-y-2">
                  <label className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      checked={analysisSettings.include_images}
                      onChange={(e) => setAnalysisSettings({ ...analysisSettings, include_images: e.target.checked })}
                      data-testid="include-images-check"
                    />
                    <span>Inclure l'analyse des images</span>
                  </label>
                  <label className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      checked={analysisSettings.deep_crawl}
                      onChange={(e) => setAnalysisSettings({ ...analysisSettings, deep_crawl: e.target.checked })}
                      data-testid="deep-crawl-check"
                    />
                    <span>Crawl en profondeur</span>
                  </label>
                  <label className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      checked={analysisSettings.competitor_analysis}
                      onChange={(e) => setAnalysisSettings({ ...analysisSettings, competitor_analysis: e.target.checked })}
                      data-testid="competitor-analysis-check"
                    />
                    <span>Analyse concurrentielle</span>
                  </label>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card data-testid="crawl-settings">
            <CardHeader>
              <CardTitle>Paramètres de Crawl</CardTitle>
              <CardDescription>Configuration du comportement du crawler</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="max-pages">Pages maximum</Label>
                  <Input
                    id="max-pages"
                    type="number"
                    value={analysisSettings.max_pages}
                    onChange={(e) => setAnalysisSettings({ ...analysisSettings, max_pages: parseInt(e.target.value) })}
                    data-testid="max-pages-input"
                  />
                  {analysisSettings.max_pages > 10000 && (
                    <p className="text-sm text-red-600 mt-1" data-testid="max-pages-validation-error">
                      Maximum 10000 pages autorisées
                    </p>
                  )}
                </div>
                <div>
                  <Label htmlFor="crawl-delay">Délai entre requêtes (ms)</Label>
                  <Input
                    id="crawl-delay"
                    type="number"
                    value={analysisSettings.crawl_delay}
                    onChange={(e) => setAnalysisSettings({ ...analysisSettings, crawl_delay: parseInt(e.target.value) })}
                    data-testid="crawl-delay-input"
                  />
                </div>
                <div className="md:col-span-2">
                  <Label htmlFor="user-agent">User Agent</Label>
                  <Input
                    id="user-agent"
                    value={analysisSettings.user_agent}
                    onChange={(e) => setAnalysisSettings({ ...analysisSettings, user_agent: e.target.value })}
                    data-testid="user-agent-input"
                  />
                </div>
              </div>
              
              <Button onClick={handleAnalysisSettingsSave} data-testid="save-crawl-settings" className="bg-orange-500 hover:bg-orange-600">
                <Save className="h-4 w-4 mr-2" />
                Sauvegarder
              </Button>
              
              {showSuccess === 'Paramètres d\'analyse sauvegardés' && (
                <div className="mt-2 p-2 bg-green-50 rounded-lg" data-testid="crawl-settings-saved">
                  <p className="text-sm text-green-800">Paramètres de crawl sauvegardés</p>
                </div>
              )}
            </CardContent>
          </Card>

          <Card data-testid="seo-preferences">
            <CardHeader>
              <CardTitle>Préférences SEO</CardTitle>
              <CardDescription>Options d'analyse SEO avancées</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-3">
                <label className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    checked={analysisSettings.keyword_density}
                    onChange={(e) => setAnalysisSettings({ ...analysisSettings, keyword_density: e.target.checked })}
                    data-testid="keyword-density-check"
                  />
                  <span>Analyse de densité des mots-clés</span>
                </label>
                <label className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    checked={analysisSettings.meta_analysis}
                    onChange={(e) => setAnalysisSettings({ ...analysisSettings, meta_analysis: e.target.checked })}
                    data-testid="meta-analysis-check"
                  />
                  <span>Analyse des meta tags</span>
                </label>
                <label className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    checked={analysisSettings.schema_validation}
                    onChange={(e) => setAnalysisSettings({ ...analysisSettings, schema_validation: e.target.checked })}
                    data-testid="schema-validation-check"
                  />
                  <span>Validation des données structurées</span>
                </label>
              </div>

              <div>
                <Label htmlFor="target-keywords">Mots-clés cibles (séparés par des virgules)</Label>
                <Input
                  id="target-keywords"
                  placeholder="plage, restaurant, marina"
                  value={analysisSettings.target_keywords}
                  onChange={(e) => setAnalysisSettings({ ...analysisSettings, target_keywords: e.target.value })}
                  data-testid="target-keywords-input"
                />
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Notifications Settings Tab */}
        <TabsContent value="notifications" className="space-y-6">
          <Card data-testid="email-notifications">
            <CardHeader>
              <CardTitle>Notifications Email</CardTitle>
              <CardDescription>Configurez vos préférences de notifications</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-3">
                <label className="flex items-center justify-between p-3 border rounded-lg">
                  <div className="flex items-center space-x-2">
                    <Mail className="h-4 w-4" />
                    <span>Analyse terminée</span>
                  </div>
                  <input
                    type="checkbox"
                    checked={notificationSettings.analysis_complete_email}
                    onChange={(e) => setNotificationSettings({ ...notificationSettings, analysis_complete_email: e.target.checked })}
                    data-testid="analysis-complete-email"
                  />
                </label>
                <label className="flex items-center justify-between p-3 border rounded-lg">
                  <div className="flex items-center space-x-2">
                    <BarChart3 className="h-4 w-4" />
                    <span>Rapport hebdomadaire</span>
                  </div>
                  <input
                    type="checkbox"
                    checked={notificationSettings.weekly_summary_email}
                    onChange={(e) => setNotificationSettings({ ...notificationSettings, weekly_summary_email: e.target.checked })}
                    data-testid="weekly-summary-email"
                  />
                </label>
                <label className="flex items-center justify-between p-3 border rounded-lg">
                  <div className="flex items-center space-x-2">
                    <AlertTriangle className="h-4 w-4" />
                    <span>Problèmes critiques</span>
                  </div>
                  <input
                    type="checkbox"
                    checked={notificationSettings.critical_issues_email}
                    onChange={(e) => setNotificationSettings({ ...notificationSettings, critical_issues_email: e.target.checked })}
                    data-testid="critical-issues-email"
                  />
                </label>
              </div>
            </CardContent>
          </Card>

          <Card data-testid="notification-frequency">
            <CardHeader>
              <CardTitle>Fréquence des Notifications</CardTitle>
              <CardDescription>À quelle fréquence recevoir les notifications</CardDescription>
            </CardHeader>
            <CardContent>
              <div>
                <Label htmlFor="frequency">Fréquence</Label>
                <select
                  id="frequency"
                  className="w-full p-2 border rounded-lg mt-1"
                  value={notificationSettings.frequency}
                  onChange={(e) => setNotificationSettings({ ...notificationSettings, frequency: e.target.value as any })}
                  data-testid="frequency-select"
                >
                  <option value="immediate" data-testid="freq-immediate">Immédiat</option>
                  <option value="daily" data-testid="freq-daily">Quotidien</option>
                  <option value="weekly" data-testid="freq-weekly">Hebdomadaire</option>
                </select>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Test des Notifications</CardTitle>
              <CardDescription>Vérifiez que vos notifications fonctionnent</CardDescription>
            </CardHeader>
            <CardContent>
              <Button onClick={handleNotificationTest} data-testid="test-notification-button">
                <Mail className="h-4 w-4 mr-2" />
                Envoyer un test
              </Button>
              
              {showSuccess === 'Email de test envoyé avec succès' && (
                <div className="mt-4 p-3 bg-green-50 border border-green-200 rounded-lg" data-testid="test-notification-sent">
                  <p className="text-sm text-green-800">Email de test envoyé à {profile.email}</p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* API Settings Tab */}
        <TabsContent value="api" className="space-y-6">
          <APIKeyManagement apiSettings={apiSettings} onUpdate={(updates) => setApiSettings({ ...apiSettings, ...updates })} />

          <Card data-testid="api-usage-stats">
            <CardHeader>
              <CardTitle>Utilisation API</CardTitle>
              <CardDescription>Statistiques d'utilisation de votre clé API</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div className="text-center p-4 bg-blue-50 rounded-lg">
                  <div className="text-2xl font-bold text-blue-600" data-testid="total-requests">{apiSettings.total_requests}</div>
                  <div className="text-sm text-gray-600">Requêtes totales</div>
                </div>
                <div className="text-center p-4 bg-green-50 rounded-lg">
                  <div className="text-2xl font-bold text-green-600" data-testid="monthly-limit">{apiSettings.monthly_limit}</div>
                  <div className="text-sm text-gray-600">Limite mensuelle</div>
                </div>
                <div className="text-center p-4 bg-orange-50 rounded-lg">
                  <div className="text-2xl font-bold text-orange-600" data-testid="remaining-requests">
                    {apiSettings.monthly_limit - apiSettings.total_requests}
                  </div>
                  <div className="text-sm text-gray-600">Requêtes restantes</div>
                </div>
              </div>

              {/* Simulated Usage Chart */}
              <div className="mt-6" data-testid="usage-chart">
                <h4 className="font-semibold mb-3">Utilisation des 30 derniers jours</h4>
                <div className="h-32 bg-gray-50 rounded-lg flex items-end justify-center space-x-1 p-4">
                  {Array.from({ length: 30 }, (_, i) => (
                    <div
                      key={i}
                      className="bg-orange-500 rounded-t"
                      style={{
                        height: `${Math.random() * 80 + 20}%`,
                        width: '3%'
                      }}
                    />
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>

          <Card data-testid="rate-limiting-config">
            <CardHeader>
              <CardTitle>Limitation de Débit</CardTitle>
              <CardDescription>Configurez les limites d'utilisation de l'API</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="requests-per-minute">Requêtes par minute</Label>
                  <Input
                    id="requests-per-minute"
                    type="number"
                    value={apiSettings.requests_per_minute}
                    onChange={(e) => setApiSettings({ ...apiSettings, requests_per_minute: parseInt(e.target.value) })}
                    data-testid="requests-per-minute"
                  />
                </div>
                <div>
                  <Label htmlFor="burst-limit">Limite de rafale</Label>
                  <Input
                    id="burst-limit"
                    type="number"
                    value={apiSettings.burst_limit}
                    onChange={(e) => setApiSettings({ ...apiSettings, burst_limit: parseInt(e.target.value) })}
                    data-testid="burst-limit"
                  />
                </div>
              </div>
              
              <Button data-testid="save-rate-limits" className="bg-orange-500 hover:bg-orange-600">
                <Save className="h-4 w-4 mr-2" />
                Sauvegarder Limites
              </Button>
              
              {showSuccess === 'Limites API sauvegardées' && (
                <div className="mt-2 p-2 bg-green-50 rounded-lg" data-testid="rate-limits-saved">
                  <p className="text-sm text-green-800">Limites de débit sauvegardées</p>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}