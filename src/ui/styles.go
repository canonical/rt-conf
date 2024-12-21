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
			Foreground(lipgloss.Color(strongOrange))
		// Padding(0, 1)

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

func CenteredSquareWithText(
	appWidth, appHeight, textWidth, textHeight int,
	content string) string {
	// Calculate dimensions for the rectangles
	// outerWidth := int(float64(appWidth) * 0.9)
	// outerHeight := int(float64(appHeight) * 0.9)
	outerWidth := appWidth
	outerHeight := appHeight

	// innerWidth := int(float64(m.width) * 0.3)
	// innerHeight := int(float64(m.height) * 0.3)
	innerWidth := textWidth
	innerHeight := textHeight

	// Define styles for the rectangles
	outerRectangleStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(lipgloss.Color(strongOrange)).
		Width(outerWidth).
		Height(outerHeight)

	innerRectangleStyle := lipgloss.NewStyle().
		// Border(lipgloss.NormalBorder(), true).
		// BorderForeground(lipgloss.Color("#E95420")).
		Width(innerWidth).
		Height(innerHeight).
		Bold(true)

	// Render the inner rectangle
	innerRectangle := innerRectangleStyle.Render(content)

	// Calculate center positions for the inner rectangle within the outer rectangle
	innerX := (outerWidth - innerWidth) / 2
	innerY := (outerHeight - innerHeight) / 2

	// Combine the rectangles
	final := lipgloss.Place(
		appWidth, appHeight,
		lipgloss.Center, lipgloss.Center,
		outerRectangleStyle.Render(
			lipgloss.Place(
				innerWidth, innerHeight,
				lipgloss.Center, lipgloss.Center,
				lipgloss.NewStyle().Margin(innerY, innerX).Render(innerRectangle),
			),
		),
	)

	return final
}
