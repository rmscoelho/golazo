package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// LeagueInfo contains league metadata for display purposes.
type LeagueInfo struct {
	ID      int
	Name    string
	Country string
}

// Region constants for organizing leagues
const (
	RegionEurope  = "Europe"
	RegionAmerica = "Americas"
	RegionGlobal  = "Global"
)

// AllSupportedLeagues contains all leagues that Golazo supports organized by region.
var AllSupportedLeagues = map[string][]LeagueInfo{
	RegionEurope: {
		// Top 5 European Leagues
		{ID: 47, Name: "Premier League", Country: "England"},
		{ID: 87, Name: "La Liga", Country: "Spain"},
		{ID: 54, Name: "Bundesliga", Country: "Germany"},
		{ID: 55, Name: "Serie A", Country: "Italy"},
		{ID: 53, Name: "Ligue 1", Country: "France"},
		// Top 5 Women's Leagues
		{ID: 9227, Name: "Women's Super League", Country: "England"},
		{ID: 9907, Name: "Liga F", Country: "Spain"},
		{ID: 9676, Name: "Frauen-Bundesliga", Country: "Germany"},
		{ID: 10178, Name: "Serie A Femminile", Country: "Italy"},
		{ID: 9667, Name: "Première Ligue Féminine", Country: "France"},
		// Other European Leagues
		{ID: 67, Name: "Allsvenskan", Country: "Sweden"},
		{ID: 38, Name: "Austrian Bundesliga", Country: "Austria"},
		{ID: 40, Name: "Belgian First Division", Country: "Belgium"},
		{ID: 48, Name: "EFL Championship", Country: "England"},
		{ID: 108, Name: "EFL League One", Country: "England"},
		{ID: 109, Name: "EFL League Two", Country: "England"},
		{ID: 196, Name: "Ekstraklasa", Country: "Poland"},
		{ID: 57, Name: "Eredivisie", Country: "Netherlands"},
		{ID: 218, Name: "League of Ireland First Division", Country: "Ireland"},
		{ID: 126, Name: "League of Ireland Premier Division", Country: "Ireland"},
		{ID: 61, Name: "Primeira Liga", Country: "Portugal"},
		{ID: 64, Name: "Scottish Premiership", Country: "Scotland"},
		{ID: 135, Name: "Super League 1", Country: "Greece"},
		{ID: 46, Name: "Superligaen", Country: "Denmark"},
		{ID: 71, Name: "Süper Lig", Country: "Turkey"},
		{ID: 69, Name: "Swiss Super League", Country: "Switzerland"},
		{ID: 63, Name: "Russian Premier League", Country: "Russia"},
		{ID: 441, Name: "Ukrainian Premier League", Country: "Ukraine"},
		// European Competitions
		{ID: 42, Name: "UEFA Champions League", Country: "Europe"},
		{ID: 10216, Name: "UEFA Conference League", Country: "Europe"},
		{ID: 73, Name: "UEFA Europa League", Country: "Europe"},
		{ID: 50, Name: "UEFA Euro", Country: "Europe"},
		{ID: 292, Name: "UEFA Women's Euro", Country: "Europe"},
		{ID: 9375, Name: "Women's UEFA Champions League", Country: "Europe"},
		// Domestic Cups
		{ID: 138, Name: "Copa del Rey", Country: "Spain"},
		{ID: 139, Name: "Supercopa de España", Country: "Spain"},
		{ID: 132, Name: "FA Cup", Country: "England"},
		{ID: 209, Name: "DFB Pokal", Country: "Germany"},
		{ID: 10650, Name: "Women's DFB Pokal", Country: "Germany"},
		{ID: 141, Name: "Coppa Italia", Country: "Italy"},
		{ID: 134, Name: "Coupe de France", Country: "France"},
	},
	RegionAmerica: {
		// South America
		{ID: 268, Name: "Brasileirão Série A", Country: "Brazil"},
		{ID: 8814, Name: "Brasileirão Série B", Country: "Brazil"},
		{ID: 44, Name: "Copa America", Country: "South America"},
		{ID: 9490, Name: "Copa Colombia", Country: "Colombia"},
		{ID: 45, Name: "Copa Libertadores", Country: "South America"},
		{ID: 299, Name: "Copa Sudamericana", Country: "South America"},
		{ID: 112, Name: "Liga Profesional", Country: "Argentina"},
		{ID: 274, Name: "Primera A", Country: "Colombia"},
		{ID: 9125, Name: "Primera B", Country: "Colombia"},
		// North America
		{ID: 130, Name: "MLS", Country: "USA"},
		{ID: 9134, Name: "NWSL", Country: "USA"},
		{ID: 230, Name: "Liga MX", Country: "Mexico"},
	},
	RegionGlobal: {
		// Middle East
		{ID: 536, Name: "Saudi Pro League", Country: "Saudi Arabia"},
		// Asia
		{ID: 9478, Name: "Indian Super League", Country: "India"},
		{ID: 223, Name: "J. League", Country: "Japan"},
		{ID: 9080, Name: "K League 1", Country: "South Korea"},
		{ID: 9137, Name: "Chinese League One", Country: "China"},
		{ID: 535, Name: "Qatar Stars League", Country: "Qatar"},
		// Oceania
		{ID: 113, Name: "A-League", Country: "Australia"},
		// Africa
		{ID: 519, Name: "Egyptian Premier League", Country: "Egypt"},
		{ID: 537, Name: "Premier Soccer League", Country: "South Africa"},
		{ID: 530, Name: "Botola Pro", Country: "Morocco"},
		// International Competitions
		{ID: 289, Name: "Africa Cup of Nations", Country: "International"},
		{ID: 77, Name: "FIFA World Cup", Country: "International"},
		{ID: 76, Name: "Women's FIFA World Cup", Country: "International"},
		{ID: 78, Name: "FIFA Club World Cup", Country: "International"},
		{ID: 9806, Name: "UEFA Nations League", Country: "International"},
		{ID: 489, Name: "Club Friendlies", Country: "International"},
		{ID: 114, Name: "International Friendlies", Country: "International"},
	},
}

