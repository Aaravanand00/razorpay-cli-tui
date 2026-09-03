package aiassist

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/ai"
	"github.com/razorpay/razorpay-cli/config"
	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

type State int

const (
	StateInput State = iota
	StateLoading
	StatePreview
	StateError
)

type SuggestionMsg struct {
	Suggestion *ai.CommandSuggestion
	Err        error
}

type ExecuteSuggestionMsg struct {
	Suggestion *ai.CommandSuggestion
}

type Model struct {
	state        *state.SessionState
	client       *ai.Client
	input        textinput.Model
	spinner      spinner.Model
	currentState State
	suggestion   *ai.CommandSuggestion
	errorMessage string
	width        int
	height       int
}

func New(s *state.SessionState, width, height int) Model {
	ti := textinput.New()
	ti.Placeholder = "e.g. find failed payments from yesterday, or create order for ₹500"
	ti.Focus()
	ti.CharLimit = 250
	ti.Width = width - 12
	if ti.Width < 30 {
		ti.Width = 30
	}
	ti.Prompt = "🤖 Prompt : "
	ti.PromptStyle = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.ColorPrimary)

	apiKey := config.AIApiKey()
	aiClient := ai.NewClientWithKey(apiKey)

	return Model{
		state:        s,
		client:       aiClient,
		input:        ti,
		spinner:      sp,
		currentState: StateInput,
		width:        width,
		height:       height,
	}
}

func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.input.Width = w - 16
	if m.input.Width < 30 {
		m.input.Width = 30
	}
}

func (m *Model) Reset() {
	m.currentState = StateInput
	m.input.SetValue("")
	m.input.Focus()
	m.suggestion = nil
	m.errorMessage = ""
}

func (m *Model) CurrentState() State {
	return m.currentState
}

func (m *Model) Suggestion() *ai.CommandSuggestion {
	return m.suggestion
}
