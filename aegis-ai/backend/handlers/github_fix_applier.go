package handlers

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "time"
)

type GitHubFixApplier struct {
    Token      string
    RepoURL    string
    Branch     string
}

type FixApplicationResult struct {
    Success    bool   `json:"success"`
    Message    string `json:"message"`
    CommitSHA  string `json:"commit_sha,omitempty"`
    PRURL      string `json:"pr_url,omitempty"`
    Error      string `json:"error,omitempty"`
}

func NewGitHubFixApplier(token, repoURL string) *GitHubFixApplier {
    return &GitHubFixApplier{
        Token:   token,
        RepoURL: repoURL,
        Branch:  fmt.Sprintf("security-fix-%d", time.Now().Unix()),
    }
}

// ConfigureGitUser sets up git configuration with user details
func (g *GitHubFixApplier) ConfigureGitUser(username, email string) error {
    fmt.Printf("🔧 Configuring git user: %s <%s>\n", username, email)
    
    // Set git user name
    cmd := exec.Command("git", "config", "user.name", username)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to set git user name: %v", err)
    }
    
    // Set git user email
    cmd = exec.Command("git", "config", "user.email", email)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to set git user email: %v", err)
    }
    
    fmt.Printf("✅ Git configured for user: %s <%s>\n", username, email)
    return nil
}

func (g *GitHubFixApplier) ApplyFixAndCommit(fix AutoFix, filePath string, lineNumber int) (*FixApplicationResult, error) {
    // Create temporary directory
    tempDir := fmt.Sprintf("/tmp/security-fix-%d", time.Now().Unix())
    defer os.RemoveAll(tempDir)
    
    // Clone repository
    if err := g.cloneRepo(tempDir); err != nil {
        return nil, fmt.Errorf("failed to clone repository: %v", err)
    }
    
    // Apply the fix to the specific file
    if err := g.applyFixToFile(tempDir, filePath, lineNumber, fix); err != nil {
        return nil, fmt.Errorf("failed to apply fix: %v", err)
    }
    
    // Create new branch
    if err := g.createBranch(tempDir); err != nil {
        return nil, fmt.Errorf("failed to create branch: %v", err)
    }
    
    // Commit changes
    commitSHA, err := g.commitChanges(tempDir, fix)
    if err != nil {
        return nil, fmt.Errorf("failed to commit changes: %v", err)
    }
    
    // Push changes
    if err := g.pushChanges(tempDir); err != nil {
        return nil, fmt.Errorf("failed to push changes: %v", err)
    }
    
    // Create pull request
    prURL, err := g.createPullRequest(fix)
    if err != nil {
        // If PR creation fails, we still consider it successful since changes are pushed
        fmt.Printf("⚠️ Failed to create PR: %v\n", err)
    }
    
    return &FixApplicationResult{
        Success:   true,
        Message:   fmt.Sprintf("Security fix applied and pushed to branch: %s", g.Branch),
        CommitSHA: commitSHA,
        PRURL:     prURL,
    }, nil
}

func (g *GitHubFixApplier) cloneRepo(tempDir string) error {
    // Clone with authentication
    authRepoURL := strings.Replace(g.RepoURL, "https://", fmt.Sprintf("https://%s@", g.Token), 1)
    cmd := exec.Command("git", "clone", "--depth", "1", authRepoURL, tempDir)
    return cmd.Run()
}

func (g *GitHubFixApplier) applyFixToFile(tempDir, filePath string, lineNumber int, fix AutoFix) error {
    fullPath := fmt.Sprintf("%s/%s", tempDir, filePath)
    
    // Read file content
    content, err := os.ReadFile(fullPath)
    if err != nil {
        return err
    }
    
    lines := strings.Split(string(content), "\n")
    
    // Apply the fix (simple line replacement for demo)
    // In production, you'd use more sophisticated code modification
    if lineNumber > 0 && lineNumber <= len(lines) {
        lines[lineNumber-1] = fix.Fixed
    }
    
    // Write modified content back
    modifiedContent := strings.Join(lines, "\n")
    return os.WriteFile(fullPath, []byte(modifiedContent), 0644)
}

func (g *GitHubFixApplier) createBranch(tempDir string) error {
    cmd := exec.Command("git", "checkout", "-b", g.Branch)
    cmd.Dir = tempDir
    return cmd.Run()
}

func (g *GitHubFixApplier) commitChanges(tempDir string, fix AutoFix) (string, error) {
    // Add all changes
    addCmd := exec.Command("git", "add", ".")
    addCmd.Dir = tempDir
    if err := addCmd.Run(); err != nil {
        return "", err
    }
    
    // Commit with meaningful message
    commitMsg := fmt.Sprintf("🔒 Security fix: %s\n\n%s\n\nRegulation: %s\nApplied by Aegis AI", 
        fix.RiskTitle, fix.Explanation, fix.Regulation)
    
    commitCmd := exec.Command("git", "commit", "-m", commitMsg)
    commitCmd.Dir = tempDir
    if err := commitCmd.Run(); err != nil {
        return "", err
    }
    
    // Get commit SHA
    shaCmd := exec.Command("git", "rev-parse", "HEAD")
    shaCmd.Dir = tempDir
    output, err := shaCmd.Output()
    if err != nil {
        return "", err
    }
    
    return strings.TrimSpace(string(output)), nil
}

func (g *GitHubFixApplier) pushChanges(tempDir string) error {
    cmd := exec.Command("git", "push", "origin", g.Branch)
    cmd.Dir = tempDir
    return cmd.Run()
}

func (g *GitHubFixApplier) createPullRequest(fix AutoFix) (string, error) {
    // This would use GitHub API to create a PR
    // For now, return the branch URL
    repoName := extractRepoName(g.RepoURL)
    return fmt.Sprintf("https://github.com/%s/pull/new/%s", repoName, g.Branch), nil
}

func extractRepoName(repoURL string) string {
    parts := strings.Split(repoURL, "/")
    if len(parts) >= 2 {
        return fmt.Sprintf("%s/%s", parts[len(parts)-2], parts[len(parts)-1])
    }
    return repoURL
}