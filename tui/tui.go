package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Start launches the interactive Razorpay Terminal UI
func Start(readOnly bool, initialMode string) error {
	p := tea.NewProgram(NewModelWithOptions(readOnly, initialMode), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
