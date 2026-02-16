package main

import (
	"log"
	"net/http"

	"wild-heart-detroit-api/internal/api"
	"wild-heart-detroit-api/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	handler := api.NewHandlerWithConfig(cfg)

	http.Handle("/media", handler)

	log.Println("Server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
