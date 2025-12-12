package api

import (
	"context"
	"time"
)

// Client defines the interface for a football API client
// This abstraction allows us to swap implementations (FotMob, other APIs, mock, etc.)
type Client interface {
	// GetMatchesByDate retrieves all matches for a specific date
	GetMatchesByDate(ctx context.Context, date time.Time) ([]Match, error)

	// GetMatchDetails retrieves detailed information about a specific match
	GetMatchDetails(ctx context.Context, matchID int) (*MatchDetails, error)

	// GetLeagues retrieves available leagues
	GetLeagues(ctx context.Context) ([]League, error)

	// GetLeagueMatches retrieves matches for a specific league
	GetLeagueMatches(ctx context.Context, leagueID int) ([]Match, error)

	// GetLeagueTable retrieves the league table/standings for a specific league
	GetLeagueTable(ctx context.Context, leagueID int) ([]LeagueTableEntry, error)
}

