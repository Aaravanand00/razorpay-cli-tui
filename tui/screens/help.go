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
		case "a":
			h.state.PopScreen()
			h.state.PushScreen(state.ScreenAIAssist, "🤖 AI Assist")
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
		{Key: "a", Desc: "AI Assist"},
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
	subtitle := styles.SubtitleStyle.PaddingLeft(1).Render("Comprehensive cheat sheet for universal keybindings, interactive forms, and environment safety modes.")

	cardWidth := (h.width - 10) / 2
	if cardWidth < 42 {
		cardWidth = 42
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorCardBg).
		Padding(1, 3).
		Width(cardWidth)

	subHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	sectionTagStyle := lipgloss.NewStyle().Foreground(styles.ColorTextDim).Bold(true)

	// Left Column: Navigation & Inspection
	col1Content := subHeaderStyle.Render("🌐 Navigation & Shortcuts") + "\n\n" +
		sectionTagStyle.Render("── NAVIGATION & APP ──") + "\n" +
		formatShortcutRow("↑ / ↓ / j / k", "Move cursor & navigate lists") + "\n" +
		formatShortcutRow("Enter", "Select / Open action / View details") + "\n" +
		formatShortcutRow("Esc", "Go back to previous screen") + "\n" +
		formatShortcutRow("q / Ctrl+C", "Quit Razorpay TUI") + "\n\n" +
		sectionTagStyle.Render("── FEATURES & SHORTCUTS ──") + "\n" +
		formatShortcutRow("a / ctrl+a", "Open AI Assistant (Natural Language)") + "\n" +
		formatShortcutRow("/", "Fuzzy search actions & modules") + "\n" +
		formatShortcutRow("c", "Open Credentials & Mode Config") + "\n" +
		formatShortcutRow("? / h", "Toggle this Help modal") + "\n" +
		formatShortcutRow("r", "Refresh table via live REST API") + "\n" +
		formatShortcutRow("Tab / v", "Toggle Summary ⇄ Raw JSON")

	col1 := cardStyle.Render(col1Content)

	// Right Column: Forms, Modes & Environments
	col2Content := subHeaderStyle.Render("📝 Forms & Environments") + "\n\n" +
		sectionTagStyle.Render("── INTERACTIVE FORMS ──") + "\n" +
		formatShortcutRow("Tab / ↓", "Focus next form input field") + "\n" +
		formatShortcutRow("Shift+Tab / ↑", "Focus previous form field") + "\n" +
		formatShortcutRow("Enter", "Submit form & execute live API") + "\n\n" +
		sectionTagStyle.Render("── ENVIRONMENT MODES & SAFETY ──") + "\n" +
		formatShortcutRow("▲ TEST MODE", "Test Sandbox (rzp_test_...)") + "\n" +
		formatShortcutRow("● LIVE MODE", "Live Production (rzp_live_...)") + "\n" +
		formatShortcutRow("🛡️ READ-ONLY", "Safe Mode (locks write actions)") + "\n" +
		formatShortcutRow("← / → / 1 / 2", "Switch active profile in Config") + "\n" +
		formatShortcutRow("go run . tui", "Launch Razorpay Terminal UI")

	col2 := cardStyle.Render(col2Content)

	gridRow := lipgloss.JoinHorizontal(lipgloss.Top, col1, "    ", col2)

	mainContent := "\n" + title + "\n" + subtitle + "\n\n" + gridRow

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
