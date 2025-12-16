package ui

import (
	"fmt"
	"strings"

	"github.com/0xjuanma/golazo/internal/api"
	"github.com/0xjuanma/golazo/internal/constants"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// RenderLiveMatchesListPanel renders the left panel using bubbletea list component.
// Note: listModel is passed by value, so SetSize must be called before this function.
func RenderLiveMatchesListPanel(width, height int, listModel list.Model) string {
	// Wrap list in panel
	title := panelTitleStyle.Width(width - 6).Render(constants.PanelLiveMatches)
	listView := listModel.View()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		listView,
	)

	panel := panelStyle.
		Width(width).
		Height(height).
		Render(content)

	return panel
}

// RenderStatsListPanel renders the left panel for stats view using bubbletea list component.
// Note: listModel is passed by value, so SetSize must be called before this function.
// Minimal design matching live view - uses list headers instead of hardcoded titles.
// List titles are only shown when there are items. Empty lists show gray messages instead.
// For 1-day view, shows both finished and upcoming lists stacked vertically.
func RenderStatsListPanel(width, height int, finishedList list.Model, upcomingList list.Model, dateRange int) string {
	// Render date range selector
	dateSelector := renderDateRangeSelector(width-6, dateRange)

	emptyStyle := lipgloss.NewStyle().
		Foreground(dimColor).
		Padding(2, 2).
		Align(lipgloss.Center).
		Width(width - 6)

	var finishedListView string
	finishedItems := finishedList.Items()
	if len(finishedItems) == 0 {
		// No items - show empty message, no list title
		finishedListView = emptyStyle.Render(constants.EmptyNoFinishedMatches + "\n\nTry selecting a different date range (h/l keys)")
	} else {
		// Has items - show list (which includes its title)
		finishedListView = finishedList.View()
	}

	// For 1-day view, show both lists stacked vertically
	if dateRange == 1 {
		var upcomingListView string
		upcomingItems := upcomingList.Items()
		if len(upcomingItems) == 0 {
			// No upcoming matches - show empty message, no list title
			upcomingListView = emptyStyle.Render("No upcoming matches scheduled for today")
		} else {
			// Has items - show list (which includes its title)
			upcomingListView = upcomingList.View()
		}

		// Combine both lists with date selector
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			dateSelector,
			"",
			finishedListView,
			"",
			upcomingListView,
		)
		panel := panelStyle.
			Width(width).
			Height(height).
			Render(content)
		return panel
	}

	// For 3-day view, only show finished matches
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		dateSelector,
		"",
		finishedListView,
	)

	panel := panelStyle.
		Width(width).
		Height(height).
		Render(content)

	return panel
}

// renderDateRangeSelector renders a horizontal date range selector (Today, 3d).
func renderDateRangeSelector(width int, selected int) string {
	options := []struct {
		days  int
		label string
	}{
		{1, "Today"},
		{3, "3d"},
	}

	items := make([]string, 0, len(options))
	for _, opt := range options {
		if opt.days == selected {
			// Selected option - use highlight color
			item := matchListItemSelectedStyle.Render(opt.label)
			items = append(items, item)
		} else {
			// Unselected option - use normal color
			item := matchListItemStyle.Render(opt.label)
			items = append(items, item)
		}
	}

	// Join items with separator
	separator := "  "
	selector := strings.Join(items, separator)

	// Center the selector
	selectorStyle := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Padding(0, 1)

	return selectorStyle.Render(selector)
}

// RenderMultiPanelViewWithList renders the live matches view with list component.
func RenderMultiPanelViewWithList(width, height int, listModel list.Model, details *api.MatchDetails, liveUpdates []string, sp spinner.Model, loading bool, randomSpinner *RandomCharSpinner, viewLoading bool) string {
	// Handle edge case: if width/height not set, use defaults
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	// Reserve 3 lines at top for spinner (always reserve to prevent layout shift)
	spinnerHeight := 3
	availableHeight := height - spinnerHeight
	if availableHeight < 10 {
		availableHeight = 10 // Minimum height for panels
	}

	// Render spinner centered in reserved space
	var spinnerArea string
	if viewLoading && randomSpinner != nil {
		spinnerView := randomSpinner.View()
		if spinnerView != "" {
			// Center the spinner horizontally using style with width and alignment
			spinnerStyle := lipgloss.NewStyle().
				Width(width).
				Height(spinnerHeight).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center)
			spinnerArea = spinnerStyle.Render(spinnerView)
		} else {
			// Fallback if spinner view is empty
			spinnerStyle := lipgloss.NewStyle().
				Width(width).
				Height(spinnerHeight).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center)
			spinnerArea = spinnerStyle.Render("Loading...")
		}
	} else {
		// Reserve space with empty lines - ensure it takes up exactly spinnerHeight lines
		spinnerArea = strings.Repeat("\n", spinnerHeight)
	}

	// Calculate panel dimensions
	leftWidth := width * 35 / 100
	if leftWidth < 25 {
		leftWidth = 25
	}
	rightWidth := width - leftWidth - 1
	if rightWidth < 35 {
		rightWidth = 35
		leftWidth = width - rightWidth - 1
	}

	// Use panelHeight similar to stats view to ensure proper spacing
	panelHeight := availableHeight - 2

	// Render left panel (matches list) - shifted down
	leftPanel := RenderLiveMatchesListPanel(leftWidth, panelHeight, listModel)

	// Render right panel (match details with live updates) - shifted down
	rightPanel := renderMatchDetailsPanel(rightWidth, panelHeight, details, liveUpdates, sp, loading)

	// Create separator
	separatorStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Height(panelHeight).
		Padding(0, 1)
	separator := separatorStyle.Render("│")

	// Combine panels
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		separator,
		rightPanel,
	)

	// Combine spinner area and panels - this shifts panels down
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		spinnerArea,
		panels,
	)

	return content
}

