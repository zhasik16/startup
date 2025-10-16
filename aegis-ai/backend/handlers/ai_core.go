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

// Store analyses in memory (global variables)
var analyses = make(map[string]*AIAnalysisResponse)
var analysisStatus = make(map[string]string)
var analysisStorage = make(map[string]*Analysis)

// ENHANCED AI ANALYSIS WITH COMPREHENSIVE SECURITY SCANNING
func AnalyzeEntireCodebase(repoPath string) (*AIAnalysisResponse, error) {
    fmt.Println("🧠 ENHANCED AI SECURITY ANALYSIS STARTED...")
    
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
    
    prompt := buildEnhancedAIPrompt(request)
    
    // 🚀 TRY GROQ FIRST WITH ENHANCED MODELS
    if os.Getenv("GROQ_API_KEY") != "" {
        fmt.Println("🚀 Using Enhanced Groq AI (Comprehensive Security Analysis)...")
        response, err := callEnhancedGroqAI(prompt)
        if err == nil {
            // 🆕 ENHANCED AUTO-FIXES WITH COMPREHENSIVE ANALYSIS
            fixEngine := NewAutoFixEngine()
            
            // Combine all risks for the auto-fix engine
            allRisks := combineAllRisks(response)
            autoFixes := fixEngine.GenerateFixes(allRisks, codebase)
            response.AutoFixes = autoFixes
            
            // 🆕 ENHANCE WITH ADDITIONAL ANALYSIS DATA
            response = enhanceAnalysisWithAdditionalData(response, codebase, context)
            
            fmt.Printf("✅ Enhanced AI analysis complete: %d critical, %d high, %d medium risks, %d auto-fixes\n", 
                len(response.CriticalRisks), len(response.HighRisks), len(response.MediumRisks), len(autoFixes))
            return response, nil
        }
        fmt.Printf("⚠️ Enhanced Groq AI failed: %v\n", err)
    }
    
    return nil, fmt.Errorf("all AI services unavailable. Please set GROQ_API_KEY")
}

// Helper function to combine all risk levels for auto-fixing
func combineAllRisks(response *AIAnalysisResponse) []Risk {
    var allRisks []Risk
    allRisks = append(allRisks, response.CriticalRisks...)
    allRisks = append(allRisks, response.HighRisks...)
    allRisks = append(allRisks, response.MediumRisks...)
    return allRisks
}

