package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
    "aegis-ai/handlers"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/joho/godotenv"
)

func main() {
    // Enhanced .env file loading with debugging
    fmt.Println("🔧 Loading environment variables...")
    
    // Try multiple methods to load .env file
    var envLoaded bool
    
    // Method 1: Current directory
    err := godotenv.Load()
    if err == nil {
        fmt.Println("✅ .env loaded from current directory")
        envLoaded = true
    } else {
        fmt.Printf("❌ Method 1 failed: %v\n", err)
        
        // Method 2: Explicit ./ path
        err = godotenv.Load("./.env")
        if err == nil {
            fmt.Println("✅ .env loaded from ./ path")
            envLoaded = true
        } else {
            fmt.Printf("❌ Method 2 failed: %v\n", err)
            
            // Method 3: Absolute path
            cwd, _ := os.Getwd()
            absPath := filepath.Join(cwd, ".env")
            err = godotenv.Load(absPath)
            if err == nil {
                fmt.Println("✅ .env loaded from absolute path:", absPath)
                envLoaded = true
            } else {
                fmt.Printf("❌ Method 3 failed: %v\n", err)
                fmt.Println("🚨 WARNING: .env file could not be loaded!")
            }
        }
    }
    
    // Debug: Show what environment variables are available
    fmt.Println("\n🔧 Environment Variables Check:")
    clientID := os.Getenv("GITHUB_CLIENT_ID")
    clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
    redirectURI := os.Getenv("GITHUB_REDIRECT_URI")
    
    fmt.Printf("GITHUB_CLIENT_ID: '%s' (length: %d)\n", maskString(clientID), len(clientID))
    fmt.Printf("GITHUB_CLIENT_SECRET: '%s' (length: %d)\n", maskString(clientSecret), len(clientSecret))
    fmt.Printf("GITHUB_REDIRECT_URI: '%s'\n", redirectURI)
    
    // Check if essential variables are set
    if clientID == "" || clientSecret == "" {
        fmt.Println("🚨 CRITICAL: GitHub OAuth credentials are missing!")
        fmt.Println("   Please check your .env file in the backend folder")
        fmt.Println("   Current working directory:", getCurrentDir())
        fmt.Println("   Looking for .env file in:", filepath.Join(getCurrentDir(), ".env"))
    } else {
        fmt.Println("✅ GitHub OAuth credentials are loaded")
    }
    
    router := gin.Default()
    
    // Add CORS debug middleware
    router.Use(func(c *gin.Context) {
        log.Printf("🌐 CORS Debug - Origin: %s", c.Request.Header.Get("Origin"))
        log.Printf("🌐 CORS Debug - Method: %s", c.Request.Method)
        log.Printf("🌐 CORS Debug - Path: %s", c.Request.URL.Path)
        log.Printf("🌐 CORS Debug - Authorization Header: %s", c.Request.Header.Get("Authorization"))
        log.Printf("🌐 CORS Debug - All Headers: %v", c.Request.Header)
        c.Next()
    })
    
    // Add request logging middleware
    router.Use(func(c *gin.Context) {
        log.Printf("📥 Incoming request: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr)
        log.Printf("📥 Origin header: %s", c.Request.Header.Get("Origin"))
        c.Next()
    })
    
    // FIXED CORS configuration - Complete fix for Authorization headers
    corsConfig := cors.Config{
        AllowOrigins: []string{
            "https://organic-system-4jj5767gwp6w2q9wq-3000.app.github.dev", // Your EXACT frontend URL
            "https://organic-system-4jj5767gwp6w2q9wq-8080.app.github.dev", // Your backend URL
            "http://localhost:3000",       // Fallback for local
            "https://localhost:3000",      // Added for NextAuth
            "http://127.0.0.1:3000",       // Fallback for local
        },
        AllowMethods: []string{
            "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD",
        },
        AllowHeaders: []string{
            "Origin",
            "Content-Type",
            "Authorization",     // CRITICAL: This allows Bearer tokens
            "Accept",
            "X-Requested-With",
            "X-Session-ID",      // For legacy session support
            "Access-Control-Request-Method",
            "Access-Control-Request-Headers",
        },
        ExposeHeaders: []string{
            "Authorization",
            "Content-Length",
            "Content-Type",
            "X-Total-Count",
            "Link",
        },
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }
    
    router.Use(cors.New(corsConfig))
    
    // Handle preflight OPTIONS requests explicitly
    router.OPTIONS("/*path", func(c *gin.Context) {
        log.Printf("🔧 Handling OPTIONS preflight for: %s", c.Request.URL.Path)
        c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
        c.Header("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowMethods, ", "))
        c.Header("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowHeaders, ", "))
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Max-Age", "43200") // 12 hours
        c.Status(204)
    })
    
    // Health check with environment info
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Aegis AI is running in Codespaces!",
            "status":  "healthy",
            "env_loaded": envLoaded,
            "github_configured": clientID != "" && clientSecret != "",
            "cors_enabled": true,
        })
    })
    
    // Webhook endpoint for GitHub
    router.POST("/webhook", handlers.HandleWebhook)
    
    // Frontend API endpoints
    router.POST("/api/analyze", handlers.HandleManualAnalysis)
    router.GET("/api/analysis/:id", handlers.GetAnalysis)
    router.GET("/api/analysis/:id/status", handlers.GetAnalysisStatus)
    router.GET("/api/analyses", handlers.GetAllAnalyses)
    
    // FIXED: Changed auth endpoints to /api/auth/ prefix to avoid conflicts with NextAuth
    router.GET("/api/auth/github", handlers.HandleGitHubAuth)
    router.GET("/api/auth/github/callback", handlers.HandleGitHubCallback)
    
    // Protected endpoints (require authentication)
    router.GET("/api/user/repos", handlers.AuthMiddleware(), handlers.HandleGetUserRepos)
    router.POST("/api/analysis/:id/fix/:fixIndex", handlers.AuthMiddleware(), handlers.ApplyFix)
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("🚀 Server starting on port %s in Codespaces", port)
    log.Printf("🌐 CORS configured for:")
    for _, origin := range corsConfig.AllowOrigins {
        log.Printf("   - %s", origin)
    }
    log.Printf("🔧 Auth endpoints at: /api/auth/*")
    log.Printf("🔧 API endpoints at: /api/*")
    log.Printf("🔧 CORS allows headers: %v", corsConfig.AllowHeaders)
    
    if clientID == "" || clientSecret == "" {
        log.Printf("⚠️  WARNING: GitHub OAuth not configured. Authentication will not work.")
    }
    
    router.Run(":" + port)
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

// Helper function to get current directory
func getCurrentDir() string {
    dir, err := os.Getwd()
    if err != nil {
        return "unknown"
    }
    return dir
}