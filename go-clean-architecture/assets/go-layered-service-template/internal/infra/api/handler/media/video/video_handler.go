package video

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	usecasevideo "example.com/service/internal/application/usecase/media/video"
)

// VideoHandler exposes HTTP handlers for video operations.
type VideoHandler struct {
	uploadVideo usecasevideo.UploadVideo
	getVideo    usecasevideo.GetVideo
}

// NewVideoHandler creates a new VideoHandler.
func NewVideoHandler(uploadVideo usecasevideo.UploadVideo, getVideo usecasevideo.GetVideo) *VideoHandler {
	return &VideoHandler{
		uploadVideo: uploadVideo,
		getVideo:    getVideo,
	}
}

type uploadVideoRequest struct {
	Title    string `json:"title"`
	FilePath string `json:"file_path"`
}

type uploadVideoResponse struct {
	ID string `json:"id"`
}

// UploadVideo handles POST /videos.
func (h *VideoHandler) UploadVideo(w http.ResponseWriter, r *http.Request) {
	var req uploadVideoRequest
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	out, err := h.uploadVideo.Execute(r.Context(), usecasevideo.UploadVideoInput{
		Title:    req.Title,
		FilePath: req.FilePath,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(uploadVideoResponse{ID: out.ID})
}

type videoResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	FilePath string `json:"file_path"`
	Status   string `json:"status"`
}

// GetVideo handles GET /videos/{id}.
func (h *VideoHandler) GetVideo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	out, err := h.getVideo.Execute(r.Context(), usecasevideo.GetVideoInput{ID: id})
	if err != nil {
		writeError(w, err)
		return
	}

	resp := videoResponse{
		ID:       out.Video.ID,
		Title:    out.Video.Title,
		FilePath: out.Video.FilePath,
		Status:   string(out.Video.Status),
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
