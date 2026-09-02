package screens

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/tui/components"
	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

type TableDataLoadedMsg struct {
	Items   []map[string]interface{}
	RawJSON string
}

type TableDataErrorMsg struct {
	Err string
}

type TableViewScreen struct {
	state       *state.SessionState
	table       table.Model
	spinner     spinner.Model
	loading     bool
	errMsg      string
	rawDataList []map[string]interface{}
	rawJSON     string
	action      state.ActionItem
	width       int
	height      int
	initialized bool
}

func NewTableViewScreen(s *state.SessionState, action state.ActionItem, width, height int) TableViewScreen {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.ColorSecondary)

	columns := getColumnsForResource(action.APIPath, width)

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(calculateTableHeight(height, false)),
	)

	sTable := table.DefaultStyles()
	sTable.Header = styles.TableHeaderStyle
	sTable.Selected = styles.TableCellSelectedStyle
	sTable.Cell = lipgloss.NewStyle().Padding(0, 1)
	t.SetStyles(sTable)

	return TableViewScreen{
		state:       s,
		table:       t,
		spinner:     sp,
		loading:     true,
		action:      action,
		width:       width,
		height:      height,
		initialized: true,
	}
}

func calculateTableHeight(screenHeight int, hasToast bool) int {
	h := screenHeight - 11
	if hasToast {
		h -= 3
	}
	if h < 5 {
		h = 5
	}
	return h
}

func (tv *TableViewScreen) Init() tea.Cmd {
	return tea.Batch(tv.spinner.Tick, tv.fetchDataCmd())
}

func (tv *TableViewScreen) SetSize(width, height int) {
	if !tv.initialized || tv.state == nil {
		return
	}
	tv.width = width
	tv.height = height
	hasToast := tv.state.Toast != nil && tv.state.Toast.Message != ""
	tv.table.SetHeight(calculateTableHeight(height, hasToast))
	tv.table.SetColumns(getColumnsForResource(tv.action.APIPath, width))
}

func (tv *TableViewScreen) fetchDataCmd() tea.Cmd {
	return func() tea.Msg {
		if tv.state.Client == nil || !tv.state.HasCredentials() {
			return TableDataErrorMsg{
				Err: "No API credentials configured. Press 'c' to enter your Razorpay API Key ID and Secret.",
			}
		}

		apiPath := tv.action.APIPath
		// If path has placeholders, default list path
		if strings.Contains(apiPath, "{") {
			parts := strings.Split(apiPath, "/{")
			apiPath = parts[0]
		}

		respBytes, err := tv.state.Client.Get(apiPath, nil)
		if err != nil {
			return TableDataErrorMsg{Err: err.Error()}
		}

		var parsed interface{}
		if err := json.Unmarshal(respBytes, &parsed); err != nil {
			return TableDataErrorMsg{Err: "Failed to parse API response: " + err.Error()}
		}

		rawBytes, _ := json.MarshalIndent(parsed, "", "  ")
		rawStr := string(rawBytes)

		var items []map[string]interface{}

		if respMap, ok := parsed.(map[string]interface{}); ok {
			if rawItems, exists := respMap["items"]; exists {
				if itemList, ok := rawItems.([]interface{}); ok {
					for _, it := range itemList {
						if itemMap, ok := it.(map[string]interface{}); ok {
							items = append(items, itemMap)
						}
					}
				}
			} else {
				// Single entity response
				items = append(items, respMap)
			}
		} else if respSlice, ok := parsed.([]interface{}); ok {
			for _, it := range respSlice {
				if itemMap, ok := it.(map[string]interface{}); ok {
					items = append(items, itemMap)
				}
			}
		}

		return TableDataLoadedMsg{
			Items:   items,
			RawJSON: rawStr,
		}
	}
}

func (tv *TableViewScreen) Update(msg tea.Msg) (TableViewScreen, tea.Cmd) {
	var cmds []tea.Cmd
	tv.SetSize(tv.state.Width, tv.state.Height)

	switch msg := msg.(type) {
	case TableDataLoadedMsg:
		tv.loading = false
		tv.errMsg = ""
		tv.rawDataList = msg.Items
		tv.rawJSON = msg.RawJSON
		tv.populateRows(msg.Items)
		return *tv, nil

	case TableDataErrorMsg:
		tv.loading = false
		tv.errMsg = msg.Err
		return *tv, nil

	case spinner.TickMsg:
		if tv.loading {
			var cmd tea.Cmd
			tv.spinner, cmd = tv.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "down", "j", "k":
			tv.state.ClearToast()
		case "r":
			tv.loading = true
			tv.errMsg = ""
			return *tv, tea.Batch(tv.spinner.Tick, tv.fetchDataCmd())

		case "c":
			tv.state.PushScreen(state.ScreenConfig, "⚙️ Configure")
			return *tv, nil

		case "enter":
			if len(tv.rawDataList) > 0 && tv.table.Cursor() >= 0 && tv.table.Cursor() < len(tv.rawDataList) {
				sel := tv.rawDataList[tv.table.Cursor()]
				tv.state.SelectedRowData = sel
				rawB, _ := json.MarshalIndent(sel, "", "  ")
				tv.state.SelectedRawJSON = string(rawB)

				rowID := fmt.Sprintf("%v", sel["id"])
				if rowID == "<nil>" || rowID == "" {
					rowID = "Record Details"
				}
				tv.state.PushScreen(state.ScreenDetail, "🔍 "+rowID)
			}
			return *tv, nil
		}
	}

	var tableCmd tea.Cmd
	tv.table, tableCmd = tv.table.Update(msg)
	cmds = append(cmds, tableCmd)

	return *tv, tea.Batch(cmds...)
}

