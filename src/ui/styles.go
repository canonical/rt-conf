package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	// weakOrange   = "#DF4A16" // WARN: Non offical CANONICAL color
	strongOrange = "#E95420"

	// RGB: (140, 51, 19) == 60% of strongOrange
	weakOrange = "#8C3313" // WARN: Non offical CANONICAL color

	green     = "#3EB34F"
	porcelain = "#F7F7F7"
	ash       = "#888888"
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
			Padding(1, 2)
		// Foreground(lipgloss.Color("#ED3146"))
		//Foreground(lipgloss.Color("#2D3748"))

	titleStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(porcelain)).
			Background(lipgloss.Color(strongOrange)).
			Padding(0, 1).
			Bold(true)

	// statusMessageStyle = lipgloss.NewStyle().
	// 			Foreground(
	// 		lipgloss.AdaptiveColor{
	// 			Light: green,
	// 			Dark:  green,
	// 		},
	// 	).
	// 	Render

	infoMessageStyle = lipgloss.NewStyle().
				Foreground(
			lipgloss.AdaptiveColor{
				Light: ash,
				Dark:  ash,
			},
		).
		Bold(true).
		Render

	errorMessageStyle = lipgloss.NewStyle().
				Foreground(
			lipgloss.AdaptiveColor{
				Light: "#ED3146",
				Dark:  "#ED3146",
			},
		).
		Bold(true).
		Render

	innerMenuStyle = lipgloss.
			NewStyle().
			Foreground(
			lipgloss.AdaptiveColor{
				Light: porcelain,
				Dark:  porcelain,
			}).
		Background(
			lipgloss.AdaptiveColor{
				Light: strongOrange,
				Dark:  strongOrange,
			}).
		Bold(true).
		Render

	focusedStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(strongOrange)).
			Padding(0, 1)

	blurredStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("240"))

	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Render("[ Apply ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Apply"))
)
