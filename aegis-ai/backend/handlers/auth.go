package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"
    "github.com/gin-gonic/gin"
)

// Remove the static oauthConfig initialization
type GitHubOAuthConfig struct {
    ClientID     string
    ClientSecret string
    RedirectURI  string
}

// Initialize when needed instead of at package level
func getOAuthConfig() *GitHubOAuthConfig {
    // Use the new /api/auth/github/callback endpoint
    redirectURI := os.Getenv("GITHUB_REDIRECT_URI")
    if redirectURI == "" {
        // Fallback: construct the redirect URI for the new endpoint
        backendURL := os.Getenv("BACKEND_URL")
        if backendURL == "" {
            backendURL = "http://localhost:8080"
        }
        redirectURI = backendURL + "/api/auth/github/callback"
    }
    
    return &GitHubOAuthConfig{
        ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
        ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
        RedirectURI:  redirectURI,
    }
}

type GitHubUser struct {
    ID        int    `json:"id"`
    Login     string `json:"login"`
    Name      string `json:"name"`
    AvatarURL string `json:"avatar_url"`
    Email     string `json:"email"`
}

type GitHubRepo struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    FullName    string `json:"full_name"`
    Description string `json:"description"`
    HTMLURL     string `json:"html_url"`
    Private     bool   `json:"private"`
    Fork        bool   `json:"fork"`
}

// Store user sessions (in production, use Redis or database)
var userSessions = make(map[string]*GitHubUser)
var userTokens = make(map[string]string) // user ID -> access token

func HandleGitHubAuth(c *gin.Context) {
    oauthConfig := getOAuthConfig()
    
    // Debug: Check if environment variables are loaded
    fmt.Printf("🔧 AUTH - ClientID: '%s'\n", oauthConfig.ClientID)
    fmt.Printf("🔧 AUTH - ClientSecret: '%s'\n", maskString(oauthConfig.ClientSecret))
    fmt.Printf("🔧 AUTH - RedirectURI: '%s'\n", oauthConfig.RedirectURI)
    
    if oauthConfig.ClientID == "" {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "GitHub OAuth not configured properly - missing Client ID",
            "details": "Check if GITHUB_CLIENT_ID is set in environment variables",
        })
        return
    }

    // Redirect to GitHub OAuth
    authURL := fmt.Sprintf(
        "https://github.com/oauth/authorize?client_id=%s&redirect_uri=%s&scope=repo,user",
        oauthConfig.ClientID,
        oauthConfig.RedirectURI,
    )
    
    fmt.Printf("🔧 Redirecting to GitHub OAuth: %s\n", authURL)
    c.Redirect(http.StatusFound, authURL)
}

func HandleGitHubCallback(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
        return
    }

    oauthConfig := getOAuthConfig()
    
    // Exchange code for access token
    token, err := exchangeCodeForToken(code, oauthConfig)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token: " + err.Error()})
        return
    }

    // Get user info
    user, err := getGitHubUser(token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info: " + err.Error()})
        return
    }

    // Store user session
    sessionID := generateSessionID()
    userSessions[sessionID] = user
    userTokens[user.Login] = token

    // Redirect to frontend with session
    frontendURL := os.Getenv("FRONTEND_URL")
    if frontendURL == "" {
        frontendURL = "http://localhost:3000"
    }
    redirectURL := fmt.Sprintf("%s/auth?session=%s", frontendURL, sessionID)
    
    fmt.Printf("🔧 Redirecting to frontend: %s\n", redirectURL)
    c.Redirect(http.StatusFound, redirectURL)
}

