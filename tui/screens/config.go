package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/config"
	"github.com/razorpay/razorpay-cli/tui/components"
	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

type ConfigScreen struct {
	state       *state.SessionState
	keyInput    textinput.Model
	secretInput textinput.Model
	focusIndex  int
	width       int
	height      int
}

func NewConfigScreen(s *state.SessionState, width, height int) ConfigScreen {
	ki := textinput.New()
	ki.Placeholder = "rzp_test_... or rzp_live_..."
	ki.SetValue(s.KeyID)
	ki.Focus()
	ki.Prompt = " Key ID:     "
	ki.PromptStyle = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	ki.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorText)
	ki.Width = 50

	si := textinput.New()
	si.Placeholder = "Your API Key Secret"
	si.SetValue(s.KeySecret)
	si.EchoMode = textinput.EchoPassword
	si.EchoCharacter = '•'
	si.Prompt = " Key Secret: "
	si.PromptStyle = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	si.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorText)
	si.Width = 50

	return ConfigScreen{
		state:       s,
		keyInput:    ki,
		secretInput: si,
		focusIndex:  0,
		width:       width,
		height:      height,
	}
}

func (c *ConfigScreen) SetSize(width, height int) {
	c.width = width
	c.height = height
}

func (c *ConfigScreen) Update(msg tea.Msg) (ConfigScreen, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			c.focusIndex = (c.focusIndex + 1) % 2
			if c.focusIndex == 0 {
				c.keyInput.Focus()
				c.secretInput.Blur()
			} else {
				c.keyInput.Blur()
				c.secretInput.Focus()
			}
		case "shift+tab", "up":
			c.focusIndex = (c.focusIndex - 1 + 2) % 2
			if c.focusIndex == 0 {
				c.keyInput.Focus()
				c.secretInput.Blur()
			} else {
				c.keyInput.Blur()
				c.secretInput.Focus()
			}
		case "enter":
			keyID := strings.TrimSpace(c.keyInput.Value())
			keySecret := strings.TrimSpace(c.secretInput.Value())

			if keyID == "" || keySecret == "" {
				c.state.SetToast("Key ID and Key Secret cannot be empty", true)
				return *c, nil
			}

			err := config.Save(keyID, keySecret)
			if err != nil {
				c.state.SetToast("Failed to save config: "+err.Error(), true)
				return *c, nil
			}

			// Update runtime session state
			c.state.KeyID = keyID
			c.state.KeySecret = keySecret
			c.state.IsLiveMode = strings.HasPrefix(keyID, "rzp_live_")
			c.state.Client = api.New(keyID, keySecret)
			c.state.SetToast("Credentials saved successfully!", false)

			// Pop back to previous screen
			c.state.PopScreen()
			return *c, nil
		}
	}

	var cmd tea.Cmd
	if c.focusIndex == 0 {
		c.keyInput, cmd = c.keyInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		c.secretInput, cmd = c.secretInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return *c, tea.Batch(cmds...)
}

func (c ConfigScreen) View() string {
	var sections []string

	// 1. Header
	headerView := components.RenderHeader(c.state, c.width)
	sections = append(sections, headerView)

	// 2. Toast
	toastView := components.RenderToast(c.state)
	if toastView != "" {
		sections = append(sections, toastView)
	}

	// 3. Footer Keys (Using Ctrl+C for Quit because q is a typeable character in inputs)
	keys := []components.KeyHelp{
		{Key: "Tab", Desc: "Switch Field"},
		{Key: "Enter", Desc: "Save Credentials"},
		{Key: "Esc", Desc: "Cancel / Back"},
		{Key: "Ctrl+C", Desc: "Quit"},
	}
	footerView := components.RenderFooter(c.state, c.width, keys)

	// Calculate Available Height to pin footer strictly to the bottom
	headerHeight := lipgloss.Height(headerView)
	footerHeight := lipgloss.Height(footerView)
	toastHeight := 0
	if toastView != "" {
		toastHeight = lipgloss.Height(toastView)
	}

	bodyHeight := c.height - headerHeight - footerHeight - toastHeight
	if bodyHeight < 8 {
		bodyHeight = 8
	}

	// Form Body Elements
	title := styles.TitleStyle.Render("⚙️  Configure Razorpay API Credentials")
	desc := styles.SubtitleStyle.Render("Enter your Razorpay Key ID and Secret. Saved securely to ~/.razorpay/config.yaml")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorderFocus).
		Padding(1, 3).
		Width(c.width - 6)

	var formContent strings.Builder
	formContent.WriteString(c.keyInput.View() + "\n\n")
	formContent.WriteString(c.secretInput.View() + "\n\n")

	hint := lipgloss.NewStyle().Foreground(styles.ColorMuted).Render("💡 Tip: Use 'rzp_test_...' for Test Mode or 'rzp_live_...' for Live Mode.")
	formContent.WriteString(hint)

	rawBody := title + "\n" + desc + "\n\n" + box.Render(formContent.String())

	// Container that stretches body to full available height, pushing footer to bottom
	bodyContainer := lipgloss.NewStyle().
		Height(bodyHeight).
		Width(c.width - 2)

	sections = append(sections, bodyContainer.Render(rawBody))

	// 4. Footer at bottom
	sections = append(sections, footerView)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
