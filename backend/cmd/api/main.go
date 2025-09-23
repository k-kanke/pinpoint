package main

import (
    "log"

    httpadapter "github.com/k-kanke/pinpoint/backend/app/adapter/http"
    gormrepo "github.com/k-kanke/pinpoint/backend/app/adapter/repository/gorm"
    "github.com/k-kanke/pinpoint/backend/app/config"
    pdb "github.com/k-kanke/pinpoint/backend/app/infra/db"
    "github.com/k-kanke/pinpoint/backend/app/usecase/pin"
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

    // DI wiring
    threadRepo := gormrepo.NewThreadRepository(gdb)
    commentRepo := gormrepo.NewCommentRepository(gdb)
    photoRepo := gormrepo.NewPhotoRepository(gdb)
    pinUC := pin.New(threadRepo, commentRepo, photoRepo)

    // Router
    r := httpadapter.NewRouter(cfg, pinUC)
    log.Printf("Starting API on :%s", cfg.APIPort)
    if err := r.Run(":" + cfg.APIPort); err != nil {
        log.Fatal(err)
    }
}
