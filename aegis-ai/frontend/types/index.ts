export interface Risk {
  file: string;
  line: number;
  title: string;
  description: string;
  impact: string;
  confidence: number;
  code_snippet: string;
}

export interface AutoFix {
  risk_title: string;
  original: string;
  fixed: string;
  explanation: string;
  regulation: string;
}

export interface AnalysisSummary {
  total_critical: number;
  total_high: number;
  total_medium: number;
  business_type: string;
  compliance_requirements: string[];
}

export interface ArchitectureAnalysis {
  overview: string;
  strengths: string[];
  concerns: string[];
  recommendations: string[];
}

export interface ComplianceAnalysis {
  standards: string[];
  gaps: string[];
  recommendations: string[];
}

export interface AIAnalysisResponse {
  critical_risks: Risk[];
  high_risks: Risk[];
  medium_risks: Risk[];
  auto_fixes: AutoFix[];
  explanations: string[];
  summary: AnalysisSummary;
  architecture?: ArchitectureAnalysis;
  compliance?: ComplianceAnalysis;
}

export interface PullRequestEvent {
  action: string;
  number: number;
  pull_request: {
    html_url: string;
    head: {
      ref: string;
      sha: string;
    };
  };
  repository: {
    clone_url: string;
    name: string;
  };
}

