package video

import (
	"errors"
	"strings"
)

// Status defines the processing status of a video.
type Status string

const (
	// StatusPending indicates the video is queued to be processed.
	StatusPending Status = "pending"
	// StatusProcessing indicates the video is being processed.
	StatusProcessing Status = "processing"
	// StatusCompleted indicates the video processing finished successfully.
	StatusCompleted Status = "completed"
	// StatusFailed indicates the video processing failed.
	StatusFailed Status = "failed"
)

var (
	// ErrInvalidID indicates an invalid video ID.
	ErrInvalidID = errors.New("invalid id")
	// ErrInvalidTitle indicates an invalid video title.
	ErrInvalidTitle = errors.New("invalid title")
	// ErrInvalidFilePath indicates an invalid file path.
	ErrInvalidFilePath = errors.New("invalid file path")
)

// Video represents a video to be processed (bounded context: media).
type Video struct {
	ID           string
	Title        string
	FilePath     string
	Status       Status
	ErrorMessage string
}

// New creates a new Video enforcing domain invariants.
func New(id, title, filePath string) (Video, error) {
	if strings.TrimSpace(id) == "" {
		return Video{}, ErrInvalidID
	}
	if strings.TrimSpace(title) == "" {
		return Video{}, ErrInvalidTitle
	}
	if strings.TrimSpace(filePath) == "" {
		return Video{}, ErrInvalidFilePath
	}

	return Video{
		ID:       id,
		Title:    title,
		FilePath: filePath,
		Status:   StatusPending,
	}, nil
}

// CanBeProcessed reports whether the video can be (re)processed.
func (v Video) CanBeProcessed() bool {
	return v.Status == StatusPending || v.Status == StatusFailed
}

