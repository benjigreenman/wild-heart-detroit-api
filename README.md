# wild-heart-detroit-api
Abstract api to access data from third party apis. Youtube, spotify, apple, and google maps

## Project Structure

- **cmd/local/**: Entry point for running the API server locally.
- **cmd/lambda/**: Entry point for AWS Lambda deployment.
- **internal/api/**: HTTP handler logic for serving API requests.
- **internal/services/**: Business logic for aggregating and combining media content from all providers.
- **internal/provider/**: Provider implementations for Spotify, YouTube, and Apple. Each provider fetches and maps data from its respective API.
- **internal/config/**: Loads configuration from a local config file (for local development only).
- **scripts/**: Utility scripts for testing, running locally, and building for Lambda.

## API Overview

The API aggregates media content (such as podcast episodes) from multiple third-party providers (YouTube, Spotify, Apple). It exposes a single endpoint:

- `GET /media` — Returns a combined list of recent episodes from all providers in a unified structure.

## Module Responsibilities

- **Provider Modules** (`internal/provider/`):
  - Each provider (Spotify, YouTube, Apple) implements a `Provider` interface to fetch and map episodes from its API.
  - Credentials/API keys are injected via config.
- **Service Layer** (`internal/services/mediaservice.go`):
  - Fetches episodes from all providers concurrently.
  - Combines and groups results for the API response.
- **API Handler** (`internal/api/handler.go`):
  - Handles HTTP requests, invokes the service layer, and returns JSON responses.
- **Config Loader** (`internal/config/config.go`):
  - Loads secrets and API keys from `config.json` for local development.

## Scripts (Run from Project Root)

- **test.bat** — Runs all Go tests in the project.
- **run_local.bat** — Builds and runs the local API server.
- **build_lambda.bat** — Builds the Lambda handler for AWS deployment.

### How to Run Scripts

1. Open a terminal and navigate to the project root:
   ```
   cd path\to\wild-heart-detroit-api
   ```
2. Run the desired script:
   - Run tests: `test.bat`
   - Run local server: `run_local.bat`
   - Build Lambda handler: `build_lambda.bat`

## Local Configuration

- The file `config.json` (in the project cmd/local folder) contains API keys and secrets for local development.
- **Do NOT commit `config.json` to GitHub**. It should be added to `.gitignore` to protect sensitive information.
- Example `config.json`:
  ```json
  {
    "SPOTIFY_CLIENT_ID": "spotify_client_id",
    "SPOTIFY_CLIENT_SECRET": "spotify_client_secret",
    "SPOTIFY_PODCAST_ID": "spotify_podcast_id",
    "YOUTUBE_API_KEY": "youtube_api_key",
    "YOUTUBE_PLAYLIST_ID": "youtube_playlist_id",
    "APPLE_PODCAST_ID": "apple_podcast_id"
  }
  ```
- The local server loads this file automatically at startup.
- For AWS Lambda, secrets are loaded from environment variables (set in the Lambda console or via deployment tools).

---

For more details, see the code in each module or contact the project maintainer.