// UPDATED: HandleGetUserRepos now supports both Bearer tokens and session IDs
func HandleGetUserRepos(c *gin.Context) {
    fmt.Printf("🔧 HandleGetUserRepos called\n")
    
    var user *GitHubUser
    var token string
    var err error

    // Try to get from Bearer token first (NextAuth)
    authHeader := c.GetHeader("Authorization")
    fmt.Printf("🔧 Authorization header: %s\n", authHeader)
    
    if strings.HasPrefix(authHeader, "Bearer ") {
        token = strings.TrimPrefix(authHeader, "Bearer ")
        fmt.Printf("🔧 Using Bearer token for repos request: %s\n", maskString(token))
        
        // Get user info from GitHub using the token
        user, err = getGitHubUser(token)
        if err != nil {
            fmt.Printf("🔧 Error getting user from token: %v\n", err)
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid GitHub token: " + err.Error()})
            return
        }
        fmt.Printf("🔧 User from token: %s (ID: %d)\n", user.Login, user.ID)
    } else {
        // Fall back to session ID (legacy system)
        sessionID := c.GetHeader("X-Session-ID")
        fmt.Printf("🔧 X-Session-ID header: %s\n", sessionID)
        
        if sessionID == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "No authentication provided. Use Authorization: Bearer <token> or X-Session-ID header"})
            return
        }

        var exists bool
        user, exists = userSessions[sessionID]
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
            return
        }

        token, exists = userTokens[user.Login]
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token found for session"})
            return
        }
        fmt.Printf("🔧 Using session ID for repos request: %s\n", sessionID)
    }

    // Get repositories using the token
    fmt.Printf("🔧 Fetching repositories for user: %s\n", user.Login)
    repos, err := getUserRepos(token)
    if err != nil {
        fmt.Printf("🔧 Error fetching repositories: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get repositories: " + err.Error()})
        return
    }

    fmt.Printf("🔧 Successfully fetched %d repositories\n", len(repos))
    c.JSON(http.StatusOK, gin.H{
        "repos": repos,
        "user":  user,
    })
}

func exchangeCodeForToken(code string, config *GitHubOAuthConfig) (string, error) {
    client := &http.Client{}
    
    // Create the request body
    reqBody := fmt.Sprintf(
        "client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
        config.ClientID,
        config.ClientSecret,
        code,
        config.RedirectURI,
    )

    // Create the request with the body
    req, err := http.NewRequest("POST", "https://github.com/oauth/access_token", strings.NewReader(reqBody))
    if err != nil {
        return "", err
    }

    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    var result struct {
        AccessToken string `json:"access_token"`
        TokenType   string `json:"token_type"`
        Scope       string `json:"scope"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "", err
    }

    return result.AccessToken, nil
}

func getGitHubUser(token string) (*GitHubUser, error) {
    client := &http.Client{}
    
    req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "token "+token)
    req.Header.Set("Accept", "application/vnd.github.v3+json")

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Check if the response is successful
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
    }

    var user GitHubUser
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }

    return &user, nil
}

func getUserRepos(token string) ([]GitHubRepo, error) {
    client := &http.Client{}
    
    req, err := http.NewRequest("GET", "https://api.github.com/user/repos?per_page=100&sort=updated", nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "token "+token)
    req.Header.Set("Accept", "application/vnd.github.v3+json")

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // Check if the response is successful
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
    }

    var repos []GitHubRepo
    if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
        return nil, err
    }

    return repos, nil
}

func generateSessionID() string {
    return fmt.Sprintf("session_%d", len(userSessions)+1)
}

// UPDATED: AuthMiddleware now supports both Bearer tokens and session IDs
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        var user *GitHubUser
        var err error

        // Try Bearer token first (NextAuth)
        authHeader := c.GetHeader("Authorization")
        if strings.HasPrefix(authHeader, "Bearer ") {
            token := strings.TrimPrefix(authHeader, "Bearer ")
            user, err = getGitHubUser(token)
            if err != nil {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid GitHub token"})
                c.Abort()
                return
            }
        } else {
            // Fall back to session ID
            sessionID := c.GetHeader("X-Session-ID")
            if sessionID == "" {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required. Use Authorization: Bearer <token> or X-Session-ID header"})
                c.Abort()
                return
            }

            var exists bool
            user, exists = userSessions[sessionID]
            if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
                c.Abort()
                return
            }
        }

        // Add user to context
        c.Set("user", user)
        c.Next()
    }
}

// Get user token for making GitHub requests
func GetUserToken(username string) (string, bool) {
    token, exists := userTokens[username]
    return token, exists
}

// Helper function to mask sensitive strings for logging
func maskString(s string) string {
    if len(s) == 0 {
        return "<empty>"
    }
    if len(s) <= 8 {
        return "***"
    }
    return s[:4] + "***" + s[len(s)-4:]
}