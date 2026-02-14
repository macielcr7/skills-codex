package video

import (
	"context"
	"errors"

	"example.com/service/internal/domain/entity/media/video"
)

// ErrVideoNotFound indicates a video does not exist in the repository.
var ErrVideoNotFound = errors.New("video not found")

// Repository defines persistence operations for the Video aggregate.
type Repository interface {
	Create(ctx context.Context, v video.Video) error
	FindByID(ctx context.Context, id string) (video.Video, error)
}
