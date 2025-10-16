package handlers

import (
    "fmt"
    "net/http"
    "os/exec"
    "strconv"
    "github.com/gin-gonic/gin"
)

type ManualAnalysisRequest struct {
    RepoURL string `json:"repo_url" binding:"required"`
}

type ManualAnalysisResponse struct {
    AnalysisID string `json:"analysis_id"`
    Status     string `json:"status"`
    Message    string `json:"message"`
}

func HandleManualAnalysis(c *gin.Context) {
    var req ManualAnalysisRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
        return
    }

    analysisID := generateAnalysisID()
    analysisStatus[analysisID] = "processing"

    // Store in analysisStorage for fix handling
    analysisStorage[analysisID] = &Analysis{
        ID:       analysisID,
        RepoURL:  req.RepoURL,
        RepoName: extractRepoName(req.RepoURL),
    }

    // Start analysis in background
    go processManualAnalysis(analysisID, req.RepoURL)

    c.JSON(http.StatusAccepted, ManualAnalysisResponse{
        AnalysisID: analysisID,
        Status:     "processing",
        Message:    "Analysis started in background",
    })
}

func processManualAnalysis(analysisID string, repoURL string) {
    fmt.Printf("🔍 [MANUAL] Starting analysis for repo: %s\n", repoURL)
    
    repoPath := "/tmp/repo_manual_scan_" + analysisID
    exec.Command("rm", "-rf", repoPath).Run()

    // Clone repository
    cloneCmd := exec.Command("git", "clone", "--depth", "1", repoURL, repoPath)
    if err := cloneCmd.Run(); err != nil {
        fmt.Printf("❌ Clone failed: %v\n", err)
        analysisStatus[analysisID] = "failed"
        return
    }

    fmt.Printf("📁 [ANALYSIS %s] Repository cloned, starting AI analysis...\n", analysisID)
    
    // Analyze with AI
    analysis, err := AnalyzeEntireCodebase(repoPath)
    if err != nil {
        fmt.Printf("❌ [ANALYSIS %s] AI Analysis failed: %v\n", analysisID, err)
        analysisStatus[analysisID] = "failed"
        return
    }

    // Store analysis in both systems
    analyses[analysisID] = analysis
    
    // Update analysisStorage with the full analysis data
    if storedAnalysis, exists := analysisStorage[analysisID]; exists {
        storedAnalysis.Risks = append(analysis.CriticalRisks, analysis.HighRisks...)
        storedAnalysis.AutoFixes = analysis.AutoFixes
        storedAnalysis.Summary = analysis.Summary
    }
    
    analysisStatus[analysisID] = "completed"
    
    fmt.Printf("✅ [ANALYSIS %s] Analysis complete: %d critical risks found\n", analysisID, len(analysis.CriticalRisks))
    
    // Cleanup
    exec.Command("rm", "-rf", repoPath).Run()
}

func GetAnalysis(c *gin.Context) {
    analysisID := c.Param("id")
    
    analysis, exists := analyses[analysisID]
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "Analysis not found"})
        return
    }

    if analysisStatus[analysisID] != "completed" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Analysis not completed"})
        return
    }

    c.JSON(http.StatusOK, analysis)
}

func GetAnalysisStatus(c *gin.Context) {
    analysisID := c.Param("id")
    
    status, exists := analysisStatus[analysisID]
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "Analysis not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": status})
}

func GetAllAnalyses(c *gin.Context) {
    // Return all completed analyses for the dashboard
    var completedAnalyses []*AIAnalysisResponse
    for id, analysis := range analyses {
        if analysisStatus[id] == "completed" {
            completedAnalyses = append(completedAnalyses, analysis)
        }
    }
    
    // Get pagination parameters from query
    page := 1
    limit := 10
    
    if pageStr := c.Query("page"); pageStr != "" {
        if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
            page = p
        }
    }
    
    if limitStr := c.Query("limit"); limitStr != "" {
        if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
            limit = l
        }
    }
    
    // Calculate pagination
    total := len(completedAnalyses)
    totalPages := (total + limit - 1) / limit // Ceiling division
    
    // Apply pagination
    start := (page - 1) * limit
    end := start + limit
    if start > total {
        start = total
    }
    if end > total {
        end = total
    }
    
    paginatedData := completedAnalyses[start:end]
    
    c.JSON(http.StatusOK, gin.H{
        "data": paginatedData,
        "pagination": gin.H{
            "page":       page,
            "limit":      limit,
            "total":      total,
            "totalPages": totalPages,
        },
    })
}

func generateAnalysisID() string {
    return fmt.Sprintf("analysis_%d", len(analyses)+1)
}
