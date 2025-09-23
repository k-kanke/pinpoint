package model

import (
    "time"
    "github.com/google/uuid"
)

// User represents application user.
type User struct {
    ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Username         string    `gorm:"size:50;uniqueIndex;not null"`
    Email            string    `gorm:"size:255;uniqueIndex;not null"`
    PasswordHash     string    `gorm:"size:255;not null"`
    SubscriptionPlan string    `gorm:"size:50;not null;default:'free'"`
    CreatedAt        time.Time `gorm:"not null;default:now()"`
}

// Thread represents a geo-located post (pin thread).
// Note: location uses PostGIS geometry(Point,4326). GORM handles via raw SQL migration.
type Thread struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Title     string    `gorm:"size:140;not null"`
    Body      string    `gorm:"type:text;not null"`
    // Location column will be added by migration as geometry(Point,4326)
    ExpiryAt  *time.Time
    CreatedAt time.Time `gorm:"not null;default:now()"`

    Photos   []Photo   `gorm:"constraint:OnDelete:CASCADE"`
    Comments []Comment `gorm:"constraint:OnDelete:CASCADE"`
}

// Photo represents a photo associated with a thread.
type Photo struct {
    ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    ThreadID uuid.UUID `gorm:"type:uuid;index;not null"`
    URL      string    `gorm:"type:text;not null"`
}

// Comment represents a comment on a thread.
type Comment struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    ThreadID  uuid.UUID `gorm:"type:uuid;index;not null"`
    Body      string    `gorm:"type:text;not null"`
    CreatedAt time.Time `gorm:"not null;default:now()"`
}

