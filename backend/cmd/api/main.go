package main

import (
    "log"
    "net/http"

    httpadapter "github.com/k-kanke/pinpoint/backend/app/adapter/http"
    "github.com/k-kanke/pinpoint/backend/app/config"
    pdb "github.com/k-kanke/pinpoint/backend/app/infra/db"
    "github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()

    // Connect DB
    gdb, err := pdb.Connect()
    if err != nil {
        log.Fatalf("db connect error: %v", err)
    }

    // Run migrations (idempotent)
    if err := pdb.Migrate(gdb); err != nil {
        log.Fatalf("db migrate error: %v", err)
    }

    // Router
    r := httpadapter.NewRouter()
    // Keep existing minimal health endpoint under root as well
    r.GET("/healthz", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    log.Printf("Starting API on :%s", cfg.APIPort)
    if err := r.Run(":" + cfg.APIPort); err != nil {
        log.Fatal(err)
    }
}
