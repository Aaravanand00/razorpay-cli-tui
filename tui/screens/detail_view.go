package screens

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/tui/components"
	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

type DetailViewScreen struct {
	state       *state.SessionState
	viewport    viewport.Model
	activeTab   int // 0: Summary Card, 1: Raw JSON Tree
	data        map[string]interface{}
	rawJSON     string
	width       int
	height      int
	initialized bool
}

func NewDetailViewScreen(s *state.SessionState, width, height int) DetailViewScreen {
	vp := viewport.New(width-6, calculateDetailViewportHeight(height, false))
	vp.SetContent(formatSyntaxJSON(s.SelectedRawJSON))

	return DetailViewScreen{
		state:       s,
		viewport:    vp,
		activeTab:   0,
		data:        s.SelectedRowData,
		rawJSON:     s.SelectedRawJSON,
		width:       width,
		height:      height,
		initialized: true,
	}
}

func calculateDetailViewportHeight(screenHeight int, hasToast bool) int {
	h := screenHeight - 16
	if hasToast {
		h -= 2
	}
	if h < 5 {
		h = 5
	}
	return h
}

func (dv *DetailViewScreen) SetSize(width, height int) {
	if !dv.initialized || dv.state == nil {
		return
	}
	dv.width = width
	dv.height = height
	hasToast := dv.state.Toast != nil && dv.state.Toast.Message != ""
	dv.viewport.Width = width - 6
	dv.viewport.Height = calculateDetailViewportHeight(height, hasToast)
}

func (dv *DetailViewScreen) Update(msg tea.Msg) (DetailViewScreen, tea.Cmd) {
	var cmds []tea.Cmd
	dv.SetSize(dv.state.Width, dv.state.Height)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "v", "right", "left":
			if dv.activeTab == 0 {
				dv.activeTab = 1
			} else {
				dv.activeTab = 0
			}
			return *dv, nil

		case "1":
			dv.activeTab = 0
			return *dv, nil

		case "2":
			dv.activeTab = 1
			return *dv, nil

		case "c":
			dv.state.PushScreen(state.ScreenConfig, "⚙️ Configure")
			return *dv, nil
		}
	}

	if dv.activeTab == 1 {
		var vpCmd tea.Cmd
		dv.viewport, vpCmd = dv.viewport.Update(msg)
		cmds = append(cmds, vpCmd)
	}

	return *dv, tea.Batch(cmds...)
}

