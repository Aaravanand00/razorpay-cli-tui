package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razorpay/razorpay-cli/tui/screens"
	"github.com/razorpay/razorpay-cli/tui/state"
)

func TestModelInitialization(t *testing.T) {
	m := NewModel()
	if m.state == nil {
		t.Fatal("expected non-nil session state")
	}
	if m.state.CurrentScreen != state.ScreenHome {
		t.Fatalf("expected initial screen to be ScreenHome, got %v", m.state.CurrentScreen)
	}

	modules := screens.GetModules()
	if len(modules) != 14 {
		t.Fatalf("expected 14 modules, got %d", len(modules))
	}
}

func TestModelWindowResize(t *testing.T) {
	m := NewModel()
	updatedModel, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	
	mod, ok := updatedModel.(Model)
	if !ok {
		t.Fatal("expected model type Model")
	}
	if !mod.ready {
		t.Fatal("expected model to be ready after WindowSizeMsg")
	}
	if mod.state.Width != 100 || mod.state.Height != 30 {
		t.Fatalf("expected 100x30 dimensions, got %dx%d", mod.state.Width, mod.state.Height)
	}

	viewOutput := mod.View()
	if viewOutput == "" {
		t.Fatal("expected non-empty view output")
	}
}
