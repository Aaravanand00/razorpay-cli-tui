package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Start launches the interactive Razorpay Terminal UI
func Start() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
