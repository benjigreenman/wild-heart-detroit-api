package api

import (
	"encoding/json"
	"net/http"
	"os"

	"wild-heart-detroit-api/internal/provider"
	"wild-heart-detroit-api/internal/services"
)

type Handler struct {
	mediaService *services.MediaService
}

func NewHandler() *Handler {
	// For AWS Lambda or env-based config
	config := map[string]string{
		"SPOTIFY_CLIENT_ID":     os.Getenv("SPOTIFY_CLIENT_ID"),
		"SPOTIFY_CLIENT_SECRET": os.Getenv("SPOTIFY_CLIENT_SECRET"),
		"SPOTIFY_PODCAST_ID":    os.Getenv("SPOTIFY_PODCAST_ID"),
		"YOUTUBE_API_KEY":       os.Getenv("YOUTUBE_API_KEY"),
		"YOUTUBE_PLAYLIST_ID":   os.Getenv("YOUTUBE_PLAYLIST_ID"),
		"APPLE_PODCAST_ID":      os.Getenv("APPLE_PODCAST_ID"),
	}
	return NewHandlerWithConfig(config)
}

func NewHandlerWithConfig(config map[string]string) *Handler {
	mediaService := services.NewMediaService(
		provider.NewProvider("youtube", config).(*provider.YouTubeProvider),
		provider.NewProvider("spotify", config).(*provider.SpotifyProvider),
		provider.NewProvider("apple", config).(*provider.AppleProvider),
	)
	return &Handler{
		mediaService: mediaService,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := h.mediaService.GetAllMediaContent(
		ctx,
		15,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
