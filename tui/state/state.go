package state

import (
	"strings"

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
	Message string
	IsError bool
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
	KeyID          string
	KeySecret      string
	IsLiveMode     bool
	Toast          *Toast
	Width          int
	Height         int
}

func NewSessionState() *SessionState {
	config.Init()
	keyID := config.KeyID()
	keySecret := config.KeySecret()

	isLive := strings.HasPrefix(keyID, "rzp_live_")

	var client *api.Client
	if keyID != "" && keySecret != "" {
		client = api.New(keyID, keySecret)
	}

	return &SessionState{
		CurrentScreen: ScreenHome,
		ScreenStack:   []ScreenType{},
		Breadcrumbs:   []string{"Razorpay"},
		Client:        client,
		KeyID:         keyID,
		KeySecret:     keySecret,
		IsLiveMode:    isLive,
	}
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

func (s *SessionState) SetToast(msg string, isError bool) {
	s.Toast = &Toast{
		Message: msg,
		IsError: isError,
	}
}

func (s *SessionState) ClearToast() {
	s.Toast = nil
}
