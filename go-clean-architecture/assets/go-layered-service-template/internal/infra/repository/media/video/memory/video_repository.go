package memory

import (
	"context"
	"sync"

	"example.com/service/internal/domain/entity/media/video"
	repovideo "example.com/service/internal/domain/repository/media/video"
)

// VideoRepository is an in-memory implementation of domain repositories.
type VideoRepository struct {
	mu     sync.RWMutex
	videos map[string]video.Video
}

// NewVideoRepository creates an in-memory repository implementation.
func NewVideoRepository() *VideoRepository {
	return &VideoRepository{
		videos: make(map[string]video.Video),
	}
}

// Create stores a new video.
func (r *VideoRepository) Create(_ context.Context, v video.Video) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.videos[v.ID] = v
	return nil
}

// FindByID retrieves a video by ID.
func (r *VideoRepository) FindByID(_ context.Context, id string) (video.Video, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	v, ok := r.videos[id]
	if !ok {
		return video.Video{}, repovideo.ErrVideoNotFound
	}

	return v, nil
}
