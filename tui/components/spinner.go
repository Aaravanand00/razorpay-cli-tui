package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

func NewSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(styles.ColorPrimary)
	return s
}

func RenderLoading(s spinner.Model, text string) string {
	if text == "" {
		text = "Communicating with Razorpay API..."
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		s.View(),
		" ",
		lipgloss.NewStyle().Foreground(styles.ColorMuted).Render(text),
	)
}
