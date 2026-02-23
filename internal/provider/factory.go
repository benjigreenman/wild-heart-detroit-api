package provider

func NewProvider(kind string, config map[string]string) Provider {
	switch kind {
	case "spotify":
		return &SpotifyProvider{
			ClientID:     config["SPOTIFY_CLIENT_ID"],
			ClientSecret: config["SPOTIFY_CLIENT_SECRET"],
			PodcastID:    config["SPOTIFY_PODCAST_ID"],
		}
	case "youtube":
		return &YouTubeProvider{
			APIKey:    config["YOUTUBE_API_KEY"],
			PodcastId: config["YOUTUBE_PLAYLIST_ID"],
		}
	case "apple":
		return &AppleProvider{
			PodcastID: config["APPLE_PODCAST_ID"],
		}
	default:
		return nil
	}
}
