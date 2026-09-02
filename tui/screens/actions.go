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

type actionItem struct {
	action state.ActionItem
}

func (i actionItem) Title() string       { return i.action.Title }
func (i actionItem) Description() string { return i.action.Description }
func (i actionItem) FilterValue() string {
	return i.action.ID + " " + i.action.Title + " " + i.action.Description + " " + i.action.CLICommand
}

type actionDelegate struct{}

func (d actionDelegate) Height() int                             { return 2 }
func (d actionDelegate) Spacing() int                            { return 1 } // Generous space between action cards
func (d actionDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d actionDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(actionItem)
	if !ok {
		return
	}

	title := i.action.Title
	desc := i.action.Description
	cmdBadge := lipgloss.NewStyle().Foreground(styles.ColorSecondary).Render("(" + i.action.CLICommand + ")")

	isSelected := index == m.Index()

	if isSelected {
		header := lipgloss.JoinHorizontal(
			lipgloss.Center,
			styles.ItemCardSelected.Render(fmt.Sprintf("▶  %-28s", title)),
			" ",
			cmdBadge,
		)
		descLine := styles.ItemDescSelected.Render(desc)
		fmt.Fprintf(w, "%s\n%s", header, descLine)
	} else {
		header := lipgloss.JoinHorizontal(
			lipgloss.Center,
			styles.ItemCardNormal.Render(fmt.Sprintf("   %-28s", title)),
			" ",
			cmdBadge,
		)
		descLine := styles.ItemDescNormal.Render(desc)
		fmt.Fprintf(w, "%s\n%s", header, descLine)
	}
}

type ActionsScreen struct {
	list        list.Model
	state       *state.SessionState
	moduleID    string
	initialized bool
}

