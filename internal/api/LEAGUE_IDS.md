# League ID Reference

This document tracks the league IDs used by the FotMob API in this application.

## Supported Leagues

The application currently supports **24 leagues/competitions**:

### Top 5 European Leagues
| League | Country | FotMob ID |
|--------|---------|-----------|
| Premier League | England | 47 |
| La Liga | Spain | 87 |
| Bundesliga | Germany | 54 |
| Serie A | Italy | 55 |
| Ligue 1 | France | 53 |

### Second Tier European Leagues
| League | Country | FotMob ID |
|--------|---------|-----------|
| Eredivisie | Netherlands | 57 |
| Primeira Liga | Portugal | 61 |
| Belgian Pro League | Belgium | 114 |
| Scottish Premiership | Scotland | 64 |
| Süper Lig | Turkey | 71 |
| Swiss Super League | Switzerland | 66 |
| Austrian Bundesliga | Austria | 109 |
| Ekstraklasa | Poland | 52 |

### European Competitions
| Competition | Type | FotMob ID |
|-------------|------|-----------|
| UEFA Champions League | Club | 42 |
| UEFA Europa League | Club | 73 |
| UEFA Euro | International | 50 |

### Domestic Cups
| Cup | Country | FotMob ID |
|-----|---------|-----------|
| Copa del Rey | Spain | 138 |

### South America
| League/Competition | Country | FotMob ID |
|--------------------|---------|-----------|
| Brasileirão Série A | Brazil | 268 |
| Liga Profesional | Argentina | 112 |
| Copa Libertadores | International | 14 |
| Copa America | International | 44 |

### North America
| League | Country | FotMob ID |
|--------|---------|-----------|
| MLS | USA | 130 |
| Liga MX | Mexico | 230 |

### International
| Competition | FotMob ID |
|-------------|-----------|
| FIFA World Cup | 77 |

## FotMob API League IDs

**Location:** `internal/data/settings.go` (AllSupportedLeagues)

```go
AllSupportedLeagues = []LeagueInfo{
    // Top 5 European Leagues
    {ID: 47, Name: "Premier League", Country: "England"},
    {ID: 87, Name: "La Liga", Country: "Spain"},
    {ID: 54, Name: "Bundesliga", Country: "Germany"},
    {ID: 55, Name: "Serie A", Country: "Italy"},
    {ID: 53, Name: "Ligue 1", Country: "France"},
    // Second Tier European Leagues
    {ID: 57, Name: "Eredivisie", Country: "Netherlands"},
    {ID: 61, Name: "Primeira Liga", Country: "Portugal"},
    {ID: 114, Name: "Belgian Pro League", Country: "Belgium"},
    {ID: 64, Name: "Scottish Premiership", Country: "Scotland"},
    {ID: 71, Name: "Süper Lig", Country: "Turkey"},
    {ID: 66, Name: "Swiss Super League", Country: "Switzerland"},
    {ID: 109, Name: "Austrian Bundesliga", Country: "Austria"},
    {ID: 52, Name: "Ekstraklasa", Country: "Poland"},
    // European Competitions
    {ID: 42, Name: "UEFA Champions League", Country: "Europe"},
    {ID: 73, Name: "UEFA Europa League", Country: "Europe"},
    {ID: 50, Name: "UEFA Euro", Country: "Europe"},
    // Domestic Cups
    {ID: 138, Name: "Copa del Rey", Country: "Spain"},
    // South America
    {ID: 268, Name: "Brasileirão Série A", Country: "Brazil"},
    {ID: 112, Name: "Liga Profesional", Country: "Argentina"},
    {ID: 14, Name: "Copa Libertadores", Country: "South America"},
    {ID: 44, Name: "Copa America", Country: "South America"},
    // North America
    {ID: 130, Name: "MLS", Country: "USA"},
    {ID: 230, Name: "Liga MX", Country: "Mexico"},
    // International
    {ID: 77, Name: "FIFA World Cup", Country: "International"},
}
```

**API Endpoint:** `https://www.fotmob.com/api/leagues?id={leagueID}&tab={tab}`

Where `tab` can be:
- `fixtures` - Upcoming matches
- `results` - Finished matches

## Notes

- **FotMob** is used for both the **Live Matches** and **Stats** views
- When adding new leagues, update `internal/data/settings.go` and this document
- Tournament data (World Cup, Euro, Copa America) only available during competition periods
- Users can select specific leagues in Settings to reduce API calls and improve performance

