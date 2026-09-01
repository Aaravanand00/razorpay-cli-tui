package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

func RenderHeader(s *state.SessionState, width int) string {
	logo := styles.BrandStyle.Render("⚡ RAZORPAY CLI")

	// Breadcrumb
	crumbs := strings.Join(s.Breadcrumbs, " > ")
	breadcrumbView := styles.BreadcrumbStyle.Render(" " + crumbs)

	// Mode badge
	var badge string
	if s.KeyID == "" {
		badge = styles.BadgeNoAuthStyle.Render("! NO CREDENTIALS")
	} else if s.IsLiveMode {
		badge = styles.BadgeLiveStyle.Render("● LIVE MODE")
	} else {
		badge = styles.BadgeTestStyle.Render("▲ TEST MODE")
	}

	leftSide := lipgloss.JoinHorizontal(lipgloss.Center, logo, breadcrumbView)
	rightSide := badge

	gap := width - lipgloss.Width(leftSide) - lipgloss.Width(rightSide) - 4
	if gap < 2 {
		gap = 2
	}

	content := lipgloss.JoinHorizontal(lipgloss.Center, leftSide, strings.Repeat(" ", gap), rightSide)
	return styles.HeaderContainer.Width(width - 2).Render(content)
}