func GetActionsForModule(modID string) []state.ActionItem {
	switch modID {
	case "orders":
		return []state.ActionItem{
			{ID: "orders-list", Title: "📋 List Orders", Description: "Fetch and list all orders with pagination & status filters", CLICommand: "razorpay orders list", HTTPMethod: "GET", APIPath: "/v1/orders"},
			{ID: "orders-create", Title: "➕ Create Order", Description: "Create a new order with amount, currency, receipt & notes", CLICommand: "razorpay orders create", HTTPMethod: "POST", APIPath: "/v1/orders", IsForm: true},
			{ID: "orders-fetch", Title: "🔍 Fetch Order", Description: "Retrieve specific order details by Order ID (order_xxx)", CLICommand: "razorpay orders fetch", HTTPMethod: "GET", APIPath: "/v1/orders/{id}", IsForm: true},
			{ID: "orders-payments", Title: "💳 Order Payments", Description: "List all payments associated with an order ID", CLICommand: "razorpay orders payments", HTTPMethod: "GET", APIPath: "/v1/orders/{id}/payments", IsForm: true},
			{ID: "orders-update", Title: "✏️  Update Order", Description: "Update notes and metadata on an existing order", CLICommand: "razorpay orders update", HTTPMethod: "PATCH", APIPath: "/v1/orders/{id}", IsForm: true},
		}

	case "payments":
		return []state.ActionItem{
			{ID: "payments-list", Title: "📋 List Payments", Description: "Fetch and list all payment transactions", CLICommand: "razorpay payments list", HTTPMethod: "GET", APIPath: "/v1/payments"},
			{ID: "payments-fetch", Title: "🔍 Fetch Payment", Description: "Retrieve payment details by Payment ID (pay_xxx)", CLICommand: "razorpay payments fetch", HTTPMethod: "GET", APIPath: "/v1/payments/{id}", IsForm: true},
			{ID: "payments-capture", Title: "💰 Capture Payment", Description: "Capture an authorized payment with amount & currency", CLICommand: "razorpay payments capture", HTTPMethod: "POST", APIPath: "/v1/payments/{id}/capture", IsForm: true},
			{ID: "payments-card", Title: "💳 Card Details", Description: "Fetch card information associated with a payment", CLICommand: "razorpay payments card", HTTPMethod: "GET", APIPath: "/v1/payments/{id}/card", IsForm: true},
			{ID: "payments-update", Title: "✏️  Update Payment", Description: "Update internal notes on a payment record", CLICommand: "razorpay payments update", HTTPMethod: "PATCH", APIPath: "/v1/payments/{id}", IsForm: true},
			{ID: "payments-downtime-list", Title: "⚠️  Payment Downtimes", Description: "List active downtime across issuers, banks & methods", CLICommand: "razorpay payments downtime list", HTTPMethod: "GET", APIPath: "/v1/payments/downtimes"},
			{ID: "payments-downtime-fetch", Title: "🔍 Fetch Downtime", Description: "Retrieve specific downtime report by ID", CLICommand: "razorpay payments downtime fetch", HTTPMethod: "GET", APIPath: "/v1/payments/downtimes/{id}", IsForm: true},
		}

	case "refunds":
		return []state.ActionItem{
			{ID: "refunds-list", Title: "📋 List Refunds", Description: "List all processed refunds with filter options", CLICommand: "razorpay refunds list", HTTPMethod: "GET", APIPath: "/v1/refunds"},
			{ID: "refunds-create", Title: "💸 Create Refund", Description: "Issue a full or partial refund on a captured payment", CLICommand: "razorpay refunds create", HTTPMethod: "POST", APIPath: "/v1/payments/{id}/refund", IsForm: true},
			{ID: "refunds-fetch", Title: "🔍 Fetch Refund", Description: "Fetch specific refund details by Refund ID (rfr_xxx)", CLICommand: "razorpay refunds fetch", HTTPMethod: "GET", APIPath: "/v1/refunds/{id}", IsForm: true},
			{ID: "refunds-update", Title: "✏️  Update Refund", Description: "Update notes on a refund record", CLICommand: "razorpay refunds update", HTTPMethod: "PATCH", APIPath: "/v1/refunds/{id}", IsForm: true},
			{ID: "refunds-payment-refunds", Title: "📜 Payment Refunds", Description: "Fetch all refunds issued for a payment ID", CLICommand: "razorpay refunds payment-refunds", HTTPMethod: "GET", APIPath: "/v1/payments/{id}/refunds", IsForm: true},
			{ID: "refunds-payment-refund", Title: "🔍 Specific Payment Refund", Description: "Fetch a specific refund ID under a payment", CLICommand: "razorpay refunds payment-refund", HTTPMethod: "GET", APIPath: "/v1/payments/{id}/refunds/{refund_id}", IsForm: true},
		}

	case "customers":
		return []state.ActionItem{
			{ID: "customers-list", Title: "📋 List Customers", Description: "List all customer accounts", CLICommand: "razorpay customers list", HTTPMethod: "GET", APIPath: "/v1/customers"},
			{ID: "customers-create", Title: "👤 Create Customer", Description: "Create a new customer profile (name, email, contact)", CLICommand: "razorpay customers create", HTTPMethod: "POST", APIPath: "/v1/customers", IsForm: true},
			{ID: "customers-fetch", Title: "🔍 Fetch Customer", Description: "Fetch customer profile by Customer ID (cust_xxx)", CLICommand: "razorpay customers fetch", HTTPMethod: "GET", APIPath: "/v1/customers/{id}", IsForm: true},
			{ID: "customers-update", Title: "✏️  Update Customer", Description: "Update customer name, email or contact details", CLICommand: "razorpay customers update", HTTPMethod: "PUT", APIPath: "/v1/customers/{id}", IsForm: true},
		}

	case "invoices":
		return []state.ActionItem{
			{ID: "invoices-list", Title: "📋 List Invoices", Description: "List all customer invoices", CLICommand: "razorpay invoices list", HTTPMethod: "GET", APIPath: "/v1/invoices"},
			{ID: "invoices-create", Title: "🧾 Create Invoice", Description: "Create a new invoice for customer billing", CLICommand: "razorpay invoices create", HTTPMethod: "POST", APIPath: "/v1/invoices", IsForm: true},
			{ID: "invoices-fetch", Title: "🔍 Fetch Invoice", Description: "Retrieve invoice details by Invoice ID (inv_xxx)", CLICommand: "razorpay invoices fetch", HTTPMethod: "GET", APIPath: "/v1/invoices/{id}", IsForm: true},
			{ID: "invoices-update", Title: "✏️  Update Invoice", Description: "Update details of a draft invoice", CLICommand: "razorpay invoices update", HTTPMethod: "PATCH", APIPath: "/v1/invoices/{id}", IsForm: true},
			{ID: "invoices-issue", Title: "📨 Issue Invoice", Description: "Issue and finalize a draft invoice", CLICommand: "razorpay invoices issue", HTTPMethod: "POST", APIPath: "/v1/invoices/{id}/issue", IsForm: true},
			{ID: "invoices-notify", Title: "🔔 Notify Customer", Description: "Resend invoice notification via SMS / Email", CLICommand: "razorpay invoices notify", HTTPMethod: "POST", APIPath: "/v1/invoices/{id}/notify", IsForm: true},
			{ID: "invoices-cancel", Title: "🚫 Cancel Invoice", Description: "Cancel an issued unpaid invoice", CLICommand: "razorpay invoices cancel", HTTPMethod: "POST", APIPath: "/v1/invoices/{id}/cancel", IsForm: true},
			{ID: "invoices-delete", Title: "🗑️  Delete Invoice", Description: "Delete an unissued draft invoice", CLICommand: "razorpay invoices delete", HTTPMethod: "DELETE", APIPath: "/v1/invoices/{id}", IsForm: true},
			{ID: "invoices-items-list", Title: "📦 List Line Items", Description: "List all invoice line items in inventory", CLICommand: "razorpay invoices items list", HTTPMethod: "GET", APIPath: "/v1/items"},
			{ID: "invoices-items-create", Title: "➕ Create Line Item", Description: "Add a new line item to inventory", CLICommand: "razorpay invoices items create", HTTPMethod: "POST", APIPath: "/v1/items", IsForm: true},
			{ID: "invoices-items-fetch", Title: "🔍 Fetch Line Item", Description: "Fetch line item details by Item ID (item_xxx)", CLICommand: "razorpay invoices items fetch", HTTPMethod: "GET", APIPath: "/v1/items/{id}", IsForm: true},
			{ID: "invoices-items-update", Title: "✏️  Update Line Item", Description: "Update line item price, name or description", CLICommand: "razorpay invoices items update", HTTPMethod: "PATCH", APIPath: "/v1/items/{id}", IsForm: true},
			{ID: "invoices-items-delete", Title: "🗑️  Delete Line Item", Description: "Delete an item from inventory", CLICommand: "razorpay invoices items delete", HTTPMethod: "DELETE", APIPath: "/v1/items/{id}", IsForm: true},
		}

	case "payment-links":
		return []state.ActionItem{
			{ID: "payment-links-list", Title: "📋 List Payment Links", Description: "List all generated standard payment links", CLICommand: "razorpay payment-links list", HTTPMethod: "GET", APIPath: "/v1/payment_links"},
			{ID: "payment-links-create", Title: "🔗 Create Payment Link", Description: "Generate a shareable payment link", CLICommand: "razorpay payment-links create", HTTPMethod: "POST", APIPath: "/v1/payment_links", IsForm: true},
			{ID: "payment-links-fetch", Title: "🔍 Fetch Payment Link", Description: "Retrieve payment link details by ID (plink_xxx)", CLICommand: "razorpay payment-links fetch", HTTPMethod: "GET", APIPath: "/v1/payment_links/{id}", IsForm: true},
			{ID: "payment-links-update", Title: "✏️  Update Payment Link", Description: "Update expiry or reminder preferences", CLICommand: "razorpay payment-links update", HTTPMethod: "PATCH", APIPath: "/v1/payment_links/{id}", IsForm: true},
			{ID: "payment-links-cancel", Title: "🚫 Cancel Payment Link", Description: "Cancel an active payment link", CLICommand: "razorpay payment-links cancel", HTTPMethod: "POST", APIPath: "/v1/payment_links/{id}/cancel", IsForm: true},
			{ID: "payment-links-notify", Title: "🔔 Resend Notification", Description: "Send payment link SMS/Email reminder", CLICommand: "razorpay payment-links notify", HTTPMethod: "POST", APIPath: "/v1/payment_links/{id}/notify_by/{medium}", IsForm: true},
		}

	case "qr-codes":
		return []state.ActionItem{
			{ID: "qr-codes-list", Title: "📋 List QR Codes", Description: "List all generated BharatQR codes", CLICommand: "razorpay qr-codes list", HTTPMethod: "GET", APIPath: "/v1/payments/qr_codes"},
			{ID: "qr-codes-create", Title: "📱 Create QR Code", Description: "Generate a new static or dynamic QR code", CLICommand: "razorpay qr-codes create", HTTPMethod: "POST", APIPath: "/v1/payments/qr_codes", IsForm: true},
			{ID: "qr-codes-fetch", Title: "🔍 Fetch QR Code", Description: "Fetch QR code details and image URL by ID (qr_xxx)", CLICommand: "razorpay qr-codes fetch", HTTPMethod: "GET", APIPath: "/v1/payments/qr_codes/{id}", IsForm: true},
			{ID: "qr-codes-payments", Title: "💳 QR Payments", Description: "List payments received through this QR code", CLICommand: "razorpay qr-codes payments", HTTPMethod: "GET", APIPath: "/v1/payments/qr_codes/{id}/payments", IsForm: true},
			{ID: "qr-codes-update", Title: "✏️  Update QR Code", Description: "Update customer details or notes on a QR code", CLICommand: "razorpay qr-codes update", HTTPMethod: "PATCH", APIPath: "/v1/payments/qr_codes/{id}", IsForm: true},
			{ID: "qr-codes-close", Title: "🔒 Close QR Code", Description: "Deactivate and close an active QR code", CLICommand: "razorpay qr-codes close", HTTPMethod: "POST", APIPath: "/v1/payments/qr_codes/{id}/close", IsForm: true},
		}

	case "subscriptions":
		return []state.ActionItem{
			{ID: "subs-list", Title: "📋 List Subscriptions", Description: "List all customer subscriptions", CLICommand: "razorpay subscriptions list", HTTPMethod: "GET", APIPath: "/v1/subscriptions"},
			{ID: "subs-create", Title: "🔄 Create Subscription", Description: "Create a recurring billing subscription", CLICommand: "razorpay subscriptions create", HTTPMethod: "POST", APIPath: "/v1/subscriptions", IsForm: true},
			{ID: "subs-fetch", Title: "🔍 Fetch Subscription", Description: "Retrieve subscription by ID (sub_xxx)", CLICommand: "razorpay subscriptions fetch", HTTPMethod: "GET", APIPath: "/v1/subscriptions/{id}", IsForm: true},
			{ID: "subs-update", Title: "✏️  Update Subscription", Description: "Change plan or schedule subscription updates", CLICommand: "razorpay subscriptions update", HTTPMethod: "PATCH", APIPath: "/v1/subscriptions/{id}", IsForm: true},
			{ID: "subs-pause", Title: "⏸️  Pause Subscription", Description: "Temporarily pause recurring charges", CLICommand: "razorpay subscriptions pause", HTTPMethod: "POST", APIPath: "/v1/subscriptions/{id}/pause", IsForm: true},
			{ID: "subs-resume", Title: "▶️  Resume Subscription", Description: "Resume a paused subscription", CLICommand: "razorpay subscriptions resume", HTTPMethod: "POST", APIPath: "/v1/subscriptions/{id}/resume", IsForm: true},
			{ID: "subs-cancel", Title: "🚫 Cancel Subscription", Description: "Cancel an active subscription", CLICommand: "razorpay subscriptions cancel", HTTPMethod: "POST", APIPath: "/v1/subscriptions/{id}/cancel", IsForm: true},
			{ID: "subs-pending-update", Title: "⏳ Pending Update", Description: "Check scheduled upcoming updates", CLICommand: "razorpay subscriptions pending-update", HTTPMethod: "GET", APIPath: "/v1/subscriptions/{id}/retrieve_scheduled_changes", IsForm: true},
			{ID: "subs-cancel-update", Title: "✖️  Cancel Update", Description: "Discard scheduled subscription update", CLICommand: "razorpay subscriptions cancel-update", HTTPMethod: "POST", APIPath: "/v1/subscriptions/{id}/cancel_scheduled_changes", IsForm: true},
			{ID: "subs-invoices", Title: "🧾 Subscription Invoices", Description: "List all billing cycle invoices for subscription", CLICommand: "razorpay subscriptions invoices", HTTPMethod: "GET", APIPath: "/v1/invoices?subscription_id={id}", IsForm: true},
			{ID: "subs-delete-offer", Title: "🏷️  Delete Offer", Description: "Unlink an offer from subscription", CLICommand: "razorpay subscriptions delete-offer", HTTPMethod: "DELETE", APIPath: "/v1/subscriptions/{id}/offer", IsForm: true},
			{ID: "subs-plans-list", Title: "📜 List Plans", Description: "List all subscription billing plans", CLICommand: "razorpay subscriptions plans list", HTTPMethod: "GET", APIPath: "/v1/plans"},
			{ID: "subs-plans-create", Title: "➕ Create Plan", Description: "Create a new billing plan (weekly/monthly/yearly)", CLICommand: "razorpay subscriptions plans create", HTTPMethod: "POST", APIPath: "/v1/plans", IsForm: true},
			{ID: "subs-plans-fetch", Title: "🔍 Fetch Plan", Description: "Fetch plan details by Plan ID (plan_xxx)", CLICommand: "razorpay subscriptions plans fetch", HTTPMethod: "GET", APIPath: "/v1/plans/{id}", IsForm: true},
		}

	case "route":
		return []state.ActionItem{
			{ID: "route-accounts-fetch", Title: "👥 Fetch Linked Account", Description: "Fetch linked merchant account by ID", CLICommand: "razorpay route accounts fetch", HTTPMethod: "GET", APIPath: "/v1/accounts/{id}", IsForm: true},
			{ID: "route-accounts-create", Title: "➕ Create Linked Account", Description: "Onboard a new seller / partner account", CLICommand: "razorpay route accounts create", HTTPMethod: "POST", APIPath: "/v1/accounts", IsForm: true},
			{ID: "route-accounts-update", Title: "✏️  Update Linked Account", Description: "Update linked account details", CLICommand: "razorpay route accounts update", HTTPMethod: "PATCH", APIPath: "/v1/accounts/{id}", IsForm: true},
			{ID: "route-accounts-payments", Title: "💳 Account Payments", Description: "List payments credited to linked account", CLICommand: "razorpay route accounts payments", HTTPMethod: "GET", APIPath: "/v1/accounts/{id}/payments", IsForm: true},
			{ID: "route-product-fetch", Title: "📦 Product Config", Description: "Fetch Route product configuration", CLICommand: "razorpay route accounts product-fetch", HTTPMethod: "GET", APIPath: "/v1/accounts/{id}/products/{product_id}", IsForm: true},
			{ID: "route-product-request", Title: "📝 Request Product", Description: "Request Route product activation", CLICommand: "razorpay route accounts product-request", HTTPMethod: "POST", APIPath: "/v1/accounts/{id}/products", IsForm: true},
			{ID: "route-product-update", Title: "⚙️ Update Product Config", Description: "Update settlement & product settings", CLICommand: "razorpay route accounts product-update", HTTPMethod: "PATCH", APIPath: "/v1/accounts/{id}/products/{product_id}", IsForm: true},
			{ID: "route-stakeholder-create", Title: "👤 Add Stakeholder", Description: "Add business stakeholder to account", CLICommand: "razorpay route accounts stakeholder-create", HTTPMethod: "POST", APIPath: "/v1/accounts/{id}/stakeholders", IsForm: true},
			{ID: "route-stakeholder-update", Title: "✏️  Update Stakeholder", Description: "Update stakeholder KYC details", CLICommand: "razorpay route accounts stakeholder-update", HTTPMethod: "PATCH", APIPath: "/v1/accounts/{id}/stakeholders/{stakeholder_id}", IsForm: true},
			{ID: "route-transfers-list", Title: "🔀 List Transfers", Description: "List all marketplace split transfers", CLICommand: "razorpay route transfers list", HTTPMethod: "GET", APIPath: "/v1/transfers"},
			{ID: "route-transfers-create", Title: "💸 Direct Transfer", Description: "Create direct transfer to linked account", CLICommand: "razorpay route transfers create", HTTPMethod: "POST", APIPath: "/v1/transfers", IsForm: true},
			{ID: "route-transfers-create-order", Title: "🛒 Order with Transfers", Description: "Create order with split transfers", CLICommand: "razorpay route transfers create-from-order", HTTPMethod: "POST", APIPath: "/v1/orders", IsForm: true},
			{ID: "route-transfers-create-payment", Title: "💳 Payment Transfers", Description: "Split captured payment into transfers", CLICommand: "razorpay route transfers create-from-payment", HTTPMethod: "POST", APIPath: "/v1/payments/{id}/transfers", IsForm: true},
			{ID: "route-transfers-fetch", Title: "🔍 Fetch Transfer", Description: "Retrieve transfer by Transfer ID (trf_xxx)", CLICommand: "razorpay route transfers fetch", HTTPMethod: "GET", APIPath: "/v1/transfers/{id}", IsForm: true},
			{ID: "route-transfers-fetch-order", Title: "📦 Order Transfers", Description: "List transfers for an order ID", CLICommand: "razorpay route transfers fetch-by-order", HTTPMethod: "GET", APIPath: "/v1/orders/{id}/transfers", IsForm: true},
			{ID: "route-transfers-fetch-payment", Title: "🧾 Payment Transfers", Description: "List transfers for a payment ID", CLICommand: "razorpay route transfers fetch-by-payment", HTTPMethod: "GET", APIPath: "/v1/payments/{id}/transfers", IsForm: true},
			{ID: "route-transfers-update", Title: "🔒 Update Settlement Hold", Description: "Hold or release transfer settlement", CLICommand: "razorpay route transfers update", HTTPMethod: "PATCH", APIPath: "/v1/transfers/{id}", IsForm: true},
			{ID: "route-transfers-reverse", Title: "↩️  Reverse Transfer", Description: "Reverse transfer funds back to merchant", CLICommand: "razorpay route transfers reverse", HTTPMethod: "POST", APIPath: "/v1/transfers/{id}/reversals", IsForm: true},
			{ID: "route-transfers-reversals", Title: "📜 Transfer Reversals", Description: "List all reversals for a transfer", CLICommand: "razorpay route transfers reversals", HTTPMethod: "GET", APIPath: "/v1/transfers/{id}/reversals", IsForm: true},
			{ID: "route-refund-with-reversal", Title: "🔄 Refund with Reversal", Description: "Refund payment & reverse transfers", CLICommand: "razorpay route refund-with-reversal", HTTPMethod: "POST", APIPath: "/v1/payments/{id}/refund", IsForm: true},
		}

	case "smart-collect":
		return []state.ActionItem{
			{ID: "sc-list", Title: "📋 List Virtual Accounts", Description: "List all Smart Collect virtual accounts", CLICommand: "razorpay smart-collect list", HTTPMethod: "GET", APIPath: "/v1/virtual_accounts"},
			{ID: "sc-create", Title: "🏢 Create Virtual Account", Description: "Create customer identifier bank/UPI account", CLICommand: "razorpay smart-collect create", HTTPMethod: "POST", APIPath: "/v1/virtual_accounts", IsForm: true},
			{ID: "sc-fetch", Title: "🔍 Fetch Virtual Account", Description: "Fetch virtual account by ID (va_xxx)", CLICommand: "razorpay smart-collect fetch", HTTPMethod: "GET", APIPath: "/v1/virtual_accounts/{id}", IsForm: true},
			{ID: "sc-update", Title: "✏️  Update Virtual Account", Description: "Update virtual account descriptor or notes", CLICommand: "razorpay smart-collect update", HTTPMethod: "PATCH", APIPath: "/v1/virtual_accounts/{id}", IsForm: true},
			{ID: "sc-close", Title: "🔒 Close Virtual Account", Description: "Close and deactivate virtual account", CLICommand: "razorpay smart-collect close", HTTPMethod: "POST", APIPath: "/v1/virtual_accounts/{id}/close", IsForm: true},
			{ID: "sc-payments", Title: "💳 Virtual Account Payments", Description: "List payments received on account", CLICommand: "razorpay smart-collect payments", HTTPMethod: "GET", APIPath: "/v1/virtual_accounts/{id}/payments", IsForm: true},
			{ID: "sc-add-receiver", Title: "➕ Add Bank/VPA Receiver", Description: "Attach custom receiver to account", CLICommand: "razorpay smart-collect add-receiver", HTTPMethod: "POST", APIPath: "/v1/virtual_accounts/{id}/receivers", IsForm: true},
			{ID: "sc-fetch-bank-transfer", Title: "🏦 Bank Transfer Details", Description: "Fetch NEFT/RTGS/IMPS details for payment", CLICommand: "razorpay smart-collect fetch-by-bank-transfer", HTTPMethod: "GET", APIPath: "/v1/payments/{id}/bank_transfer", IsForm: true},
			{ID: "sc-fetch-utr", Title: "🔍 Fetch by UTR Number", Description: "Lookup payment via banking UTR number", CLICommand: "razorpay smart-collect fetch-by-utr", HTTPMethod: "GET", APIPath: "/v1/payments?utr={utr}", IsForm: true},
			{ID: "sc-fetch-upi", Title: "📱 UPI Transfer Details", Description: "Fetch Smart Collect 2 UPI details for payment", CLICommand: "razorpay smart-collect fetch-upi-payment", HTTPMethod: "GET", APIPath: "/v1/payments/{id}/upi_transfer", IsForm: true},
			{ID: "sc-tpv-create", Title: "🛡️ Create TPV Account", Description: "Create third-party validated account", CLICommand: "razorpay smart-collect tpv-create", HTTPMethod: "POST", APIPath: "/v1/virtual_accounts", IsForm: true},
			{ID: "sc-tpv-add-payer", Title: "➕ Add TPV Payer", Description: "Register allowed payer bank account", CLICommand: "razorpay smart-collect tpv-add-payer", HTTPMethod: "POST", APIPath: "/v1/virtual_accounts/{id}/allowed_payers", IsForm: true},
			{ID: "sc-tpv-delete-payer", Title: "🗑️  Delete TPV Payer", Description: "Remove allowed payer bank account", CLICommand: "razorpay smart-collect tpv-delete-payer", HTTPMethod: "DELETE", APIPath: "/v1/virtual_accounts/{id}/allowed_payers/{payer_id}", IsForm: true},
		}

	case "settlements":
		return []state.ActionItem{
			{ID: "settle-list", Title: "📋 List Settlements", Description: "List all bank payouts and settlements", CLICommand: "razorpay settlements list", HTTPMethod: "GET", APIPath: "/v1/settlements"},
			{ID: "settle-fetch", Title: "🔍 Fetch Settlement", Description: "Fetch settlement details by Settlement ID (setl_xxx)", CLICommand: "razorpay settlements fetch", HTTPMethod: "GET", APIPath: "/v1/settlements/{id}", IsForm: true},
			{ID: "settle-recon", Title: "📊 Settlement Recon", Description: "Fetch settlement reconciliation report", CLICommand: "razorpay settlements recon", HTTPMethod: "GET", APIPath: "/v1/settlements/recon/combined"},
			{ID: "settle-instant-list", Title: "⚡ List Instant Payouts", Description: "List all on-demand instant settlements", CLICommand: "razorpay settlements instant-list", HTTPMethod: "GET", APIPath: "/v1/settlements/ondemand"},
			{ID: "settle-instant-fetch", Title: "🔍 Fetch Instant Payout", Description: "Fetch instant settlement status by ID", CLICommand: "razorpay settlements instant-fetch", HTTPMethod: "GET", APIPath: "/v1/settlements/ondemand/{id}", IsForm: true},
			{ID: "settle-instant-create", Title: "💸 Request Instant Payout", Description: "Trigger immediate on-demand payout", CLICommand: "razorpay settlements instant-create", HTTPMethod: "POST", APIPath: "/v1/settlements/ondemand", IsForm: true},
		}

	case "disputes":
		return []state.ActionItem{
			{ID: "disputes-list", Title: "📋 List Disputes", Description: "List all customer chargebacks & disputes", CLICommand: "razorpay disputes list", HTTPMethod: "GET", APIPath: "/v1/disputes"},
			{ID: "disputes-fetch", Title: "🔍 Fetch Dispute", Description: "Fetch dispute details by Dispute ID (disp_xxx)", CLICommand: "razorpay disputes fetch", HTTPMethod: "GET", APIPath: "/v1/disputes/{id}", IsForm: true},
			{ID: "disputes-accept", Title: "✅ Accept Dispute", Description: "Accept chargeback liability and close dispute", CLICommand: "razorpay disputes accept", HTTPMethod: "POST", APIPath: "/v1/disputes/{id}/accept", IsForm: true},
			{ID: "disputes-contest", Title: "⚔️ Contest Dispute", Description: "Submit defense summary & proof documents", CLICommand: "razorpay disputes contest", HTTPMethod: "PATCH", APIPath: "/v1/disputes/{id}/contest", IsForm: true},
		}

	case "documents":
		return []state.ActionItem{
			{ID: "docs-create", Title: "📤 Upload Document", Description: "Upload KYC, dispute evidence or tax doc", CLICommand: "razorpay documents create", HTTPMethod: "POST", APIPath: "/v1/documents", IsForm: true},
			{ID: "docs-fetch", Title: "🔍 Fetch Document", Description: "Fetch document metadata by ID (doc_xxx)", CLICommand: "razorpay documents fetch", HTTPMethod: "GET", APIPath: "/v1/documents/{id}", IsForm: true},
			{ID: "docs-fetch-content", Title: "📥 Download Content", Description: "Download and view document file bytes", CLICommand: "razorpay documents fetch-content", HTTPMethod: "GET", APIPath: "/v1/documents/{id}/content", IsForm: true},
		}

	case "configure":
		return []state.ActionItem{
			{ID: "config-setup", Title: "⚙️ Configure Credentials", Description: "Configure API Key ID, Secret and Test/Live Mode", CLICommand: "razorpay configure", HTTPMethod: "LOCAL", APIPath: "config", IsForm: true},
		}

	default:
		return []state.ActionItem{}
	}
}

