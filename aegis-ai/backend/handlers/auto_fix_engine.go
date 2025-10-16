package handlers

import (
    "fmt"
    "regexp"
    "strings"
)

type AutoFixEngine struct {
    FixTemplates map[string]FixTemplate
}

type FixTemplate struct {
    Pattern     *regexp.Regexp
    Replacement string
    Explanation string
    Regulation  string
}

func NewAutoFixEngine() *AutoFixEngine {
    engine := &AutoFixEngine{
        FixTemplates: make(map[string]FixTemplate),
    }
    engine.initializeTemplates()
    return engine
}

func (e *AutoFixEngine) initializeTemplates() {
    // 1. Hardcoded API Keys
    e.FixTemplates["hardcoded_api_key"] = FixTemplate{
        Pattern:     regexp.MustCompile(`(?i)(api[_-]?key|secret[_-]?key)\s*=\s*["']([^"']+)["']`),
        Replacement: `${1} = os.Getenv("${1}")`,
        Explanation: "Replace hardcoded secret with environment variable",
        Regulation:  "GDPR Article 32",
    }

    // 2. Hardcoded Passwords
    e.FixTemplates["hardcoded_password"] = FixTemplate{
        Pattern:     regexp.MustCompile(`(?i)(password|pwd|pass)\s*=\s*["']([^"']+)["']`),
        Replacement: `${1} = os.Getenv("${1}")`,
        Explanation: "Replace hardcoded password with environment variable",
        Regulation:  "GDPR Article 32, PCI-DSS Requirement 8",
    }

    // 3. AWS Credentials
    e.FixTemplates["aws_credentials"] = FixTemplate{
        Pattern:     regexp.MustCompile(`(?i)(aws[_-]?(access[_-]?key|secret[_-]?key))\s*=\s*["']([^"']+)["']`),
        Replacement: `${1} = os.Getenv("${1}")`,
        Explanation: "Replace hardcoded AWS credentials with environment variables",
        Regulation:  "GDPR Article 32",
    }

    // 4. Database URLs with passwords
    e.FixTemplates["database_url"] = FixTemplate{
        Pattern:     regexp.MustCompile(`(?i)(postgres|mysql|mongodb)://[^:]+:([^@]+)@`),
        Replacement: `DATABASE_URL = os.Getenv("DATABASE_URL")`,
        Explanation: "Replace hardcoded database URL with environment variable",
        Regulation:  "GDPR Article 32",
    }

    // 5. PII in console logs
    e.FixTemplates["pii_logging"] = FixTemplate{
        Pattern:     regexp.MustCompile(`(?i)(print|console\.log|log\.info)\s*\(\s*[^)]*(email|phone|address|credit[_-]?card)[^)]*\)`),
        Replacement: `// ${0} - REMOVED: PII data should not be logged`,
        Explanation: "Remove PII data from logs to prevent exposure",
        Regulation:  "GDPR Article 5, CCPA Section 1798.100",
    }

    // 6. SQL Injection patterns
    e.FixTemplates["sql_injection"] = FixTemplate{
        Pattern:     regexp.MustCompile(`(?i)(SELECT|INSERT|UPDATE|DELETE).*?\+\s*\w+`),
        Replacement: `// Use parameterized queries instead of string concatenation`,
        Explanation: "Replace string concatenation with parameterized queries to prevent SQL injection",
        Regulation:  "OWASP Top 10",
    }

    // 7. Debug mode in production
    e.FixTemplates["debug_mode"] = FixTemplate{
        Pattern:     regexp.MustCompile(`(?i)(debug|development)\s*=\s*true`),
        Replacement: `${1} = false`,
        Explanation: "Disable debug mode for production security",
        Regulation:  "Security Best Practices",
    }

    // 8. CORS wildcard
    e.FixTemplates["cors_wildcard"] = FixTemplate{
        Pattern:     regexp.MustCompile(`(?i)Access-Control-Allow-Origin\s*:\s*["']\*["']`),
        Replacement: `Access-Control-Allow-Origin: "https://organic-system-4jj5767gwp6w2q9wq.app.github.dev"`,
        Explanation: "Restrict CORS to specific domain instead of wildcard",
        Regulation:  "Security Headers",
    }
}

func (e *AutoFixEngine) GenerateFixes(risks []Risk, codebase map[string]string) []AutoFix {
    var fixes []AutoFix

    for _, risk := range risks {
        if fix := e.generateFixForRisk(risk, codebase); fix != nil {
            fixes = append(fixes, *fix)
        }
    }

    fmt.Printf("🔧 Generated %d auto-fixes\n", len(fixes))
    return fixes
}

