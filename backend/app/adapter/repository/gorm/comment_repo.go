package gormrepo

import (
    "context"

    "github.com/k-kanke/pinpoint/backend/app/domain/model"
    "gorm.io/gorm"
)

type CommentRepository struct{ db *gorm.DB }

func NewCommentRepository(db *gorm.DB) *CommentRepository { return &CommentRepository{db: db} }

func (r *CommentRepository) Create(ctx context.Context, c *model.Comment) error {
    return r.db.WithContext(ctx).Create(c).Error
}