func NewActionsScreen(s *state.SessionState, moduleID string, width, height int) ActionsScreen {
	actions := GetActionsForModule(moduleID)
	items := make([]list.Item, len(actions))
	for idx, a := range actions {
		items[idx] = actionItem{action: a}
	}

	bodyHeight := height - 6
	if bodyHeight < 5 {
		bodyHeight = 5
	}

	l := list.New(items, actionDelegate{}, width-2, bodyHeight)
	l.Title = fmt.Sprintf("⚡ Available Actions for %s", s.SelectedModule.Title)
	l.Styles.Title = styles.TitleStyle
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)
	l.KeyMap.Quit.Unbind()
	l.KeyMap.ForceQuit.Unbind()
	l.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(styles.ColorSecondary)

	return ActionsScreen{
		list:        l,
		state:       s,
		moduleID:    moduleID,
		initialized: true,
	}
}

func (a *ActionsScreen) IsInitialized() bool {
	return a.initialized
}

func (a *ActionsScreen) IsFiltering() bool {
	return a.list.FilterState() == list.Filtering
}

func (a *ActionsScreen) ResetFilter() {
	a.list.ResetFilter()
}

func (a *ActionsScreen) SetSize(width, height int) {
	if a.initialized {
		bodyHeight := height - 6
		if a.state.Toast != nil && a.state.Toast.Message != "" {
			bodyHeight -= 3
		}
		if bodyHeight < 5 {
			bodyHeight = 5
		}
		a.list.SetSize(width-2, bodyHeight)
	}
}

