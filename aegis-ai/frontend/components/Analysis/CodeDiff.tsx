'use client';

import { useState } from 'react';

interface CodeDiffProps {
  original: string;
  fixed: string;
  language?: string;
}

export default function CodeDiff({ original, fixed, language = 'javascript' }: CodeDiffProps) {
  const [view, setView] = useState<'split' | 'unified'>('split');

  return (
    <div className="bg-slate-900 rounded-xl border border-white/10 overflow-hidden">
      {/* Header */}
      <div className="flex justify-between items-center px-4 py-3 bg-slate-800 border-b border-white/10">
        <div className="flex space-x-4">
          <button
            onClick={() => setView('split')}
            className={`px-3 py-1 rounded text-sm font-medium transition-colors ${
              view === 'split'
                ? 'bg-blue-500 text-white'
                : 'text-gray-400 hover:text-white'
            }`}
          >
            Split View
          </button>
          <button
            onClick={() => setView('unified')}
            className={`px-3 py-1 rounded text-sm font-medium transition-colors ${
              view === 'unified'
                ? 'bg-blue-500 text-white'
                : 'text-gray-400 hover:text-white'
            }`}
          >
            Unified View
          </button>
        </div>
        <div className="text-sm text-gray-400 capitalize">{language}</div>
      </div>

      {/* Code Content */}
      <div className={`${view === 'split' ? 'grid md:grid-cols-2' : ''} divide-x divide-white/10`}>
        {/* Original Code */}
        <div className="relative">
          <div className="absolute top-0 left-0 px-3 py-1 bg-red-500/20 text-red-300 text-xs font-medium rounded-br">
            Original
          </div>
          <pre className="p-4 text-sm text-gray-300 overflow-x-auto pt-8">
            <code>{original}</code>
          </pre>
        </div>

        {/* Fixed Code */}
        <div className="relative">
          <div className="absolute top-0 left-0 px-3 py-1 bg-green-500/20 text-green-300 text-xs font-medium rounded-br">
            Fixed
          </div>
          <pre className="p-4 text-sm text-gray-300 overflow-x-auto pt-8">
            <code>{fixed}</code>
          </pre>
        </div>
      </div>

      {/* Actions */}
      <div className="px-4 py-3 bg-slate-800 border-t border-white/10 flex justify-end space-x-3">
        <button className="px-4 py-2 bg-white/10 hover:bg-white/20 text-white rounded-lg text-sm font-medium transition-colors">
          Copy Fix
        </button>
        <button className="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg text-sm font-medium transition-colors">
          Apply Fix
        </button>
      </div>
    </div>
  );
}