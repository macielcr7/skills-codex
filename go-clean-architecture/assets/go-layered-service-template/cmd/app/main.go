package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"example.com/service/internal/infra/api"
	"example.com/service/internal/infra/config"
	idshared "example.com/service/internal/infra/id/shared"
	usecasevideo "example.com/service/internal/application/usecase/media/video"
	handlervideo "example.com/service/internal/infra/api/handler/media/video"
	mediarepo "example.com/service/internal/infra/repository/media/video/memory"
)

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg := config.Load()

	videoRepo := mediarepo.NewVideoRepository()
	idGen := idshared.UUIDGenerator{}

	uploadVideo := usecasevideo.NewUploadVideo(videoRepo, idGen)
	getVideo := usecasevideo.NewGetVideo(videoRepo)

	videoHandler := handlervideo.NewVideoHandler(uploadVideo, getVideo)
	router := api.NewRouter(videoHandler)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("http server listening", "port", cfg.HTTPPort)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server error", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("shutting down")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutdown error", "error", err)
		os.Exit(1)
	}
}
