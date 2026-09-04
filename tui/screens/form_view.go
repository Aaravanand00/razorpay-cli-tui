package screens

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/razorpay/razorpay-cli/tui/components"
	"github.com/razorpay/razorpay-cli/tui/state"
	"github.com/razorpay/razorpay-cli/tui/styles"
)

type FormExecutionResultMsg struct {
	Success bool
	Data    map[string]interface{}
	RawJSON string
	Err     string
}

type FormField struct {
	Key         string
	Label       string
	Placeholder string
	Required    bool
	IsAmount    bool
	Input       textinput.Model
}

type FormState int

const (
	FormStateEditing FormState = iota
	FormStateSubmitting
	FormStateSuccess
	FormStateError
)

type FormViewScreen struct {
	state       *state.SessionState
	action      state.ActionItem
	fields      []FormField
	focusedIdx  int
	formState   FormState
	spinner     spinner.Model
	resultData  map[string]interface{}
	resultRaw   string
	resultErr   string
	width       int
	height      int
	initialized bool
}

func NewFormViewScreen(s *state.SessionState, action state.ActionItem, width, height int) FormViewScreen {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(styles.ColorSecondary)

	fields := generateFieldsForAction(action, width)

	if len(fields) > 0 {
		fields[0].Input.Focus()
	}

	return FormViewScreen{
		state:       s,
		action:      action,
		fields:      fields,
		focusedIdx:  0,
		formState:   FormStateEditing,
		spinner:     sp,
		width:       width,
		height:      height,
		initialized: true,
	}
}

// PopulateWithFlags populates form fields with flag values (e.g. from AI suggestion).
func (fv *FormViewScreen) PopulateWithFlags(flags map[string]string) {
	if len(flags) == 0 {
		return
	}
	for i := range fv.fields {
		k := fv.fields[i].Key
		val, exists := flags[k]
		if !exists {
			val, exists = flags[strings.ReplaceAll(k, "_", "-")]
		}
		if exists {
			// If this form field represents Rupees (₹ INR), convert CLI paise value to Rupees
			if fv.fields[i].IsAmount {
				if amtPaise, err := strconv.ParseFloat(val, 64); err == nil {
					if amtPaise >= 100 && int64(amtPaise)%100 == 0 {
						val = fmt.Sprintf("%d", int64(amtPaise/100))
					} else if amtPaise >= 100 {
						val = fmt.Sprintf("%.2f", amtPaise/100.0)
					}
				}
			}
			fv.fields[i].Input.SetValue(val)
		}
	}
}

func (fv FormViewScreen) IsSuccess() bool {
	return fv.formState == FormStateSuccess
}

func (fv FormViewScreen) IsError() bool {
	return fv.formState == FormStateError
}

func (fv *FormViewScreen) SetSize(width, height int) {
	if !fv.initialized || fv.state == nil {
		return
	}
	fv.width = width
	fv.height = height
	for i := range fv.fields {
		fv.fields[i].Input.Width = width - 32
	}
}

