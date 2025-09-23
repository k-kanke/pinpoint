package httpadapter

import (
    "net/http"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/k-kanke/pinpoint/backend/app/config"
    "github.com/k-kanke/pinpoint/backend/app/usecase/pin"
)

// NewRouter configures the Gin engine with CORS and routes.
func NewRouter(cfg config.Config, pinUC pin.UseCase) *gin.Engine {
    r := gin.Default()

    // CORS
    c := cors.Config{
        AllowOrigins:     []string{cfg.FrontendOrigin},
        AllowMethods:     []string{"GET", "POST", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }
    r.Use(cors.New(c))

    // Health
    r.GET("/healthz", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    // Static files (uploaded images)
    r.Static("/static", cfg.UploadDir)

    // API routes
    api := r.Group("/api")
    {
        api.GET("/pins", GetNearbyPinsHandler(pinUC))
        api.POST("/pins", CreatePinHandler(pinUC))
        api.GET("/pins/:id", GetPinByIDHandler(pinUC))
        api.POST("/pins/:id/comments", CreateCommentHandler(pinUC))
        api.POST("/pins/:id/photos", UploadPhotoHandler(pinUC, cfg))
    }

    return r
}
