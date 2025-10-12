import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  // appDir is now stable, no need for experimental config
  images: {
    domains: ['github.com'],
  },
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${process.env.BACKEND_URL || 'http://localhost:8080'}/:path*`,
      },
    ];
  },
  // Enable strict mode for better development experience
  reactStrictMode: true,
  // Compiler options for smaller bundle size
  compiler: {
    removeConsole: process.env.NODE_ENV === 'production',
  },
  // Environment variables that should be available in the browser
  env: {
    CUSTOM_KEY: process.env.CUSTOM_KEY,
  },
  // Enable SWC minification (faster than Terser)
  swcMinify: true,
  // Optimize package imports
  modularizeImports: {
    '@heroicons/react': {
      transform: '@heroicons/react/{{member}}',
    },
  },
  // Optional: Enable trailing slashes
  trailingSlash: false,
  // Optional: Enable export for static deployment
  output: 'standalone',
}

export default nextConfig