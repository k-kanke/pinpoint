package httpadapter

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
    r := gin.Default()

    // Health
    r.GET("/healthz", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    // TODO: mount API routes under /api
    // e.g., r.GET("/api/pins", handler.GetNearby)

    return r
}

