'use client';

import { signIn, useSession } from "next-auth/react"
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

export default function Home() {
  const { data: session, status } = useSession()
  const router = useRouter()
  const [error, setError] = useState('')

  useEffect(() => {
    if (session) {
      router.push('/repositories')
    }
  }, [session, router])

  const handleSignIn = async () => {
    try {
      setError('')
      await signIn('github', { 
        callbackUrl: '/repositories',
        redirect: true 
      })
    } catch (err) {
      setError('Failed to sign in with GitHub')
      console.error('Sign in error:', err)
    }
  }

  // Show loading state
  if (status === 'loading') {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-purple-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-gray-300">Loading...</p>
        </div>
      </div>
    )
  }

  // Show authenticated state (redirecting)
  if (session) {
    return (
      <div className="min-h-screen flex items-center justify-center p-8">
        <div className="max-w-4xl w-full space-y-8 text-center">
          <h1 className="text-4xl font-bold text-white">Welcome to Aegis AI</h1>
          <p className="text-xl text-gray-300">Redirecting to repositories...</p>
          <div className="w-16 h-16 border-4 border-purple-500 border-t-transparent rounded-full animate-spin mx-auto"></div>
        </div>
      </div>
    )
  }

  // Show unauthenticated state (sign in button)
  return (
    <div className="min-h-screen flex items-center justify-center p-8">
      <div className="max-w-4xl w-full space-y-12 text-center">
        {/* Hero Section */}
        <div className="space-y-6">
          <h1 className="text-5xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
            Aegis AI
          </h1>
          <p className="text-xl text-gray-300 max-w-2xl mx-auto">
            Secure your code with AI-powered security analysis and automated fixes
          </p>
        </div>

        {/* Error Message */}
        {error && (
          <div className="bg-red-500/20 border border-red-500/30 rounded-xl p-4">
            <p className="text-red-300">{error}</p>
            <p className="text-red-400 text-sm mt-2">
              Check your GitHub OAuth App configuration
            </p>
          </div>
        )}

        {/* Auth Button */}
        <div className="space-y-4">
          <button
            onClick={handleSignIn}
            className="bg-gray-800 hover:bg-gray-700 border border-gray-600 text-white px-8 py-4 rounded-xl font-semibold text-lg transition-all duration-200 flex items-center justify-center space-x-3 mx-auto"
          >
            <span>🔗</span>
            <span>Connect with GitHub</span>
          </button>
          <p className="text-gray-400 text-sm">
            We'll only access your repositories to analyze and fix security issues
          </p>
        </div>

        {/* Features */}
        <div className="grid md:grid-cols-3 gap-6 mt-12">
          {[
            { icon: '🔐', title: 'Secure OAuth', desc: 'Your credentials are safe' },
            { icon: '📊', title: 'Repo Access', desc: 'Only access what you allow' },
            { icon: '🚀', title: 'Auto Commits', desc: 'Fix and commit automatically' },
          ].map((feature, index) => (
            <div key={index} className="p-6 bg-white/5 rounded-xl border border-white/10">
              <div className="text-2xl mb-3">{feature.icon}</div>
              <h3 className="font-semibold text-white mb-2">{feature.title}</h3>
              <p className="text-gray-400 text-sm">{feature.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}