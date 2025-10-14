package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"  
    "path/filepath"
    "strings"
    "time"
)

// UPDATED TYPE DEFINITIONS - Fixed JSON field names for frontend
type AIAnalysisRequest struct {
    Codebase map[string]string `json:"codebase"`
    Context  AnalysisContext   `json:"context"`
}

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

type Risk struct {
    File        string  `json:"file"`
    Line        int     `json:"line"`
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Impact      string  `json:"impact"`
    Confidence  float64 `json:"confidence"`
    CodeSnippet string  `json:"code_snippet"`
}

type AutoFix struct {
    RiskTitle    string `json:"risk_title"`
    Original     string `json:"original"`
    Fixed        string `json:"fixed"`
    Explanation  string `json:"explanation"`
    Regulation   string `json:"regulation"`
}

type AnalysisSummary struct {
    TotalCritical int      `json:"total_critical"`
    TotalHigh     int      `json:"total_high"`
    TotalMedium   int      `json:"total_medium"`
    BusinessType  string   `json:"business_type"`
    Compliance    []string `json:"compliance_requirements"`
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

type HuggingFaceResponse []struct {
    GeneratedText string `json:"generated_text"`
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

// Store analyses in memory (global variables)
var analyses = make(map[string]*AIAnalysisResponse)
var analysisStatus = make(map[string]string)

// THE MAIN AI FUNCTION - GROQ FIRST
func AnalyzeEntireCodebase(repoPath string) (*AIAnalysisResponse, error) {
    fmt.Println("🧠 AI ANALYSIS STARTED...")
    
    codebase, languages, err := extractEntireCodebase(repoPath)
    if err != nil {
        return nil, fmt.Errorf("failed to extract codebase: %v", err)
    }
    
    fmt.Printf("📁 Found %d files in languages: %v\n", len(codebase), languages)
    
    context := AnalysisContext{
        Languages:    languages,
        BusinessType: detectBusinessType(codebase),
        Requirements: detectComplianceRequirements(codebase),
    }
    
    request := AIAnalysisRequest{
        Codebase: codebase,
        Context:  context,
    }
    
    prompt := buildAIPrompt(request)
    
    // 🚀 TRY GROQ FIRST
    if os.Getenv("GROQ_API_KEY") != "" {
        fmt.Println("🚀 Trying Groq AI (fast & reliable)...")
        response, err := callGroqAI(prompt)
        if err == nil {
            // 🆕 ADD AUTO-FIXES HERE
            fixEngine := NewAutoFixEngine()
            autoFixes := fixEngine.GenerateFixes(response.CriticalRisks, codebase)
            response.AutoFixes = autoFixes
            
            fmt.Printf("✅ Groq AI analysis complete: %d critical risks found, %d auto-fixes generated\n", 
                len(response.CriticalRisks), len(autoFixes))
            return response, nil
        }
        fmt.Printf("⚠️ Groq AI failed: %v\n", err)
    }
    
    return nil, fmt.Errorf("all AI services unavailable. Please set GROQ_API_KEY or HUGGINGFACE_API_KEY")
}

// GROQ AI CALL - FAST & RELIABLE
func callGroqAI(prompt string) (*AIAnalysisResponse, error) {
    apiKey := os.Getenv("GROQ_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("GROQ_API_KEY not set")
    }
    
    fmt.Printf("🔑 Using Groq API key: %s...\n", apiKey[:8]) // Log first 8 chars for debugging
    
    // Groq models that are fast and free
    models := []string{
        "llama-3.1-8b-instant",  // Newer model name
        "llama-3.2-1b-preview",  // Smaller, faster
        "llama-3.2-3b-preview",  // Medium size
    }
    
    var lastError error
    
    for _, model := range models {
        fmt.Printf("🤖 Trying Groq model: %s\n", model)
        
        groqRequest := GroqRequest{
            Messages: []GroqMessage{
                {
                    Role:    "user",
                    Content: prompt,
                },
            },
            Model:       model,
            Temperature: 0.1,
            MaxTokens:   4000,
            TopP:        0.9,
        }
        
        jsonData, err := json.Marshal(groqRequest)
        if err != nil {
            lastError = fmt.Errorf("failed to marshal request: %v", err)
            continue
        }
        
        req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
        if err != nil {
            lastError = fmt.Errorf("failed to create request: %v", err)
            continue
        }
        
        req.Header.Set("Authorization", "Bearer "+apiKey)
        req.Header.Set("Content-Type", "application/json")
        
        client := &http.Client{Timeout: 60 * time.Second} // Increased timeout
        resp, err := client.Do(req)
        if err != nil {
            lastError = err
            fmt.Printf("❌ Groq model %s connection failed: %v\n", model, err)
            continue
        }
        defer resp.Body.Close()
        
        body, _ := io.ReadAll(resp.Body)
        
        // Debug: Print response status and body for troubleshooting
        fmt.Printf("📡 Response status: %d\n", resp.StatusCode)
        if resp.StatusCode != 200 {
            fmt.Printf("📡 Response body: %s\n", string(body))
            lastError = fmt.Errorf("Groq model %s failed with status: %s", model, resp.Status)
            fmt.Printf("❌ %s\n", lastError)
            continue
        }
        
        var groqResp GroqResponse
        if err := json.Unmarshal(body, &groqResp); err != nil {
            lastError = fmt.Errorf("failed to parse Groq response: %v", err)
            continue
        }
        
        if len(groqResp.Choices) > 0 && groqResp.Choices[0].Message.Content != "" {
            fmt.Printf("✅ Success with Groq model: %s\n", model)
            return parseAIResponse(groqResp.Choices[0].Message.Content)
        }
        
        lastError = fmt.Errorf("Groq model %s returned empty response", model)
    }
    
    return nil, fmt.Errorf("all Groq models failed: %v", lastError)
}

// HUGGING FACE CALL (FALLBACK)
func callHuggingFaceAI(request AIAnalysisRequest) (*AIAnalysisResponse, error) {
    apiKey := os.Getenv("HUGGINGFACE_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("HUGGINGFACE_API_KEY not set")
    }
    
    prompt := buildAIPrompt(request)
    
    // Try smaller, faster models that are usually loaded
    models := []string{
        "distilgpt2",                           // Small & fast, usually ready
        "microsoft/DialoGPT-medium",           // Usually ready
        "gpt2",                                 // Usually ready
    }
    
    var lastError error
    for _, model := range models {
        fmt.Printf("🤖 Trying Hugging Face model: %s\n", model)
        
        url := fmt.Sprintf("https://api-inference.huggingface.co/models/%s", model)
        
        requestBody := map[string]interface{}{
            "inputs":     prompt,
            "parameters": map[string]interface{}{
                "max_new_tokens":  1000,
                "temperature":     0.3,
                "top_p":           0.95,
                "do_sample":       true,
                "return_full_text": false,
            },
        }
        
        jsonData, _ := json.Marshal(requestBody)
        
        req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
        req.Header.Set("Authorization", "Bearer "+apiKey)
        req.Header.Set("Content-Type", "application/json")
        
        client := &http.Client{Timeout: 30 * time.Second}
        resp, err := client.Do(req)
        if err != nil {
            lastError = err
            continue
        }
        defer resp.Body.Close()
        
        if resp.StatusCode == 200 {
            body, _ := io.ReadAll(resp.Body)
            var response HuggingFaceResponse
            if err := json.Unmarshal(body, &response); err == nil && len(response) > 0 {
                if response[0].GeneratedText != "" {
                    fmt.Printf("✅ Success with Hugging Face model: %s\n", model)
                    return parseAIResponse(response[0].GeneratedText)
                }
            }
        }
        
        lastError = fmt.Errorf("model %s failed with status: %s", model, resp.Status)
    }
    
    return nil, fmt.Errorf("all Hugging Face models failed: %v", lastError)
}

// Extract EVERY code file
func extractEntireCodebase(repoPath string) (map[string]string, []string, error) {
    codebase := make(map[string]string)
    languages := make(map[string]bool)
    
    // 🚀 ONLY SCAN KEY FILES - skip tests, docs, etc.
    priorityPatterns := []string{
        "*.py", "*.js", "*.ts", "*.java", "*.go", "*.rb", "*.php", 
        "*.cpp", "*.c", "*.cs", "*.swift", "*.kt", "*.rs",
        "config.*", "*.config", "*.env*", "*.json", "*.yaml", "*.yml",
        "package.json", "requirements.txt", "pom.xml", "build.gradle",
    }
    
    var allFiles []string
    
    // Use fast file finding
    for _, pattern := range priorityPatterns {
        findCmd := exec.Command("find", repoPath, "-name", pattern, 
            "-not", "-path", "*/node_modules/*",
            "-not", "-path", "*/.git/*", 
            "-not", "-path", "*/test/*",
            "-not", "-path", "*/tests/*",
            "-not", "-path", "*/__pycache__/*",
            "-not", "-path", "*/dist/*",
            "-not", "-path", "*/build/*",
            "-not", "-path", "*/target/*")
        
        output, err := findCmd.Output()
        if err == nil {
            files := strings.Split(strings.TrimSpace(string(output)), "\n")
            allFiles = append(allFiles, files...)
        }
    }
    
    // 🚀 LIMIT TO 20 FILES MAX for speed
    fileCount := 0
    for _, file := range allFiles {
        if file == "" || fileCount >= 20 {
            break
        }
        
        content, err := os.ReadFile(file)
        if err != nil {
            continue
        }
        
        // 🚀 SKIP LARGE FILES (>100KB)
        if len(content) > 100000 {
            continue
        }
        
        relativePath := strings.TrimPrefix(file, repoPath+"/")
        codebase[relativePath] = string(content)
        
        ext := strings.ToLower(filepath.Ext(file))
        languages[ext] = true
        fileCount++
    }
    
    langSlice := make([]string, 0, len(languages))
    for lang := range languages {
        langSlice = append(langSlice, lang)
    }
    
    fmt.Printf("📁 Scanning %d/%d files for AI analysis\n", len(codebase), len(allFiles))
    return codebase, langSlice, nil
}

// Build comprehensive AI prompt
func buildAIPrompt(request AIAnalysisRequest) string {
    var codebaseStr strings.Builder
    codebaseStr.WriteString("SECURITY ANALYSIS - Focus on CRITICAL risks only:\n\n")
    
    // 🚀 Only include first 5000 chars per file to avoid token limits
    for file, content := range request.Codebase {
        if len(content) > 5000 {
            content = content[:5000] + "\n// ... [truncated for analysis]"
        }
        codebaseStr.WriteString(fmt.Sprintf("=== FILE: %s ===\n%s\n\n", file, content))
    }
    
    return fmt.Sprintf(`You are a security expert. QUICKLY analyze this code for CRITICAL security risks only.

FOCUS ON:
- Hardcoded secrets (API keys, passwords, tokens)
- PII data exposure  
- Payment data (credit cards)
- Database credentials
- AWS/cloud keys

IGNORE:
- Code style issues
- Minor best practices
- Test files

FORMAT RESPONSE AS JSON:
{
    "critical_risks": [
        {
            "file": "filename", 
            "line": 123,
            "title": "Brief risk title",
            "description": "1-sentence explanation", 
            "impact": "Business impact",
            "confidence": 0.95,
            "code_snippet": "exact problematic code line"
        }
    ],
    "high_risks": [...],
    "medium_risks": [...],
    "explanations": ["overall analysis notes"]
}

IMPORTANT: For code_snippet, provide the EXACT line of code that contains the issue.

CODE:
%s

Be concise and focus on critical risks only.`, codebaseStr.String())
}

// Existing OpenRouter request function
func makeOpenRouterRequest(aiRequest OpenRouterRequest, apiKey string) (*AIAnalysisResponse, error) {
    jsonData, _ := json.Marshal(aiRequest)
    
    req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("HTTP-Referer", "https://aegis-ai.com")
    req.Header.Set("X-Title", "Aegis AI Security Platform")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("API call failed: %v", err)
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    
    var aiResp OpenRouterResponse
    if err := json.Unmarshal(body, &aiResp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %v", err)
    }
    
    if aiResp.Error.Message != "" {
        return nil, fmt.Errorf("AI error: %s", aiResp.Error.Message)
    }
    
    if len(aiResp.Choices) == 0 {
        return nil, fmt.Errorf("no response received")
    }
    
    return parseAIResponse(aiResp.Choices[0].Message.Content)
}

// Check if file contains code
func isCodeFile(filename string) bool {
    codeExtensions := map[string]bool{
        ".py": true, ".js": true, ".ts": true, ".jsx": true, ".tsx": true,
        ".java": true, ".go": true, ".rb": true, ".php": true, ".cpp": true,
        ".c": true, ".h": true, ".hpp": true, ".cs": true, ".swift": true,
        ".kt": true, ".rs": true, ".scala": true, ".pl": true, ".r": true,
        ".m": true, ".sql": true, ".sh": true, ".yaml": true, ".yml": true,
        ".json": true, ".xml": true, ".html": true, ".css": true, ".scss": true,
        ".config": true, ".env": true, ".properties": true,
    }
    
    ext := strings.ToLower(filepath.Ext(filename))
    return codeExtensions[ext]
}

// Skip unnecessary directories
func shouldSkipDirectory(path string) bool {
    skipDirs := []string{
        ".git", "node_modules", "vendor", "dist", "build", "target",
        "__pycache__", ".next", ".nuxt", ".output", "coverage",
    }
    
    base := filepath.Base(path)
    for _, dir := range skipDirs {
        if base == dir {
            return true
        }
    }
    return false
}

// Detect business type from code patterns
func detectBusinessType(codebase map[string]string) string {
    // Analyze codebase to determine business domain
    for file, content := range codebase {
        if strings.Contains(strings.ToLower(file+content), "patient") || 
           strings.Contains(strings.ToLower(file+content), "medical") {
            return "healthcare"
        }
        if strings.Contains(strings.ToLower(file+content), "payment") ||
           strings.Contains(strings.ToLower(file+content), "invoice") {
            return "fintech"
        }
        if strings.Contains(strings.ToLower(file+content), "user") ||
           strings.Contains(strings.ToLower(file+content), "customer") {
            return "ecommerce"
        }
    }
    return "technology"
}

// Detect compliance requirements
func detectComplianceRequirements(codebase map[string]string) []string {
    var requirements []string
    
    for _, content := range codebase {
        contentLower := strings.ToLower(content)
        
        if strings.Contains(contentLower, "gdpr") || strings.Contains(contentLower, "europe") {
            requirements = append(requirements, "GDPR")
        }
        if strings.Contains(contentLower, "hipaa") || strings.Contains(contentLower, "medical") {
            requirements = append(requirements, "HIPAA")
        }
        if strings.Contains(contentLower, "pci") || strings.Contains(contentLower, "payment") {
            requirements = append(requirements, "PCI_DSS")
        }
        if strings.Contains(contentLower, "ccpa") || strings.Contains(contentLower, "california") {
            requirements = append(requirements, "CCPA")
        }
    }
    
    return unique(requirements)
}

func unique(slice []string) []string {
    keys := make(map[string]bool)
    list := []string{}
    for _, entry := range slice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}

// Parse AI response JSON
func parseAIResponse(aiResponse string) (*AIAnalysisResponse, error) {
    // Extract JSON from AI response (it might have text around it)
    start := strings.Index(aiResponse, "{")
    end := strings.LastIndex(aiResponse, "}") + 1
    
    if start == -1 || end == -1 {
        return nil, fmt.Errorf("no JSON found in AI response")
    }
    
    jsonStr := aiResponse[start:end]
    
    var response AIAnalysisResponse
    if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
        // If JSON parsing fails, create a basic response with the raw content
        fmt.Printf("⚠️  AI returned non-JSON response, using fallback\n")
        return createFallbackResponse(aiResponse), nil
    }
    
    return &response, nil
}

// Fallback if AI doesn't return proper JSON
func createFallbackResponse(aiResponse string) *AIAnalysisResponse {
    return &AIAnalysisResponse{
        CriticalRisks: []Risk{
            {
                File:        "AI Analysis",
                Line:        1,
                Title:       "AI Analysis Completed",
                Description: "The AI has analyzed your codebase. Raw response: " + aiResponse,
                Impact:      "Review required",
                Confidence:  0.9,
                CodeSnippet: aiResponse,
            },
        },
        Explanations: []string{"AI analysis completed. Review the findings above."},
        Summary: AnalysisSummary{
            TotalCritical: 1,
            BusinessType:  "unknown",
        },
    }
}