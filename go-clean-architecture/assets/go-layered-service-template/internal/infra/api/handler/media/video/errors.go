package video

import (
	"errors"
	"net/http"

	"example.com/service/internal/domain/entity/media/video"
	repovideo "example.com/service/internal/domain/repository/media/video"
)

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, video.ErrInvalidID),
		errors.Is(err, video.ErrInvalidTitle),
		errors.Is(err, video.ErrInvalidFilePath):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, repovideo.ErrVideoNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
