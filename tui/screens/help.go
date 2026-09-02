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

	// 4. Main Body Content (Organized 2x2 Grid or Stacked Cards)
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

	// Card 1: Global Navigation
	card1Content := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary).Render("🌐 Global Navigation") + "\n\n" +
		formatShortcutRow("↑ / ↓ / j / k", "Move cursor & navigate lists") + "\n" +
		formatShortcutRow("Enter", "Select / Open / Submit") + "\n" +
		formatShortcutRow("Esc", "Go back to previous screen") + "\n" +
		formatShortcutRow("q / Ctrl+C", "Quit Razorpay TUI") + "\n" +
		formatShortcutRow("? / h", "Toggle this Help modal") + "\n" +
		formatShortcutRow("c", "Open Credentials Configuration")

	card1 := cardStyle.Render(card1Content)

	// Card 2: Data Table & Inspection
	card2Content := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary).Render("📊 Data Table & JSON Inspector") + "\n\n" +
		formatShortcutRow("↑ / ↓", "Select row in table") + "\n" +
		formatShortcutRow("Enter", "Inspect record (Screen 4)") + "\n" +
		formatShortcutRow("r", "Refresh table via live REST API") + "\n" +
		formatShortcutRow("Tab / v / 1 / 2", "Toggle Summary ⇄ Raw JSON") + "\n" +
		formatShortcutRow("PgUp / PgDn", "Scroll large JSON payloads") + "\n" +
		formatShortcutRow("/", "Search / filter actions")

	card2 := cardStyle.Render(card2Content)

	// Card 3: Form Builder & Execution
	card3Content := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary).Render("📝 Forms & API Execution") + "\n\n" +
		formatShortcutRow("Tab / ↓", "Focus next input field") + "\n" +
		formatShortcutRow("Shift+Tab / ↑", "Focus previous input field") + "\n" +
		formatShortcutRow("Enter", "Submit & execute live API call") + "\n" +
		formatShortcutRow("r", "Reset form / new transaction") + "\n" +
		formatShortcutRow("Enter (on success)", "View created record details")

	card3 := cardStyle.Render(card3Content)

	// Card 4: Environments & Flags
	card4Content := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary).Render("🛡️ Environments & CLI Flags") + "\n\n" +
		formatShortcutRow("▲ TEST MODE", "Test Sandbox (rzp_test_...)") + "\n" +
		formatShortcutRow("● LIVE MODE", "Live Production (rzp_live_...)") + "\n" +
		formatShortcutRow("🛡️ READ-ONLY", "--read-only (blocks mutations)") + "\n" +
		formatShortcutRow("← / → / 1 / 2", "Switch mode in Config Screen") + "\n" +
		formatShortcutRow("go run . tui", "Launch Razorpay Terminal UI")

	card4 := cardStyle.Render(card4Content)

	row1 := lipgloss.JoinHorizontal(lipgloss.Top, card1, "  ", card2)
	row2 := lipgloss.JoinHorizontal(lipgloss.Top, card3, "  ", card4)

	mainContent := "\n" + title + "\n\n" + row1 + "\n\n" + row2

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
