package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

func RenderHeader(s *state.SessionState, width int) string {
	logo := styles.BrandStyle.Render("⚡ RAZORPAY TUI")

	// Dynamic Navigation Steps Indicator
	var step1, step2, step3 string
	arrow := styles.StepArrow.Render("➔")

	switch s.CurrentScreen {
	case state.ScreenHome:
		step1 = styles.StepActiveStyle.Render("● 1. MODULES")
		step2 = styles.StepInactiveStyle.Render("○ 2. ACTIONS")
		step3 = styles.StepInactiveStyle.Render("○ 3. DATA/FORM")
	case state.ScreenActions:
		step1 = styles.StepDoneStyle.Render("✓ 1. MODULES")
		step2 = styles.StepActiveStyle.Render("● 2. ACTIONS")
		step3 = styles.StepInactiveStyle.Render("○ 3. DATA/FORM")
	case state.ScreenConfig:
		step1 = styles.StepDoneStyle.Render("✓ 1. MODULES")
		step2 = styles.StepActiveStyle.Render("● 2. CONFIGURATION")
		step3 = styles.StepInactiveStyle.Render("○ 3. CREDENTIALS")
	default:
		step1 = styles.StepDoneStyle.Render("✓ 1. MODULES")
		step2 = styles.StepDoneStyle.Render("✓ 2. ACTIONS")
		step3 = styles.StepActiveStyle.Render("● 3. ACTIVE VIEW")
	}

	stepper := lipgloss.JoinHorizontal(lipgloss.Center, " ", step1, arrow, step2, arrow, step3)

	// Auth, Mode & Read-Only Badge
	var badge string
	if !s.HasCredentials() {
		if s.IsReadOnly {
			badge = styles.BadgeNoAuthStyle.Render("🛡️ READ-ONLY (NO KEYS)")
		} else if s.ActiveMode == "live" {
			badge = styles.BadgeLiveStyle.Render("● LIVE (NO KEYS CONFIGURED)")
		} else {
			badge = styles.BadgeTestStyle.Render("▲ TEST (NO KEYS CONFIGURED)")
		}
	} else if s.IsReadOnly {
		if s.IsLiveMode {
			badge = styles.BadgeLiveStyle.Render("🛡️ LIVE (READ-ONLY): " + s.MaskedKey())
		} else {
			badge = styles.BadgeTestStyle.Render("🛡️ TEST (READ-ONLY): " + s.MaskedKey())
		}
	} else if s.IsLiveMode {
		badge = styles.BadgeLiveStyle.Render("● LIVE MODE: " + s.MaskedKey())
	} else {
		badge = styles.BadgeTestStyle.Render("▲ TEST MODE: " + s.MaskedKey())
	}

	leftTop := lipgloss.JoinHorizontal(lipgloss.Center, logo, stepper)
	gapTop := width - lipgloss.Width(leftTop) - lipgloss.Width(badge) - 4
	if gapTop < 2 {
		gapTop = 2
	}
	line1 := lipgloss.JoinHorizontal(lipgloss.Center, leftTop, strings.Repeat(" ", gapTop), badge)

	// Location / Breadcrumb Line
	crumbs := strings.Join(s.Breadcrumbs, "  ❯  ")
	locationText := lipgloss.JoinHorizontal(
		lipgloss.Center,
		styles.LocationBar.Render("📍 Navigation: "),
		styles.LocationActive.Render(crumbs),
	)

	fullHeader := lipgloss.JoinVertical(lipgloss.Left, line1, locationText)
	return styles.HeaderContainer.Width(width - 2).Render(fullHeader)
}
