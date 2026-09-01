package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

type KeyHelp struct {
	Key  string
	Desc string
}

func RenderFooter(s *state.SessionState, width int, keys []KeyHelp) string {
	if len(keys) == 0 {
		keys = []KeyHelp{
			{Key: "↑/↓", Desc: "Navigate"},
			{Key: "Enter", Desc: "Select"},
			{Key: "/", Desc: "Filter"},
			{Key: "Esc", Desc: "Back"},
			{Key: "q", Desc: "Quit"},
		}
	}

	var items []string
	for _, k := range keys {
		item := lipgloss.JoinHorizontal(
			lipgloss.Center,
			styles.KeyStyle.Render("["+k.Key+"]"),
			styles.KeyDescStyle.Render(" "+k.Desc),
		)
		items = append(items, item)
	}

	content := strings.Join(items, "  •  ")
	return styles.FooterContainer.Width(width - 2).Render(content)
}