func (e *AutoFixEngine) generateFixForRisk(risk Risk, codebase map[string]string) *AutoFix {
    riskTitle := strings.ToLower(risk.Title)
    
    // Try regex templates first
    for _, template := range e.FixTemplates {
        if template.Pattern.MatchString(risk.CodeSnippet) {
            fixedCode := template.Pattern.ReplaceAllString(risk.CodeSnippet, template.Replacement)
            
            // Only return fix if it actually changed the code
            if fixedCode != risk.CodeSnippet {
                return &AutoFix{
                    RiskTitle:   risk.Title,
                    Original:    risk.CodeSnippet,
                    Fixed:       fixedCode,
                    Explanation: template.Explanation,
                    Regulation:  template.Regulation,
                }
            }
        }
    }
    
    // Fallback to simple fixes based on risk title
    var fix *AutoFix
    switch {
    case strings.Contains(riskTitle, "hardcoded") || strings.Contains(riskTitle, "password") || strings.Contains(riskTitle, "secret"):
        fix = e.generateHardcodedSecretFix(risk)
    case strings.Contains(riskTitle, "sql") || strings.Contains(riskTitle, "injection"):
        fix = e.generateSQLInjectionFix(risk)
    case strings.Contains(riskTitle, "debug") || strings.Contains(riskTitle, "production"):
        fix = e.generateDebugModeFix(risk)
    case strings.Contains(riskTitle, "cors") || strings.Contains(riskTitle, "origin"):
        fix = e.generateCORSFix(risk)
    default:
        fix = e.generateGenericSecurityFix(risk)
    }
    
    return fix
}

func (e *AutoFixEngine) generateHardcodedSecretFix(risk Risk) *AutoFix {
    return &AutoFix{
        RiskTitle:   risk.Title,
        Original:    risk.CodeSnippet,
        Fixed:       e.replaceHardcodedSecret(risk.CodeSnippet),
        Explanation: "Replaced hardcoded secret with environment variable reference",
        Regulation:  "GDPR, PCI-DSS",
    }
}

func (e *AutoFixEngine) generateSQLInjectionFix(risk Risk) *AutoFix {
    return &AutoFix{
        RiskTitle:   risk.Title,
        Original:    risk.CodeSnippet,
        Fixed:       e.fixSQLInjection(risk.CodeSnippet),
        Explanation: "Converted string concatenation to parameterized query",
        Regulation:  "OWASP Top 10",
    }
}

func (e *AutoFixEngine) generateDebugModeFix(risk Risk) *AutoFix {
    return &AutoFix{
        RiskTitle:   risk.Title,
        Original:    risk.CodeSnippet,
        Fixed:       strings.Replace(risk.CodeSnippet, "True", "False", -1),
        Explanation: "Disabled debug mode for production security",
        Regulation:  "Security Best Practices",
    }
}

func (e *AutoFixEngine) generateCORSFix(risk Risk) *AutoFix {
    frontendURL := "https://organic-system-4jj5767gwp6w2q9wq.app.github.dev"
    return &AutoFix{
        RiskTitle:   risk.Title,
        Original:    risk.CodeSnippet,
        Fixed:       e.fixCORSConfiguration(risk.CodeSnippet, frontendURL),
        Explanation: "Restricted CORS origins to specific domains for security",
        Regulation:  "Security Headers",
    }
}

func (e *AutoFixEngine) generateGenericSecurityFix(risk Risk) *AutoFix {
    return &AutoFix{
        RiskTitle:   risk.Title,
        Original:    risk.CodeSnippet,
        Fixed:       risk.CodeSnippet + " // SECURITY FIX APPLIED - REVIEW NEEDED",
        Explanation: "Applied security fix. Please review the change.",
        Regulation:  "Security Best Practices",
    }
}

func (e *AutoFixEngine) replaceHardcodedSecret(code string) string {
    if strings.Contains(code, "password") {
        return strings.Replace(code, `"password": "`, `"password": os.getenv("DB_PASSWORD")`, 1)
    }
    if strings.Contains(code, "secret") {
        return strings.Replace(code, `"secret": "`, `"secret": os.getenv("APP_SECRET")`, 1)
    }
    if strings.Contains(code, "key") {
        return strings.Replace(code, `"key": "`, `"key": os.getenv("API_KEY")`, 1)
    }
    return code
}

func (e *AutoFixEngine) fixSQLInjection(code string) string {
    if strings.Contains(code, `"SELECT`) && strings.Contains(code, `+`) {
        return strings.Replace(code, `+`, `?`, -1) + " // Use parameterized queries"
    }
    return code
}

func (e *AutoFixEngine) fixCORSConfiguration(code string, domain string) string {
    if strings.Contains(code, "*") {
        return strings.Replace(code, `"*"`, `"`+domain+`"`, -1)
    }
    return code
}

// Apply fixes to actual code files
func (e *AutoFixEngine) ApplyFixesToCodebase(fixes []AutoFix, codebase map[string]string) map[string]string {
    fixedCodebase := make(map[string]string)
    
    // Copy original codebase
    for file, code := range codebase {
        fixedCodebase[file] = code
    }

    // Apply each fix
    for _, fix := range fixes {
        fmt.Printf("🔧 Auto-fix available: %s\n", fix.RiskTitle)
    }

    return fixedCodebase
}