'use client';

import { useState } from 'react';

export default function Settings() {
  const [settings, setSettings] = useState({
    apiKey: '',
    autoScan: true,
    notifications: true,
    complianceStandards: ['GDPR', 'CCPA', 'HIPAA'],
    riskThreshold: 'high'
  });

  return (
    <div className="min-h-screen p-8 max-w-4xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-white">Settings</h1>
        <p className="text-gray-400">Configure your Aegis AI preferences</p>
      </div>

      <div className="space-y-8">
        {/* API Configuration */}
        <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
          <h2 className="text-xl font-bold text-white mb-4">API Configuration</h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-200 mb-2">
                API Key
              </label>
              <input
                type="password"
                value={settings.apiKey}
                onChange={(e) => setSettings({ ...settings, apiKey: e.target.value })}
                placeholder="Enter your API key"
                className="w-full px-4 py-3 bg-white/5 border border-white/20 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              />
            </div>
          </div>
        </div>

        {/* Scan Settings */}
        <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
          <h2 className="text-xl font-bold text-white mb-4">Scan Settings</h2>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <div className="font-medium text-white">Automatic Scanning</div>
                <div className="text-sm text-gray-400">Automatically scan new pull requests</div>
              </div>
              <label className="relative inline-flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  checked={settings.autoScan}
                  onChange={(e) => setSettings({ ...settings, autoScan: e.target.checked })}
                  className="sr-only peer"
                />
                <div className="w-11 h-6 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-purple-600"></div>
              </label>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-200 mb-2">
                Risk Threshold
              </label>
              <select
                value={settings.riskThreshold}
                onChange={(e) => setSettings({ ...settings, riskThreshold: e.target.value })}
                className="w-full px-4 py-3 bg-white/5 border border-white/20 rounded-xl text-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              >
                <option value="critical">Critical Only</option>
                <option value="high">High and Above</option>
                <option value="medium">Medium and Above</option>
                <option value="low">All Issues</option>
              </select>
            </div>
          </div>
        </div>

        {/* Compliance Standards */}
        <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
          <h2 className="text-xl font-bold text-white mb-4">Compliance Standards</h2>
          <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
            {['GDPR', 'CCPA', 'HIPAA', 'PCI-DSS', 'SOC2', 'ISO27001'].map((standard) => (
              <label key={standard} className="flex items-center space-x-3 p-3 bg-white/5 rounded-lg border border-white/10 hover:bg-white/10 transition-colors cursor-pointer">
                <input
                  type="checkbox"
                  checked={settings.complianceStandards.includes(standard)}
                  onChange={(e) => {
                    const newStandards = e.target.checked
                      ? [...settings.complianceStandards, standard]
                      : settings.complianceStandards.filter(s => s !== standard);
                    setSettings({ ...settings, complianceStandards: newStandards });
                  }}
                  className="w-4 h-4 text-purple-600 bg-white/5 border-white/20 rounded focus:ring-purple-500 focus:ring-2"
                />
                <span className="text-white text-sm">{standard}</span>
              </label>
            ))}
          </div>
        </div>

        {/* Actions */}
        <div className="flex justify-end space-x-4">
          <button className="px-6 py-3 bg-white/10 hover:bg-white/20 text-white rounded-xl font-medium transition-colors">
            Reset to Defaults
          </button>
          <button className="px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white rounded-xl font-medium transition-colors">
            Save Settings
          </button>
        </div>
      </div>
    </div>
  );
}