package aiassist

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/config"
	"github.com/razorpay/razorpay-cli/tui/components"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

func (m Model) View() string {
	var sections []string

	// 1. Top Header
	headerView := components.RenderHeader(m.state, m.width)
	sections = append(sections, headerView)

	// 2. Toast Alert
	toastView := components.RenderToast(m.state)
	if toastView != "" {
		sections = append(sections, toastView)
	}

	// 3. Footer Keys
	var keys []components.KeyHelp
	apiKey := config.AIApiKey()
	if apiKey == "" {
		keys = []components.KeyHelp{
			{Key: "c", Desc: "Configure API Key"},
			{Key: "Esc", Desc: "Back"},
			{Key: "Ctrl+C", Desc: "Quit"},
		}
	} else {
		switch m.currentState {
		case StateInput:
			keys = []components.KeyHelp{
				{Key: "Enter", Desc: "Generate Command"},
				{Key: "c", Desc: "Config"},
				{Key: "?/h", Desc: "Help"},
				{Key: "Esc", Desc: "Back"},
				{Key: "Ctrl+C", Desc: "Quit"},
			}
		case StateLoading:
			keys = []components.KeyHelp{
				{Key: "Esc", Desc: "Cancel"},
				{Key: "Ctrl+C", Desc: "Quit"},
			}
		case StatePreview:
			keys = []components.KeyHelp{
				{Key: "Enter", Desc: "Confirm & Execute"},
				{Key: "Esc / r", Desc: "Edit Prompt"},
				{Key: "c", Desc: "Config"},
				{Key: "Ctrl+C", Desc: "Quit"},
			}
		case StateError:
			keys = []components.KeyHelp{
				{Key: "Enter / r", Desc: "Retry"},
				{Key: "c", Desc: "Config"},
				{Key: "Esc", Desc: "Back"},
				{Key: "Ctrl+C", Desc: "Quit"},
			}
		}
	}
	footerView := components.RenderFooter(m.state, m.width, keys)

	// Height math
	headerHeight := lipgloss.Height(headerView)
	footerHeight := lipgloss.Height(footerView)
	toastHeight := 0
	if toastView != "" {
		toastHeight = lipgloss.Height(toastView)
	}

	bodyHeight := m.height - headerHeight - footerHeight - toastHeight
	if bodyHeight < 10 {
		bodyHeight = 10
	}

	// 4. Body Content
	var mainContent string
	if apiKey == "" {
		mainContent = m.renderNotConfiguredCard()
	} else {
		switch m.currentState {
		case StateInput:
			mainContent = m.renderInputView()
		case StateLoading:
			mainContent = m.renderLoadingView()
		case StatePreview:
			mainContent = m.renderPreviewView()
		case StateError:
			mainContent = m.renderErrorView()
		}
	}

	bodyContainer := lipgloss.NewStyle().
		Height(bodyHeight).
		Width(m.width - 2)

	sections = append(sections, bodyContainer.Render(mainContent))
	sections = append(sections, footerView)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) renderNotConfiguredCard() string {
	title := styles.TitleStyle.PaddingLeft(1).Render("🤖 Razorpay AI Assistant — Natural Language CLI")

	cardContent := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorWarning).Render("🔒 AI Assist Not Configured") + "\n\n" +
		lipgloss.NewStyle().Foreground(styles.ColorText).Render(
			"To translate natural-language prompts (e.g., 'find failed payments from yesterday') into\n"+
				"exact CLI commands and flags with Claude Haiku, please configure your Anthropic API Key.\n\n"+
				"👉 Press [c] to open Configuration and enter your API Key,\n"+
				"   or set the environment variable: export RAZORPAY_AI_API_KEY=sk-ant-...\n")

	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorWarning).
		Background(styles.ColorCardBg).
		Padding(2, 3).
		Width(m.width - 6).
		Render(cardContent)

	return "\n" + title + "\n\n" + card
}

