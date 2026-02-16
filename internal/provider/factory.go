package provider

func NewProvider(kind string, config map[string]string) Provider {
	switch kind {
	case "spotify":
		return &SpotifyProvider{
			ClientID:     config["SPOTIFY_CLIENT_ID"],
			ClientSecret: config["SPOTIFY_CLIENT_SECRET"],
		}
	case "youtube":
		return &YouTubeProvider{
			APIKey: config["YOUTUBE_API_KEY"],
		}
	case "apple":
		return &AppleProvider{
			APIKey: config["APPLE_API_KEY"],
		}
	default:
		return nil
	}
}