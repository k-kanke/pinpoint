package pin

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/k-kanke/pinpoint/backend/app/domain/model"
)

// Ports (interfaces) that repositories must implement.
type ThreadRepo interface {
    Create(ctx context.Context, t *model.Thread, lat, lon float64) error
    FindNearby(ctx context.Context, lat, lon float64, radiusMeters int, limit int) ([]model.Thread, error)
    GetByID(ctx context.Context, id uuid.UUID) (*model.Thread, error)
}

type CommentRepo interface {
    Create(ctx context.Context, c *model.Comment) error
}

type PhotoRepo interface {
    Create(ctx context.Context, p *model.Photo) error
}

// UseCase interface exposed to handlers.
type UseCase interface {
    GetNearby(ctx context.Context, lat, lon float64, radiusMeters, limit int) ([]model.Thread, error)
    CreatePin(ctx context.Context, cmd CreatePinCommand) (uuid.UUID, error)
    GetByID(ctx context.Context, id uuid.UUID) (*model.Thread, error)
    CreateComment(ctx context.Context, threadID uuid.UUID, body string) (uuid.UUID, error)
    AttachPhoto(ctx context.Context, threadID uuid.UUID, url string) error
}

type useCase struct {
    threads ThreadRepo
    comments CommentRepo
    photos PhotoRepo
}

func New(threads ThreadRepo, comments CommentRepo, photos PhotoRepo) UseCase {
    return &useCase{threads: threads, comments: comments, photos: photos}
}

type CreatePinCommand struct {
    Title string
    Body string
    Lat float64
    Lon float64
    ExpiryAt *time.Time
    PhotoURL string
}

func (u *useCase) GetNearby(ctx context.Context, lat, lon float64, radiusMeters, limit int) ([]model.Thread, error) {
    return u.threads.FindNearby(ctx, lat, lon, radiusMeters, limit)
}

func (u *useCase) CreatePin(ctx context.Context, cmd CreatePinCommand) (uuid.UUID, error) {
    if cmd.Title == "" || cmd.Body == "" {
        return uuid.Nil, errors.New("title and body required")
    }
    t := &model.Thread{
        Title: cmd.Title,
        Body:  cmd.Body,
        ExpiryAt: cmd.ExpiryAt,
    }
    if err := u.threads.Create(ctx, t, cmd.Lat, cmd.Lon); err != nil {
        return uuid.Nil, err
    }
    if cmd.PhotoURL != "" {
        _ = u.photos.Create(ctx, &model.Photo{ThreadID: t.ID, URL: cmd.PhotoURL})
    }
    return t.ID, nil
}

func (u *useCase) GetByID(ctx context.Context, id uuid.UUID) (*model.Thread, error) {
    return u.threads.GetByID(ctx, id)
}

func (u *useCase) CreateComment(ctx context.Context, threadID uuid.UUID, body string) (uuid.UUID, error) {
    if body == "" {
        return uuid.Nil, errors.New("body required")
    }
    c := &model.Comment{ThreadID: threadID, Body: body}
    if err := u.comments.Create(ctx, c); err != nil {
        return uuid.Nil, err
    }
    return c.ID, nil
}

func (u *useCase) AttachPhoto(ctx context.Context, threadID uuid.UUID, url string) error {
    if url == "" {
        return errors.New("url required")
    }
    p := &model.Photo{ThreadID: threadID, URL: url}
    return u.photos.Create(ctx, p)
}
