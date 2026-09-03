package cmd

import (
	"github.com/razorpay/razorpay-cli/ai"
	"github.com/razorpay/razorpay-cli/tui"
	"github.com/spf13/cobra"
)

var (
	tuiReadOnly bool
	tuiTestMode bool
	tuiLiveMode bool
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start the interactive Terminal User Interface (TUI)",
	Long: `The 'tui' command launches a full-featured interactive terminal dashboard
for the Razorpay API. It lets you browse, search, create and manage all
Razorpay resources (orders, payments, refunds, invoices, etc.) without having
to remember individual CLI flags or command names.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ai.SetGlobalRootCommand(rootCmd)
		initialMode := ""
		if tuiLiveMode {
			initialMode = "live"
		} else if tuiTestMode {
			initialMode = "test"
		}
		return tui.Start(tuiReadOnly, initialMode)
	},
}

func init() {
	tuiCmd.Flags().BoolVar(&tuiReadOnly, "read-only", false, "Interact with the TUI in safe read-only mode")
	tuiCmd.Flags().BoolVar(&tuiTestMode, "test", false, "Start directly in Test / Sandbox mode")
	tuiCmd.Flags().BoolVar(&tuiLiveMode, "live", false, "Start directly in Live / Production mode")
	rootCmd.AddCommand(tuiCmd)
}