func (a *ActionsScreen) Update(msg tea.Msg) (ActionsScreen, tea.Cmd) {
	var cmd tea.Cmd
	a.SetSize(a.state.Width, a.state.Height)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if a.list.FilterState() == list.Filtering {
			if msg.String() == "esc" || (msg.String() == "/" && a.list.FilterValue() == "") {
				a.list.ResetFilter()
				return *a, nil
			}
			break
		}
		switch msg.String() {
		case "enter":
			if sel, ok := a.list.SelectedItem().(actionItem); ok {
				a.state.SelectedAction = sel.action
				if sel.action.IsForm {
					// Enforce write/edit permission check
					allowed, errMsg := a.state.CanPerformAction(true)
					if !allowed {
						toastCmd := a.state.SetToast(errMsg, true)
						return *a, toastCmd
					}
					a.state.PushScreen(state.ScreenForm, sel.action.Title)
				} else {
					a.state.PushScreen(state.ScreenTable, sel.action.Title)
				}
			}
		}
	}

	a.list, cmd = a.list.Update(msg)
	return *a, cmd
}

func (a ActionsScreen) View() string {
	var sections []string

	// 1. Dynamic Header with Stepper & Location
	sections = append(sections, components.RenderHeader(a.state, a.state.Width))

	// 2. Toast (if any)
	toast := components.RenderToast(a.state)
	if toast != "" {
		sections = append(sections, toast)
	}

	// 3. Spacious Actions List
	sections = append(sections, a.list.View())

	// 4. Contextual Footer
	var keys []components.KeyHelp
	if a.list.FilterState() == list.Filtering {
		keys = []components.KeyHelp{
			{Key: "Enter", Desc: "Apply Filter"},
			{Key: "Esc", Desc: "Cancel Search"},
		}
	} else {
		keys = []components.KeyHelp{
			{Key: "↑/↓", Desc: "Navigate"},
			{Key: "Enter", Desc: "Open Action"},
			{Key: "/", Desc: "Search"},
			{Key: "Esc", Desc: "Back"},
			{Key: "q", Desc: "Quit"},
		}
	}
	sections = append(sections, components.RenderFooter(a.state, a.state.Width, keys))

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
