package ui

import "github.com/charmbracelet/lipgloss"

const (
	// weakOrange   = "#DF4A16" // WARN: Non offical CANONICAL color
	strongOrange = "#E95420"

	// RGB: (140, 51, 19) == 60% of strongOrange
	weakOrange = "#8C3313" // WARN: Non offical CANONICAL color

	green = "#3EB34F"
)

var (
	// ------------------------------- DefaultItemStyles ------------------------
	NormalTitle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{
			Light: "#1a1a1a",
			Dark:  "#dddddd",
		}).
		Padding(0, 0, 0, 2). //nolint:mnd
		Bold(true)

	NormalDesc = NormalTitle.
			Foreground(
			lipgloss.AdaptiveColor{
				Light: "#A49FA5",
				Dark:  "#777777",
			})

	SelectedTitle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(
			lipgloss.AdaptiveColor{
				Light: strongOrange,
				Dark:  strongOrange,
			}).
		Foreground(lipgloss.AdaptiveColor{
			Light: strongOrange,
			Dark:  strongOrange,
		}).
		Padding(0, 0, 0, 1).
		Bold(true)

	SelectedDesc = SelectedTitle.
			Foreground(
			lipgloss.AdaptiveColor{
				Light: weakOrange,
				Dark:  weakOrange,
			})

	DimmedTitle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{
			Light: "#A49FA5",
			Dark:  "#777777",
		}).
		Padding(0, 0, 0, 2) //nolint:mnd

	DimmedDesc = DimmedTitle.
			Foreground(
			lipgloss.AdaptiveColor{
				Light: strongOrange,
				Dark:  "#4D4D4D",
			})

	FilterMatch = lipgloss.NewStyle().Underline(true)

	appStyle = lipgloss.
			NewStyle().
			Padding(1, 2).
			Foreground(lipgloss.Color("#ED3146"))
		//Foreground(lipgloss.Color("#2D3748"))

	titleStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("#F7F7F7")).
			Background(lipgloss.Color(strongOrange)).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(
			lipgloss.AdaptiveColor{
				Light: green,
				Dark:  green,
			},
		).
		Render
)