func (dv DetailViewScreen) View() string {
	var sections []string

	// 1. Header
	headerView := components.RenderHeader(dv.state, dv.width)
	sections = append(sections, headerView)

	// 2. Toast
	toastView := components.RenderToast(dv.state)
	if toastView != "" {
		sections = append(sections, toastView)
	}

	// 3. Footer
	keys := []components.KeyHelp{
		{Key: "Tab / v", Desc: "Switch Mode"},
		{Key: "↑/↓", Desc: "Scroll JSON"},
		{Key: "a", Desc: "AI Assist"},
		{Key: "c", Desc: "Config"},
		{Key: "?/h", Desc: "Help"},
		{Key: "Esc", Desc: "Back to Table"},
		{Key: "q", Desc: "Quit"},
	}
	footerView := components.RenderFooter(dv.state, dv.width, keys)

	// Height math
	headerHeight := lipgloss.Height(headerView)
	footerHeight := lipgloss.Height(footerView)
	toastHeight := 0
	if toastView != "" {
		toastHeight = lipgloss.Height(toastView)
	}

	bodyHeight := dv.height - headerHeight - footerHeight - toastHeight
	if bodyHeight < 10 {
		bodyHeight = 10
	}

	// Tab Bar
	var tab1, tab2 string
	if dv.activeTab == 0 {
		tab1 = styles.TabActiveStyle.Render("● 1. Summary Card")
		tab2 = styles.TabInactiveStyle.Render("○ 2. Raw JSON Tree")
	} else {
		tab1 = styles.TabInactiveStyle.Render("○ 1. Summary Card")
		tab2 = styles.TabActiveStyle.Render("● 2. Raw JSON Tree")
	}
	tabsRow := lipgloss.JoinHorizontal(lipgloss.Center, "  ", tab1, "  ", tab2)

	// Body Content
	var mainContent string
	if dv.activeTab == 0 {
		mainContent = dv.renderSummaryCard()
	} else {
		vpContainer := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorBorder).
			Padding(0, 1).
			Width(dv.width - 6).
			Render(dv.viewport.View())

		scrollInfo := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			PaddingLeft(1).
			Render(fmt.Sprintf("Scroll: %3.f%% • Use [↑/↓] or [PgUp/PgDn] to scroll JSON payload", dv.viewport.ScrollPercent()*100))

		mainContent = vpContainer + "\n" + scrollInfo
	}

	rawBody := "\n" + tabsRow + "\n\n" + mainContent

	bodyContainer := lipgloss.NewStyle().
		Height(bodyHeight).
		Width(dv.width - 2)

	sections = append(sections, bodyContainer.Render(rawBody))
	sections = append(sections, footerView)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (dv DetailViewScreen) renderSummaryCard() string {
	if len(dv.data) == 0 {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorBorder).
			Padding(2, 3).
			Width(dv.width - 6).
			Render("📭 No detailed record data available. Press [Tab] to view Raw JSON.")
	}

	var fields []string

	// Primary ID and Entity
	id := fmt.Sprintf("%v", dv.data["id"])
	if id != "" && id != "<nil>" {
		fields = append(fields, formatDetailRow("Record ID", id))
	}

	entity := fmt.Sprintf("%v", dv.data["entity"])
	if entity != "" && entity != "<nil>" {
		fields = append(fields, formatDetailRow("Entity Type", entity))
	}

	// Status
	status := fmt.Sprintf("%v", dv.data["status"])
	if status != "" && status != "<nil>" {
		fields = append(fields, formatDetailRow("Status", formatStatus(status)))
	}

	// Amount & Currency
	if amt, exists := dv.data["amount"]; exists {
		curr := fmt.Sprintf("%v", dv.data["currency"])
		fields = append(fields, formatDetailRow("Amount", formatAmount(amt, curr)))
	}

	// Receipt
	if rcpt, exists := dv.data["receipt"]; exists && fmt.Sprintf("%v", rcpt) != "<nil>" {
		fields = append(fields, formatDetailRow("Receipt", fmt.Sprintf("%v", rcpt)))
	}

	// Customer Info
	if email, exists := dv.data["email"]; exists && fmt.Sprintf("%v", email) != "<nil>" {
		fields = append(fields, formatDetailRow("Email", fmt.Sprintf("%v", email)))
	}
	if contact, exists := dv.data["contact"]; exists && fmt.Sprintf("%v", contact) != "<nil>" {
		fields = append(fields, formatDetailRow("Contact", fmt.Sprintf("%v", contact)))
	}
	if method, exists := dv.data["method"]; exists && fmt.Sprintf("%v", method) != "<nil>" {
		fields = append(fields, formatDetailRow("Payment Method", fmt.Sprintf("%v", method)))
	}

	// Timestamps
	if createdAt, exists := dv.data["created_at"]; exists {
		fields = append(fields, formatDetailRow("Created At", formatDate(createdAt)))
	}
	if issuedAt, exists := dv.data["issued_at"]; exists {
		fields = append(fields, formatDetailRow("Issued At", formatDate(issuedAt)))
	}

	// Other Attributes
	var remainingKeys []string
	excludedKeys := map[string]bool{
		"id": true, "entity": true, "status": true, "amount": true, "currency": true,
		"receipt": true, "email": true, "contact": true, "method": true, "created_at": true,
		"issued_at": true, "notes": true,
	}

	for k, v := range dv.data {
		if !excludedKeys[k] && v != nil && fmt.Sprintf("%v", v) != "<nil>" && fmt.Sprintf("%v", v) != "[]" && fmt.Sprintf("%v", v) != "map[]" {
			remainingKeys = append(remainingKeys, k)
		}
	}
	sort.Strings(remainingKeys)

	for _, k := range remainingKeys {
		val := fmt.Sprintf("%v", dv.data[k])
		// Format amount fields (like amount_due, amount_paid, refund_amount)
		if strings.HasPrefix(strings.ToLower(k), "amount_") || strings.HasSuffix(strings.ToLower(k), "_amount") {
			curr := fmt.Sprintf("%v", dv.data["currency"])
			val = formatAmount(dv.data[k], curr)
		} else if len(val) > 60 {
			val = val[:57] + "..."
		}
		label := strings.Title(strings.ReplaceAll(k, "_", " "))
		fields = append(fields, formatDetailRow(label, val))
	}

	// Notes card if present
	var notesCard string
	if rawNotes, ok := dv.data["notes"]; ok && rawNotes != nil {
		if notesMap, ok := rawNotes.(map[string]interface{}); ok && len(notesMap) > 0 {
			var noteLines []string
			for nk, nv := range notesMap {
				noteLines = append(noteLines, fmt.Sprintf("  • %-18s: %v", nk, nv))
			}
			notesCard = "\n\n" + lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary).Render("📝 Metadata / Notes:") + "\n" + strings.Join(noteLines, "\n")
		}
	}

	cardContent := strings.Join(fields, "\n") + notesCard

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Padding(1, 3).
		Width(dv.width - 6).
		Render(cardContent)
}

func formatDetailRow(label, val string) string {
	keyStyled := styles.DetailKeyStyle.Render(fmt.Sprintf("%-20s :", label))
	valStyled := styles.DetailValStyle.Render(val)
	return lipgloss.JoinHorizontal(lipgloss.Center, keyStyled, "  ", valStyled)
}

func formatSyntaxJSON(raw string) string {
	if raw == "" {
		return "{\n  \"message\": \"No JSON payload available\"\n}"
	}

	var parsed interface{}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return raw
	}

	prettyBytes, _ := json.MarshalIndent(parsed, "", "  ")
	lines := strings.Split(string(prettyBytes), "\n")

	var highlightedLines []string
	for idx, line := range lines {
		lineNum := lipgloss.NewStyle().Foreground(styles.ColorTextDim).Render(fmt.Sprintf("%3d │ ", idx+1))

		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "\"") && strings.Contains(trimmed, "\":") {
			parts := strings.SplitN(line, "\":", 2)
			keyPart := styles.JSONKeyStyle.Render(parts[0] + "\"")
			valPart := colorizeJSONValue(parts[1])
			highlightedLines = append(highlightedLines, lineNum+keyPart+":"+valPart)
		} else {
			highlightedLines = append(highlightedLines, lineNum+colorizeJSONValue(line))
		}
	}

	return strings.Join(highlightedLines, "\n")
}

func colorizeJSONValue(val string) string {
	trimmed := strings.TrimSpace(val)

	if strings.HasPrefix(trimmed, "\"") {
		return styles.JSONStringStyle.Render(val)
	}
	if trimmed == "true" || trimmed == "false" {
		return styles.JSONBoolStyle.Render(val)
	}
	if trimmed == "null" {
		return styles.JSONNullStyle.Render(val)
	}
	if _, err := strconv.ParseFloat(strings.TrimRight(trimmed, ","), 64); err == nil {
		return styles.JSONNumberStyle.Render(val)
	}

	return val
}
