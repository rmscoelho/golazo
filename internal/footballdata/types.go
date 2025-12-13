package footballdata

import (
	"time"

	"github.com/0xjuanma/golazo/internal/api"
)

// footballdataMatch represents a match in Football-Data.org API format
type footballdataMatch struct {
	ID       int    `json:"id"`
	UTCDate  string `json:"utcDate"`
	Status   string `json:"status"`
	Matchday int    `json:"matchday"`
	Stage    string `json:"stage"`
	Venue    string `json:"venue,omitempty"`

	Competition struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"competition"`

	HomeTeam struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		ShortName string `json:"shortName"`
		TLA       string `json:"tla"`
		Crest     string `json:"crest,omitempty"`
	} `json:"homeTeam"`

	AwayTeam struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		ShortName string `json:"shortName"`
		TLA       string `json:"tla"`
		Crest     string `json:"crest,omitempty"`
	} `json:"awayTeam"`

	Score struct {
		Winner   string `json:"winner,omitempty"`
		Duration string `json:"duration,omitempty"`
		FullTime struct {
			HomeTeam *int `json:"homeTeam,omitempty"`
			AwayTeam *int `json:"awayTeam,omitempty"`
		} `json:"fullTime"`
		HalfTime struct {
			HomeTeam *int `json:"homeTeam,omitempty"`
			AwayTeam *int `json:"awayTeam,omitempty"`
		} `json:"halfTime"`
	} `json:"score"`
}

// footballdataMatchesResponse represents the response from /v4/matches endpoint
type footballdataMatchesResponse struct {
	Matches []footballdataMatch `json:"matches"`
}

// footballdataMatchDetails represents detailed match information
type footballdataMatchDetails struct {
	footballdataMatch
	Attendance int `json:"attendance,omitempty"`
	Minute     int `json:"minute,omitempty"`
}

// toAPIMatch converts a footballdataMatch to api.Match
func (m footballdataMatch) toAPIMatch() api.Match {
	match := api.Match{
		ID: m.ID,
		League: api.League{
			ID:   m.Competition.ID,
			Name: m.Competition.Name,
		},
		HomeTeam: api.Team{
			ID:        m.HomeTeam.ID,
			Name:      m.HomeTeam.Name,
			ShortName: m.HomeTeam.ShortName,
			Logo:      m.HomeTeam.Crest,
		},
		AwayTeam: api.Team{
			ID:        m.AwayTeam.ID,
			Name:      m.AwayTeam.Name,
			ShortName: m.AwayTeam.ShortName,
			Logo:      m.AwayTeam.Crest,
		},
		Round: m.Stage,
	}

	// Parse match time
	if m.UTCDate != "" {
		if t, err := time.Parse(time.RFC3339, m.UTCDate); err == nil {
			match.MatchTime = &t
		}
	}

	// Map status
	switch m.Status {
	case "FINISHED":
		match.Status = api.MatchStatusFinished
	case "LIVE", "IN_PLAY":
		match.Status = api.MatchStatusLive
	case "SCHEDULED", "TIMED":
		match.Status = api.MatchStatusNotStarted
	case "POSTPONED":
		match.Status = api.MatchStatusPostponed
	case "CANCELLED":
		match.Status = api.MatchStatusCancelled
	default:
		match.Status = api.MatchStatusNotStarted
	}

	// Set scores if available
	if m.Score.FullTime.HomeTeam != nil {
		homeScore := *m.Score.FullTime.HomeTeam
		match.HomeScore = &homeScore
	}
	if m.Score.FullTime.AwayTeam != nil {
		awayScore := *m.Score.FullTime.AwayTeam
		match.AwayScore = &awayScore
	}

	return match
}

