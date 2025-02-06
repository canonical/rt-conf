package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	// ------------------------------- DefaultItemStyles ------------------------

	StatusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

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
				Light: StrongOrange,
				Dark:  StrongOrange,
			}).
		Foreground(lipgloss.AdaptiveColor{
			Light: StrongOrange,
			Dark:  StrongOrange,
		}).
		Padding(0, 0, 0, 1).
		Bold(true)

	SelectedDesc = SelectedTitle.
			Foreground(
			lipgloss.AdaptiveColor{
				Light: WeakOrange,
				Dark:  WeakOrange,
			})

	Section = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{
			Light: "#A49FA5",
			Dark:  "#777777",
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
				Light: StrongOrange,
				Dark:  "#4D4D4D",
			})

	FilterMatch = lipgloss.NewStyle().Underline(true)

	AppStyle = lipgloss.
			NewStyle().
			Padding(1, 2)

	TitleStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(Porcelain)).
			Background(lipgloss.Color(StrongOrange)).
			Padding(0, 1).
			Bold(true)

	InfoMessageStyle = lipgloss.NewStyle().
				Foreground(
			lipgloss.AdaptiveColor{
				Light: Ash,
				Dark:  Ash,
			},
		).
		Bold(true).
		Render

	ErrorMessageStyle = lipgloss.NewStyle().
				Foreground(
			lipgloss.AdaptiveColor{
				Light: "#ED3146",
				Dark:  "#ED3146",
			},
		).
		Bold(true).
		Render

	InnerMenuStyle = lipgloss.
			NewStyle().
			Foreground(
			lipgloss.AdaptiveColor{
				Light: Porcelain,
				Dark:  Porcelain,
			}).
		Background(
			lipgloss.AdaptiveColor{
				Light: StrongOrange,
				Dark:  StrongOrange,
			}).
		Bold(true).
		Render

	FocusedStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color(StrongOrange))

	BlurredStyle = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("240"))

	CursorStyle         = FocusedStyle
	NoStyle             = lipgloss.NewStyle()
	HelpStyle           = BlurredStyle
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	ButtonPadding = []int{0, 2, 0, 0}

	// TODO: refactor this in to be scalable (Is there a more scalable way?)
	FocusedBackButton = FocusedStyle.Padding(ButtonPadding...).
				Render("[ Back ]")

	BlurredBackButton = lipgloss.NewStyle().
				Padding(ButtonPadding...).
				Render(
			fmt.Sprintf("[ %s ]",
				BlurredStyle.Render("Back"),
			),
		)

	FocusedApplyButton = FocusedStyle.Padding(ButtonPadding...).
				Render("[ Apply ]")

	BlurredApplyButton = lipgloss.NewStyle().
				Padding(ButtonPadding...).
				Render(
			fmt.Sprintf("[ %s ]",
				BlurredStyle.Render("Apply"),
			),
		)
)
