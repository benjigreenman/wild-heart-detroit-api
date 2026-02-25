package main

import (
	"context"
	"encoding/json"
	"log"

	"wild-heart-detroit-api/internal/config"
	"wild-heart-detroit-api/internal/provider"
	"wild-heart-detroit-api/internal/services"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// WHEN SETTING UP IN AWS
// (1) Create a new Lambda function
// (2) Set the handler to "bootstrap"
// (3) In the api gateway, add the desired route/path (api/media) and set the integration to the Lambda function created in step 1

func handler(ctx context.Context) (events.APIGatewayV2HTTPResponse, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	mediaService := services.NewMediaService(
		provider.NewProvider("youtube", cfg).(*provider.YouTubeProvider),
		provider.NewProvider("spotify", cfg).(*provider.SpotifyProvider),
		provider.NewProvider("apple", cfg).(*provider.AppleProvider),
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