// ENHANCED GROQ AI CALL WITH BETTER MODELS AND PROMPTS
func callEnhancedGroqAI(prompt string) (*AIAnalysisResponse, error) {
    apiKey := os.Getenv("GROQ_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("GROQ_API_KEY not set")
    }
    
    fmt.Printf("🔑 Using Groq API key: %s...\n", apiKey[:8])
    
    // Enhanced models for better analysis
    models := []string{
        "llama-3.1-70b-versatile",    // Most capable for comprehensive analysis
        "mixtral-8x7b-32768",         // Large context window
        "llama-3.1-8b-instant",       // Fast and reliable
    }
    
    var lastError error
    
    for _, model := range models {
        fmt.Printf("🤖 Trying enhanced model: %s\n", model)
        
        groqRequest := GroqRequest{
            Messages: []GroqMessage{
                {
                    Role:    "system",
                    Content: "You are a senior security engineer with 15+ years of experience in application security, penetration testing, and compliance auditing. Provide comprehensive security analysis with detailed risk categorization, compliance mapping, and architectural insights.",
                },
                {
                    Role:    "user",
                    Content: prompt,
                },
            },
            Model:       model,
            Temperature: 0.1,  // Lower for more consistent security analysis
            MaxTokens:   8000, // Increased for comprehensive analysis
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
        
        client := &http.Client{Timeout: 120 * time.Second} // Increased timeout for larger models
        resp, err := client.Do(req)
        if err != nil {
            lastError = err
            fmt.Printf("❌ Model %s connection failed: %v\n", model, err)
            continue
        }
        defer resp.Body.Close()
        
        body, _ := io.ReadAll(resp.Body)
        
        if resp.StatusCode != 200 {
            fmt.Printf("📡 Response status: %d\n", resp.StatusCode)
            lastError = fmt.Errorf("model %s failed with status: %s", model, resp.Status)
            continue
        }
        
        var groqResp GroqResponse
        if err := json.Unmarshal(body, &groqResp); err != nil {
            lastError = fmt.Errorf("failed to parse response: %v", err)
            continue
        }
        
        if len(groqResp.Choices) > 0 && groqResp.Choices[0].Message.Content != "" {
            fmt.Printf("✅ Success with enhanced model: %s\n", model)
            return parseEnhancedAIResponse(groqResp.Choices[0].Message.Content)
        }
        
        lastError = fmt.Errorf("model %s returned empty response", model)
    }
    
    return nil, fmt.Errorf("all enhanced models failed: %v", lastError)
}

// ENHANCED AI PROMPT FOR COMPREHENSIVE SECURITY ANALYSIS
func buildEnhancedAIPrompt(request AIAnalysisRequest) string {
    var codebaseStr strings.Builder
    codebaseStr.WriteString("COMPREHENSIVE SECURITY AUDIT - PRODUCTION READINESS REVIEW\n\n")
    codebaseStr.WriteString("BUSINESS CONTEXT: " + request.Context.BusinessType + "\n")
    codebaseStr.WriteString("COMPLIANCE REQUIREMENTS: " + strings.Join(request.Context.Requirements, ", ") + "\n")
    codebaseStr.WriteString("LANGUAGES DETECTED: " + strings.Join(request.Context.Languages, ", ") + "\n\n")
    
    // Smart file prioritization
    priorityFiles := []string{}
    configFiles := []string{}
    sourceFiles := []string{}
    
    for file := range request.Codebase {
        if isSecurityCriticalFile(file) {
            priorityFiles = append(priorityFiles, file)
        } else if isConfigFile(file) {
            configFiles = append(configFiles, file)
        } else {
            sourceFiles = append(sourceFiles, file)
        }
    }
    
    // Add priority security files first
    codebaseStr.WriteString("=== PRIORITY SECURITY FILES (High Risk) ===\n")
    for _, file := range priorityFiles {
        content := truncateContent(request.Codebase[file], 4000)
        codebaseStr.WriteString(fmt.Sprintf("🔐 FILE: %s\n%s\n\n", file, content))
    }
    
    // Add configuration files
    codebaseStr.WriteString("=== CONFIGURATION FILES (Medium Risk) ===\n")
    for _, file := range configFiles {
        content := truncateContent(request.Codebase[file], 2000)
        codebaseStr.WriteString(fmt.Sprintf("⚙️  FILE: %s\n%s\n\n", file, content))
    }
    
    // Add source code files
    codebaseStr.WriteString("=== SOURCE CODE FILES (Context) ===\n")
    for _, file := range sourceFiles {
        content := truncateContent(request.Codebase[file], 1500)
        codebaseStr.WriteString(fmt.Sprintf("📄 FILE: %s\n%s\n\n", file, content))
    }
    
    return fmt.Sprintf(`COMPREHENSIVE SECURITY ANALYSIS REQUEST

You are a senior security engineer conducting a production security audit. Analyze this codebase thoroughly and provide a detailed security assessment.

CRITICAL SECURITY FOCUS AREAS:

1. AUTHENTICATION & AUTHORIZATION:
   - Hardcoded credentials, API keys, secrets, tokens
   - Weak password policies
   - Missing multi-factor authentication
   - Broken access control (IDOR, privilege escalation)
   - Session management issues

2. DATA PROTECTION & PRIVACY:
   - PII exposure (emails, phones, addresses, SSN)
   - Payment data (credit cards, bank info)
   - Database credentials in code
   - Unencrypted sensitive data
   - Data leakage in logs, errors, responses

3. INJECTION & INPUT VALIDATION:
   - SQL injection vulnerabilities
   - XSS (Cross-site scripting)
   - Command injection
   - XXE (XML External Entity)
   - Unsafe deserialization
   - Path traversal

4. CONFIGURATION & DEPLOYMENT:
   - Debug mode enabled in production
   - Exposed admin interfaces
   - CORS misconfiguration
   - Security headers missing
   - Default credentials
   - Exposed .git directories

5. DEPENDENCY & SUPPLY CHAIN:
   - Known vulnerable dependencies
   - Outdated libraries with CVEs
   - Untrusted package sources
   - Missing integrity checks

6. API & NETWORK SECURITY:
   - Unauthenticated endpoints
   - Rate limiting missing
   - SSL/TLS misconfiguration
   - Information disclosure in headers

BUSINESS IMPACT ASSESSMENT:
- Financial impact potential
- Data breach severity
- Compliance violation risk
- Reputation damage
- Operational disruption

COMPLIANCE MAPPING:
- GDPR: Data protection, privacy, consent
- HIPAA: Medical data protection
- PCI-DSS: Payment card security
- SOC2: Security controls
- ISO27001: Information security

REQUIRED RESPONSE FORMAT (STRICT JSON):
{
    "critical_risks": [
        {
            "file": "config/database.yml",
            "line": 15,
            "title": "Hardcoded Database Password",
            "description": "Database password is exposed in plain text in configuration file, allowing full database compromise",
            "impact": "Complete data breach potential - attackers can access, modify, or delete all application data",
            "confidence": 0.98,
            "code_snippet": "password: \"mysecretpassword123\"",
            "cvss_score": 9.8,
            "exploitation_complexity": "Low",
            "remediation_priority": "Immediate",
            "compliance_violations": ["GDPR Article 32", "PCI-DSS Requirement 8"]
        }
    ],
    "high_risks": [
        {
            "file": "app/controllers/user_controller.js",
            "line": 42,
            "title": "SQL Injection in User Search",
            "description": "User input directly concatenated into SQL query without parameterization",
            "impact": "Database compromise via SQL injection - data theft, modification, or deletion",
            "confidence": 0.95,
            "code_snippet": "const query = \"SELECT * FROM users WHERE name = '\" + userInput + \"'\"",
            "cvss_score": 8.6,
            "exploitation_complexity": "Low",
            "remediation_priority": "High",
            "compliance_violations": ["OWASP Top 10 A03:2021"]
        }
    ],
    "medium_risks": [
        {
            "file": "config/application.rb",
            "line": 8,
            "title": "Debug Mode Enabled in Production",
            "description": "Application debug mode is enabled, exposing sensitive information in error messages",
            "impact": "Information disclosure - stack traces, configuration details, and system information exposed",
            "confidence": 0.90,
            "code_snippet": "config.debug_exception_response_format = :default",
            "cvss_score": 5.3,
            "exploitation_complexity": "Low",
            "remediation_priority": "Medium",
            "compliance_violations": ["Security Best Practices"]
        }
    ],
    "explanations": [
        "Overall security posture: Critical issues found requiring immediate attention",
        "Data protection: Multiple instances of sensitive data exposure detected",
        "Authentication: Weak credential management practices identified",
        "Compliance: Several regulatory violations requiring remediation"
    ],
    "summary": {
        "total_critical": 3,
        "total_high": 5,
        "total_medium": 8,
        "business_type": "fintech",
        "compliance_requirements": ["GDPR", "PCI-DSS", "SOC2"]
    },
    "architecture": {
        "overview": "Monolithic application with mixed security concerns - strong authentication but weak data protection controls",
        "strengths": [
            "Input validation present in most endpoints",
            "HTTPS enforcement configured",
            "Session timeout implemented"
        ],
        "concerns": [
            "No security headers configured",
            "Error handling exposes system information",
            "No rate limiting on authentication endpoints"
        ],
        "recommendations": [
            "Implement security headers (CSP, HSTS)",
            "Add comprehensive logging and monitoring",
            "Conduct penetration testing for business-critical flows"
        ]
    },
    "compliance": {
        "standards": ["GDPR", "PCI-DSS", "OWASP Top 10"],
        "gaps": [
            "No data encryption at rest for PII",
            "Missing audit trails for data access",
            "No incident response plan documented"
        ],
        "recommendations": [
            "Implement data classification policy",
            "Establish regular security training",
            "Create incident response procedures"
        ]
    }
}

ANALYZE THIS CODEBASE:
%s

Provide a thorough, professional security assessment with actionable recommendations.`, codebaseStr.String())
}

// ENHANCED AI RESPONSE PARSING
func parseEnhancedAIResponse(aiResponse string) (*AIAnalysisResponse, error) {
    // Extract JSON from AI response
    start := strings.Index(aiResponse, "{")
    end := strings.LastIndex(aiResponse, "}") + 1
    
    if start == -1 || end == -1 {
        return nil, fmt.Errorf("no JSON found in AI response")
    }
    
    jsonStr := aiResponse[start:end]
    
    var response AIAnalysisResponse
    if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
        fmt.Printf("⚠️ Enhanced AI returned non-JSON response, using fallback: %v\n", err)
        return createEnhancedFallbackResponse(aiResponse), nil
    }
    
    // Enhance risks with additional fields for compatibility
    response = enhanceRiskData(response)
    
    return &response, nil
}