// RenderStatsViewWithList renders the stats view with list component.
// Rebuilt to match live view structure exactly: spinner at top, left panel (matches), right panel (details).
func RenderStatsViewWithList(width, height int, finishedList list.Model, upcomingList list.Model, details *api.MatchDetails, randomSpinner *RandomCharSpinner, viewLoading bool, dateRange int) string {
	// Handle edge case: if width/height not set, use defaults
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	// Reserve 3 lines at top for spinner (always reserve to prevent layout shift)
	// Match live view exactly
	spinnerHeight := 3
	availableHeight := height - spinnerHeight
	if availableHeight < 10 {
		availableHeight = 10 // Minimum height for panels
	}

	// Render spinner centered in reserved space - match live view exactly
	var spinnerArea string
	if viewLoading && randomSpinner != nil {
		spinnerView := randomSpinner.View()
		if spinnerView != "" {
			// Center the spinner horizontally using style with width and alignment
			spinnerStyle := lipgloss.NewStyle().
				Width(width).
				Height(spinnerHeight).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center)
			spinnerArea = spinnerStyle.Render(spinnerView)
		} else {
			// Fallback if spinner view is empty
			spinnerStyle := lipgloss.NewStyle().
				Width(width).
				Height(spinnerHeight).
				Align(lipgloss.Center).
				AlignVertical(lipgloss.Center)
			spinnerArea = spinnerStyle.Render("Loading...")
		}
	} else {
		// Reserve space with empty lines - ensure it takes up exactly spinnerHeight lines
		spinnerArea = strings.Repeat("\n", spinnerHeight)
	}

	// Calculate panel dimensions - match live view exactly (35% left, 65% right)
	leftWidth := width * 35 / 100
	if leftWidth < 25 {
		leftWidth = 25
	}
	rightWidth := width - leftWidth - 1
	if rightWidth < 35 {
		rightWidth = 35
		leftWidth = width - rightWidth - 1
	}

	// Use panelHeight similar to live view to ensure proper spacing
	panelHeight := availableHeight - 2

	// Render left panel (finished matches list) - match live view structure
	// For 1-day view, combine finished and upcoming lists vertically
	leftPanel := RenderStatsListPanel(leftWidth, panelHeight, finishedList, upcomingList, dateRange)

	// Render right panel (match details) - use dedicated stats panel renderer
	rightPanel := renderStatsMatchDetailsPanel(rightWidth, panelHeight, details)

	// Create separator - match live view exactly
	separatorStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Height(panelHeight).
		Padding(0, 1)
	separator := separatorStyle.Render("│")

	// Combine panels - match live view exactly
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPanel,
		separator,
		rightPanel,
	)

	// Combine spinner area and panels - this shifts panels down
	// Match live view exactly - use lipgloss.Left
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		spinnerArea,
		panels,
	)

	return content
}

