# Razorpay CLI & Terminal UI (TUI)

Command-line interface and interactive terminal user interface for the [Razorpay API](https://razorpay.com/docs/api/). Manage payments, orders, customers, invoices, refunds, settlements, disputes, payment links, QR codes, subscriptions, Route, and Smart Collect directly from your terminal.

---

## ⚡ Interactive Terminal UI (TUI)

Launch a full-screen, keyboard-driven interactive Terminal User Interface built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

```bash
# Launch interactive TUI
razorpay tui

# Launch directly in specific environment mode
razorpay tui --test       # Test Sandbox Mode
razorpay tui --live       # Live Production Mode
razorpay tui --read-only  # Safe Read-Only Mode (Locks all mutation operations)
```

### ✨ Key Capabilities & Features

* **🧭 14-Module Navigation Hub:** Seamlessly navigate across all 14 Razorpay services (Orders, Payments, Refunds, Customers, Invoices, Payment Links, QR Codes, Subscriptions, Settlements, Disputes, Smart Collect, Route, Documents, Configuration) with real-time fuzzy search (`/`).
* **📊 Live Interactive Data Tables:** Direct real-time REST API data fetching, dynamic resource-tailored columns, paise-to-Rupees currency conversion (`₹1,500.00`), relative timestamps, color-coded status pills (`● paid`, `● captured`, `▲ created`, `✖ failed`), and live refresh (`r`).
* **🔍 Detail View & Syntax-Highlighted JSON Inspector:** Dual-mode tab switcher (`Tab` / `v`) between a structured 2-column Summary Grid and a syntax-highlighted Raw JSON Tree with smooth scroll navigation.
* **📝 Dynamic Form Builder & API Executor:** Interactive forms for resource creation/updates with required field validators, automatic currency unit conversion, smart phone number normalization (auto `+91` prefix for 10 digits while preserving international codes), and real-time execution feedback.
* **⚙️ Dual Profile Credential Manager:** Manage independent Test Sandbox (`rzp_test_...`) and Live Production (`rzp_live_...`) credentials with safety isolation, masked secret inputs, and 1-key instant environment toggle (`c`).
* **🤖 Natural Language AI Assist (`a` / `ctrl+a`):** Type requests in plain English or Hindi (e.g. *"find failed payments from yesterday"*, *"create order for ₹500"*) and let Google Gemini / Claude Haiku generate the exact CLI command and flags with an interactive preview, smart amount conversion, and safety confirmations before execution.
* **🛡️ Safe Read-Only Mode:** Explore live account data with dynamic action-specific permission locks that prevent accidental charges or mutations.
* **❓ Global Help & Cheat Sheet Modal:** Press `?` or `h` from anywhere in the TUI to open an instant keyboard shortcuts and environment guide.

### ⌨️ Universal TUI Keybindings

| Key | Action |
| :--- | :--- |
| `↑` / `↓` / `j` / `k` | Navigate lists, tables, and form fields |
| `Enter` | Select module / Inspect table record / Submit form |
| `a` / `ctrl+a` | Open Natural Language AI Assistant (Available on ALL screens) |
| `Tab` / `Shift+Tab` | Switch tabs (Summary ⇄ JSON Tree) / Navigate form fields |
| `r` | Refresh live table data / Reset form |
| `c` | Open Credentials & Environment Configuration (Test ⇄ Live ⇄ AI) |
| `/` | Real-time fuzzy search across modules and actions |
| `?` / `h` | Toggle Global Help & Keybindings reference modal |
| `Esc` | Return to previous screen / Actions menu |
| `q` / `Ctrl+C` | Quit Razorpay TUI |

---

## 🤖 Natural Language AI Assistant

The TUI includes an integrated AI Assistant powered by **Google Gemini** (`gemini-flash-latest`, `gemini-2.5-flash-lite`) and **Anthropic Claude** (`claude-haiku-4-5-20251001`):

1. **How it works:** Press `a` from anywhere in the TUI. Type your request in natural language. The AI dynamically introspects the entire Cobra command tree and generates the exact command, flags, and an explanation.
2. **Safe Preview Model:** The AI **never** executes commands directly. It shows a full preview card (`Command`, `Explanation`, `Confidence Badge`, and destructive warnings for mutating operations).
3. **Smart Currency & Pre-Fill:** AI-suggested values (like amounts in paise) are automatically converted into friendly format (₹ Rupees) and pre-filled into interactive forms.
4. **Unified Execution:** Pressing `Enter` routes the suggested command through the exact same execution engine as manual TUI navigation, ensuring `--read-only` mode and credential checks apply identically.
5. **Configuration:** Provide your Gemini or Anthropic API Key via `~/.razorpay/config.yaml` (`ai_api_key`), inside the TUI via `[c] Config`, or with `export RAZORPAY_AI_API_KEY=AQ...` / `export RAZORPAY_AI_API_KEY=sk-ant-...`.

---

## Installation

Install the latest release for your platform:

```bash
curl -fsSL https://razorpay.com/cli/latest/install.sh | bash
```

The script downloads the right tarball for your OS/arch, extracts the `razorpay` binary, and places it at `/usr/local/bin/razorpay`. Confirm the install:

```bash
$ razorpay --version
razorpay version v1.0.9
```

Other options (Homebrew-style manual download, `go install`, building from source) are in [docs/install.md](docs/install.md).

## Configuration

Run `configure` interactively. The current value (if any) appears in brackets and Enter keeps it. The secret is masked while you type.

```bash
$ razorpay configure
Razorpay Key ID [None]: rzp_test_1DP5mmOlF5G5ag
Razorpay Key Secret [None]:
Credentials saved to /Users/you/.razorpay/config.yaml
```

Or pass credentials non-interactively (any flag you omit is prompted for):

```bash
razorpay configure --key-id rzp_test_xxxxxxxxxxxx --key-secret xxxxxxxxxxxxxxxxxxxx
```

Credentials live in `~/.razorpay/config.yaml`. Environment variables override the file:

```bash
export RAZORPAY_KEY_ID=rzp_test_xxxxxxxxxxxx
export RAZORPAY_KEY_SECRET=xxxxxxxxxxxxxxxxxxxx
```

Generate keys from the [Razorpay Dashboard](https://dashboard.razorpay.com/app/website-app-settings/api-keys): `rzp_test_` for development, `rzp_live_` for production.

## Examples

All amounts are in the smallest currency unit (paise for INR: `50000` = ₹500.00).

### Create an order

```bash
$ razorpay orders create --amount 50000 --currency INR --receipt order-001
{
  "id": "order_RB58MiP5SPFYyM",
  "entity": "order",
  "amount": 50000,
  "amount_paid": 0,
  "amount_due": 50000,
  "currency": "INR",
  "receipt": "order-001",
  "status": "created",
  "attempts": 0,
  "notes": [],
  "created_at": 1756455561
}
```

### Fetch an order

```bash
$ razorpay orders fetch order_RB58MiP5SPFYyM
{
  "id": "order_RB58MiP5SPFYyM",
  "entity": "order",
  "amount": 50000,
  "amount_paid": 0,
  "amount_due": 50000,
  "currency": "INR",
  "receipt": "order-001",
  "status": "created",
  "attempts": 0,
  "notes": [],
  "created_at": 1756455561
}
```

### List recent payments

```bash
$ razorpay payments list --count 2
{
  "entity": "collection",
  "count": 2,
  "items": [
    {
      "id": "pay_29QQoUBi66xm2f",
      "entity": "payment",
      "amount": 50000,
      "currency": "INR",
      "status": "captured",
      "order_id": "order_RB58MiP5SPFYyM",
      "method": "card",
      "captured": true,
      "created_at": 1756455622
    },
    {
      "id": "pay_29QQjUNkbFbiY8",
      "entity": "payment",
      "amount": 12500,
      "currency": "INR",
      "status": "authorized",
      "order_id": "order_RAxxJYS92AKL8z",
      "method": "upi",
      "captured": false,
      "created_at": 1756342140
    }
  ]
}
```

### Create a customer

```bash
$ razorpay customers create --name "Ada Lovelace" --email ada@example.com --contact 9876543210
{
  "id": "cust_1Aa00000000004",
  "entity": "customer",
  "name": "Ada Lovelace",
  "email": "ada@example.com",
  "contact": "9876543210",
  "gstin": null,
  "notes": [],
  "created_at": 1756455800
}
```

### Create a payment link

```bash
$ razorpay payment-links create \
    --amount 30000 --currency INR \
    --customer-name "Ada Lovelace" --customer-email ada@example.com --customer-contact 9876543210 \
    --description "Invoice 0042" --reference-id INV-0042
{
  "id": "plink_LFhRgkn8nfx6vY",
  "entity": "payment_link",
  "amount": 30000,
  "currency": "INR",
  "status": "created",
  "short_url": "https://rzp.io/i/abc123",
  "reference_id": "INV-0042",
  "description": "Invoice 0042",
  "customer": {
    "name": "Ada Lovelace",
    "email": "ada@example.com",
    "contact": "9876543210"
  },
  "created_at": 1756455900
}
```

### Errors

Errors are surfaced with the HTTP status and the API's response body, so they're easy to grep for in CI:

```bash
$ razorpay orders fetch order_does_not_exist
Error: API request failed with status 400: {"error":{"code":"BAD_REQUEST_ERROR","description":"The id provided does not exist"}}
```

Run `razorpay <group> --help` to list subcommands, or `razorpay <group> <subcommand> --help` for flags and examples.

## Commands

| Group           | What it manages                                       |
| --------------- | ----------------------------------------------------- |
| `configure`     | Save API credentials                                  |
| `payments`      | Payments (capture, card details, downtime, transfers) |
| `orders`        | Orders                                                |
| `customers`     | Customers                                             |
| `refunds`       | Refunds                                               |
| `invoices`      | Invoices and line items                               |
| `payment-links` | Payment Links                                         |
| `qr-codes`      | QR Codes                                              |
| `subscriptions` | Subscriptions and plans                               |
| `route`         | Route (linked accounts and transfers)                 |
| `smart-collect` | Smart Collect (virtual accounts)                      |
| `settlements`   | Settlements and reconciliation                        |
| `disputes`      | Disputes                                              |
| `documents`     | Documents                                             |

## Repository layout

```
cmd/
  <resource>/        one package per resource (payments, orders, ...)
    <resource>.go      parent Cmd + subcommand registration
    <subcommand>.go    one file per subcommand (list.go, fetch.go, ...)
  cmdutil/           shared client + error helpers
  configure.go       `configure` command (lives in package cmd directly)
  root.go            root command + subpackage wiring
api/                 HTTP client (auth, JSON pretty-print, multipart upload)
config/              config file + env-var loader (viper)
tests/               end-to-end test suite (build tag `e2e`)
docs/                per-resource usage guides
install.sh           platform installer
AGENTS.md            agent guidance for this repo
CHANGELOG.md         release notes
```

To add a subcommand to an existing resource, drop a new file into `cmd/<resource>/` and register it on `Cmd` inside that package's `<resource>.go`. To add a new top-level resource, mirror an existing folder and add `rootCmd.AddCommand(<pkg>.Cmd)` in `cmd/root.go`.

## Build and test

Every dev workflow has a `make` target. Prefer those over invoking `go` directly so flags, build tags, and lint versions stay in one place.

```bash
make setup                # download Go module dependencies
make build                # compile the razorpay binary (output in ./razorpay)
make build-all-platforms  # cross-compile for darwin / linux / windows
make fmt                  # gofmt -s every .go file
make lint                 # run golangci-lint (auto-installs the pinned version)
make test                 # run the end-to-end suite against the live API
make ci                   # build + test + lint + go mod tidy check
make clean                # remove build artefacts
```

The end-to-end suite lives in [`tests/`](tests/) and is gated by the `e2e` build tag. `make test` exits with a non-zero status unless `RAZORPAY_KEY_ID` and `RAZORPAY_KEY_SECRET` are exported in the environment:

```bash
export RAZORPAY_KEY_ID=rzp_test_xxxxxxxxxxxx
export RAZORPAY_KEY_SECRET=xxxxxxxxxxxxxxxxxxxx
make test
```

See [tests/README.md](tests/README.md) for the layout of the test suite and how each resource is exercised.

## Documentation

- [docs/install.md](docs/install.md): platform-specific install
- Per-resource guides: [payments](docs/payments.md) · [orders](docs/orders.md) · [customers](docs/customers.md) · [refunds](docs/refunds.md) · [settlements](docs/settlements.md) · [disputes](docs/disputes.md)
- [CHANGELOG.md](CHANGELOG.md): release notes
- [AGENTS.md](AGENTS.md): agent guidance for working in this repository
