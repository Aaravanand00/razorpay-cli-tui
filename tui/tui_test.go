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

func TestActionsForModules(t *testing.T) {
	modules := screens.GetModules()
	totalActions := 0
	for _, mod := range modules {
		actions := screens.GetActionsForModule(mod.ID)
		if len(actions) == 0 {
			t.Fatalf("module %s has 0 actions", mod.ID)
		}
		totalActions += len(actions)
	}

	// Verify all 14 modules have their actions mapped
	if totalActions != 108 {
		t.Fatalf("expected 108 total actions across modules, got %d", totalActions)
	}
}

func TestNavigationToActionsAndBack(t *testing.T) {
	var model tea.Model = NewModel()
	model, _ = model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// Press Enter to go into first module (Orders)
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	m := model.(Model)
	if m.state.CurrentScreen != state.ScreenActions {
		t.Fatalf("expected ScreenActions, got %v", m.state.CurrentScreen)
	}

	// Press Esc to go back to Home
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	m = model.(Model)
	if m.state.CurrentScreen != state.ScreenHome {
		t.Fatalf("expected ScreenHome after Esc, got %v", m.state.CurrentScreen)
	}
}

func TestDualProfileAndModeToggle(t *testing.T) {
	m := NewModel()
	m.state.TestKeyID = "rzp_test_123456789"
	m.state.TestKeySecret = "test_sec_123"
	m.state.LiveKeyID = "rzp_live_987654321"
	m.state.LiveKeySecret = "live_sec_987"
	m.state.ActiveMode = "test"
	m.state.SyncActiveCredentials()

	if m.state.KeyID != "rzp_test_123456789" {
		t.Fatalf("expected test key active, got %s", m.state.KeyID)
	}

	// Toggle mode to Live via state
	m.state.ToggleMode()
	if m.state.ActiveMode != "live" {
		t.Fatalf("expected active mode 'live' after toggle, got %s", m.state.ActiveMode)
	}
	if m.state.KeyID != "rzp_live_987654321" {
		t.Fatalf("expected live key active, got %s", m.state.KeyID)
	}

	// Toggle back to Test
	m.state.ToggleMode()
	if m.state.ActiveMode != "test" {
		t.Fatalf("expected active mode 'test' after second toggle, got %s", m.state.ActiveMode)
	}
}

func TestReadOnlyOption(t *testing.T) {
	m := NewModelWithOptions(true, "test")
	if !m.state.IsReadOnly {
		t.Fatal("expected IsReadOnly to be true")
	}
	if m.state.ActiveMode != "test" {
		t.Fatalf("expected activeMode test, got %s", m.state.ActiveMode)
	}
}

func TestPermissionLockWithoutCredentials(t *testing.T) {
	sess := state.NewSessionState()
	sess.TestKeyID = ""
	sess.TestKeySecret = ""
	sess.LiveKeyID = ""
	sess.LiveKeySecret = ""
	sess.SyncActiveCredentials()

	// Test write permission lock
	allowed, errMsg := sess.CanPerformAction(true)
	if allowed {
		t.Fatal("expected write action to be blocked when no credentials exist")
	}
	if errMsg == "" {
		t.Fatal("expected non-empty error message")
	}
}

func TestTableViewInitializationAndColumns(t *testing.T) {
	sess := state.NewSessionState()
	action := state.ActionItem{
		ID:         "orders-list",
		Title:      "📋 List Orders",
		CLICommand: "razorpay orders list",
		HTTPMethod: "GET",
		APIPath:    "/v1/orders",
	}

	tv := screens.NewTableViewScreen(sess, action, 100, 30)
	view := tv.View()
	if view == "" {
		t.Fatal("expected non-empty table view")
	}

	// Simulate loaded items
	mockItems := []map[string]interface{}{
		{
			"id":         "order_DBJOWzybf0sJbb",
			"amount":     50000,
			"currency":   "INR",
			"status":     "paid",
			"receipt":    "Receipt #101",
			"created_at": 1672531199,
		},
	}

	tv, _ = tv.Update(screens.TableDataLoadedMsg{
		Items:   mockItems,
		RawJSON: `[{"id":"order_DBJOWzybf0sJbb"}]`,
	})

	updatedView := tv.View()
	if updatedView == "" {
		t.Fatal("expected non-empty updated table view")
	}
}
