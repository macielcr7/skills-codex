package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	usecasevideo "example.com/service/internal/application/usecase/media/video"
	handlervideo "example.com/service/internal/infra/api/handler/media/video"
	mediarepo "example.com/service/internal/infra/repository/media/video/memory"
)

type fixedIDGenerator struct{}

func (fixedIDGenerator) NewID() string { return "v1" }

func TestRouter_VideoFlow(t *testing.T) {
	videoRepo := mediarepo.NewVideoRepository()
	uploadVideo := usecasevideo.NewUploadVideo(videoRepo, fixedIDGenerator{})
	getVideo := usecasevideo.NewGetVideo(videoRepo)

	videoHandler := handlervideo.NewVideoHandler(uploadVideo, getVideo)
	router := NewRouter(videoHandler)

	postReq := httptest.NewRequest(http.MethodPost, "/videos", strings.NewReader(`{"title":"t","file_path":"/tmp/a.mp4"}`))
	postReq.Header.Set("Content-Type", "application/json")
	postRec := httptest.NewRecorder()
	router.ServeHTTP(postRec, postReq)

	require.Equal(t, http.StatusCreated, postRec.Code)

	var postResp struct {
		ID string `json:"id"`
	}
	require.NoError(t, json.NewDecoder(postRec.Body).Decode(&postResp))
	require.Equal(t, "v1", postResp.ID)

	getReq := httptest.NewRequest(http.MethodGet, "/videos/v1", nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	require.Equal(t, http.StatusOK, getRec.Code)

	var getResp struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		FilePath string `json:"file_path"`
		Status   string `json:"status"`
	}
	require.NoError(t, json.NewDecoder(getRec.Body).Decode(&getResp))
	require.Equal(t, "v1", getResp.ID)
	require.Equal(t, "t", getResp.Title)
	require.Equal(t, "/tmp/a.mp4", getResp.FilePath)
	require.Equal(t, "pending", getResp.Status)
}

func TestRouter_Health(t *testing.T) {
	videoRepo := mediarepo.NewVideoRepository()
	uploadVideo := usecasevideo.NewUploadVideo(videoRepo, fixedIDGenerator{})
	getVideo := usecasevideo.NewGetVideo(videoRepo)
	videoHandler := handlervideo.NewVideoHandler(uploadVideo, getVideo)
	router := NewRouter(videoHandler)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "ok", strings.TrimSpace(rec.Body.String()))
}
