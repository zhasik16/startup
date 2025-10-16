package handlers

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "github.com/gin-gonic/gin"
)

// ApplyFix handles applying a specific fix to the repository
func ApplyFix(c *gin.Context) {
    analysisID := c.Param("id")
    fixIndexStr := c.Param("fixIndex")
    
    fmt.Printf("🔧 Applying fix - AnalysisID: %s, FixIndex: %s\n", analysisID, fixIndexStr)
    
    // Get user from Bearer token (NextAuth)
    authHeader := c.GetHeader("Authorization")
    if !strings.HasPrefix(authHeader, "Bearer ") {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "No Bearer token provided"})
        return
    }
    
    token := strings.TrimPrefix(authHeader, "Bearer ")
    fmt.Printf("🔧 Using Bearer token for fix: %s\n", maskString(token))
    
    // Get user info from GitHub using the token
    user, err := getGitHubUser(token)
    if err != nil {
        fmt.Printf("❌ Failed to get user from token: %v\n", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid GitHub token: " + err.Error()})
        return
    }
    
    fmt.Printf("🔧 User applying fix: %s\n", user.Login)
    
    // Convert fixIndex to integer
    fixIndex, err := strconv.Atoi(fixIndexStr)
    if err != nil || fixIndex < 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fix index"})
        return
    }
    
    // Get analysis from storage
    analysis, exists := analysisStorage[analysisID]
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "Analysis not found"})
        return
    }
    
    // Check if fix index is valid
    if fixIndex >= len(analysis.AutoFixes) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Fix index out of range"})
        return
    }
    
    fix := analysis.AutoFixes[fixIndex]
    fmt.Printf("🔧 Fix to apply: %s\n", fix.RiskTitle)
    fmt.Printf("🔧 File: '%s', line: %d\n", fix.FilePath, fix.LineNumber)
    
    // ENHANCED: Always ensure we have a valid file path
    if fix.FilePath == "" || fix.FilePath == "unknown" || fix.FilePath == ":0" {
        resolvedFilePath, resolvedLineNumber := resolveFilePath(analysis, fixIndex)
        if resolvedFilePath == "" {
            resolvedFilePath = getSmartFilePath(fix.RiskTitle)
            resolvedLineNumber = 1
            fmt.Printf("🔧 Using smart file path fallback: %s\n", resolvedFilePath)
        }
        fix.FilePath = resolvedFilePath
        fix.LineNumber = resolvedLineNumber
        fmt.Printf("🔧 Resolved file path: %s:%d\n", fix.FilePath, fix.LineNumber)
    }
    
    // Apply the fix using GitHubFixApplier
    applier := NewGitHubFixApplier(token, analysis.RepoURL)
    
    // Configure git before applying fix
    if err := applier.ConfigureGitUser(user.Login, user.Email); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to configure git: " + err.Error(),
        })
        return
    }
    
    // Apply the fix with file path and line number
    result, err := applier.ApplyFixAndCommit(fix, fix.FilePath, fix.LineNumber)
    
    if err != nil {
        fmt.Printf("❌ Failed to apply fix: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to apply fix: " + err.Error(),
        })
        return
    }
    
    fmt.Printf("✅ Fix applied successfully: %s\n", result.Message)
    
    c.JSON(http.StatusOK, gin.H{
        "success":    result.Success,
        "message":    result.Message,
        "commit_sha": result.CommitSHA,
        "pr_url":     result.PRURL,
        "branch":     applier.Branch,
    })
}

// resolveFilePath attempts to find a valid file path for the fix
func resolveFilePath(analysis *Analysis, fixIndex int) (string, int) {
    fix := analysis.AutoFixes[fixIndex]
    
    fmt.Printf("🔍 Resolving file path for fix: %s\n", fix.RiskTitle)
    
    // Strategy 1: Look for matching risks with the same title
    for _, risk := range analysis.Risks {
        if risk.Title == fix.RiskTitle {
            if risk.FilePath != "" && risk.FilePath != "unknown" {
                fmt.Printf("✅ Found matching risk with file path: %s:%d\n", risk.FilePath, risk.LineNumber)
                return risk.FilePath, risk.LineNumber
            }
        }
    }
    
    // Strategy 2: Look for risks with similar titles
    for _, risk := range analysis.Risks {
        if strings.Contains(strings.ToLower(risk.Title), strings.ToLower(fix.RiskTitle)) ||
           strings.Contains(strings.ToLower(fix.RiskTitle), strings.ToLower(risk.Title)) {
            if risk.FilePath != "" && risk.FilePath != "unknown" {
                fmt.Printf("✅ Found similar risk with file path: %s:%d\n", risk.FilePath, risk.LineNumber)
                return risk.FilePath, risk.LineNumber
            }
        }
    }
    
    // Strategy 3: Look for any risk with a valid file path
    for _, risk := range analysis.Risks {
        if risk.FilePath != "" && risk.FilePath != "unknown" {
            fmt.Printf("⚠️ Using fallback file path from any risk: %s:%d\n", risk.FilePath, risk.LineNumber)
            return risk.FilePath, risk.LineNumber
        }
    }
    
    fmt.Printf("❌ No valid file path found in risks for fix: %s\n", fix.RiskTitle)
    return "", 0
}

// getSmartFilePath returns an appropriate file path based on the risk type
func getSmartFilePath(riskTitle string) string {
    riskTitle = strings.ToLower(riskTitle)
    
    switch {
    case strings.Contains(riskTitle, "database") || strings.Contains(riskTitle, "password"):
        return "config/database.yml"
    case strings.Contains(riskTitle, "api") || strings.Contains(riskTitle, "key"):
        return "config/application.yml"
    case strings.Contains(riskTitle, "environment") || strings.Contains(riskTitle, "env"):
        return ".env.example"
    case strings.Contains(riskTitle, "docker"):
        return "Dockerfile"
    case strings.Contains(riskTitle, "package") || strings.Contains(riskTitle, "dependency"):
        return "package.json"
    case strings.Contains(riskTitle, "requirement"):
        return "requirements.txt"
    case strings.Contains(riskTitle, "python"):
        return "app.py"
    case strings.Contains(riskTitle, "javascript") || strings.Contains(riskTitle, "node"):
        return "index.js"
    case strings.Contains(riskTitle, "java"):
        return "src/main/java/Application.java"
    case strings.Contains(riskTitle, "go"):
        return "main.go"
    default:
        commonFiles := []string{
            ".env.example",
            "config.yml",
            "settings.py",
            "config.py",
            "app.config",
            "application.properties",
            "README.md",
        }
        
        for _, file := range commonFiles {
            if (strings.Contains(riskTitle, "config") && strings.Contains(file, "config")) ||
               (strings.Contains(riskTitle, "setting") && strings.Contains(file, "setting")) {
                return file
            }
        }
        
        return "config/security_fixes.yml"
    }
}