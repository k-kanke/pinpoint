package gormrepo

import (
    "context"

    "github.com/k-kanke/pinpoint/backend/app/domain/model"
    "gorm.io/gorm"
)

type PhotoRepository struct{ db *gorm.DB }

func NewPhotoRepository(db *gorm.DB) *PhotoRepository { return &PhotoRepository{db: db} }

func (r *PhotoRepository) Create(ctx context.Context, p *model.Photo) error {
    return r.db.WithContext(ctx).Create(p).Error
}

