package gormrepo

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/k-kanke/pinpoint/backend/app/domain/model"
    "gorm.io/gorm"
)

type ThreadRepository struct {
    db *gorm.DB
}

func NewThreadRepository(db *gorm.DB) *ThreadRepository {
    return &ThreadRepository{db: db}
}

func (r *ThreadRepository) Create(ctx context.Context, t *model.Thread, lat, lon float64) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(t).Error; err != nil {
            return err
        }
        // Set geometry location after insert
        q := `UPDATE threads SET location = ST_SetSRID(ST_MakePoint(?, ?), 4326) WHERE id = ?`
        if err := tx.Exec(q, lon, lat, t.ID).Error; err != nil {
            return fmt.Errorf("set location failed: %w", err)
        }
        return nil
    })
}

func (r *ThreadRepository) FindNearby(ctx context.Context, lat, lon float64, radiusMeters int, limit int) ([]model.Thread, error) {
    var rows []model.Thread
    q := `SELECT id, title, body, expiry_at, created_at
          FROM threads
          WHERE location IS NOT NULL
            AND ST_DWithin(location, ST_SetSRID(ST_MakePoint(?, ?), 4326), ?)
          ORDER BY created_at DESC
          LIMIT ?`
    if err := r.db.WithContext(ctx).Raw(q, lon, lat, radiusMeters, limit).Scan(&rows).Error; err != nil {
        return nil, err
    }
    return rows, nil
}

func (r *ThreadRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Thread, error) {
    var t model.Thread
    if err := r.db.WithContext(ctx).
        Preload("Photos").
        Preload("Comments").
        First(&t, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &t, nil
}

