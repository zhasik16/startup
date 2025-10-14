import { 
  AIAnalysisResponse,
  Risk,
  AutoFix,
  AnalysisSummary
} from '@/types';

// Response interfaces
export interface AnalysisTriggerResponse {
  analysis_id: string;
  status: 'processing' | 'completed' | 'failed';
  message: string;
}

export interface AnalysisStatusResponse {
  status: 'processing' | 'completed' | 'failed';
}

export interface ApplyFixResponse {
  success: boolean;
  message: string;
  fix_applied?: string;
  details?: string;
  next_steps?: string;
  simulated?: boolean;
  committed?: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface GitHubRepo {
  id: number;
  name: string;
  full_name: string;
  description: string;
  html_url: string;
  private: boolean;
  fork: boolean;
}

export interface GitHubUser {
  id: number;
  login: string;
  name: string;
  avatar_url: string;
  email: string;
}

export interface UserReposResponse {
  repos: GitHubRepo[];
  user: GitHubUser;
}

// Custom error class
export class ApiError extends Error {
  constructor(
    public code: string,
    message: string,
    public details?: any
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'https://organic-system-4jj5767gwp6w2q9wq-8080.app.github.dev';

export class AegisApi {
  private static getAccessToken(): string | null {
    // This should be provided by the component using NextAuth
    return null; // We'll pass the token explicitly
  }

  private static async fetchWithErrorHandling<T>(
    url: string, 
    options: RequestInit = {},
    accessToken?: string // Add accessToken parameter
  ): Promise<T> {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    // Use Bearer token if provided (NextAuth)
    if (accessToken) {
      headers['Authorization'] = `Bearer ${accessToken}`;
    }

    if (options.headers) {
      Object.assign(headers, options.headers);
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      // Debug logging
      console.log('🔗 API Request:', {
        url,
        method: options.method,
        hasToken: !!accessToken,
        status: response.status
      });

      // Handle connection errors
      if (!response.ok) {
        // Try to parse error response, but fallback if it's not JSON
        let errorData = {};
        try {
          errorData = await response.json();
        } catch {
          errorData = { error: `HTTP ${response.status}: ${response.statusText}` };
        }
        
        throw new ApiError(
          `HTTP_${response.status}`,
          (errorData as any).error || `HTTP ${response.status}: ${response.statusText}`,
          errorData
        );
      }

      return await response.json();
    } catch (error) {
      // Handle network errors specifically
      if (error instanceof TypeError && error.message.includes('fetch')) {
        throw new ApiError(
          'NETWORK_ERROR',
          `Cannot connect to backend at ${url}. Make sure the backend server is running on ${API_BASE_URL}`
        );
      }
      
      if (error instanceof ApiError) {
        throw error;
      }
      
      throw new ApiError(
        'UNKNOWN_ERROR',
        error instanceof Error ? error.message : 'An unexpected error occurred'
      );
    }
  }

  // Auth methods
  static getGitHubAuthURL(): string {
    return `${API_BASE_URL}/auth/github`;
  }

  static async getUserRepos(accessToken?: string): Promise<UserReposResponse> {
    return this.fetchWithErrorHandling<UserReposResponse>(
      `${API_BASE_URL}/api/user/repos`,
      {},
      accessToken // Pass the access token
    );
  }

  static isAuthenticated(): boolean {
    // Check if we have an access token (simplified)
    return typeof window !== 'undefined' && !!localStorage.getItem('github_access_token');
  }

  static logout(): void {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('github_session');
      localStorage.removeItem('github_access_token');
    }
  }

  // Analysis methods - updated to accept accessToken
  static async triggerAnalysis(repoUrl: string, accessToken?: string): Promise<AnalysisTriggerResponse> {
    return this.fetchWithErrorHandling<AnalysisTriggerResponse>(`${API_BASE_URL}/api/analyze`, {
      method: 'POST',
      body: JSON.stringify({
        repo_url: repoUrl
      }),
    }, accessToken);
  }

  static async getAnalysis(analysisId: string, accessToken?: string): Promise<AIAnalysisResponse> {
    return this.fetchWithErrorHandling<AIAnalysisResponse>(`${API_BASE_URL}/api/analysis/${analysisId}`, {}, accessToken);
  }

  static async getAnalyses(page = 1, limit = 10, accessToken?: string): Promise<PaginatedResponse<AIAnalysisResponse>> {
    return this.fetchWithErrorHandling<PaginatedResponse<AIAnalysisResponse>>(
      `${API_BASE_URL}/api/analyses?page=${page}&limit=${limit}`,
      {},
      accessToken
    );
  }

  static async getAnalysisStatus(analysisId: string, accessToken?: string): Promise<AnalysisStatusResponse> {
    return this.fetchWithErrorHandling<AnalysisStatusResponse>(
      `${API_BASE_URL}/api/analysis/${analysisId}/status`,
      {},
      accessToken
    );
  }

  static async applyFix(analysisId: string, fixIndex: number, accessToken?: string): Promise<ApplyFixResponse> {
    return this.fetchWithErrorHandling<ApplyFixResponse>(
      `${API_BASE_URL}/api/analysis/${analysisId}/fix/${fixIndex}`,
      { method: 'POST' },
      accessToken
    );
  }

  // Health check method
  static async healthCheck(): Promise<{ message: string; status: string }> {
    return this.fetchWithErrorHandling<{ message: string; status: string }>(
      `${API_BASE_URL}/`
    );
  }
}