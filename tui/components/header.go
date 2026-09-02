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

	// Auth & Mode Badge
	var badge string
	if s.KeyID == "" {
		badge = styles.BadgeNoAuthStyle.Render("! NO CREDENTIALS")
	} else if s.IsLiveMode {
		maskedKey := s.KeyID
		if len(maskedKey) > 12 {
			maskedKey = maskedKey[:9] + "..." + maskedKey[len(maskedKey)-4:]
		}
		badge = styles.BadgeLiveStyle.Render("● LIVE: " + maskedKey)
	} else {
		maskedKey := s.KeyID
		if len(maskedKey) > 12 {
			maskedKey = maskedKey[:9] + "..." + maskedKey[len(maskedKey)-4:]
		}
		badge = styles.BadgeTestStyle.Render("▲ TEST: " + maskedKey)
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
