declare namespace NodeJS {
  interface ProcessEnv {
    // Backend Configuration
    NEXT_PUBLIC_API_URL: string;
    BACKEND_URL: string;
    
    // API Keys (optional)
    GROQ_API_KEY?: string;
    HUGGINGFACE_API_KEY?: string;
    GITHUB_TOKEN?: string;
    
    // Application Settings
    NODE_ENV: 'development' | 'production' | 'test';
    
    // Additional environment variables
    [key: string]: string | undefined;
  }
}