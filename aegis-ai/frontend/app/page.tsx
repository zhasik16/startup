'use client';

import { useState } from 'react';
import { AegisApi } from '@/lib/api';

export default function Home() {
  const [repoUrl, setRepoUrl] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [analysisId, setAnalysisId] = useState<string | null>(null);

  const handleAnalyze = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!repoUrl) return;

    setIsLoading(true);
    try {
      const result = await AegisApi.triggerAnalysis(repoUrl);
      setAnalysisId(result.analysis_id);
      // In real implementation, you'd redirect to analysis page
      window.location.href = `/analysis/${result.analysis_id}`;
    } catch (error) {
      console.error('Analysis failed:', error);
      alert('Analysis failed. Please check the repository URL and try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-8">
      <div className="max-w-4xl w-full space-y-12">
        {/* Hero Section */}
        <div className="text-center space-y-6">
          <div className="flex justify-center items-center space-x-4 mb-8">
            <div className="w-12 h-12 bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl flex items-center justify-center">
              <span className="text-2xl">🛡️</span>
            </div>
            <h1 className="text-5xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
              Aegis AI
            </h1>
          </div>
          
          <p className="text-xl text-gray-300 max-w-2xl mx-auto leading-relaxed">
            Automated security expert that integrates into your SDLC. 
            <span className="block text-purple-300 font-semibold mt-2">
              Find and fix compliance violations before they reach production.
            </span>
          </p>
        </div>

        {/* Analysis Form */}
        <div className="bg-white/10 backdrop-blur-lg rounded-2xl p-8 border border-white/20 shadow-2xl">
          <form onSubmit={handleAnalyze} className="space-y-6">
            <div>
              <label htmlFor="repoUrl" className="block text-sm font-medium text-gray-200 mb-2">
                Repository URL
              </label>
              <input
                type="url"
                id="repoUrl"
                value={repoUrl}
                onChange={(e) => setRepoUrl(e.target.value)}
                placeholder="https://github.com/username/repository"
                className="w-full px-4 py-3 bg-white/5 border border-white/20 rounded-xl text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all duration-200"
                required
              />
            </div>

            <button
              type="submit"
              disabled={isLoading || !repoUrl}
              className="w-full py-4 px-6 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 disabled:from-gray-600 disabled:to-gray-600 disabled:cursor-not-allowed text-white font-semibold rounded-xl transition-all duration-200 transform hover:scale-[1.02] focus:scale-[0.98] focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:ring-offset-slate-900"
            >
              {isLoading ? (
                <div className="flex items-center justify-center space-x-2">
                  <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  <span>Analyzing Codebase...</span>
                </div>
              ) : (
                '🚀 Start Security Analysis'
              )}
            </button>
          </form>

          {/* Features Grid */}
          <div className="grid md:grid-cols-3 gap-6 mt-12 pt-8 border-t border-white/10">
            {[
              {
                icon: '🔍',
                title: 'Code & Data Flow Analysis',
                description: 'Dynamic mapping of PII data flow through your entire application'
              },
              {
                icon: '⚖️',
                title: 'Regulation-as-Code',
                description: 'Machine-readable rules for global privacy regulations'
              },
              {
                icon: '🤖',
                title: 'Automated Detection & Fixes',
                description: 'AI-powered violation detection with automatic fix suggestions'
              }
            ].map((feature, index) => (
              <div key={index} className="text-center p-4 rounded-lg bg-white/5 hover:bg-white/10 transition-colors duration-200">
                <div className="text-2xl mb-3">{feature.icon}</div>
                <h3 className="font-semibold text-white mb-2">{feature.title}</h3>
                <p className="text-sm text-gray-300">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
          {[
            { label: 'Critical Issues', value: '0' },
            { label: 'Auto-Fixes', value: '0' },
            { label: 'Compliance', value: '0%' },
            { label: 'Analysis Time', value: '<30s' }
          ].map((stat, index) => (
            <div key={index} className="bg-white/5 rounded-lg p-4 border border-white/10">
              <div className="text-2xl font-bold text-white">{stat.value}</div>
              <div className="text-xs text-gray-400 mt-1">{stat.label}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}