func (tv *TableViewScreen) populateRows(items []map[string]interface{}) {
	var rows []table.Row

	for _, item := range items {
		row := extractRowValues(item, tv.action.APIPath)
		rows = append(rows, row)
	}

	tv.table.SetRows(rows)
}

func getColumnsForResource(apiPath string, width int) []table.Column {
	availableWidth := width - 10
	if availableWidth < 60 {
		availableWidth = 60
	}

	if strings.Contains(apiPath, "orders") {
		return []table.Column{
			{Title: "ORDER ID", Width: 22},
			{Title: "AMOUNT", Width: 15},
			{Title: "STATUS", Width: 14},
			{Title: "RECEIPT", Width: 20},
			{Title: "CREATED AT", Width: 24},
		}
	}

	if strings.Contains(apiPath, "payments") {
		return []table.Column{
			{Title: "PAYMENT ID", Width: 22},
			{Title: "AMOUNT", Width: 15},
			{Title: "STATUS", Width: 14},
			{Title: "METHOD", Width: 12},
			{Title: "EMAIL / CONTACT", Width: 25},
			{Title: "CREATED AT", Width: 22},
		}
	}

	if strings.Contains(apiPath, "refunds") {
		return []table.Column{
			{Title: "REFUND ID", Width: 22},
			{Title: "PAYMENT ID", Width: 22},
			{Title: "AMOUNT", Width: 15},
			{Title: "STATUS", Width: 14},
			{Title: "SPEED", Width: 12},
			{Title: "CREATED AT", Width: 22},
		}
	}

	if strings.Contains(apiPath, "customers") {
		return []table.Column{
			{Title: "CUSTOMER ID", Width: 22},
			{Title: "NAME", Width: 24},
			{Title: "EMAIL", Width: 28},
			{Title: "CONTACT", Width: 18},
			{Title: "CREATED AT", Width: 22},
		}
	}

	if strings.Contains(apiPath, "invoices") {
		return []table.Column{
			{Title: "INVOICE ID", Width: 22},
			{Title: "AMOUNT", Width: 15},
			{Title: "STATUS", Width: 14},
			{Title: "CUSTOMER", Width: 24},
			{Title: "ISSUED AT", Width: 22},
		}
	}

	if strings.Contains(apiPath, "payment_links") {
		return []table.Column{
			{Title: "LINK ID", Width: 22},
			{Title: "AMOUNT", Width: 15},
			{Title: "STATUS", Width: 14},
			{Title: "DESCRIPTION", Width: 28},
			{Title: "SHORT URL", Width: 28},
		}
	}

	if strings.Contains(apiPath, "qr_codes") {
		return []table.Column{
			{Title: "QR ID", Width: 22},
			{Title: "USAGE", Width: 14},
			{Title: "STATUS", Width: 14},
			{Title: "PAYMENTS COUNT", Width: 16},
			{Title: "CREATED AT", Width: 22},
		}
	}

	if strings.Contains(apiPath, "subscriptions") {
		return []table.Column{
			{Title: "SUBSCRIPTION ID", Width: 22},
			{Title: "PLAN ID", Width: 22},
			{Title: "STATUS", Width: 14},
			{Title: "TOTAL COUNT", Width: 14},
			{Title: "PAID COUNT", Width: 14},
		}
	}

	if strings.Contains(apiPath, "settlements") {
		return []table.Column{
			{Title: "SETTLEMENT ID", Width: 22},
			{Title: "AMOUNT", Width: 15},
			{Title: "STATUS", Width: 14},
			{Title: "UTR", Width: 24},
			{Title: "CREATED AT", Width: 22},
		}
	}

	if strings.Contains(apiPath, "disputes") {
		return []table.Column{
			{Title: "DISPUTE ID", Width: 22},
			{Title: "PAYMENT ID", Width: 22},
			{Title: "AMOUNT", Width: 15},
			{Title: "STATUS", Width: 14},
			{Title: "REASON", Width: 24},
		}
	}

	// Default fallback columns
	return []table.Column{
		{Title: "ID", Width: 24},
		{Title: "STATUS", Width: 16},
		{Title: "AMOUNT / DETAILS", Width: 28},
		{Title: "CREATED AT", Width: 24},
	}
}

