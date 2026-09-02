package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/config"
	"github.com/razorpay/razorpay-cli/tui/components"
	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

type ConfigScreen struct {
	state        *state.SessionState
	activeMode   string // "test" or "live"
	testKeyInput textinput.Model
	testSecInput textinput.Model
	liveKeyInput textinput.Model
	liveSecInput textinput.Model
	focusIndex   int // 0: mode switch, 1: test key, 2: test secret, 3: live key, 4: live secret
	width        int
	height       int
}

func NewConfigScreen(s *state.SessionState, width, height int) ConfigScreen {
	tki := textinput.New()
	tki.Placeholder = "rzp_test_..."
	tki.SetValue(s.TestKeyID)
	tki.Prompt = " Key ID:     "
	tki.PromptStyle = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorWarning)
	tki.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorText)
	tki.Width = 38

	tsi := textinput.New()
	tsi.Placeholder = "Test Key Secret"
	tsi.SetValue(s.TestKeySecret)
	tsi.EchoMode = textinput.EchoPassword
	tsi.EchoCharacter = '•'
	tsi.Prompt = " Key Secret: "
	tsi.PromptStyle = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorWarning)
	tsi.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorText)
	tsi.Width = 38

	lki := textinput.New()
	lki.Placeholder = "rzp_live_..."
	lki.SetValue(s.LiveKeyID)
	lki.Prompt = " Key ID:     "
	lki.PromptStyle = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSuccess)
	lki.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorText)
	lki.Width = 38

	lsi := textinput.New()
	lsi.Placeholder = "Live Key Secret"
	lsi.SetValue(s.LiveKeySecret)
	lsi.EchoMode = textinput.EchoPassword
	lsi.EchoCharacter = '•'
	lsi.Prompt = " Key Secret: "
	lsi.PromptStyle = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSuccess)
	lsi.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorText)
	lsi.Width = 38

	// Default focus on test key
	tki.Focus()

	return ConfigScreen{
		state:        s,
		activeMode:   s.ActiveMode,
		testKeyInput: tki,
		testSecInput: tsi,
		liveKeyInput: lki,
		liveSecInput: lsi,
		focusIndex:   1,
		width:        width,
		height:       height,
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
		case "tab":
			c.setFocus((c.focusIndex + 1) % 5)
			return *c, nil
		case "shift+tab":
			c.setFocus((c.focusIndex - 1 + 5) % 5)
			return *c, nil
		case "down":
			if c.focusIndex == 0 {
				c.setFocus(1)
			} else if c.focusIndex == 1 {
				c.setFocus(2)
			} else if c.focusIndex == 2 {
				c.setFocus(0)
			} else if c.focusIndex == 3 {
				c.setFocus(4)
			} else if c.focusIndex == 4 {
				c.setFocus(0)
			}
			return *c, nil
		case "up":
			if c.focusIndex == 0 {
				c.setFocus(2)
			} else if c.focusIndex == 1 {
				c.setFocus(0)
			} else if c.focusIndex == 2 {
				c.setFocus(1)
			} else if c.focusIndex == 3 {
				c.setFocus(0)
			} else if c.focusIndex == 4 {
				c.setFocus(3)
			}
			return *c, nil
		case "right":
			// Jump from Test Box -> Live Box
			if c.focusIndex == 1 {
				c.setFocus(3)
				return *c, nil
			} else if c.focusIndex == 2 {
				c.setFocus(4)
				return *c, nil
			}
		case "left":
			// Jump from Live Box -> Test Box
			if c.focusIndex == 3 {
				c.setFocus(1)
				return *c, nil
			} else if c.focusIndex == 4 {
				c.setFocus(2)
				return *c, nil
			}
		case " ", "m":
			if c.focusIndex == 0 {
				if c.activeMode == "test" {
					c.activeMode = "live"
				} else {
					c.activeMode = "test"
				}
				return *c, nil
			}
		case "enter":
			testKey := strings.TrimSpace(c.testKeyInput.Value())
			testSec := strings.TrimSpace(c.testSecInput.Value())
			liveKey := strings.TrimSpace(c.liveKeyInput.Value())
			liveSec := strings.TrimSpace(c.liveSecInput.Value())

			err := config.SaveDualConfig(c.activeMode, testKey, testSec, liveKey, liveSec)
			if err != nil {
				c.state.SetToast("Failed to save config: "+err.Error(), true)
				return *c, nil
			}

			// Update state
			c.state.ActiveMode = c.activeMode
			c.state.TestKeyID = testKey
			c.state.TestKeySecret = testSec
			c.state.LiveKeyID = liveKey
			c.state.LiveKeySecret = liveSec
			c.state.SyncActiveCredentials()

			c.state.SetToast("Credentials & Active Profile saved successfully!", false)
			c.state.PopScreen()
			return *c, nil
		}
	}

	var cmd tea.Cmd
	switch c.focusIndex {
	case 1:
		c.testKeyInput, cmd = c.testKeyInput.Update(msg)
		cmds = append(cmds, cmd)
	case 2:
		c.testSecInput, cmd = c.testSecInput.Update(msg)
		cmds = append(cmds, cmd)
	case 3:
		c.liveKeyInput, cmd = c.liveKeyInput.Update(msg)
		cmds = append(cmds, cmd)
	case 4:
		c.liveSecInput, cmd = c.liveSecInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return *c, tea.Batch(cmds...)
}

