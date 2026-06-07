package repository

import (
	"context"
	"errors"
	"go-fake-flix/internal/modules/video/entities"
	"sync"
)

type VideoRepository interface {
	Get(ctx context.Context, id string) (entities.Video, error)
	Add(ctx context.Context, v entities.Video) error
}

var (
	ErrNotFound      = errors.New("video not found")
	ErrAlreadyExists = errors.New("video already exists")
)

type Repository struct {
	mu    sync.RWMutex
	store map[string]entities.Video
}

func New() *Repository {
	return &Repository{store: make(map[string]entities.Video)}
}

func (r *Repository) Get(ctx context.Context, id string) (entities.Video, error) {
	if err := ctx.Err(); err != nil {
		return entities.Video{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.store[id]
	if !ok {
		return entities.Video{}, ErrNotFound
	}
	return v, nil
}

func (r *Repository) Add(ctx context.Context, v entities.Video) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.store[v.FileName]; ok {
		return ErrAlreadyExists
	}
	r.store[v.FileName] = v
	return nil
}
