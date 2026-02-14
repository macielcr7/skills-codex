package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	handlervideo "example.com/service/internal/infra/api/handler/media/video"
)

// NewRouter builds the HTTP router.
func NewRouter(mediaVideoHandler *handlervideo.VideoHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Route("/videos", func(r chi.Router) {
		r.Post("/", mediaVideoHandler.UploadVideo)
		r.Get("/{id}", mediaVideoHandler.GetVideo)
	})

	return r
}
