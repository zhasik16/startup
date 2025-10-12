package main

import (
    "log"
    "os"
    "aegis-ai/handlers"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Aegis AI is running!",
            "status":  "healthy",
        })
    })
    
    router.POST("/webhook", handlers.HandleWebhook)
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("🚀 Server starting on port %s", port)
    router.Run(":" + port)
}