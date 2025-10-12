package handlers

import (
    "fmt"
    "os/exec"
    "github.com/gin-gonic/gin"
)

type PullRequestEvent struct {
    Action      string `json:"action"`
    Number      int    `json:"number"`
    PullRequest struct {
        HTMLURL string `json:"html_url"`
        Head    struct {
            Ref string `json:"ref"`
            SHA string `json:"sha"`  // 🚨 FIXED: Added closing quote
        } `json:"head"`  // 🚨 FIXED: Added closing brace
    } `json:"pull_request"`
    Repository struct {
        CloneURL string `json:"clone_url"`
        Name     string `json:"name"`
    } `json:"repository"`
}

func HandleWebhook(c *gin.Context) {
    fmt.Println("✅ WEBHOOK RECEIVED!")
    
    var event PullRequestEvent
    if err := c.ShouldBindJSON(&event); err != nil {
        fmt.Printf("❌ Error parsing JSON: %v\n", err)
        c.JSON(400, gin.H{"error": "Invalid JSON"})
        return
    }
    
    // 🚀 IMMEDIATE RESPONSE - process async
    if event.Action == "opened" {
        fmt.Printf("🎯 Starting ASYNC analysis for PR #%d\n", event.Number)
        
        // Process in background goroutine
        go processWithAIAsync(event)
        
        // Respond immediately to prevent timeout
        c.JSON(202, gin.H{
            "status":  "accepted", 
            "message": "AI analysis started in background",
            "pr":      event.Number,
        })
        return
    }
    
    c.JSON(200, gin.H{"status": "ignored"})
}

// Async processing
func processWithAIAsync(event PullRequestEvent) {
    fmt.Printf("🔍 [ASYNC] Starting analysis for PR #%d\n", event.Number)
    
    repoPath := "/tmp/repo_ai_scan_" + fmt.Sprintf("%d", event.Number)
    exec.Command("rm", "-rf", repoPath).Run()
    
    // Fast clone with minimal data
    cloneCmd := exec.Command("git", "clone", "--depth", "1", "--filter=blob:none", event.Repository.CloneURL, repoPath)
    if err := cloneCmd.Run(); err != nil {
        fmt.Printf("❌ Clone failed: %v\n", err)
        return
    }
    
    fmt.Printf("📁 [PR #%d] Repository cloned, starting AI analysis...\n", event.Number)
    
    // ANALYZE WITH AI
    analysis, err := AnalyzeEntireCodebase(repoPath)
    if err != nil {
        fmt.Printf("❌ [PR #%d] AI Analysis failed: %v\n", event.Number, err)
        return
    }
    
    fmt.Printf("✅ [PR #%d] Analysis complete: %d critical risks found\n", event.Number, len(analysis.CriticalRisks))
    
    // Post to GitHub
    if err := PostAIResultsToPR(event.PullRequest.HTMLURL, analysis); err != nil {
        fmt.Printf("❌ Failed to post to GitHub: %v\n", err)
    }
    
    // Cleanup
    exec.Command("rm", "-rf", repoPath).Run()
}