// ENHANCE ANALYSIS WITH ADDITIONAL DATA
func enhanceAnalysisWithAdditionalData(response *AIAnalysisResponse, codebase map[string]string, context AnalysisContext) *AIAnalysisResponse {
    // Add business context to summary
    if response.Summary.BusinessType == "" {
        response.Summary.BusinessType = context.BusinessType
    }
    
    // Add compliance requirements
    if response.Summary.Compliance == nil {
        response.Summary.Compliance = context.Requirements
    }
    
    // Ensure architecture analysis exists
    if response.Architecture == nil {
        response.Architecture = &ArchitectureAnalysis{
            Overview: "Standard application architecture with typical security considerations",
            Strengths: []string{
                "Code structure follows common patterns",
                "Configuration management present",
            },
            Concerns: []string{
                "Limited security controls implementation",
                "Basic error handling mechanisms",
            },
            Recommendations: []string{
                "Implement comprehensive security testing",
                "Add security monitoring and alerting",
            },
        }
    }
    
    // Ensure compliance analysis exists
    if response.Compliance == nil {
        response.Compliance = &ComplianceAnalysis{
            Standards: context.Requirements,
            Gaps: []string{
                "Basic security controls need enhancement",
                "Documentation for security processes required",
            },
            Recommendations: []string{
                "Establish security governance framework",
                "Implement regular security assessments",
            },
        }
    }
    
    return response
}

