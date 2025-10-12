package handlers

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"
)

// GitHub API structures
type GitHubComment struct {
    Body string `json:"body"`
}

type GitHubAppAuth struct {
    Token string `json:"token"`
}

// Post AI results to GitHub PR
func PostAIResultsToPR(prHTMLURL string, analysis *AIAnalysisResponse) error {
    fmt.Printf("📤 Posting AI results to GitHub PR: %s\n", prHTMLURL)
    
    // Extract PR number from URL (we don't need repo name for logging)
    _, prNumber, err := extractRepoAndPR(prHTMLURL) // Use _ to ignore repo name
    if err != nil {
        return fmt.Errorf("failed to parse PR URL: %v", err)
    }
    
    // Create the comment content
    commentBody := createGitHubComment(analysis, prNumber)
    
    // For now, we'll log it (next step: actual API call)
    fmt.Printf("💬 GITHUB COMMENT READY:\n%s\n", commentBody)
    
    return nil
}

// Create beautiful GitHub comment with AI findings
func createGitHubComment(analysis *AIAnalysisResponse, prNumber int) string {
    var comment strings.Builder
    
    // Header
    comment.WriteString("## 🛡️ Aegis AI Security Analysis\n\n")
    comment.WriteString("🤖 **AI-Powered Security Scan Results**\n\n")
    
    // Summary
    comment.WriteString("### 📊 Executive Summary\n")
    comment.WriteString(fmt.Sprintf("- **Critical Risks**: %d\n", len(analysis.CriticalRisks)))
    comment.WriteString(fmt.Sprintf("- **High Risks**: %d\n", len(analysis.HighRisks)))
    comment.WriteString(fmt.Sprintf("- **Medium Risks**: %d\n", len(analysis.MediumRisks)))
    comment.WriteString(fmt.Sprintf("- **Auto-Fixes Provided**: %d\n\n", len(analysis.AutoFixes)))
    
    // Compliance Score
    complianceScore := calculateComplianceScore(analysis)
    comment.WriteString(fmt.Sprintf("### 📈 Compliance Score: %d/100\n\n", complianceScore))
    
    // Critical Risks
    if len(analysis.CriticalRisks) > 0 {
        comment.WriteString("### 🚨 Critical Security Risks\n\n")
        for i, risk := range analysis.CriticalRisks {
            comment.WriteString(fmt.Sprintf("#### %d. %s\n", i+1, risk.Title))
            comment.WriteString(fmt.Sprintf("- **File**: `%s:%d`\n", risk.File, risk.Line))
            comment.WriteString(fmt.Sprintf("- **Confidence**: %.0f%%\n", risk.Confidence*100))
            comment.WriteString(fmt.Sprintf("- **Impact**: %s\n", risk.Impact))
            comment.WriteString(fmt.Sprintf("- **Description**: %s\n", risk.Description))
            
            // Code snippet
            if risk.CodeSnippet != "" {
                comment.WriteString("```\n")
                comment.WriteString(risk.CodeSnippet)
                comment.WriteString("\n```\n")
            }
            comment.WriteString("\n")
        }
    }
    
    // High Risks
    if len(analysis.HighRisks) > 0 {
        comment.WriteString("### ⚠️ High Security Risks\n\n")
        for i, risk := range analysis.HighRisks {
            comment.WriteString(fmt.Sprintf("%d. **%s** - `%s:%d` (%.0f%% confidence)\n", 
                i+1, risk.Title, risk.File, risk.Line, risk.Confidence*100))
            comment.WriteString(fmt.Sprintf("   - %s\n", risk.Description))
        }
        comment.WriteString("\n")
    }
    
    // Auto-Fixes
    if len(analysis.AutoFixes) > 0 {
        comment.WriteString("### 🛠️ Auto-Fix Suggestions\n\n")
        for i, fix := range analysis.AutoFixes {
            comment.WriteString(fmt.Sprintf("#### Fix %d: %s\n", i+1, fix.RiskTitle))
            comment.WriteString("**Original Code:**\n")
            comment.WriteString("```\n")
            comment.WriteString(fix.Original)
            comment.WriteString("\n```\n")
            comment.WriteString("**Fixed Code:**\n")
            comment.WriteString("```\n")
            comment.WriteString(fix.Fixed)
            comment.WriteString("\n```\n")
            comment.WriteString(fmt.Sprintf("**Explanation**: %s\n\n", fix.Explanation))
        }
    }
    
    // Architecture & Compliance
    if analysis.Architecture != nil {
        comment.WriteString("### 🏗️ Architecture Analysis\n")
        comment.WriteString(analysis.Architecture.Overview + "\n\n")
    }
    
    if analysis.Compliance != nil {
        comment.WriteString("### 📋 Compliance Report\n")
        comment.WriteString(strings.Join(analysis.Compliance.Standards, ", ") + "\n\n")
    }
    
    // Footer
    comment.WriteString("---\n")
    comment.WriteString("🔍 **Powered by Aegis AI** - Automated security scanning for modern development teams\n")
    
    return comment.String()
}

// Calculate compliance score based on findings
func calculateComplianceScore(analysis *AIAnalysisResponse) int {
    baseScore := 100
    baseScore -= len(analysis.CriticalRisks) * 25
    baseScore -= len(analysis.HighRisks) * 15
    baseScore -= len(analysis.MediumRisks) * 5
    if baseScore < 0 {
        return 0
    }
    return baseScore
}

// Extract repo and PR number from GitHub URL
func extractRepoAndPR(prHTMLURL string) (string, int, error) {
    // Example: https://github.com/owner/repo/pull/123
    parts := strings.Split(prHTMLURL, "/")
    if len(parts) < 7 {
        return "", 0, fmt.Errorf("invalid PR URL: %s", prHTMLURL)
    }
    
    owner := parts[3]
    repoName := parts[4]
    prNumber, err := strconv.Atoi(parts[6])
    if err != nil {
        return "", 0, fmt.Errorf("invalid PR number: %s", parts[6])
    }
    
    return owner + "/" + repoName, prNumber, nil
}

// Actual GitHub API call (commented until we have auth)
func postCommentToGitHub(repo string, prNumber int, commentBody string) error {
    // This requires GitHub App authentication
    // We'll implement this after we set up the GitHub App properly
    
    url := fmt.Sprintf("https://api.github.com/repos/%s/issues/%d/comments", repo, prNumber)
    
    comment := GitHubComment{Body: commentBody}
    jsonData, _ := json.Marshal(comment)
    
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    req.Header.Set("Authorization", "Bearer YOUR_GITHUB_TOKEN")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/vnd.github.v3+json")
    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("GitHub API call failed: %v", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 201 {
        return fmt.Errorf("GitHub API error: %s", resp.Status)
    }
    
    fmt.Printf("✅ Comment posted to GitHub PR #%d\n", prNumber)
    return nil
}