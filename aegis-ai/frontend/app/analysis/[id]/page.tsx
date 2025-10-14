'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { AegisApi } from '@/lib/api';
import { AIAnalysisResponse, Risk, AutoFix } from '@/types';
import { ApiError } from '@/lib/api';

export default function AnalysisPage() {
  const params = useParams();
  const [analysis, setAnalysis] = useState<AIAnalysisResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [status, setStatus] = useState<'processing' | 'completed' | 'failed'>('processing');
  const [selectedRisk, setSelectedRisk] = useState<Risk | null>(null);
  const [selectedFix, setSelectedFix] = useState<AutoFix | null>(null);

  useEffect(() => {
    const fetchAnalysis = async () => {
      try {
        const analysisId = params.id as string;
        
        // First check status
        const statusData = await AegisApi.getAnalysisStatus(analysisId);
        setStatus(statusData.status);

        if (statusData.status === 'completed') {
          const data = await AegisApi.getAnalysis(analysisId);
          setAnalysis(data);
          setLoading(false);
        } else if (statusData.status === 'failed') {
          setError('Analysis failed. Please try again.');
          setLoading(false);
        }

        // If still processing, set up polling
        if (statusData.status === 'processing') {
          const interval = setInterval(async () => {
            try {
              const newStatus = await AegisApi.getAnalysisStatus(analysisId);
              setStatus(newStatus.status);

              if (newStatus.status === 'completed') {
                const data = await AegisApi.getAnalysis(analysisId);
                setAnalysis(data);
                clearInterval(interval);
                setLoading(false);
              } else if (newStatus.status === 'failed') {
                setError('Analysis failed. Please try again.');
                clearInterval(interval);
                setLoading(false);
              }
            } catch (err) {
              console.error('Status check failed:', err);
            }
          }, 3000); // Poll every 3 seconds

          return () => clearInterval(interval);
        }
      } catch (err) {
        console.error('Failed to fetch analysis:', err);
        setError(err instanceof ApiError ? err.message : 'An unexpected error occurred');
        setLoading(false);
      }
    };

    fetchAnalysis();
  }, [params.id]);

  const handleApplyFix = async (fixIndex: number) => {
    if (!analysis) return;
  
    try {
      const result = await AegisApi.applyFix(params.id as string, fixIndex);
      if (result.success) {
        setSelectedFix(null);
        
        // Show success message with details (safely handle optional fields)
        const message = `✅ ${result.message}`;
        const details = result.details ? `\n\n${result.details}` : '';
        const nextSteps = result.next_steps ? `\n\n${result.next_steps}` : '';
        
        alert(message + details + nextSteps);
        
        // Refresh analysis data to show updated state
        const updatedAnalysis = await AegisApi.getAnalysis(params.id as string);
        setAnalysis(updatedAnalysis);
      } else {
        alert(`❌ Failed to apply fix: ${result.message}`);
      }
    } catch (err) {
      alert('❌ Error applying fix. Please try again.');
    }
  };

  const handleCopyFix = (fix: AutoFix) => {
    navigator.clipboard.writeText(fix.fixed);
    alert('Fix copied to clipboard!');
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center space-y-4">
          <div className="w-16 h-16 border-4 border-purple-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
          <div>
            <p className="text-gray-300 text-lg">Analyzing your codebase</p>
            <p className="text-gray-500 text-sm">This may take a few minutes...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center space-y-4">
          <div className="text-red-400 text-6xl">⚠️</div>
          <p className="text-red-400 text-xl">{error}</p>
          <button 
            onClick={() => window.location.href = '/'}
            className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white px-6 py-3 rounded-xl font-semibold transition-all duration-200"
          >
            Try Again
          </button>
        </div>
      </div>
    );
  }

  if (status === 'processing') {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center space-y-4">
          <div className="w-16 h-16 border-4 border-purple-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
          <div>
            <p className="text-gray-300 text-lg">AI Analysis in Progress</p>
            <p className="text-gray-500 text-sm">Scanning codebase for security issues...</p>
            <div className="mt-4 flex justify-center space-x-2">
              <div className="w-2 h-2 bg-purple-500 rounded-full animate-bounce" style={{ animationDelay: '0s' }}></div>
              <div className="w-2 h-2 bg-purple-500 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></div>
              <div className="w-2 h-2 bg-purple-500 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!analysis) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-400 text-xl">Analysis not found</p>
        </div>
      </div>
    );
  }

  const complianceScore = calculateComplianceScore(analysis);

  return (
    <div className="min-h-screen p-8 max-w-7xl mx-auto">
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <div className="flex items-center space-x-4">
          <div className="w-10 h-10 bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
            <span className="text-xl">🛡️</span>
          </div>
          <div>
            <h1 className="text-3xl font-bold text-white">Security Analysis</h1>
            <p className="text-gray-400">AI-powered code security assessment</p>
          </div>
        </div>
        <div className="text-right">
          <div className="text-2xl font-bold text-white">{complianceScore}/100</div>
          <div className="text-sm text-gray-400">Compliance Score</div>
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="bg-red-500/20 border border-red-500/30 rounded-xl p-6">
          <div className="text-2xl font-bold text-white">{(analysis.critical_risks || []).length}</div>
          <div className="text-red-300">Critical Risks</div>
        </div>
        <div className="bg-orange-500/20 border border-orange-500/30 rounded-xl p-6">
          <div className="text-2xl font-bold text-white">{(analysis.high_risks || []).length}</div>
          <div className="text-orange-300">High Risks</div>
        </div>
        <div className="bg-yellow-500/20 border border-yellow-500/30 rounded-xl p-6">
          <div className="text-2xl font-bold text-white">{(analysis.medium_risks || []).length}</div>
          <div className="text-yellow-300">Medium Risks</div>
        </div>
        <div className="bg-green-500/20 border border-green-500/30 rounded-xl p-6">
          <div className="text-2xl font-bold text-white">{(analysis.auto_fixes || []).length}</div>
          <div className="text-green-300">Auto-Fixes</div>
        </div>
      </div>

      <div className="grid lg:grid-cols-2 gap-8">
        {/* Risks Section */}
        <div className="space-y-6">
          {/* Critical Risks */}
          {analysis.critical_risks.length > 0 && (
            <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-red-500/30">
              <h2 className="text-xl font-bold text-white mb-4 flex items-center">
                <span className="w-3 h-3 bg-red-500 rounded-full mr-2"></span>
                Critical Risks ({(analysis.critical_risks || []).length})
              </h2>
              <div className="space-y-4">
                {(analysis.critical_risks || []).map((risk, index) => (
                  <div
                    key={index}
                    className="p-4 bg-red-500/10 rounded-lg border border-red-500/20 cursor-pointer hover:bg-red-500/15 transition-colors"
                    onClick={() => setSelectedRisk(risk)}
                  >
                    <div className="flex justify-between items-start mb-2">
                      <h3 className="font-semibold text-white">{risk.title}</h3>
                      <span className="text-red-300 text-sm font-medium">
                        {Math.round(risk.confidence * 100)}% confidence
                      </span>
                    </div>
                    <p className="text-gray-300 text-sm mb-2">{risk.description}</p>
                    <div className="text-xs text-gray-400">
                      {risk.file}:{risk.line} • {risk.impact}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* High Risks */}
          {(analysis.high_risks || []).length > 0 && (
            <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-orange-500/30">
              <h2 className="text-xl font-bold text-white mb-4 flex items-center">
                <span className="w-3 h-3 bg-orange-500 rounded-full mr-2"></span>
                High Risks ({(analysis.high_risks || []).length})
              </h2>
              <div className="space-y-3">
                {(analysis.high_risks || []).map((risk, index) => (
                  <div
                    key={index}
                    className="p-3 bg-orange-500/10 rounded-lg border border-orange-500/20 cursor-pointer hover:bg-orange-500/15 transition-colors"
                    onClick={() => setSelectedRisk(risk)}
                  >
                    <h3 className="font-semibold text-white text-sm">{risk.title}</h3>
                    <div className="text-xs text-gray-400 mt-1">
                      {risk.file}:{risk.line}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Medium Risks */}
          {(analysis.medium_risks || []).length > 0 && (
            <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-yellow-500/30">
              <h2 className="text-xl font-bold text-white mb-4 flex items-center">
                <span className="w-3 h-3 bg-yellow-500 rounded-full mr-2"></span>
                Medium Risks ({(analysis.medium_risks || []).length})
              </h2>
              <div className="space-y-3">
                {(analysis.medium_risks || []).map((risk, index) => (
                  <div
                    key={index}
                    className="p-3 bg-yellow-500/10 rounded-lg border border-yellow-500/20 cursor-pointer hover:bg-yellow-500/15 transition-colors"
                    onClick={() => setSelectedRisk(risk)}
                  >
                    <h3 className="font-semibold text-white text-sm">{risk.title}</h3>
                    <div className="text-xs text-gray-400 mt-1">
                      {risk.file}:{risk.line}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Auto-Fixes & Compliance Section */}
        <div className="space-y-6">
          {/* Auto-Fixes Section */}
          <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-green-500/30">
            <h2 className="text-xl font-bold text-white mb-4 flex items-center">
              <span className="text-green-400 mr-2">🛠️</span>
              Auto-Fix Suggestions ({(analysis.auto_fixes || []).length})
            </h2>
            <div className="space-y-4">
              {(analysis.auto_fixes || []).map((fix, index) => (
                <div
                  key={index}
                  className="p-4 bg-green-500/10 rounded-lg border border-green-500/20 cursor-pointer hover:bg-green-500/15 transition-colors"
                  onClick={() => setSelectedFix(fix)}
                >
                  <h3 className="font-semibold text-white mb-2">{fix.risk_title}</h3>
                  <p className="text-gray-300 text-sm mb-3">{fix.explanation}</p>
                  <div className="text-xs text-green-400 font-medium">{fix.regulation}</div>
                </div>
              ))}
            </div>
          </div>

          {/* Compliance */}
          {analysis.compliance && (
            <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-blue-500/30">
              <h2 className="text-xl font-bold text-white mb-4">Compliance Standards</h2>
              <div className="space-y-2">
                {(analysis.compliance.standards || []).map((standard, index) => (
                  <div key={index} className="flex items-center text-sm">
                    <span className="text-green-400 mr-2">✓</span>
                    <span className="text-gray-300">{standard}</span>
                  </div>
                ))}
              </div>
              {analysis.compliance.gaps && analysis.compliance.gaps.length > 0 && (
                <div className="mt-4">
                  <h3 className="font-semibold text-white mb-2">Compliance Gaps</h3>
                  <div className="space-y-1">
                    {(analysis.compliance.gaps || []).map((gap, index) => (
                      <div key={index} className="flex items-center text-sm">
                        <span className="text-red-400 mr-2">⚠</span>
                        <span className="text-gray-300">{gap}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Architecture Analysis */}
          {analysis.architecture && (
            <div className="bg-white/5 backdrop-blur-lg rounded-2xl p-6 border border-purple-500/30">
              <h2 className="text-xl font-bold text-white mb-4">Architecture Analysis</h2>
              <p className="text-gray-300 text-sm mb-4">{analysis.architecture.overview}</p>
              
              {analysis.architecture.strengths && analysis.architecture.strengths.length > 0 && (
                <div className="mb-4">
                  <h3 className="font-semibold text-white mb-2 text-sm">Strengths</h3>
                  <div className="space-y-1">
                    {analysis.architecture.strengths.map((strength, index) => (
                      <div key={index} className="flex items-center text-sm">
                        <span className="text-green-400 mr-2">✓</span>
                        <span className="text-gray-300">{strength}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {analysis.architecture.recommendations && analysis.architecture.recommendations.length > 0 && (
                <div>
                  <h3 className="font-semibold text-white mb-2 text-sm">Recommendations</h3>
                  <div className="space-y-1">
                    {analysis.architecture.recommendations.map((recommendation, index) => (
                      <div key={index} className="flex items-center text-sm">
                        <span className="text-blue-400 mr-2">💡</span>
                        <span className="text-gray-300">{recommendation}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Risk Detail Modal */}
      {selectedRisk && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
          <div className="bg-slate-800 rounded-2xl p-6 max-w-2xl w-full max-h-[80vh] overflow-y-auto">
            <div className="flex justify-between items-start mb-4">
              <h3 className="text-xl font-bold text-white">{selectedRisk.title}</h3>
              <button
                onClick={() => setSelectedRisk(null)}
                className="text-gray-400 hover:text-white text-xl"
              >
                ✕
              </button>
            </div>
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <h4 className="font-semibold text-white mb-2">File</h4>
                  <p className="text-gray-300 text-sm">{selectedRisk.file}</p>
                </div>
                <div>
                  <h4 className="font-semibold text-white mb-2">Line</h4>
                  <p className="text-gray-300 text-sm">{selectedRisk.line}</p>
                </div>
              </div>
              <div>
                <h4 className="font-semibold text-white mb-2">Description</h4>
                <p className="text-gray-300">{selectedRisk.description}</p>
              </div>
              <div>
                <h4 className="font-semibold text-white mb-2">Impact</h4>
                <p className="text-gray-300">{selectedRisk.impact}</p>
              </div>
              <div>
                <h4 className="font-semibold text-white mb-2">Confidence</h4>
                <div className="flex items-center space-x-2">
                  <div className="w-full bg-gray-700 rounded-full h-2">
                    <div 
                      className="bg-purple-500 h-2 rounded-full" 
                      style={{ width: `${selectedRisk.confidence * 100}%` }}
                    ></div>
                  </div>
                  <span className="text-gray-300 text-sm">{Math.round(selectedRisk.confidence * 100)}%</span>
                </div>
              </div>
              <div>
                <h4 className="font-semibold text-white mb-2">Code Snippet</h4>
                <pre className="bg-slate-900 rounded-lg p-4 text-sm text-gray-300 overflow-x-auto">
                  {selectedRisk.code_snippet}
                </pre>
              </div>
            </div>
            <div className="mt-6 flex justify-end">
              <button
                onClick={() => setSelectedRisk(null)}
                className="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Fix Detail Modal */}
      {selectedFix && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
          <div className="bg-slate-800 rounded-2xl p-6 max-w-4xl w-full max-h-[80vh] overflow-y-auto">
            <div className="flex justify-between items-start mb-6">
              <h3 className="text-xl font-bold text-white">{selectedFix.risk_title}</h3>
              <button
                onClick={() => setSelectedFix(null)}
                className="text-gray-400 hover:text-white text-xl"
              >
                ✕
              </button>
            </div>
            
            <div className="grid md:grid-cols-2 gap-6 mb-6">
              <div>
                <h4 className="font-semibold text-white mb-3">Original Code</h4>
                <pre className="bg-red-500/10 border border-red-500/30 rounded-lg p-4 text-sm text-gray-300 overflow-x-auto">
                  {selectedFix.original}
                </pre>
              </div>
              <div>
                <h4 className="font-semibold text-white mb-3">Fixed Code</h4>
                <pre className="bg-green-500/10 border border-green-500/30 rounded-lg p-4 text-sm text-gray-300 overflow-x-auto">
                  {selectedFix.fixed}
                </pre>
              </div>
            </div>

            <div className="mb-6 p-4 bg-blue-500/10 rounded-lg border border-blue-500/30">
              <h4 className="font-semibold text-white mb-2">Explanation</h4>
              <p className="text-gray-300">{selectedFix.explanation}</p>
              <div className="mt-2 text-sm text-blue-400">{selectedFix.regulation}</div>
            </div>

            <div className="flex justify-end space-x-3">
              <button 
                onClick={() => handleCopyFix(selectedFix)}
                className="px-4 py-2 bg-white/10 hover:bg-white/20 text-white rounded-lg text-sm font-medium transition-colors"
              >
                Copy Fix
              </button>
              <button 
                onClick={() => handleApplyFix(analysis.auto_fixes.indexOf(selectedFix))}
                className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg text-sm font-medium transition-colors flex items-center space-x-2"
              >
                <span>🚀 Apply Fix & Commit</span>
              </button>
            </div>
            
            <div className="mt-4 text-xs text-gray-400">
              <p>📝 This will create a new branch, apply the fix, commit changes, and push to GitHub.</p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

function calculateComplianceScore(analysis: AIAnalysisResponse): number {
  let baseScore = 100;
  
  // Safe data access with fallbacks
  const criticalRisks = analysis.critical_risks || [];
  const highRisks = analysis.high_risks || [];
  const mediumRisks = analysis.medium_risks || [];
  
  baseScore -= criticalRisks.length * 25;
  baseScore -= highRisks.length * 15;
  baseScore -= mediumRisks.length * 5;
  
  return Math.max(0, baseScore);
}