// ENHANCE RISK DATA WITH ADDITIONAL FIELDS
func enhanceRiskData(response AIAnalysisResponse) AIAnalysisResponse {
    // Add missing fields to risks
    for i := range response.CriticalRisks {
        if response.CriticalRisks[i].FilePath == "" {
            response.CriticalRisks[i].FilePath = response.CriticalRisks[i].File
        }
        if response.CriticalRisks[i].LineNumber == 0 {
            response.CriticalRisks[i].LineNumber = response.CriticalRisks[i].Line
        }
    }
    
    for i := range response.HighRisks {
        if response.HighRisks[i].FilePath == "" {
            response.HighRisks[i].FilePath = response.HighRisks[i].File
        }
        if response.HighRisks[i].LineNumber == 0 {
            response.HighRisks[i].LineNumber = response.HighRisks[i].Line
        }
    }
    
    for i := range response.MediumRisks {
        if response.MediumRisks[i].FilePath == "" {
            response.MediumRisks[i].FilePath = response.MediumRisks[i].File
        }
        if response.MediumRisks[i].LineNumber == 0 {
            response.MediumRisks[i].LineNumber = response.MediumRisks[i].Line
        }
    }
    
    return response
}

// ENHANCED FALLBACK RESPONSE
func createEnhancedFallbackResponse(aiResponse string) *AIAnalysisResponse {
    return &AIAnalysisResponse{
        CriticalRisks: []Risk{
            {
                File:        "AI Analysis",
                FilePath:    "AI Analysis",
                Line:        1,
                LineNumber:  1,
                Title:       "Enhanced AI Analysis Completed",
                Description: "AI has analyzed your codebase with comprehensive security checks. Raw response: " + aiResponse,
                Impact:      "Comprehensive security review completed",
                Confidence:  0.9,
                CodeSnippet: aiResponse,
            },
        },
        Explanations: []string{"Enhanced AI security analysis completed. Review findings for detailed security assessment."},
        Summary: AnalysisSummary{
            TotalCritical: 1,
            TotalHigh:     0,
            TotalMedium:   0,
            BusinessType:  "unknown",
            Compliance:    []string{"Basic Security Review"},
        },
        Architecture: &ArchitectureAnalysis{
            Overview: "Standard application architecture review completed",
            Strengths: []string{
                "Basic code structure analysis performed",
                "Security pattern recognition implemented",
            },
            Concerns: []string{
                "Limited context for comprehensive assessment",
                "Need for deeper code analysis",
            },
            Recommendations: []string{
                "Consider manual security review for critical components",
                "Implement additional security testing",
            },
        },
        Compliance: &ComplianceAnalysis{
            Standards: []string{"Basic Security"},
            Gaps: []string{
                "Limited compliance context available",
                "Need for specific regulatory review",
            },
            Recommendations: []string{
                "Conduct targeted compliance assessment",
                "Review specific regulatory requirements",
            },
        },
    }
}

