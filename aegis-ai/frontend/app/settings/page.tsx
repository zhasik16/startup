'use client';

import { useState } from 'react';

export default function Settings() {
  const [settings, setSettings] = useState({
    apiKey: '',
    autoScan: true,
    notifications: true,
    emailReports: false,
    complianceStandards: ['GDPR', 'CCPA', 'HIPAA'],
    riskThreshold: 'high',
    language: 'english',
    theme: 'dark'
  });

  const [saved, setSaved] = useState(false);

  const handleSave = () => {
    setSaved(true);
    setTimeout(() => setSaved(false), 3000);
  };

  const handleReset = () => {
    setSettings({
      apiKey: '',
      autoScan: true,
      notifications: true,
      emailReports: false,
      complianceStandards: ['GDPR', 'CCPA', 'HIPAA'],
      riskThreshold: 'high',
      language: 'english',
      theme: 'dark'
    });
  };

  return (
    <div className="min-h-screen p-6 max-w-6xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white mb-2">Settings</h1>
            <p className="text-gray-300">Configure your Aegis AI security preferences</p>
          </div>
          {saved && (
            <div className="flex items-center space-x-2 px-4 py-2 bg-green-500/20 border border-green-500/30 rounded-lg">
              <span className="text-green-400">✓</span>
              <span className="text-green-300 text-sm">Settings saved successfully!</span>
            </div>
          )}
        </div>
      </div>

      <div className="grid lg:grid-cols-3 gap-8">
        {/* Main Settings Column */}
        <div className="lg:col-span-2 space-y-6">
          {/* API Configuration */}
          <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
            <div className="flex items-center space-x-3 mb-4">
              <div className="w-8 h-8 bg-blue-500/20 rounded-lg flex items-center justify-center">
                <span className="text-blue-400">🔑</span>
              </div>
              <h2 className="text-xl font-bold text-white">API Configuration</h2>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-200 mb-2">
                  GitHub Access Token
                </label>
                <input
                  type="password"
                  value={settings.apiKey}
                  onChange={(e) => setSettings({ ...settings, apiKey: e.target.value })}
                  placeholder="ghp_your_token_here"
                  className="w-full px-4 py-3 bg-black/30 border border-gray-600 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                <p className="text-xs text-gray-400 mt-2">
                  Required for accessing private repositories and auto-committing fixes
                </p>
              </div>
            </div>
          </div>

          {/* Scan Settings */}
          <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
            <div className="flex items-center space-x-3 mb-4">
              <div className="w-8 h-8 bg-green-500/20 rounded-lg flex items-center justify-center">
                <span className="text-green-400">🔍</span>
              </div>
              <h2 className="text-xl font-bold text-white">Scan Settings</h2>
            </div>
            <div className="space-y-6">
              <div className="flex items-center justify-between p-4 bg-white/5 rounded-lg border border-white/10">
                <div>
                  <div className="font-medium text-white">Automatic PR Scanning</div>
                  <div className="text-sm text-gray-300">Automatically scan new pull requests</div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    checked={settings.autoScan}
                    onChange={(e) => setSettings({ ...settings, autoScan: e.target.checked })}
                    className="sr-only peer"
                  />
                  <div className="w-12 h-6 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-6 peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-green-600"></div>
                </label>
              </div>

              <div className="flex items-center justify-between p-4 bg-white/5 rounded-lg border border-white/10">
                <div>
                  <div className="font-medium text-white">Push Notifications</div>
                  <div className="text-sm text-gray-300">Get notified for critical security issues</div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    checked={settings.notifications}
                    onChange={(e) => setSettings({ ...settings, notifications: e.target.checked })}
                    className="sr-only peer"
                  />
                  <div className="w-12 h-6 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-6 peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-green-600"></div>
                </label>
              </div>

              <div className="flex items-center justify-between p-4 bg-white/5 rounded-lg border border-white/10">
                <div>
                  <div className="font-medium text-white">Email Reports</div>
                  <div className="text-sm text-gray-300">Receive weekly security reports via email</div>
                </div>
                <label className="relative inline-flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    checked={settings.emailReports}
                    onChange={(e) => setSettings({ ...settings, emailReports: e.target.checked })}
                    className="sr-only peer"
                  />
                  <div className="w-12 h-6 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-6 peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-green-600"></div>
                </label>
              </div>

              <div>
                <label className="block text-sm font-medium text-white mb-3">
                  Risk Threshold
                </label>
                <select
                  value={settings.riskThreshold}
                  onChange={(e) => setSettings({ ...settings, riskThreshold: e.target.value })}
                  className="w-full px-4 py-3 bg-black/30 border border-gray-600 rounded-xl text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                >
                  <option value="critical">🚨 Critical Only</option>
                  <option value="high">⚠️ High and Above</option>
                  <option value="medium">📝 Medium and Above</option>
                  <option value="low">ℹ️ All Issues</option>
                </select>
                <p className="text-xs text-gray-400 mt-2">
                  Determines which security issues are reported
                </p>
              </div>
            </div>
          </div>

          {/* Compliance Standards */}
          <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
            <div className="flex items-center space-x-3 mb-4">
              <div className="w-8 h-8 bg-purple-500/20 rounded-lg flex items-center justify-center">
                <span className="text-purple-400">⚖️</span>
              </div>
              <h2 className="text-xl font-bold text-white">Compliance Standards</h2>
            </div>
            <p className="text-gray-300 mb-4 text-sm">
              Select the compliance standards you need to adhere to
            </p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {[
                { id: 'GDPR', name: 'GDPR', description: 'General Data Protection Regulation' },
                { id: 'CCPA', name: 'CCPA', description: 'California Consumer Privacy Act' },
                { id: 'HIPAA', name: 'HIPAA', description: 'Health Insurance Portability Act' },
                { id: 'PCI-DSS', name: 'PCI DSS', description: 'Payment Card Industry Standards' },
                { id: 'SOC2', name: 'SOC 2', description: 'Service Organization Control' },
                { id: 'ISO27001', name: 'ISO 27001', description: 'Information Security Management' },
              ].map((standard) => (
                <label key={standard.id} className="flex items-start space-x-3 p-4 bg-white/5 rounded-lg border border-white/10 hover:bg-white/10 transition-colors cursor-pointer">
                  <input
                    type="checkbox"
                    checked={settings.complianceStandards.includes(standard.id)}
                    onChange={(e) => {
                      const newStandards = e.target.checked
                        ? [...settings.complianceStandards, standard.id]
                        : settings.complianceStandards.filter(s => s !== standard.id);
                      setSettings({ ...settings, complianceStandards: newStandards });
                    }}
                    className="w-4 h-4 text-purple-600 bg-black/30 border-gray-600 rounded focus:ring-purple-500 focus:ring-2 mt-1"
                  />
                  <div className="flex-1">
                    <div className="font-medium text-white text-sm">{standard.name}</div>
                    <div className="text-gray-400 text-xs">{standard.description}</div>
                  </div>
                </label>
              ))}
            </div>
          </div>
        </div>

        {/* Sidebar Column */}
        <div className="space-y-6">
          {/* Preferences */}
          <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
            <h3 className="font-bold text-white mb-4">Preferences</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-200 mb-2">
                  Language
                </label>
                <select
                  value={settings.language}
                  onChange={(e) => setSettings({ ...settings, language: e.target.value })}
                  className="w-full px-3 py-2 bg-black/30 border border-gray-600 rounded-lg text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="english">English</option>
                  <option value="spanish">Spanish</option>
                  <option value="french">French</option>
                  <option value="german">German</option>
                </select>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-200 mb-2">
                  Theme
                </label>
                <select
                  value={settings.theme}
                  onChange={(e) => setSettings({ ...settings, theme: e.target.value })}
                  className="w-full px-3 py-2 bg-black/30 border border-gray-600 rounded-lg text-white text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="dark">Dark</option>
                  <option value="light">Light</option>
                  <option value="system">System</option>
                </select>
              </div>
            </div>
          </div>

          {/* Quick Stats */}
          <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
            <h3 className="font-bold text-white mb-4">Security Stats</h3>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-gray-300 text-sm">Scans This Month</span>
                <span className="text-white font-medium">24</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-gray-300 text-sm">Critical Issues Found</span>
                <span className="text-red-400 font-medium">3</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-gray-300 text-sm">Auto-Fixes Applied</span>
                <span className="text-green-400 font-medium">18</span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-gray-300 text-sm">Compliance Score</span>
                <span className="text-blue-400 font-medium">92%</span>
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
            <h3 className="font-bold text-white mb-4">Actions</h3>
            <div className="space-y-3">
              <button
                onClick={handleSave}
                className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white py-3 rounded-lg font-medium transition-all duration-200 transform hover:scale-105"
              >
                Save Changes
              </button>
              <button
                onClick={handleReset}
                className="w-full bg-white/10 hover:bg-white/20 text-white py-3 rounded-lg font-medium transition-colors"
              >
                Reset to Defaults
              </button>
              <button className="w-full bg-red-500/20 hover:bg-red-500/30 text-red-300 py-3 rounded-lg font-medium transition-colors border border-red-500/30">
                Delete All Data
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}