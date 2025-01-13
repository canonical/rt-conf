package components

import (
	"fmt"

	"github.com/canonical/rt-conf/src/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

type Button struct {
	text   string
	button string
}

func NewButton(text string) *Button {
	return &Button{
		text: text,
	}
}

func (b *Button) SetFocused() {
	b.button = FocusedButton(b.text)
}

func (b *Button) SetBlurred() {
	b.button = BlurredButton(b.text)
}
func (b *Button) Render() *string {
	return &b.button
}

// BlurredButton returns a blurred button
func BlurredButton(text string) string {
	s := lipgloss.NewStyle().
		Padding(styles.ButtonPadding...).
		Render(
			fmt.Sprintf("[ %s ]",
				styles.BlurredStyle.Render(text),
			),
		)
	return s
}

// FocusedButton returns a focused button
func FocusedButton(text string) string {
	s := styles.FocusedStyle.
		Padding(styles.ButtonPadding...).
		Render(fmt.Sprintf("[ %s ]", text))
	return s
}
