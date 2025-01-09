package styles

import "github.com/charmbracelet/lipgloss"

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
		BorderForeground(lipgloss.Color(StrongOrange)).
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