func extractRowValues(item map[string]interface{}, apiPath string) table.Row {
	id := fmt.Sprintf("%v", item["id"])
	if id == "<nil>" {
		id = "-"
	}

	currency := fmt.Sprintf("%v", item["currency"])
	if currency == "<nil>" || currency == "" {
		currency = "INR"
	}

	amountFormatted := formatAmount(item["amount"], currency)
	status := formatStatus(fmt.Sprintf("%v", item["status"]))
	createdAt := formatDate(item["created_at"])

	if strings.Contains(apiPath, "orders") {
		receipt := fmt.Sprintf("%v", item["receipt"])
		if receipt == "<nil>" {
			receipt = "-"
		}
		return table.Row{id, amountFormatted, status, receipt, createdAt}
	}

	if strings.Contains(apiPath, "payments") {
		method := fmt.Sprintf("%v", item["method"])
		if method == "<nil>" {
			method = "-"
		}
		email := fmt.Sprintf("%v", item["email"])
		if email == "<nil>" || email == "" {
			email = fmt.Sprintf("%v", item["contact"])
		}
		return table.Row{id, amountFormatted, status, method, email, createdAt}
	}

	if strings.Contains(apiPath, "refunds") {
		paymentID := fmt.Sprintf("%v", item["payment_id"])
		speed := fmt.Sprintf("%v", item["speed_processed"])
		if speed == "<nil>" || speed == "" {
			speed = "normal"
		}
		return table.Row{id, paymentID, amountFormatted, status, speed, createdAt}
	}

	if strings.Contains(apiPath, "customers") {
		name := fmt.Sprintf("%v", item["name"])
		if name == "<nil>" {
			name = "-"
		}
		email := fmt.Sprintf("%v", item["email"])
		contact := fmt.Sprintf("%v", item["contact"])
		return table.Row{id, name, email, contact, createdAt}
	}

	if strings.Contains(apiPath, "invoices") {
		cust := fmt.Sprintf("%v", item["customer_id"])
		if cust == "<nil>" {
			cust = "-"
		}
		issuedAt := formatDate(item["issued_at"])
		return table.Row{id, amountFormatted, status, cust, issuedAt}
	}

	if strings.Contains(apiPath, "payment_links") {
		desc := fmt.Sprintf("%v", item["description"])
		if desc == "<nil>" {
			desc = "-"
		}
		url := fmt.Sprintf("%v", item["short_url"])
		if url == "<nil>" {
			url = "-"
		}
		return table.Row{id, amountFormatted, status, desc, url}
	}

	if strings.Contains(apiPath, "qr_codes") {
		usage := fmt.Sprintf("%v", item["usage"])
		paymentsCount := fmt.Sprintf("%v", item["payments_amount_received"])
		return table.Row{id, usage, status, paymentsCount, createdAt}
	}

	if strings.Contains(apiPath, "subscriptions") {
		planID := fmt.Sprintf("%v", item["plan_id"])
		totalCount := fmt.Sprintf("%v", item["total_count"])
		paidCount := fmt.Sprintf("%v", item["paid_count"])
		return table.Row{id, planID, status, totalCount, paidCount}
	}

	if strings.Contains(apiPath, "settlements") {
		utr := fmt.Sprintf("%v", item["utr"])
		if utr == "<nil>" {
			utr = "-"
		}
		return table.Row{id, amountFormatted, status, utr, createdAt}
	}

	if strings.Contains(apiPath, "disputes") {
		paymentID := fmt.Sprintf("%v", item["payment_id"])
		reason := fmt.Sprintf("%v", item["reason_code"])
		return table.Row{id, paymentID, amountFormatted, status, reason}
	}

	return table.Row{id, status, amountFormatted, createdAt}
}

func formatAmount(val interface{}, currency string) string {
	if val == nil {
		return "-"
	}

	var paise float64
	switch v := val.(type) {
	case float64:
		paise = v
	case int64:
		paise = float64(v)
	case int:
		paise = float64(v)
	case string:
		p, _ := strconv.ParseFloat(v, 64)
		paise = p
	default:
		return fmt.Sprintf("%v", val)
	}

	rupees := paise / 100.0
	sym := "₹"
	if currency != "INR" && currency != "" {
		sym = currency + " "
	}

	return fmt.Sprintf("%s%.2f", sym, rupees)
}