func (fv *FormViewScreen) Update(msg tea.Msg) (FormViewScreen, tea.Cmd) {
	var cmds []tea.Cmd
	fv.SetSize(fv.state.Width, fv.state.Height)

	switch msg := msg.(type) {
	case FormExecutionResultMsg:
		if msg.Success {
			fv.formState = FormStateSuccess
			fv.resultData = msg.Data
			fv.resultRaw = msg.RawJSON
			fv.state.SelectedRowData = msg.Data
			fv.state.SelectedRawJSON = msg.RawJSON
		} else {
			fv.formState = FormStateError
			fv.resultErr = msg.Err
		}
		return *fv, nil

	case spinner.TickMsg:
		if fv.formState == FormStateSubmitting {
			var cmd tea.Cmd
			fv.spinner, cmd = fv.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			if fv.formState == FormStateEditing && len(fv.fields) > 0 {
				fv.fields[fv.focusedIdx].Input.Blur()
				fv.focusedIdx = (fv.focusedIdx + 1) % len(fv.fields)
				fv.fields[fv.focusedIdx].Input.Focus()
				fv.state.ClearToast()
				return *fv, textinput.Blink
			}

		case "shift+tab", "up":
			if fv.formState == FormStateEditing && len(fv.fields) > 0 {
				fv.fields[fv.focusedIdx].Input.Blur()
				fv.focusedIdx--
				if fv.focusedIdx < 0 {
					fv.focusedIdx = len(fv.fields) - 1
				}
				fv.fields[fv.focusedIdx].Input.Focus()
				fv.state.ClearToast()
				return *fv, textinput.Blink
			}

		case "q":
			if fv.formState == FormStateError || fv.formState == FormStateSuccess {
				return *fv, tea.Quit
			}

		case "a":
			if fv.formState == FormStateError || fv.formState == FormStateSuccess {
				fv.state.PushScreen(state.ScreenAIAssist, "🤖 AI Assist")
				return *fv, nil
			}

		case "r":
			if fv.formState == FormStateError || fv.formState == FormStateSuccess {
				fv.formState = FormStateEditing
				return *fv, nil
			}

		case "c":
			if fv.formState == FormStateError || fv.formState == FormStateSuccess {
				fv.state.PushScreen(state.ScreenConfig, "⚙️ Configure")
				return *fv, nil
			}

		case "enter":
			if fv.formState == FormStateSuccess {
				// Transition to Detail View
				rowID := fmt.Sprintf("%v", fv.resultData["id"])
				if rowID == "<nil>" || rowID == "" {
					rowID = "Created Record"
				}
				fv.state.PushScreen(state.ScreenDetail, "🔍 "+rowID)
				return *fv, nil
			}

			if fv.formState == FormStateError {
				fv.formState = FormStateEditing
				return *fv, nil
			}

			if fv.formState == FormStateEditing {
				// Validate required fields
				for _, f := range fv.fields {
					if f.Required && strings.TrimSpace(f.Input.Value()) == "" {
						fv.state.SetToast(fmt.Sprintf("⚠️ '%s' is required. Please fill this field.", f.Label), true)
						return *fv, nil
					}
				}

				// Check permission
				allowed, errMsg := fv.state.CanPerformActionItem(fv.action)
				if !allowed {
					fv.state.SetToast(errMsg, true)
					return *fv, nil
				}

				fv.formState = FormStateSubmitting
				return *fv, tea.Batch(fv.spinner.Tick, fv.executeActionCmd())
			}
		}
	}

	if fv.formState == FormStateEditing && len(fv.fields) > 0 {
		var inputCmd tea.Cmd
		fv.fields[fv.focusedIdx].Input, inputCmd = fv.fields[fv.focusedIdx].Input.Update(msg)
		cmds = append(cmds, inputCmd)
	}

	return *fv, tea.Batch(cmds...)
}