// renderStatsMatchDetailsPanel renders the right panel for stats view with match details.
// Clean, minimalistic design with bordered sections inspired by modern TUI dashboards.
// Uses the app's cyan/red brand theme throughout.
func renderStatsMatchDetailsPanel(width, height int, details *api.MatchDetails) string {
	if details == nil {
		emptyMessage := lipgloss.NewStyle().
			Foreground(dimColor).
			Align(lipgloss.Center).
			Width(width - 4).
			PaddingTop(height / 3).
			Render("Select a match to view details")

		return panelStyle.
			Width(width).
			Height(height).
			Render(emptyMessage)
	}

	// Brand colors: cyan (51) for accents/headers, red (196) for highlights/scores
	contentWidth := width - 4

	// Section box style - cyan borders
	sectionStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor). // cyan borders
		Padding(0, 1).
		Width(contentWidth - 2)

	// Section header style - cyan
	headerStyle := lipgloss.NewStyle().
		Foreground(accentColor). // cyan headers
		Bold(true)

	// Label style for key-value pairs
	labelStyle := lipgloss.NewStyle().
		Foreground(dimColor).
		Width(14)

	// Value style
	valueStyle := lipgloss.NewStyle().
		Foreground(textColor)

	// Score/highlight style - red for emphasis
	scoreStyle := lipgloss.NewStyle().
		Foreground(secondaryColor). // red for scores
		Bold(true)

	// Team name style - cyan accent
	teamNameStyle := lipgloss.NewStyle().
		Foreground(accentColor). // cyan for team names
		Bold(true)

	// Helper to create key-value row
	kvRow := func(label, value string) string {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			labelStyle.Render(label+":"),
			valueStyle.Render(" "+value),
		)
	}

	var sections []string

	// === MATCH INFO SECTION ===
	homeTeam := details.HomeTeam.ShortName
	if homeTeam == "" {
		homeTeam = details.HomeTeam.Name
	}
	awayTeam := details.AwayTeam.ShortName
	if awayTeam == "" {
		awayTeam = details.AwayTeam.Name
	}

	// Status text with brand colors
	var statusText string
	var statusStyle lipgloss.Style
	switch details.Status {
	case api.MatchStatusFinished:
		statusText = "FT"
		statusStyle = lipgloss.NewStyle().Foreground(accentColor) // cyan for completed
	case api.MatchStatusLive:
		if details.LiveTime != nil {
			statusText = *details.LiveTime
		} else {
			statusText = "LIVE"
		}
		statusStyle = lipgloss.NewStyle().Foreground(secondaryColor).Bold(true) // red for live
	case api.MatchStatusNotStarted:
		if details.MatchTime != nil {
			statusText = details.MatchTime.Format("15:04")
		} else {
			statusText = "TBD"
		}
		statusStyle = lipgloss.NewStyle().Foreground(dimColor)
	default:
		statusText = string(details.Status)
		statusStyle = lipgloss.NewStyle().Foreground(dimColor)
	}

	matchInfoLines := []string{
		headerStyle.Render("Match Info"),
	}

	// Score line - red highlight for scores, cyan for team names
	if details.HomeScore != nil && details.AwayScore != nil {
		scoreDisplay := lipgloss.JoinHorizontal(lipgloss.Center,
			teamNameStyle.Render(homeTeam),
			scoreStyle.Render(fmt.Sprintf(" %d - %d ", *details.HomeScore, *details.AwayScore)),
			teamNameStyle.Render(awayTeam),
		)
		matchInfoLines = append(matchInfoLines, lipgloss.NewStyle().Width(contentWidth-4).Align(lipgloss.Center).Render(scoreDisplay))
	} else {
		matchupDisplay := lipgloss.JoinHorizontal(lipgloss.Center,
			teamNameStyle.Render(homeTeam),
			valueStyle.Render(" vs "),
			teamNameStyle.Render(awayTeam),
		)
		matchInfoLines = append(matchInfoLines, lipgloss.NewStyle().Width(contentWidth-4).Align(lipgloss.Center).Render(matchupDisplay))
	}

	matchInfoLines = append(matchInfoLines,
		kvRow("Status", statusStyle.Render(statusText)),
	)

	if details.League.Name != "" {
		matchInfoLines = append(matchInfoLines, kvRow("League", details.League.Name))
	}

	if details.Venue != "" {
		matchInfoLines = append(matchInfoLines, kvRow("Venue", details.Venue))
	}

	if details.MatchTime != nil {
		matchInfoLines = append(matchInfoLines, kvRow("Date", details.MatchTime.Format("02 Jan 2006")))
	}

	if details.HalfTimeScore != nil && details.HalfTimeScore.Home != nil && details.HalfTimeScore.Away != nil {
		htScore := fmt.Sprintf("%d - %d", *details.HalfTimeScore.Home, *details.HalfTimeScore.Away)
		matchInfoLines = append(matchInfoLines, kvRow("Half-Time", htScore))
	}

	sections = append(sections, sectionStyle.Render(strings.Join(matchInfoLines, "\n")))

	// === GOALS SECTION ===
	var homeGoals, awayGoals []api.MatchEvent
	for _, event := range details.Events {
		if event.Type == "goal" {
			if event.Team.ID == details.HomeTeam.ID {
				homeGoals = append(homeGoals, event)
			} else {
				awayGoals = append(awayGoals, event)
			}
		}
	}

	if len(homeGoals) > 0 || len(awayGoals) > 0 {
		goalsLines := []string{
			headerStyle.Render("Goals"),
		}

		// Goal minute style - red accent
		minuteStyle := lipgloss.NewStyle().Foreground(secondaryColor).Bold(true)

		// Home team goals
		if len(homeGoals) > 0 {
			goalsLines = append(goalsLines, teamNameStyle.Render(homeTeam))
			for _, goal := range homeGoals {
				player := "Unknown"
				if goal.Player != nil {
					player = *goal.Player
				}
				goalLine := lipgloss.JoinHorizontal(lipgloss.Left,
					minuteStyle.Render(fmt.Sprintf("  %d'", goal.Minute)),
					valueStyle.Render("  "+player),
				)
				goalsLines = append(goalsLines, goalLine)
			}
		}

		// Away team goals
		if len(awayGoals) > 0 {
			goalsLines = append(goalsLines, teamNameStyle.Render(awayTeam))
			for _, goal := range awayGoals {
				player := "Unknown"
				if goal.Player != nil {
					player = *goal.Player
				}
				goalLine := lipgloss.JoinHorizontal(lipgloss.Left,
					minuteStyle.Render(fmt.Sprintf("  %d'", goal.Minute)),
					valueStyle.Render("  "+player),
				)
				goalsLines = append(goalsLines, goalLine)
			}
		}

		sections = append(sections, sectionStyle.Render(strings.Join(goalsLines, "\n")))
	}

	// === CARDS SECTION ===
	var homeYellow, homeRed, awayYellow, awayRed int
	for _, event := range details.Events {
		if event.Type == "card" {
			isHome := event.Team.ID == details.HomeTeam.ID
			if event.EventType != nil {
				switch *event.EventType {
				case "yellow":
					if isHome {
						homeYellow++
					} else {
						awayYellow++
					}
				case "red":
					if isHome {
						homeRed++
					} else {
						awayRed++
					}
				}
			}
		}
	}

	totalCards := homeYellow + homeRed + awayYellow + awayRed
	if totalCards > 0 {
		cardsLines := []string{
			headerStyle.Render("Cards"),
		}

		// Card visual styles - cyan for yellow cards, red for red cards (brand colors)
		yellowCardStyle := lipgloss.NewStyle().Foreground(accentColor) // cyan blocks for yellow
		redCardStyle := lipgloss.NewStyle().Foreground(secondaryColor) // red blocks for red

		// Team card summary with visual bars
		homeCardsStr := ""
		if homeYellow > 0 {
			homeCardsStr += yellowCardStyle.Render(strings.Repeat("▪", homeYellow))
		}
		if homeRed > 0 {
			if homeCardsStr != "" {
				homeCardsStr += " "
			}
			homeCardsStr += redCardStyle.Render(strings.Repeat("▪", homeRed))
		}
		if homeCardsStr == "" {
			homeCardsStr = lipgloss.NewStyle().Foreground(dimColor).Render("-")
		}

		awayCardsStr := ""
		if awayYellow > 0 {
			awayCardsStr += yellowCardStyle.Render(strings.Repeat("▪", awayYellow))
		}
		if awayRed > 0 {
			if awayCardsStr != "" {
				awayCardsStr += " "
			}
			awayCardsStr += redCardStyle.Render(strings.Repeat("▪", awayRed))
		}
		if awayCardsStr == "" {
			awayCardsStr = lipgloss.NewStyle().Foreground(dimColor).Render("-")
		}

		cardsLines = append(cardsLines, kvRow(homeTeam, homeCardsStr))
		cardsLines = append(cardsLines, kvRow(awayTeam, awayCardsStr))

		// Summary row
		summaryStyle := lipgloss.NewStyle().Foreground(dimColor).Italic(true)
		totalSummary := fmt.Sprintf("Total: %d yellow, %d red", homeYellow+awayYellow, homeRed+awayRed)
		cardsLines = append(cardsLines, summaryStyle.Render(totalSummary))

		sections = append(sections, sectionStyle.Render(strings.Join(cardsLines, "\n")))
	}

	// === MATCH STATS SECTION ===
	statsLines := []string{
		headerStyle.Render("Match Stats"),
	}

	// Event counts as simple stats
	var subs int
	for _, event := range details.Events {
		if event.Type == "substitution" {
			subs++
		}
	}

	// Stats with cyan accent for numbers
	statNumStyle := lipgloss.NewStyle().Foreground(accentColor)
	statsLines = append(statsLines, kvRow("Events", statNumStyle.Render(fmt.Sprintf("%d", len(details.Events)))))
	if subs > 0 {
		statsLines = append(statsLines, kvRow("Substitutions", statNumStyle.Render(fmt.Sprintf("%d", subs))))
	}

	// Only show stats section if there's meaningful data
	if len(details.Events) > 0 {
		sections = append(sections, sectionStyle.Render(strings.Join(statsLines, "\n")))
	}

	// Combine all sections with spacing
	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	return panelStyle.
		Width(width).
		Height(height).
		Render(content)
}
