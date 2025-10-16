package handlers


// Unified Risk type used across all files
type Risk struct {
    File        string  `json:"file"`
    Line        int     `json:"line"`
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Impact      string  `json:"impact"`
    Confidence  float64 `json:"confidence"`
    CodeSnippet string  `json:"code_snippet"`
    FilePath    string  `json:"file_path,omitempty"`
    LineNumber  int     `json:"line_number,omitempty"`
}

// Unified AutoFix type with all required fields - SIMPLIFIED to match ai_core.go
type AutoFix struct {
    RiskTitle    string `json:"risk_title"`
    Original     string `json:"original"`
    Fixed        string `json:"fixed"`
    Explanation  string `json:"explanation"`
    Regulation   string `json:"regulation"`
    FilePath     string `json:"file_path"`     // ADD THIS
    LineNumber   int    `json:"line_number"`   // ADD THIS
    CommitMessage string `json:"commit_message,omitempty"`
}

// Analysis storage structure
type Analysis struct {
    ID        string    `json:"id"`
    RepoURL   string    `json:"repo_url"`
    RepoName  string    `json:"repo_name"`
    Risks     []Risk    `json:"risks"`
    AutoFixes []AutoFix `json:"auto_fixes"`
    Summary   AnalysisSummary `json:"summary"`
    CreatedAt string    `json:"created_at"`
}

// Analysis Summary
type AnalysisSummary struct {
    TotalCritical int      `json:"total_critical"`
    TotalHigh     int      `json:"total_high"`
    TotalMedium   int      `json:"total_medium"`
    BusinessType  string   `json:"business_type"`
    Compliance    []string `json:"compliance_requirements"`
}

// AI Analysis Response
type AIAnalysisResponse struct {
    CriticalRisks []Risk         `json:"critical_risks"`
    HighRisks     []Risk         `json:"high_risks"`
    MediumRisks   []Risk         `json:"medium_risks"`
    AutoFixes     []AutoFix      `json:"auto_fixes"`
    Explanations  []string       `json:"explanations"`
    Summary       AnalysisSummary `json:"summary"`
    Architecture  *ArchitectureAnalysis `json:"architecture,omitempty"`
    Compliance    *ComplianceAnalysis   `json:"compliance,omitempty"`
}

// AI Analysis Request
type AIAnalysisRequest struct {
    Codebase map[string]string `json:"codebase"`
    Context  AnalysisContext   `json:"context"`
}

type AnalysisContext struct {
    Languages    []string `json:"languages"`
    BusinessType string   `json:"business_type"`
    Requirements []string `json:"requirements"`
}

type ArchitectureAnalysis struct {
    Overview    string   `json:"overview"`
    Strengths   []string `json:"strengths"`
    Concerns    []string `json:"concerns"`
    Recommendations []string `json:"recommendations"`
}

type ComplianceAnalysis struct {
    Standards   []string `json:"standards"`
    Gaps        []string `json:"gaps"`
    Recommendations []string `json:"recommendations"`
}

// API Request/Response types
type OpenRouterRequest struct {
    Model       string                  `json:"model"`
    Messages    []OpenRouterMessage     `json:"messages"`
    Temperature float64                 `json:"temperature,omitempty"`
    MaxTokens   int                     `json:"max_tokens,omitempty"`
}

type OpenRouterMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type OpenRouterResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
    Error struct {
        Message string `json:"message"`
    } `json:"error"`
}

type GroqRequest struct {
    Messages    []GroqMessage `json:"messages"`
    Model       string        `json:"model"`
    Temperature float64       `json:"temperature,omitempty"`
    MaxTokens   int           `json:"max_tokens,omitempty"`
    TopP        float64       `json:"top_p,omitempty"`
}

type GroqMessage struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type GroqResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
    Error struct {
        Message string `json:"message"`
    } `json:"error"`
}