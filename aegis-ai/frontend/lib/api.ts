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

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export class AegisApi {
  private static async fetchWithErrorHandling<T>(
    url: string, 
    options: RequestInit = {}
  ): Promise<T> {
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...options.headers,
        },
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new ApiError(
          `HTTP_${response.status}`,
          errorData.error || `HTTP ${response.status}: ${response.statusText}`,
          errorData
        );
      }

      return await response.json();
    } catch (error) {
      if (error instanceof ApiError) {
        throw error;
      }
      throw new ApiError(
        'NETWORK_ERROR',
        error instanceof Error ? error.message : 'Network request failed'
      );
    }
  }

  static async triggerAnalysis(repoUrl: string): Promise<AnalysisTriggerResponse> {
    return this.fetchWithErrorHandling<AnalysisTriggerResponse>(`${API_BASE_URL}/webhook`, {
      method: 'POST',
      body: JSON.stringify({
        action: 'analysis_triggered',
        repository: { clone_url: repoUrl },
        pull_request: { html_url: repoUrl }
      }),
    });
  }

  static async getAnalysis(analysisId: string): Promise<AIAnalysisResponse> {
    return this.fetchWithErrorHandling<AIAnalysisResponse>(`${API_BASE_URL}/analysis/${analysisId}`);
  }

  static async getAnalyses(page = 1, limit = 10): Promise<PaginatedResponse<AIAnalysisResponse>> {
    return this.fetchWithErrorHandling<PaginatedResponse<AIAnalysisResponse>>(
      `${API_BASE_URL}/analyses?page=${page}&limit=${limit}`
    );
  }

  static async getAnalysisStatus(analysisId: string): Promise<AnalysisStatusResponse> {
    return this.fetchWithErrorHandling<AnalysisStatusResponse>(
      `${API_BASE_URL}/analysis/${analysisId}/status`
    );
  }

  static async applyFix(analysisId: string, fixIndex: number): Promise<ApplyFixResponse> {
    return this.fetchWithErrorHandling<ApplyFixResponse>(
      `${API_BASE_URL}/analysis/${analysisId}/fix/${fixIndex}`,
      { method: 'POST' }
    );
  }
}