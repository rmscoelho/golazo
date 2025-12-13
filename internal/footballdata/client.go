package footballdata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/0xjuanma/golazo/internal/api"
)

const (
	baseURL = "https://api.football-data.org/v4"
)

// Client implements the api.Client interface for Football-Data.org API
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// NewClient creates a new Football-Data.org API client.
// apiKey is required for authentication (get one at https://www.football-data.org/)
func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}

// FinishedMatchesByDateRange retrieves finished matches within a date range.
// This is used for the stats view to show completed matches.
func (c *Client) FinishedMatchesByDateRange(ctx context.Context, dateFrom, dateTo time.Time) ([]api.Match, error) {
	dateFromStr := dateFrom.Format("2006-01-02")
	dateToStr := dateTo.Format("2006-01-02")

	url := fmt.Sprintf("%s/matches?dateFrom=%s&dateTo=%s&status=FINISHED", c.baseURL, dateFromStr, dateToStr)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Auth-Token", c.apiKey)
	req.Header.Set("User-Agent", "golazo/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch matches: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response footballdataMatchesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	matches := make([]api.Match, 0, len(response.Matches))
	for _, m := range response.Matches {
		matches = append(matches, m.toAPIMatch())
	}

	return matches, nil
}

// RecentFinishedMatches retrieves finished matches from the last N days.
// Defaults to last 7 days if days is 0.
func (c *Client) RecentFinishedMatches(ctx context.Context, days int) ([]api.Match, error) {
	if days == 0 {
		days = 7
	}

	dateTo := time.Now()
	dateFrom := dateTo.AddDate(0, 0, -days)

	return c.FinishedMatchesByDateRange(ctx, dateFrom, dateTo)
}

// MatchesByDate retrieves all matches for a specific date.
// Implements api.Client interface.
func (c *Client) MatchesByDate(ctx context.Context, date time.Time) ([]api.Match, error) {
	dateStr := date.Format("2006-01-02")
	url := fmt.Sprintf("%s/matches?dateFrom=%s&dateTo=%s", c.baseURL, dateStr, dateStr)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Auth-Token", c.apiKey)
	req.Header.Set("User-Agent", "golazo/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch matches: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response footballdataMatchesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	matches := make([]api.Match, 0, len(response.Matches))
	for _, m := range response.Matches {
		matches = append(matches, m.toAPIMatch())
	}

	return matches, nil
}

// MatchDetails retrieves detailed information about a specific match.
// Implements api.Client interface.
func (c *Client) MatchDetails(ctx context.Context, matchID int) (*api.MatchDetails, error) {
	url := fmt.Sprintf("%s/matches/%d", c.baseURL, matchID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Auth-Token", c.apiKey)
	req.Header.Set("User-Agent", "golazo/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch match details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var match footballdataMatchDetails
	if err := json.NewDecoder(resp.Body).Decode(&match); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	baseMatch := match.footballdataMatch.toAPIMatch()
	details := &api.MatchDetails{
		Match: baseMatch,
		Events: []api.MatchEvent{}, // Football-Data.org doesn't provide events in basic match details
	}

	return details, nil
}

// Leagues retrieves available leagues.
// Implements api.Client interface.
func (c *Client) Leagues(ctx context.Context) ([]api.League, error) {
	// Football-Data.org doesn't have a simple leagues endpoint
	// Would need to query competitions endpoint
	return []api.League{}, nil
}

// LeagueMatches retrieves matches for a specific league.
// Implements api.Client interface.
func (c *Client) LeagueMatches(ctx context.Context, leagueID int) ([]api.Match, error) {
	url := fmt.Sprintf("%s/competitions/%d/matches", c.baseURL, leagueID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Auth-Token", c.apiKey)
	req.Header.Set("User-Agent", "golazo/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch league matches: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Matches []footballdataMatch `json:"matches"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	matches := make([]api.Match, 0, len(response.Matches))
	for _, m := range response.Matches {
		matches = append(matches, m.toAPIMatch())
	}

	return matches, nil
}

// LeagueTable retrieves the league table/standings for a specific league.
// Implements api.Client interface.
func (c *Client) LeagueTable(ctx context.Context, leagueID int) ([]api.LeagueTableEntry, error) {
	url := fmt.Sprintf("%s/competitions/%d/standings", c.baseURL, leagueID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-Auth-Token", c.apiKey)
	req.Header.Set("User-Agent", "golazo/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch league table: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Standings []struct {
			Table []struct {
				Position int `json:"position"`
				Team     struct {
					ID        int    `json:"id"`
					Name      string `json:"name"`
					ShortName string `json:"shortName"`
					Crest     string `json:"crest,omitempty"`
				} `json:"team"`
				PlayedGames int `json:"playedGames"`
				Won         int `json:"won"`
				Draw        int `json:"draw"`
				Lost        int `json:"lost"`
				GoalsFor    int `json:"goalsFor"`
				GoalsAgainst int `json:"goalsAgainst"`
				GoalDifference int `json:"goalDifference"`
				Points      int `json:"points"`
			} `json:"table"`
		} `json:"standings"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Standings) == 0 {
		return []api.LeagueTableEntry{}, nil
	}

	entries := make([]api.LeagueTableEntry, 0, len(response.Standings[0].Table))
	for _, row := range response.Standings[0].Table {
		entries = append(entries, api.LeagueTableEntry{
			Position:       row.Position,
			Team: api.Team{
				ID:        row.Team.ID,
				Name:      row.Team.Name,
				ShortName: row.Team.ShortName,
				Logo:      row.Team.Crest,
			},
			Played:         row.PlayedGames,
			Won:            row.Won,
			Drawn:          row.Draw,
			Lost:           row.Lost,
			GoalsFor:       row.GoalsFor,
			GoalsAgainst:   row.GoalsAgainst,
			GoalDifference: row.GoalDifference,
			Points:         row.Points,
		})
	}

	return entries, nil
}

