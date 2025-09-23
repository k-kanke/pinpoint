package httpadapter

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/k-kanke/pinpoint/backend/app/usecase/pin"
)

type createPinRequest struct {
    Title    string  `json:"title" binding:"required,min=1,max=140"`
    Body     string  `json:"body" binding:"required"`
    Lat      float64 `json:"lat" binding:"required"`
    Lon      float64 `json:"lon" binding:"required"`
    ExpiryAt string  `json:"expiry_at"` // optional ISO8601
    PhotoURL string  `json:"photo_url"` // optional for MVP (URL-based)
}

type createCommentRequest struct {
    Body string `json:"body" binding:"required"`
}

func GetNearbyPinsHandler(uc pin.UseCase) gin.HandlerFunc {
    return func(c *gin.Context) {
        var q struct {
            Lat    float64 `form:"lat" binding:"required"`
            Lon    float64 `form:"lon" binding:"required"`
            Radius int     `form:"radius"`
            Limit  int     `form:"limit"`
        }
        if err := c.BindQuery(&q); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        if q.Radius <= 0 { q.Radius = 1000 } // meters
        if q.Limit <= 0 || q.Limit > 200 { q.Limit = 100 }

        items, err := uc.GetNearby(c.Request.Context(), q.Lat, q.Lon, q.Radius, q.Limit)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, items)
    }
}

func CreatePinHandler(uc pin.UseCase) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req createPinRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        var expiry *time.Time
        if req.ExpiryAt != "" {
            if t, err := time.Parse(time.RFC3339, req.ExpiryAt); err == nil {
                expiry = &t
            } else {
                c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expiry_at format, use RFC3339"})
                return
            }
        }
        id, err := uc.CreatePin(c.Request.Context(), pin.CreatePinCommand{
            Title: req.Title,
            Body: req.Body,
            Lat: req.Lat,
            Lon: req.Lon,
            ExpiryAt: expiry,
            PhotoURL: req.PhotoURL,
        })
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, gin.H{"id": id})
    }
}

func GetPinByIDHandler(uc pin.UseCase) gin.HandlerFunc {
    return func(c *gin.Context) {
        idStr := c.Param("id")
        id, err := uuid.Parse(idStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
            return
        }
        res, err := uc.GetByID(c.Request.Context(), id)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, res)
    }
}

func CreateCommentHandler(uc pin.UseCase) gin.HandlerFunc {
    return func(c *gin.Context) {
        idStr := c.Param("id")
        threadID, err := uuid.Parse(idStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
            return
        }
        var req createCommentRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        cid, err := uc.CreateComment(c.Request.Context(), threadID, req.Body)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, gin.H{"id": cid})
    }
}

