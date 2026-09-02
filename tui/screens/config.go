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
	tki.Width = 36

	tsi := textinput.New()
	tsi.Placeholder = "Test Key Secret"
	tsi.SetValue(s.TestKeySecret)
	tsi.EchoMode = textinput.EchoPassword
	tsi.EchoCharacter = '•'
	tsi.Prompt = " Key Secret: "
	tsi.PromptStyle = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorWarning)
	tsi.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorText)
	tsi.Width = 36

	lki := textinput.New()
	lki.Placeholder = "rzp_live_..."
	lki.SetValue(s.LiveKeyID)
	lki.Prompt = " Key ID:     "
	lki.PromptStyle = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSuccess)
	lki.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorText)
	lki.Width = 36

	lsi := textinput.New()
	lsi.Placeholder = "Live Key Secret"
	lsi.SetValue(s.LiveKeySecret)
	lsi.EchoMode = textinput.EchoPassword
	lsi.EchoCharacter = '•'
	lsi.Prompt = " Key Secret: "
	lsi.PromptStyle = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSuccess)
	lsi.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorText)
	lsi.Width = 36

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
				c.setFocus(3)
			} else if c.focusIndex == 3 {
				c.setFocus(4)
			} else if c.focusIndex == 4 {
				c.setFocus(0)
			}
			return *c, nil
		case "up":
			if c.focusIndex == 0 {
				c.setFocus(4)
			} else if c.focusIndex == 1 {
				c.setFocus(0)
			} else if c.focusIndex == 2 {
				c.setFocus(1)
			} else if c.focusIndex == 3 {
				c.setFocus(2)
			} else if c.focusIndex == 4 {
				c.setFocus(3)
			}
			return *c, nil
		case "left":
			if c.focusIndex == 0 {
				c.activeMode = "test"
				c.state.ClearToast()
				return *c, nil
			}
		case "right":
			if c.focusIndex == 0 {
				c.activeMode = "live"
				c.state.ClearToast()
				return *c, nil
			}
		case " ", "m":
			if c.focusIndex == 0 {
				if c.activeMode == "test" {
					c.activeMode = "live"
				} else {
					c.activeMode = "test"
				}
				c.state.ClearToast()
				return *c, nil
			}
		case "1", "t":
			if c.focusIndex == 0 {
				c.activeMode = "test"
				c.state.ClearToast()
				return *c, nil
			}
		case "2", "l":
			if c.focusIndex == 0 {
				c.activeMode = "live"
				c.state.ClearToast()
				return *c, nil
			}
		case "enter":
			testKey := strings.TrimSpace(c.testKeyInput.Value())
			testSec := strings.TrimSpace(c.testSecInput.Value())
			liveKey := strings.TrimSpace(c.liveKeyInput.Value())
			liveSec := strings.TrimSpace(c.liveSecInput.Value())

			// 1. Check if completely empty
			if testKey == "" && testSec == "" && liveKey == "" && liveSec == "" {
				cmd := c.state.SetToast("Please enter at least one valid Razorpay API Key ID and Secret", true)
				return *c, cmd
			}

			// 2. Strict Focus-First Validation
			var validationErr string

			if c.focusIndex == 1 || c.focusIndex == 2 {
				// User is explicitly inside Test Sandbox box
				if testKey != "" && !strings.HasPrefix(testKey, "rzp_test_") {
					validationErr = "Invalid Test Key ID format: Must start with 'rzp_test_...'"
				} else if testKey != "" && testSec == "" {
					validationErr = "Please enter the Test Key Secret for your Test Key ID"
				} else if liveKey != "" && !strings.HasPrefix(liveKey, "rzp_live_") {
					validationErr = "Invalid Live Key ID format: Must start with 'rzp_live_...'"
				} else if liveKey != "" && liveSec == "" {
					validationErr = "Please enter the Live Key Secret for your Live Key ID"
				}
			} else if c.focusIndex == 3 || c.focusIndex == 4 {
				// User is explicitly inside Live Production box
				if liveKey != "" && !strings.HasPrefix(liveKey, "rzp_live_") {
					validationErr = "Invalid Live Key ID format: Must start with 'rzp_live_...'"
				} else if liveKey != "" && liveSec == "" {
					validationErr = "Please enter the Live Key Secret for your Live Key ID"
				} else if testKey != "" && !strings.HasPrefix(testKey, "rzp_test_") {
					validationErr = "Invalid Test Key ID format: Must start with 'rzp_test_...'"
				} else if testKey != "" && testSec == "" {
					validationErr = "Please enter the Test Key Secret for your Test Key ID"
				}
			}

			// 3. Active Mode Lock Validation
			if validationErr == "" {
				if c.activeMode == "live" && (liveKey == "" || liveSec == "") {
					if testKey != "" {
						validationErr = "Active mode is set to LIVE, but Live keys are empty. Switch Active Environment to Test or enter Live keys."
					} else {
						validationErr = "Please enter valid Live Key ID (rzp_live_...) and Secret."
					}
				} else if c.activeMode == "test" && (testKey == "" || testSec == "") {
					if liveKey != "" {
						validationErr = "Active mode is set to TEST, but Test keys are empty. Switch Active Environment to Live or enter Test keys."
					} else {
						validationErr = "Please enter valid Test Key ID (rzp_test_...) and Secret."
					}
				}
			}

			if validationErr != "" {
				cmd := c.state.SetToast(validationErr, true)
				return *c, cmd
			}

			// 4. Save to config
			err := config.SaveDualConfig(c.activeMode, testKey, testSec, liveKey, liveSec)
			if err != nil {
				cmd := c.state.SetToast("Failed to save config: "+err.Error(), true)
				return *c, cmd
			}

			// 5. Update state
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
	c.state.ClearToast() // Instantly clears previous box warning when switching fields
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

	// 1. Top Header
	headerView := components.RenderHeader(c.state, c.width)
	sections = append(sections, headerView)

	// 2. Toast Alert (Auto-clears after 5s)
	toastView := components.RenderToast(c.state)
	if toastView != "" {
		sections = append(sections, toastView)
	}

	// 3. Footer Keys
	keys := []components.KeyHelp{
		{Key: "Tab / Shift+Tab", Desc: "Switch Field"},
		{Key: "← / →", Desc: "Select Mode"},
		{Key: "Enter", Desc: "Save All"},
		{Key: "?/h", Desc: "Help"},
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
	if bodyHeight < 14 {
		bodyHeight = 14
	}

	// Active Mode Switcher Row with Spacious Padding
	var testPill, livePill string
	if c.activeMode == "test" {
		testPill = styles.BadgeTestStyle.Render("● 1. TEST SANDBOX (ACTIVE)")
		livePill = lipgloss.NewStyle().Foreground(styles.ColorTextDim).Background(styles.ColorNavy).Padding(0, 1).Render("○ 2. LIVE PRODUCTION")
	} else {
		testPill = lipgloss.NewStyle().Foreground(styles.ColorTextDim).Background(styles.ColorNavy).Padding(0, 1).Render("○ 1. TEST SANDBOX")
		livePill = styles.BadgeLiveStyle.Render("● 2. LIVE PRODUCTION (ACTIVE)")
	}

	modeSwitchContent := lipgloss.JoinHorizontal(lipgloss.Center,
		"  Active Environment:  ",
		testPill,
		"     ",
		livePill,
	)

	var modeSwitch string
	if c.focusIndex == 0 {
		modeSwitch = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorBorderFocus).
			Background(styles.ColorCardBg).
			Padding(1, 2).
			Render("▶ " + modeSwitchContent + "    [Use ← / → Arrows or 1/2 to Switch]")
	} else {
		modeSwitch = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorBorder).
			Padding(1, 2).
			Render("  " + modeSwitchContent)
	}

	// Dynamic Box Widths
	boxWidth := (c.width - 12) / 2
	if boxWidth < 42 {
		boxWidth = 42
	}

	// Test Box Styling (Spacious with 1, 3 padding)
	testBorderColor := styles.ColorBorder
	testHeaderBadge := styles.BadgeTestStyle.Render("▲ TEST SANDBOX")
	if c.focusIndex == 1 || c.focusIndex == 2 {
		testBorderColor = styles.ColorWarning
		testHeaderBadge = styles.BadgeTestStyle.Render("▶ ▲ TEST SANDBOX (TYPING HERE)")
	}

	testBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(testBorderColor).
		Padding(1, 3).
		Width(boxWidth)

	testContent := testHeaderBadge + "\n\n" +
		c.testKeyInput.View() + "\n\n" +
		c.testSecInput.View()

	// Live Box Styling (Spacious with 1, 3 padding)
	liveBorderColor := styles.ColorBorder
	liveHeaderBadge := styles.BadgeLiveStyle.Render("● LIVE PRODUCTION")
	if c.focusIndex == 3 || c.focusIndex == 4 {
		liveBorderColor = styles.ColorSuccess
		liveHeaderBadge = styles.BadgeLiveStyle.Render("▶ ● LIVE PRODUCTION (TYPING HERE)")
	}

	liveBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(liveBorderColor).
		Padding(1, 3).
		Width(boxWidth)

	liveContent := liveHeaderBadge + "\n\n" +
		c.liveKeyInput.View() + "\n\n" +
		c.liveSecInput.View()

	boxesRow := lipgloss.JoinHorizontal(lipgloss.Top, testBoxStyle.Render(testContent), "    ", liveBoxStyle.Render(liveContent))

	tip := lipgloss.NewStyle().Foreground(styles.ColorMuted).PaddingLeft(1).Render("💡 Format: Test Key starts with 'rzp_test_...', Live Key with 'rzp_live_...'. Press 'Tab' to move.")

	title := styles.TitleStyle.PaddingLeft(1).Render("⚙️  Razorpay API Credentials & Dual Profile Manager")
	desc := styles.SubtitleStyle.PaddingLeft(1).Render("Manage Sandbox & Production keys. Saved securely to ~/.razorpay/config.yaml")

	rawBody := "\n" + title + "\n" + desc + "\n\n" + modeSwitch + "\n\n" + boxesRow + "\n\n" + tip

	bodyContainer := lipgloss.NewStyle().
		Height(bodyHeight).
		Width(c.width - 2)

	sections = append(sections, bodyContainer.Render(rawBody))
	sections = append(sections, footerView)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
