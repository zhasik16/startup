'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { AegisApi, ApiError } from '@/lib/api';
import { AIAnalysisResponse } from '@/types';

export default function Dashboard() {
  const [analyses, setAnalyses] = useState<AIAnalysisResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [stats, setStats] = useState({
    totalScans: 0,
    criticalIssues: 0,
    autoFixes: 0,
    complianceScore: 0
  });

  useEffect(() => {
    const fetchData = async () => {
      try {
        // First check if backend is reachable
        await AegisApi.healthCheck();
        
        const response = await AegisApi.getAnalyses();
        
        // Debug the response structure
        console.log('📊 Backend response:', response);
        
        // Handle different response structures safely
        const analysesData = response.data || response || [];
        console.log('📊 Analyses data:', analysesData);
        
        setAnalyses(Array.isArray(analysesData) ? analysesData : []);
        
        // Calculate stats safely
        const totalScans = Array.isArray(analysesData) ? analysesData.length : 0;
        const criticalIssues = Array.isArray(analysesData) 
          ? analysesData.reduce((sum, analysis) => sum + (analysis.summary?.total_critical || 0), 0)
          : 0;
        const autoFixes = Array.isArray(analysesData)
          ? analysesData.reduce((sum, analysis) => sum + (analysis.auto_fixes?.length || 0), 0)
          : 0;
        const complianceScore = Array.isArray(analysesData) && analysesData.length > 0 
          ? Math.round(analysesData.reduce((sum, analysis) => sum + calculateComplianceScore(analysis), 0) / analysesData.length)
          : 0;

        setStats({ totalScans, criticalIssues, autoFixes, complianceScore });
        setError(null);
      } catch (err) {
        console.error('Failed to fetch analyses:', err);
        if (err instanceof ApiError) {
          setError(err.message);
        } else {
          setError('Failed to connect to the backend server');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen p-8">
        <div className="animate-pulse space-y-8">
          <div className="h-8 bg-white/10 rounded w-1/4"></div>
          <div className="grid grid-cols-4 gap-6">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-32 bg-white/10 rounded-xl"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen p-8 max-w-7xl mx-auto">
        <div className="bg-red-500/20 border border-red-500/30 rounded-2xl p-6 mb-6">
          <h2 className="text-xl font-bold text-white mb-2">Connection Error</h2>
          <p className="text-gray-300 mb-4">{error}</p>
          <div className="text-sm text-gray-400">
            <p>Make sure your backend server is running:</p>
            <code className="block bg-black/30 p-2 rounded mt-2">
              cd backend && go run webhook.go
            </code>
          </div>
        </div>
        
        <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
          <h2 className="text-xl font-bold text-white mb-4">Backend Status</h2>
          <div className="space-y-2 text-gray-300">
            <p>• Backend URL: {process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}</p>
            <p>• Expected endpoint: {process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/analyses</p>
            <p>• Check if backend is running on port 8080</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen p-8 max-w-7xl mx-auto">
      {/* Header */}
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-bold text-white">Security Dashboard</h1>
          <p className="text-gray-400">Overview of your codebase security</p>
        </div>
        <Link
          href="/"
          className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white px-6 py-3 rounded-xl font-semibold transition-all duration-200 transform hover:scale-105"
        >
          New Analysis
        </Link>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
          <div className="text-2xl font-bold text-white">{stats.totalScans}</div>
          <div className="text-gray-400">Total Scans</div>
        </div>
        <div className="bg-red-500/20 border border-red-500/30 rounded-2xl p-6">
          <div className="text-2xl font-bold text-white">{stats.criticalIssues}</div>
          <div className="text-red-300">Critical Issues</div>
        </div>
        <div className="bg-green-500/20 border border-green-500/30 rounded-2xl p-6">
          <div className="text-2xl font-bold text-white">{stats.autoFixes}</div>
          <div className="text-green-300">Auto-Fixes</div>
        </div>
        <div className="bg-blue-500/20 border border-blue-500/30 rounded-2xl p-6">
          <div className="text-2xl font-bold text-white">{stats.complianceScore}%</div>
          <div className="text-blue-300">Avg Compliance</div>
        </div>
      </div>

      {/* Recent Analyses */}
      <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
        <h2 className="text-xl font-bold text-white mb-6">Recent Analyses</h2>
        
        {analyses.length === 0 ? (
          <div className="text-center py-12">
            <div className="text-6xl mb-4">🛡️</div>
            <p className="text-gray-400 mb-4">No analyses yet</p>
            <p className="text-gray-500 text-sm mb-4">
              The backend is connected but no analyses have been run yet.
            </p>
            <Link
              href="/"
              className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white px-6 py-3 rounded-xl font-semibold inline-block"
            >
              Run Your First Analysis
            </Link>
          </div>
        ) : (
          <div className="space-y-4">
            {analyses.map((analysis, index) => (
              <div
                key={index}
                className="flex items-center justify-between p-4 bg-white/5 rounded-lg border border-white/10 hover:bg-white/10 transition-colors"
              >
                <div className="flex items-center space-x-4">
                  <div className={`w-3 h-3 rounded-full ${
                    analysis.summary?.total_critical > 0 ? 'bg-red-500' : 
                    analysis.summary?.total_high > 0 ? 'bg-orange-500' : 'bg-green-500'
                  }`}></div>
                  <div>
                    <div className="font-semibold text-white">
                      Analysis #{analyses.length - index}
                    </div>
                    <div className="text-sm text-gray-400">
                      {analysis.summary?.business_type || 'Unknown'} • {analysis.compliance?.standards?.join(', ') || 'No standards'}
                    </div>
                  </div>
                </div>
                
                <div className="flex items-center space-x-6">
                  <div className="text-right">
                    <div className="flex space-x-4 text-sm">
                      <span className="text-red-400">{analysis.summary?.total_critical || 0} Critical</span>
                      <span className="text-orange-400">{analysis.summary?.total_high || 0} High</span>
                      <span className="text-yellow-400">{analysis.summary?.total_medium || 0} Medium</span>
                    </div>
                    <div className="text-xs text-gray-400 mt-1">
                      {analysis.auto_fixes?.length || 0} auto-fixes available
                    </div>
                  </div>
                  
                  <Link
                    href={`/analysis/${index}`}
                    className="bg-white/10 hover:bg-white/20 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors"
                  >
                    View Details
                  </Link>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Quick Actions */}
      <div className="grid md:grid-cols-3 gap-6 mt-8">
        <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
          <div className="text-2xl mb-3">📊</div>
          <h3 className="font-semibold text-white mb-2">Compliance Reports</h3>
          <p className="text-gray-400 text-sm mb-4">
            Generate detailed compliance reports for GDPR, HIPAA, PCI-DSS
          </p>
          <button className="w-full bg-white/10 hover:bg-white/20 text-white py-2 rounded-lg text-sm font-medium transition-colors">
            Generate Report
          </button>
        </div>

        <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
          <div className="text-2xl mb-3">⚙️</div>
          <h3 className="font-semibold text-white mb-2">Rule Configuration</h3>
          <p className="text-gray-400 text-sm mb-4">
            Customize security rules and compliance requirements
          </p>
          <button className="w-full bg-white/10 hover:bg-white/20 text-white py-2 rounded-lg text-sm font-medium transition-colors">
            Configure Rules
          </button>
        </div>

        <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-white/10">
          <div className="text-2xl mb-3">🔔</div>
          <h3 className="font-semibold text-white mb-2">Alert Settings</h3>
          <p className="text-gray-400 text-sm mb-4">
            Set up notifications for critical security findings
          </p>
          <button className="w-full bg-white/10 hover:bg-white/20 text-white py-2 rounded-lg text-sm font-medium transition-colors">
            Manage Alerts
          </button>
        </div>
      </div>
    </div>
  );
}

function calculateComplianceScore(analysis: AIAnalysisResponse): number {
  let baseScore = 100;
  baseScore -= (analysis.critical_risks?.length || 0) * 25;
  baseScore -= (analysis.high_risks?.length || 0) * 15;
  baseScore -= (analysis.medium_risks?.length || 0) * 5;
  return Math.max(0, baseScore);
}