func formatDate(val interface{}) string {
	if val == nil {
		return "-"
	}

	var sec int64
	switch v := val.(type) {
	case float64:
		sec = int64(v)
	case int64:
		sec = v
	case int:
		sec = int64(v)
	case string:
		s, _ := strconv.ParseInt(v, 10, 64)
		sec = s
	default:
		return fmt.Sprintf("%v", val)
	}

	if sec <= 0 {
		return "-"
	}

	t := time.Unix(sec, 0)
	return t.Format("02 Jan 2006, 03:04 PM")
}

func formatStatus(status string) string {
	st := strings.ToLower(strings.TrimSpace(status))
	switch st {
	case "paid", "captured", "active", "issued", "processed":
		return styles.StatusPillPaid
	case "created", "authorized", "pending", "draft":
		return styles.StatusPillCreated
	case "failed", "cancelled", "expired", "rejected", "closed":
		return styles.StatusPillFailed
	default:
		return st
	}
}

func (tv TableViewScreen) View() string {
	var sections []string

	// 1. Header
	headerView := components.RenderHeader(tv.state, tv.width)
	sections = append(sections, headerView)

	// 2. Toast
	toastView := components.RenderToast(tv.state)
	if toastView != "" {
		sections = append(sections, toastView)
	}

	// 3. Footer Keys
	keys := []components.KeyHelp{
		{Key: "↑/↓", Desc: "Navigate"},
		{Key: "Enter", Desc: "Inspect Record"},
		{Key: "r", Desc: "Refresh"},
		{Key: "c", Desc: "Config"},
		{Key: "?/h", Desc: "Help"},
		{Key: "Esc", Desc: "Back"},
		{Key: "q", Desc: "Quit"},
	}
	footerView := components.RenderFooter(tv.state, tv.width, keys)

	// Dynamic Available Height calculation to pin footer strictly to bottom
	headerHeight := lipgloss.Height(headerView)
	footerHeight := lipgloss.Height(footerView)
	toastHeight := 0
	if toastView != "" {
		toastHeight = lipgloss.Height(toastView)
	}

	bodyHeight := tv.height - headerHeight - footerHeight - toastHeight
	if bodyHeight < 10 {
		bodyHeight = 10
	}

	// 4. Main Body Content
	var bodyContent string
	title := styles.TitleStyle.PaddingLeft(1).Render(fmt.Sprintf("%s (%s)", tv.action.Title, tv.action.CLICommand))

	if !tv.state.HasCredentials() {
		modeName := "Test Sandbox"
		if tv.state.IsLiveMode {
			modeName = "Live Production"
		}
		authCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorSecondary).
			Background(styles.ColorCardBg).
			Padding(2, 4).
			Width(tv.width - 6).
			Render(fmt.Sprintf(
				"🔑  Razorpay API Credentials Required (%s)\n\n"+
					"To fetch and view live %s records from Razorpay Cloud, please enter your API Key ID and Key Secret.\n\n"+
					"👉 Press [c] to Configure Credentials\n"+
					"👉 Press [Esc] to go back to Actions",
				modeName, tv.action.Title,
			))
		bodyContent = "\n" + title + "\n\n" + authCard
	} else if tv.loading {
		loadingMsg := lipgloss.NewStyle().
			Padding(4, 2).
			Render(fmt.Sprintf("%s Fetching live records from Razorpay API (%s)...", tv.spinner.View(), tv.action.APIPath))
		bodyContent = "\n" + title + "\n\n" + loadingMsg
	} else if tv.errMsg != "" {
		errCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorError).
			Padding(1, 2).
			Width(tv.width - 6).
			Render(fmt.Sprintf("⚠️  API Request Failed:\n%s\n\n💡 Tip: Press 'r' to Retry or 'c' to Configure Credentials.", tv.errMsg))
		bodyContent = "\n" + title + "\n\n" + errCard
	} else if len(tv.rawDataList) == 0 {
		emptyCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorBorder).
			Padding(2, 3).
			Width(tv.width - 6).
			Render("📭 No records found for this resource in your active environment.\n\n💡 Tip: Press 'r' to refresh or 'Esc' to go back.")
		bodyContent = "\n" + title + "\n\n" + emptyCard
	} else {
		countBar := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			PaddingLeft(1).
			Render(fmt.Sprintf("Showing %d records • Select row and press [Enter] to inspect full JSON details", len(tv.rawDataList)))

		tableView := styles.TableContainerStyle.
			Width(tv.width - 4).
			Render(tv.table.View())

		bodyContent = "\n" + title + "\n" + countBar + "\n" + tableView
	}

	bodyContainer := lipgloss.NewStyle().
		Height(bodyHeight).
		Width(tv.width - 2)

	sections = append(sections, bodyContainer.Render(bodyContent))
	sections = append(sections, footerView)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
