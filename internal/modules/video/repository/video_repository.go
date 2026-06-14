package repository

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"go-fake-flix/internal/modules/video/entities"
)

type VideoRepository interface {
	Get(ctx context.Context, id int64) (entities.Video, error)
	Add(ctx context.Context, v entities.Video) error
	NextID() int64
}

var (
	ErrNotFound = errors.New("video not found")
)

type Repository struct {
	mu     sync.RWMutex
	store  map[int64]entities.Video
	nextID atomic.Int64
}

func New() *Repository {
	return &Repository{store: make(map[int64]entities.Video)}
}

func (r *Repository) NextID() int64 {
	return r.nextID.Add(1)
}

func (r *Repository) Get(ctx context.Context, id int64) (entities.Video, error) {
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

	r.store[v.ID] = v
	return nil
}
