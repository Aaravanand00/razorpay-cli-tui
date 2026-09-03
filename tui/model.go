package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razorpay/razorpay-cli/tui/aiassist"
	"github.com/razorpay/razorpay-cli/tui/screens"
	"github.com/razorpay/razorpay-cli/tui/state"
)

type Model struct {
	state    *state.SessionState
	home     screens.HomeScreen
	actions  screens.ActionsScreen
	table    screens.TableViewScreen
	detail   screens.DetailViewScreen
	form     screens.FormViewScreen
	config   screens.ConfigScreen
	help     screens.HelpScreen
	aiAssist aiassist.Model
	ready    bool
}

func NewModel() Model {
	return NewModelWithOptions(false, "")
}

func NewModelWithOptions(readOnly bool, initialMode string) Model {
	sess := state.NewSessionState()
	sess.IsReadOnly = readOnly
	sess.ExplicitModeFlag = initialMode
	if initialMode == "live" || initialMode == "test" {
		sess.ActiveMode = initialMode
		sess.SyncActiveCredentials()
	}

	if readOnly {
		sess.SetToast("🛡️ Safe Read-Only Mode Active (Write operations locked)", false)
	} else if initialMode == "live" && !sess.HasCredentials() {
		sess.SetToast("● Live Production Mode (No keys configured). Press 'c' to add Live keys.", false)
	} else if initialMode == "test" && !sess.HasCredentials() {
		sess.SetToast("▲ Test Sandbox Mode (No keys configured). Press 'c' to add Test keys.", false)
	}

	return Model{
		state:    sess,
		home:     screens.NewHomeScreen(sess, 80, 24),
		config:   screens.NewConfigScreen(sess, 80, 24),
		help:     screens.NewHelpScreen(sess, 80, 24),
		aiAssist: aiassist.New(sess, 80, 24),
		ready:    false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case state.ClearToastMsg:
		m.state.HandleClearToast(msg.ToastID)
		return m, nil

	case tea.WindowSizeMsg:
		m.state.Width = msg.Width
		m.state.Height = msg.Height
		m.home.SetSize(msg.Width, msg.Height)
		m.actions.SetSize(msg.Width, msg.Height)
		m.config.SetSize(msg.Width, msg.Height)
		m.table.SetSize(msg.Width, msg.Height)
		m.detail.SetSize(msg.Width, msg.Height)
		m.form.SetSize(msg.Width, msg.Height)
		m.help.SetSize(msg.Width, msg.Height)
		m.aiAssist.SetSize(msg.Width, msg.Height)
		m.ready = true
		return m, nil

	case aiassist.ExecuteSuggestionMsg:
		if msg.Suggestion == nil {
			return m, nil
		}
		// Unified Action Execution Path: find matching module & action
		targetMod := strings.ToLower(strings.TrimSpace(msg.Suggestion.Resource))
		targetSub := strings.ToLower(strings.TrimSpace(msg.Suggestion.Subcommand))

		availableActions := screens.GetActionsForModule(targetMod)
		var matchedAction *state.ActionItem
		for _, a := range availableActions {
			if strings.Contains(strings.ToLower(a.CLICommand), targetSub) ||
				strings.Contains(strings.ToLower(a.ID), targetSub) {
				matchedAction = &a
				break
			}
		}

		// Match Module Item
		for _, mod := range screens.GetModules() {
			if strings.EqualFold(mod.ID, targetMod) {
				m.state.SelectedModule = mod
				break
			}
		}

		if matchedAction != nil {
			// Enforce existing safety rules (Read-Only Mode lock + Credentials validation)
			allowed, errMsg := m.state.CanPerformActionItem(*matchedAction)
			if !allowed {
				m.state.SetToast(errMsg, true)
				return m, nil
			}

			// Cleanly pop ScreenAIAssist and establish proper Actions hierarchy:
			// Stack: [ScreenHome] -> [ScreenActions] -> [ScreenTable / ScreenForm]
			// So pressing Esc returns to the Actions menu where the create/list options are!
			m.state.PopScreen()
			m.state.PushScreen(state.ScreenActions, m.state.SelectedModule.Title)
			m.actions = screens.NewActionsScreen(m.state, m.state.SelectedModule.ID, m.state.Width, m.state.Height)

			m.state.SelectedAction = *matchedAction
			if !matchedAction.IsForm {
				// Table List Action
				m.state.PushScreen(state.ScreenTable, matchedAction.Title)
				m.table = screens.NewTableViewScreen(m.state, *matchedAction, m.state.Width, m.state.Height)
				cmds = append(cmds, m.table.Init())
			} else {
				// Form / Mutation Action
				m.state.PushScreen(state.ScreenForm, matchedAction.Title)
				m.form = screens.NewFormViewScreen(m.state, *matchedAction, m.state.Width, m.state.Height)
				m.form.PopulateWithFlags(msg.Suggestion.Flags)
			}
			return m, tea.Batch(cmds...)
		}

	case tea.KeyMsg:
		// Global Help Toggle ('?' or 'h' key)
		if msg.String() == "?" || msg.String() == "h" {
			if m.state.CurrentScreen == state.ScreenHelp {
				m.state.PopScreen()
				return m, nil
			} else if m.state.CurrentScreen != state.ScreenForm &&
				m.state.CurrentScreen != state.ScreenConfig &&
				m.state.CurrentScreen != state.ScreenAIAssist &&
				!m.home.IsFiltering() &&
				!m.actions.IsFiltering() {
				m.state.PushScreen(state.ScreenHelp, "❓ Help")
				m.help = screens.NewHelpScreen(m.state, m.state.Width, m.state.Height)
				return m, nil
			}
		}

		// Global AI Assist Open ('a' or 'ctrl+a' key)
		if msg.String() == "a" || msg.String() == "ctrl+a" {
			if m.state.CurrentScreen != state.ScreenAIAssist &&
				m.state.CurrentScreen != state.ScreenForm &&
				m.state.CurrentScreen != state.ScreenConfig &&
				!m.home.IsFiltering() &&
				!m.actions.IsFiltering() {
				m.state.PushScreen(state.ScreenAIAssist, "🤖 AI Assist")
				m.aiAssist = aiassist.New(m.state, m.state.Width, m.state.Height)
				return m, nil
			}
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.state.CurrentScreen == state.ScreenHome ||
				m.state.CurrentScreen == state.ScreenActions ||
				m.state.CurrentScreen == state.ScreenTable ||
				m.state.CurrentScreen == state.ScreenDetail {
				return m, tea.Quit
			}
		case "esc":
			if m.state.CurrentScreen == state.ScreenHome {
				if m.home.IsFiltering() {
					m.home.ResetFilter()
					return m, nil
				}
			} else if m.state.CurrentScreen == state.ScreenActions {
				if m.actions.IsFiltering() {
					m.actions.ResetFilter()
					return m, nil
				}
				m.state.PopScreen()
				return m, nil
			} else {
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

		if m.state.CurrentScreen == state.ScreenActions && prevScreen == state.ScreenHome {
			m.actions = screens.NewActionsScreen(m.state, m.state.SelectedModule.ID, m.state.Width, m.state.Height)
		} else if m.state.CurrentScreen == state.ScreenConfig && prevScreen == state.ScreenHome {
			m.config = screens.NewConfigScreen(m.state, m.state.Width, m.state.Height)
		}

	case state.ScreenActions:
		var cmd tea.Cmd
		prevScreen := m.state.CurrentScreen
		m.actions, cmd = m.actions.Update(msg)
		cmds = append(cmds, cmd)

		if m.state.CurrentScreen == state.ScreenTable && prevScreen == state.ScreenActions {
			m.table = screens.NewTableViewScreen(m.state, m.state.SelectedAction, m.state.Width, m.state.Height)
			cmds = append(cmds, m.table.Init())
		} else if m.state.CurrentScreen == state.ScreenForm && prevScreen == state.ScreenActions {
			m.form = screens.NewFormViewScreen(m.state, m.state.SelectedAction, m.state.Width, m.state.Height)
		}

	case state.ScreenTable:
		var cmd tea.Cmd
		prevScreen := m.state.CurrentScreen
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)

		if m.state.CurrentScreen == state.ScreenDetail && prevScreen == state.ScreenTable {
			m.detail = screens.NewDetailViewScreen(m.state, m.state.Width, m.state.Height)
		}

	case state.ScreenDetail:
		var cmd tea.Cmd
		m.detail, cmd = m.detail.Update(msg)
		cmds = append(cmds, cmd)

	case state.ScreenForm:
		var cmd tea.Cmd
		prevScreen := m.state.CurrentScreen
		m.form, cmd = m.form.Update(msg)
		cmds = append(cmds, cmd)

		if m.state.CurrentScreen == state.ScreenDetail && prevScreen == state.ScreenForm {
			m.detail = screens.NewDetailViewScreen(m.state, m.state.Width, m.state.Height)
		}

	case state.ScreenConfig:
		var cmd tea.Cmd
		prevScreen := m.state.CurrentScreen
		m.config, cmd = m.config.Update(msg)
		cmds = append(cmds, cmd)

		if prevScreen == state.ScreenConfig && m.state.CurrentScreen == state.ScreenTable {
			cmds = append(cmds, m.table.Init())
		}

	case state.ScreenHelp:
		var cmd tea.Cmd
		m.help, cmd = m.help.Update(msg)
		cmds = append(cmds, cmd)

	case state.ScreenAIAssist:
		var cmd tea.Cmd
		m.aiAssist, cmd = m.aiAssist.Update(msg)
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
	case state.ScreenTable:
		return m.table.View()
	case state.ScreenDetail:
		return m.detail.View()
	case state.ScreenForm:
		return m.form.View()
	case state.ScreenConfig:
		return m.config.View()
	case state.ScreenHelp:
		return m.help.View()
	case state.ScreenAIAssist:
		return m.aiAssist.View()
	default:
		return m.home.View()
	}
}