func (fv *FormViewScreen) executeActionCmd() tea.Cmd {
	return func() tea.Msg {
		if fv.state.Client == nil || !fv.state.HasCredentials() {
			return FormExecutionResultMsg{
				Success: false,
				Err:     "No active API credentials configured. Press 'c' to enter your Key ID and Secret.",
			}
		}

		apiPath := fv.action.APIPath
		payload := make(map[string]interface{})

		// Extract values
		for _, f := range fv.fields {
			val := strings.TrimSpace(f.Input.Value())
			if val == "" {
				continue
			}

			// Replace URL params like {id}
			if strings.Contains(apiPath, "{"+f.Key+"}") || (f.Key == "id" && strings.Contains(apiPath, "{id}")) {
				apiPath = strings.ReplaceAll(apiPath, "{"+f.Key+"}", val)
				apiPath = strings.ReplaceAll(apiPath, "{id}", val)
				continue
			}

			// Parse Amounts to Paise
			if f.IsAmount {
				amtRupees, err := strconv.ParseFloat(val, 64)
				if err == nil {
					payload[f.Key] = int64(amtRupees * 100) // Convert to paise
				} else {
					payload[f.Key] = val
				}
				continue
			}

			// Parse Notes as JSON dictionary if passed as key=val or JSON
			if f.Key == "notes" {
				if strings.HasPrefix(val, "{") {
					var notesMap map[string]interface{}
					if err := json.Unmarshal([]byte(val), &notesMap); err == nil {
						payload["notes"] = notesMap
						continue
					}
				}
				payload["notes"] = map[string]string{"note": val}
				continue
			}

			// Smart Phone / Contact Normalizer (+91 default for 10 digits, preserves international)
			if f.Key == "contact" || f.Key == "customer_contact" || f.Key == "phone" {
				payload[f.Key] = normalizePhoneNumber(val)
				continue
			}

			payload[f.Key] = val
		}

		// Payment Links payload normalization (Razorpay API expects nested customer object)
		if strings.Contains(fv.action.ID, "payment-link") && strings.Contains(fv.action.ID, "create") {
			customer := make(map[string]interface{})
			if name, ok := payload["customer_name"]; ok && fmt.Sprintf("%v", name) != "" {
				customer["name"] = name
				delete(payload, "customer_name")
			}
			if email, ok := payload["customer_email"]; ok && fmt.Sprintf("%v", email) != "" {
				customer["email"] = email
				delete(payload, "customer_email")
			}
			if contact, ok := payload["customer_contact"]; ok && fmt.Sprintf("%v", contact) != "" {
				customer["contact"] = contact
				delete(payload, "customer_contact")
			}
			if len(customer) > 0 {
				payload["customer"] = customer
			}
		}

		// QR Codes payload normalization (Razorpay API expects type=upi_qr and payment_amount in paise)
		if strings.Contains(fv.action.ID, "qr") && strings.Contains(fv.action.ID, "create") {
			payload["type"] = "upi_qr"
			if amt, ok := payload["amount"]; ok {
				payload["payment_amount"] = amt
				payload["fixed_amount"] = true
				delete(payload, "amount")
			}
		}

		// Plans payload normalization
		if strings.Contains(fv.action.ID, "plans") && strings.Contains(fv.action.ID, "create") {
			item := make(map[string]interface{})
			if name, ok := payload["name"]; ok {
				item["name"] = name
				delete(payload, "name")
			}
			if amt, ok := payload["amount"]; ok {
				item["amount"] = amt
				delete(payload, "amount")
			}
			if curr, ok := payload["currency"]; ok {
				item["currency"] = curr
				delete(payload, "currency")
			} else {
				item["currency"] = "INR"
			}
			if desc, ok := payload["description"]; ok {
				item["description"] = desc
				delete(payload, "description")
			}
			payload["item"] = item
			if intervalVal, ok := payload["interval"]; ok {
				if intervalInt, err := strconv.Atoi(fmt.Sprintf("%v", intervalVal)); err == nil {
					payload["interval"] = intervalInt
				}
			}
		}

		// Subscriptions create total_count & quantity integer conversion
		if strings.Contains(fv.action.ID, "subs-create") {
			if tc, ok := payload["total_count"]; ok {
				if tcInt, err := strconv.Atoi(fmt.Sprintf("%v", tc)); err == nil {
					payload["total_count"] = tcInt
				}
			}
			if qty, ok := payload["quantity"]; ok {
				if qtyInt, err := strconv.Atoi(fmt.Sprintf("%v", qty)); err == nil {
					payload["quantity"] = qtyInt
				}
			}
		}

		var respBytes []byte
		var err error

		method := strings.ToUpper(fv.action.HTTPMethod)
		switch method {
		case "POST":
			respBytes, err = fv.state.Client.Post(apiPath, payload)
		case "PATCH":
			respBytes, err = fv.state.Client.Patch(apiPath, payload)
		case "PUT":
			respBytes, err = fv.state.Client.Put(apiPath, payload)
		case "DELETE":
			respBytes, err = fv.state.Client.Delete(apiPath)
		case "GET":
			respBytes, err = fv.state.Client.Get(apiPath, nil)
		default:
			respBytes, err = fv.state.Client.Post(apiPath, payload)
		}

		if err != nil {
			return FormExecutionResultMsg{
				Success: false,
				Err:     err.Error(),
			}
		}

		var resultData map[string]interface{}
		_ = json.Unmarshal(respBytes, &resultData)

		rawBytes, _ := json.MarshalIndent(resultData, "", "  ")

		return FormExecutionResultMsg{
			Success: true,
			Data:    resultData,
			RawJSON: string(rawBytes),
		}
	}
}

