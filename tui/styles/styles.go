package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Razorpay Theme Colors
var (
	ColorPrimary     = lipgloss.Color("#0B72E7") // Razorpay Royal Blue
	ColorNavy        = lipgloss.Color("#0C2340") // Dark Navy
	ColorSecondary   = lipgloss.Color("#53B5FD") // Sky Blue
	ColorBackground  = lipgloss.Color("#0B0F19") // Deep Black Slate
	ColorCardBg      = lipgloss.Color("#131D2F") // Elevated Card Slate
	ColorText        = lipgloss.Color("#F8FAFC") // Off-white Text
	ColorMuted       = lipgloss.Color("#64748B") // Slate Subtext
	ColorMutedLight  = lipgloss.Color("#94A3B8") // Gray Light
	ColorBorder      = lipgloss.Color("#1E293B") // Subtle Border
	ColorBorderFocus = lipgloss.Color("#3B82F6") // Focused Blue
	ColorSuccess     = lipgloss.Color("#10B981") // Green
	ColorWarning     = lipgloss.Color("#F59E0B") // Amber
	ColorError       = lipgloss.Color("#EF4444") // Red
	ColorAccentTag   = lipgloss.Color("#1D4ED8") // Tag Blue
)

// Lipgloss Styles
var (
	// Brand & Logo
	BrandStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPrimary).
			Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Mode Badges
	BadgeTestStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(ColorWarning).
			Padding(0, 1)

	BadgeLiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorSuccess).
			Padding(0, 1)

	BadgeNoAuthStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(ColorError).
				Padding(0, 1)

	BadgeCountStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			Background(ColorNavy).
			Padding(0, 1)

	// Header & Breadcrumbs (Height = 2 lines total including border)
	HeaderContainer = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	BreadcrumbStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary)

	// Items (Spacious & Clean)
	ItemNormal = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(ColorText)

	ItemSelected = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPrimary).
			Padding(0, 2)

	ItemDesc = lipgloss.NewStyle().
			Foreground(ColorMuted).
			PaddingLeft(2)

	ItemDescSelected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E0F2FE")).
				Background(ColorPrimary).
				PaddingLeft(2)

	// Footer & Keybindings (Height = 2 lines total including border)
	FooterContainer = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	KeyDescStyle = lipgloss.NewStyle().
			Foreground(ColorMutedLight)

	// Toasts & Notifications
	ToastSuccessStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(ColorSuccess).
				Padding(0, 2)

	ToastErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorError).
			Padding(0, 2)
)
