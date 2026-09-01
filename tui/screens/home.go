package screens

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/tui/components"
	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

type moduleItem struct {
	module state.ModuleItem
}

func (i moduleItem) Title() string       { return i.module.Title }
func (i moduleItem) Description() string { return i.module.Description }
func (i moduleItem) FilterValue() string {
	return i.module.ID + " " + i.module.Title + " " + i.module.Description
}

type moduleDelegate struct{}

func (d moduleDelegate) Height() int                             { return 2 }
func (d moduleDelegate) Spacing() int                            { return 1 }
func (d moduleDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d moduleDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(moduleItem)
	if !ok {
		return
	}

	title := i.module.Title
	desc := i.module.Description
	countBadge := styles.BadgeCountStyle.Render(fmt.Sprintf("%d cmds", i.module.CommandCount))

	isSelected := index == m.Index()

	if isSelected {
		header := lipgloss.JoinHorizontal(
			lipgloss.Center,
			styles.ItemSelected.Render(fmt.Sprintf("▶ %-26s", title)),
			" ",
			countBadge,
		)
		descLine := styles.ItemDescSelected.Render(desc)
		fmt.Fprintf(w, "%s\n%s", header, descLine)
	} else {
		header := lipgloss.JoinHorizontal(
			lipgloss.Center,
			styles.ItemNormal.Render(fmt.Sprintf("  %-26s", title)),
			" ",
			countBadge,
		)
		descLine := styles.ItemDesc.Render(desc)
		fmt.Fprintf(w, "%s\n%s", header, descLine)
	}
}

type HomeScreen struct {
	list  list.Model
	state *state.SessionState
}

func GetModules() []state.ModuleItem {
	return []state.ModuleItem{
		{ID: "orders", Title: "🛒 Orders", Description: "Create and manage orders, checkout sessions and order payments", CommandCount: 5},
		{ID: "payments", Title: "💳 Payments", Description: "Fetch, capture, update payments and check method downtimes", CommandCount: 7},
		{ID: "refunds", Title: "💸 Refunds", Description: "Issue and track full/partial payment refunds and batches", CommandCount: 6},
		{ID: "customers", Title: "👥 Customers", Description: "Create, fetch, update customer records and tokens", CommandCount: 4},
		{ID: "invoices", Title: "🧾 Invoices", Description: "Manage invoices, issue drafts, line items and notifications", CommandCount: 13},
		{ID: "payment-links", Title: "🔗 Payment Links", Description: "Generate, notify, update and cancel standard payment links", CommandCount: 6},
		{ID: "qr-codes", Title: "📱 QR Codes", Description: "Create BharatQR codes, monitor payments and manage status", CommandCount: 6},
		{ID: "subscriptions", Title: "🔄 Subscriptions", Description: "Manage recurring plans, subscriptions, pauses and invoices", CommandCount: 14},
		{ID: "route", Title: "🔀 Route", Description: "Manage linked merchant accounts, transfers, and reversals", CommandCount: 20},
		{ID: "smart-collect", Title: "🏢 Smart Collect", Description: "Customer virtual bank accounts, UPI IDs, and TPV payers", CommandCount: 13},
		{ID: "settlements", Title: "🏦 Settlements", Description: "View settlement transfers, instant payouts, and recon reports", CommandCount: 6},
		{ID: "disputes", Title: "⚖️ Disputes", Description: "Track chargebacks, contest disputes, and upload evidence", CommandCount: 4},
		{ID: "documents", Title: "📄 Documents", Description: "Upload KYC/dispute documents and fetch file contents", CommandCount: 3},
		{ID: "configure", Title: "⚙️ Configure", Description: "Configure API credentials, Live/Test mode, and defaults", CommandCount: 1},
	}
}

func NewHomeScreen(s *state.SessionState, width, height int) HomeScreen {
	rawModules := GetModules()
	items := make([]list.Item, len(rawModules))
	for idx, m := range rawModules {
		items[idx] = moduleItem{module: m}
	}

	listHeight := height - 8
	if listHeight < 5 {
		listHeight = 5
	}

	l := list.New(items, moduleDelegate{}, width-4, listHeight)
	l.Title = "📦 Razorpay API Modules  (Select a module & press Enter)"
	l.Styles.Title = styles.TitleStyle
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(styles.ColorMuted)

	return HomeScreen{
		list:  l,
		state: s,
	}
}

func (h *HomeScreen) SetSize(width, height int) {
	listHeight := height - 8
	if listHeight < 5 {
		listHeight = 5
	}
	h.list.SetSize(width-4, listHeight)
}

func (h *HomeScreen) Update(msg tea.Msg) (HomeScreen, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if h.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "enter":
			if sel, ok := h.list.SelectedItem().(moduleItem); ok {
				h.state.SelectedModule = sel.module
				if sel.module.ID == "configure" {
					h.state.PushScreen(state.ScreenConfig, "⚙️ Configure")
				} else {
					h.state.PushScreen(state.ScreenActions, sel.module.Title)
				}
			}
		case "c":
			h.state.PushScreen(state.ScreenConfig, "⚙️ Configure")
		}
	}

	h.list, cmd = h.list.Update(msg)
	return *h, cmd
}

func (h HomeScreen) View() string {
	var sb strings.Builder

	// Header (ALWAYS Line 1)
	sb.WriteString(components.RenderHeader(h.state, h.state.Width))
	sb.WriteString("\n")

	// Toast (if any)
	toast := components.RenderToast(h.state)
	if toast != "" {
		sb.WriteString(toast + "\n")
	}

	// Module List
	sb.WriteString(h.list.View())
	sb.WriteString("\n")

	// Contextual Footer
	keys := []components.KeyHelp{
		{Key: "↑/↓", Desc: "Navigate"},
		{Key: "Enter", Desc: "Select Module"},
		{Key: "/", Desc: "Search"},
		{Key: "c", Desc: "Configure Keys"},
		{Key: "q", Desc: "Quit"},
	}
	sb.WriteString(components.RenderFooter(h.state, h.state.Width, keys))

	return sb.String()
}