func generateFieldsForAction(action state.ActionItem, width int) []FormField {
	inputWidth := width - 34
	if inputWidth < 40 {
		inputWidth = 40
	}

	createInput := func(placeholder string, isPassword bool) textinput.Model {
		ti := textinput.New()
		ti.Placeholder = placeholder
		ti.Width = inputWidth
		ti.CharLimit = 120
		if isPassword {
			ti.EchoMode = textinput.EchoPassword
			ti.EchoCharacter = '•'
		}
		return ti
	}

	id := action.ID

	// 1. ORDERS MODULE
	if id == "orders-create" {
		return []FormField{
			{Key: "amount", Label: "Amount (₹ INR)", Placeholder: "e.g. 500 (will auto-convert to 50000 paise)", Required: true, IsAmount: true, Input: createInput("500.00", false)},
			{Key: "currency", Label: "Currency", Placeholder: "INR", Required: true, Input: createInput("INR", false)},
			{Key: "receipt", Label: "Receipt Number", Placeholder: "e.g. Receipt #101", Input: createInput("Receipt #101", false)},
			{Key: "notes", Label: "Notes / Description", Placeholder: "e.g. Order for premium subscription", Input: createInput("Order note", false)},
		}
	}
	if id == "orders-update" {
		return []FormField{
			{Key: "id", Label: "Order ID", Placeholder: "e.g. order_xxx", Required: true, Input: createInput("order_xxx", false)},
			{Key: "notes", Label: "Notes / Metadata", Placeholder: "e.g. Updated customer delivery tag", Input: createInput("Updated notes", false)},
		}
	}
	if id == "orders-fetch" || id == "orders-payments" {
		return []FormField{
			{Key: "id", Label: "Order ID", Placeholder: "e.g. order_xxx", Required: true, Input: createInput("order_xxx", false)},
		}
	}

	// 2. PAYMENTS MODULE
	if id == "payments-capture" {
		return []FormField{
			{Key: "id", Label: "Payment ID", Placeholder: "e.g. pay_xxx", Required: true, Input: createInput("pay_xxx", false)},
			{Key: "amount", Label: "Amount (₹ INR)", Placeholder: "e.g. 500.00", Required: true, IsAmount: true, Input: createInput("500.00", false)},
			{Key: "currency", Label: "Currency", Placeholder: "INR", Required: true, Input: createInput("INR", false)},
		}
	}
	if id == "payments-update" {
		return []FormField{
			{Key: "id", Label: "Payment ID", Placeholder: "e.g. pay_xxx", Required: true, Input: createInput("pay_xxx", false)},
			{Key: "notes", Label: "Notes / Metadata", Placeholder: "e.g. Internal tracking notes", Input: createInput("Updated notes", false)},
		}
	}
	if id == "payments-fetch" || id == "payments-card" {
		return []FormField{
			{Key: "id", Label: "Payment ID", Placeholder: "e.g. pay_xxx", Required: true, Input: createInput("pay_xxx", false)},
		}
	}
	if id == "payments-downtime-fetch" {
		return []FormField{
			{Key: "id", Label: "Downtime ID", Placeholder: "e.g. down_xxx", Required: true, Input: createInput("down_xxx", false)},
		}
	}

	// 3. REFUNDS MODULE
	if id == "refunds-create" {
		return []FormField{
			{Key: "id", Label: "Payment ID", Placeholder: "e.g. pay_xxx", Required: true, Input: createInput("pay_xxx", false)},
			{Key: "amount", Label: "Amount (₹ INR)", Placeholder: "e.g. 500 (leave blank for full refund)", IsAmount: true, Input: createInput("", false)},
			{Key: "speed", Label: "Speed", Placeholder: "normal / optimum", Input: createInput("normal", false)},
			{Key: "notes", Label: "Refund Reason / Notes", Placeholder: "e.g. Customer return request", Input: createInput("Customer return", false)},
		}
	}
	if id == "refunds-update" {
		return []FormField{
			{Key: "id", Label: "Refund ID", Placeholder: "e.g. rfr_xxx", Required: true, Input: createInput("rfr_xxx", false)},
			{Key: "notes", Label: "Notes / Metadata", Placeholder: "e.g. Updated refund note", Input: createInput("Updated notes", false)},
		}
	}
	if id == "refunds-fetch" {
		return []FormField{
			{Key: "id", Label: "Refund ID", Placeholder: "e.g. rfr_xxx", Required: true, Input: createInput("rfr_xxx", false)},
		}
	}
	if id == "refunds-payment-refunds" {
		return []FormField{
			{Key: "id", Label: "Payment ID", Placeholder: "e.g. pay_xxx", Required: true, Input: createInput("pay_xxx", false)},
		}
	}
	if id == "refunds-payment-refund" {
		return []FormField{
			{Key: "id", Label: "Payment ID", Placeholder: "e.g. pay_xxx", Required: true, Input: createInput("pay_xxx", false)},
			{Key: "refund_id", Label: "Refund ID", Placeholder: "e.g. rfr_xxx", Required: true, Input: createInput("rfr_xxx", false)},
		}
	}

	// 4. CUSTOMERS MODULE
	if id == "customers-create" {
		return []FormField{
			{Key: "name", Label: "Customer Name", Placeholder: "e.g. Rahul Sharma", Required: true, Input: createInput("Rahul Sharma", false)},
			{Key: "email", Label: "Email Address", Placeholder: "e.g. rahul@example.com", Required: true, Input: createInput("rahul@example.com", false)},
			{Key: "contact", Label: "Phone / Contact", Placeholder: "e.g. 9876543210 (auto +91) or +14155552671", Input: createInput("9876543210", false)},
			{Key: "notes", Label: "Customer Notes", Placeholder: "e.g. VIP Customer", Input: createInput("VIP Customer", false)},
		}
	}
	if id == "customers-update" {
		return []FormField{
			{Key: "id", Label: "Customer ID", Placeholder: "e.g. cust_xxx", Required: true, Input: createInput("cust_xxx", false)},
			{Key: "name", Label: "Customer Name", Placeholder: "e.g. Rahul Sharma (Updated)", Input: createInput("", false)},
			{Key: "email", Label: "Email Address", Placeholder: "e.g. rahul.new@example.com", Input: createInput("", false)},
			{Key: "contact", Label: "Phone / Contact", Placeholder: "e.g. 9876543210 (auto +91)", Input: createInput("", false)},
			{Key: "notes", Label: "Customer Notes", Placeholder: "e.g. Updated metadata", Input: createInput("", false)},
		}
	}
	if id == "customers-fetch" {
		return []FormField{
			{Key: "id", Label: "Customer ID", Placeholder: "e.g. cust_xxx", Required: true, Input: createInput("cust_xxx", false)},
		}
	}

	// 5. PAYMENT LINKS MODULE
	if id == "payment-links-create" {
		return []FormField{
			{Key: "amount", Label: "Amount (₹ INR)", Placeholder: "e.g. 999.00", Required: true, IsAmount: true, Input: createInput("999.00", false)},
			{Key: "description", Label: "Description", Placeholder: "e.g. Payment for Invoice #102", Required: true, Input: createInput("Payment for order", false)},
			{Key: "customer_name", Label: "Customer Name", Placeholder: "e.g. Amit Kumar", Input: createInput("Amit Kumar", false)},
			{Key: "customer_email", Label: "Customer Email", Placeholder: "e.g. amit@example.com", Input: createInput("amit@example.com", false)},
			{Key: "customer_contact", Label: "Customer Contact", Placeholder: "e.g. 9876543210 (auto +91) or +14155552671", Input: createInput("9876543210", false)},
		}
	}
	if id == "payment-links-update" {
		return []FormField{
			{Key: "id", Label: "Payment Link ID", Placeholder: "e.g. plink_xxx", Required: true, Input: createInput("plink_xxx", false)},
			{Key: "reference_id", Label: "Reference ID", Placeholder: "e.g. ref_12345", Input: createInput("", false)},
			{Key: "notes", Label: "Notes / Metadata", Placeholder: "e.g. Updated link notes", Input: createInput("", false)},
		}
	}
	if id == "payment-links-fetch" || id == "payment-links-cancel" {
		return []FormField{
			{Key: "id", Label: "Payment Link ID", Placeholder: "e.g. plink_xxx", Required: true, Input: createInput("plink_xxx", false)},
		}
	}
	if id == "payment-links-notify" {
		return []FormField{
			{Key: "id", Label: "Payment Link ID", Placeholder: "e.g. plink_xxx", Required: true, Input: createInput("plink_xxx", false)},
			{Key: "medium", Label: "Notify Medium", Placeholder: "sms or email", Required: true, Input: createInput("sms", false)},
		}
	}

	// 6. INVOICES MODULE
	if id == "invoices-create" {
		return []FormField{
			{Key: "customer_id", Label: "Customer ID", Placeholder: "e.g. cust_xxx", Required: true, Input: createInput("cust_xxx", false)},
			{Key: "amount", Label: "Amount (₹ INR)", Placeholder: "e.g. 1500.00", Required: true, IsAmount: true, Input: createInput("1500.00", false)},
			{Key: "currency", Label: "Currency", Placeholder: "INR", Required: true, Input: createInput("INR", false)},
			{Key: "description", Label: "Invoice Description", Placeholder: "e.g. Web development services", Input: createInput("Consulting services", false)},
		}
	}
	if id == "invoices-update" {
		return []FormField{
			{Key: "id", Label: "Invoice ID", Placeholder: "e.g. inv_xxx", Required: true, Input: createInput("inv_xxx", false)},
			{Key: "description", Label: "Description", Placeholder: "e.g. Updated invoice description", Input: createInput("", false)},
			{Key: "notes", Label: "Notes / Metadata", Placeholder: "e.g. Updated terms", Input: createInput("", false)},
		}
	}
	if id == "invoices-fetch" || id == "invoices-issue" || id == "invoices-notify" || id == "invoices-cancel" || id == "invoices-delete" {
		return []FormField{
			{Key: "id", Label: "Invoice ID", Placeholder: "e.g. inv_xxx", Required: true, Input: createInput("inv_xxx", false)},
		}
	}
	if id == "invoices-items-create" {
		return []FormField{
			{Key: "name", Label: "Item Name", Placeholder: "e.g. Cloud Hosting (1 Month)", Required: true, Input: createInput("Cloud Hosting", false)},
			{Key: "amount", Label: "Unit Price (₹ INR)", Placeholder: "e.g. 499.00", Required: true, IsAmount: true, Input: createInput("499.00", false)},
			{Key: "currency", Label: "Currency", Placeholder: "INR", Required: true, Input: createInput("INR", false)},
			{Key: "description", Label: "Item Description", Placeholder: "e.g. Monthly server hosting package", Input: createInput("Standard server hosting", false)},
		}
	}
	if id == "invoices-items-update" {
		return []FormField{
			{Key: "id", Label: "Item ID", Placeholder: "e.g. item_xxx", Required: true, Input: createInput("item_xxx", false)},
			{Key: "name", Label: "Item Name", Placeholder: "e.g. Cloud Hosting Pro", Input: createInput("", false)},
			{Key: "amount", Label: "Unit Price (₹ INR)", Placeholder: "e.g. 799.00", IsAmount: true, Input: createInput("", false)},
			{Key: "description", Label: "Description", Placeholder: "e.g. Upgraded package", Input: createInput("", false)},
		}
	}
	if id == "invoices-items-fetch" || id == "invoices-items-delete" {
		return []FormField{
			{Key: "id", Label: "Item ID", Placeholder: "e.g. item_xxx", Required: true, Input: createInput("item_xxx", false)},
		}
	}

	// 7. QR CODES MODULE
	if id == "qr-codes-create" {
		return []FormField{
			{Key: "name", Label: "QR Code Name", Placeholder: "e.g. Store Front Counter QR", Required: true, Input: createInput("Store Front QR", false)},
			{Key: "usage", Label: "Usage Type", Placeholder: "single_use or multiple_use", Required: true, Input: createInput("single_use", false)},
			{Key: "amount", Label: "Payment Amount (₹ INR)", Placeholder: "e.g. 100.00 (leave blank for dynamic)", IsAmount: true, Input: createInput("100.00", false)},
			{Key: "description", Label: "QR Description", Placeholder: "e.g. Store checkout counter QR", Input: createInput("Store counter QR", false)},
		}
	}
	if id == "qr-codes-update" {
		return []FormField{
			{Key: "id", Label: "QR ID", Placeholder: "e.g. qr_xxx", Required: true, Input: createInput("qr_xxx", false)},
			{Key: "notes", Label: "Notes / Metadata", Placeholder: "e.g. Updated QR notes", Input: createInput("", false)},
		}
	}
	if id == "qr-codes-fetch" || id == "qr-codes-payments" || id == "qr-codes-close" {
		return []FormField{
			{Key: "id", Label: "QR ID", Placeholder: "e.g. qr_xxx", Required: true, Input: createInput("qr_xxx", false)},
		}
	}

	// 8. SUBSCRIPTIONS MODULE
	if id == "subs-create" {
		return []FormField{
			{Key: "plan_id", Label: "Plan ID", Placeholder: "e.g. plan_xxx", Required: true, Input: createInput("plan_xxx", false)},
			{Key: "total_count", Label: "Total Billing Cycles", Placeholder: "e.g. 12", Required: true, Input: createInput("12", false)},
			{Key: "quantity", Label: "Quantity", Placeholder: "e.g. 1", Input: createInput("1", false)},
			{Key: "notes", Label: "Notes / Metadata", Placeholder: "e.g. SaaS Pro Membership", Input: createInput("SaaS Pro Plan", false)},
		}
	}
	if id == "subs-update" {
		return []FormField{
			{Key: "id", Label: "Subscription ID", Placeholder: "e.g. sub_xxx", Required: true, Input: createInput("sub_xxx", false)},
			{Key: "plan_id", Label: "New Plan ID", Placeholder: "e.g. plan_xxx", Input: createInput("", false)},
			{Key: "quantity", Label: "Quantity", Placeholder: "e.g. 2", Input: createInput("", false)},
		}
	}
	if id == "subs-plans-create" {
		return []FormField{
			{Key: "period", Label: "Period", Placeholder: "weekly / monthly / yearly", Required: true, Input: createInput("monthly", false)},
			{Key: "interval", Label: "Interval", Placeholder: "e.g. 1", Required: true, Input: createInput("1", false)},
			{Key: "amount", Label: "Plan Price (₹ INR)", Placeholder: "e.g. 999.00", Required: true, IsAmount: true, Input: createInput("999.00", false)},
			{Key: "name", Label: "Plan Name", Placeholder: "e.g. Silver Monthly Plan", Required: true, Input: createInput("Silver Monthly Plan", false)},
			{Key: "description", Label: "Plan Description", Placeholder: "e.g. Recurring monthly membership", Input: createInput("Monthly membership", false)},
		}
	}
	if id == "subs-plans-fetch" {
		return []FormField{
			{Key: "id", Label: "Plan ID", Placeholder: "e.g. plan_xxx", Required: true, Input: createInput("plan_xxx", false)},
		}
	}
	if strings.Contains(id, "subs-") {
		return []FormField{
			{Key: "id", Label: "Subscription ID", Placeholder: "e.g. sub_xxx", Required: true, Input: createInput("sub_xxx", false)},
		}
	}

	// Fallback for Fetch / Single ID operations
	if strings.Contains(id, "fetch") || strings.Contains(id, "delete") || strings.Contains(id, "close") || strings.Contains(id, "accept") || strings.Contains(id, "contest") {
		resName := "Entity"
		if strings.Contains(id, "order") {
			resName = "Order"
		} else if strings.Contains(id, "payment") {
			resName = "Payment"
		} else if strings.Contains(id, "refund") {
			resName = "Refund"
		} else if strings.Contains(id, "customer") {
			resName = "Customer"
		} else if strings.Contains(id, "invoice") {
			resName = "Invoice"
		} else if strings.Contains(id, "dispute") {
			resName = "Dispute"
		} else if strings.Contains(id, "document") || strings.Contains(id, "doc") {
			resName = "Document"
		}
		return []FormField{
			{Key: "id", Label: fmt.Sprintf("%s ID", resName), Placeholder: fmt.Sprintf("e.g. %s_xxx", strings.ToLower(resName[:4])), Required: true, Input: createInput(fmt.Sprintf("%s_xxx", strings.ToLower(resName[:4])), false)},
		}
	}

	// Generic Form with ID and Notes
	return []FormField{
		{Key: "id", Label: "Entity ID", Placeholder: "e.g. id_xxx", Required: strings.Contains(action.APIPath, "{"), Input: createInput("id_xxx", false)},
		{Key: "notes", Label: "Notes / Payload", Placeholder: "e.g. optional metadata", Input: createInput("", false)},
	}
}

