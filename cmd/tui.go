package cmd

import (
	"github.com/razorpay/razorpay-cli/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start the interactive Terminal User Interface (TUI)",
	Long: `The 'tui' command launches a full-featured interactive terminal dashboard
for the Razorpay API. It lets you browse, search, create, and manage all
Razorpay resources (orders, payments, refunds, invoices, etc.) without having
to remember individual CLI flags or command names.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Start()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
