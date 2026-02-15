package video

import (
	"errors"
	"strings"

	"github.com/asaskevich/govalidator"
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
	ID           string `valid:"required"`
	Title        string `valid:"required"`
	FilePath     string `valid:"required"`
	Status       Status
	ErrorMessage string
}

// New creates a new Video enforcing domain invariants.
func New(id, title, filePath string) (Video, error) {
	v := Video{
		ID:       id,
		Title:    title,
		FilePath: filePath,
		Status:   StatusPending,
	}

	if err := v.Validate(); err != nil {
		return Video{}, err
	}

	return v, nil
}

// Validate validates the entity fields and domain invariants.
func (v Video) Validate() error {
	if strings.TrimSpace(v.ID) == "" {
		return ErrInvalidID
	}
	if strings.TrimSpace(v.Title) == "" {
		return ErrInvalidTitle
	}
	if strings.TrimSpace(v.FilePath) == "" {
		return ErrInvalidFilePath
	}

	_, err := govalidator.ValidateStruct(v)
	return err
}

// CanBeProcessed reports whether the video can be (re)processed.
func (v Video) CanBeProcessed() bool {
	return v.Status == StatusPending || v.Status == StatusFailed
}
