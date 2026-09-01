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

type item struct {
	module state.ModuleItem
}

func (i item) Title() string       { return i.module.Title }
func (i item) Description() string { return i.module.Description }
func (i item) FilterValue() string { return i.module.ID + " " + i.module.Title + " " + i.module.Description }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	title := i.module.Title
	desc := i.module.Description
	countBadge := fmt.Sprintf("[%d cmds]", i.module.CommandCount)

	isSelected := index == m.Index()

	if isSelected {
		titleLine := styles.ItemSelected.Render(fmt.Sprintf("▶ %-26s %s", title, countBadge))
		descLine := styles.ItemDescSelected.Render("  " + desc)
		fmt.Fprintf(w, "%s\n%s", titleLine, descLine)
	} else {
		titleLine := styles.ItemNormal.Render(fmt.Sprintf("  %-26s %s", title, countBadge))
		descLine := styles.ItemDesc.Render("  " + desc)
		fmt.Fprintf(w, "%s\n%s", titleLine, descLine)
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
		items[idx] = item{module: m}
	}

	l := list.New(items, itemDelegate{}, width-4, height-7)
	l.Title = "📦 Razorpay API Modules (Select a module to view actions)"
	l.Styles.Title = styles.TitleStyle
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(styles.ColorMuted)
	l.Styles.HelpStyle = lipgloss.NewStyle().Foreground(styles.ColorMuted)

	return HomeScreen{
		list:  l,
		state: s,
	}
}

func (h *HomeScreen) SetSize(width, height int) {
	h.list.SetSize(width-4, height-7)
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
			if sel, ok := h.list.SelectedItem().(item); ok {
				h.state.SelectedModule = sel.module
				h.state.PushScreen(state.ScreenActions, sel.module.Title)
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

	// Header
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

	// Contextual Footer Keys
	keys := []components.KeyHelp{
		{Key: "↑/↓", Desc: "Navigate"},
		{Key: "Enter", Desc: "Select Module"},
		{Key: "/", Desc: "Filter"},
		{Key: "c", Desc: "Config"},
		{Key: "q", Desc: "Quit"},
	}
	sb.WriteString(components.RenderFooter(h.state, h.state.Width, keys))

	return sb.String()
}
