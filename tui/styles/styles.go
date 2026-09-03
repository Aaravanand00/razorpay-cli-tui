package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Razorpay Theme Colors
var (
	ColorPrimary       = lipgloss.Color("#0B72E7") // Razorpay Royal Blue
	ColorNavy          = lipgloss.Color("#0C2340") // Dark Navy
	ColorSecondary     = lipgloss.Color("#53B5FD") // Sky Blue
	ColorBackground    = lipgloss.Color("#080D1A") // Deep Premium Dark
	ColorCardBg        = lipgloss.Color("#111C30") // Elevated Card Slate
	ColorCardBgHover   = lipgloss.Color("#1E3A8A") // Highlighted Blue
	ColorText          = lipgloss.Color("#FFFFFF") // Pure White
	ColorTextMuted     = lipgloss.Color("#94A3B8") // Gray Subtext
	ColorTextDim       = lipgloss.Color("#64748B") // Dark Gray
	ColorBorder        = lipgloss.Color("#1E293B") // Border Gray
	ColorBorderFocus   = lipgloss.Color("#3B82F6") // Focused Blue
	ColorSuccess       = lipgloss.Color("#10B981") // Green
	ColorWarning       = lipgloss.Color("#F59E0B") // Amber
	ColorError         = lipgloss.Color("#EF4444") // Red
	ColorMuted         = ColorTextMuted            // Backward compatibility alias
)

// Lipgloss Styles
var (
	// Brand & Logo
	BrandStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			Background(ColorNavy).
			Padding(0, 1)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorTextMuted)

	// Step Navigation Badges
	StepActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPrimary).
			Padding(0, 1)

	StepDoneStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Background(ColorNavy).
			Padding(0, 1)

	StepInactiveStyle = lipgloss.NewStyle().
				Foreground(ColorTextDim).
				Background(ColorCardBg).
				Padding(0, 1)

	StepArrow = lipgloss.NewStyle().
			Foreground(ColorTextDim).
			Padding(0, 1)

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

	BadgeReadOnlyStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#0284C7")).
				Padding(0, 1)

	BadgeCountStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			Background(ColorNavy).
			Padding(0, 1)

	BadgeCountSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(ColorNavy).
				Padding(0, 1)

	// Header & Breadcrumbs
	HeaderContainer = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	LocationBar = lipgloss.NewStyle().
			Foreground(ColorTextMuted).
			Padding(0, 1)

	LocationActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary)

	// Spacious Cards for Items
	ItemCardNormal = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(ColorText)

	ItemCardSelected = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPrimary).
			Padding(0, 2)

	ItemDescNormal = lipgloss.NewStyle().
			Foreground(ColorTextMuted).
			PaddingLeft(4)

	ItemDescSelected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E0F2FE")).
				Background(ColorPrimary).
				PaddingLeft(4)

	// Footer & Keybindings
	FooterContainer = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	KeyDescStyle = lipgloss.NewStyle().
			Foreground(ColorTextMuted)

	// Toasts & Notifications
	ToastSuccessStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(ColorSuccess).
				Padding(0, 2)

	ToastWarningStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#000000")).
				Background(ColorWarning).
				Padding(0, 2)

	ToastErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorError).
			Padding(0, 2)

	// Table Styles
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorSecondary).
				Background(ColorNavy).
				Padding(0, 1)

	TableCellSelectedStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(ColorPrimary)

	TableContainerStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorBorder).
				Padding(0, 1)

	// Status Pills
	StatusPillPaid = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#10B981")).
			Render("● paid")

	StatusPillCaptured = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#10B981")).
				Render("● captured")

	StatusPillCreated = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#F59E0B")).
				Render("▲ created")

	StatusPillActive = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#10B981")).
				Render("● active")

	StatusPillFailed = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#EF4444")).
				Render("✖ failed")

	// Detail & JSON Inspector Styles
	TabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorPrimary).
			Padding(0, 2)

	TabInactiveStyle = lipgloss.NewStyle().
				Foreground(ColorTextMuted).
				Background(ColorNavy).
				Padding(0, 2)

	DetailKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary).
			Width(22)

	DetailValStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	JSONKeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#38BDF8"))

	JSONStringStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#34D399"))

	JSONNumberStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FBBF24"))

	JSONBoolStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F472B6"))

	JSONNullStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(ColorTextDim)
)
