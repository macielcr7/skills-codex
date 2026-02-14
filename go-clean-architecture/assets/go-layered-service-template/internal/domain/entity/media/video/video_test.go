package video

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew_Valid(t *testing.T) {
	v, err := New("v1", "title", "/tmp/file.mp4")
	require.NoError(t, err)
	require.Equal(t, "v1", v.ID)
	require.Equal(t, StatusPending, v.Status)
}

func TestNew_Invalid(t *testing.T) {
	_, err := New("", "title", "/tmp/file.mp4")
	require.ErrorIs(t, err, ErrInvalidID)

	_, err = New("v1", "", "/tmp/file.mp4")
	require.ErrorIs(t, err, ErrInvalidTitle)

	_, err = New("v1", "title", "")
	require.ErrorIs(t, err, ErrInvalidFilePath)
}

func TestVideo_CanBeProcessed(t *testing.T) {
	v, err := New("v1", "title", "/tmp/file.mp4")
	require.NoError(t, err)
	require.True(t, v.CanBeProcessed())

	v.Status = StatusFailed
	require.True(t, v.CanBeProcessed())

	v.Status = StatusProcessing
	require.False(t, v.CanBeProcessed())
}

