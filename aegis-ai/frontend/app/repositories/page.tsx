'use client';

import { useSession } from "next-auth/react"
import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { signIn } from "next-auth/react"

interface GitHubRepo {
  id: number;
  name: string;
  full_name: string;
  description: string;
  html_url: string;
  private: boolean;
}

interface ReposResponse {
  repos: GitHubRepo[];
  user: any;
}

export default function Repositories() {
  const { data: session, status } = useSession()
  const router = useRouter()
  const [repos, setRepos] = useState<GitHubRepo[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string>('')
  const [analyzingRepo, setAnalyzingRepo] = useState<string | null>(null)

  useEffect(() => {
    const fetchRepos = async () => {
      if (status === 'loading') {
        return
      }

      if (status === 'unauthenticated') {
        setError('Please sign in with GitHub to view your repositories')
        setLoading(false)
        return
      }

      const accessToken = session?.accessToken
      
      if (accessToken) {
        try {
          const backendUrl = 'https://organic-system-4jj5767gwp6w2q9wq-8080.app.github.dev'
          
          const response = await fetch(`${backendUrl}/api/user/repos`, {
            headers: {
              'Authorization': `Bearer ${accessToken}`,
              'Content-Type': 'application/json',
            },
          })
          
          if (!response.ok) {
            const errorText = await response.text()
            throw new Error(`Failed to load repositories: ${response.status}`)
          }
          
          const data: ReposResponse = await response.json()
          setRepos(data.repos || [])
          
        } catch (error) {
          console.error('Failed to fetch repositories:', error)
          setError('Failed to load repositories. Please try again.')
        } finally {
          setLoading(false)
        }
      } else {
        setError('Not authenticated. Please sign in with GitHub.')
        setLoading(false)
      }
    }

    fetchRepos()
  }, [session, status])

  const handleAnalyze = async (repo: GitHubRepo) => {
    const accessToken = session?.accessToken
    
    if (!accessToken) {
      setError('Not authenticated. Please sign in.')
      return
    }
    
    setAnalyzingRepo(repo.full_name)
    
    try {
      const backendUrl = 'https://organic-system-4jj5767gwp6w2q9wq-8080.app.github.dev'
      
      const response = await fetch(`${backendUrl}/api/analyze`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${accessToken}`,
        },
        body: JSON.stringify({
          repo_url: repo.html_url,
          repo_name: repo.full_name,
          repo_id: repo.id.toString()
        }),
      })

      if (!response.ok) {
        throw new Error(`Analysis failed: ${response.status}`)
      }

      const analysisData = await response.json()
      
      router.push(`/analysis/${analysisData.analysis_id}`)
      
    } catch (error) {
      console.error('Analysis failed:', error)
      setError('Failed to start analysis. Please try again.')
      setAnalyzingRepo(null)
    }
  }

  const handleSignIn = async () => {
    await signIn('github')
  }

  const handleRetry = () => {
    setError('')
    setLoading(true)
    setTimeout(() => {
      window.location.reload()
    }, 1000)
  }

  if (status === 'loading' || loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-purple-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-300">Loading your repositories...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-400 text-lg mb-4">{error}</p>
          <button 
            onClick={handleSignIn}
            className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg font-medium transition-colors"
          >
            Sign In with GitHub
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen p-8 max-w-6xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-white">Your Repositories</h1>
        <p className="text-gray-400">Select a repository to analyze for security issues</p>
        {repos.length > 0 && (
          <p className="text-gray-500 text-sm mt-2">
            Showing {repos.length} repository{repos.length !== 1 ? 's' : ''}
          </p>
        )}
      </div>

      {repos.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-400 text-lg">No repositories found</p>
          <p className="text-gray-500 mt-2">Make sure you have access to some repositories on GitHub</p>
        </div>
      ) : (
        <div className="grid gap-4">
          {repos.map((repo) => (
            <div
              key={repo.id}
              className="flex items-center justify-between p-6 bg-white/5 rounded-xl border border-white/10 hover:bg-white/10 transition-colors"
            >
              <div className="flex-1">
                <div className="flex items-center space-x-3">
                  <h3 className="font-semibold text-white text-lg">{repo.full_name}</h3>
                  {analyzingRepo === repo.full_name && (
                    <div className="w-4 h-4 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
                  )}
                </div>
                <p className="text-gray-400 mt-1">{repo.description || 'No description'}</p>
                <div className="flex items-center space-x-4 mt-2 text-sm text-gray-500">
                  <span>{repo.private ? '🔒 Private' : '🌍 Public'}</span>
                  <span>📅 Updated recently</span>
                </div>
              </div>
              <button 
                onClick={() => handleAnalyze(repo)}
                disabled={analyzingRepo !== null}
                className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 disabled:from-gray-600 disabled:to-gray-700 disabled:cursor-not-allowed text-white px-6 py-2 rounded-lg font-medium transition-colors flex items-center space-x-2"
              >
                {analyzingRepo === repo.full_name ? (
                  <>
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                    <span>Analyzing...</span>
                  </>
                ) : (
                  <span>Analyze</span>
                )}
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}