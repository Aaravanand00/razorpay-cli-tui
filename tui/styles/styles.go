package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Razorpay Theme Colors
var (
	ColorPrimary     = lipgloss.Color("#0B72E7") // Razorpay Royal Blue
	ColorNavy        = lipgloss.Color("#0C2340") // Dark Navy
	ColorSecondary   = lipgloss.Color("#53B5FD") // Sky Blue
	ColorBackground  = lipgloss.Color("#0F172A") // Slate Background
	ColorCardBg      = lipgloss.Color("#1E293B") // Dark Card Slate
	ColorActiveCard  = lipgloss.Color("#1E3A8A") // Highlighted Card Navy
	ColorText        = lipgloss.Color("#F8FAFC") // Off-white Text
	ColorMuted       = lipgloss.Color("#94A3B8") // Gray Subtext
	ColorBorder      = lipgloss.Color("#334155") // Subtle Border
	ColorBorderFocus = lipgloss.Color("#3B82F6") // Focused Border Blue
	ColorSuccess     = lipgloss.Color("#10B981") // Green
	ColorWarning     = lipgloss.Color("#F59E0B") // Amber / Yellow
	ColorError       = lipgloss.Color("#EF4444") // Red
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
			Foreground(ColorPrimary)

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

	// Header & Breadcrumbs
	HeaderContainer = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(ColorBorder).
			Padding(0, 1).
			MarginBottom(1)

	BreadcrumbStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary)

	// Items
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

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2).
			Margin(0, 1)

	CardActiveStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorderFocus).
			Padding(1, 2).
			Margin(0, 1)

	// Footer & Keybindings
	FooterContainer = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorBorder).
			Padding(0, 1).
			MarginTop(1)

	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	KeyDescStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Toasts & Notifications
	ToastSuccessStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(ColorSuccess).
				Padding(0, 2).
				MarginBottom(1)

	ToastErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorError).
			Padding(0, 2).
			MarginBottom(1)
)
