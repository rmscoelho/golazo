package data

import (
	"fmt"
	"os"
)

// GetFootballDataAPIKey retrieves the API-Football.com (RapidAPI) API key from environment variable.
// The API key must be set via the FOOTBALL_DATA_API_KEY environment variable.
func GetFootballDataAPIKey() (string, error) {
	apiKey := os.Getenv("FOOTBALL_DATA_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("FOOTBALL_DATA_API_KEY environment variable not set")
	}

	return apiKey, nil
}