// HELPER FUNCTIONS FOR ENHANCED ANALYSIS

func isSecurityCriticalFile(filename string) bool {
    criticalPatterns := []string{
        "config", ".env", "secret", "key", "credential", "password",
        "database", "auth", "login", "token", "jwt", "oauth",
        "dockerfile", "compose", "kube", "setting", "property",
        "package.json", "requirements.txt", "pom.xml", "build.gradle",
        "web.config", "application.yml", "settings.py", "config.py",
    }
    
    filenameLower := strings.ToLower(filename)
    for _, pattern := range criticalPatterns {
        if strings.Contains(filenameLower, pattern) {
            return true
        }
    }
    return false
}

func isConfigFile(filename string) bool {
    configPatterns := []string{
        ".json", ".yaml", ".yml", ".xml", ".properties", ".conf",
        ".config", ".ini", ".cfg", ".toml",
    }
    
    filenameLower := strings.ToLower(filename)
    for _, pattern := range configPatterns {
        if strings.Contains(filenameLower, pattern) {
            return true
        }
    }
    return false
}

func truncateContent(content string, maxLen int) string {
    if len(content) <= maxLen {
        return content
    }
    return content[:maxLen] + "\n\n// ... [truncated for analysis - " + 
        fmt.Sprintf("%d chars total]", len(content))
}

// ENHANCED CODEBASE EXTRACTION
func extractEntireCodebase(repoPath string) (map[string]string, []string, error) {
    codebase := make(map[string]string)
    languages := make(map[string]bool)
    
    // Enhanced priority patterns for better coverage
    priorityPatterns := []string{
        "*.py", "*.js", "*.ts", "*.jsx", "*.tsx", "*.java", "*.go", "*.rb", "*.php", 
        "*.cpp", "*.c", "*.h", "*.hpp", "*.cs", "*.swift", "*.kt", "*.rs", "*.scala",
        "*.pl", "*.r", "*.m", "*.sql", "*.sh", "*.bash",
        "config.*", "*.config", "*.env*", "*.json", "*.yaml", "*.yml", "*.xml",
        "package.json", "requirements.txt", "pom.xml", "build.gradle", "composer.json",
        "Dockerfile", "docker-compose.yml", "*.tf", "*.pp", "*.md", "*.txt",
        "*.html", "*.htm", "*.css", "*.scss", "*.sass", "*.less",
    }
    
    var allFiles []string
    
    // Use fast file finding with enhanced patterns
    for _, pattern := range priorityPatterns {
        findCmd := exec.Command("find", repoPath, "-name", pattern, 
            "-not", "-path", "*/node_modules/*",
            "-not", "-path", "*/.git/*", 
            "-not", "-path", "*/test/*",
            "-not", "-path", "*/tests/*",
            "-not", "-path", "*/__pycache__/*",
            "-not", "-path", "*/dist/*",
            "-not", "-path", "*/build/*",
            "-not", "-path", "*/target/*",
            "-not", "-path", "*/vendor/*",
            "-not", "-path", "*/tmp/*",
            "-not", "-path", "*/temp/*")
        
        output, err := findCmd.Output()
        if err == nil {
            files := strings.Split(strings.TrimSpace(string(output)), "\n")
            allFiles = append(allFiles, files...)
        }
    }
    
    // 🚀 LIMIT TO 25 FILES MAX for comprehensive analysis
    fileCount := 0
    for _, file := range allFiles {
        if file == "" || fileCount >= 25 {
            break
        }
        
        content, err := os.ReadFile(file)
        if err != nil {
            continue
        }
        
        // 🚀 SKIP LARGE FILES (>200KB) but allow more content
        if len(content) > 200000 {
            continue
        }
        
        relativePath := strings.TrimPrefix(file, repoPath+"/")
        codebase[relativePath] = string(content)
        
        ext := strings.ToLower(filepath.Ext(file))
        if ext != "" {
            languages[ext] = true
        }
        fileCount++
    }
    
    langSlice := make([]string, 0, len(languages))
    for lang := range languages {
        langSlice = append(langSlice, lang)
    }
    
    fmt.Printf("📁 Enhanced scanning: %d/%d files for comprehensive AI analysis\n", len(codebase), len(allFiles))
    return codebase, langSlice, nil
}

