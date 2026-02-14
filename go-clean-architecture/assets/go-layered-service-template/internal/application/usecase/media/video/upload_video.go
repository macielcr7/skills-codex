package video

import (
	"context"

	"example.com/service/internal/application/service/shared"
	domainvideo "example.com/service/internal/domain/entity/media/video"
	repovideo "example.com/service/internal/domain/repository/media/video"
)

// UploadVideo is the use case responsible for registering a new video.
type UploadVideo struct {
	repo  repovideo.Repository
	idGen shared.IDGenerator
}

// NewUploadVideo builds the UploadVideo use case with its dependencies.
func NewUploadVideo(repo repovideo.Repository, idGen shared.IDGenerator) UploadVideo {
	return UploadVideo{
		repo:  repo,
		idGen: idGen,
	}
}

// UploadVideoInput defines the input for the UploadVideo use case.
type UploadVideoInput struct {
	Title    string
	FilePath string
}

// UploadVideoOutput defines the output for the UploadVideo use case.
type UploadVideoOutput struct {
	ID string
}

// Execute runs the UploadVideo use case.
func (uc UploadVideo) Execute(ctx context.Context, in UploadVideoInput) (UploadVideoOutput, error) {
	id := uc.idGen.NewID()

	videoEntity, err := domainvideo.New(id, in.Title, in.FilePath)
	if err != nil {
		return UploadVideoOutput{}, err
	}

	if err := uc.repo.Create(ctx, videoEntity); err != nil {
		return UploadVideoOutput{}, err
	}

	return UploadVideoOutput{ID: id}, nil
}
