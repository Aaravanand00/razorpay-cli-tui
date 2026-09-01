package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razorpay/razorpay-cli/tui/screens"
	"github.com/razorpay/razorpay-cli/tui/state"
)

type Model struct {
	state   *state.SessionState
	home    screens.HomeScreen
	actions screens.ActionsScreen
	ready   bool
}

func NewModel() Model {
	sess := state.NewSessionState()
	return Model{
		state: sess,
		home:  screens.NewHomeScreen(sess, 80, 24),
		ready: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.state.Width = msg.Width
		m.state.Height = msg.Height
		m.home.SetSize(msg.Width, msg.Height)
		m.actions.SetSize(msg.Width, msg.Height)
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			// Allow quit from Home and Actions screen directly
			if m.state.CurrentScreen == state.ScreenHome || m.state.CurrentScreen == state.ScreenActions {
				return m, tea.Quit
			}
		case "esc":
			if m.state.CurrentScreen != state.ScreenHome {
				m.state.PopScreen()
				return m, nil
			}
		}
	}

	// Screen Router
	switch m.state.CurrentScreen {
	case state.ScreenHome:
		var cmd tea.Cmd
		prevScreen := m.state.CurrentScreen
		m.home, cmd = m.home.Update(msg)
		cmds = append(cmds, cmd)

		// If user selected a module, initialize the actions screen for that module
		if m.state.CurrentScreen == state.ScreenActions && prevScreen == state.ScreenHome {
			m.actions = screens.NewActionsScreen(m.state, m.state.SelectedModule.ID, m.state.Width, m.state.Height)
		}

	case state.ScreenActions:
		var cmd tea.Cmd
		m.actions, cmd = m.actions.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing Razorpay Terminal UI..."
	}

	switch m.state.CurrentScreen {
	case state.ScreenHome:
		return m.home.View()
	case state.ScreenActions:
		return m.actions.View()
	default:
		return m.home.View()
	}
}
