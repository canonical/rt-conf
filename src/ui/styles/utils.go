package styles

import "github.com/charmbracelet/lipgloss"

func JoinHorizontal(left, right *string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, *left, *right)
}

func CenteredSquareWithText(
	appWidth, appHeight, textWidth, textHeight int,
	content string) string {
	// Calculate dimensions for the rectangles
	outerWidth := appWidth
	outerHeight := appHeight

	innerWidth := textWidth
	innerHeight := textHeight

	// Define styles for the rectangles
	outerRectangleStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true).
		BorderForeground(lipgloss.Color(StrongOrange)).
		Width(outerWidth).
		Height(outerHeight)

	innerRectangleStyle := lipgloss.NewStyle().
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
