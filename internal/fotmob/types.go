package fotmob

import (
	"sort"
	"strconv"
	"time"

	"github.com/0xjuanma/golazo/internal/api"
)

// fotmobMatch represents a match in FotMob's API format
// Note: FotMob uses string IDs in JSON, but we convert them to ints
type fotmobMatch struct {
	ID     string `json:"id"` // FotMob returns string IDs
	Round  string `json:"round"`
	Home   team   `json:"home"`
	Away   team   `json:"away"`
	Status status `json:"status"`
	League league `json:"league"`
}

type team struct {
	ID        string `json:"id"` // FotMob returns string IDs
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
}

type league struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
}

type status struct {
	UTCTime   string    `json:"utcTime"`   // Can be null/empty
	Started   *bool     `json:"started"`   // Can be null
	Finished  *bool     `json:"finished"`  // Can be null
	Cancelled *bool     `json:"cancelled"` // Can be null
	LiveTime  *liveTime `json:"liveTime,omitempty"`
	Score     *score    `json:"score,omitempty"`
}

type liveTime struct {
	Short string `json:"short"`
}

type score struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

// toAPIMatch converts a fotmobMatch to api.Match
func (m fotmobMatch) toAPIMatch() api.Match {
	// Convert string IDs to ints
	matchID := parseInt(m.ID)
	homeID := parseInt(m.Home.ID)
	awayID := parseInt(m.Away.ID)

	match := api.Match{
		ID: matchID,
		League: api.League{
			ID:          m.League.ID,
			Name:        m.League.Name,
			Country:     m.League.Country,
			CountryCode: m.League.CountryCode,
		},
		HomeTeam: api.Team{
			ID:        homeID,
			Name:      m.Home.Name,
			ShortName: m.Home.ShortName,
		},
		AwayTeam: api.Team{
			ID:        awayID,
			Name:      m.Away.Name,
			ShortName: m.Away.ShortName,
		},
		Round: m.Round,
	}

	// Parse match time - FotMob uses .000Z format sometimes
	if m.Status.UTCTime != "" {
		var t time.Time
		var err error
		t, err = time.Parse(time.RFC3339, m.Status.UTCTime)
		if err != nil {
			// Try alternative format with milliseconds
			t, err = time.Parse("2006-01-02T15:04:05.000Z", m.Status.UTCTime)
		}
		if err == nil {
			match.MatchTime = &t
		}
	}

	// Determine status - handle null boolean values
	if m.Status.Cancelled != nil && *m.Status.Cancelled {
		match.Status = api.MatchStatusCancelled
	} else if m.Status.Finished != nil && *m.Status.Finished {
		match.Status = api.MatchStatusFinished
	} else if m.Status.Started != nil && *m.Status.Started {
		match.Status = api.MatchStatusLive
		if m.Status.LiveTime != nil {
			match.LiveTime = &m.Status.LiveTime.Short
		}
	} else {
		match.Status = api.MatchStatusNotStarted
	}

	// Set scores if available
	if m.Status.Score != nil {
		match.HomeScore = &m.Status.Score.Home
		match.AwayScore = &m.Status.Score.Away
	}

	return match
}

// fotmobMatchDetails represents detailed match information from FotMob
type fotmobMatchDetails struct {
	ID     string  `json:"id"` // FotMob returns string IDs
	Round  string  `json:"round"`
	Home   team    `json:"home"`
	Away   team    `json:"away"`
	Status status  `json:"status"`
	League league  `json:"league"`
	Events []event `json:"events"`
}

type event struct {
	ID        int    `json:"id"`
	Minute    int    `json:"minute"`
	Type      string `json:"type"`
	TeamID    string `json:"teamId"` // FotMob returns teamId as string
	Player    string `json:"player,omitempty"`
	Assist    string `json:"assist,omitempty"`
	EventType string `json:"eventType,omitempty"`
}

// toAPIMatchDetails converts fotmobMatchDetails to api.MatchDetails
func (m fotmobMatchDetails) toAPIMatchDetails() *api.MatchDetails {
	baseMatch := fotmobMatch{
		ID:     m.ID,
		Round:  m.Round,
		Home:   m.Home,
		Away:   m.Away,
		Status: m.Status,
		League: m.League,
	}.toAPIMatch()

	details := &api.MatchDetails{
		Match:  baseMatch,
		Events: make([]api.MatchEvent, 0, len(m.Events)),
	}

	// Convert events and sort by minute for chronological order
	events := make([]api.MatchEvent, 0, len(m.Events))
	for _, e := range m.Events {
		event := api.MatchEvent{
			ID:        e.ID,
			Minute:    e.Minute,
			Type:      e.Type,
			Timestamp: time.Now(), // FotMob doesn't always provide timestamp
		}

		if e.Player != "" {
			event.Player = &e.Player
		}
		if e.Assist != "" {
			event.Assist = &e.Assist
		}
		if e.EventType != "" {
			event.EventType = &e.EventType
		}

		// Set team based on teamId - match with home or away team
		// Convert team IDs to int for comparison
		homeIDInt := parseInt(m.Home.ID)
		awayIDInt := parseInt(m.Away.ID)
		eventTeamIDInt := parseInt(e.TeamID)

		switch eventTeamIDInt {
		case homeIDInt:
			event.Team = api.Team{
				ID:        homeIDInt,
				Name:      m.Home.Name,
				ShortName: m.Home.ShortName,
			}
		case awayIDInt:
			event.Team = api.Team{
				ID:        awayIDInt,
				Name:      m.Away.Name,
				ShortName: m.Away.ShortName,
			}
		default:
			// Fallback if team ID doesn't match
			event.Team = api.Team{
				ID: eventTeamIDInt,
			}
		}

		events = append(events, event)
	}

	// Sort events by minute (chronological order)
	sort.Slice(events, func(i, j int) bool {
		return events[i].Minute < events[j].Minute
	})

	details.Events = events
	return details
}

// fotmobTableRow represents a single row in the league table from FotMob
type fotmobTableRow struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	ShortName      string `json:"shortName"`
	Rank           int    `json:"rank"`
	Played         int    `json:"played"`
	Wins           int    `json:"wins"`
	Draws          int    `json:"draws"`
	Losses         int    `json:"losses"`
	GoalsFor       int    `json:"goalsFor"`
	GoalsAgainst   int    `json:"goalsAgainst"`
	GoalDifference int    `json:"goalDifference"`
	Points         int    `json:"points"`
}

// toAPITableEntry converts fotmobTableRow to api.LeagueTableEntry
func (r fotmobTableRow) toAPITableEntry() api.LeagueTableEntry {
	return api.LeagueTableEntry{
		Position: r.Rank,
		Team: api.Team{
			ID:        r.ID,
			Name:      r.Name,
			ShortName: r.ShortName,
		},
		Played:         r.Played,
		Won:            r.Wins,
		Drawn:          r.Draws,
		Lost:           r.Losses,
		GoalsFor:       r.GoalsFor,
		GoalsAgainst:   r.GoalsAgainst,
		GoalDifference: r.GoalDifference,
		Points:         r.Points,
	}
}

// Helper function to parse time from various formats
func parseTime(timeStr string) *time.Time {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return &t
		}
	}
	return nil
}

// Helper function to parse int from string
// Returns 0 if parsing fails (for required fields)
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}