func (m Model) renderInputView() string {
	title := styles.TitleStyle.PaddingLeft(1).Render("🤖 Razorpay AI Assistant — Natural Language CLI Builder")

	inputCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorderFocus).
		Background(styles.ColorCardBg).
		Padding(1, 2).
		Width(m.width - 6).
		Render(m.input.View())

	examplesContent := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary).Render("💡 Example Prompts:") + "\n\n" +
		lipgloss.NewStyle().Foreground(styles.ColorTextMuted).Render(
			" • 'find failed payments from yesterday'\n"+
				" • 'create a ₹500 order for VIP customer with receipt #101'\n"+
				" • 'list all customers created last week'\n"+
				" • 'issue full refund for payment pay_Lxxxxxxxx'\n"+
				" • 'show recent settlements'")

	examplesCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorCardBg).
		Padding(1, 2).
		Width(m.width - 6).
		Render(examplesContent)

	hint := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		PaddingLeft(1).
		Render("Type your request in plain English or Hindi • Press [Enter] to generate CLI command")

	return "\n" + title + "\n\n" + inputCard + "\n\n" + examplesCard + "\n\n" + hint
}

func (m Model) renderLoadingView() string {
	title := styles.TitleStyle.PaddingLeft(1).Render("🤖 Razorpay AI Assistant — Translating Request")

	loadingText := fmt.Sprintf("%s  Translating prompt into Razorpay CLI command using Claude Haiku AI...", m.spinner.View())

	loadingCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Background(styles.ColorCardBg).
		Padding(2, 3).
		Width(m.width - 6).
		Render(loadingText)

	return "\n" + title + "\n\n" + loadingCard
}

func (m Model) renderPreviewView() string {
	if m.suggestion == nil {
		return m.renderInputView()
	}

	title := styles.TitleStyle.PaddingLeft(1).Render("🤖 AI Suggested CLI Command — Preview & Confirm")

	cmdFormatted := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#38BDF8")).Render(m.suggestion.FullCommandString())
	explFormatted := lipgloss.NewStyle().Foreground(styles.ColorText).Render(m.suggestion.Explanation)

	confBadge := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSuccess).Render("● High Confidence")
	if m.suggestion.IsLowConfidence() {
		confBadge = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorWarning).Render("▲ Low Confidence / Destructive Operation")
	}

	content := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary).Render("💻 Generated Command:") + "\n" +
		"   " + cmdFormatted + "\n\n" +
		lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary).Render("📖 Explanation:") + "\n" +
		"   " + explFormatted + "\n\n" +
		lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary).Render("🎯 Review Assessment:") + "\n" +
		"   " + confBadge

	previewCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorderFocus).
		Background(styles.ColorCardBg).
		Padding(1, 3).
		Width(m.width - 6).
		Render(content)

	var warningCard string
	if m.suggestion.IsLowConfidence() {
		warningText := "⚠️  This action modifies data or is ambiguous. Please review carefully before confirming."
		warningCard = "\n\n" + lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorWarning).
			Background(styles.ColorCardBg).
			Foreground(styles.ColorWarning).
			Bold(true).
			Padding(1, 2).
			Width(m.width - 6).
			Render(warningText)
	}

	confirmHint := lipgloss.NewStyle().
		Foreground(styles.ColorSuccess).
		Bold(true).
		PaddingLeft(1).
		Render("👉 Press [Enter] to Execute Command via Safe Unified Path • Press [Esc] or [r] to Edit Prompt")

	return "\n" + title + "\n\n" + previewCard + warningCard + "\n\n" + confirmHint
}

func (m Model) renderErrorView() string {
	title := styles.TitleStyle.PaddingLeft(1).Render("🤖 Razorpay AI Assistant — Error")

	errContent := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorError).Render("✖ AI Suggestion Failed:") + "\n\n" +
		lipgloss.NewStyle().Foreground(styles.ColorText).Render(m.errorMessage) + "\n\n" +
		lipgloss.NewStyle().Foreground(styles.ColorTextMuted).Render("👉 Press [Enter] or [r] to edit your prompt • Press [c] to check API configuration")

	errCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorError).
		Background(styles.ColorCardBg).
		Padding(2, 3).
		Width(m.width - 6).
		Render(errContent)

	return "\n" + title + "\n\n" + errCard
}
