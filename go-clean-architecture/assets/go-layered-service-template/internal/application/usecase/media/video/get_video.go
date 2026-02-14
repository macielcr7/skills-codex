package video

import (
	"context"

	"example.com/service/internal/domain/entity/media/video"
	repovideo "example.com/service/internal/domain/repository/media/video"
)

// GetVideo is the use case responsible for retrieving a video by ID.
type GetVideo struct {
	repo repovideo.Repository
}

// NewGetVideo builds the GetVideo use case with its dependencies.
func NewGetVideo(repo repovideo.Repository) GetVideo {
	return GetVideo{repo: repo}
}

// GetVideoInput defines the input for the GetVideo use case.
type GetVideoInput struct {
	ID string
}

// GetVideoOutput defines the output for the GetVideo use case.
type GetVideoOutput struct {
	Video video.Video
}

// Execute runs the GetVideo use case.
func (uc GetVideo) Execute(ctx context.Context, in GetVideoInput) (GetVideoOutput, error) {
	video, err := uc.repo.FindByID(ctx, in.ID)
	if err != nil {
		return GetVideoOutput{}, err
	}

	return GetVideoOutput{
		Video: video,
	}, nil
}
