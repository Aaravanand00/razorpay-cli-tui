package state

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/config"
)

type ScreenType int

const (
	ScreenHome ScreenType = iota
	ScreenActions
	ScreenTable
	ScreenDetail
	ScreenForm
	ScreenConfig
	ScreenHelp
)

type Toast struct {
	ID      int
	Message string
	IsError bool
}

type ClearToastMsg struct {
	ToastID int
}

type ModuleItem struct {
	ID           string
	Title        string
	Description  string
	CommandCount int
}

type ActionItem struct {
	ID          string
	Title       string
	Description string
	CLICommand  string
	HTTPMethod  string
	APIPath     string
	IsForm      bool
}

type SessionState struct {
	CurrentScreen  ScreenType
	ScreenStack    []ScreenType
	SelectedModule ModuleItem
	SelectedAction ActionItem
	Breadcrumbs    []string
	Client         *api.Client

	// Dual Mode & Credentials
	ActiveMode    string // "test" or "live"
	TestKeyID     string
	TestKeySecret string
	LiveKeyID     string
	LiveKeySecret string

	// Active Credentials
	KeyID      string
	KeySecret  string
	IsLiveMode bool
	IsReadOnly bool

	Toast        *Toast
	toastCounter int
	Width        int
	Height       int
}

func NewSessionState() *SessionState {
	config.Init()

	testKeyID := config.TestKeyID()
	testKeySecret := config.TestKeySecret()
	liveKeyID := config.LiveKeyID()
	liveKeySecret := config.LiveKeySecret()

	// Legacy fallback
	legacyKeyID := config.KeyID()
	legacyKeySecret := config.KeySecret()

	if testKeyID == "" && strings.HasPrefix(legacyKeyID, "rzp_test_") {
		testKeyID = legacyKeyID
		testKeySecret = legacyKeySecret
	}
	if liveKeyID == "" && strings.HasPrefix(legacyKeyID, "rzp_live_") {
		liveKeyID = legacyKeyID
		liveKeySecret = legacyKeySecret
	}

	activeMode := config.ActiveMode()

	s := &SessionState{
		CurrentScreen: ScreenHome,
		ScreenStack:   []ScreenType{},
		Breadcrumbs:   []string{"Razorpay"},
		ActiveMode:    activeMode,
		TestKeyID:     testKeyID,
		TestKeySecret: testKeySecret,
		LiveKeyID:     liveKeyID,
		LiveKeySecret: liveKeySecret,
	}

	s.SyncActiveCredentials()
	return s
}

func (s *SessionState) SyncActiveCredentials() {
	if s.ActiveMode == "live" {
		s.KeyID = s.LiveKeyID
		s.KeySecret = s.LiveKeySecret
		s.IsLiveMode = true
	} else {
		s.KeyID = s.TestKeyID
		s.KeySecret = s.TestKeySecret
		s.IsLiveMode = false
	}

	if s.KeyID != "" && s.KeySecret != "" {
		s.Client = api.New(s.KeyID, s.KeySecret)
	} else {
		s.Client = nil
	}
}

func (s *SessionState) ToggleMode() string {
	if s.ActiveMode == "test" {
		if s.LiveKeyID == "" {
			s.SetToast("No Live API keys configured. Press 'c' to add Live keys.", true)
			return "No Live keys configured"
		}
		s.ActiveMode = "live"
		config.SetActiveMode("live")
		s.SyncActiveCredentials()
		msg := fmt.Sprintf("Switched to ● LIVE MODE (%s)", s.MaskedKey())
		s.SetToast(msg, false)
		return msg
	} else {
		if s.TestKeyID == "" {
			s.SetToast("No Test API keys configured. Press 'c' to add Test keys.", true)
			return "No Test keys configured"
		}
		s.ActiveMode = "test"
		config.SetActiveMode("test")
		s.SyncActiveCredentials()
		msg := fmt.Sprintf("Switched to ▲ TEST MODE (%s)", s.MaskedKey())
		s.SetToast(msg, false)
		return msg
	}
}

func (s *SessionState) MaskedKey() string {
	if s.KeyID == "" {
		return "No Key"
	}
	if len(s.KeyID) <= 12 {
		return s.KeyID
	}
	return s.KeyID[:9] + "..." + s.KeyID[len(s.KeyID)-4:]
}

func (s *SessionState) HasCredentials() bool {
	return s.KeyID != "" && s.KeySecret != ""
}

func (s *SessionState) PushScreen(next ScreenType, crumb string) {
	s.ScreenStack = append(s.ScreenStack, s.CurrentScreen)
	s.CurrentScreen = next
	if crumb != "" {
		s.Breadcrumbs = append(s.Breadcrumbs, crumb)
	}
}

func (s *SessionState) PopScreen() bool {
	if len(s.ScreenStack) == 0 {
		return false
	}
	lastIdx := len(s.ScreenStack) - 1
	s.CurrentScreen = s.ScreenStack[lastIdx]
	s.ScreenStack = s.ScreenStack[:lastIdx]

	if len(s.Breadcrumbs) > 1 {
		s.Breadcrumbs = s.Breadcrumbs[:len(s.Breadcrumbs)-1]
	}
	return true
}

func (s *SessionState) SetToast(msg string, isError bool) tea.Cmd {
	s.toastCounter++
	currentID := s.toastCounter
	s.Toast = &Toast{
		ID:      currentID,
		Message: msg,
		IsError: isError,
	}

	// 5-second automatic dismiss timer
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return ClearToastMsg{ToastID: currentID}
	})
}

func (s *SessionState) ClearToast() {
	s.Toast = nil
}

func (s *SessionState) HandleClearToast(id int) {
	if s.Toast != nil && s.Toast.ID == id {
		s.Toast = nil
	}
}
