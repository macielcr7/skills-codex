package video

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"example.com/service/internal/domain/entity/media/video"
	repovideo "example.com/service/internal/domain/repository/media/video"
)

type fakeVideoRepo struct {
	created []video.Video
}

func (f *fakeVideoRepo) Create(_ context.Context, v video.Video) error {
	f.created = append(f.created, v)
	return nil
}

func (f *fakeVideoRepo) FindByID(_ context.Context, id string) (video.Video, error) {
	for _, v := range f.created {
		if v.ID == id {
			return v, nil
		}
	}
	return video.Video{}, repovideo.ErrVideoNotFound
}

type fixedIDGenerator struct {
	id string
}

func (g fixedIDGenerator) NewID() string {
	return g.id
}

func TestUploadVideo_Valid(t *testing.T) {
	repo := &fakeVideoRepo{}
	idGen := fixedIDGenerator{id: "v1"}

	uc := NewUploadVideo(repo, idGen)
	out, err := uc.Execute(context.Background(), UploadVideoInput{
		Title:    "title",
		FilePath: "/tmp/file.mp4",
	})
	require.NoError(t, err)
	require.Equal(t, "v1", out.ID)
	require.Len(t, repo.created, 1)
	require.Equal(t, "v1", repo.created[0].ID)
}

func TestUploadVideo_Invalid(t *testing.T) {
	repo := &fakeVideoRepo{}
	idGen := fixedIDGenerator{id: "v1"}

	uc := NewUploadVideo(repo, idGen)
	_, err := uc.Execute(context.Background(), UploadVideoInput{
		Title:    "",
		FilePath: "/tmp/file.mp4",
	})
	require.ErrorIs(t, err, video.ErrInvalidTitle)
	require.Empty(t, repo.created)
}
