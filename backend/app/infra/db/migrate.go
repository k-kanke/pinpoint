package db

import (
    "fmt"

    "github.com/k-kanke/pinpoint/backend/app/domain/model"
    "gorm.io/gorm"
)

// Migrate runs database migrations for all models and applies PostGIS-specific changes.
func Migrate(gdb *gorm.DB) error {
    // Base tables via AutoMigrate (except geometry column specifics)
    if err := gdb.AutoMigrate(&model.User{}, &model.Thread{}, &model.Photo{}, &model.Comment{}); err != nil {
        return fmt.Errorf("automigrate failed: %w", err)
    }

    // Ensure PostGIS geometry(Point,4326) column exists on threads.location
    // and a GIST index for spatial queries
    alter := `ALTER TABLE threads
              ADD COLUMN IF NOT EXISTS location geometry(Point,4326);`
    if err := gdb.Exec(alter).Error; err != nil {
        return fmt.Errorf("add location column failed: %w", err)
    }

    // Create spatial index
    idx := `CREATE INDEX IF NOT EXISTS idx_threads_location
            ON threads USING GIST (location);`
    if err := gdb.Exec(idx).Error; err != nil {
        return fmt.Errorf("create spatial index failed: %w", err)
    }

    // Helpful composite indexes
    if err := gdb.Exec(`CREATE INDEX IF NOT EXISTS idx_comments_thread_id ON comments(thread_id);`).Error; err != nil {
        return fmt.Errorf("create idx_comments_thread_id failed: %w", err)
    }
    if err := gdb.Exec(`CREATE INDEX IF NOT EXISTS idx_photos_thread_id ON photos(thread_id);`).Error; err != nil {
        return fmt.Errorf("create idx_photos_thread_id failed: %w", err)
    }

    return nil
}