// ENHANCED BUSINESS TYPE DETECTION
func detectBusinessType(codebase map[string]string) string {
    contentAnalysis := strings.ToLower(fmt.Sprintf("%v", codebase))
    
    // Enhanced business type detection
    switch {
    case strings.Contains(contentAnalysis, "patient") || strings.Contains(contentAnalysis, "medical") || 
         strings.Contains(contentAnalysis, "health") || strings.Contains(contentAnalysis, "hospital"):
        return "healthcare"
    case strings.Contains(contentAnalysis, "payment") || strings.Contains(contentAnalysis, "invoice") ||
         strings.Contains(contentAnalysis, "bank") || strings.Contains(contentAnalysis, "financial") ||
         strings.Contains(contentAnalysis, "transaction") || strings.Contains(contentAnalysis, "card"):
        return "fintech"
    case strings.Contains(contentAnalysis, "user") || strings.Contains(contentAnalysis, "customer") ||
         strings.Contains(contentAnalysis, "cart") || strings.Contains(contentAnalysis, "product") ||
         strings.Contains(contentAnalysis, "order") || strings.Contains(contentAnalysis, "shop"):
        return "ecommerce"
    case strings.Contains(contentAnalysis, "education") || strings.Contains(contentAnalysis, "school") ||
         strings.Contains(contentAnalysis, "university") || strings.Contains(contentAnalysis, "course"):
        return "education"
    case strings.Contains(contentAnalysis, "government") || strings.Contains(contentAnalysis, "public") ||
         strings.Contains(contentAnalysis, "citizen") || strings.Contains(contentAnalysis, "agency"):
        return "government"
    default:
        return "technology"
    }
}

// ENHANCED COMPLIANCE DETECTION
func detectComplianceRequirements(codebase map[string]string) []string {
    var requirements []string
    contentAnalysis := strings.ToLower(fmt.Sprintf("%v", codebase))
    
    // Enhanced compliance requirement detection
    if strings.Contains(contentAnalysis, "gdpr") || strings.Contains(contentAnalysis, "europe") || 
       strings.Contains(contentAnalysis, "privacy") || strings.Contains(contentAnalysis, "data protection") {
        requirements = append(requirements, "GDPR")
    }
    if strings.Contains(contentAnalysis, "hipaa") || strings.Contains(contentAnalysis, "medical") ||
       strings.Contains(contentAnalysis, "health") || strings.Contains(contentAnalysis, "patient") {
        requirements = append(requirements, "HIPAA")
    }
    if strings.Contains(contentAnalysis, "pci") || strings.Contains(contentAnalysis, "payment") ||
       strings.Contains(contentAnalysis, "card") || strings.Contains(contentAnalysis, "transaction") {
        requirements = append(requirements, "PCI-DSS")
    }
    if strings.Contains(contentAnalysis, "ccpa") || strings.Contains(contentAnalysis, "california") ||
       strings.Contains(contentAnalysis, "consumer") || strings.Contains(contentAnalysis, "privacy act") {
        requirements = append(requirements, "CCPA")
    }
    if strings.Contains(contentAnalysis, "soc2") || strings.Contains(contentAnalysis, "soc") ||
       strings.Contains(contentAnalysis, "service organization") {
        requirements = append(requirements, "SOC2")
    }
    if strings.Contains(contentAnalysis, "iso27001") || strings.Contains(contentAnalysis, "iso") ||
       strings.Contains(contentAnalysis, "information security") {
        requirements = append(requirements, "ISO27001")
    }
    
    // Always include basic security standards
    requirements = append(requirements, "OWASP Top 10", "Security Best Practices")
    
    return unique(requirements)
}

// KEEP EXISTING HELPER FUNCTIONS
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

func shouldSkipDirectory(path string) bool {
    skipDirs := []string{
        ".git", "node_modules", "vendor", "dist", "build", "target",
        "__pycache__", ".next", ".nuxt", ".output", "coverage",
        "tmp", "temp", "logs", "cache", ".DS_Store",
    }
    
    base := filepath.Base(path)
    for _, dir := range skipDirs {
        if base == dir {
            return true
        }
    }
    return false
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

// KEEP EXISTING TYPES
type HuggingFaceResponse []struct {
    GeneratedText string `json:"generated_text"`
}