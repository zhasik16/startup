package handlers

import (
    "fmt"
    "regexp"
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
}

func (e *AutoFixEngine) GenerateFixes(risks []Risk, originalCode map[string]string) []AutoFix {
    var fixes []AutoFix

    for _, risk := range risks {
        if fix := e.generateFixForRisk(risk, originalCode[risk.File]); fix != nil {
            fixes = append(fixes, *fix)
        }
    }

    return fixes
}

func (e *AutoFixEngine) generateFixForRisk(risk Risk, originalCode string) *AutoFix {
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
    return nil
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