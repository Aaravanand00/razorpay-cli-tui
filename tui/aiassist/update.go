package aiassist

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razorpay/razorpay-cli/ai"
	"github.com/razorpay/razorpay-cli/config"
	"github.com/razorpay/razorpay-cli/tui/state"
)

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	m.SetSize(m.state.Width, m.state.Height)

	// Sync API key in case user updated it in Config screen
	apiKey := config.AIApiKey()
	m.client = ai.NewClientWithKey(apiKey)

	switch msg := msg.(type) {
	case SuggestionMsg:
		if msg.Err != nil {
			m.currentState = StateError
			m.errorMessage = msg.Err.Error()
		} else {
			m.currentState = StatePreview
			m.suggestion = msg.Suggestion
		}
		return m, nil

	case spinner.TickMsg:
		if m.currentState == StateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.KeyMsg:
		// When AI API Key is not configured
		if apiKey == "" {
			switch msg.String() {
			case "c":
				m.state.PushScreen(state.ScreenConfig, "⚙️ Configure")
				return m, nil
			case "esc":
				m.state.PopScreen()
				return m, nil
			}
			return m, nil
		}

		switch m.currentState {
		case StateInput:
			switch msg.String() {
			case "enter":
				val := strings.TrimSpace(m.input.Value())
				if val == "" {
					m.state.SetToast("⚠️ Please enter a prompt for the AI Assistant", true)
					return m, nil
				}
				m.currentState = StateLoading
				return m, tea.Batch(m.spinner.Tick, fetchSuggestionCmd(m.client, val, m.state.ActiveMode))
			case "esc":
				m.state.PopScreen()
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}

			var inputCmd tea.Cmd
			m.input, inputCmd = m.input.Update(msg)
			cmds = append(cmds, inputCmd)

		case StateLoading:
			switch msg.String() {
			case "esc":
				m.currentState = StateInput
				m.input.Focus()
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}

		case StatePreview:
			switch msg.String() {
			case "enter":
				// Forward execution message to main TUI router to execute with existing safety checks
				return m, func() tea.Msg {
					return ExecuteSuggestionMsg{Suggestion: m.suggestion}
				}
			case "esc", "r", "e":
				m.currentState = StateInput
				m.input.Focus()
				return m, nil
			case "c":
				m.state.PushScreen(state.ScreenConfig, "⚙️ Configure")
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}

		case StateError:
			switch msg.String() {
			case "enter", "r", "esc":
				m.currentState = StateInput
				m.input.Focus()
				return m, nil
			case "c":
				m.state.PushScreen(state.ScreenConfig, "⚙️ Configure")
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func fetchSuggestionCmd(client *ai.Client, userInput, activeMode string) tea.Cmd {
	return func() tea.Msg {
		context := "Active environment: " + activeMode
		sugg, err := client.GetSuggestion(userInput, context)
		return SuggestionMsg{
			Suggestion: sugg,
			Err:        err,
		}
	}
}