func fetchLeagueName(leagueID int) (string, error) {
	url := fmt.Sprintf("https://www.fotmob.com/leagues/%d/overview", leagueID)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return extractLeagueName(string(body))
}

func extractLeagueName(html string) (string, error) {
	// Look for the league name in the title tag
	titleRegex := regexp.MustCompile(`<title[^>]*>([^<]+)</title>`)
	matches := titleRegex.FindStringSubmatch(html)
	if len(matches) > 1 {
		title := matches[1]
		// Extract league name from title (remove "matches, tables and news" suffix)
		parts := strings.Split(title, " matches, tables and news")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0]), nil
		}
	}

	return "", fmt.Errorf("league name not found in HTML")
}

func validateLeague(league LeagueInfo) string {
	fetchedName, err := fetchLeagueName(league.ID)
	if err != nil {
		return fmt.Sprintf("ERROR: %v", err)
	}

	if strings.TrimSpace(fetchedName) == strings.TrimSpace(league.Name) {
		return "VALID"
	}

	return fmt.Sprintf("INVALID (got: %s)", fetchedName)
}

func main() {
	fmt.Printf("%-8s %-40s %s\n", "ID", "TRACKED NAME", "RESULT")
	fmt.Println(strings.Repeat("-", 80))

	for region, leagues := range AllSupportedLeagues {
		fmt.Printf("\n%s:\n", region)
		for _, league := range leagues {
			result := validateLeague(league)
			fmt.Printf("%-8d %-40s %s\n", league.ID, league.Name, result)

			// Small delay to be respectful to the server
			time.Sleep(100 * time.Millisecond)
		}
	}
}
