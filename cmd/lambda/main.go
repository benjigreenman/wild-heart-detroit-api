package main

import (
	"context"
	"encoding/json"
	"os"

	"wild-heart-detroit-api/internal/provider"
	"wild-heart-detroit-api/internal/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context) (events.APIGatewayV2HTTPResponse, error) {
	config := map[string]string{
		"SPOTIFY_CLIENT_ID":     os.Getenv("SPOTIFY_CLIENT_ID"),
		"SPOTIFY_CLIENT_SECRET": os.Getenv("SPOTIFY_CLIENT_SECRET"),
		"YOUTUBE_API_KEY":       os.Getenv("YOUTUBE_API_KEY"),
		"APPLE_API_KEY":         os.Getenv("APPLE_API_KEY"),
	}

	mediaService := services.NewMediaService(
		provider.NewProvider("youtube", config).(*provider.YouTubeProvider),
		provider.NewProvider("spotify", config).(*provider.SpotifyProvider),
		provider.NewProvider("apple", config).(*provider.AppleProvider),
	)

	items, err := mediaService.GetAllMediaContent(ctx, 15)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: 500}, err
	}

	body, _ := json.Marshal(items)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
