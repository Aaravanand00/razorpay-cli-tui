package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/tui/components"
	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

type HelpScreen struct {
	state       *state.SessionState
	width       int
	height      int
	initialized bool
}

func NewHelpScreen(s *state.SessionState, width, height int) HelpScreen {
	return HelpScreen{
		state:       s,
		width:       width,
		height:      height,
		initialized: true,
	}
}

func (h *HelpScreen) SetSize(width, height int) {
	if !h.initialized || h.state == nil {
		return
	}
	h.width = width
	h.height = height
}

func (h *HelpScreen) Update(msg tea.Msg) (HelpScreen, tea.Cmd) {
	h.SetSize(h.state.Width, h.state.Height)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "?", "h", "q":
			h.state.PopScreen()
			return *h, nil
		}
	}

	return *h, nil
}

func (h HelpScreen) View() string {
	var sections []string

	// 1. Header
	headerView := components.RenderHeader(h.state, h.width)
	sections = append(sections, headerView)

	// 2. Toast
	toastView := components.RenderToast(h.state)
	if toastView != "" {
		sections = append(sections, toastView)
	}

	// 3. Footer Keys
	keys := []components.KeyHelp{
		{Key: "Esc / ? / q", Desc: "Close Help"},
		{Key: "c", Desc: "Config"},
		{Key: "Ctrl+C", Desc: "Quit TUI"},
	}
	footerView := components.RenderFooter(h.state, h.width, keys)

	// Height math
	headerHeight := lipgloss.Height(headerView)
	footerHeight := lipgloss.Height(footerView)
	toastHeight := 0
	if toastView != "" {
		toastHeight = lipgloss.Height(toastView)
	}

	bodyHeight := h.height - headerHeight - footerHeight - toastHeight
	if bodyHeight < 10 {
		bodyHeight = 10
	}

	// 4. Main Body Content (2 Side-by-Side Comprehensive Cards)
	title := styles.TitleStyle.PaddingLeft(1).Render("❓ Razorpay Terminal UI — Keyboard Shortcuts & Navigation Guide")

	cardWidth := (h.width - 8) / 2
	if cardWidth < 36 {
		cardWidth = 36
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorCardBg).
		Padding(1, 2).
		Width(cardWidth)

	// Left Column: Navigation & Inspection
	col1Content := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary).Render("🌐 Navigation & Inspection Shortcuts") + "\n\n" +
		formatShortcutRow("↑ / ↓ / j / k", "Move cursor & navigate lists") + "\n" +
		formatShortcutRow("Enter", "Select / Open action / View details") + "\n" +
		formatShortcutRow("Esc", "Go back to previous screen") + "\n" +
		formatShortcutRow("q / Ctrl+C", "Quit Razorpay TUI") + "\n" +
		formatShortcutRow("? / h", "Toggle this Help modal") + "\n" +
		formatShortcutRow("c", "Open Credentials Configuration") + "\n" +
		formatShortcutRow("/", "Fuzzy search actions & modules") + "\n" +
		formatShortcutRow("r", "Refresh table via live REST API") + "\n" +
		formatShortcutRow("Tab / v / 1 / 2", "Toggle Summary ⇄ Raw JSON") + "\n" +
		formatShortcutRow("PgUp / PgDn", "Scroll large JSON payloads")

	col1 := cardStyle.Render(col1Content)

	// Right Column: Forms & Environments
	col2Content := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary).Render("📝 Forms & Environment Modes") + "\n\n" +
		formatShortcutRow("Tab / ↓", "Focus next form input field") + "\n" +
		formatShortcutRow("Shift+Tab / ↑", "Focus previous form field") + "\n" +
		formatShortcutRow("Enter", "Submit form & execute live API") + "\n" +
		formatShortcutRow("▲ TEST MODE", "Test Sandbox (rzp_test_...)") + "\n" +
		formatShortcutRow("● LIVE MODE", "Live Production (rzp_live_...)") + "\n" +
		formatShortcutRow("🛡️ READ-ONLY", "--read-only (locks write ops)") + "\n" +
		formatShortcutRow("← / → / 1 / 2", "Switch mode in Config Screen") + "\n" +
		formatShortcutRow("go run . tui", "Launch Razorpay Terminal UI")

	col2 := cardStyle.Render(col2Content)

	gridRow := lipgloss.JoinHorizontal(lipgloss.Top, col1, "  ", col2)

	mainContent := "\n" + title + "\n\n" + gridRow

	bodyContainer := lipgloss.NewStyle().
		Height(bodyHeight).
		Width(h.width - 2)

	sections = append(sections, bodyContainer.Render(mainContent))
	sections = append(sections, footerView)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func formatShortcutRow(key, desc string) string {
	k := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#38BDF8")).Render(fmt.Sprintf("%-16s", key))
	d := lipgloss.NewStyle().Foreground(styles.ColorText).Render(desc)
	return k + " : " + d
}