func (fv FormViewScreen) View() string {
	var sections []string

	// 1. Header
	headerView := components.RenderHeader(fv.state, fv.width)
	sections = append(sections, headerView)

	// 2. Toast
	toastView := components.RenderToast(fv.state)
	if toastView != "" {
		sections = append(sections, toastView)
	}

	// 3. Footer Keys
	var keys []components.KeyHelp
	if fv.formState == FormStateSuccess {
		keys = []components.KeyHelp{
			{Key: "Enter", Desc: "Inspect Full Details (Screen 4)"},
			{Key: "a", Desc: "AI Assist"},
			{Key: "r", Desc: "New Entry"},
			{Key: "Esc", Desc: "Back to Actions"},
			{Key: "q", Desc: "Quit"},
		}
	} else if fv.formState == FormStateError {
		keys = []components.KeyHelp{
			{Key: "Enter / r", Desc: "Retry"},
			{Key: "a", Desc: "AI Assist"},
			{Key: "c", Desc: "Config"},
			{Key: "Esc", Desc: "Back"},
			{Key: "q", Desc: "Quit"},
		}
	} else {
		keys = []components.KeyHelp{
			{Key: "Tab / ↓", Desc: "Next Field"},
			{Key: "Shift+Tab", Desc: "Prev Field"},
			{Key: "Enter", Desc: "Submit Action"},
			{Key: "Esc", Desc: "Back to Actions"},
			{Key: "Ctrl+C", Desc: "Quit"},
		}
	}
	footerView := components.RenderFooter(fv.state, fv.width, keys)

	// Height math
	headerHeight := lipgloss.Height(headerView)
	footerHeight := lipgloss.Height(footerView)
	toastHeight := 0
	if toastView != "" {
		toastHeight = lipgloss.Height(toastView)
	}

	bodyHeight := fv.height - headerHeight - footerHeight - toastHeight
	if bodyHeight < 10 {
		bodyHeight = 10
	}

	// 4. Main Body Content
	var mainContent string
	title := styles.TitleStyle.PaddingLeft(1).Render(fmt.Sprintf("%s (%s %s)", fv.action.Title, fv.action.HTTPMethod, fv.action.CLICommand))

	switch fv.formState {
	case FormStateSubmitting:
		loadingCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorSecondary).
			Background(styles.ColorCardBg).
			Padding(3, 4).
			Width(fv.width - 6).
			Render(fmt.Sprintf(
				"%s  Submitting request to Razorpay Cloud API (%s %s)...\n\nPlease wait while your transaction is processed in %s.",
				fv.spinner.View(), fv.action.HTTPMethod, fv.action.APIPath, fv.state.ActiveMode,
			))
		mainContent = "\n" + title + "\n\n" + loadingCard

	case FormStateSuccess:
		resID := fmt.Sprintf("%v", fv.resultData["id"])
		if resID == "<nil>" || resID == "" {
			resID = "Operation Successful"
		}
		successCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorSuccess).
			Background(styles.ColorCardBg).
			Padding(2, 4).
			Width(fv.width - 6).
			Render(fmt.Sprintf(
				"🎉  Action Successfully Executed!\n\n"+
					"• Status       : %s\n"+
					"• Resource ID  : %s\n"+
					"• Environment  : %s\n\n"+
					"👉 Press [Enter] to inspect full JSON details in Screen 4\n"+
					"👉 Press [r] to create another entry\n"+
					"👉 Press [Esc] to return to Actions Menu",
				styles.StatusPillPaid, resID, fv.state.ActiveMode,
			))
		mainContent = "\n" + title + "\n\n" + successCard

	case FormStateError:
		errCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorError).
			Background(styles.ColorCardBg).
			Padding(2, 4).
			Width(fv.width - 6).
			Render(fmt.Sprintf(
				"⚠️  API Execution Failed:\n\n%s\n\n"+
					"👉 Press [Enter] or [r] to edit parameters and retry\n"+
					"👉 Press [c] to check API credentials\n"+
					"👉 Press [Esc] to go back",
				fv.resultErr,
			))
		mainContent = "\n" + title + "\n\n" + errCard

	case FormStateEditing:
		var formLines []string
		for i, field := range fv.fields {
			labelColor := styles.ColorTextMuted
			if i == fv.focusedIdx {
				labelColor = styles.ColorSecondary
			}
			reqBadge := ""
			if field.Required {
				reqBadge = lipgloss.NewStyle().Foreground(styles.ColorError).Render(" *")
			}
			label := lipgloss.NewStyle().
				Bold(i == fv.focusedIdx).
				Foreground(labelColor).
				Width(24).
				Render(field.Label + reqBadge + " :")

			formLines = append(formLines, lipgloss.JoinHorizontal(lipgloss.Center, label, "  ", field.Input.View()))
			formLines = append(formLines, "")
		}

		formCard := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorBorderFocus).
			Background(styles.ColorCardBg).
			Padding(1, 3).
			Width(fv.width - 6).
			Render(strings.Join(formLines, "\n"))

		submitHint := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			PaddingLeft(1).
			Render("💡 Use [Tab] / [Shift+Tab] to navigate fields • Press [Enter] to submit to Razorpay API")

		mainContent = "\n" + title + "\n\n" + formCard + "\n" + submitHint
	}

	bodyContainer := lipgloss.NewStyle().
		Height(bodyHeight).
		Width(fv.width - 2)

	sections = append(sections, bodyContainer.Render(mainContent))
	sections = append(sections, footerView)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func normalizePhoneNumber(phone string) string {
	cleaned := strings.TrimSpace(phone)
	if cleaned == "" {
		return ""
	}

	// Remove common spacing, hyphens, parentheses
	r := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "")
	cleaned = r.Replace(cleaned)

	// If already starts with '+', keep international format
	if strings.HasPrefix(cleaned, "+") {
		return cleaned
	}

	// If starts with '91' and has 12 digits, prepend '+'
	if strings.HasPrefix(cleaned, "91") && len(cleaned) == 12 {
		return "+" + cleaned
	}

	// If 10-digit standard mobile number, auto-prepend '+91'
	if len(cleaned) == 10 {
		return "+91" + cleaned
	}

	// If starts with '0' and has 11 digits (e.g. 09876543210)
	if strings.HasPrefix(cleaned, "0") && len(cleaned) == 11 {
		return "+91" + cleaned[1:]
	}

	return cleaned
}