func (c *ConfigScreen) setFocus(idx int) {
	c.focusIndex = idx
	c.testKeyInput.Blur()
	c.testSecInput.Blur()
	c.liveKeyInput.Blur()
	c.liveSecInput.Blur()

	switch idx {
	case 1:
		c.testKeyInput.Focus()
	case 2:
		c.testSecInput.Focus()
	case 3:
		c.liveKeyInput.Focus()
	case 4:
		c.liveSecInput.Focus()
	}
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

	// 3. Footer Keys
	keys := []components.KeyHelp{
		{Key: "Tab / Shift+Tab", Desc: "Next/Prev Field"},
		{Key: "← / →", Desc: "Switch Box"},
		{Key: "Space", Desc: "Toggle Mode"},
		{Key: "Enter", Desc: "Save All"},
		{Key: "Esc", Desc: "Back"},
		{Key: "Ctrl+C", Desc: "Quit"},
	}
	footerView := components.RenderFooter(c.state, c.width, keys)

	// Calculate Available Height
	headerHeight := lipgloss.Height(headerView)
	footerHeight := lipgloss.Height(footerView)
	toastHeight := 0
	if toastView != "" {
		toastHeight = lipgloss.Height(toastView)
	}

	bodyHeight := c.height - headerHeight - footerHeight - toastHeight
	if bodyHeight < 12 {
		bodyHeight = 12
	}

	// Active Mode Switcher Row
	var modeSwitch string
	if c.activeMode == "test" {
		testPill := styles.BadgeTestStyle.Render("● TEST MODE (Sandbox Active)")
		livePill := lipgloss.NewStyle().Foreground(styles.ColorTextDim).Render("○ LIVE MODE (Production)")
		modeSwitch = lipgloss.JoinHorizontal(lipgloss.Center, "   Active Environment:  ", testPill, "   ", livePill)
	} else {
		testPill := lipgloss.NewStyle().Foreground(styles.ColorTextDim).Render("○ TEST MODE (Sandbox)")
		livePill := styles.BadgeLiveStyle.Render("● LIVE MODE (Production Active)")
		modeSwitch = lipgloss.JoinHorizontal(lipgloss.Center, "   Active Environment:  ", testPill, "   ", livePill)
	}

	if c.focusIndex == 0 {
		modeSwitch = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorBorderFocus).
			Background(styles.ColorNavy).
			Padding(0, 1).
			Render("▶ " + modeSwitch + "  [Press Space to Toggle Active Mode]")
	} else {
		modeSwitch = lipgloss.NewStyle().
			Padding(0, 1).
			Render("  " + modeSwitch)
	}

	// Dynamic Box Widths
	boxWidth := (c.width - 8) / 2
	if boxWidth < 42 {
		boxWidth = 42
	}

	// Test Box Styling (Highlight when focused)
	testBorderColor := styles.ColorBorder
	testHeaderBadge := styles.BadgeTestStyle.Render("▲ TEST SANDBOX")
	if c.focusIndex == 1 || c.focusIndex == 2 {
		testBorderColor = styles.ColorWarning
		testHeaderBadge = styles.BadgeTestStyle.Render("▶ ▲ TEST SANDBOX (ACTIVE INPUT)")
	}

	testBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(testBorderColor).
		Padding(1, 2).
		Width(boxWidth)

	testContent := testHeaderBadge + "\n\n" +
		c.testKeyInput.View() + "\n\n" +
		c.testSecInput.View()

	// Live Box Styling (Highlight when focused)
	liveBorderColor := styles.ColorBorder
	liveHeaderBadge := styles.BadgeLiveStyle.Render("● LIVE PRODUCTION")
	if c.focusIndex == 3 || c.focusIndex == 4 {
		liveBorderColor = styles.ColorSuccess
		liveHeaderBadge = styles.BadgeLiveStyle.Render("▶ ● LIVE PRODUCTION (ACTIVE INPUT)")
	}

	liveBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(liveBorderColor).
		Padding(1, 2).
		Width(boxWidth)

	liveContent := liveHeaderBadge + "\n\n" +
		c.liveKeyInput.View() + "\n\n" +
		c.liveSecInput.View()

	boxesRow := lipgloss.JoinHorizontal(lipgloss.Top, testBoxStyle.Render(testContent), "  ", liveBoxStyle.Render(liveContent))

	tip := lipgloss.NewStyle().Foreground(styles.ColorMuted).Render("💡 Use 'Tab' or '← / →' arrows to jump between Test and Live boxes. Enter your keys once and press Enter to save!")

	title := styles.TitleStyle.Render("⚙️  Razorpay API Credentials & Dual Profile Manager")
	desc := styles.SubtitleStyle.Render("Manage Sandbox & Production keys. Enter your keys once — switch anytime in 1 sec!")

	rawBody := title + "\n" + desc + "\n\n" + modeSwitch + "\n\n" + boxesRow + "\n\n" + tip

	bodyContainer := lipgloss.NewStyle().
		Height(bodyHeight).
		Width(c.width - 2)

	sections = append(sections, bodyContainer.Render(rawBody))
	sections = append(sections, footerView)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
