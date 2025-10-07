package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    
    "github.com/gin-gonic/gin"
)

func main() {
    // For development, use Gin's debug mode
    gin.SetMode(gin.DebugMode)
    
    router := gin.Default()
    
    // Basic health check
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Aegis AI webhook server is running!",
            "status":  "healthy",
        })
    })
    
    // Webhook endpoint for GitHub
    router.POST("/webhook", func(c *gin.Context) {
        fmt.Println("📨 Received webhook from GitHub!")
        
        // For now, just log that we received something
        c.JSON(200, gin.H{
            "message": "Webhook received successfully",
        })
    })
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Default port for development
    }
    
    log.Printf("🚀 Server starting on port %s", port)
    log.Printf("📍 Health check: http://localhost:%s/", port)
    
    if err := router.Run(":" + port); err != nil {
        log.Fatal("Could not start server: ", err)
    }
}