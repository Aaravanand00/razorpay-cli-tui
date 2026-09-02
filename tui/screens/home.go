package screens

import (
	"fmt"
	"io"

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
func (d moduleDelegate) Spacing() int                            { return 1 } // Generous space between cards
func (d moduleDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d moduleDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(moduleItem)
	if !ok {
		return
	}

	title := i.module.Title
	desc := i.module.Description

	isSelected := index == m.Index()

	if isSelected {
		countBadge := styles.BadgeCountSelectedStyle.Render(fmt.Sprintf("%d commands", i.module.CommandCount))
		header := lipgloss.JoinHorizontal(
			lipgloss.Center,
			styles.ItemCardSelected.Render(fmt.Sprintf("▶  %-26s", title)),
			" ",
			countBadge,
		)
		descLine := styles.ItemDescSelected.Render(desc)
		fmt.Fprintf(w, "%s\n%s", header, descLine)
	} else {
		countBadge := styles.BadgeCountStyle.Render(fmt.Sprintf("%d commands", i.module.CommandCount))
		header := lipgloss.JoinHorizontal(
			lipgloss.Center,
			styles.ItemCardNormal.Render(fmt.Sprintf("   %-26s", title)),
			" ",
			countBadge,
		)
		descLine := styles.ItemDescNormal.Render(desc)
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
		{ID: "route", Title: "🔀 Route", Description: "Manage linked merchant accounts, transfers and reversals", CommandCount: 20},
		{ID: "smart-collect", Title: "🏢 Smart Collect", Description: "Customer virtual bank accounts, UPI IDs and TPV payers", CommandCount: 13},
		{ID: "settlements", Title: "🏦 Settlements", Description: "View settlement transfers, instant payouts and recon reports", CommandCount: 6},
		{ID: "disputes", Title: "⚖️ Disputes", Description: "Track chargebacks, contest disputes and upload evidence", CommandCount: 4},
		{ID: "documents", Title: "📄 Documents", Description: "Upload KYC/dispute documents and fetch file contents", CommandCount: 3},
		{ID: "configure", Title: "⚙️ Configure", Description: "Configure API credentials, Live/Test mode and defaults", CommandCount: 1},
	}
}

func NewHomeScreen(s *state.SessionState, width, height int) HomeScreen {
	rawModules := GetModules()
	items := make([]list.Item, len(rawModules))
	for idx, m := range rawModules {
		items[idx] = moduleItem{module: m}
	}

	bodyHeight := height - 6
	if bodyHeight < 5 {
		bodyHeight = 5
	}

	l := list.New(items, moduleDelegate{}, width-2, bodyHeight)
	l.Title = "📦 Razorpay API Modules  (Choose a module & press Enter)"
	l.Styles.Title = styles.TitleStyle
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.KeyMap.Quit.Unbind()
	l.KeyMap.ForceQuit.Unbind()
	l.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(styles.ColorSecondary)

	return HomeScreen{
		list:  l,
		state: s,
	}
}

func (h *HomeScreen) IsFiltering() bool {
	return h.list.FilterState() == list.Filtering
}

func (h *HomeScreen) ResetFilter() {
	h.list.ResetFilter()
}

func (h *HomeScreen) SetSize(width, height int) {
	bodyHeight := height - 6
	if h.state.Toast != nil && h.state.Toast.Message != "" {
		bodyHeight -= 3
	}
	if bodyHeight < 5 {
		bodyHeight = 5
	}
	h.list.SetSize(width-2, bodyHeight)
}

func (h *HomeScreen) Update(msg tea.Msg) (HomeScreen, tea.Cmd) {
	var cmd tea.Cmd
	h.SetSize(h.state.Width, h.state.Height)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if h.list.FilterState() == list.Filtering {
			if msg.String() == "esc" || (msg.String() == "/" && h.list.FilterValue() == "") {
				h.list.ResetFilter()
				return *h, nil
			}
			break
		}
		switch msg.String() {
		case "esc":
			// On Home Screen, Esc does nothing (prevents quitting). Only q or Ctrl+C quits.
			return *h, nil
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
	var sections []string

	// 1. Dynamic Header with Stepper & Location
	sections = append(sections, components.RenderHeader(h.state, h.state.Width))

	// 2. Toast (if any)
	toast := components.RenderToast(h.state)
	if toast != "" {
		sections = append(sections, toast)
	}

	// 3. Spacious Module List
	sections = append(sections, h.list.View())

	// 4. Contextual Footer
	var keys []components.KeyHelp
	if h.list.FilterState() == list.Filtering {
		keys = []components.KeyHelp{
			{Key: "Enter", Desc: "Apply Filter"},
			{Key: "Esc", Desc: "Cancel Search"},
		}
	} else {
		keys = []components.KeyHelp{
			{Key: "↑/↓", Desc: "Navigate"},
			{Key: "Enter", Desc: "Select Module"},
			{Key: "/", Desc: "Search"},
			{Key: "c", Desc: "Config (Test/Live)"},
			{Key: "q", Desc: "Quit"},
		}
	}
	sections = append(sections, components.RenderFooter(h.state, h.state.Width, keys))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
